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
	"github.com/aws/aws-sdk-go/service/elbv2"
	rgapi "github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
)

func TestELBName(t *testing.T) {
	tests := []struct {
		name       string
		awsCluster infrav1.AWSCluster
		expected   string
	}{
		{
			name: "name is not defined by user, so generate the default",
			awsCluster: infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "example",
					Namespace: metav1.NamespaceDefault,
				},
			},
			expected: "example-apiserver",
		},
		{
			name: "name is defined by user, so use it",
			awsCluster: infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "example",
					Namespace: metav1.NamespaceDefault,
				},
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						Name: pointer.String("myapiserver"),
					},
				},
			},
			expected: "myapiserver",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheme := runtime.NewScheme()
			_ = infrav1.AddToScheme(scheme)
			client := fake.NewClientBuilder().WithScheme(scheme).Build()

			scope, err := scope.NewClusterScope(scope.ClusterScopeParams{
				Client: client,
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Name:      tt.awsCluster.Name,
						Namespace: tt.awsCluster.Namespace,
					},
				},
				AWSCluster: &tt.awsCluster,
			})
			if err != nil {
				t.Fatalf("failed to create scope: %s", err)
			}

			elbName, err := ELBName(scope)
			if err != nil {
				t.Fatalf("unable to get ELB name: %v", err)
			}
			if elbName != tt.expected {
				t.Fatalf("expected ELB name: %v, got name: %v", tt.expected, elbName)
			}
		})
	}
}

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

