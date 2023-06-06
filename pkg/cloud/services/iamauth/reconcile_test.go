/*
Copyright 2020 The Kubernetes Authors.

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

package iamauth

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	expclusterv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
)

func TestReconcileIAMAuth(t *testing.T) {
	var (
		mockCtrl *gomock.Controller
		ctx      context.Context
	)
	setup := func(t *testing.T) {
		t.Helper()
		mockCtrl = gomock.NewController(t)
		ctx = context.TODO()
	}

	teardown := func() {
		mockCtrl.Finish()
	}
	t.Run("Should successfully find roles for MachineDeployments and MachinePools", func(t *testing.T) {
		g := NewWithT(t)
		setup(t)
		namespace, err := testEnv.CreateNamespace(ctx, fmt.Sprintf("integ-test-%s", util.RandomString(5)))
		g.Expect(err).To(BeNil())
		ns := namespace.Name
		name := "default"
		eksCluster := createEKSCluster(name, ns)
		g.Expect(testEnv.Create(ctx, eksCluster)).To(Succeed())
		awsMP := createAWSMachinePoolForClusterWithInstanceProfile(name, ns, eksCluster.Name, "nodes.cluster-api-provider-aws.sigs.k8s.io")
		infraRef := corev1.ObjectReference{
			Kind:       awsMP.TypeMeta.Kind,
			Name:       awsMP.Name,
			Namespace:  awsMP.Namespace,
			APIVersion: awsMP.TypeMeta.APIVersion,
		}
		g.Expect(testEnv.Create(ctx, awsMP)).To(Succeed())
		mp := createMachinepoolForCluster(name, ns, eksCluster.Name, infraRef)
		g.Expect(testEnv.Create(ctx, mp)).To(Succeed())

		awsMachineTemplate := createAWSMachineTemplateForClusterWithInstanceProfile(name, ns, eksCluster.Name, "eks-nodes.cluster-api-provider-aws.sigs.k8s.io")
		infraRefForMD := corev1.ObjectReference{
			Kind:       awsMachineTemplate.TypeMeta.Kind,
			Name:       awsMachineTemplate.Name,
			Namespace:  awsMachineTemplate.Namespace,
			APIVersion: awsMachineTemplate.TypeMeta.APIVersion,
		}
		g.Expect(testEnv.Create(ctx, awsMachineTemplate)).To(Succeed())
		md := createMachineDeploymentForCluster(name, ns, eksCluster.Name, infraRefForMD)
		g.Expect(testEnv.Create(ctx, md)).To(Succeed())

		expectedRoles := map[string]struct{}{
			"nodes.cluster-api-provider-aws.sigs.k8s.io":     {},
			"eks-nodes.cluster-api-provider-aws.sigs.k8s.io": {},
		}

		controllerIdentity := createControllerIdentity()
		g.Expect(testEnv.Create(ctx, controllerIdentity)).To(Succeed())
		managedScope, err := scope.NewManagedControlPlaneScope(scope.ManagedControlPlaneScopeParams{
			Client:       testEnv,
			ControlPlane: eksCluster,
			Cluster: &clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: ns,
				},
			},
		})
		g.Expect(err).To(BeNil(), "failed to create managedScope")
		authService := NewService(managedScope, BackendTypeConfigMap, managedScope.Client)
		gotRoles, err := authService.getRolesForWorkers(ctx)
		g.Expect(err).To(BeNil(), "failed to get roles for workers")
		g.Expect(gotRoles).To(BeEquivalentTo(expectedRoles), "did not get correct roles for workers")
		defer teardown()
		defer t.Cleanup(func() {
			g.Expect(testEnv.Cleanup(ctx, namespace, eksCluster, awsMP, mp, awsMachineTemplate, md, controllerIdentity)).To(Succeed())
		})
	})
}

func createEKSCluster(name, namespace string) *ekscontrolplanev1.AWSManagedControlPlane {
	eksCluster := &ekscontrolplanev1.AWSManagedControlPlane{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				clusterv1.ClusterNameLabel: name,
			},
		},
		Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{},
	}
	return eksCluster
}

func createAWSMachinePoolForClusterWithInstanceProfile(name, namespace, clusterName, instanceProfile string) *expinfrav1.AWSMachinePool {
	awsMP := &expinfrav1.AWSMachinePool{
		TypeMeta: metav1.TypeMeta{
			Kind: "AWSMachinePool",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				clusterv1.ClusterNameLabel: clusterName,
			},
		},
		Spec: expinfrav1.AWSMachinePoolSpec{
			AWSLaunchTemplate: expinfrav1.AWSLaunchTemplate{
				IamInstanceProfile: instanceProfile,
			},
			MaxSize: 1,
		},
	}
	return awsMP
}

func createMachinepoolForCluster(name, namespace, clusterName string, infrastructureRef corev1.ObjectReference) *expclusterv1.MachinePool {
	mp := &expclusterv1.MachinePool{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				clusterv1.ClusterNameLabel: clusterName,
			},
		},
		Spec: expclusterv1.MachinePoolSpec{
			ClusterName: clusterName,
			Template: clusterv1.MachineTemplateSpec{
				Spec: clusterv1.MachineSpec{
					ClusterName:       clusterName,
					InfrastructureRef: infrastructureRef,
				},
			},
		},
	}
	return mp
}

func createAWSMachineTemplateForClusterWithInstanceProfile(name, namespace, clusterName, instanceProfile string) *infrav1.AWSMachineTemplate {
	mt := &infrav1.AWSMachineTemplate{
		TypeMeta: metav1.TypeMeta{
			Kind: "AWSMachineTemplate",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				clusterv1.ClusterNameLabel: clusterName,
			},
		},
		Spec: infrav1.AWSMachineTemplateSpec{
			Template: infrav1.AWSMachineTemplateResource{
				Spec: infrav1.AWSMachineSpec{
					IAMInstanceProfile: instanceProfile,
					InstanceType:       "m5.xlarge",
				},
			},
		},
	}
	return mt
}

func createMachineDeploymentForCluster(name, namespace, clusterName string, infrastructureRef corev1.ObjectReference) *clusterv1.MachineDeployment {
	md := &clusterv1.MachineDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				clusterv1.ClusterNameLabel: clusterName,
			},
		},
		Spec: clusterv1.MachineDeploymentSpec{
			ClusterName: clusterName,
			Template: clusterv1.MachineTemplateSpec{
				Spec: clusterv1.MachineSpec{
					ClusterName:       clusterName,
					InfrastructureRef: infrastructureRef,
				},
			},
			Replicas: pointer.Int32(2),
		},
	}
	return md
}

func createControllerIdentity() *infrav1.AWSClusterControllerIdentity {
	controllerIdentity := &infrav1.AWSClusterControllerIdentity{
		TypeMeta: metav1.TypeMeta{
			Kind: string(infrav1.ControllerIdentityKind),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
		},
		Spec: infrav1.AWSClusterControllerIdentitySpec{
			AWSClusterIdentitySpec: infrav1.AWSClusterIdentitySpec{
				AllowedNamespaces: &infrav1.AllowedNamespaces{},
			},
		},
	}
	return controllerIdentity
}
