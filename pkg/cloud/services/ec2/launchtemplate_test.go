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

package ec2

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/aws/smithy-go"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/ssm/mock_ssmiface"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/userdata"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
)

const (
	testUserData = `## template: jinja
#cloud-config

write_files:
-   path: /tmp/kubeadm-join-config.yaml
	owner: root:root
	permissions: '0640'
	content: |
	  ---
	  apiVersion: kubeadm.k8s.io/v1beta2
	  discovery:
		bootstrapToken:
		  apiServerEndpoint: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
		  caCertHashes:
		  - sha256:xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
		  token: xxxxxx.xxxxxxxxxxxxxxxx
		  unsafeSkipCAVerification: false
	  kind: JoinConfiguration
	  nodeRegistration:
		kubeletExtraArgs:
		  cloud-provider: aws
		name: '{{ ds.meta_data.local_hostname }}'

runcmd:
  - kubeadm join --config /tmp/kubeadm-join-config.yaml
users:
  - name: xxxxxxxx
	sudo: ALL=(ALL) NOPASSWD:ALL
	ssh_authorized_keys:
	  - ssh-rsa xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx user@example.com
`
)

var testUserDataHash = userdata.ComputeHash([]byte(testUserData))

var (
	testBootstrapData     = []byte("different from testUserData since bootstrap data may be in S3 while EC2 user data points to that S3 object")
	testBootstrapDataHash = userdata.ComputeHash(testBootstrapData)
)

func defaultEC2AndDataTags(name string, clusterName string, userDataSecretKey types.NamespacedName, bootstrapDataHash string) []ec2types.Tag {
	tags := defaultEC2Tags(name, clusterName)
	tags = append(
		tags,
		ec2types.Tag{
			Key:   aws.String(infrav1.LaunchTemplateBootstrapDataSecret),
			Value: aws.String(userDataSecretKey.String()),
		},
		ec2types.Tag{
			Key:   aws.String(infrav1.LaunchTemplateBootstrapDataHash),
			Value: aws.String(bootstrapDataHash),
		})

	sortTags(tags)
	return tags
}

func TestGetLaunchTemplate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name               string
		launchTemplateName string
		expect             func(m *mocks.MockEC2APIMockRecorder)
		check              func(g *WithT, launchTemplate *expinfrav1.AWSLaunchTemplate, userDataHash string, err error)
	}{
		{
			name: "Should not return launch template if empty launch template name passed",
			check: func(g *WithT, launchTemplate *expinfrav1.AWSLaunchTemplate, userDataHash string, err error) {
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(userDataHash).Should(BeEmpty())
				g.Expect(launchTemplate).Should(BeNil())
			},
		},
		{
			name:               "Should not return error if no launch template exist with given name",
			launchTemplateName: "foo",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeLaunchTemplateVersions(context.TODO(), gomock.Eq(&ec2.DescribeLaunchTemplateVersionsInput{
					LaunchTemplateName: aws.String("foo"),
					Versions:           []string{"$Latest"},
				})).
					Return(nil, &smithy.GenericAPIError{
						Code:    awserrors.LaunchTemplateNameNotFound,
						Message: "The specified launch template, with template name foo, does not exist.",
					})
			},
			check: func(g *WithT, launchTemplate *expinfrav1.AWSLaunchTemplate, userDataHash string, err error) {
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(userDataHash).Should(BeEmpty())
				g.Expect(launchTemplate).Should(BeNil())
			},
		},
		{
			name:               "Should return error if AWS failed during launch template fetching",
			launchTemplateName: "foo",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeLaunchTemplateVersions(context.TODO(), gomock.Eq(&ec2.DescribeLaunchTemplateVersionsInput{
					LaunchTemplateName: aws.String("foo"),
					Versions:           []string{"$Latest"},
				})).Return(nil, awserrors.NewFailedDependency("dependency failure"))
			},
			check: func(g *WithT, launchTemplate *expinfrav1.AWSLaunchTemplate, userDataHash string, err error) {
				g.Expect(err).To(HaveOccurred())
				g.Expect(userDataHash).Should(BeEmpty())
				g.Expect(launchTemplate).Should(BeNil())
			},
		},
		{
			name:               "Should not return with error if no launch template versions received from AWS",
			launchTemplateName: "foo",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeLaunchTemplateVersions(context.TODO(), gomock.Eq(&ec2.DescribeLaunchTemplateVersionsInput{
					LaunchTemplateName: aws.String("foo"),
					Versions:           []string{"$Latest"},
				})).Return(nil, nil)
			},
			check: func(g *WithT, launchTemplate *expinfrav1.AWSLaunchTemplate, userDataHash string, err error) {
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(userDataHash).Should(BeEmpty())
				g.Expect(launchTemplate).Should(BeNil())
			},
		},
		{
			name:               "Should successfully return launch template if exist with given name",
			launchTemplateName: "foo",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeLaunchTemplateVersions(context.TODO(), gomock.Eq(&ec2.DescribeLaunchTemplateVersionsInput{
					LaunchTemplateName: aws.String("foo"),
					Versions:           []string{"$Latest"},
				})).Return(&ec2.DescribeLaunchTemplateVersionsOutput{
					LaunchTemplateVersions: []ec2types.LaunchTemplateVersion{
						{
							LaunchTemplateId:   aws.String("lt-12345"),
							LaunchTemplateName: aws.String("foo"),
							LaunchTemplateData: &ec2types.ResponseLaunchTemplateData{
								SecurityGroupIds: []string{"sg-id"},
								ImageId:          aws.String("foo-image"),
								IamInstanceProfile: &ec2types.LaunchTemplateIamInstanceProfileSpecification{
									Arn: aws.String("instance-profile/foo-profile"),
								},
								KeyName: aws.String("foo-keyname"),
								BlockDeviceMappings: []ec2types.LaunchTemplateBlockDeviceMapping{
									{
										DeviceName: aws.String("foo-device"),
										Ebs: &ec2types.LaunchTemplateEbsBlockDevice{
											Encrypted:  aws.Bool(true),
											VolumeSize: aws.Int32(16),
											VolumeType: ec2types.VolumeTypeGp2,
										},
									},
								},
								NetworkInterfaces: []ec2types.LaunchTemplateInstanceNetworkInterfaceSpecification{
									{
										DeviceIndex: aws.Int32(1),
										Groups:      []string{"foo-group"},
									},
								},
								UserData: aws.String(base64.StdEncoding.EncodeToString([]byte(testUserData))),
							},
							VersionNumber: aws.Int64(1),
						},
					},
				}, nil)
			},
			check: func(g *WithT, launchTemplate *expinfrav1.AWSLaunchTemplate, userDataHash string, err error) {
				wantLT := &expinfrav1.AWSLaunchTemplate{
					Name: "foo",
					AMI: infrav1.AMIReference{
						ID: aws.String("foo-image"),
					},
					IamInstanceProfile:       "foo-profile",
					SSHKeyName:               aws.String("foo-keyname"),
					VersionNumber:            aws.Int64(1),
					AdditionalSecurityGroups: []infrav1.AWSResourceReference{{ID: aws.String("sg-id")}},
				}

				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(userDataHash).Should(Equal(testUserDataHash))
				g.Expect(launchTemplate).Should(Equal(wantLT))
			},
		},
		{
			name:               "Should return computed userData if AWS returns empty userData",
			launchTemplateName: "foo",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeLaunchTemplateVersions(context.TODO(), gomock.Eq(&ec2.DescribeLaunchTemplateVersionsInput{
					LaunchTemplateName: aws.String("foo"),
					Versions:           []string{"$Latest"},
				})).Return(&ec2.DescribeLaunchTemplateVersionsOutput{
					LaunchTemplateVersions: []ec2types.LaunchTemplateVersion{
						{
							LaunchTemplateId:   aws.String("lt-12345"),
							LaunchTemplateName: aws.String("foo"),
							LaunchTemplateData: &ec2types.ResponseLaunchTemplateData{
								SecurityGroupIds: []string{"sg-id"},
								ImageId:          aws.String("foo-image"),
								IamInstanceProfile: &ec2types.LaunchTemplateIamInstanceProfileSpecification{
									Arn: aws.String("instance-profile/foo-profile"),
								},
								KeyName: aws.String("foo-keyname"),
								BlockDeviceMappings: []ec2types.LaunchTemplateBlockDeviceMapping{
									{
										DeviceName: aws.String("foo-device"),
										Ebs: &ec2types.LaunchTemplateEbsBlockDevice{
											Encrypted:  aws.Bool(true),
											VolumeSize: aws.Int32(16),
											VolumeType: ec2types.VolumeTypeGp2,
										},
									},
								},
								NetworkInterfaces: []ec2types.LaunchTemplateInstanceNetworkInterfaceSpecification{
									{
										DeviceIndex: aws.Int32(1),
										Groups:      []string{"foo-group"},
									},
								},
							},
							VersionNumber: aws.Int64(1),
						},
					},
				}, nil)
			},
			check: func(g *WithT, launchTemplate *expinfrav1.AWSLaunchTemplate, userDataHash string, err error) {
				wantLT := &expinfrav1.AWSLaunchTemplate{
					Name: "foo",
					AMI: infrav1.AMIReference{
						ID: aws.String("foo-image"),
					},
					IamInstanceProfile:       "foo-profile",
					SSHKeyName:               aws.String("foo-keyname"),
					VersionNumber:            aws.Int64(1),
					AdditionalSecurityGroups: []infrav1.AWSResourceReference{{ID: aws.String("sg-id")}},
				}

				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(userDataHash).Should(Equal(userdata.ComputeHash(nil)))
				g.Expect(launchTemplate).Should(Equal(wantLT))
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			scheme, err := setupScheme()
			g.Expect(err).NotTo(HaveOccurred())
			client := fake.NewClientBuilder().WithScheme(scheme).Build()

			cs, err := setupClusterScope(client)
			g.Expect(err).NotTo(HaveOccurred())
			mockEC2Client := mocks.NewMockEC2API(mockCtrl)

			s := NewService(cs)
			s.EC2Client = mockEC2Client

			if tc.expect != nil {
				tc.expect(mockEC2Client.EXPECT())
			}

			launchTemplate, userData, _, _, err := s.GetLaunchTemplate(tc.launchTemplateName)
			tc.check(g, launchTemplate, userData, err)
		})
	}
}

