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

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"fmt"

	"github.com/sirupsen/logrus"

	reg "sigs.k8s.io/promo-tools/v4/internal/legacy/dockerregistry"
	"sigs.k8s.io/promo-tools/v4/internal/legacy/dockerregistry/registry"
	"sigs.k8s.io/promo-tools/v4/internal/legacy/dockerregistry/schema"
	impl "sigs.k8s.io/promo-tools/v4/internal/promoter/image"
	"sigs.k8s.io/promo-tools/v4/promoter/image/checkresults"
	options "sigs.k8s.io/promo-tools/v4/promoter/image/options"
)

var AllowedOutputFormats = []string{
	"csv",
	"yaml",
}

type Promoter struct {
	Options *options.Options
	impl    promoterImplementation
}

func New(opts *options.Options) *Promoter {
	return &Promoter{
		Options: options.DefaultOptions,
		impl:    impl.NewDefaultPromoterImplementation(opts),
	}
}

func (p *Promoter) SetImplementation(pi promoterImplementation) {
	p.impl = pi
}

//counterfeiter:generate . promoterImplementation

// promoterImplementation handles all the functionality in the promoter
// modes of operation.
type promoterImplementation interface {
	// General methods common to all modes of the promoter
	ValidateOptions(*options.Options) error
	ActivateServiceAccounts(*options.Options) error
	PrecheckAndExit(*options.Options, []schema.Manifest) error

	// Methods for promotion mode:
	ParseManifests(*options.Options) ([]schema.Manifest, error)
	MakeSyncContext(*options.Options, []schema.Manifest) (*reg.SyncContext, error)
	GetPromotionEdges(*reg.SyncContext, []schema.Manifest) (map[reg.PromotionEdge]interface{}, error)
	PromoteImages(*reg.SyncContext, map[reg.PromotionEdge]interface{}) error

	// Methods for snapshot mode:
	GetSnapshotSourceRegistry(*options.Options) (*registry.Context, error)
	GetSnapshotManifests(*options.Options) ([]schema.Manifest, error)
	AppendManifestToSnapshot(*options.Options, []schema.Manifest) ([]schema.Manifest, error)
	GetRegistryImageInventory(*options.Options, []schema.Manifest) (registry.RegInvImage, error)
	Snapshot(*options.Options, registry.RegInvImage) error

	// Methods for image vulnerability scans:
	ScanEdges(*options.Options, *reg.SyncContext, map[reg.PromotionEdge]interface{}) error

	// Methods for image signing
	PrewarmTUFCache() error
	ValidateStagingSignatures(map[reg.PromotionEdge]interface{}) (map[reg.PromotionEdge]interface{}, error)
	SignImages(*options.Options, *reg.SyncContext, map[reg.PromotionEdge]interface{}) error
	WriteSBOMs(*options.Options, *reg.SyncContext, map[reg.PromotionEdge]interface{}) error

	// Methods for checking signatures
	GetLatestImages(*options.Options) ([]string, error)
	GetSignatureStatus(*options.Options, []string) (checkresults.Signature, error)
	FixMissingSignatures(*options.Options, checkresults.Signature) error
	FixPartialSignatures(*options.Options, checkresults.Signature) error

	// Utility functions
	PrintVersion()
	PrintSecDisclaimer()
	PrintSection(string, bool)
}

// PromoteImages is the main method for image promotion
// it runs by taking all its parameters from a set of options.
func (p *Promoter) PromoteImages(opts *options.Options) (err error) {
	logrus.Infof("PromoteImages start")
	// Validate the options. Perhaps another image-specific
	// validation function may be needed.
	if err := p.impl.ValidateOptions(opts); err != nil {
		return fmt.Errorf("validating options: %w", err)
	}

	if err := p.impl.ActivateServiceAccounts(opts); err != nil {
		return fmt.Errorf("activating service accounts: %w", err)
	}

	// Prewarm the TUF cache with the targets and keys. This is done
	// to avoid collisions when signing and verifying in parallel
	if err := p.impl.PrewarmTUFCache(); err != nil {
		return fmt.Errorf("prewarming TUF cache: %w", err)
	}

	logrus.Infof("Parsing manifests")
	mfests, err := p.impl.ParseManifests(opts)
	if err != nil {
		return fmt.Errorf("parsing manifests: %w", err)
	}

	p.impl.PrintVersion()
	p.impl.PrintSection("START (PROMOTION)", opts.Confirm)

	logrus.Infof("Creating sync context manifests")
	sc, err := p.impl.MakeSyncContext(opts, mfests)
	if err != nil {
		return fmt.Errorf("creating sync context: %w", err)
	}

	logrus.Infof("Getting promotion edges")
	promotionEdges, err := p.impl.GetPromotionEdges(sc, mfests)
	if err != nil {
		return fmt.Errorf("filtering edges: %w", err)
	}

	// TODO: Let's rethink this option
	if opts.ParseOnly {
		logrus.Info("Manifests parsed, exiting as ParseOnly is set")
		return nil
	}

	// Verify any signatures in staged images
	logrus.Infof("Validating staging signatures")
	if _, err := p.impl.ValidateStagingSignatures(promotionEdges); err != nil {
		return fmt.Errorf("checking signtaures in staging images: %w", err)
	}

	// Check the pull request
	if !opts.Confirm {
		return p.impl.PrecheckAndExit(opts, mfests)
	}

	logrus.Infof("Promoting images")
	if err := p.impl.PromoteImages(sc, promotionEdges); err != nil {
		return fmt.Errorf("running promotion: %w", err)
	}

	logrus.Infof("Signing images")
	if err := p.impl.SignImages(opts, sc, promotionEdges); err != nil {
		return fmt.Errorf("signing images: %w", err)
	}

	logrus.Infof("Writing SBOMs")
	if err := p.impl.WriteSBOMs(opts, sc, promotionEdges); err != nil {
		return fmt.Errorf("writing SBOMs: %w", err)
	}

	logrus.Infof("Finish")
	return nil
}

