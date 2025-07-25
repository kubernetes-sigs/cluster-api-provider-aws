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
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/filter"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

var (
	testSecurityGroupRoles = []infrav1.SecurityGroupRole{
		infrav1.SecurityGroupBastion,
		infrav1.SecurityGroupAPIServerLB,
		infrav1.SecurityGroupLB,
		infrav1.SecurityGroupControlPlane,
		infrav1.SecurityGroupNode,
	}
)

func TestReconcileSecurityGroups(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name       string
		input      *infrav1.NetworkSpec
		expect     func(m *mocks.MockEC2APIMockRecorder)
		err        error
		awsCluster func(acl infrav1.AWSCluster) infrav1.AWSCluster
	}{
		{
			name: "no existing",
			awsCluster: func(acl infrav1.AWSCluster) infrav1.AWSCluster {
				return acl
			},
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID:                "vpc-securitygroups",
					InternetGatewayID: aws.String("igw-01"),
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
					EmptyRoutesDefaultVPCSecurityGroup: true,
				},
				Subnets: infrav1.Subnets{
					infrav1.SubnetSpec{
						ID:               "subnet-securitygroups-private",
						IsPublic:         false,
						AvailabilityZone: "us-east-1a",
					},
					infrav1.SubnetSpec{
						ID:               "subnet-securitygroups-public",
						IsPublic:         true,
						NatGatewayID:     aws.String("nat-01"),
						AvailabilityZone: "us-east-1a",
					},
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSecurityGroups(context.TODO(), &ec2.DescribeSecurityGroupsInput{
					Filters: []types.Filter{
						filter.EC2.VPC("vpc-securitygroups"),
						filter.EC2.SecurityGroupName("default"),
					},
				}).
					Return(&ec2.DescribeSecurityGroupsOutput{
						SecurityGroups: []types.SecurityGroup{
							{
								Description: aws.String("default VPC security group"),
								GroupName:   aws.String("default"),
								GroupId:     aws.String("sg-default"),
							},
						},
					}, nil)
				m.RevokeSecurityGroupIngress(context.TODO(), gomock.AssignableToTypeOf(&ec2.RevokeSecurityGroupIngressInput{
					GroupId: aws.String("sg-default"),
				}))

				m.RevokeSecurityGroupEgress(context.TODO(), gomock.AssignableToTypeOf(&ec2.RevokeSecurityGroupEgressInput{
					GroupId: aws.String("sg-default"),
				}))

				m.DescribeSecurityGroups(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeSecurityGroupsInput{})).
					Return(&ec2.DescribeSecurityGroupsOutput{}, nil)

				securityGroupBastion := m.CreateSecurityGroup(context.TODO(), gomock.Eq(&ec2.CreateSecurityGroupInput{
					VpcId:       aws.String("vpc-securitygroups"),
					GroupName:   aws.String("test-cluster-bastion"),
					Description: aws.String("Kubernetes cluster test-cluster: bastion"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSecurityGroup,
							Tags: []types.Tag{
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

				m.AuthorizeSecurityGroupIngress(context.TODO(), gomock.AssignableToTypeOf(&ec2.AuthorizeSecurityGroupIngressInput{
					GroupId: aws.String("sg-bastion"),
				})).
					Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil).
					After(securityGroupBastion)

				securityGroupAPIServerLb := m.CreateSecurityGroup(context.TODO(), gomock.Eq(&ec2.CreateSecurityGroupInput{
					VpcId:       aws.String("vpc-securitygroups"),
					GroupName:   aws.String("test-cluster-apiserver-lb"),
					Description: aws.String("Kubernetes cluster test-cluster: apiserver-lb"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSecurityGroup,
							Tags: []types.Tag{
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

				m.AuthorizeSecurityGroupIngress(context.TODO(), gomock.AssignableToTypeOf(&ec2.AuthorizeSecurityGroupIngressInput{
					GroupId: aws.String("sg-apiserver-lb"),
				})).
					Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil).
					After(securityGroupAPIServerLb)

				m.CreateSecurityGroup(context.TODO(), gomock.Eq(&ec2.CreateSecurityGroupInput{
					VpcId:       aws.String("vpc-securitygroups"),
					GroupName:   aws.String("test-cluster-lb"),
					Description: aws.String("Kubernetes cluster test-cluster: lb"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSecurityGroup,
							Tags: []types.Tag{
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

				securityGroupControl := m.CreateSecurityGroup(context.TODO(), gomock.Eq(&ec2.CreateSecurityGroupInput{
					VpcId:       aws.String("vpc-securitygroups"),
					GroupName:   aws.String("test-cluster-controlplane"),
					Description: aws.String("Kubernetes cluster test-cluster: controlplane"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSecurityGroup,
							Tags: []types.Tag{
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

				m.AuthorizeSecurityGroupIngress(context.TODO(), gomock.AssignableToTypeOf(&ec2.AuthorizeSecurityGroupIngressInput{
					GroupId: aws.String("sg-control"),
				})).
					Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil).
					After(securityGroupControl)

				securityGroupNode := m.CreateSecurityGroup(context.TODO(), gomock.Eq(&ec2.CreateSecurityGroupInput{
					VpcId:       aws.String("vpc-securitygroups"),
					GroupName:   aws.String("test-cluster-node"),
					Description: aws.String("Kubernetes cluster test-cluster: node"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSecurityGroup,
							Tags: []types.Tag{
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

				m.AuthorizeSecurityGroupIngress(context.TODO(), gomock.AssignableToTypeOf(&ec2.AuthorizeSecurityGroupIngressInput{
					GroupId: aws.String("sg-node"),
				})).
					Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil).
					After(securityGroupNode)
			},
		},
		{
			name: "NLB is defined with preserve client IP disabled",
			awsCluster: func(acl infrav1.AWSCluster) infrav1.AWSCluster {
				acl.Spec.ControlPlaneLoadBalancer = &infrav1.AWSLoadBalancerSpec{
					LoadBalancerType: infrav1.LoadBalancerTypeNLB,
				}
				return acl
			},
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID:                "vpc-securitygroups",
					InternetGatewayID: aws.String("igw-01"),
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
					CidrBlock: "10.0.0.0/16",
				},
				Subnets: infrav1.Subnets{
					infrav1.SubnetSpec{
						ID:               "subnet-securitygroups-private",
						IsPublic:         false,
						AvailabilityZone: "us-east-1a",
					},
					infrav1.SubnetSpec{
						ID:               "subnet-securitygroups-public",
						IsPublic:         true,
						NatGatewayID:     aws.String("nat-01"),
						AvailabilityZone: "us-east-1a",
					},
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSecurityGroups(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeSecurityGroupsInput{})).
					Return(&ec2.DescribeSecurityGroupsOutput{}, nil)

				securityGroupBastion := m.CreateSecurityGroup(context.TODO(), gomock.Eq(&ec2.CreateSecurityGroupInput{
					VpcId:       aws.String("vpc-securitygroups"),
					GroupName:   aws.String("test-cluster-bastion"),
					Description: aws.String("Kubernetes cluster test-cluster: bastion"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSecurityGroup,
							Tags: []types.Tag{
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

				m.AuthorizeSecurityGroupIngress(context.TODO(), gomock.AssignableToTypeOf(&ec2.AuthorizeSecurityGroupIngressInput{
					GroupId: aws.String("sg-bastion"),
				})).
					Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil).
					After(securityGroupBastion)

				securityGroupAPIServerLb := m.CreateSecurityGroup(context.TODO(), gomock.Eq(&ec2.CreateSecurityGroupInput{
					VpcId:       aws.String("vpc-securitygroups"),
					GroupName:   aws.String("test-cluster-apiserver-lb"),
					Description: aws.String("Kubernetes cluster test-cluster: apiserver-lb"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSecurityGroup,
							Tags: []types.Tag{
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

				m.AuthorizeSecurityGroupIngress(context.TODO(), gomock.AssignableToTypeOf(&ec2.AuthorizeSecurityGroupIngressInput{
					GroupId: aws.String("sg-apiserver-lb"),
				})).
					Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil).
					After(securityGroupAPIServerLb)

				lbSecurityGroup := m.CreateSecurityGroup(context.TODO(), gomock.Eq(&ec2.CreateSecurityGroupInput{
					VpcId:       aws.String("vpc-securitygroups"),
					GroupName:   aws.String("test-cluster-lb"),
					Description: aws.String("Kubernetes cluster test-cluster: lb"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSecurityGroup,
							Tags: []types.Tag{
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
				})).Return(&ec2.CreateSecurityGroupOutput{GroupId: aws.String("sg-lb")}, nil)

				m.AuthorizeSecurityGroupIngress(context.TODO(), gomock.AssignableToTypeOf(&ec2.AuthorizeSecurityGroupIngressInput{
					GroupId: aws.String("sg-lb"),
				})).
					Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil).
					After(lbSecurityGroup)

				securityGroupControl := m.CreateSecurityGroup(context.TODO(), gomock.Eq(&ec2.CreateSecurityGroupInput{
					VpcId:       aws.String("vpc-securitygroups"),
					GroupName:   aws.String("test-cluster-controlplane"),
					Description: aws.String("Kubernetes cluster test-cluster: controlplane"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSecurityGroup,
							Tags: []types.Tag{
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

				m.AuthorizeSecurityGroupIngress(context.TODO(), gomock.AssignableToTypeOf(&ec2.AuthorizeSecurityGroupIngressInput{
					GroupId: aws.String("sg-control"),
				})).
					Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil).
					After(securityGroupControl)

				securityGroupNode := m.CreateSecurityGroup(context.TODO(), gomock.Eq(&ec2.CreateSecurityGroupInput{
					VpcId:       aws.String("vpc-securitygroups"),
					GroupName:   aws.String("test-cluster-node"),
					Description: aws.String("Kubernetes cluster test-cluster: node"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSecurityGroup,
							Tags: []types.Tag{
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

				m.AuthorizeSecurityGroupIngress(context.TODO(), gomock.AssignableToTypeOf(&ec2.AuthorizeSecurityGroupIngressInput{
					GroupId: aws.String("sg-node"),
				})).
					Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil).
					After(securityGroupNode)
			},
		},
		{
			name: "all overrides defined, do not tag",
			awsCluster: func(acl infrav1.AWSCluster) infrav1.AWSCluster {
				return acl
			},
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID:                "vpc-securitygroups",
					InternetGatewayID: aws.String("igw-01"),
				},
				Subnets: infrav1.Subnets{
					infrav1.SubnetSpec{
						ID:               "subnet-securitygroups-private",
						IsPublic:         false,
						AvailabilityZone: "us-east-1a",
					},
					infrav1.SubnetSpec{
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
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSecurityGroups(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeSecurityGroupsInput{})).
					Return(&ec2.DescribeSecurityGroupsOutput{
						SecurityGroups: []types.SecurityGroup{
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
			name: "additional tags includes cloud provider tag, only tag lb",
			awsCluster: func(acl infrav1.AWSCluster) infrav1.AWSCluster {
				acl.Spec.AdditionalTags = infrav1.Tags{
					infrav1.ClusterAWSCloudProviderTagKey("test-cluster"): "owned",
				}
				return acl
			},
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID:                "vpc-securitygroups",
					InternetGatewayID: aws.String("igw-01"),
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
				},
				Subnets: infrav1.Subnets{
					infrav1.SubnetSpec{
						ID:               "subnet-securitygroups-private",
						IsPublic:         false,
						AvailabilityZone: "us-east-1a",
					},
					infrav1.SubnetSpec{
						ID:               "subnet-securitygroups-public",
						IsPublic:         true,
						NatGatewayID:     aws.String("nat-01"),
						AvailabilityZone: "us-east-1a",
					},
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSecurityGroups(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeSecurityGroupsInput{})).
					Return(&ec2.DescribeSecurityGroupsOutput{}, nil)

				securityGroupBastion := m.CreateSecurityGroup(context.TODO(), gomock.Eq(&ec2.CreateSecurityGroupInput{
					VpcId:       aws.String("vpc-securitygroups"),
					GroupName:   aws.String("test-cluster-bastion"),
					Description: aws.String("Kubernetes cluster test-cluster: bastion"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSecurityGroup,
							Tags: []types.Tag{
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

				m.AuthorizeSecurityGroupIngress(context.TODO(), gomock.AssignableToTypeOf(&ec2.AuthorizeSecurityGroupIngressInput{
					GroupId: aws.String("sg-bastion"),
				})).
					Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil).
					After(securityGroupBastion)

				securityGroupAPIServerLb := m.CreateSecurityGroup(context.TODO(), gomock.Eq(&ec2.CreateSecurityGroupInput{
					VpcId:       aws.String("vpc-securitygroups"),
					GroupName:   aws.String("test-cluster-apiserver-lb"),
					Description: aws.String("Kubernetes cluster test-cluster: apiserver-lb"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSecurityGroup,
							Tags: []types.Tag{
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

				m.AuthorizeSecurityGroupIngress(context.TODO(), gomock.AssignableToTypeOf(&ec2.AuthorizeSecurityGroupIngressInput{
					GroupId: aws.String("sg-apiserver-lb"),
				})).
					Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil).
					After(securityGroupAPIServerLb)

				lbSecurityGroup := m.CreateSecurityGroup(context.TODO(), gomock.Eq(&ec2.CreateSecurityGroupInput{
					VpcId:       aws.String("vpc-securitygroups"),
					GroupName:   aws.String("test-cluster-lb"),
					Description: aws.String("Kubernetes cluster test-cluster: lb"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSecurityGroup,
							Tags: []types.Tag{
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
				})).Return(&ec2.CreateSecurityGroupOutput{GroupId: aws.String("sg-lb")}, nil)

				m.AuthorizeSecurityGroupIngress(context.TODO(), gomock.AssignableToTypeOf(&ec2.AuthorizeSecurityGroupIngressInput{
					GroupId: aws.String("sg-lb"),
				})).
					Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil).
					After(lbSecurityGroup)

				securityGroupControl := m.CreateSecurityGroup(context.TODO(), gomock.Eq(&ec2.CreateSecurityGroupInput{
					VpcId:       aws.String("vpc-securitygroups"),
					GroupName:   aws.String("test-cluster-controlplane"),
					Description: aws.String("Kubernetes cluster test-cluster: controlplane"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSecurityGroup,
							Tags: []types.Tag{
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

				m.AuthorizeSecurityGroupIngress(context.TODO(), gomock.AssignableToTypeOf(&ec2.AuthorizeSecurityGroupIngressInput{
					GroupId: aws.String("sg-control"),
				})).
					Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil).
					After(securityGroupControl)

				m.CreateSecurityGroup(context.TODO(), gomock.Eq(&ec2.CreateSecurityGroupInput{
					VpcId:       aws.String("vpc-securitygroups"),
					GroupName:   aws.String("test-cluster-node"),
					Description: aws.String("Kubernetes cluster test-cluster: node"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSecurityGroup,
							Tags: []types.Tag{
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
			},
		},
		{
			name: "managed vpc with overrides, returns error",
			awsCluster: func(acl infrav1.AWSCluster) infrav1.AWSCluster {
				return acl
			},
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID:                "vpc-securitygroups",
					InternetGatewayID: aws.String("igw-01"),
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
				},
				Subnets: infrav1.Subnets{
					infrav1.SubnetSpec{
						ID:               "subnet-securitygroups-private",
						IsPublic:         false,
						AvailabilityZone: "us-east-1a",
					},
					infrav1.SubnetSpec{
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
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSecurityGroups(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeSecurityGroupsInput{})).
					Return(&ec2.DescribeSecurityGroupsOutput{
						SecurityGroups: []types.SecurityGroup{
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
		{
			name: "when VPC default security group has no rules then no errors are returned",
			awsCluster: func(acl infrav1.AWSCluster) infrav1.AWSCluster {
				return acl
			},
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID:                "vpc-securitygroups",
					InternetGatewayID: aws.String("igw-01"),
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
					EmptyRoutesDefaultVPCSecurityGroup: true,
				},
				Subnets: infrav1.Subnets{
					infrav1.SubnetSpec{
						ID:               "subnet-securitygroups-private",
						IsPublic:         false,
						AvailabilityZone: "us-east-1a",
					},
					infrav1.SubnetSpec{
						ID:               "subnet-securitygroups-public",
						IsPublic:         true,
						NatGatewayID:     aws.String("nat-01"),
						AvailabilityZone: "us-east-1a",
					},
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSecurityGroups(context.TODO(), &ec2.DescribeSecurityGroupsInput{
					Filters: []types.Filter{
						filter.EC2.VPC("vpc-securitygroups"),
						filter.EC2.SecurityGroupName("default"),
					},
				}).
					Return(&ec2.DescribeSecurityGroupsOutput{
						SecurityGroups: []types.SecurityGroup{
							{
								Description: aws.String("default VPC security group"),
								GroupName:   aws.String("default"),
								GroupId:     aws.String("sg-default"),
							},
						},
					}, nil)

				m.RevokeSecurityGroupIngress(context.TODO(), gomock.AssignableToTypeOf(&ec2.RevokeSecurityGroupIngressInput{
					GroupId: aws.String("sg-default"),
				})).Return(&ec2.RevokeSecurityGroupIngressOutput{}, awserr.New("InvalidPermission.NotFound", "rules not found in security group", nil))

				m.RevokeSecurityGroupEgress(context.TODO(), gomock.AssignableToTypeOf(&ec2.RevokeSecurityGroupEgressInput{
					GroupId: aws.String("sg-default"),
				})).Return(&ec2.RevokeSecurityGroupEgressOutput{}, awserr.New("InvalidPermission.NotFound", "rules not found in security group", nil))

				m.DescribeSecurityGroups(context.TODO(), &ec2.DescribeSecurityGroupsInput{
					Filters: []types.Filter{
						filter.EC2.VPC("vpc-securitygroups"),
						filter.EC2.Cluster("test-cluster"),
					},
				}).Return(&ec2.DescribeSecurityGroupsOutput{}, nil)

				m.CreateSecurityGroup(context.TODO(), gomock.AssignableToTypeOf(&ec2.CreateSecurityGroupInput{})).
					Return(&ec2.CreateSecurityGroupOutput{GroupId: aws.String("sg-node")}, nil).AnyTimes()

				m.AuthorizeSecurityGroupIngress(context.TODO(), gomock.AssignableToTypeOf(&ec2.AuthorizeSecurityGroupIngressInput{})).
					Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil).AnyTimes()
			},
		},
		{
			name: "authorized target ingress rules are not revoked",
			awsCluster: func(acl infrav1.AWSCluster) infrav1.AWSCluster {
				return acl
			},
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID:                "vpc-securitygroups",
					InternetGatewayID: aws.String("igw-01"),
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
					EmptyRoutesDefaultVPCSecurityGroup: true,
				},
				Subnets: infrav1.Subnets{
					infrav1.SubnetSpec{
						ID:               "subnet-securitygroups-private",
						IsPublic:         false,
						AvailabilityZone: "us-east-1a",
					},
					infrav1.SubnetSpec{
						ID:               "subnet-securitygroups-public",
						IsPublic:         true,
						NatGatewayID:     aws.String("nat-01"),
						AvailabilityZone: "us-east-1a",
					},
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSecurityGroups(context.TODO(), &ec2.DescribeSecurityGroupsInput{
					Filters: []types.Filter{
						filter.EC2.VPC("vpc-securitygroups"),
						filter.EC2.SecurityGroupName("default"),
					},
				}).
					Return(&ec2.DescribeSecurityGroupsOutput{
						SecurityGroups: []types.SecurityGroup{
							{
								Description: aws.String("default VPC security group"),
								GroupName:   aws.String("default"),
								GroupId:     aws.String("sg-default"),
							},
						},
					}, nil)

				m.RevokeSecurityGroupIngress(context.TODO(), gomock.Eq(&ec2.RevokeSecurityGroupIngressInput{
					GroupId: aws.String("sg-default"),
					IpPermissions: []types.IpPermission{
						{
							IpProtocol: aws.String("-1"),
							UserIdGroupPairs: []types.UserIdGroupPair{
								{
									GroupId: aws.String("sg-default"),
								},
							},
						},
					},
				})).Times(1)

				m.RevokeSecurityGroupEgress(context.TODO(), gomock.AssignableToTypeOf(&ec2.RevokeSecurityGroupEgressInput{
					GroupId: aws.String("sg-default"),
				}))

				securityGroupBastion := types.SecurityGroup{
					Description: aws.String("Kubernetes cluster test-cluster: bastion"),
					GroupName:   aws.String("test-cluster-bastion"),
					GroupId:     aws.String("sg-bastion"),
					Tags: []types.Tag{
						{
							Key:   aws.String("Name"),
							Value: aws.String("test-cluster-bastion"),
						}, {
							Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
							Value: aws.String("owned"),
						}, {
							Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
							Value: aws.String("bastion"),
						},
					},
				}

				securityGroupLB := types.SecurityGroup{
					Description: aws.String("Kubernetes cluster test-cluster: lb"),
					GroupName:   aws.String("test-cluster-lb"),
					GroupId:     aws.String("sg-lb"),
					Tags: []types.Tag{
						{
							Key:   aws.String("Name"),
							Value: aws.String("test-cluster-lb"),
						}, {
							Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
							Value: aws.String("owned"),
						}, {
							Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
							Value: aws.String("lb"),
						}, {
							Key:   aws.String("kubernetes.io/cluster/test-cluster"),
							Value: aws.String("owned"),
						},
					},
				}

				securityGroupAPIServerLB := types.SecurityGroup{
					Description: aws.String("Kubernetes cluster test-cluster: apiserver-lb"),
					GroupName:   aws.String("test-cluster-apiserver-lb"),
					GroupId:     aws.String("sg-apiserver-lb"),
					Tags: []types.Tag{
						{
							Key:   aws.String("Name"),
							Value: aws.String("test-cluster-apiserver-lb"),
						}, {
							Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
							Value: aws.String("owned"),
						}, {
							Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
							Value: aws.String("apiserver-lb"),
						},
					},
					IpPermissions: []types.IpPermission{
						{
							FromPort:   aws.Int32(6443),
							IpProtocol: aws.String("tcp"),
							IpRanges: []types.IpRange{
								{
									CidrIp:      aws.String("0.0.0.0/0"),
									Description: aws.String("Kubernetes API"),
								},
							},
							ToPort: aws.Int32(6443),
						},
						// Extra rule to be revoked
						{
							FromPort:   aws.Int32(22),
							IpProtocol: aws.String("tcp"),
							ToPort:     aws.Int32(22),
							IpRanges: []types.IpRange{
								{
									CidrIp:      aws.String("0.0.0.0/0"),
									Description: aws.String("SSH"),
								},
							},
						},
					},
				}

				securityGroupControl := types.SecurityGroup{
					Description: aws.String("Kubernetes cluster test-cluster: controlplane"),
					GroupName:   aws.String("test-cluster-controlplane"),
					GroupId:     aws.String("sg-control"),
					Tags: []types.Tag{
						{
							Key:   aws.String("Name"),
							Value: aws.String("test-cluster-controlplane"),
						}, {
							Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
							Value: aws.String("owned"),
						}, {
							Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
							Value: aws.String("controlplane"),
						},
					},
					IpPermissions: []types.IpPermission{
						{
							FromPort:   aws.Int32(6443),
							IpProtocol: aws.String("tcp"),
							ToPort:     aws.Int32(6443),
							UserIdGroupPairs: []types.UserIdGroupPair{
								{
									Description: aws.String("Kubernetes API"),
									GroupId:     aws.String("sg-apiserver-lb"),
								}, {
									Description: aws.String("Kubernetes API"),
									GroupId:     aws.String("sg-control"),
								}, {
									Description: aws.String("Kubernetes API"),
									GroupId:     aws.String("sg-node"),
								},
							},
						},
						{
							FromPort:   aws.Int32(2379),
							IpProtocol: aws.String("tcp"),
							ToPort:     aws.Int32(2379),
							UserIdGroupPairs: []types.UserIdGroupPair{
								{
									Description: aws.String("etcd"),
									GroupId:     aws.String("sg-control"),
								},
							},
						},
						{
							FromPort:   aws.Int32(2380),
							IpProtocol: aws.String("tcp"),
							ToPort:     aws.Int32(2380),
							UserIdGroupPairs: []types.UserIdGroupPair{
								{
									Description: aws.String("etcd peer"),
									GroupId:     aws.String("sg-control"),
								},
							},
						},
					},
				}

				securityGroupNode := types.SecurityGroup{
					Description: aws.String("Kubernetes cluster test-cluster: node"),
					GroupName:   aws.String("test-cluster-node"),
					GroupId:     aws.String("sg-node"),
					Tags: []types.Tag{
						{
							Key:   aws.String("Name"),
							Value: aws.String("test-cluster-node"),
						}, {
							Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
							Value: aws.String("owned"),
						}, {
							Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
							Value: aws.String("node"),
						},
					},
					IpPermissions: []types.IpPermission{
						{
							FromPort:   aws.Int32(30000),
							ToPort:     aws.Int32(32767),
							IpProtocol: aws.String("tcp"),
							IpRanges: []types.IpRange{
								{
									CidrIp:      aws.String("0.0.0.0/0"),
									Description: aws.String("Node Port Services"),
								},
							},
						}, {
							FromPort:   aws.Int32(10250),
							IpProtocol: aws.String("tcp"),
							ToPort:     aws.Int32(10250),
							UserIdGroupPairs: []types.UserIdGroupPair{
								{
									Description: aws.String("Kubelet API"),
									GroupId:     aws.String("sg-control"),
								}, {
									Description: aws.String("Kubelet API"),
									GroupId:     aws.String("sg-node"),
								},
							},
						},
					},
				}

				m.DescribeSecurityGroups(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeSecurityGroupsInput{})).
					Return(&ec2.DescribeSecurityGroupsOutput{
						SecurityGroups: []types.SecurityGroup{
							securityGroupBastion,
							securityGroupLB,
							securityGroupAPIServerLB,
							securityGroupControl,
							securityGroupNode,
						},
					}, nil)

				m.RevokeSecurityGroupIngress(context.TODO(), gomock.Eq(&ec2.RevokeSecurityGroupIngressInput{
					GroupId: aws.String("sg-apiserver-lb"),
					IpPermissions: []types.IpPermission{
						{
							FromPort:   aws.Int32(22),
							ToPort:     aws.Int32(22),
							IpProtocol: aws.String("tcp"),
							IpRanges: []types.IpRange{
								{
									CidrIp:      aws.String("0.0.0.0/0"),
									Description: aws.String("SSH"),
								},
							},
						},
					},
				})).Times(1)

				m.AuthorizeSecurityGroupIngress(context.TODO(), gomock.AssignableToTypeOf(&ec2.AuthorizeSecurityGroupIngressInput{
					GroupId: aws.String("sg-bastion"),
					IpPermissions: []types.IpPermission{
						{
							ToPort:     aws.Int32(22),
							FromPort:   aws.Int32(22),
							IpProtocol: aws.String("tcp"),
						},
					},
				})).
					Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil).AnyTimes()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ec2Mock := mocks.NewMockEC2API(mockCtrl)

			scheme := runtime.NewScheme()
			_ = infrav1.AddToScheme(scheme)
			client := fake.NewClientBuilder().WithScheme(scheme).Build()
			cluster := &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: *tc.input,
				},
			}
			awsCluster := tc.awsCluster(*cluster)
			cs, err := scope.NewClusterScope(scope.ClusterScopeParams{
				Client: client,
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"},
				},
				AWSCluster: &awsCluster,
			})
			if err != nil {
				t.Fatalf("Failed to create test context: %v", err)
			}

			tc.expect(ec2Mock.EXPECT())

			s := NewService(cs, testSecurityGroupRoles)
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
	scheme := runtime.NewScheme()
	_ = infrav1.AddToScheme(scheme)
	client := fake.NewClientBuilder().WithScheme(scheme).Build()
	cs, err := scope.NewClusterScope(scope.ClusterScopeParams{
		Client: client,
		Cluster: &clusterv1.Cluster{
			ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"},
		},
		AWSCluster: &infrav1.AWSCluster{},
	})
	if err != nil {
		t.Fatalf("Failed to create test context: %v", err)
	}

	s := NewService(cs, testSecurityGroupRoles)
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

func TestAdditionalControlPlaneSecurityGroup(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = infrav1.AddToScheme(scheme)

	testCases := []struct {
		name                          string
		networkSpec                   infrav1.NetworkSpec
		networkStatus                 infrav1.NetworkStatus
		expectedAdditionalIngressRule infrav1.IngressRule
		wantErr                       bool
	}{
		{
			name: "default control plane security group is used",
			networkSpec: infrav1.NetworkSpec{
				AdditionalControlPlaneIngressRules: []infrav1.IngressRule{
					{
						Description: "test",
						Protocol:    infrav1.SecurityGroupProtocolTCP,
						FromPort:    9345,
						ToPort:      9345,
					},
				},
			},
			networkStatus: infrav1.NetworkStatus{
				SecurityGroups: map[infrav1.SecurityGroupRole]infrav1.SecurityGroup{
					infrav1.SecurityGroupControlPlane: {
						ID: "cp-sg-id",
					},
					infrav1.SecurityGroupNode: {
						ID: "node-sg-id",
					},
				},
			},
			expectedAdditionalIngressRule: infrav1.IngressRule{
				Description:            "test",
				Protocol:               infrav1.SecurityGroupProtocolTCP,
				FromPort:               9345,
				ToPort:                 9345,
				SourceSecurityGroupIDs: []string{"cp-sg-id"},
			},
		},
		{
			name: "custom security group id is used",
			networkSpec: infrav1.NetworkSpec{
				AdditionalControlPlaneIngressRules: []infrav1.IngressRule{
					{
						Description:            "test",
						Protocol:               infrav1.SecurityGroupProtocolTCP,
						FromPort:               9345,
						ToPort:                 9345,
						SourceSecurityGroupIDs: []string{"test"},
					},
				},
			},
			networkStatus: infrav1.NetworkStatus{
				SecurityGroups: map[infrav1.SecurityGroupRole]infrav1.SecurityGroup{
					infrav1.SecurityGroupControlPlane: {
						ID: "cp-sg-id",
					},
					infrav1.SecurityGroupNode: {
						ID: "node-sg-id",
					},
				},
			},
			expectedAdditionalIngressRule: infrav1.IngressRule{
				Description:            "test",
				Protocol:               infrav1.SecurityGroupProtocolTCP,
				FromPort:               9345,
				ToPort:                 9345,
				SourceSecurityGroupIDs: []string{"test"},
			},
		},
		{
			name: "another security group role is used",
			networkSpec: infrav1.NetworkSpec{
				AdditionalControlPlaneIngressRules: []infrav1.IngressRule{
					{
						Description:              "test",
						Protocol:                 infrav1.SecurityGroupProtocolTCP,
						FromPort:                 9345,
						ToPort:                   9345,
						SourceSecurityGroupRoles: []infrav1.SecurityGroupRole{infrav1.SecurityGroupNode},
					},
				},
			},
			networkStatus: infrav1.NetworkStatus{
				SecurityGroups: map[infrav1.SecurityGroupRole]infrav1.SecurityGroup{
					infrav1.SecurityGroupControlPlane: {
						ID: "cp-sg-id",
					},
					infrav1.SecurityGroupNode: {
						ID: "node-sg-id",
					},
				},
			},
			expectedAdditionalIngressRule: infrav1.IngressRule{
				Description:            "test",
				Protocol:               infrav1.SecurityGroupProtocolTCP,
				FromPort:               9345,
				ToPort:                 9345,
				SourceSecurityGroupIDs: []string{"node-sg-id"},
			},
		},
		{
			name: "another security group role and a custom security group id is used",
			networkSpec: infrav1.NetworkSpec{
				AdditionalControlPlaneIngressRules: []infrav1.IngressRule{
					{
						Description:              "test",
						Protocol:                 infrav1.SecurityGroupProtocolTCP,
						FromPort:                 9345,
						ToPort:                   9345,
						SourceSecurityGroupIDs:   []string{"test"},
						SourceSecurityGroupRoles: []infrav1.SecurityGroupRole{infrav1.SecurityGroupNode},
					},
				},
			},
			networkStatus: infrav1.NetworkStatus{
				SecurityGroups: map[infrav1.SecurityGroupRole]infrav1.SecurityGroup{
					infrav1.SecurityGroupControlPlane: {
						ID: "cp-sg-id",
					},
					infrav1.SecurityGroupNode: {
						ID: "node-sg-id",
					},
				},
			},
			expectedAdditionalIngressRule: infrav1.IngressRule{
				Description:            "test",
				Protocol:               infrav1.SecurityGroupProtocolTCP,
				FromPort:               9345,
				ToPort:                 9345,
				SourceSecurityGroupIDs: []string{"test", "node-sg-id"},
			},
		},
		{
			name: "don't set source security groups if cidr blocks are set",
			networkSpec: infrav1.NetworkSpec{
				AdditionalControlPlaneIngressRules: []infrav1.IngressRule{
					{
						Description: "test",
						Protocol:    infrav1.SecurityGroupProtocolTCP,
						FromPort:    9345,
						ToPort:      9345,
						CidrBlocks:  []string{"test-cidr-block"},
					},
				},
			},
			networkStatus: infrav1.NetworkStatus{
				SecurityGroups: map[infrav1.SecurityGroupRole]infrav1.SecurityGroup{
					infrav1.SecurityGroupControlPlane: {
						ID: "cp-sg-id",
					},
					infrav1.SecurityGroupNode: {
						ID: "node-sg-id",
					},
				},
			},
			expectedAdditionalIngressRule: infrav1.IngressRule{
				Description: "test",
				Protocol:    infrav1.SecurityGroupProtocolTCP,
				FromPort:    9345,
				ToPort:      9345,
			},
		},
		{
			name: "set nat gateway IPs cidr as source if specified",
			networkSpec: infrav1.NetworkSpec{
				AdditionalControlPlaneIngressRules: []infrav1.IngressRule{
					{
						Description:          "test",
						Protocol:             infrav1.SecurityGroupProtocolTCP,
						FromPort:             9345,
						ToPort:               9345,
						NatGatewaysIPsSource: true,
					},
				},
			},
			networkStatus: infrav1.NetworkStatus{
				SecurityGroups: map[infrav1.SecurityGroupRole]infrav1.SecurityGroup{
					infrav1.SecurityGroupControlPlane: {
						ID: "cp-sg-id",
					},
					infrav1.SecurityGroupNode: {
						ID: "node-sg-id",
					},
				},
				NatGatewaysIPs: []string{"test-ip"},
			},
			expectedAdditionalIngressRule: infrav1.IngressRule{
				Description: "test",
				Protocol:    infrav1.SecurityGroupProtocolTCP,
				CidrBlocks:  []string{"test-ip/32"},
				FromPort:    9345,
				ToPort:      9345,
			},
		},
		{
			name: "error if nat gateway IPs cidr as source are specified but not available",
			networkSpec: infrav1.NetworkSpec{
				AdditionalControlPlaneIngressRules: []infrav1.IngressRule{
					{
						Description:          "test",
						Protocol:             infrav1.SecurityGroupProtocolTCP,
						FromPort:             9345,
						ToPort:               9345,
						NatGatewaysIPsSource: true,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cs, err := scope.NewClusterScope(scope.ClusterScopeParams{
				Client: fake.NewClientBuilder().WithScheme(scheme).Build(),
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"},
				},
				AWSCluster: &infrav1.AWSCluster{
					Spec: infrav1.AWSClusterSpec{
						NetworkSpec: tc.networkSpec,
					},
					Status: infrav1.AWSClusterStatus{
						Network: tc.networkStatus,
					},
				},
			})
			if err != nil {
				t.Fatalf("Failed to create test context: %v", err)
			}

			s := NewService(cs, testSecurityGroupRoles)
			rules, err := s.getSecurityGroupIngressRules(infrav1.SecurityGroupControlPlane)
			if err != nil {
				if tc.wantErr {
					return
				}
				t.Fatalf("Failed to lookup controlplane security group ingress rules: %v, wantErr %v", err, tc.wantErr)
			}

			found := false
			for _, r := range rules {
				if r.Description != "test" {
					continue
				}
				found = true

				if r.Protocol != tc.expectedAdditionalIngressRule.Protocol {
					t.Fatalf("Expected protocol %s, got %s", tc.expectedAdditionalIngressRule.Protocol, r.Protocol)
				}

				if r.FromPort != tc.expectedAdditionalIngressRule.FromPort {
					t.Fatalf("Expected from port %d, got %d", tc.expectedAdditionalIngressRule.FromPort, r.FromPort)
				}

				if r.ToPort != tc.expectedAdditionalIngressRule.ToPort {
					t.Fatalf("Expected to port %d, got %d", tc.expectedAdditionalIngressRule.ToPort, r.ToPort)
				}

				if !sets.New[string](tc.expectedAdditionalIngressRule.SourceSecurityGroupIDs...).Equal(sets.New[string](tc.expectedAdditionalIngressRule.SourceSecurityGroupIDs...)) {
					t.Fatalf("Expected source security group IDs %v, got %v", tc.expectedAdditionalIngressRule.SourceSecurityGroupIDs, r.SourceSecurityGroupIDs)
				}
			}

			if !found {
				t.Fatal("Additional ingress rule was not found")
			}
		})
	}
}

func TestAdditionalManagedControlPlaneSecurityGroup(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = ekscontrolplanev1.AddToScheme(scheme)

	testCases := []struct {
		name                          string
		networkSpec                   infrav1.NetworkSpec
		expectedAdditionalIngressRule infrav1.IngressRule
	}{
		{
			name: "default control plane security group is used",
			networkSpec: infrav1.NetworkSpec{
				AdditionalControlPlaneIngressRules: []infrav1.IngressRule{
					{
						Description: "test",
						Protocol:    infrav1.SecurityGroupProtocolTCP,
						FromPort:    9345,
						ToPort:      9345,
					},
				},
			},
			expectedAdditionalIngressRule: infrav1.IngressRule{
				Description:            "test",
				Protocol:               infrav1.SecurityGroupProtocolTCP,
				FromPort:               9345,
				ToPort:                 9345,
				SourceSecurityGroupIDs: []string{"cp-sg-id"},
			},
		},
		{
			name: "don't set source security groups if cidr blocks are set",
			networkSpec: infrav1.NetworkSpec{
				AdditionalControlPlaneIngressRules: []infrav1.IngressRule{
					{
						Description: "test",
						Protocol:    infrav1.SecurityGroupProtocolTCP,
						FromPort:    9345,
						ToPort:      9345,
						CidrBlocks:  []string{"test-cidr-block"},
					},
				},
			},
			expectedAdditionalIngressRule: infrav1.IngressRule{
				Description: "test",
				Protocol:    infrav1.SecurityGroupProtocolTCP,
				FromPort:    9345,
				ToPort:      9345,
				CidrBlocks:  []string{"test-cidr-block"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cs, err := scope.NewManagedControlPlaneScope(scope.ManagedControlPlaneScopeParams{
				Client: fake.NewClientBuilder().WithScheme(scheme).Build(),
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"},
				},
				ControlPlane: &ekscontrolplanev1.AWSManagedControlPlane{
					Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
						NetworkSpec: tc.networkSpec,
					},
					Status: ekscontrolplanev1.AWSManagedControlPlaneStatus{
						Network: infrav1.NetworkStatus{
							SecurityGroups: map[infrav1.SecurityGroupRole]infrav1.SecurityGroup{
								infrav1.SecurityGroupControlPlane: {
									ID: "cp-sg-id",
								},
								infrav1.SecurityGroupNode: {
									ID: "node-sg-id",
								},
							},
						},
					},
				},
			})
			if err != nil {
				t.Fatalf("Failed to create test context: %v", err)
			}

			s := NewService(cs, testSecurityGroupRoles)
			rules, err := s.getSecurityGroupIngressRules(infrav1.SecurityGroupControlPlane)
			if err != nil {
				t.Fatalf("Failed to lookup controlplane security group ingress rules: %v", err)
			}

			found := false
			for _, r := range rules {
				if r.Description == "test" {
					found = true

					if !reflect.DeepEqual(r, tc.expectedAdditionalIngressRule) {
						t.Fatalf("Expected ingress rule %#v, got %#v", tc.expectedAdditionalIngressRule, r)
					}
				}
			}

			if !found {
				t.Fatal("Additional ingress rule was not found")
			}
		})
	}
}

func TestControlPlaneLoadBalancerIngressRules(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = infrav1.AddToScheme(scheme)

	testCases := []struct {
		name                string
		awsCluster          *infrav1.AWSCluster
		expectedIngresRules infrav1.IngressRules
	}{
		{
			name: "when no ingress rules are passed and nat gateway IPs are not available, the default is set",
			awsCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{},
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							CidrBlock: "10.0.0.0/16",
						},
					},
				},
				Status: infrav1.AWSClusterStatus{},
			},
			expectedIngresRules: infrav1.IngressRules{
				infrav1.IngressRule{
					Description: "Kubernetes API",
					Protocol:    infrav1.SecurityGroupProtocolTCP,
					FromPort:    6443,
					ToPort:      6443,
					CidrBlocks:  []string{services.AnyIPv4CidrBlock},
				},
			},
		},
		{
			name: "when no ingress rules are passed and nat gateway IPs are not available, the default for IPv6 is set",
			awsCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{},
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							CidrBlock: "10.0.0.0/16",
							IPv6:      &infrav1.IPv6{},
						},
					},
				},
				Status: infrav1.AWSClusterStatus{},
			},
			expectedIngresRules: infrav1.IngressRules{
				infrav1.IngressRule{
					Description:    "Kubernetes API IPv6",
					Protocol:       infrav1.SecurityGroupProtocolTCP,
					FromPort:       6443,
					ToPort:         6443,
					IPv6CidrBlocks: []string{services.AnyIPv6CidrBlock},
				},
			},
		},
		{
			name: "when no ingress rules are passed, allow the Nat Gateway IPs and default to allow all",
			awsCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{},
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							CidrBlock: "10.0.0.0/16",
						},
					},
				},
				Status: infrav1.AWSClusterStatus{
					Network: infrav1.NetworkStatus{
						NatGatewaysIPs: []string{"1.2.3.4"},
					},
				},
			},
			expectedIngresRules: infrav1.IngressRules{
				infrav1.IngressRule{
					Description: "Kubernetes API",
					Protocol:    infrav1.SecurityGroupProtocolTCP,
					FromPort:    6443,
					ToPort:      6443,
					CidrBlocks:  []string{"1.2.3.4/32"},
				},
				infrav1.IngressRule{
					Description: "Kubernetes API",
					Protocol:    infrav1.SecurityGroupProtocolTCP,
					FromPort:    6443,
					ToPort:      6443,
					CidrBlocks:  []string{services.AnyIPv4CidrBlock},
				},
			},
		},
		{
			name: "defined rules are used",
			awsCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						IngressRules: infrav1.IngressRules{
							{
								Description: "My custom ingress rule",
								Protocol:    infrav1.SecurityGroupProtocolTCP,
								FromPort:    1234,
								ToPort:      1234,
								CidrBlocks:  []string{"172.126.1.1/0"},
							},
						},
					},
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							CidrBlock: "10.0.0.0/16",
						},
					},
				},
				Status: infrav1.AWSClusterStatus{
					Network: infrav1.NetworkStatus{
						NatGatewaysIPs: []string{"1.2.3.4"},
					},
				},
			},
			expectedIngresRules: infrav1.IngressRules{
				infrav1.IngressRule{
					Description: "Kubernetes API",
					Protocol:    infrav1.SecurityGroupProtocolTCP,
					FromPort:    6443,
					ToPort:      6443,
					CidrBlocks:  []string{"1.2.3.4/32"},
				},
				infrav1.IngressRule{
					Description: "My custom ingress rule",
					Protocol:    infrav1.SecurityGroupProtocolTCP,
					FromPort:    1234,
					ToPort:      1234,
					CidrBlocks:  []string{"172.126.1.1/0"},
				},
			},
		},
		{
			name: "when no ingress rules are passed while using internal LB",
			awsCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						Scheme: &infrav1.ELBSchemeInternal,
					},
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							CidrBlock: "10.0.0.0/16",
						},
					},
				},
			},
			expectedIngresRules: infrav1.IngressRules{
				infrav1.IngressRule{
					Description: "Kubernetes API",
					Protocol:    infrav1.SecurityGroupProtocolTCP,
					FromPort:    6443,
					ToPort:      6443,
					CidrBlocks:  []string{"10.0.0.0/16"},
				},
				infrav1.IngressRule{
					Description: "Kubernetes API",
					Protocol:    infrav1.SecurityGroupProtocolTCP,
					FromPort:    6443,
					ToPort:      6443,
					CidrBlocks:  []string{services.AnyIPv4CidrBlock},
				},
			},
		},
		{
			name: "when no ingress rules are passed while using internal LB and IPv6",
			awsCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						Scheme: &infrav1.ELBSchemeInternal,
					},
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							IPv6: &infrav1.IPv6{
								CidrBlock: "10.0.0.0/16",
							},
						},
					},
				},
			},
			expectedIngresRules: infrav1.IngressRules{
				infrav1.IngressRule{
					Description:    "Kubernetes API IPv6",
					Protocol:       infrav1.SecurityGroupProtocolTCP,
					FromPort:       6443,
					ToPort:         6443,
					IPv6CidrBlocks: []string{"10.0.0.0/16"},
				},
				infrav1.IngressRule{
					Description:    "Kubernetes API IPv6",
					Protocol:       infrav1.SecurityGroupProtocolTCP,
					FromPort:       6443,
					ToPort:         6443,
					IPv6CidrBlocks: []string{services.AnyIPv6CidrBlock},
				},
			},
		},
		{
			name: "defined rules are used while using internal LB",
			awsCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						IngressRules: infrav1.IngressRules{
							{
								Description: "My custom ingress rule",
								Protocol:    infrav1.SecurityGroupProtocolTCP,
								FromPort:    1234,
								ToPort:      1234,
								CidrBlocks:  []string{"172.126.1.1/0"},
							},
						},
						Scheme: &infrav1.ELBSchemeInternal,
					},
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							CidrBlock: "10.0.0.0/16",
						},
					},
				},
			},
			expectedIngresRules: infrav1.IngressRules{
				infrav1.IngressRule{
					Description: "Kubernetes API",
					Protocol:    infrav1.SecurityGroupProtocolTCP,
					FromPort:    6443,
					ToPort:      6443,
					CidrBlocks:  []string{"10.0.0.0/16"},
				},
				infrav1.IngressRule{
					Description: "My custom ingress rule",
					Protocol:    infrav1.SecurityGroupProtocolTCP,
					FromPort:    1234,
					ToPort:      1234,
					CidrBlocks:  []string{"172.126.1.1/0"},
				},
			},
		},
		{
			name: "defined rules are used when using internal and external LB",
			awsCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						IngressRules: []infrav1.IngressRule{
							{
								Description: "My custom ingress rule",
								Protocol:    infrav1.SecurityGroupProtocolTCP,
								FromPort:    1234,
								ToPort:      1234,
								CidrBlocks:  []string{"172.126.1.1/0"},
							},
						},
						Scheme: &infrav1.ELBSchemeInternal,
					},
					SecondaryControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						IngressRules: []infrav1.IngressRule{
							{
								Description: "Another custom ingress rule",
								Protocol:    infrav1.SecurityGroupProtocolTCP,
								FromPort:    2345,
								ToPort:      2345,
								CidrBlocks:  []string{"0.0.0.0/0"},
							},
						},
					},
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							CidrBlock: "10.0.0.0/16",
						},
					},
				},
			},
			expectedIngresRules: infrav1.IngressRules{
				infrav1.IngressRule{
					Description: "Kubernetes API",
					Protocol:    infrav1.SecurityGroupProtocolTCP,
					FromPort:    6443,
					ToPort:      6443,
					CidrBlocks:  []string{"10.0.0.0/16"},
				},
				infrav1.IngressRule{
					Description: "My custom ingress rule",
					Protocol:    infrav1.SecurityGroupProtocolTCP,
					FromPort:    1234,
					ToPort:      1234,
					CidrBlocks:  []string{"172.126.1.1/0"},
				},
				infrav1.IngressRule{
					Description: "Another custom ingress rule",
					Protocol:    infrav1.SecurityGroupProtocolTCP,
					FromPort:    2345,
					ToPort:      2345,
					CidrBlocks:  []string{"0.0.0.0/0"},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cs, err := scope.NewClusterScope(scope.ClusterScopeParams{
				Client: fake.NewClientBuilder().WithScheme(scheme).Build(),
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"},
				},
				AWSCluster: tc.awsCluster,
			})
			if err != nil {
				t.Fatalf("Failed to create test context: %v", err)
			}

			s := NewService(cs, testSecurityGroupRoles)
			rules, err := s.getSecurityGroupIngressRules(infrav1.SecurityGroupAPIServerLB)
			if err != nil {
				t.Fatalf("Failed to lookup controlplane load balancer security group ingress rules: %v", err)
			}

			g := NewGomegaWithT(t)
			g.Expect(rules).To(Equal(tc.expectedIngresRules))
		})
	}
}

func TestDeleteSecurityGroups(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name    string
		input   *infrav1.NetworkSpec
		expect  func(m *mocks.MockEC2APIMockRecorder)
		wantErr bool
	}{
		{
			name: "do not delete security groups provided as overrides",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID:                "vpc-securitygroups",
					InternetGatewayID: aws.String("igw-01"),
				},
				Subnets: infrav1.Subnets{
					infrav1.SubnetSpec{
						ID:               "subnet-securitygroups-private",
						IsPublic:         false,
						AvailabilityZone: "us-east-1a",
					},
					infrav1.SubnetSpec{
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
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSecurityGroups(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeSecurityGroupsInput{}), gomock.Any()).Return(&ec2.DescribeSecurityGroupsOutput{}, nil)
			},
		},
		{
			name: "Should skip SG deletion if VPC ID not present",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{},
			},
		},
		{
			name: "Should return error if unable to find cluster-owned security groups in vpc",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{ID: "vpc-id"},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSecurityGroups(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeSecurityGroupsInput{}), gomock.Any()).Return(&ec2.DescribeSecurityGroupsOutput{}, awserrors.NewFailedDependency("dependency-failure"))
			},
			wantErr: true,
		},
		{
			name: "Should return error if unable to describe any SG present in VPC and owned by cluster",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{ID: "vpc-id"},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSecurityGroups(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeSecurityGroupsInput{}), gomock.Any()).Return(nil, awserr.New("dependency-failure", "dependency-failure", errors.Errorf("dependency-failure")))
			},
			wantErr: true,
		},
		{
			name: "Should not revoke Ingress rules for a SG if IP permissions are not set and able to delete the SG",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{ID: "vpc-id"},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSecurityGroups(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeSecurityGroupsInput{}), gomock.Any()).Return(&ec2.DescribeSecurityGroupsOutput{
					SecurityGroups: []types.SecurityGroup{
						{
							GroupId:   aws.String("group-id"),
							GroupName: aws.String("group-name"),
						},
					},
				}, nil).Times(2)

				m.DeleteSecurityGroup(context.TODO(), gomock.AssignableToTypeOf(&ec2.DeleteSecurityGroupInput{}), gomock.Any()).Return(nil, nil)
			},
		},
		{
			name: "Should return error if failed to revoke Ingress rules for a SG",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{ID: "vpc-id"},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSecurityGroups(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeSecurityGroupsInput{}), gomock.Any()).Return(&ec2.DescribeSecurityGroupsOutput{
					SecurityGroups: []types.SecurityGroup{
						{
							GroupId:   aws.String("group-id"),
							GroupName: aws.String("group-name"),
							IpPermissions: []types.IpPermission{
								{
									ToPort: aws.Int32(4),
								},
							},
						},
					},
				}, nil).Times(2)
				m.RevokeSecurityGroupIngress(context.TODO(), gomock.AssignableToTypeOf(&ec2.RevokeSecurityGroupIngressInput{}), gomock.Any()).Return(nil, awserr.New("failure", "failure", errors.Errorf("failure")))
			},
			wantErr: true,
		},
		{
			name: "Should delete SG successfully",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{ID: "vpc-id"},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSecurityGroups(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeSecurityGroupsInput{}), gomock.Any()).Return(&ec2.DescribeSecurityGroupsOutput{
					SecurityGroups: []types.SecurityGroup{
						{
							GroupId:   aws.String("group-id"),
							GroupName: aws.String("group-name"),
							IpPermissions: []types.IpPermission{
								{
									ToPort: aws.Int32(4),
								},
							},
						},
					},
				}, nil).Times(2)
				m.RevokeSecurityGroupIngress(context.TODO(), gomock.AssignableToTypeOf(&ec2.RevokeSecurityGroupIngressInput{})).Return(nil, nil)
				m.DeleteSecurityGroup(context.TODO(), gomock.AssignableToTypeOf(&ec2.DeleteSecurityGroupInput{})).Return(nil, nil)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			ec2Mock := mocks.NewMockEC2API(mockCtrl)

			scheme := runtime.NewScheme()
			g.Expect(infrav1.AddToScheme(scheme)).NotTo(HaveOccurred())

			awsCluster := &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: *tc.input,
				},
			}

			client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(awsCluster).WithStatusSubresource(awsCluster).Build()

			cs, err := scope.NewClusterScope(scope.ClusterScopeParams{
				Client: client,
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"},
				},
				AWSCluster: awsCluster,
			})
			g.Expect(err).NotTo(HaveOccurred())

			if tc.expect != nil {
				tc.expect(ec2Mock.EXPECT())
			}

			s := NewService(cs, testSecurityGroupRoles)
			s.EC2Client = ec2Mock

			err = s.DeleteSecurityGroups()
			if tc.wantErr {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).NotTo(HaveOccurred())
		})
	}
}

