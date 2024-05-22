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
	"context"
	"fmt"
	"os"
	"strings"

	credentials "cloud.google.com/go/iam/credentials/apiv1"
	"cloud.google.com/go/iam/credentials/apiv1/credentialspb"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/gcrane"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/nozzle/throttler"
	"github.com/sigstore/sigstore/pkg/tuf"
	"github.com/sirupsen/logrus"
	gopts "google.golang.org/api/option"

	"sigs.k8s.io/promo-tools/v4/image/consts"
	reg "sigs.k8s.io/promo-tools/v4/internal/legacy/dockerregistry"
	"sigs.k8s.io/promo-tools/v4/internal/legacy/gcloud"
	options "sigs.k8s.io/promo-tools/v4/promoter/image/options"
	"sigs.k8s.io/promo-tools/v4/promoter/image/ratelimit"
	"sigs.k8s.io/promo-tools/v4/types/image"
	"sigs.k8s.io/release-sdk/sign"
	"sigs.k8s.io/release-utils/version"
)

const (
	oidcTokenAudience  = "sigstore"
	signatureTagSuffix = ".sig"

	TestSigningAccount = "k8s-infra-promoter-test-signer@k8s-cip-test-prod.iam.gserviceaccount.com"
)

// ValidateStagingSignatures checks if edges (images) have a signature
// applied during its staging run. If they do it verifies them and
// returns an error if they are not valid.
func (di *DefaultPromoterImplementation) ValidateStagingSignatures(
	edges map[reg.PromotionEdge]interface{},
) (map[reg.PromotionEdge]interface{}, error) {
	refsToEdges := map[string]reg.PromotionEdge{}
	for edge := range edges {
		ref := edge.SrcReference()
		refsToEdges[ref] = edge
	}

	refs := []string{}
	for ref := range refsToEdges {
		refs = append(refs, ref)
	}

	res, err := di.signer.VerifyImages(refs...)
	if err != nil {
		return nil, fmt.Errorf("verify images: %w", err)
	}

	signedEdges := map[reg.PromotionEdge]interface{}{}
	res.Range(func(key, _ any) bool {
		ref, ok := key.(string)
		if !ok {
			logrus.Errorf("Interface conversion failed: key is not a string: %v", key)
			return false
		}

		edge, ok := refsToEdges[ref]
		if !ok {
			logrus.Errorf("Reference %s is not in edge map", ref)
			return true
		}
		signedEdges[edge] = nil
		return true
	})

	return signedEdges, nil
}

// SignImages signs the promoted images and stores their signatures in
// the registry
func (di *DefaultPromoterImplementation) SignImages(
	opts *options.Options, _ *reg.SyncContext, edges map[reg.PromotionEdge]interface{},
) error {
	if !opts.SignImages {
		logrus.Info("Not signing images (--sign=false)")
		return nil
	}
	if len(edges) == 0 {
		logrus.Info("No images were promoted. Nothing to sign.")
		return nil
	}

	// Options for the new signer
	signOpts := defaultSignerOptions(opts)

	// Get the identity token we will use
	token, err := di.GetIdentityToken(opts, opts.SignerAccount)
	if err != nil {
		return fmt.Errorf("generating identity token: %w", err)
	}
	signOpts.IdentityToken = token

	// Creating a new Signer after setting the identity token is MANDATORY
	// because that's the only way to propagate the identity token to the
	// internal Signer structs. Without that, the identity token wouldn't be
	// used at all and images would be signed with a wrong identity.
	di.signer = sign.New(signOpts)

	// We only sign the first normalized image per digest of each edge.
	type key struct {
		identity string
		digest   image.Digest
	}
	sortedEdges := map[key][]reg.PromotionEdge{}
	for edge := range edges {
		// Skip signing the signature, sbom and attestation layers
		if strings.HasSuffix(string(edge.DstImageTag.Tag), ".sig") ||
			strings.HasSuffix(string(edge.DstImageTag.Tag), ".att") ||
			edge.DstImageTag.Tag == "" {
			continue
		}

		k := key{identity: targetIdentity(&edge), digest: edge.Digest}
		if _, ok := sortedEdges[k]; !ok {
			sortedEdges[k] = []reg.PromotionEdge{}
		}
		sortedEdges[k] = append(sortedEdges[k], edge)
	}

	t := throttler.New(opts.MaxSignatureOps, len(sortedEdges))
	// Sign the required edges
	for d := range sortedEdges {
		d := d
		go func(k key) {
			t.Done(di.signAndReplicate(signOpts, k.identity, sortedEdges[k]))
		}(d)
		if t.Throttle() > 0 {
			break
		}
	}

	return t.Err()
}