func TestServiceSDKToLaunchTemplate(t *testing.T) {
	tests := []struct {
		name                  string
		input                 ec2types.LaunchTemplateVersion
		wantLT                *expinfrav1.AWSLaunchTemplate
		wantUserDataHash      string
		wantDataSecretKey     *types.NamespacedName
		wantBootstrapDataHash *string
		wantErr               bool
	}{
		{
			name: "lots of input",
			input: ec2types.LaunchTemplateVersion{
				LaunchTemplateId:   aws.String("lt-12345"),
				LaunchTemplateName: aws.String("foo"),
				LaunchTemplateData: &ec2types.ResponseLaunchTemplateData{
					ImageId: aws.String("foo-image"),
					IamInstanceProfile: &ec2types.LaunchTemplateIamInstanceProfileSpecification{
						Arn: aws.String("instance-profile/foo-profile"),
					},
					KeyName: aws.String("foo-keyname"),
					BlockDeviceMappings: []ec2types.LaunchTemplateBlockDeviceMapping{
						{
							DeviceName: aws.String("foo-device"),
							Ebs: &ec2types.LaunchTemplateEbsBlockDevice{
								Encrypted:  aws.Bool(true),
								VolumeSize: aws.Int32(16),
								VolumeType: ec2types.VolumeTypeGp2,
							},
						},
					},
					NetworkInterfaces: []ec2types.LaunchTemplateInstanceNetworkInterfaceSpecification{
						{
							DeviceIndex: aws.Int32(1),
							Groups:      []string{"foo-group"},
						},
					},
					UserData: aws.String(base64.StdEncoding.EncodeToString([]byte(testUserData))),
				},
				VersionNumber: aws.Int64(1),
			},
			wantLT: &expinfrav1.AWSLaunchTemplate{
				Name: "foo",
				AMI: infrav1.AMIReference{
					ID: aws.String("foo-image"),
				},
				IamInstanceProfile: "foo-profile",
				SSHKeyName:         aws.String("foo-keyname"),
				VersionNumber:      aws.Int64(1),
			},
			wantUserDataHash:      testUserDataHash,
			wantDataSecretKey:     nil, // respective tag is not given
			wantBootstrapDataHash: nil, // respective tag is not given
		},
		{
			name: "spot market options",
			input: ec2types.LaunchTemplateVersion{
				LaunchTemplateId:   aws.String("lt-12345"),
				LaunchTemplateName: aws.String("foo"),
				LaunchTemplateData: &ec2types.ResponseLaunchTemplateData{
					ImageId: aws.String("foo-image"),
					IamInstanceProfile: &ec2types.LaunchTemplateIamInstanceProfileSpecification{
						Arn: aws.String("instance-profile/foo-profile"),
					},
					KeyName: aws.String("foo-keyname"),
					InstanceMarketOptions: &ec2types.LaunchTemplateInstanceMarketOptions{
						MarketType: ec2types.MarketTypeSpot,
						SpotOptions: &ec2types.LaunchTemplateSpotMarketOptions{
							MaxPrice: aws.String("0.05"),
						},
					},
					UserData: aws.String(base64.StdEncoding.EncodeToString([]byte(testUserData))),
				},
				VersionNumber: aws.Int64(1),
			},
			wantLT: &expinfrav1.AWSLaunchTemplate{
				Name: "foo",
				AMI: infrav1.AMIReference{
					ID: aws.String("foo-image"),
				},
				IamInstanceProfile: "foo-profile",
				SSHKeyName:         aws.String("foo-keyname"),
				VersionNumber:      aws.Int64(1),
				SpotMarketOptions: &infrav1.SpotMarketOptions{
					MaxPrice: aws.String("0.05"),
				},
			},
			wantUserDataHash:  testUserDataHash,
			wantDataSecretKey: nil,
		},
		{
			name: "spot market options with no max price",
			input: ec2types.LaunchTemplateVersion{
				LaunchTemplateId:   aws.String("lt-12345"),
				LaunchTemplateName: aws.String("foo"),
				LaunchTemplateData: &ec2types.ResponseLaunchTemplateData{
					ImageId: aws.String("foo-image"),
					IamInstanceProfile: &ec2types.LaunchTemplateIamInstanceProfileSpecification{
						Arn: aws.String("instance-profile/foo-profile"),
					},
					KeyName: aws.String("foo-keyname"),
					InstanceMarketOptions: &ec2types.LaunchTemplateInstanceMarketOptions{
						MarketType:  ec2types.MarketTypeSpot,
						SpotOptions: &ec2types.LaunchTemplateSpotMarketOptions{},
					},
					UserData: aws.String(base64.StdEncoding.EncodeToString([]byte(testUserData))),
				},
				VersionNumber: aws.Int64(1),
			},
			wantLT: &expinfrav1.AWSLaunchTemplate{
				Name: "foo",
				AMI: infrav1.AMIReference{
					ID: aws.String("foo-image"),
				},
				IamInstanceProfile: "foo-profile",
				SSHKeyName:         aws.String("foo-keyname"),
				VersionNumber:      aws.Int64(1),
				SpotMarketOptions:  &infrav1.SpotMarketOptions{},
			},
			wantUserDataHash:  testUserDataHash,
			wantDataSecretKey: nil,
		},
		{
			name: "spot market options without SpotOptions",
			input: ec2types.LaunchTemplateVersion{
				LaunchTemplateId:   aws.String("lt-12345"),
				LaunchTemplateName: aws.String("foo"),
				LaunchTemplateData: &ec2types.ResponseLaunchTemplateData{
					ImageId: aws.String("foo-image"),
					IamInstanceProfile: &ec2types.LaunchTemplateIamInstanceProfileSpecification{
						Arn: aws.String("instance-profile/foo-profile"),
					},
					KeyName: aws.String("foo-keyname"),
					InstanceMarketOptions: &ec2types.LaunchTemplateInstanceMarketOptions{
						MarketType: ec2types.MarketTypeSpot,
					},
					UserData: aws.String(base64.StdEncoding.EncodeToString([]byte(testUserData))),
				},
				VersionNumber: aws.Int64(1),
			},
			wantLT: &expinfrav1.AWSLaunchTemplate{
				Name: "foo",
				AMI: infrav1.AMIReference{
					ID: aws.String("foo-image"),
				},
				IamInstanceProfile: "foo-profile",
				SSHKeyName:         aws.String("foo-keyname"),
				VersionNumber:      aws.Int64(1),
				SpotMarketOptions:  &infrav1.SpotMarketOptions{},
			},
			wantUserDataHash:  testUserDataHash,
			wantDataSecretKey: nil,
		},
		{
			name: "non-spot market type",
			input: ec2types.LaunchTemplateVersion{
				LaunchTemplateId:   aws.String("lt-12345"),
				LaunchTemplateName: aws.String("foo"),
				LaunchTemplateData: &ec2types.ResponseLaunchTemplateData{
					ImageId: aws.String("foo-image"),
					IamInstanceProfile: &ec2types.LaunchTemplateIamInstanceProfileSpecification{
						Arn: aws.String("instance-profile/foo-profile"),
					},
					KeyName: aws.String("foo-keyname"),
					InstanceMarketOptions: &ec2types.LaunchTemplateInstanceMarketOptions{
						MarketType: ec2types.MarketTypeCapacityBlock,
						SpotOptions: &ec2types.LaunchTemplateSpotMarketOptions{
							MaxPrice: aws.String("0.05"),
						},
					},
					UserData: aws.String(base64.StdEncoding.EncodeToString([]byte(testUserData))),
				},
				VersionNumber: aws.Int64(1),
			},
			wantLT: &expinfrav1.AWSLaunchTemplate{
				Name: "foo",
				AMI: infrav1.AMIReference{
					ID: aws.String("foo-image"),
				},
				IamInstanceProfile: "foo-profile",
				SSHKeyName:         aws.String("foo-keyname"),
				VersionNumber:      aws.Int64(1),
				SpotMarketOptions:  nil, // Should be nil since market type is not "spot"
			},
			wantUserDataHash:  testUserDataHash,
			wantDataSecretKey: nil,
		},
		{
			name: "tag of bootstrap secret",
			input: ec2types.LaunchTemplateVersion{
				LaunchTemplateId:   aws.String("lt-12345"),
				LaunchTemplateName: aws.String("foo"),
				LaunchTemplateData: &ec2types.ResponseLaunchTemplateData{
					ImageId: aws.String("foo-image"),
					IamInstanceProfile: &ec2types.LaunchTemplateIamInstanceProfileSpecification{
						Arn: aws.String("instance-profile/foo-profile"),
					},
					KeyName: aws.String("foo-keyname"),
					BlockDeviceMappings: []ec2types.LaunchTemplateBlockDeviceMapping{
						{
							DeviceName: aws.String("foo-device"),
							Ebs: &ec2types.LaunchTemplateEbsBlockDevice{
								Encrypted:  aws.Bool(true),
								VolumeSize: aws.Int32(16),
								VolumeType: ec2types.VolumeTypeGp2,
							},
						},
					},
					NetworkInterfaces: []ec2types.LaunchTemplateInstanceNetworkInterfaceSpecification{
						{
							DeviceIndex: aws.Int32(1),
							Groups:      []string{"foo-group"},
						},
					},
					TagSpecifications: []ec2types.LaunchTemplateTagSpecification{
						{
							ResourceType: ec2types.ResourceTypeInstance,
							Tags: []ec2types.Tag{
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/bootstrap-data-secret"),
									Value: aws.String("bootstrap-secret-ns/bootstrap-secret"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/bootstrap-data-hash"),
									Value: aws.String(testBootstrapDataHash),
								},
							},
						},
					},
					UserData: aws.String(base64.StdEncoding.EncodeToString([]byte(testUserData))),
				},
				VersionNumber: aws.Int64(1),
			},
			wantLT: &expinfrav1.AWSLaunchTemplate{
				Name: "foo",
				AMI: infrav1.AMIReference{
					ID: aws.String("foo-image"),
				},
				IamInstanceProfile: "foo-profile",
				SSHKeyName:         aws.String("foo-keyname"),
				VersionNumber:      aws.Int64(1),
			},
			wantUserDataHash:      testUserDataHash,
			wantDataSecretKey:     &types.NamespacedName{Namespace: "bootstrap-secret-ns", Name: "bootstrap-secret"},
			wantBootstrapDataHash: &testBootstrapDataHash,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			s := &Service{}
			gotLT, gotUserDataHash, gotDataSecretKey, gotBootstrapDataHash, err := s.SDKToLaunchTemplate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error mismatch: got %v, wantErr %v", err, tt.wantErr)
			}
			if !cmp.Equal(gotLT, tt.wantLT) {
				t.Fatalf("launchTemplate mismatch: got %v, want %v", gotLT, tt.wantLT)
			}
			if !cmp.Equal(gotUserDataHash, tt.wantUserDataHash) {
				t.Fatalf("userDataHash mismatch: got %v, want %v", gotUserDataHash, tt.wantUserDataHash)
			}
			if !cmp.Equal(gotDataSecretKey, tt.wantDataSecretKey) {
				t.Fatalf("userDataSecretKey mismatch: got %v, want %v", gotDataSecretKey, tt.wantDataSecretKey)
			}
			g.Expect(gotBootstrapDataHash).To(Equal(tt.wantBootstrapDataHash))
		})
	}
}

