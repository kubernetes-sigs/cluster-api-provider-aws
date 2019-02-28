package main

import (
	"fmt"
	"github.com/openshift/cluster-autoscaler-operator/pkg/operator"
	"time"

	"context"

	"github.com/golang/glog"
	osconfigv1 "github.com/openshift/api/config/v1"
	autoscalingv1alpha1 "github.com/openshift/cluster-autoscaler-operator/pkg/apis/autoscaling/v1alpha1"
	cvorm "github.com/openshift/cluster-version-operator/lib/resourcemerge"
	kappsapi "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const nodeGroupSize = 6

func clusterAutoscalerResource() *autoscalingv1alpha1.ClusterAutoscaler {
	var (
		PodPriorityThreshold int32 = -10
		MaxPodGracePeriod    int32 = 60
		MaxNodesTotal        int32 = 100
		CoresMin             int32 = 16
		CoresMax             int32 = 32
		MemoryMin            int32 = 32
		MemoryMax            int32 = 64
	)

	return &autoscalingv1alpha1.ClusterAutoscaler{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterAutoscaler",
			APIVersion: "autoscaling.openshift.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      caName,
			Namespace: namespace,
		},
		Spec: autoscalingv1alpha1.ClusterAutoscalerSpec{
			MaxPodGracePeriod:    &MaxPodGracePeriod,
			PodPriorityThreshold: &PodPriorityThreshold,
			ResourceLimits: &autoscalingv1alpha1.ResourceLimits{
				MaxNodesTotal: &MaxNodesTotal,
				Cores: &autoscalingv1alpha1.ResourceRange{
					Min: CoresMin,
					Max: CoresMax,
				},
				Memory: &autoscalingv1alpha1.ResourceRange{
					Min: MemoryMin,
					Max: MemoryMax,
				},
			},
		},
	}
}

func machineAutoscalerResource() *autoscalingv1alpha1.MachineAutoscaler {
	return &autoscalingv1alpha1.MachineAutoscaler{
		TypeMeta: metav1.TypeMeta{
			Kind:       "MachineAutoscaler",
			APIVersion: "autoscaling.openshift.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "compute-pool",
			Namespace: namespace,
		},
		Spec: autoscalingv1alpha1.MachineAutoscalerSpec{
			MinReplicas: 1,
			MaxReplicas: nodeGroupSize,
			ScaleTargetRef: autoscalingv1alpha1.CrossVersionObjectReference{
				Kind:       "MachineSet",
				APIVersion: "machine.openshift.io/v1beta1",
				Name:       "kubemark-machine-pool",
			},
		},
	}
}

