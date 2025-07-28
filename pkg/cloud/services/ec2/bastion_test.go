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
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/filter"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func TestServiceDeleteBastion(t *testing.T) {
	clusterName := "cluster"

	describeInput := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			filter.EC2.ProviderRole(infrav1.BastionRoleTagValue),
			filter.EC2.Cluster(clusterName),
			filter.EC2.InstanceStates(
				types.InstanceStateNamePending,
				types.InstanceStateNameRunning,
				types.InstanceStateNameStopping,
				types.InstanceStateNameStopped,
			),
		},
	}

	foundOutput := &ec2.DescribeInstancesOutput{
		Reservations: []types.Reservation{
			{
				Instances: []types.Instance{
					{
						InstanceId: aws.String("id123"),
						State: &types.InstanceState{
							Name: types.InstanceStateNameRunning,
						},
						Placement: &types.Placement{
							AvailabilityZone: aws.String("us-east-1"),
						},
					},
				},
			},
		},
	}

	tests := []struct {
		name          string
		expect        func(m *mocks.MockEC2APIMockRecorder)
		expectError   bool
		bastionStatus *infrav1.Instance
	}{
		{
			name: "instance not found",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.
					DescribeInstances(context.TODO(), gomock.Eq(describeInput)).
					Return(&ec2.DescribeInstancesOutput{}, nil)
			},
			expectError: false,
		},
		{
			name: "describe error",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.
					DescribeInstances(context.TODO(), gomock.Eq(describeInput)).
					Return(nil, errors.New("some error"))
			},
			expectError: true,
		},
		{
			name: "terminate fails",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.
					DescribeInstances(context.TODO(), gomock.Eq(describeInput)).
					Return(foundOutput, nil)
				m.
					TerminateInstances(context.TODO(),
						gomock.Eq(&ec2.TerminateInstancesInput{
							InstanceIds: ([]string{"id123"}),
						}),
					).
					Return(nil, errors.New("some error"))
			},
			expectError: true,
		},
		{
			name: "wait after terminate fails",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.
					DescribeInstances(context.TODO(), gomock.Eq(describeInput)).
					Return(foundOutput, nil).Times(1)
				m.
					TerminateInstances(context.TODO(),
						gomock.Eq(&ec2.TerminateInstancesInput{
							InstanceIds: []string{"id123"},
						}),
					).
					Return(&ec2.TerminateInstancesOutput{
						TerminatingInstances: []types.InstanceStateChange{
							{
								InstanceId: aws.String("id123"),
								CurrentState: &types.InstanceState{
									Name: types.InstanceStateNameShuttingDown,
								},
								PreviousState: &types.InstanceState{
									Name: types.InstanceStateNameRunning,
								},
							},
						},
					}, nil)
				// Simulate an error when waiting for the instance to be terminated by returning a state of "pending" instead of "terminated".
				m.DescribeInstances(gomock.Any(), gomock.Eq(&ec2.DescribeInstancesInput{
					InstanceIds: []string{"id123"},
				}), gomock.Any()).Return(&ec2.DescribeInstancesOutput{
					Reservations: []types.Reservation{
						{
							Instances: []types.Instance{
								{
									InstanceId: aws.String("id123"),
									State: &types.InstanceState{
										Name: types.InstanceStateNamePending,
									},
								},
							},
						},
					},
				}, nil)
			},
			expectError: true,
		},
		{
			name: "success",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.
					DescribeInstances(context.TODO(), gomock.Eq(describeInput)).
					Return(foundOutput, nil).Times(1)
				m.
					TerminateInstances(context.TODO(),
						gomock.Eq(&ec2.TerminateInstancesInput{
							InstanceIds: []string{"id123"},
						}),
					).
					Return(&ec2.TerminateInstancesOutput{
						TerminatingInstances: []types.InstanceStateChange{
							{
								InstanceId: aws.String("id123"),
								CurrentState: &types.InstanceState{
									Name: types.InstanceStateNameShuttingDown,
								},
								PreviousState: &types.InstanceState{
									Name: types.InstanceStateNameRunning,
								},
							},
						},
					}, nil)
				m.DescribeInstances(gomock.Any(), gomock.Eq(&ec2.DescribeInstancesInput{
					InstanceIds: []string{"id123"},
				}), gomock.Any()).
					Return(&ec2.DescribeInstancesOutput{
						Reservations: []types.Reservation{
							{
								Instances: []types.Instance{
									{
										InstanceId: aws.String("id123"),
										State: &types.InstanceState{
											Name: types.InstanceStateNameTerminated,
										},
									},
								},
							},
						},
					}, nil)
			},
			expectError:   false,
			bastionStatus: nil,
		},
	}

	for _, tc := range tests {
		managedValues := []bool{false, true}
		for i := range managedValues {
			managed := managedValues[i]

			t.Run(fmt.Sprintf("managed=%t %s", managed, tc.name), func(t *testing.T) {
				g := NewWithT(t)

				mockControl := gomock.NewController(t)
				defer mockControl.Finish()

				ec2Mock := mocks.NewMockEC2API(mockControl)

				scheme, err := setupScheme()
				g.Expect(err).To(BeNil())

				awsCluster := &infrav1.AWSCluster{
					ObjectMeta: metav1.ObjectMeta{Name: "test"},
					Spec: infrav1.AWSClusterSpec{
						NetworkSpec: infrav1.NetworkSpec{
							VPC: infrav1.VPCSpec{
								ID: "vpcID",
							},
						},
					},
				}

				client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(awsCluster).WithStatusSubresource(awsCluster).Build()

				scope, err := scope.NewClusterScope(scope.ClusterScopeParams{
					Cluster: &clusterv1.Cluster{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "ns",
							Name:      clusterName,
						},
					},
					AWSCluster: awsCluster,
					Client:     client,
				})
				g.Expect(err).To(BeNil())

				if managed {
					scope.AWSCluster.Spec.NetworkSpec.VPC.Tags = infrav1.Tags{
						infrav1.ClusterTagKey(clusterName): string(infrav1.ResourceLifecycleOwned),
					}
				}

				tc.expect(ec2Mock.EXPECT())
				s := NewService(scope)
				s.EC2Client = ec2Mock

				err = s.DeleteBastion()
				if tc.expectError {
					g.Expect(err).NotTo(BeNil())
					return
				}

				g.Expect(err).To(BeNil())

				g.Expect(scope.AWSCluster.Status.Bastion).To(BeEquivalentTo(tc.bastionStatus))
			})
		}
	}
}

