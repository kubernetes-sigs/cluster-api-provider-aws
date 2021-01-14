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
	"github.com/pkg/errors"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cmd/flags"
	ec2service "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/ec2"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/cmd"
	"sigs.k8s.io/cluster-api/util"
)

var (
	kmsKeyID    string
	kmsKeyIDPtr *string
)

func EncryptedCopyAMICmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "encrypted-copy",
		Short: "Encrypt and copy AMI snapshot, then create an AMI with that snapshot",
		Long: cmd.LongDesc(`
			Find the AMI based on Kubernetes version, OS, region in the AWS account where AMIs are stored.
			Encrypt and copy the snapshot of the AMI to the current AWS account.
			Create an AMI with that snapshot.
		`),
		Example: cmd.Examples(`
		# Create an encrypted AMI:
		# Available os options: centos-7, ubuntu-18.04, ubuntu-20.04, amazon-2
		clusterawsadm ami encrypted-copy --kubernetes-version=v1.18.12 --os=ubuntu-20.04  --region=us-west-2

		# owner-id and dry-run flags are optional. region can be set via flag or env
		clusterawsadm ami encrypted-copy --os centos-7 --kubernetes-version=v1.19.4 --owner-id=111111111111 --dry-run

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

			image, err := ec2service.DefaultAMILookup(ec2Client, ownerID, opSystem, kubernetesVersion, "")
			if err != nil {
				return err
			}

			if kmsKeyID == "" {
				kmsKeyIDPtr = nil
			}
			kmsKeyIDPtr = &kmsKeyID

			if len(image.BlockDeviceMappings) == 0 || image.BlockDeviceMappings[0].Ebs == nil {
				return errors.New("image does not have EBS attached")
			}

			copyInput := &ec2.CopySnapshotInput{
				Description:       image.Description,
				DestinationRegion: &region,
				DryRun:            &dryRun,
				Encrypted:         pointer.BoolPtr(true),
				SourceRegion:      &region,
				KmsKeyId:          kmsKeyIDPtr,
				SourceSnapshotId:  image.BlockDeviceMappings[0].Ebs.SnapshotId,
			}

			out, err := ec2Client.CopySnapshot(copyInput)
			if err != nil {
				fmt.Printf("Failed copying snapshot %q\n", err)
				return err
			}
			fmt.Printf("Copying snapshot %v as snapshot %v, this may take a couple of minutes ...\n", *image.BlockDeviceMappings[0].Ebs.SnapshotId, *out.SnapshotId)

			err = ec2Client.WaitUntilSnapshotCompleted(&ec2.DescribeSnapshotsInput{
				DryRun:      &dryRun,
				SnapshotIds: []*string{out.SnapshotId},
			})
			if err != nil {
				fmt.Printf("Failed waiting for encrypted snapshot copy completion: %q\n", *out.SnapshotId)
				return err
			}

			fmt.Println("Completed!")

			ebsMapping := &ec2.BlockDeviceMapping{
				DeviceName: image.BlockDeviceMappings[0].DeviceName,
				Ebs: &ec2.EbsBlockDevice{
					SnapshotId: out.SnapshotId,
				},
			}

			imgName := *image.Name + util.RandomString(3) + strconv.Itoa(int(time.Now().Unix()))
			fmt.Printf("Creating AMI %s\n", imgName)

			out2, err := ec2Client.RegisterImage(&ec2.RegisterImageInput{
				Architecture:        image.Architecture,
				BlockDeviceMappings: []*ec2.BlockDeviceMapping{ebsMapping},
				Description:         image.Description,
				DryRun:              &dryRun,
				EnaSupport:          image.EnaSupport,
				KernelId:            image.KernelId,
				Name:                &imgName,
				RamdiskId:           image.RamdiskId,
				RootDeviceName:      image.RootDeviceName,
				SriovNetSupport:     image.SriovNetSupport,
				VirtualizationType:  image.VirtualizationType,
			})
			if err != nil {
				fmt.Printf("Failed to create AMI from encrypted snapshot: %q\n", *out.SnapshotId)
				return err
			}

			fmt.Printf("Created AMI %v from encrypted snapshot: %q\n", *out2.ImageId, *out.SnapshotId)

			return nil
		},
	}

	flags.AddRegionFlag(newCmd)
	addOsFlag(newCmd)
	addKubernetesVersionFlag(newCmd)
	addDryRunFlag(newCmd)
	addOwnerIDFlag(newCmd)
	addKmsKeyIDFlag(newCmd)
	return newCmd
}

func addKmsKeyIDFlag(c *cobra.Command) {
	c.Flags().StringVar(&kmsKeyID, "kms-key-id", "", "The ID of the KMS key for Amazon EBS encryption")
}
