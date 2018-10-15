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

package ec2

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	awsec2 "github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/client"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/ec2"
)

func RootCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "ec2",
		Short: "ec2 commands",
		Long:  `EC2 related comamnds`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	newCmd.AddCommand(consoleOutputCmd())
	return newCmd
}

func consoleOutputCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "console-output [EC2 Machine ID]",
		Short: "Get console output for host",
		Long: `Get console output for host.
You must enter a Machine ID for machine you want to get console output. Machine ID can be obtained by following https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-metadata.html
`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				fmt.Printf("Error: requires EC2 Machine ID as an argument\n\n")
				cmd.Help()
				os.Exit(202)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			sess, err := session.NewSession()
			if err != nil {
				return err
			}

			instanceID, err := client.MachineInstanceID(args[0])
			if err != nil {
				return err
			}

			svc := ec2.NewService(awsec2.New(sess))
			out, err := svc.GetConsoleOutput(instanceID)
			if err != nil {
				return err
			}

			fmt.Println(out)
			return nil
		},
	}
	return newCmd
}
