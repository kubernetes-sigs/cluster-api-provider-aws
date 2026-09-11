/*
Copyright 2025 The Kubernetes Authors.

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

package utils

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	. "github.com/onsi/gomega"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"

	rosacontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
)

func TestNodePoolToRosaMachinePoolSpec(t *testing.T) {
	g := NewWithT(t)

	builder := cmv1.NewNodePool().
		ID("test-nodepool").
		Version(cmv1.NewVersion().ID("openshift-v4.15.0")).
		AvailabilityZone("us-east-1a").
		Subnet("subnet-12345").
		Labels(map[string]string{"role": "worker"}).
		AutoRepair(true).
		TuningConfigs("tuning1").
		AWSNodePool(
			cmv1.NewAWSNodePool().
				InstanceType("m5.large").
				AdditionalSecurityGroupIds("sg-123", "sg-456").
				RootVolume(cmv1.NewAWSVolume().Size(120)),
		).
		Autoscaling(
			cmv1.NewNodePoolAutoscaling().
				MinReplica(2).
				MaxReplica(5),
		).
		Taints(
			cmv1.NewTaint().Key("dedicated").Value("gpu").Effect(string(corev1.TaintEffectNoSchedule)),
		).
		NodeDrainGracePeriod(
			cmv1.NewValue().Value(10),
		).
		ManagementUpgrade(
			cmv1.NewNodePoolManagementUpgrade().
				MaxSurge("1").
				MaxUnavailable("2"),
		)

	nodePool, err := builder.Build()
	g.Expect(err).ToNot(HaveOccurred())

	actualSpec := NodePoolToRosaMachinePoolSpec(nodePool)
	expectedSpec := expinfrav1.RosaMachinePoolSpec{
		NodePoolName:             "test-nodepool",
		Version:                  "4.15.0",
		AvailabilityZone:         "us-east-1a",
		Subnet:                   "subnet-12345",
		Labels:                   map[string]string{"role": "worker"},
		AutoRepair:               true,
		InstanceType:             "m5.large",
		TuningConfigs:            []string{"tuning1"},
		AdditionalSecurityGroups: []string{"sg-123", "sg-456"},
		VolumeSize:               120,
		Autoscaling: &rosacontrolplanev1.AutoScaling{
			MinReplicas: 2,
			MaxReplicas: 5,
		},
		Taints: []expinfrav1.RosaTaint{
			{
				Key:    "dedicated",
				Value:  "gpu",
				Effect: corev1.TaintEffectNoSchedule,
			},
		},
		NodeDrainGracePeriod: &metav1.Duration{Duration: 10 * time.Minute},
		UpdateConfig: &expinfrav1.RosaUpdateConfig{
			RollingUpdate: &expinfrav1.RollingUpdate{
				MaxSurge:       ptr.To(intstr.FromInt32(1)),
				MaxUnavailable: ptr.To(intstr.FromInt32(2)),
			},
		},
	}

	g.Expect(expectedSpec).To(Equal(actualSpec))
}

func TestNodePoolToRosaMachinePoolSpec_DefaultsWhenOCMFieldsAbsent(t *testing.T) {
	tests := []struct {
		name     string
		builder  *cmv1.NodePoolBuilder
		expected expinfrav1.RosaMachinePoolSpec
	}{
		{
			name: "no NodeDrainGracePeriod or ManagementUpgrade from OCM",
			builder: cmv1.NewNodePool().
				ID("test-nodepool").
				AutoRepair(true).
				AWSNodePool(cmv1.NewAWSNodePool().InstanceType("m5.large")),
			expected: expinfrav1.RosaMachinePoolSpec{
				NodePoolName:         "test-nodepool",
				AutoRepair:           true,
				InstanceType:         "m5.large",
				NodeDrainGracePeriod: &metav1.Duration{},
				UpdateConfig: &expinfrav1.RosaUpdateConfig{
					RollingUpdate: &expinfrav1.RollingUpdate{
						MaxUnavailable: ptr.To(intstr.FromInt32(0)),
						MaxSurge:       ptr.To(intstr.FromInt32(1)),
					},
				},
			},
		},
		{
			name: "ManagementUpgrade present but MaxSurge and MaxUnavailable empty",
			builder: cmv1.NewNodePool().
				ID("test-nodepool").
				AutoRepair(true).
				AWSNodePool(cmv1.NewAWSNodePool().InstanceType("m5.large")).
				ManagementUpgrade(cmv1.NewNodePoolManagementUpgrade()),
			expected: expinfrav1.RosaMachinePoolSpec{
				NodePoolName:         "test-nodepool",
				AutoRepair:           true,
				InstanceType:         "m5.large",
				NodeDrainGracePeriod: &metav1.Duration{},
				UpdateConfig: &expinfrav1.RosaUpdateConfig{
					RollingUpdate: &expinfrav1.RollingUpdate{
						MaxUnavailable: ptr.To(intstr.FromInt32(0)),
						MaxSurge:       ptr.To(intstr.FromInt32(1)),
					},
				},
			},
		},
		{
			name: "ManagementUpgrade with only MaxSurge set",
			builder: cmv1.NewNodePool().
				ID("test-nodepool").
				AutoRepair(true).
				AWSNodePool(cmv1.NewAWSNodePool().InstanceType("m5.large")).
				ManagementUpgrade(
					cmv1.NewNodePoolManagementUpgrade().MaxSurge("3"),
				),
			expected: expinfrav1.RosaMachinePoolSpec{
				NodePoolName:         "test-nodepool",
				AutoRepair:           true,
				InstanceType:         "m5.large",
				NodeDrainGracePeriod: &metav1.Duration{},
				UpdateConfig: &expinfrav1.RosaUpdateConfig{
					RollingUpdate: &expinfrav1.RollingUpdate{
						MaxUnavailable: ptr.To(intstr.FromInt32(0)),
						MaxSurge:       ptr.To(intstr.FromInt32(3)),
					},
				},
			},
		},
		{
			name: "ManagementUpgrade with only MaxUnavailable set",
			builder: cmv1.NewNodePool().
				ID("test-nodepool").
				AutoRepair(true).
				AWSNodePool(cmv1.NewAWSNodePool().InstanceType("m5.large")).
				ManagementUpgrade(
					cmv1.NewNodePoolManagementUpgrade().MaxUnavailable("2"),
				),
			expected: expinfrav1.RosaMachinePoolSpec{
				NodePoolName:         "test-nodepool",
				AutoRepair:           true,
				InstanceType:         "m5.large",
				NodeDrainGracePeriod: &metav1.Duration{},
				UpdateConfig: &expinfrav1.RosaUpdateConfig{
					RollingUpdate: &expinfrav1.RollingUpdate{
						MaxUnavailable: ptr.To(intstr.FromInt32(2)),
						MaxSurge:       ptr.To(intstr.FromInt32(1)),
					},
				},
			},
		},
		{
			name: "NodeDrainGracePeriod with zero value from OCM",
			builder: cmv1.NewNodePool().
				ID("test-nodepool").
				AutoRepair(true).
				AWSNodePool(cmv1.NewAWSNodePool().InstanceType("m5.large")).
				NodeDrainGracePeriod(cmv1.NewValue().Value(0)),
			expected: expinfrav1.RosaMachinePoolSpec{
				NodePoolName:         "test-nodepool",
				AutoRepair:           true,
				InstanceType:         "m5.large",
				NodeDrainGracePeriod: &metav1.Duration{},
				UpdateConfig: &expinfrav1.RosaUpdateConfig{
					RollingUpdate: &expinfrav1.RollingUpdate{
						MaxUnavailable: ptr.To(intstr.FromInt32(0)),
						MaxSurge:       ptr.To(intstr.FromInt32(1)),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			nodePool, err := tt.builder.Build()
			g.Expect(err).ToNot(HaveOccurred())

			actualSpec := NodePoolToRosaMachinePoolSpec(nodePool)
			g.Expect(actualSpec).To(Equal(tt.expected))
		})
	}
}

func TestNodePoolToRosaMachinePoolSpec_NoPhantomDiff(t *testing.T) {
	tests := []struct {
		name    string
		builder *cmv1.NodePoolBuilder
	}{
		{
			name: "OCM returns no NodeDrainGracePeriod or ManagementUpgrade",
			builder: cmv1.NewNodePool().
				ID("test-nodepool").
				AutoRepair(true).
				AWSNodePool(cmv1.NewAWSNodePool().InstanceType("m5.large")),
		},
		{
			name: "OCM returns ManagementUpgrade with empty MaxSurge/MaxUnavailable",
			builder: cmv1.NewNodePool().
				ID("test-nodepool").
				AutoRepair(true).
				AWSNodePool(cmv1.NewAWSNodePool().InstanceType("m5.large")).
				ManagementUpgrade(cmv1.NewNodePoolManagementUpgrade()),
		},
		{
			name: "OCM returns NodeDrainGracePeriod with zero value",
			builder: cmv1.NewNodePool().
				ID("test-nodepool").
				AutoRepair(true).
				AWSNodePool(cmv1.NewAWSNodePool().InstanceType("m5.large")).
				NodeDrainGracePeriod(cmv1.NewValue().Value(0)),
		},
		{
			name: "OCM returns ManagementUpgrade with default values",
			builder: cmv1.NewNodePool().
				ID("test-nodepool").
				AutoRepair(true).
				AWSNodePool(cmv1.NewAWSNodePool().InstanceType("m5.large")).
				ManagementUpgrade(
					cmv1.NewNodePoolManagementUpgrade().MaxSurge("1").MaxUnavailable("0"),
				),
		},
		{
			name: "fully populated OCM response",
			builder: cmv1.NewNodePool().
				ID("test-nodepool").
				AutoRepair(true).
				AWSNodePool(cmv1.NewAWSNodePool().InstanceType("m5.large")).
				NodeDrainGracePeriod(cmv1.NewValue().Value(10)).
				ManagementUpgrade(
					cmv1.NewNodePoolManagementUpgrade().MaxSurge("3").MaxUnavailable("1"),
				),
		},
	}

	ignoredFields := []string{
		"ProviderIDList",
		"Version",
		"AdditionalTags",
		"AdditionalSecurityGroups",
		"VolumeSize",
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			nodePool, err := tt.builder.Build()
			g.Expect(err).ToNot(HaveOccurred())

			currentSpec := NodePoolToRosaMachinePoolSpec(nodePool)

			desiredPool := &expinfrav1.ROSAMachinePool{
				Spec: currentSpec,
			}
			desiredPool.Default()
			desiredSpec := desiredPool.Spec

			diff := cmp.Diff(desiredSpec, currentSpec,
				cmpopts.EquateEmpty(),
				cmpopts.IgnoreFields(currentSpec, ignoredFields...))
			g.Expect(diff).To(BeEmpty(), "phantom diff detected between defaulted spec and OCM conversion: %s", diff)
		})
	}
}
