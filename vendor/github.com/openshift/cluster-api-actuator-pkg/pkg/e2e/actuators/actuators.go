package actuators

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang/glog"
	g "github.com/onsi/ginkgo"
	o "github.com/onsi/gomega"
	e2e "github.com/openshift/cluster-api-actuator-pkg/pkg/e2e/framework"
	mapiv1beta1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	corev1 "k8s.io/api/core/v1"
	kpolicyapi "k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/apimachinery/pkg/util/wait"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

var nodeDrainLabels = map[string]string{
	e2e.WorkerRoleLabel:  "",
	"node-draining-test": string(uuid.NewUUID()),
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
			if _, exists := node.Labels[e2e.WorkerRoleLabel]; !exists {
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
			if _, exists := node.Labels[e2e.WorkerRoleLabel]; !exists {
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

func replicationControllerWorkload(namespace string) *corev1.ReplicationController {
	var replicas int32 = 20
	return &corev1.ReplicationController{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pdb-workload",
			Namespace: namespace,
		},
		Spec: corev1.ReplicationControllerSpec{
			Replicas: &replicas,
			Selector: map[string]string{
				"app": "nginx",
			},
			Template: &corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: "nginx",
					Labels: map[string]string{
						"app": "nginx",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "work",
							Image:   "busybox",
							Command: []string{"sleep", "10h"},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									"cpu":    resource.MustParse("50m"),
									"memory": resource.MustParse("50Mi"),
								},
							},
						},
					},
					NodeSelector: nodeDrainLabels,
					Tolerations: []corev1.Toleration{
						{
							Key:      "kubemark",
							Operator: corev1.TolerationOpExists,
						},
					},
				},
			},
		},
	}
}

func podDisruptionBudget(namespace string) *kpolicyapi.PodDisruptionBudget {
	maxUnavailable := intstr.FromInt(1)
	return &kpolicyapi.PodDisruptionBudget{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx-pdb",
			Namespace: namespace,
		},
		Spec: kpolicyapi.PodDisruptionBudgetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx",
				},
			},
			MaxUnavailable: &maxUnavailable,
		},
	}
}

