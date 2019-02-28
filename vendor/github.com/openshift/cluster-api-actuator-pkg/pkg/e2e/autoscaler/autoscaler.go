package autoscaler

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
	caov1alpha1 "github.com/openshift/cluster-autoscaler-operator/pkg/apis/autoscaling/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/apimachinery/pkg/util/wait"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func newWorkLoad() *batchv1.Job {
	backoffLimit := int32(4)
	completions := int32(50)
	parallelism := int32(50)
	activeDeadlineSeconds := int64(100)
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "workload",
			Namespace: "default",
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
								"300",
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
						"node-role.kubernetes.io/worker": "",
					},
					Tolerations: []corev1.Toleration{
						{
							Key:      "kubemark",
							Operator: corev1.TolerationOpExists,
						},
					},
				},
			},
			ActiveDeadlineSeconds: &activeDeadlineSeconds,
			BackoffLimit:          &backoffLimit,
			Completions:           &completions,
			Parallelism:           &parallelism,
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

var _ = g.Describe("[Feature:Machines] Autoscaler should", func() {
	defer g.GinkgoRecover()

	g.It("scale out", func() {
		var err error
		client, err := e2e.LoadClient()
		o.Expect(err).NotTo(o.HaveOccurred())
		glog.Info("Get one machineSet")
		machineSetList := mapiv1beta1.MachineSetList{}
		err = wait.PollImmediate(1*time.Second, e2e.WaitMedium, func() (bool, error) {
			if err := client.List(context.TODO(), runtimeclient.InNamespace(e2e.TestContext.MachineApiNamespace), &machineSetList); err != nil {
				glog.Errorf("error querying api for nodeList object: %v, retrying...", err)
				return false, err
			}
			return len(machineSetList.Items) > 0, nil
		})
		o.Expect(err).NotTo(o.HaveOccurred())

		// When we add support for machineDeployments on the installer, cluster-autoscaler and cluster-autoscaler-operator
		// we need to test against deployments instead so we skip this test.
		targetMachineSet := machineSetList.Items[0]
		if ownerReferences := targetMachineSet.GetOwnerReferences(); len(ownerReferences) > 0 {
			// glog.Infof("MachineSet %s is owned by a machineDeployment. Please run tests against machineDeployment instead", targetMachineSet.Name)
			g.Skip(fmt.Sprintf("MachineSet %s is owned by a machineDeployment. Please run tests against machineDeployment instead", targetMachineSet.Name))
		}

		glog.Infof("Create ClusterAutoscaler and MachineAutoscaler objects. Targeting machineSet %s", targetMachineSet.Name)
		initialNumberOfReplicas := targetMachineSet.Spec.Replicas
		tenSecondString := "10s"
		clusterAutoscaler := caov1alpha1.ClusterAutoscaler{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "default",
				Namespace: e2e.TestContext.MachineApiNamespace,
			},
			TypeMeta: metav1.TypeMeta{
				Kind:       "ClusterAutoscaler",
				APIVersion: "autoscaling.openshift.io/v1alpha1",
			},
			Spec: caov1alpha1.ClusterAutoscalerSpec{
				ScaleDown: &caov1alpha1.ScaleDownConfig{
					Enabled:           true,
					DelayAfterAdd:     &tenSecondString,
					DelayAfterDelete:  &tenSecondString,
					DelayAfterFailure: &tenSecondString,
					UnneededTime:      &tenSecondString,
				},
			},
		}
		machineAutoscaler := caov1alpha1.MachineAutoscaler{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: fmt.Sprintf("autoscale-%s", targetMachineSet.Name),
				Namespace:    e2e.TestContext.MachineApiNamespace,
			},
			TypeMeta: metav1.TypeMeta{
				Kind:       "MachineAutoscaler",
				APIVersion: "autoscaling.openshift.io/v1alpha1",
			},
			Spec: caov1alpha1.MachineAutoscalerSpec{
				MaxReplicas: *initialNumberOfReplicas + 1,
				MinReplicas: 1,
				ScaleTargetRef: caov1alpha1.CrossVersionObjectReference{
					Name:       targetMachineSet.Name,
					Kind:       "MachineSet",
					APIVersion: "machine.openshift.io/v1beta1",
				},
			},
		}
		err = wait.PollImmediate(1*time.Second, e2e.WaitMedium, func() (bool, error) {
			if err := client.Create(context.TODO(), &clusterAutoscaler); err != nil {
				if !strings.Contains(err.Error(), "already exists") {
					glog.Errorf("error querying api for clusterAutoscaler object: %v, retrying...", err)
					return false, err
				}
			}
			if err := client.Create(context.TODO(), &machineAutoscaler); err != nil {
				if !strings.Contains(err.Error(), "already exists") {
					glog.Errorf("error querying api for machineAutoscaler object: %v, retrying...", err)
					return false, err
				}
			}
			return true, nil
		})
		o.Expect(err).NotTo(o.HaveOccurred())

		workLoad := newWorkLoad()

		// We want to clean up these objects on any subsequent error.
		defer func() {
			if workLoad != nil {
				cascadeDelete := metav1.DeletePropagationForeground
				wait.PollImmediate(1*time.Second, e2e.WaitShort, func() (bool, error) {
					if err := client.Delete(context.TODO(), workLoad, func(opt *runtimeclient.DeleteOptions) {
						opt.PropagationPolicy = &cascadeDelete
					}); err != nil {
						glog.Errorf("error querying api for workLoad object: %v, retrying...", err)
						return false, nil
					}
					return true, nil
				})
				glog.Info("Deleted workload object")
			}

			wait.PollImmediate(1*time.Second, e2e.WaitShort, func() (bool, error) {
				if err := client.Delete(context.TODO(), &machineAutoscaler); err != nil {
					glog.Errorf("error querying api for machineAutoscaler object: %v, retrying...", err)
					return false, nil
				}
				return true, nil
			})
			glog.Info("Deleted machineAutoscaler object")

			wait.PollImmediate(1*time.Second, e2e.WaitShort, func() (bool, error) {
				if err := client.Delete(context.TODO(), &clusterAutoscaler); err != nil {
					glog.Errorf("error querying api for clusterAutoscaler object: %v, retrying...", err)
					return false, nil
				}
				return true, nil
			})
			glog.Info("Deleted clusterAutoscaler object")
		}()

		nodeTestLabel := fmt.Sprintf("machine.openshift.io/autoscaling-test-%v", string(uuid.NewUUID()))

		// Label all nodes belonging to the machineset (before scale up phase)
		err = labelMachineSetNodes(client, &targetMachineSet, nodeTestLabel)
		o.Expect(err).NotTo(o.HaveOccurred())

		glog.Info("Get nodeList")
		nodeList := corev1.NodeList{}
		err = wait.PollImmediate(1*time.Second, e2e.WaitMedium, func() (bool, error) {
			if err := client.List(context.TODO(), runtimeclient.MatchingLabels(map[string]string{nodeTestLabel: ""}), &nodeList); err != nil {
				glog.Errorf("error querying api for nodeList object: %v, retrying...", err)
				return false, err
			}
			return true, nil
		})
		o.Expect(err).NotTo(o.HaveOccurred())

		nodeGroupInitialTotalNodes := len(nodeList.Items)
		glog.Infof("Cluster initial number of nodes in node group %v is %d", targetMachineSet.Name, nodeGroupInitialTotalNodes)

		glog.Info("Create workload")

		err = wait.PollImmediate(1*time.Second, e2e.WaitMedium, func() (bool, error) {
			if err := client.Create(context.TODO(), workLoad); err != nil {
				glog.Errorf("error querying api for workLoad object: %v, retrying...", err)
				return false, err
			}
			return true, nil
		})
		o.Expect(err).NotTo(o.HaveOccurred())

		glog.Info("Waiting for cluster to scale out number of replicas")
		err = wait.PollImmediate(5*time.Second, e2e.WaitLong, func() (bool, error) {
			msKey := types.NamespacedName{
				Namespace: e2e.TestContext.MachineApiNamespace,
				Name:      targetMachineSet.Name,
			}
			ms := &mapiv1beta1.MachineSet{}
			if err := client.Get(context.TODO(), msKey, ms); err != nil {
				glog.Errorf("error querying api for clusterAutoscaler object: %v, retrying...", err)
				return false, nil
			}
			glog.Infof("MachineSet %s. Initial number of replicas: %d. Current number of replicas: %d", targetMachineSet.Name, *initialNumberOfReplicas, *ms.Spec.Replicas)
			return *ms.Spec.Replicas > *initialNumberOfReplicas, nil
		})
		o.Expect(err).NotTo(o.HaveOccurred())

		glog.Info("Waiting for cluster to scale up nodes")
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
				nodeCounter++
			}

			glog.Infof("Expecting at least one new node to come up. Initial number of node group nodes: %d. Current number of nodes in the group: %d", nodeGroupInitialTotalNodes, nodeCounter)
			return nodeCounter > nodeGroupInitialTotalNodes, nil
		})
		o.Expect(err).NotTo(o.HaveOccurred())

		// Label all nodes belonging to the machineset (after scale up phase)
		err = labelMachineSetNodes(client, &targetMachineSet, nodeTestLabel)
		o.Expect(err).NotTo(o.HaveOccurred())

		glog.Info("Delete workload")
		err = wait.PollImmediate(5*time.Second, e2e.WaitMedium, func() (bool, error) {
			cascadeDelete := metav1.DeletePropagationForeground
			if err := client.Delete(context.TODO(), workLoad, func(opt *runtimeclient.DeleteOptions) {
				opt.PropagationPolicy = &cascadeDelete
			}); err != nil {
				glog.Errorf("error querying api for workLoad object: %v, retrying...", err)
				return false, nil
			}
			workLoad = nil
			return true, nil
		})
		o.Expect(err).NotTo(o.HaveOccurred())

		// As we have just deleted the workload the autoscaler will
		// start to scale down the unneeded nodes. We wait for that
		// condition; if successful we assert that (a smoke test of)
		// scale down is functional.
		glog.Info("Wait for cluster to have at most initial number of replicas")
		err = wait.PollImmediate(5*time.Second, e2e.WaitLong, func() (bool, error) {
			msKey := types.NamespacedName{
				Namespace: e2e.TestContext.MachineApiNamespace,
				Name:      targetMachineSet.Name,
			}
			ms := &mapiv1beta1.MachineSet{}
			if err := client.Get(context.TODO(), msKey, ms); err != nil {
				glog.Errorf("error querying api for machineSet object: %v, retrying...", err)
				return false, nil
			}
			glog.Infof("Initial number of replicas: %d. Current number of replicas: %d", *initialNumberOfReplicas, *ms.Spec.Replicas)
			if *ms.Spec.Replicas > *initialNumberOfReplicas {
				return false, nil
			}

			// Make sure all scaled down nodes are really gove (so they don't affect tests run after this one)
			scaledNodes := corev1.NodeList{}
			if err := client.List(context.TODO(), runtimeclient.MatchingLabels(map[string]string{nodeTestLabel: ""}), &scaledNodes); err != nil {
				glog.Errorf("Error querying api for node objects: %v, retrying...", err)
				return false, nil
			}
			scaledNodesLen := int32(len(scaledNodes.Items))
			glog.Infof("Current number of replicas: %d. Current number of nodes: %d", *ms.Spec.Replicas, scaledNodesLen)
			// get all linked nodes (so we can wait later on their deletion)
			return scaledNodesLen <= *ms.Spec.Replicas && scaledNodesLen <= *initialNumberOfReplicas, nil
		})
		o.Expect(err).NotTo(o.HaveOccurred())
	})

})
