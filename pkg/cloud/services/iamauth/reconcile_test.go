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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
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
		infraRef := clusterv1.ContractVersionedObjectReference{
			Kind:     awsMP.TypeMeta.Kind,
			Name:     awsMP.Name,
			APIGroup: awsMP.TypeMeta.GroupVersionKind().Group,
		}
		configRef := clusterv1.ContractVersionedObjectReference{
			Kind:     "EKSConfig",
			Name:     awsMP.Name,
			APIGroup: awsMP.TypeMeta.GroupVersionKind().Group,
		}
		g.Expect(testEnv.Create(ctx, awsMP)).To(Succeed())
		mp := createMachinepoolForCluster(name, ns, eksCluster.Name, infraRef, configRef)
		g.Expect(testEnv.Create(ctx, mp)).To(Succeed())

		awsMachineTemplate := createAWSMachineTemplateForClusterWithInstanceProfile(name, ns, eksCluster.Name, "eks-nodes.cluster-api-provider-aws.sigs.k8s.io")
		infraRefForMD := clusterv1.ContractVersionedObjectReference{
			Kind:     awsMachineTemplate.TypeMeta.Kind,
			Name:     awsMachineTemplate.Name,
			APIGroup: awsMachineTemplate.TypeMeta.GroupVersionKind().Group,
		}
		configRefForMD := clusterv1.ContractVersionedObjectReference{
			Kind:     "EKSConfig",
			Name:     awsMachineTemplate.Name,
			APIGroup: awsMachineTemplate.TypeMeta.GroupVersionKind().Group,
		}
		g.Expect(testEnv.Create(ctx, awsMachineTemplate)).To(Succeed())
		md := createMachineDeploymentForCluster(name, ns, eksCluster.Name, infraRefForMD, configRefForMD)
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
			Kind:       "AWSMachinePool",
			APIVersion: "infrastructure.cluster.x-k8s.io/v1beta2",
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

func createMachinepoolForCluster(name, namespace, clusterName string, infrastructureRef, configRef clusterv1.ContractVersionedObjectReference) *clusterv1.MachinePool {
	mp := &clusterv1.MachinePool{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				clusterv1.ClusterNameLabel: clusterName,
			},
		},
		Spec: clusterv1.MachinePoolSpec{
			ClusterName: clusterName,
			Template: clusterv1.MachineTemplateSpec{
				Spec: clusterv1.MachineSpec{
					ClusterName:       clusterName,
					InfrastructureRef: infrastructureRef,
					Bootstrap: clusterv1.Bootstrap{
						ConfigRef: configRef,
					},
				},
			},
		},
	}
	return mp
}

func createAWSMachineTemplateForClusterWithInstanceProfile(name, namespace, clusterName, instanceProfile string) *infrav1.AWSMachineTemplate {
	mt := &infrav1.AWSMachineTemplate{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AWSMachineTemplate",
			APIVersion: "bootstrap.cluster.x-k8s.io/v1beta2",
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

func createMachineDeploymentForCluster(name, namespace, clusterName string, infrastructureRef, configRefForMD clusterv1.ContractVersionedObjectReference) *clusterv1.MachineDeployment {
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
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "test-app"},
			},
			Template: clusterv1.MachineTemplateSpec{
				Spec: clusterv1.MachineSpec{
					ClusterName:       clusterName,
					InfrastructureRef: infrastructureRef,
					Bootstrap: clusterv1.Bootstrap{
						ConfigRef: configRefForMD,
					},
				},
			},
			Replicas: ptr.To[int32](2),
		},
	}
	return md
}

// TestDedupAndSortRoles pins two properties that keep ReconcileMappings
// deterministic on repeated reconciles:
//   - duplicate RoleARNs collapse to a single entry, with later entries
//     (user-configured mappings) winning over earlier ones (node-role
//     discovery), and
//   - output is sorted by RoleARN so aws-auth ConfigMap key order and CRD
//     backend Create/Delete order do not depend on nodeRoles map iteration.
func TestDedupAndSortRoles(t *testing.T) {
	g := NewWithT(t)

	nodeRole := ekscontrolplanev1.RoleMapping{
		RoleARN: "arn:aws:iam::000000000000:role/KubernetesNode",
		KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
			UserName: "system:node:{{EC2PrivateDNSName}}",
			Groups:   []string{"system:bootstrappers", "system:nodes"},
		},
	}
	adminRole := ekscontrolplanev1.RoleMapping{
		RoleARN: "arn:aws:iam::000000000000:role/KubernetesAdmin",
		KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
			UserName: "admin:{{SessionName}}",
			Groups:   []string{"system:masters"},
		},
	}
	// nodeRoleOverride shares an ARN with nodeRole but sets a different
	// UserName/Groups — simulates a user-configured mapping that collides with
	// the discovered node role. The user-configured entry (appended later)
	// must win.
	nodeRoleOverride := ekscontrolplanev1.RoleMapping{
		RoleARN: nodeRole.RoleARN,
		KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
			UserName: "custom-node-user",
			Groups:   []string{"system:masters"},
		},
	}

	got := dedupAndSortRoles([]ekscontrolplanev1.RoleMapping{nodeRole, adminRole, nodeRoleOverride})

	g.Expect(got).To(HaveLen(2))
	g.Expect(got[0].RoleARN).To(Equal(adminRole.RoleARN), "output must be sorted by RoleARN ascending")
	g.Expect(got[1].RoleARN).To(Equal(nodeRole.RoleARN))
	g.Expect(got[1].UserName).To(Equal("custom-node-user"), "later duplicate entry (user-configured) must override earlier one (node-role discovery)")
	g.Expect(got[1].Groups).To(Equal([]string{"system:masters"}))
}

// TestDedupAndSortUsers mirrors TestDedupAndSortRoles for UserMappings.
func TestDedupAndSortUsers(t *testing.T) {
	g := NewWithT(t)

	alice := ekscontrolplanev1.UserMapping{
		UserARN: "arn:aws:iam::000000000000:user/Alice",
		KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
			UserName: "alice",
			Groups:   []string{"system:masters"},
		},
	}
	bob := ekscontrolplanev1.UserMapping{
		UserARN: "arn:aws:iam::000000000000:user/Bob",
		KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
			UserName: "bob",
			Groups:   []string{"viewers"},
		},
	}
	aliceUpdated := ekscontrolplanev1.UserMapping{
		UserARN: alice.UserARN,
		KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
			UserName: "alice-updated",
			Groups:   []string{"editors"},
		},
	}

	got := dedupAndSortUsers([]ekscontrolplanev1.UserMapping{bob, alice, aliceUpdated})

	g.Expect(got).To(HaveLen(2))
	g.Expect(got[0].UserARN).To(Equal(alice.UserARN), "output must be sorted by UserARN ascending")
	g.Expect(got[0].UserName).To(Equal("alice-updated"), "later duplicate entry must override earlier one")
	g.Expect(got[1].UserARN).To(Equal(bob.UserARN))
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
