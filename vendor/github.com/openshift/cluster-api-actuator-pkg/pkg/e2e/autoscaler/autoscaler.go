package autoscaler

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	g "github.com/onsi/ginkgo"
	o "github.com/onsi/gomega"
	e2e "github.com/openshift/cluster-api-actuator-pkg/pkg/e2e/framework"
	mapiv1beta1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	caov1 "github.com/openshift/cluster-autoscaler-operator/pkg/apis/autoscaling/v1"
	caov1beta1 "github.com/openshift/cluster-autoscaler-operator/pkg/apis/autoscaling/v1beta1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/utils/pointer"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	autoscalingTestLabel = "test.autoscaling.label"
)

func newWorkLoad() *batchv1.Job {
	backoffLimit := int32(4)
	completions := int32(50)
	parallelism := int32(50)
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "workload",
			Namespace: "default",
			Labels:    map[string]string{autoscalingTestLabel: ""},
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "workload",
							Image: "busybox",
							Command: []string{
								"sleep",
								"86400", // 1 day
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									"memory": resource.MustParse("500Mi"),
									"cpu":    resource.MustParse("500m"),
								},
							},
						},
					},
					RestartPolicy: corev1.RestartPolicy("Never"),
					NodeSelector: map[string]string{
						e2e.WorkerNodeRoleLabel: "",
					},
					Tolerations: []corev1.Toleration{
						{
							Key:      "kubemark",
							Operator: corev1.TolerationOpExists,
						},
					},
				},
			},
			BackoffLimit: &backoffLimit,
			Completions:  &completions,
			Parallelism:  &parallelism,
		},
	}
}

func labelMachineSetNodes(client runtimeclient.Client, ms *mapiv1beta1.MachineSet, nodeTestLabel string) error {
	return wait.PollImmediate(e2e.RetryMedium, e2e.WaitShort, func() (bool, error) {
		scaledMachines := mapiv1beta1.MachineList{}
		if err := client.List(context.TODO(), runtimeclient.MatchingLabels(ms.Spec.Selector.MatchLabels), &scaledMachines); err != nil {
			glog.Errorf("Error querying api for machineset object: %v, retrying...", err)
			return false, nil
		}

		// get all linked nodes and label them
		for _, machine := range scaledMachines.Items {
			if machine.Status.NodeRef == nil {
				glog.Errorf("Machine %q does not have node reference set", machine.Name)
				return false, nil
			}
			node := corev1.Node{}
			if err := client.Get(context.TODO(), types.NamespacedName{Name: machine.Status.NodeRef.Name}, &node); err != nil {
				glog.Errorf("error querying api for node object: %v, retrying...", err)
				return false, nil
			}

			labelNode := false
			if node.Labels == nil {
				labelNode = true
			} else if _, exists := node.Labels[nodeTestLabel]; !exists {
				labelNode = true
			}

			if labelNode {
				nodeCopy := node.DeepCopy()
				if nodeCopy.Labels == nil {
					nodeCopy.Labels = make(map[string]string)
				}
				nodeCopy.Labels[nodeTestLabel] = ""
				if err := client.Update(context.TODO(), nodeCopy); err != nil {
					glog.Errorf("error updating api for node object: %v, retrying...", err)
					return false, nil
				}
				glog.Infof("Labeling node %q with %q label", nodeCopy.Name, nodeTestLabel)
			}
		}
		return true, nil
	})
}

// Build default CA resource to allow fast scaling up and down
func clusterAutoscalerResource() *caov1.ClusterAutoscaler {
	tenSecondString := "10s"
	return &caov1.ClusterAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default",
			Namespace: e2e.TestContext.MachineApiNamespace,
			Labels: map[string]string{
				autoscalingTestLabel: "",
			},
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterAutoscaler",
			APIVersion: "autoscaling.openshift.io/v1",
		},
		Spec: caov1.ClusterAutoscalerSpec{
			ScaleDown: &caov1.ScaleDownConfig{
				Enabled:           true,
				DelayAfterAdd:     &tenSecondString,
				DelayAfterDelete:  &tenSecondString,
				DelayAfterFailure: &tenSecondString,
				UnneededTime:      &tenSecondString,
			},
		},
	}
}

// Build MA resource from targeted machineset
func machineAutoscalerResource(targetMachineSet *mapiv1beta1.MachineSet, minReplicas, maxReplicas int32) *caov1beta1.MachineAutoscaler {
	return &caov1beta1.MachineAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("autoscale-%s", targetMachineSet.Name),
			Namespace:    e2e.TestContext.MachineApiNamespace,
			Labels: map[string]string{
				autoscalingTestLabel: "",
			},
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "MachineAutoscaler",
			APIVersion: "autoscaling.openshift.io/v1beta1",
		},
		Spec: caov1beta1.MachineAutoscalerSpec{
			MaxReplicas: maxReplicas,
			MinReplicas: minReplicas,
			ScaleTargetRef: caov1beta1.CrossVersionObjectReference{
				Name:       targetMachineSet.Name,
				Kind:       "MachineSet",
				APIVersion: "machine.openshift.io/v1beta1",
			},
		},
	}
}

