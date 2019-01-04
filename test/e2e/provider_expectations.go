package main

import (
	"time"

	"context"
	"errors"
	"fmt"

	"github.com/golang/glog"
	kappsapi "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	capiv1alpha1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
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

func (tc *testConfig) ExpectOneClusterObject() error {
	listOptions := client.ListOptions{
		Namespace: namespace,
	}
	clusterList := capiv1alpha1.ClusterList{}

	err := wait.PollImmediate(1*time.Second, waitShort, func() (bool, error) {
		if err := tc.client.List(context.TODO(), &listOptions, &clusterList); err != nil {
			glog.Errorf("error querying api for clusterList object: %v, retrying...", err)
			return false, nil
		}
		if len(clusterList.Items) != 1 {
			return false, errors.New("more than one cluster object found")
		}
		return true, nil
	})
	return err
}

func (tc *testConfig) ExpectAllMachinesLinkedToANode() error {
	machineAnnotationKey := "cluster.k8s.io/machine"
	listOptions := client.ListOptions{
		Namespace: namespace,
	}
	machineList := capiv1alpha1.MachineList{}
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
	machineList := capiv1alpha1.MachineList{}
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
	var triagedWorkerMachine capiv1alpha1.Machine
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
