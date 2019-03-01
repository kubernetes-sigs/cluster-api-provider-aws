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

package config

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type configuration struct {
	ClusterName             string
	AwsRegion               string
	SSHKeyName              string
	VpcID                   string
	ControlPlaneMachineType string
	NodeMachineType         string
}

const (
	rootIdentifier = "aws"
)

// RootCmd is the root of the `alpha bootstrap command`
func RootCmd() *cobra.Command {

	force := false

	newCmd := &cobra.Command{
		Use:   "config",
		Short: "configure a cluster",
		Long:  `Cluster configuration commands`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}
			return nil
		},
	}

	newCmd.PersistentFlags().BoolVarP(&force, "force", "f", force, "Overwrite if files already exist")

	newCmd.AddCommand(InitCmd(&force))

	return newCmd
}

// InitCmd is the root of the `alpha bootstrap command`
func InitCmd(force *bool) *cobra.Command {

	var clusterName string
	var sshKeyName string
	var vpcID string
	var controlPlaneMachineType string
	var nodeMachineType string

	newCmd := &cobra.Command{
		Use:   "init",
		Short: "initialize new cluster",
		Long:  `Create directory with editable manifests`,
		RunE: func(cmd *cobra.Command, args []string) error {

			// empty, err := outDirisEmpty()
			// if err != nil {
			// 	fmt.Println("Error checking directory.")
			// 	return errors.WithStack(err)
			// }

			// if !empty && !*force {
			// 	fmt.Println("cannot write to empty directory without --force flag")
			// 	return errors.WithStack(err)
			// }

			fmt.Println("Initializing this directory for cluster-api-provider-aws configuration using kustomize")

			c := &configuration{
				ClusterName:             clusterName,
				SSHKeyName:              sshKeyName,
				VpcID:                   vpcID,
				ControlPlaneMachineType: controlPlaneMachineType,
				NodeMachineType:         nodeMachineType,
			}
			//err = c.initDir()
			// if err != nil {
			// 	return errors.Wrap(err, "cannot walk asset hierarchy and initialize directory")
			// }

			fmt.Println("\nDirectory initialized for Cluster API Provider AWS.")

			return nil
		},
	}

	newCmd.Flags().StringVarP(&clusterName, "cluster-name", "n", "default", "Name of cluster")
	newCmd.Flags().StringVarP(&sshKeyName, "ssh-key-name", "k", "default", "Name of EC2 SSH Key")
	newCmd.Flags().StringVarP(&vpcID, "vpc-id", "v", "", "VPC ID (Optional)")
	newCmd.Flags().StringVarP(&controlPlaneMachineType, "control-plane-instance-type", "c", "t3.medium", "Control Plane Instance Type")
	newCmd.Flags().StringVarP(&nodeMachineType, "node-instance-type", "c", "t3.medium", "Worker Node Instance Type")

	return newCmd
}

func outDirisEmpty() (bool, error) {
	cwd, err := os.Open("out")

	if err != nil {
		errors.WithStack(err)
	}

	entries, err := cwd.Readdir(1)

	if err != nil {
		errors.WithStack(err)
	}

	if len(entries) == 1 {
		return false, nil
	}
	return true, nil
}
