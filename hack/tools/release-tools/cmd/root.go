/*
Copyright 2026 The Kubernetes Authors.

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

// Package cmd wires the top-level `release-tool` cobra command and registers
// each subcommand domain (currently just `ami`).
package cmd

import (
	"github.com/spf13/cobra"

	amicmd "sigs.k8s.io/cluster-api-provider-aws/hack/tools/release-tools/ami/cmd"
)

// Root returns the top-level `release-tool` cobra command.
func Root() *cobra.Command {
	root := &cobra.Command{
		Use:           "release-tool",
		Short:         "CAPA release helper utilities",
		Long:          "release-tool collects automation helpers used by the CAPA release process.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(amicmd.Cmd())
	return root
}
