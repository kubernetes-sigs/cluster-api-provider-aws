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
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	elbtypes "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing/types"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbv2types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	rgapi "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	rgapitypes "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/helpers"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
)

var stubInfraV1TargetGroupSpecAPI = infrav1.TargetGroupSpec{
	Name:     "name",
	Port:     infrav1.DefaultAPIServerPort,
	Protocol: "TCP",
	HealthCheck: &infrav1.TargetGroupHealthCheck{
		IntervalSeconds:         aws.Int64(10),
		TimeoutSeconds:          aws.Int64(5),
		ThresholdCount:          aws.Int64(5),
		UnhealthyThresholdCount: aws.Int64(3),
	},
}

const maxWaitActiveUpdateDelete = time.Minute

func TestELBName(t *testing.T) {
	tests := []struct {
		name       string
		awsCluster *infrav1.AWSCluster
		expected   string
	}{
		{
			name: "name is not defined by user, so generate the default",
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "example",
					Namespace: metav1.NamespaceDefault,
				},
			},
			expected: "example-apiserver",
		},
		{
			name: "name is defined by user, so use it",
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "example",
					Namespace: metav1.NamespaceDefault,
				},
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						Name: ptr.To[string]("myapiserver"),
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
				AWSCluster: tt.awsCluster,
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
				m.DescribeSubnets(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{
					SubnetIds: []string{
						"subnet-1",
						"subnet-2",
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []ec2types.Subnet{
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
				HealthCheckProtocol: &infrav1.ELBProtocolUDP,
			},
			mocks: func(m *mocks.MockEC2APIMockRecorder) {},
			expect: func(t *testing.T, g *WithT, res *infrav1.LoadBalancer) {
				t.Helper()
				expectedTarget := fmt.Sprintf("%v:%d", infrav1.ELBProtocolUDP, infrav1.DefaultAPIServerPort)
				g.Expect(expectedTarget).To(Equal(res.HealthCheck.Target))
			},
		},
		{
			name:  "Should create load balancer spec with default elb health check protocol",
			lb:    &infrav1.AWSLoadBalancerSpec{},
			mocks: func(m *mocks.MockEC2APIMockRecorder) {},
			expect: func(t *testing.T, g *WithT, res *infrav1.LoadBalancer) {
				t.Helper()
				expectedTarget := fmt.Sprintf("%v:%d", infrav1.ELBProtocolTCP, infrav1.DefaultAPIServerPort)
				g.Expect(expectedTarget).To(Equal(res.HealthCheck.Target))
			},
		},
	}

	ctx := context.TODO()

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

			spec, err := s.getAPIServerClassicELBSpec(ctx, clusterScope.Name())
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
				m.DescribeSubnets(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{
					SubnetIds: []string{
						"subnet-1",
						"subnet-2",
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []ec2types.Subnet{
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
		{
			name: "A base listener is set up for NLB, with additional listeners",
			lb: &infrav1.AWSLoadBalancerSpec{
				LoadBalancerType: infrav1.LoadBalancerTypeNLB,
				AdditionalListeners: []infrav1.AdditionalListenerSpec{
					{
						Port:     443,
						Protocol: infrav1.ELBProtocolTCP,
					},
				},
			},
			mocks: func(m *mocks.MockEC2APIMockRecorder) {},
			expect: func(t *testing.T, g *WithT, res *infrav1.LoadBalancer) {
				t.Helper()
				if len(res.ELBListeners) != 2 {
					t.Errorf("Expected 2 listener to be configured, got %v listener(s)", len(res.ELBListeners))
				}
			},
		},
	}

	ctx := context.TODO()

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

			spec, err := s.getAPIServerLBSpec(ctx, clusterScope.Name(), clusterScope.ControlPlaneLoadBalancer())
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
				m.DescribeLoadBalancers(gomock.Any(), gomock.Eq(&elb.DescribeLoadBalancersInput{
					LoadBalancerNames: []string{elbName},
				})).
					Return(&elb.DescribeLoadBalancersOutput{
						LoadBalancerDescriptions: []elbtypes.LoadBalancerDescription{
							{
								LoadBalancerName: aws.String(elbName),
								Scheme:           aws.String(string(infrav1.ELBSchemeInternetFacing)),
								Subnets:          []string{clusterSubnetID},
							},
						},
					}, nil)
				m.DescribeLoadBalancerAttributes(gomock.Any(), gomock.Eq(&elb.DescribeLoadBalancerAttributesInput{
					LoadBalancerName: aws.String(elbName),
				})).
					Return(&elb.DescribeLoadBalancerAttributesOutput{
						LoadBalancerAttributes: &elbtypes.LoadBalancerAttributes{
							CrossZoneLoadBalancing: &elbtypes.CrossZoneLoadBalancing{
								Enabled: false,
							},
						},
					}, nil)
				m.DescribeTags(gomock.Any(), &elb.DescribeTagsInput{LoadBalancerNames: []string{elbName}}).Return(
					&elb.DescribeTagsOutput{
						TagDescriptions: []elbtypes.TagDescription{
							{
								LoadBalancerName: aws.String(elbName),
								Tags: []elbtypes.Tag{{
									Key:   aws.String(infrav1.ClusterTagKey(clusterName)),
									Value: aws.String(string(infrav1.ResourceLifecycleOwned)),
								}},
							},
						},
					}, nil)

				m.RegisterInstancesWithLoadBalancer(gomock.Any(), gomock.Eq(&elb.RegisterInstancesWithLoadBalancerInput{
					Instances:        []elbtypes.Instance{{InstanceId: aws.String(instanceID)}},
					LoadBalancerName: aws.String(elbName),
				})).
					Return(&elb.RegisterInstancesWithLoadBalancerOutput{
						Instances: []elbtypes.Instance{{InstanceId: aws.String(instanceID)}},
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
				m.DescribeLoadBalancers(gomock.Any(), gomock.Eq(&elb.DescribeLoadBalancersInput{
					LoadBalancerNames: []string{elbName},
				})).
					Return(&elb.DescribeLoadBalancersOutput{
						LoadBalancerDescriptions: []elbtypes.LoadBalancerDescription{
							{
								Scheme:            aws.String(string(infrav1.ELBSchemeInternetFacing)),
								Subnets:           []string{elbSubnetID},
								AvailabilityZones: []string{az},
							},
						},
					}, nil)
				m.DescribeLoadBalancerAttributes(gomock.Any(), gomock.Eq(&elb.DescribeLoadBalancerAttributesInput{
					LoadBalancerName: aws.String(elbName),
				})).
					Return(&elb.DescribeLoadBalancerAttributesOutput{
						LoadBalancerAttributes: &elbtypes.LoadBalancerAttributes{
							CrossZoneLoadBalancing: &elbtypes.CrossZoneLoadBalancing{
								Enabled: false,
							},
						},
					}, nil)
				m.DescribeTags(gomock.Any(), &elb.DescribeTagsInput{LoadBalancerNames: []string{elbName}}).Return(
					&elb.DescribeTagsOutput{
						TagDescriptions: []elbtypes.TagDescription{
							{
								LoadBalancerName: aws.String(elbName),
								Tags: []elbtypes.Tag{{
									Key:   aws.String(infrav1.ClusterTagKey(clusterName)),
									Value: aws.String(string(infrav1.ResourceLifecycleOwned)),
								}},
							},
						},
					}, nil)

				m.RegisterInstancesWithLoadBalancer(gomock.Any(), gomock.Eq(&elb.RegisterInstancesWithLoadBalancerInput{
					Instances:        []elbtypes.Instance{{InstanceId: aws.String(instanceID)}},
					LoadBalancerName: aws.String(elbName),
				})).
					Return(&elb.RegisterInstancesWithLoadBalancerOutput{
						Instances: []elbtypes.Instance{{InstanceId: aws.String(instanceID)}},
					}, nil)
			},
			ec2Mocks: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSubnets(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{
					SubnetIds: []string{
						elbSubnetID,
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []ec2types.Subnet{
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
				m.DescribeLoadBalancers(gomock.Any(), gomock.Eq(&elb.DescribeLoadBalancersInput{
					LoadBalancerNames: []string{elbName},
				})).
					Return(&elb.DescribeLoadBalancersOutput{
						LoadBalancerDescriptions: []elbtypes.LoadBalancerDescription{
							{
								Scheme:            aws.String(string(infrav1.ELBSchemeInternetFacing)),
								Subnets:           []string{elbSubnetID},
								AvailabilityZones: []string{differentAZ},
							},
						},
					}, nil)
				m.DescribeLoadBalancerAttributes(gomock.Any(), gomock.Eq(&elb.DescribeLoadBalancerAttributesInput{
					LoadBalancerName: aws.String(elbName),
				})).
					Return(&elb.DescribeLoadBalancerAttributesOutput{
						LoadBalancerAttributes: &elbtypes.LoadBalancerAttributes{
							CrossZoneLoadBalancing: &elbtypes.CrossZoneLoadBalancing{
								Enabled: false,
							},
						},
					}, nil)
				m.DescribeTags(gomock.Any(), &elb.DescribeTagsInput{LoadBalancerNames: []string{elbName}}).Return(
					&elb.DescribeTagsOutput{
						TagDescriptions: []elbtypes.TagDescription{
							{
								LoadBalancerName: aws.String(elbName),
								Tags: []elbtypes.Tag{{
									Key:   aws.String(infrav1.ClusterTagKey(clusterName)),
									Value: aws.String(string(infrav1.ResourceLifecycleOwned)),
								}},
							},
						},
					}, nil)
			},
			ec2Mocks: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSubnets(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{
					SubnetIds: []string{
						elbSubnetID,
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []ec2types.Subnet{
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

	ctx := context.TODO()

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

			err = s.RegisterInstanceWithAPIServerELB(ctx, instance)
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
		tgArn           = "arn::target-group"
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
				m.DescribeLoadBalancers(gomock.Any(), gomock.Eq(&elbv2.DescribeLoadBalancersInput{
					Names: []string{elbName},
				})).
					Return(&elbv2.DescribeLoadBalancersOutput{
						LoadBalancers: []elbv2types.LoadBalancer{
							{
								LoadBalancerArn:  aws.String(elbArn),
								LoadBalancerName: aws.String(elbName),
								Scheme:           SchemeToSDKScheme(infrav1.ELBSchemeInternetFacing),
								AvailabilityZones: []elbv2types.AvailabilityZone{
									{
										SubnetId: aws.String(clusterSubnetID),
									},
								},
							},
						},
					}, nil)
				m.DescribeLoadBalancerAttributes(gomock.Any(), gomock.Eq(&elbv2.DescribeLoadBalancerAttributesInput{
					LoadBalancerArn: aws.String(elbArn),
				})).
					Return(&elbv2.DescribeLoadBalancerAttributesOutput{
						Attributes: []elbv2types.LoadBalancerAttribute{
							{
								Key:   aws.String("load_balancing.cross_zone.enabled"),
								Value: aws.String("true"),
							},
						},
					}, nil)
				m.DescribeTags(gomock.Any(), &elbv2.DescribeTagsInput{ResourceArns: []string{elbArn}}).Return(
					&elbv2.DescribeTagsOutput{
						TagDescriptions: []elbv2types.TagDescription{
							{
								ResourceArn: aws.String(elbArn),
								Tags: []elbv2types.Tag{{
									Key:   aws.String(infrav1.ClusterTagKey(clusterName)),
									Value: aws.String(string(infrav1.ResourceLifecycleOwned)),
								}},
							},
						},
					}, nil)
				m.DescribeTargetGroups(gomock.Any(), &elbv2.DescribeTargetGroupsInput{
					LoadBalancerArn: aws.String(elbArn),
				}).Return(&elbv2.DescribeTargetGroupsOutput{
					TargetGroups: []elbv2types.TargetGroup{
						{
							HealthCheckEnabled:  aws.Bool(true),
							HealthCheckPort:     aws.String(infrav1.DefaultAPIServerPortString),
							HealthCheckProtocol: elbv2types.ProtocolEnumTcp,
							LoadBalancerArns:    []string{elbArn},
							Port:                aws.Int32(infrav1.DefaultAPIServerPort),
							Protocol:            elbv2types.ProtocolEnumTcp,
							TargetGroupArn:      aws.String(tgArn),
							TargetGroupName:     aws.String("something-generated"),
							VpcId:               aws.String("vpc-id"),
						},
					},
				}, nil)
				m.RegisterTargets(gomock.Any(), gomock.Eq(&elbv2.RegisterTargetsInput{
					TargetGroupArn: aws.String(tgArn),
					Targets: []elbv2types.TargetDescription{
						{
							Id:   aws.String(instanceID),
							Port: aws.Int32(infrav1.DefaultAPIServerPort),
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
			name: "multiple listeners",
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: clusterName},
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						Name:             aws.String(elbName),
						LoadBalancerType: infrav1.LoadBalancerTypeNLB,
						AdditionalListeners: []infrav1.AdditionalListenerSpec{
							{
								Port:     443,
								Protocol: infrav1.ELBProtocolTCP,
							},
							{
								Port:     8443,
								Protocol: infrav1.ELBProtocolTCP,
							},
						},
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
				m.DescribeLoadBalancers(gomock.Any(), gomock.Eq(&elbv2.DescribeLoadBalancersInput{
					Names: []string{elbName},
				})).
					Return(&elbv2.DescribeLoadBalancersOutput{
						LoadBalancers: []elbv2types.LoadBalancer{
							{
								LoadBalancerArn:  aws.String(elbArn),
								LoadBalancerName: aws.String(elbName),
								Scheme:           SchemeToSDKScheme(infrav1.ELBSchemeInternetFacing),
								AvailabilityZones: []elbv2types.AvailabilityZone{
									{
										SubnetId: aws.String(clusterSubnetID),
									},
								},
							},
						},
					}, nil)
				m.DescribeLoadBalancerAttributes(gomock.Any(), gomock.Eq(&elbv2.DescribeLoadBalancerAttributesInput{
					LoadBalancerArn: aws.String(elbArn),
				})).
					Return(&elbv2.DescribeLoadBalancerAttributesOutput{
						Attributes: []elbv2types.LoadBalancerAttribute{
							{
								Key:   aws.String("load_balancing.cross_zone.enabled"),
								Value: aws.String("true"),
							},
						},
					}, nil)
				m.DescribeTags(gomock.Any(), &elbv2.DescribeTagsInput{ResourceArns: []string{elbArn}}).Return(
					&elbv2.DescribeTagsOutput{
						TagDescriptions: []elbv2types.TagDescription{
							{
								ResourceArn: aws.String(elbArn),
								Tags: []elbv2types.Tag{{
									Key:   aws.String(infrav1.ClusterTagKey(clusterName)),
									Value: aws.String(string(infrav1.ResourceLifecycleOwned)),
								}},
							},
						},
					}, nil)
				m.DescribeTargetGroups(gomock.Any(), &elbv2.DescribeTargetGroupsInput{
					LoadBalancerArn: aws.String(elbArn),
				}).Return(&elbv2.DescribeTargetGroupsOutput{
					TargetGroups: []elbv2types.TargetGroup{
						{
							HealthCheckEnabled:  aws.Bool(true),
							HealthCheckPort:     aws.String(infrav1.DefaultAPIServerPortString),
							HealthCheckProtocol: elbv2types.ProtocolEnumTcp,
							LoadBalancerArns:    []string{elbArn},
							Port:                aws.Int32(infrav1.DefaultAPIServerPort),
							Protocol:            elbv2types.ProtocolEnumTcp,
							TargetGroupArn:      aws.String(tgArn),
							TargetGroupName:     aws.String("something-generated"),
							VpcId:               aws.String("vpc-id"),
						},
						{
							HealthCheckEnabled:  aws.Bool(true),
							HealthCheckPort:     aws.String("443"),
							HealthCheckProtocol: elbv2types.ProtocolEnumTcp,
							LoadBalancerArns:    []string{elbArn},
							Port:                aws.Int32(443),
							Protocol:            elbv2types.ProtocolEnumTcp,
							TargetGroupArn:      aws.String("target-group::arn::443"),
							TargetGroupName:     aws.String("something-generated-443"),
							VpcId:               aws.String("vpc-id"),
						},
						{
							HealthCheckEnabled:  aws.Bool(true),
							HealthCheckPort:     aws.String("8443"),
							HealthCheckProtocol: elbv2types.ProtocolEnumTcp,
							LoadBalancerArns:    []string{elbArn},
							Port:                aws.Int32(8443),
							Protocol:            elbv2types.ProtocolEnumTcp,
							TargetGroupArn:      aws.String("target-group::arn::8443"),
							TargetGroupName:     aws.String("something-generated-8443"),
							VpcId:               aws.String("vpc-id"),
						},
					},
				}, nil)
				m.RegisterTargets(gomock.Any(), gomock.Eq(&elbv2.RegisterTargetsInput{
					TargetGroupArn: aws.String(tgArn),
					Targets: []elbv2types.TargetDescription{
						{
							Id:   aws.String(instanceID),
							Port: aws.Int32(infrav1.DefaultAPIServerPort),
						},
					},
				})).Return(&elbv2.RegisterTargetsOutput{}, nil)
				m.RegisterTargets(gomock.Any(), gomock.Eq(&elbv2.RegisterTargetsInput{
					TargetGroupArn: aws.String("target-group::arn::443"),
					Targets: []elbv2types.TargetDescription{
						{
							Id:   aws.String(instanceID),
							Port: aws.Int32(443),
						},
					},
				})).Return(&elbv2.RegisterTargetsOutput{}, nil)
				m.RegisterTargets(gomock.Any(), gomock.Eq(&elbv2.RegisterTargetsInput{
					TargetGroupArn: aws.String("target-group::arn::8443"),
					Targets: []elbv2types.TargetDescription{
						{
							Id:   aws.String(instanceID),
							Port: aws.Int32(8443),
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
				m.DescribeLoadBalancers(gomock.Any(), gomock.Eq(&elbv2.DescribeLoadBalancersInput{
					Names: []string{elbName},
				})).
					Return(&elbv2.DescribeLoadBalancersOutput{
						LoadBalancers: []elbv2types.LoadBalancer{
							{
								LoadBalancerArn:  aws.String(elbArn),
								LoadBalancerName: aws.String(elbName),
								Scheme:           SchemeToSDKScheme(infrav1.ELBSchemeInternetFacing),
								AvailabilityZones: []elbv2types.AvailabilityZone{
									{
										SubnetId: aws.String(clusterSubnetID),
									},
								},
							},
						},
					}, nil)
				m.DescribeLoadBalancerAttributes(gomock.Any(), gomock.Eq(&elbv2.DescribeLoadBalancerAttributesInput{
					LoadBalancerArn: aws.String(elbArn),
				})).
					Return(&elbv2.DescribeLoadBalancerAttributesOutput{
						Attributes: []elbv2types.LoadBalancerAttribute{
							{
								Key:   aws.String("load_balancing.cross_zone.enabled"),
								Value: aws.String("true"),
							},
						},
					}, nil)
				m.DescribeTags(gomock.Any(), &elbv2.DescribeTagsInput{ResourceArns: []string{elbArn}}).Return(
					&elbv2.DescribeTagsOutput{
						TagDescriptions: []elbv2types.TagDescription{
							{
								ResourceArn: aws.String(elbArn),
								Tags: []elbv2types.Tag{{
									Key:   aws.String(infrav1.ClusterTagKey(clusterName)),
									Value: aws.String(string(infrav1.ResourceLifecycleOwned)),
								}},
							},
						},
					}, nil)
				m.DescribeTargetGroups(gomock.Any(), &elbv2.DescribeTargetGroupsInput{
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

	ctx := context.TODO()

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

			err = s.RegisterInstanceWithAPIServerLB(ctx, instance, clusterScope.ControlPlaneLoadBalancer())
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
				m.CreateLoadBalancer(gomock.Any(), gomock.Eq(&elbv2.CreateLoadBalancerInput{
					Name:    aws.String(elbName),
					Scheme:  elbv2types.LoadBalancerSchemeEnumInternetFacing,
					Type:    elbv2types.LoadBalancerTypeEnumNetwork,
					Subnets: []string{clusterSubnetID},
					Tags: []elbv2types.Tag{
						{
							Key:   aws.String("test"),
							Value: aws.String("tag"),
						},
					},
				})).Return(&elbv2.CreateLoadBalancerOutput{
					LoadBalancers: []elbv2types.LoadBalancer{
						{
							LoadBalancerArn:  aws.String(elbArn),
							LoadBalancerName: aws.String(elbName),
							Scheme:           SchemeToSDKScheme(infrav1.ELBSchemeInternetFacing),
							DNSName:          aws.String(dns),
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
				m.CreateLoadBalancer(gomock.Any(), gomock.Eq(&elbv2.CreateLoadBalancerInput{
					Name:          aws.String(elbName),
					IpAddressType: elbv2types.IpAddressTypeDualstack,
					Scheme:        elbv2types.LoadBalancerSchemeEnumInternetFacing,
					Type:          elbv2types.LoadBalancerTypeEnumNetwork,
					Subnets:       []string{clusterSubnetID},
					Tags: []elbv2types.Tag{
						{
							Key:   aws.String("test"),
							Value: aws.String("tag"),
						},
					},
				})).Return(&elbv2.CreateLoadBalancerOutput{
					LoadBalancers: []elbv2types.LoadBalancer{
						{
							LoadBalancerArn:  aws.String(elbArn),
							LoadBalancerName: aws.String(elbName),
							Scheme:           SchemeToSDKScheme(infrav1.ELBSchemeInternetFacing),
							DNSName:          aws.String(dns),
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
				m.CreateLoadBalancer(gomock.Any(), gomock.Eq(&elbv2.CreateLoadBalancerInput{
					Name:    aws.String(elbName),
					Scheme:  elbv2types.LoadBalancerSchemeEnumInternetFacing,
					Type:    elbv2types.LoadBalancerTypeEnumNetwork,
					Subnets: []string{clusterSubnetID},
					Tags: []elbv2types.Tag{
						{
							Key:   aws.String("test"),
							Value: aws.String("tag"),
						},
					},
				})).Return(nil, errors.New("nope"))
			},
			check: func(t *testing.T, _ *infrav1.LoadBalancer, err error) {
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
			name: "PreserveClientIP is enabled",
			spec: func(spec infrav1.LoadBalancer) infrav1.LoadBalancer {
				return spec
			},
			awsCluster: func(acl infrav1.AWSCluster) infrav1.AWSCluster {
				acl.Spec.ControlPlaneLoadBalancer.PreserveClientIP = true
				return acl
			},
			elbV2APIMocks: func(m *mocks.MockELBV2APIMockRecorder) {
				m.CreateLoadBalancer(gomock.Any(), gomock.Eq(&elbv2.CreateLoadBalancerInput{
					Name:    aws.String(elbName),
					Scheme:  elbv2types.LoadBalancerSchemeEnumInternetFacing,
					Type:    elbv2types.LoadBalancerTypeEnumNetwork,
					Subnets: []string{clusterSubnetID},
					Tags: []elbv2types.Tag{
						{
							Key:   aws.String("test"),
							Value: aws.String("tag"),
						},
					},
				})).Return(&elbv2.CreateLoadBalancerOutput{
					LoadBalancers: []elbv2types.LoadBalancer{
						{
							LoadBalancerArn:  aws.String(elbArn),
							LoadBalancerName: aws.String(elbName),
							Scheme:           SchemeToSDKScheme(infrav1.ELBSchemeInternetFacing),
							DNSName:          aws.String(dns),
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
				m.CreateLoadBalancer(gomock.Any(), gomock.Eq(&elbv2.CreateLoadBalancerInput{
					Name:    aws.String(elbName),
					Scheme:  elbv2types.LoadBalancerSchemeEnumInternetFacing,
					Type:    elbv2types.LoadBalancerTypeEnumApplication,
					Subnets: []string{clusterSubnetID},
					Tags: []elbv2types.Tag{
						{
							Key:   aws.String("test"),
							Value: aws.String("tag"),
						},
					},
					SecurityGroups: []string{"sg-id"},
				})).Return(&elbv2.CreateLoadBalancerOutput{
					LoadBalancers: []elbv2types.LoadBalancer{
						{
							LoadBalancerArn:  aws.String(elbArn),
							LoadBalancerName: aws.String(elbName),
							Scheme:           SchemeToSDKScheme(infrav1.ELBSchemeInternetFacing),
							DNSName:          aws.String(dns),
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

	ctx := context.TODO()

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
								Port:     aws.String(infrav1.DefaultAPIServerPortString),
							},
						},
					},
				},
				LoadBalancerType: infrav1.LoadBalancerTypeNLB,
				SubnetIDs:        []string{clusterSubnetID},
			}

			spec := tc.spec(*loadBalancerSpec)
			lb, err := s.createLB(ctx, &spec, clusterScope.ControlPlaneLoadBalancer())
			tc.check(t, lb, err)
		})
	}
}

func TestReconcileTargetGroupsAndListeners(t *testing.T) {
	const (
		namespace       = "foo"
		clusterName     = "bar"
		clusterSubnetID = "subnet-1"
		elbName         = "bar-apiserver"
		elbArn          = "arn::apiserver"
		tgArn           = "arn::target-group"
		vpcID           = "vpc-id"
		dns             = "asdf:9999/asdf"
	)

	tests := []struct {
		name          string
		elbV2APIMocks func(m *mocks.MockELBV2APIMockRecorder)
		check         func(t *testing.T, tgs []*elbv2types.TargetGroup, listeners []*elbv2types.Listener, err error)
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
				m.DescribeTargetGroups(gomock.Any(), gomock.Eq(&elbv2.DescribeTargetGroupsInput{
					LoadBalancerArn: aws.String(elbArn),
				})).Return(&elbv2.DescribeTargetGroupsOutput{
					TargetGroups: []elbv2types.TargetGroup{},
				}, nil)
				m.CreateTargetGroup(gomock.Any(), gomock.Eq(&elbv2.CreateTargetGroupInput{
					Name:     aws.String("name"),
					Port:     aws.Int32(infrav1.DefaultAPIServerPort),
					Protocol: elbv2types.ProtocolEnumTcp,
					VpcId:    aws.String(vpcID),
					Tags: []elbv2types.Tag{
						{
							Key:   aws.String("test"),
							Value: aws.String("tag"),
						},
					},
					HealthCheckEnabled:         aws.Bool(true),
					HealthCheckPort:            aws.String(infrav1.DefaultAPIServerPortString),
					HealthCheckProtocol:        elbv2types.ProtocolEnumTcp,
					HealthyThresholdCount:      aws.Int32(infrav1.DefaultAPIServerHealthThresholdCount),
					UnhealthyThresholdCount:    aws.Int32(infrav1.DefaultAPIServerUnhealthThresholdCount),
					HealthCheckIntervalSeconds: aws.Int32(infrav1.DefaultAPIServerHealthCheckIntervalSec),
					HealthCheckTimeoutSeconds:  aws.Int32(infrav1.DefaultAPIServerHealthCheckTimeoutSec),
				})).Return(&elbv2.CreateTargetGroupOutput{
					TargetGroups: []elbv2types.TargetGroup{
						{
							TargetGroupArn:             aws.String(tgArn),
							TargetGroupName:            aws.String("name"),
							VpcId:                      aws.String(vpcID),
							HealthyThresholdCount:      aws.Int32(infrav1.DefaultAPIServerHealthThresholdCount),
							UnhealthyThresholdCount:    aws.Int32(infrav1.DefaultAPIServerUnhealthThresholdCount),
							HealthCheckIntervalSeconds: aws.Int32(infrav1.DefaultAPIServerHealthCheckIntervalSec),
							HealthCheckTimeoutSeconds:  aws.Int32(infrav1.DefaultAPIServerHealthCheckTimeoutSec),
						},
					},
				}, nil)
				m.ModifyTargetGroupAttributes(gomock.Any(), gomock.Eq(&elbv2.ModifyTargetGroupAttributesInput{
					TargetGroupArn: aws.String(tgArn),
					Attributes: []elbv2types.TargetGroupAttribute{
						{
							Key:   aws.String(infrav1.TargetGroupAttributeEnableConnectionTermination),
							Value: aws.String("false"),
						},
						{
							Key:   aws.String(infrav1.TargetGroupAttributeUnhealthyDrainingIntervalSeconds),
							Value: aws.String("300"),
						},
						{
							Key:   aws.String(infrav1.TargetGroupAttributeEnablePreserveClientIP),
							Value: aws.String("false"),
						},
					},
				})).Return(nil, nil)
				m.DescribeListeners(gomock.Any(), gomock.Eq(&elbv2.DescribeListenersInput{
					LoadBalancerArn: aws.String(elbArn),
				})).Return(&elbv2.DescribeListenersOutput{
					Listeners: []elbv2types.Listener{},
				}, nil)
				m.CreateListener(gomock.Any(), gomock.Eq(&elbv2.CreateListenerInput{
					DefaultActions: []elbv2types.Action{
						{
							TargetGroupArn: aws.String(tgArn),
							Type:           elbv2types.ActionTypeEnumForward,
						},
					},
					LoadBalancerArn: aws.String(elbArn),
					Port:            aws.Int32(infrav1.DefaultAPIServerPort),
					Protocol:        elbv2types.ProtocolEnumTcp,
					Tags: []elbv2types.Tag{
						{
							Key:   aws.String("test"),
							Value: aws.String("tag"),
						},
					},
				})).Return(&elbv2.CreateListenerOutput{
					Listeners: []elbv2types.Listener{
						{
							DefaultActions: []elbv2types.Action{
								{
									TargetGroupArn: aws.String(tgArn),
									Type:           elbv2types.ActionTypeEnumForward,
								},
							},
							ListenerArn: aws.String("listener::arn"),
							Port:        aws.Int32(infrav1.DefaultAPIServerPort),
							Protocol:    elbv2types.ProtocolEnumTcp,
						},
					},
				}, nil)
			},
			check: func(t *testing.T, tgs []*elbv2types.TargetGroup, listeners []*elbv2types.Listener, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}
				if len(tgs) != 1 {
					t.Fatalf("no target groups created")
				}
				if len(listeners) != 1 {
					t.Fatalf("no listeners created")
				}

				if len(listeners[0].DefaultActions) != 1 {
					t.Fatalf("no default actions created")
				}

				if *tgs[0].TargetGroupArn != *listeners[0].DefaultActions[0].TargetGroupArn {
					t.Fatalf("target group and listener did not have matching arns. target group ARN: %q. listener's target group ARN: %q", *tgs[0].TargetGroupArn, *listeners[0].DefaultActions[0].TargetGroupArn)
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
				m.DescribeTargetGroups(gomock.Any(), &elbv2.DescribeTargetGroupsInput{
					LoadBalancerArn: aws.String(elbArn),
				}).Return(&elbv2.DescribeTargetGroupsOutput{
					TargetGroups: []elbv2types.TargetGroup{},
				}, nil)
				m.CreateTargetGroup(gomock.Any(), &elbv2.CreateTargetGroupInput{
					Name:          aws.String("name"),
					Port:          aws.Int32(infrav1.DefaultAPIServerPort),
					Protocol:      elbv2types.ProtocolEnumTcp,
					VpcId:         aws.String(vpcID),
					IpAddressType: elbv2types.TargetGroupIpAddressTypeEnumIpv6,
					Tags: []elbv2types.Tag{
						{
							Key:   aws.String("test"),
							Value: aws.String("tag"),
						},
					},
					HealthCheckEnabled:         aws.Bool(true),
					HealthCheckPort:            aws.String(infrav1.DefaultAPIServerPortString),
					HealthCheckProtocol:        elbv2types.ProtocolEnumTcp,
					HealthyThresholdCount:      aws.Int32(infrav1.DefaultAPIServerHealthThresholdCount),
					UnhealthyThresholdCount:    aws.Int32(infrav1.DefaultAPIServerUnhealthThresholdCount),
					HealthCheckIntervalSeconds: aws.Int32(infrav1.DefaultAPIServerHealthCheckIntervalSec),
					HealthCheckTimeoutSeconds:  aws.Int32(infrav1.DefaultAPIServerHealthCheckTimeoutSec),
				}).Return(&elbv2.CreateTargetGroupOutput{
					TargetGroups: []elbv2types.TargetGroup{
						{
							TargetGroupArn:             aws.String(tgArn),
							TargetGroupName:            aws.String("name"),
							VpcId:                      aws.String(vpcID),
							HealthyThresholdCount:      aws.Int32(infrav1.DefaultAPIServerHealthThresholdCount),
							UnhealthyThresholdCount:    aws.Int32(infrav1.DefaultAPIServerUnhealthThresholdCount),
							HealthCheckIntervalSeconds: aws.Int32(infrav1.DefaultAPIServerHealthCheckIntervalSec),
							HealthCheckTimeoutSeconds:  aws.Int32(infrav1.DefaultAPIServerHealthCheckTimeoutSec),
							IpAddressType:              elbv2types.TargetGroupIpAddressTypeEnumIpv6,
						},
					},
				}, nil)
				m.ModifyTargetGroupAttributes(gomock.Any(), &elbv2.ModifyTargetGroupAttributesInput{
					TargetGroupArn: aws.String(tgArn),
					Attributes: []elbv2types.TargetGroupAttribute{
						{
							Key:   aws.String(infrav1.TargetGroupAttributeEnableConnectionTermination),
							Value: aws.String("false"),
						},
						{
							Key:   aws.String(infrav1.TargetGroupAttributeUnhealthyDrainingIntervalSeconds),
							Value: aws.String("300"),
						},
						{
							Key:   aws.String(infrav1.TargetGroupAttributeEnablePreserveClientIP),
							Value: aws.String("false"),
						},
					},
				}).Return(nil, nil)
				m.DescribeListeners(gomock.Any(), &elbv2.DescribeListenersInput{
					LoadBalancerArn: aws.String(elbArn),
				}).Return(&elbv2.DescribeListenersOutput{
					Listeners: []elbv2types.Listener{},
				}, nil)
				m.CreateListener(gomock.Any(), &elbv2.CreateListenerInput{
					DefaultActions: []elbv2types.Action{
						{
							TargetGroupArn: aws.String(tgArn),
							Type:           elbv2types.ActionTypeEnumForward,
						},
					},
					LoadBalancerArn: aws.String(elbArn),
					Port:            aws.Int32(infrav1.DefaultAPIServerPort),
					Protocol:        elbv2types.ProtocolEnumTcp,
					Tags: []elbv2types.Tag{
						{
							Key:   aws.String("test"),
							Value: aws.String("tag"),
						},
					},
				}).Return(&elbv2.CreateListenerOutput{
					Listeners: []elbv2types.Listener{
						{
							ListenerArn: aws.String("listener::arn"),
						},
					},
				}, nil)
			},
			check: func(t *testing.T, tgs []*elbv2types.TargetGroup, _ []*elbv2types.Listener, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}
				tg := tgs[0]
				got := tg.IpAddressType
				want := elbv2types.TargetGroupIpAddressTypeEnumIpv6
				if got != want {
					t.Fatalf("did not set ip address type to ipv6")
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
				m.DescribeTargetGroups(gomock.Any(), &elbv2.DescribeTargetGroupsInput{
					LoadBalancerArn: aws.String(elbArn),
				}).Return(&elbv2.DescribeTargetGroupsOutput{
					TargetGroups: []elbv2types.TargetGroup{},
				}, nil)
				m.CreateTargetGroup(gomock.Any(), &elbv2.CreateTargetGroupInput{
					Name:     aws.String("name"),
					Port:     aws.Int32(infrav1.DefaultAPIServerPort),
					Protocol: elbv2types.ProtocolEnumTcp,
					VpcId:    aws.String(vpcID),
					Tags: []elbv2types.Tag{
						{
							Key:   aws.String("test"),
							Value: aws.String("tag"),
						},
					},
					HealthyThresholdCount:      aws.Int32(infrav1.DefaultAPIServerHealthThresholdCount),
					UnhealthyThresholdCount:    aws.Int32(infrav1.DefaultAPIServerUnhealthThresholdCount),
					HealthCheckIntervalSeconds: aws.Int32(infrav1.DefaultAPIServerHealthCheckIntervalSec),
					HealthCheckTimeoutSeconds:  aws.Int32(infrav1.DefaultAPIServerHealthCheckTimeoutSec),
				}).Return(&elbv2.CreateTargetGroupOutput{
					TargetGroups: []elbv2types.TargetGroup{
						{
							TargetGroupArn:             aws.String(tgArn),
							TargetGroupName:            aws.String("name"),
							VpcId:                      aws.String(vpcID),
							HealthyThresholdCount:      aws.Int32(infrav1.DefaultAPIServerHealthThresholdCount),
							UnhealthyThresholdCount:    aws.Int32(infrav1.DefaultAPIServerUnhealthThresholdCount),
							HealthCheckIntervalSeconds: aws.Int32(infrav1.DefaultAPIServerHealthCheckIntervalSec),
							HealthCheckTimeoutSeconds:  aws.Int32(infrav1.DefaultAPIServerHealthCheckTimeoutSec),
							HealthCheckEnabled:         aws.Bool(false),
						},
					},
				}, nil)
				m.ModifyTargetGroupAttributes(gomock.Any(), &elbv2.ModifyTargetGroupAttributesInput{
					TargetGroupArn: aws.String(tgArn),
					Attributes: []elbv2types.TargetGroupAttribute{
						{
							Key:   aws.String(infrav1.TargetGroupAttributeEnableConnectionTermination),
							Value: aws.String("false"),
						},
						{
							Key:   aws.String(infrav1.TargetGroupAttributeUnhealthyDrainingIntervalSeconds),
							Value: aws.String("300"),
						},
						{
							Key:   aws.String(infrav1.TargetGroupAttributeEnablePreserveClientIP),
							Value: aws.String("false"),
						},
					},
				}).Return(nil, nil)
				m.DescribeListeners(gomock.Any(), &elbv2.DescribeListenersInput{
					LoadBalancerArn: aws.String(elbArn),
				}).Return(&elbv2.DescribeListenersOutput{
					Listeners: []elbv2types.Listener{},
				}, nil)
				m.CreateListener(gomock.Any(), &elbv2.CreateListenerInput{
					DefaultActions: []elbv2types.Action{
						{
							TargetGroupArn: aws.String(tgArn),
							Type:           elbv2types.ActionTypeEnumForward,
						},
					},
					LoadBalancerArn: aws.String(elbArn),
					Port:            aws.Int32(infrav1.DefaultAPIServerPort),
					Protocol:        elbv2types.ProtocolEnumTcp,
					Tags: []elbv2types.Tag{
						{
							Key:   aws.String("test"),
							Value: aws.String("tag"),
						},
					},
				}).Return(&elbv2.CreateListenerOutput{
					Listeners: []elbv2types.Listener{
						{
							ListenerArn: aws.String("listener::arn"),
						},
					},
				}, nil)
			},
			check: func(t *testing.T, tgs []*elbv2types.TargetGroup, _ []*elbv2types.Listener, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}
				got := *tgs[0].HealthCheckEnabled
				want := false
				if got != want {
					t.Fatalf("health check not disabled on target group")
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
				m.DescribeTargetGroups(gomock.Any(), &elbv2.DescribeTargetGroupsInput{
					LoadBalancerArn: aws.String(elbArn),
				}).Return(&elbv2.DescribeTargetGroupsOutput{
					TargetGroups: []elbv2types.TargetGroup{},
				}, nil)
				m.CreateTargetGroup(gomock.Any(), &elbv2.CreateTargetGroupInput{
					HealthCheckEnabled:  aws.Bool(true),
					HealthCheckPort:     aws.String(infrav1.DefaultAPIServerPortString),
					HealthCheckProtocol: elbv2types.ProtocolEnumTcp,
					Name:                aws.String("name"),
					Port:                aws.Int32(infrav1.DefaultAPIServerPort),
					Protocol:            elbv2types.ProtocolEnumTcp,
					VpcId:               aws.String(vpcID),
					Tags: []elbv2types.Tag{
						{
							Key:   aws.String("test"),
							Value: aws.String("tag"),
						},
					},
					HealthyThresholdCount:      aws.Int32(infrav1.DefaultAPIServerHealthThresholdCount),
					UnhealthyThresholdCount:    aws.Int32(infrav1.DefaultAPIServerUnhealthThresholdCount),
					HealthCheckIntervalSeconds: aws.Int32(infrav1.DefaultAPIServerHealthCheckIntervalSec),
					HealthCheckTimeoutSeconds:  aws.Int32(infrav1.DefaultAPIServerHealthCheckTimeoutSec),
				}).Return(&elbv2.CreateTargetGroupOutput{
					TargetGroups: []elbv2types.TargetGroup{
						{
							TargetGroupArn:             aws.String(tgArn),
							TargetGroupName:            aws.String("name"),
							VpcId:                      aws.String(vpcID),
							HealthyThresholdCount:      aws.Int32(infrav1.DefaultAPIServerHealthThresholdCount),
							UnhealthyThresholdCount:    aws.Int32(infrav1.DefaultAPIServerUnhealthThresholdCount),
							HealthCheckIntervalSeconds: aws.Int32(infrav1.DefaultAPIServerHealthCheckIntervalSec),
							HealthCheckTimeoutSeconds:  aws.Int32(infrav1.DefaultAPIServerHealthCheckTimeoutSec),
						},
					},
				}, nil)
				m.ModifyTargetGroupAttributes(gomock.Any(), gomock.Eq(&elbv2.ModifyTargetGroupAttributesInput{
					TargetGroupArn: aws.String(tgArn),
					Attributes: []elbv2types.TargetGroupAttribute{
						{
							Key:   aws.String(infrav1.TargetGroupAttributeEnableConnectionTermination),
							Value: aws.String("false"),
						},
						{
							Key:   aws.String(infrav1.TargetGroupAttributeUnhealthyDrainingIntervalSeconds),
							Value: aws.String("300"),
						},
					},
				})).Return(nil, nil)
				m.DescribeListeners(gomock.Any(), &elbv2.DescribeListenersInput{
					LoadBalancerArn: aws.String(elbArn),
				}).Return(&elbv2.DescribeListenersOutput{
					Listeners: []elbv2types.Listener{},
				}, nil)
				m.CreateListener(gomock.Any(), &elbv2.CreateListenerInput{
					DefaultActions: []elbv2types.Action{
						{
							TargetGroupArn: aws.String(tgArn),
							Type:           elbv2types.ActionTypeEnumForward,
						},
					},
					LoadBalancerArn: aws.String(elbArn),
					Port:            aws.Int32(infrav1.DefaultAPIServerPort),
					Protocol:        elbv2types.ProtocolEnumTcp,
					Tags: []elbv2types.Tag{
						{
							Key:   aws.String("test"),
							Value: aws.String("tag"),
						},
					},
				}).Return(&elbv2.CreateListenerOutput{
					Listeners: []elbv2types.Listener{
						{
							ListenerArn: aws.String("listener::arn"),
						},
					},
				}, nil)
			},
			check: func(t *testing.T, tgs []*elbv2types.TargetGroup, listeners []*elbv2types.Listener, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}

				if len(tgs) != 1 {
					t.Fatalf("did not create target groups")
				}

				if len(listeners) != 1 {
					t.Fatalf("did not create any listeners")
				}
			},
		},
		{
			name: "NLB with HTTP health check",
			awsCluster: func(acl infrav1.AWSCluster) infrav1.AWSCluster {
				acl.Spec.ControlPlaneLoadBalancer.Scheme = &infrav1.ELBSchemeInternetFacing
				acl.Spec.ControlPlaneLoadBalancer.LoadBalancerType = infrav1.LoadBalancerTypeNLB
				acl.Spec.ControlPlaneLoadBalancer.HealthCheckProtocol = &infrav1.ELBProtocolHTTP
				return acl
			},
			spec: func(spec infrav1.LoadBalancer) infrav1.LoadBalancer {
				tg := stubInfraV1TargetGroupSpecAPI
				tg.VpcID = vpcID
				tg.HealthCheck.Protocol = aws.String("HTTP")
				tg.HealthCheck.Port = aws.String(infrav1.DefaultAPIServerPortString)
				tg.HealthCheck.Path = aws.String("/readyz")
				spec.ELBListeners = []infrav1.Listener{
					{
						Protocol:    "TCP",
						Port:        infrav1.DefaultAPIServerPort,
						TargetGroup: tg,
					},
				}
				return spec
			},
			elbV2APIMocks: func(m *mocks.MockELBV2APIMockRecorder) {
				m.DescribeTargetGroups(gomock.Any(), &elbv2.DescribeTargetGroupsInput{
					LoadBalancerArn: aws.String(elbArn),
				}).Return(&elbv2.DescribeTargetGroupsOutput{
					TargetGroups: []elbv2types.TargetGroup{},
				}, nil)
				m.CreateTargetGroup(gomock.Any(), &elbv2.CreateTargetGroupInput{
					Name:                       aws.String("name"),
					Port:                       aws.Int32(infrav1.DefaultAPIServerPort),
					Protocol:                   elbv2types.ProtocolEnumTcp,
					VpcId:                      aws.String(vpcID),
					HealthCheckEnabled:         aws.Bool(true),
					HealthCheckPort:            aws.String(infrav1.DefaultAPIServerPortString),
					HealthCheckProtocol:        elbv2types.ProtocolEnumHttp,
					HealthCheckPath:            aws.String("/readyz"),
					HealthCheckIntervalSeconds: aws.Int32(10),
					HealthCheckTimeoutSeconds:  aws.Int32(5),
					HealthyThresholdCount:      aws.Int32(5),
					UnhealthyThresholdCount:    aws.Int32(3),
					Tags: []elbv2types.Tag{
						{
							Key:   aws.String("test"),
							Value: aws.String("tag"),
						},
					},
				}).Return(&elbv2.CreateTargetGroupOutput{
					TargetGroups: []elbv2types.TargetGroup{
						{
							TargetGroupArn:             aws.String(tgArn),
							TargetGroupName:            aws.String("name"),
							VpcId:                      aws.String(vpcID),
							HealthCheckEnabled:         aws.Bool(true),
							HealthCheckPort:            aws.String(infrav1.DefaultAPIServerPortString),
							HealthCheckProtocol:        elbv2types.ProtocolEnumHttp,
							HealthCheckPath:            aws.String("/readyz"),
							HealthCheckIntervalSeconds: aws.Int32(10),
							HealthCheckTimeoutSeconds:  aws.Int32(5),
							HealthyThresholdCount:      aws.Int32(5),
							UnhealthyThresholdCount:    aws.Int32(3),
						},
					},
				}, nil)
				m.DescribeListeners(gomock.Any(), &elbv2.DescribeListenersInput{
					LoadBalancerArn: aws.String(elbArn),
				}).Return(&elbv2.DescribeListenersOutput{
					Listeners: []elbv2types.Listener{},
				}, nil)
				m.CreateListener(gomock.Any(), &elbv2.CreateListenerInput{
					DefaultActions: []elbv2types.Action{
						{
							TargetGroupArn: aws.String(tgArn),
							Type:           elbv2types.ActionTypeEnumForward,
						},
					},
					LoadBalancerArn: aws.String(elbArn),
					Port:            aws.Int32(infrav1.DefaultAPIServerPort),
					Protocol:        elbv2types.ProtocolEnumTcp,
					Tags: []elbv2types.Tag{
						{
							Key:   aws.String("test"),
							Value: aws.String("tag"),
						},
					},
				}).Return(&elbv2.CreateListenerOutput{
					Listeners: []elbv2types.Listener{
						{
							ListenerArn: aws.String("listener::arn"),
						},
					},
				}, nil)
				m.ModifyTargetGroupAttributes(gomock.Any(), &elbv2.ModifyTargetGroupAttributesInput{
					TargetGroupArn: aws.String(tgArn),
					Attributes: []elbv2types.TargetGroupAttribute{
						{
							Key:   aws.String(infrav1.TargetGroupAttributeEnableConnectionTermination),
							Value: aws.String("false"),
						},
						{
							Key:   aws.String(infrav1.TargetGroupAttributeUnhealthyDrainingIntervalSeconds),
							Value: aws.String("300"),
						},
						{
							Key:   aws.String(infrav1.TargetGroupAttributeEnablePreserveClientIP),
							Value: aws.String("false"),
						},
					},
				}).Return(nil, nil)
			},
			check: func(t *testing.T, tgs []*elbv2types.TargetGroup, _ []*elbv2types.Listener, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}
				got := tgs[0].HealthCheckProtocol
				want := elbv2types.ProtocolEnumHttp
				if got != want {
					t.Fatalf("Health Check protocol for the API Target group did not equal expected value: %s; was: '%s'", want, got)
				}
			},
		},
		{
			name: "NLB with HTTPS health check",
			awsCluster: func(acl infrav1.AWSCluster) infrav1.AWSCluster {
				acl.Spec.ControlPlaneLoadBalancer.Scheme = &infrav1.ELBSchemeInternetFacing
				acl.Spec.ControlPlaneLoadBalancer.LoadBalancerType = infrav1.LoadBalancerTypeNLB
				acl.Spec.ControlPlaneLoadBalancer.HealthCheckProtocol = &infrav1.ELBProtocolHTTPS
				return acl
			},
			spec: func(spec infrav1.LoadBalancer) infrav1.LoadBalancer {
				tg := stubInfraV1TargetGroupSpecAPI
				tg.VpcID = vpcID
				tg.HealthCheck.Protocol = aws.String("HTTPS")
				tg.HealthCheck.Port = aws.String(infrav1.DefaultAPIServerPortString)
				tg.HealthCheck.Path = aws.String("/readyz")
				spec.ELBListeners = []infrav1.Listener{
					{
						Protocol:    "TCP",
						Port:        infrav1.DefaultAPIServerPort,
						TargetGroup: tg,
					},
				}
				return spec
			},
			elbV2APIMocks: func(m *mocks.MockELBV2APIMockRecorder) {
				m.DescribeTargetGroups(gomock.Any(), &elbv2.DescribeTargetGroupsInput{
					LoadBalancerArn: aws.String(elbArn),
				}).Return(&elbv2.DescribeTargetGroupsOutput{
					TargetGroups: []elbv2types.TargetGroup{},
				}, nil)
				m.CreateTargetGroup(gomock.Any(), &elbv2.CreateTargetGroupInput{
					Name:                       aws.String("name"),
					Port:                       aws.Int32(infrav1.DefaultAPIServerPort),
					Protocol:                   elbv2types.ProtocolEnumTcp,
					VpcId:                      aws.String(vpcID),
					HealthCheckEnabled:         aws.Bool(true),
					HealthCheckPort:            aws.String(infrav1.DefaultAPIServerPortString),
					HealthCheckProtocol:        elbv2types.ProtocolEnumHttps,
					HealthCheckPath:            aws.String("/readyz"),
					HealthCheckIntervalSeconds: aws.Int32(10),
					HealthCheckTimeoutSeconds:  aws.Int32(5),
					HealthyThresholdCount:      aws.Int32(5),
					UnhealthyThresholdCount:    aws.Int32(3),
					Tags: []elbv2types.Tag{
						{
							Key:   aws.String("test"),
							Value: aws.String("tag"),
						},
					},
				}).Return(&elbv2.CreateTargetGroupOutput{
					TargetGroups: []elbv2types.TargetGroup{
						{
							TargetGroupArn:             aws.String(tgArn),
							TargetGroupName:            aws.String("name"),
							VpcId:                      aws.String(vpcID),
							HealthCheckEnabled:         aws.Bool(true),
							HealthCheckPort:            aws.String(infrav1.DefaultAPIServerPortString),
							HealthCheckProtocol:        elbv2types.ProtocolEnumHttps,
							HealthCheckPath:            aws.String("/readyz"),
							HealthCheckIntervalSeconds: aws.Int32(10),
							HealthCheckTimeoutSeconds:  aws.Int32(5),
							HealthyThresholdCount:      aws.Int32(5),
							UnhealthyThresholdCount:    aws.Int32(3),
						},
					},
				}, nil)
				m.DescribeListeners(gomock.Any(), &elbv2.DescribeListenersInput{
					LoadBalancerArn: aws.String(elbArn),
				}).Return(&elbv2.DescribeListenersOutput{
					Listeners: []elbv2types.Listener{},
				}, nil)
				m.CreateListener(gomock.Any(), &elbv2.CreateListenerInput{
					DefaultActions: []elbv2types.Action{
						{
							TargetGroupArn: aws.String(tgArn),
							Type:           elbv2types.ActionTypeEnumForward,
						},
					},
					LoadBalancerArn: aws.String(elbArn),
					Port:            aws.Int32(infrav1.DefaultAPIServerPort),
					Protocol:        elbv2types.ProtocolEnumTcp,
					Tags: []elbv2types.Tag{
						{
							Key:   aws.String("test"),
							Value: aws.String("tag"),
						},
					},
				}).Return(&elbv2.CreateListenerOutput{
					Listeners: []elbv2types.Listener{
						{
							ListenerArn: aws.String("listener::arn"),
						},
					},
				}, nil)
				m.ModifyTargetGroupAttributes(gomock.Any(), &elbv2.ModifyTargetGroupAttributesInput{
					TargetGroupArn: aws.String(tgArn),
					Attributes: []elbv2types.TargetGroupAttribute{
						{
							Key:   aws.String(infrav1.TargetGroupAttributeEnableConnectionTermination),
							Value: aws.String("false"),
						},
						{
							Key:   aws.String(infrav1.TargetGroupAttributeUnhealthyDrainingIntervalSeconds),
							Value: aws.String("300"),
						},
						{
							Key:   aws.String(infrav1.TargetGroupAttributeEnablePreserveClientIP),
							Value: aws.String("false"),
						},
					},
				}).Return(nil, nil)
			},
			check: func(t *testing.T, tgs []*elbv2types.TargetGroup, _ []*elbv2types.Listener, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}
				got := tgs[0].HealthCheckProtocol
				want := elbv2types.ProtocolEnumHttps
				if got != want {
					t.Fatalf("Health Check protocol for the API Target group did not equal expected value: %s; was: '%s'", want, got)
				}
			},
		},
	}

	ctx := context.TODO()

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
								Port:     aws.String(infrav1.DefaultAPIServerPortString),
							},
						},
					},
				},
				LoadBalancerType: infrav1.LoadBalancerTypeNLB,
				SubnetIDs:        []string{clusterSubnetID},
			}

			spec := tc.spec(*loadBalancerSpec)
			tgs, listeners, err := s.reconcileTargetGroupsAndListeners(ctx, spec.ARN, &spec, clusterScope.ControlPlaneLoadBalancer())
			tc.check(t, tgs, listeners, err)
		})
	}
}

func TestReconcileV2LB(t *testing.T) {
	const (
		namespace       = "foo"
		clusterName     = "bar"
		clusterSubnetID = "subnet-1"
		elbName         = "bar-apiserver"
		elbArn          = "arn::apiserver"
		tgArn           = "arn::target-group"
		vpcID           = "vpc-id"
		az              = "us-west-1a"
	)

	tests := []struct {
		name          string
		elbV2APIMocks func(m *mocks.MockELBV2APIMockRecorder)
		check         func(t *testing.T, lb *infrav1.LoadBalancer, err error)
		awsCluster    func(acl infrav1.AWSCluster) infrav1.AWSCluster
		spec          func(spec infrav1.LoadBalancer) infrav1.LoadBalancer
	}{
		{
			name: "ensure status populated with BYO NLB",
			spec: func(spec infrav1.LoadBalancer) infrav1.LoadBalancer {
				return spec
			},
			awsCluster: func(acl infrav1.AWSCluster) infrav1.AWSCluster {
				acl.Spec.ControlPlaneLoadBalancer.Name = aws.String(elbName)
				return acl
			},
			elbV2APIMocks: func(m *mocks.MockELBV2APIMockRecorder) {
				m.DescribeLoadBalancers(gomock.Any(), &elbv2.DescribeLoadBalancersInput{
					Names: []string{elbName},
				}).
					Return(&elbv2.DescribeLoadBalancersOutput{
						LoadBalancers: []elbv2types.LoadBalancer{
							{
								LoadBalancerArn:  aws.String(elbArn),
								LoadBalancerName: aws.String(elbName),
								Scheme:           SchemeToSDKScheme(infrav1.ELBSchemeInternetFacing),
								AvailabilityZones: []elbv2types.AvailabilityZone{
									{
										SubnetId: aws.String(clusterSubnetID),
										ZoneName: aws.String(az),
									},
								},
								VpcId: aws.String(vpcID),
							},
						},
					}, nil)
				m.DescribeLoadBalancerAttributes(gomock.Any(), &elbv2.DescribeLoadBalancerAttributesInput{LoadBalancerArn: aws.String(elbArn)}).Return(
					&elbv2.DescribeLoadBalancerAttributesOutput{
						Attributes: []elbv2types.LoadBalancerAttribute{
							{
								Key:   aws.String("load_balancing.cross_zone.enabled"),
								Value: aws.String("false"),
							},
						},
					},
					nil,
				)
				m.DescribeTags(gomock.Any(), &elbv2.DescribeTagsInput{ResourceArns: []string{elbArn}}).Return(
					&elbv2.DescribeTagsOutput{
						TagDescriptions: []elbv2types.TagDescription{
							{
								ResourceArn: aws.String(elbArn),
								Tags:        []elbv2types.Tag{},
							},
						},
					},
					nil,
				)
				m.WaitUntilLoadBalancerAvailable(gomock.Any(), gomock.Eq(&elbv2.DescribeLoadBalancersInput{
					LoadBalancerArns: []string{elbArn},
				}), maxWaitActiveUpdateDelete).Return(nil)
			},
			check: func(t *testing.T, lb *infrav1.LoadBalancer, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}
				if len(lb.AvailabilityZones) != 1 {
					t.Errorf("Expected LB to contain 1 availability zone, got %v", len(lb.AvailabilityZones))
				}
			},
		},
		{
			name: "ensure NLB without SGs doesn't attempt to add new SGs",
			spec: func(spec infrav1.LoadBalancer) infrav1.LoadBalancer {
				return spec
			},
			awsCluster: func(acl infrav1.AWSCluster) infrav1.AWSCluster {
				acl.Spec.ControlPlaneLoadBalancer.Name = aws.String(elbName)
				acl.Spec.ControlPlaneLoadBalancer.LoadBalancerType = infrav1.LoadBalancerTypeNLB
				acl.Spec.ControlPlaneLoadBalancer.AdditionalSecurityGroups = []string{"sg-001"}
				return acl
			},
			elbV2APIMocks: func(m *mocks.MockELBV2APIMockRecorder) {
				m.DescribeLoadBalancers(gomock.Any(), gomock.Eq(&elbv2.DescribeLoadBalancersInput{
					Names: []string{elbName},
				})).
					Return(&elbv2.DescribeLoadBalancersOutput{
						LoadBalancers: []elbv2types.LoadBalancer{
							{
								LoadBalancerArn:  aws.String(elbArn),
								LoadBalancerName: aws.String(elbName),
								Scheme:           SchemeToSDKScheme(infrav1.ELBSchemeInternetFacing),
								AvailabilityZones: []elbv2types.AvailabilityZone{
									{
										SubnetId: aws.String(clusterSubnetID),
										ZoneName: aws.String(az),
									},
								},
								VpcId: aws.String(vpcID),
							},
						},
					}, nil)
				m.DescribeTargetGroups(gomock.Any(), gomock.Eq(&elbv2.DescribeTargetGroupsInput{
					LoadBalancerArn: aws.String(elbArn),
				})).
					Return(&elbv2.DescribeTargetGroupsOutput{
						NextMarker: new(string),
						TargetGroups: []elbv2types.TargetGroup{
							{
								HealthCheckEnabled: aws.Bool(true),
								LoadBalancerArns:   []string{elbArn},
								Matcher:            &elbv2types.Matcher{},
								TargetGroupArn:     aws.String(tgArn),
								TargetGroupName:    aws.String("targetGroup"),
							},
						},
					}, nil)
				m.ModifyLoadBalancerAttributes(gomock.Any(), &elbv2.ModifyLoadBalancerAttributesInput{
					LoadBalancerArn: aws.String(elbArn),
					Attributes: []elbv2types.LoadBalancerAttribute{
						{
							Key:   aws.String("load_balancing.cross_zone.enabled"),
							Value: aws.String("false"),
						},
					},
				}).
					Return(&elbv2.ModifyLoadBalancerAttributesOutput{}, nil)

				m.CreateTargetGroup(gomock.Any(), helpers.PartialMatchCreateTargetGroupInput(t, &elbv2.CreateTargetGroupInput{
					HealthCheckEnabled:         aws.Bool(true),
					HealthCheckIntervalSeconds: aws.Int32(infrav1.DefaultAPIServerHealthCheckIntervalSec),
					HealthCheckPort:            aws.String(infrav1.DefaultAPIServerPortString),
					HealthCheckProtocol:        elbv2types.ProtocolEnumTcp,
					HealthCheckTimeoutSeconds:  aws.Int32(infrav1.DefaultAPIServerHealthCheckTimeoutSec),
					HealthyThresholdCount:      aws.Int32(infrav1.DefaultAPIServerHealthThresholdCount),
					// Note: this is treated as a prefix with the partial matcher.
					Name:     aws.String("apiserver-target"),
					Port:     aws.Int32(infrav1.DefaultAPIServerPort),
					Protocol: elbv2types.ProtocolEnumTcp,
					Tags: []elbv2types.Tag{
						{
							Key:   aws.String("Name"),
							Value: aws.String("bar-apiserver"),
						},
						{
							Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/bar"),
							Value: aws.String("owned"),
						},
						{
							Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
							Value: aws.String("apiserver"),
						},
					},
					UnhealthyThresholdCount: aws.Int32(infrav1.DefaultAPIServerUnhealthThresholdCount),
					VpcId:                   aws.String(vpcID),
				})).Return(&elbv2.CreateTargetGroupOutput{
					TargetGroups: []elbv2types.TargetGroup{
						{
							TargetGroupArn:             aws.String(tgArn),
							VpcId:                      aws.String(vpcID),
							HealthyThresholdCount:      aws.Int32(infrav1.DefaultAPIServerHealthThresholdCount),
							UnhealthyThresholdCount:    aws.Int32(infrav1.DefaultAPIServerUnhealthThresholdCount),
							HealthCheckIntervalSeconds: aws.Int32(infrav1.DefaultAPIServerHealthCheckIntervalSec),
							HealthCheckTimeoutSeconds:  aws.Int32(infrav1.DefaultAPIServerHealthCheckTimeoutSec),
						},
					},
				}, nil)

				m.ModifyTargetGroupAttributes(gomock.Any(), gomock.Eq(&elbv2.ModifyTargetGroupAttributesInput{
					TargetGroupArn: aws.String(tgArn),
					Attributes: []elbv2types.TargetGroupAttribute{
						{
							Key:   aws.String(infrav1.TargetGroupAttributeEnableConnectionTermination),
							Value: aws.String("false"),
						},
						{
							Key:   aws.String(infrav1.TargetGroupAttributeUnhealthyDrainingIntervalSeconds),
							Value: aws.String("300"),
						},
						{
							Key:   aws.String(infrav1.TargetGroupAttributeEnablePreserveClientIP),
							Value: aws.String("false"),
						},
					},
				})).Return(nil, nil)

				m.DescribeListeners(gomock.Any(), gomock.Eq(&elbv2.DescribeListenersInput{
					LoadBalancerArn: aws.String(elbArn),
				})).
					Return(&elbv2.DescribeListenersOutput{
						Listeners: []elbv2types.Listener{{
							DefaultActions: []elbv2types.Action{{
								TargetGroupArn: aws.String("arn::targetgroup"),
							}},
							ListenerArn:     aws.String("arn::listener"),
							LoadBalancerArn: aws.String(elbArn),
						}},
					}, nil)
				m.CreateListener(gomock.Any(), gomock.Eq(&elbv2.CreateListenerInput{
					DefaultActions: []elbv2types.Action{
						{
							TargetGroupArn: aws.String(tgArn),
							Type:           elbv2types.ActionTypeEnumForward,
						},
					},
					LoadBalancerArn: aws.String(elbArn),
					Port:            aws.Int32(infrav1.DefaultAPIServerPort),
					Protocol:        elbv2types.ProtocolEnumTcp,
					Tags: []elbv2types.Tag{
						{
							Key:   aws.String("Name"),
							Value: aws.String("bar-apiserver"),
						},
						{
							Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/bar"),
							Value: aws.String("owned"),
						},
						{
							Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
							Value: aws.String("apiserver"),
						},
					},
				})).Return(&elbv2.CreateListenerOutput{
					Listeners: []elbv2types.Listener{
						{
							DefaultActions: []elbv2types.Action{
								{
									TargetGroupArn: aws.String(tgArn),
									Type:           elbv2types.ActionTypeEnumForward,
								},
							},
							ListenerArn: aws.String("listener::arn"),
							Port:        aws.Int32(infrav1.DefaultAPIServerPort),
							Protocol:    elbv2types.ProtocolEnumTcp,
						},
					},
				}, nil)
				m.DescribeLoadBalancerAttributes(gomock.Any(), &elbv2.DescribeLoadBalancerAttributesInput{LoadBalancerArn: aws.String(elbArn)}).Return(
					&elbv2.DescribeLoadBalancerAttributesOutput{
						Attributes: []elbv2types.LoadBalancerAttribute{
							{
								Key:   aws.String("load_balancing.cross_zone.enabled"),
								Value: aws.String("false"),
							},
							{
								Key:   aws.String(infrav1.ClusterTagKey(clusterName)),
								Value: aws.String(string(infrav1.ResourceLifecycleOwned)),
							},
						},
					}, nil)
				m.DescribeTags(gomock.Any(), &elbv2.DescribeTagsInput{ResourceArns: []string{elbArn}}).Return(
					&elbv2.DescribeTagsOutput{
						TagDescriptions: []elbv2types.TagDescription{
							{
								ResourceArn: aws.String(elbArn),
								Tags: []elbv2types.Tag{
									{
										Key:   aws.String(infrav1.ClusterTagKey(clusterName)),
										Value: aws.String(string(infrav1.ResourceLifecycleOwned)),
									},
								},
							},
						},
					}, nil)

				// Avoid the need to sort the AddTagsInput.Tags slice
				m.AddTags(gomock.Any(), gomock.AssignableToTypeOf(&elbv2.AddTagsInput{})).Return(&elbv2.AddTagsOutput{}, nil)

				m.SetSubnets(gomock.Any(), gomock.Eq(&elbv2.SetSubnetsInput{
					LoadBalancerArn: aws.String(elbArn),
				})).Return(&elbv2.SetSubnetsOutput{}, nil)

				m.WaitUntilLoadBalancerAvailable(gomock.Any(), gomock.Eq(&elbv2.DescribeLoadBalancersInput{
					LoadBalancerArns: []string{elbArn},
				}), maxWaitActiveUpdateDelete).Return(nil)
			},
			check: func(t *testing.T, lb *infrav1.LoadBalancer, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}
				if len(lb.SecurityGroupIDs) != 0 {
					t.Errorf("Expected LB to contain 0 security groups, got %v", len(lb.SecurityGroupIDs))
				}
			},
		},
	}

	ctx := context.TODO()

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
				AWSCluster:                &cluster,
				MaxWaitActiveUpdateDelete: maxWaitActiveUpdateDelete,
			})
			if err != nil {
				t.Fatal(err)
			}

			tc.elbV2APIMocks(elbV2APIMocks.EXPECT())

			s := &Service{
				scope:       clusterScope,
				ELBV2Client: elbV2APIMocks,
			}
			reconciler, err := s.getOrCreateV2LB(ctx, clusterScope.ControlPlaneLoadBalancer())
			if err == nil {
				err = reconciler()
			}
			lb := s.scope.Network().APIServerELB

			tc.check(t, &lb, err)
		})
	}
}

func TestReconcileLoadbalancers(t *testing.T) {
	const (
		namespace       = "foo"
		clusterName     = "bar"
		clusterSubnetID = "subnet-1"
		elbName         = "bar-apiserver"
		elbArn          = "arn::apiserver"
		secondElbName   = "bar-apiserver2"
		secondElbArn    = "arn::apiserver2"
		vpcID           = "vpc-id"
		az              = "us-west-1a"
	)

	primaryELB := func() elbv2types.LoadBalancer {
		return elbv2types.LoadBalancer{
			LoadBalancerArn:  aws.String(elbArn),
			LoadBalancerName: aws.String(elbName),
			DNSName:          aws.String(elbName),
			Scheme:           SchemeToSDKScheme(infrav1.ELBSchemeInternetFacing),
			AvailabilityZones: []elbv2types.AvailabilityZone{
				{
					SubnetId: aws.String(clusterSubnetID),
					ZoneName: aws.String(az),
				},
			},
			VpcId: aws.String(vpcID),
		}
	}

	secondaryELB := func() elbv2types.LoadBalancer {
		return elbv2types.LoadBalancer{
			LoadBalancerArn:  aws.String(secondElbArn),
			LoadBalancerName: aws.String(secondElbName),
			DNSName:          aws.String(secondElbName),
			Scheme:           SchemeToSDKScheme(infrav1.ELBSchemeInternal),
			AvailabilityZones: []elbv2types.AvailabilityZone{
				{
					SubnetId: aws.String(clusterSubnetID),
					ZoneName: aws.String(az),
				},
			},
			VpcId: aws.String(vpcID),
		}
	}

	tests := []struct {
		name          string
		elbV2APIMocks func(*mocks.MockELBV2APIMockRecorder)
		check         func(t *testing.T, firstLB, secondLB *infrav1.LoadBalancer, err error)
		awsCluster    func(acl infrav1.AWSCluster) infrav1.AWSCluster
		spec          func(spec infrav1.LoadBalancer) infrav1.LoadBalancer
	}{
		{
			name: "ensure two load balancers are reconciled",
			awsCluster: func(acl infrav1.AWSCluster) infrav1.AWSCluster {
				acl.Spec.ControlPlaneLoadBalancer.Name = aws.String(elbName)
				acl.Spec.SecondaryControlPlaneLoadBalancer = &infrav1.AWSLoadBalancerSpec{
					Name:             aws.String(secondElbName),
					Scheme:           &infrav1.ELBSchemeInternal,
					LoadBalancerType: infrav1.LoadBalancerTypeNLB,
				}
				return acl
			},
			elbV2APIMocks: func(m *mocks.MockELBV2APIMockRecorder) {
				m.DescribeLoadBalancers(gomock.Any(), gomock.Eq(&elbv2.DescribeLoadBalancersInput{
					Names: []string{elbName},
				})).
					Return(&elbv2.DescribeLoadBalancersOutput{
						LoadBalancers: []elbv2types.LoadBalancer{primaryELB()},
					}, nil)
				m.DescribeLoadBalancerAttributes(gomock.Any(), &elbv2.DescribeLoadBalancerAttributesInput{LoadBalancerArn: aws.String(elbArn)}).Return(
					&elbv2.DescribeLoadBalancerAttributesOutput{
						Attributes: []elbv2types.LoadBalancerAttribute{
							{
								Key:   aws.String("load_balancing.cross_zone.enabled"),
								Value: aws.String("false"),
							},
						},
					},
					nil,
				)
				m.DescribeTags(gomock.Any(), &elbv2.DescribeTagsInput{ResourceArns: []string{elbArn}}).Return(
					&elbv2.DescribeTagsOutput{
						TagDescriptions: []elbv2types.TagDescription{
							{
								ResourceArn: aws.String(elbArn),
								Tags:        []elbv2types.Tag{},
							},
						},
					},
					nil,
				)

				m.DescribeLoadBalancers(gomock.Any(), gomock.Eq(&elbv2.DescribeLoadBalancersInput{
					Names: []string{secondElbName},
				})).
					Return(&elbv2.DescribeLoadBalancersOutput{
						LoadBalancers: []elbv2types.LoadBalancer{secondaryELB()},
					}, nil)
				m.DescribeLoadBalancerAttributes(gomock.Any(), &elbv2.DescribeLoadBalancerAttributesInput{LoadBalancerArn: aws.String(secondElbArn)}).Return(
					&elbv2.DescribeLoadBalancerAttributesOutput{
						Attributes: []elbv2types.LoadBalancerAttribute{
							{
								Key:   aws.String("load_balancing.cross_zone.enabled"),
								Value: aws.String("false"),
							},
						},
					},
					nil,
				)
				m.DescribeTags(gomock.Any(), &elbv2.DescribeTagsInput{ResourceArns: []string{secondElbArn}}).Return(
					&elbv2.DescribeTagsOutput{
						TagDescriptions: []elbv2types.TagDescription{
							{
								ResourceArn: aws.String(secondElbArn),
								Tags:        []elbv2types.Tag{},
							},
						},
					},
					nil,
				)
				m.WaitUntilLoadBalancerAvailable(gomock.Any(), &elbv2.DescribeLoadBalancersInput{
					LoadBalancerArns: []string{elbArn},
				}, maxWaitActiveUpdateDelete).Return(nil)
				m.WaitUntilLoadBalancerAvailable(gomock.Any(), &elbv2.DescribeLoadBalancersInput{
					LoadBalancerArns: []string{secondElbArn},
				}, maxWaitActiveUpdateDelete).Return(nil)
			},
			check: func(t *testing.T, firstLB *infrav1.LoadBalancer, secondLB *infrav1.LoadBalancer, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}

				if len(firstLB.AvailabilityZones) != 1 {
					t.Errorf("Expected first LB to contain 1 availability zone, got %v", len(firstLB.AvailabilityZones))
				}
				if secondLB == nil {
					t.Errorf("Expected second LB to be populated, was nil")
				}
			},
		},
		{
			name: "ensure two load balancers are created concurrently",
			awsCluster: func(acl infrav1.AWSCluster) infrav1.AWSCluster {
				acl.Spec.ControlPlaneLoadBalancer.Name = aws.String(elbName)
				acl.Spec.SecondaryControlPlaneLoadBalancer = &infrav1.AWSLoadBalancerSpec{
					Name:             aws.String(secondElbName),
					Scheme:           &infrav1.ELBSchemeInternal,
					LoadBalancerType: infrav1.LoadBalancerTypeNLB,
				}
				return acl
			},
			elbV2APIMocks: func(m *mocks.MockELBV2APIMockRecorder) {
				// Initial DescribeLoadBalancers return empty
				m.DescribeLoadBalancers(gomock.Any(), gomock.Eq(&elbv2.DescribeLoadBalancersInput{
					Names: []string{elbName},
				})).Return(&elbv2.DescribeLoadBalancersOutput{}, nil)
				m.DescribeLoadBalancers(gomock.Any(), gomock.Eq(&elbv2.DescribeLoadBalancersInput{
					Names: []string{secondElbName},
				})).Return(&elbv2.DescribeLoadBalancersOutput{}, nil)

				// Create both load balancers
				createLB1 := m.CreateLoadBalancer(gomock.Any(), gomock.Eq(&elbv2.CreateLoadBalancerInput{
					Name:           aws.String(elbName),
					Scheme:         SchemeToSDKScheme(infrav1.ELBSchemeInternetFacing),
					SecurityGroups: []string{""},
					Tags: []elbv2types.Tag{
						{
							Key:   aws.String("Name"),
							Value: aws.String(elbName),
						},
						{
							Key:   aws.String(infrav1.ClusterTagKey(clusterName)),
							Value: aws.String(string(infrav1.ResourceLifecycleOwned)),
						},
						{
							Key:   aws.String(infrav1.NameAWSClusterAPIRole),
							Value: aws.String(infrav1.APIServerRoleTagValue),
						},
					},
					Type: elbv2types.LoadBalancerTypeEnumNetwork,
				})).Return(&elbv2.CreateLoadBalancerOutput{
					LoadBalancers: []elbv2types.LoadBalancer{primaryELB()},
				}, nil)

				createLB2 := m.CreateLoadBalancer(gomock.Any(), gomock.Eq(&elbv2.CreateLoadBalancerInput{
					Name:           aws.String(secondElbName),
					Scheme:         SchemeToSDKScheme(infrav1.ELBSchemeInternal),
					SecurityGroups: []string{""},
					Tags: []elbv2types.Tag{
						{
							Key:   aws.String("Name"),
							Value: aws.String(secondElbName),
						},
						{
							Key:   aws.String(infrav1.ClusterTagKey(clusterName)),
							Value: aws.String(string(infrav1.ResourceLifecycleOwned)),
						},
						{
							Key:   aws.String(infrav1.NameAWSClusterAPIRole),
							Value: aws.String(infrav1.APIServerRoleTagValue),
						},
					},
					Type: elbv2types.LoadBalancerTypeEnumNetwork,
				})).Return(&elbv2.CreateLoadBalancerOutput{
					LoadBalancers: []elbv2types.LoadBalancer{secondaryELB()},
				}, nil)

				// Assert that we don't wait for either load balancer to be
				// available until both load balancers have been created
				m.WaitUntilLoadBalancerAvailable(gomock.Any(), gomock.Eq(&elbv2.DescribeLoadBalancersInput{
					LoadBalancerArns: []string{elbArn},
				}), maxWaitActiveUpdateDelete).Return(nil).After(createLB1).After(createLB2)
				m.WaitUntilLoadBalancerAvailable(gomock.Any(), gomock.Eq(&elbv2.DescribeLoadBalancersInput{
					LoadBalancerArns: []string{secondElbArn},
				}), maxWaitActiveUpdateDelete).Return(nil).After(createLB1).After(createLB2)

				// Make minimal assertions on other calls not under test
				m.DescribeTargetGroups(gomock.Any(), gomock.Any()).Return(&elbv2.DescribeTargetGroupsOutput{}, nil).AnyTimes()
				m.DescribeListeners(gomock.Any(), gomock.Any()).Return(&elbv2.DescribeListenersOutput{}, nil).AnyTimes()
				m.ModifyTargetGroupAttributes(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
				m.ModifyListener(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
				m.SetSecurityGroups(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()

				// These calls return the same info for both load balancers, but it's not important to this test
				m.CreateTargetGroup(gomock.Any(), gomock.Any()).Return(&elbv2.CreateTargetGroupOutput{
					TargetGroups: []elbv2types.TargetGroup{
						{
							TargetGroupName: aws.String(elbName),
							TargetGroupArn:  aws.String(elbArn),
						},
					},
				}, nil).AnyTimes()
				m.CreateListener(gomock.Any(), gomock.Any()).Return(&elbv2.CreateListenerOutput{
					Listeners: []elbv2types.Listener{
						{
							ListenerArn: aws.String(elbArn),
						},
					},
				}, nil).AnyTimes()
			},
			check: func(t *testing.T, firstLB, secondLB *infrav1.LoadBalancer, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}
				if firstLB == nil {
					t.Errorf("Expected first LB to be populated, was nil")
				}
				if secondLB == nil {
					t.Errorf("Expected second LB to be populated, was nil")
				}
			},
		},
	}

	ctx := context.TODO()

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
				AWSCluster:                &cluster,
				MaxWaitActiveUpdateDelete: maxWaitActiveUpdateDelete,
			})
			if err != nil {
				t.Fatal(err)
			}

			tc.elbV2APIMocks(elbV2APIMocks.EXPECT())

			s := &Service{
				scope:       clusterScope,
				ELBV2Client: elbV2APIMocks,
			}
			err = s.ReconcileLoadbalancers(ctx)
			firstLB := s.scope.Network().APIServerELB
			secondLB := s.scope.Network().SecondaryAPIServerELB
			tc.check(t, &firstLB, &secondLB, err)
		})
	}
}

func TestDeleteAPIServerELB(t *testing.T) {
	clusterName := "bar"
	elbName := "bar-apiserver"
	tests := []struct {
		name             string
		elbAPIMocks      func(m *mocks.MockELBAPIMockRecorder)
		verifyAWSCluster func(*infrav1.AWSCluster)
	}{
		{
			name: "if control plane ELB is not found, do nothing",
			elbAPIMocks: func(m *mocks.MockELBAPIMockRecorder) {
				m.DescribeLoadBalancers(gomock.Any(), &elb.DescribeLoadBalancersInput{
					LoadBalancerNames: []string{elbName},
				}).Return(nil, &elbtypes.AccessPointNotFoundException{})
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
				m.DescribeLoadBalancers(gomock.Any(), &elb.DescribeLoadBalancersInput{LoadBalancerNames: []string{elbName}}).Return(
					&elb.DescribeLoadBalancersOutput{
						LoadBalancerDescriptions: []elbtypes.LoadBalancerDescription{
							{
								LoadBalancerName: aws.String(elbName),
								Scheme:           aws.String(string(infrav1.ELBSchemeInternetFacing)),
							},
						},
					},
					nil,
				)

				m.DescribeLoadBalancerAttributes(gomock.Any(), &elb.DescribeLoadBalancerAttributesInput{LoadBalancerName: aws.String(elbName)}).Return(
					&elb.DescribeLoadBalancerAttributesOutput{
						LoadBalancerAttributes: &elbtypes.LoadBalancerAttributes{
							CrossZoneLoadBalancing: &elbtypes.CrossZoneLoadBalancing{
								Enabled: false,
							},
						},
					},
					nil,
				)

				m.DescribeTags(gomock.Any(), &elb.DescribeTagsInput{LoadBalancerNames: []string{elbName}}).Return(
					&elb.DescribeTagsOutput{
						TagDescriptions: []elbtypes.TagDescription{
							{
								LoadBalancerName: aws.String(elbName),
								Tags:             []elbtypes.Tag{},
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
				m.DescribeLoadBalancers(gomock.Any(), &elb.DescribeLoadBalancersInput{LoadBalancerNames: []string{elbName}}).Return(
					&elb.DescribeLoadBalancersOutput{
						LoadBalancerDescriptions: []elbtypes.LoadBalancerDescription{
							{
								LoadBalancerName: aws.String(elbName),
								Scheme:           aws.String(string(infrav1.ELBSchemeInternetFacing)),
							},
						},
					},
					nil,
				)

				m.DescribeLoadBalancerAttributes(gomock.Any(), &elb.DescribeLoadBalancerAttributesInput{LoadBalancerName: aws.String(elbName)}).Return(
					&elb.DescribeLoadBalancerAttributesOutput{
						LoadBalancerAttributes: &elbtypes.LoadBalancerAttributes{
							CrossZoneLoadBalancing: &elbtypes.CrossZoneLoadBalancing{
								Enabled: false,
							},
						},
					},
					nil,
				)

				m.DescribeTags(gomock.Any(), &elb.DescribeTagsInput{LoadBalancerNames: []string{elbName}}).Return(
					&elb.DescribeTagsOutput{
						TagDescriptions: []elbtypes.TagDescription{
							{
								LoadBalancerName: aws.String(elbName),
								Tags: []elbtypes.Tag{{
									Key:   aws.String(infrav1.ClusterTagKey(clusterName)),
									Value: aws.String(string(infrav1.ResourceLifecycleOwned)),
								}},
							},
						},
					},
					nil,
				)

				m.DeleteLoadBalancer(gomock.Any(), &elb.DeleteLoadBalancerInput{LoadBalancerName: aws.String(elbName)}).Return(
					&elb.DeleteLoadBalancerOutput{}, nil)

				m.DescribeLoadBalancers(gomock.Any(), &elb.DescribeLoadBalancersInput{LoadBalancerNames: []string{elbName}}).Return(
					&elb.DescribeLoadBalancersOutput{
						LoadBalancerDescriptions: []elbtypes.LoadBalancerDescription{},
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

	ctx := context.TODO()

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

			client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(awsCluster).WithStatusSubresource(awsCluster).Build()

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

			err = s.deleteAPIServerELB(ctx)
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
	tgArn := "arn::target-group"
	tests := []struct {
		name         string
		elbv2ApiMock func(m *mocks.MockELBV2APIMockRecorder)
	}{
		{
			name: "if control plane NLB is not found, do nothing",
			elbv2ApiMock: func(m *mocks.MockELBV2APIMockRecorder) {
				m.DescribeLoadBalancers(gomock.Any(), &elbv2.DescribeLoadBalancersInput{
					Names: []string{elbName},
				}).Return(nil, &elbtypes.AccessPointNotFoundException{})
			},
		},
		{
			name: "if control plane NLB is found, and it is not managed, do nothing",
			elbv2ApiMock: func(m *mocks.MockELBV2APIMockRecorder) {
				m.DescribeLoadBalancers(gomock.Any(), &elbv2.DescribeLoadBalancersInput{Names: []string{elbName}}).Return(
					&elbv2.DescribeLoadBalancersOutput{
						LoadBalancers: []elbv2types.LoadBalancer{
							{
								LoadBalancerArn:  aws.String(elbArn),
								LoadBalancerName: aws.String(elbName),
								Scheme:           SchemeToSDKScheme(infrav1.ELBSchemeInternetFacing),
							},
						},
					},
					nil,
				)

				m.DescribeLoadBalancerAttributes(gomock.Any(), &elbv2.DescribeLoadBalancerAttributesInput{LoadBalancerArn: aws.String(elbArn)}).Return(
					&elbv2.DescribeLoadBalancerAttributesOutput{
						Attributes: []elbv2types.LoadBalancerAttribute{
							{
								Key:   aws.String("load_balancing.cross_zone.enabled"),
								Value: aws.String("false"),
							},
						},
					},
					nil,
				)

				m.DescribeTags(gomock.Any(), &elbv2.DescribeTagsInput{ResourceArns: []string{elbArn}}).Return(
					&elbv2.DescribeTagsOutput{
						TagDescriptions: []elbv2types.TagDescription{
							{
								ResourceArn: aws.String(elbArn),
								Tags:        []elbv2types.Tag{},
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
				m.DescribeLoadBalancers(gomock.Any(), &elbv2.DescribeLoadBalancersInput{Names: []string{elbName}}).Return(
					&elbv2.DescribeLoadBalancersOutput{
						LoadBalancers: []elbv2types.LoadBalancer{
							{
								LoadBalancerArn:  aws.String(elbArn),
								LoadBalancerName: aws.String(elbName),
								Scheme:           SchemeToSDKScheme(infrav1.ELBSchemeInternetFacing),
							},
						},
					},
					nil,
				)

				m.DescribeLoadBalancerAttributes(gomock.Any(), &elbv2.DescribeLoadBalancerAttributesInput{LoadBalancerArn: aws.String(elbArn)}).Return(
					&elbv2.DescribeLoadBalancerAttributesOutput{
						Attributes: []elbv2types.LoadBalancerAttribute{
							{
								Key:   aws.String("load_balancing.cross_zone.enabled"),
								Value: aws.String("false"),
							},
						},
					},
					nil,
				)

				m.DescribeTags(gomock.Any(), &elbv2.DescribeTagsInput{ResourceArns: []string{elbArn}}).Return(
					&elbv2.DescribeTagsOutput{
						TagDescriptions: []elbv2types.TagDescription{
							{
								ResourceArn: aws.String(elbArn),
								Tags: []elbv2types.Tag{{
									Key:   aws.String(infrav1.ClusterTagKey(clusterName)),
									Value: aws.String(string(infrav1.ResourceLifecycleOwned)),
								}},
							},
						},
					},
					nil,
				)

				// delete listeners
				m.DescribeListeners(gomock.Any(), &elbv2.DescribeListenersInput{LoadBalancerArn: aws.String(elbArn)}).Return(&elbv2.DescribeListenersOutput{
					Listeners: []elbv2types.Listener{
						{
							ListenerArn: aws.String("listener::arn"),
						},
					},
				}, nil)
				m.DeleteListener(gomock.Any(), &elbv2.DeleteListenerInput{ListenerArn: aws.String("listener::arn")}).Return(&elbv2.DeleteListenerOutput{}, nil)
				// delete target groups
				m.DescribeTargetGroups(gomock.Any(), &elbv2.DescribeTargetGroupsInput{LoadBalancerArn: aws.String(elbArn)}).Return(&elbv2.DescribeTargetGroupsOutput{
					TargetGroups: []elbv2types.TargetGroup{
						{
							TargetGroupArn: aws.String(tgArn),
						},
					},
				}, nil)
				m.DeleteTargetGroup(gomock.Any(), &elbv2.DeleteTargetGroupInput{TargetGroupArn: aws.String(tgArn)}).Return(&elbv2.DeleteTargetGroupOutput{}, nil)
				// delete the load balancer

				m.DeleteLoadBalancer(gomock.Any(), &elbv2.DeleteLoadBalancerInput{LoadBalancerArn: aws.String(elbArn)}).Return(
					&elbv2.DeleteLoadBalancerOutput{}, nil)

				m.DescribeLoadBalancers(gomock.Any(), &elbv2.DescribeLoadBalancersInput{Names: []string{elbName}}).Return(
					&elbv2.DescribeLoadBalancersOutput{
						LoadBalancers: []elbv2types.LoadBalancer{},
					},
					nil,
				)
			},
		},
	}

	ctx := context.TODO()

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

			client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(awsCluster).WithStatusSubresource(awsCluster).Build()

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

			err = s.deleteExistingNLBs(ctx)
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
				m.GetResourcesPages(gomock.Any(), &rgapi.GetResourcesInput{
					ResourceTypeFilters: []string{elbResourceType},
					TagFilters: []rgapitypes.TagFilter{
						{
							Key:    aws.String(infrav1.ClusterAWSCloudProviderTagKey(clusterName)),
							Values: []string{string(infrav1.ResourceLifecycleOwned)},
						},
					},
				}, gomock.Any()).Do(func(_, _, y interface{}) {
					funct := y.(func(output *rgapi.GetResourcesOutput))
					funct(&rgapi.GetResourcesOutput{
						ResourceTagMappingList: []rgapitypes.ResourceTagMapping{
							{
								ResourceARN: aws.String("arn:aws:elasticloadbalancing:eu-west-2:1234567890:loadbalancer/lb-service-name"),
								Tags: []rgapitypes.Tag{{
									Key:   aws.String(infrav1.ClusterAWSCloudProviderTagKey(clusterName)),
									Value: aws.String(string(infrav1.ResourceLifecycleOwned)),
								}},
							},
						},
					})
				}).Return(nil)
			},
			elbAPIMocks: func(m *mocks.MockELBAPIMockRecorder) {
				m.DeleteLoadBalancer(gomock.Any(), &elb.DeleteLoadBalancerInput{LoadBalancerName: aws.String("lb-service-name")}).Return(nil, nil)
			},
			postDeleteRGAPIMocks: func(m *mocks.MockResourceGroupsTaggingAPIAPIMockRecorder) {
				m.GetResourcesPages(gomock.Any(), &rgapi.GetResourcesInput{
					ResourceTypeFilters: []string{elbResourceType},
					TagFilters: []rgapitypes.TagFilter{
						{
							Key:    aws.String(infrav1.ClusterAWSCloudProviderTagKey(clusterName)),
							Values: []string{string(infrav1.ResourceLifecycleOwned)},
						},
					},
				}, gomock.Any()).Do(func(_, _, y interface{}) {
					funct := y.(func(output *rgapi.GetResourcesOutput))
					funct(&rgapi.GetResourcesOutput{
						ResourceTagMappingList: []rgapitypes.ResourceTagMapping{},
					})
				}).Return(nil)
			},
		},
		{
			name: "fall back to ELB API when Resource Groups Tagging API fails and then delete successfully",
			rgAPIMocks: func(m *mocks.MockResourceGroupsTaggingAPIAPIMockRecorder) {
				m.GetResourcesPages(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.Errorf("connection failure")).AnyTimes()
			},
			elbAPIMocks: func(m *mocks.MockELBAPIMockRecorder) {
				m.DescribeLoadBalancersPages(gomock.Any(), gomock.Any(), gomock.Any()).Do(func(_, _, y interface{}) {
					funct := y.(func(output *elb.DescribeLoadBalancersOutput))
					funct(&elb.DescribeLoadBalancersOutput{
						LoadBalancerDescriptions: []elbtypes.LoadBalancerDescription{
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
					})
				}).Return(nil)
				m.DescribeTags(gomock.Any(), &elb.DescribeTagsInput{LoadBalancerNames: []string{"lb-service-name", "another-service-not-owned", "service-without-tags"}}).Return(&elb.DescribeTagsOutput{
					TagDescriptions: []elbtypes.TagDescription{
						{
							LoadBalancerName: aws.String("lb-service-name"),
							Tags: []elbtypes.Tag{{
								Key:   aws.String(infrav1.ClusterAWSCloudProviderTagKey(clusterName)),
								Value: aws.String(string(infrav1.ResourceLifecycleOwned)),
							}},
						},
						{
							LoadBalancerName: aws.String("another-service-not-owned"),
							Tags: []elbtypes.Tag{{
								Key:   aws.String("some-tag-key"),
								Value: aws.String("some-tag-value"),
							}},
						},
						{
							LoadBalancerName: aws.String("service-without-tags"),
							Tags:             []elbtypes.Tag{},
						},
					},
				}, nil)
				m.DeleteLoadBalancer(gomock.Any(), &elb.DeleteLoadBalancerInput{LoadBalancerName: aws.String("lb-service-name")}).Return(nil, nil)
			},
			postDeleteElbAPIMocks: func(m *mocks.MockELBAPIMockRecorder) {
				m.DescribeLoadBalancersPages(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
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

			err = s.deleteAWSCloudProviderELBs(ctx)
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
				m.GetResourcesPages(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			},
			DescribeElbAPIMocks: func(m *mocks.MockELBAPIMockRecorder) {
				m.DescribeLoadBalancers(gomock.Any(), &elb.DescribeLoadBalancersInput{
					LoadBalancerNames: []string{"bar-apiserver"},
				}).Return(&elb.DescribeLoadBalancersOutput{LoadBalancerDescriptions: []elbtypes.LoadBalancerDescription{{Scheme: ptr.To[string](string(infrav1.ELBSchemeInternal))}}}, nil)
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

			_, err = s.describeClassicELB(ctx, tc.lbName)
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
				m.GetResourcesPages(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			},
			DescribeElbV2APIMocks: func(m *mocks.MockELBV2APIMockRecorder) {
				m.DescribeLoadBalancers(gomock.Any(), &elbv2.DescribeLoadBalancersInput{
					Names: []string{"bar-apiserver"},
				}).Return(&elbv2.DescribeLoadBalancersOutput{LoadBalancers: []elbv2types.LoadBalancer{{Scheme: SchemeToSDKScheme(infrav1.ELBSchemeInternal)}}}, nil)
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

			_, err = s.describeLB(ctx, tc.lbName, clusterScope.ControlPlaneLoadBalancer())
			if err == nil {
				t.Fatal(err)
			}
		})
	}
}

func TestChunkELBs(t *testing.T) {
	base := "loadbalancer"
	names := make([]string, 0, 25)
	for i := range 25 {
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
	testHTTP := infrav1.ELBProtocol("HTTP")
	testHTTPS := infrav1.ELBProtocol("HTTPS")
	testTCP := infrav1.ELBProtocol("TCP")

	tests := []struct {
		testName                  string
		lbSpec                    *infrav1.AWSLoadBalancerSpec
		expectedHealthCheckTarget string
	}{
		{
			"default case",
			&infrav1.AWSLoadBalancerSpec{},
			"TCP:6443",
		},
		{
			"protocol http",
			&infrav1.AWSLoadBalancerSpec{
				HealthCheckProtocol: &testHTTP,
			},
			"HTTP:6443/readyz",
		},
		{
			"protocol https",
			&infrav1.AWSLoadBalancerSpec{
				HealthCheckProtocol: &testHTTPS,
			},
			"HTTPS:6443/readyz",
		},
		{
			"protocol tcp",
			&infrav1.AWSLoadBalancerSpec{
				HealthCheckProtocol: &testTCP,
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

func stubGetBaseService(t *testing.T, clusterName string) *Service {
	t.Helper()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	rgapiMock := mocks.NewMockResourceGroupsTaggingAPIAPI(mockCtrl)
	elbV2ApiMock := mocks.NewMockELBV2API(mockCtrl)

	scheme, err := setupScheme()
	if err != nil {
		t.Fatal(err)
	}
	awsCluster := &infrav1.AWSCluster{
		ObjectMeta: metav1.ObjectMeta{Name: clusterName},
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

	return &Service{
		scope:                 clusterScope,
		ResourceTaggingClient: rgapiMock,
		ELBV2Client:           elbV2ApiMock,
	}
}

func TestService_getAPITargetGroupHealthCheck(t *testing.T) {
	tests := []struct {
		name   string
		lbSpec *infrav1.AWSLoadBalancerSpec
		want   *infrav1.TargetGroupHealthCheck
	}{
		{
			name:   "default config",
			lbSpec: nil,
			want: &infrav1.TargetGroupHealthCheck{
				Protocol:                aws.String("TCP"),
				Port:                    aws.String("6443"),
				Path:                    nil,
				IntervalSeconds:         aws.Int64(infrav1.DefaultAPIServerHealthCheckIntervalSec),
				TimeoutSeconds:          aws.Int64(infrav1.DefaultAPIServerHealthCheckTimeoutSec),
				ThresholdCount:          aws.Int64(infrav1.DefaultAPIServerHealthThresholdCount),
				UnhealthyThresholdCount: aws.Int64(infrav1.DefaultAPIServerUnhealthThresholdCount),
			},
		},
		{
			name:   "default attributes, API health check TCP",
			lbSpec: &infrav1.AWSLoadBalancerSpec{},
			want: &infrav1.TargetGroupHealthCheck{
				Protocol:                aws.String("TCP"),
				Port:                    aws.String("6443"),
				Path:                    nil,
				IntervalSeconds:         aws.Int64(infrav1.DefaultAPIServerHealthCheckIntervalSec),
				TimeoutSeconds:          aws.Int64(infrav1.DefaultAPIServerHealthCheckTimeoutSec),
				ThresholdCount:          aws.Int64(infrav1.DefaultAPIServerHealthThresholdCount),
				UnhealthyThresholdCount: aws.Int64(infrav1.DefaultAPIServerUnhealthThresholdCount),
			},
		},
		{
			name: "default attributes, API health check HTTP",
			lbSpec: &infrav1.AWSLoadBalancerSpec{
				HealthCheckProtocol: &infrav1.ELBProtocolHTTP,
			},
			want: &infrav1.TargetGroupHealthCheck{
				Protocol:                aws.String("HTTP"),
				Port:                    aws.String("6443"),
				Path:                    aws.String("/readyz"),
				IntervalSeconds:         aws.Int64(infrav1.DefaultAPIServerHealthCheckIntervalSec),
				TimeoutSeconds:          aws.Int64(infrav1.DefaultAPIServerHealthCheckTimeoutSec),
				ThresholdCount:          aws.Int64(infrav1.DefaultAPIServerHealthThresholdCount),
				UnhealthyThresholdCount: aws.Int64(infrav1.DefaultAPIServerUnhealthThresholdCount),
			},
		},
		{
			name: "default attributes, API health check HTTPS",
			lbSpec: &infrav1.AWSLoadBalancerSpec{
				HealthCheckProtocol: &infrav1.ELBProtocolHTTPS,
			},
			want: &infrav1.TargetGroupHealthCheck{
				Protocol:                aws.String("HTTPS"),
				Port:                    aws.String("6443"),
				Path:                    aws.String("/readyz"),
				IntervalSeconds:         aws.Int64(infrav1.DefaultAPIServerHealthCheckIntervalSec),
				TimeoutSeconds:          aws.Int64(infrav1.DefaultAPIServerHealthCheckTimeoutSec),
				ThresholdCount:          aws.Int64(infrav1.DefaultAPIServerHealthThresholdCount),
				UnhealthyThresholdCount: aws.Int64(infrav1.DefaultAPIServerUnhealthThresholdCount),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := stubGetBaseService(t, "foo")
			if got := s.getAPITargetGroupHealthCheck(tt.lbSpec); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.getAPITargetGroupHealthCheck() Got unexpected result:\n%v", cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestService_getAdditionalTargetGroupHealthCheck(t *testing.T) {
	tests := []struct {
		name     string
		listener infrav1.AdditionalListenerSpec
		want     *infrav1.TargetGroupHealthCheck
		wantErr  bool
	}{
		{
			name: "TCP defaults",
			listener: infrav1.AdditionalListenerSpec{
				Protocol: "TCP",
				Port:     22623,
			},
			want: &infrav1.TargetGroupHealthCheck{
				Protocol:                aws.String("TCP"),
				Port:                    aws.String("22623"),
				IntervalSeconds:         aws.Int64(infrav1.DefaultAPIServerHealthCheckIntervalSec),
				TimeoutSeconds:          aws.Int64(infrav1.DefaultAPIServerHealthCheckTimeoutSec),
				ThresholdCount:          aws.Int64(infrav1.DefaultAPIServerHealthThresholdCount),
				UnhealthyThresholdCount: aws.Int64(infrav1.DefaultAPIServerUnhealthThresholdCount),
			},
		},
		{
			name: "Listener TCP, Health check protocol TCP, probe defaults",
			listener: infrav1.AdditionalListenerSpec{
				Port:     22623,
				Protocol: infrav1.ELBProtocolTCP,
			},
			want: &infrav1.TargetGroupHealthCheck{
				Protocol:                aws.String("TCP"),
				Port:                    aws.String("22623"),
				IntervalSeconds:         aws.Int64(infrav1.DefaultAPIServerHealthCheckIntervalSec),
				TimeoutSeconds:          aws.Int64(infrav1.DefaultAPIServerHealthCheckTimeoutSec),
				ThresholdCount:          aws.Int64(infrav1.DefaultAPIServerHealthThresholdCount),
				UnhealthyThresholdCount: aws.Int64(infrav1.DefaultAPIServerUnhealthThresholdCount),
			},
		},
		{
			name: "Listener TCP, Health check protocol HTTP, probe defaults",
			listener: infrav1.AdditionalListenerSpec{
				Port:     22623,
				Protocol: infrav1.ELBProtocolTCP,
				HealthCheck: &infrav1.TargetGroupHealthCheckAdditionalSpec{
					Protocol: aws.String("HTTP"),
					Path:     aws.String("/healthz"),
				},
			},
			want: &infrav1.TargetGroupHealthCheck{
				Protocol:                aws.String("HTTP"),
				Path:                    aws.String("/healthz"),
				Port:                    aws.String("22623"),
				IntervalSeconds:         aws.Int64(infrav1.DefaultAPIServerHealthCheckIntervalSec),
				TimeoutSeconds:          aws.Int64(infrav1.DefaultAPIServerHealthCheckTimeoutSec),
				ThresholdCount:          aws.Int64(infrav1.DefaultAPIServerHealthThresholdCount),
				UnhealthyThresholdCount: aws.Int64(infrav1.DefaultAPIServerUnhealthThresholdCount),
			},
		},
		{
			name: "Listener TCP, Health check protocol HTTP, probe customized",
			listener: infrav1.AdditionalListenerSpec{
				Port:     22623,
				Protocol: infrav1.ELBProtocolTCP,
				HealthCheck: &infrav1.TargetGroupHealthCheckAdditionalSpec{
					Protocol:                aws.String("HTTP"),
					Path:                    aws.String("/healthz"),
					IntervalSeconds:         aws.Int64(5),
					TimeoutSeconds:          aws.Int64(5),
					ThresholdCount:          aws.Int64(2),
					UnhealthyThresholdCount: aws.Int64(2),
				},
			},
			want: &infrav1.TargetGroupHealthCheck{
				Protocol:                aws.String("HTTP"),
				Port:                    aws.String("22623"),
				Path:                    aws.String("/healthz"),
				IntervalSeconds:         aws.Int64(5),
				TimeoutSeconds:          aws.Int64(5),
				ThresholdCount:          aws.Int64(2),
				UnhealthyThresholdCount: aws.Int64(2),
			},
		},
		{
			name: "Listener TCP, Health check protocol HTTPS, custom health check port and probes",
			listener: infrav1.AdditionalListenerSpec{
				Port:     22623,
				Protocol: infrav1.ELBProtocolTCP,
				HealthCheck: &infrav1.TargetGroupHealthCheckAdditionalSpec{
					Protocol:                aws.String("HTTPS"),
					Port:                    aws.String("22624"),
					Path:                    aws.String("/healthz"),
					IntervalSeconds:         aws.Int64(5),
					TimeoutSeconds:          aws.Int64(5),
					ThresholdCount:          aws.Int64(2),
					UnhealthyThresholdCount: aws.Int64(2),
				},
			},
			want: &infrav1.TargetGroupHealthCheck{
				Protocol:                aws.String("HTTPS"),
				Port:                    aws.String("22624"),
				Path:                    aws.String("/healthz"),
				IntervalSeconds:         aws.Int64(5),
				TimeoutSeconds:          aws.Int64(5),
				ThresholdCount:          aws.Int64(2),
				UnhealthyThresholdCount: aws.Int64(2),
			},
		},
		{
			name: "Listener TCP, Health check protocol TCP, custom health check port and probes, missing UnhealthyThresholdCount, want default",
			listener: infrav1.AdditionalListenerSpec{
				Port:     22623,
				Protocol: infrav1.ELBProtocolTCP,
				HealthCheck: &infrav1.TargetGroupHealthCheckAdditionalSpec{
					IntervalSeconds: aws.Int64(5),
					TimeoutSeconds:  aws.Int64(5),
					ThresholdCount:  aws.Int64(2),
				},
			},
			want: &infrav1.TargetGroupHealthCheck{
				Protocol:                aws.String("TCP"),
				Port:                    aws.String("22623"),
				IntervalSeconds:         aws.Int64(5),
				TimeoutSeconds:          aws.Int64(5),
				ThresholdCount:          aws.Int64(2),
				UnhealthyThresholdCount: aws.Int64(infrav1.DefaultAPIServerUnhealthThresholdCount),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := stubGetBaseService(t, "bar")
			if got := s.getAdditionalTargetGroupHealthCheck(tt.listener); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.getAdditionalTargetGroupHealthCheck() Got unexpected result:\n %v", cmp.Diff(got, tt.want))
			}
		})
	}
}
