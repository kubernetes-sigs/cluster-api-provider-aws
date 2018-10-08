// Copyright © 2018 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cmd/alpha"
)

func RootCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "clusterawsadm",
		Short: "cluster api aws management",
		Long:  `Cluster API Provider AWS commands`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
			cmd.Help()
		},
	}
	newCmd.AddCommand(alpha.AlphaCmd())
	return newCmd
}

// Execute starts the process
func Execute() {
	if err := RootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	flag.CommandLine.Parse([]string{})
	flag.Set("v", "2")

	// Honor glog flags for verbosity control
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
}
