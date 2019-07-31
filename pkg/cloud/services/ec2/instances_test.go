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
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/infrastructure/v1alpha2"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/ec2/mock_ec2iface"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/elb/mock_elbiface"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha2"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestInstanceIfExists(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name       string
		instanceID string
		expect     func(m *mock_ec2iface.MockEC2APIMockRecorder)
		check      func(instance *v1alpha2.Instance, err error)
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
			check: func(instance *v1alpha2.Instance, err error) {
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
			check: func(instance *v1alpha2.Instance, err error) {
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
			check: func(i *v1alpha2.Instance, err error) {
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
				AWSCluster: &v1alpha2.AWSCluster{
					Spec: v1alpha2.AWSClusterSpec{
						NetworkSpec: v1alpha2.NetworkSpec{
							VPC: v1alpha2.VPCSpec{
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
				AWSCluster: &v1alpha2.AWSCluster{},
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
		machineConfig *v1alpha2.AWSMachineSpec
		awsCluster    *v1alpha2.AWSCluster
		expect        func(m *mock_ec2iface.MockEC2APIMockRecorder)
		check         func(instance *v1alpha2.Instance, err error)
	}{
		{
			name: "simple",
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						Data: pointer.StringPtr("user-data"),
					},
				},
			},
			machineConfig: &v1alpha2.AWSMachineSpec{
				AMI: v1alpha2.AWSResourceReference{
					ID: aws.String("abc"),
				},
				InstanceType: "m5.large",
			},
			awsCluster: &v1alpha2.AWSCluster{
				Spec: v1alpha2.AWSClusterSpec{
					NetworkSpec: v1alpha2.NetworkSpec{
						Subnets: v1alpha2.Subnets{
							&v1alpha2.SubnetSpec{
								ID:       "subnet-1",
								IsPublic: false,
							},
							&v1alpha2.SubnetSpec{
								IsPublic: false,
							},
						},
					},
				},
				Status: v1alpha2.AWSClusterStatus{
					Network: v1alpha2.Network{
						SecurityGroups: map[v1alpha2.SecurityGroupRole]v1alpha2.SecurityGroup{
							v1alpha2.SecurityGroupControlPlane: {
								ID: "1",
							},
							v1alpha2.SecurityGroupNode: {
								ID: "2",
							},
							v1alpha2.SecurityGroupLB: {
								ID: "3",
							},
						},
						APIServerELB: v1alpha2.ClassicELB{
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
			check: func(instance *v1alpha2.Instance, err error) {
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
						Data: pointer.StringPtr("user-data"),
					},
				},
			},
			machineConfig: &v1alpha2.AWSMachineSpec{
				AMI: v1alpha2.AWSResourceReference{
					ID: aws.String("abc"),
				},
				InstanceType:     "m5.2xlarge",
				AvailabilityZone: aws.String("us-east-1c"),
			},
			awsCluster: &v1alpha2.AWSCluster{
				Spec: v1alpha2.AWSClusterSpec{
					NetworkSpec: v1alpha2.NetworkSpec{
						Subnets: v1alpha2.Subnets{
							&v1alpha2.SubnetSpec{
								ID:               "subnet-1",
								AvailabilityZone: "us-east-1a",
								IsPublic:         false,
							},
							&v1alpha2.SubnetSpec{
								ID:               "subnet-2",
								AvailabilityZone: "us-east-1b",
								IsPublic:         false,
							},
							&v1alpha2.SubnetSpec{
								ID:               "subnet-3",
								AvailabilityZone: "us-east-1c",
								IsPublic:         false,
							},
							&v1alpha2.SubnetSpec{
								ID:               "subnet-3-public",
								AvailabilityZone: "us-east-1c",
								IsPublic:         true,
							},
						},
					},
				},
				Status: v1alpha2.AWSClusterStatus{
					Network: v1alpha2.Network{
						SecurityGroups: map[v1alpha2.SecurityGroupRole]v1alpha2.SecurityGroup{
							v1alpha2.SecurityGroupControlPlane: {
								ID: "1",
							},
							v1alpha2.SecurityGroupNode: {
								ID: "2",
							},
							v1alpha2.SecurityGroupLB: {
								ID: "3",
							},
						},
						APIServerELB: v1alpha2.ClassicELB{
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
			check: func(instance *v1alpha2.Instance, err error) {
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}

				if instance.SubnetID != "subnet-3" {
					t.Fatalf("expected subnet-3 from availability zone us-east-1c, got %q", instance.SubnetID)
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
					ClusterNetwork: &clusterv1.ClusterNetworkingConfig{
						ServiceDomain: "cluster.local",
						Services: clusterv1.NetworkRanges{
							CIDRBlocks: []string{"192.168.0.0/16"},
						},
						Pods: clusterv1.NetworkRanges{
							CIDRBlocks: []string{"192.168.0.0/16"},
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
						Data: pointer.StringPtr("user-data"),
					},
				},
			}

			awsMachine := &v1alpha2.AWSMachine{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: clusterv1.SchemeGroupVersion.String(),
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

			s := NewService(clusterScope)
			instance, err := s.CreateInstance(machineScope)
			tc.check(instance, err)
		})
	}
}
