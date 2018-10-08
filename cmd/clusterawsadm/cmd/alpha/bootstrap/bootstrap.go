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

package bootstrap

import (
	"fmt"
	"os"

	"github.com/golang/glog"

	"github.com/aws/aws-sdk-go/aws/session"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	awssts "github.com/aws/aws-sdk-go/service/sts"
	"github.com/spf13/cobra"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/cloudformation"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/sts"
)

func RootCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "bootstrap cloudformation",
		Long:  `Create and apply bootstrap AWS CloudFormation template to create IAM permissions for the Cluster API`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	newCmd.AddCommand(generateCmd())
	newCmd.AddCommand(createStackCmd())
	return newCmd
}

func generateCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "generate-cloudformation [AWS Account ID]",
		Short: "Generate bootstrap AWS CloudFormation template",
		Long: `Generate bootstrap AWS CloudFormation template with initial IAM policies.
 You must enter an AWS account ID to generate the CloudFormation template.
`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			if !sts.ValidateAccountID(args[0]) {
				glog.Error("Not entered a valid AWS Account ID")
				os.Exit(500)
			}

			template := cloudformation.BootstrapTemplate(args[0])
			j, err := template.YAML()
			if err != nil {
				glog.Error(err)
				os.Exit(1)
			}
			fmt.Print(string(j))
		},
	}
	return newCmd
}

func createStackCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "create-stack",
		Short: "Create a new AWS CloudFormation stack using the bootstrap template",
		Long:  "Create a new AWS CloudFormation stack using the bootstrap template",
		Run: func(cmd *cobra.Command, args []string) {

			stackName := "cluster-api-provider-aws-sigs-k8s-io"
			sess, err := session.NewSession()
			if err != nil {
				glog.Error(err)
				os.Exit(403)
			}

			stsSvc := sts.NewService(awssts.New(sess))

			accountID, stsErr := stsSvc.AccountID()

			if stsErr != nil {
				glog.Error(stsErr)
				os.Exit(1)
			}

			cfnSvc := cloudformation.NewService(cfn.New(sess))

			err = cfnSvc.ReconcileBootstrapStack(stackName, accountID)

			showErr := cfnSvc.ShowStackResources(stackName)

			if showErr != nil {
				glog.Error(showErr)
				os.Exit(1)
			}

			if err != nil {
				glog.Error(err)
				os.Exit(1)
			}
		},
	}

	return newCmd
}
