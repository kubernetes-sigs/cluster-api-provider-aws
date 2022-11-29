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
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
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

				lbSecurityGroup := m.CreateSecurityGroup(gomock.Eq(&ec2.CreateSecurityGroupInput{
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
				})).Return(&ec2.CreateSecurityGroupOutput{GroupId: aws.String("sg-lb")}, nil)

				m.AuthorizeSecurityGroupIngress(gomock.AssignableToTypeOf(&ec2.AuthorizeSecurityGroupIngressInput{
					GroupId: aws.String("sg-lb"),
				})).
					Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil).
					After(lbSecurityGroup)

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
				m.DescribeSecurityGroupsPages(gomock.AssignableToTypeOf(&ec2.DescribeSecurityGroupsInput{}), gomock.Any()).Return(nil)
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
				m.DescribeSecurityGroupsPages(gomock.AssignableToTypeOf(&ec2.DescribeSecurityGroupsInput{}), gomock.Any()).Return(awserrors.NewFailedDependency("dependency-failure"))
			},
			wantErr: true,
		},
		{
			name: "Should return error if unable to describe any SG present in VPC and owned by cluster",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{ID: "vpc-id"},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSecurityGroupsPages(gomock.AssignableToTypeOf(&ec2.DescribeSecurityGroupsInput{}), gomock.Any()).
					Do(processSecurityGroupsPage).Return(nil)
				m.DescribeSecurityGroups(gomock.AssignableToTypeOf(&ec2.DescribeSecurityGroupsInput{})).Return(nil, awserr.New("dependency-failure", "dependency-failure", errors.Errorf("dependency-failure")))
			},
			wantErr: true,
		},
		{
			name: "Should not revoke Ingress rules for a SG if IP permissions are not set and able to delete the SG",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{ID: "vpc-id"},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSecurityGroupsPages(gomock.AssignableToTypeOf(&ec2.DescribeSecurityGroupsInput{}), gomock.Any()).
					Do(processSecurityGroupsPage).Return(nil)
				m.DescribeSecurityGroups(gomock.AssignableToTypeOf(&ec2.DescribeSecurityGroupsInput{})).Return(&ec2.DescribeSecurityGroupsOutput{
					SecurityGroups: []*ec2.SecurityGroup{
						{
							GroupId:   aws.String("group-id"),
							GroupName: aws.String("group-name"),
						},
					},
				}, nil)
				m.DeleteSecurityGroup(gomock.AssignableToTypeOf(&ec2.DeleteSecurityGroupInput{})).Return(nil, nil)
			},
		},
		{
			name: "Should return error if failed to revoke Ingress rules for a SG",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{ID: "vpc-id"},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSecurityGroupsPages(gomock.AssignableToTypeOf(&ec2.DescribeSecurityGroupsInput{}), gomock.Any()).
					Do(processSecurityGroupsPage).Return(nil)
				m.DescribeSecurityGroups(gomock.AssignableToTypeOf(&ec2.DescribeSecurityGroupsInput{})).Return(&ec2.DescribeSecurityGroupsOutput{
					SecurityGroups: []*ec2.SecurityGroup{
						{
							GroupId:   aws.String("group-id"),
							GroupName: aws.String("group-name"),
							IpPermissions: []*ec2.IpPermission{
								{
									ToPort: aws.Int64(4),
								},
							},
						},
					},
				}, nil)
				m.RevokeSecurityGroupIngress(gomock.AssignableToTypeOf(&ec2.RevokeSecurityGroupIngressInput{})).Return(nil, awserr.New("failure", "failure", errors.Errorf("failure")))
			},
			wantErr: true,
		},
		{
			name: "Should delete SG successfully",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{ID: "vpc-id"},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSecurityGroupsPages(gomock.AssignableToTypeOf(&ec2.DescribeSecurityGroupsInput{}), gomock.Any()).
					Do(processSecurityGroupsPage).Return(nil)
				m.DescribeSecurityGroups(gomock.AssignableToTypeOf(&ec2.DescribeSecurityGroupsInput{})).Return(&ec2.DescribeSecurityGroupsOutput{
					SecurityGroups: []*ec2.SecurityGroup{
						{
							GroupId:   aws.String("group-id"),
							GroupName: aws.String("group-name"),
							IpPermissions: []*ec2.IpPermission{
								{
									ToPort: aws.Int64(4),
								},
							},
						},
					},
				}, nil)
				m.RevokeSecurityGroupIngress(gomock.AssignableToTypeOf(&ec2.RevokeSecurityGroupIngressInput{})).Return(nil, nil)
				m.DeleteSecurityGroup(gomock.AssignableToTypeOf(&ec2.DeleteSecurityGroupInput{})).Return(nil, nil)
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
				TypeMeta: metav1.TypeMeta{
					APIVersion: infrav1.GroupVersion.String(),
					Kind:       "AWSCluster",
				},
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: *tc.input,
				},
			}

			client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(awsCluster).Build()

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
		input    *ec2.IpPermission
		expected infrav1.IngressRules
	}{
		{
			name: "Two group pairs",
			input: &ec2.IpPermission{
				IpProtocol: aws.String("tcp"),
				FromPort:   aws.Int64(10250),
				ToPort:     aws.Int64(10250),
				UserIdGroupPairs: []*ec2.UserIdGroupPair{
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
					SourceSecurityGroupIDs: []string{"sg-source-1", "sg-source-2"},
				},
			},
		},
		{
			name: "Mix of group pairs and cidr blocks",
			input: &ec2.IpPermission{
				IpProtocol: aws.String("tcp"),
				FromPort:   aws.Int64(22),
				ToPort:     aws.Int64(22),
				IpRanges: []*ec2.IpRange{
					{
						CidrIp:      aws.String("0.0.0.0/0"),
						Description: aws.String("MY-SSH"),
					},
				},
				UserIdGroupPairs: []*ec2.UserIdGroupPair{
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

var processSecurityGroupsPage = func(_, y interface{}) {
	funcType := y.(func(out *ec2.DescribeSecurityGroupsOutput, last bool) bool)
	funcType(&ec2.DescribeSecurityGroupsOutput{
		SecurityGroups: []*ec2.SecurityGroup{
			{
				GroupId:   aws.String("group-id"),
				GroupName: aws.String("group-name"),
			},
		},
	}, true)
}
