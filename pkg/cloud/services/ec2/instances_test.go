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
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha2"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/ec2/mock_ec2iface"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/elb/mock_elbiface"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha2"
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
						// echo "user-data" | base64
						Data: pointer.StringPtr("dXNlci1kYXRhCg=="),
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
								InstanceId:   aws.String("two"),
								InstanceType: aws.String("m5.large"),
								SubnetId:     aws.String("subnet-1"),
								ImageId:      aws.String("ami-1"),
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
						// echo "user-data" | base64
						Data: pointer.StringPtr("dXNlci1kYXRhCg=="),
					},
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AWSResourceReference{
					ID: aws.String("abc"),
				},
				InstanceType:     "m5.2xlarge",
				AvailabilityZone: aws.String("us-east-1c"),
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
								InstanceId:   aws.String("two"),
								InstanceType: aws.String("m5.large"),
								SubnetId:     aws.String("subnet-3"),
								ImageId:      aws.String("ami-1"),
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
			name: "with machine SSH Key Pair",
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						// echo "user-data" | base64
						Data: pointer.StringPtr("dXNlci1kYXRhCg=="),
					},
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AWSResourceReference{
					ID: aws.String("abc"),
				},
				InstanceType:     "m5.2xlarge",
				AvailabilityZone: aws.String("us-east-1c"),
				SSHKeyName:       aws.String("machine-key-pair"),
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
								InstanceId:   aws.String("two"),
								InstanceType: aws.String("m5.large"),
								SubnetId:     aws.String("subnet-3"),
								ImageId:      aws.String("ami-1"),
								KeyName: aws.String("machine-key-pair"),
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

				if instance.SSHKeyName == nil {
					t.Fatal("expected ssh key pair name machine-key-pair, got nil")
				}
				if *instance.SSHKeyName != "machine-key-pair" {
					t.Fatalf("expected ssh key pair name machine-key-pair, got %q", *instance.SSHKeyName)
				}
			},
		},
		{
			name: "with cluster SSH Key Pair",
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						// echo "user-data" | base64
						Data: pointer.StringPtr("dXNlci1kYXRhCg=="),
					},
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AWSResourceReference{
					ID: aws.String("abc"),
				},
				InstanceType:     "m5.2xlarge",
				AvailabilityZone: aws.String("us-east-1c"),
			},
			awsCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					SSHKeyName: aws.String("cluster-key-pair"),
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
								InstanceId:   aws.String("two"),
								InstanceType: aws.String("m5.large"),
								SubnetId:     aws.String("subnet-3"),
								ImageId:      aws.String("ami-1"),
								KeyName: aws.String("cluster-key-pair"),
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

				if instance.SSHKeyName == nil {
					t.Fatal("expected ssh key pair name cluster-key-pair, got nil")
				}

				if *instance.SSHKeyName != "cluster-key-pair" {
					t.Fatalf("expected ssh key pair name cluster-key-pair, got %q", *instance.SSHKeyName)
				}
			},
		},
		{
			name: "with NIL machine SSH Key Pair",
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						// echo "user-data" | base64
						Data: pointer.StringPtr("dXNlci1kYXRhCg=="),
					},
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AWSResourceReference{
					ID: aws.String("abc"),
				},
				InstanceType:     "m5.2xlarge",
				AvailabilityZone: aws.String("us-east-1c"),
				SSHKeyName:       nil,
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
								InstanceId:   aws.String("two"),
								InstanceType: aws.String("m5.large"),
								SubnetId:     aws.String("subnet-3"),
								ImageId:      aws.String("ami-1"),
								KeyName: aws.String(defaultSSHKeyName),
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

				if instance.SSHKeyName == nil {
					t.Fatalf("expected ssh key pair name %s, got nil",defaultSSHKeyName)
				}

				if *instance.SSHKeyName != defaultSSHKeyName {
					t.Fatalf("expected ssh key pair name %s, got %q", defaultSSHKeyName, *instance.SSHKeyName)
				}
			},
		},
		{
			name: "with NIL cluster SSH Key Pair",
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						// echo "user-data" | base64
						Data: pointer.StringPtr("dXNlci1kYXRhCg=="),
					},
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AWSResourceReference{
					ID: aws.String("abc"),
				},
				InstanceType:     "m5.2xlarge",
				AvailabilityZone: aws.String("us-east-1c"),
			},
			awsCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					SSHKeyName: nil,
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
								InstanceId:   aws.String("two"),
								InstanceType: aws.String("m5.large"),
								SubnetId:     aws.String("subnet-3"),
								ImageId:      aws.String("ami-1"),
								KeyName: aws.String(defaultSSHKeyName),
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

				if instance.SSHKeyName == nil {
					t.Fatalf("expected ssh key pair name %s, got nil", defaultSSHKeyName)
				}

				if *instance.SSHKeyName != defaultSSHKeyName {
					t.Fatalf("expected ssh key pair name %s, got %q", defaultSSHKeyName, *instance.SSHKeyName)
				}
			},
		},
		{
			name: "with empty string machine SSH Key Pair",
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						// echo "user-data" | base64
						Data: pointer.StringPtr("dXNlci1kYXRhCg=="),
					},
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AWSResourceReference{
					ID: aws.String("abc"),
				},
				InstanceType:     "m5.2xlarge",
				AvailabilityZone: aws.String("us-east-1c"),
				SSHKeyName:       aws.String(""),
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
								InstanceId:   aws.String("two"),
								InstanceType: aws.String("m5.large"),
								SubnetId:     aws.String("subnet-3"),
								ImageId:      aws.String("ami-1"),
								KeyName: aws.String(""),
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

				if instance.SSHKeyName == nil {
					t.Fatal("expected ssh key pair name \"\", got nil")
				}

				if *instance.SSHKeyName != "" {
					t.Fatalf("expected ssh key pair name \"\", got %q", *instance.SSHKeyName)
				}
			},
		},
		{
			name: "with empty string cluster SSH Key Pair",
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						// echo "user-data" | base64
						Data: pointer.StringPtr("dXNlci1kYXRhCg=="),
					},
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AWSResourceReference{
					ID: aws.String("abc"),
				},
				InstanceType:     "m5.2xlarge",
				AvailabilityZone: aws.String("us-east-1c"),
			},
			awsCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					SSHKeyName: aws.String(""),
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
								InstanceId:   aws.String("two"),
								InstanceType: aws.String("m5.large"),
								SubnetId:     aws.String("subnet-3"),
								ImageId:      aws.String("ami-1"),
								KeyName: aws.String(""),
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

				if instance.SSHKeyName == nil {
					t.Fatal("expected ssh key pair name \"\", got nil")
				}

				if *instance.SSHKeyName != "" {
					t.Fatalf("expected ssh key pair name \"\", got %q", *instance.SSHKeyName)
				}
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			// defer mockCtrl.Finish()
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

			awsCluster := &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: clusterv1.GroupVersion.String(),
							Kind:       "Cluster",
							Name:       "test1",
						},
					},
				},
			}

			machine := &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test1",
					Labels: map[string]string{
						"set":                             "node",
						clusterv1.MachineClusterLabelName: "test1",
					},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						// echo "user-data" | base64
						Data: pointer.StringPtr("dXNlci1kYXRhCg=="),
					},
				},
			}

			awsMachine := &infrav1.AWSMachine{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: clusterv1.GroupVersion.String(),
							Kind:       "Machine",
							Name:       "test1",
						},
					},
				},
			}

			machineScope, err := scope.NewMachineScope(scope.MachineScopeParams{
				Client: fake.NewFakeClient(cluster, machine),
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
				Client: fake.NewFakeClient(cluster, machine),
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
			instance, err := s.CreateInstance(machineScope)
			tc.check(instance, err)
		})
	}
}
