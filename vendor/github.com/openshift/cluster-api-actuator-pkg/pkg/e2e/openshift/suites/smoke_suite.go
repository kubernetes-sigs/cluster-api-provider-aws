package suites

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"context"

	"github.com/golang/glog"
	osconfigv1 "github.com/openshift/api/config/v1"
	mapiv1beta1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	caov1alpha1 "github.com/openshift/cluster-autoscaler-operator/pkg/apis/autoscaling/v1alpha1"
	cvoresourcemerge "github.com/openshift/cluster-version-operator/lib/resourcemerge"
	kappsapi "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	waitShort       = 1 * time.Minute
	waitMedium      = 3 * time.Minute
	waitLong        = 10 * time.Minute
	namespace       = "openshift-machine-api"
	workerRoleLabel = "node-role.kubernetes.io/worker"
)

func smokeTestingSuite(tc *testConfig) []struct {
	expect func(ctx context.Context) error
	desc   string
} {
	return []struct {
		expect func(ctx context.Context) error
		desc   string
	}{
		{
			expect: tc.ExpectMachineAPIOperatorAvailable,
			desc:   "[openshift] Expect Machine API operator to be available",
		},
		{
			expect: tc.ExpectMachineAPIClusterOperatorStatusAvailable,
			desc:   "[openshift] Expect Machine API ClusterOperator status to be available",
		},
		{
			expect: tc.ExpectMachineAPIOperatorReconcileControllersDeployment,
			desc:   "[openshift] Expect Machine API Operator to reconcile controllers",
		},
		{
			expect: tc.ExpectClusterAutoscalerOperatorAvailable,
			desc:   "[openshift] Expect Cluster Autoscaler operator to be available",
		},
		{
			expect: tc.ExpectClusterAutoscalerClusterOperatorStatusAvailable,
			desc:   "[openshift] Expect Cluster Autoscaler ClusterOperator status to be available",
		},
		{
			expect: tc.ExpectAllMachinesLinkedToANode,
			desc:   "[openshift] Expect each machines linked to a node",
		},
		{
			expect: tc.ExpectAdditiveReconcileMachineTaints,
			desc:   "[openshift] Expect machines to propagate taints to node in an additive manner",
		},
		{
			expect: tc.ExpectNewNodeWhenDeletingMachine,
			desc:   "[openshift] Expect a new node to come up after deleting a machine from a MachineSet",
		},
		{
			expect: tc.ExpectAutoscalerScalesOut,
			desc:   "[openshift] Expect cluster to scale out when workload increases",
		},
		{
			expect: tc.ExpectNodeToBeDrainedBeforeMachineIsDeleted,
			desc:   "[openshift] Expect node to be drained before its machine is deleted",
		},
	}
}

func (tc *testConfig) ExpectMachineAPIOperatorAvailable(ctx context.Context) error {
	name := "machine-api-operator"
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
		if d.Status.ReadyReplicas < 1 {
			return false, nil
		}
		return true, nil
	})
	return err
}

func (tc *testConfig) ExpectMachineAPIClusterOperatorStatusAvailable(ctx context.Context) error {
	name := "machine-api"
	key := types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}
	clusterOperator := &osconfigv1.ClusterOperator{}

	err := wait.PollImmediate(1*time.Second, waitShort, func() (bool, error) {
		if err := tc.client.Get(context.TODO(), key, clusterOperator); err != nil {
			glog.Errorf("error querying api for OperatorStatus object: %v, retrying...", err)
			return false, nil
		}
		if available := cvoresourcemerge.FindOperatorStatusCondition(clusterOperator.Status.Conditions, osconfigv1.OperatorAvailable); available != nil {
			if available.Status == osconfigv1.ConditionTrue {
				return true, nil
			}
		}
		return false, nil
	})
	return err
}

