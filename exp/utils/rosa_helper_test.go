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

	// Build a NodePool with various fields populated
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
	// Build the expected spec
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
