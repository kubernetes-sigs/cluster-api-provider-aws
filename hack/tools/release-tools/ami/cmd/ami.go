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

// Package cmd assembles the `release-tool ami` subcommand group and its
// leaf subcommands.
package cmd

import (
	"github.com/spf13/cobra"

	k8sreleasecmd "sigs.k8s.io/cluster-api-provider-aws/hack/tools/release-tools/ami/cmd/k8srelease"
)

// Cmd returns the `ami` subcommand group.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ami",
		Short: "AMI-related release helpers",
	}
	cmd.AddCommand(k8sreleasecmd.Cmd())
	return cmd
}
