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
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/blang/semver"
	"github.com/pkg/errors"
	ec2service "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/ec2"
	"sigs.k8s.io/cluster-api/test/framework/kubernetesversions"
)

const latestStableReleaseURL = "https://dl.k8s.io/release/stable%s.txt"

func getSupportedOsList() []string {
	return []string{"centos-7", "ubuntu-18.04", "ubuntu-20.04", "amazon-2"}
}

func getimageRegionList() []string {
	return []string{
		"ap-northeast-1",
		"ap-northeast-2",
		"ap-south-1",
		"ap-southeast-1",
		"ap-southeast-2",
		"ca-central-1",
		"eu-central-1",
		"eu-west-1",
		"eu-west-2",
		"eu-west-3",
		"sa-east-1",
		"us-east-1",
		"us-east-2",
		"us-west-1",
		"us-west-2",
	}
}

func getSupportedKubernetesVersions() ([]string, error) {
	supportedVersions := make([]string, 0)
	latestVersion, err := latestStableRelease()
	if err != nil {
		return nil, err
	}

	supportedVersions = append(supportedVersions, latestVersion)
	nMinusOne, err := kubernetesversions.PreviousMinorRelease(latestVersion)
	if err != nil {
		return nil, err
	}
	supportedVersions = append(supportedVersions, nMinusOne)

	nMinusTwo, err := kubernetesversions.PreviousMinorRelease(nMinusOne)
	if err != nil {
		return nil, err
	}
	supportedVersions = append(supportedVersions, nMinusTwo)

	return supportedVersions, nil
}

// latestStableRelease fetches the latest stable Kubernetes version
// If it is a x.x.0 release, it gets the previous minor version.
func latestStableRelease() (string, error) {
	resp, err := http.Get(fmt.Sprintf(latestStableReleaseURL, ""))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	latestVersion := strings.TrimSpace(string(b))
	tagPrefix := "v"
	latestVersionSemVer, err := semver.Make(strings.TrimPrefix(latestVersion, tagPrefix))
	if err != nil {
		return "", err
	}

	// If it is the first release, use the previous version instead
	if latestVersionSemVer.Patch == 0 {
		latestVersionSemVer.Minor = latestVersionSemVer.Minor - 1
		// Address to get stable release for a particular version is: https://dl.k8s.io/release/stable-1.19.txt"
		olderVersion := fmt.Sprintf("-%v.%v", latestVersionSemVer.Major, latestVersionSemVer.Minor)
		resp, err = http.Get(fmt.Sprintf(latestStableReleaseURL, olderVersion))
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		latestVersion = strings.TrimSpace(string(b))
	}
	return latestVersion, nil
}

func getAllImages(ec2Client ec2iface.EC2API, ownerID string) (map[string][]*ec2.Image, error) {
	if ownerID == "" {
		ownerID = ec2service.DefaultMachineAMIOwnerID
	}

	describeImageInput := &ec2.DescribeImagesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("owner-id"),
				Values: []*string{aws.String(ownerID)},
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

	out, err := ec2Client.DescribeImages(describeImageInput)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch AMIs")
	}
	if len(out.Images) == 0 {
		return nil, errors.Errorf("no AMIs in the account: %q", ownerID)
	}

	imagesMap := make(map[string][]*ec2.Image)
	for _, image := range out.Images {
		arr := strings.Split(aws.StringValue(image.Name), "-")
		arr = arr[:len(arr)-2]
		name := strings.Join(arr, "-")
		images, ok := imagesMap[name]
		if !ok {
			images = make([]*ec2.Image, 0)
		}
		imagesMap[name] = append(images, image)
	}

	return imagesMap, nil
}

func findAMI(imagesMap map[string][]*ec2.Image, baseOS, kubernetesVersion string) (*ec2.Image, error) {
	amiNameFormat := "capa-ami-{{.BaseOS}}-{{.K8sVersion}}"
	amiName, err := ec2service.GenerateAmiName(amiNameFormat, baseOS, kubernetesVersion)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to process ami format: %q", amiNameFormat)
	}

	if val, ok := imagesMap[amiName]; ok {
		latestImage, err := ec2service.GetLatestImage(val)
		if err != nil {
			return nil, err
		}
		return latestImage, nil
	}

	return nil, errors.Errorf("failed to find ami %s", amiName)
}