func (tc *testConfig) ExpectMachineAPIOperatorReconcileControllersDeployment(ctx context.Context) error {
	key := types.NamespacedName{
		Namespace: namespace,
		Name:      "clusterapi-manager-controllers",
	}
	d := &kappsapi.Deployment{}

	glog.Info("Get deployment")
	err := wait.PollImmediate(1*time.Second, waitShort, func() (bool, error) {
		if err := tc.client.Get(context.TODO(), key, d); err != nil {
			glog.Errorf("error querying api for Deployment object: %v, retrying...", err)
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return err
	}

	glog.Info("Delete deployment")
	err = wait.PollImmediate(1*time.Second, waitShort, func() (bool, error) {
		if err := tc.client.Delete(context.TODO(), d); err != nil {
			glog.Errorf("error querying api for Deployment object: %v, retrying...", err)
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return err
	}

	glog.Info("Verify deployment is recreated")
	err = wait.PollImmediate(1*time.Second, waitLong, func() (bool, error) {
		if err := tc.client.Get(context.TODO(), key, d); err != nil {
			glog.Errorf("error querying api for Deployment object: %v, retrying...", err)
			return false, nil
		}
		if d.Status.ReadyReplicas < 1 || !d.DeletionTimestamp.IsZero() {
			return false, nil
		}
		return true, nil
	})
	return err
}

func (tc *testConfig) ExpectClusterAutoscalerOperatorAvailable(ctx context.Context) error {
	name := "cluster-autoscaler-operator"
	key := types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}
	d := &kappsapi.Deployment{}

	err := wait.PollImmediate(1*time.Second, 1*time.Minute, func() (bool, error) {
		if err := tc.client.Get(context.TODO(), key, d); err != nil {
			glog.Errorf("error querying api for Deployment object: %v, retrying...", err)
			return false, nil
		}
		if d.Status.ReadyReplicas < 1 {
			return false, nil
		}
		return true, nil
	})
	return err
}

func (tc *testConfig) ExpectClusterAutoscalerClusterOperatorStatusAvailable(ctx context.Context) error {
	name := "cluster-autoscaler"
	key := types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}
	clusterOperator := &osconfigv1.ClusterOperator{}

	err := wait.PollImmediate(1*time.Second, 1*time.Minute, func() (bool, error) {
		if err := tc.client.Get(context.TODO(), key, clusterOperator); err != nil {
			glog.Errorf("error querying api for OperatorStatus object: %v, retrying...", err)
			return false, nil
		}
		if cvoresourcemerge.IsOperatorStatusConditionTrue(clusterOperator.Status.Conditions, osconfigv1.OperatorAvailable) {
			return true, nil
		}
		return false, nil
	})
	return err
}

func (tc *testConfig) ExpectAllMachinesLinkedToANode(ctx context.Context) error {
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

func (tc *testConfig) ExpectAdditiveReconcileMachineTaints(ctx context.Context) error {
	glog.Info("Verify machine taints are getting applied to node")
	listOptions := client.ListOptions{
		Namespace: namespace,
	}
	machineList := mapiv1beta1.MachineList{}

	if err := tc.client.List(context.TODO(), &listOptions, &machineList); err != nil {
		return fmt.Errorf("error querying api for machineList object: %v", err)

	}
	glog.Info("Got the machine list")
	machine := machineList.Items[0]
	if machine.Status.NodeRef == nil {
		return fmt.Errorf("machine %s has no NodeRef", machine.Name)
	}
	glog.Infof("Got the machine, %s", machine.Name)
	nodeName := machine.Status.NodeRef.Name
	nodeKey := types.NamespacedName{
		Namespace: namespace,
		Name:      nodeName,
	}
	node := &corev1.Node{}

	if err := tc.client.Get(context.TODO(), nodeKey, node); err != nil {
		return fmt.Errorf("error querying api for node object: %v", err)
	}
	glog.Infof("Got the node, %s, from machine, %s", node.Name, machine.Name)
	nodeTaint := corev1.Taint{
		Key:    "not-from-machine",
		Value:  "true",
		Effect: corev1.TaintEffectNoSchedule,
	}
	node.Spec.Taints = []corev1.Taint{nodeTaint}
	if err := tc.client.Update(context.TODO(), node); err != nil {
		return fmt.Errorf("error updating node object with non-machine taint: %v", err)
	}
	glog.Info("Updated node object with taint")
	machineTaint := corev1.Taint{
		Key:    "from-machine",
		Value:  "true",
		Effect: corev1.TaintEffectNoSchedule,
	}
	machine.Spec.Taints = []corev1.Taint{machineTaint}
	if err := tc.client.Update(context.TODO(), &machine); err != nil {
		return fmt.Errorf("error updating machine object with taint: %v", err)
	}
	glog.Info("Updated machine object with taint")
	var expectedTaints = sets.NewString("not-from-machine", "from-machine")
	err := wait.PollImmediate(1*time.Second, waitLong, func() (bool, error) {
		if err := tc.client.Get(context.TODO(), nodeKey, node); err != nil {
			glog.Errorf("error querying api for node object: %v", err)
			return false, nil
		}
		glog.Info("Got the node again for verification of taints")
		var observedTaints = sets.NewString()
		for _, taint := range node.Spec.Taints {
			observedTaints.Insert(taint.Key)
		}
		if expectedTaints.Difference(observedTaints).HasAny("not-from-machine", "from-machine") == false {
			glog.Infof("expected : %v, observed %v , difference %v, ", expectedTaints, observedTaints, expectedTaints.Difference(observedTaints))
			return true, nil
		}
		glog.Infof("All expected taints not found on node. Missing: %v", expectedTaints.Difference(observedTaints))
		return false, nil
	})
	return err
}

func (tc *testConfig) ExpectNewNodeWhenDeletingMachine(ctx context.Context) error {
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
		if m.Labels["machine.openshift.io/cluster-api-machine-role"] == "worker" {
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

// ExpectAutoscalerScalesOut is an smoke test for the autoscaling feature
// Create a clusterAutoscaler object
// Create a machineAutoscaler object
// Create a workLoad to force autoscaling
// Validate the targeted machineSet scales out the field for the expected number of replicas
// Validate the number of nodes in the cluster is growing
// Delete the workLoad and so provoke scale down
// Validate the targeted machineSet scales down its replica count
// Validate the number of nodes scales down to the initial number before scale out
// Delete the machineAutoscaler object
// Delete the clusterAutoscaler object
// ExpectAutoscalerScalesOut is an smoke test for the autoscaling feature
// Create a clusterAutoscaler object
// Create a machineAutoscaler object
// Create a workLoad to force autoscaling
// Validate the targeted machineSet scales out the field for the expected number of replicas
// Validate the number of nodes in the cluster is growing
// Delete the workLoad and so provoke scale down
// Validate the targeted machineSet scales down its replica count
// Validate the number of nodes scales down to the initial number before scale out
// Delete the machineAutoscaler object
// Delete the clusterAutoscaler object
func (tc *testConfig) ExpectAutoscalerScalesOut(ctx context.Context) error {
	listOptions := client.ListOptions{
		Namespace: namespace,
	}
	glog.Info("Get one machineSet")
	machineSetList := mapiv1beta1.MachineSetList{}
	if err := wait.PollImmediate(1*time.Second, waitMedium, func() (bool, error) {
		if err := tc.client.List(ctx, &listOptions, &machineSetList); err != nil {
			glog.Errorf("error querying api for nodeList object: %v, retrying...", err)
			return false, ctx.Err()
		}
		return len(machineSetList.Items) > 0, nil
	}); err != nil {
		return err
	}

	// When we add support for machineDeployments on the installer, cluster-autoscaler and cluster-autoscaler-operator
	// we need to test against deployments instead so we skip this test.
	targetMachineSet := machineSetList.Items[0]
	if ownerReferences := targetMachineSet.GetOwnerReferences(); len(ownerReferences) > 0 {
		glog.Infof("MachineSet %s is owned by a machineDeployment. Please run tests against machineDeployment instead", targetMachineSet.Name)
		return nil
	}

	glog.Infof("Create ClusterAutoscaler and MachineAutoscaler objects. Targeting machineSet %s", targetMachineSet.Name)
	initialNumberOfReplicas := targetMachineSet.Spec.Replicas
	tenSeconds := "10s"
	clusterAutoscaler := caov1alpha1.ClusterAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default",
			Namespace: namespace,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterAutoscaler",
			APIVersion: "autoscaling.openshift.io/v1alpha1",
		},
		Spec: caov1alpha1.ClusterAutoscalerSpec{
			ScaleDown: &caov1alpha1.ScaleDownConfig{
				Enabled:           true,
				DelayAfterAdd:     &tenSeconds,
				DelayAfterDelete:  &tenSeconds,
				DelayAfterFailure: &tenSeconds,
				UnneededTime:      &tenSeconds,
			},
		},
	}
	machineAutoscaler := caov1alpha1.MachineAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("autoscale-%s", targetMachineSet.Name),
			Namespace:    namespace,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "MachineAutoscaler",
			APIVersion: "autoscaling.openshift.io/v1alpha1",
		},
		Spec: caov1alpha1.MachineAutoscalerSpec{
			MaxReplicas: 2,
			MinReplicas: 1,
			ScaleTargetRef: caov1alpha1.CrossVersionObjectReference{
				Name:       targetMachineSet.Name,
				Kind:       "MachineSet",
				APIVersion: "machine.openshift.io/v1beta1",
			},
		},
	}
	if err := wait.PollImmediate(1*time.Second, waitMedium, func() (bool, error) {
		if err := tc.client.Create(ctx, &clusterAutoscaler); err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				glog.Errorf("error querying api for clusterAutoscaler object: %v, retrying...", err)
				return false, ctx.Err()
			}
		}
		if err := tc.client.Create(ctx, &machineAutoscaler); err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				glog.Errorf("error querying api for machineAutoscaler object: %v, retrying...", err)
				return false, ctx.Err()
			}
		}
		return true, nil
	}); err != nil {
		return err
	}

	workLoad := newWorkLoad()

	// We want to clean up these objects on any subsequent error.

	defer func() {
		if workLoad != nil {
			cascadeDelete := metav1.DeletePropagationForeground
			wait.PollImmediate(1*time.Second, waitShort, func() (bool, error) {
				if err := tc.client.Delete(context.TODO(), workLoad, func(opt *client.DeleteOptions) {
					opt.PropagationPolicy = &cascadeDelete
				}); err != nil {
					glog.Errorf("error querying api for workLoad object: %v, retrying...", err)
					return false, nil
				}
				return true, nil
			})
			glog.Info("Deleted workload object")
		}

		wait.PollImmediate(1*time.Second, waitShort, func() (bool, error) {
			if err := tc.client.Delete(context.TODO(), &machineAutoscaler); err != nil {
				glog.Errorf("error querying api for machineAutoscaler object: %v, retrying...", err)
				return false, nil
			}
			return true, nil
		})
		glog.Info("Deleted machineAutoscaler object")

		wait.PollImmediate(1*time.Second, waitShort, func() (bool, error) {
			if err := tc.client.Delete(context.TODO(), &clusterAutoscaler); err != nil {
				glog.Errorf("error querying api for clusterAutoscaler object: %v, retrying...", err)
				return false, nil
			}
			return true, nil
		})
		glog.Info("Deleted clusterAutoscaler object")
	}()

	glog.Info("Get nodeList")
	nodeList := corev1.NodeList{}
	if err := wait.PollImmediate(1*time.Second, waitMedium, func() (bool, error) {
		if err := tc.client.List(ctx, &listOptions, &nodeList); err != nil {
			glog.Errorf("error querying api for nodeList object: %v, retrying...", err)
			return false, ctx.Err()
		}
		return true, nil
	}); err != nil {
		return err
	}

	clusterInitialTotalNodes := len(nodeList.Items)
	glog.Infof("Cluster initial number of nodes is %d", clusterInitialTotalNodes)

	glog.Info("Create workload")

	if err := wait.PollImmediate(1*time.Second, waitMedium, func() (bool, error) {
		if err := tc.client.Create(ctx, workLoad); err != nil {
			glog.Errorf("error querying api for workLoad object: %v, retrying...", err)
			return false, ctx.Err()
		}
		return true, nil
	}); err != nil {
		return err
	}

	glog.Info("Wait for cluster to scale out number of replicas")
	if err := wait.PollImmediate(1*time.Second, waitLong, func() (bool, error) {
		msKey := types.NamespacedName{
			Namespace: namespace,
			Name:      targetMachineSet.Name,
		}
		ms := &mapiv1beta1.MachineSet{}
		if err := tc.client.Get(ctx, msKey, ms); err != nil {
			glog.Errorf("error querying api for clusterAutoscaler object: %v, retrying...", err)
			return false, ctx.Err()
		}
		glog.Infof("MachineSet %s. Initial number of replicas: %d. New number of replicas: %d", targetMachineSet.Name, *initialNumberOfReplicas, *ms.Spec.Replicas)
		return *ms.Spec.Replicas > *initialNumberOfReplicas, nil
	}); err != nil {
		return err
	}

	glog.Info("Wait for cluster to scale out nodes")
	if err := wait.PollImmediate(1*time.Second, waitLong, func() (bool, error) {
		nodeList := corev1.NodeList{}
		if err := tc.client.List(ctx, &listOptions, &nodeList); err != nil {
			glog.Errorf("error querying api for nodeList object: %v, retrying...", err)
			return false, ctx.Err()
		}
		glog.Info("Expect at least a new node to come up")
		glog.Infof("Initial number of nodes: %d. New number of nodes: %d", clusterInitialTotalNodes, len(nodeList.Items))
		return len(nodeList.Items) > clusterInitialTotalNodes, nil
	}); err != nil {
		return err
	}

	glog.Info("Delete workload")
	if err := wait.PollImmediate(1*time.Second, waitMedium, func() (bool, error) {
		cascadeDelete := metav1.DeletePropagationForeground
		if err := tc.client.Delete(ctx, workLoad, func(opt *client.DeleteOptions) {
			opt.PropagationPolicy = &cascadeDelete
		}); err != nil {
			glog.Errorf("error querying api for workLoad object: %v, retrying...", err)
			return false, ctx.Err()
		}
		workLoad = nil
		return true, nil
	}); err != nil {
		return err
	}

	// As we have just deleted the workload the autoscaler will
	// start to scale down the unneeded nodes. We wait for that
	// condition; if successful we assert that (a smoke test of)
	// scale down is functional.

	glog.Info("Wait for cluster to match initial number of replicas")
	if err := wait.PollImmediate(1*time.Second, waitLong, func() (bool, error) {
		msKey := types.NamespacedName{
			Namespace: namespace,
			Name:      targetMachineSet.Name,
		}
		ms := &mapiv1beta1.MachineSet{}
		if err := tc.client.Get(ctx, msKey, ms); err != nil {
			glog.Errorf("error querying api for machineSet object: %v, retrying...", err)
			return false, nil
		}
		glog.Infof("Initial number of replicas: %d. Current number of replicas: %d", *initialNumberOfReplicas, *ms.Spec.Replicas)
		return *ms.Spec.Replicas == *initialNumberOfReplicas, nil
	}); err != nil {
		return err
	}

	glog.Info("Wait for cluster to match initial number of nodes")
	return wait.PollImmediate(1*time.Second, waitLong, func() (bool, error) {
		nodeList := corev1.NodeList{}
		if err := tc.client.List(ctx, &listOptions, &nodeList); err != nil {
			glog.Errorf("error querying api for nodeList object: %v, retrying...", err)
			return false, ctx.Err()
		}
		glog.Infof("Initial number of nodes: %d. Current number of nodes: %d", clusterInitialTotalNodes, len(nodeList.Items))
		return len(nodeList.Items) == clusterInitialTotalNodes, nil
	})
}

