/*
Copyright 2025 The Kubernetes Authors.

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
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
)

func TestDeleteOrphanedENIs(t *testing.T) {
	testCases := []struct {
		name          string
		input         ScopeBuilder
		expect        func(m *mocks.MockEC2APIMockRecorder)
		errorExpected bool
	}{
		{
			name: "No orphaned ENIs found",
			input: NewClusterScope().
				WithNetwork(&infrav1.NetworkSpec{
					VPC: infrav1.VPCSpec{
						ID: subnetsVPCID,
						Tags: infrav1.Tags{
							infrav1.ClusterTagKey("test-cluster"): "owned",
						},
					},
				}).
				WithSecurityGroups(map[infrav1.SecurityGroupRole]infrav1.SecurityGroup{
					infrav1.SecurityGroupLB:   {ID: "sg-lb"},
					infrav1.SecurityGroupNode: {ID: "sg-node"},
				}),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeNetworkInterfaces(context.TODO(), gomock.Any()).
					Return(&ec2.DescribeNetworkInterfacesOutput{
						NetworkInterfaces: []types.NetworkInterface{},
					}, nil)
			},
			errorExpected: false,
		},
		{
			name: "Available ENIs found and deleted",
			input: NewClusterScope().
				WithNetwork(&infrav1.NetworkSpec{
					VPC: infrav1.VPCSpec{
						ID: subnetsVPCID,
						Tags: infrav1.Tags{
							infrav1.ClusterTagKey("test-cluster"): "owned",
						},
					},
				}).
				WithSecurityGroups(map[infrav1.SecurityGroupRole]infrav1.SecurityGroup{
					infrav1.SecurityGroupLB:   {ID: "sg-lb"},
					infrav1.SecurityGroupNode: {ID: "sg-node"},
				}),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeNetworkInterfaces(context.TODO(), gomock.Any()).
					Return(&ec2.DescribeNetworkInterfacesOutput{
						NetworkInterfaces: []types.NetworkInterface{
							{NetworkInterfaceId: aws.String("eni-111")},
							{NetworkInterfaceId: aws.String("eni-222")},
						},
					}, nil)
				m.DeleteNetworkInterface(context.TODO(), &ec2.DeleteNetworkInterfaceInput{
					NetworkInterfaceId: aws.String("eni-111"),
				}).Return(&ec2.DeleteNetworkInterfaceOutput{}, nil)
				m.DeleteNetworkInterface(context.TODO(), &ec2.DeleteNetworkInterfaceInput{
					NetworkInterfaceId: aws.String("eni-222"),
				}).Return(&ec2.DeleteNetworkInterfaceOutput{}, nil)
			},
			errorExpected: false,
		},
		{
			name: "One deletion fails, other still attempted",
			input: NewClusterScope().
				WithNetwork(&infrav1.NetworkSpec{
					VPC: infrav1.VPCSpec{
						ID: subnetsVPCID,
						Tags: infrav1.Tags{
							infrav1.ClusterTagKey("test-cluster"): "owned",
						},
					},
				}).
				WithSecurityGroups(map[infrav1.SecurityGroupRole]infrav1.SecurityGroup{
					infrav1.SecurityGroupLB: {ID: "sg-lb"},
				}),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeNetworkInterfaces(context.TODO(), gomock.Any()).
					Return(&ec2.DescribeNetworkInterfacesOutput{
						NetworkInterfaces: []types.NetworkInterface{
							{NetworkInterfaceId: aws.String("eni-111")},
							{NetworkInterfaceId: aws.String("eni-222")},
						},
					}, nil)
				m.DeleteNetworkInterface(context.TODO(), &ec2.DeleteNetworkInterfaceInput{
					NetworkInterfaceId: aws.String("eni-111"),
				}).Return(nil, fmt.Errorf("InvalidNetworkInterfaceID.NotFound"))
				m.DeleteNetworkInterface(context.TODO(), &ec2.DeleteNetworkInterfaceInput{
					NetworkInterfaceId: aws.String("eni-222"),
				}).Return(&ec2.DeleteNetworkInterfaceOutput{}, nil)
			},
			errorExpected: true,
		},
		{
			name: "No security groups, skip entirely",
			input: NewClusterScope().
				WithNetwork(&infrav1.NetworkSpec{
					VPC: infrav1.VPCSpec{
						ID: subnetsVPCID,
						Tags: infrav1.Tags{
							infrav1.ClusterTagKey("test-cluster"): "owned",
						},
					},
				}),
			expect:        func(m *mocks.MockEC2APIMockRecorder) {},
			errorExpected: false,
		},
		{
			name: "Unmanaged VPC, skip entirely",
			input: NewClusterScope().
				WithNetwork(&infrav1.NetworkSpec{
					VPC: infrav1.VPCSpec{
						ID: subnetsVPCID,
					},
				}).
				WithSecurityGroups(map[infrav1.SecurityGroupRole]infrav1.SecurityGroup{
					infrav1.SecurityGroupLB: {ID: "sg-lb"},
				}),
			expect:        func(m *mocks.MockEC2APIMockRecorder) {},
			errorExpected: false,
		},
		{
			name: "DescribeNetworkInterfaces fails",
			input: NewClusterScope().
				WithNetwork(&infrav1.NetworkSpec{
					VPC: infrav1.VPCSpec{
						ID: subnetsVPCID,
						Tags: infrav1.Tags{
							infrav1.ClusterTagKey("test-cluster"): "owned",
						},
					},
				}).
				WithSecurityGroups(map[infrav1.SecurityGroupRole]infrav1.SecurityGroup{
					infrav1.SecurityGroupLB: {ID: "sg-lb"},
				}),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeNetworkInterfaces(context.TODO(), gomock.Any()).
					Return(nil, fmt.Errorf("aws api error"))
			},
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			ec2Mock := mocks.NewMockEC2API(mockCtrl)

			scope, err := tc.input.Build()
			g.Expect(err).NotTo(HaveOccurred())

			s := NewService(scope)
			s.EC2Client = ec2Mock

			tc.expect(ec2Mock.EXPECT())
			err = s.deleteOrphanedENIs()

			if tc.errorExpected {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).NotTo(HaveOccurred())
			}
		})
	}
}
