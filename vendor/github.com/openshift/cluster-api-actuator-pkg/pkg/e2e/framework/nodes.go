package framework

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// GetWorkerNodes returns all nodes with the nodeWorkerRoleLabel label
func GetWorkerNodes(client runtimeclient.Client) ([]corev1.Node, error) {
	listOptions := runtimeclient.ListOptions{
		Namespace: TestContext.MachineApiNamespace,
	}
	listOptions.MatchingLabels(map[string]string{WorkerNodeRoleLabel: ""})
	workerNodes := &corev1.NodeList{}
	err := client.List(context.TODO(), &listOptions, workerNodes)
	if err != nil {
		return nil, err
	}
	return workerNodes.Items, nil
}

// FilterReadyNodes fileter the list of nodes and returns the list with ready nodes
func FilterReadyNodes(nodes []corev1.Node) []corev1.Node {
	var readyNodes []corev1.Node
	for _, n := range nodes {
		if IsNodeReady(&n) {
			readyNodes = append(readyNodes, n)
		}
	}
	return readyNodes
}
