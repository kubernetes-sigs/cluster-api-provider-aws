package infra

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang/glog"
	e2e "github.com/openshift/cluster-api-actuator-pkg/pkg/e2e/framework"
	mapiv1beta1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	controllernode "github.com/openshift/cluster-api/pkg/controller/node"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/scale"
	"k8s.io/utils/pointer"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

const (
	machineRoleLabel = "machine.openshift.io/cluster-api-machine-role"
	machineAPIGroup  = "machine.openshift.io"
)

func isOneMachinePerNode(client runtimeclient.Client) bool {
	listOptions := runtimeclient.ListOptions{
		Namespace: e2e.TestContext.MachineApiNamespace,
	}
	machineList := mapiv1beta1.MachineList{}
	nodeList := corev1.NodeList{}

	if err := wait.PollImmediate(5*time.Second, e2e.WaitShort, func() (bool, error) {
		if err := client.List(context.TODO(), &listOptions, &machineList); err != nil {
			glog.Errorf("Error querying api for machineList object: %v, retrying...", err)
			return false, nil
		}
		if err := client.List(context.TODO(), &listOptions, &nodeList); err != nil {
			glog.Errorf("Error querying api for nodeList object: %v, retrying...", err)
			return false, nil
		}

		glog.Infof("Expecting the same number of machines and nodes, have %d nodes and %d machines", len(nodeList.Items), len(machineList.Items))
		if len(machineList.Items) != len(nodeList.Items) {
			return false, nil
		}

		nodeNameToMachineAnnotation := make(map[string]string)
		for _, node := range nodeList.Items {
			if _, ok := node.Annotations[controllernode.MachineAnnotationKey]; !ok {
				glog.Errorf("Node %q does not have a MachineAnnotationKey %q, retrying...", node.Name, controllernode.MachineAnnotationKey)
				return false, nil
			}
			nodeNameToMachineAnnotation[node.Name] = node.Annotations[controllernode.MachineAnnotationKey]
		}
		for _, machine := range machineList.Items {
			if machine.Status.NodeRef == nil {
				glog.Errorf("Machine %q has no NodeRef, retrying...", machine.Name)
				return false, nil
			}
			nodeName := machine.Status.NodeRef.Name
			if nodeNameToMachineAnnotation[nodeName] != fmt.Sprintf("%s/%s", e2e.TestContext.MachineApiNamespace, machine.Name) {
				glog.Errorf("Node name %q does not match expected machine name %q, retrying...", nodeName, machine.Name)
				return false, nil
			}
			glog.Infof("Machine %q is linked to node %q", machine.Name, nodeName)
		}
		return true, nil
	}); err != nil {
		glog.Errorf("Error checking isOneMachinePerNode: %v", err)
		return false
	}
	return true
}

// getClusterSize returns the number of nodes of the cluster
func getClusterSize(client runtimeclient.Client) (int, error) {
	nodes, err := e2e.GetNodes(client)
	if err != nil {
		return 0, fmt.Errorf("error getting nodes: %v", err)
	}

	glog.Infof("Cluster size is %d nodes", len(nodes))
	return len(nodes), nil
}

// machineSetsSnapShotLogs logs the state of all the machineSets in the cluster
func machineSetsSnapShotLogs(client runtimeclient.Client) error {
	machineSets, err := e2e.GetMachineSets(context.TODO(), client)
	if err != nil {
		return fmt.Errorf("error getting machines: %v", err)
	}

	for _, machineset := range machineSets {
		glog.Infof("MachineSet %q replicas %d. Ready: %d, available %d",
			machineset.Name,
			pointer.Int32PtrDerefOr(machineset.Spec.Replicas, e2e.DefaultMachineSetReplicas),
			machineset.Status.ReadyReplicas,
			machineset.Status.AvailableReplicas)
	}
	return nil
}

