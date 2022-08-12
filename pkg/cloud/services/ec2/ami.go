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

	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	capierrors "sigs.k8s.io/cluster-api/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/blang/semver"
	"github.com/pkg/errors"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/record"
)

const (
	Amazon2         = "amazon-2"          // self, eks
	Amazon2GPU      = "amazon-2-gpu"      // eks
	Ubuntu1804      = "ubuntu-18.04"      // self, eks
	Ubuntu2004      = "ubuntu-20.04"      // self, eks
	CentOS7         = "centos-7"          // self
	FlatcarStable   = "flatcar-stable"    // self
	BottleRocket    = "bottlerocket"      // eks
	Windows2019Core = "windows-2019-core" // eks
	Windows2019Full = "windows-2019-full" // eks
	Windows2022Core = "windows-2022-core" // eks (coming soon)
	Windows2022Full = "windows-2022-full" // eks (coming soon)

	AMD64 = "x86_64"
	ARM64 = "arm64"

	// DefaultMachineAMILookupArch is the default architecture to use when looking up machine AMIs.
	DefaultMachineAMILookupArch = AMD64

	// DefaultMachineAMIOwnerID is a heptio/VMware owned account. Please see:
	// https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/487
	DefaultMachineAMIOwnerID = "258751437250"

	// ubuntuOwnerID is Ubuntu owned account. Please see:
	// https://ubuntu.com/server/docs/cloud-images/amazon-ec2
	ubuntuOwnerID = "099720109477"

	// Description regex for fetching Ubuntu AMIs for bastion host.
	ubuntuImageDescription = "Canonical??Ubuntu??20.04?LTS??amd64?focal?image*"

	// DefaultMachineAMILookupBaseOS is the default base operating system to use
	// when looking up machine AMIs.
	DefaultMachineAMILookupBaseOS = Ubuntu1804

	// DefaultEKSMachineAMILookupBaseOS is the default base operating system to use
	// when looking up machine AMIs.
	DefaultEKSMachineAMILookupBaseOS = Ubuntu1804

	// DefaultAmiNameFormat is defined in the build/ directory of this project.
	// The pattern is:
	// 1. the string value `capa-ami-`
	// 2. the baseOS of the AMI, for example: ubuntu-18.04, ubuntu-20.04, centos-7, amazon-2, flatcar-stable
	// 3. the kubernetes version as defined by the packages produced by kubernetes/release with or without v as a prefix,
	// for example: 1.13.0, 1.12.5-mybuild.1, v1.17.3
	// 4. a `-` followed by any additional characters.
	DefaultAmiNameFormat = "capa-ami-{{.BaseOS}}-?{{.K8sVersion}}-*"

	// Amazon's AMI timestamp format.
	createDateTimestampFormat = "2006-01-02T15:04:05.000Z"

	// EKS AMI ID SSM Parameter name.
	eksAmiSSMParameterFormat = "/aws/service/eks/optimized-ami/%s/amazon-linux-2/recommended/image_id"

	// EKS GPU AMI ID SSM Parameter name.
	eksGPUAmiSSMParameterFormat = "/aws/service/eks/optimized-ami/%s/amazon-linux-2-gpu/recommended/image_id"
)

var machineAMINameFormats = map[string]string{
	Amazon2:         DefaultAmiNameFormat,
	Amazon2GPU:      "",
	Ubuntu1804:      DefaultAmiNameFormat,
	Ubuntu2004:      DefaultAmiNameFormat,
	BottleRocket:    "",
	Windows2019Core: "",
	Windows2019Full: "",
	Windows2022Core: "",
	Windows2022Full: "",
	CentOS7:         "",
	FlatcarStable:   DefaultAmiNameFormat,
}

