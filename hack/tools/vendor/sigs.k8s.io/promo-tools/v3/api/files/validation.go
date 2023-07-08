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

package files

import (
	"encoding/hex"
	"fmt"
	"strings"
)

// Validate checks for semantic errors in the yaml fields (the structure of the
// yaml is checked during unmarshaling).
func (m *Manifest) Validate() error {
	if err := ValidateFilestores(m.Filestores); err != nil {
		return err
	}

	return ValidateFiles(m.Files)
}

// ValidateFilestores validates the Filestores field of the manifest.
func ValidateFilestores(filestores []Filestore) error {
	if len(filestores) == 0 {
		return fmt.Errorf("at least one filestore must be specified")
	}

	var source *Filestore
	destinationCount := 0

	for i := range filestores {
		filestore := &filestores[i]

		if filestore.Base == "" {
			return fmt.Errorf("filestore did not have base set")
		}

		// Currently we support GCS and s3 backends.
		if strings.HasPrefix(filestore.Base, GCSScheme+"://") { //nolint: revive
			// ok
		} else if strings.HasPrefix(filestore.Base, S3Scheme+"://") { //nolint: revive
			// ok
		} else {
			return fmt.Errorf(
				"filestore has unsupported scheme in base %q",
				filestore.Base)
		}

		if filestore.Src {
			if source != nil {
				return fmt.Errorf("found multiple source filestores")
			}
			source = filestore
		} else {
			destinationCount++
		}
	}
	if source == nil {
		return fmt.Errorf("source filestore not found")
	}

	if destinationCount == 0 {
		return fmt.Errorf("no destination filestores found")
	}

	return nil
}

// ValidateFiles validates the Files field of the manifest.
func ValidateFiles(files []File) error {
	if len(files) == 0 {
		return fmt.Errorf("at least one file must be specified")
	}

	for i := range files {
		f := &files[i]

		if f.Name == "" {
			return fmt.Errorf("name is required for file")
		}

		if f.SHA256 == "" {
			return fmt.Errorf("sha256 is required for file")
		}

		sha256, err := hex.DecodeString(f.SHA256)
		if err != nil {
			return fmt.Errorf("sha256 was not valid (not hex): %q", f.SHA256)
		}

		if len(sha256) != 32 {
			return fmt.Errorf("sha256 was not valid (bad length): %q", f.SHA256)
		}
	}

	return nil
}
