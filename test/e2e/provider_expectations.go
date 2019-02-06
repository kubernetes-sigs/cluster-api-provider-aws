package main

import (
	"time"

	"context"
	"errors"
	"fmt"

	"github.com/golang/glog"
	mapiv1beta1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	kappsapi "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	waitShort  = 1 * time.Minute
	waitMedium = 3 * time.Minute
	waitLong   = 10 * time.Minute
)

func (tc *testConfig) ExpectProviderAvailable() error {
	name := "clusterapi-manager-controllers"
	key := types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}
	d := &kappsapi.Deployment{}

	err := wait.PollImmediate(1*time.Second, waitShort, func() (bool, error) {
		if err := tc.client.Get(context.TODO(), key, d); err != nil {
			glog.Errorf("error querying api for Deployment object: %v, retrying...", err)
			return false, nil
		}
		return d.Status.ReadyReplicas > 0, nil
	})
	return err
}

func (tc *testConfig) ExpectNoClusterObject() error {
	listOptions := client.ListOptions{
		Namespace: namespace,
	}
	clusterList := mapiv1beta1.ClusterList{}

	err := wait.PollImmediate(1*time.Second, waitShort, func() (bool, error) {
		if err := tc.client.List(context.TODO(), &listOptions, &clusterList); err != nil {
			glog.Errorf("error querying api for clusterList object: %v, retrying...", err)
			return false, nil
		}
		if len(clusterList.Items) > 0 {
			return false, errors.New("a cluster object was found")
		}
		return true, nil
	})
	return err
}

func (tc *testConfig) ExpectAllMachinesLinkedToANode() error {
	machineAnnotationKey := "machine.openshift.io/machine"
	listOptions := client.ListOptions{
		Namespace: namespace,
	}
	machineList := mapiv1beta1.MachineList{}
	nodeList := corev1.NodeList{}

	err := wait.PollImmediate(1*time.Second, waitShort, func() (bool, error) {
		if err := tc.client.List(context.TODO(), &listOptions, &machineList); err != nil {
			glog.Errorf("error querying api for machineList object: %v, retrying...", err)
			return false, nil
		}
		if err := tc.client.List(context.TODO(), &listOptions, &nodeList); err != nil {
			glog.Errorf("error querying api for nodeList object: %v, retrying...", err)
			return false, nil
		}
		glog.Infof("Waiting for %d machines to become nodes", len(machineList.Items))
		return len(machineList.Items) == len(nodeList.Items), nil
	})
	if err != nil {
		return err
	}

	return wait.PollImmediate(1*time.Second, waitShort, func() (bool, error) {
		nodeNameToMachineAnnotation := make(map[string]string)
		for _, node := range nodeList.Items {
			nodeNameToMachineAnnotation[node.Name] = node.Annotations[machineAnnotationKey]
		}
		for _, machine := range machineList.Items {
			if machine.Status.NodeRef == nil {
				glog.Errorf("machine %s has no NodeRef, retrying...", machine.Name)
				return false, nil
			}
			nodeName := machine.Status.NodeRef.Name
			if nodeNameToMachineAnnotation[nodeName] != fmt.Sprintf("%s/%s", namespace, machine.Name) {
				glog.Errorf("node name %s does not match expected machine name %s, retrying...", nodeName, machine.Name)
				return false, nil
			}
		}
		return true, nil
	})
}

