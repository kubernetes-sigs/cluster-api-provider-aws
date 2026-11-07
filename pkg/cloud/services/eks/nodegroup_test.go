/*
Copyright 2020 The Kubernetes Authors.

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

package eks

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/go-logr/logr"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/eks/iam"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/eks/mock_eksiface"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
)

func TestScalingConfig(t *testing.T) {
	tests := []struct {
		name                      string
		replicas                  *int32
		minSize                   *int32
		maxSize                   *int32
		externalAutoscalerManaged bool
		expectDesiredSize         *int32
		expectMinSize             *int32
		expectMaxSize             *int32
	}{
		{
			name:                      "external autoscaler ignores non-nil replicas",
			replicas:                  aws.Int32(3),
			minSize:                   aws.Int32(1),
			maxSize:                   aws.Int32(5),
			externalAutoscalerManaged: true,
			expectDesiredSize:         nil,
			expectMinSize:             aws.Int32(1),
			expectMaxSize:             aws.Int32(5),
		},
		{
			name:              "replicas within bounds",
			replicas:          aws.Int32(3),
			minSize:           aws.Int32(1),
			maxSize:           aws.Int32(5),
			expectDesiredSize: aws.Int32(3),
			expectMinSize:     aws.Int32(1),
			expectMaxSize:     aws.Int32(5),
		},
		{
			name:              "replicas below minSize",
			replicas:          aws.Int32(0),
			minSize:           aws.Int32(2),
			maxSize:           aws.Int32(5),
			expectDesiredSize: aws.Int32(2),
			expectMinSize:     aws.Int32(2),
			expectMaxSize:     aws.Int32(5),
		},
		{
			name:              "replicas above maxSize",
			replicas:          aws.Int32(10),
			minSize:           aws.Int32(1),
			maxSize:           aws.Int32(5),
			expectDesiredSize: aws.Int32(5),
			expectMinSize:     aws.Int32(1),
			expectMaxSize:     aws.Int32(5),
		},
		{
			name:              "nil replicas defaults to 1",
			replicas:          nil,
			minSize:           aws.Int32(0),
			maxSize:           aws.Int32(5),
			expectDesiredSize: aws.Int32(1),
			expectMinSize:     aws.Int32(0),
			expectMaxSize:     aws.Int32(5),
		},
		{
			name:              "nil replicas clamped to minSize",
			replicas:          nil,
			minSize:           aws.Int32(3),
			maxSize:           aws.Int32(5),
			expectDesiredSize: aws.Int32(3),
			expectMinSize:     aws.Int32(3),
			expectMaxSize:     aws.Int32(5),
		},
		{
			name:              "no scaling config",
			replicas:          aws.Int32(3),
			minSize:           nil,
			maxSize:           nil,
			expectDesiredSize: aws.Int32(3),
			expectMinSize:     nil,
			expectMaxSize:     nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			machinePool := &clusterv1.MachinePool{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pool",
					Namespace: "default",
				},
				Spec: clusterv1.MachinePoolSpec{
					Replicas: tc.replicas,
				},
			}

			if tc.externalAutoscalerManaged {
				machinePool.ObjectMeta.Annotations = map[string]string{
					clusterv1.ReplicasManagedByAnnotation: "cluster-autoscaler",
				}
			}

			var scaling *expinfrav1.ManagedMachinePoolScaling
			if tc.minSize != nil || tc.maxSize != nil {
				scaling = &expinfrav1.ManagedMachinePoolScaling{
					MinSize: tc.minSize,
					MaxSize: tc.maxSize,
				}
			}

			managedMachinePool := &expinfrav1.AWSManagedMachinePool{
				Spec: expinfrav1.AWSManagedMachinePoolSpec{
					Scaling: scaling,
				},
			}

			mockScope := &scope.ManagedMachinePoolScope{
				MachinePool:        machinePool,
				ManagedMachinePool: managedMachinePool,
			}

			service := &NodegroupService{
				scope: mockScope,
			}

			cfg := service.scalingConfig()

			if tc.expectDesiredSize == nil {
				g.Expect(cfg.DesiredSize).To(BeNil())
			} else {
				g.Expect(cfg.DesiredSize).ToNot(BeNil())
				g.Expect(*cfg.DesiredSize).To(Equal(*tc.expectDesiredSize))
			}

			if tc.expectMinSize == nil {
				g.Expect(cfg.MinSize).To(BeNil())
			} else {
				g.Expect(cfg.MinSize).ToNot(BeNil())
				g.Expect(*cfg.MinSize).To(Equal(*tc.expectMinSize))
			}

			if tc.expectMaxSize == nil {
				g.Expect(cfg.MaxSize).To(BeNil())
			} else {
				g.Expect(cfg.MaxSize).ToNot(BeNil())
				g.Expect(*cfg.MaxSize).To(Equal(*tc.expectMaxSize))
			}
		})
	}
}

func TestReconcileNodegroupConfigWithExternalAutoscaler(t *testing.T) {
	tests := []struct {
		name               string
		scaling            *expinfrav1.ManagedMachinePoolScaling
		nodegroupScaling   *ekstypes.NodegroupScalingConfig
		expectConfigUpdate bool
	}{
		{
			name: "ignores desired size drift",
			nodegroupScaling: &ekstypes.NodegroupScalingConfig{
				DesiredSize: aws.Int32(4),
				MinSize:     aws.Int32(1),
				MaxSize:     aws.Int32(5),
			},
		},
		{
			name: "reconciles min and max size drift",
			scaling: &expinfrav1.ManagedMachinePoolScaling{
				MinSize: aws.Int32(2),
				MaxSize: aws.Int32(6),
			},
			nodegroupScaling: &ekstypes.NodegroupScalingConfig{
				DesiredSize: aws.Int32(4),
				MinSize:     aws.Int32(1),
				MaxSize:     aws.Int32(5),
			},
			expectConfigUpdate: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			mockControl := gomock.NewController(t)
			eksMock := mock_eksiface.NewMockEKSAPI(mockControl)

			if tc.expectConfigUpdate {
				eksMock.EXPECT().
					UpdateNodegroupConfig(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, input *eks.UpdateNodegroupConfigInput, _ ...func(*eks.Options)) (*eks.UpdateNodegroupConfigOutput, error) {
						g.Expect(input.ScalingConfig).ToNot(BeNil())
						g.Expect(input.ScalingConfig.DesiredSize).To(BeNil())
						g.Expect(input.ScalingConfig.MinSize).To(Equal(aws.Int32(2)))
						g.Expect(input.ScalingConfig.MaxSize).To(Equal(aws.Int32(6)))
						return &eks.UpdateNodegroupConfigOutput{}, nil
					})
			}

			machinePool := &clusterv1.MachinePool{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						clusterv1.ReplicasManagedByAnnotation: "cluster-autoscaler",
					},
				},
				Spec: clusterv1.MachinePoolSpec{
					Replicas: aws.Int32(3),
				},
			}
			managedMachinePool := &expinfrav1.AWSManagedMachinePool{
				Spec: expinfrav1.AWSManagedMachinePoolSpec{
					EKSNodegroupName: "test-nodegroup",
					Scaling:          tc.scaling,
				},
			}
			mockScope := &scope.ManagedMachinePoolScope{
				Logger:             *logger.NewLogger(logr.Discard()),
				MachinePool:        machinePool,
				ManagedMachinePool: managedMachinePool,
				ControlPlane: &ekscontrolplanev1.AWSManagedControlPlane{
					Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
						EKSClusterName: "test-cluster",
					},
				},
			}
			service := &NodegroupService{
				scope:      mockScope,
				EKSClient:  eksMock,
				IAMService: iam.IAMService{Wrapper: &mockScope.Logger},
			}
			nodegroup := &ekstypes.Nodegroup{
				NodegroupName:    aws.String("test-nodegroup"),
				ScalingConfig:    tc.nodegroupScaling,
				NodeRepairConfig: &ekstypes.NodeRepairConfig{Enabled: aws.Bool(false)},
			}

			g.Expect(service.reconcileNodegroupConfig(context.Background(), nodegroup)).To(Succeed())
		})
	}
}