// go template for AMI Name Lookup
var eksAMINameFormats = map[string]string{
	Amazon2:         "amazon-eks-node-{{.K8sVersion}}-v*",
	Amazon2GPU:      "amazon-eks-gpu-node-{{.K8sVersion}}-v*",
	Ubuntu1804:      "ubuntu-eks/k8s_{{.K8sVersion}}/images/*18.04*",
	Ubuntu2004:      "ubuntu-eks/k8s_{{.K8sVersion}}/images/*20.04*",
	BottleRocket:    "bottlerocket-aws-k8s-{{.K8sVersion}}*",
	Windows2019Core: "Windows_Server-2019-English-Core-EKS_Optimized-{{.K8sVersion}}-*", // coming soon
	Windows2019Full: "Windows_Server-2019-English-Full-EKS_Optimized-{{.K8sVersion}}-*", // coming soon
	Windows2022Core: "Windows_Server-2022-English-Core-EKS_Optimized-{{.K8sVersion}}-*", // currently unavailable
	Windows2022Full: "Windows_Server-2022-English-Full-EKS_Optimized-{{.K8sVersion}}-*", // currently unavailable
	CentOS7:         "",                                                                 // unavailable, leave empty so ok check fails and return default, add pattern if they become available
	FlatcarStable:   "",                                                                 // unavailable, leave empty so ok check fails and return default, add pattern if they become available
}

// AMILookup contains the parameters used to template AMI names used for lookup.
type AMILookup struct {
	BaseOS     string
	K8sVersion string
	Arch       string
}

// GenerateAmiName will generate an AMI name.
func GenerateAmiName(amiNameFormat, baseOS, arch, kubernetesVersion string) (string, error) {
	// revert to default if not specified
	if amiNameFormat == "" {
		amiNameFormat = DefaultAmiNameFormat
	}
	if baseOS == "" {
		baseOS = DefaultMachineAMILookupBaseOS
	}
	amiNameParameters := AMILookup{
		BaseOS:     baseOS,
		K8sVersion: strings.TrimPrefix(kubernetesVersion, "v"),
		Arch:       arch,
	}

	var templateBytes bytes.Buffer
	tpl, err := template.New("amiName").Parse(amiNameFormat)
	if err != nil {
		return amiNameFormat, errors.Wrapf(err, "failed create template from string: %q", amiNameFormat)
	}
	err = tpl.Execute(&templateBytes, amiNameParameters)
	if err != nil {
		return amiNameFormat, errors.Wrapf(err, "failed to substitute string: %q", amiNameFormat)
	}
	return templateBytes.String(), nil
}

// cleanArch since this is kube the arch options are amd64/arm64 so trying to account for what AWS is looking for
func cleanArch(arch string) string {
	switch arch {
	case AMD64:
		return arch
	case ARM64:
		return arch
	case "amd64":
		return AMD64 // only accepted in query filter is "x86_64"
	case "aarch64":
		return ARM64 // for a query filter needs to be "arm64"
	}
	return arch
}

// DefaultAMILookup will do a default AMI lookup.
func DefaultAMILookup(ec2Client ec2iface.EC2API, ownerID, baseOS, arch, kubernetesVersion, amiNameFormat string) (*ec2.Image, error) {
	if amiNameFormat == "" {
		amiNameFormat = DefaultAmiNameFormat
	}
	if ownerID == "" {
		ownerID = DefaultMachineAMIOwnerID
	}
	if baseOS == "" {
		baseOS = DefaultMachineAMILookupBaseOS
	}
	if arch == "" {
		arch = DefaultMachineAMILookupArch
	}

	// generate the name before we clean Arch in case the user expected their version expanded in the template
	amiName, err := GenerateAmiName(amiNameFormat, baseOS, arch, kubernetesVersion)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to process ami format: %q", amiNameFormat)
	}

	arch = cleanArch(arch) // account for what AWS wants in a filter

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
				Values: []*string{aws.String(arch)},
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

	out, err := ec2Client.DescribeImages(describeImageInput)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find ami: %q", amiName)
	}
	if out == nil || len(out.Images) == 0 {
		return nil, errors.Errorf("found no AMIs with the name: %q", amiName)
	}
	latestImage, err := GetLatestImage(out.Images)
	if err != nil {
		return nil, err
	}

	return latestImage, nil
}

