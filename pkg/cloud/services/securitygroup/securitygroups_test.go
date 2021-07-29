/*
Copyright 2019 The Kubernetes Authors.

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

package securitygroup

import (
	"context"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/ec2/mock_ec2iface"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestReconcileSecurityGroups(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name   string
		input  *infrav1.NetworkSpec
		expect func(m *mock_ec2iface.MockEC2APIMockRecorder)
		err    error
	}{
		{
			name: "no existing",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID:                "vpc-securitygroups",
					InternetGatewayID: aws.String("igw-01"),
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
				},
				Subnets: infrav1.Subnets{
					&infrav1.SubnetSpec{
						ID:               "subnet-securitygroups-private",
						IsPublic:         false,
						AvailabilityZone: "us-east-1a",
					},
					&infrav1.SubnetSpec{
						ID:               "subnet-securitygroups-public",
						IsPublic:         true,
						NatGatewayID:     aws.String("nat-01"),
						AvailabilityZone: "us-east-1a",
					},
				},
			},
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.DescribeSecurityGroups(gomock.AssignableToTypeOf(&ec2.DescribeSecurityGroupsInput{})).
					Return(&ec2.DescribeSecurityGroupsOutput{}, nil)

				securityGroupBastion := m.CreateSecurityGroup(gomock.Eq(&ec2.CreateSecurityGroupInput{
					VpcId:       aws.String("vpc-securitygroups"),
					GroupName:   aws.String("test-cluster-bastion"),
					Description: aws.String("Kubernetes cluster test-cluster: bastion"),
					TagSpecifications: []*ec2.TagSpecification{
						{
							ResourceType: aws.String("security-group"),
							Tags: []*ec2.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-bastion"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
									Value: aws.String("owned"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
									Value: aws.String("bastion"),
								},
							},
						},
					},
				})).
					Return(&ec2.CreateSecurityGroupOutput{GroupId: aws.String("sg-bastion")}, nil)

				m.AuthorizeSecurityGroupIngress(gomock.AssignableToTypeOf(&ec2.AuthorizeSecurityGroupIngressInput{
					GroupId: aws.String("sg-bastion"),
				})).
					Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil).
					After(securityGroupBastion)

				////////////////////////

				securityGroupAPIServerLb := m.CreateSecurityGroup(gomock.Eq(&ec2.CreateSecurityGroupInput{
					VpcId:       aws.String("vpc-securitygroups"),
					GroupName:   aws.String("test-cluster-apiserver-lb"),
					Description: aws.String("Kubernetes cluster test-cluster: apiserver-lb"),
					TagSpecifications: []*ec2.TagSpecification{
						{
							ResourceType: aws.String("security-group"),
							Tags: []*ec2.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-apiserver-lb"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
									Value: aws.String("owned"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
									Value: aws.String("apiserver-lb"),
								},
							},
						},
					},
				})).
					Return(&ec2.CreateSecurityGroupOutput{GroupId: aws.String("sg-apiserver-lb")}, nil)

				m.AuthorizeSecurityGroupIngress(gomock.AssignableToTypeOf(&ec2.AuthorizeSecurityGroupIngressInput{
					GroupId: aws.String("sg-apiserver-lb"),
				})).
					Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil).
					After(securityGroupAPIServerLb)

				////////////////////////

				m.CreateSecurityGroup(gomock.Eq(&ec2.CreateSecurityGroupInput{
					VpcId:       aws.String("vpc-securitygroups"),
					GroupName:   aws.String("test-cluster-lb"),
					Description: aws.String("Kubernetes cluster test-cluster: lb"),
					TagSpecifications: []*ec2.TagSpecification{
						{
							ResourceType: aws.String("security-group"),
							Tags: []*ec2.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-lb"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("owned"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
									Value: aws.String("owned"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
									Value: aws.String("lb"),
								},
							},
						},
					},
				})).
					Return(&ec2.CreateSecurityGroupOutput{GroupId: aws.String("sg-lb")}, nil)

				////////////////////////

				securityGroupControl := m.CreateSecurityGroup(gomock.Eq(&ec2.CreateSecurityGroupInput{
					VpcId:       aws.String("vpc-securitygroups"),
					GroupName:   aws.String("test-cluster-controlplane"),
					Description: aws.String("Kubernetes cluster test-cluster: controlplane"),
					TagSpecifications: []*ec2.TagSpecification{
						{
							ResourceType: aws.String("security-group"),
							Tags: []*ec2.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-controlplane"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
									Value: aws.String("owned"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
									Value: aws.String("controlplane"),
								},
							},
						},
					},
				})).
					Return(&ec2.CreateSecurityGroupOutput{GroupId: aws.String("sg-control")}, nil)

				m.AuthorizeSecurityGroupIngress(gomock.AssignableToTypeOf(&ec2.AuthorizeSecurityGroupIngressInput{
					GroupId: aws.String("sg-control"),
				})).
					Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil).
					After(securityGroupControl)

				//////////////////////////////////////////////

				securityGroupNode := m.CreateSecurityGroup(gomock.Eq(&ec2.CreateSecurityGroupInput{
					VpcId:       aws.String("vpc-securitygroups"),
					GroupName:   aws.String("test-cluster-node"),
					Description: aws.String("Kubernetes cluster test-cluster: node"),
					TagSpecifications: []*ec2.TagSpecification{
						{
							ResourceType: aws.String("security-group"),
							Tags: []*ec2.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-node"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
									Value: aws.String("owned"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
									Value: aws.String("node"),
								},
							},
						},
					},
				})).
					Return(&ec2.CreateSecurityGroupOutput{GroupId: aws.String("sg-node")}, nil)

				m.AuthorizeSecurityGroupIngress(gomock.AssignableToTypeOf(&ec2.AuthorizeSecurityGroupIngressInput{
					GroupId: aws.String("sg-node"),
				})).
					Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil).
					After(securityGroupNode)
			},
		},
		{
			name: "all overridden, do not tag",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID:                "vpc-securitygroups",
					InternetGatewayID: aws.String("igw-01"),
				},
				Subnets: infrav1.Subnets{
					&infrav1.SubnetSpec{
						ID:               "subnet-securitygroups-private",
						IsPublic:         false,
						AvailabilityZone: "us-east-1a",
					},
					&infrav1.SubnetSpec{
						ID:               "subnet-securitygroups-public",
						IsPublic:         true,
						NatGatewayID:     aws.String("nat-01"),
						AvailabilityZone: "us-east-1a",
					},
				},
				SecurityGroupOverrides: map[infrav1.SecurityGroupRole]string{
					infrav1.SecurityGroupBastion:      "sg-bastion",
					infrav1.SecurityGroupAPIServerLB:  "sg-apiserver-lb",
					infrav1.SecurityGroupLB:           "sg-lb",
					infrav1.SecurityGroupControlPlane: "sg-control",
					infrav1.SecurityGroupNode:         "sg-node",
				},
			},
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.DescribeSecurityGroups(gomock.AssignableToTypeOf(&ec2.DescribeSecurityGroupsInput{})).
					Return(&ec2.DescribeSecurityGroupsOutput{
						SecurityGroups: []*ec2.SecurityGroup{
							{GroupId: aws.String("sg-bastion"), GroupName: aws.String("Bastion Security Group")},
							{GroupId: aws.String("sg-apiserver-lb"), GroupName: aws.String("API load balancer Security Group")},
							{GroupId: aws.String("sg-lb"), GroupName: aws.String("Load balancer Security Group")},
							{GroupId: aws.String("sg-control"), GroupName: aws.String("Control plane Security Group")},
							{GroupId: aws.String("sg-node"), GroupName: aws.String("Node Security Group")},
						},
					}, nil).AnyTimes()

			},
		},
		{
			name: "managed vpc with overrides, returns error",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID:                "vpc-securitygroups",
					InternetGatewayID: aws.String("igw-01"),
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
				},
				Subnets: infrav1.Subnets{
					&infrav1.SubnetSpec{
						ID:               "subnet-securitygroups-private",
						IsPublic:         false,
						AvailabilityZone: "us-east-1a",
					},
					&infrav1.SubnetSpec{
						ID:               "subnet-securitygroups-public",
						IsPublic:         true,
						NatGatewayID:     aws.String("nat-01"),
						AvailabilityZone: "us-east-1a",
					},
				},
				SecurityGroupOverrides: map[infrav1.SecurityGroupRole]string{
					infrav1.SecurityGroupBastion:      "sg-bastion",
					infrav1.SecurityGroupAPIServerLB:  "sg-apiserver-lb",
					infrav1.SecurityGroupLB:           "sg-lb",
					infrav1.SecurityGroupControlPlane: "sg-control",
					infrav1.SecurityGroupNode:         "sg-node",
				},
			},
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.DescribeSecurityGroups(gomock.AssignableToTypeOf(&ec2.DescribeSecurityGroupsInput{})).
					Return(&ec2.DescribeSecurityGroupsOutput{
						SecurityGroups: []*ec2.SecurityGroup{
							{GroupId: aws.String("sg-bastion"), GroupName: aws.String("Bastion Security Group")},
							{GroupId: aws.String("sg-apiserver-lb"), GroupName: aws.String("API load balancer Security Group")},
							{GroupId: aws.String("sg-lb"), GroupName: aws.String("Load balancer Security Group")},
							{GroupId: aws.String("sg-control"), GroupName: aws.String("Control plane Security Group")},
							{GroupId: aws.String("sg-node"), GroupName: aws.String("Node Security Group")},
						},
					}, nil).AnyTimes()

			},
			err: errors.New(`security group overrides provided for managed vpc "test-cluster"`),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ec2Mock := mock_ec2iface.NewMockEC2API(mockCtrl)

			scope, err := scope.NewClusterScope(scope.ClusterScopeParams{
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"},
				},
				AWSCluster: &infrav1.AWSCluster{
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

			if err := s.ReconcileSecurityGroups(); err != nil && tc.err != nil {
				if !strings.Contains(err.Error(), tc.err.Error()) {
					t.Fatalf("was expecting error to look like '%v', but got '%v'", tc.err, err)
				}
			} else if err != nil {
				t.Fatalf("got an unexpected error: %v", err)
			}
		})
	}
}

func TestControlPlaneSecurityGroupNotOpenToAnyCIDR(t *testing.T) {
	scope, err := scope.NewClusterScope(scope.ClusterScopeParams{
		Cluster: &clusterv1.Cluster{
			ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"},
		},
		AWSCluster: &infrav1.AWSCluster{},
	})
	if err != nil {
		t.Fatalf("Failed to create test context: %v", err)
	}

	s := NewService(scope)
	rules, err := s.getSecurityGroupIngressRules(infrav1.SecurityGroupControlPlane)
	if err != nil {
		t.Fatalf("Failed to lookup controlplane security group ingress rules: %v", err)
	}

	for _, r := range rules {
		if sets.NewString(r.CidrBlocks...).Has(services.AnyIPv4CidrBlock) {
			t.Fatal("Ingress rule allows any CIDR block")
		}
	}
}

func TestDeleteSecurityGroups(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name   string
		input  *infrav1.NetworkSpec
		expect func(m *mock_ec2iface.MockEC2APIMockRecorder)
		err    error
	}{
		{
			name: "do not delete overridden security groups, only delete 'owned' SGs",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID:                "vpc-securitygroups",
					InternetGatewayID: aws.String("igw-01"),
				},
				Subnets: infrav1.Subnets{
					&infrav1.SubnetSpec{
						ID:               "subnet-securitygroups-private",
						IsPublic:         false,
						AvailabilityZone: "us-east-1a",
					},
					&infrav1.SubnetSpec{
						ID:               "subnet-securitygroups-public",
						IsPublic:         true,
						NatGatewayID:     aws.String("nat-01"),
						AvailabilityZone: "us-east-1a",
					},
				},
				SecurityGroupOverrides: map[infrav1.SecurityGroupRole]string{
					infrav1.SecurityGroupBastion:      "sg-bastion",
					infrav1.SecurityGroupAPIServerLB:  "sg-apiserver-lb",
					infrav1.SecurityGroupLB:           "sg-lb",
					infrav1.SecurityGroupControlPlane: "sg-control",
					infrav1.SecurityGroupNode:         "sg-node",
				},
			},
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.DescribeSecurityGroupsPages(gomock.Any(), gomock.Any()).Do(func(_, y interface{}) {
					funct := y.(func(output *ec2.DescribeSecurityGroupsOutput, lastPage bool) bool)
					funct(&ec2.DescribeSecurityGroupsOutput{
						SecurityGroups: []*ec2.SecurityGroup{
							{
								GroupName: aws.String("sg-bastion"),
								GroupId:   aws.String("sg-bastion"),

								Tags: []*ec2.Tag{
									{
										Key:   aws.String("Name"),
										Value: aws.String("test-cluster-nat"),
									},
									{
										Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
										Value: aws.String("owned"),
									},
									{
										Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
										Value: aws.String("common"),
									},
								},
							},
						},
					}, true)
				}).Return(nil)

				m.DescribeSecurityGroups(gomock.Any()).Return(&ec2.DescribeSecurityGroupsOutput{}, nil)
				m.DeleteSecurityGroup(gomock.Any()).Return(nil, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ec2Mock := mock_ec2iface.NewMockEC2API(mockCtrl)

			scheme := runtime.NewScheme()
			_ = infrav1.AddToScheme(scheme)
			awsCluster := &infrav1.AWSCluster{
				TypeMeta: metav1.TypeMeta{
					APIVersion: infrav1.GroupVersion.String(),
					Kind:       "AWSCluster",
				},
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: *tc.input,
				},
			}
			client := fake.NewFakeClientWithScheme(scheme, awsCluster)

			ctx := context.TODO()
			client.Create(ctx, awsCluster)

			scope, err := scope.NewClusterScope(scope.ClusterScopeParams{
				Client: client,
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"},
				},
				AWSCluster: awsCluster,
			})
			if err != nil {
				t.Fatalf("Failed to create test context: %v", err)
			}

			tc.expect(ec2Mock.EXPECT())

			s := NewService(scope)
			s.EC2Client = ec2Mock

			if err := s.DeleteSecurityGroups(); err != nil && tc.err != nil {
				if !strings.Contains(err.Error(), tc.err.Error()) {
					t.Fatalf("was expecting error to look like '%v', but got '%v'", tc.err, err)
				}
			} else if err != nil {
				t.Fatalf("got an unexpected error: %v", err)
			}
		})
	}
}