// getMachinesFromMachineSet returns an array of machines owned by a given machineSet
func getMachinesFromMachineSet(client runtimeclient.Client, machineSet mapiv1beta1.MachineSet) ([]mapiv1beta1.Machine, error) {
	machines, err := e2e.GetMachines(context.TODO(), client)
	if err != nil {
		return nil, fmt.Errorf("error getting machines: %v", err)
	}
	var machinesForSet []mapiv1beta1.Machine
	for key := range machines {
		if metav1.IsControlledBy(&machines[key], &machineSet) {
			machinesForSet = append(machinesForSet, machines[key])
		}
	}
	return machinesForSet, nil
}

// deleteMachine deletes a specific machine and returns an error otherwise
func deleteMachine(client runtimeclient.Client, machine *mapiv1beta1.Machine) error {
	return wait.PollImmediate(1*time.Second, time.Minute, func() (bool, error) {
		if err := client.Delete(context.TODO(), machine); err != nil {
			glog.Errorf("Error querying api for machine object %q: %v, retrying...", machine.Name, err)
			return false, err
		}
		return true, nil
	})
}

// getNodesFromMachineSet returns an array of nodes backed by machines owned by a given machineSet
func getNodesFromMachineSet(client runtimeclient.Client, machineSet mapiv1beta1.MachineSet) ([]*corev1.Node, error) {
	machines, err := getMachinesFromMachineSet(client, machineSet)
	if err != nil {
		return nil, fmt.Errorf("error calling getMachinesFromMachineSet %v", err)
	}

	var nodes []*corev1.Node
	for key := range machines {
		node, err := getNodeFromMachine(client, &machines[key])
		if err != nil {
			return nil, fmt.Errorf("error getting node from machine %q: %v", machines[key].Name, err)
		}
		nodes = append(nodes, node)
	}
	glog.Infof("MachineSet %q have %d nodes", machineSet.Name, len(nodes))
	return nodes, nil
}

// getNodeFromMachine returns the node object referenced by machine.Status.NodeRef
func getNodeFromMachine(client runtimeclient.Client, machine *mapiv1beta1.Machine) (*corev1.Node, error) {
	var node corev1.Node
	if machine.Status.NodeRef == nil {
		glog.Errorf("Machine %q has no NodeRef", machine.Name)
		return nil, fmt.Errorf("machine %q has no NodeRef", machine.Name)
	}
	key := runtimeclient.ObjectKey{Namespace: machine.Status.NodeRef.Namespace, Name: machine.Status.NodeRef.Name}
	if err := client.Get(context.Background(), key, &node); err != nil {
		return nil, fmt.Errorf("error getting node %q: %v", node.Name, err)
	}

	glog.Infof("Machine %q is backing node %q", machine.Name, node.Name)
	return &node, nil
}

// nodesAreReady returns true if an array of nodes are all ready
func nodesAreReady(nodes []*corev1.Node) bool {
	// All nodes needs to be ready
	for key := range nodes {
		if !e2e.IsNodeReady(nodes[key]) {
			glog.Errorf("Node %q is not ready. Conditions are: %v", nodes[key].Name, nodes[key].Status.Conditions)
			return false
		}
		glog.Infof("Node %q is ready. Conditions are: %v", nodes[key].Name, nodes[key].Status.Conditions)
	}
	return true
}

// scaleMachineSet scales a machineSet with a given name to the given number of replicas
func scaleMachineSet(name string, replicas int) error {
	scaleClient, err := getScaleClient()
	if err != nil {
		return fmt.Errorf("error calling getScaleClient %v", err)
	}

	scale, err := scaleClient.Scales(e2e.TestContext.MachineApiNamespace).Get(schema.GroupResource{Group: machineAPIGroup, Resource: "MachineSet"}, name)
	if err != nil {
		return fmt.Errorf("error calling scaleClient.Scales get: %v", err)
	}

	scaleUpdate := scale.DeepCopy()
	scaleUpdate.Spec.Replicas = int32(replicas)
	_, err = scaleClient.Scales(e2e.TestContext.MachineApiNamespace).Update(schema.GroupResource{Group: machineAPIGroup, Resource: "MachineSet"}, scaleUpdate)
	if err != nil {
		return fmt.Errorf("error calling scaleClient.Scales update: %v", err)
	}
	return nil
}