func TestServiceLaunchTemplateNeedsUpdate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name                  string
		incoming              *expinfrav1.AWSLaunchTemplate
		existing              *expinfrav1.AWSLaunchTemplate
		expect                func(m *mocks.MockEC2APIMockRecorder)
		want                  bool
		wantNeedsUpdateReason string
		wantErr               bool
	}{
		{
			name: "only core security groups, order shouldn't matter",
			incoming: &expinfrav1.AWSLaunchTemplate{
				AdditionalSecurityGroups: []infrav1.AWSResourceReference{},
			},
			existing: &expinfrav1.AWSLaunchTemplate{
				AdditionalSecurityGroups: []infrav1.AWSResourceReference{
					{ID: aws.String("sg-222")},
					{ID: aws.String("sg-111")},
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "the same security groups",
			incoming: &expinfrav1.AWSLaunchTemplate{
				AdditionalSecurityGroups: []infrav1.AWSResourceReference{
					{ID: aws.String("sg-999")},
				},
			},
			existing: &expinfrav1.AWSLaunchTemplate{
				AdditionalSecurityGroups: []infrav1.AWSResourceReference{
					{ID: aws.String("sg-111")},
					{ID: aws.String("sg-222")},
					{ID: aws.String("sg-999")},
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "core security group removed externally",
			incoming: &expinfrav1.AWSLaunchTemplate{
				AdditionalSecurityGroups: []infrav1.AWSResourceReference{
					{ID: aws.String("sg-999")},
				},
			},
			existing: &expinfrav1.AWSLaunchTemplate{
				AdditionalSecurityGroups: []infrav1.AWSResourceReference{
					{ID: aws.String("sg-222")},
					{ID: aws.String("sg-999")},
				},
			},
			want:                  true,
			wantNeedsUpdateReason: "AdditionalSecurityGroupsIDs",
			wantErr:               false,
		},
		{
			name: "new additional security group",
			incoming: &expinfrav1.AWSLaunchTemplate{
				AdditionalSecurityGroups: []infrav1.AWSResourceReference{
					{ID: aws.String("sg-999")},
					{ID: aws.String("sg-000")},
				},
			},
			existing: &expinfrav1.AWSLaunchTemplate{
				AdditionalSecurityGroups: []infrav1.AWSResourceReference{
					{ID: aws.String("sg-111")},
					{ID: aws.String("sg-222")},
					{ID: aws.String("sg-999")},
				},
			},
			want:                  true,
			wantNeedsUpdateReason: "AdditionalSecurityGroupsIDs",
			wantErr:               false,
		},
		{
			name: "Should return true if incoming IamInstanceProfile is not same as existing IamInstanceProfile",
			incoming: &expinfrav1.AWSLaunchTemplate{
				IamInstanceProfile: DefaultAmiNameFormat,
			},
			existing: &expinfrav1.AWSLaunchTemplate{
				AdditionalSecurityGroups: []infrav1.AWSResourceReference{
					{ID: aws.String("sg-111")},
					{ID: aws.String("sg-222")},
				},
				IamInstanceProfile: "some-other-profile",
			},
			want:                  true,
			wantNeedsUpdateReason: "IamInstanceProfile",
		},
		{
			name: "Should return true if incoming InstanceType is not same as existing InstanceType",
			incoming: &expinfrav1.AWSLaunchTemplate{
				InstanceType: "t3.micro",
			},
			existing: &expinfrav1.AWSLaunchTemplate{
				AdditionalSecurityGroups: []infrav1.AWSResourceReference{
					{ID: aws.String("sg-111")},
					{ID: aws.String("sg-222")},
				},
				InstanceType: "t3.large",
			},
			want:                  true,
			wantNeedsUpdateReason: "InstanceType",
		},
		{
			name: "new additional security group with filters",
			incoming: &expinfrav1.AWSLaunchTemplate{
				AdditionalSecurityGroups: []infrav1.AWSResourceReference{
					{Filters: []infrav1.Filter{{Name: "sg-1", Values: []string{"test-1"}}}},
				},
			},
			existing: &expinfrav1.AWSLaunchTemplate{
				AdditionalSecurityGroups: []infrav1.AWSResourceReference{
					{Filters: []infrav1.Filter{{Name: "sg-2", Values: []string{"test-2"}}}},
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSecurityGroups(context.TODO(), gomock.Eq(&ec2.DescribeSecurityGroupsInput{Filters: []ec2types.Filter{{Name: aws.String("sg-1"), Values: []string{"test-1"}}}})).
					Return(&ec2.DescribeSecurityGroupsOutput{SecurityGroups: []ec2types.SecurityGroup{{GroupId: aws.String("sg-1")}}}, nil)
				m.DescribeSecurityGroups(context.TODO(), gomock.Eq(&ec2.DescribeSecurityGroupsInput{Filters: []ec2types.Filter{{Name: aws.String("sg-2"), Values: []string{"test-2"}}}})).
					Return(&ec2.DescribeSecurityGroupsOutput{SecurityGroups: []ec2types.SecurityGroup{{GroupId: aws.String("sg-2")}}}, nil)
			},
			want:                  true,
			wantNeedsUpdateReason: "AdditionalSecurityGroupsIDs",
			wantErr:               false,
		},
		{
			name: "new launch template instance metadata options, requiring IMDSv2",
			incoming: &expinfrav1.AWSLaunchTemplate{
				InstanceMetadataOptions: &infrav1.InstanceMetadataOptions{
					HTTPPutResponseHopLimit: 1,
					HTTPTokens:              infrav1.HTTPTokensStateRequired,
				},
			},
			existing: &expinfrav1.AWSLaunchTemplate{
				AdditionalSecurityGroups: []infrav1.AWSResourceReference{
					{ID: aws.String("sg-111")},
					{ID: aws.String("sg-222")},
				},
			},
			want:                  true,
			wantNeedsUpdateReason: "InstanceMetadataOptions",
			wantErr:               false,
		},
		{
			name:     "new launch template instance metadata options, removing IMDSv2 requirement",
			incoming: &expinfrav1.AWSLaunchTemplate{},
			existing: &expinfrav1.AWSLaunchTemplate{
				InstanceMetadataOptions: &infrav1.InstanceMetadataOptions{
					HTTPPutResponseHopLimit: 1,
					HTTPTokens:              infrav1.HTTPTokensStateRequired,
				},
			},
			want:                  true,
			wantNeedsUpdateReason: "InstanceMetadataOptions",
			wantErr:               false,
		},
		{
			name: "Should return true if incoming SpotMarketOptions is different from existing SpotMarketOptions",
			incoming: &expinfrav1.AWSLaunchTemplate{
				SpotMarketOptions: &infrav1.SpotMarketOptions{
					MaxPrice: aws.String("0.10"),
				},
			},
			existing: &expinfrav1.AWSLaunchTemplate{
				AdditionalSecurityGroups: []infrav1.AWSResourceReference{
					{ID: aws.String("sg-111")},
					{ID: aws.String("sg-222")},
				},
				SpotMarketOptions: &infrav1.SpotMarketOptions{
					MaxPrice: aws.String("0.05"),
				},
			},
			want:                  true,
			wantNeedsUpdateReason: "SpotMarketOptions",
			wantErr:               false,
		},
		{
			name: "Should return true if incoming adds SpotMarketOptions and existing has none",
			incoming: &expinfrav1.AWSLaunchTemplate{
				SpotMarketOptions: &infrav1.SpotMarketOptions{
					MaxPrice: aws.String("0.10"),
				},
			},
			existing: &expinfrav1.AWSLaunchTemplate{
				AdditionalSecurityGroups: []infrav1.AWSResourceReference{
					{ID: aws.String("sg-111")},
					{ID: aws.String("sg-222")},
				},
				SpotMarketOptions: nil,
			},
			want:                  true,
			wantNeedsUpdateReason: "SpotMarketOptions",
			wantErr:               false,
		},
		{
			name: "Should return true if incoming removes SpotMarketOptions and existing has some",
			incoming: &expinfrav1.AWSLaunchTemplate{
				SpotMarketOptions: nil,
			},
			existing: &expinfrav1.AWSLaunchTemplate{
				AdditionalSecurityGroups: []infrav1.AWSResourceReference{
					{ID: aws.String("sg-111")},
					{ID: aws.String("sg-222")},
				},
				SpotMarketOptions: &infrav1.SpotMarketOptions{
					MaxPrice: aws.String("0.05"),
				},
			},
			want:                  true,
			wantNeedsUpdateReason: "SpotMarketOptions",
			wantErr:               false,
		},
		{
			name: "Should return true if SSH key names are different",
			incoming: &expinfrav1.AWSLaunchTemplate{
				SSHKeyName: aws.String("new-key"),
			},
			existing: &expinfrav1.AWSLaunchTemplate{
				AdditionalSecurityGroups: []infrav1.AWSResourceReference{
					{ID: aws.String("sg-111")},
					{ID: aws.String("sg-222")},
				},
				SSHKeyName: aws.String("old-key"),
			},
			want:                  true,
			wantNeedsUpdateReason: "SSHKeyName",
			wantErr:               false,
		},
		{
			name: "Should return true if one has SSH key name and other doesn't",
			incoming: &expinfrav1.AWSLaunchTemplate{
				SSHKeyName: aws.String("new-key"),
			},
			existing: &expinfrav1.AWSLaunchTemplate{
				AdditionalSecurityGroups: []infrav1.AWSResourceReference{
					{ID: aws.String("sg-111")},
					{ID: aws.String("sg-222")},
				},
				SSHKeyName: nil,
			},
			want:                  true,
			wantNeedsUpdateReason: "SSHKeyName",
			wantErr:               false,
		},
		{
			name: "Should return false if no SSH key is set in the spec and AWS returns no key pair as well",
			incoming: &expinfrav1.AWSLaunchTemplate{
				SSHKeyName: aws.String(""), // explicit empty string
			},
			existing: &expinfrav1.AWSLaunchTemplate{
				AdditionalSecurityGroups: []infrav1.AWSResourceReference{
					{ID: aws.String("sg-111")},
					{ID: aws.String("sg-222")},
				},
				SSHKeyName: nil,
			},
			want:                  false,
			wantNeedsUpdateReason: "",
			wantErr:               false,
		},
		{
			name: "Should return true if incoming PrivateDNSName is different from existing PrivateDNSName",
			incoming: &expinfrav1.AWSLaunchTemplate{
				PrivateDNSName: &infrav1.PrivateDNSName{
					EnableResourceNameDNSARecord:    aws.Bool(true),
					EnableResourceNameDNSAAAARecord: aws.Bool(false),
					HostnameType:                    aws.String("resource-name"),
				},
			},
			existing: &expinfrav1.AWSLaunchTemplate{
				AdditionalSecurityGroups: []infrav1.AWSResourceReference{
					{ID: aws.String("sg-111")},
					{ID: aws.String("sg-222")},
				},
				PrivateDNSName: &infrav1.PrivateDNSName{
					EnableResourceNameDNSARecord:    aws.Bool(false),
					EnableResourceNameDNSAAAARecord: aws.Bool(false),
					HostnameType:                    aws.String("ip-name"),
				},
			},
			want:                  true,
			wantNeedsUpdateReason: "PrivateDNSName",
			wantErr:               false,
		},
		{
			name: "Should return true if incoming adds PrivateDNSName and existing has none",
			incoming: &expinfrav1.AWSLaunchTemplate{
				PrivateDNSName: &infrav1.PrivateDNSName{
					EnableResourceNameDNSARecord:    aws.Bool(true),
					EnableResourceNameDNSAAAARecord: aws.Bool(false),
					HostnameType:                    aws.String("resource-name"),
				},
			},
			existing: &expinfrav1.AWSLaunchTemplate{
				AdditionalSecurityGroups: []infrav1.AWSResourceReference{
					{ID: aws.String("sg-111")},
					{ID: aws.String("sg-222")},
				},
				PrivateDNSName: nil,
			},
			want:                  true,
			wantNeedsUpdateReason: "PrivateDNSName",
			wantErr:               false,
		},
		{
			name: "Should return true if incoming removes PrivateDNSName and existing has some",
			incoming: &expinfrav1.AWSLaunchTemplate{
				PrivateDNSName: nil,
			},
			existing: &expinfrav1.AWSLaunchTemplate{
				AdditionalSecurityGroups: []infrav1.AWSResourceReference{
					{ID: aws.String("sg-111")},
					{ID: aws.String("sg-222")},
				},
				PrivateDNSName: &infrav1.PrivateDNSName{
					EnableResourceNameDNSARecord:    aws.Bool(true),
					EnableResourceNameDNSAAAARecord: aws.Bool(false),
					HostnameType:                    aws.String("resource-name"),
				},
			},
			want:                  true,
			wantNeedsUpdateReason: "PrivateDNSName",
			wantErr:               false,
		},
		{
			name: "Should return true if capacity reservation IDs are different",
			incoming: &expinfrav1.AWSLaunchTemplate{
				CapacityReservationID: aws.String("new-reservation"),
			},
			existing: &expinfrav1.AWSLaunchTemplate{
				AdditionalSecurityGroups: []infrav1.AWSResourceReference{
					{ID: aws.String("sg-111")},
					{ID: aws.String("sg-222")},
				},
				CapacityReservationID: aws.String("old-reservation"),
			},
			want:                  true,
			wantNeedsUpdateReason: "CapacityReservationID",
			wantErr:               false,
		},
		{
			name: "Should return true if one has capacity reservation ID and other doesn't",
			incoming: &expinfrav1.AWSLaunchTemplate{
				CapacityReservationID: aws.String("new-reservation"),
			},
			existing: &expinfrav1.AWSLaunchTemplate{
				AdditionalSecurityGroups: []infrav1.AWSResourceReference{
					{ID: aws.String("sg-111")},
					{ID: aws.String("sg-222")},
				},
				CapacityReservationID: nil,
			},
			want:                  true,
			wantNeedsUpdateReason: "CapacityReservationID",
			wantErr:               false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			ac := &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Status: infrav1.AWSClusterStatus{
					Network: infrav1.NetworkStatus{
						SecurityGroups: map[infrav1.SecurityGroupRole]infrav1.SecurityGroup{
							infrav1.SecurityGroupNode: {
								ID: "sg-111",
							},
							infrav1.SecurityGroupLB: {
								ID: "sg-222",
							},
						},
					},
				},
			}
			s := &Service{
				scope: &scope.ClusterScope{
					AWSCluster: ac,
				},
			}
			machinePoolScope := &scope.MachinePoolScope{
				InfraCluster: &scope.ClusterScope{
					AWSCluster: ac,
				},
			}
			mockEC2Client := mocks.NewMockEC2API(mockCtrl)
			s.EC2Client = mockEC2Client

			if tt.expect != nil {
				tt.expect(mockEC2Client.EXPECT())
			}

			got, gotNeedsUpdateReason, err := s.LaunchTemplateNeedsUpdate(machinePoolScope, tt.incoming, tt.existing)
			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(gotNeedsUpdateReason).Should(Equal(tt.wantNeedsUpdateReason))
			g.Expect(got).Should(Equal(tt.want))
		})
	}
}

func TestGetLaunchTemplateID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name               string
		launchTemplateName string
		expect             func(m *mocks.MockEC2APIMockRecorder)
		check              func(g *WithT, launchTemplateID string, err error)
	}{
		{
			name:   "Should return with no error if empty launch template name passed",
			expect: func(m *mocks.MockEC2APIMockRecorder) {},
			check: func(g *WithT, launchTemplateID string, err error) {
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(launchTemplateID).Should(BeEmpty())
			},
		},
		{
			name:               "Should not return error if launch template does not exist",
			launchTemplateName: "foo",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeLaunchTemplateVersions(context.TODO(), gomock.Eq(&ec2.DescribeLaunchTemplateVersionsInput{
					LaunchTemplateName: aws.String("foo"),
					Versions:           []string{"$Latest"},
				})).Return(nil, &smithy.GenericAPIError{
					Code:    awserrors.LaunchTemplateNameNotFound,
					Message: "The specified launch template, with template name foo, does not exist.",
				})
			},
			check: func(g *WithT, launchTemplateID string, err error) {
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(launchTemplateID).Should(BeEmpty())
			},
		},
		{
			name:               "Should return with error if AWS failed to fetch launch template",
			launchTemplateName: "foo",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeLaunchTemplateVersions(context.TODO(), gomock.Eq(&ec2.DescribeLaunchTemplateVersionsInput{
					LaunchTemplateName: aws.String("foo"),
					Versions:           []string{"$Latest"},
				})).Return(nil, awserrors.NewFailedDependency("Dependency issue from AWS"))
			},
			check: func(g *WithT, launchTemplateID string, err error) {
				g.Expect(err).To(HaveOccurred())
				g.Expect(launchTemplateID).Should(BeEmpty())
			},
		},
		{
			name:               "Should not return error if AWS returns no launch template versions info in output",
			launchTemplateName: "foo",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeLaunchTemplateVersions(context.TODO(), gomock.Eq(&ec2.DescribeLaunchTemplateVersionsInput{
					LaunchTemplateName: aws.String("foo"),
					Versions:           []string{"$Latest"},
				})).Return(nil, nil)
			},
			check: func(g *WithT, launchTemplateID string, err error) {
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(launchTemplateID).Should(BeEmpty())
			},
		},
		{
			name:               "Should successfully return launch template ID for given name if exists",
			launchTemplateName: "foo",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeLaunchTemplateVersions(context.TODO(), gomock.Eq(&ec2.DescribeLaunchTemplateVersionsInput{
					LaunchTemplateName: aws.String("foo"),
					Versions:           []string{"$Latest"},
				})).Return(&ec2.DescribeLaunchTemplateVersionsOutput{
					LaunchTemplateVersions: []ec2types.LaunchTemplateVersion{
						{
							LaunchTemplateId:   aws.String("lt-12345"),
							LaunchTemplateName: aws.String("foo"),
							LaunchTemplateData: &ec2types.ResponseLaunchTemplateData{
								ImageId: aws.String("foo-image"),
								IamInstanceProfile: &ec2types.LaunchTemplateIamInstanceProfileSpecification{
									Arn: aws.String("instance-profile/foo-profile"),
								},
								KeyName: aws.String("foo-keyname"),
								BlockDeviceMappings: []ec2types.LaunchTemplateBlockDeviceMapping{
									{
										DeviceName: aws.String("foo-device"),
										Ebs: &ec2types.LaunchTemplateEbsBlockDevice{
											Encrypted:  aws.Bool(true),
											VolumeSize: aws.Int32(16),
											VolumeType: ec2types.VolumeTypeGp2,
										},
									},
								},
								NetworkInterfaces: []ec2types.LaunchTemplateInstanceNetworkInterfaceSpecification{
									{
										DeviceIndex: aws.Int32(1),
										Groups:      []string{"foo-group"},
									},
								},
								UserData: aws.String(base64.StdEncoding.EncodeToString([]byte(testUserData))),
							},
							VersionNumber: aws.Int64(1),
						},
					},
				}, nil)
			},
			check: func(g *WithT, launchTemplateID string, err error) {
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(launchTemplateID).Should(Equal("lt-12345"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			scheme, err := setupScheme()
			g.Expect(err).NotTo(HaveOccurred())
			client := fake.NewClientBuilder().WithScheme(scheme).Build()

			cs, err := setupClusterScope(client)
			g.Expect(err).NotTo(HaveOccurred())
			mockEC2Client := mocks.NewMockEC2API(mockCtrl)

			s := NewService(cs)
			s.EC2Client = mockEC2Client

			if tc.expect != nil {
				tc.expect(mockEC2Client.EXPECT())
			}
			launchTemplate, err := s.GetLaunchTemplateID(tc.launchTemplateName)
			tc.check(g, launchTemplate, err)
		})
	}
}

