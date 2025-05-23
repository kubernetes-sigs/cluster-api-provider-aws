/*
Copyright 2024 The Kubernetes Authors.

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

package network

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func TestReconcileCarrierGateway(t *testing.T) {
	testCases := []struct {
		name   string
		input  *infrav1.NetworkSpec
		expect func(m *mocks.MockEC2APIMockRecorder)
	}{
		{
			name: "has cagw",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID: "vpc-cagw",
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeCarrierGateways(context.TODO(), gomock.Eq(&ec2.DescribeCarrierGatewaysInput{
					Filters: []types.Filter{
						{
							Name:   aws.String("vpc-id"),
							Values: []string{"vpc-cagw"},
						},
					},
				})).
					Return(&ec2.DescribeCarrierGatewaysOutput{
						CarrierGateways: []types.CarrierGateway{
							{
								CarrierGatewayId: ptr.To("cagw-01"),
							},
						},
					}, nil).AnyTimes()

				m.CreateTags(context.TODO(), gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
					Return(nil, nil).AnyTimes()
			},
		},
		{
			name: "no cagw attached, creates one",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID: "vpc-cagw",
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeCarrierGateways(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeCarrierGatewaysInput{})).
					Return(&ec2.DescribeCarrierGatewaysOutput{}, nil).AnyTimes()

				m.CreateCarrierGateway(context.TODO(), gomock.AssignableToTypeOf(&ec2.CreateCarrierGatewayInput{})).
					Return(&ec2.CreateCarrierGatewayOutput{
						CarrierGateway: &types.CarrierGateway{
							CarrierGatewayId: aws.String("cagw-1"),
							VpcId:            aws.String("vpc-cagw"),
							Tags: []types.Tag{
								{
									Key:   aws.String(infrav1.ClusterTagKey("test-cluster")),
									Value: aws.String("owned"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
									Value: aws.String("common"),
								},
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-cagw"),
								},
							},
						},
					}, nil).AnyTimes()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			ec2Mock := mocks.NewMockEC2API(mockCtrl)

			scheme := runtime.NewScheme()
			_ = infrav1.AddToScheme(scheme)
			client := fake.NewClientBuilder().WithScheme(scheme).Build()
			scope, err := scope.NewClusterScope(scope.ClusterScopeParams{
				Client: client,
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"},
				},
				AWSCluster: &infrav1.AWSCluster{
					ObjectMeta: metav1.ObjectMeta{Name: "test"},
					Spec: infrav1.AWSClusterSpec{
						NetworkSpec: *tc.input,
					},
				},
			})
			if err != nil {
				t.Fatalf("Failed to create test context: %v", err)
			}

			tc.expect(ec2Mock.EXPECT())

			s := NewService(scope)
			s.EC2Client = ec2Mock

			if err := s.reconcileCarrierGateway(); err != nil {
				t.Fatalf("got an unexpected error: %v", err)
			}
			mockCtrl.Finish()
		})
	}
}

func TestDeleteCarrierGateway(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name    string
		input   *infrav1.NetworkSpec
		expect  func(m *mocks.MockEC2APIMockRecorder)
		wantErr bool
	}{
		{
			name: "Should ignore deletion if vpc is unmanaged",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID: "vpc-cagw",
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {},
		},
		{
			name: "Should ignore deletion if carrier gateway is not found",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID: "vpc-cagw",
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeCarrierGateways(context.TODO(), gomock.Eq(&ec2.DescribeCarrierGatewaysInput{
					Filters: []types.Filter{
						{
							Name:   aws.String("vpc-id"),
							Values: []string{"vpc-cagw"},
						},
					},
				})).Return(&ec2.DescribeCarrierGatewaysOutput{}, nil)
			},
		},
		{
			name: "Should successfully delete the carrier gateway",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID: "vpc-cagw",
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeCarrierGateways(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeCarrierGatewaysInput{})).
					Return(&ec2.DescribeCarrierGatewaysOutput{
						CarrierGateways: []types.CarrierGateway{
							{
								CarrierGatewayId: aws.String("cagw-0"),
								VpcId:            aws.String("vpc-gateways"),
							},
						},
					}, nil)

				m.DeleteCarrierGateway(context.TODO(), &ec2.DeleteCarrierGatewayInput{
					CarrierGatewayId: aws.String("cagw-0"),
				}).Return(&ec2.DeleteCarrierGatewayOutput{}, nil)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			ec2Mock := mocks.NewMockEC2API(mockCtrl)

			scheme := runtime.NewScheme()
			err := infrav1.AddToScheme(scheme)
			g.Expect(err).NotTo(HaveOccurred())
			client := fake.NewClientBuilder().WithScheme(scheme).Build()

			scope, err := scope.NewClusterScope(scope.ClusterScopeParams{
				Client: client,
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"},
				},
				AWSCluster: &infrav1.AWSCluster{
					ObjectMeta: metav1.ObjectMeta{Name: "test"},
					Spec: infrav1.AWSClusterSpec{
						NetworkSpec: *tc.input,
					},
				},
			})
			g.Expect(err).NotTo(HaveOccurred())

			tc.expect(ec2Mock.EXPECT())

			s := NewService(scope)
			s.EC2Client = ec2Mock

			err = s.deleteCarrierGateway()
			if tc.wantErr {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).NotTo(HaveOccurred())
		})
	}
}
