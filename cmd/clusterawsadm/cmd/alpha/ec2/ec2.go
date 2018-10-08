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
	"github.com/golang/glog"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	awsec2 "github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/ec2"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/client"
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
		Use:   "console-output",
		Short: "Get console output for host",
		Long:  "Get console output for host",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			sess, err := session.NewSession()
			if err != nil {
				glog.Error(err)
				os.Exit(403)
			}

			instanceID, err := client.MachineInstanceID(args[0])

			if err != nil {
				glog.Error(err)
				os.Exit(404)
			}

			svc := ec2.NewService(awsec2.New(sess))

			out, err := svc.GetConsole(instanceID)

			if err != nil {
				glog.Error(err)
				os.Exit(500)
			}

			fmt.Println(out)

		},
	}
	return newCmd
}
