/*
Copyright 2021 The Kubernetes Authors.

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

package manifest

import (
	"errors"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"sigs.k8s.io/promo-tools/v4/promobot"
)

// validateCmd takes a set of manifests and checks them
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate a promotion manifest",
	Long: `validate reads and verifies the promoter manifests in a given directory.

There are two modes of operation:

 kpromo manifests validate manifestsPath/
 kpromo manifests validate filestores.yaml [ files.yaml | files/ ]

The first invocation reads a directory containing multiple projects. kpromo will 
detect all manifests in the directory by searching for the following structure:

	├── filestores
	│   ├── project1
	│   │   └── filepromoter-manifest.yaml
	│   └── project2
	│       └── filepromoter-manifest.yaml
	└── manifests
	    ├── project1
	    │   ├── blue.yaml
	    │   ├── green.yaml
	    │   └── red.yaml
	    └── project2
	        ├── blue.yaml
	        ├── green.yaml
	        └── red.yaml

kpromo will read every manifest and validate it, failing on the first error.

The second invocation takes two or more arguments: the path to a yaml file defining
the  filestores for a project and one or more paths to files or directories where
the yaml file manifests can be found. For example:

	├── filepromoter-manifest.yaml
	└── files
		├── blue.yaml
		├── green.yaml
		└── red.yaml



`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Use = "validate [ manifestsPath | filestores.yaml files1.yaml ]"
			return cmd.Help()
		}
		if err := runValidate(args); err != nil {
			return fmt.Errorf("validating manifests: %w", err)
		}
		return nil
	},
}

func init() {
	ManifestCmd.AddCommand(validateCmd)
}

// runValidate checks one or more promoter manifest files
func runValidate(manifestPaths []string) error {
	// If we only got one argument, we assume it is a directory
	if len(manifestPaths) == 1 {
		return validateDirectory(manifestPaths[0])
	}

	// Otherwise it is a filestores + files combo:
	return validateSingle(manifestPaths[0], manifestPaths[1])
}

// validateSingle takes a filestores manifest and a path to read files from
func validateSingle(filestoresPath, filesPath string) error {
	i, err := os.Stat(filestoresPath)
	if err != nil {
		return fmt.Errorf("opening filestores manifest file: %w", err)
	}
	if i.IsDir() {
		return errors.New("first argument has to be a yaml file defining the filestores")
	}

	m, err := promobot.ReadManifest(promobot.PromoteFilesOptions{
		FilestoresPath: filestoresPath,
		FilesPath:      filesPath,
	})
	if err != nil {
		return fmt.Errorf("reading manifest from %s and %s: %w", filestoresPath, filesPath, err)
	}

	if err := m.Validate(); err != nil {
		return err
	}
	logrus.Infof("Manifest correctly validated from:\n - FileStores: %s\n - Files: %s", filestoresPath, filesPath)
	return nil
}

// validateDirectory validates a directory containing multiple projects
func validateDirectory(mPath string) error {
	i, err := os.Stat(mPath)
	if err != nil {
		return fmt.Errorf("opening manifest path %s: %w", mPath, err)
	}
	if !i.IsDir() {
		return fmt.Errorf("path is not a directory: %s", mPath)
	}

	manifests, err := promobot.ReadManifests(promobot.PromoteFilesOptions{
		ManifestsPath: mPath,
	})
	if err != nil {
		return fmt.Errorf("reading manifests from %s: %w", mPath, err)
	}
	for i, m := range manifests {
		if err := m.Validate(); err != nil {
			return fmt.Errorf("validating manifest %d from %s: %w", i, mPath, err)
		}
	}
	logrus.Infof("%d correct manifests found in %s", len(manifests), mPath)
	return nil
}
