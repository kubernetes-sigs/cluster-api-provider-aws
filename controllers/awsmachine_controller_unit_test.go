/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2/klogr"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/handler"
)

func newMachine(clusterName, machineName string) *clusterv1.Machine {
	return &clusterv1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				clusterv1.ClusterLabelName: clusterName,
			},
			Name:      machineName,
			Namespace: "default",
		},
	}
}

func newMachineWithInfrastructureRef(clusterName, machineName string) *clusterv1.Machine {
	m := newMachine(clusterName, machineName)
	m.Spec.InfrastructureRef = v1.ObjectReference{
		Kind:       "AWSMachine",
		Namespace:  "",
		Name:       "aws" + machineName,
		APIVersion: infrav1.GroupVersion.String(),
	}
	return m
}

func newCluster(name string) *clusterv1.Cluster {
	return &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
	}
}

func TestAWSMachineReconciler_AWSClusterToAWSMachines(t *testing.T) {
	clusterName := "my-cluster"
	ctx := context.TODO()
	g := NewWithT(t)
	g.Expect(testEnv.Create(ctx, newCluster(clusterName))).To(Succeed())
	g.Expect(testEnv.Create(ctx, newMachineWithInfrastructureRef(clusterName, "my-machine-0"))).To(Succeed())
	g.Expect(testEnv.Create(ctx, newMachineWithInfrastructureRef(clusterName, "my-machine-1"))).To(Succeed())
	g.Expect(testEnv.Create(ctx, newMachine(clusterName, "my-machine-2"))).To(Succeed())

	reconciler := &AWSMachineReconciler{
		Client: testEnv.Client,
		Log:    klogr.New(),
	}

	g.Expect(reconciler.AWSClusterToAWSMachines(handler.MapObject{
		Object: &infrav1.AWSCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:      clusterName,
				Namespace: "default",
				OwnerReferences: []metav1.OwnerReference{
					{
						Name:       clusterName,
						Kind:       "Cluster",
						APIVersion: clusterv1.GroupVersion.String(),
					},
				},
			},
		},
	})).Should(HaveLen(2))
}
