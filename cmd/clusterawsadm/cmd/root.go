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

// Package cmd implements the clusterawsadm command line utility.
package cmd

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog/v2"
	"k8s.io/kubectl/pkg/util/templates"
	ctrl "sigs.k8s.io/controller-runtime"

	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/cmd/ami"
	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/cmd/bootstrap"
	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/cmd/controller"
	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/cmd/eks"
	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/cmd/gc"
	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/cmd/resource"
	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/cmd/version"
)

var (
	verbosity *int
)

// RootCmd is the Cobra root command.
func RootCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "clusterawsadm",
		Short: "Kubernetes Cluster API Provider AWS Management Utility",
		Long: templates.LongDesc(`
			clusterawsadm provides helpers for bootstrapping Kubernetes Cluster
			API Provider AWS. Use clusterawsadm to view required AWS Identity and Access Management
			(IAM) policies as JSON docs, or create IAM roles and instance profiles automatically
			using AWS CloudFormation.

			clusterawsadm additionally helps provide credentials for use with clusterctl.
		`),
		Example: templates.Examples(`
			# Create AWS Identity and Access Management (IAM) roles for use with
			# Kubernetes Cluster API Provider AWS.
			clusterawsadm bootstrap iam create-cloudformation-stack

			# Encode credentials for use with clusterctl init
			export AWS_B64ENCODED_CREDENTIALS=$(clusterawsadm bootstrap credentials encode-as-profile)
			clusterctl init --infrastructure aws
		`),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	newCmd.AddCommand(bootstrap.RootCmd())
	newCmd.AddCommand(version.Cmd(os.Stdout))
	newCmd.AddCommand(ami.RootCmd())
	newCmd.AddCommand(eks.RootCmd())
	newCmd.AddCommand(controller.RootCmd())
	newCmd.AddCommand(resource.RootCmd())
	newCmd.AddCommand(gc.RootCmd())

	return newCmd
}

// Execute starts the process.
func Execute() {
	if err := flag.CommandLine.Parse([]string{}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, "")
		os.Exit(1)
	}

	if err := RootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	verbosity = flag.CommandLine.Int("v", 2, "Set the log level verbosity.")
	_ = flag.Set("v", strconv.Itoa(*verbosity))
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	ctrl.SetLogger(klog.NewKlogr().V(*verbosity))
}
