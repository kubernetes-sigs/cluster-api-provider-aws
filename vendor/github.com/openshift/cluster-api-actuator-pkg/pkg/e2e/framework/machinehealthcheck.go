package framework

import (
	"context"

	"github.com/ghodss/yaml"
	healthcheckingv1alpha1 "github.com/openshift/machine-api-operator/pkg/apis/healthchecking/v1alpha1"
	"github.com/openshift/machine-api-operator/pkg/util/conditions"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/rand"
)

const (
	// KubeletKillerPodName contains the name of the pod that stops kubelet process
	KubeletKillerPodName = "kubelet-killer"
	// MachineHealthCheckName contains the name of the machinehealthcheck used for tests
	MachineHealthCheckName = "workers-check"
)

// CreateUnhealthyConditionsConfigMap creates node-unhealthy-conditions configmap with relevant conditions
func CreateUnhealthyConditionsConfigMap(unhealthyConditions *conditions.UnhealthyConditions) error {
	client, err := LoadClient()
	if err != nil {
		return err
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: TestContext.MachineApiNamespace,
			Name:      healthcheckingv1alpha1.ConfigMapNodeUnhealthyConditions,
		},
	}

	conditionsData, err := yaml.Marshal(unhealthyConditions)
	if err != nil {
		return err
	}

	cm.Data = map[string]string{"conditions": string(conditionsData)}
	return client.Create(context.TODO(), cm)
}

// DeleteUnhealthyConditionsConfigMap deletes node-unhealthy-conditions configmap
func DeleteUnhealthyConditionsConfigMap() error {
	client, err := LoadClient()
	if err != nil {
		return err
	}

	key := types.NamespacedName{
		Name:      healthcheckingv1alpha1.ConfigMapNodeUnhealthyConditions,
		Namespace: TestContext.MachineApiNamespace,
	}
	cm := &corev1.ConfigMap{}
	err = client.Get(context.TODO(), key, cm)
	if err != nil {
		return err
	}

	return client.Delete(context.TODO(), cm)
}

// CreateMachineHealthCheck will create MachineHealthCheck CR with the relevant selector
func CreateMachineHealthCheck(labels map[string]string) error {
	client, err := LoadClient()
	if err != nil {
		return err
	}

	mhc := &healthcheckingv1alpha1.MachineHealthCheck{
		ObjectMeta: metav1.ObjectMeta{
			Name:      MachineHealthCheckName,
			Namespace: TestContext.MachineApiNamespace,
		},
		Spec: healthcheckingv1alpha1.MachineHealthCheckSpec{
			Selector: metav1.LabelSelector{
				MatchLabels: labels,
			},
		},
	}
	return client.Create(context.TODO(), mhc)
}

// StopKubelet creates pod in the node PID namespace that stops kubelet process
func StopKubelet(nodeName string) error {
	client, err := LoadClient()
	if err != nil {
		return err
	}

	_true := true
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      KubeletKillerPodName + rand.String(5),
			Namespace: TestContext.MachineApiNamespace,
			Labels: map[string]string{
				KubeletKillerPodName: "",
			},
		},
		Spec: corev1.PodSpec{
			RestartPolicy: corev1.RestartPolicyNever,
			Containers: []corev1.Container{
				{
					Name:    KubeletKillerPodName,
					Image:   "busybox",
					Command: []string{"pkill", "-STOP", "hyperkube"},
					SecurityContext: &corev1.SecurityContext{
						Privileged: &_true,
					},
				},
			},
			NodeName: nodeName,
			HostPID:  true,
		},
	}
	return client.Create(context.TODO(), pod)
}
