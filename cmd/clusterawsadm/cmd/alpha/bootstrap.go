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

package alpha

import (
	"fmt"
	"os"

	"github.com/golang/glog"

	"github.com/spf13/cobra"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/cloudformation"
)

var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "bootstrap cloudformation",
	Long:  `Create and apply bootstrap AWS CloudFormation template to create IAM permissions for the Cluster API`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		cmd.Help()
	},
}

var bootstrapGenerateCmd = &cobra.Command{
	Use:   "generate-cloudformation",
	Short: "Generate bootstrap AWS CloudFormation template",
	Long:  "Generate bootstrap AWS CloudFormation template with initial IAM policies",
	Run: func(cmd *cobra.Command, args []string) {
		template := cloudformation.BootstrapTemplate()
		j, err := template.YAML()
		if err != nil {
			glog.Error(err)
			os.Exit(1)
		}
		fmt.Print(string(j))
	},
}

var bootstrapCreateStack = &cobra.Command{
	Use:   "create-stack",
	Short: "Create a new AWS CloudFormation stack using the bootstrap template",
	Long:  "Create a new AWS CloudFormation stack using the bootstrap template",
	Run: func(cmd *cobra.Command, args []string) {
		err := cloudformation.CreateBootstrapStack()
		if err != nil {
			glog.Error(err)
			os.Exit(1)
		}
	},
}

func initBootstrap() {
	alphaCmd.AddCommand(bootstrapCmd)
	bootstrapCmd.AddCommand(bootstrapGenerateCmd)
	bootstrapCmd.AddCommand(bootstrapCreateStack)
}
