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

package bootstrap

import (
	"github.com/spf13/cobra"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cmd/bootstrap/credentials"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cmd/bootstrap/iam"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/cmd"
)

// RootCmd is the root of the `alpha bootstrap command`.
func RootCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "bootstrap [command]",
		Short: "bootstrap commands",
		Args:  cobra.NoArgs,
		Long: cmd.LongDesc(`
			In order to use Kubernetes Cluster API Provider AWS, an AWS account needs to
			be prepared with AWS Identity and Access Management (IAM) roles to be used by
			clusters as well as provide Kubernetes Cluster API Provider AWS with credentials
			to use to provision infrastructure.
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	newCmd.AddCommand(iam.RootCmd())
	newCmd.AddCommand(credentials.RootCmd())

	return newCmd
}
