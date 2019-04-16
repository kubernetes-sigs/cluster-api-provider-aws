/*
Copyright 2019 The Kubernetes Authors.

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

package ec2

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
)

const (
	// machineAMIOwnerID is a heptio/VMware owned account. Please see:
	// https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/487
	machineAMIOwnerID = "258751437250"

	// amiNameFormat is defined in the build/ directory of this project.
	// The pattern is:
	// 1. the string value `ami-`
	// 2. the baseOS of the AMI, for example: ubuntu, centos, amazon
	// 3. the version of the baseOS, for example: 18.04 (ubuntu), 7 (centos), 2 (amazon)
	// 4. the kubernetes version as defined by the packages produced by kubernetes/release, for example: 1.13.0-00, 1.12.5-01
	// 5. the timestamp that the AMI was built
	amiNameFormat = "capa-ami-%s-%s-%s-??-??????????"

	// Amazon's AMI timestamp format
	createDateTimestampFormat = "2006-01-02T15:04:05.000Z"
)

func amiName(baseOS, baseOSVersion, kubernetesVersion string) string {
	return fmt.Sprintf(amiNameFormat, baseOS, baseOSVersion, strings.TrimPrefix(kubernetesVersion, "v"))
}

// defaultAMILookup returns the default AMI based on region
func (s *Service) defaultAMILookup(baseOS, baseOSVersion, kubernetesVersion string) (string, error) {
	describeImageInput := &ec2.DescribeImagesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("owner-id"),
				Values: []*string{aws.String(machineAMIOwnerID)},
			},
			{
				Name:   aws.String("name"),
				Values: []*string{aws.String(amiName(baseOS, baseOSVersion, kubernetesVersion))},
			},
			{
				Name:   aws.String("architecture"),
				Values: []*string{aws.String("x86_64")},
			},
			{
				Name:   aws.String("state"),
				Values: []*string{aws.String("available")},
			},
			{
				Name:   aws.String("virtualization-type"),
				Values: []*string{aws.String("hvm")},
			},
		},
	}

	out, err := s.scope.EC2.DescribeImages(describeImageInput)
	if err != nil {
		return "", errors.Wrapf(err, "failed to find ami: %q", amiName(baseOS, baseOSVersion, kubernetesVersion))
	}
	if len(out.Images) == 0 {
		return "", errors.Errorf("found no AMIs with the name: %q", amiName(baseOS, baseOSVersion, kubernetesVersion))
	}
	latestImage := getLatestImage(out.Images)
	s.scope.V(2).Info("Found and using an existing AMI", "ami-id", aws.StringValue(latestImage.ImageId))
	return aws.StringValue(latestImage.ImageId), nil
}

type images []*ec2.Image

// Len is the number of elements in the collection.
func (i images) Len() int {
	return len(i)
}

// Less reports whether the element with
// index i should sort before the element with index j.
// TODO(chuckha) Ignoring errors until this causes a problem
func (i images) Less(k, j int) bool {
	firstTime, err := time.Parse(createDateTimestampFormat, aws.StringValue(i[k].CreationDate))
	if err != nil {
		return false
	}
	secondTime, err := time.Parse(createDateTimestampFormat, aws.StringValue(i[j].CreationDate))
	if err != nil {
		return false
	}
	return firstTime.Before(secondTime)
}

// Swap swaps the elements with indexes i and j.
func (i images) Swap(k, j int) {
	i[k], i[j] = i[j], i[k]
}

// getLatestImage assumes imgs is not empty. Responsibility of the caller to check.
func getLatestImage(imgs []*ec2.Image) *ec2.Image {
	// old to new (newest one is last)
	sort.Sort(images(imgs))
	return imgs[len(imgs)-1]
}

func (s *Service) defaultBastionAMILookup(region string) string {
	switch region {
	case "ap-northeast-1":
		return "ami-d39a02b5"
	case "ap-northeast-2":
		return "ami-67973709"
	case "ap-south-1":
		return "ami-5d055232"
	case "ap-southeast-1":
		return "ami-325d2e4e"
	case "ap-southeast-2":
		return "ami-37df2255"
	case "ca-central-1":
		return "ami-f0870294"
	case "eu-central-1":
		return "ami-af79ebc0"
	case "eu-west-1":
		return "ami-4d46d534"
	case "eu-west-2":
		return "ami-d7aab2b3"
	case "eu-west-3":
		return "ami-5e0eb923"
	case "sa-east-1":
		return "ami-1157157d"
	case "us-east-1":
		return "ami-41e0b93b"
	case "us-east-2":
		return "ami-2581aa40"
	case "us-west-1":
		return "ami-79aeae19"
	case "us-west-2":
		return "ami-1ee65166"
	default:
		return "unknown region"
	}
}