func TestGetAPIServerClassicELBSpecControlPlaneLoadBalancer(t *testing.T) {
	tests := []struct {
		name   string
		lb     *infrav1.AWSLoadBalancerSpec
		mocks  func(m *mocks.MockEC2APIMockRecorder)
		expect func(t *testing.T, g *WithT, res *infrav1.LoadBalancer)
	}{
		{
			name:  "nil load balancer config",
			lb:    nil,
			mocks: func(m *mocks.MockEC2APIMockRecorder) {},
			expect: func(t *testing.T, g *WithT, res *infrav1.LoadBalancer) {
				t.Helper()
				if res.ClassicElbAttributes.CrossZoneLoadBalancing {
					t.Error("Expected load balancer not to have cross-zone load balancing enabled")
				}
			},
		},
		{
			name: "load balancer config with cross zone enabled",
			lb: &infrav1.AWSLoadBalancerSpec{
				CrossZoneLoadBalancing: true,
			},
			mocks: func(m *mocks.MockEC2APIMockRecorder) {},
			expect: func(t *testing.T, g *WithT, res *infrav1.LoadBalancer) {
				t.Helper()
				if !res.ClassicElbAttributes.CrossZoneLoadBalancing {
					t.Error("Expected load balancer to have cross-zone load balancing enabled")
				}
			},
		},
		{
			name: "load balancer config with subnets specified",
			lb: &infrav1.AWSLoadBalancerSpec{
				Subnets: []string{"subnet-1", "subnet-2"},
			},
			mocks: func(m *mocks.MockEC2APIMockRecorder) {
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
			expect: func(t *testing.T, g *WithT, res *infrav1.LoadBalancer) {
				t.Helper()
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
			mocks: func(m *mocks.MockEC2APIMockRecorder) {},
			expect: func(t *testing.T, g *WithT, res *infrav1.LoadBalancer) {
				t.Helper()
				if len(res.SecurityGroupIDs) != 3 {
					t.Errorf("Expected load balancer to be configured for 3 security groups, got %v", len(res.SecurityGroupIDs))
				}
			},
		},
		{
			name: "Should create load balancer spec if elb health check protocol specified in config",
			lb: &infrav1.AWSLoadBalancerSpec{
				HealthCheckProtocol: &infrav1.ELBProtocolTCP,
			},
			mocks: func(m *mocks.MockEC2APIMockRecorder) {},
			expect: func(t *testing.T, g *WithT, res *infrav1.LoadBalancer) {
				t.Helper()
				expectedTarget := fmt.Sprintf("%v:%d", infrav1.ELBProtocolTCP, infrav1.DefaultAPIServerPort)
				g.Expect(expectedTarget, res.HealthCheck.Target)
			},
		},
		{
			name:  "Should create load balancer spec with default elb health check protocol",
			lb:    &infrav1.AWSLoadBalancerSpec{},
			mocks: func(m *mocks.MockEC2APIMockRecorder) {},
			expect: func(t *testing.T, g *WithT, res *infrav1.LoadBalancer) {
				t.Helper()
				expectedTarget := fmt.Sprintf("%v:%d", infrav1.ELBProtocolTCP, infrav1.DefaultAPIServerPort)
				g.Expect(expectedTarget, res.HealthCheck.Target)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			ec2Mock := mocks.NewMockEC2API(mockCtrl)

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

			spec, err := s.getAPIServerClassicELBSpec(clusterScope.Name())
			if err != nil {
				t.Fatal(err)
			}

			tc.expect(t, g, spec)
		})
	}
}

func TestGetAPIServerV2ELBSpecControlPlaneLoadBalancer(t *testing.T) {
	tests := []struct {
		name   string
		lb     *infrav1.AWSLoadBalancerSpec
		mocks  func(m *mocks.MockEC2APIMockRecorder)
		expect func(t *testing.T, g *WithT, res *infrav1.LoadBalancer)
	}{
		{
			name:  "nil load balancer config",
			lb:    nil,
			mocks: func(m *mocks.MockEC2APIMockRecorder) {},
			expect: func(t *testing.T, g *WithT, res *infrav1.LoadBalancer) {
				t.Helper()
				if _, ok := res.ELBAttributes["load_balancing.cross_zone.enabled"]; ok {
					t.Error("Expected load balancer not to have cross-zone load balancing enabled")
				}
			},
		},
		{
			name: "load balancer config with cross zone enabled",
			lb: &infrav1.AWSLoadBalancerSpec{
				CrossZoneLoadBalancing: true,
			},
			mocks: func(m *mocks.MockEC2APIMockRecorder) {},
			expect: func(t *testing.T, g *WithT, res *infrav1.LoadBalancer) {
				t.Helper()
				if _, ok := res.ELBAttributes["load_balancing.cross_zone.enabled"]; !ok {
					t.Error("Expected load balancer to have cross-zone load balancing enabled")
				}
			},
		},
		{
			name: "load balancer config with subnets specified",
			lb: &infrav1.AWSLoadBalancerSpec{
				Subnets: []string{"subnet-1", "subnet-2"},
			},
			mocks: func(m *mocks.MockEC2APIMockRecorder) {
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
			expect: func(t *testing.T, g *WithT, res *infrav1.LoadBalancer) {
				t.Helper()
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
				LoadBalancerType:         infrav1.LoadBalancerTypeALB,
			},
			mocks: func(m *mocks.MockEC2APIMockRecorder) {},
			expect: func(t *testing.T, g *WithT, res *infrav1.LoadBalancer) {
				t.Helper()
				if len(res.SecurityGroupIDs) != 3 {
					t.Errorf("Expected load balancer to be configured for 3 security groups, got %v", len(res.SecurityGroupIDs))
				}
			},
		},
		{
			name: "load balancer config with additional security groups specified for NLB has no security groups",
			lb: &infrav1.AWSLoadBalancerSpec{
				AdditionalSecurityGroups: []string{"sg-00001", "sg-00002"},
				LoadBalancerType:         infrav1.LoadBalancerTypeNLB,
			},
			mocks: func(m *mocks.MockEC2APIMockRecorder) {},
			expect: func(t *testing.T, g *WithT, res *infrav1.LoadBalancer) {
				t.Helper()
				if len(res.SecurityGroupIDs) != 0 {
					t.Errorf("Expected load balancer to be configured for 0 security groups, got %v", len(res.SecurityGroupIDs))
				}
			},
		},
		{
			name: "A base listener is set up for NLB",
			lb: &infrav1.AWSLoadBalancerSpec{
				LoadBalancerType: infrav1.LoadBalancerTypeNLB,
			},
			mocks: func(m *mocks.MockEC2APIMockRecorder) {},
			expect: func(t *testing.T, g *WithT, res *infrav1.LoadBalancer) {
				t.Helper()
				if len(res.ELBListeners) != 1 {
					t.Errorf("Expected 1 listener to be configured by default, got %v listener(s)", len(res.ELBListeners))
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			ec2Mock := mocks.NewMockEC2API(mockCtrl)

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

			spec, err := s.getAPIServerLBSpec(clusterScope.Name())
			if err != nil {
				t.Fatal(err)
			}

			tc.expect(t, g, spec)
		})
	}
}

func TestRegisterInstanceWithAPIServerELB(t *testing.T) {
	const (
		namespace       = "foo"
		clusterName     = "bar"
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
		elbAPIMocks func(m *mocks.MockELBAPIMockRecorder)
		ec2Mocks    func(m *mocks.MockEC2APIMockRecorder)
		check       func(t *testing.T, err error)
	}{
		{
			name: "no load balancer subnets specified",
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: clusterName},
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						Name: aws.String(elbName),
					},
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: infrav1.Subnets{{
							ID:               clusterSubnetID,
							AvailabilityZone: az,
						}},
					},
				},
			},
			elbAPIMocks: func(m *mocks.MockELBAPIMockRecorder) {
				m.DescribeLoadBalancers(gomock.Eq(&elb.DescribeLoadBalancersInput{
					LoadBalancerNames: aws.StringSlice([]string{elbName}),
				})).
					Return(&elb.DescribeLoadBalancersOutput{
						LoadBalancerDescriptions: []*elb.LoadBalancerDescription{
							{
								LoadBalancerName: aws.String(elbName),
								Scheme:           aws.String(string(infrav1.ELBSchemeInternetFacing)),
								Subnets:          []*string{aws.String(clusterSubnetID)},
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
				m.DescribeTags(&elb.DescribeTagsInput{LoadBalancerNames: []*string{aws.String(elbName)}}).Return(
					&elb.DescribeTagsOutput{
						TagDescriptions: []*elb.TagDescription{
							{
								LoadBalancerName: aws.String(elbName),
								Tags: []*elb.Tag{{
									Key:   aws.String(infrav1.ClusterTagKey(clusterName)),
									Value: aws.String(string(infrav1.ResourceLifecycleOwned)),
								}},
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
			ec2Mocks: func(m *mocks.MockEC2APIMockRecorder) {},
			check: func(t *testing.T, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}
			},
		},
		{
			name: "load balancer subnets specified in the same az from the instance",
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: clusterName},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: infrav1.Subnets{{
							ID:               clusterSubnetID,
							AvailabilityZone: az,
						}},
					},
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						Name:    aws.String("bar-apiserver"),
						Subnets: []string{elbSubnetID},
					},
				},
			},
			elbAPIMocks: func(m *mocks.MockELBAPIMockRecorder) {
				m.DescribeLoadBalancers(gomock.Eq(&elb.DescribeLoadBalancersInput{
					LoadBalancerNames: aws.StringSlice([]string{elbName}),
				})).
					Return(&elb.DescribeLoadBalancersOutput{
						LoadBalancerDescriptions: []*elb.LoadBalancerDescription{
							{
								Scheme:            aws.String(string(infrav1.ELBSchemeInternetFacing)),
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
				m.DescribeTags(&elb.DescribeTagsInput{LoadBalancerNames: []*string{aws.String(elbName)}}).Return(
					&elb.DescribeTagsOutput{
						TagDescriptions: []*elb.TagDescription{
							{
								LoadBalancerName: aws.String(elbName),
								Tags: []*elb.Tag{{
									Key:   aws.String(infrav1.ClusterTagKey(clusterName)),
									Value: aws.String(string(infrav1.ResourceLifecycleOwned)),
								}},
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
			ec2Mocks: func(m *mocks.MockEC2APIMockRecorder) {
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
				t.Helper()
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}
			},
		},
		{
			name: "load balancer subnets specified in a different az from the instance",
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: clusterName},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: infrav1.Subnets{{
							ID:               clusterSubnetID,
							AvailabilityZone: az,
						}},
					},
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						Name:    aws.String(elbName),
						Subnets: []string{elbSubnetID},
					},
				},
			},
			elbAPIMocks: func(m *mocks.MockELBAPIMockRecorder) {
				m.DescribeLoadBalancers(gomock.Eq(&elb.DescribeLoadBalancersInput{
					LoadBalancerNames: aws.StringSlice([]string{elbName}),
				})).
					Return(&elb.DescribeLoadBalancersOutput{
						LoadBalancerDescriptions: []*elb.LoadBalancerDescription{
							{
								Scheme:            aws.String(string(infrav1.ELBSchemeInternetFacing)),
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
				m.DescribeTags(&elb.DescribeTagsInput{LoadBalancerNames: []*string{aws.String(elbName)}}).Return(
					&elb.DescribeTagsOutput{
						TagDescriptions: []*elb.TagDescription{
							{
								LoadBalancerName: aws.String(elbName),
								Tags: []*elb.Tag{{
									Key:   aws.String(infrav1.ClusterTagKey(clusterName)),
									Value: aws.String(string(infrav1.ResourceLifecycleOwned)),
								}},
							},
						},
					}, nil)
			},
			ec2Mocks: func(m *mocks.MockEC2APIMockRecorder) {
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
				t.Helper()
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
			elbAPIMocks := mocks.NewMockELBAPI(mockCtrl)
			ec2Mock := mocks.NewMockEC2API(mockCtrl)

			scheme, err := setupScheme()
			if err != nil {
				t.Fatal(err)
			}

			client := fake.NewClientBuilder().WithScheme(scheme).Build()
			clusterScope, err := scope.NewClusterScope(scope.ClusterScopeParams{
				Client: client,
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: namespace,
						Name:      clusterName,
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

func TestRegisterInstanceWithAPIServerNLB(t *testing.T) {
	const (
		namespace       = "foo"
		clusterName     = "bar"
		clusterSubnetID = "subnet-1"
		elbName         = "bar-apiserver"
		elbArn          = "arn::apiserver"
		elbSubnetID     = "elb-subnet"
		instanceID      = "test-instance"
		az              = "us-west-1a"
		differentAZ     = "us-east-2c"
	)

	tests := []struct {
		name          string
		awsCluster    *infrav1.AWSCluster
		elbV2APIMocks func(m *mocks.MockELBV2APIMockRecorder)
		ec2Mocks      func(m *mocks.MockEC2APIMockRecorder)
		check         func(t *testing.T, err error)
	}{
		{
			name: "no load balancer subnets specified",
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: clusterName},
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						Name:             aws.String(elbName),
						LoadBalancerType: infrav1.LoadBalancerTypeNLB,
					},
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: infrav1.Subnets{{
							ID:               clusterSubnetID,
							AvailabilityZone: az,
						}},
					},
				},
			},
			elbV2APIMocks: func(m *mocks.MockELBV2APIMockRecorder) {
				m.DescribeLoadBalancers(gomock.Eq(&elbv2.DescribeLoadBalancersInput{
					Names: aws.StringSlice([]string{elbName}),
				})).
					Return(&elbv2.DescribeLoadBalancersOutput{
						LoadBalancers: []*elbv2.LoadBalancer{
							{
								LoadBalancerArn:  aws.String(elbArn),
								LoadBalancerName: aws.String(elbName),
								Scheme:           aws.String(string(infrav1.ELBSchemeInternetFacing)),
								AvailabilityZones: []*elbv2.AvailabilityZone{
									{
										SubnetId: aws.String(clusterSubnetID),
									},
								},
							},
						},
					}, nil)
				m.DescribeLoadBalancerAttributes(gomock.Eq(&elbv2.DescribeLoadBalancerAttributesInput{
					LoadBalancerArn: aws.String(elbArn),
				})).
					Return(&elbv2.DescribeLoadBalancerAttributesOutput{
						Attributes: []*elbv2.LoadBalancerAttribute{
							{
								Key:   aws.String("load_balancing.cross_zone.enabled"),
								Value: aws.String("true"),
							},
						},
					}, nil)
				m.DescribeTags(&elbv2.DescribeTagsInput{ResourceArns: []*string{aws.String(elbArn)}}).Return(
					&elbv2.DescribeTagsOutput{
						TagDescriptions: []*elbv2.TagDescription{
							{
								ResourceArn: aws.String(elbArn),
								Tags: []*elbv2.Tag{{
									Key:   aws.String(infrav1.ClusterTagKey(clusterName)),
									Value: aws.String(string(infrav1.ResourceLifecycleOwned)),
								}},
							},
						},
					}, nil)
				m.DescribeTargetGroups(&elbv2.DescribeTargetGroupsInput{
					LoadBalancerArn: aws.String(elbArn),
				}).Return(&elbv2.DescribeTargetGroupsOutput{
					TargetGroups: []*elbv2.TargetGroup{
						{
							HealthCheckEnabled:  aws.Bool(true),
							HealthCheckPort:     aws.String("infrav1.DefaultAPIServerPort"),
							HealthCheckProtocol: aws.String("TCP"),
							LoadBalancerArns:    aws.StringSlice([]string{elbArn}),
							Port:                aws.Int64(infrav1.DefaultAPIServerPort),
							Protocol:            aws.String("TCP"),
							TargetGroupArn:      aws.String("target-group::arn"),
							TargetGroupName:     aws.String("something-generated"),
							VpcId:               aws.String("vpc-id"),
						},
					},
				}, nil)
				m.RegisterTargets(gomock.Eq(&elbv2.RegisterTargetsInput{
					TargetGroupArn: aws.String("target-group::arn"),
					Targets: []*elbv2.TargetDescription{
						{
							Id:   aws.String(instanceID),
							Port: aws.Int64(infrav1.DefaultAPIServerPort),
						},
					},
				})).Return(&elbv2.RegisterTargetsOutput{}, nil)
			},
			ec2Mocks: func(m *mocks.MockEC2APIMockRecorder) {},
			check: func(t *testing.T, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}
			},
		},
		{
			name: "there are no target groups to register the instance into",
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: clusterName},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: infrav1.Subnets{{
							ID:               clusterSubnetID,
							AvailabilityZone: az,
						}},
					},
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						Name:             aws.String(elbName),
						Subnets:          []string{elbSubnetID},
						LoadBalancerType: infrav1.LoadBalancerTypeNLB,
					},
				},
			},
			elbV2APIMocks: func(m *mocks.MockELBV2APIMockRecorder) {
				m.DescribeLoadBalancers(gomock.Eq(&elbv2.DescribeLoadBalancersInput{
					Names: aws.StringSlice([]string{elbName}),
				})).
					Return(&elbv2.DescribeLoadBalancersOutput{
						LoadBalancers: []*elbv2.LoadBalancer{
							{
								LoadBalancerArn:  aws.String(elbArn),
								LoadBalancerName: aws.String(elbName),
								Scheme:           aws.String(string(infrav1.ELBSchemeInternetFacing)),
								AvailabilityZones: []*elbv2.AvailabilityZone{
									{
										SubnetId: aws.String(clusterSubnetID),
									},
								},
							},
						},
					}, nil)
				m.DescribeLoadBalancerAttributes(gomock.Eq(&elbv2.DescribeLoadBalancerAttributesInput{
					LoadBalancerArn: aws.String(elbArn),
				})).
					Return(&elbv2.DescribeLoadBalancerAttributesOutput{
						Attributes: []*elbv2.LoadBalancerAttribute{
							{
								Key:   aws.String("load_balancing.cross_zone.enabled"),
								Value: aws.String("true"),
							},
						},
					}, nil)
				m.DescribeTags(&elbv2.DescribeTagsInput{ResourceArns: []*string{aws.String(elbArn)}}).Return(
					&elbv2.DescribeTagsOutput{
						TagDescriptions: []*elbv2.TagDescription{
							{
								ResourceArn: aws.String(elbArn),
								Tags: []*elbv2.Tag{{
									Key:   aws.String(infrav1.ClusterTagKey(clusterName)),
									Value: aws.String(string(infrav1.ResourceLifecycleOwned)),
								}},
							},
						},
					}, nil)
				m.DescribeTargetGroups(&elbv2.DescribeTargetGroupsInput{
					LoadBalancerArn: aws.String(elbArn),
				}).Return(&elbv2.DescribeTargetGroupsOutput{}, nil)
			},
			ec2Mocks: func(m *mocks.MockEC2APIMockRecorder) {},
			check: func(t *testing.T, err error) {
				t.Helper()
				expectedErrMsg := fmt.Sprintf("no target groups found for load balancer with arn '%s'", elbArn)
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
			elbV2APIMocks := mocks.NewMockELBV2API(mockCtrl)
			ec2Mock := mocks.NewMockEC2API(mockCtrl)

			scheme, err := setupScheme()
			if err != nil {
				t.Fatal(err)
			}

			client := fake.NewClientBuilder().WithScheme(scheme).Build()
			clusterScope, err := scope.NewClusterScope(scope.ClusterScopeParams{
				Client: client,
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: namespace,
						Name:      clusterName,
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

			tc.elbV2APIMocks(elbV2APIMocks.EXPECT())
			tc.ec2Mocks(ec2Mock.EXPECT())

			s := &Service{
				scope:       clusterScope,
				EC2Client:   ec2Mock,
				ELBV2Client: elbV2APIMocks,
			}

			err = s.RegisterInstanceWithAPIServerLB(instance)
			tc.check(t, err)
		})
	}
}

func TestCreateNLB(t *testing.T) {
	const (
		namespace       = "foo"
		clusterName     = "bar"
		clusterSubnetID = "subnet-1"
		elbName         = "bar-apiserver"
		elbArn          = "arn::apiserver"
		vpcID           = "vpc-id"
		dns             = "asdf:9999/asdf"
	)

	tests := []struct {
		name          string
		elbV2APIMocks func(m *mocks.MockELBV2APIMockRecorder)
		check         func(t *testing.T, lb *infrav1.LoadBalancer, err error)
		awsCluster    func(acl infrav1.AWSCluster) infrav1.AWSCluster
		spec          func(spec infrav1.LoadBalancer) infrav1.LoadBalancer
	}{
		{
			name: "main create flow",
			spec: func(spec infrav1.LoadBalancer) infrav1.LoadBalancer {
				return spec
			},
			awsCluster: func(acl infrav1.AWSCluster) infrav1.AWSCluster {
				return acl
			},
			elbV2APIMocks: func(m *mocks.MockELBV2APIMockRecorder) {
				m.CreateLoadBalancer(gomock.Eq(&elbv2.CreateLoadBalancerInput{
					Name:    aws.String(elbName),
					Scheme:  aws.String("internet-facing"),
					Type:    aws.String("network"),
					Subnets: aws.StringSlice([]string{clusterSubnetID}),
					Tags: []*elbv2.Tag{
						{
							Key:   aws.String("test"),
							Value: aws.String("tag"),
						},
					},
				})).Return(&elbv2.CreateLoadBalancerOutput{
					LoadBalancers: []*elbv2.LoadBalancer{
						{
							LoadBalancerArn:  aws.String(elbArn),
							LoadBalancerName: aws.String(elbName),
							Scheme:           aws.String(string(infrav1.ELBSchemeInternetFacing)),
							DNSName:          aws.String(dns),
						},
					},
				}, nil)
				m.CreateTargetGroup(gomock.Eq(&elbv2.CreateTargetGroupInput{
					HealthCheckEnabled:  aws.Bool(true),
					HealthCheckPort:     aws.String("infrav1.DefaultAPIServerPort"),
					HealthCheckProtocol: aws.String("tcp"),
					Name:                aws.String("name"),
					Port:                aws.Int64(infrav1.DefaultAPIServerPort),
					Protocol:            aws.String("TCP"),
					VpcId:               aws.String(vpcID),
				})).Return(&elbv2.CreateTargetGroupOutput{
					TargetGroups: []*elbv2.TargetGroup{
						{
							TargetGroupArn:  aws.String("target-group::arn"),
							TargetGroupName: aws.String("name"),
							VpcId:           aws.String(vpcID),
						},
					},
				}, nil)
				m.ModifyTargetGroupAttributes(gomock.Eq(&elbv2.ModifyTargetGroupAttributesInput{
					TargetGroupArn: aws.String("target-group::arn"),
					Attributes: []*elbv2.TargetGroupAttribute{
						{
							Key:   aws.String(infrav1.TargetGroupAttributeEnablePreserveClientIP),
							Value: aws.String("false"),
						},
					},
				})).Return(nil, nil)
				m.CreateListener(gomock.Eq(&elbv2.CreateListenerInput{
					DefaultActions: []*elbv2.Action{
						{
							TargetGroupArn: aws.String("target-group::arn"),
							Type:           aws.String(elbv2.ActionTypeEnumForward),
						},
					},
					LoadBalancerArn: aws.String(elbArn),
					Port:            aws.Int64(infrav1.DefaultAPIServerPort),
					Protocol:        aws.String("TCP"),
					Tags: []*elbv2.Tag{
						{
							Key:   aws.String("test"),
							Value: aws.String("tag"),
						},
					},
				})).Return(&elbv2.CreateListenerOutput{
					Listeners: []*elbv2.Listener{
						{
							ListenerArn: aws.String("listener::arn"),
						},
					},
				}, nil)
			},
			check: func(t *testing.T, lb *infrav1.LoadBalancer, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}
				if lb.DNSName != dns {
					t.Fatalf("DNSName did not equal expected value; was: '%s'", lb.DNSName)
				}
			},
		},
		{
			name: "created with ipv6 vpc",
			spec: func(spec infrav1.LoadBalancer) infrav1.LoadBalancer {
				return spec
			},
			awsCluster: func(acl infrav1.AWSCluster) infrav1.AWSCluster {
				acl.Spec.NetworkSpec.VPC.IPv6 = &infrav1.IPv6{
					CidrBlock: "2022:1234::/64",
					PoolID:    "pool-id",
				}
				return acl
			},
			elbV2APIMocks: func(m *mocks.MockELBV2APIMockRecorder) {
				m.CreateLoadBalancer(gomock.Eq(&elbv2.CreateLoadBalancerInput{
					Name:          aws.String(elbName),
					IpAddressType: aws.String("dualstack"),
					Scheme:        aws.String("internet-facing"),
					Type:          aws.String("network"),
					Subnets:       aws.StringSlice([]string{clusterSubnetID}),
					Tags: []*elbv2.Tag{
						{
							Key:   aws.String("test"),
							Value: aws.String("tag"),
						},
					},
				})).Return(&elbv2.CreateLoadBalancerOutput{
					LoadBalancers: []*elbv2.LoadBalancer{
						{
							LoadBalancerArn:  aws.String(elbArn),
							LoadBalancerName: aws.String(elbName),
							Scheme:           aws.String(string(infrav1.ELBSchemeInternetFacing)),
							DNSName:          aws.String(dns),
						},
					},
				}, nil)
				m.CreateTargetGroup(gomock.Eq(&elbv2.CreateTargetGroupInput{
					HealthCheckEnabled:  aws.Bool(true),
					HealthCheckPort:     aws.String("infrav1.DefaultAPIServerPort"),
					HealthCheckProtocol: aws.String("tcp"),
					Name:                aws.String("name"),
					Port:                aws.Int64(infrav1.DefaultAPIServerPort),
					Protocol:            aws.String("TCP"),
					VpcId:               aws.String(vpcID),
					IpAddressType:       aws.String("ipv6"),
				})).Return(&elbv2.CreateTargetGroupOutput{
					TargetGroups: []*elbv2.TargetGroup{
						{
							TargetGroupArn:  aws.String("target-group::arn"),
							TargetGroupName: aws.String("name"),
							VpcId:           aws.String(vpcID),
						},
					},
				}, nil)
				m.ModifyTargetGroupAttributes(gomock.Eq(&elbv2.ModifyTargetGroupAttributesInput{
					TargetGroupArn: aws.String("target-group::arn"),
					Attributes: []*elbv2.TargetGroupAttribute{
						{
							Key:   aws.String(infrav1.TargetGroupAttributeEnablePreserveClientIP),
							Value: aws.String("false"),
						},
					},
				})).Return(nil, nil)
				m.CreateListener(gomock.Eq(&elbv2.CreateListenerInput{
					DefaultActions: []*elbv2.Action{
						{
							TargetGroupArn: aws.String("target-group::arn"),
							Type:           aws.String(elbv2.ActionTypeEnumForward),
						},
					},
					LoadBalancerArn: aws.String(elbArn),
					Port:            aws.Int64(infrav1.DefaultAPIServerPort),
					Protocol:        aws.String("TCP"),
					Tags: []*elbv2.Tag{
						{
							Key:   aws.String("test"),
							Value: aws.String("tag"),
						},
					},
				})).Return(&elbv2.CreateListenerOutput{
					Listeners: []*elbv2.Listener{
						{
							ListenerArn: aws.String("listener::arn"),
						},
					},
				}, nil)
			},
			check: func(t *testing.T, lb *infrav1.LoadBalancer, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}
				if lb.DNSName != dns {
					t.Fatalf("DNSName did not equal expected value; was: '%s'", lb.DNSName)
				}
			},
		},
		{
			name: "creating a load balancer fails",
			spec: func(spec infrav1.LoadBalancer) infrav1.LoadBalancer {
				return spec
			},
			awsCluster: func(acl infrav1.AWSCluster) infrav1.AWSCluster {
				return acl
			},
			elbV2APIMocks: func(m *mocks.MockELBV2APIMockRecorder) {
				m.CreateLoadBalancer(gomock.Eq(&elbv2.CreateLoadBalancerInput{
					Name:    aws.String(elbName),
					Scheme:  aws.String("internet-facing"),
					Type:    aws.String("network"),
					Subnets: aws.StringSlice([]string{clusterSubnetID}),
					Tags: []*elbv2.Tag{
						{
							Key:   aws.String("test"),
							Value: aws.String("tag"),
						},
					},
				})).Return(nil, errors.New("nope"))
			},
			check: func(t *testing.T, lb *infrav1.LoadBalancer, err error) {
				t.Helper()
				if err == nil {
					t.Fatal("expected error, got nothing")
				}
				if !strings.Contains(err.Error(), "nope") {
					t.Fatalf("expected error to contain 'nope' was instead: %s", err)
				}
			},
		},
		{
			name: "no health check",
			spec: func(spec infrav1.LoadBalancer) infrav1.LoadBalancer {
				spec.ELBListeners = []infrav1.Listener{
					{
						Protocol: "TCP",
						Port:     infrav1.DefaultAPIServerPort,
						TargetGroup: infrav1.TargetGroupSpec{
							Name:     "name",
							Port:     infrav1.DefaultAPIServerPort,
							Protocol: "TCP",
							VpcID:    vpcID,
						},
					},
				}
				return spec
			},
			awsCluster: func(acl infrav1.AWSCluster) infrav1.AWSCluster {
				return acl
			},
			elbV2APIMocks: func(m *mocks.MockELBV2APIMockRecorder) {
				m.CreateLoadBalancer(gomock.Eq(&elbv2.CreateLoadBalancerInput{
					Name:    aws.String(elbName),
					Scheme:  aws.String("internet-facing"),
					Type:    aws.String("network"),
					Subnets: aws.StringSlice([]string{clusterSubnetID}),
					Tags: []*elbv2.Tag{
						{
							Key:   aws.String("test"),
							Value: aws.String("tag"),
						},
					},
				})).Return(&elbv2.CreateLoadBalancerOutput{
					LoadBalancers: []*elbv2.LoadBalancer{
						{
							LoadBalancerArn:  aws.String(elbArn),
							LoadBalancerName: aws.String(elbName),
							Scheme:           aws.String(string(infrav1.ELBSchemeInternetFacing)),
							DNSName:          aws.String(dns),
						},
					},
				}, nil)
				m.CreateTargetGroup(gomock.Eq(&elbv2.CreateTargetGroupInput{
					Name:     aws.String("name"),
					Port:     aws.Int64(infrav1.DefaultAPIServerPort),
					Protocol: aws.String("TCP"),
					VpcId:    aws.String(vpcID),
				})).Return(&elbv2.CreateTargetGroupOutput{
					TargetGroups: []*elbv2.TargetGroup{
						{
							TargetGroupArn:  aws.String("target-group::arn"),
							TargetGroupName: aws.String("name"),
							VpcId:           aws.String(vpcID),
						},
					},
				}, nil)
				m.ModifyTargetGroupAttributes(gomock.Eq(&elbv2.ModifyTargetGroupAttributesInput{
					TargetGroupArn: aws.String("target-group::arn"),
					Attributes: []*elbv2.TargetGroupAttribute{
						{
							Key:   aws.String(infrav1.TargetGroupAttributeEnablePreserveClientIP),
							Value: aws.String("false"),
						},
					},
				})).Return(nil, nil)
				m.CreateListener(gomock.Eq(&elbv2.CreateListenerInput{
					DefaultActions: []*elbv2.Action{
						{
							TargetGroupArn: aws.String("target-group::arn"),
							Type:           aws.String(elbv2.ActionTypeEnumForward),
						},
					},
					LoadBalancerArn: aws.String(elbArn),
					Port:            aws.Int64(infrav1.DefaultAPIServerPort),
					Protocol:        aws.String("TCP"),
					Tags: []*elbv2.Tag{
						{
							Key:   aws.String("test"),
							Value: aws.String("tag"),
						},
					},
				})).Return(&elbv2.CreateListenerOutput{
					Listeners: []*elbv2.Listener{
						{
							ListenerArn: aws.String("listener::arn"),
						},
					},
				}, nil)
			},
			check: func(t *testing.T, lb *infrav1.LoadBalancer, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}
				if lb.DNSName != dns {
					t.Fatalf("DNSName did not equal expected value; was: '%s'", lb.DNSName)
				}
			},
		},
		{
			name: "PreserveClientIP is enabled",
			spec: func(spec infrav1.LoadBalancer) infrav1.LoadBalancer {
				return spec
			},
			awsCluster: func(acl infrav1.AWSCluster) infrav1.AWSCluster {
				acl.Spec.ControlPlaneLoadBalancer.PreserveClientIP = true
				return acl
			},
			elbV2APIMocks: func(m *mocks.MockELBV2APIMockRecorder) {
				m.CreateLoadBalancer(gomock.Eq(&elbv2.CreateLoadBalancerInput{
					Name:    aws.String(elbName),
					Scheme:  aws.String("internet-facing"),
					Type:    aws.String("network"),
					Subnets: aws.StringSlice([]string{clusterSubnetID}),
					Tags: []*elbv2.Tag{
						{
							Key:   aws.String("test"),
							Value: aws.String("tag"),
						},
					},
				})).Return(&elbv2.CreateLoadBalancerOutput{
					LoadBalancers: []*elbv2.LoadBalancer{
						{
							LoadBalancerArn:  aws.String(elbArn),
							LoadBalancerName: aws.String(elbName),
							Scheme:           aws.String(string(infrav1.ELBSchemeInternetFacing)),
							DNSName:          aws.String(dns),
						},
					},
				}, nil)
				m.CreateTargetGroup(gomock.Eq(&elbv2.CreateTargetGroupInput{
					HealthCheckEnabled:  aws.Bool(true),
					HealthCheckPort:     aws.String("infrav1.DefaultAPIServerPort"),
					HealthCheckProtocol: aws.String("tcp"),
					Name:                aws.String("name"),
					Port:                aws.Int64(infrav1.DefaultAPIServerPort),
					Protocol:            aws.String("TCP"),
					VpcId:               aws.String(vpcID),
				})).Return(&elbv2.CreateTargetGroupOutput{
					TargetGroups: []*elbv2.TargetGroup{
						{
							TargetGroupArn:  aws.String("target-group::arn"),
							TargetGroupName: aws.String("name"),
							VpcId:           aws.String(vpcID),
						},
					},
				}, nil)
				m.CreateListener(gomock.Eq(&elbv2.CreateListenerInput{
					DefaultActions: []*elbv2.Action{
						{
							TargetGroupArn: aws.String("target-group::arn"),
							Type:           aws.String(elbv2.ActionTypeEnumForward),
						},
					},
					LoadBalancerArn: aws.String(elbArn),
					Port:            aws.Int64(infrav1.DefaultAPIServerPort),
					Protocol:        aws.String("TCP"),
					Tags: []*elbv2.Tag{
						{
							Key:   aws.String("test"),
							Value: aws.String("tag"),
						},
					},
				})).Return(&elbv2.CreateListenerOutput{
					Listeners: []*elbv2.Listener{
						{
							ListenerArn: aws.String("listener::arn"),
						},
					},
				}, nil)
			},
			check: func(t *testing.T, lb *infrav1.LoadBalancer, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}
				if lb.DNSName != dns {
					t.Fatalf("DNSName did not equal expected value; was: '%s'", lb.DNSName)
				}
			},
		},
		{
			name: "load balancer is not an NLB scope security groups will be added",
			spec: func(spec infrav1.LoadBalancer) infrav1.LoadBalancer {
				spec.SecurityGroupIDs = []string{"sg-id"}
				return spec
			},
			awsCluster: func(acl infrav1.AWSCluster) infrav1.AWSCluster {
				acl.Spec.ControlPlaneLoadBalancer.LoadBalancerType = infrav1.LoadBalancerTypeALB
				return acl
			},
			elbV2APIMocks: func(m *mocks.MockELBV2APIMockRecorder) {
				m.CreateLoadBalancer(gomock.Eq(&elbv2.CreateLoadBalancerInput{
					Name:    aws.String(elbName),
					Scheme:  aws.String("internet-facing"),
					Type:    aws.String("application"),
					Subnets: aws.StringSlice([]string{clusterSubnetID}),
					Tags: []*elbv2.Tag{
						{
							Key:   aws.String("test"),
							Value: aws.String("tag"),
						},
					},
					SecurityGroups: aws.StringSlice([]string{"sg-id"}),
				})).Return(&elbv2.CreateLoadBalancerOutput{
					LoadBalancers: []*elbv2.LoadBalancer{
						{
							LoadBalancerArn:  aws.String(elbArn),
							LoadBalancerName: aws.String(elbName),
							Scheme:           aws.String(string(infrav1.ELBSchemeInternetFacing)),
							DNSName:          aws.String(dns),
						},
					},
				}, nil)
				m.CreateTargetGroup(gomock.Eq(&elbv2.CreateTargetGroupInput{
					HealthCheckEnabled:  aws.Bool(true),
					HealthCheckPort:     aws.String("infrav1.DefaultAPIServerPort"),
					HealthCheckProtocol: aws.String("tcp"),
					Name:                aws.String("name"),
					Port:                aws.Int64(infrav1.DefaultAPIServerPort),
					Protocol:            aws.String("TCP"),
					VpcId:               aws.String(vpcID),
				})).Return(&elbv2.CreateTargetGroupOutput{
					TargetGroups: []*elbv2.TargetGroup{
						{
							TargetGroupArn:  aws.String("target-group::arn"),
							TargetGroupName: aws.String("name"),
							VpcId:           aws.String(vpcID),
						},
					},
				}, nil)
				m.ModifyTargetGroupAttributes(gomock.Eq(&elbv2.ModifyTargetGroupAttributesInput{
					TargetGroupArn: aws.String("target-group::arn"),
					Attributes: []*elbv2.TargetGroupAttribute{
						{
							Key:   aws.String(infrav1.TargetGroupAttributeEnablePreserveClientIP),
							Value: aws.String("false"),
						},
					},
				})).Return(nil, nil)
				m.CreateListener(gomock.Eq(&elbv2.CreateListenerInput{
					DefaultActions: []*elbv2.Action{
						{
							TargetGroupArn: aws.String("target-group::arn"),
							Type:           aws.String(elbv2.ActionTypeEnumForward),
						},
					},
					LoadBalancerArn: aws.String(elbArn),
					Port:            aws.Int64(infrav1.DefaultAPIServerPort),
					Protocol:        aws.String("TCP"),
					Tags: []*elbv2.Tag{
						{
							Key:   aws.String("test"),
							Value: aws.String("tag"),
						},
					},
				})).Return(&elbv2.CreateListenerOutput{
					Listeners: []*elbv2.Listener{
						{
							ListenerArn: aws.String("listener::arn"),
						},
					},
				}, nil)
			},
			check: func(t *testing.T, lb *infrav1.LoadBalancer, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}
				if lb.DNSName != dns {
					t.Fatalf("DNSName did not equal expected value; was: '%s'", lb.DNSName)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			elbV2APIMocks := mocks.NewMockELBV2API(mockCtrl)

			scheme, err := setupScheme()
			if err != nil {
				t.Fatal(err)
			}
			awsCluster := &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: clusterName},
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						Name:             aws.String(elbName),
						LoadBalancerType: infrav1.LoadBalancerTypeNLB,
					},
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							ID: vpcID,
						},
					},
				},
			}
			client := fake.NewClientBuilder().WithScheme(scheme).Build()
			cluster := tc.awsCluster(*awsCluster)
			clusterScope, err := scope.NewClusterScope(scope.ClusterScopeParams{
				Client: client,
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: namespace,
						Name:      clusterName,
					},
				},
				AWSCluster: &cluster,
			})
			if err != nil {
				t.Fatal(err)
			}

			tc.elbV2APIMocks(elbV2APIMocks.EXPECT())

			s := &Service{
				scope:       clusterScope,
				ELBV2Client: elbV2APIMocks,
			}

			loadBalancerSpec := &infrav1.LoadBalancer{
				ARN:    elbArn,
				Name:   elbName,
				Scheme: infrav1.ELBSchemeInternetFacing,
				Tags: map[string]string{
					"test": "tag",
				},
				ELBListeners: []infrav1.Listener{
					{
						Protocol: "TCP",
						Port:     infrav1.DefaultAPIServerPort,
						TargetGroup: infrav1.TargetGroupSpec{
							Name:     "name",
							Port:     infrav1.DefaultAPIServerPort,
							Protocol: "TCP",
							VpcID:    vpcID,
							HealthCheck: &infrav1.TargetGroupHealthCheck{
								Protocol: aws.String("tcp"),
								Port:     aws.String("infrav1.DefaultAPIServerPort"),
							},
						},
					},
				},
				LoadBalancerType: infrav1.LoadBalancerTypeNLB,
				SubnetIDs:        []string{clusterSubnetID},
			}

			spec := tc.spec(*loadBalancerSpec)
			lb, err := s.createLB(&spec)
			tc.check(t, lb, err)
		})
	}
}

