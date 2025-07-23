/*
Copyright 2022 The Kubernetes Authors.

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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func TestReconcileEgressOnlyInternetGateways(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name              string
		input             *infrav1.NetworkSpec
		expect            func(m *mocks.MockEC2APIMockRecorder)
		wantErrContaining *string
	}{
		{
			name: "has eigw",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID:   "vpc-egress-only-gateways",
					IPv6: &infrav1.IPv6{},
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeEgressOnlyInternetGateways(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeEgressOnlyInternetGatewaysInput{})).
					Return(&ec2.DescribeEgressOnlyInternetGatewaysOutput{
						EgressOnlyInternetGateways: []types.EgressOnlyInternetGateway{
							{
								EgressOnlyInternetGatewayId: aws.String("eigw-0"),
								Attachments: []types.InternetGatewayAttachment{
									{
										State: types.AttachmentStatusAttached,
										VpcId: aws.String("vpc-egress-only-gateways"),
									},
								},
							},
						},
					}, nil)

				m.CreateTags(context.TODO(), gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
					Return(nil, nil)
			},
		},
		{
			name: "has more than 1 eigw, should return error",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID:   "vpc-egress-only-gateways",
					IPv6: &infrav1.IPv6{},
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
				},
			},
			wantErrContaining: aws.String("expected 1 egress only internet gateway in vpc \"vpc-egress-only-gateways\", but found 2: [eigw-0 eigw-1]"),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeEgressOnlyInternetGateways(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeEgressOnlyInternetGatewaysInput{})).
					Return(&ec2.DescribeEgressOnlyInternetGatewaysOutput{
						EgressOnlyInternetGateways: []types.EgressOnlyInternetGateway{
							{
								EgressOnlyInternetGatewayId: aws.String("eigw-0"),
								Attachments: []types.InternetGatewayAttachment{
									{
										State: types.AttachmentStatusAttached,
										VpcId: aws.String("vpc-egress-only-gateways"),
									},
								},
							},
							{
								EgressOnlyInternetGatewayId: aws.String("eigw-1"),
								Attachments: []types.InternetGatewayAttachment{
									{
										State: types.AttachmentStatusAttached,
										VpcId: aws.String("vpc-egress-only-gateways"),
									},
								},
							},
						},
					}, nil)
			},
		},
		{
			name: "no eigw attached, creates one",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					IPv6: &infrav1.IPv6{},
					ID:   "vpc-egress-only-gateways",
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeEgressOnlyInternetGateways(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeEgressOnlyInternetGatewaysInput{})).
					Return(&ec2.DescribeEgressOnlyInternetGatewaysOutput{}, nil)

				m.CreateEgressOnlyInternetGateway(context.TODO(), gomock.AssignableToTypeOf(&ec2.CreateEgressOnlyInternetGatewayInput{})).
					Return(&ec2.CreateEgressOnlyInternetGatewayOutput{
						EgressOnlyInternetGateway: &types.EgressOnlyInternetGateway{
							EgressOnlyInternetGatewayId: aws.String("igw-1"),
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
									Value: aws.String("test-cluster-eigw"),
								},
							},
							Attachments: []types.InternetGatewayAttachment{
								{
									State: types.AttachmentStatusAttached,
									VpcId: aws.String("vpc-egress-only-gateways"),
								},
							},
						},
					}, nil)
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

			err = s.reconcileEgressOnlyInternetGateways()
			if tc.wantErrContaining != nil {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(ContainSubstring(*tc.wantErrContaining))
				return
			}
			g.Expect(err).NotTo(HaveOccurred())
		})
	}
}

func TestDeleteEgressOnlyInternetGateways(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name    string
		input   *infrav1.NetworkSpec
		expect  func(m *mocks.MockEC2APIMockRecorder)
		wantErr bool
	}{
		{
			name: "Should ignore deletion if vpc is not ipv6",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID: "vpc-gateways",
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {},
		},
		{
			name: "Should ignore deletion if vpc is unmanaged",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					IPv6: &infrav1.IPv6{},
					ID:   "vpc-gateways",
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {},
		},
		{
			name: "Should ignore deletion if egress only internet gateway is not found",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					IPv6: &infrav1.IPv6{},
					ID:   "vpc-gateways",
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeEgressOnlyInternetGateways(context.TODO(), gomock.Eq(&ec2.DescribeEgressOnlyInternetGatewaysInput{
					Filters: []types.Filter{
						{
							Name:   aws.String("tag-key"),
							Values: []string{infrav1.ClusterTagKey("test-cluster")},
						},
					},
				})).Return(&ec2.DescribeEgressOnlyInternetGatewaysOutput{}, nil)
			},
		},
		{
			name: "Should successfully delete the egress only internet gateway",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID: "vpc-gateways",
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
					IPv6: &infrav1.IPv6{},
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeEgressOnlyInternetGateways(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeEgressOnlyInternetGatewaysInput{})).
					Return(&ec2.DescribeEgressOnlyInternetGatewaysOutput{
						EgressOnlyInternetGateways: []types.EgressOnlyInternetGateway{
							{
								EgressOnlyInternetGatewayId: aws.String("eigw-0"),
								Attachments: []types.InternetGatewayAttachment{
									{
										State: types.AttachmentStatusAttached,
										VpcId: aws.String("vpc-gateways"),
									},
								},
							},
						},
					}, nil)
				m.DeleteEgressOnlyInternetGateway(context.TODO(), &ec2.DeleteEgressOnlyInternetGatewayInput{
					EgressOnlyInternetGatewayId: aws.String("eigw-0"),
				}).Return(&ec2.DeleteEgressOnlyInternetGatewayOutput{}, nil)
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

			err = s.deleteEgressOnlyInternetGateways()
			if tc.wantErr {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).NotTo(HaveOccurred())
		})
	}
}