func (tc *testConfig) ExpectNodeToBeDrainedBeforeMachineIsDeleted(ctx context.Context) error {
	delObjects := make(map[string]runtime.Object)

	defer func() {
		// Remove resources
		for key := range delObjects {
			glog.Infof("Deleting object %q", key)
			if err := tc.client.Delete(context.TODO(), delObjects[key]); err != nil {
				glog.Errorf("Unable to delete object %q: %v", key, err)
			}
		}
	}()

	// Take the first worker machineset (assuming only worker machines are backed by machinesets)
	machinesets := mapiv1beta1.MachineSetList{}
	if err := wait.PollImmediate(1*time.Second, 1*time.Minute, func() (bool, error) {
		if err := tc.client.List(context.TODO(), &client.ListOptions{}, &machinesets); err != nil {
			glog.Errorf("Error querying api for machineset object: %v, retrying...", err)
			return false, nil
		}
		if len(machinesets.Items) < 1 {
			glog.Errorf("Expected at least one machineset, have none")
			return false, nil
		}
		return true, nil
	}); err != nil {
		return err
	}

	// Create two machines
	machine1 := machineFromMachineset(&machinesets.Items[0])
	machine1.Name = "machine1"

	if err := tc.client.Create(context.TODO(), machine1); err != nil {
		return fmt.Errorf("unable to create machine %q: %v", machine1.Name, err)
	}

	delObjects["machine1"] = machine1

	machine2 := machineFromMachineset(&machinesets.Items[0])
	machine2.Name = "machine2"

	if err := tc.client.Create(context.TODO(), machine2); err != nil {
		return fmt.Errorf("unable to create machine %q: %v", machine2.Name, err)
	}

	delObjects["machine2"] = machine2

	// Wait until both new nodes are ready
	if err := wait.PollImmediate(1*time.Second, 10*time.Minute, func() (bool, error) {
		nodes := corev1.NodeList{}
		listOpt := &client.ListOptions{}
		listOpt.MatchingLabels(nodeDrainLabels)
		if err := tc.client.List(context.TODO(), listOpt, &nodes); err != nil {
			glog.Errorf("Error querying api for Node object: %v, retrying...", err)
			return false, nil
		}
		// expecting nodeGroupSize nodes
		nodeCounter := 0
		for _, node := range nodes.Items {
			if _, exists := node.Labels[workerRoleLabel]; !exists {
				continue
			}

			if !isNodeReady(&node) {
				continue
			}

			nodeCounter++
		}

		if nodeCounter < 2 {
			glog.Errorf("Expecting 2 nodes with %#v labels in Ready state, got %v", nodeDrainLabels, nodeCounter)
			return false, nil
		}

		glog.Infof("Expected number (2) of nodes with %v label in Ready state found", nodeDrainLabels)
		return true, nil
	}); err != nil {
		return err
	}

	rc := replicationControllerWorkload(namespace)
	if err := tc.client.Create(context.TODO(), rc); err != nil {
		return fmt.Errorf("unable to create RC %q: %v", rc.Name, err)
	}

	delObjects["rc"] = rc

	pdb := podDisruptionBudget(namespace)
	if err := tc.client.Create(context.TODO(), pdb); err != nil {
		return fmt.Errorf("unable to create PDB %q: %v", pdb.Name, err)
	}

	delObjects["pdb"] = pdb

	// Wait until all replicas are ready
	if err := wait.PollImmediate(1*time.Second, 10*time.Minute, func() (bool, error) {
		rcObj := corev1.ReplicationController{}
		key := types.NamespacedName{
			Namespace: rc.Namespace,
			Name:      rc.Name,
		}
		if err := tc.client.Get(context.TODO(), key, &rcObj); err != nil {
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
	}); err != nil {
		return err
	}

	// All pods are distributed evenly among all nodes so it's fine to drain
	// random node and observe reconciliation of pods on the other one.
	if err := tc.client.Delete(context.TODO(), machine1); err != nil {
		return fmt.Errorf("unable to delete machine %q: %v", machine1.Name, err)
	}

	delete(delObjects, "machine1")

	// We still should be able to list the machine as until rc.replicas-1 are running on the other node
	var drainedNodeName string
	if err := wait.PollImmediate(1*time.Second, 10*time.Minute, func() (bool, error) {
		machine := mapiv1beta1.Machine{}

		key := types.NamespacedName{
			Namespace: machine1.Namespace,
			Name:      machine1.Name,
		}
		if err := tc.client.Get(context.TODO(), key, &machine); err != nil {
			glog.Errorf("Error querying api machine %q object: %v, retrying...", machine1.Name, err)
			return false, nil
		}
		if machine.Status.NodeRef == nil || machine.Status.NodeRef.Kind != "Node" {
			glog.Error("Machine %q not linked to a node", machine.Name)
		}

		drainedNodeName = machine.Status.NodeRef.Name
		node := corev1.Node{}

		if err := tc.client.Get(context.TODO(), types.NamespacedName{Name: drainedNodeName}, &node); err != nil {
			glog.Errorf("Error querying api node %q object: %v, retrying...", drainedNodeName, err)
			return false, nil
		}

		if !node.Spec.Unschedulable {
			glog.Errorf("Node %q is expected to be marked as unschedulable, it is not", node.Name)
			return false, nil
		}

		glog.Infof("Node %q is mark unschedulable as expected", node.Name)

		pods := corev1.PodList{}
		listOpt := &client.ListOptions{}
		listOpt.MatchingLabels(rc.Spec.Selector)
		if err := tc.client.List(context.TODO(), listOpt, &pods); err != nil {
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
		if err := tc.client.Get(context.TODO(), key, &rcObj); err != nil {
			glog.Errorf("Error querying api RC %q object: %v, retrying...", rc.Name, err)
			return false, nil
		}

		// The point of the test is to make sure majority of the pods is rescheduled
		// to other nodes. Pod disruption budget makes sure at most one pod
		// owned by the RC is not Ready. So no need to test it. Though, usefull to have it printed.
		glog.Infof("RC ReadyReplicas/Replicas: %v/%v", rcObj.Status.ReadyReplicas, rcObj.Status.Replicas)

		// This makes sure at most one replica is not ready
		if rcObj.Status.Replicas-rcObj.Status.ReadyReplicas > 1 {
			return false, fmt.Errorf("pod disruption budget not respecpted, node was not properly drained")
		}

		// Depends on timing though a machine can be deleted even before there is only
		// one pod left on the node (that is being evicted).
		if podCounter > 2 {
			glog.Infof("Expecting at most 2 pods to be scheduled to drained node %v, got %v", machine.Status.NodeRef.Name, podCounter)
			return false, nil
		}

		glog.Info("Expected result: all pods from the RC up to last one or two got scheduled to a different node while respecting PDB")
		return true, nil
	}); err != nil {
		return err
	}

	// Validate the machine is deleted
	if err := wait.PollImmediate(5*time.Second, 1*time.Minute, func() (bool, error) {
		machine := mapiv1beta1.Machine{}

		key := types.NamespacedName{
			Namespace: machine1.Namespace,
			Name:      machine1.Name,
		}
		err := tc.client.Get(context.TODO(), key, &machine)
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
	}); err != nil {
		return err
	}

	// Validate underlying node is removed as well
	if err := wait.PollImmediate(5*time.Second, waitLong, func() (bool, error) {
		node := corev1.Node{}

		key := types.NamespacedName{
			Name: drainedNodeName,
		}
		err := tc.client.Get(context.TODO(), key, &node)
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
	}); err != nil {
		return err
	}

	return nil
}
