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

// Package ami provides a way to generate AMI commands.
package ami

import (
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	cm "sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/cmd/ami/common"
	ls "sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/cmd/ami/list"
)

// RootCmd is the root of the `ami command`.
func RootCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "ami [command]",
		Short: "AMI commands",
		Args:  cobra.NoArgs,
		Long: templates.LongDesc(`
			All AMI related actions such as:
			# Copy AMIs based on Kubernetes version, OS etc from an AWS account where AMIs are stored
            to the current AWS account (use case: air-gapped deployments)
			# (to be implemented) List available AMIs
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	newCmd.AddCommand(cm.CopyAMICmd())
	newCmd.AddCommand(cm.EncryptedCopyAMICmd())
	newCmd.AddCommand(ls.ListAMICmd())

	return newCmd
}
