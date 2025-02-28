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
	"k8s.io/kubectl/pkg/util/templates"

	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/cloudformation/bootstrap"
	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/converters"
)

var errInvalidDocumentName = fmt.Errorf("invalid document name, use one of: %+v", bootstrap.ManagedIAMPolicyNames)

func printPolicyCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "print-policy",
		Short: "Generate and show an IAM policy",
		Long: templates.LongDesc(`
			Generate and show an AWS Identity and Access Management (IAM) policy for
			Kubernetes Cluster API Provider AWS.
		`),
		Example: templates.Examples(`
		# Print out all the IAM policies for the Kubernetes CLuster API Provider AWS.
		clusterawsadm bootstrap iam print-policy

		# Print out the IAM policy for the Kubernetes Cluster API Provider AWS Controller.
		clusterawsadm bootstrap iam print-policy --document AWSIAMManagedPolicyControllers

		# Print out the IAM policy for the Kubernetes Cluster API Provider AWS Controller using a given configuration file.
		clusterawsadm bootstrap iam print-policy --document AWSIAMManagedPolicyControllers --config bootstrap_config.yaml

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

			if policyName == "" {
				return template.PrintPolicyDocs()
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
	addConfigFlag(newCmd)
	newCmd.Flags().String("document", "", fmt.Sprintf("which document to show: %+v", bootstrap.ManagedIAMPolicyNames))
	return newCmd
}

func getDocumentName(cmd *cobra.Command) (bootstrap.PolicyName, error) {
	val := bootstrap.PolicyName(cmd.Flags().Lookup("document").Value.String())

	if val == "" {
		return "", nil
	}

	if !val.IsValid() {
		return "", errInvalidDocumentName
	}

	return val, nil
}