func TestDeleteLaunchTemplate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name      string
		versionID string
		expect    func(m *mocks.MockEC2APIMockRecorder)
		wantErr   bool
	}{
		{
			name:      "Should not return error if successfully deletes given launch template ID",
			versionID: "1",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DeleteLaunchTemplate(context.TODO(), gomock.Eq(&ec2.DeleteLaunchTemplateInput{
					LaunchTemplateId: aws.String("1"),
				})).Return(&ec2.DeleteLaunchTemplateOutput{}, nil)
			},
		},
		{
			name:      "Should return error if failed to delete given launch template ID",
			versionID: "1",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DeleteLaunchTemplate(context.TODO(), gomock.Eq(&ec2.DeleteLaunchTemplateInput{
					LaunchTemplateId: aws.String("1"),
				})).Return(nil, awserrors.NewFailedDependency("dependency failure"))
			},
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			scheme, err := setupScheme()
			g.Expect(err).NotTo(HaveOccurred())
			client := fake.NewClientBuilder().WithScheme(scheme).Build()

			cs, err := setupClusterScope(client)
			g.Expect(err).NotTo(HaveOccurred())
			mockEC2Client := mocks.NewMockEC2API(mockCtrl)

			s := NewService(cs)
			s.EC2Client = mockEC2Client
			tc.expect(mockEC2Client.EXPECT())

			err = s.DeleteLaunchTemplate(tc.versionID)
			if tc.wantErr {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).NotTo(HaveOccurred())
		})
	}
}

