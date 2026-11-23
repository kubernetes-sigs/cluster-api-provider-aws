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

package changelog

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"sigs.k8s.io/cluster-api-provider-aws/hack/tools/release-tools/internal/changelog"
)

// bodyCmd returns the `body` cobra command.
func bodyCmd() *cobra.Command {
	var input string

	cmd := &cobra.Command{
		Use:   "body",
		Short: "Print a committed CHANGELOG file's body, with any contract front-matter stripped",
		RunE: func(cmd *cobra.Command, _ []string) error {
			content, err := os.ReadFile(input) //nolint:gosec // input is an operator-supplied CLI flag
			if err != nil {
				return fmt.Errorf("reading %s: %w", input, err)
			}
			_, body := changelog.ParseFrontMatter(content)
			_, err = cmd.OutOrStdout().Write(body)
			return err
		},
	}

	cmd.Flags().StringVar(&input, "input", "", "path to the committed CHANGELOG/vX.Y.Z.md file (required)")
	_ = cmd.MarkFlagRequired("input")

	return cmd
}
