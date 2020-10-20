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
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/blang/semver"
	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/record"
)

const (
	// defaultMachineAMIOwnerID is a heptio/VMware owned account. Please see:
	// https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/487
	defaultMachineAMIOwnerID = "258751437250"

	// defaultMachineAMILookupBaseOS is the default base operating system to use
	// when looking up machine AMIs
	defaultMachineAMILookupBaseOS = "ubuntu-18.04"

	// defaultAmiNameFormat is defined in the build/ directory of this project.
	// The pattern is:
	// 1. the string value `capa-ami-`
	// 2. the baseOS of the AMI, for example: ubuntu-18.04, centos-7, amazon-2
	// 3. the kubernetes version as defined by the packages produced by kubernetes/release with or without v as a prefix, for example: 1.13.0, 1.12.5-mybuild.1, v1.17.3
	// 4. a `-` followed by any additional characters
	defaultAmiNameFormat = "capa-ami-{{.BaseOS}}-?{{.K8sVersion}}-*"

	// Amazon's AMI timestamp format
	createDateTimestampFormat = "2006-01-02T15:04:05.000Z"

	// EKS AMI ID SSM Parameter name
	eksAmiSSMParameterFormat = "/aws/service/eks/optimized-ami/%s/amazon-linux-2/recommended/image_id"
)

// AMILookup contains the parameters used to template AMI names used for lookup.
type AMILookup struct {
	BaseOS     string
	K8sVersion string
}

func amiName(amiNameFormat, baseOS, kubernetesVersion string) (string, error) {
	amiNameParameters := AMILookup{baseOS, strings.TrimPrefix(kubernetesVersion, "v")}
	// revert to default if not specified
	if amiNameFormat == "" {
		amiNameFormat = defaultAmiNameFormat
	}
	var templateBytes bytes.Buffer
	template, err := template.New("amiName").Parse(amiNameFormat)
	if err != nil {
		return amiNameFormat, errors.Wrapf(err, "failed create template from string: %q", amiNameFormat)
	}
	err = template.Execute(&templateBytes, amiNameParameters)
	if err != nil {
		return amiNameFormat, errors.Wrapf(err, "failed to substitute string: %q", amiNameFormat)
	}
	return templateBytes.String(), nil
}

// defaultAMILookup returns the default AMI based on region
func (s *Service) defaultAMILookup(amiNameFormat, ownerID, baseOS, kubernetesVersion string) (string, error) {
	if amiNameFormat == "" {
		amiNameFormat = defaultAmiNameFormat
	}
	if ownerID == "" {
		ownerID = defaultMachineAMIOwnerID
	}
	if baseOS == "" {
		baseOS = defaultMachineAMILookupBaseOS
	}
	amiName, err := amiName(amiNameFormat, baseOS, kubernetesVersion)
	if err != nil {
		return "", errors.Wrapf(err, "failed to process ami format: %q", amiNameFormat)
	}
	describeImageInput := &ec2.DescribeImagesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("owner-id"),
				Values: []*string{aws.String(ownerID)},
			},
			{
				Name:   aws.String("name"),
				Values: []*string{aws.String(amiName)},
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

	out, err := s.EC2Client.DescribeImages(describeImageInput)
	if err != nil {
		record.Eventf(s.scope.InfraCluster(), "FailedDescribeImages", "Failed to find ami %q: %v", amiName, err)
		return "", errors.Wrapf(err, "failed to find ami: %q", amiName)
	}
	if len(out.Images) == 0 {
		return "", errors.Errorf("found no AMIs with the name: %q", amiName)
	}
	latestImage, err := getLatestImage(out.Images)
	if err != nil {
		s.scope.Error(err, "failed getting latest image from AMI list")
		return "", err
	}
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
// At this point all CreationDates have been checked for errors so ignoring the error is ok.
func (i images) Less(k, j int) bool {
	firstTime, _ := time.Parse(createDateTimestampFormat, aws.StringValue(i[k].CreationDate))
	secondTime, _ := time.Parse(createDateTimestampFormat, aws.StringValue(i[j].CreationDate))
	return firstTime.Before(secondTime)
}

// Swap swaps the elements with indexes i and j.
func (i images) Swap(k, j int) {
	i[k], i[j] = i[j], i[k]
}

// getLatestImage assumes imgs is not empty. Responsibility of the caller to check.
func getLatestImage(imgs []*ec2.Image) (*ec2.Image, error) {
	for _, img := range imgs {
		if _, err := time.Parse(createDateTimestampFormat, aws.StringValue(img.CreationDate)); err != nil {
			return nil, err
		}
	}
	// old to new (newest one is last)
	sort.Sort(images(imgs))
	return imgs[len(imgs)-1], nil
}

func (s *Service) defaultBastionAMILookup(region string) string {
	switch region {
	case "ap-northeast-1":
		return "ami-02b658ac34935766f"
	case "ap-northeast-2":
		return "ami-064ab8637cf33f1bb"
	case "ap-south-1":
		return "ami-03f0fd1a2ba530e75"
	case "ap-southeast-1":
		return "ami-0c8e97a27be37adfd"
	case "ap-southeast-2":
		return "ami-099c1869f33464fde"
	case "ca-central-1":
		return "ami-0c27a26eca5dc74fc"
	case "eu-central-1":
		return "ami-092391a11f8aa4b7b"
	case "eu-west-1":
		return "ami-0823c236601fef765"
	case "eu-west-2":
		return "ami-09a1e275e350acf38"
	case "eu-west-3":
		return "ami-014d8dccd70fd2632"
	case "sa-east-1":
		return "ami-0fd2c3d373788b726"
	case "us-east-1":
		return "ami-0817d428a6fb68645"
	case "us-east-2":
		return "ami-0e82959d4ed12de3f"
	case "us-west-1":
		return "ami-03fac5402e10ea93b"
	case "us-west-2":
		return "ami-07a29e5e945228fa1"
	case "eu-north-1":
		return "ami-0ede7f804d699ea83"
	case "eu-south-1":
		return "ami-06a72c0e11b785451"
	default:
		return "unknown region"
	}
}

func (s *Service) eksAMILookup(kubernetesVersion string) (string, error) {
	// format ssm parameter path properly
	formattedVersion, err := formatVersionForEKS(kubernetesVersion)
	if err != nil {
		return "", err
	}

	paramName := fmt.Sprintf(eksAmiSSMParameterFormat, formattedVersion)

	input := &ssm.GetParameterInput{
		Name: aws.String(paramName),
	}

	out, err := s.SSMClient.GetParameter(input)
	if err != nil {
		record.Eventf(s.scope.InfraCluster(), "FailedGetParameter", "Failed to get ami SSM parameter %q: %v", paramName, err)

		return "", errors.Wrapf(err, "failed to get ami SSM parameter: %q", paramName)
	}

	if out.Parameter.Value == nil {
		return "", errors.Errorf("SSM parameter returned with nil value: %q", paramName)
	}

	id := aws.StringValue(out.Parameter.Value)
	s.scope.Info("found AMI", "id", id, "version", formattedVersion)

	return id, nil
}

func formatVersionForEKS(version string) (string, error) {
	parsed, err := semver.ParseTolerant(version)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d.%d", parsed.Major, parsed.Minor), nil
}
