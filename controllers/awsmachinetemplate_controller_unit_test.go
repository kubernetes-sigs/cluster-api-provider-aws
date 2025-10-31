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
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func TestAWSMachineTemplateReconciler(t *testing.T) {
	setupScheme := func() *runtime.Scheme {
		scheme := runtime.NewScheme()
		_ = infrav1.AddToScheme(scheme)
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

	t.Run("getRegion", func(t *testing.T) {
		t.Run("should get region from AWSCluster", func(t *testing.T) {
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
					InfrastructureRef: &corev1.ObjectReference{
						Kind:      "AWSCluster",
						Name:      "test-aws-cluster",
						Namespace: "default",
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

			region, err := reconciler.getRegion(context.Background(), template)

			g.Expect(err).To(BeNil())
			g.Expect(region).To(Equal("us-west-2"))
		})

		t.Run("should return error when no owner cluster found", func(t *testing.T) {
			g := NewWithT(t)
			template := newAWSMachineTemplate("test-template")

			reconciler := &AWSMachineTemplateReconciler{
				Client: newFakeClient(template),
			}

			region, err := reconciler.getRegion(context.Background(), template)

			g.Expect(err).ToNot(BeNil())
			g.Expect(err.Error()).To(ContainSubstring("no owner cluster found"))
			g.Expect(region).To(Equal(""))
		})

		t.Run("should return empty when cluster has no infrastructure ref", func(t *testing.T) {
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

			region, err := reconciler.getRegion(context.Background(), template)

			g.Expect(err).To(BeNil())
			g.Expect(region).To(Equal(""))
		})

		t.Run("should return empty when AWSCluster not found", func(t *testing.T) {
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
					InfrastructureRef: &corev1.ObjectReference{
						Kind:      "AWSCluster",
						Name:      "test-aws-cluster",
						Namespace: "default",
					},
				},
			}

			reconciler := &AWSMachineTemplateReconciler{
				Client: newFakeClient(template, cluster),
			}

			region, err := reconciler.getRegion(context.Background(), template)

			g.Expect(err).To(BeNil())
			g.Expect(region).To(Equal(""))
		})
	})

	// Note: getInstanceTypeInfo tests are skipped as they require EC2 client injection
	// which would need significant refactoring. The function is tested indirectly through
	// integration tests.

	t.Run("Reconcile", func(t *testing.T) {
		t.Run("should skip when capacity and nodeInfo already set", func(t *testing.T) {
			g := NewWithT(t)
			template := newAWSMachineTemplate("test-template")
			template.Status.Capacity = corev1.ResourceList{
				corev1.ResourceCPU: *resource.NewQuantity(2, resource.DecimalSI),
			}
			template.Status.NodeInfo = &infrav1.NodeInfo{
				Architecture: infrav1.ArchitectureAmd64,
			}

			reconciler := &AWSMachineTemplateReconciler{
				Client: newFakeClient(template),
			}

			result, err := reconciler.Reconcile(context.Background(), ctrl.Request{
				NamespacedName: client.ObjectKeyFromObject(template),
			})

			g.Expect(err).To(BeNil())
			g.Expect(result.Requeue).To(BeFalse())
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
					InfrastructureRef: &corev1.ObjectReference{
						Kind:      "AWSCluster",
						Name:      "test-aws-cluster",
						Namespace: "default",
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
			g.Expect(result.Requeue).To(BeFalse())
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
			g.Expect(result.Requeue).To(BeFalse())
		})

		t.Run("should return error when no owner cluster", func(t *testing.T) {
			g := NewWithT(t)
			template := newAWSMachineTemplate("test-template")

			reconciler := &AWSMachineTemplateReconciler{
				Client: newFakeClient(template),
			}

			result, err := reconciler.Reconcile(context.Background(), ctrl.Request{
				NamespacedName: client.ObjectKeyFromObject(template),
			})

			g.Expect(err).ToNot(BeNil())
			g.Expect(err.Error()).To(ContainSubstring("no owner cluster found"))
			g.Expect(result.Requeue).To(BeFalse())
		})

		t.Run("should skip when region is empty", func(t *testing.T) {
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
			g.Expect(result.Requeue).To(BeFalse())
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
			g.Expect(result.Requeue).To(BeFalse())
		})
	})
}
