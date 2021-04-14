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
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	amiv1 "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/api/ami/v1alpha1"
)

type ListInput struct {
	Region            string
	KubernetesVersion string
	OperatingSystem   string
}

func List(input ListInput) (*amiv1.AWSAMIList, error) {
	supportedOsList := []string{}
	if input.OperatingSystem == "" {
		supportedOsList = getSupportedOsList()
	} else {
		supportedOsList = append(supportedOsList, input.OperatingSystem)
	}
	imageRegionList := []string{}
	if input.Region == "" {
		imageRegionList = getimageRegionList()
	} else {
		imageRegionList = append(imageRegionList, input.Region)
	}

	supportedVersions := []string{}
	if input.KubernetesVersion == "" {
		var err error
		supportedVersions, err = getSupportedKubernetesVersions()
		if err != nil {
			fmt.Println("Failed to calculate supported Kubernetes versions")
			return nil, err
		}
	} else {
		supportedVersions = append(supportedVersions, input.KubernetesVersion)
	}
	listByVersion := amiv1.AWSAMIList{
		TypeMeta: metav1.TypeMeta{
			Kind:       amiv1.AWSAMIListKind,
			APIVersion: amiv1.SchemeGroupVersion.String(),
		},
		Items: []amiv1.AWSAMI{},
	}
	for _, region := range imageRegionList {
		imageMap := make(map[string][]*ec2.Image)
		sess, err := session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
			Config:            aws.Config{Region: aws.String(region)},
		})
		if err != nil {
			return nil, err
		}

		ec2Client := ec2.New(sess)
		imagesForRegion, err := getAllImages(ec2Client, "")
		if err != nil {
			return nil, err
		}

		for key, image := range imagesForRegion {
			images, ok := imageMap[key]
			if !ok {
				images = make([]*ec2.Image, 0)
			}
			imageMap[key] = append(images, image...)
		}
		for _, version := range supportedVersions {
			for _, os := range supportedOsList {
				image, err := findAMI(imageMap, os, version)
				if err != nil {
					return nil, err
				}
				creationTimestamp, err := time.Parse(time.RFC3339, aws.StringValue(image.CreationDate))
				if err != nil {
					return nil, err
				}

				listByVersion.Items = append(listByVersion.Items, amiv1.AWSAMI{
					TypeMeta: metav1.TypeMeta{
						Kind:       amiv1.AWSAMIKind,
						APIVersion: amiv1.SchemeGroupVersion.String(),
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:              aws.StringValue(image.Name),
						CreationTimestamp: metav1.NewTime(creationTimestamp),
					},
					Spec: amiv1.AWSAMISpec{
						OS:                os,
						Region:            region,
						ImageID:           aws.StringValue(image.ImageId),
						KubernetesVersion: version,
					},
				})
			}
		}
	}

	return &listByVersion, nil
}