// getScaleClient returns a ScalesGetter object to manipulate scale subresources
func getScaleClient() (scale.ScalesGetter, error) {
	cfg, err := e2e.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("error getting config %v", err)
	}
	mapper, err := apiutil.NewDiscoveryRESTMapper(cfg)
	if err != nil {
		return nil, fmt.Errorf("error calling NewDiscoveryRESTMapper %v", err)
	}

	discovery := discovery.NewDiscoveryClientForConfigOrDie(cfg)
	scaleKindResolver := scale.NewDiscoveryScaleKindResolver(discovery)
	scaleClient, err := scale.NewForConfig(cfg, mapper, dynamic.LegacyAPIPathResolverFunc, scaleKindResolver)
	if err != nil {
		return nil, fmt.Errorf("error calling building scale client %v", err)
	}
	return scaleClient, nil
}

// nodesSnapShotLogs logs the state of all the nodes in the cluster
func nodesSnapShotLogs(client runtimeclient.Client) error {
	nodes, err := e2e.GetNodes(client)
	if err != nil {
		return fmt.Errorf("error getting nodes: %v", err)
	}

	for key, node := range nodes {
		glog.Infof("Node %q. Ready: %t. Unschedulable: %t", node.Name, e2e.IsNodeReady(&nodes[key]), node.Spec.Unschedulable)
	}
	return nil
}

func waitForClusterSizeToBeHealthy(client runtimeclient.Client, targetSize int) error {
	if err := wait.PollImmediate(5*time.Second, e2e.WaitLong, func() (bool, error) {
		glog.Infof("Cluster size expected to be %d nodes", targetSize)
		if err := machineSetsSnapShotLogs(client); err != nil {
			return false, err
		}

		if err := nodesSnapShotLogs(client); err != nil {
			return false, err
		}

		finalClusterSize, err := getClusterSize(client)
		if err != nil {
			return false, err
		}
		return finalClusterSize == targetSize, nil
	}); err != nil {
		return fmt.Errorf("Did not reach expected number of nodes: %v", err)
	}

	glog.Infof("waiting for all nodes to be ready")
	if err := e2e.WaitUntilAllNodesAreReady(client); err != nil {
		return err
	}

	glog.Infof("waiting for all nodes to be schedulable")
	if err := waitUntilAllNodesAreSchedulable(client); err != nil {
		return err
	}

	glog.Infof("waiting for each node to be backed by a machine")
	if !isOneMachinePerNode(client) {
		return fmt.Errorf("One machine per node condition violated")
	}

	return nil
}

// waitUntilAllNodesAreSchedulable waits for all cluster nodes to be schedulable and returns an error otherwise
func waitUntilAllNodesAreSchedulable(client runtimeclient.Client) error {
	return wait.PollImmediate(1*time.Second, time.Minute, func() (bool, error) {
		nodeList := corev1.NodeList{}
		if err := client.List(context.TODO(), &runtimeclient.ListOptions{}, &nodeList); err != nil {
			glog.Errorf("error querying api for nodeList object: %v, retrying...", err)
			return false, nil
		}
		// All nodes needs to be schedulable
		for _, node := range nodeList.Items {
			if node.Spec.Unschedulable == true {
				glog.Errorf("Node %q is unschedulable", node.Name)
				return false, nil
			}
			glog.Infof("Node %q is schedulable", node.Name)
		}
		return true, nil
	})
}

func machineFromMachineset(machineset *mapiv1beta1.MachineSet) *mapiv1beta1.Machine {
	randomUUID := string(uuid.NewUUID())

	machine := &mapiv1beta1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: machineset.Namespace,
			Name:      "machine-" + randomUUID[:6],
			Labels:    machineset.Labels,
		},
		Spec: machineset.Spec.Template.Spec,
	}
	if machine.Spec.ObjectMeta.Labels == nil {
		machine.Spec.ObjectMeta.Labels = map[string]string{}
	}
	for key := range nodeDrainLabels {
		if _, exists := machine.Spec.ObjectMeta.Labels[key]; exists {
			continue
		}
		machine.Spec.ObjectMeta.Labels[key] = nodeDrainLabels[key]
	}
	return machine
}

