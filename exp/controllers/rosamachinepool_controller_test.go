package controllers

import (
	"testing"

	. "github.com/onsi/gomega"
	"k8s.io/utils/ptr"

	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	expclusterv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
)

func TestNodePoolToRosaMachinePoolSpec(t *testing.T) {
	g := NewWithT(t)

	rosaMachinePoolSpec := expinfrav1.RosaMachinePoolSpec{
		NodePoolName:  "test-nodepool",
		Version:       "4.14.5",
		Subnet:        "subnet-id",
		AutoRepair:    true,
		InstanceType:  "m5.large",
		TuningConfigs: []string{"config1"},
	}

	machinePoolSpec := expclusterv1.MachinePoolSpec{
		Replicas: ptr.To[int32](2),
	}

	nodePoolBuilder := nodePoolBuilder(rosaMachinePoolSpec, machinePoolSpec)

	nodePoolSpec, err := nodePoolBuilder.Build()
	g.Expect(err).ToNot(HaveOccurred())

	expectedSpec := nodePoolToRosaMachinePoolSpec(nodePoolSpec)

	g.Expect(expectedSpec).To(Equal(rosaMachinePoolSpec))
}