func TestServiceReconcileBastion(t *testing.T) {
	clusterName := "cluster"

	describeInput := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			filter.EC2.ProviderRole(infrav1.BastionRoleTagValue),
			filter.EC2.Cluster(clusterName),
			filter.EC2.InstanceStates(
				types.InstanceStateNamePending,
				types.InstanceStateNameRunning,
				types.InstanceStateNameStopping,
				types.InstanceStateNameStopped,
			),
		},
	}

	foundOutput := &ec2.DescribeInstancesOutput{
		Reservations: []types.Reservation{
			{
				Instances: []types.Instance{
					{
						InstanceId: aws.String("id123"),
						State: &types.InstanceState{
							Name: types.InstanceStateNameRunning,
						},
						Placement: &types.Placement{
							AvailabilityZone: aws.String("us-east-1"),
						},
					},
				},
			},
		},
	}

	tests := []struct {
		name           string
		bastionEnabled bool
		expect         func(m *mocks.MockEC2APIMockRecorder)
		expectError    bool
		bastionStatus  *infrav1.Instance
	}{
		{
			name: "Should ignore reconciliation if instance not found",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.
					DescribeInstances(context.TODO(), gomock.Eq(describeInput)).
					Return(&ec2.DescribeInstancesOutput{}, nil)
			},
			expectError: false,
		},
		{
			name: "Should fail reconcile if describe instance fails",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.
					DescribeInstances(context.TODO(), gomock.Eq(describeInput)).
					Return(nil, errors.New("some error"))
			},
			expectError: true,
		},
		{
			name: "Should fail reconcile if terminate instance fails",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.
					DescribeInstances(context.TODO(), gomock.Eq(describeInput)).
					Return(foundOutput, nil).MinTimes(1)
				m.
					TerminateInstances(context.TODO(),
						gomock.Eq(&ec2.TerminateInstancesInput{
							InstanceIds: []string{"id123"},
						}),
					).
					Return(nil, errors.New("some error"))
			},
			expectError: true,
		},
		{
			name: "Should create bastion successfully",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeInstances(context.TODO(), gomock.Eq(describeInput)).
					Return(&ec2.DescribeInstancesOutput{}, nil).MinTimes(1)
				m.DescribeImages(context.TODO(), gomock.Eq(&ec2.DescribeImagesInput{Filters: []types.Filter{
					{
						Name:   aws.String("architecture"),
						Values: []string{"x86_64"},
					},
					{
						Name:   aws.String("state"),
						Values: []string{"available"},
					},
					{
						Name:   aws.String("virtualization-type"),
						Values: []string{"hvm"},
					},
					{
						Name:   aws.String("description"),
						Values: []string{ubuntuImageDescription},
					},
					{
						Name:   aws.String("owner-id"),
						Values: []string{ubuntuOwnerID},
					},
				}})).Return(&ec2.DescribeImagesOutput{Images: images{
					{
						ImageId:      aws.String("ubuntu-ami-id-latest"),
						CreationDate: aws.String("2019-02-08T17:02:31.000Z"),
					},
					{
						ImageId:      aws.String("ubuntu-ami-id-old"),
						CreationDate: aws.String("2014-02-08T17:02:31.000Z"),
					},
				}}, nil)
				m.RunInstances(context.TODO(), gomock.Any()).
					Return(&ec2.RunInstancesOutput{
						Instances: []types.Instance{
							{
								State: &types.InstanceState{
									Name: types.InstanceStateNameRunning,
								},
								IamInstanceProfile: &types.IamInstanceProfile{
									Arn: aws.String("arn:aws:iam::123456789012:instance-profile/foo"),
								},
								InstanceId:     aws.String("id123"),
								InstanceType:   types.InstanceTypeT3Micro,
								SubnetId:       aws.String("subnet-1"),
								ImageId:        aws.String("ubuntu-ami-id-latest"),
								RootDeviceName: aws.String("device-1"),
								BlockDeviceMappings: []types.InstanceBlockDeviceMapping{
									{
										DeviceName: aws.String("device-1"),
										Ebs: &types.EbsInstanceBlockDevice{
											VolumeId: aws.String("volume-1"),
										},
									},
								},
								Placement: &types.Placement{
									AvailabilityZone: aws.String("us-east-1"),
								},
							},
						},
					}, nil)
			},
			bastionEnabled: true,
			expectError:    false,
			bastionStatus: &infrav1.Instance{
				ID:               "id123",
				State:            "running",
				Type:             "t3.micro",
				SubnetID:         "subnet-1",
				ImageID:          "ubuntu-ami-id-latest",
				IAMProfile:       "foo",
				Addresses:        []clusterv1.MachineAddress{},
				AvailabilityZone: "us-east-1",
				VolumeIDs:        []string{"volume-1"},
			},
		},
	}

	for _, tc := range tests {
		managedValues := []bool{false, true}
		for i := range managedValues {
			managed := managedValues[i]

			t.Run(fmt.Sprintf("managed=%t %s", managed, tc.name), func(t *testing.T) {
				g := NewWithT(t)

				mockControl := gomock.NewController(t)
				defer mockControl.Finish()

				ec2Mock := mocks.NewMockEC2API(mockControl)

				scheme, err := setupScheme()
				g.Expect(err).To(BeNil())

				awsCluster := &infrav1.AWSCluster{
					ObjectMeta: metav1.ObjectMeta{Name: "test"},
					Spec: infrav1.AWSClusterSpec{
						NetworkSpec: infrav1.NetworkSpec{
							VPC: infrav1.VPCSpec{
								ID: "vpcID",
							},
							Subnets: infrav1.Subnets{
								{
									ID: "subnet-1",
								},
								{
									ID:       "subnet-2",
									IsPublic: true,
								},
							},
						},
						Bastion: infrav1.Bastion{Enabled: tc.bastionEnabled},
					},
				}

				client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(awsCluster).WithStatusSubresource(awsCluster).Build()

				scope, err := scope.NewClusterScope(scope.ClusterScopeParams{
					Cluster: &clusterv1.Cluster{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "ns",
							Name:      clusterName,
						},
					},
					AWSCluster: awsCluster,
					Client:     client,
				})
				g.Expect(err).To(BeNil())

				if managed {
					scope.AWSCluster.Spec.NetworkSpec.VPC.Tags = infrav1.Tags{
						infrav1.ClusterTagKey(clusterName): string(infrav1.ResourceLifecycleOwned),
					}
				}

				tc.expect(ec2Mock.EXPECT())
				s := NewService(scope)
				s.EC2Client = ec2Mock

				err = s.ReconcileBastion()
				if tc.expectError {
					g.Expect(err).NotTo(BeNil())
					return
				}

				g.Expect(err).To(BeNil())

				g.Expect(scope.AWSCluster.Status.Bastion).To(BeEquivalentTo(tc.bastionStatus))
			})
		}
	}
}

