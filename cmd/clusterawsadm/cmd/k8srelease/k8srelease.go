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

package k8srelease

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	detect "sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/k8srelease"
	cmdout "sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/printers"
)

// Cmd builds the detect-k8s-release command.
func Cmd() *cobra.Command {
	return CmdMulti()
}

// CmdMulti builds a standalone CLI command to detect Kubernetes release versions for CAPA AMI build policy and explicit Kubernetes version inputs.
//
// Returns:
// Configured *cobra.Command for `detect-k8s-release`, including argument validation and output flags.
func CmdMulti() *cobra.Command {
	var token string
	var output string

	newCmd := &cobra.Command{
		Use:   "detect-k8s-release <version(s)|capa>",
		Short: "Detect supported Kubernetes release versions",
		Long: templates.LongDesc(`
			Query the kubernetes/kubernetes GitHub repository for stable release tags.
			Pass "capa" detect release versions for the latest 3 minor Kubernetes versions supported by CAPA according
			to the AMI publication policy, or pass one (e.g. 1.36) or more minor versions 
			(e.g. 1.32 1.30 1.28) in any order.
		`),
		Example: templates.Examples(`
			# CAPA policy mode (top 3 latest minors)
			clusterawsadm detect-k8s-release capa

			# Explicit minor mode
			clusterawsadm detect-k8s-release 1.36 1.32 1.33
		`),
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ValidateMultiArgs(args); err != nil {
				return err
			}

			printer, err := cmdout.New(output, os.Stdout)
			if err != nil {
				return fmt.Errorf("failed creating output printer: %w", err)
			}

			result, err := detect.DetectK8sVersions(token, args...)
			if err != nil {
				return err
			}

			if output == string(cmdout.PrinterTypeTable) {
				return printer.Print(result.ToTable())
			}
			return printer.Print(result)
		},
	}

	newCmd.Flags().StringVar(&token, "token", "", "GitHub personal access token (increases API rate limit)")
	newCmd.Flags().StringVarP(&output, "output", "o", "table", "Output format: table, json, or yaml")
	return newCmd
}

func ValidateMultiArgs(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("at least one argument is required")
	}
	if args[0] == "capa" && len(args) > 1 {
		return fmt.Errorf("invalid arguments: \"capa\" cannot be combined with specific minor versions")
	}
	return nil
}