func TestCreateLaunchTemplate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	formatTagsInput := func(arg *ec2.CreateLaunchTemplateInput) {
		sortTags(arg.TagSpecifications[0].Tags)

		for index := range arg.LaunchTemplateData.TagSpecifications {
			sortTags(arg.LaunchTemplateData.TagSpecifications[index].Tags)
		}
	}

	userDataSecretKey := types.NamespacedName{
		Namespace: "bootstrap-secret-ns",
		Name:      "bootstrap-secret",
	}
	userData := []byte{1, 0, 0}
	testCases := []struct {
		name                 string
		awsResourceReference []infrav1.AWSResourceReference
		expect               func(g *WithT, m *mocks.MockEC2APIMockRecorder)
		check                func(g *WithT, s string, e error)
	}{
		{
			name:                 "Should not return error if successfully created launch template id",
			awsResourceReference: []infrav1.AWSResourceReference{{ID: aws.String("1")}},
			expect: func(g *WithT, m *mocks.MockEC2APIMockRecorder) {
				sgMap := make(map[infrav1.SecurityGroupRole]infrav1.SecurityGroup)
				sgMap[infrav1.SecurityGroupNode] = infrav1.SecurityGroup{ID: "1"}
				sgMap[infrav1.SecurityGroupLB] = infrav1.SecurityGroup{ID: "2"}

				expectedInput := &ec2.CreateLaunchTemplateInput{
					LaunchTemplateData: &ec2types.RequestLaunchTemplateData{
						InstanceType: ec2types.InstanceTypeT3Large,
						IamInstanceProfile: &ec2types.LaunchTemplateIamInstanceProfileSpecificationRequest{
							Name: aws.String("instance-profile"),
						},
						KeyName:          aws.String("default"),
						UserData:         ptr.To[string](base64.StdEncoding.EncodeToString(userData)),
						SecurityGroupIds: []string{"nodeSG", "lbSG", "1"},
						ImageId:          aws.String("imageID"),
						InstanceMarketOptions: &ec2types.LaunchTemplateInstanceMarketOptionsRequest{
							MarketType: ec2types.MarketTypeSpot,
							SpotOptions: &ec2types.LaunchTemplateSpotMarketOptionsRequest{
								MaxPrice: aws.String("0.9"),
							},
						},
						TagSpecifications: []ec2types.LaunchTemplateTagSpecificationRequest{
							{
								ResourceType: ec2types.ResourceTypeInstance,
								Tags:         defaultEC2AndDataTags("aws-mp-name", "cluster-name", userDataSecretKey, testBootstrapDataHash),
							},
							{
								ResourceType: ec2types.ResourceTypeVolume,
								Tags:         defaultEC2Tags("aws-mp-name", "cluster-name"),
							},
						},
					},
					LaunchTemplateName: aws.String("aws-mp-name"),
					TagSpecifications: []ec2types.TagSpecification{
						{
							ResourceType: ec2types.ResourceTypeLaunchTemplate,
							Tags:         defaultEC2Tags("aws-mp-name", "cluster-name"),
						},
					},
				}
				m.CreateLaunchTemplate(context.TODO(), gomock.AssignableToTypeOf(expectedInput)).Return(&ec2.CreateLaunchTemplateOutput{
					LaunchTemplate: &ec2types.LaunchTemplate{
						LaunchTemplateId: aws.String("launch-template-id"),
					},
				}, nil).Do(func(ctx context.Context, arg *ec2.CreateLaunchTemplateInput, requestOptions ...ec2.Options) {
					// formatting added to match arrays during cmp.Equal
					formatTagsInput(arg)
					if !cmp.Equal(expectedInput, arg, cmpopts.IgnoreUnexported(
						ec2types.LaunchTemplateIamInstanceProfileSpecificationRequest{},
						ec2types.LaunchTemplateSpotMarketOptionsRequest{},
						ec2types.LaunchTemplateInstanceMarketOptionsRequest{},
						ec2types.Tag{},
						ec2types.LaunchTemplateTagSpecificationRequest{},
						ec2types.RequestLaunchTemplateData{},
						ec2types.TagSpecification{},
						ec2.CreateLaunchTemplateInput{},
					)) {
						t.Fatalf("mismatch in input expected: %+v, got: %+v, diff: %s", expectedInput, arg, cmp.Diff(expectedInput, arg))
					}
				})
			},
			check: func(g *WithT, id string, err error) {
				g.Expect(id).Should(Equal("launch-template-id"))
				g.Expect(err).NotTo(HaveOccurred())
			},
		},
		{
			name:                 "Should successfully create launch template id with AdditionalSecurityGroups Filter",
			awsResourceReference: []infrav1.AWSResourceReference{{Filters: []infrav1.Filter{{Name: "sg-1", Values: []string{"test"}}}}},
			expect: func(g *WithT, m *mocks.MockEC2APIMockRecorder) {
				sgMap := make(map[infrav1.SecurityGroupRole]infrav1.SecurityGroup)
				sgMap[infrav1.SecurityGroupNode] = infrav1.SecurityGroup{ID: "1"}
				sgMap[infrav1.SecurityGroupLB] = infrav1.SecurityGroup{ID: "2"}

				expectedInput := &ec2.CreateLaunchTemplateInput{
					LaunchTemplateData: &ec2types.RequestLaunchTemplateData{
						InstanceType: ec2types.InstanceTypeT3Large,
						IamInstanceProfile: &ec2types.LaunchTemplateIamInstanceProfileSpecificationRequest{
							Name: aws.String("instance-profile"),
						},
						KeyName:          aws.String("default"),
						UserData:         ptr.To[string](base64.StdEncoding.EncodeToString(userData)),
						SecurityGroupIds: []string{"nodeSG", "lbSG", "sg-1"},
						ImageId:          aws.String("imageID"),
						InstanceMarketOptions: &ec2types.LaunchTemplateInstanceMarketOptionsRequest{
							MarketType: ec2types.MarketTypeSpot,
							SpotOptions: &ec2types.LaunchTemplateSpotMarketOptionsRequest{
								MaxPrice: aws.String("0.9"),
							},
						},
						TagSpecifications: []ec2types.LaunchTemplateTagSpecificationRequest{
							{
								ResourceType: ec2types.ResourceTypeInstance,
								Tags:         defaultEC2AndDataTags("aws-mp-name", "cluster-name", userDataSecretKey, testBootstrapDataHash),
							},
							{
								ResourceType: ec2types.ResourceTypeVolume,
								Tags:         defaultEC2Tags("aws-mp-name", "cluster-name"),
							},
						},
					},
					LaunchTemplateName: aws.String("aws-mp-name"),
					TagSpecifications: []ec2types.TagSpecification{
						{
							ResourceType: ec2types.ResourceTypeLaunchTemplate,
							Tags:         defaultEC2Tags("aws-mp-name", "cluster-name"),
						},
					},
				}
				m.CreateLaunchTemplate(context.TODO(), gomock.AssignableToTypeOf(expectedInput)).Return(&ec2.CreateLaunchTemplateOutput{
					LaunchTemplate: &ec2types.LaunchTemplate{
						LaunchTemplateId: aws.String("launch-template-id"),
					},
				}, nil).Do(func(ctx context.Context, arg *ec2.CreateLaunchTemplateInput, requestOptions ...ec2.Options) {
					// formatting added to match arrays during reflect.DeepEqual
					formatTagsInput(arg)
					if !cmp.Equal(expectedInput, arg, cmpopts.IgnoreUnexported(
						ec2types.LaunchTemplateIamInstanceProfileSpecificationRequest{},
						ec2types.LaunchTemplateSpotMarketOptionsRequest{},
						ec2types.LaunchTemplateInstanceMarketOptionsRequest{},
						ec2types.Tag{},
						ec2types.LaunchTemplateTagSpecificationRequest{},
						ec2types.RequestLaunchTemplateData{},
						ec2types.TagSpecification{},
						ec2.CreateLaunchTemplateInput{},
					)) {
						t.Fatalf("mismatch in input expected: %+v, got: %+v", expectedInput, arg)
					}
				})
				m.DescribeSecurityGroups(context.TODO(), gomock.Eq(&ec2.DescribeSecurityGroupsInput{Filters: []ec2types.Filter{{Name: aws.String("sg-1"), Values: []string{"test"}}}})).
					Return(&ec2.DescribeSecurityGroupsOutput{SecurityGroups: []ec2types.SecurityGroup{{GroupId: aws.String("sg-1")}}}, nil)
			},
			check: func(g *WithT, id string, err error) {
				g.Expect(id).Should(Equal("launch-template-id"))
				g.Expect(err).NotTo(HaveOccurred())
			},
		},
		{
			name:                 "Should return with error if failed to create launch template id",
			awsResourceReference: []infrav1.AWSResourceReference{{ID: aws.String("1")}},
			expect: func(g *WithT, m *mocks.MockEC2APIMockRecorder) {
				sgMap := make(map[infrav1.SecurityGroupRole]infrav1.SecurityGroup)
				sgMap[infrav1.SecurityGroupNode] = infrav1.SecurityGroup{ID: "1"}
				sgMap[infrav1.SecurityGroupLB] = infrav1.SecurityGroup{ID: "2"}

				expectedInput := &ec2.CreateLaunchTemplateInput{
					LaunchTemplateData: &ec2types.RequestLaunchTemplateData{
						InstanceType: ec2types.InstanceTypeT3Large,
						IamInstanceProfile: &ec2types.LaunchTemplateIamInstanceProfileSpecificationRequest{
							Name: aws.String("instance-profile"),
						},
						KeyName:          aws.String("default"),
						UserData:         ptr.To[string](base64.StdEncoding.EncodeToString(userData)),
						SecurityGroupIds: []string{"nodeSG", "lbSG", "1"},
						ImageId:          aws.String("imageID"),
						InstanceMarketOptions: &ec2types.LaunchTemplateInstanceMarketOptionsRequest{
							MarketType: ec2types.MarketTypeSpot,
							SpotOptions: &ec2types.LaunchTemplateSpotMarketOptionsRequest{
								MaxPrice: aws.String("0.9"),
							},
						},
						TagSpecifications: []ec2types.LaunchTemplateTagSpecificationRequest{
							{
								ResourceType: ec2types.ResourceTypeInstance,
								Tags:         defaultEC2AndDataTags("aws-mp-name", "cluster-name", userDataSecretKey, testBootstrapDataHash),
							},
							{
								ResourceType: ec2types.ResourceTypeVolume,
								Tags:         defaultEC2Tags("aws-mp-name", "cluster-name"),
							},
						},
					},
					LaunchTemplateName: aws.String("aws-mp-name"),
					TagSpecifications: []ec2types.TagSpecification{
						{
							ResourceType: ec2types.ResourceTypeLaunchTemplate,
							Tags:         defaultEC2Tags("aws-mp-name", "cluster-name"),
						},
					},
				}
				m.CreateLaunchTemplate(context.TODO(), gomock.AssignableToTypeOf(expectedInput)).Return(nil,
					awserrors.NewFailedDependency("dependency failure")).Do(func(ctx context.Context, arg *ec2.CreateLaunchTemplateInput, requestOptions ...ec2.Options) {
					// formatting added to match arrays during cmp.Equal
					formatTagsInput(arg)
					if !cmp.Equal(expectedInput, arg, cmpopts.IgnoreUnexported(
						ec2types.LaunchTemplateIamInstanceProfileSpecificationRequest{},
						ec2types.LaunchTemplateSpotMarketOptionsRequest{},
						ec2types.LaunchTemplateInstanceMarketOptionsRequest{},
						ec2types.Tag{},
						ec2types.LaunchTemplateTagSpecificationRequest{},
						ec2types.RequestLaunchTemplateData{},
						ec2types.TagSpecification{},
						ec2.CreateLaunchTemplateInput{},
					)) {
						t.Fatalf("mismatch in input expected: %+v, got: %+v", expectedInput, arg)
					}
				})
			},
			check: func(g *WithT, id string, err error) {
				g.Expect(id).Should(BeEmpty())
				g.Expect(err).To(HaveOccurred())
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			scheme, err := setupScheme()
			g.Expect(err).NotTo(HaveOccurred())
			client := fake.NewClientBuilder().WithScheme(scheme).Build()

			cs, err := setupClusterScope(client)
			g.Expect(err).NotTo(HaveOccurred())
			mockEC2Client := mocks.NewMockEC2API(mockCtrl)

			ms, err := setupMachinePoolScope(client, cs)
			g.Expect(err).NotTo(HaveOccurred())

			ms.AWSMachinePool.Spec.AWSLaunchTemplate.AdditionalSecurityGroups = tc.awsResourceReference

			s := NewService(cs)
			s.EC2Client = mockEC2Client

			if tc.expect != nil {
				tc.expect(g, mockEC2Client.EXPECT())
			}

			launchTemplate, err := s.CreateLaunchTemplate(ms, aws.String("imageID"), userDataSecretKey, userData, testBootstrapDataHash)
			tc.check(g, launchTemplate, err)
		})
	}
}