func waitUntilNodesAreReady(client runtimeclient.Client, listOpt *runtimeclient.ListOptions, nodeCount int) error {
	return wait.PollImmediate(e2e.RetryMedium, e2e.WaitLong, func() (bool, error) {
		nodes := corev1.NodeList{}
		if err := client.List(context.TODO(), listOpt, &nodes); err != nil {
			glog.Errorf("Error querying api for Node object: %v, retrying...", err)
			return false, nil
		}
		// expecting nodeGroupSize nodes
		readyNodes := 0
		for _, node := range nodes.Items {
			if _, exists := node.Labels[e2e.WorkerNodeRoleLabel]; !exists {
				continue
			}

			if !e2e.IsNodeReady(&node) {
				continue
			}

			readyNodes++
		}

		if readyNodes < nodeCount {
			glog.Errorf("Expecting %v nodes with %#v labels in Ready state, got %v", nodeCount, nodeDrainLabels, readyNodes)
			return false, nil
		}

		glog.Infof("Expected number (%v) of nodes with %v label in Ready state found", nodeCount, nodeDrainLabels)
		return true, nil
	})
}

func waitUntilNodesAreDeleted(client runtimeclient.Client, listOpt *runtimeclient.ListOptions) error {
	return wait.PollImmediate(e2e.RetryMedium, e2e.WaitLong, func() (bool, error) {
		nodes := corev1.NodeList{}
		if err := client.List(context.TODO(), listOpt, &nodes); err != nil {
			glog.Errorf("Error querying api for Node object: %v, retrying...", err)
			return false, nil
		}
		// expecting nodeGroupSize nodes
		nodeCounter := 0
		for _, node := range nodes.Items {
			if _, exists := node.Labels[e2e.WorkerNodeRoleLabel]; !exists {
				continue
			}

			if !e2e.IsNodeReady(&node) {
				continue
			}

			nodeCounter++
		}

		if nodeCounter > 0 {
			glog.Errorf("Expecting to found 0 nodes with %#v labels , got %v", nodeDrainLabels, nodeCounter)
			return false, nil
		}

		glog.Infof("Found 0 number of nodes with %v label as expected", nodeDrainLabels)
		return true, nil
	})
}

func waitUntilAllRCPodsAreReady(client runtimeclient.Client, rc *corev1.ReplicationController) error {
	return wait.PollImmediate(e2e.RetryMedium, e2e.WaitLong, func() (bool, error) {
		rcObj := corev1.ReplicationController{}
		key := types.NamespacedName{
			Namespace: rc.Namespace,
			Name:      rc.Name,
		}
		if err := client.Get(context.TODO(), key, &rcObj); err != nil {
			glog.Errorf("Error querying api RC %q object: %v, retrying...", rc.Name, err)
			return false, nil
		}
		if rcObj.Status.ReadyReplicas == 0 {
			glog.Infof("Waiting for at least one RC ready replica, ReadyReplicas: %v, Replicas: %v", rcObj.Status.ReadyReplicas, rcObj.Status.Replicas)
			return false, nil
		}
		glog.Infof("Waiting for RC ready replicas, ReadyReplicas: %v, Replicas: %v", rcObj.Status.ReadyReplicas, rcObj.Status.Replicas)
		return rcObj.Status.Replicas == rcObj.Status.ReadyReplicas, nil
	})
}

