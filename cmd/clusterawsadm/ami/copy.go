/*
Copyright 2021 The Kubernetes Authors.

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

package ami

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	amiv1 "sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/api/ami/v1beta1"
	ec2service "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/ec2"
	"sigs.k8s.io/cluster-api/util"
)

// CopyInput defines input that can be copied to create an AWSAMI.
type CopyInput struct {
	SourceRegion      string
	DestinationRegion string
	OwnerID           string
	OperatingSystem   string
	KubernetesVersion string
	KmsKeyID          string
	DryRun            bool
	Encrypted         bool
	Log               logr.Logger
}

// Copy will create an AWSAMI from a CopyInput.
func Copy(input CopyInput) (*amiv1.AWSAMI, error) {
	sourceCfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(input.SourceRegion))
	if err != nil {
		return nil, err
	}
	ec2Client := ec2.NewFromConfig(sourceCfg)

	image, err := ec2service.DefaultAMILookup(ec2Client, input.OwnerID, input.OperatingSystem, input.KubernetesVersion, ec2service.Amd64ArchitectureTag, "")
	if err != nil {
		return nil, err
	}

	var newImageID, newImageName string

	destCfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(input.DestinationRegion))
	if err != nil {
		return nil, err
	}

	if input.Encrypted {
		newImageName, newImageID, err = copyWithSnapshot(copyWithSnapshotInput{
			sourceRegion:      input.SourceRegion,
			image:             image,
			destinationRegion: input.DestinationRegion,
			encrypted:         input.Encrypted,
			kmsKeyID:          input.KmsKeyID,
			cfg:               destCfg,
			log:               input.Log,
		})
	} else {
		newImageName, newImageID, err = copyWithoutSnapshot(copyWithoutSnapshotInput{
			sourceRegion: input.SourceRegion,
			image:        image,
			dryRun:       input.DryRun,
			cfg:          destCfg,
			log:          input.Log,
		})
	}

	if err != nil {
		return nil, err
	}

	ami := amiv1.AWSAMI{
		ObjectMeta: metav1.ObjectMeta{
			Name:              newImageName,
			CreationTimestamp: metav1.NewTime(time.Now()),
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       amiv1.AWSAMIKind,
			APIVersion: amiv1.SchemeGroupVersion.String(),
		},
		Spec: amiv1.AWSAMISpec{
			OS:                input.OperatingSystem,
			Region:            input.DestinationRegion,
			ImageID:           newImageID,
			KubernetesVersion: input.KubernetesVersion,
		},
	}

	input.Log.Info("Completed!")

	return &ami, err
}

type copyWithoutSnapshotInput struct {
	sourceRegion string
	dryRun       bool
	log          logr.Logger
	cfg          aws.Config
	image        *types.Image
}

func copyWithoutSnapshot(input copyWithoutSnapshotInput) (string, string, error) {
	imgName := aws.ToString(input.image.Name)
	ec2Client := ec2.NewFromConfig(input.cfg)
	in2 := &ec2.CopyImageInput{
		Description:   input.image.Description,
		DryRun:        aws.Bool(input.dryRun),
		Name:          input.image.Name,
		SourceImageId: input.image.ImageId,
		SourceRegion:  aws.String(input.sourceRegion),
	}
	log := input.log.WithValues("imageName", imgName)
	log.Info("Copying the retrieved image", "imageID", aws.ToString(input.image.ImageId), "ownerID", aws.ToString(input.image.OwnerId))
	out, err := ec2Client.CopyImage(context.TODO(), in2)
	if err != nil {
		fmt.Printf("version %v\n", out)
		return imgName, "", err
	}

	return imgName, aws.ToString(out.ImageId), nil
}

type copyWithSnapshotInput struct {
	sourceRegion      string
	destinationRegion string
	kmsKeyID          string
	dryRun            bool
	encrypted         bool
	log               logr.Logger
	image             *types.Image
	cfg               aws.Config
}

func copyWithSnapshot(input copyWithSnapshotInput) (string, string, error) {
	ec2Client := ec2.NewFromConfig(input.cfg)
	imgName := *input.image.Name + util.RandomString(3) + strconv.Itoa(int(time.Now().Unix()))
	log := input.log.WithValues("imageName", imgName)

	if len(input.image.BlockDeviceMappings) == 0 || input.image.BlockDeviceMappings[0].Ebs == nil {
		return imgName, "", errors.New("image does not have EBS attached")
	}

	kmsKeyIDPtr := aws.String(input.kmsKeyID)
	if input.kmsKeyID == "" {
		kmsKeyIDPtr = nil
	}

	copyInput := &ec2.CopySnapshotInput{
		Description:      input.image.Description,
		DryRun:           aws.Bool(input.dryRun),
		Encrypted:        aws.Bool(input.encrypted),
		SourceRegion:     aws.String(input.sourceRegion),
		KmsKeyId:         kmsKeyIDPtr,
		SourceSnapshotId: input.image.BlockDeviceMappings[0].Ebs.SnapshotId,
	}

	// Generate a presigned url from the CopySnapshotInput
	scl := ec2.NewPresignClient(ec2Client)
	str, err := scl.PresignCopySnapshot(context.TODO(), copyInput, ec2.WithPresignClientFromClientOptions(
		func(o *ec2.Options) {
			o.Region = input.destinationRegion
		},
	))
	if err != nil {
		return imgName, "", errors.Wrap(err, "Failed to generate presigned url")
	}
	copyInput.PresignedUrl = aws.String(str.URL)

	out, err := ec2Client.CopySnapshot(context.TODO(), copyInput, func(o *ec2.Options) {
		o.Region = input.destinationRegion
	})
	if err != nil {
		return imgName, "", errors.Wrap(err, "Failed copying snapshot")
	}
	log.Info("Copying snapshot, this may take a couple of minutes...",
		"sourceSnapshot", aws.ToString(input.image.BlockDeviceMappings[0].Ebs.SnapshotId),
		"destinationSnapshot", aws.ToString(out.SnapshotId),
	)

	err = ec2.NewSnapshotCompletedWaiter(ec2Client).Wait(context.TODO(), &ec2.DescribeSnapshotsInput{
		DryRun:      aws.Bool(input.dryRun),
		SnapshotIds: []string{aws.ToString(out.SnapshotId)},
	}, time.Hour*1)
	if err != nil {
		return imgName, "", errors.Wrap(err, fmt.Sprintf("Failed waiting for encrypted snapshot copy completion: %q\n", aws.ToString(out.SnapshotId)))
	}

	ebsMapping := types.BlockDeviceMapping{
		DeviceName: input.image.BlockDeviceMappings[0].DeviceName,
		Ebs: &types.EbsBlockDevice{
			SnapshotId: out.SnapshotId,
		},
	}

	log.Info("Registering AMI")

	registerOut, err := ec2Client.RegisterImage(context.TODO(), &ec2.RegisterImageInput{
		Architecture:        input.image.Architecture,
		BlockDeviceMappings: []types.BlockDeviceMapping{ebsMapping},
		Description:         input.image.Description,
		DryRun:              aws.Bool(input.dryRun),
		EnaSupport:          input.image.EnaSupport,
		KernelId:            input.image.KernelId,
		Name:                aws.String(imgName),
		RamdiskId:           input.image.RamdiskId,
		RootDeviceName:      input.image.RootDeviceName,
		SriovNetSupport:     input.image.SriovNetSupport,
		VirtualizationType:  aws.String(string(input.image.VirtualizationType)),
	})

	if err != nil {
		return imgName, "", err
	}

	return imgName, aws.ToString(registerOut.ImageId), err
}
