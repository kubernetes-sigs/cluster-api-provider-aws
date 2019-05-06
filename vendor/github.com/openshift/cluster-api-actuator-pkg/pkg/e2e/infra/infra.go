package infra

import (
	"context"
	"fmt"
	"strings"
	"time"

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
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/utils/pointer"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

var nodeDrainLabels = map[string]string{
	e2e.WorkerNodeRoleLabel: "",
	"node-draining-test":    string(uuid.NewUUID()),
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

func invalidMachinesetWithEmptyProviderConfig() *mapiv1beta1.MachineSet {
	var oneReplicas int32 = 1
	return &mapiv1beta1.MachineSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "invalid-machineset",
			Namespace: e2e.TestContext.MachineApiNamespace,
		},
		Spec: mapiv1beta1.MachineSetSpec{
			Replicas: &oneReplicas,
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"little-kitty": "i-am-little-kitty",
				},
			},
			Template: mapiv1beta1.MachineTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"big-kitty": "i-am-bit-kitty",
					},
				},
				Spec: mapiv1beta1.MachineSpec{
					// Empty providerSpec!!! we don't want to provision real instances.
					// Just to observe how many machine replicas get created.
					ProviderSpec: mapiv1beta1.ProviderSpec{},
				},
			},
		},
	}
}