func TestLaunchTemplateDataCreation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	t.Run("Should return error if failed to create launch template data", func(t *testing.T) {
		g := NewWithT(t)
		scheme, err := setupScheme()
		g.Expect(err).NotTo(HaveOccurred())
		client := fake.NewClientBuilder().WithScheme(scheme).Build()

		cs, err := setupClusterScope(client)
		g.Expect(err).NotTo(HaveOccurred())
		cs.AWSCluster.Status.Network.SecurityGroups[infrav1.SecurityGroupBastion] = infrav1.SecurityGroup{ID: "1"}

		ms, err := setupMachinePoolScope(client, cs)
		g.Expect(err).NotTo(HaveOccurred())

		s := NewService(cs)

		userDataSecretKey := types.NamespacedName{
			Namespace: "bootstrap-secret-ns",
			Name:      "bootstrap-secret",
		}
		launchTemplate, err := s.CreateLaunchTemplate(ms, aws.String("imageID"), userDataSecretKey, nil, "")
		g.Expect(err).To(HaveOccurred())
		g.Expect(launchTemplate).Should(BeEmpty())
	})
}

var LaunchTemplateVersionIgnoreUnexported = cmpopts.IgnoreUnexported(
	ec2types.CapacityReservationTarget{},
	ec2types.LaunchTemplateCapacityReservationSpecificationRequest{},
	ec2types.LaunchTemplateIamInstanceProfileSpecificationRequest{},
	ec2types.LaunchTemplateSpotMarketOptionsRequest{},
	ec2types.LaunchTemplateInstanceMarketOptionsRequest{},
	ec2types.Tag{},
	ec2types.LaunchTemplateTagSpecificationRequest{},
	ec2types.RequestLaunchTemplateData{},
	ec2.CreateLaunchTemplateVersionInput{},
)

func TestCreateLaunchTemplateVersion(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	formatTagsInput := func(arg *ec2.CreateLaunchTemplateVersionInput) {
		for index := range arg.LaunchTemplateData.TagSpecifications {
			sortTags(arg.LaunchTemplateData.TagSpecifications[index].Tags)
		}
	}
	userDataSecretKey := types.NamespacedName{
		Namespace: "bootstrap-secret-ns",
		Name:      "bootstrap-secret",
	}
	userData := []byte{1, 0, 0}
	testCases := []struct {
		name                 string
		imageID              *string
		awsResourceReference []infrav1.AWSResourceReference
		expect               func(m *mocks.MockEC2APIMockRecorder)
		wantErr              bool
		mpScopeUpdater       func(*scope.MachinePoolScope)
		marketType           ec2types.MarketType
	}{
		{
			name:                 "Should successfully creates launch template version",
			awsResourceReference: []infrav1.AWSResourceReference{{ID: aws.String("1")}},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				sgMap := make(map[infrav1.SecurityGroupRole]infrav1.SecurityGroup)
				sgMap[infrav1.SecurityGroupNode] = infrav1.SecurityGroup{ID: "1"}
				sgMap[infrav1.SecurityGroupLB] = infrav1.SecurityGroup{ID: "2"}

				expectedInput := &ec2.CreateLaunchTemplateVersionInput{
					LaunchTemplateData: &ec2types.RequestLaunchTemplateData{
						InstanceType: ec2types.InstanceTypeT3Large,
						IamInstanceProfile: &ec2types.LaunchTemplateIamInstanceProfileSpecificationRequest{
							Name: aws.String("instance-profile"),
						},
						KeyName:          aws.String("default"),
						UserData:         ptr.To[string](base64.StdEncoding.EncodeToString(userData)),
						SecurityGroupIds: []string{"nodeSG", "lbSG", "1"},
						ImageId:          aws.String("imageID"),
						InstanceMarketOptions: &ec2types.LaunchTemplateInstanceMarketOptionsRequest{
							MarketType: ec2types.MarketTypeSpot,
							SpotOptions: &ec2types.LaunchTemplateSpotMarketOptionsRequest{
								MaxPrice: aws.String("0.9"),
							},
						},
						TagSpecifications: []ec2types.LaunchTemplateTagSpecificationRequest{
							{
								ResourceType: ec2types.ResourceTypeInstance,
								Tags:         defaultEC2AndDataTags("aws-mp-name", "cluster-name", userDataSecretKey, testBootstrapDataHash),
							},
							{
								ResourceType: ec2types.ResourceTypeVolume,
								Tags:         defaultEC2Tags("aws-mp-name", "cluster-name"),
							},
						},
					},
					LaunchTemplateId: aws.String("launch-template-id"),
				}
				m.CreateLaunchTemplateVersion(context.TODO(), gomock.AssignableToTypeOf(expectedInput)).Return(&ec2.CreateLaunchTemplateVersionOutput{
					LaunchTemplateVersion: &ec2types.LaunchTemplateVersion{
						LaunchTemplateId: aws.String("launch-template-id"),
					},
				}, nil).Do(
					func(ctx context.Context, arg *ec2.CreateLaunchTemplateVersionInput, requestOptions ...ec2.Options) {
						// formatting added to match tags slice during cmp.Equal()
						formatTagsInput(arg)
						if !cmp.Equal(expectedInput, arg, LaunchTemplateVersionIgnoreUnexported) {
							t.Fatalf("mismatch in input expected: %+v, but got %+v, diff: %s", expectedInput, arg, cmp.Diff(expectedInput, arg, LaunchTemplateVersionIgnoreUnexported))
						}
					})
			},
		},
		{
			name:                 "Should successfully create launch template version with capacity-block",
			awsResourceReference: []infrav1.AWSResourceReference{{ID: aws.String("1")}},
			mpScopeUpdater: func(mps *scope.MachinePoolScope) {
				spec := mps.AWSMachinePool.Spec
				spec.AWSLaunchTemplate.CapacityReservationID = aws.String("cr-12345678901234567")
				spec.AWSLaunchTemplate.MarketType = infrav1.MarketTypeCapacityBlock
				spec.AWSLaunchTemplate.SpotMarketOptions = nil
				mps.AWSMachinePool.Spec = spec
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				sgMap := make(map[infrav1.SecurityGroupRole]infrav1.SecurityGroup)
				sgMap[infrav1.SecurityGroupNode] = infrav1.SecurityGroup{ID: "1"}
				sgMap[infrav1.SecurityGroupLB] = infrav1.SecurityGroup{ID: "2"}

				expectedInput := &ec2.CreateLaunchTemplateVersionInput{
					LaunchTemplateData: &ec2types.RequestLaunchTemplateData{
						InstanceType: ec2types.InstanceTypeT3Large,
						IamInstanceProfile: &ec2types.LaunchTemplateIamInstanceProfileSpecificationRequest{
							Name: aws.String("instance-profile"),
						},
						KeyName:          aws.String("default"),
						UserData:         ptr.To[string](base64.StdEncoding.EncodeToString(userData)),
						SecurityGroupIds: []string{"nodeSG", "lbSG", "1"},
						ImageId:          aws.String("imageID"),
						InstanceMarketOptions: &ec2types.LaunchTemplateInstanceMarketOptionsRequest{
							MarketType: ec2types.MarketTypeCapacityBlock,
						},
						CapacityReservationSpecification: &ec2types.LaunchTemplateCapacityReservationSpecificationRequest{
							CapacityReservationTarget: &ec2types.CapacityReservationTarget{
								CapacityReservationId: aws.String("cr-12345678901234567"),
							},
						},
						TagSpecifications: []ec2types.LaunchTemplateTagSpecificationRequest{
							{
								ResourceType: ec2types.ResourceTypeInstance,
								Tags:         defaultEC2AndDataTags("aws-mp-name", "cluster-name", userDataSecretKey, testBootstrapDataHash),
							},
							{
								ResourceType: ec2types.ResourceTypeVolume,
								Tags:         defaultEC2Tags("aws-mp-name", "cluster-name"),
							},
						},
					},
					LaunchTemplateId: aws.String("launch-template-id"),
				}
				m.CreateLaunchTemplateVersion(context.TODO(), gomock.AssignableToTypeOf(expectedInput)).Return(&ec2.CreateLaunchTemplateVersionOutput{
					LaunchTemplateVersion: &ec2types.LaunchTemplateVersion{
						LaunchTemplateId: aws.String("launch-template-id"),
					},
				}, nil).Do(
					func(ctx context.Context, arg *ec2.CreateLaunchTemplateVersionInput, requestOptions ...ec2.Options) {
						// formatting added to match tags slice during cmp.Equal()
						formatTagsInput(arg)
						if !cmp.Equal(expectedInput, arg, LaunchTemplateVersionIgnoreUnexported) {
							t.Fatalf("mismatch in input expected: %+v, but got %+v, diff: %s", expectedInput, arg, cmp.Diff(expectedInput, arg, LaunchTemplateVersionIgnoreUnexported))
						}
					})
			},
		},
		{
			name:                 "Should successfully create launch template version with capacity reservation ID and preference",
			awsResourceReference: []infrav1.AWSResourceReference{{ID: aws.String("1")}},
			mpScopeUpdater: func(mps *scope.MachinePoolScope) {
				spec := mps.AWSMachinePool.Spec
				spec.AWSLaunchTemplate.CapacityReservationID = aws.String("cr-12345678901234567")
				spec.AWSLaunchTemplate.CapacityReservationPreference = infrav1.CapacityReservationPreferenceOnly
				spec.AWSLaunchTemplate.SpotMarketOptions = nil
				mps.AWSMachinePool.Spec = spec
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				sgMap := make(map[infrav1.SecurityGroupRole]infrav1.SecurityGroup)
				sgMap[infrav1.SecurityGroupNode] = infrav1.SecurityGroup{ID: "1"}
				sgMap[infrav1.SecurityGroupLB] = infrav1.SecurityGroup{ID: "2"}

				expectedInput := &ec2.CreateLaunchTemplateVersionInput{
					LaunchTemplateData: &ec2types.RequestLaunchTemplateData{
						InstanceType: ec2types.InstanceTypeT3Large,
						IamInstanceProfile: &ec2types.LaunchTemplateIamInstanceProfileSpecificationRequest{
							Name: aws.String("instance-profile"),
						},
						KeyName:          aws.String("default"),
						UserData:         ptr.To[string](base64.StdEncoding.EncodeToString(userData)),
						SecurityGroupIds: []string{"nodeSG", "lbSG", "1"},
						ImageId:          aws.String("imageID"),
						CapacityReservationSpecification: &ec2types.LaunchTemplateCapacityReservationSpecificationRequest{
							CapacityReservationTarget: &ec2types.CapacityReservationTarget{
								CapacityReservationId: aws.String("cr-12345678901234567"),
							},
							CapacityReservationPreference: ec2types.CapacityReservationPreferenceCapacityReservationsOnly,
						},
						TagSpecifications: []ec2types.LaunchTemplateTagSpecificationRequest{
							{
								ResourceType: ec2types.ResourceTypeInstance,
								Tags:         defaultEC2AndDataTags("aws-mp-name", "cluster-name", userDataSecretKey, testBootstrapDataHash),
							},
							{
								ResourceType: ec2types.ResourceTypeVolume,
								Tags:         defaultEC2Tags("aws-mp-name", "cluster-name"),
							},
						},
					},
					LaunchTemplateId: aws.String("launch-template-id"),
				}
				m.CreateLaunchTemplateVersion(context.TODO(), gomock.AssignableToTypeOf(expectedInput)).Return(&ec2.CreateLaunchTemplateVersionOutput{
					LaunchTemplateVersion: &ec2types.LaunchTemplateVersion{
						LaunchTemplateId: aws.String("launch-template-id"),
					},
				}, nil).Do(
					func(ctx context.Context, arg *ec2.CreateLaunchTemplateVersionInput, requestOptions ...ec2.Options) {
						// formatting added to match tags slice during cmp.Equal()
						formatTagsInput(arg)
						if !cmp.Equal(expectedInput, arg, LaunchTemplateVersionIgnoreUnexported) {
							t.Fatalf("mismatch in input expected: %+v, but got %+v, diff: %s", expectedInput, arg, cmp.Diff(expectedInput, arg, LaunchTemplateVersionIgnoreUnexported))
						}
					})
			},
		},
		{
			name:                 "Should return error if AWS failed during launch template version creation",
			awsResourceReference: []infrav1.AWSResourceReference{{ID: aws.String("1")}},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				sgMap := make(map[infrav1.SecurityGroupRole]infrav1.SecurityGroup)
				sgMap[infrav1.SecurityGroupNode] = infrav1.SecurityGroup{ID: "1"}
				sgMap[infrav1.SecurityGroupLB] = infrav1.SecurityGroup{ID: "2"}

				expectedInput := &ec2.CreateLaunchTemplateVersionInput{
					LaunchTemplateData: &ec2types.RequestLaunchTemplateData{
						InstanceType: ec2types.InstanceTypeT3Large,
						IamInstanceProfile: &ec2types.LaunchTemplateIamInstanceProfileSpecificationRequest{
							Name: aws.String("instance-profile"),
						},
						KeyName:          aws.String("default"),
						UserData:         ptr.To[string](base64.StdEncoding.EncodeToString(userData)),
						SecurityGroupIds: []string{"nodeSG", "lbSG", "1"},
						ImageId:          aws.String("imageID"),
						InstanceMarketOptions: &ec2types.LaunchTemplateInstanceMarketOptionsRequest{
							MarketType: ec2types.MarketTypeSpot,
							SpotOptions: &ec2types.LaunchTemplateSpotMarketOptionsRequest{
								MaxPrice: aws.String("0.9"),
							},
						},
						TagSpecifications: []ec2types.LaunchTemplateTagSpecificationRequest{
							{
								ResourceType: ec2types.ResourceTypeInstance,
								Tags:         defaultEC2AndDataTags("aws-mp-name", "cluster-name", userDataSecretKey, testBootstrapDataHash),
							},
							{
								ResourceType: ec2types.ResourceTypeVolume,
								Tags:         defaultEC2Tags("aws-mp-name", "cluster-name"),
							},
						},
					},
					LaunchTemplateId: aws.String("launch-template-id"),
				}
				m.CreateLaunchTemplateVersion(context.TODO(), gomock.AssignableToTypeOf(expectedInput)).Return(nil,
					awserrors.NewFailedDependency("dependency failure")).Do(
					func(ctx context.Context, arg *ec2.CreateLaunchTemplateVersionInput, requestOptions ...ec2.Options) {
						// formatting added to match tags slice during cmp.Equal()
						formatTagsInput(arg)
						if !cmp.Equal(expectedInput, arg, LaunchTemplateVersionIgnoreUnexported) {
							t.Fatalf("mismatch in input expected: %+v, got: %+v", expectedInput, arg)
						}
					})
			},
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			scheme, err := setupScheme()
			g.Expect(err).NotTo(HaveOccurred())
			client := fake.NewClientBuilder().WithScheme(scheme).Build()

			cs, err := setupClusterScope(client)
			g.Expect(err).NotTo(HaveOccurred())

			ms, err := setupMachinePoolScope(client, cs)
			g.Expect(err).NotTo(HaveOccurred())
			if updateScope := tc.mpScopeUpdater; updateScope != nil {
				updateScope(ms)
			}

			ms.AWSMachinePool.Spec.AWSLaunchTemplate.AdditionalSecurityGroups = tc.awsResourceReference

			mockEC2Client := mocks.NewMockEC2API(mockCtrl)
			s := NewService(cs)
			s.EC2Client = mockEC2Client

			if tc.expect != nil {
				tc.expect(mockEC2Client.EXPECT())
			}
			if tc.wantErr {
				g.Expect(s.CreateLaunchTemplateVersion("launch-template-id", ms, aws.String("imageID"), userDataSecretKey, userData, testBootstrapDataHash)).To(HaveOccurred())
				return
			}
			g.Expect(s.CreateLaunchTemplateVersion("launch-template-id", ms, aws.String("imageID"), userDataSecretKey, userData, testBootstrapDataHash)).NotTo(HaveOccurred())
		})
	}
}

