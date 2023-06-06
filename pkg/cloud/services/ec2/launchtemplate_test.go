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
	"encoding/base64"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/ssm/mock_ssmiface"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/userdata"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
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

var (
	testUserDataHash = userdata.ComputeHash([]byte(testUserData))
)

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
				m.DescribeLaunchTemplateVersions(gomock.Eq(&ec2.DescribeLaunchTemplateVersionsInput{
					LaunchTemplateName: aws.String("foo"),
					Versions:           []*string{aws.String("$Latest")},
				})).
					Return(nil, awserr.New(
						awserrors.LaunchTemplateNameNotFound,
						"The specified launch template, with template name foo, does not exist.",
						nil,
					))
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
				m.DescribeLaunchTemplateVersions(gomock.Eq(&ec2.DescribeLaunchTemplateVersionsInput{
					LaunchTemplateName: aws.String("foo"),
					Versions:           []*string{aws.String("$Latest")},
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
				m.DescribeLaunchTemplateVersions(gomock.Eq(&ec2.DescribeLaunchTemplateVersionsInput{
					LaunchTemplateName: aws.String("foo"),
					Versions:           []*string{aws.String("$Latest")},
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
				m.DescribeLaunchTemplateVersions(gomock.Eq(&ec2.DescribeLaunchTemplateVersionsInput{
					LaunchTemplateName: aws.String("foo"),
					Versions:           []*string{aws.String("$Latest")},
				})).Return(&ec2.DescribeLaunchTemplateVersionsOutput{
					LaunchTemplateVersions: []*ec2.LaunchTemplateVersion{
						{
							LaunchTemplateId:   aws.String("lt-12345"),
							LaunchTemplateName: aws.String("foo"),
							LaunchTemplateData: &ec2.ResponseLaunchTemplateData{
								SecurityGroupIds: []*string{aws.String("sg-id")},
								ImageId:          aws.String("foo-image"),
								IamInstanceProfile: &ec2.LaunchTemplateIamInstanceProfileSpecification{
									Arn: aws.String("instance-profile/foo-profile"),
								},
								KeyName: aws.String("foo-keyname"),
								BlockDeviceMappings: []*ec2.LaunchTemplateBlockDeviceMapping{
									{
										DeviceName: aws.String("foo-device"),
										Ebs: &ec2.LaunchTemplateEbsBlockDevice{
											Encrypted:  aws.Bool(true),
											VolumeSize: aws.Int64(16),
											VolumeType: aws.String("cool"),
										},
									},
								},
								NetworkInterfaces: []*ec2.LaunchTemplateInstanceNetworkInterfaceSpecification{
									{
										DeviceIndex: aws.Int64(1),
										Groups:      []*string{aws.String("foo-group")},
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
				m.DescribeLaunchTemplateVersions(gomock.Eq(&ec2.DescribeLaunchTemplateVersionsInput{
					LaunchTemplateName: aws.String("foo"),
					Versions:           []*string{aws.String("$Latest")},
				})).Return(&ec2.DescribeLaunchTemplateVersionsOutput{
					LaunchTemplateVersions: []*ec2.LaunchTemplateVersion{
						{
							LaunchTemplateId:   aws.String("lt-12345"),
							LaunchTemplateName: aws.String("foo"),
							LaunchTemplateData: &ec2.ResponseLaunchTemplateData{
								SecurityGroupIds: []*string{aws.String("sg-id")},
								ImageId:          aws.String("foo-image"),
								IamInstanceProfile: &ec2.LaunchTemplateIamInstanceProfileSpecification{
									Arn: aws.String("instance-profile/foo-profile"),
								},
								KeyName: aws.String("foo-keyname"),
								BlockDeviceMappings: []*ec2.LaunchTemplateBlockDeviceMapping{
									{
										DeviceName: aws.String("foo-device"),
										Ebs: &ec2.LaunchTemplateEbsBlockDevice{
											Encrypted:  aws.Bool(true),
											VolumeSize: aws.Int64(16),
											VolumeType: aws.String("cool"),
										},
									},
								},
								NetworkInterfaces: []*ec2.LaunchTemplateInstanceNetworkInterfaceSpecification{
									{
										DeviceIndex: aws.Int64(1),
										Groups:      []*string{aws.String("foo-group")},
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

			launchTemplate, userData, err := s.GetLaunchTemplate(tc.launchTemplateName)
			tc.check(g, launchTemplate, userData, err)
		})
	}
}

func TestServiceSDKToLaunchTemplate(t *testing.T) {
	tests := []struct {
		name     string
		input    *ec2.LaunchTemplateVersion
		wantLT   *expinfrav1.AWSLaunchTemplate
		wantHash string
		wantErr  bool
	}{
		{
			name: "lots of input",
			input: &ec2.LaunchTemplateVersion{
				LaunchTemplateId:   aws.String("lt-12345"),
				LaunchTemplateName: aws.String("foo"),
				LaunchTemplateData: &ec2.ResponseLaunchTemplateData{
					ImageId: aws.String("foo-image"),
					IamInstanceProfile: &ec2.LaunchTemplateIamInstanceProfileSpecification{
						Arn: aws.String("instance-profile/foo-profile"),
					},
					KeyName: aws.String("foo-keyname"),
					BlockDeviceMappings: []*ec2.LaunchTemplateBlockDeviceMapping{
						{
							DeviceName: aws.String("foo-device"),
							Ebs: &ec2.LaunchTemplateEbsBlockDevice{
								Encrypted:  aws.Bool(true),
								VolumeSize: aws.Int64(16),
								VolumeType: aws.String("cool"),
							},
						},
					},
					NetworkInterfaces: []*ec2.LaunchTemplateInstanceNetworkInterfaceSpecification{
						{
							DeviceIndex: aws.Int64(1),
							Groups:      []*string{aws.String("foo-group")},
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
			wantHash: testUserDataHash,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{}
			gotLT, gotHash, err := s.SDKToLaunchTemplate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error mismatch: got %v, wantErr %v", err, tt.wantErr)
			}
			if !cmp.Equal(gotLT, tt.wantLT) {
				t.Fatalf("launchTemplate mismatch: got %v, want %v", gotLT, tt.wantLT)
			}
			if !cmp.Equal(gotHash, tt.wantHash) {
				t.Fatalf("userDataHash mismatch: got %v, want %v", gotHash, tt.wantHash)
			}
		})
	}
}

func TestServiceLaunchTemplateNeedsUpdate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name     string
		incoming *expinfrav1.AWSLaunchTemplate
		existing *expinfrav1.AWSLaunchTemplate
		expect   func(m *mocks.MockEC2APIMockRecorder)
		want     bool
		wantErr  bool
	}{
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
			want:    true,
			wantErr: false,
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
			want:    true,
			wantErr: false,
		},
		{
			name: "Should return true if incoming IamInstanceProfile is not same as existing IamInstanceProfile",
			incoming: &expinfrav1.AWSLaunchTemplate{
				IamInstanceProfile: DefaultAmiNameFormat,
			},
			existing: &expinfrav1.AWSLaunchTemplate{
				IamInstanceProfile: "some-other-profile",
			},
			want: true,
		},
		{
			name: "Should return true if incoming InstanceType is not same as existing InstanceType",
			incoming: &expinfrav1.AWSLaunchTemplate{
				InstanceType: "t3.micro",
			},
			existing: &expinfrav1.AWSLaunchTemplate{
				InstanceType: "t3.large",
			},
			want: true,
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
				m.DescribeSecurityGroups(gomock.Eq(&ec2.DescribeSecurityGroupsInput{Filters: []*ec2.Filter{{Name: aws.String("sg-1"), Values: aws.StringSlice([]string{"test-1"})}}})).
					Return(&ec2.DescribeSecurityGroupsOutput{SecurityGroups: []*ec2.SecurityGroup{{GroupId: aws.String("sg-1")}}}, nil)
				m.DescribeSecurityGroups(gomock.Eq(&ec2.DescribeSecurityGroupsInput{Filters: []*ec2.Filter{{Name: aws.String("sg-2"), Values: aws.StringSlice([]string{"test-2"})}}})).
					Return(&ec2.DescribeSecurityGroupsOutput{SecurityGroups: []*ec2.SecurityGroup{{GroupId: aws.String("sg-2")}}}, nil)
			},
			want:    true,
			wantErr: false,
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

			got, err := s.LaunchTemplateNeedsUpdate(machinePoolScope, tt.incoming, tt.existing)
			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).NotTo(HaveOccurred())
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
				m.DescribeLaunchTemplateVersions(gomock.Eq(&ec2.DescribeLaunchTemplateVersionsInput{
					LaunchTemplateName: aws.String("foo"),
					Versions:           []*string{aws.String("$Latest")},
				})).Return(nil, awserr.New(
					awserrors.LaunchTemplateNameNotFound,
					"The specified launch template, with template name foo, does not exist.",
					nil,
				))
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
				m.DescribeLaunchTemplateVersions(gomock.Eq(&ec2.DescribeLaunchTemplateVersionsInput{
					LaunchTemplateName: aws.String("foo"),
					Versions:           []*string{aws.String("$Latest")},
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
				m.DescribeLaunchTemplateVersions(gomock.Eq(&ec2.DescribeLaunchTemplateVersionsInput{
					LaunchTemplateName: aws.String("foo"),
					Versions:           []*string{aws.String("$Latest")},
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
				m.DescribeLaunchTemplateVersions(gomock.Eq(&ec2.DescribeLaunchTemplateVersionsInput{
					LaunchTemplateName: aws.String("foo"),
					Versions:           []*string{aws.String("$Latest")},
				})).Return(&ec2.DescribeLaunchTemplateVersionsOutput{
					LaunchTemplateVersions: []*ec2.LaunchTemplateVersion{
						{
							LaunchTemplateId:   aws.String("lt-12345"),
							LaunchTemplateName: aws.String("foo"),
							LaunchTemplateData: &ec2.ResponseLaunchTemplateData{
								ImageId: aws.String("foo-image"),
								IamInstanceProfile: &ec2.LaunchTemplateIamInstanceProfileSpecification{
									Arn: aws.String("instance-profile/foo-profile"),
								},
								KeyName: aws.String("foo-keyname"),
								BlockDeviceMappings: []*ec2.LaunchTemplateBlockDeviceMapping{
									{
										DeviceName: aws.String("foo-device"),
										Ebs: &ec2.LaunchTemplateEbsBlockDevice{
											Encrypted:  aws.Bool(true),
											VolumeSize: aws.Int64(16),
											VolumeType: aws.String("cool"),
										},
									},
								},
								NetworkInterfaces: []*ec2.LaunchTemplateInstanceNetworkInterfaceSpecification{
									{
										DeviceIndex: aws.Int64(1),
										Groups:      []*string{aws.String("foo-group")},
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
				m.DeleteLaunchTemplate(gomock.Eq(&ec2.DeleteLaunchTemplateInput{
					LaunchTemplateId: aws.String("1"),
				})).Return(&ec2.DeleteLaunchTemplateOutput{}, nil)
			},
		},
		{
			name:      "Should return error if failed to delete given launch template ID",
			versionID: "1",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DeleteLaunchTemplate(gomock.Eq(&ec2.DeleteLaunchTemplateInput{
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

	var formatTagsInput = func(arg *ec2.CreateLaunchTemplateInput) {
		sortTags(arg.TagSpecifications[0].Tags)

		for index := range arg.LaunchTemplateData.TagSpecifications {
			sortTags(arg.LaunchTemplateData.TagSpecifications[index].Tags)
		}
	}

	var userData = []byte{1, 0, 0}
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

				var expectedInput = &ec2.CreateLaunchTemplateInput{
					LaunchTemplateData: &ec2.RequestLaunchTemplateData{
						InstanceType: aws.String("t3.large"),
						IamInstanceProfile: &ec2.LaunchTemplateIamInstanceProfileSpecificationRequest{
							Name: aws.String("instance-profile"),
						},
						KeyName:          aws.String("default"),
						UserData:         pointer.String(base64.StdEncoding.EncodeToString(userData)),
						SecurityGroupIds: aws.StringSlice([]string{"nodeSG", "lbSG", "1"}),
						ImageId:          aws.String("imageID"),
						InstanceMarketOptions: &ec2.LaunchTemplateInstanceMarketOptionsRequest{
							MarketType: aws.String("spot"),
							SpotOptions: &ec2.LaunchTemplateSpotMarketOptionsRequest{
								MaxPrice: aws.String("0.9"),
							},
						},
						TagSpecifications: []*ec2.LaunchTemplateTagSpecificationRequest{
							{
								ResourceType: aws.String(ec2.ResourceTypeInstance),
								Tags:         defaultEC2Tags("aws-mp-name", "cluster-name"),
							},
							{
								ResourceType: aws.String(ec2.ResourceTypeVolume),
								Tags:         defaultEC2Tags("aws-mp-name", "cluster-name"),
							},
						},
					},
					LaunchTemplateName: aws.String("aws-mp-name"),
					TagSpecifications: []*ec2.TagSpecification{
						{
							ResourceType: aws.String(ec2.ResourceTypeLaunchTemplate),
							Tags:         defaultEC2Tags("aws-mp-name", "cluster-name"),
						},
					},
				}
				m.CreateLaunchTemplate(gomock.AssignableToTypeOf(expectedInput)).Return(&ec2.CreateLaunchTemplateOutput{
					LaunchTemplate: &ec2.LaunchTemplate{
						LaunchTemplateId: aws.String("launch-template-id"),
					},
				}, nil).Do(func(arg *ec2.CreateLaunchTemplateInput) {
					// formatting added to match arrays during cmp.Equal
					formatTagsInput(arg)
					if !cmp.Equal(expectedInput, arg) {
						t.Fatalf("mismatch in input expected: %+v, got: %+v", expectedInput, arg)
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

				var expectedInput = &ec2.CreateLaunchTemplateInput{
					LaunchTemplateData: &ec2.RequestLaunchTemplateData{
						InstanceType: aws.String("t3.large"),
						IamInstanceProfile: &ec2.LaunchTemplateIamInstanceProfileSpecificationRequest{
							Name: aws.String("instance-profile"),
						},
						KeyName:          aws.String("default"),
						UserData:         pointer.String(base64.StdEncoding.EncodeToString(userData)),
						SecurityGroupIds: aws.StringSlice([]string{"nodeSG", "lbSG", "sg-1"}),
						ImageId:          aws.String("imageID"),
						InstanceMarketOptions: &ec2.LaunchTemplateInstanceMarketOptionsRequest{
							MarketType: aws.String("spot"),
							SpotOptions: &ec2.LaunchTemplateSpotMarketOptionsRequest{
								MaxPrice: aws.String("0.9"),
							},
						},
						TagSpecifications: []*ec2.LaunchTemplateTagSpecificationRequest{
							{
								ResourceType: aws.String(ec2.ResourceTypeInstance),
								Tags:         defaultEC2Tags("aws-mp-name", "cluster-name"),
							},
							{
								ResourceType: aws.String(ec2.ResourceTypeVolume),
								Tags:         defaultEC2Tags("aws-mp-name", "cluster-name"),
							},
						},
					},
					LaunchTemplateName: aws.String("aws-mp-name"),
					TagSpecifications: []*ec2.TagSpecification{
						{
							ResourceType: aws.String(ec2.ResourceTypeLaunchTemplate),
							Tags:         defaultEC2Tags("aws-mp-name", "cluster-name"),
						},
					},
				}
				m.CreateLaunchTemplate(gomock.AssignableToTypeOf(expectedInput)).Return(&ec2.CreateLaunchTemplateOutput{
					LaunchTemplate: &ec2.LaunchTemplate{
						LaunchTemplateId: aws.String("launch-template-id"),
					},
				}, nil).Do(func(arg *ec2.CreateLaunchTemplateInput) {
					// formatting added to match arrays during reflect.DeepEqual
					formatTagsInput(arg)
					if !cmp.Equal(expectedInput, arg) {
						t.Fatalf("mismatch in input expected: %+v, got: %+v", expectedInput, arg)
					}
				})
				m.DescribeSecurityGroups(gomock.Eq(&ec2.DescribeSecurityGroupsInput{Filters: []*ec2.Filter{{Name: aws.String("sg-1"), Values: aws.StringSlice([]string{"test"})}}})).
					Return(&ec2.DescribeSecurityGroupsOutput{SecurityGroups: []*ec2.SecurityGroup{{GroupId: aws.String("sg-1")}}}, nil)
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

				var expectedInput = &ec2.CreateLaunchTemplateInput{
					LaunchTemplateData: &ec2.RequestLaunchTemplateData{
						InstanceType: aws.String("t3.large"),
						IamInstanceProfile: &ec2.LaunchTemplateIamInstanceProfileSpecificationRequest{
							Name: aws.String("instance-profile"),
						},
						KeyName:          aws.String("default"),
						UserData:         pointer.String(base64.StdEncoding.EncodeToString(userData)),
						SecurityGroupIds: aws.StringSlice([]string{"nodeSG", "lbSG", "1"}),
						ImageId:          aws.String("imageID"),
						InstanceMarketOptions: &ec2.LaunchTemplateInstanceMarketOptionsRequest{
							MarketType: aws.String("spot"),
							SpotOptions: &ec2.LaunchTemplateSpotMarketOptionsRequest{
								MaxPrice: aws.String("0.9"),
							},
						},
						TagSpecifications: []*ec2.LaunchTemplateTagSpecificationRequest{
							{
								ResourceType: aws.String(ec2.ResourceTypeInstance),
								Tags:         defaultEC2Tags("aws-mp-name", "cluster-name"),
							},
							{
								ResourceType: aws.String(ec2.ResourceTypeVolume),
								Tags:         defaultEC2Tags("aws-mp-name", "cluster-name"),
							},
						},
					},
					LaunchTemplateName: aws.String("aws-mp-name"),
					TagSpecifications: []*ec2.TagSpecification{
						{
							ResourceType: aws.String(ec2.ResourceTypeLaunchTemplate),
							Tags:         defaultEC2Tags("aws-mp-name", "cluster-name"),
						},
					},
				}
				m.CreateLaunchTemplate(gomock.AssignableToTypeOf(expectedInput)).Return(nil,
					awserrors.NewFailedDependency("dependency failure")).Do(func(arg *ec2.CreateLaunchTemplateInput) {
					// formatting added to match arrays during cmp.Equal
					formatTagsInput(arg)
					if !cmp.Equal(expectedInput, arg) {
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

			launchTemplate, err := s.CreateLaunchTemplate(ms, aws.String("imageID"), userData)
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

		launchTemplate, err := s.CreateLaunchTemplate(ms, aws.String("imageID"), nil)
		g.Expect(err).To(HaveOccurred())
		g.Expect(launchTemplate).Should(BeEmpty())
	})
}

func TestCreateLaunchTemplateVersion(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	var formatTagsInput = func(arg *ec2.CreateLaunchTemplateVersionInput) {
		for index := range arg.LaunchTemplateData.TagSpecifications {
			sortTags(arg.LaunchTemplateData.TagSpecifications[index].Tags)
		}
	}
	var userData = []byte{1, 0, 0}
	testCases := []struct {
		name                 string
		imageID              *string
		awsResourceReference []infrav1.AWSResourceReference
		expect               func(m *mocks.MockEC2APIMockRecorder)
		wantErr              bool
	}{
		{
			name:                 "Should successfully creates launch template version",
			awsResourceReference: []infrav1.AWSResourceReference{{ID: aws.String("1")}},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				sgMap := make(map[infrav1.SecurityGroupRole]infrav1.SecurityGroup)
				sgMap[infrav1.SecurityGroupNode] = infrav1.SecurityGroup{ID: "1"}
				sgMap[infrav1.SecurityGroupLB] = infrav1.SecurityGroup{ID: "2"}

				var expectedInput = &ec2.CreateLaunchTemplateVersionInput{
					LaunchTemplateData: &ec2.RequestLaunchTemplateData{
						InstanceType: aws.String("t3.large"),
						IamInstanceProfile: &ec2.LaunchTemplateIamInstanceProfileSpecificationRequest{
							Name: aws.String("instance-profile"),
						},
						KeyName:          aws.String("default"),
						UserData:         pointer.String(base64.StdEncoding.EncodeToString(userData)),
						SecurityGroupIds: aws.StringSlice([]string{"nodeSG", "lbSG", "1"}),
						ImageId:          aws.String("imageID"),
						InstanceMarketOptions: &ec2.LaunchTemplateInstanceMarketOptionsRequest{
							MarketType: aws.String("spot"),
							SpotOptions: &ec2.LaunchTemplateSpotMarketOptionsRequest{
								MaxPrice: aws.String("0.9"),
							},
						},
						TagSpecifications: []*ec2.LaunchTemplateTagSpecificationRequest{
							{
								ResourceType: aws.String(ec2.ResourceTypeInstance),
								Tags:         defaultEC2Tags("aws-mp-name", "cluster-name"),
							},
							{
								ResourceType: aws.String(ec2.ResourceTypeVolume),
								Tags:         defaultEC2Tags("aws-mp-name", "cluster-name"),
							},
						},
					},
					LaunchTemplateId: aws.String("launch-template-id"),
				}
				m.CreateLaunchTemplateVersion(gomock.AssignableToTypeOf(expectedInput)).Return(&ec2.CreateLaunchTemplateVersionOutput{
					LaunchTemplateVersion: &ec2.LaunchTemplateVersion{
						LaunchTemplateId: aws.String("launch-template-id"),
					},
				}, nil).Do(
					func(arg *ec2.CreateLaunchTemplateVersionInput) {
						// formatting added to match tags slice during cmp.Equal()
						formatTagsInput(arg)
						if !cmp.Equal(expectedInput, arg) {
							t.Fatalf("mismatch in input expected: %+v, but got %+v", expectedInput, arg)
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

				var expectedInput = &ec2.CreateLaunchTemplateVersionInput{
					LaunchTemplateData: &ec2.RequestLaunchTemplateData{
						InstanceType: aws.String("t3.large"),
						IamInstanceProfile: &ec2.LaunchTemplateIamInstanceProfileSpecificationRequest{
							Name: aws.String("instance-profile"),
						},
						KeyName:          aws.String("default"),
						UserData:         pointer.String(base64.StdEncoding.EncodeToString(userData)),
						SecurityGroupIds: aws.StringSlice([]string{"nodeSG", "lbSG", "1"}),
						ImageId:          aws.String("imageID"),
						InstanceMarketOptions: &ec2.LaunchTemplateInstanceMarketOptionsRequest{
							MarketType: aws.String("spot"),
							SpotOptions: &ec2.LaunchTemplateSpotMarketOptionsRequest{
								MaxPrice: aws.String("0.9"),
							},
						},
						TagSpecifications: []*ec2.LaunchTemplateTagSpecificationRequest{
							{
								ResourceType: aws.String(ec2.ResourceTypeInstance),
								Tags:         defaultEC2Tags("aws-mp-name", "cluster-name"),
							},
							{
								ResourceType: aws.String(ec2.ResourceTypeVolume),
								Tags:         defaultEC2Tags("aws-mp-name", "cluster-name"),
							},
						},
					},
					LaunchTemplateId: aws.String("launch-template-id"),
				}
				m.CreateLaunchTemplateVersion(gomock.AssignableToTypeOf(expectedInput)).Return(nil,
					awserrors.NewFailedDependency("dependency failure")).Do(
					func(arg *ec2.CreateLaunchTemplateVersionInput) {
						// formatting added to match tags slice during cmp.Equal()
						formatTagsInput(arg)
						if !cmp.Equal(expectedInput, arg) {
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

			ms.AWSMachinePool.Spec.AWSLaunchTemplate.AdditionalSecurityGroups = tc.awsResourceReference

			mockEC2Client := mocks.NewMockEC2API(mockCtrl)
			s := NewService(cs)
			s.EC2Client = mockEC2Client

			if tc.expect != nil {
				tc.expect(mockEC2Client.EXPECT())
			}
			if tc.wantErr {
				g.Expect(s.CreateLaunchTemplateVersion("launch-template-id", ms, aws.String("imageID"), userData)).To(HaveOccurred())
				return
			}
			g.Expect(s.CreateLaunchTemplateVersion("launch-template-id", ms, aws.String("imageID"), userData)).NotTo(HaveOccurred())
		})
	}
}

func TestBuildLaunchTemplateTagSpecificationRequest(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name  string
		check func(g *WithT, m []*ec2.LaunchTemplateTagSpecificationRequest)
	}{
		{
			name: "Should create tag specification request for building Launch template tags",
			check: func(g *WithT, res []*ec2.LaunchTemplateTagSpecificationRequest) {
				expected := []*ec2.LaunchTemplateTagSpecificationRequest{
					{
						ResourceType: aws.String(ec2.ResourceTypeInstance),
						Tags:         defaultEC2Tags("aws-mp-name", "cluster-name"),
					},
					{
						ResourceType: aws.String(ec2.ResourceTypeVolume),
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
			tc.check(g, s.buildLaunchTemplateTagSpecificationRequest(ms))
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
					Version: aws.String(DefaultAmiNameFormat),
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeImages(gomock.AssignableToTypeOf(&ec2.DescribeImagesInput{})).
					Return(&ec2.DescribeImagesOutput{
						Images: []*ec2.Image{
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
				m.DescribeInstanceTypes(gomock.Eq(&ec2.DescribeInstanceTypesInput{
					InstanceTypes: []*string{
						aws.String("m5.large"),
					},
				})).
					Return(&ec2.DescribeInstanceTypesOutput{
						InstanceTypes: []*ec2.InstanceTypeInfo{
							{
								ProcessorInfo: &ec2.ProcessorInfo{
									SupportedArchitectures: []*string{
										aws.String("x86_64"),
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
					Version: aws.String(DefaultAmiNameFormat),
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeImages(gomock.AssignableToTypeOf(&ec2.DescribeImagesInput{})).
					Return(&ec2.DescribeImagesOutput{
						Images: []*ec2.Image{
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
				m.DescribeInstanceTypes(gomock.Eq(&ec2.DescribeInstanceTypesInput{
					InstanceTypes: []*string{
						aws.String("m5.large"),
					},
				})).
					Return(&ec2.DescribeInstanceTypesOutput{
						InstanceTypes: []*ec2.InstanceTypeInfo{
							{
								ProcessorInfo: &ec2.ProcessorInfo{
									SupportedArchitectures: []*string{
										aws.String("x86_64"),
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
					Version: aws.String(DefaultAmiNameFormat),
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeImages(gomock.AssignableToTypeOf(&ec2.DescribeImagesInput{})).
					Return(&ec2.DescribeImagesOutput{
						Images: []*ec2.Image{
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
				m.DescribeInstanceTypes(gomock.Eq(&ec2.DescribeInstanceTypesInput{
					InstanceTypes: []*string{
						aws.String("t4g.large"),
					},
				})).
					Return(&ec2.DescribeInstanceTypesOutput{
						InstanceTypes: []*ec2.InstanceTypeInfo{
							{
								ProcessorInfo: &ec2.ProcessorInfo{
									SupportedArchitectures: []*string{
										aws.String("arm64"),
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
					Version: aws.String(DefaultAmiNameFormat),
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeImages(gomock.AssignableToTypeOf(&ec2.DescribeImagesInput{})).
					Return(nil, awserrors.NewFailedDependency("dependency-failure"))
				m.DescribeInstanceTypes(gomock.Eq(&ec2.DescribeInstanceTypesInput{
					InstanceTypes: []*string{
						aws.String("m5.large"),
					},
				})).
					Return(&ec2.DescribeInstanceTypesOutput{
						InstanceTypes: []*ec2.InstanceTypeInfo{
							{
								ProcessorInfo: &ec2.ProcessorInfo{
									SupportedArchitectures: []*string{
										aws.String("x86_64"),
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

			id, err := s.DiscoverLaunchTemplateAMI(ms)
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
				m.DescribeInstanceTypes(gomock.Any()).
					Return(&ec2.DescribeInstanceTypesOutput{
						InstanceTypes: []*ec2.InstanceTypeInfo{
							{
								ProcessorInfo: &ec2.ProcessorInfo{
									SupportedArchitectures: []*string{
										aws.String("x86_64"),
									},
								},
							},
						},
					}, nil)
			},
			expectSSM: func(m *mock_ssmiface.MockSSMAPIMockRecorder) {
				m.GetParameter(gomock.AssignableToTypeOf(&ssm.GetParameterInput{})).
					Return(&ssm.GetParameterOutput{
						Parameter: &ssm.Parameter{
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

			id, err := s.DiscoverLaunchTemplateAMI(ms)
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
				m.DeleteLaunchTemplateVersions(gomock.Eq(
					&ec2.DeleteLaunchTemplateVersionsInput{
						LaunchTemplateId: aws.String("id"),
						Versions:         aws.StringSlice([]string{"12"}),
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
				m.DeleteLaunchTemplateVersions(gomock.Eq(
					&ec2.DeleteLaunchTemplateVersionsInput{
						LaunchTemplateId: aws.String("id"),
						Versions:         aws.StringSlice([]string{"12"}),
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
