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

package elb

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/ec2/mock_ec2iface"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/elb/mock_elbiface"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/elb/mock_resourcegroupstaggingapiiface"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestGenerateELBName(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "test",
			expected: "test-apiserver",
		},
		{
			name:     "0123456789012345678901",
			expected: "0123456789012345678901-apiserver",
		},
		{
			name:     "01234567890123456789012",
			expected: "26o3cjil5at5qn27vukn5x09b3ql-k8s",
		},
		{
			name:     "anotherverylongtoolongname",
			expected: "t8gnrbbifaaf5d0k4xmwui3xwvip-k8s",
		},
		{
			name:     "anotherverylongtoolongnameanotherverylongtoolongname",
			expected: "tph1huzox1f10z9ow1inrootjws8-k8s",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elbName, err := GenerateELBName(tt.name)
			if err != nil {
				t.Error(err)
			}

			if elbName != tt.expected {
				t.Errorf("expected ELB name: %v, got name: %v", tt.expected, elbName)
			}

			if len(elbName) > 32 {
				t.Errorf("ELB name too long: %v vs. %s", len(elbName), "32")
			}
		})
	}
}

func TestGetAPIServerClassicELBSpec_ControlPlaneLoadBalancer(t *testing.T) {
	tests := []struct {
		name   string
		lb     *infrav1.AWSLoadBalancerSpec
		mocks  func(m *mock_ec2iface.MockEC2APIMockRecorder)
		expect func(t *testing.T, res *infrav1.ClassicELB)
	}{
		{
			name:  "nil load balancer config",
			lb:    nil,
			mocks: func(m *mock_ec2iface.MockEC2APIMockRecorder) {},
			expect: func(t *testing.T, res *infrav1.ClassicELB) {
				if res.Attributes.CrossZoneLoadBalancing {
					t.Error("Expected load balancer not to have cross-zone load balancing enabled")
				}
			},
		},
		{
			name: "load balancer config with cross zone enabled",
			lb: &infrav1.AWSLoadBalancerSpec{
				CrossZoneLoadBalancing: true,
			},
			mocks: func(m *mock_ec2iface.MockEC2APIMockRecorder) {},
			expect: func(t *testing.T, res *infrav1.ClassicELB) {
				if !res.Attributes.CrossZoneLoadBalancing {
					t.Error("Expected load balancer to have cross-zone load balancing enabled")
				}
			},
		},
		{
			name: "load balancer config with subnets specified",
			lb: &infrav1.AWSLoadBalancerSpec{
				Subnets: []string{"subnet-1", "subnet-2"},
			},
			mocks: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.DescribeSubnets(gomock.Eq(&ec2.DescribeSubnetsInput{
					SubnetIds: []*string{
						aws.String("subnet-1"),
						aws.String("subnet-2"),
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []*ec2.Subnet{
							{
								SubnetId:         aws.String("subnet-1"),
								AvailabilityZone: aws.String("us-east-1a"),
							},
							{
								SubnetId:         aws.String("subnet-2"),
								AvailabilityZone: aws.String("us-east-1b"),
							},
						},
					}, nil)
			},
			expect: func(t *testing.T, res *infrav1.ClassicELB) {
				if len(res.SubnetIDs) != 2 {
					t.Errorf("Expected load balancer to be configured for 2 subnets, got %v", len(res.SubnetIDs))
				}
				if len(res.AvailabilityZones) != 2 {
					t.Errorf("Expected load balancer to be configured for 2 availability zones, got %v", len(res.AvailabilityZones))
				}
			},
		},
		{
			name: "load balancer config with additional security groups specified",
			lb: &infrav1.AWSLoadBalancerSpec{
				AdditionalSecurityGroups: []string{"sg-00001", "sg-00002"},
			},
			mocks: func(m *mock_ec2iface.MockEC2APIMockRecorder) {},
			expect: func(t *testing.T, res *infrav1.ClassicELB) {
				if len(res.SecurityGroupIDs) != 3 {
					t.Errorf("Expected load balancer to be configured for 3 security groups, got %v", len(res.SecurityGroupIDs))
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			ec2Mock := mock_ec2iface.NewMockEC2API(mockCtrl)

			scheme := runtime.NewScheme()
			_ = infrav1.AddToScheme(scheme)
			client := fake.NewClientBuilder().WithScheme(scheme).Build()
			clusterScope, err := scope.NewClusterScope(scope.ClusterScopeParams{
				Client: client,
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "foo",
						Name:      "bar",
					},
				},
				AWSCluster: &infrav1.AWSCluster{
					ObjectMeta: metav1.ObjectMeta{Name: "test"},
					Spec: infrav1.AWSClusterSpec{
						ControlPlaneLoadBalancer: tc.lb,
					},
				},
			})
			if err != nil {
				t.Fatal(err)
			}

			tc.mocks(ec2Mock.EXPECT())

			s := &Service{
				scope:     clusterScope,
				EC2Client: ec2Mock,
			}

			spec, err := s.getAPIServerClassicELBSpec()
			if err != nil {
				t.Fatal(err)
			}

			tc.expect(t, spec)
		})
	}
}

