/*
Copyright 2018 The Kubernetes Authors.

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
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/spf13/cobra"
	bootstrapv1 "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/api/bootstrap/v1alpha1"
	cfnBootstrap "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cloudformation/bootstrap"
	cloudformation "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cloudformation/service"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cmd/flags"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/sts"

	creds "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/credentials"
)

const backupAWSRegion = "us-east-1"

var (
	extraControlPlanePolicies []string
	extraNodePolicies         []string
)

// RootCmd is the root of the `alpha bootstrap command`.
func RootCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "bootstrap cloudformation",
		Long:  `Create and apply bootstrap AWS CloudFormation template to create IAM permissions for the Cluster API`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	newCmd.AddCommand(generateCmd())
	newCmd.AddCommand(createStackCmd())
	newCmd.AddCommand(generateIAMPolicyDocJSON())
	newCmd.AddCommand(encodeAWSSecret())
	newCmd.AddCommand(generateAWSDefaultProfileWithChain())
	newCmd.PersistentFlags().String("partition", "aws", "AWS partition, for AWS GovCloud (US) it is aws-us-gov")
	flags.MarkAlphaDeprecated(newCmd)

	return newCmd
}

func bootstrapTemplateFromCmdLine() cfnBootstrap.Template {
	conf := bootstrapv1.NewAWSIAMConfiguration()
	conf.Spec.BootstrapUser.Enable = true
	conf.Spec.ControlPlane.ExtraPolicyAttachments = extraControlPlanePolicies
	conf.Spec.Nodes.ExtraPolicyAttachments = extraNodePolicies
	return cfnBootstrap.Template{
		Spec: &conf.Spec,
	}
}

func generateCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "generate-cloudformation [AWS Account ID]",
		Short: "Generate bootstrap AWS CloudFormation template",
		Long: `Generate bootstrap AWS CloudFormation template with initial IAM policies.
You must enter an AWS account ID to generate the CloudFormation template.

Instructions for obtaining the AWS account ID can be found on https://docs.aws.amazon.com/IAM/latest/UserGuide/console_account-alias.html
`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				fmt.Printf("Error: requires AWS Account ID as an argument\n\n")
				if err := cmd.Help(); err != nil {
					return err
				}
				os.Exit(200)
			}
			if !sts.ValidateAccountID(args[0]) {
				fmt.Printf("Error: provided AWS Account ID is invalid\n\n")
				if err := cmd.Help(); err != nil {
					return err
				}
				os.Exit(201)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			template := bootstrapTemplateFromCmdLine().RenderCloudFormation()
			j, err := template.YAML()
			if err != nil {
				return err
			}

			fmt.Print(string(j))
			return nil
		},
	}

	newCmd.Flags().StringSliceVar(&extraControlPlanePolicies, "extra-controlplane-policies", []string{}, "Comma-separated list of extra policies (ARNs) to add to the created control plane role (must already exist)")
	newCmd.Flags().StringSliceVar(&extraNodePolicies, "extra-node-policies", []string{}, "Comma-separated list of extra policies (ARNs) to add to the created nodes role (must already exist)")
	flags.MarkAlphaDeprecated(newCmd)
	return newCmd
}

func createStackCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "create-stack",
		Short: "Create a new AWS CloudFormation stack using the bootstrap template",
		Long:  "Create a new AWS CloudFormation stack using the bootstrap template",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			t := bootstrapTemplateFromCmdLine()
			fmt.Printf("Attempting to create CloudFormation stack %s\n", t.Spec.StackName)
			sess, err := session.NewSessionWithOptions(session.Options{
				SharedConfigState: session.SharedConfigEnable,
			})
			if err != nil {
				fmt.Printf("Error: %v", err)
				return err
			}

			cfnSvc := cloudformation.NewService(cfn.New(sess))

			err = cfnSvc.ReconcileBootstrapStack(t.Spec.StackName, *t.RenderCloudFormation())
			if err != nil {
				fmt.Printf("Error: %v", err)
				return err
			}

			return cfnSvc.ShowStackResources(t.Spec.StackName)
		},
	}

	newCmd.Flags().StringSliceVar(&extraControlPlanePolicies, "extra-controlplane-policies", []string{}, "Comma-separated list of extra policies (ARNs) to add to the created control plane role (must already exist)")
	newCmd.Flags().StringSliceVar(&extraNodePolicies, "extra-node-policies", []string{}, "Comma-separated list of extra policies (ARNs) to add to the created nodes role (must already exist)")
	flags.MarkAlphaDeprecated(newCmd)
	return newCmd
}

func generateIAMPolicyDocJSON() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "generate-iam-policy-docs [AWS Account ID] [Directory for JSON]",
		Short: "Generate PolicyDocument JSON for all ManagedIAMPolicies",
		Long:  `Generate PolicyDocument JSON for all ManagedIAMPolicies`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				fmt.Printf("Error: requires, as arguments, an AWS Account ID and a directory for the exported JSON\n\n")
				if err := cmd.Help(); err != nil {
					return err
				}
				os.Exit(300)
			}
			accountID := args[0]
			policyDocDir := args[1]

			var err error
			if !sts.ValidateAccountID(accountID) {
				fmt.Printf("Error: provided AWS Account ID is invalid\n\n")
				if err := cmd.Help(); err != nil {
					return err
				}
				os.Exit(301)
			}

			if _, err = os.Stat(policyDocDir); os.IsNotExist(err) {
				err = os.Mkdir(policyDocDir, 0o755)
				if err != nil {
					fmt.Printf("Error: failed to make directory %q, %v", policyDocDir, err)
					if err := cmd.Help(); err != nil {
						return err
					}
					os.Exit(302)
				}
			}
			if err != nil {
				fmt.Printf("Error: failed to stat directory %q, %v", policyDocDir, err)
				if err := cmd.Help(); err != nil {
					return err
				}
				os.Exit(303)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			policyDocDir := args[1]

			t := bootstrapTemplateFromCmdLine()
			err := t.GenerateManagedIAMPolicyDocuments(policyDocDir)
			if err != nil {
				return fmt.Errorf("failed to generate PolicyDocument for all ManagedIAMPolicies: %w", err)
			}

			fmt.Printf("PolicyDocument for all ManagedIAMPolicies successfully generated in JSON at %q\n", policyDocDir)
			return nil
		},
	}
	flags.MarkAlphaDeprecated(newCmd)
	return newCmd
}

func encodeAWSSecret() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "encode-aws-credentials",
		Short: "Encode AWS credentials as a base64 encoded Kubernetes secret",
		Long:  "Encode AWS credentials as a base64 encoded Kubernetes secret",
		RunE: func(cmd *cobra.Command, args []string) error {
			region, err := flags.GetRegion(cmd)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not resolve AWS region, defaulting to %s.\n", backupAWSRegion)
				region = backupAWSRegion
			}

			awsCreds, err := creds.NewAWSCredentialFromDefaultChain(region)
			if err != nil {
				return err
			}

			str, err := awsCreds.RenderBase64EncodedAWSDefaultProfile()
			if err != nil {
				return err
			}

			fmt.Println(str)

			return nil
		},
	}
	flags.MarkAlphaDeprecated(newCmd)
	flags.AddRegionFlag(newCmd)
	return newCmd
}

func generateAWSDefaultProfileWithChain() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "generate-aws-default-profile",
		Short: "Generate an AWS profile from the current environment",
		Long:  "Generate an AWS profile from the current environment for the ephemeral bootstrap cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.CredentialWarning(cmd)

			region, err := flags.GetRegion(cmd)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not resolve AWS region, defaulting to %s.\n", backupAWSRegion)
				region = backupAWSRegion
			}

			awsCreds, err := creds.NewAWSCredentialFromDefaultChain(region)
			if err != nil {
				return err
			}

			profile, err := awsCreds.RenderAWSDefaultProfile()
			if err != nil {
				return err
			}

			fmt.Println(profile)

			return nil
		},
	}
	flags.MarkAlphaDeprecated(newCmd)
	flags.AddRegionFlag(newCmd)
	return newCmd
}