func TestBuildLaunchTemplateTagSpecificationRequest(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userDataSecretKey := types.NamespacedName{
		Namespace: "bootstrap-secret-ns",
		Name:      "bootstrap-secret",
	}
	bootstrapDataHash := userdata.ComputeHash([]byte("shell-script"))
	testCases := []struct {
		name  string
		check func(g *WithT, m []ec2types.LaunchTemplateTagSpecificationRequest)
	}{
		{
			name: "Should create tag specification request for building Launch template tags",
			check: func(g *WithT, res []ec2types.LaunchTemplateTagSpecificationRequest) {
				expected := []ec2types.LaunchTemplateTagSpecificationRequest{
					{
						ResourceType: ec2types.ResourceTypeInstance,
						Tags:         defaultEC2AndDataTags("aws-mp-name", "cluster-name", userDataSecretKey, bootstrapDataHash),
					},
					{
						ResourceType: ec2types.ResourceTypeVolume,
						Tags:         defaultEC2Tags("aws-mp-name", "cluster-name"),
					},
				}
				// sorting tags for comparing each request tags during cmp.Equal()
				for _, each := range res {
					sortTags(each.Tags)
				}
				g.Expect(res).Should(Equal(expected))
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			scheme, err := setupScheme()
			g.Expect(err).NotTo(HaveOccurred())
			client := fake.NewClientBuilder().WithScheme(scheme).Build()

			cs, err := setupClusterScope(client)
			g.Expect(err).NotTo(HaveOccurred())

			ms, err := setupMachinePoolScope(client, cs)
			g.Expect(err).NotTo(HaveOccurred())

			s := NewService(cs)
			tc.check(g, s.buildLaunchTemplateTagSpecificationRequest(ms, userDataSecretKey, bootstrapDataHash))
		})
	}
}

