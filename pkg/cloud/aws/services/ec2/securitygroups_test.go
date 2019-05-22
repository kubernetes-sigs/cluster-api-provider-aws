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

package ec2

import (
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/ec2/mock_ec2iface"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/elb/mock_elbiface"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/tags"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

func TestReconcileSecurityGroups(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name   string
		input  *v1alpha1.NetworkSpec
		expect func(m *mock_ec2iface.MockEC2APIMockRecorder)
		err    error
	}{
		{
			name: "no existing",
			input: &v1alpha1.NetworkSpec{
				VPC: v1alpha1.VPCSpec{
					ID:                "vpc-securitygroups",
					InternetGatewayID: aws.String("igw-01"),
					Tags: tags.Map{
						tags.ClusterKey("test-cluster"): "owned",
					},
				},
				Subnets: v1alpha1.Subnets{
					&v1alpha1.SubnetSpec{
						ID:               "subnet-securitygroups-private",
						IsPublic:         false,
						AvailabilityZone: "us-east-1a",
					},
					&v1alpha1.SubnetSpec{
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

				m.CreateTags(gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
					Return(nil, nil)

				m.AuthorizeSecurityGroupIngress(gomock.AssignableToTypeOf(&ec2.AuthorizeSecurityGroupIngressInput{
					GroupId: aws.String("sg-bastion"),
				})).
					Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil).
					After(securityGroupBastion)

				securityGroupControl := m.CreateSecurityGroup(gomock.Eq(&ec2.CreateSecurityGroupInput{
					VpcId:       aws.String("vpc-securitygroups"),
					GroupName:   aws.String("test-cluster-controlplane"),
					Description: aws.String("Kubernetes cluster test-cluster: controlplane"),
				})).
					Return(&ec2.CreateSecurityGroupOutput{GroupId: aws.String("sg-control")}, nil)

				m.CreateTags(gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
					Return(nil, nil)

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

				m.CreateTags(gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
					Return(nil, nil)

				m.AuthorizeSecurityGroupIngress(gomock.AssignableToTypeOf(&ec2.AuthorizeSecurityGroupIngressInput{
					GroupId: aws.String("sg-node"),
				})).
					Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil).
					After(securityGroupNode)

				////////////////////////

				m.CreateSecurityGroup(gomock.Eq(&ec2.CreateSecurityGroupInput{
					VpcId:       aws.String("vpc-securitygroups"),
					GroupName:   aws.String("test-cluster-lb"),
					Description: aws.String("Kubernetes cluster test-cluster: lb"),
				})).
					Return(&ec2.CreateSecurityGroupOutput{GroupId: aws.String("sg-lb")}, nil)

				m.CreateTags(gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
					Return(nil, nil)

			},
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
