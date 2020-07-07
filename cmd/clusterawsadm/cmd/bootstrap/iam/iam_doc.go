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

package iam

import (
	"fmt"

	"github.com/spf13/cobra"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cloudformation/bootstrap"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/converters"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/cmd"
)

func printPolicyCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "print-policy",
		Short: "Generate and show an IAM policy",
		Long: cmd.LongDesc(`
			Generate and show an AWS Identity and Access Management (IAM) policy for
			Kubernetes Cluster API Provider AWS.
		`),
		Example: cmd.Examples(`
		# Print out the IAM policy for the Kubernetes Cluster API Provider AWS Controller.
		clusterawsadm bootstrap iam print-policy --document AWSIAMManagedPolicyControllers

		# Print out the IAM policy for the Kubernetes AWS Cloud Provider for the control plane.
		clusterawsadm bootstrap iam print-policy --document AWSIAMManagedPolicyCloudProviderControlPlane

		# Print out the IAM policy for the Kubernetes AWS Cloud Provider for all nodes.
		clusterawsadm bootstrap iam print-policy --document AWSIAMManagedPolicyCloudProviderNodes

		# Print out the IAM policy for the Kubernetes AWS EBS CSI Driver Controller.
		clusterawsadm bootstrap iam print-policy --document AWSEBSCSIPolicyController
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			template, err := getBootstrapTemplate(cmd)
			if err != nil {
				return err
			}

			policyName, err := getDocumentName(cmd)
			if err != nil {
				return err
			}

			policyDocument := template.GetPolicyDocFromPolicyName(policyName)
			str, err := converters.IAMPolicyDocumentToJSON(*policyDocument)
			if err != nil {
				return err
			}

			fmt.Println(str)
			return nil
		},
	}

	newCmd.Flags().String("document", "", fmt.Sprintf("which document to show: %+v", bootstrap.ManagedIAMPolicyNames))
	return newCmd
}

func getDocumentName(cmd *cobra.Command) (bootstrap.PolicyName, error) {
	val := bootstrap.PolicyName(cmd.Flags().Lookup("document").Value.String())
	if !val.IsValid() {
		return "", fmt.Errorf("invalid document name, use one of: %+v", bootstrap.ManagedIAMPolicyNames)
	}

	return val, nil
}
