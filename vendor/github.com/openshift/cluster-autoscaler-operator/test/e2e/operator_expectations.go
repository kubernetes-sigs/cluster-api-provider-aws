package main

import (
	"fmt"
	"time"

	"context"

	"github.com/golang/glog"
	autoscalingv1alpha1 "github.com/openshift/cluster-autoscaler-operator/pkg/apis/autoscaling/v1alpha1"
	kappsapi "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
)

var (
	PodPriorityThreshold int32 = -10
	MaxPodGracePeriod    int32 = 60
	MaxNodesTotal        int32 = 100
	CoresMin             int32 = 16
	CoresMax             int32 = 32
	MemoryMin            int32 = 32
	MemoryMax            int32 = 64
	NvidiaGPUMin         int32 = 4
	NvidiaGPUMax         int32 = 8
)

func CreateClusterAutoscaler() error {
	ca := &autoscalingv1alpha1.ClusterAutoscaler{
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

	return F.Client.Create(context.TODO(), ca)
}

func ExpectOperatorAvailable() error {
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

func ExpectClusterAutoscalerAvailable() error {
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