func TestDeleteAPIServerELB(t *testing.T) {
	clusterName := "bar" //nolint:goconst // does not need to be a package-level const
	elbName := "bar-apiserver"
	tests := []struct {
		name             string
		elbAPIMocks      func(m *mocks.MockELBAPIMockRecorder)
		verifyAWSCluster func(*infrav1.AWSCluster)
	}{
		{
			name: "if control plane ELB is not found, do nothing",
			elbAPIMocks: func(m *mocks.MockELBAPIMockRecorder) {
				m.DescribeLoadBalancers(gomock.Eq(&elb.DescribeLoadBalancersInput{
					LoadBalancerNames: aws.StringSlice([]string{elbName}),
				})).Return(nil, awserr.New(elb.ErrCodeAccessPointNotFoundException, "", nil))
			},
			verifyAWSCluster: func(awsCluster *infrav1.AWSCluster) {
				loadBalancerConditionReady := conditions.IsTrue(awsCluster, infrav1.LoadBalancerReadyCondition)
				if loadBalancerConditionReady {
					t.Fatalf("Expected LoadBalancerReady condition to be False, but was True")
				}
				loadBalancerConditionReason := conditions.GetReason(awsCluster, infrav1.LoadBalancerReadyCondition)
				if loadBalancerConditionReason != clusterv1.DeletedReason {
					t.Fatalf("Expected LoadBalancerReady condition reason to be Deleted, but was %s", loadBalancerConditionReason)
				}
			},
		},
		{
			name: "if control plane ELB is found, and it is not managed, do nothing",
			elbAPIMocks: func(m *mocks.MockELBAPIMockRecorder) {
				m.DescribeLoadBalancers(&elb.DescribeLoadBalancersInput{LoadBalancerNames: []*string{aws.String(elbName)}}).Return(
					&elb.DescribeLoadBalancersOutput{
						LoadBalancerDescriptions: []*elb.LoadBalancerDescription{
							{
								LoadBalancerName: aws.String(elbName),
								Scheme:           aws.String(string(infrav1.ELBSchemeInternetFacing)),
							},
						},
					},
					nil,
				)

				m.DescribeLoadBalancerAttributes(&elb.DescribeLoadBalancerAttributesInput{LoadBalancerName: aws.String(elbName)}).Return(
					&elb.DescribeLoadBalancerAttributesOutput{
						LoadBalancerAttributes: &elb.LoadBalancerAttributes{
							CrossZoneLoadBalancing: &elb.CrossZoneLoadBalancing{
								Enabled: aws.Bool(false),
							},
						},
					},
					nil,
				)

				m.DescribeTags(&elb.DescribeTagsInput{LoadBalancerNames: []*string{aws.String(elbName)}}).Return(
					&elb.DescribeTagsOutput{
						TagDescriptions: []*elb.TagDescription{
							{
								LoadBalancerName: aws.String(elbName),
								Tags:             []*elb.Tag{},
							},
						},
					},
					nil,
				)
			},
			verifyAWSCluster: func(awsCluster *infrav1.AWSCluster) {
				loadBalancerConditionReady := conditions.IsTrue(awsCluster, infrav1.LoadBalancerReadyCondition)
				if loadBalancerConditionReady {
					t.Fatalf("Expected LoadBalancerReady condition to be False, but was True")
				}
				loadBalancerConditionReason := conditions.GetReason(awsCluster, infrav1.LoadBalancerReadyCondition)
				if loadBalancerConditionReason != clusterv1.DeletedReason {
					t.Fatalf("Expected LoadBalancerReady condition reason to be Deleted, but was %s", loadBalancerConditionReason)
				}
			},
		},
		{
			name: "if control plane ELB is found, and it is managed, delete the ELB",
			elbAPIMocks: func(m *mocks.MockELBAPIMockRecorder) {
				m.DescribeLoadBalancers(&elb.DescribeLoadBalancersInput{LoadBalancerNames: []*string{aws.String(elbName)}}).Return(
					&elb.DescribeLoadBalancersOutput{
						LoadBalancerDescriptions: []*elb.LoadBalancerDescription{
							{
								LoadBalancerName: aws.String(elbName),
								Scheme:           aws.String(string(infrav1.ELBSchemeInternetFacing)),
							},
						},
					},
					nil,
				)

				m.DescribeLoadBalancerAttributes(&elb.DescribeLoadBalancerAttributesInput{LoadBalancerName: aws.String(elbName)}).Return(
					&elb.DescribeLoadBalancerAttributesOutput{
						LoadBalancerAttributes: &elb.LoadBalancerAttributes{
							CrossZoneLoadBalancing: &elb.CrossZoneLoadBalancing{
								Enabled: aws.Bool(false),
							},
						},
					},
					nil,
				)

				m.DescribeTags(&elb.DescribeTagsInput{LoadBalancerNames: []*string{aws.String(elbName)}}).Return(
					&elb.DescribeTagsOutput{
						TagDescriptions: []*elb.TagDescription{
							{
								LoadBalancerName: aws.String(elbName),
								Tags: []*elb.Tag{{
									Key:   aws.String(infrav1.ClusterTagKey(clusterName)),
									Value: aws.String(string(infrav1.ResourceLifecycleOwned)),
								}},
							},
						},
					},
					nil,
				)

				m.DeleteLoadBalancer(&elb.DeleteLoadBalancerInput{LoadBalancerName: aws.String(elbName)}).Return(
					&elb.DeleteLoadBalancerOutput{}, nil)

				m.DescribeLoadBalancers(&elb.DescribeLoadBalancersInput{LoadBalancerNames: []*string{aws.String(elbName)}}).Return(
					&elb.DescribeLoadBalancersOutput{
						LoadBalancerDescriptions: []*elb.LoadBalancerDescription{},
					},
					nil,
				)
			},
			verifyAWSCluster: func(awsCluster *infrav1.AWSCluster) {
				loadBalancerConditionReady := conditions.IsTrue(awsCluster, infrav1.LoadBalancerReadyCondition)
				if loadBalancerConditionReady {
					t.Fatalf("Expected LoadBalancerReady condition to be False, but was True")
				}
				loadBalancerConditionReason := conditions.GetReason(awsCluster, infrav1.LoadBalancerReadyCondition)
				if loadBalancerConditionReason != clusterv1.DeletedReason {
					t.Fatalf("Expected LoadBalancerReady condition reason to be Deleted, but was %s", loadBalancerConditionReason)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			rgapiMock := mocks.NewMockResourceGroupsTaggingAPIAPI(mockCtrl)
			elbapiMock := mocks.NewMockELBAPI(mockCtrl)

			scheme, err := setupScheme()
			if err != nil {
				t.Fatal(err)
			}

			awsCluster := &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						Name: aws.String(elbName),
					},
				},
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

			tc.elbAPIMocks(elbapiMock.EXPECT())

			s := &Service{
				scope:                 clusterScope,
				ResourceTaggingClient: rgapiMock,
				ELBClient:             elbapiMock,
			}

			err = s.deleteAPIServerELB()
			if err != nil {
				t.Fatal(err)
			}

			tc.verifyAWSCluster(awsCluster)
		})
	}
}

