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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/spf13/cobra"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cloudformation/bootstrap"
	cloudformation "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cloudformation/service"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cmd/bootstrap/credentials"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cmd/flags"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/cmd"
)

func printCloudFormationTemplateCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "print-cloudformation-template",
		Short: "Print cloudformation template",
		Long: cmd.LongDesc(`
			Generate and print out a CloudFormation template that can be used to
			provision AWS Identity and Access Management (IAM) policies and roles for use
			with Kubernetes Cluster API Provider AWS.
		`),
		Example: cmd.Examples(`
		# Print out the default CloudFormation template.
		clusterawsadm bootstrap iam print-cloudformation-template

		# Print out a CloudFormation template using a custom configuration.
		clusterawsadm bootstrap iam print-cloudformation-template --config bootstrap_config.yaml
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := getBootstrapTemplate(cmd)
			if err != nil {
				return err
			}

			cfnTemplate := t.RenderCloudFormation()
			yml, err := cfnTemplate.YAML()
			if err != nil {
				return err
			}

			fmt.Println(string(yml))
			return nil
		},
	}
	addConfigFlag(newCmd)

	return newCmd
}

func createCloudFormationStackCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:     "create-cloudformation-stack",
		Aliases: []string{"update-cloudformation-stack"},
		Short:   "Create or update an AWS CloudFormation stack",
		Args:    cobra.NoArgs,
		Long: cmd.LongDesc(`
	Create or update an AWS CloudFormation stack for bootstrapping Kubernetes Cluster
	API and Kubernetes AWS Identity and Access Management (IAM) permissions. To use this
	command, there must be AWS credentials loaded in this environment.
		` + credentials.CredentialHelp),
		Example: cmd.Examples(`
		# Create or update IAM roles and policies for Kubernetes using a AWS CloudFormation stack.
		clusterawsadm bootstrap iam create-cloudformation-stack

		# Create or update IAM roles and policies for Kubernetes using a AWS CloudFormation stack with a custom configuration.
		clusterawsadm bootstrap iam create-cloudformation-stack --config bootstrap_config.yaml
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := getBootstrapTemplate(cmd)
			if err != nil {
				return err
			}

			fmt.Printf("Attempting to create AWS CloudFormation stack %s\n", t.Spec.StackName)
			sess, err := session.NewSessionWithOptions(session.Options{
				SharedConfigState: session.SharedConfigEnable,
			})
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return err
			}

			cfnSvc := cloudformation.NewService(cfn.New(sess))

			err = cfnSvc.ReconcileBootstrapStack(t.Spec.StackName, *t.RenderCloudFormation())
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return err
			}

			return cfnSvc.ShowStackResources(t.Spec.StackName)
		},
	}
	addConfigFlag(newCmd)
	flags.AddRegionFlag(newCmd)
	return newCmd
}

func deleteCloudFormationStackCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:     "delete-cloudformation-stack",
		Aliases: []string{"update-cloudformation-stack"},
		Short:   "Delete an AWS CloudFormation stack",
		Args:    cobra.NoArgs,
		Long: cmd.LongDesc(`
			Delete the AWS CloudFormation stack that created AWS Identity and Access
			Management (IAM) resources for use with Kubernetes Cluster API Provider
			AWS.
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := getBootstrapTemplate(cmd)
			if err != nil {
				return err
			}

			if err := resolveTemplateRegion(t, cmd); err != nil {
				return err
			}

			fmt.Printf("Attempting to delete AWS CloudFormation stack %s\n", t.Spec.StackName)
			sess, err := session.NewSessionWithOptions(session.Options{
				SharedConfigState: session.SharedConfigEnable,
				Config:            aws.Config{Region: aws.String(t.Spec.Region)},
			})
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return err
			}

			cfnSvc := cloudformation.NewService(cfn.New(sess))

			err = cfnSvc.DeleteStack(t.Spec.StackName, nil)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return err
			}

			return cfnSvc.ShowStackResources(t.Spec.StackName)
		},
	}
	addConfigFlag(newCmd)
	flags.AddRegionFlag(newCmd)
	return newCmd
}

func resolveTemplateRegion(t *bootstrap.Template, cmd *cobra.Command) error {
	cmdLineRegion, err := flags.GetRegion(cmd)
	if t.Spec.Region == "" && err != nil {
		return err
	}
	t.Spec.Region = cmdLineRegion
	return nil
}

func renderTemplate(t *bootstrap.Template) (string, error) {
	cfnTemplate := t.RenderCloudFormation()
	yml, err := cfnTemplate.YAML()
	if err != nil {
		return "", err
	}
	return string(yml), nil
}
