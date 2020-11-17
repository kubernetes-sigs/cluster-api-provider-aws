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

package copy

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cmd/flags"
	ec2service "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/ec2"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/cmd"
)

func CopyAMICmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "copy",
		Short: "Copy AMI",
		Long:  cmd.LongDesc("Copy AMI"),
		Example: cmd.Examples(`
		# Copy AMI from the default AWS account where AMIs are stored.
		# Available os options: centos-7, ubuntu-18.04, ubuntu-20.04, amazon-2
		clusterawsadm ami copy --kubernetes-version=v1.18.12 --os=ubuntu-20.04  --region=us-west-2

		# source-account and dry-run flags are optional. region can be set via flag or env
		clusterawsadm ami copy --os centos-7 --kubernetes-version=1.19.4 --source-account=111111111111 --dry-run
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			region, err := flags.GetRegion(cmd)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Could not resolve AWS region, define it with --region flag or as an environment variable.")
				return err
			}

			sess, err := session.NewSessionWithOptions(session.Options{
				SharedConfigState: session.SharedConfigEnable,
				Config:            aws.Config{Region: aws.String(region)},
			})
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return err
			}
			ec2Client := ec2.New(sess)
			dryRun, err := cmd.Flags().GetBool("dry-run")
			if err != nil {
				fmt.Printf("Failed to parse dry-run value: %v. Defaulting to --dry-run=false\n", err)
			}

			os := cmd.Flags().Lookup("os").Value.String()
			if os == "" {
				return errors.Errorf("missing --os flag")
			}

			kubernetesVersion := cmd.Flags().Lookup("kubernetes-version").Value.String()
			if kubernetesVersion == "" {
				return errors.Errorf("missing --kubernetes-version flag")
			}

			ownerID := cmd.Flags().Lookup("source-account").Value.String()
			if ownerID == "" {
				fmt.Printf("Missing source-account value. Defaulting to %s\n", ec2service.DefaultMachineAMIOwnerID)
				ownerID = ec2service.DefaultMachineAMIOwnerID
			}

			image, err := ec2service.DefaultAMILookup(ec2Client, ownerID, os, kubernetesVersion, "")
			if err != nil {
				return err
			}
			in2 := &ec2.CopyImageInput{
				ClientToken:   nil,
				Description:   nil,
				DryRun:        &dryRun,
				Name:          image.Name,
				SourceImageId: image.ImageId,
				SourceRegion:  &region,
			}
			out, err := ec2Client.CopyImage(in2)
			if err != nil {
				fmt.Printf("version %q\n", out)
				return err
			}
			return nil
		},
	}

	flags.AddRegionFlag(newCmd)
	addOsFlag(newCmd)
	addKubernetesVersionFlag(newCmd)
	addDryRunFlag(newCmd)
	addSourceAccountFlag(newCmd)
	return newCmd
}

func addOsFlag(c *cobra.Command) {
	c.Flags().String("os", "", "Operating system of the AMI to be copied")
}

func addKubernetesVersionFlag(c *cobra.Command) {
	c.Flags().String("kubernetes-version", "", "Kubernetes version of the AMI to be copied")
}

func addSourceAccountFlag(c *cobra.Command) {
	c.Flags().String("source-account", "", "The source AWS account ID, where the AMI will be copied from")
}

func addDryRunFlag(c *cobra.Command) {
	c.Flags().Bool("dry-run", false, "Check if AMI exists and can be copied")
}