var _ = g.Describe("[Feature:Machines] Actuator should", func() {
	defer g.GinkgoRecover()

	g.It("drain node before removing machine resource", func() {
		var err error
		client, err := e2e.LoadClient()
		o.Expect(err).NotTo(o.HaveOccurred())

		delObjects := make(map[string]runtime.Object)

		defer func() {
			// Remove resources
			for key := range delObjects {
				glog.Infof("Deleting object %q", key)
				if err := client.Delete(context.TODO(), delObjects[key]); err != nil {
					glog.Errorf("Unable to delete object %q: %v", key, err)
				}
			}

			listOpt := &runtimeclient.ListOptions{}
			listOpt.MatchingLabels(nodeDrainLabels)
			// TODO(jchaloup): we need to make sure this gets called no matter what
			// and waits until all labeled nodes are gone. Though, it it does not
			// happend in the timeout set, it will not happen ever.
			err := waitUntilNodesAreDeleted(client, listOpt)
			o.Expect(err).NotTo(o.HaveOccurred())
		}()

		g.By("Taking the first worker machineset (assuming only worker machines are backed by machinesets)")
		machinesets := mapiv1beta1.MachineSetList{}
		err = wait.PollImmediate(e2e.RetryMedium, e2e.WaitShort, func() (bool, error) {
			if err := client.List(context.TODO(), &runtimeclient.ListOptions{}, &machinesets); err != nil {
				glog.Errorf("Error querying api for machineset object: %v, retrying...", err)
				return false, nil
			}
			if len(machinesets.Items) < 1 {
				glog.Errorf("Expected at least one machineset, have none")
				return false, nil
			}
			return true, nil
		})
		o.Expect(err).NotTo(o.HaveOccurred())

		g.By("Creating two new machines, one for node about to be drained, other for moving workload from drained node")
		// Create two machines
		machine1 := machineFromMachineset(&machinesets.Items[0])
		machine1.Name = "machine1"

		err = func() error {
			if err := client.Create(context.TODO(), machine1); err != nil {
				return fmt.Errorf("unable to create machine %q: %v", machine1.Name, err)
			}
			delObjects["machine1"] = machine1

			machine2 := machineFromMachineset(&machinesets.Items[0])
			machine2.Name = "machine2"

			if err := client.Create(context.TODO(), machine2); err != nil {
				return fmt.Errorf("unable to create machine %q: %v", machine2.Name, err)
			}
			delObjects["machine2"] = machine2

			return nil
		}()
		o.Expect(err).NotTo(o.HaveOccurred())

		g.By("Waiting until both new nodes are ready")
		listOpt := &runtimeclient.ListOptions{}
		listOpt.MatchingLabels(nodeDrainLabels)
		err = waitUntilNodesAreReady(client, listOpt, 2)
		o.Expect(err).NotTo(o.HaveOccurred())

		g.By("Creating RC with workload")
		rc := replicationControllerWorkload("default")
		err = client.Create(context.TODO(), rc)
		o.Expect(err).NotTo(o.HaveOccurred())
		delObjects["rc"] = rc

		g.By("Creating PDB for RC")
		pdb := podDisruptionBudget("default")
		err = client.Create(context.TODO(), pdb)
		o.Expect(err).NotTo(o.HaveOccurred())
		delObjects["pdb"] = pdb

		g.By("Wait until all replicas are ready")
		err = wait.PollImmediate(e2e.RetryMedium, e2e.WaitLong, func() (bool, error) {
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
				glog.Infof("Waiting for at least one RC ready replica (%v/%v)", rcObj.Status.ReadyReplicas, rcObj.Status.Replicas)
				return false, nil
			}
			glog.Infof("Waiting for RC ready replicas (%v/%v)", rcObj.Status.ReadyReplicas, rcObj.Status.Replicas)
			if rcObj.Status.Replicas != rcObj.Status.ReadyReplicas {
				return false, nil
			}
			return true, nil
		})
		o.Expect(err).NotTo(o.HaveOccurred())

		// TODO(jchaloup): delete machine that has at least half of the RC pods

		// All pods are distributed evenly among all nodes so it's fine to drain
		// random node and observe reconciliation of pods on the other one.
		g.By("Delete machine to trigger node draining")
		err = client.Delete(context.TODO(), machine1)
		o.Expect(err).NotTo(o.HaveOccurred())
		delete(delObjects, "machine1")

		// We still should be able to list the machine as until rc.replicas-1 are running on the other node
		g.By("Observing and verifying node draining")
		var drainedNodeName string
		err = wait.PollImmediate(e2e.RetryMedium, e2e.WaitLong, func() (bool, error) {
			machine := mapiv1beta1.Machine{}

			key := types.NamespacedName{
				Namespace: machine1.Namespace,
				Name:      machine1.Name,
			}
			if err := client.Get(context.TODO(), key, &machine); err != nil {
				glog.Errorf("Error querying api machine %q object: %v, retrying...", machine1.Name, err)
				return false, nil
			}
			if machine.Status.NodeRef == nil || machine.Status.NodeRef.Kind != "Node" {
				glog.Error("Machine %q not linked to a node", machine.Name)
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

			// expecting nodeGroupSize nodes
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

			// The point of the test is to make sure majority of the pods is rescheduled
			// to other nodes. Pod disruption budget makes sure at most one pod
			// owned by the RC is not Ready. So no need to test it. Though, usefull to have it printed.
			glog.Infof("RC ReadyReplicas/Replicas: %v/%v", rcObj.Status.ReadyReplicas, rcObj.Status.Replicas)

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
		o.Expect(err).NotTo(o.HaveOccurred())

		g.By("Validating the machine is deleted")
		err = wait.PollImmediate(e2e.RetryMedium, e2e.WaitShort, func() (bool, error) {
			machine := mapiv1beta1.Machine{}

			key := types.NamespacedName{
				Namespace: machine1.Namespace,
				Name:      machine1.Name,
			}
			err := client.Get(context.TODO(), key, &machine)
			if err == nil {
				glog.Errorf("Machine %q not yet deleted", machine1.Name)
				return false, nil
			}

			if !strings.Contains(err.Error(), "not found") {
				glog.Errorf("Error querying api machine %q object: %v, retrying...", machine1.Name, err)
				return false, nil
			}

			glog.Infof("Machine %q successfully deleted", machine1.Name)
			return true, nil
		})
		o.Expect(err).NotTo(o.HaveOccurred())

		// Validate underlying node is removed as well
		err = wait.PollImmediate(e2e.RetryMedium, e2e.WaitLong, func() (bool, error) {
			node := corev1.Node{}

			key := types.NamespacedName{
				Name: drainedNodeName,
			}
			err := client.Get(context.TODO(), key, &node)
			if err == nil {
				glog.Errorf("Node %q not yet deleted", drainedNodeName)
				return false, nil
			}

			if !strings.Contains(err.Error(), "not found") {
				glog.Errorf("Error querying api node %q object: %v, retrying...", drainedNodeName, err)
				return false, nil
			}

			glog.Infof("Node %q successfully deleted", drainedNodeName)
			return true, nil
		})
		o.Expect(err).NotTo(o.HaveOccurred())

	})

})
