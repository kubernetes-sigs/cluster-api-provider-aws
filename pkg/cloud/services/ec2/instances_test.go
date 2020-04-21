/*
Copyright 2018 The Kubernetes Authors.

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
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/ec2/mock_ec2iface"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/elb/mock_elbiface"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestInstanceIfExists(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name       string
		instanceID string
		expect     func(m *mock_ec2iface.MockEC2APIMockRecorder)
		check      func(instance *infrav1.Instance, err error)
	}{
		{
			name:       "does not exist",
			instanceID: "hello",
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.DescribeInstances(gomock.Eq(&ec2.DescribeInstancesInput{
					InstanceIds: []*string{aws.String("hello")},
				})).
					Return(nil, awserrors.NewNotFound(errors.New("not found")))
			},
			check: func(instance *infrav1.Instance, err error) {
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}

				if instance != nil {
					t.Fatalf("Did not expect anything but got something: %+v", instance)
				}
			},
		},
		{
			name:       "does not exist with bad request error",
			instanceID: "hello-does-not-exist",
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.DescribeInstances(gomock.Eq(&ec2.DescribeInstancesInput{
					InstanceIds: []*string{aws.String("hello-does-not-exist")},
				})).
					Return(nil, awserr.New(awserrors.InvalidInstanceID, "does not exist", nil))
			},
			check: func(instance *infrav1.Instance, err error) {
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}

				if instance != nil {
					t.Fatalf("Did not expect anything but got something: %+v", instance)
				}
			},
		},
		{
			name:       "instance exists",
			instanceID: "id-1",
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.DescribeInstances(gomock.Eq(&ec2.DescribeInstancesInput{
					InstanceIds: []*string{aws.String("id-1")},
				})).
					Return(&ec2.DescribeInstancesOutput{
						Reservations: []*ec2.Reservation{
							{
								Instances: []*ec2.Instance{
									{
										InstanceId:   aws.String("id-1"),
										InstanceType: aws.String("m5.large"),
										SubnetId:     aws.String("subnet-1"),
										ImageId:      aws.String("ami-1"),
										IamInstanceProfile: &ec2.IamInstanceProfile{
											Arn: aws.String("arn:aws:iam::123456789012:instance-profile/foo"),
										},
										State: &ec2.InstanceState{
											Code: aws.Int64(16),
											Name: aws.String(ec2.StateAvailable),
										},
										RootDeviceName: aws.String("device-1"),
										BlockDeviceMappings: []*ec2.InstanceBlockDeviceMapping{
											{
												DeviceName: aws.String("device-1"),
												Ebs: &ec2.EbsInstanceBlockDevice{
													VolumeId: aws.String("volume-1"),
												},
											},
										},
									},
								},
							},
						},
					}, nil)

			},
			check: func(instance *infrav1.Instance, err error) {
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}

				if instance == nil {
					t.Fatalf("expected instance but got nothing")
				}

				if instance.ID != "id-1" {
					t.Fatalf("expected id-1 but got: %v", instance.ID)
				}
			},
		},
		{
			name:       "error describing instances",
			instanceID: "one",
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.DescribeInstances(&ec2.DescribeInstancesInput{
					InstanceIds: []*string{aws.String("one")},
				}).
					Return(nil, errors.New("some unknown error"))
			},
			check: func(i *infrav1.Instance, err error) {
				if err == nil {
					t.Fatalf("expected an error but got none.")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ec2Mock := mock_ec2iface.NewMockEC2API(mockCtrl)
			elbMock := mock_elbiface.NewMockELBAPI(mockCtrl)

			scope, err := scope.NewClusterScope(scope.ClusterScopeParams{
				Cluster: &clusterv1.Cluster{},
				AWSClients: scope.AWSClients{
					EC2: ec2Mock,
					ELB: elbMock,
				},
				AWSCluster: &infrav1.AWSCluster{
					Spec: infrav1.AWSClusterSpec{
						NetworkSpec: infrav1.NetworkSpec{
							VPC: infrav1.VPCSpec{
								ID: "test-vpc",
							},
						},
					},
				},
			})
			if err != nil {
				t.Fatalf("Failed to create test context: %v", err)
			}

			tc.expect(ec2Mock.EXPECT())

			s := NewService(scope)
			instance, err := s.InstanceIfExists(&tc.instanceID)
			tc.check(instance, err)
		})
	}
}

func TestTerminateInstance(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	instanceNotFoundError := errors.New("instance not found")

	testCases := []struct {
		name       string
		instanceID string
		expect     func(m *mock_ec2iface.MockEC2APIMockRecorder)
		check      func(err error)
	}{
		{
			name:       "instance exists",
			instanceID: "i-exist",
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.TerminateInstances(gomock.Eq(&ec2.TerminateInstancesInput{
					InstanceIds: []*string{aws.String("i-exist")},
				})).
					Return(&ec2.TerminateInstancesOutput{}, nil)
			},
			check: func(err error) {
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}
			},
		},
		{
			name:       "instance does not exist",
			instanceID: "i-donotexist",
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.TerminateInstances(gomock.Eq(&ec2.TerminateInstancesInput{
					InstanceIds: []*string{aws.String("i-donotexist")},
				})).
					Return(&ec2.TerminateInstancesOutput{}, instanceNotFoundError)
			},
			check: func(err error) {
				if err == nil {
					t.Fatalf("did not expect error: %v", err)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ec2Mock := mock_ec2iface.NewMockEC2API(mockCtrl)
			elbMock := mock_elbiface.NewMockELBAPI(mockCtrl)

			scope, err := scope.NewClusterScope(scope.ClusterScopeParams{
				AWSClients: scope.AWSClients{
					EC2: ec2Mock,
					ELB: elbMock,
				},
				Cluster:    &clusterv1.Cluster{},
				AWSCluster: &infrav1.AWSCluster{},
			})

			if err != nil {
				t.Fatalf("Failed to create test context: %v", err)
			}

			tc.expect(ec2Mock.EXPECT())

			s := NewService(scope)
			err = s.TerminateInstance(tc.instanceID)
			tc.check(err)
		})
	}
}

func TestCreateInstance(t *testing.T) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "bootstrap-data",
		},
		Data: map[string][]byte{
			"value": []byte("data"),
		},
	}

	testcases := []struct {
		name          string
		machine       clusterv1.Machine
		machineConfig *infrav1.AWSMachineSpec
		awsCluster    *infrav1.AWSCluster
		expect        func(m *mock_ec2iface.MockEC2APIMockRecorder)
		check         func(instance *infrav1.Instance, err error)
	}{
		{
			name: "simple",
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: pointer.StringPtr("bootstrap-data"),
					},
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AWSResourceReference{
					ID: aws.String("abc"),
				},
				InstanceType: "m5.large",
			},
			awsCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: infrav1.Subnets{
							&infrav1.SubnetSpec{
								ID:       "subnet-1",
								IsPublic: false,
							},
							&infrav1.SubnetSpec{
								IsPublic: false,
							},
						},
					},
				},
				Status: infrav1.AWSClusterStatus{
					Network: infrav1.Network{
						SecurityGroups: map[infrav1.SecurityGroupRole]infrav1.SecurityGroup{
							infrav1.SecurityGroupControlPlane: {
								ID: "1",
							},
							infrav1.SecurityGroupNode: {
								ID: "2",
							},
							infrav1.SecurityGroupLB: {
								ID: "3",
							},
						},
						APIServerELB: infrav1.ClassicELB{
							DNSName: "test-apiserver.us-east-1.aws",
						},
					},
				},
			},
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.
					DescribeImages(gomock.Any()).
					Return(&ec2.DescribeImagesOutput{
						Images: []*ec2.Image{
							{
								Name: aws.String("ami-1"),
							},
						},
					}, nil)
				m. // TODO: Restore these parameters, but with the tags as well
					RunInstances(gomock.Any()).
					Return(&ec2.Reservation{
						Instances: []*ec2.Instance{
							{
								State: &ec2.InstanceState{
									Name: aws.String(ec2.InstanceStateNamePending),
								},
								IamInstanceProfile: &ec2.IamInstanceProfile{
									Arn: aws.String("arn:aws:iam::123456789012:instance-profile/foo"),
								},
								InstanceId:     aws.String("two"),
								InstanceType:   aws.String("m5.large"),
								SubnetId:       aws.String("subnet-1"),
								ImageId:        aws.String("ami-1"),
								RootDeviceName: aws.String("device-1"),
								BlockDeviceMappings: []*ec2.InstanceBlockDeviceMapping{
									{
										DeviceName: aws.String("device-1"),
										Ebs: &ec2.EbsInstanceBlockDevice{
											VolumeId: aws.String("volume-1"),
										},
									},
								},
							},
						},
					}, nil)
				m.WaitUntilInstanceRunningWithContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)

			},
			check: func(instance *infrav1.Instance, err error) {
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}
			},
		},
		{
			name: "with availability zone",
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: pointer.StringPtr("bootstrap-data"),
					},
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AWSResourceReference{
					ID: aws.String("abc"),
				},
				InstanceType:  "m5.2xlarge",
				FailureDomain: aws.String("us-east-1c"),
			},
			awsCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: infrav1.Subnets{
							&infrav1.SubnetSpec{
								ID:               "subnet-1",
								AvailabilityZone: "us-east-1a",
								IsPublic:         false,
							},
							&infrav1.SubnetSpec{
								ID:               "subnet-2",
								AvailabilityZone: "us-east-1b",
								IsPublic:         false,
							},
							&infrav1.SubnetSpec{
								ID:               "subnet-3",
								AvailabilityZone: "us-east-1c",
								IsPublic:         false,
							},
							&infrav1.SubnetSpec{
								ID:               "subnet-3-public",
								AvailabilityZone: "us-east-1c",
								IsPublic:         true,
							},
						},
					},
				},
				Status: infrav1.AWSClusterStatus{
					Network: infrav1.Network{
						SecurityGroups: map[infrav1.SecurityGroupRole]infrav1.SecurityGroup{
							infrav1.SecurityGroupControlPlane: {
								ID: "1",
							},
							infrav1.SecurityGroupNode: {
								ID: "2",
							},
							infrav1.SecurityGroupLB: {
								ID: "3",
							},
						},
						APIServerELB: infrav1.ClassicELB{
							DNSName: "test-apiserver.us-east-1.aws",
						},
					},
				},
			},
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.
					DescribeImages(gomock.Any()).
					Return(&ec2.DescribeImagesOutput{
						Images: []*ec2.Image{
							{
								Name: aws.String("ami-1"),
							},
						},
					}, nil)

				m.
					RunInstances(gomock.Any()).
					Return(&ec2.Reservation{
						Instances: []*ec2.Instance{
							{
								State: &ec2.InstanceState{
									Name: aws.String(ec2.InstanceStateNamePending),
								},
								IamInstanceProfile: &ec2.IamInstanceProfile{
									Arn: aws.String("arn:aws:iam::123456789012:instance-profile/foo"),
								},
								InstanceId:     aws.String("two"),
								InstanceType:   aws.String("m5.large"),
								SubnetId:       aws.String("subnet-3"),
								ImageId:        aws.String("ami-1"),
								RootDeviceName: aws.String("device-1"),
								BlockDeviceMappings: []*ec2.InstanceBlockDeviceMapping{
									{
										DeviceName: aws.String("device-1"),
										Ebs: &ec2.EbsInstanceBlockDevice{
											VolumeId: aws.String("volume-1"),
										},
									},
								},
							},
						},
					}, nil)

				m.WaitUntilInstanceRunningWithContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)

			},
			check: func(instance *infrav1.Instance, err error) {
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}

				if instance.SubnetID != "subnet-3" {
					t.Fatalf("expected subnet-3 from availability zone us-east-1c, got %q", instance.SubnetID)
				}
			},
		},
		{
			name: "with ImageLookupOrg specified at the machine level",
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: pointer.StringPtr("bootstrap-data"),
					},
					Version: pointer.StringPtr("v1.16.1"),
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				ImageLookupOrg: "test-org-123",
				InstanceType:   "m5.large",
			},
			awsCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: infrav1.Subnets{
							&infrav1.SubnetSpec{
								ID:       "subnet-1",
								IsPublic: false,
							},
							&infrav1.SubnetSpec{
								IsPublic: false,
							},
						},
					},
				},
				Status: infrav1.AWSClusterStatus{
					Network: infrav1.Network{
						SecurityGroups: map[infrav1.SecurityGroupRole]infrav1.SecurityGroup{
							infrav1.SecurityGroupControlPlane: {
								ID: "1",
							},
							infrav1.SecurityGroupNode: {
								ID: "2",
							},
							infrav1.SecurityGroupLB: {
								ID: "3",
							},
						},
						APIServerELB: infrav1.ClassicELB{
							DNSName: "test-apiserver.us-east-1.aws",
						},
					},
				},
			},
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				// verify that the ImageLookupOrg is used when finding AMIs
				m.
					DescribeImages(gomock.Eq(&ec2.DescribeImagesInput{
						Filters: []*ec2.Filter{
							{
								Name:   aws.String("owner-id"),
								Values: []*string{aws.String("test-org-123")},
							},
							{
								Name:   aws.String("name"),
								Values: []*string{aws.String(amiName("capa-ami-${BASE_OS}-?${K8S_VERSION}-*","ubuntu-18.04", "v1.16.1"))},
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
					})).
					Return(&ec2.DescribeImagesOutput{
						Images: []*ec2.Image{
							{
								Name:         aws.String("ami-1"),
								CreationDate: aws.String("2006-01-02T15:04:05.000Z"),
							},
						},
					}, nil)
				m. // TODO: Restore these parameters, but with the tags as well
					RunInstances(gomock.Any()).
					Return(&ec2.Reservation{
						Instances: []*ec2.Instance{
							{
								State: &ec2.InstanceState{
									Name: aws.String(ec2.InstanceStateNamePending),
								},
								IamInstanceProfile: &ec2.IamInstanceProfile{
									Arn: aws.String("arn:aws:iam::123456789012:instance-profile/foo"),
								},
								InstanceId:     aws.String("two"),
								InstanceType:   aws.String("m5.large"),
								SubnetId:       aws.String("subnet-1"),
								ImageId:        aws.String("ami-1"),
								RootDeviceName: aws.String("device-1"),
								BlockDeviceMappings: []*ec2.InstanceBlockDeviceMapping{
									{
										DeviceName: aws.String("device-1"),
										Ebs: &ec2.EbsInstanceBlockDevice{
											VolumeId: aws.String("volume-1"),
										},
									},
								},
							},
						},
					}, nil)

				m.WaitUntilInstanceRunningWithContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)

			},
			check: func(instance *infrav1.Instance, err error) {
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}
			},
		},
		{
			name: "with ImageLookupOrg specified at the cluster-level",
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: pointer.StringPtr("bootstrap-data"),
					},
					Version: pointer.StringPtr("v1.16.1"),
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				InstanceType: "m5.large",
			},
			awsCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: infrav1.Subnets{
							&infrav1.SubnetSpec{
								ID:       "subnet-1",
								IsPublic: false,
							},
							&infrav1.SubnetSpec{
								IsPublic: false,
							},
						},
					},
					ImageLookupOrg: "cluster-level-image-lookup-org",
				},
				Status: infrav1.AWSClusterStatus{
					Network: infrav1.Network{
						SecurityGroups: map[infrav1.SecurityGroupRole]infrav1.SecurityGroup{
							infrav1.SecurityGroupControlPlane: {
								ID: "1",
							},
							infrav1.SecurityGroupNode: {
								ID: "2",
							},
							infrav1.SecurityGroupLB: {
								ID: "3",
							},
						},
						APIServerELB: infrav1.ClassicELB{
							DNSName: "test-apiserver.us-east-1.aws",
						},
					},
				},
			},
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				// verify that the ImageLookupOrg is used when finding AMIs
				m.
					DescribeImages(gomock.Eq(&ec2.DescribeImagesInput{
						Filters: []*ec2.Filter{
							{
								Name:   aws.String("owner-id"),
								Values: []*string{aws.String("cluster-level-image-lookup-org")},
							},
							{
								Name:   aws.String("name"),
								Values: []*string{aws.String(amiName("capa-ami-${BASE_OS}-?${K8S_VERSION}-*","ubuntu-18.04", "v1.16.1"))},
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
					})).
					Return(&ec2.DescribeImagesOutput{
						Images: []*ec2.Image{
							{
								Name:         aws.String("ami-1"),
								CreationDate: aws.String("2006-01-02T15:04:05.000Z"),
							},
						},
					}, nil)
				m. // TODO: Restore these parameters, but with the tags as well
					RunInstances(gomock.Any()).
					Return(&ec2.Reservation{
						Instances: []*ec2.Instance{
							{
								State: &ec2.InstanceState{
									Name: aws.String(ec2.InstanceStateNamePending),
								},
								IamInstanceProfile: &ec2.IamInstanceProfile{
									Arn: aws.String("arn:aws:iam::123456789012:instance-profile/foo"),
								},
								InstanceId:     aws.String("two"),
								InstanceType:   aws.String("m5.large"),
								SubnetId:       aws.String("subnet-1"),
								ImageId:        aws.String("ami-1"),
								RootDeviceName: aws.String("device-1"),
								BlockDeviceMappings: []*ec2.InstanceBlockDeviceMapping{
									{
										DeviceName: aws.String("device-1"),
										Ebs: &ec2.EbsInstanceBlockDevice{
											VolumeId: aws.String("volume-1"),
										},
									},
								},
							},
						},
					}, nil)

				m.WaitUntilInstanceRunningWithContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)

			},
			check: func(instance *infrav1.Instance, err error) {
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}
			},
		},
		{
			name: "AWSMachine ImageLookupOrg overrides AWSCluster ImageLookupOrg",
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: pointer.StringPtr("bootstrap-data"),
					},
					Version: pointer.StringPtr("v1.16.1"),
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				InstanceType:   "m5.large",
				ImageLookupOrg: "machine-level-image-lookup-org",
			},
			awsCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: infrav1.Subnets{
							&infrav1.SubnetSpec{
								ID:       "subnet-1",
								IsPublic: false,
							},
							&infrav1.SubnetSpec{
								IsPublic: false,
							},
						},
					},
					ImageLookupOrg: "cluster-level-image-lookup-org",
				},
				Status: infrav1.AWSClusterStatus{
					Network: infrav1.Network{
						SecurityGroups: map[infrav1.SecurityGroupRole]infrav1.SecurityGroup{
							infrav1.SecurityGroupControlPlane: {
								ID: "1",
							},
							infrav1.SecurityGroupNode: {
								ID: "2",
							},
							infrav1.SecurityGroupLB: {
								ID: "3",
							},
						},
						APIServerELB: infrav1.ClassicELB{
							DNSName: "test-apiserver.us-east-1.aws",
						},
					},
				},
			},
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				// verify that the ImageLookupOrg is used when finding AMIs
				m.
					DescribeImages(gomock.Eq(&ec2.DescribeImagesInput{
						Filters: []*ec2.Filter{
							{
								Name:   aws.String("owner-id"),
								Values: []*string{aws.String("machine-level-image-lookup-org")},
							},
							{
								Name:   aws.String("name"),
								Values: []*string{aws.String(amiName("capa-ami-${BASE_OS}-?${K8S_VERSION}-*","ubuntu-18.04", "v1.16.1"))},
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
					})).
					Return(&ec2.DescribeImagesOutput{
						Images: []*ec2.Image{
							{
								Name:         aws.String("ami-1"),
								CreationDate: aws.String("2006-01-02T15:04:05.000Z"),
							},
						},
					}, nil)
				m. // TODO: Restore these parameters, but with the tags as well
					RunInstances(gomock.Any()).
					Return(&ec2.Reservation{
						Instances: []*ec2.Instance{
							{
								State: &ec2.InstanceState{
									Name: aws.String(ec2.InstanceStateNamePending),
								},
								IamInstanceProfile: &ec2.IamInstanceProfile{
									Arn: aws.String("arn:aws:iam::123456789012:instance-profile/foo"),
								},
								InstanceId:     aws.String("two"),
								InstanceType:   aws.String("m5.large"),
								SubnetId:       aws.String("subnet-1"),
								ImageId:        aws.String("ami-1"),
								RootDeviceName: aws.String("device-1"),
								BlockDeviceMappings: []*ec2.InstanceBlockDeviceMapping{
									{
										DeviceName: aws.String("device-1"),
										Ebs: &ec2.EbsInstanceBlockDevice{
											VolumeId: aws.String("volume-1"),
										},
									},
								},
							},
						},
					}, nil)

				m.WaitUntilInstanceRunningWithContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)

			},
			check: func(instance *infrav1.Instance, err error) {
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			ec2Mock := mock_ec2iface.NewMockEC2API(mockCtrl)
			elbMock := mock_elbiface.NewMockELBAPI(mockCtrl)

			cluster := &clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test1",
				},
				Spec: clusterv1.ClusterSpec{
					ClusterNetwork: &clusterv1.ClusterNetwork{
						ServiceDomain: "cluster.local",
						Services: &clusterv1.NetworkRanges{
							CIDRBlocks: []string{"192.168.0.0/16"},
						},
						Pods: &clusterv1.NetworkRanges{
							CIDRBlocks: []string{"192.168.0.0/16"},
						},
					},
				},
			}

			awsCluster := tc.awsCluster

			machine := &tc.machine

			awsMachine := &infrav1.AWSMachine{
				ObjectMeta: metav1.ObjectMeta{
					Name: "aws-test1",
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: clusterv1.GroupVersion.String(),
							Kind:       "Machine",
							Name:       "test1",
						},
					},
				},
			}

			client := fake.NewFakeClient(secret, cluster, machine)

			machineScope, err := scope.NewMachineScope(scope.MachineScopeParams{
				Client: client,
				AWSClients: scope.AWSClients{
					EC2: ec2Mock,
					ELB: elbMock,
				},
				Cluster:    cluster,
				Machine:    machine,
				AWSMachine: awsMachine,
				AWSCluster: awsCluster,
			})
			if err != nil {
				t.Fatalf("Failed to create test context: %v", err)
			}
			machineScope.AWSMachine.Spec = *tc.machineConfig
			tc.expect(ec2Mock.EXPECT())

			clusterScope, err := scope.NewClusterScope(scope.ClusterScopeParams{
				Client: client,
				AWSClients: scope.AWSClients{
					EC2: ec2Mock,
					ELB: elbMock,
				},
				Cluster:    cluster,
				AWSCluster: tc.awsCluster,
			})
			if err != nil {
				t.Fatalf("Failed to create test context: %v", err)
			}

			s := NewService(clusterScope)
			instance, err := s.CreateInstance(machineScope, []byte("userData"))
			tc.check(instance, err)
		})
	}
}