func verifyNodeDraining(client runtimeclient.Client, targetMachine *mapiv1beta1.Machine, rc *corev1.ReplicationController) (string, error) {
	var drainedNodeName string
	err := wait.PollImmediate(e2e.RetryMedium, e2e.WaitLong, func() (bool, error) {
		machine := mapiv1beta1.Machine{}

		key := types.NamespacedName{
			Namespace: targetMachine.Namespace,
			Name:      targetMachine.Name,
		}
		if err := client.Get(context.TODO(), key, &machine); err != nil {
			glog.Errorf("Error querying api machine %q object: %v, retrying...", targetMachine.Name, err)
			return false, nil
		}
		if machine.Status.NodeRef == nil || machine.Status.NodeRef.Kind != "Node" {
			glog.Errorf("Machine %q not linked to a node", machine.Name)
			return false, nil
		}

		drainedNodeName = machine.Status.NodeRef.Name
		node := corev1.Node{}

		if err := client.Get(context.TODO(), types.NamespacedName{Name: drainedNodeName}, &node); err != nil {
			glog.Errorf("Error querying api node %q object: %v, retrying...", drainedNodeName, err)
			return false, nil
		}

		if !node.Spec.Unschedulable {
			glog.Errorf("Node %q is expected to be marked as unschedulable, it is not", node.Name)
			return false, nil
		}

		glog.Infof("Node %q is mark unschedulable as expected", node.Name)

		pods := corev1.PodList{}
		listOpt := &runtimeclient.ListOptions{}
		listOpt.MatchingLabels(rc.Spec.Selector)
		if err := client.List(context.TODO(), listOpt, &pods); err != nil {
			glog.Errorf("Error querying api for Pods object: %v, retrying...", err)
			return false, nil
		}

		podCounter := 0
		for _, pod := range pods.Items {
			if pod.Spec.NodeName != machine.Status.NodeRef.Name {
				continue
			}
			if !pod.DeletionTimestamp.IsZero() {
				continue
			}
			podCounter++
		}

		glog.Infof("Have %v pods scheduled to node %q", podCounter, machine.Status.NodeRef.Name)

		// Verify we have enough pods running as well
		rcObj := corev1.ReplicationController{}
		key = types.NamespacedName{
			Namespace: rc.Namespace,
			Name:      rc.Name,
		}
		if err := client.Get(context.TODO(), key, &rcObj); err != nil {
			glog.Errorf("Error querying api RC %q object: %v, retrying...", rc.Name, err)
			return false, nil
		}

		// The point of the test is to make sure majority of the pods are rescheduled
		// to other nodes. Pod disruption budget makes sure at most one pod
		// owned by the RC is not Ready. So no need to test it. Though, useful to have it printed.
		glog.Infof("RC ReadyReplicas: %v, Replicas: %v", rcObj.Status.ReadyReplicas, rcObj.Status.Replicas)

		// This makes sure at most one replica is not ready
		if rcObj.Status.Replicas-rcObj.Status.ReadyReplicas > 1 {
			return false, fmt.Errorf("pod disruption budget not respected, node was not properly drained")
		}

		// Depends on timing though a machine can be deleted even before there is only
		// one pod left on the node (that is being evicted).
		if podCounter > 2 {
			glog.Infof("Expecting at most 2 pods to be scheduled to drained node %q, got %v", machine.Status.NodeRef.Name, podCounter)
			return false, nil
		}

		glog.Info("Expected result: all pods from the RC up to last one or two got scheduled to a different node while respecting PDB")
		return true, nil
	})

	return drainedNodeName, err
}

func waitUntilNodeDoesNotExists(client runtimeclient.Client, nodeName string) error {
	return wait.PollImmediate(e2e.RetryMedium, e2e.WaitLong, func() (bool, error) {
		node := corev1.Node{}

		key := types.NamespacedName{
			Name: nodeName,
		}
		err := client.Get(context.TODO(), key, &node)
		if err == nil {
			glog.Errorf("Node %q not yet deleted", nodeName)
			return false, nil
		}

		if !strings.Contains(err.Error(), "not found") {
			glog.Errorf("Error querying api node %q object: %v, retrying...", nodeName, err)
			return false, nil
		}

		glog.Infof("Node %q successfully deleted", nodeName)
		return true, nil
	})
}