func TestServiceReconcileBastionUSGOV(t *testing.T) {
	clusterName := "cluster-us-gov"

	describeInput := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			filter.EC2.ProviderRole(infrav1.BastionRoleTagValue),
			filter.EC2.Cluster(clusterName),
			filter.EC2.InstanceStates(
				types.InstanceStateNamePending,
				types.InstanceStateNameRunning,
				types.InstanceStateNameStopping,
				types.InstanceStateNameStopped,
			),
		},
	}

	foundOutput := &ec2.DescribeInstancesOutput{
		Reservations: []types.Reservation{
			{
				Instances: []types.Instance{
					{
						InstanceId: aws.String("id123"),
						State: &types.InstanceState{
							Name: types.InstanceStateNameRunning,
						},
						Placement: &types.Placement{
							AvailabilityZone: aws.String("us-gov-east-1"),
						},
					},
				},
			},
		},
	}

	tests := []struct {
		name           string
		bastionEnabled bool
		expect         func(m *mocks.MockEC2APIMockRecorder)
		expectError    bool
		bastionStatus  *infrav1.Instance
	}{
		{
			name: "Should ignore reconciliation if instance not found",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.
					DescribeInstances(context.TODO(), gomock.Eq(describeInput)).
					Return(&ec2.DescribeInstancesOutput{}, nil)
			},
			expectError: false,
		},
		{
			name: "Should fail reconcile if describe instance fails",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.
					DescribeInstances(context.TODO(), gomock.Eq(describeInput)).
					Return(nil, errors.New("some error"))
			},
			expectError: true,
		},
		{
			name: "Should fail reconcile if terminate instance fails",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.
					DescribeInstances(context.TODO(), gomock.Eq(describeInput)).
					Return(foundOutput, nil).MinTimes(1)
				m.
					TerminateInstances(context.TODO(),
						gomock.Eq(&ec2.TerminateInstancesInput{
							InstanceIds: []string{"id123"},
						}),
					).
					Return(nil, errors.New("some error"))
			},
			expectError: true,
		},
		{
			name: "Should create bastion successfully",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeInstances(context.TODO(), gomock.Eq(describeInput)).
					Return(&ec2.DescribeInstancesOutput{}, nil).MinTimes(1)
				m.DescribeImages(context.TODO(), gomock.Eq(&ec2.DescribeImagesInput{Filters: []types.Filter{
					{
						Name:   aws.String("architecture"),
						Values: []string{"x86_64"},
					},
					{
						Name:   aws.String("state"),
						Values: []string{"available"},
					},
					{
						Name:   aws.String("virtualization-type"),
						Values: []string{"hvm"},
					},
					{
						Name:   aws.String("description"),
						Values: []string{ubuntuImageDescription},
					},
					{
						Name:   aws.String("owner-id"),
						Values: []string{ubuntuOwnerIDUsGov},
					},
				}})).Return(&ec2.DescribeImagesOutput{Images: images{
					{
						ImageId:      aws.String("ubuntu-ami-id-latest"),
						CreationDate: aws.String("2019-02-08T17:02:31.000Z"),
					},
					{
						ImageId:      aws.String("ubuntu-ami-id-old"),
						CreationDate: aws.String("2014-02-08T17:02:31.000Z"),
					},
				}}, nil)
				m.RunInstances(context.TODO(), gomock.Any()).
					Return(&ec2.RunInstancesOutput{
						Instances: []types.Instance{
							{
								State: &types.InstanceState{
									Name: types.InstanceStateNameRunning,
								},
								IamInstanceProfile: &types.IamInstanceProfile{
									Arn: aws.String("arn:aws-us-gov:iam::123456789012:instance-profile/foo"),
								},
								InstanceId:     aws.String("id123"),
								InstanceType:   types.InstanceTypeT3Micro,
								SubnetId:       aws.String("subnet-1"),
								ImageId:        aws.String("ubuntu-ami-id-latest"),
								RootDeviceName: aws.String("device-1"),
								BlockDeviceMappings: []types.InstanceBlockDeviceMapping{
									{
										DeviceName: aws.String("device-1"),
										Ebs: &types.EbsInstanceBlockDevice{
											VolumeId: aws.String("volume-1"),
										},
									},
								},
								Placement: &types.Placement{
									AvailabilityZone: aws.String("us-gov-east-1"),
								},
							},
						},
					}, nil)
			},
			bastionEnabled: true,
			expectError:    false,
			bastionStatus: &infrav1.Instance{
				ID:               "id123",
				State:            "running",
				Type:             "t3.micro",
				SubnetID:         "subnet-1",
				ImageID:          "ubuntu-ami-id-latest",
				IAMProfile:       "foo",
				Addresses:        []clusterv1.MachineAddress{},
				AvailabilityZone: "us-gov-east-1",
				VolumeIDs:        []string{"volume-1"},
			},
		},
	}

	for _, tc := range tests {
		managedValues := []bool{false, true}
		for i := range managedValues {
			managed := managedValues[i]

			t.Run(fmt.Sprintf("managed=%t %s", managed, tc.name), func(t *testing.T) {
				g := NewWithT(t)

				mockControl := gomock.NewController(t)
				defer mockControl.Finish()

				ec2Mock := mocks.NewMockEC2API(mockControl)

				scheme, err := setupScheme()
				g.Expect(err).To(BeNil())

				awsCluster := &infrav1.AWSCluster{
					ObjectMeta: metav1.ObjectMeta{Name: "test"},
					Spec: infrav1.AWSClusterSpec{
						NetworkSpec: infrav1.NetworkSpec{
							VPC: infrav1.VPCSpec{
								ID: "vpcID",
							},
							Subnets: infrav1.Subnets{
								{
									ID: "subnet-1",
								},
								{
									ID:       "subnet-2",
									IsPublic: true,
								},
							},
						},
						Bastion: infrav1.Bastion{Enabled: tc.bastionEnabled},
						Region:  "us-gov-east-1",
					},
				}

				client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(awsCluster).WithStatusSubresource(awsCluster).Build()

				scope, err := scope.NewClusterScope(scope.ClusterScopeParams{
					Cluster: &clusterv1.Cluster{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "ns",
							Name:      clusterName,
						},
					},
					AWSCluster: awsCluster,
					Client:     client,
				})
				g.Expect(err).To(BeNil())

				if managed {
					scope.AWSCluster.Spec.NetworkSpec.VPC.Tags = infrav1.Tags{
						infrav1.ClusterTagKey(clusterName): string(infrav1.ResourceLifecycleOwned),
					}
				}

				tc.expect(ec2Mock.EXPECT())
				s := NewService(scope)
				s.EC2Client = ec2Mock

				err = s.ReconcileBastion()
				if tc.expectError {
					g.Expect(err).NotTo(BeNil())
					return
				}

				g.Expect(err).To(BeNil())

				g.Expect(scope.AWSCluster.Status.Bastion).To(BeEquivalentTo(tc.bastionStatus))
			})
		}
	}
}
