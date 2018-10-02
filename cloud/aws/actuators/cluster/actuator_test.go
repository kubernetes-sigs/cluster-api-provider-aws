// Copyright © 2018 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cluster_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	providerconfig "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	clientv1 "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1"

	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/actuators/cluster"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/actuators/cluster/mock_clusteriface"
	ec2svc "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/ec2"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/ec2/mock_ec2iface"
)

type clusterGetter struct {
	ci *mock_clusteriface.MockClusterInterface
}

func (c *clusterGetter) Clusters(ns string) clientv1.ClusterInterface {
	return c.ci
}

func TestReconcile(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	cg := &clusterGetter{
		ci: mock_clusteriface.NewMockClusterInterface(mockCtrl),
	}
	me := mock_ec2iface.NewMockEC2API(mockCtrl)
	defer mockCtrl.Finish()

	cg.ci.EXPECT().
		UpdateStatus(gomock.AssignableToTypeOf(&clusterv1.Cluster{})).
		Return(&clusterv1.Cluster{}, nil)

	gomock.InOrder(
		me.EXPECT().
			DescribeVpcs(&ec2.DescribeVpcsInput{
				Filters: []*ec2.Filter{&ec2.Filter{
					Name:   aws.String("tag-key"),
					Values: aws.StringSlice([]string{"kubernetes.io/cluster/test"}),
				}},
			}).
			Return(&ec2.DescribeVpcsOutput{
				Vpcs: []*ec2.Vpc{},
			}, nil),
		me.EXPECT().
			CreateVpc(&ec2.CreateVpcInput{
				CidrBlock: aws.String("10.0.0.0/16"),
			}).
			Return(&ec2.CreateVpcOutput{
				Vpc: &ec2.Vpc{
					VpcId:     aws.String("1234"),
					CidrBlock: aws.String("10.0.0.0/16"),
				},
			}, nil),
		me.EXPECT().
			WaitUntilVpcAvailable(&ec2.DescribeVpcsInput{
				VpcIds: []*string{aws.String("1234")},
			}).
			Return(nil),
		me.EXPECT().
			CreateTags(&ec2.CreateTagsInput{
				Resources: aws.StringSlice([]string{"1234"}),
				Tags: []*ec2.Tag{&ec2.Tag{
					Key:   aws.String("kubernetes.io/cluster/test"),
					Value: aws.String("owned"),
				}},
			}).
			Return(nil, nil),
		me.EXPECT().
			DescribeSubnets(&ec2.DescribeSubnetsInput{
				Filters: []*ec2.Filter{
					&ec2.Filter{
						Name: aws.String("vpc-id"),
						Values: []*string{
							aws.String("1234"),
						},
					},
				},
			}).
			Return(&ec2.DescribeSubnetsOutput{
				Subnets: []*ec2.Subnet{
					&ec2.Subnet{
						SubnetId:            aws.String("snow"),
						VpcId:               aws.String("1234"),
						AvailabilityZone:    aws.String("antarctica"),
						CidrBlock:           aws.String("10.0.0.0/24"),
						MapPublicIpOnLaunch: aws.Bool(false),
					},
					&ec2.Subnet{
						SubnetId:            aws.String("ice"),
						VpcId:               aws.String("1234"),
						AvailabilityZone:    aws.String("antarctica"),
						CidrBlock:           aws.String("10.0.1.0/24"),
						MapPublicIpOnLaunch: aws.Bool(true),
					},
				},
			}, nil),
		me.EXPECT().
			DescribeAvailabilityZones(&ec2.DescribeAvailabilityZonesInput{
				Filters: []*ec2.Filter{
					&ec2.Filter{
						Name:   aws.String("state"),
						Values: []*string{aws.String("available")},
					},
				},
			}).
			Return(&ec2.DescribeAvailabilityZonesOutput{
				AvailabilityZones: []*ec2.AvailabilityZone{
					&ec2.AvailabilityZone{ZoneName: aws.String("antarctica")},
				},
			}, nil),
		me.EXPECT().
			DescribeInternetGateways(&ec2.DescribeInternetGatewaysInput{
				Filters: []*ec2.Filter{
					&ec2.Filter{
						Name:   aws.String("attachment.vpc-id"),
						Values: []*string{aws.String("1234")},
					},
				},
			}).
			Return(&ec2.DescribeInternetGatewaysOutput{
				InternetGateways: []*ec2.InternetGateway{
					&ec2.InternetGateway{
						InternetGatewayId: aws.String("carrot"),
					},
				},
			}, nil),
		me.EXPECT().
			DescribeNatGatewaysPages(gomock.Any(), gomock.Any()).
			Return(nil),
		me.EXPECT().
			AllocateAddress(&ec2.AllocateAddressInput{Domain: aws.String("vpc")}).
			Return(&ec2.AllocateAddressOutput{AllocationId: aws.String("scarf")}, nil),
		me.EXPECT().
			CreateNatGateway(&ec2.CreateNatGatewayInput{
				AllocationId: aws.String("scarf"),
				SubnetId:     aws.String("ice"),
			}).
			Return(&ec2.CreateNatGatewayOutput{
				NatGateway: &ec2.NatGateway{
					NatGatewayId: aws.String("nat-ice1"),
				},
			}, nil),
		me.EXPECT().
			WaitUntilNatGatewayAvailable(&ec2.DescribeNatGatewaysInput{NatGatewayIds: []*string{aws.String("nat-ice1")}}).
			Return(nil),
		me.EXPECT().
			CreateTags(&ec2.CreateTagsInput{
				Resources: aws.StringSlice([]string{"nat-ice1"}),
				Tags: []*ec2.Tag{&ec2.Tag{
					Key:   aws.String("kubernetes.io/cluster/test"),
					Value: aws.String("owned"),
				}},
			}).
			Return(nil, nil),
		me.EXPECT().
			DescribeRouteTables(&ec2.DescribeRouteTablesInput{
				Filters: []*ec2.Filter{
					&ec2.Filter{
						Name: aws.String("vpc-id"),
						Values: []*string{
							aws.String("1234"),
						},
					},
				},
			}).Return(&ec2.DescribeRouteTablesOutput{}, nil),
		me.EXPECT().
			CreateRouteTable(&ec2.CreateRouteTableInput{VpcId: aws.String("1234")}).
			Return(&ec2.CreateRouteTableOutput{RouteTable: &ec2.RouteTable{RouteTableId: aws.String("rt-1")}}, nil),
		me.EXPECT().
			CreateTags(&ec2.CreateTagsInput{
				Resources: aws.StringSlice([]string{"rt-1"}),
				Tags: []*ec2.Tag{&ec2.Tag{
					Key:   aws.String("kubernetes.io/cluster/test"),
					Value: aws.String("owned"),
				}},
			}).
			Return(nil, nil),
		me.EXPECT().
			CreateRoute(&ec2.CreateRouteInput{
				RouteTableId:         aws.String("rt-1"),
				DestinationCidrBlock: aws.String("0.0.0.0/0"),
				NatGatewayId:         aws.String("nat-ice1"),
			}).
			Return(&ec2.CreateRouteOutput{}, nil),
		me.EXPECT().
			AssociateRouteTable(&ec2.AssociateRouteTableInput{RouteTableId: aws.String("rt-1"), SubnetId: aws.String("snow")}).
			Return(&ec2.AssociateRouteTableOutput{}, nil),
		me.EXPECT().
			CreateRouteTable(&ec2.CreateRouteTableInput{VpcId: aws.String("1234")}).
			Return(&ec2.CreateRouteTableOutput{RouteTable: &ec2.RouteTable{RouteTableId: aws.String("rt-2")}}, nil),
		me.EXPECT().
			CreateTags(&ec2.CreateTagsInput{
				Resources: aws.StringSlice([]string{"rt-2"}),
				Tags: []*ec2.Tag{&ec2.Tag{
					Key:   aws.String("kubernetes.io/cluster/test"),
					Value: aws.String("owned"),
				}},
			}).
			Return(nil, nil),
		me.EXPECT().
			CreateRoute(&ec2.CreateRouteInput{
				RouteTableId:         aws.String("rt-2"),
				DestinationCidrBlock: aws.String("0.0.0.0/0"),
				GatewayId:            aws.String("carrot"),
			}).
			Return(&ec2.CreateRouteOutput{}, nil),
		me.EXPECT().
			AssociateRouteTable(&ec2.AssociateRouteTableInput{RouteTableId: aws.String("rt-2"), SubnetId: aws.String("ice")}).
			Return(&ec2.AssociateRouteTableOutput{}, nil),
		me.EXPECT().
			DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
				Filters: []*ec2.Filter{
					{
						Name:   aws.String("vpc-id"),
						Values: []*string{aws.String("1234")},
					},
					{
						Name:   aws.String("tag-key"),
						Values: []*string{aws.String("kubernetes.io/cluster/test")},
					},
				},
			}).
			Return(&ec2.DescribeSecurityGroupsOutput{
				SecurityGroups: []*ec2.SecurityGroup{
					&ec2.SecurityGroup{
						GroupId:   aws.String("sg-cp1"),
						GroupName: aws.String("test-bastion"),
						IpPermissions: []*ec2.IpPermission{
							&ec2.IpPermission{
								FromPort:   aws.Int64(22),
								ToPort:     aws.Int64(22),
								IpProtocol: aws.String("tcp"),
								IpRanges: []*ec2.IpRange{
									&ec2.IpRange{
										CidrIp:      aws.String("0.0.0.0/0"),
										Description: aws.String("SSH"),
									},
								},
							},
						},
					},
					&ec2.SecurityGroup{
						GroupId:   aws.String("sg-cp1"),
						GroupName: aws.String("test-controlplane"),
						IpPermissions: []*ec2.IpPermission{
							&ec2.IpPermission{
								FromPort:   aws.Int64(22),
								ToPort:     aws.Int64(22),
								IpProtocol: aws.String("tcp"),
								IpRanges: []*ec2.IpRange{
									&ec2.IpRange{
										CidrIp:      aws.String("0.0.0.0/0"),
										Description: aws.String("SSH"),
									},
								},
							},
							&ec2.IpPermission{
								FromPort:   aws.Int64(6443),
								ToPort:     aws.Int64(6443),
								IpProtocol: aws.String("tcp"),
								IpRanges: []*ec2.IpRange{
									&ec2.IpRange{
										CidrIp:      aws.String("0.0.0.0/0"),
										Description: aws.String("Kubernetes API"),
									},
								},
							},
							&ec2.IpPermission{
								FromPort:   aws.Int64(2379),
								ToPort:     aws.Int64(2379),
								IpProtocol: aws.String("tcp"),
								UserIdGroupPairs: []*ec2.UserIdGroupPair{
									&ec2.UserIdGroupPair{
										GroupId:     aws.String("sg-cp1"),
										Description: aws.String("etcd"),
									},
								},
							},
							&ec2.IpPermission{
								FromPort:   aws.Int64(2380),
								ToPort:     aws.Int64(2380),
								IpProtocol: aws.String("tcp"),
								UserIdGroupPairs: []*ec2.UserIdGroupPair{
									&ec2.UserIdGroupPair{
										GroupId:     aws.String("sg-cp1"),
										Description: aws.String("etcd peer"),
									},
								},
							},
						},
					},
					&ec2.SecurityGroup{
						GroupId:   aws.String("sg-nd1"),
						GroupName: aws.String("test-node"),
						IpPermissions: []*ec2.IpPermission{
							&ec2.IpPermission{
								FromPort:   aws.Int64(22),
								ToPort:     aws.Int64(22),
								IpProtocol: aws.String("tcp"),
								IpRanges: []*ec2.IpRange{
									&ec2.IpRange{
										CidrIp:      aws.String("0.0.0.0/0"),
										Description: aws.String("SSH"),
									},
								},
							},
							&ec2.IpPermission{
								FromPort:   aws.Int64(30000),
								ToPort:     aws.Int64(32767),
								IpProtocol: aws.String("tcp"),
								IpRanges: []*ec2.IpRange{
									&ec2.IpRange{
										CidrIp:      aws.String("0.0.0.0/0"),
										Description: aws.String("Node Port Services"),
									},
								},
							},
							&ec2.IpPermission{
								FromPort:   aws.Int64(10250),
								ToPort:     aws.Int64(10250),
								IpProtocol: aws.String("tcp"),
								UserIdGroupPairs: []*ec2.UserIdGroupPair{
									&ec2.UserIdGroupPair{
										GroupId:     aws.String("sg-cp1"),
										Description: aws.String("Kubelet API"),
									},
								},
							},
						},
					},
				},
			}, nil),
	)

	c, err := providerconfig.NewCodec()
	if err != nil {
		t.Fatalf("failed to create codec: %v", err)
	}
	ap := cluster.ActuatorParams{
		Codec: c,
		EC2Service: &ec2svc.Service{
			EC2: me,
		},
		ClustersGetter: cg,
	}

	a, err := cluster.NewActuator(ap)
	if err != nil {
		t.Fatalf("could not create an actuator: %v", err)
	}

	if err := a.Reconcile(&clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "test", ClusterName: "test"}}); err != nil {
		t.Fatalf("failed to reconcile cluster: %v", err)
	}
}
