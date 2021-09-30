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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/filter"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/ec2/mock_ec2iface"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestDeleteBastion(t *testing.T) {
	clusterName := "cluster"

	describeInput := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			filter.EC2.ProviderRole(infrav1.BastionRoleTagValue),
			filter.EC2.Cluster(clusterName),
			filter.EC2.InstanceStates(
				ec2.InstanceStateNamePending,
				ec2.InstanceStateNameRunning,
				ec2.InstanceStateNameStopping,
				ec2.InstanceStateNameStopped,
			),
		},
	}

	foundOutput := &ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{
			{
				Instances: []*ec2.Instance{
					{
						InstanceId: aws.String("id123"),
						State: &ec2.InstanceState{
							Name: aws.String(ec2.InstanceStateNameRunning),
						},
						Placement: &ec2.Placement{
							AvailabilityZone: aws.String("us-east-1"),
						},
					},
				},
			},
		},
	}

	tests := []struct {
		name        string
		expect      func(m *mock_ec2iface.MockEC2APIMockRecorder)
		expectError bool
	}{
		{
			name: "instance not found",
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.
					DescribeInstances(gomock.Eq(describeInput)).
					Return(&ec2.DescribeInstancesOutput{}, nil)
			},
			expectError: false,
		},
		{
			name: "describe error",
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.
					DescribeInstances(gomock.Eq(describeInput)).
					Return(nil, errors.New("some error"))
			},
			expectError: true,
		},
		{
			name: "terminate fails",
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.
					DescribeInstances(gomock.Eq(describeInput)).
					Return(foundOutput, nil)
				m.
					TerminateInstances(
						gomock.Eq(&ec2.TerminateInstancesInput{
							InstanceIds: aws.StringSlice([]string{"id123"}),
						}),
					).
					Return(nil, errors.New("some error"))
			},
			expectError: true,
		},
		{
			name: "wait after terminate fails",
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.
					DescribeInstances(gomock.Eq(describeInput)).
					Return(foundOutput, nil)
				m.
					TerminateInstances(
						gomock.Eq(&ec2.TerminateInstancesInput{
							InstanceIds: aws.StringSlice([]string{"id123"}),
						}),
					).
					Return(nil, nil)
				m.
					WaitUntilInstanceTerminated(
						gomock.Eq(&ec2.DescribeInstancesInput{
							InstanceIds: aws.StringSlice([]string{"id123"}),
						}),
					).
					Return(errors.New("some error"))
			},
			expectError: true,
		},
		{
			name: "success",
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.
					DescribeInstances(gomock.Eq(describeInput)).
					Return(foundOutput, nil)
				m.
					TerminateInstances(
						gomock.Eq(&ec2.TerminateInstancesInput{
							InstanceIds: aws.StringSlice([]string{"id123"}),
						}),
					).
					Return(nil, nil)
				m.
					WaitUntilInstanceTerminated(
						gomock.Eq(&ec2.DescribeInstancesInput{
							InstanceIds: aws.StringSlice([]string{"id123"}),
						}),
					).
					Return(nil)
			},
			expectError: false,
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

				ec2Mock := mock_ec2iface.NewMockEC2API(mockControl)

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

				client := fake.NewClientBuilder().WithScheme(scheme).Build()
				ctx := context.TODO()
				client.Create(ctx, awsCluster)

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
			})
		}
	}
}