func (di *DefaultPromoterImplementation) signAndReplicate(signOpts *sign.Options, identity string, edges []reg.PromotionEdge) error {
	// Build the reference we will use
	firstEdge := &edges[0]
	imageRef := firstEdge.DstReference()

	// Make a shallow copy so we can safely modify the options per go routine
	signOptsCopy := *signOpts

	// Update the production container identity (".critical.identity.docker-reference")
	signOptsCopy.SignContainerIdentity = identity
	logrus.Infof("Using new production registry reference for %s: %v", imageRef, identity)

	// Add an annotation recording the kpromo version to ensure we
	// get a 2nd signature, otherwise cosign will not resign a signed image:
	signOptsCopy.Annotations = []string{
		fmt.Sprintf("org.kubernetes.kpromo.version=kpromo-%s", version.GetVersionInfo().GitVersion),
	}

	logrus.Infof("Signing image %s", imageRef)

	// Carry over existing signatures from the staging repo
	if err := di.copyAttachedObjects(firstEdge); err != nil {
		return fmt.Errorf("copying staging signatures: %w", err)
	}

	// Sign the first promoted image in the edges list:
	if _, err := di.signer.SignImageWithOptions(&signOptsCopy, imageRef); err != nil {
		return fmt.Errorf("signing image %s: %w", imageRef, err)
	}

	// If the same digest was promoted to more than one
	// registry, copy the signature from the first one
	if len(edges) == 1 {
		logrus.WithField("image", string(edges[0].Digest)).Debug(
			"Not replicating signatures, image promoted to single registry",
		)
		return nil
	}

	if err := di.replicateSignatures(
		firstEdge, edges[1:],
	); err != nil {
		return fmt.Errorf("replicating signatures: %w", err)
	}
	return nil
}

// targetIdentity returns the production identity for a promotion edge.
//
// This means we will substitute the .critical.identity.docker-reference within
// an example signature:
// 'us-west2-docker.pkg.dev/k8s-artifacts-prod/images/kubernetes/conformance-arm64'
//
// to match the production registry:
// 'registry.k8s.io/kubernetes/conformance-arm64'
func targetIdentity(edge *reg.PromotionEdge) string {
	identity := fmt.Sprintf("%s/%s", edge.DstRegistry.Name, edge.DstImageTag.Name)

	if !strings.Contains(string(edge.DstRegistry.Name), productionRepositoryPath) {
		logrus.Infof(
			"No production registry path %q used in image, not modifying target signature reference",
			productionRepositoryPath,
		)
		return identity
	}

	idx := strings.Index(identity, productionRepositoryPath) + len(productionRepositoryPath)
	newRef := consts.ProdRegistry + identity[idx:]

	return newRef
}

// copyAttachedObjects copies any attached signatures from the staging registry to
// the production registry. The function is called copyAttachedObjects as it will
// move attestations and SBOMs too once we stabilize the signing code.
func (di *DefaultPromoterImplementation) copyAttachedObjects(edge *reg.PromotionEdge) error {
	sigTag := digestToSignatureTag(edge.Digest)
	srcRefString := fmt.Sprintf(
		"%s/%s:%s", edge.SrcRegistry.Name, edge.SrcImageTag.Name, sigTag,
	)
	srcRef, err := name.ParseReference(srcRefString)
	if err != nil {
		return fmt.Errorf("parsing signed source reference %s: %w", srcRefString, err)
	}

	dstRefString := fmt.Sprintf(
		"%s/%s:%s", edge.DstRegistry.Name, edge.DstImageTag.Name, sigTag,
	)
	dstRef, err := name.ParseReference(dstRefString)
	if err != nil {
		return fmt.Errorf("parsing reference: %w", err)
	}

	logrus.Infof("Signature pre copy: %s to %s", srcRefString, dstRefString)
	craneOpts := []crane.Option{
		crane.WithAuthFromKeychain(gcrane.Keychain),
		crane.WithUserAgent(image.UserAgent),
		crane.WithTransport(ratelimit.Limiter),
	}

	if err := crane.Copy(srcRef.String(), dstRef.String(), craneOpts...); err != nil {
		// If the signature layer does not exist it means that the src image
		// is not signed, so we catch the error and return nil
		if strings.Contains(err.Error(), "MANIFEST_UNKNOWN") {
			logrus.Debugf("Reference %s is not signed, not copying", srcRef.String())
			return nil
		}
		return fmt.Errorf(
			"copying signature %s to %s: %w", srcRef.String(), dstRef.String(), err,
		)
	}
	return nil
}

// digestToSignatureTag takes a digest and infers the tag name where
// its signature can be found
func digestToSignatureTag(dg image.Digest) string {
	return strings.ReplaceAll(string(dg), "sha256:", "sha256-") + signatureTagSuffix
}

