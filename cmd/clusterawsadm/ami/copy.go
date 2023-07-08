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
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
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
	sourceSession, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            aws.Config{Region: aws.String(input.SourceRegion)},
	})
	if err != nil {
		return nil, err
	}
	ec2Client := ec2.New(sourceSession)

	image, err := ec2service.DefaultAMILookup(ec2Client, input.OwnerID, input.OperatingSystem, input.KubernetesVersion, ec2service.Amd64ArchitectureTag, "")
	if err != nil {
		return nil, err
	}

	var newImageID, newImageName string

	destSession, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            aws.Config{Region: aws.String(input.DestinationRegion)},
	})
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
			sess:              destSession,
			log:               input.Log,
		})
	} else {
		newImageName, newImageID, err = copyWithoutSnapshot(copyWithoutSnapshotInput{
			sourceRegion: input.SourceRegion,
			image:        image,
			dryRun:       input.DryRun,
			sess:         destSession,
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

	if err == nil {
		input.Log.Info("Completed!")
	}

	return &ami, err
}

type copyWithoutSnapshotInput struct {
	sourceRegion string
	dryRun       bool
	log          logr.Logger
	sess         *session.Session
	image        *ec2.Image
}

func copyWithoutSnapshot(input copyWithoutSnapshotInput) (string, string, error) {
	imgName := aws.StringValue(input.image.Name)
	ec2Client := ec2.New(input.sess)
	in2 := &ec2.CopyImageInput{
		Description:   input.image.Description,
		DryRun:        aws.Bool(input.dryRun),
		Name:          input.image.Name,
		SourceImageId: input.image.ImageId,
		SourceRegion:  aws.String(input.sourceRegion),
	}
	log := input.log.WithValues("imageName", imgName)
	log.Info("Copying the retrieved image", "imageID", aws.StringValue(input.image.ImageId), "ownerID", aws.StringValue(input.image.OwnerId))
	out, err := ec2Client.CopyImage(in2)
	if err != nil {
		fmt.Printf("version %q\n", out)
		return imgName, "", err
	}

	return imgName, aws.StringValue(out.ImageId), nil
}

type copyWithSnapshotInput struct {
	sourceRegion      string
	destinationRegion string
	kmsKeyID          string
	dryRun            bool
	encrypted         bool
	log               logr.Logger
	image             *ec2.Image
	sess              *session.Session
}

func copyWithSnapshot(input copyWithSnapshotInput) (string, string, error) {
	ec2Client := ec2.New(input.sess)
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
		Description:       input.image.Description,
		DestinationRegion: aws.String(input.destinationRegion),
		DryRun:            aws.Bool(input.dryRun),
		Encrypted:         aws.Bool(input.encrypted),
		SourceRegion:      aws.String(input.sourceRegion),
		KmsKeyId:          kmsKeyIDPtr,
		SourceSnapshotId:  input.image.BlockDeviceMappings[0].Ebs.SnapshotId,
	}

	// Generate a presigned url from the CopySnapshotInput
	req, _ := ec2Client.CopySnapshotRequest(copyInput)
	str, err := req.Presign(15 * time.Minute)
	if err != nil {
		return imgName, "", errors.Wrap(err, "Failed to generate presigned url")
	}
	copyInput.PresignedUrl = aws.String(str)

	out, err := ec2Client.CopySnapshot(copyInput)
	if err != nil {
		return imgName, "", errors.Wrap(err, "Failed copying snapshot")
	}
	log.Info("Copying snapshot, this may take a couple of minutes...",
		"sourceSnapshot", aws.StringValue(input.image.BlockDeviceMappings[0].Ebs.SnapshotId),
		"destinationSnapshot", aws.StringValue(out.SnapshotId),
	)

	err = ec2Client.WaitUntilSnapshotCompleted(&ec2.DescribeSnapshotsInput{
		DryRun:      aws.Bool(input.dryRun),
		SnapshotIds: []*string{out.SnapshotId},
	})
	if err != nil {
		return imgName, "", errors.Wrap(err, fmt.Sprintf("Failed waiting for encrypted snapshot copy completion: %q\n", aws.StringValue(out.SnapshotId)))
	}

	ebsMapping := &ec2.BlockDeviceMapping{
		DeviceName: input.image.BlockDeviceMappings[0].DeviceName,
		Ebs: &ec2.EbsBlockDevice{
			SnapshotId: out.SnapshotId,
		},
	}

	log.Info("Registering AMI")

	registerOut, err := ec2Client.RegisterImage(&ec2.RegisterImageInput{
		Architecture:        input.image.Architecture,
		BlockDeviceMappings: []*ec2.BlockDeviceMapping{ebsMapping},
		Description:         input.image.Description,
		DryRun:              aws.Bool(input.dryRun),
		EnaSupport:          input.image.EnaSupport,
		KernelId:            input.image.KernelId,
		Name:                aws.String(imgName),
		RamdiskId:           input.image.RamdiskId,
		RootDeviceName:      input.image.RootDeviceName,
		SriovNetSupport:     input.image.SriovNetSupport,
		VirtualizationType:  input.image.VirtualizationType,
	})

	if err != nil {
		return imgName, "", err
	}

	return imgName, aws.StringValue(registerOut.ImageId), err
}
