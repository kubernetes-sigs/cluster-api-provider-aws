package suites

import (
	mapiv1beta1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	kpolicyapi "k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/uuid"
)

var nodeDrainLabels = map[string]string{
	workerRoleLabel:      "",
	"node-draining-test": "",
}

func newWorkLoad() *batchv1.Job {
	backoffLimit := int32(4)
	completions := int32(50)
	parallelism := int32(50)
	activeDeadlineSeconds := int64(100)
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "workload",
			Namespace: namespace,
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
				},
			},
			ActiveDeadlineSeconds: &activeDeadlineSeconds,
			BackoffLimit:          &backoffLimit,
			Completions:           &completions,
			Parallelism:           &parallelism,
		},
	}
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

// 1. create two machines (without machineset) and wait until nodes are registered and ready
// 1. create rc
// 1. create pdb
// 1. pick a node that has at least half of the rc pods
// 1. drain node
// 1. observe the machine object is not deleted before a node is drained,
//    i.e. as long as there is at least one pod running on the drained node,
//    the machine object can not be deleted

func isNodeReady(node *corev1.Node) bool {
	for _, c := range node.Status.Conditions {
		if c.Type == corev1.NodeReady {
			return c.Status == corev1.ConditionTrue
		}
	}
	return false
}