// defaultAMIIDLookup returns the default AMI based on region.
func (s *Service) defaultAMIIDLookup(amiNameFormat, ownerID, baseOS, arch, kubernetesVersion string) (string, error) {
	latestImage, err := DefaultAMILookup(s.EC2Client, ownerID, baseOS, arch, kubernetesVersion, amiNameFormat)
	if err != nil {
		record.Eventf(s.scope.InfraCluster(), "FailedDescribeImages", "Failed to find ami for OS=%s and Kubernetes-version=%s: %v", baseOS, kubernetesVersion, err)
		return "", errors.Wrapf(err, "failed to find ami")
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

// GetLatestImage assumes imgs is not empty. Responsibility of the caller to check.
func GetLatestImage(imgs []*ec2.Image) (*ec2.Image, error) {
	for _, img := range imgs {
		if _, err := time.Parse(createDateTimestampFormat, aws.StringValue(img.CreationDate)); err != nil {
			return nil, err
		}
	}
	// old to new (newest one is last)
	sort.Sort(images(imgs))
	return imgs[len(imgs)-1], nil
}

func (s *Service) defaultBastionAMILookup() (string, error) {
	describeImageInput := &ec2.DescribeImagesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("owner-id"),
				Values: []*string{aws.String(ubuntuOwnerID)},
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
			{
				Name:   aws.String("description"),
				Values: aws.StringSlice([]string{ubuntuImageDescription}),
			},
		},
	}
	out, err := s.EC2Client.DescribeImages(describeImageInput)
	if err != nil {
		return "", errors.Wrapf(err, "failed to describe images within region: %q", s.scope.Region())
	}
	if len(out.Images) == 0 {
		return "", errors.Errorf("found no AMIs within the region: %q", s.scope.Region())
	}
	latestImage, err := GetLatestImage(out.Images)
	if err != nil {
		return "", err
	}
	return *latestImage.ImageId, nil
}

// SSMParameterLookup takes an SSM parameter name and returns the value which is presumed to be an AMI ID
func (s *Service) SSMParameterLookup(name string) (string, error) {
	input := &ssm.GetParameterInput{
		Name: aws.String(name),
	}

	out, err := s.SSMClient.GetParameter(input)
	if err != nil {
		record.Eventf(s.scope.InfraCluster(), "FailedGetParameter", "Failed to get ami SSM parameter %q: %v", name, err)

		return "", errors.Wrapf(err, "failed to get ami SSM parameter: %q", name)
	}

	if out.Parameter == nil || out.Parameter.Value == nil {
		return "", errors.Errorf("SSM parameter returned with nil value: %q", name)
	}

	id := aws.StringValue(out.Parameter.Value)
	s.scope.Info("found AMI", "id", id)

	return id, nil
}

func (s *Service) GetMachineAMIID(scope *scope.MachineScope) (string, error) {
	// Pick image from the machine configuration, or use a default one.
	if scope.AWSMachine.Spec.AMI.ID != nil { //nolint:nestif
		return *scope.AWSMachine.Spec.AMI.ID, nil
	}

	if scope.AWSMachine.Spec.ImageLookupSSMParameterName != "" {
		return s.SSMParameterLookup(scope.AWSMachine.Spec.ImageLookupSSMParameterName)
	}

	if scope.Machine.Spec.Version == nil {
		err := errors.New("Either AWSMachine's spec.ami.id or Machine's spec.version must be defined")
		scope.SetFailureReason(capierrors.CreateMachineError)
		scope.SetFailureMessage(err)
		return "", err
	}

	if *scope.AWSMachine.Spec.AMI.EKSOptimizedLookupType != "" {
		lookupAMI, err := s.eksAMILookup(*scope.Machine.Spec.Version, scope.AWSMachine.Spec.AMI.EKSOptimizedLookupType)
		if err != nil {
			return "", err
		}
		return lookupAMI, nil
	}

	imageLookupFormat := scope.AWSMachine.Spec.ImageLookupFormat
	if imageLookupFormat == "" {
		imageLookupFormat = scope.InfraCluster.ImageLookupFormat()
	}

	imageLookupOrg := scope.AWSMachine.Spec.ImageLookupOrg
	if imageLookupOrg == "" {
		imageLookupOrg = scope.InfraCluster.ImageLookupOrg()
	}

	imageLookupBaseOS := scope.AWSMachine.Spec.ImageLookupBaseOS
	if imageLookupBaseOS == "" {
		imageLookupBaseOS = scope.InfraCluster.ImageLookupBaseOS()
	}

	imageLookupArch := scope.AWSMachine.Spec.ImageLookupArch
	if imageLookupArch == "" {
		imageLookupArch = scope.InfraCluster.ImageLookupArch()
	}

	return s.defaultAMIIDLookup(imageLookupFormat, imageLookupOrg, imageLookupBaseOS, imageLookupArch, *scope.Machine.Spec.Version)
}