var _ = g.Describe("[Feature:Machines] Managed cluster should", func() {
	defer g.GinkgoRecover()

	g.It("have machines linked with nodes", func() {
		var err error
		client, err := e2e.LoadClient()
		o.Expect(err).NotTo(o.HaveOccurred())

		o.Expect(isOneMachinePerNode(client)).To(o.BeTrue())
		o.Expect(err).NotTo(o.HaveOccurred())
	})

	g.It("have ability to additively reconcile taints from machine to nodes", func() {
		var err error
		client, err := e2e.LoadClient()
		o.Expect(err).NotTo(o.HaveOccurred())

		machines, err := e2e.GetMachines(context.TODO(), client)
		o.Expect(err).NotTo(o.HaveOccurred())
		o.Expect(len(machines)).To(o.BeNumerically(">", 0))
		machine := &machines[0]
		originalMachineTaints := machine.Spec.Taints
		g.By(fmt.Sprintf("getting machine %q", machine.Name))

		node, err := getNodeFromMachine(client, machine)
		o.Expect(err).NotTo(o.HaveOccurred())
		originalNodeTaints := node.Spec.Taints
		g.By(fmt.Sprintf("getting the backed node %q", node.Name))

		nodeTaint := corev1.Taint{
			Key:    "not-from-machine",
			Value:  "true",
			Effect: corev1.TaintEffectNoSchedule,
		}
		g.By(fmt.Sprintf("updating node %q with taint: %v", node.Name, nodeTaint))
		node.Spec.Taints = append(node.Spec.Taints, nodeTaint)
		err = client.Update(context.TODO(), node)
		o.Expect(err).NotTo(o.HaveOccurred())

		machineTaint := corev1.Taint{
			Key:    fmt.Sprintf("from-machine-%v", string(uuid.NewUUID())),
			Value:  "true",
			Effect: corev1.TaintEffectNoSchedule,
		}
		g.By(fmt.Sprintf("updating machine %q with taint: %v", machine.Name, machineTaint))
		machine.Spec.Taints = append(machine.Spec.Taints, machineTaint)
		err = client.Update(context.TODO(), machine)
		o.Expect(err).NotTo(o.HaveOccurred())

		var expectedTaints = sets.NewString("not-from-machine", machineTaint.Key)
		o.Eventually(func() bool {
			glog.Info("Getting node from machine again for verification of taints")
			node, err := getNodeFromMachine(client, machine)
			if err != nil {
				return false
			}
			var observedTaints = sets.NewString()
			for _, taint := range node.Spec.Taints {
				observedTaints.Insert(taint.Key)
			}
			if expectedTaints.Difference(observedTaints).HasAny("not-from-machine", machineTaint.Key) == false {
				glog.Infof("Expected : %v, observed %v , difference %v, ", expectedTaints, observedTaints, expectedTaints.Difference(observedTaints))
				return true
			}
			glog.Infof("Did not find all expected taints on the node. Missing: %v", expectedTaints.Difference(observedTaints))
			return false
		}, e2e.WaitMedium, 5*time.Second).Should(o.BeTrue())

		g.By("Getting the latest version of the original machine")
		machine, err = e2e.GetMachine(context.TODO(), client, machine.Name)
		o.Expect(err).NotTo(o.HaveOccurred())

		g.By("Setting back the original machine taints")
		machine.Spec.Taints = originalMachineTaints
		err = client.Update(context.TODO(), machine)
		o.Expect(err).NotTo(o.HaveOccurred())

		g.By("Getting the latest version of the node")
		node, err = getNodeFromMachine(client, machine)
		o.Expect(err).NotTo(o.HaveOccurred())

		g.By("Setting back the original node taints")
		node.Spec.Taints = originalNodeTaints
		err = client.Update(context.TODO(), node)
		o.Expect(err).NotTo(o.HaveOccurred())
	})

	g.It("recover from deleted worker machines", func() {
		var err error
		client, err := e2e.LoadClient()
		o.Expect(err).NotTo(o.HaveOccurred())

		// Expect for cluster to cool down from previous tests
		err = e2e.WaitUntilAllNodesAreReady(client)
		o.Expect(err).NotTo(o.HaveOccurred())

		g.By("checking initial cluster state")
		initialClusterSize, err := getClusterSize(client)
		err = waitForClusterSizeToBeHealthy(client, initialClusterSize)
		o.Expect(err).NotTo(o.HaveOccurred())

		g.By("getting worker node")
		workerNodes, err := e2e.GetWorkerNodes(client)
		o.Expect(err).NotTo(o.HaveOccurred())
		o.Expect(workerNodes).ToNot(o.BeEmpty())

		workerNode := &workerNodes[0]
		workerMachine, err := e2e.GetMachineFromNode(client, workerNode)
		o.Expect(err).NotTo(o.HaveOccurred())
		g.By(fmt.Sprintf("deleting machine object %q", workerMachine.Name))
		err = deleteMachine(client, workerMachine)
		o.Expect(err).NotTo(o.HaveOccurred())

		g.By(fmt.Sprintf("waiting for node object %q to go away", workerNode.Name))
		nodeList := corev1.NodeList{}
		o.Eventually(func() bool {
			if err := client.List(context.TODO(), nil, &nodeList); err != nil {
				glog.Errorf("Error querying api for nodeList object: %v, retrying...", err)
				return false
			}
			for _, n := range nodeList.Items {
				if n.Name == workerNode.Name {
					glog.Infof("Node %q still exists. Node conditions are: %v", workerNode.Name, workerNode.Status.Conditions)
					return false
				}
			}
			return true
		}, e2e.WaitLong, 5*time.Second).Should(o.BeTrue())

		g.By(fmt.Sprintf("waiting for new node object to come up"))
		err = waitForClusterSizeToBeHealthy(client, initialClusterSize)
		o.Expect(err).NotTo(o.HaveOccurred())
	})

	g.It("grow or decrease when scaling out or in", func() {
		g.By("checking initial cluster state")
		client, err := e2e.LoadClient()
		o.Expect(err).NotTo(o.HaveOccurred())

		initialClusterSize, err := getClusterSize(client)
		err = waitForClusterSizeToBeHealthy(client, initialClusterSize)
		o.Expect(err).NotTo(o.HaveOccurred())

		machineSets, err := e2e.GetMachineSets(context.TODO(), client)
		o.Expect(err).NotTo(o.HaveOccurred())
		o.Expect(len(machineSets)).To(o.BeNumerically(">=", 2))
		machineSet := machineSets[0]
		initialReplicasMachineSet := int(pointer.Int32PtrDerefOr(machineSet.Spec.Replicas, e2e.DefaultMachineSetReplicas))
		scaleOut := 3
		scaleIn := initialReplicasMachineSet
		originalReplicas := initialReplicasMachineSet
		clusterGrowth := scaleOut - originalReplicas
		clusterDecrease := scaleOut - scaleIn
		intermediateClusterSize := initialClusterSize + clusterGrowth
		finalClusterSize := initialClusterSize + clusterGrowth - clusterDecrease

		g.By(fmt.Sprintf("scaling out %q machineSet to %d replicas", machineSet.Name, scaleOut))
		err = scaleMachineSet(machineSet.Name, scaleOut)
		o.Expect(err).NotTo(o.HaveOccurred())

		g.By(fmt.Sprintf("waiting for cluster to grow %d nodes. Size should be %d", clusterGrowth, intermediateClusterSize))
		err = waitForClusterSizeToBeHealthy(client, intermediateClusterSize)
		o.Expect(err).NotTo(o.HaveOccurred())

		g.By(fmt.Sprintf("scaling in %q machineSet to %d replicas", machineSet.Name, scaleIn))
		err = scaleMachineSet(machineSet.Name, scaleIn)
		o.Expect(err).NotTo(o.HaveOccurred())

		g.By(fmt.Sprintf("waiting for cluster to decrease %d nodes. Final size should be %d nodes", clusterDecrease, finalClusterSize))
		err = waitForClusterSizeToBeHealthy(client, finalClusterSize)
		o.Expect(err).NotTo(o.HaveOccurred())
	})

	g.It("grow and decrease when scaling different machineSets simultaneously", func() {
		client, err := e2e.LoadClient()
		o.Expect(err).NotTo(o.HaveOccurred())
		scaleOut := 3

		g.By("checking initial cluster size")
		initialClusterSize, err := getClusterSize(client)
		o.Expect(err).NotTo(o.HaveOccurred())

		g.By("getting worker machineSets")
		machineSets, err := e2e.GetMachineSets(context.TODO(), client)
		o.Expect(err).NotTo(o.HaveOccurred())
		o.Expect(len(machineSets)).To(o.BeNumerically(">=", 2))
		machineSet0 := machineSets[0]
		initialReplicasMachineSet0 := int(pointer.Int32PtrDerefOr(machineSet0.Spec.Replicas, e2e.DefaultMachineSetReplicas))
		machineSet1 := machineSets[1]
		initialReplicasMachineSet1 := int(pointer.Int32PtrDerefOr(machineSet1.Spec.Replicas, e2e.DefaultMachineSetReplicas))

		g.By(fmt.Sprintf("scaling %q from %d to %d replicas", machineSet0.Name, initialReplicasMachineSet0, scaleOut))
		err = scaleMachineSet(machineSet0.Name, scaleOut)
		o.Expect(err).NotTo(o.HaveOccurred())

		g.By(fmt.Sprintf("scaling %q from %d to %d replicas", machineSet1.Name, initialReplicasMachineSet1, scaleOut))
		err = scaleMachineSet(machineSet1.Name, scaleOut)
		o.Expect(err).NotTo(o.HaveOccurred())

		o.Eventually(func() bool {
			nodes, err := getNodesFromMachineSet(client, machineSet0)
			if err != nil {
				return false
			}
			return len(nodes) == scaleOut && nodesAreReady(nodes)
		}, e2e.WaitLong, 5*time.Second).Should(o.BeTrue())

		o.Eventually(func() bool {
			nodes, err := getNodesFromMachineSet(client, machineSet1)
			if err != nil {
				return false
			}
			return len(nodes) == scaleOut && nodesAreReady(nodes)
		}, e2e.WaitLong, 5*time.Second).Should(o.BeTrue())

		g.By(fmt.Sprintf("scaling %q from %d to %d replicas", machineSet0.Name, scaleOut, initialReplicasMachineSet0))
		err = scaleMachineSet(machineSet0.Name, initialReplicasMachineSet0)
		o.Expect(err).NotTo(o.HaveOccurred())

		g.By(fmt.Sprintf("scaling %q from %d to %d replicas", machineSet1.Name, scaleOut, initialReplicasMachineSet1))
		err = scaleMachineSet(machineSet1.Name, initialReplicasMachineSet1)
		o.Expect(err).NotTo(o.HaveOccurred())

		g.By(fmt.Sprintf("waiting for cluster to get back to original size. Final size should be %d nodes", initialClusterSize))
		err = waitForClusterSizeToBeHealthy(client, initialClusterSize)
		o.Expect(err).NotTo(o.HaveOccurred())
	})

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
		err = waitUntilAllRCPodsAreReady(client, rc)
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
		drainedNodeName, err := verifyNodeDraining(client, machine1, rc)
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

		g.By("Validate underlying node is removed as well")
		err = waitUntilNodeDoesNotExists(client, drainedNodeName)
		o.Expect(err).NotTo(o.HaveOccurred())

	})

	g.It("reject invalid machinesets", func() {
		var err error
		client, err := e2e.LoadClient()
		o.Expect(err).NotTo(o.HaveOccurred())

		g.By("Creating invalid machineset")
		invalidMachineSet := invalidMachinesetWithEmptyProviderConfig()

		err = client.Create(context.TODO(), invalidMachineSet)
		o.Expect(err).NotTo(o.HaveOccurred())

		g.By("Waiting for ReconcileError MachineSet event")
		err = wait.PollImmediate(e2e.RetryMedium, e2e.WaitShort, func() (bool, error) {
			eventList := corev1.EventList{}
			if err := client.List(context.TODO(), nil, &eventList); err != nil {
				glog.Errorf("error querying api for eventList object: %v, retrying...", err)
				return false, nil
			}

			glog.Infof("Fetching ReconcileError MachineSet invalid-machineset event")
			for _, event := range eventList.Items {
				if event.Reason != "ReconcileError" || event.InvolvedObject.Kind != "MachineSet" || event.InvolvedObject.Name != invalidMachineSet.Name {
					continue
				}

				glog.Infof("Found ReconcileError event for %q machine set with the following message: %v", event.InvolvedObject.Name, event.Message)
				return true, nil
			}

			return false, nil
		})
		o.Expect(err).NotTo(o.HaveOccurred())

		// Verify the number of machines does not grow over time.
		// The assumption is once the ReconcileError event is recorded and caught,
		// the machineset is not reconciled again until it's updated.
		machineList := &mapiv1beta1.MachineList{}
		err = client.List(context.TODO(), runtimeclient.MatchingLabels(invalidMachineSet.Spec.Template.Labels), machineList)
		o.Expect(err).NotTo(o.HaveOccurred())

		g.By(fmt.Sprintf("Verify no machine from %q machineset were created", invalidMachineSet.Name))
		glog.Infof("Have %v machines generated from %q machineset", len(machineList.Items), invalidMachineSet.Name)
		o.Expect(len(machineList.Items)).To(o.BeNumerically("==", 0))

		g.By("Deleting invalid machineset")
		err = client.Delete(context.TODO(), invalidMachineSet)
		o.Expect(err).NotTo(o.HaveOccurred())
	})
})