// Snapshot runs the steps to output a representation in json or yaml of a registry
func (p *Promoter) Snapshot(opts *options.Options) (err error) {
	if err := p.impl.ValidateOptions(opts); err != nil {
		return fmt.Errorf("validating options: %w", err)
	}

	if err := p.impl.ActivateServiceAccounts(opts); err != nil {
		return fmt.Errorf("activating service accounts: %w", err)
	}

	p.impl.PrintVersion()
	p.impl.PrintSection("START (SNAPSHOT)", opts.Confirm)

	mfests, err := p.impl.GetSnapshotManifests(opts)
	if err != nil {
		return fmt.Errorf("getting snapshot manifests: %w", err)
	}

	mfests, err = p.impl.AppendManifestToSnapshot(opts, mfests)
	if err != nil {
		return fmt.Errorf("adding the specified manifest to the snapshot context: %w", err)
	}

	rii, err := p.impl.GetRegistryImageInventory(opts, mfests)
	if err != nil {
		return fmt.Errorf("getting registry image inventory: %w", err)
	}

	if err := p.impl.Snapshot(opts, rii); err != nil {
		return fmt.Errorf("generating snapshot: %w", err)
	}
	return nil
}

// SecurityScan runs just like an image promotion, but instead of
// actually copying the new detected images, it will run a vulnerability
// scan on them
func (p *Promoter) SecurityScan(opts *options.Options) error {
	if err := p.impl.ValidateOptions(opts); err != nil {
		return fmt.Errorf("validating options: %w", err)
	}

	if err := p.impl.ActivateServiceAccounts(opts); err != nil {
		return fmt.Errorf("activating service accounts: %w", err)
	}

	mfests, err := p.impl.ParseManifests(opts)
	if err != nil {
		return fmt.Errorf("parsing manifests: %w", err)
	}

	p.impl.PrintVersion()
	p.impl.PrintSection("START (VULN CHECK)", opts.Confirm)
	p.impl.PrintSecDisclaimer()

	sc, err := p.impl.MakeSyncContext(opts, mfests)
	if err != nil {
		return fmt.Errorf("creating sync context: %w", err)
	}

	promotionEdges, err := p.impl.GetPromotionEdges(sc, mfests)
	if err != nil {
		return fmt.Errorf("filtering edges: %w", err)
	}

	// TODO: Let's rethink this option
	if opts.ParseOnly {
		logrus.Info("Manifests parsed, exiting as ParseOnly is set")
		return nil
	}

	// Check the pull request
	if !opts.Confirm {
		return p.impl.PrecheckAndExit(opts, mfests)
	}

	if err := p.impl.ScanEdges(opts, sc, promotionEdges); err != nil {
		return fmt.Errorf("running vulnerability scan: %w", err)
	}
	return nil
}

// CheckSignatures checks the consistency of a set of images
func (p *Promoter) CheckSignatures(opts *options.Options) error {
	logrus.Info("Fetching latest promoted images")
	images, err := p.impl.GetLatestImages(opts)
	if err != nil {
		return fmt.Errorf("getting latest promoted images: %w", err)
	}

	logrus.Info("Checking signatures")
	results, err := p.impl.GetSignatureStatus(opts, images)
	if err != nil {
		return fmt.Errorf("checking signature status in images: %w", err)
	}

	if results.TotalPartial() == 0 && results.TotalUnsigned() == 0 {
		logrus.Info("Signature consistency OK!")
		return nil
	}

	logrus.Infof("Fixing %d unsigned images", results.TotalUnsigned())
	if err := p.impl.FixMissingSignatures(opts, results); err != nil {
		return fmt.Errorf("fixing missing signatures: %w", err)
	}

	logrus.Infof("Fixing %d images with partial signatures", results.TotalPartial())
	if err := p.impl.FixPartialSignatures(opts, results); err != nil {
		return fmt.Errorf("fixing partial signatures: %w", err)
	}

	return nil
}
