/*
Copyright 2022 The Kubernetes Authors.

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

package gc

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	rgapi "github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/test/mocks"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func TestReconcileDelete(t *testing.T) {
	testCases := []struct {
		name         string
		clusterScope cloud.ClusterScoper
		elbMocks     func(m *mocks.MockELBAPIMockRecorder)
		elbv2Mocks   func(m *mocks.MockELBV2APIMockRecorder)
		rgAPIMocks   func(m *mocks.MockResourceGroupsTaggingAPIAPIMockRecorder)
		ec2Mocks     func(m *mocks.MockEC2APIMockRecorder)
		expectErr    bool
	}{
		{
			name:         "eks with no Service load balances",
			clusterScope: createManageScope(t),
			rgAPIMocks: func(m *mocks.MockResourceGroupsTaggingAPIAPIMockRecorder) {
				m.GetResourcesWithContext(gomock.Any(), &rgapi.GetResourcesInput{
					TagFilters: []*rgapi.TagFilter{
						{
							Key:    aws.String("kubernetes.io/cluster/eks-test-cluster"),
							Values: []*string{aws.String("owned")},
						},
					},
				}).DoAndReturn(func(awsCtx context.Context, input *rgapi.GetResourcesInput, opts ...request.Option) (*rgapi.GetResourcesOutput, error) {
					return &rgapi.GetResourcesOutput{
						ResourceTagMappingList: []*rgapi.ResourceTagMapping{},
					}, nil
				})
			},
			elbMocks:   func(m *mocks.MockELBAPIMockRecorder) {},
			elbv2Mocks: func(m *mocks.MockELBV2APIMockRecorder) {},
			ec2Mocks:   func(m *mocks.MockEC2APIMockRecorder) {},
			expectErr:  false,
		},
		{
			name:         "ec2 cluster with no Service load balances",
			clusterScope: createUnManageScope(t),
			rgAPIMocks: func(m *mocks.MockResourceGroupsTaggingAPIAPIMockRecorder) {
				m.GetResourcesWithContext(gomock.Any(), &rgapi.GetResourcesInput{
					TagFilters: []*rgapi.TagFilter{
						{
							Key:    aws.String("kubernetes.io/cluster/cluster1"),
							Values: []*string{aws.String("owned")},
						},
					},
				}).DoAndReturn(func(awsCtx context.Context, input *rgapi.GetResourcesInput, opts ...request.Option) (*rgapi.GetResourcesOutput, error) {
					return &rgapi.GetResourcesOutput{
						ResourceTagMappingList: []*rgapi.ResourceTagMapping{},
					}, nil
				})
			},
			elbMocks:   func(m *mocks.MockELBAPIMockRecorder) {},
			elbv2Mocks: func(m *mocks.MockELBV2APIMockRecorder) {},
			ec2Mocks:   func(m *mocks.MockEC2APIMockRecorder) {},
			expectErr:  false,
		},
		{
			name:         "eks with non Service load balancer",
			clusterScope: createManageScope(t),
			rgAPIMocks: func(m *mocks.MockResourceGroupsTaggingAPIAPIMockRecorder) {
				m.GetResourcesWithContext(gomock.Any(), &rgapi.GetResourcesInput{
					TagFilters: []*rgapi.TagFilter{
						{
							Key:    aws.String("kubernetes.io/cluster/eks-test-cluster"),
							Values: []*string{aws.String("owned")},
						},
					},
				}).DoAndReturn(func(awsCtx context.Context, input *rgapi.GetResourcesInput, opts ...request.Option) (*rgapi.GetResourcesOutput, error) {
					return &rgapi.GetResourcesOutput{
						ResourceTagMappingList: []*rgapi.ResourceTagMapping{
							{
								ResourceARN: aws.String("arn:aws:elasticloadbalancing:eu-west-2:1234567890:loadbalancer/aec24434cd2ce4630bd14a955413ee37"),
								Tags: []*rgapi.Tag{
									{
										Key:   aws.String("kubernetes.io/cluster/eks-test-cluster"),
										Value: aws.String("owned"),
									},
								},
							},
						},
					}, nil
				})
			},
			elbMocks:   func(m *mocks.MockELBAPIMockRecorder) {},
			elbv2Mocks: func(m *mocks.MockELBV2APIMockRecorder) {},
			ec2Mocks:   func(m *mocks.MockEC2APIMockRecorder) {},
			expectErr:  false,
		},
		{
			name:         "ec2 cluster with non Service load balancer",
			clusterScope: createUnManageScope(t),
			rgAPIMocks: func(m *mocks.MockResourceGroupsTaggingAPIAPIMockRecorder) {
				m.GetResourcesWithContext(gomock.Any(), &rgapi.GetResourcesInput{
					TagFilters: []*rgapi.TagFilter{
						{
							Key:    aws.String("kubernetes.io/cluster/cluster1"),
							Values: []*string{aws.String("owned")},
						},
					},
				}).DoAndReturn(func(awsCtx context.Context, input *rgapi.GetResourcesInput, opts ...request.Option) (*rgapi.GetResourcesOutput, error) {
					return &rgapi.GetResourcesOutput{
						ResourceTagMappingList: []*rgapi.ResourceTagMapping{
							{
								ResourceARN: aws.String("arn:aws:elasticloadbalancing:eu-west-2:1234567890:loadbalancer/aec24434cd2ce4630bd14a955413ee37"),
								Tags: []*rgapi.Tag{
									{
										Key:   aws.String("kubernetes.io/cluster/cluster1"),
										Value: aws.String("owned"),
									},
								},
							},
						},
					}, nil
				})
			},
			elbMocks:   func(m *mocks.MockELBAPIMockRecorder) {},
			elbv2Mocks: func(m *mocks.MockELBV2APIMockRecorder) {},
			ec2Mocks:   func(m *mocks.MockEC2APIMockRecorder) {},
			expectErr:  false,
		},
		{
			name:         "eks with ELB Service load balancer",
			clusterScope: createManageScope(t),
			rgAPIMocks: func(m *mocks.MockResourceGroupsTaggingAPIAPIMockRecorder) {
				m.GetResourcesWithContext(gomock.Any(), &rgapi.GetResourcesInput{
					TagFilters: []*rgapi.TagFilter{
						{
							Key:    aws.String("kubernetes.io/cluster/eks-test-cluster"),
							Values: []*string{aws.String("owned")},
						},
					},
				}).DoAndReturn(func(awsCtx context.Context, input *rgapi.GetResourcesInput, opts ...request.Option) (*rgapi.GetResourcesOutput, error) {
					return &rgapi.GetResourcesOutput{
						ResourceTagMappingList: []*rgapi.ResourceTagMapping{
							{
								ResourceARN: aws.String("arn:aws:elasticloadbalancing:eu-west-2:1234567890:loadbalancer/aec24434cd2ce4630bd14a955413ee37"),
								Tags: []*rgapi.Tag{
									{
										Key:   aws.String("kubernetes.io/cluster/eks-test-cluster"),
										Value: aws.String("owned"),
									},
									{
										Key:   aws.String(serviceNameTag),
										Value: aws.String("default/svc1"),
									},
								},
							},
						},
					}, nil
				})
			},
			elbMocks: func(m *mocks.MockELBAPIMockRecorder) {
				m.DeleteLoadBalancerWithContext(gomock.Any(), &elb.DeleteLoadBalancerInput{
					LoadBalancerName: aws.String("aec24434cd2ce4630bd14a955413ee37"),
				}).Return(&elb.DeleteLoadBalancerOutput{}, nil)
			},
			elbv2Mocks: func(m *mocks.MockELBV2APIMockRecorder) {},
			ec2Mocks:   func(m *mocks.MockEC2APIMockRecorder) {},
			expectErr:  false,
		},
		{
			name:         "ec2 cluster with ELB Service load balancer",
			clusterScope: createUnManageScope(t),
			rgAPIMocks: func(m *mocks.MockResourceGroupsTaggingAPIAPIMockRecorder) {
				m.GetResourcesWithContext(gomock.Any(), &rgapi.GetResourcesInput{
					TagFilters: []*rgapi.TagFilter{
						{
							Key:    aws.String("kubernetes.io/cluster/cluster1"),
							Values: []*string{aws.String("owned")},
						},
					},
				}).DoAndReturn(func(awsCtx context.Context, input *rgapi.GetResourcesInput, opts ...request.Option) (*rgapi.GetResourcesOutput, error) {
					return &rgapi.GetResourcesOutput{
						ResourceTagMappingList: []*rgapi.ResourceTagMapping{
							{
								ResourceARN: aws.String("arn:aws:elasticloadbalancing:eu-west-2:1234567890:loadbalancer/aec24434cd2ce4630bd14a955413ee37"),
								Tags: []*rgapi.Tag{
									{
										Key:   aws.String("kubernetes.io/cluster/cluster1"),
										Value: aws.String("owned"),
									},
									{
										Key:   aws.String(serviceNameTag),
										Value: aws.String("default/svc1"),
									},
								},
							},
						},
					}, nil
				})
			},
			elbMocks: func(m *mocks.MockELBAPIMockRecorder) {
				m.DeleteLoadBalancerWithContext(gomock.Any(), &elb.DeleteLoadBalancerInput{
					LoadBalancerName: aws.String("aec24434cd2ce4630bd14a955413ee37"),
				}).Return(&elb.DeleteLoadBalancerOutput{}, nil)
			},
			elbv2Mocks: func(m *mocks.MockELBV2APIMockRecorder) {},
			ec2Mocks:   func(m *mocks.MockEC2APIMockRecorder) {},
			expectErr:  false,
		},
		{
			name:         "eks with NLB Service load balancer",
			clusterScope: createManageScope(t),
			rgAPIMocks: func(m *mocks.MockResourceGroupsTaggingAPIAPIMockRecorder) {
				m.GetResourcesWithContext(gomock.Any(), &rgapi.GetResourcesInput{
					TagFilters: []*rgapi.TagFilter{
						{
							Key:    aws.String("kubernetes.io/cluster/eks-test-cluster"),
							Values: []*string{aws.String("owned")},
						},
					},
				}).DoAndReturn(func(awsCtx context.Context, input *rgapi.GetResourcesInput, opts ...request.Option) (*rgapi.GetResourcesOutput, error) {
					return &rgapi.GetResourcesOutput{
						ResourceTagMappingList: []*rgapi.ResourceTagMapping{
							{
								ResourceARN: aws.String("arn:aws:elasticloadbalancing:eu-west-2:1234567890:loadbalancer/net/aec24434cd2ce4630bd14a955413ee37"),
								Tags: []*rgapi.Tag{
									{
										Key:   aws.String("kubernetes.io/cluster/eks-test-cluster"),
										Value: aws.String("owned"),
									},
									{
										Key:   aws.String(serviceNameTag),
										Value: aws.String("default/svc1"),
									},
								},
							},
						},
					}, nil
				})
			},
			elbMocks: func(m *mocks.MockELBAPIMockRecorder) {},
			elbv2Mocks: func(m *mocks.MockELBV2APIMockRecorder) {
				m.DeleteLoadBalancerWithContext(gomock.Any(), &elbv2.DeleteLoadBalancerInput{
					LoadBalancerArn: aws.String("arn:aws:elasticloadbalancing:eu-west-2:1234567890:loadbalancer/net/aec24434cd2ce4630bd14a955413ee37"),
				}).Return(&elbv2.DeleteLoadBalancerOutput{}, nil)
			},
			ec2Mocks:  func(m *mocks.MockEC2APIMockRecorder) {},
			expectErr: false,
		},
		{
			name:         "ec2 cluster with NLB Service load balancer",
			clusterScope: createUnManageScope(t),
			rgAPIMocks: func(m *mocks.MockResourceGroupsTaggingAPIAPIMockRecorder) {
				m.GetResourcesWithContext(gomock.Any(), &rgapi.GetResourcesInput{
					TagFilters: []*rgapi.TagFilter{
						{
							Key:    aws.String("kubernetes.io/cluster/cluster1"),
							Values: []*string{aws.String("owned")},
						},
					},
				}).DoAndReturn(func(awsCtx context.Context, input *rgapi.GetResourcesInput, opts ...request.Option) (*rgapi.GetResourcesOutput, error) {
					return &rgapi.GetResourcesOutput{
						ResourceTagMappingList: []*rgapi.ResourceTagMapping{
							{
								ResourceARN: aws.String("arn:aws:elasticloadbalancing:eu-west-2:1234567890:loadbalancer/net/aec24434cd2ce4630bd14a955413ee37"),
								Tags: []*rgapi.Tag{
									{
										Key:   aws.String("kubernetes.io/cluster/cluster1"),
										Value: aws.String("owned"),
									},
									{
										Key:   aws.String(serviceNameTag),
										Value: aws.String("default/svc1"),
									},
								},
							},
						},
					}, nil
				})
			},
			elbMocks: func(m *mocks.MockELBAPIMockRecorder) {},
			elbv2Mocks: func(m *mocks.MockELBV2APIMockRecorder) {
				m.DeleteLoadBalancerWithContext(gomock.Any(), &elbv2.DeleteLoadBalancerInput{
					LoadBalancerArn: aws.String("arn:aws:elasticloadbalancing:eu-west-2:1234567890:loadbalancer/net/aec24434cd2ce4630bd14a955413ee37"),
				}).Return(&elbv2.DeleteLoadBalancerOutput{}, nil)
			},
			ec2Mocks:  func(m *mocks.MockEC2APIMockRecorder) {},
			expectErr: false,
		},
		{
			name:         "eks with ALB Service load balancer",
			clusterScope: createManageScope(t),
			rgAPIMocks: func(m *mocks.MockResourceGroupsTaggingAPIAPIMockRecorder) {
				m.GetResourcesWithContext(gomock.Any(), &rgapi.GetResourcesInput{
					TagFilters: []*rgapi.TagFilter{
						{
							Key:    aws.String("kubernetes.io/cluster/eks-test-cluster"),
							Values: []*string{aws.String("owned")},
						},
					},
				}).DoAndReturn(func(awsCtx context.Context, input *rgapi.GetResourcesInput, opts ...request.Option) (*rgapi.GetResourcesOutput, error) {
					return &rgapi.GetResourcesOutput{
						ResourceTagMappingList: []*rgapi.ResourceTagMapping{
							{
								ResourceARN: aws.String("arn:aws:elasticloadbalancing:eu-west-2:1234567890:loadbalancer/app/aec24434cd2ce4630bd14a955413ee37"),
								Tags: []*rgapi.Tag{
									{
										Key:   aws.String("kubernetes.io/cluster/eks-test-cluster"),
										Value: aws.String("owned"),
									},
									{
										Key:   aws.String(serviceNameTag),
										Value: aws.String("default/svc1"),
									},
								},
							},
						},
					}, nil
				})
			},
			elbMocks: func(m *mocks.MockELBAPIMockRecorder) {},
			elbv2Mocks: func(m *mocks.MockELBV2APIMockRecorder) {
				m.DeleteLoadBalancerWithContext(gomock.Any(), &elbv2.DeleteLoadBalancerInput{
					LoadBalancerArn: aws.String("arn:aws:elasticloadbalancing:eu-west-2:1234567890:loadbalancer/app/aec24434cd2ce4630bd14a955413ee37"),
				}).Return(&elbv2.DeleteLoadBalancerOutput{}, nil)
			},
			ec2Mocks:  func(m *mocks.MockEC2APIMockRecorder) {},
			expectErr: false,
		},
		{
			name:         "ec2 cluster with ALB Service load balancer",
			clusterScope: createUnManageScope(t),
			rgAPIMocks: func(m *mocks.MockResourceGroupsTaggingAPIAPIMockRecorder) {
				m.GetResourcesWithContext(gomock.Any(), &rgapi.GetResourcesInput{
					TagFilters: []*rgapi.TagFilter{
						{
							Key:    aws.String("kubernetes.io/cluster/cluster1"),
							Values: []*string{aws.String("owned")},
						},
					},
				}).DoAndReturn(func(awsCtx context.Context, input *rgapi.GetResourcesInput, opts ...request.Option) (*rgapi.GetResourcesOutput, error) {
					return &rgapi.GetResourcesOutput{
						ResourceTagMappingList: []*rgapi.ResourceTagMapping{
							{
								ResourceARN: aws.String("arn:aws:elasticloadbalancing:eu-west-2:1234567890:loadbalancer/app/aec24434cd2ce4630bd14a955413ee37"),
								Tags: []*rgapi.Tag{
									{
										Key:   aws.String("kubernetes.io/cluster/cluster1"),
										Value: aws.String("owned"),
									},
									{
										Key:   aws.String(serviceNameTag),
										Value: aws.String("default/svc1"),
									},
								},
							},
						},
					}, nil
				})
			},
			elbMocks: func(m *mocks.MockELBAPIMockRecorder) {},
			elbv2Mocks: func(m *mocks.MockELBV2APIMockRecorder) {
				m.DeleteLoadBalancerWithContext(gomock.Any(), &elbv2.DeleteLoadBalancerInput{
					LoadBalancerArn: aws.String("arn:aws:elasticloadbalancing:eu-west-2:1234567890:loadbalancer/app/aec24434cd2ce4630bd14a955413ee37"),
				}).Return(&elbv2.DeleteLoadBalancerOutput{}, nil)
			},
			ec2Mocks:  func(m *mocks.MockEC2APIMockRecorder) {},
			expectErr: false,
		},
		{
			name:         "eks cluster full test",
			clusterScope: createManageScope(t),
			rgAPIMocks: func(m *mocks.MockResourceGroupsTaggingAPIAPIMockRecorder) {
				m.GetResourcesWithContext(gomock.Any(), &rgapi.GetResourcesInput{
					TagFilters: []*rgapi.TagFilter{
						{
							Key:    aws.String("kubernetes.io/cluster/eks-test-cluster"),
							Values: []*string{aws.String("owned")},
						},
					},
				}).DoAndReturn(func(awsCtx context.Context, input *rgapi.GetResourcesInput, opts ...request.Option) (*rgapi.GetResourcesOutput, error) {
					return &rgapi.GetResourcesOutput{
						ResourceTagMappingList: []*rgapi.ResourceTagMapping{
							{
								ResourceARN: aws.String("arn:aws:elasticloadbalancing:eu-west-2:1234567890:targetgroup/k8s-default-podinfo-2c868b281a/e979fe9bd6825433"),
								Tags: []*rgapi.Tag{
									{
										Key:   aws.String("kubernetes.io/cluster/cluster1"),
										Value: aws.String("owned"),
									},
									{
										Key:   aws.String(serviceNameTag),
										Value: aws.String("default/svc1"),
									},
								},
							},
							{
								ResourceARN: aws.String("arn:aws:elasticloadbalancing:eu-west-2:1234567890:loadbalancer/aec24434cd2ce4630bd14a955413ee37"),
								Tags: []*rgapi.Tag{
									{
										Key:   aws.String("kubernetes.io/cluster/cluster1"),
										Value: aws.String("owned"),
									},
									{
										Key:   aws.String(serviceNameTag),
										Value: aws.String("default/svc1"),
									},
								},
							},
						},
					}, nil
				})
			},
			elbMocks: func(m *mocks.MockELBAPIMockRecorder) {
				m.DeleteLoadBalancerWithContext(gomock.Any(), &elb.DeleteLoadBalancerInput{
					LoadBalancerName: aws.String("aec24434cd2ce4630bd14a955413ee37"),
				}).Return(&elb.DeleteLoadBalancerOutput{}, nil)
			},
			elbv2Mocks: func(m *mocks.MockELBV2APIMockRecorder) {
				m.DeleteTargetGroupWithContext(gomock.Any(), &elbv2.DeleteTargetGroupInput{
					TargetGroupArn: aws.String("arn:aws:elasticloadbalancing:eu-west-2:1234567890:targetgroup/k8s-default-podinfo-2c868b281a/e979fe9bd6825433"),
				})
			},
			ec2Mocks:  func(m *mocks.MockEC2APIMockRecorder) {},
			expectErr: false,
		},
		{
			name:         "eks should ignore unhandled resources",
			clusterScope: createManageScope(t),
			rgAPIMocks: func(m *mocks.MockResourceGroupsTaggingAPIAPIMockRecorder) {
				m.GetResourcesWithContext(gomock.Any(), &rgapi.GetResourcesInput{
					TagFilters: []*rgapi.TagFilter{
						{
							Key:    aws.String("kubernetes.io/cluster/eks-test-cluster"),
							Values: []*string{aws.String("owned")},
						},
					},
				}).DoAndReturn(func(awsCtx context.Context, input *rgapi.GetResourcesInput, opts ...request.Option) (*rgapi.GetResourcesOutput, error) {
					return &rgapi.GetResourcesOutput{
						ResourceTagMappingList: []*rgapi.ResourceTagMapping{
							{
								ResourceARN: aws.String("arn:aws:ec2:eu-west-2:217426147237:s3/somebucket"),
								Tags: []*rgapi.Tag{
									{
										Key:   aws.String("kubernetes.io/cluster/eks-test-cluster"),
										Value: aws.String("owned"),
									},
									{
										Key:   aws.String(serviceNameTag),
										Value: aws.String("default/svc1"),
									},
									{
										Key:   aws.String("Name"),
										Value: aws.String("eks-cluster-sg-default_capi-managed-test-control-plane-10156951"),
									},
								},
							},
						},
					}, nil
				})
			},
			elbMocks:   func(m *mocks.MockELBAPIMockRecorder) {},
			elbv2Mocks: func(m *mocks.MockELBV2APIMockRecorder) {},
			ec2Mocks:   func(m *mocks.MockEC2APIMockRecorder) {},
			expectErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			rgapiMock := mocks.NewMockResourceGroupsTaggingAPIAPI(mockCtrl)
			elbapiMock := mocks.NewMockELBAPI(mockCtrl)
			elbv2Mock := mocks.NewMockELBV2API(mockCtrl)
			ec2Mock := mocks.NewMockEC2API(mockCtrl)

			tc.rgAPIMocks(rgapiMock.EXPECT())
			tc.elbMocks(elbapiMock.EXPECT())
			tc.elbv2Mocks(elbv2Mock.EXPECT())
			tc.ec2Mocks(ec2Mock.EXPECT())

			ctx := context.TODO()

			opts := []ServiceOption{
				WithELBClient(elbapiMock),
				WithELBv2Client(elbv2Mock),
				WithResourceTaggingClient(rgapiMock),
				WithEC2Client(ec2Mock),
			}
			wkSvc := NewService(tc.clusterScope, opts...)
			err := wkSvc.ReconcileDelete(ctx)

			if tc.expectErr {
				g.Expect(err).NotTo(BeNil())
				return
			}

			g.Expect(err).To(BeNil())
		})
	}
}

func createManageScope(t *testing.T) *scope.ManagedControlPlaneScope {
	t.Helper()
	g := NewWithT(t)

	cluster := createEKSCluster()
	cp := createManagedControlPlane()
	objs := []client.Object{cluster, cp}

	scheme := createScheme()
	client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build()

	managedScope, err := scope.NewManagedControlPlaneScope(scope.ManagedControlPlaneScopeParams{
		Client:         client,
		Cluster:        cluster,
		ControlPlane:   cp,
		ControllerName: "test-controller",
	})
	g.Expect(err).NotTo(HaveOccurred())

	return managedScope
}

func createUnManageScope(t *testing.T) *scope.ClusterScope {
	t.Helper()
	g := NewWithT(t)

	cluster := createUnmanagedCluster()
	awsCluster := createAWSCluser()
	objs := []client.Object{cluster, awsCluster}

	scheme := createScheme()
	client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build()

	clusterScope, err := scope.NewClusterScope(scope.ClusterScopeParams{
		Client:         client,
		Cluster:        cluster,
		AWSCluster:     awsCluster,
		ControllerName: "test-controller",
	})
	g.Expect(err).NotTo(HaveOccurred())

	return clusterScope
}

func createScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()
	_ = corev1.AddToScheme(scheme)
	_ = ekscontrolplanev1.AddToScheme(scheme)
	_ = infrav1.AddToScheme(scheme)
	_ = clusterv1.AddToScheme(scheme)

	return scheme
}

func createEKSCluster() *clusterv1.Cluster {
	return &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cluster1",
			Namespace: "default",
		},
		Spec: clusterv1.ClusterSpec{
			InfrastructureRef: &corev1.ObjectReference{
				Kind:       "AWSManagedControlPlane",
				APIVersion: ekscontrolplanev1.GroupVersion.String(),
				Name:       "cp1",
				Namespace:  "default",
			},
		},
	}
}

func createManagedControlPlane() *ekscontrolplanev1.AWSManagedControlPlane {
	return &ekscontrolplanev1.AWSManagedControlPlane{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AWSManagedControlPlane",
			APIVersion: ekscontrolplanev1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cp1",
			Namespace: "default",
		},
		Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
			EKSClusterName: "eks-test-cluster",
		},
	}
}

func createAWSCluser() *infrav1.AWSCluster {
	return &infrav1.AWSCluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AWSCluster",
			APIVersion: infrav1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cluster1",
			Namespace: "default",
		},
		Spec: infrav1.AWSClusterSpec{},
	}
}

func createUnmanagedCluster() *clusterv1.Cluster {
	return &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cluster1",
			Namespace: "default",
		},
		Spec: clusterv1.ClusterSpec{
			InfrastructureRef: &corev1.ObjectReference{
				Kind:       "AWSCluster",
				APIVersion: infrav1.GroupVersion.String(),
				Name:       "cluster1",
				Namespace:  "default",
			},
		},
	}
}
