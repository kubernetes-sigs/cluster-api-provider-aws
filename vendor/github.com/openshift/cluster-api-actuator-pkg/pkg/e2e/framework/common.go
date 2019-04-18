package framework

import (
	"context"
	"fmt"
	"time"

	mapiv1beta1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	caov1 "github.com/openshift/cluster-autoscaler-operator/pkg/apis/autoscaling/v1"
	caov1beta1 "github.com/openshift/cluster-autoscaler-operator/pkg/apis/autoscaling/v1beta1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	WorkerNodeRoleLabel = "node-role.kubernetes.io/worker"
	WaitShort           = 1 * time.Minute
	WaitMedium          = 3 * time.Minute
	WaitLong            = 10 * time.Minute
	RetryMedium         = 5 * time.Second

	// DefaultMachineSetReplicas is the default number of replicas of a machineset
	// if MachineSet.Spec.Replicas field is set to nil
	DefaultMachineSetReplicas = 0
)

// GetNodes gets a list of nodes from a running cluster
// Optionaly, labels may be used to constrain listed nodes.
func GetNodes(client runtimeclient.Client, labels ...map[string]string) ([]corev1.Node, error) {
	nodeList := corev1.NodeList{}
	listOptions := &runtimeclient.ListOptions{}
	if len(labels) > 0 && len(labels[0]) > 0 {
		listOptions.MatchingLabels(labels[0])
	}
	if err := client.List(context.TODO(), listOptions, &nodeList); err != nil {
		return nil, fmt.Errorf("error querying api for nodeList object: %v", err)
	}
	return nodeList.Items, nil
}

// GetMachineSets gets a list of machinesets from the default machine API namespace.
// Optionaly, labels may be used to constrain listed machinesets.
func GetMachineSets(ctx context.Context, client runtimeclient.Client, labels ...map[string]string) ([]mapiv1beta1.MachineSet, error) {
	machineSetList := &mapiv1beta1.MachineSetList{}
	listOptions := runtimeclient.InNamespace(TestContext.MachineApiNamespace)
	if len(labels) > 0 && len(labels[0]) > 0 {
		listOptions.MatchingLabels(labels[0])
	}
	if err := client.List(ctx, listOptions, machineSetList); err != nil {
		return nil, fmt.Errorf("error querying api for machineSetList object: %v", err)
	}
	return machineSetList.Items, nil
}

// GetMachineSet gets a machineset by its name from the default machine API namespace.
func GetMachineSet(ctx context.Context, client runtimeclient.Client, machineSetName string) (*mapiv1beta1.MachineSet, error) {
	machineSet := &mapiv1beta1.MachineSet{}
	if err := client.Get(ctx, runtimeclient.ObjectKey{Namespace: TestContext.MachineApiNamespace, Name: machineSetName}, machineSet); err != nil {
		return nil, fmt.Errorf("error querying api for machineSet object: %v", err)
	}
	return machineSet, nil
}

// GetMachines gets a list of machinesets from the default machine API namespace.
// Optionaly, labels may be used to constrain listed machinesets.
func GetMachines(ctx context.Context, client runtimeclient.Client, labels ...map[string]string) ([]mapiv1beta1.Machine, error) {
	machineList := &mapiv1beta1.MachineList{}
	listOptions := runtimeclient.InNamespace(TestContext.MachineApiNamespace)
	if len(labels) > 0 && len(labels[0]) > 0 {
		listOptions.MatchingLabels(labels[0])
	}
	if err := client.List(ctx, listOptions, machineList); err != nil {
		return nil, fmt.Errorf("error querying api for machineList object: %v", err)
	}
	return machineList.Items, nil
}

// GetMachine get a machine by its name from the default machine API namespace.
func GetMachine(ctx context.Context, client runtimeclient.Client, machineName string) (*mapiv1beta1.Machine, error) {
	machine := &mapiv1beta1.Machine{}
	if err := client.Get(ctx, runtimeclient.ObjectKey{Namespace: TestContext.MachineApiNamespace, Name: machineName}, machine); err != nil {
		return nil, fmt.Errorf("error querying api for machine object: %v", err)
	}
	return machine, nil
}

// DeleteObjectsByLabels list all objects of a given kind by labels and deletes them.
// Currently supported kinds:
// - caov1beta1.MachineAutoscalerList
// - caov1.ClusterAutoscalerList
// - batchv1.JobList
func DeleteObjectsByLabels(ctx context.Context, client runtimeclient.Client, labels map[string]string, list runtime.Object) error {
	if err := client.List(ctx, runtimeclient.MatchingLabels(labels), list); err != nil {
		return fmt.Errorf("Unable to list objects: %v", err)
	}

	// TODO(jchaloup): find a way how to list the items independent of a kind
	var objs []runtime.Object
	switch d := list.(type) {
	case *caov1beta1.MachineAutoscalerList:
		for _, item := range d.Items {
			objs = append(objs, runtime.Object(&item))
		}
	case *caov1.ClusterAutoscalerList:
		for _, item := range d.Items {
			objs = append(objs, runtime.Object(&item))
		}
	case *batchv1.JobList:
		for _, item := range d.Items {
			objs = append(objs, runtime.Object(&item))
		}

	default:
		return fmt.Errorf("List type %#v not recognized", list)
	}

	cascadeDelete := metav1.DeletePropagationForeground
	for _, obj := range objs {
		if err := client.Delete(ctx, obj, func(opt *runtimeclient.DeleteOptions) {
			opt.PropagationPolicy = &cascadeDelete
		}); err != nil {
			return fmt.Errorf("error deleting object: %v", err)
		}
	}

	return nil
}
