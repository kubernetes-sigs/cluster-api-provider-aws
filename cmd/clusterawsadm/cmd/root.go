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

package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cmd/alpha"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/versioninfo"
)

// RootCmd is the Cobra root command
func RootCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "clusterawsadm",
		Short: "cluster api aws management",
		Long:  `Cluster API Provider AWS commands`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}
			return nil
		},
	}
	newCmd.AddCommand(alpha.AlphaCmd())
	newCmd.AddCommand(versioninfo.VersionCmd())
	return newCmd
}

// Execute starts the process
func Execute() {
	if err := flag.CommandLine.Parse([]string{}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := RootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Honor glog flags for verbosity control
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
}
