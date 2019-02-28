package infra

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang/glog"
	g "github.com/onsi/ginkgo"
	o "github.com/onsi/gomega"
	e2e "github.com/openshift/cluster-api-actuator-pkg/pkg/e2e/framework"
	mapiv1beta1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	controllernode "github.com/openshift/cluster-api/pkg/controller/node"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/apimachinery/pkg/util/wait"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = g.Describe("[Feature:Machines] Machines should", func() {
	defer g.GinkgoRecover()

	g.It("be linked with nodes", func() {
		var err error
		client, err := e2e.LoadClient()
		o.Expect(err).NotTo(o.HaveOccurred())

		listOptions := runtimeclient.ListOptions{
			Namespace: e2e.TestContext.MachineApiNamespace,
		}
		machineList := mapiv1beta1.MachineList{}
		nodeList := corev1.NodeList{}

		err = wait.PollImmediate(5*time.Second, e2e.WaitShort, func() (bool, error) {
			if err := client.List(context.TODO(), &listOptions, &machineList); err != nil {
				glog.Errorf("error querying api for machineList object: %v, retrying...", err)
				return false, nil
			}
			if err := client.List(context.TODO(), &listOptions, &nodeList); err != nil {
				glog.Errorf("error querying api for nodeList object: %v, retrying...", err)
				return false, nil
			}

			// Every machine needs to be linked to a node, though some nodes do not have to linked to any machines
			glog.Infof("Expecting at least the same number of machines as nodes, have %v nodes and %v machines", len(nodeList.Items), len(machineList.Items))
			if len(machineList.Items) > len(nodeList.Items) {
				return false, nil
			}

			nodeNameToMachineAnnotation := make(map[string]string)
			for _, node := range nodeList.Items {
				nodeNameToMachineAnnotation[node.Name] = node.Annotations[controllernode.MachineAnnotationKey]
			}
			for _, machine := range machineList.Items {
				if machine.Status.NodeRef == nil {
					glog.Errorf("machine %s has no NodeRef, retrying...", machine.Name)
					return false, nil
				}
				nodeName := machine.Status.NodeRef.Name
				if nodeNameToMachineAnnotation[nodeName] != fmt.Sprintf("%s/%s", e2e.TestContext.MachineApiNamespace, machine.Name) {
					glog.Errorf("node name %s does not match expected machine name %s, retrying...", nodeName, machine.Name)
					return false, nil
				}
				glog.Infof("Machine %q is linked to node %q", machine.Name, nodeName)
			}
			return true, nil
		})
		o.Expect(err).NotTo(o.HaveOccurred())
	})

	g.It("have ability to additively reconcile taints", func() {
		var err error
		client, err := e2e.LoadClient()
		o.Expect(err).NotTo(o.HaveOccurred())

		g.By("Verify machine taints are getting applied to node")
		err = func() error {
			listOptions := runtimeclient.ListOptions{
				Namespace: e2e.TestContext.MachineApiNamespace,
			}
			machineList := mapiv1beta1.MachineList{}

			if err := client.List(context.TODO(), &listOptions, &machineList); err != nil {
				return fmt.Errorf("error querying api for machineList object: %v", err)
			}
			g.By("Got the machine list")
			machine := machineList.Items[0]
			if machine.Status.NodeRef == nil {
				return fmt.Errorf("machine %s has no NodeRef", machine.Name)
			}
			g.By(fmt.Sprintf("Got the machine %s", machine.Name))
			nodeName := machine.Status.NodeRef.Name
			nodeKey := types.NamespacedName{
				Namespace: e2e.TestContext.MachineApiNamespace,
				Name:      nodeName,
			}
			node := &corev1.Node{}

			if err := client.Get(context.TODO(), nodeKey, node); err != nil {
				return fmt.Errorf("error querying api for node object: %v", err)
			}
			g.By(fmt.Sprintf("Got the node %s from machine, %s", node.Name, machine.Name))
			nodeTaint := corev1.Taint{
				Key:    "not-from-machine",
				Value:  "true",
				Effect: corev1.TaintEffectNoSchedule,
			}
			// Do not remove any taint, just extend the list
			// The test removes the nodes anyway, so the list will not grow over time much
			node.Spec.Taints = append(node.Spec.Taints, nodeTaint)
			if err := client.Update(context.TODO(), node); err != nil {
				return fmt.Errorf("error updating node object with non-machine taint: %v", err)
			}
			g.By("Updated node object with taint")
			machineTaint := corev1.Taint{
				Key:    fmt.Sprintf("from-machine-%v", string(uuid.NewUUID())),
				Value:  "true",
				Effect: corev1.TaintEffectNoSchedule,
			}

			// Do not remove any taint, just extend the list
			// The test removes the machine anyway, so the list will not grow over time much
			machine.Spec.Taints = append(machine.Spec.Taints, machineTaint)
			if err := client.Update(context.TODO(), &machine); err != nil {
				return fmt.Errorf("error updating machine object with taint: %v", err)
			}

			g.By("Updated machine object with taint")
			var expectedTaints = sets.NewString("not-from-machine", machineTaint.Key)
			err := wait.PollImmediate(1*time.Second, e2e.WaitLong, func() (bool, error) {
				if err := client.Get(context.TODO(), nodeKey, node); err != nil {
					glog.Errorf("error querying api for node object: %v", err)
					return false, nil
				}
				glog.Info("Got the node again for verification of taints")
				var observedTaints = sets.NewString()
				for _, taint := range node.Spec.Taints {
					observedTaints.Insert(taint.Key)
				}
				if expectedTaints.Difference(observedTaints).HasAny("not-from-machine", machineTaint.Key) == false {
					glog.Infof("expected : %v, observed %v , difference %v, ", expectedTaints, observedTaints, expectedTaints.Difference(observedTaints))
					return true, nil
				}
				glog.Infof("Did not find all expected taints on the node. Missing: %v", expectedTaints.Difference(observedTaints))
				return false, nil
			})
			return err
		}()
		o.Expect(err).NotTo(o.HaveOccurred())
	})

	g.It("be replaced with new one when deleted", func() {
		var err error
		client, err := e2e.LoadClient()
		o.Expect(err).NotTo(o.HaveOccurred())

		// Assume cluster state
		err = e2e.WaitUntilAllNodesAreReady(client)
		o.Expect(err).NotTo(o.HaveOccurred())

		glog.Info("Get machineList")
		machineList := mapiv1beta1.MachineList{}
		err = wait.PollImmediate(1*time.Second, e2e.WaitShort, func() (bool, error) {
			if err := client.List(context.TODO(), runtimeclient.InNamespace(e2e.TestContext.MachineApiNamespace), &machineList); err != nil {
				glog.Errorf("error querying api for machineList object: %v, retrying...", err)
				return false, nil
			}
			return true, nil
		})
		o.Expect(err).NotTo(o.HaveOccurred())

		glog.Info("Get nodeList")
		nodeList := corev1.NodeList{}
		err = wait.PollImmediate(1*time.Second, e2e.WaitShort, func() (bool, error) {
			if err := client.List(context.TODO(), nil, &nodeList); err != nil {
				glog.Errorf("error querying api for nodeList object: %v, retrying...", err)
				return false, nil
			}
			return true, nil
		})
		o.Expect(err).NotTo(o.HaveOccurred())

		clusterInitialTotalNodes := len(nodeList.Items)
		clusterInitialTotalMachines := len(machineList.Items)
		var triagedWorkerMachine *mapiv1beta1.Machine
		var triagedWorkerNode *corev1.Node

		glog.Infof("Initial number of nodes: %v, initial number of machines: %v", clusterInitialTotalNodes, clusterInitialTotalMachines)

		err = func() error {
			for _, machine := range machineList.Items {
				if machine.Status.NodeRef == nil {
					glog.Infof("Machine %q is not linked to any node", machine.Name)
					continue
				}
				node := corev1.Node{}
				if err := client.Get(context.TODO(), types.NamespacedName{Name: machine.Status.NodeRef.Name}, &node); err != nil {
					glog.Errorf("error querying api for node object: %v, retrying...", err)
					continue
				}

				if node.Labels == nil {
					continue
				}

				if _, exists := node.Labels["node-role.kubernetes.io/worker"]; !exists {
					continue
				}

				triagedWorkerNode = node.DeepCopy()
				triagedWorkerMachine = machine.DeepCopy()
				break
			}
			if triagedWorkerMachine == nil {
				return errors.New("Unable to find any worker machine")
			}
			return nil
		}()
		o.Expect(err).NotTo(o.HaveOccurred())

		glog.Infof("Delete machine %v", triagedWorkerMachine.Name)
		err = wait.PollImmediate(1*time.Second, e2e.WaitShort, func() (bool, error) {
			if err := client.Delete(context.TODO(), triagedWorkerMachine); err != nil {
				glog.Errorf("error querying api for Deployment object: %v, retrying...", err)
				return false, nil
			}
			return true, nil
		})
		o.Expect(err).NotTo(o.HaveOccurred())

		err = wait.PollImmediate(5*time.Second, e2e.WaitMedium, func() (bool, error) {
			if err := client.List(context.TODO(), runtimeclient.InNamespace(e2e.TestContext.MachineApiNamespace), &machineList); err != nil {
				glog.Errorf("error querying api for machineList object: %v, retrying...", err)
				return false, nil
			}
			glog.Info("Expect new machine to come up")
			return len(machineList.Items) == clusterInitialTotalMachines, nil
		})
		o.Expect(err).NotTo(o.HaveOccurred())

		err = wait.PollImmediate(5*time.Second, e2e.WaitLong, func() (bool, error) {
			if err := client.List(context.TODO(), nil, &nodeList); err != nil {
				glog.Errorf("error querying api for nodeList object: %v, retrying...", err)
				return false, nil
			}
			glog.Info("Expect deleted machine node to go away")
			for _, n := range nodeList.Items {
				if n.Name == triagedWorkerNode.Name {
					return false, nil
				}
			}
			return true, nil
		})
		o.Expect(err).NotTo(o.HaveOccurred())

		err = wait.PollImmediate(5*time.Second, e2e.WaitLong, func() (bool, error) {
			if err := client.List(context.TODO(), nil, &nodeList); err != nil {
				glog.Errorf("error querying api for nodeList object: %v, retrying...", err)
				return false, nil
			}
			glog.Info("Expect new node to come up")
			return len(nodeList.Items) == clusterInitialTotalNodes, nil
		})
		o.Expect(err).NotTo(o.HaveOccurred())
	})

})