func workloadResource() *batchv1.Job {
	var backoffLimit int32 = 4
	var completions int32 = 40
	var parallelism int32 = 40
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "workload",
			Namespace: "default",
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "work",
							Image: "busybox",
							Command: []string{
								"sleep",
								"120",
							},
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									"cpu":    resource.MustParse("500m"),
									"memory": resource.MustParse("500Mi"),
								},
							},
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
					NodeSelector: map[string]string{
						"node-role.kubernetes.io/compute": "",
					},
					Tolerations: []v1.Toleration{
						{
							Key:      "kubemark",
							Operator: v1.TolerationOpExists,
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

func verifyClusterAutoscalerOperatorIsReady() error {
	name := "cluster-autoscaler-operator"
	key := types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}
	d := &kappsapi.Deployment{}

	err := wait.PollImmediate(1*time.Second, 1*time.Minute, func() (bool, error) {
		if err := F.Client.Get(context.TODO(), key, d); err != nil {
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

func verifyClusterAutoscalerIsReady() error {
	name := fmt.Sprintf("cluster-autoscaler-%s", caName)
	key := types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}
	d := &kappsapi.Deployment{}

	err := wait.PollImmediate(1*time.Second, 1*time.Minute, func() (bool, error) {
		if err := F.Client.Get(context.TODO(), key, d); err != nil {
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

func ExpectOperatorAvailable() error {
	return verifyClusterAutoscalerOperatorIsReady()
}

func ExpectClusterOperatorStatusAvailable() error {
	name := operator.OperatorName
	key := types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}
	clusterOperator := &osconfigv1.ClusterOperator{}

	err := wait.PollImmediate(1*time.Second, 1*time.Minute, func() (bool, error) {
		if err := F.Client.Get(context.TODO(), key, clusterOperator); err != nil {
			glog.Errorf("error querying api for OperatorStatus object: %v, retrying...", err)
			return false, nil
		}
		if cvorm.IsOperatorStatusConditionTrue(clusterOperator.Status.Conditions, osconfigv1.OperatorAvailable) {
			return true, nil
		}
		return false, nil
	})
	return err
}

func CreateClusterAutoscaler() error {
	ca := clusterAutoscalerResource()
	return F.Client.Create(context.TODO(), ca)
}

func ExpectClusterAutoscalerAvailable() error {
	return verifyClusterAutoscalerIsReady()
}

func isNodeReady(node *v1.Node) bool {
	for _, c := range node.Status.Conditions {
		if c.Type == v1.NodeReady {
			return c.Status == v1.ConditionTrue
		}
	}
	return false
}

// ExpectToScaleUpAndDown assumes the cluster autoscaler is already deployed
func ExpectToScaleUpAndDown() error {
	tenSeconds := "10s"

	ca := *clusterAutoscalerResource()
	// change the min core and memory to 1 so the autoscaler scales down as much as it can
	ca.Spec.ResourceLimits.Cores.Min = 1
	ca.Spec.ResourceLimits.Memory.Min = 1
	ca.Spec.ScaleDown = &autoscalingv1alpha1.ScaleDownConfig{
		Enabled:           true,
		DelayAfterAdd:     &tenSeconds,
		DelayAfterDelete:  &tenSeconds,
		DelayAfterFailure: &tenSeconds,
		UnneededTime:      &tenSeconds,
	}
	if err := F.Client.Create(context.TODO(), &ca); err != nil {
		return fmt.Errorf("unable to deploy ClusterAutoscaler resource: %v", err)
	}

	if err := verifyClusterAutoscalerIsReady(); err != nil {
		return fmt.Errorf("cluster autoscaler not ready: %v", err)
	}

	ma := machineAutoscalerResource()
	if err := F.Client.Create(context.TODO(), ma); err != nil {
		return fmt.Errorf("unable to deploy MachineAutoscaler resource: %v", err)
	}

	workload := workloadResource()
	if err := F.Client.Create(context.TODO(), workload); err != nil {
		return fmt.Errorf("unable to deploy workload resource: %v", err)
	}

	err := wait.PollImmediate(1*time.Second, 1*time.Minute, func() (bool, error) {
		nodes := v1.NodeList{}
		if err := F.Client.List(context.TODO(), &client.ListOptions{}, &nodes); err != nil {
			glog.Errorf("error querying api for Node object: %v, retrying...", err)
			return false, nil
		}
		// expecting nodeGroupSize nodes
		nodeCounter := 0
		for _, node := range nodes.Items {
			if _, exists := node.Labels["node-role.kubernetes.io/compute"]; !exists {
				continue
			}

			if !isNodeReady(&node) {
				continue
			}

			nodeCounter++
		}

		if nodeCounter < 6 {
			glog.Errorf("Expecting %v nodes with 'node-role.kubernetes.io/compute' label in Ready state, got %v", nodeGroupSize, nodeCounter)
			return false, nil
		}

		glog.Infof("Expected number (%v) of nodes with 'node-role.kubernetes.io/compute' label in Ready state found", nodeGroupSize)
		return true, nil
	})

	if err != nil {
		return err
	}

	glog.Infof("Succesfully scaled up")

	cascadeDelete := metav1.DeletePropagationForeground
	if err := F.Client.Delete(context.TODO(), workload, func(opt *client.DeleteOptions) { opt.PropagationPolicy = &cascadeDelete }); err != nil {
		return fmt.Errorf("unable to delete workload resource: %v", err)
	}

	if err = wait.PollImmediate(1*time.Second, 10*time.Minute, func() (bool, error) {
		nodes := v1.NodeList{}
		if err := F.Client.List(context.TODO(), &client.ListOptions{}, &nodes); err != nil {
			glog.Errorf("error querying api for Node object: %v, retrying...", err)
			return false, nil
		}
		// expecting nodeGroupSize nodes
		nodeCounter := 0
		for _, node := range nodes.Items {
			if _, exists := node.Labels["node-role.kubernetes.io/compute"]; !exists {
				continue
			}

			if !isNodeReady(&node) {
				continue
			}

			nodeCounter++
		}

		if nodeCounter > 1 {
			glog.Errorf("Expecting 1 node with 'node-role.kubernetes.io/compute' label in Ready state, got %v", nodeCounter)
			return false, nil
		}

		glog.Info("Expected number (1) of nodes with 'node-role.kubernetes.io/compute' label in Ready state found")
		return true, nil
	}); err != nil {
		return err
	}

	glog.Infof("Succesfully scaled down")

	return nil
}
