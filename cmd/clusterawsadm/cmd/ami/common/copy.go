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

package common

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/ami"
	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/cmd/flags"
	cmdout "sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/printers"
	logf "sigs.k8s.io/cluster-api/cmd/clusterctl/log"
)

// CopyAMICmd will copy AMIs from an AWS account to the AWS account which credentials are provided.
func CopyAMICmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "copy",
		Short: "Copy AMIs from an AWS account to the AWS account which credentials are provided",
		Long: templates.LongDesc(`
			Copy AMIs based on Kubernetes version, OS, region from an AWS account where AMIs are stored
            to the current AWS account (use case: air-gapped deployments)
		`),
		Example: templates.Examples(`
		# Copy AMI from the default AWS account where AMIs are stored.
		# Available os options: centos-7, ubuntu-24.04, ubuntu-22.04, amazon-2, flatcar-stable
		clusterawsadm ami copy --kubernetes-version=v1.30.1 --os=ubuntu-22.04  --region=us-west-2

		# owner-id and dry-run flags are optional. region can be set via flag or env
		clusterawsadm ami copy --os centos-7 --kubernetes-version=v1.19.4 --owner-id=111111111111 --dry-run

		# copy from us-east-1 to us-east-2
		clusterawsadm ami copy --os centos-7 --kubernetes-version=v1.19.4 --region us-east-2 --source-region us-east-1
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			printer, err := cmdout.New("yaml", os.Stdout)
			if err != nil {
				return fmt.Errorf("failed creating output printer: %w", err)
			}
			region, err := flags.GetRegionWithError(cmd)
			if err != nil {
				return err
			}
			sourceRegion, err := GetSourceRegion(cmd)
			if err != nil {
				return err
			}

			dryRun, err := cmd.Flags().GetBool("dry-run")
			if err != nil {
				fmt.Printf("Failed to parse dry-run value: %v. Defaulting to --dry-run=false\n", err)
			}

			log := logf.Log

			ami, err := ami.Copy(ami.CopyInput{
				DestinationRegion: region,
				DryRun:            dryRun,
				KubernetesVersion: kubernetesVersion,
				Log:               log,
				OperatingSystem:   opSystem,
				OwnerID:           ownerID,
				SourceRegion:      sourceRegion,
			},
			)
			if err != nil {
				fmt.Print(err)
				return err
			}

			printer.Print(ami)

			return nil
		},
	}

	flags.AddRegionFlag(newCmd)
	addOsFlag(newCmd)
	addKubernetesVersionFlag(newCmd)
	addDryRunFlag(newCmd)
	addOwnerIDFlag(newCmd)
	addSourceRegion(newCmd)
	return newCmd
}