var _ = g.Describe("[Feature:Machines] Autoscaler should", func() {
	defer g.GinkgoRecover()
	scaleUpFunc := func(client runtimeclient.Client, targetMachineSet mapiv1beta1.MachineSet, nodeTestLabel string, initialNumberOfReplicas int32, expectedReplicas int32) {
		g.By(fmt.Sprintf("Creating MachineAutoscaler object for targetMachineSet %q", targetMachineSet.Name))
		machineAutoscaler := machineAutoscalerResource(&targetMachineSet, 1, expectedReplicas)
		err := client.Create(context.TODO(), machineAutoscaler)
		o.Expect(err).NotTo(o.HaveOccurred())

		g.By(fmt.Sprintf("Waiting for cluster to scale out number of replicas of targetMachineSet %q", targetMachineSet.Name))
		err = wait.PollImmediate(5*time.Second, e2e.WaitLong, func() (bool, error) {
			ms, err := e2e.GetMachineSet(context.TODO(), client, targetMachineSet.Name)
			if err != nil {
				glog.Errorf("error getting machineset object: %v, retrying...", err)
				return false, nil
			}
			glog.Infof("MachineSet %q. Initial number of replicas: %d. Current number of replicas: %d", targetMachineSet.Name, initialNumberOfReplicas, pointer.Int32PtrDerefOr(ms.Spec.Replicas, e2e.DefaultMachineSetReplicas))
			return pointer.Int32PtrDerefOr(ms.Spec.Replicas, e2e.DefaultMachineSetReplicas) == expectedReplicas, nil
		})
		o.Expect(err).NotTo(o.HaveOccurred())

		g.By(fmt.Sprintf("Waiting for cluster to scale up nodes for targetMachineSet %q", targetMachineSet.Name))
		err = wait.PollImmediate(5*time.Second, e2e.WaitLong, func() (bool, error) {
			scaledMachines := mapiv1beta1.MachineList{}
			if err := client.List(context.TODO(), runtimeclient.MatchingLabels(targetMachineSet.Spec.Selector.MatchLabels), &scaledMachines); err != nil {
				glog.Errorf("Error querying api for machineset object: %v, retrying...", err)
				return false, nil
			}

			// get all linked nodes and label them
			nodeCounter := 0
			for _, machine := range scaledMachines.Items {
				if machine.Status.NodeRef == nil {
					glog.Errorf("Machine %q does not have node reference set", machine.Name)
					return false, nil
				}
				glog.Infof("Machine %q is linked to node %q", machine.Name, machine.Status.NodeRef.Name)
				nodeCounter++
			}

			glog.Infof("Expecting %d nodes. Current number of nodes in the group: %d", expectedReplicas, nodeCounter)
			return nodeCounter == int(expectedReplicas), nil
		})
		o.Expect(err).NotTo(o.HaveOccurred())

		g.By(fmt.Sprintf("Labeling all nodes belonging to the machineset %q (after scale up phase)", targetMachineSet.Name))
		err = labelMachineSetNodes(client, &targetMachineSet, nodeTestLabel)
		o.Expect(err).NotTo(o.HaveOccurred())
	}

	scaleDownFunc := func(client runtimeclient.Client, targetMachineSet string, nodeTestLabel string, initialNumberOfReplicas int32) {
		g.By(fmt.Sprintf("Waiting for machineset %q to have at most initial number of replicas", targetMachineSet))
		err := wait.PollImmediate(5*time.Second, e2e.WaitLong, func() (bool, error) {
			ms, err := e2e.GetMachineSet(context.TODO(), client, targetMachineSet)
			if err != nil {
				glog.Errorf("error getting machineset object: %v, retrying...", err)
				return false, nil
			}
			msReplicas := pointer.Int32PtrDerefOr(ms.Spec.Replicas, e2e.DefaultMachineSetReplicas)
			glog.Infof("Initial number of replicas: %d. Current number of replicas: %d", initialNumberOfReplicas, msReplicas)
			if msReplicas > initialNumberOfReplicas {
				return false, nil
			}

			// Make sure all scaled down nodes are really gone (so they don't affect tests to be run next)
			scaledNodes := corev1.NodeList{}
			if err := client.List(context.TODO(), runtimeclient.MatchingLabels(map[string]string{nodeTestLabel: ""}), &scaledNodes); err != nil {
				glog.Errorf("Error querying api for node objects: %v, retrying...", err)
				return false, nil
			}
			scaledNodesLen := int32(len(scaledNodes.Items))
			glog.Infof("Current number of replicas: %d. Current number of nodes: %d", msReplicas, scaledNodesLen)
			return scaledNodesLen <= msReplicas && scaledNodesLen <= initialNumberOfReplicas, nil
		})
		o.Expect(err).NotTo(o.HaveOccurred())
	}

	g.It("scale out", func() {
		var err error
		client, err := e2e.LoadClient()
		o.Expect(err).NotTo(o.HaveOccurred())

		nodeTestLabel0 := fmt.Sprintf("machine.openshift.io/autoscaling-test-%v", string(uuid.NewUUID()))
		nodeTestLabel1 := fmt.Sprintf("machine.openshift.io/autoscaling-test-%v", string(uuid.NewUUID()))

		// We want to clean up these objects on any subsequent error.
		defer func() {
			err = e2e.DeleteObjectsByLabels(context.TODO(), client, map[string]string{autoscalingTestLabel: ""}, &batchv1.JobList{})
			if err != nil {
				// if this one fails, there are still other resources to be deleted.
				glog.Warning(err)
			} else {
				glog.Info("Deleted workload object")
			}

			err = e2e.DeleteObjectsByLabels(context.TODO(), client, map[string]string{autoscalingTestLabel: ""}, &caov1beta1.MachineAutoscalerList{})
			if err != nil {
				// if this one fails, there are still other resources to be deleted.
				glog.Warning(err)
			} else {
				glog.Info("Deleted machineAutoscaler object")
			}

			err = e2e.DeleteObjectsByLabels(context.TODO(), client, map[string]string{autoscalingTestLabel: ""}, &caov1.ClusterAutoscalerList{})
			if err != nil {
				// if this one fails, there is no point of returning an error as this is the last resource deletion action
				glog.Warning(err)
			} else {
				glog.Info("Deleted clusterAutoscaler object")
			}
		}()

		g.By("Getting target machineSet")
		machinesets, err := e2e.GetMachineSets(context.TODO(), client)
		o.Expect(err).NotTo(o.HaveOccurred())
		o.Expect(len(machinesets)).To(o.BeNumerically(">", 1))

		targetMachineSet0 := machinesets[0]
		glog.Infof("Target machineSet0 %q", targetMachineSet0.Name)
		targetMachineSet1 := machinesets[1]
		glog.Infof("Target machineSet1 %q", targetMachineSet1.Name)

		// When we add support for machineDeployments on the installer, cluster-autoscaler and cluster-autoscaler-operator
		// we need to test against deployments instead so we skip this test.
		if ownerReferences0 := targetMachineSet0.GetOwnerReferences(); len(ownerReferences0) > 0 {
			// glog.Infof("MachineSet %s is owned by a machineDeployment. Please run tests against machineDeployment instead", targetMachineSet.Name)
			g.Skip(fmt.Sprintf("MachineSet %q is owned by a machineDeployment. Please run tests against machineDeployment instead", targetMachineSet0.Name))
		}

		if ownerReferences1 := targetMachineSet1.GetOwnerReferences(); len(ownerReferences1) > 0 {
			g.Skip(fmt.Sprintf("MachineSet %q is owned by a machineDeployment. Please run tests against machineDeployment instead", targetMachineSet1.Name))
		}

		g.By("Create ClusterAutoscaler object")
		clusterAutoscaler := clusterAutoscalerResource()
		err = client.Create(context.TODO(), clusterAutoscaler)
		o.Expect(err).NotTo(o.HaveOccurred())

		initialNumberOfReplicas0 := pointer.Int32PtrDerefOr(targetMachineSet0.Spec.Replicas, e2e.DefaultMachineSetReplicas)
		initialNumberOfReplicas1 := pointer.Int32PtrDerefOr(targetMachineSet1.Spec.Replicas, e2e.DefaultMachineSetReplicas)
		glog.Infof("initialNumberOfReplicas0 %d, initialNumberOfReplicas1 %d", initialNumberOfReplicas0, initialNumberOfReplicas1)

		g.By("Creating workload")
		err = client.Create(context.TODO(), newWorkLoad())
		o.Expect(err).NotTo(o.HaveOccurred())

		scaleUpFunc(client, targetMachineSet1, nodeTestLabel1, initialNumberOfReplicas1, initialNumberOfReplicas1+2)
		scaleUpFunc(client, targetMachineSet0, nodeTestLabel0, initialNumberOfReplicas0, initialNumberOfReplicas1+1)

		// Delete workload
		g.By("Deleting workload")
		err = e2e.DeleteObjectsByLabels(context.TODO(), client, map[string]string{autoscalingTestLabel: ""}, &batchv1.JobList{})
		o.Expect(err).NotTo(o.HaveOccurred())

		// As we have just deleted the workload the autoscaler will
		// start to scale down the unneeded nodes. We wait for that
		// condition; if successful we assert that (a smoke test of)
		// scale down is functional.
		scaleDownFunc(client, targetMachineSet0.Name, nodeTestLabel0, initialNumberOfReplicas0)
		scaleDownFunc(client, targetMachineSet1.Name, nodeTestLabel1, initialNumberOfReplicas1)
	})

})
