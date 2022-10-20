/*
Copyright 2022 The Kubernetes Authors.

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
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/gomega"
)

func TestAllPatchVersions(t *testing.T) {
	tests := []struct {
		name          string
		latestVersion string
		want          []string
	}{
		{
			name:          "Should return 3 different patch versions for v1.23.3",
			latestVersion: "v1.23.3",
			want:          []string{"v1.23.2", "v1.23.1", "v1.23.0"},
		},
		{
			name:          "Should return empty patch versions for v1.23.0",
			latestVersion: "v1.23.0",
			want:          []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := allPatchesForVersion(tt.latestVersion)
			if err != nil {
				t.Fatalf("error while fetching patch versions %+v", err)
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("allPatchesForVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindAMI(t *testing.T) {
	tests := []struct {
		name      string
		imagesMap map[string][]*ec2.Image
		want      *ec2.Image
	}{
		{
			name: "find AMI based on the latest ami format",
			imagesMap: map[string][]*ec2.Image{
				"capa-ami-amazon-2-v1.23.5-*": {
					{
						ImageId:      aws.String("capa-ami-amazon-2-v1.25.3-1664536077"),
						CreationDate: aws.String("2011-02-08T17:02:31.000Z"),
					},
				},
			},
			want: &ec2.Image{
				ImageId:      aws.String("capa-ami-amazon-2-v1.25.3-1664536077"),
				CreationDate: aws.String("2011-02-08T17:02:31.000Z"),
			},
		},
		{
			name: "find AMI based on the old ami format",
			imagesMap: map[string][]*ec2.Image{
				"capa-ami-amazon-2-1.23.5": {
					{
						ImageId:      aws.String("capa-ami-amazon-2-1.25.3-00-1664536077"),
						CreationDate: aws.String("2011-02-08T17:02:31.000Z"),
					},
				},
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			got, err := findAMI(tt.imagesMap, "amazon-2", "v1.23.5")
			if err != nil {
				t.Fatalf("error while finding AMI %+v", err)
			}
			g.Expect(got).Should(Equal(tt.want))
		})
	}
}