func TestDiscoverLaunchTemplateAMI(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name              string
		awsLaunchTemplate expinfrav1.AWSLaunchTemplate
		machineTemplate   clusterv1.MachineTemplateSpec
		expect            func(m *mocks.MockEC2APIMockRecorder)
		check             func(*WithT, *string, error)
	}{
		{
			name: "Should return default AMI for non EKS managed cluster if Image lookup format, org and BaseOS passed",
			awsLaunchTemplate: expinfrav1.AWSLaunchTemplate{
				Name:              "aws-launch-tmpl",
				ImageLookupFormat: "ilf",
				ImageLookupOrg:    "ilo",
				ImageLookupBaseOS: "ilbo",
				InstanceType:      "m5.large",
			},
			machineTemplate: clusterv1.MachineTemplateSpec{
				Spec: clusterv1.MachineSpec{
					Version: DefaultAmiNameFormat,
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeImages(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeImagesInput{})).
					Return(&ec2.DescribeImagesOutput{
						Images: []ec2types.Image{
							{
								ImageId:      aws.String("ancient"),
								CreationDate: aws.String("2011-02-08T17:02:31.000Z"),
							},
							{
								ImageId:      aws.String("latest"),
								CreationDate: aws.String("2019-02-08T17:02:31.000Z"),
							},
							{
								ImageId:      aws.String("oldest"),
								CreationDate: aws.String("2014-02-08T17:02:31.000Z"),
							},
						},
					}, nil)
				m.DescribeInstanceTypes(context.TODO(), gomock.Eq(&ec2.DescribeInstanceTypesInput{
					InstanceTypes: []ec2types.InstanceType{
						ec2types.InstanceTypeM5Large,
					},
				})).
					Return(&ec2.DescribeInstanceTypesOutput{
						InstanceTypes: []ec2types.InstanceTypeInfo{
							{
								ProcessorInfo: &ec2types.ProcessorInfo{
									SupportedArchitectures: []ec2types.ArchitectureType{
										ec2types.ArchitectureTypeX8664,
									},
								},
							},
						},
					}, nil)
			},
			check: func(g *WithT, res *string, err error) {
				g.Expect(res).Should(Equal(aws.String("latest")))
				g.Expect(err).NotTo(HaveOccurred())
			},
		},
		{
			name: "Should return AMI and use infra cluster image details, if not passed in aws launchtemplate",
			awsLaunchTemplate: expinfrav1.AWSLaunchTemplate{
				Name:         "aws-launch-tmpl",
				InstanceType: "m5.large",
			},
			machineTemplate: clusterv1.MachineTemplateSpec{
				Spec: clusterv1.MachineSpec{
					Version: DefaultAmiNameFormat,
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeImages(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeImagesInput{})).
					Return(&ec2.DescribeImagesOutput{
						Images: []ec2types.Image{
							{
								ImageId:      aws.String("ancient"),
								CreationDate: aws.String("2011-02-08T17:02:31.000Z"),
							},
							{
								ImageId:      aws.String("latest"),
								CreationDate: aws.String("2019-02-08T17:02:31.000Z"),
							},
							{
								ImageId:      aws.String("oldest"),
								CreationDate: aws.String("2014-02-08T17:02:31.000Z"),
							},
						},
					}, nil)
				m.DescribeInstanceTypes(context.TODO(), gomock.Eq(&ec2.DescribeInstanceTypesInput{
					InstanceTypes: []ec2types.InstanceType{
						ec2types.InstanceTypeM5Large,
					},
				})).
					Return(&ec2.DescribeInstanceTypesOutput{
						InstanceTypes: []ec2types.InstanceTypeInfo{
							{
								ProcessorInfo: &ec2types.ProcessorInfo{
									SupportedArchitectures: []ec2types.ArchitectureType{
										ec2types.ArchitectureTypeX8664,
									},
								},
							},
						},
					}, nil)
			},
			check: func(g *WithT, res *string, err error) {
				g.Expect(res).Should(Equal(aws.String("latest")))
				g.Expect(err).NotTo(HaveOccurred())
			},
		},
		{
			name: "Should return arm64 AMI and use infra cluster image details, if not passed in aws launchtemplate",
			awsLaunchTemplate: expinfrav1.AWSLaunchTemplate{
				Name:         "aws-launch-tmpl",
				InstanceType: "t4g.large",
			},
			machineTemplate: clusterv1.MachineTemplateSpec{
				Spec: clusterv1.MachineSpec{
					Version: DefaultAmiNameFormat,
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeImages(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeImagesInput{})).
					Return(&ec2.DescribeImagesOutput{
						Images: []ec2types.Image{
							{
								ImageId:      aws.String("ancient"),
								CreationDate: aws.String("2011-02-08T17:02:31.000Z"),
							},
							{
								ImageId:      aws.String("latest"),
								CreationDate: aws.String("2019-02-08T17:02:31.000Z"),
							},
							{
								ImageId:      aws.String("oldest"),
								CreationDate: aws.String("2014-02-08T17:02:31.000Z"),
							},
						},
					}, nil)
				m.DescribeInstanceTypes(context.TODO(), gomock.Eq(&ec2.DescribeInstanceTypesInput{
					InstanceTypes: []ec2types.InstanceType{
						ec2types.InstanceTypeT4gLarge,
					},
				})).
					Return(&ec2.DescribeInstanceTypesOutput{
						InstanceTypes: []ec2types.InstanceTypeInfo{
							{
								ProcessorInfo: &ec2types.ProcessorInfo{
									SupportedArchitectures: []ec2types.ArchitectureType{
										ec2types.ArchitectureTypeArm64,
									},
								},
							},
						},
					}, nil)
			},
			check: func(g *WithT, res *string, err error) {
				g.Expect(res).Should(Equal(aws.String("latest")))
				g.Expect(err).NotTo(HaveOccurred())
			},
		},
		{
			name: "Should return AWSlaunchtemplate ID if provided",
			awsLaunchTemplate: expinfrav1.AWSLaunchTemplate{
				Name: "aws-launch-tmpl",
				AMI:  infrav1.AMIReference{ID: aws.String("id")},
			},
			check: func(g *WithT, res *string, err error) {
				g.Expect(res).Should(Equal(aws.String("id")))
				g.Expect(err).NotTo(HaveOccurred())
			},
		},
		{
			name: "Should return with error if both AWSlaunchtemplate ID and machinePool version is not provided",
			awsLaunchTemplate: expinfrav1.AWSLaunchTemplate{
				Name: "aws-launch-tmpl",
			},
			check: func(g *WithT, res *string, err error) {
				g.Expect(err).To(HaveOccurred())
				g.Expect(res).To(BeNil())
			},
		},
		{
			name: "Should return error if AWS failed while describing images",
			awsLaunchTemplate: expinfrav1.AWSLaunchTemplate{
				Name:         "aws-launch-tmpl",
				InstanceType: "m5.large",
			},
			machineTemplate: clusterv1.MachineTemplateSpec{
				Spec: clusterv1.MachineSpec{
					Version: DefaultAmiNameFormat,
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeImages(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeImagesInput{})).
					Return(nil, awserrors.NewFailedDependency("dependency-failure"))
				m.DescribeInstanceTypes(context.TODO(), gomock.Eq(&ec2.DescribeInstanceTypesInput{
					InstanceTypes: []ec2types.InstanceType{
						ec2types.InstanceTypeM5Large,
					},
				})).
					Return(&ec2.DescribeInstanceTypesOutput{
						InstanceTypes: []ec2types.InstanceTypeInfo{
							{
								ProcessorInfo: &ec2types.ProcessorInfo{
									SupportedArchitectures: []ec2types.ArchitectureType{
										ec2types.ArchitectureTypeX8664,
									},
								},
							},
						},
					}, nil)
			},
			check: func(g *WithT, res *string, err error) {
				g.Expect(res).To(BeNil())
				g.Expect(err).To(HaveOccurred())
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			ec2Mock := mocks.NewMockEC2API(mockCtrl)

			scheme, err := setupScheme()
			g.Expect(err).NotTo(HaveOccurred())
			client := fake.NewClientBuilder().WithScheme(scheme).Build()

			cs, err := setupClusterScope(client)
			g.Expect(err).NotTo(HaveOccurred())

			ms, err := setupMachinePoolScope(client, cs)
			g.Expect(err).NotTo(HaveOccurred())

			ms.AWSMachinePool.Spec.AWSLaunchTemplate = tc.awsLaunchTemplate
			ms.MachinePool.Spec.Template = tc.machineTemplate

			if tc.expect != nil {
				tc.expect(ec2Mock.EXPECT())
			}

			s := NewService(cs)
			s.EC2Client = ec2Mock

			id, err := s.DiscoverLaunchTemplateAMI(context.TODO(), ms)
			tc.check(g, id, err)
		})
	}
}

func TestDiscoverLaunchTemplateAMIForEKS(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name              string
		awsLaunchTemplate expinfrav1.AWSLaunchTemplate
		machineTemplate   clusterv1.MachineTemplateSpec
		expectEC2         func(m *mocks.MockEC2APIMockRecorder)
		expectSSM         func(m *mock_ssmiface.MockSSMAPIMockRecorder)
		check             func(*WithT, *string, error)
	}{
		{
			name: "Should return AMI and use EKS infra cluster image details, if not passed in aws launch template",
			expectEC2: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeInstanceTypes(context.TODO(), gomock.Any()).
					Return(&ec2.DescribeInstanceTypesOutput{
						InstanceTypes: []ec2types.InstanceTypeInfo{
							{
								ProcessorInfo: &ec2types.ProcessorInfo{
									SupportedArchitectures: []ec2types.ArchitectureType{
										ec2types.ArchitectureTypeX8664,
									},
								},
							},
						},
					}, nil)
			},
			expectSSM: func(m *mock_ssmiface.MockSSMAPIMockRecorder) {
				m.GetParameter(context.TODO(), gomock.AssignableToTypeOf(&ssm.GetParameterInput{})).
					Return(&ssm.GetParameterOutput{
						Parameter: &ssmtypes.Parameter{
							Value: aws.String("latest"),
						},
					}, nil)
			},
			check: func(g *WithT, res *string, err error) {
				g.Expect(res).Should(Equal(aws.String("latest")))
				g.Expect(err).NotTo(HaveOccurred())
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			ssmMock := mock_ssmiface.NewMockSSMAPI(mockCtrl)
			ec2Mock := mocks.NewMockEC2API(mockCtrl)

			scheme, err := setupScheme()
			g.Expect(err).NotTo(HaveOccurred())
			client := fake.NewClientBuilder().WithScheme(scheme).Build()

			mcps, err := setupNewManagedControlPlaneScope(client)
			g.Expect(err).NotTo(HaveOccurred())

			ms, err := setupMachinePoolScope(client, mcps)
			g.Expect(err).NotTo(HaveOccurred())

			if tc.expectEC2 != nil {
				tc.expectEC2(ec2Mock.EXPECT())
			}

			if tc.expectSSM != nil {
				tc.expectSSM(ssmMock.EXPECT())
			}

			s := NewService(mcps)
			s.EC2Client = ec2Mock
			s.SSMClient = ssmMock

			id, err := s.DiscoverLaunchTemplateAMI(context.TODO(), ms)
			tc.check(g, id, err)
		})
	}
}

func TestDeleteLaunchTemplateVersion(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type args struct {
		id      string
		version *int64
	}
	testCases := []struct {
		name    string
		args    args
		expect  func(m *mocks.MockEC2APIMockRecorder)
		wantErr bool
	}{
		{
			name:    "Should return error if version is nil",
			wantErr: true,
		},
		{
			name: "Should return error if AWS unable to delete launch template version",
			args: args{
				id:      "id",
				version: aws.Int64(12),
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DeleteLaunchTemplateVersions(context.TODO(), gomock.Eq(
					&ec2.DeleteLaunchTemplateVersionsInput{
						LaunchTemplateId: aws.String("id"),
						Versions:         []string{"12"},
					},
				)).Return(nil, awserrors.NewFailedDependency("dependency-failure"))
			},
			wantErr: true,
		},
		{
			name: "Should successfully deletes launch template version if AWS call passed",
			args: args{
				id:      "id",
				version: aws.Int64(12),
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DeleteLaunchTemplateVersions(context.TODO(), gomock.Eq(
					&ec2.DeleteLaunchTemplateVersionsInput{
						LaunchTemplateId: aws.String("id"),
						Versions:         []string{"12"},
					},
				)).Return(nil, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			scheme, err := setupScheme()
			g.Expect(err).NotTo(HaveOccurred())
			client := fake.NewClientBuilder().WithScheme(scheme).Build()

			cs, err := setupClusterScope(client)
			g.Expect(err).NotTo(HaveOccurred())

			ec2Mock := mocks.NewMockEC2API(mockCtrl)
			s := NewService(cs)
			s.EC2Client = ec2Mock

			if tc.expect != nil {
				tc.expect(ec2Mock.EXPECT())
			}

			if tc.wantErr {
				g.Expect(s.deleteLaunchTemplateVersion(tc.args.id, tc.args.version)).To(HaveOccurred())
				return
			}
			g.Expect(s.deleteLaunchTemplateVersion(tc.args.id, tc.args.version)).NotTo(HaveOccurred())
		})
	}
}
