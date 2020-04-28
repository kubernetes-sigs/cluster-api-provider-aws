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

package network

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/ec2/mock_ec2iface"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/elb/mock_elbiface"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
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
				})).
					Return(&ec2.CreateSecurityGroupOutput{GroupId: aws.String("sg-bastion")}, nil)

				m.CreateTags(matchesTags(&ec2.CreateTagsInput{
					Resources: []*string{aws.String("sg-bastion")},
					Tags: []*ec2.Tag{
						{
							Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
							Value: aws.String("owned"),
						},
						{
							Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
							Value: aws.String("bastion"),
						},
						{
							Key:   aws.String("Name"),
							Value: aws.String("test-cluster-bastion"),
						},
					},
				})).
					Return(nil, nil).
					After(securityGroupBastion)

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
				})).
					Return(&ec2.CreateSecurityGroupOutput{GroupId: aws.String("sg-apiserver-lb")}, nil)

				m.CreateTags(matchesTags(&ec2.CreateTagsInput{
					Resources: []*string{aws.String("sg-apiserver-lb")},
					Tags: []*ec2.Tag{
						{
							Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
							Value: aws.String("owned"),
						},
						{
							Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
							Value: aws.String("apiserver-lb"),
						},
						{
							Key:   aws.String("Name"),
							Value: aws.String("test-cluster-apiserver-lb"),
						},
					},
				})).
					Return(nil, nil).
					After(securityGroupAPIServerLb)

				m.AuthorizeSecurityGroupIngress(gomock.AssignableToTypeOf(&ec2.AuthorizeSecurityGroupIngressInput{
					GroupId: aws.String("sg-apiserver-lb"),
				})).
					Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil).
					After(securityGroupAPIServerLb)

				////////////////////////

				securityGroupLb := m.CreateSecurityGroup(gomock.Eq(&ec2.CreateSecurityGroupInput{
					VpcId:       aws.String("vpc-securitygroups"),
					GroupName:   aws.String("test-cluster-lb"),
					Description: aws.String("Kubernetes cluster test-cluster: lb"),
				})).
					Return(&ec2.CreateSecurityGroupOutput{GroupId: aws.String("sg-lb")}, nil)

				m.CreateTags(matchesTags(&ec2.CreateTagsInput{
					Resources: []*string{aws.String("sg-lb")},
					Tags: []*ec2.Tag{
						{
							Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
							Value: aws.String("owned"),
						},
						{
							Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
							Value: aws.String("lb"),
						},
						{
							Key:   aws.String("Name"),
							Value: aws.String("test-cluster-lb"),
						},
						{
							Key:   aws.String("kubernetes.io/cluster/test-cluster"),
							Value: aws.String("owned"),
						},
					},
				})).
					Return(nil, nil).
					After(securityGroupLb)

				////////////////////////

				securityGroupControl := m.CreateSecurityGroup(gomock.Eq(&ec2.CreateSecurityGroupInput{
					VpcId:       aws.String("vpc-securitygroups"),
					GroupName:   aws.String("test-cluster-controlplane"),
					Description: aws.String("Kubernetes cluster test-cluster: controlplane"),
				})).
					Return(&ec2.CreateSecurityGroupOutput{GroupId: aws.String("sg-control")}, nil)

				m.CreateTags(matchesTags(&ec2.CreateTagsInput{
					Resources: []*string{aws.String("sg-control")},
					Tags: []*ec2.Tag{
						{
							Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
							Value: aws.String("owned"),
						}, {
							Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
							Value: aws.String("controlplane"),
						}, {
							Key:   aws.String("Name"),
							Value: aws.String("test-cluster-controlplane"),
						},
					},
				})).
					Return(nil, nil).
					After(securityGroupControl)

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
				})).
					Return(&ec2.CreateSecurityGroupOutput{GroupId: aws.String("sg-node")}, nil)

				m.CreateTags(matchesTags(&ec2.CreateTagsInput{
					Resources: []*string{aws.String("sg-node")},
					Tags: []*ec2.Tag{
						{
							Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
							Value: aws.String("owned"),
						}, {
							Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
							Value: aws.String("node"),
						}, {
							Key:   aws.String("Name"),
							Value: aws.String("test-cluster-node"),
						},
					},
				})).
					Return(nil, nil).
					After(securityGroupNode)

				m.AuthorizeSecurityGroupIngress(gomock.AssignableToTypeOf(&ec2.AuthorizeSecurityGroupIngressInput{
					GroupId: aws.String("sg-node"),
				})).
					Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil).
					After(securityGroupNode)

			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ec2Mock := mock_ec2iface.NewMockEC2API(mockCtrl)
			elbMock := mock_elbiface.NewMockELBAPI(mockCtrl)

			scope, err := scope.NewClusterScope(scope.ClusterScopeParams{
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"},
				},
				AWSClients: scope.AWSClients{
					EC2: ec2Mock,
					ELB: elbMock,
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

			s := NewService(scope.NetworkScope)
			if err := s.reconcileSecurityGroups(); err != nil && tc.err != nil {
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

	s := NewService(scope.NetworkScope)
	rules, err := s.getSecurityGroupIngressRules(infrav1.SecurityGroupControlPlane)
	if err != nil {
		t.Fatalf("Failed to lookup controlplane security group ingress rules: %v", err)
	}

	for _, r := range rules {
		if sets.NewString(r.CidrBlocks...).Has(anyIPv4CidrBlock) {
			t.Fatal("Ingress rule allows any CIDR block")
		}
	}
}

func matchesTags(input *ec2.CreateTagsInput) gomock.Matcher {
	return tagMatcher{input}
}

type tagMatcher struct {
	*ec2.CreateTagsInput
}

func (t tagMatcher) Matches(x interface{}) bool {
	actual, ok := x.(*ec2.CreateTagsInput)
	if !ok {
		return false
	}

	if !reflect.DeepEqual(actual.Resources, t.Resources) {
		return false
	}

	if len(actual.Tags) != len(t.Tags) {
		return false
	}

	for _, at := range actual.Tags {
		match := false
		for _, tt := range t.Tags {
			if reflect.DeepEqual(at, tt) {
				match = true
				break
			}
		}
		if !match {
			return false
		}
	}

	return true
}

func (t tagMatcher) String() string {
	return fmt.Sprintf("matches %v", t.CreateTagsInput)
}
