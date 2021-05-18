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
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/ec2/mock_ec2iface"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/userdata"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
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
		expect             func(m *mock_ec2iface.MockEC2APIMockRecorder)
		check              func(launchtemplate *expinfrav1.AWSLaunchTemplate, userdatahash string, err error)
	}{
		{
			name:               "does not exist",
			launchTemplateName: "foo",
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
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
			check: func(launchtemplate *expinfrav1.AWSLaunchTemplate, userdatahash string, err error) {
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}

				if userdatahash != "" {
					t.Fatalf("Did not expect a userdata hash, but got something: %s", userdatahash)
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

			launchtemplate, userdatahash, err := s.GetLaunchTemplate(tc.launchTemplateName)
			tc.check(launchtemplate, userdatahash, err)
		})
	}
}

func TestService_SDKToLaunchTemplate(t *testing.T) {
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
				AMI: infrav1.AWSResourceReference{
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
			if !reflect.DeepEqual(gotLT, tt.wantLT) {
				t.Fatalf("launchtemplate mismatch: got %v, want %v", gotLT, tt.wantLT)
			}
			if !reflect.DeepEqual(gotHash, tt.wantHash) {
				t.Fatalf("userdatahash mismatch: got %v, want %v", gotHash, tt.wantHash)
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
