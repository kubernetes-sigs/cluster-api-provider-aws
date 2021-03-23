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
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/ec2/mock_ec2iface"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
)

func TestGetLaunchTemplate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name               string
		launchTemplateName string
		expect             func(m *mock_ec2iface.MockEC2APIMockRecorder)
		check              func(launchtemplate *expinfrav1.AWSLaunchTemplate, err error)
	}{
		{
			name:               "does not exist",
			launchTemplateName: "foo",
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.DescribeLaunchTemplateVersions(gomock.Eq(&ec2.DescribeLaunchTemplateVersionsInput{
					LaunchTemplateName: aws.String("foo"),
					Versions:           []*string{aws.String("$Latest")},
				})).
					Return(nil, awserrors.NewNotFound("not found"))
			},
			check: func(launchtemplate *expinfrav1.AWSLaunchTemplate, err error) {
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}

				if launchtemplate != nil {
					t.Fatalf("Did not expect anything but got something: %+v", launchtemplate)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ec2Mock := mock_ec2iface.NewMockEC2API(mockCtrl)

			scope, err := scope.NewClusterScope(scope.ClusterScopeParams{
				Cluster:    &clusterv1.Cluster{},
				AWSCluster: &infrav1.AWSCluster{},
			})
			if err != nil {
				t.Fatalf("Failed to create test context: %v", err)
			}

			tc.expect(ec2Mock.EXPECT())

			s := NewService(scope)
			s.EC2Client = ec2Mock

			launchtemplate, err := s.GetLaunchTemplate(tc.launchTemplateName)
			tc.check(launchtemplate, err)
		})
	}
}

func TestService_SDKToLaunchTemplate(t *testing.T) {
	tests := []struct {
		name    string
		input   *ec2.LaunchTemplateVersion
		want    *expinfrav1.AWSLaunchTemplate
		wantErr bool
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
				},
				VersionNumber: aws.Int64(1),
			},
			want: &expinfrav1.AWSLaunchTemplate{
				ID:   "lt-12345",
				Name: "foo",
				AMI: infrav1.AWSResourceReference{
					ID: aws.String("foo-image"),
				},
				IamInstanceProfile: "foo-profile",
				SSHKeyName:         aws.String("foo-keyname"),
				VersionNumber:      aws.Int64(1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{}
			got, err := s.SDKToLaunchTemplate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error mismatch: got %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("launchtemplate mismatch: got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_LaunchTemplateNeedsUpdate(t *testing.T) {
	tests := []struct {
		name     string
		incoming *expinfrav1.AWSLaunchTemplate
		existing *expinfrav1.AWSLaunchTemplate
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac := &infrav1.AWSCluster{
				Status: infrav1.AWSClusterStatus{
					Network: infrav1.Network{
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
			got, err := s.LaunchTemplateNeedsUpdate(machinePoolScope, tt.incoming, tt.existing)
			if (err != nil) != tt.wantErr {
				t.Errorf("error mismatch: got = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("result mismatch: got = %v, want %v", got, tt.want)
			}
		})
	}
}
