package quill

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/anchore/quill/internal/bus"
	"github.com/anchore/quill/internal/log"
	"github.com/anchore/quill/quill/event"
	"github.com/anchore/quill/quill/notary"
)

type NotarizeConfig struct {
	StatusConfig notary.StatusConfig
	HTTPTimeout  time.Duration
	TokenConfig  notary.TokenConfig
}

func NewNotarizeConfig(issuer, privateKeyID, privateKey string) *NotarizeConfig {
	timeout := 15 * time.Minute
	return &NotarizeConfig{
		StatusConfig: notary.StatusConfig{
			Timeout: timeout,
			Poll:    10 * time.Second,
			Wait:    true,
		},
		HTTPTimeout: 30 * time.Second,
		TokenConfig: notary.TokenConfig{
			Issuer:        issuer,
			PrivateKeyID:  privateKeyID,
			TokenLifetime: timeout + (2 * time.Minute),
			PrivateKey:    privateKey,
		},
	}
}

func (c *NotarizeConfig) WithStatusConfig(cfg notary.StatusConfig) *NotarizeConfig {
	c.StatusConfig = cfg
	c.TokenConfig.TokenLifetime = cfg.Timeout + (2 * time.Minute)
	return c
}

/*

Source: https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution

Apple's notary service requires you to adopt the following protections:

- Enable code-signing for all of the executables you distribute, and ensure that executables have valid code signatures,
  as described in Ensure a valid code signature.
  See https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution/resolving_common_notarization_issues#3087735

- Use a “Developer ID” application, kernel extension, system extension, or installer certificate for your code-signing
  signature. (Don't use a Mac Distribution, ad hoc, Apple Developer, or local development certificate.) Verify the
  certificate type before submitting, as described in Use a valid Developer ID certificate.
  See https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution/resolving_common_notarization_issues#3087721
  For more information, see Create, export, and delete signing certificates: https://help.apple.com/xcode/mac/current/#/dev154b28f09

- Enable the Hardened Runtime capability for your app and command line targets, as described in Enable hardened runtime.
  See https://help.apple.com/xcode/mac/current/#/devf87a2ac8f

- Include a secure timestamp with your code-signing signature. (The Xcode distribution workflow includes a secure
  timestamp by default. For custom workflows, see Include a secure timestamp.)
  See https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution/resolving_common_notarization_issues#3087733

- Don’t include the com.apple.security.get-task-allow entitlement with the value set to any variation of true. If your
  software hosts third-party plug-ins and needs this entitlement to debug the plug-in in the context of a host
  executable, see Avoid the get-task-allow entitlement.
  See https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution/resolving_common_notarization_issues#3087731

- Link against the macOS 10.9 or later SDK, as described in Use the macOS 10.9 SDK or later.
  See https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution/resolving_common_notarization_issues#3087723

- Ensure your processes have properly-formatted XML, ASCII-encoded entitlements as described in Ensure properly
  formatted entitlements.
  See https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution/resolving_common_notarization_issues#3561456

*/

func Notarize(path string, cfg NotarizeConfig) (notary.SubmissionStatus, error) {
	log.WithFields("binary", path).Info("notarizing binary")

	mon := bus.PublishTask(
		event.Title{
			Default:      "Notarize binary",
			WhileRunning: "Notarizing binary",
			OnSuccess:    "Notarized binary",
		},
		path,
		-1,
	)

	defer mon.SetCompleted()

	mon.Stage.Current = "validating binary"

	if isSigned, err := IsSigned(path); err != nil {
		return "", fmt.Errorf("unable to determine if binary is signed: %+v", err)
	} else if !isSigned {
		return "", fmt.Errorf("binary is not signed thus will not pass notarization")
	}

	mon.Stage.Current = "initializing client"

	token, err := notary.NewSignedToken(cfg.TokenConfig)
	if err != nil {
		return "", err
	}

	a := notary.NewAPIClient(token, cfg.HTTPTimeout)

	mon.Stage.Current = "processing payload"

	bin, err := notary.NewPayload(path)
	if err != nil {
		return "", err
	}

	mon.Stage.Current = "submitting"

	sub := notary.NewSubmission(a, bin)

	if err := sub.Start(context.Background()); err != nil {
		return "", fmt.Errorf("unable to start submission: %+v", err)
	}

	if !cfg.StatusConfig.Wait {
		log.WithFields("id", sub.ID()).Infof("Submission started but configured to not wait for the results")
		return "", nil
	}

	statusCfg := cfg.StatusConfig.WithProgress(&mon.Stage)

	status, err := notary.PollStatus(context.Background(), sub, *statusCfg)

	mon.Stage.Current = strings.ToLower(fmt.Sprintf("status %q", string(status)))

	return status, err
}
