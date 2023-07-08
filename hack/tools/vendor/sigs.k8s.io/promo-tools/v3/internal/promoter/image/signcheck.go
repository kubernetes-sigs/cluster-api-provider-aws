/*
Copyright 2023 The Kubernetes Authors.

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
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	yaml "gopkg.in/yaml.v2"
	"sigs.k8s.io/release-sdk/sign"
	"sigs.k8s.io/release-utils/http"
	"sigs.k8s.io/release-utils/version"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/gcrane"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/google"
	"github.com/sirupsen/logrus"

	"github.com/sigstore/sigstore/pkg/cryptoutils"

	checkresults "sigs.k8s.io/promo-tools/v3/promoter/image/checkresults"
	options "sigs.k8s.io/promo-tools/v3/promoter/image/options"
	"sigs.k8s.io/promo-tools/v3/types/image"
)

var mirrorsList []string

const (
	repositoryPath = "k8s-artifacts-prod/images"

	// scanRegistry is the one we index to search for images
	scanRegistry = "us-central1-docker.pkg.dev"
)

func (di *DefaultPromoterImplementation) GetLatestImages(opts *options.Options) ([]string, error) {
	// If there is a list of images to check in the options
	// we default to checking those.
	if len(opts.SignCheckReferences) > 0 {
		for _, refString := range opts.SignCheckReferences {
			_, err := name.ParseReference(refString)
			if err != nil {
				return nil, fmt.Errorf("invalid image reference %s: %w", refString, err)
			}
		}
		return opts.SignCheckReferences, nil
	}

	images, err := di.readLatestImages(opts)
	if err != nil {
		return nil, fmt.Errorf("fetching latest images: %w", err)
	}
	logrus.Infof("Images to check: +%v", images)
	return images, nil
}

func (di *DefaultPromoterImplementation) getMirrors() ([]string, error) {
	if mirrorsList != nil {
		return mirrorsList, nil
	}
	urls := []string{}
	iurls := map[string]string{}
	manifest, err := http.NewAgent().Get(
		"https://github.com/kubernetes/k8s.io/raw/main/k8s.gcr.io/manifests/k8s-staging-kubernetes/promoter-manifest.yaml",
	)
	if err != nil {
		return nil, fmt.Errorf("downloading promoter manifest: %w", err)
	}

	type entriesList struct {
		Registries []struct {
			Name string `yaml:"name,omitempty"`
			Src  bool   `yaml:"src,omitempty"`
		} `yaml:"registries"`
	}

	entries := entriesList{}
	if err := yaml.Unmarshal(manifest, &entries); err != nil {
		return nil, fmt.Errorf("unmarshalling promoter manifest: %w", err)
	}

	for _, e := range entries.Registries {
		if e.Src {
			continue
		}
		u, err := url.Parse("https://" + e.Name)
		if err != nil {
			return nil, fmt.Errorf("parsing url %s: %w", u, err)
		}
		iurls[u.Hostname()] = u.Hostname()
	}

	for u := range iurls {
		urls = append(urls, u)
	}
	mirrorsList = urls
	return urls, nil
}

func (di *DefaultPromoterImplementation) GetSignatureStatus(
	opts *options.Options, images []string,
) (checkresults.Signature, error) {
	results := checkresults.Signature{}
	mirrors, err := di.getMirrors()
	if err != nil {
		return results, fmt.Errorf("reading mirrors: %w", err)
	}
	logrus.Infof(
		"Checking %d images for signatures, each in %d mirrors",
		len(images), len(mirrors),
	)
	for _, refString := range images {
		ref, err := name.ParseReference(refString)
		if err != nil {
			return results, fmt.Errorf("parsing reference: %w", err)
		}

		digest, err := crane.Digest(refString)
		if err != nil {
			return results, fmt.Errorf("getting digest for %s: %w", refString, err)
		}

		targetImages := []string{}
		for _, mirror := range mirrors {
			rpath := repositoryPath
			if strings.HasSuffix(mirror, ".gcr.io") {
				rpath = "k8s-artifacts-prod"
			}

			targetImages = append(targetImages, fmt.Sprintf("%s/%s/%s:%s.sig",
				mirror, rpath, ref.Context().RepositoryStr(),
				strings.ReplaceAll(digest, ":", "-"),
			))
		}

		logrus.Infof("Checking %s for signatures in %d mirrors", refString, len(targetImages))
		existing, missing, err := di.CheckSignatureLayers(opts, targetImages)
		if err != nil {
			return results, fmt.Errorf("checking objects: %w", err)
		}
		results[refString] = checkresults.CheckList{
			Signed:  existing,
			Missing: missing,
		}
	}
	return results, nil
}

// miniManifest is a minimal representation of the sigstore signature manifest
type miniManifest struct {
	Layers []struct {
		MediaType   string
		Annotations map[string]string `json:"annotations"`
	} `json:"layers"`
}

// CheckSignatureLayers checks a list of signature layers to ensure
func (di *DefaultPromoterImplementation) CheckSignatureLayers(opts *options.Options, oList []string) (existing, missing []string, err error) {
	// TODO: Parallelize this check
	existing = []string{}
	missing = []string{}
	for _, s := range oList {
		e, err := objectExists(opts, s)
		if err != nil {
			return existing, missing, fmt.Errorf("checking reference: %w", err)
		}

		if !e {
			missing = append(missing, s)
			continue
		}

		existing = append(existing, s)
	}
	return existing, missing, nil
}

func objectExists(opts *options.Options, refString string) (bool, error) {
	// Check
	manifestData, err := crane.Manifest(refString)
	if err != nil {
		if strings.Contains(err.Error(), "MANIFEST_UNKNOWN") {
			logrus.WithField("image", refString).Info("No signature found")
			return false, nil
		}
		return false, fmt.Errorf("pulling signature manifest: %w", err)
	}

	manifest := &miniManifest{}
	if err := json.Unmarshal(manifestData, manifest); err != nil {
		return false, fmt.Errorf("parsing .sig image manifest: %w", err)
	}

	// Get the certificate
	if manifest.Layers == nil || len(manifest.Layers) == 0 {
		return false, nil
	}
	signedLayers := 0
	for _, layer := range manifest.Layers {
		if layer.MediaType != "application/vnd.dev.cosign.simplesigning.v1+json" {
			continue
		}

		certData, ok := layer.Annotations["dev.sigstore.cosign/certificate"]
		if !ok {
			continue
		}

		var b bytes.Buffer
		b.Write([]byte(certData))

		certs, err := cryptoutils.LoadCertificatesFromPEM(&b)
		if err != nil {
			return false, err
		}

		names := cryptoutils.GetSubjectAlternateNames(certs[0])
		for _, n := range names {
			if n == opts.SignCheckIdentity {
				return true, nil
			}
		}
		signedLayers++
	}

	if signedLayers == 0 {
		logrus.WithField("image", refString).Debugf("No certificates found")
	} else {
		logrus.WithField("image", refString).Debugf("Image signed, but not with expected identity")
	}

	return false, nil
}

// FixMissingSignatures signs an image that has no signatures at all
func (di *DefaultPromoterImplementation) FixMissingSignatures(opts *options.Options, results checkresults.Signature) error {
	for mainImg, res := range results {
		if len(res.Signed) > 0 {
			continue
		}

		logrus.Infof("Signing and replicating first mirror (%s)", mainImg)

		// Build the digest of the first missing one
		digestRef := strings.TrimSuffix(strings.ReplaceAll(res.Missing[0], ":sha256-", "@sha256:"), ".sig")
		if err := di.signReference(opts, digestRef); err != nil {
			return fmt.Errorf("signing first mirror reference %s: %w", digestRef, err)
		}

		logrus.Infof("Replicating image to %d mirrors", len(res.Missing[1:]))
		for _, targetRef := range res.Missing[1:] {
			if err := di.replicateReference(opts, res.Missing[0], targetRef); err != nil {
				return fmt.Errorf("replicating signature: %w", err)
			}
		}
	}
	return nil
}

// FixPartialSignatures fixes images that had some signatures but some mirrors
// are missing some signatures
func (di *DefaultPromoterImplementation) FixPartialSignatures(opts *options.Options, results checkresults.Signature) error {
	for mainImg, res := range results {
		if len(res.Missing) == 0 || len(res.Signed) == 0 {
			continue
		}

		logrus.Infof("%s has %d signed copies, %d are missing", mainImg, len(res.Signed), len(res.Missing))

		sourceRef := res.Signed[0]
		for _, targetRef := range res.Missing {
			// Copy the first signature to the target ref
			if err := di.replicateReference(opts, sourceRef, targetRef); err != nil {
				return fmt.Errorf("replicating signature: %w", err)
			}
		}
	}
	return nil
}

// replicateReference copies an image reference to another mirror
func (di *DefaultPromoterImplementation) replicateReference(opts *options.Options, srcRef, dstRef string) error {
	craneOpts := []crane.Option{
		crane.WithAuthFromKeychain(gcrane.Keychain),
		crane.WithUserAgent(image.UserAgent),
	}

	if !opts.SignCheckFix {
		logrus.Infof(" (NOOP) replicating %s to %s ", srcRef, dstRef)
		return nil
	}

	logrus.Infof(" replicating %s to %s ", srcRef, dstRef)

	if err := crane.Copy(srcRef, dstRef, craneOpts...); err != nil {
		return fmt.Errorf(
			"copying signature %s to %s: %w", srcRef, dstRef, err,
		)
	}
	return nil
}

// signReference takes a reference and signs it
func (di *DefaultPromoterImplementation) signReference(opts *options.Options, refString string) error {
	if !opts.SignCheckFix {
		logrus.Infof(" (NOOP) signing %s", refString)
		return nil
	}
	logrus.Infof(" signing %s", refString)

	// Options for the new signer
	signOpts := sign.Default()

	// Get the identity token we will use
	token, err := di.GetIdentityToken(opts, opts.SignerAccount)
	if err != nil {
		return fmt.Errorf("generating identity token: %w", err)
	}
	signOpts.IdentityToken = token

	// We want to sign all entities for multi-arch images
	signOpts.Recursive = true

	di.signer = sign.New(signOpts)

	// Add an annotation recording the kpromo version to ensure we
	// get a 2nd signature, otherwise cosign will not resign a signed image:
	signOpts.Annotations = map[string]interface{}{
		"org.kubernetes.kpromo.version": fmt.Sprintf(
			"kpromo-%s", version.GetVersionInfo().GitVersion,
		),
	}

	if _, err := di.signer.SignImageWithOptions(signOpts, refString); err != nil {
		return fmt.Errorf("signing image %s: %w", refString, err)
	}

	return nil
}

// readLatestImages returns the latest images uploaded to the registry.
// Note that this function uses the google GCR/AR extensions so it will
// not work on other non-GCP registries.
func (di *DefaultPromoterImplementation) readLatestImages(opts *options.Options) ([]string, error) {
	creds := google.WithAuthFromKeychain(authn.NewMultiKeychain(
		authn.DefaultKeychain,
		google.Keychain,
	))

	dateCutOff := time.Now().AddDate(0, 0, opts.SignCheckFromDays*-1)
	dateCutOffTo := time.Now()
	if opts.SignCheckToDays > 0 {
		dateCutOffTo = time.Now().AddDate(0, 0, opts.SignCheckToDays*-1)
	}
	logrus.Infof("Checking images from %s to %s",
		dateCutOff.Local().Format(time.RFC822),
		dateCutOffTo.Local().Format(time.RFC822),
	)
	images := []string{}

	repo, err := name.NewRepository(scanRegistry+"/"+repositoryPath, name.WeakValidation)
	if err != nil {
		return nil, fmt.Errorf("creating repo: %w", err)
	}

	var mt sync.Mutex
	walkFn := func(repo name.Repository, tags *google.Tags, err error) error {
		if tags == nil {
			return nil
		}

		// We ignore the -arch repositories as the promoter currently
		// ignores them and does not sign them
		if strings.HasSuffix(repo.String(), "-amd64") || strings.HasSuffix(repo.String(), "-arm") ||
			strings.HasSuffix(repo.String(), "-arm64") || strings.HasSuffix(repo.String(), "-ppc64le") ||
			strings.HasSuffix(repo.String(), "-s390x") {
			return nil
		}
		logrus.Infof("Indexing %d images from %s", len(tags.Manifests), repo)
		// First var (_) is the digest
		for _, manifest := range tags.Manifests {
			// Ignore if there are no tags
			if len(manifest.Tags) == 0 {
				continue
			}
			// Ignore signature tags
			if strings.HasSuffix(manifest.Tags[0], ".sig") {
				continue
			}

			// Ignore if uploaded before our date
			if manifest.Uploaded.Before(dateCutOff) {
				continue
			}

			if opts.SignCheckToDays > 0 && manifest.Uploaded.After(dateCutOffTo) {
				continue
			}

			mt.Lock()
			images = append(images, strings.ReplaceAll(
				fmt.Sprintf("%s:%s", repo, manifest.Tags[0]),
				scanRegistry+"/"+repositoryPath, "registry.k8s.io"),
			)
			mt.Unlock()
		}
		return nil
	}

	if err := google.Walk(repo, walkFn, creds); err != nil {
		return nil, fmt.Errorf("walking repo: %w", err)
	}

	if opts.SignCheckMaxImages != 0 && len(images) > opts.SignCheckMaxImages {
		images = images[0:opts.SignCheckMaxImages]
	}

	return images, nil
}
