/*
Copyright 2020 The Kubernetes Authors.

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
	"github.com/spf13/cobra"
)

// ManifestCmd is a kpromo subcommand which just holds further subcommands
var ManifestCmd = &cobra.Command{
	Use:           "manifest",
	Short:         "Generate/modify a manifest for artifact promotion",
	SilenceUsage:  true,
	SilenceErrors: true,
}