func (tc *testConfig) ExpectNewNodeWhenDeletingMachine() error {
	listOptions := client.ListOptions{
		Namespace: namespace,
	}
	machineList := mapiv1beta1.MachineList{}
	nodeList := corev1.NodeList{}

	glog.Info("Get machineList")
	err := wait.PollImmediate(1*time.Second, waitShort, func() (bool, error) {
		if err := tc.client.List(context.TODO(), &listOptions, &machineList); err != nil {
			glog.Errorf("error querying api for machineList object: %v, retrying...", err)
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return err
	}

	glog.Info("Get nodeList")
	err = wait.PollImmediate(1*time.Second, waitShort, func() (bool, error) {
		if err := tc.client.List(context.TODO(), &listOptions, &nodeList); err != nil {
			glog.Errorf("error querying api for nodeList object: %v, retrying...", err)
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return err
	}

	clusterInitialTotalNodes := len(nodeList.Items)
	clusterInitialTotalMachines := len(machineList.Items)
	var triagedWorkerMachine mapiv1beta1.Machine
	var triagedWorkerNode corev1.Node
MachineLoop:
	for _, m := range machineList.Items {
		if m.Labels["sigs.k8s.io/cluster-api-machine-role"] == "worker" {
			for _, n := range nodeList.Items {
				if m.Status.NodeRef == nil {
					glog.Errorf("no NodeRef found in machine %v", m.Name)
					return errors.New("no NodeRef found in machine")
				}
				if n.Name == m.Status.NodeRef.Name {
					triagedWorkerMachine = m
					triagedWorkerNode = n
					break MachineLoop
				}
			}
		}
	}

	glog.Info("Delete machine")
	err = wait.PollImmediate(1*time.Second, waitShort, func() (bool, error) {
		if err := tc.client.Delete(context.TODO(), &triagedWorkerMachine); err != nil {
			glog.Errorf("error querying api for Deployment object: %v, retrying...", err)
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return err
	}

	err = wait.PollImmediate(1*time.Second, waitMedium, func() (bool, error) {
		if err := tc.client.List(context.TODO(), &listOptions, &machineList); err != nil {
			glog.Errorf("error querying api for machineList object: %v, retrying...", err)
			return false, nil
		}
		glog.Info("Expect new machine to come up")
		return len(machineList.Items) == clusterInitialTotalMachines, nil
	})
	if err != nil {
		return err
	}

	err = wait.PollImmediate(1*time.Second, waitLong, func() (bool, error) {
		if err := tc.client.List(context.TODO(), &listOptions, &nodeList); err != nil {
			glog.Errorf("error querying api for nodeList object: %v, retrying...", err)
			return false, nil
		}
		glog.Info("Expect deleted machine node to go away")
		for _, n := range nodeList.Items {
			if n.Name == triagedWorkerNode.Name {
				return false, nil
			}
		}
		glog.Info("Expect new node to come up")
		return len(nodeList.Items) == clusterInitialTotalNodes, nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (tc *testConfig) ExpectNodeToBeDrainedBeforeDeletingMachine() error {
	listOptions := client.ListOptions{
		Namespace: namespace,
	}

	var machine mapiv1beta1.Machine
	var nodeName string
	var node *corev1.Node

	glog.Info("Get machineList with at least one machine with NodeRef set")
	if err := wait.PollImmediate(1*time.Second, waitShort, func() (bool, error) {
		machineList := mapiv1beta1.MachineList{}
		if err := tc.client.List(context.TODO(), &listOptions, &machineList); err != nil {
			glog.Errorf("error querying api for machineList object: %v, retrying...", err)
			return false, nil
		}
		for _, machineItem := range machineList.Items {
			// empty or non-worker role skipped
			if machineItem.Labels["sigs.k8s.io/cluster-api-machine-role"] == "worker" {
				if machineItem.Status.NodeRef != nil && machineItem.Status.NodeRef.Name != "" {
					machine = machineItem
					nodeName = machineItem.Status.NodeRef.Name
					return true, nil
				}
			}
		}
		return false, fmt.Errorf("no machine found with NodeRef not set")
	}); err != nil {
		return err
	}

	glog.Info("Get nodeList")
	if err := wait.PollImmediate(1*time.Second, waitShort, func() (bool, error) {
		nodeList := corev1.NodeList{}
		if err := tc.client.List(context.TODO(), &listOptions, &nodeList); err != nil {
			glog.Errorf("error querying api for nodeList object: %v, retrying...", err)
			return false, nil
		}
		for _, nodeItem := range nodeList.Items {
			if nodeItem.Name == nodeName {
				node = &nodeItem
				break
			}
		}
		if node == nil {
			return false, fmt.Errorf("node %q not found", nodeName)
		}
		return true, nil
	}); err != nil {
		return err
	}

	glog.Info("Annotate machine with `openshift.io/drain-node` annotation")
	if machine.ObjectMeta.Annotations == nil {
		machine.ObjectMeta.Annotations = make(map[string]string)
	}
	machine.ObjectMeta.Annotations["openshift.io/drain-node"] = "True"
	if err := tc.client.Update(context.TODO(), &machine); err != nil {
		return fmt.Errorf("unable to set `drain-node` annotation")
	}

	glog.Info("Delete machine and observe node draining")
	if err := tc.client.Delete(context.TODO(), &machine); err != nil {
		return fmt.Errorf("unable to delete machine %q", machine.Name)
	}

	return wait.PollImmediate(time.Second, waitShort, func() (bool, error) {
		eventList := corev1.EventList{}
		if err := tc.client.List(context.TODO(), &listOptions, &eventList); err != nil {
			glog.Errorf("error querying api for eventList object: %v, retrying...", err)
			return false, nil
		}

		glog.Infof("Fetching delete machine and node drained events")
		var nodeDrainedEvent *corev1.Event
		var machineDeletedEvent *corev1.Event
		for _, eventItem := range eventList.Items {
			if eventItem.Reason == "Deleted" && eventItem.Message == fmt.Sprintf("Node %q drained", nodeName) {
				nodeDrainedEvent = &eventItem
				continue
			}
			// always take the newest 'machine deleted' event
			if eventItem.Reason == "Deleted" && eventItem.Message == fmt.Sprintf("Deleted machine %v", machine.Name) {
				machineDeletedEvent = &eventItem
			}
		}

		if nodeDrainedEvent == nil {
			glog.Infof("Unable to find %q node drained event", nodeName)
			return false, nil
		}

		if machineDeletedEvent == nil {
			glog.Infof("Unable to find %q machine deleted event", machine.Name)
			return false, nil
		}

		glog.Infof("Node %q drained event recorded: %#v", nodeName, *nodeDrainedEvent)

		if machineDeletedEvent.FirstTimestamp.Before(&nodeDrainedEvent.FirstTimestamp) {
			err := fmt.Errorf("machine %q deleted before node %q got drained", machine.Name, nodeName)
			glog.Error(err)
			return true, err
		}

		return true, nil
	})
}