func TestRegisterInstanceWithAPIServerELB(t *testing.T) {
	const (
		clusterSubnetID = "subnet-1"
		elbName         = "bar-apiserver"
		elbSubnetID     = "elb-subnet"
		instanceID      = "test-instance"
		az              = "us-west-1a"
		differentAZ     = "us-east-2c"
	)

	tests := []struct {
		name        string
		awsCluster  *infrav1.AWSCluster
		elbAPIMocks func(m *mock_elbiface.MockELBAPIMockRecorder)
		ec2Mocks    func(m *mock_ec2iface.MockEC2APIMockRecorder)
		check       func(t *testing.T, err error)
	}{
		{
			name: "no load balancer subnets specified",
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: infrav1.Subnets{{
							ID:               clusterSubnetID,
							AvailabilityZone: az,
						}},
					},
				},
			},
			elbAPIMocks: func(m *mock_elbiface.MockELBAPIMockRecorder) {
				m.DescribeLoadBalancers(gomock.Eq(&elb.DescribeLoadBalancersInput{
					LoadBalancerNames: aws.StringSlice([]string{elbName}),
				})).
					Return(&elb.DescribeLoadBalancersOutput{
						LoadBalancerDescriptions: []*elb.LoadBalancerDescription{
							{
								Scheme:  aws.String(string(infrav1.ClassicELBSchemeInternetFacing)),
								Subnets: []*string{aws.String(clusterSubnetID)},
							}},
					}, nil)
				m.DescribeLoadBalancerAttributes(gomock.Eq(&elb.DescribeLoadBalancerAttributesInput{
					LoadBalancerName: aws.String(elbName),
				})).
					Return(&elb.DescribeLoadBalancerAttributesOutput{
						LoadBalancerAttributes: &elb.LoadBalancerAttributes{
							CrossZoneLoadBalancing: &elb.CrossZoneLoadBalancing{
								Enabled: aws.Bool(false),
							},
						},
					}, nil)
				m.RegisterInstancesWithLoadBalancer(gomock.Eq(&elb.RegisterInstancesWithLoadBalancerInput{
					Instances:        []*elb.Instance{{InstanceId: aws.String(instanceID)}},
					LoadBalancerName: aws.String(elbName),
				})).
					Return(&elb.RegisterInstancesWithLoadBalancerOutput{
						Instances: []*elb.Instance{{InstanceId: aws.String(instanceID)}},
					}, nil)
			},
			ec2Mocks: func(m *mock_ec2iface.MockEC2APIMockRecorder) {},
			check: func(t *testing.T, err error) {
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}
			},
		},
		{
			name: "load balancer subnets specified in the same az from the instance",
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: infrav1.Subnets{{
							ID:               clusterSubnetID,
							AvailabilityZone: az,
						}},
					},
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						Subnets: []string{elbSubnetID},
					},
				},
			},
			elbAPIMocks: func(m *mock_elbiface.MockELBAPIMockRecorder) {
				m.DescribeLoadBalancers(gomock.Eq(&elb.DescribeLoadBalancersInput{
					LoadBalancerNames: aws.StringSlice([]string{elbName}),
				})).
					Return(&elb.DescribeLoadBalancersOutput{
						LoadBalancerDescriptions: []*elb.LoadBalancerDescription{
							{
								Scheme:            aws.String(string(infrav1.ClassicELBSchemeInternetFacing)),
								Subnets:           []*string{aws.String(elbSubnetID)},
								AvailabilityZones: []*string{aws.String(az)},
							},
						},
					}, nil)
				m.DescribeLoadBalancerAttributes(gomock.Eq(&elb.DescribeLoadBalancerAttributesInput{
					LoadBalancerName: aws.String(elbName),
				})).
					Return(&elb.DescribeLoadBalancerAttributesOutput{
						LoadBalancerAttributes: &elb.LoadBalancerAttributes{
							CrossZoneLoadBalancing: &elb.CrossZoneLoadBalancing{
								Enabled: aws.Bool(false),
							},
						},
					}, nil)

				m.RegisterInstancesWithLoadBalancer(gomock.Eq(&elb.RegisterInstancesWithLoadBalancerInput{
					Instances:        []*elb.Instance{{InstanceId: aws.String(instanceID)}},
					LoadBalancerName: aws.String(elbName),
				})).
					Return(&elb.RegisterInstancesWithLoadBalancerOutput{
						Instances: []*elb.Instance{{InstanceId: aws.String(instanceID)}},
					}, nil)
			},
			ec2Mocks: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.DescribeSubnets(gomock.Eq(&ec2.DescribeSubnetsInput{
					SubnetIds: []*string{
						aws.String(elbSubnetID),
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []*ec2.Subnet{
							{
								SubnetId:         aws.String(elbSubnetID),
								AvailabilityZone: aws.String(az),
							},
						},
					}, nil)
			},
			check: func(t *testing.T, err error) {
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}
			},
		},
		{
			name: "load balancer subnets specified in a different az from the instance",
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: infrav1.Subnets{{
							ID:               clusterSubnetID,
							AvailabilityZone: az,
						}},
					},
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						Subnets: []string{elbSubnetID},
					},
				},
			},
			elbAPIMocks: func(m *mock_elbiface.MockELBAPIMockRecorder) {
				m.DescribeLoadBalancers(gomock.Eq(&elb.DescribeLoadBalancersInput{
					LoadBalancerNames: aws.StringSlice([]string{elbName}),
				})).
					Return(&elb.DescribeLoadBalancersOutput{
						LoadBalancerDescriptions: []*elb.LoadBalancerDescription{
							{
								Scheme:            aws.String(string(infrav1.ClassicELBSchemeInternetFacing)),
								Subnets:           []*string{aws.String(elbSubnetID)},
								AvailabilityZones: []*string{aws.String(differentAZ)},
							},
						},
					}, nil)
				m.DescribeLoadBalancerAttributes(gomock.Eq(&elb.DescribeLoadBalancerAttributesInput{
					LoadBalancerName: aws.String(elbName),
				})).
					Return(&elb.DescribeLoadBalancerAttributesOutput{
						LoadBalancerAttributes: &elb.LoadBalancerAttributes{
							CrossZoneLoadBalancing: &elb.CrossZoneLoadBalancing{
								Enabled: aws.Bool(false),
							},
						},
					}, nil)
			},
			ec2Mocks: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.DescribeSubnets(gomock.Eq(&ec2.DescribeSubnetsInput{
					SubnetIds: []*string{
						aws.String(elbSubnetID),
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []*ec2.Subnet{
							{
								SubnetId:         aws.String(elbSubnetID),
								AvailabilityZone: aws.String(differentAZ),
							},
						},
					}, nil)
			},
			check: func(t *testing.T, err error) {
				expectedErrMsg := "failed to register instance with APIServer ELB \"bar-apiserver\": instance is in availability zone \"us-west-1a\", no public subnets attached to the ELB in the same zone"
				if err == nil {
					t.Fatalf("Expected error, but got nil")
				}

				if !strings.Contains(err.Error(), expectedErrMsg) {
					t.Fatalf("Expected error: %s\nInstead got: %s", expectedErrMsg, err.Error())
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			elbAPIMocks := mock_elbiface.NewMockELBAPI(mockCtrl)
			ec2Mock := mock_ec2iface.NewMockEC2API(mockCtrl)

			scheme, err := setupScheme()
			if err != nil {
				t.Fatal(err)
			}

			client := fake.NewClientBuilder().WithScheme(scheme).Build()
			clusterScope, err := scope.NewClusterScope(scope.ClusterScopeParams{
				Client: client,
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "foo",
						Name:      "bar",
					},
				},
				AWSCluster: tc.awsCluster,
			})
			if err != nil {
				t.Fatal(err)
			}

			instance := &infrav1.Instance{
				ID:       instanceID,
				SubnetID: clusterSubnetID,
			}

			tc.elbAPIMocks(elbAPIMocks.EXPECT())
			tc.ec2Mocks(ec2Mock.EXPECT())

			s := &Service{
				scope:     clusterScope,
				EC2Client: ec2Mock,
				ELBClient: elbAPIMocks,
			}

			err = s.RegisterInstanceWithAPIServerELB(instance)
			tc.check(t, err)
		})
	}
}

