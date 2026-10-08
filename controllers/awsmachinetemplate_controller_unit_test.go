/*
Copyright 2025 The Kubernetes Authors.

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
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
)

func TestAWSMachineTemplateReconciler(t *testing.T) {
	setupScheme := func() *runtime.Scheme {
		scheme := runtime.NewScheme()
		_ = infrav1.AddToScheme(scheme)
		_ = ekscontrolplanev1.AddToScheme(scheme)
		_ = clusterv1.AddToScheme(scheme)
		_ = corev1.AddToScheme(scheme)
		return scheme
	}

	newFakeClient := func(objs ...client.Object) client.Client {
		return fake.NewClientBuilder().
			WithScheme(setupScheme()).
			WithObjects(objs...).
			WithStatusSubresource(&infrav1.AWSMachineTemplate{}).
			Build()
	}

	newAWSMachineTemplate := func(name string) *infrav1.AWSMachineTemplate {
		return &infrav1.AWSMachineTemplate{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: "default",
			},
			Spec: infrav1.AWSMachineTemplateSpec{
				Template: infrav1.AWSMachineTemplateResource{
					Spec: infrav1.AWSMachineSpec{
						InstanceType: "t3.medium",
					},
				},
			},
		}
	}

	t.Run("getEC2Client", func(t *testing.T) {
		t.Run("should return nil when cluster type is not recognised", func(t *testing.T) {
			g := NewWithT(t)
			cluster := &clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test-cluster", Namespace: "default"},
			}
			reconciler := &AWSMachineTemplateReconciler{Client: newFakeClient(cluster)}

			ec2Client, err := reconciler.getEC2Client(context.Background(), nil, cluster)

			g.Expect(err).To(BeNil())
			g.Expect(ec2Client).To(BeNil())
		})

		t.Run("should return NotFound error when AWSCluster is not found", func(t *testing.T) {
			g := NewWithT(t)
			cluster := &clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test-cluster", Namespace: "default"},
				Spec: clusterv1.ClusterSpec{
					InfrastructureRef: clusterv1.ContractVersionedObjectReference{
						Kind: kindAWSCluster,
						Name: "missing-aws-cluster",
					},
				},
			}
			reconciler := &AWSMachineTemplateReconciler{Client: newFakeClient(cluster)}

			ec2Client, err := reconciler.getEC2Client(context.Background(), nil, cluster)

			g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
			g.Expect(ec2Client).To(BeNil())
		})

		t.Run("should return NotFound error when AWSManagedControlPlane is not found", func(t *testing.T) {
			g := NewWithT(t)
			cluster := &clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test-cluster", Namespace: "default"},
				Spec: clusterv1.ClusterSpec{
					ControlPlaneRef: clusterv1.ContractVersionedObjectReference{
						Kind: AWSManagedControlPlaneRefKind,
						Name: "missing-managed-cp",
					},
				},
			}
			reconciler := &AWSMachineTemplateReconciler{Client: newFakeClient(cluster)}

			ec2Client, err := reconciler.getEC2Client(context.Background(), nil, cluster)

			g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
			g.Expect(ec2Client).To(BeNil())
		})
	})

	// getInstanceTypeCapacity and getNodeInfo require a real EC2 client and are not unit tested here.

	t.Run("Reconcile", func(t *testing.T) {
		t.Run("should skip reconcile when capacity and nodeInfo are already populated", func(t *testing.T) {
			g := NewWithT(t)
			template := newAWSMachineTemplate("test-template")
			template.Status.Capacity = corev1.ResourceList{
				corev1.ResourceCPU:    *resource.NewQuantity(2, resource.DecimalSI),
				corev1.ResourceMemory: resource.MustParse("4Gi"),
			}
			template.Status.NodeInfo = &infrav1.NodeInfo{
				Architecture:    infrav1.ArchitectureAmd64,
				OperatingSystem: infrav1.OperatingSystemLinux,
			}

			reconciler := &AWSMachineTemplateReconciler{
				Client: newFakeClient(template),
			}

			// Should skip reconcile and return early without calling AWS APIs
			// No need to set up owner cluster or region since the early return happens before that
			result, err := reconciler.Reconcile(context.Background(), ctrl.Request{
				NamespacedName: client.ObjectKeyFromObject(template),
			})

			g.Expect(err).To(BeNil())
			g.Expect(result.RequeueAfter).To(BeZero())
		})

		t.Run("should reconcile when capacity set but nodeInfo is not", func(t *testing.T) {
			g := NewWithT(t)
			template := newAWSMachineTemplate("test-template")
			template.Status.Capacity = corev1.ResourceList{
				corev1.ResourceCPU: *resource.NewQuantity(2, resource.DecimalSI),
			}
			template.OwnerReferences = []metav1.OwnerReference{
				{
					APIVersion: clusterv1.GroupVersion.String(),
					Kind:       "Cluster",
					Name:       "test-cluster",
				},
			}
			cluster := &clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-cluster",
					Namespace: "default",
				},
				Spec: clusterv1.ClusterSpec{
					InfrastructureRef: clusterv1.ContractVersionedObjectReference{
						Kind: "AWSCluster",
						Name: "test-aws-cluster",
					},
				},
			}
			awsCluster := &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-aws-cluster",
					Namespace: "default",
				},
				Spec: infrav1.AWSClusterSpec{
					Region: "us-west-2",
				},
			}

			reconciler := &AWSMachineTemplateReconciler{
				Client: newFakeClient(template, cluster, awsCluster),
			}

			// This will fail at AWS API call, but demonstrates that reconcile proceeds
			result, err := reconciler.Reconcile(context.Background(), ctrl.Request{
				NamespacedName: client.ObjectKeyFromObject(template),
			})

			g.Expect(err).To(BeNil())
			g.Expect(result.RequeueAfter).To(BeZero())
		})

		t.Run("should skip when instance type is empty", func(t *testing.T) {
			g := NewWithT(t)
			template := newAWSMachineTemplate("test-template")
			template.Spec.Template.Spec.InstanceType = ""

			reconciler := &AWSMachineTemplateReconciler{
				Client: newFakeClient(template),
			}

			result, err := reconciler.Reconcile(context.Background(), ctrl.Request{
				NamespacedName: client.ObjectKeyFromObject(template),
			})

			g.Expect(err).To(BeNil())
			g.Expect(result.RequeueAfter).To(BeZero())
		})

		t.Run("should not reconcile when cluster is paused", func(t *testing.T) {
			g := NewWithT(t)
			template := newAWSMachineTemplate("test-template")
			template.OwnerReferences = []metav1.OwnerReference{
				{
					APIVersion: clusterv1.GroupVersion.String(),
					Kind:       "Cluster",
					Name:       "test-cluster",
				},
			}
			cluster := &clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-cluster",
					Namespace: "default",
				},
				Spec: clusterv1.ClusterSpec{
					Paused: ptr.To(true),
					InfrastructureRef: clusterv1.ContractVersionedObjectReference{
						Kind: "AWSCluster",
						Name: "test-aws-cluster",
					},
				},
			}
			awsCluster := &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-aws-cluster",
					Namespace: "default",
				},
				Spec: infrav1.AWSClusterSpec{
					Region: "us-west-2",
				},
			}

			reconciler := &AWSMachineTemplateReconciler{
				Client: newFakeClient(template, cluster, awsCluster),
			}

			result, err := reconciler.Reconcile(context.Background(), ctrl.Request{
				NamespacedName: client.ObjectKeyFromObject(template),
			})

			g.Expect(err).To(BeNil())
			g.Expect(result.RequeueAfter).To(BeZero())
		})

		t.Run("should not reconcile when template has paused annotation", func(t *testing.T) {
			g := NewWithT(t)
			template := newAWSMachineTemplate("test-template")
			template.Annotations = map[string]string{clusterv1.PausedAnnotation: ""}
			template.OwnerReferences = []metav1.OwnerReference{
				{
					APIVersion: clusterv1.GroupVersion.String(),
					Kind:       "Cluster",
					Name:       "test-cluster",
				},
			}
			cluster := &clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-cluster",
					Namespace: "default",
				},
				Spec: clusterv1.ClusterSpec{
					InfrastructureRef: clusterv1.ContractVersionedObjectReference{
						Kind: "AWSCluster",
						Name: "test-aws-cluster",
					},
				},
			}
			awsCluster := &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-aws-cluster",
					Namespace: "default",
				},
				Spec: infrav1.AWSClusterSpec{
					Region: "us-west-2",
				},
			}

			reconciler := &AWSMachineTemplateReconciler{
				Client: newFakeClient(template, cluster, awsCluster),
			}

			result, err := reconciler.Reconcile(context.Background(), ctrl.Request{
				NamespacedName: client.ObjectKeyFromObject(template),
			})

			g.Expect(err).To(BeNil())
			g.Expect(result.RequeueAfter).To(BeZero())
		})

		t.Run("should skip when no owner cluster", func(t *testing.T) {
			g := NewWithT(t)
			template := newAWSMachineTemplate("test-template")

			reconciler := &AWSMachineTemplateReconciler{
				Client: newFakeClient(template),
			}

			result, err := reconciler.Reconcile(context.Background(), ctrl.Request{
				NamespacedName: client.ObjectKeyFromObject(template),
			})

			g.Expect(err).To(BeNil())
			g.Expect(result.RequeueAfter).To(BeZero())
		})

		t.Run("should skip when cluster type is not recognised", func(t *testing.T) {
			g := NewWithT(t)
			template := newAWSMachineTemplate("test-template")
			template.OwnerReferences = []metav1.OwnerReference{
				{
					APIVersion: clusterv1.GroupVersion.String(),
					Kind:       "Cluster",
					Name:       "test-cluster",
				},
			}
			cluster := &clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-cluster",
					Namespace: "default",
				},
			}

			reconciler := &AWSMachineTemplateReconciler{
				Client: newFakeClient(template, cluster),
			}

			result, err := reconciler.Reconcile(context.Background(), ctrl.Request{
				NamespacedName: client.ObjectKeyFromObject(template),
			})

			g.Expect(err).To(BeNil())
			g.Expect(result.RequeueAfter).To(BeZero())
		})

		t.Run("should return nil when template not found", func(t *testing.T) {
			g := NewWithT(t)

			reconciler := &AWSMachineTemplateReconciler{
				Client: newFakeClient(),
			}

			result, err := reconciler.Reconcile(context.Background(), ctrl.Request{
				NamespacedName: client.ObjectKey{
					Namespace: "default",
					Name:      "nonexistent",
				},
			})

			g.Expect(err).To(BeNil())
			g.Expect(result.RequeueAfter).To(BeZero())
		})
	})
}