func TestDeleteNLB(t *testing.T) {
	clusterName := "bar"
	elbName := "bar-apiserver"
	elbArn := "apiserver::arn"
	tests := []struct {
		name         string
		elbv2ApiMock func(m *mocks.MockELBV2APIMockRecorder)
	}{
		{
			name: "if control plane NLB is not found, do nothing",
			elbv2ApiMock: func(m *mocks.MockELBV2APIMockRecorder) {
				m.DescribeLoadBalancers(gomock.Eq(&elbv2.DescribeLoadBalancersInput{
					Names: aws.StringSlice([]string{elbName}),
				})).Return(nil, awserr.New(elb.ErrCodeAccessPointNotFoundException, "", nil))
			},
		},
		{
			name: "if control plane NLB is found, and it is not managed, do nothing",
			elbv2ApiMock: func(m *mocks.MockELBV2APIMockRecorder) {
				m.DescribeLoadBalancers(&elbv2.DescribeLoadBalancersInput{Names: []*string{aws.String(elbName)}}).Return(
					&elbv2.DescribeLoadBalancersOutput{
						LoadBalancers: []*elbv2.LoadBalancer{
							{
								LoadBalancerArn:  aws.String(elbArn),
								LoadBalancerName: aws.String(elbName),
								Scheme:           aws.String(string(infrav1.ELBSchemeInternetFacing)),
							},
						},
					},
					nil,
				)

				m.DescribeLoadBalancerAttributes(&elbv2.DescribeLoadBalancerAttributesInput{LoadBalancerArn: aws.String(elbArn)}).Return(
					&elbv2.DescribeLoadBalancerAttributesOutput{
						Attributes: []*elbv2.LoadBalancerAttribute{
							{
								Key:   aws.String("load_balancing.cross_zone.enabled"),
								Value: aws.String("false"),
							},
						},
					},
					nil,
				)

				m.DescribeTags(&elbv2.DescribeTagsInput{ResourceArns: []*string{aws.String(elbArn)}}).Return(
					&elbv2.DescribeTagsOutput{
						TagDescriptions: []*elbv2.TagDescription{
							{
								ResourceArn: aws.String(elbArn),
								Tags:        []*elbv2.Tag{},
							},
						},
					},
					nil,
				)
			},
		},
		{
			name: "if control plane ELB is found, and it is managed, delete the ELB",
			elbv2ApiMock: func(m *mocks.MockELBV2APIMockRecorder) {
				m.DescribeLoadBalancers(&elbv2.DescribeLoadBalancersInput{Names: []*string{aws.String(elbName)}}).Return(
					&elbv2.DescribeLoadBalancersOutput{
						LoadBalancers: []*elbv2.LoadBalancer{
							{
								LoadBalancerArn:  aws.String(elbArn),
								LoadBalancerName: aws.String(elbName),
								Scheme:           aws.String(string(infrav1.ELBSchemeInternetFacing)),
							},
						},
					},
					nil,
				)

				m.DescribeLoadBalancerAttributes(&elbv2.DescribeLoadBalancerAttributesInput{LoadBalancerArn: aws.String(elbArn)}).Return(
					&elbv2.DescribeLoadBalancerAttributesOutput{
						Attributes: []*elbv2.LoadBalancerAttribute{
							{
								Key:   aws.String("load_balancing.cross_zone.enabled"),
								Value: aws.String("false"),
							},
						},
					},
					nil,
				)

				m.DescribeTags(&elbv2.DescribeTagsInput{ResourceArns: []*string{aws.String(elbArn)}}).Return(
					&elbv2.DescribeTagsOutput{
						TagDescriptions: []*elbv2.TagDescription{
							{
								ResourceArn: aws.String(elbArn),
								Tags: []*elbv2.Tag{{
									Key:   aws.String(infrav1.ClusterTagKey(clusterName)),
									Value: aws.String(string(infrav1.ResourceLifecycleOwned)),
								}},
							},
						},
					},
					nil,
				)

				// delete listeners
				m.DescribeListeners(&elbv2.DescribeListenersInput{LoadBalancerArn: aws.String(elbArn)}).Return(&elbv2.DescribeListenersOutput{
					Listeners: []*elbv2.Listener{
						{
							ListenerArn: aws.String("listener::arn"),
						},
					},
				}, nil)
				m.DeleteListener(&elbv2.DeleteListenerInput{ListenerArn: aws.String("listener::arn")}).Return(&elbv2.DeleteListenerOutput{}, nil)
				// delete target groups
				m.DescribeTargetGroups(&elbv2.DescribeTargetGroupsInput{LoadBalancerArn: aws.String(elbArn)}).Return(&elbv2.DescribeTargetGroupsOutput{
					TargetGroups: []*elbv2.TargetGroup{
						{
							TargetGroupArn: aws.String("target-group::arn"),
						},
					},
				}, nil)
				m.DeleteTargetGroup(&elbv2.DeleteTargetGroupInput{TargetGroupArn: aws.String("target-group::arn")}).Return(&elbv2.DeleteTargetGroupOutput{}, nil)
				// delete the load balancer

				m.DeleteLoadBalancer(&elbv2.DeleteLoadBalancerInput{LoadBalancerArn: aws.String(elbArn)}).Return(
					&elbv2.DeleteLoadBalancerOutput{}, nil)

				m.DescribeLoadBalancers(&elbv2.DescribeLoadBalancersInput{Names: []*string{aws.String(elbName)}}).Return(
					&elbv2.DescribeLoadBalancersOutput{
						LoadBalancers: []*elbv2.LoadBalancer{},
					},
					nil,
				)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			rgapiMock := mocks.NewMockResourceGroupsTaggingAPIAPI(mockCtrl)
			elbv2ApiMock := mocks.NewMockELBV2API(mockCtrl)

			scheme, err := setupScheme()
			if err != nil {
				t.Fatal(err)
			}

			awsCluster := &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						Name:             aws.String(elbName),
						LoadBalancerType: infrav1.LoadBalancerTypeNLB,
					},
				},
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

			tc.elbv2ApiMock(elbv2ApiMock.EXPECT())

			s := &Service{
				scope:                 clusterScope,
				ResourceTaggingClient: rgapiMock,
				ELBV2Client:           elbv2ApiMock,
			}

			err = s.deleteExistingNLBs()
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestDeleteAWSCloudProviderELBs(t *testing.T) {
	clusterName := "bar"
	tests := []struct {
		name                  string
		rgAPIMocks            func(m *mocks.MockResourceGroupsTaggingAPIAPIMockRecorder)
		elbAPIMocks           func(m *mocks.MockELBAPIMockRecorder)
		postDeleteRGAPIMocks  func(m *mocks.MockResourceGroupsTaggingAPIAPIMockRecorder)
		postDeleteElbAPIMocks func(m *mocks.MockELBAPIMockRecorder)
	}{
		{
			name: "discover ELBs with Resource Groups Tagging API and then delete successfully",
			rgAPIMocks: func(m *mocks.MockResourceGroupsTaggingAPIAPIMockRecorder) {
				m.GetResourcesPages(&rgapi.GetResourcesInput{
					ResourceTypeFilters: aws.StringSlice([]string{elbResourceType}),
					TagFilters: []*rgapi.TagFilter{
						{
							Key:    aws.String(infrav1.ClusterAWSCloudProviderTagKey(clusterName)),
							Values: aws.StringSlice([]string{string(infrav1.ResourceLifecycleOwned)}),
						},
					},
				}, gomock.Any()).Do(func(_, y interface{}) {
					funct := y.(func(output *rgapi.GetResourcesOutput, lastPage bool) bool)
					funct(&rgapi.GetResourcesOutput{
						ResourceTagMappingList: []*rgapi.ResourceTagMapping{
							{
								ResourceARN: aws.String("arn:aws:elasticloadbalancing:eu-west-2:1234567890:loadbalancer/lb-service-name"),
								Tags: []*rgapi.Tag{{
									Key:   aws.String(infrav1.ClusterAWSCloudProviderTagKey(clusterName)),
									Value: aws.String(string(infrav1.ResourceLifecycleOwned)),
								}},
							},
						},
					}, true)
				}).Return(nil)
			},
			elbAPIMocks: func(m *mocks.MockELBAPIMockRecorder) {
				m.DeleteLoadBalancer(gomock.Eq(&elb.DeleteLoadBalancerInput{LoadBalancerName: aws.String("lb-service-name")})).Return(nil, nil)
			},
			postDeleteRGAPIMocks: func(m *mocks.MockResourceGroupsTaggingAPIAPIMockRecorder) {
				m.GetResourcesPages(&rgapi.GetResourcesInput{
					ResourceTypeFilters: aws.StringSlice([]string{elbResourceType}),
					TagFilters: []*rgapi.TagFilter{
						{
							Key:    aws.String(infrav1.ClusterAWSCloudProviderTagKey(clusterName)),
							Values: aws.StringSlice([]string{string(infrav1.ResourceLifecycleOwned)}),
						},
					},
				}, gomock.Any()).Do(func(_, y interface{}) {
					funct := y.(func(output *rgapi.GetResourcesOutput, lastPage bool) bool)
					funct(&rgapi.GetResourcesOutput{
						ResourceTagMappingList: []*rgapi.ResourceTagMapping{},
					}, true)
				}).Return(nil)
			},
		},
		{
			name: "fall back to ELB API when Resource Groups Tagging API fails and then delete successfully",
			rgAPIMocks: func(m *mocks.MockResourceGroupsTaggingAPIAPIMockRecorder) {
				m.GetResourcesPages(gomock.Any(), gomock.Any()).Return(errors.Errorf("connection failure")).AnyTimes()
			},
			elbAPIMocks: func(m *mocks.MockELBAPIMockRecorder) {
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
				m.DeleteLoadBalancer(gomock.Eq(&elb.DeleteLoadBalancerInput{LoadBalancerName: aws.String("lb-service-name")})).Return(nil, nil)
			},
			postDeleteElbAPIMocks: func(m *mocks.MockELBAPIMockRecorder) {
				m.DescribeLoadBalancersPages(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			rgapiMock := mocks.NewMockResourceGroupsTaggingAPIAPI(mockCtrl)
			elbapiMock := mocks.NewMockELBAPI(mockCtrl)

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
			if tc.postDeleteElbAPIMocks != nil {
				tc.postDeleteElbAPIMocks(elbapiMock.EXPECT())
			}
			if tc.postDeleteRGAPIMocks != nil {
				tc.postDeleteRGAPIMocks(rgapiMock.EXPECT())
			}

			s := &Service{
				scope:                 clusterScope,
				ResourceTaggingClient: rgapiMock,
				ELBClient:             elbapiMock,
			}

			err = s.deleteAWSCloudProviderELBs()
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
		rgAPIMocks          func(m *mocks.MockResourceGroupsTaggingAPIAPIMockRecorder)
		DescribeElbAPIMocks func(m *mocks.MockELBAPIMockRecorder)
	}{
		{
			name:   "Error if existing loadbalancer with same name doesn't have same scheme",
			lbName: "bar-apiserver",
			rgAPIMocks: func(m *mocks.MockResourceGroupsTaggingAPIAPIMockRecorder) {
				m.GetResourcesPages(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			},
			DescribeElbAPIMocks: func(m *mocks.MockELBAPIMockRecorder) {
				m.DescribeLoadBalancers(gomock.Eq(&elb.DescribeLoadBalancersInput{
					LoadBalancerNames: aws.StringSlice([]string{"bar-apiserver"}),
				})).Return(&elb.DescribeLoadBalancersOutput{LoadBalancerDescriptions: []*elb.LoadBalancerDescription{{Scheme: pointer.String(string(infrav1.ELBSchemeInternal))}}}, nil)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			rgapiMock := mocks.NewMockResourceGroupsTaggingAPIAPI(mockCtrl)
			elbapiMock := mocks.NewMockELBAPI(mockCtrl)

			scheme, err := setupScheme()
			if err != nil {
				t.Fatal(err)
			}
			awsCluster := &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
					Scheme: &infrav1.ELBSchemeInternetFacing,
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

func TestDescribeV2Loadbalancers(t *testing.T) {
	clusterName := "bar"
	tests := []struct {
		name                  string
		lbName                string
		rgAPIMocks            func(m *mocks.MockResourceGroupsTaggingAPIAPIMockRecorder)
		DescribeElbV2APIMocks func(m *mocks.MockELBV2APIMockRecorder)
	}{
		{
			name:   "Error if existing loadbalancer with same name doesn't have same scheme",
			lbName: "bar-apiserver",
			rgAPIMocks: func(m *mocks.MockResourceGroupsTaggingAPIAPIMockRecorder) {
				m.GetResourcesPages(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			},
			DescribeElbV2APIMocks: func(m *mocks.MockELBV2APIMockRecorder) {
				m.DescribeLoadBalancers(gomock.Eq(&elbv2.DescribeLoadBalancersInput{
					Names: aws.StringSlice([]string{"bar-apiserver"}),
				})).Return(&elbv2.DescribeLoadBalancersOutput{LoadBalancers: []*elbv2.LoadBalancer{{Scheme: pointer.String(string(infrav1.ELBSchemeInternal))}}}, nil)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			rgapiMock := mocks.NewMockResourceGroupsTaggingAPIAPI(mockCtrl)
			elbV2ApiMock := mocks.NewMockELBV2API(mockCtrl)

			scheme, err := setupScheme()
			if err != nil {
				t.Fatal(err)
			}
			awsCluster := &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
					Scheme:           &infrav1.ELBSchemeInternetFacing,
					LoadBalancerType: infrav1.LoadBalancerTypeNLB,
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
			tc.DescribeElbV2APIMocks(elbV2ApiMock.EXPECT())

			s := &Service{
				scope:                 clusterScope,
				ResourceTaggingClient: rgapiMock,
				ELBV2Client:           elbV2ApiMock,
			}

			_, err = s.describeLB(tc.lbName)
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

func TestGetHealthCheckProtocol(t *testing.T) {
	tests := []struct {
		testName                  string
		lbSpec                    *infrav1.AWSLoadBalancerSpec
		expectedHealthCheckTarget string
	}{
		{
			"default case",
			&infrav1.AWSLoadBalancerSpec{},
			"SSL:6443",
		},
		{
			"protocol http",
			&infrav1.AWSLoadBalancerSpec{
				HealthCheckProtocol: &infrav1.ELBProtocolHTTP,
			},
			"HTTP:6443/readyz",
		},
		{
			"protocol https",
			&infrav1.AWSLoadBalancerSpec{
				HealthCheckProtocol: &infrav1.ELBProtocolHTTPS,
			},
			"HTTPS:6443/readyz",
		},
		{
			"protocol tcp",
			&infrav1.AWSLoadBalancerSpec{
				HealthCheckProtocol: &infrav1.ELBProtocolTCP,
			},
			"TCP:6443",
		},
	}
	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			scheme := runtime.NewScheme()
			_ = infrav1.AddToScheme(scheme)
			client := fake.NewClientBuilder().WithScheme(scheme).Build()

			scope, err := scope.NewClusterScope(scope.ClusterScopeParams{
				Client: client,
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-elb",
						Namespace: "default",
					},
				},
				AWSCluster: &infrav1.AWSCluster{
					Spec: infrav1.AWSClusterSpec{
						ControlPlaneLoadBalancer: tc.lbSpec,
					},
				},
			})
			if err != nil {
				t.Fatal(err)
			}
			s := &Service{
				scope: scope,
			}
			healthCheck := s.getHealthCheckTarget()
			if healthCheck != tc.expectedHealthCheckTarget {
				t.Errorf("got %s, want %s", healthCheck, tc.expectedHealthCheckTarget)
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