func TestDeleteLoadbalancers(t *testing.T) {
	clusterName := "bar"
	tests := []struct {
		name                  string
		rgAPIMocks            func(m *mock_resourcegroupstaggingapiiface.MockResourceGroupsTaggingAPIAPIMockRecorder)
		elbAPIMocks           func(m *mock_elbiface.MockELBAPIMockRecorder)
		postDeleteElbAPIMocks func(m *mock_elbiface.MockELBAPIMockRecorder)
	}{
		{
			name: "deletes ELBs successfully",
			rgAPIMocks: func(m *mock_resourcegroupstaggingapiiface.MockResourceGroupsTaggingAPIAPIMockRecorder) {
				m.GetResourcesPages(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			},
			elbAPIMocks: func(m *mock_elbiface.MockELBAPIMockRecorder) {
				m.DeleteLoadBalancer(gomock.Eq(&elb.DeleteLoadBalancerInput{LoadBalancerName: aws.String("bar-apiserver")})).Return(nil, nil)
			},
			postDeleteElbAPIMocks: func(m *mock_elbiface.MockELBAPIMockRecorder) {
				m.DescribeLoadBalancers(gomock.Eq(&elb.DescribeLoadBalancersInput{
					LoadBalancerNames: aws.StringSlice([]string{"bar-apiserver"}),
				})).Return(nil, awserr.New(elb.ErrCodeAccessPointNotFoundException, "", nil))
			},
		},
		{
			name: "successful delete. falls back to listing all ELBs when listing by tag fails",
			rgAPIMocks: func(m *mock_resourcegroupstaggingapiiface.MockResourceGroupsTaggingAPIAPIMockRecorder) {
				m.GetResourcesPages(gomock.Any(), gomock.Any()).Return(errors.Errorf("connection failure")).AnyTimes()
			},
			elbAPIMocks: func(m *mock_elbiface.MockELBAPIMockRecorder) {
				m.DescribeLoadBalancersPages(gomock.Any(), gomock.Any()).Do(func(_, y interface{}) {
					funct := y.(func(output *elb.DescribeLoadBalancersOutput, lastPage bool) bool)
					funct(&elb.DescribeLoadBalancersOutput{
						LoadBalancerDescriptions: []*elb.LoadBalancerDescription{
							{
								LoadBalancerName: aws.String("lb-service-name"),
							},
							{
								LoadBalancerName: aws.String("another-service-not-owned"),
							},
							{
								LoadBalancerName: aws.String("service-without-tags"),
							},
						},
					}, true)
				}).Return(nil)
				m.DescribeTags(&elb.DescribeTagsInput{LoadBalancerNames: []*string{aws.String("lb-service-name"), aws.String("another-service-not-owned"), aws.String("service-without-tags")}}).Return(&elb.DescribeTagsOutput{
					TagDescriptions: []*elb.TagDescription{
						{
							LoadBalancerName: aws.String("lb-service-name"),
							Tags: []*elb.Tag{{
								Key:   aws.String(infrav1.ClusterAWSCloudProviderTagKey(clusterName)),
								Value: aws.String(string(infrav1.ResourceLifecycleOwned)),
							}},
						},
						{
							LoadBalancerName: aws.String("another-service-not-owned"),
							Tags: []*elb.Tag{{
								Key:   aws.String("some-tag-key"),
								Value: aws.String("some-tag-value"),
							}},
						},
						{
							LoadBalancerName: aws.String("service-without-tags"),
							Tags:             []*elb.Tag{},
						},
					},
				}, nil)
				m.DeleteLoadBalancer(gomock.Eq(&elb.DeleteLoadBalancerInput{LoadBalancerName: aws.String("bar-apiserver")})).Return(nil, nil)
				m.DeleteLoadBalancer(gomock.Eq(&elb.DeleteLoadBalancerInput{LoadBalancerName: aws.String("lb-service-name")})).Return(nil, nil)
			},
			postDeleteElbAPIMocks: func(m *mock_elbiface.MockELBAPIMockRecorder) {
				m.DescribeLoadBalancers(gomock.Eq(&elb.DescribeLoadBalancersInput{
					LoadBalancerNames: aws.StringSlice([]string{"bar-apiserver"}),
				})).Return(nil, awserr.New(elb.ErrCodeAccessPointNotFoundException, "", nil))
				m.DescribeLoadBalancersPages(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			rgapiMock := mock_resourcegroupstaggingapiiface.NewMockResourceGroupsTaggingAPIAPI(mockCtrl)
			elbapiMock := mock_elbiface.NewMockELBAPI(mockCtrl)

			scheme, err := setupScheme()
			if err != nil {
				t.Fatal(err)
			}

			awsCluster := &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec:       infrav1.AWSClusterSpec{},
			}

			client := fake.NewClientBuilder().WithScheme(scheme).Build()
			ctx := context.TODO()
			client.Create(ctx, awsCluster)

			clusterScope, err := scope.NewClusterScope(scope.ClusterScopeParams{
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "foo",
						Name:      clusterName,
					},
				},
				AWSCluster: awsCluster,
				Client:     client,
			})
			if err != nil {
				t.Fatal(err)
			}

			tc.rgAPIMocks(rgapiMock.EXPECT())
			tc.elbAPIMocks(elbapiMock.EXPECT())
			tc.postDeleteElbAPIMocks(elbapiMock.EXPECT())

			s := &Service{
				scope:                 clusterScope,
				ResourceTaggingClient: rgapiMock,
				ELBClient:             elbapiMock,
			}

			err = s.DeleteLoadbalancers()
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestDescribeLoadbalancers(t *testing.T) {
	clusterName := "bar"
	tests := []struct {
		name                string
		lbName              string
		rgAPIMocks          func(m *mock_resourcegroupstaggingapiiface.MockResourceGroupsTaggingAPIAPIMockRecorder)
		DescribeElbAPIMocks func(m *mock_elbiface.MockELBAPIMockRecorder)
	}{
		{
			name:   "Error if existing loadbalancer with same name doesn't have same scheme",
			lbName: "bar-apiserver",
			rgAPIMocks: func(m *mock_resourcegroupstaggingapiiface.MockResourceGroupsTaggingAPIAPIMockRecorder) {
				m.GetResourcesPages(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			},
			DescribeElbAPIMocks: func(m *mock_elbiface.MockELBAPIMockRecorder) {
				m.DescribeLoadBalancers(gomock.Eq(&elb.DescribeLoadBalancersInput{
					LoadBalancerNames: aws.StringSlice([]string{"bar-apiserver"}),
				})).Return(&elb.DescribeLoadBalancersOutput{LoadBalancerDescriptions: []*elb.LoadBalancerDescription{{Scheme: pointer.StringPtr(string(infrav1.ClassicELBSchemeInternal))}}}, nil)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			rgapiMock := mock_resourcegroupstaggingapiiface.NewMockResourceGroupsTaggingAPIAPI(mockCtrl)
			elbapiMock := mock_elbiface.NewMockELBAPI(mockCtrl)

			scheme, err := setupScheme()
			if err != nil {
				t.Fatal(err)
			}
			awsCluster := &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
					Scheme: &infrav1.ClassicELBSchemeInternetFacing,
				}},
			}

			client := fake.NewClientBuilder().WithScheme(scheme).Build()
			ctx := context.TODO()
			client.Create(ctx, awsCluster)

			clusterScope, err := scope.NewClusterScope(scope.ClusterScopeParams{
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "foo",
						Name:      clusterName,
					},
				},
				AWSCluster: awsCluster,
				Client:     client,
			})
			if err != nil {
				t.Fatal(err)
			}

			tc.rgAPIMocks(rgapiMock.EXPECT())
			tc.DescribeElbAPIMocks(elbapiMock.EXPECT())

			s := &Service{
				scope:                 clusterScope,
				ResourceTaggingClient: rgapiMock,
				ELBClient:             elbapiMock,
			}

			_, err = s.describeClassicELB(tc.lbName)
			if err == nil {
				t.Fatal(err)
			}
		})
	}
}

func TestChunkELBs(t *testing.T) {
	base := "loadbalancer"
	var names []string
	for i := 0; i < 25; i++ {
		names = append(names, fmt.Sprintf("%s+%d", base, i))
	}
	tests := []struct {
		testName              string
		names                 []string
		expectedChunkArrayLen int
	}{
		{
			testName:              "When the user has more the 20 ELBs",
			names:                 names,
			expectedChunkArrayLen: 2,
		}, {
			testName:              "When the user has less than 20 ELBs",
			names:                 []string{"loadBalancer-00"},
			expectedChunkArrayLen: 1,
		},
	}
	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			ans := chunkELBs(tc.names)
			if len(ans) != tc.expectedChunkArrayLen {
				t.Errorf("got %d, want %d", len(ans), tc.expectedChunkArrayLen)
			}
		})
	}
}

func setupScheme() (*runtime.Scheme, error) {
	scheme := runtime.NewScheme()
	if err := clusterv1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := infrav1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	return scheme, nil
}
