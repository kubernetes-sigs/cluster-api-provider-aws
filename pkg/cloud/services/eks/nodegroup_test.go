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
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
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
			name:                      "external autoscaler",
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
