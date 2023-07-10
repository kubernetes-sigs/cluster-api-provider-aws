/*
Copyright 2019 The Kubernetes Authors.

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

package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime/debug"

	"cloud.google.com/go/errorreporting"
	"github.com/sirupsen/logrus"

	reg "sigs.k8s.io/promo-tools/v3/internal/legacy/dockerregistry"
	"sigs.k8s.io/promo-tools/v3/internal/legacy/dockerregistry/registry"
	"sigs.k8s.io/promo-tools/v3/internal/legacy/dockerregistry/schema"
	"sigs.k8s.io/promo-tools/v3/internal/legacy/logclient"
	"sigs.k8s.io/promo-tools/v3/internal/legacy/remotemanifest"
	"sigs.k8s.io/promo-tools/v3/internal/legacy/report"
)

// InitRealServerContext creates a ServerContext with facilities that are meant
// for production use (going over the network to fetch actual official promoter
// manifests from GitHub, for example).
func InitRealServerContext(
	gcpProjectID, repoURLStr, branch, path, uuid string,
) (*ServerContext, error) {
	remoteManifestFacility, err := remotemanifest.NewGit(
		repoURLStr,
		branch,
		path)
	if err != nil {
		return nil, err
	}
	reportingFacility := report.NewGcpErrorReportingClient(
		gcpProjectID, "cip-auditor")
	loggingFacility, err := logclient.NewGcpLogClient(
		gcpProjectID, LogName)
	if err != nil {
		return nil, err
	}

	serverContext := ServerContext{
		ID:                     uuid,
		RemoteManifestFacility: remoteManifestFacility,
		ErrorReportingFacility: reportingFacility,
		LoggingFacility:        loggingFacility,
		GcrReadingFacility: GcrReadingFacility{
			ReadRepo:         reg.MkReadRepositoryCmdReal,
			ReadManifestList: reg.MkReadManifestListCmdReal,
		},
	}

	return &serverContext, nil
}

// RunAuditor runs an HTTP server.
func (s *ServerContext) RunAuditor() {
	logrus.Debug("Running on Process ID (PID):", os.Getgid())
	logrus.Info("Starting Auditor")
	logrus.Infoln(s)

	defer s.LoggingFacility.Close()
	defer s.ErrorReportingFacility.Close()

	http.HandleFunc(
		"/",
		func(w http.ResponseWriter, r *http.Request) {
			s.Audit(w, r)
		},
	)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		logrus.Infof("Defaulting to port %s", port)
	}
	// Start HTTP server.
	logrus.Infof("Listening on port %s", port)
	logrus.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

// ParsePubSubMessage parses an HTTP request body into a reg.GCRPubSubPayload.
func ParsePubSubMessage(body io.Reader) (*reg.GCRPubSubPayload, error) {
	// Handle basic errors (malformed requests).
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("iotuil.ReadAll: %v", err)
	}

	return ParsePubSubMessageBody(bodyBytes)
}

// ParsePubSubMessageBody parses the body of an HTTP request to be a
// GCRPubSubPayload.
func ParsePubSubMessageBody(
	body []byte,
) (*reg.GCRPubSubPayload, error) {
	var psm PubSubMessage
	var gcrPayload reg.GCRPubSubPayload

	if err := json.Unmarshal(body, &psm); err != nil {
		return nil, fmt.Errorf("json.Unmarshal (request body): %v", err)
	}

	if err := json.Unmarshal(psm.Message.Data, &gcrPayload); err != nil {
		return nil, fmt.Errorf("json.Unmarshal (message data): %v", err)
	}

	return &gcrPayload, nil
}

// ValidatePayload ensures that the payload is well-formed, per our
// business-logic needs.
func ValidatePayload(gcrPayload *reg.GCRPubSubPayload) error {
	if gcrPayload.FQIN == "" && gcrPayload.PQIN == "" {
		return fmt.Errorf(
			"%v: neither 'digest' nor 'tag' was specified", gcrPayload)
	}

	if err := gcrPayload.PopulateExtraFields(); err != nil {
		return err
	}

	switch gcrPayload.Action {
	case "":
		return fmt.Errorf("%v: Action not specified", gcrPayload)

	// All deletions will for now be treated as an error. If it's an insertion,
	// it can either have "digest" with FQIN, or "digest" + "tag" with PQIN. So
	// we always verify FQIN, and if there is PQIN, verify that as well.
	case "DELETE":
		// Even though this is an error, we successfully processed this message,
		// so exit with an error.
		return fmt.Errorf("%v: deletions are prohibited", gcrPayload)
	case "INSERT":
		break
	default:
		return fmt.Errorf(
			"%v: unknown action %q", gcrPayload, gcrPayload.Action)
	}

	return nil
}

// Audit receives and processes a Pub/Sub push message. It has 3 parts: (1)
// parse the request body to understand the GCR state change, (2) update the Git
// repo of the promoter manifests, and (3) reconcile these two against each
// other.
func (s *ServerContext) Audit(w http.ResponseWriter, r *http.Request) {
	logInfo := s.LoggingFacility.GetInfoLogger()
	logError := s.LoggingFacility.GetErrorLogger()
	logAlert := s.LoggingFacility.GetAlertLogger()

	defer func() {
		if msg := recover(); msg != nil {
			//nolint:errcheck // TODO: Check result of type assertion
			panicStr := msg.(string)

			stacktrace := debug.Stack()

			s.ErrorReportingFacility.Report(
				errorreporting.Entry{
					Req:   r,
					Error: fmt.Errorf("%s", panicStr),
					Stack: stacktrace,
				},
			)

			logAlert.Printf("%s\n%s\n", panicStr, string(stacktrace))
			logrus.Errorln(panicStr)
		}
	}()
	// (1) Parse request payload.
	gcrPayload, err := ParsePubSubMessage(r.Body)
	if err != nil {
		// It's important to fail any message we cannot parse, because this
		// notifies us of any changes in how the messages are created in the
		// first place.
		msg := fmt.Sprintf("(%s) TRANSACTION REJECTED: parse failure: %v", s.ID, err)

		// TODO: Properly catch write errors
		//nolint:errcheck
		_, _ = w.Write([]byte(msg))

		// TODO(panic): Don't panic!
		panic(msg)
	}

	// Additionally, fail any message that we cannot validate. This is where we
	// catch things like "DELETE" actions and warn them outright as all
	// deletions are prohibited.
	if err := ValidatePayload(gcrPayload); err != nil {
		msg := fmt.Sprintf("(%s) TRANSACTION REJECTED: validation failure: %v", s.ID, err)

		// TODO: Properly catch write errors
		//nolint:errcheck
		_, _ = w.Write([]byte(msg))

		// TODO(panic): Don't panic!
		panic(msg)
	}

	msg := fmt.Sprintf(
		"(%s) HANDLING MESSAGE: %v\n", s.ID, gcrPayload)
	logInfo.Println(msg)

	// (2) Clone fresh repo (or use one already on disk).
	logrus.Debug("Cloning GCR repo...")
	manifests, err := s.RemoteManifestFacility.Fetch()
	if err != nil {
		logError.Println(err)
		// If there is an error, return an HTTP error so that the Pub/Sub
		// message may be retried (this is a behavior of Cloud Run's handling of
		// Pub/Sub messages that are converted into HTTP messages).
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Debug info.
	logInfo.Printf("(%s) gcrPayload: %v", s.ID, gcrPayload)
	logInfo.Printf("(%s) RemoteManifestFacility: %v", s.ID, s.RemoteManifestFacility)

	// (3) Compare GCR state change with the intent of the promoter manifests
	// (exact digest match on a tagged or tagless image).
	for _, manifest := range manifests {
		m := gcrPayload.Match(&manifest)
		if (m.DigestMatch || m.TagMatch) &&
			!m.TagMismatch {
			msg := fmt.Sprintf(
				"(%s) TRANSACTION VERIFIED: %v: agrees with manifest\n",
				s.ID,
				gcrPayload,
			)

			logInfo.Println(msg)
			logrus.Infoln(msg)

			// TODO: Properly catch write errors
			//nolint:errcheck
			_, _ = w.Write([]byte(msg))
			return
		}
	}

	logInfo.Printf("(%s): could not find direct manifest entry for %v; assuming child manifest", s.ID, gcrPayload)

	// (4) It could be that the manifest is a child manifest (part of a fat
	// manifest). This is the case where the user only specifies the digest of
	// the parent image, but not the child image. When the promoter copies over
	// the entirety of the fat manifest, it will necessarily copy over the child
	// images as part of the transaction. To validate child images, we have to
	// first scan the source repository (from where the child image is being
	// promoted from) and then run reg.ReadGCRManifestLists to populate the
	// parent/child relationship maps of all relevant fat manifests.
	//
	// Because the subproject is gcr.io/k8s-artifacts-prod/<subproject>/foo...,
	// we can search for the matching subproject and run
	// reg.ReadGCRManifestLists.
	sc, err := reg.MakeSyncContext(
		manifests,
		// threads
		10,
		// dry run (although not necessary as we'll only be doing image reads,
		// it doesn't hurt)
		true,
		// useServiceAccount
		false)
	if err != nil {
		// Retry Pub/Sub message if the above fails (it shouldn't because
		// MakeSyncContext can only error out if the useServiceAccount bool is
		// set to True).
		logError.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	// Find the subproject repository responsible for the GCRPubSubPayload. This
	// is so that we can query the subproject to figure out all digests that
	// belong there, so that we can validate the child manifest in the
	// GCRPubSubPayload.
	//
	// If we can't find any source registry for this image, then reject the
	// transaction.
	srcRegistries, err := GetMatchingSourceRegistries(&manifests, gcrPayload)
	if err != nil {
		msg := fmt.Sprintf("(%s) TRANSACTION REJECTED: %v", s.ID, err)

		// TODO: Properly catch write errors
		//nolint:errcheck
		_, _ = w.Write([]byte(msg))

		// TODO(panic): Don't panic!
		panic(msg)
	}

	logInfo.Printf("(%s): reading srcRegistries %v for %q", s.ID, srcRegistries, gcrPayload)
	logrus.Debug("Querying GCR repository...")

	sc.ReadRegistries(
		srcRegistries,
		true,
		s.GcrReadingFacility.ReadRepo,
	)

	sc.ReadGCRManifestLists(s.GcrReadingFacility.ReadManifestList)
	if gcrPayload.Digest == "" {
		msg := fmt.Sprintf("(%s) TRANSACTION REJECTED: digest missing from payload --- cannot check parent digest: %v", s.ID, gcrPayload.Digest)

		// TODO: Properly catch write errors
		//nolint:errcheck
		_, _ = w.Write([]byte(msg))

		// TODO(panic): Don't panic!
		panic(msg)
	}

	logrus.Infof("(%s): looking for child digest %v", s.ID, gcrPayload.Digest)
	if parentDigest, hasParent := sc.ParentDigest[gcrPayload.Digest]; hasParent {
		msg := fmt.Sprintf(
			"(%s) TRANSACTION VERIFIED: %v: agrees with manifest (parent digest %v)\n", s.ID, gcrPayload, parentDigest)
		logInfo.Println(msg)
		logrus.Infoln(msg)

		// TODO: Properly catch write errors
		//nolint:errcheck
		_, _ = w.Write([]byte(msg))

		return
	}

	// (5) If all of the above checks fail, then this transaction is unable to be
	// verified.
	msg = fmt.Sprintf(
		"(%s) TRANSACTION REJECTED: %v: could not validate", s.ID, gcrPayload)
	// Return 200 OK, because we don't want to re-process this transaction.
	// "Terminating" the auditing here simplifies debugging as well, because the
	// same message is not repeated over and over again in the logs.

	// TODO: Properly catch write errors
	//nolint:errcheck
	_, _ = w.Write([]byte(msg))

	// TODO(panic): Don't panic!
	panic(msg)
}

// GetMatchingSourceRegistries gets the first source repository that matches the
// image information inside a GCRPubSubPayload.
func GetMatchingSourceRegistries(
	manifests *[]schema.Manifest,
	gcrPayload *reg.GCRPubSubPayload,
) ([]registry.Context, error) {
	rcs := []registry.Context{}

	for _, manifest := range *manifests {
		if !gcrPayload.Match(&manifest).PathMatch {
			continue
		}
		// Now that there is a match (at least a PathMatch), just fish out
		// the source registry in the manifest.
		for _, rc := range manifest.Registries {
			if rc.Src {
				rcs = append(rcs, rc)
			}
		}
	}

	if len(rcs) > 0 {
		return rcs, nil
	}

	return rcs,
		fmt.Errorf(
			"could not find matching source registry for %v",
			gcrPayload.FQIN)
}
