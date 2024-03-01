/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package imagepromoter

import (
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	reg "sigs.k8s.io/promo-tools/v4/internal/legacy/dockerregistry"
	"sigs.k8s.io/promo-tools/v4/internal/legacy/dockerregistry/registry"
	"sigs.k8s.io/promo-tools/v4/internal/legacy/dockerregistry/schema"
	"sigs.k8s.io/promo-tools/v4/internal/legacy/gcloud"
	"sigs.k8s.io/promo-tools/v4/internal/legacy/stream"
	"sigs.k8s.io/promo-tools/v4/internal/version"
	options "sigs.k8s.io/promo-tools/v4/promoter/image/options"
	"sigs.k8s.io/promo-tools/v4/types/image"
	"sigs.k8s.io/release-sdk/sign"
)

const vulnerabilityDiscalimer = `DISCLAIMER: Vulnerabilities are found as issues with package
binaries within image layers, not necessarily with the image layers themselves.
So a 'fixable' vulnerability may not necessarily be immediately actionable. For
example, even though a fixed version of the binary is available, it doesn't
necessarily mean that a new version of the image layer is available.`

// streamProducerFunc is a function that gets the required fields to
// construct a promotion stream producer
type StreamProducerFunc func(
	srcRegistry image.Registry, srcImageName image.Name,
	destRC registry.Context, imageName image.Name,
	digest image.Digest, tag image.Tag, tp reg.TagOp,
) stream.Producer

type DefaultPromoterImplementation struct {
	signer *sign.Signer
}

// NewDefaultPromoterImplementation creates a new DefaultPromoterImplementation instance.
func NewDefaultPromoterImplementation(opts *options.Options) *DefaultPromoterImplementation {
	return &DefaultPromoterImplementation{
		signer: sign.New(defaultSignerOptions(opts)),
	}
}

// defaultSignerOptions returns a new *sign.Options with default values applied.
func defaultSignerOptions(opts *options.Options) *sign.Options {
	signOpts := sign.Default()

	// We want to sign all entities for multi-arch images
	signOpts.Recursive = true

	// Recursive signing can take a bit longer than usual
	signOpts.Timeout = 15 * time.Minute

	// The Certificate Identity to be used to check the images signatures
	signOpts.CertIdentity = opts.SignCheckIdentity

	// The Certificate OICD Issuer to be used to check the images signatures
	signOpts.CertOidcIssuer = opts.SignCheckIssuer

	// A regex Certificate Identity to be used to check the images signatures
	signOpts.CertIdentityRegexp = opts.SignCheckIdentityRegexp

	// A regex to match a Certificate OICD Issuer to be used to check the images signatures
	signOpts.CertOidcIssuerRegexp = opts.SignCheckIssuerRegexp

	return signOpts
}

// ValidateOptions checks an options set
func (di *DefaultPromoterImplementation) ValidateOptions(opts *options.Options) error {
	if opts.Snapshot == "" && opts.ManifestBasedSnapshotOf == "" {
		if opts.Manifest == "" && opts.ThinManifestDir == "" {
			return errors.New("either a manifest ot a thin manifest dir have to be set")
		}
	}
	return nil
}

// ActivateServiceAccounts gets key files and activates service accounts
func (di *DefaultPromoterImplementation) ActivateServiceAccounts(opts *options.Options) error {
	if !opts.UseServiceAcct {
		logrus.Warn("Not setting a service account")
	}
	if err := gcloud.ActivateServiceAccounts(opts.KeyFiles); err != nil {
		return fmt.Errorf("activating service accounts: %w", err)
	}
	// TODO: Output to log the accout used
	return nil
}

// PrecheckAndExit run simple prechecks to exit before promotions
// or security scans
func (di *DefaultPromoterImplementation) PrecheckAndExit(
	opts *options.Options, mfests []schema.Manifest,
) error {
	// Make the sync context tu run the prechecks:
	sc, err := di.MakeSyncContext(opts, mfests)
	if err != nil {
		return fmt.Errorf("generatinng sync context for prechecks: %w", err)
	}

	// Run the prechecks, these will be run and the calling
	// mode of operation should exit.
	if err := sc.RunChecks([]reg.PreCheck{}); err != nil {
		return fmt.Errorf("running prechecks before promotion: %w", err)
	}
	return nil
}

func (di *DefaultPromoterImplementation) PrintVersion() {
	logrus.Info(version.Get())
}

// printSection handles the start/finish labels in the
// former legacy cli/run code
func (di *DefaultPromoterImplementation) PrintSection(message string, confirm bool) {
	dryRunLabel := ""
	if !confirm {
		dryRunLabel = "(DRY RUN) "
	}
	logrus.Infof("********** %s %s**********", message, dryRunLabel)
}

// printSecDisclaimer prints a disclaimer about false positives
// that may be found in container image lauyers.
func (di *DefaultPromoterImplementation) PrintSecDisclaimer() {
	logrus.Info(vulnerabilityDiscalimer)
}
