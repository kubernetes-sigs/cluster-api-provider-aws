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

// Package iam provides a way to generate IAM policies and roles.
package iam

import (
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"
)

// RootCmd is the root of the `bootstrap iam command`.
func RootCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "iam [command]",
		Short: "View required AWS IAM policies and create/update IAM roles using AWS CloudFormation",
		Long: templates.LongDesc(`
			View/output AWS Identity and Access Management (IAM) policy documents required for
			configuring Kubernetes Cluster API Provider AWS as well as create/update AWS IAM
			resources using AWS CloudFormation.
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	newCmd.AddCommand(printPolicyCmd())
	newCmd.AddCommand(printConfigCmd())
	newCmd.AddCommand(printCloudFormationTemplateCmd())
	newCmd.AddCommand(createCloudFormationStackCmd())
	newCmd.AddCommand(deleteCloudFormationStackCmd())
	return newCmd
}