func TestIngressRulesFromSDKType(t *testing.T) {
	tests := []struct {
		name     string
		input    types.IpPermission
		expected infrav1.IngressRules
	}{
		{
			name: "two ingress rules",
			input: types.IpPermission{
				IpProtocol: aws.String("tcp"),
				FromPort:   aws.Int32(6443),
				ToPort:     aws.Int32(6443),
				IpRanges: []types.IpRange{
					{
						CidrIp:      aws.String("0.0.0.0/0"),
						Description: aws.String("Kubernetes API"),
					},
					{
						CidrIp:      aws.String("192.168.1.1/32"),
						Description: aws.String("My VPN"),
					},
				},
			},
			expected: infrav1.IngressRules{
				{
					Description: "Kubernetes API",
					Protocol:    "tcp",
					FromPort:    6443,
					ToPort:      6443,
					CidrBlocks:  []string{"0.0.0.0/0"},
				},
				{
					Description: "My VPN",
					Protocol:    "tcp",
					FromPort:    6443,
					ToPort:      6443,
					CidrBlocks:  []string{"192.168.1.1/32"},
				},
			},
		},
		{
			name: "Two group pairs",
			input: types.IpPermission{
				IpProtocol: aws.String("tcp"),
				FromPort:   aws.Int32(10250),
				ToPort:     aws.Int32(10250),
				UserIdGroupPairs: []types.UserIdGroupPair{
					{
						Description: aws.String("Kubelet API"),
						UserId:      aws.String("aws-user-id-1"),
						GroupId:     aws.String("sg-source-1"),
					},
					{
						Description: aws.String("Kubelet API"),
						UserId:      aws.String("aws-user-id-1"),
						GroupId:     aws.String("sg-source-2"),
					},
				},
			},
			expected: infrav1.IngressRules{
				{
					Description:            "Kubelet API",
					Protocol:               "tcp",
					FromPort:               10250,
					ToPort:                 10250,
					SourceSecurityGroupIDs: []string{"sg-source-1"},
				},
				{
					Description:            "Kubelet API",
					Protocol:               "tcp",
					FromPort:               10250,
					ToPort:                 10250,
					SourceSecurityGroupIDs: []string{"sg-source-2"},
				},
			},
		},
		{
			name: "Mix of group pairs and cidr blocks",
			input: types.IpPermission{
				IpProtocol: aws.String("tcp"),
				FromPort:   aws.Int32(22),
				ToPort:     aws.Int32(22),
				IpRanges: []types.IpRange{
					{
						CidrIp:      aws.String("0.0.0.0/0"),
						Description: aws.String("MY-SSH"),
					},
				},
				UserIdGroupPairs: []types.UserIdGroupPair{
					{
						UserId:      aws.String("aws-user-id-1"),
						GroupId:     aws.String("sg-source-1"),
						Description: aws.String("SSH"),
					},
				},
			},
			expected: infrav1.IngressRules{
				{
					Description: "MY-SSH",
					Protocol:    "tcp",
					FromPort:    22,
					ToPort:      22,
					CidrBlocks:  []string{"0.0.0.0/0"},
				},
				{
					Description:            "SSH",
					Protocol:               "tcp",
					FromPort:               22,
					ToPort:                 22,
					SourceSecurityGroupIDs: []string{"sg-source-1"},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGomegaWithT(t)
			output := ingressRulesFromSDKType(tc.input)

			g.Expect(output).To(Equal(tc.expected))
		})
	}
}

func TestExpandIngressRules(t *testing.T) {
	tests := []struct {
		name     string
		input    infrav1.IngressRules
		expected infrav1.IngressRules
	}{
		{
			name: "nothing to expand, nothing to do",
			input: infrav1.IngressRules{
				{
					Description: "SSH",
					Protocol:    infrav1.SecurityGroupProtocolTCP,
					FromPort:    22,
					ToPort:      22,
				},
			},
			expected: infrav1.IngressRules{
				{
					Description: "SSH",
					Protocol:    infrav1.SecurityGroupProtocolTCP,
					FromPort:    22,
					ToPort:      22,
				},
			},
		},
		{
			name: "nothing to expand, security group roles is removed",
			input: infrav1.IngressRules{
				{
					Description: "SSH",
					Protocol:    infrav1.SecurityGroupProtocolTCP,
					FromPort:    22,
					ToPort:      22,
					SourceSecurityGroupRoles: []infrav1.SecurityGroupRole{
						infrav1.SecurityGroupControlPlane,
					},
				},
			},
			expected: infrav1.IngressRules{
				{
					Description: "SSH",
					Protocol:    infrav1.SecurityGroupProtocolTCP,
					FromPort:    22,
					ToPort:      22,
				},
			},
		},
		{
			name: "cidr blocks expand",
			input: infrav1.IngressRules{
				{
					Description:    "SSH",
					Protocol:       infrav1.SecurityGroupProtocolTCP,
					FromPort:       22,
					ToPort:         22,
					CidrBlocks:     []string{"0.0.0.0/0", "1.1.1.1/0"},
					IPv6CidrBlocks: []string{"::/0", "::/1"},
				},
			},
			expected: infrav1.IngressRules{
				{
					Description: "SSH",
					Protocol:    infrav1.SecurityGroupProtocolTCP,
					FromPort:    22,
					ToPort:      22,
					CidrBlocks:  []string{"0.0.0.0/0"},
				},
				{
					Description: "SSH",
					Protocol:    infrav1.SecurityGroupProtocolTCP,
					FromPort:    22,
					ToPort:      22,
					CidrBlocks:  []string{"1.1.1.1/0"},
				},
				{
					Description:    "SSH",
					Protocol:       infrav1.SecurityGroupProtocolTCP,
					FromPort:       22,
					ToPort:         22,
					IPv6CidrBlocks: []string{"::/0"},
				},
				{
					Description:    "SSH",
					Protocol:       infrav1.SecurityGroupProtocolTCP,
					FromPort:       22,
					ToPort:         22,
					IPv6CidrBlocks: []string{"::/1"},
				},
			},
		},
		{
			name: "security group ids expand, security group roles removed",
			input: infrav1.IngressRules{
				{
					Description:            "SSH",
					Protocol:               infrav1.SecurityGroupProtocolTCP,
					FromPort:               22,
					ToPort:                 22,
					SourceSecurityGroupIDs: []string{"sg-1", "sg-2"},
					SourceSecurityGroupRoles: []infrav1.SecurityGroupRole{
						infrav1.SecurityGroupControlPlane,
						infrav1.SecurityGroupNode,
					},
				},
			},
			expected: infrav1.IngressRules{
				{
					Description:            "SSH",
					Protocol:               infrav1.SecurityGroupProtocolTCP,
					FromPort:               22,
					ToPort:                 22,
					SourceSecurityGroupIDs: []string{"sg-1"},
				},
				{
					Description:            "SSH",
					Protocol:               infrav1.SecurityGroupProtocolTCP,
					FromPort:               22,
					ToPort:                 22,
					SourceSecurityGroupIDs: []string{"sg-2"},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGomegaWithT(t)
			output := expandIngressRules(tc.input)

			g.Expect(output).To(Equal(tc.expected))
		})
	}
}

func TestNodePortServicesIngressRules(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = infrav1.AddToScheme(scheme)

	testCases := []struct {
		name                string
		cidrBlocks          []string
		expectedIngresRules infrav1.IngressRules
	}{
		{
			name:       "default node ports services ingress rules, no node port cidr block provided",
			cidrBlocks: nil,
			expectedIngresRules: infrav1.IngressRules{
				{
					Description: "Node Port Services",
					Protocol:    infrav1.SecurityGroupProtocolTCP,
					FromPort:    30000,
					ToPort:      32767,
					CidrBlocks:  []string{services.AnyIPv4CidrBlock},
				},
				{
					Description:            "Kubelet API",
					Protocol:               infrav1.SecurityGroupProtocolTCP,
					FromPort:               10250,
					ToPort:                 10250,
					SourceSecurityGroupIDs: []string{"Id1", "Id2"},
				},
			},
		},
		{
			name:       "node port cidr block provided, no default cidr block used for node port services ingress rule",
			cidrBlocks: []string{"10.0.0.0/16"},
			expectedIngresRules: infrav1.IngressRules{
				{
					Description: "Node Port Services",
					Protocol:    infrav1.SecurityGroupProtocolTCP,
					FromPort:    30000,
					ToPort:      32767,
					CidrBlocks:  []string{"10.0.0.0/16"},
				},
				{
					Description:            "Kubelet API",
					Protocol:               infrav1.SecurityGroupProtocolTCP,
					FromPort:               10250,
					ToPort:                 10250,
					SourceSecurityGroupIDs: []string{"Id1", "Id2"},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cs, err := scope.NewClusterScope(scope.ClusterScopeParams{
				Client: fake.NewClientBuilder().WithScheme(scheme).Build(),
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"},
				},
				AWSCluster: &infrav1.AWSCluster{
					Spec: infrav1.AWSClusterSpec{
						ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{},
						NetworkSpec: infrav1.NetworkSpec{
							VPC: infrav1.VPCSpec{
								CidrBlock: "10.0.0.0/16",
							},
							NodePortIngressRuleCidrBlocks: tc.cidrBlocks,
						},
					},
					Status: infrav1.AWSClusterStatus{
						Network: infrav1.NetworkStatus{
							SecurityGroups: map[infrav1.SecurityGroupRole]infrav1.SecurityGroup{
								infrav1.SecurityGroupControlPlane: {ID: "Id1"},
								infrav1.SecurityGroupNode:         {ID: "Id2"},
							},
						},
					},
				},
			})
			if err != nil {
				t.Fatalf("Failed to create test context: %v", err)
			}

			s := NewService(cs, testSecurityGroupRoles)
			rules, err := s.getSecurityGroupIngressRules(infrav1.SecurityGroupNode)
			if err != nil {
				t.Fatalf("Failed to lookup node security group ingress rules: %v", err)
			}

			g := NewGomegaWithT(t)
			g.Expect(rules).To(Equal(tc.expectedIngresRules))
		})
	}
}
