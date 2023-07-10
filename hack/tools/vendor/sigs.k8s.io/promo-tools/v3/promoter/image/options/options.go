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
)

// Options capture the switches available to run the image promoter
type Options struct {
	// Threads determines how many promotion threads will run
	Threads int

	// Confirm captures a cli flag with the same name. It runs the security
	// scan and promotion when set. If false, the promoter will exit before\
	// making any modifications.
	Confirm bool

	// Use a service account when true
	UseServiceAcct bool

	// Use only the latest diff for the manifests. Works only when running in prow.
	UseProwManifestDiff bool

	// Manifest is the path of a manifest file
	Manifest string

	// ThinManifestDir is a directory of thin manifests
	ThinManifestDir string

	// Snapshot takes a registry reference and renders a textual representation of
	// how the imagtes stored there look like to the promoter.
	Snapshot string

	// SnapshotSvcAcct is the service account we use when snapshotting.
	// TODO(puerco): Check as we can simplify to just one account
	SnapshotSvcAcct string

	// ManifestBasedSnapshotOf performs a snapshot from the given manifests
	// as opposed of Snapshot which will snapshot a registry across the network
	ManifestBasedSnapshotOf string

	// KeyFiles is a string that points to file of service account keys
	KeyFiles string

	// SeverityThreshold is the level of security vulns to search for.
	SeverityThreshold int

	// JSONLogSummary signals to the promoter if it should print a JSON summary of the operation
	JSONLogSummary bool

	// OutputFormat is the format we will use for snapshots json/yaml
	OutputFormat string

	// MinimalSnapshot is used in snapshots. but im not sure
	MinimalSnapshot bool

	// SnapshotTag when set, only images with this tag will be snapshotted
	SnapshotTag string

	// ParseOnly is an options that causes the promoter to exit
	// before promoting or generating a snapshot when set to true
	ParseOnly bool

	// When tru, sign the container images using the sigstore cosign libraries
	SignImages bool

	// SignerAccount is a service account that will provide the identity
	// when signing promoted images
	SignerAccount string

	// SignerCredentials is a credentials json file to initialize the identity
	// of the signer before running. If specified, the promoter will
	// initialize its API client with the identity in the file and use it
	// to request tokens of the signer account.
	//
	// If this credentials file is not set, the promoter will attempt to generate
	// the OIDC tokens getting its identity from the default application credentials.
	SignerInitCredentials string

	// SignCheckReferences list of image references to check for signatures
	SignCheckReferences []string

	// SignCheckFix when true, fix missing signatures
	SignCheckFix bool

	// SignCheckFromDays number of days back to check for signatrures
	SignCheckFromDays int

	// SignCheckToDays complements SignCheckFromDays to enable date ranges
	SignCheckToDays int

	// SignCheckMaxImages limits the number of images to look when verifying
	SignCheckMaxImages int

	// SignCheckIdentity is the account we expect to sign all imges
	SignCheckIdentity string

	// SignCheckIssuer is the iisuer of the OIDC tokens used to identify the signer
	SignCheckIssuer string

	// MaxSignatureCopies maximum number of concurrent signature copies
	MaxSignatureCopies int

	// MaxSignatureOps maximum number of concurrent signature operations
	MaxSignatureOps int
}

var DefaultOptions = &Options{
	OutputFormat:        "yaml",
	Threads:             10,
	SeverityThreshold:   -1,
	SignImages:          true,
	SignerAccount:       "krel-trust@k8s-releng-prod.iam.gserviceaccount.com",
	SignCheckFix:        false,
	SignCheckReferences: []string{},
	SignCheckFromDays:   5,
	SignCheckIdentity:   "krel-trust@k8s-releng-prod.iam.gserviceaccount.com",
	SignCheckIssuer:     "https://accounts.google.com",
	MaxSignatureCopies:  50, // Maximum number of concurrent signature copies
	MaxSignatureOps:     50, // Maximum number of concurrent signature operations
}

func (o *Options) Validate() error {
	// If one of the snapshot options is set, manifests will not be checked
	if o.Snapshot == "" && o.ManifestBasedSnapshotOf == "" {
		if o.Manifest == "" && o.ThinManifestDir == "" {
			return errors.New("at least a manifest file or thin manifest directory have to be specified")
		}
	}
	return nil
}

// RunOptions capture the options of a run
type RunOptions struct {
	// Confirm
	Confirm bool

	// Use a service account when true
	UseServiceAcct bool
}
