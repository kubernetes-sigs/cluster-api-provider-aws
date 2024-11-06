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
	ctrl "sigs.k8s.io/controller-runtime"

	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/ami"
	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/cmd/flags"
	cmdout "sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/printers"
)

var (
	kmsKeyID string
)

// EncryptedCopyAMICmd is a command to encrypt and copy AMI snapshots, then create an AMI with that snapshot.
func EncryptedCopyAMICmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "encrypted-copy",
		Short: "Encrypt and copy AMI snapshot, then create an AMI with that snapshot",
		Long: templates.LongDesc(`
			Find the AMI based on Kubernetes version, OS, region in the AWS account where AMIs are stored.
			Encrypt and copy the snapshot of the AMI to the current AWS account.
			Create an AMI with that snapshot.
		`),
		Example: templates.Examples(`
		# Create an encrypted AMI:
		# Available os options: centos-7, ubuntu-24.04, ubuntu-22.04, amazon-2, flatcar-stable
		clusterawsadm ami encrypted-copy --kubernetes-version=v1.18.12 --os=ubuntu-20.04  --region=us-west-2

		# owner-id and dry-run flags are optional. region can be set via flag or env
		clusterawsadm ami encrypted-copy --os centos-7 --kubernetes-version=v1.19.4 --owner-id=111111111111 --dry-run

		# copy from us-east-1 to us-east-2
		clusterawsadm ami encrypted-copy --os centos-7 --kubernetes-version=v1.19.4 --owner-id=111111111111 --region us-east-2 --source-region us-east-1

		# Encrypt using a non-default KmsKeyId specified using Key ID:
		clusterawsadm ami encrypted-copy --os centos-7 --kubernetes-version=v1.19.4 --kms-key-id=key/1234abcd-12ab-34cd-56ef-1234567890ab

		# Encrypt using a non-default KmsKeyId specified using Key alias:
		clusterawsadm ami encrypted-copy --os centos-7 --kubernetes-version=v1.19.4 --kms-key-id=alias/ExampleAlias

		# Encrypt using a non-default KmsKeyId specified using Key ARN:
		clusterawsadm ami encrypted-copy --os centos-7 --kubernetes-version=v1.19.4 --kms-key-id=arn:aws:kms:us-east-1:012345678910:key/abcd1234-a123-456a-a12b-a123b4cd56ef

		# Encrypt using a non-default KmsKeyId specified using Alias ARN:
		clusterawsadm ami encrypted-copy --os centos-7 --kubernetes-version=v1.19.4 --kms-key-id=arn:aws:kms:us-east-1:012345678910:alias/ExampleAlias
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

			log := ctrl.Log

			ami, err := ami.Copy(ami.CopyInput{
				DestinationRegion: region,
				DryRun:            dryRun,
				Encrypted:         true,
				KmsKeyID:          kmsKeyID,
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
	addKmsKeyIDFlag(newCmd)
	addSourceRegion(newCmd)
	return newCmd
}

func addKmsKeyIDFlag(c *cobra.Command) {
	c.Flags().StringVar(&kmsKeyID, "kms-key-id", "", "The ID of the KMS key for Amazon EBS encryption")
}