// replicateSignatures takes a source edge (an image) and a list of destinations
// and copies the signature to all of them
func (di *DefaultPromoterImplementation) replicateSignatures(
	src *reg.PromotionEdge, dsts []reg.PromotionEdge,
) error {
	sigTag := digestToSignatureTag(src.Digest)
	sourceRefStr := fmt.Sprintf(
		"%s/%s:%s", src.DstRegistry.Name, src.DstImageTag.Name, sigTag,
	)
	logrus.WithField("src", sourceRefStr).Infof("Replicating signature to %d images", len(dsts))
	srcRef, err := name.ParseReference(sourceRefStr)
	if err != nil {
		return fmt.Errorf("parsing reference %q: %w", sourceRefStr, err)
	}

	dstRefs := []struct {
		reference name.Reference
		token     gcloud.Token
	}{}

	for i := range dsts {
		ref, err := name.ParseReference(fmt.Sprintf(
			"%s/%s:%s", dsts[i].DstRegistry.Name, dsts[i].DstImageTag.Name, sigTag,
		))
		if err != nil {
			return fmt.Errorf("parsing signature destination referece: %w", err)
		}
		dstRefs = append(dstRefs, struct {
			reference name.Reference
			token     gcloud.Token
		}{ref, dsts[i].DstRegistry.Token})
	}

	// Copy the signatures to the missing registries
	for _, dstRef := range dstRefs {
		logrus.WithField("src", srcRef.String()).Infof("replication > %s", dstRef.reference.String())
		opts := []crane.Option{
			crane.WithAuthFromKeychain(gcrane.Keychain),
			crane.WithUserAgent(image.UserAgent),
			crane.WithTransport(ratelimit.Limiter),
		}
		if err := crane.Copy(srcRef.String(), dstRef.reference.String(), opts...); err != nil {
			return fmt.Errorf(
				"copying signature %s to %s: %w",
				srcRef.String(), dstRef.reference.String(), err,
			)
		}
	}

	return nil
}

// WriteSBOMs writes SBOMs to each of the newly promoted images and stores
// them along the signatures in the registry
func (di *DefaultPromoterImplementation) WriteSBOMs(
	_ *options.Options, _ *reg.SyncContext, _ map[reg.PromotionEdge]interface{},
) error {
	return nil
}

// GetIdentityToken returns an identity token for the selected service account
// in order for this function to work, an account has to be already logged. This
// can be achieved using the
func (di *DefaultPromoterImplementation) GetIdentityToken(
	opts *options.Options, serviceAccount string,
) (tok string, err error) {
	credOptions := []gopts.ClientOption{}
	// If the test signer file is found switch to test credentials
	if os.Getenv("CIP_E2E_KEY_FILE") != "" {
		logrus.Info("Test keyfile set using e2e test credentials")
		// ... and also use the e2e signing identity
		serviceAccount = TestSigningAccount
		credOptions = []gopts.ClientOption{
			gopts.WithCredentialsFile(os.Getenv("CIP_E2E_KEY_FILE")),
		}
	}

	// If SignerInitCredentials, initialize the iam client using
	// the identityu in that file instead of Default Application Credentials
	if opts.SignerInitCredentials != "" {
		logrus.Infof("Using credentials from %s", opts.SignerInitCredentials)
		credOptions = []gopts.ClientOption{
			gopts.WithCredentialsFile(opts.SignerInitCredentials),
		}
	}
	ctx := context.Background()
	c, err := credentials.NewIamCredentialsClient(
		ctx, credOptions...,
	)
	if err != nil {
		return tok, fmt.Errorf("creating credentials token: %w", err)
	}
	defer c.Close()
	logrus.Infof("Signing identity for images will be %s", serviceAccount)
	req := &credentialspb.GenerateIdTokenRequest{
		Name:         fmt.Sprintf("projects/-/serviceAccounts/%s", serviceAccount),
		Audience:     oidcTokenAudience, // Should be set to "sigstore"
		IncludeEmail: true,
	}

	resp, err := c.GenerateIdToken(ctx, req)
	if err != nil {
		return tok, fmt.Errorf("getting error account: %w", err)
	}

	return resp.Token, nil
}

// PrewarmTUFCache initializes the TUF cache so that threads do not have to compete
// against each other creating the TUF database.
func (di *DefaultPromoterImplementation) PrewarmTUFCache() error {
	if err := tuf.Initialize(
		context.Background(), tuf.DefaultRemoteRoot, nil,
	); err != nil {
		return fmt.Errorf("initializing TUF client: %w", err)
	}
	return nil
}