func (s *Service) GetMachinePoolAMIID(scope *scope.MachinePoolScope) (string, error) {
	lt := scope.AWSMachinePool.Spec.AWSLaunchTemplate

	if lt.AMI.ID != nil {
		return *lt.AMI.ID, nil
	}

	if lt.ImageLookupSSMParameterName != "" {
		return s.SSMParameterLookup(lt.ImageLookupSSMParameterName)
	}

	if scope.MachinePool.Spec.Template.Spec.Version == nil {
		err := errors.New("Either AWSMachinePool's spec.awslaunchtemplate.ami.id or MachinePool's spec.template.spec.version must be defined")
		scope.Error(err, "")
		return "", err
	}

	// deprecated method covers the case someone set this field it will fall back to the old func
	if *lt.AMI.EKSOptimizedLookupType != "" {
		return s.eksAMILookup(*scope.MachinePool.Spec.Template.Spec.Version, lt.AMI.EKSOptimizedLookupType)
	}

	imageLookupFormat := lt.ImageLookupFormat
	if imageLookupFormat == "" {
		imageLookupFormat = scope.InfraCluster.ImageLookupFormat()
	}

	imageLookupOrg := lt.ImageLookupOrg
	if imageLookupOrg == "" {
		imageLookupOrg = scope.InfraCluster.ImageLookupOrg()
	}

	imageLookupBaseOS := lt.ImageLookupBaseOS
	if imageLookupBaseOS == "" {
		imageLookupBaseOS = scope.InfraCluster.ImageLookupBaseOS()
	}

	imageLookupArch := lt.ImageLookupArch
	if imageLookupArch == "" {
		imageLookupArch = scope.InfraCluster.ImageLookupArch()
	}

	// if this is EKS we can't use CAPA defaults and this is where we have scope
	if scope.IsEKSManaged() {
		// if EKS and trying to use AMI Lookup we need to convert from the capa-ami name formats to EKS ones.
		if imageLookupOrg == "" {
			imageLookupOrg = DefaultMachineAMIOwnerID
		}
		if imageLookupBaseOS == "" {
			imageLookupBaseOS = DefaultEKSMachineAMILookupBaseOS
		}
		if imageLookupArch == "" {
			imageLookupArch = DefaultMachineAMILookupArch
		}
		if imageLookupFormat == "" {
			imageLookupFormat = s.eksGetAMINameFormat(imageLookupBaseOS)
		}
	}

	return s.defaultAMIIDLookup(imageLookupFormat, imageLookupOrg, imageLookupBaseOS, imageLookupArch, *scope.MachinePool.Spec.Template.Spec.Version)
}

func (s *Service) eksAMILookup(kubernetesVersion string, amiType *infrav1.EKSAMILookupType) (string, error) {
	// format ssm parameter path properly
	formattedVersion, err := formatVersionForEKS(kubernetesVersion)
	if err != nil {
		return "", err
	}
	if amiType == nil {
		amiType = new(infrav1.EKSAMILookupType)
	}

	var paramName string
	switch *amiType {
	case infrav1.AmazonLinuxGPU:
		paramName = fmt.Sprintf(eksGPUAmiSSMParameterFormat, formattedVersion)
	default:
		paramName = fmt.Sprintf(eksAmiSSMParameterFormat, formattedVersion)
	}

	return s.SSMParameterLookup(paramName)
}

// EKSGetAMINameFormat takes a baseOS and will return the proper name search format
func (s *Service) eksGetAMINameFormat(baseOS string) string {
	r := eksAMINameFormats[DefaultEKSMachineAMILookupBaseOS]
	if format, ok := eksAMINameFormats[baseOS]; ok {
		if format == "" {
			return r
		}
		return format
	}
	return r
}

func formatVersionForEKS(version string) (string, error) {
	parsed, err := semver.ParseTolerant(version)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d.%d", parsed.Major, parsed.Minor), nil
}
