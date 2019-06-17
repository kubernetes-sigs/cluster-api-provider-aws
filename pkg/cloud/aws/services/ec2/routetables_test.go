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
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/ec2/mock_ec2iface"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/elb/mock_elbiface"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

func TestReconcileRouteTables(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name   string
		input  *v1alpha1.NetworkSpec
		expect func(m *mock_ec2iface.MockEC2APIMockRecorder)
		err    error
	}{
		{
			name: "no routes existing, single private and single public, same AZ",
			input: &v1alpha1.NetworkSpec{
				VPC: v1alpha1.VPCSpec{
					ID:                "vpc-routetables",
					InternetGatewayID: aws.String("igw-01"),
					Tags: v1alpha1.Tags{
						v1alpha1.ClusterTagKey("test-cluster"): "owned",
					},
				},
				Subnets: v1alpha1.Subnets{
					&v1alpha1.SubnetSpec{
						ID:               "subnet-routetables-private",
						IsPublic:         false,
						AvailabilityZone: "us-east-1a",
					},
					&v1alpha1.SubnetSpec{
						ID:               "subnet-routetables-public",
						IsPublic:         true,
						NatGatewayID:     aws.String("nat-01"),
						AvailabilityZone: "us-east-1a",
					},
				},
			},
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.DescribeRouteTables(gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				privateRouteTable := m.CreateRouteTable(gomock.Eq(&ec2.CreateRouteTableInput{VpcId: aws.String("vpc-routetables")})).
					Return(&ec2.CreateRouteTableOutput{RouteTable: &ec2.RouteTable{RouteTableId: aws.String("rt-1")}}, nil)

				m.CreateTags(gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
					Return(nil, nil)

				m.CreateRoute(gomock.Eq(&ec2.CreateRouteInput{
					NatGatewayId:         aws.String("nat-01"),
					DestinationCidrBlock: aws.String("0.0.0.0/0"),
					RouteTableId:         aws.String("rt-1"),
				})).
					After(privateRouteTable)

				m.AssociateRouteTable(gomock.Eq(&ec2.AssociateRouteTableInput{
					RouteTableId: aws.String("rt-1"),
					SubnetId:     aws.String("subnet-routetables-private"),
				})).
					Return(&ec2.AssociateRouteTableOutput{}, nil).
					After(privateRouteTable)

				publicRouteTable := m.CreateRouteTable(gomock.Eq(&ec2.CreateRouteTableInput{VpcId: aws.String("vpc-routetables")})).
					Return(&ec2.CreateRouteTableOutput{RouteTable: &ec2.RouteTable{RouteTableId: aws.String("rt-2")}}, nil)

				m.CreateRoute(gomock.Eq(&ec2.CreateRouteInput{
					GatewayId:            aws.String("igw-01"),
					DestinationCidrBlock: aws.String("0.0.0.0/0"),
					RouteTableId:         aws.String("rt-2"),
				})).
					After(publicRouteTable)

				m.CreateTags(gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
					Return(nil, nil)

				m.AssociateRouteTable(gomock.Eq(&ec2.AssociateRouteTableInput{
					RouteTableId: aws.String("rt-2"),
					SubnetId:     aws.String("subnet-routetables-public"),
				})).
					Return(&ec2.AssociateRouteTableOutput{}, nil).
					After(publicRouteTable)
			},
		},
		{
			name: "subnets in different availability zones, returns error",
			input: &v1alpha1.NetworkSpec{
				VPC: v1alpha1.VPCSpec{
					InternetGatewayID: aws.String("igw-01"),
					ID:                "vpc-routetables",
					Tags: v1alpha1.Tags{
						v1alpha1.ClusterTagKey("test-cluster"): "owned",
					},
				},
				Subnets: v1alpha1.Subnets{
					&v1alpha1.SubnetSpec{
						ID:               "subnet-routetables-private",
						IsPublic:         false,
						AvailabilityZone: "us-east-1a",
					},
					&v1alpha1.SubnetSpec{
						ID:               "subnet-routetables-public",
						IsPublic:         true,
						NatGatewayID:     aws.String("nat-01"),
						AvailabilityZone: "us-east-1b",
					},
				},
			},
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.DescribeRouteTables(gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)
			},
			err: errors.New(`no nat gateways available in "us-east-1a"`),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ec2Mock := mock_ec2iface.NewMockEC2API(mockCtrl)
			elbMock := mock_elbiface.NewMockELBAPI(mockCtrl)

			scope, err := actuators.NewScope(actuators.ScopeParams{
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"},
				},
				AWSClients: actuators.AWSClients{
					EC2: ec2Mock,
					ELB: elbMock,
				},
			})

			if err != nil {
				t.Fatalf("Failed to create test context: %v", err)
			}

			scope.ClusterConfig = &v1alpha1.AWSClusterProviderSpec{
				NetworkSpec: *tc.input,
			}

			tc.expect(ec2Mock.EXPECT())

			s := NewService(scope)
			if err := s.reconcileRouteTables(); err != nil && tc.err != nil {
				if !strings.Contains(err.Error(), tc.err.Error()) {
					t.Fatalf("was expecting error to look like '%v', but got '%v'", tc.err, err)
				}
			} else if err != nil {
				t.Fatalf("got an unexpected error: %v", err)
			}
		})
	}
}
