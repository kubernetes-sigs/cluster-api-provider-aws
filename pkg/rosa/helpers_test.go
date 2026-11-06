package rosa

import (
	"testing"

	. "github.com/onsi/gomega"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

func TestIsNodePoolReady(t *testing.T) {
	tests := []struct {
		name          string
		nodePool      *cmv1.NodePool
		expectedReady bool
	}{
		{
			name:          "fixed replicas:0 without autoscaling is ready when no replicas are running",
			nodePool:      nodePool(t, 0, nil, "", 0),
			expectedReady: true,
		},
		{
			name:          "autoscaling with minReplicas:0 is ready when no replicas are running",
			nodePool:      nodePool(t, 0, autoscalingBuilder(0, 1), "", 0),
			expectedReady: true,
		},
		{
			name:          "fixed replicas are ready only once all replicas are running",
			nodePool:      nodePool(t, 2, nil, "", 1),
			expectedReady: false,
		},
		{
			name:          "fixed replicas are ready when all replicas are running",
			nodePool:      nodePool(t, 2, nil, "", 2),
			expectedReady: true,
		},
		{
			name:          "a provisioning error message keeps the node pool not ready",
			nodePool:      nodePool(t, 0, nil, "NoReplicas: ", 0),
			expectedReady: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			g.Expect(IsNodePoolReady(tt.nodePool)).To(Equal(tt.expectedReady))
		})
	}
}

func autoscalingBuilder(minReplicas, maxReplicas int) *cmv1.NodePoolAutoscalingBuilder {
	return cmv1.NewNodePoolAutoscaling().MinReplica(minReplicas).MaxReplica(maxReplicas)
}

func nodePool(t *testing.T, replicas int, autoscaling *cmv1.NodePoolAutoscalingBuilder, message string, currentReplicas int) *cmv1.NodePool {
	t.Helper()
	builder := cmv1.NewNodePool().
		Replicas(replicas).
		Status(cmv1.NewNodePoolStatus().
			CurrentReplicas(currentReplicas).
			Message(message))
	if autoscaling != nil {
		builder = builder.Autoscaling(autoscaling)
	}
	np, err := builder.Build()
	NewWithT(t).Expect(err).ToNot(HaveOccurred())
	return np
}
