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
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2/klogr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
)

func TestAWSMachineReconciler_AWSClusterToAWSMachines(t *testing.T) {
	testCases := []struct {
		name         string
		ownerCluster *clusterv1.Cluster
		awsCluster   *infrav1.AWSCluster
		awsMachine   *clusterv1.Machine
		requests     []reconcile.Request
	}{
		{
			name:         "Should create reconcile request successfully",
			ownerCluster: &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "capi-test-6"}},
			awsMachine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name: "aws-test-6",
					Labels: map[string]string{
						clusterv1.ClusterLabelName: "capi-test-6",
					},
				},
				Spec: clusterv1.MachineSpec{
					InfrastructureRef: corev1.ObjectReference{
						Kind:       "AWSMachine",
						Name:       "aws-machine-6",
						APIVersion: infrav1.GroupVersion.String(),
					}},
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "aws-test-6",
					OwnerReferences: []metav1.OwnerReference{
						{
							Name:       "capi-test-6",
							Kind:       "Cluster",
							APIVersion: clusterv1.GroupVersion.String(),
						},
					},
				},
			},
			requests: []reconcile.Request{
				{
					NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "aws-machine-6",
					},
				},
			},
		},
		{
			name:         "Should not create reconcile request for deleted clusters",
			ownerCluster: &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "capi-test-1", DeletionTimestamp: &metav1.Time{Time: time.Now()}}},
			awsMachine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						clusterv1.ClusterLabelName: "aws-test-1",
					},
					Name: "aws-test-1",
				},
				Spec: clusterv1.MachineSpec{
					InfrastructureRef: corev1.ObjectReference{
						Kind:       "AWSMachine",
						Name:       "aws-machine-1",
						APIVersion: infrav1.GroupVersion.String(),
					}},
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "aws-test-1",
					OwnerReferences: []metav1.OwnerReference{
						{
							Name:       "capi-test-1",
							Kind:       "Cluster",
							APIVersion: clusterv1.GroupVersion.String(),
						},
					},
					DeletionTimestamp: &metav1.Time{Time: time.Now()},
				},
			},
		},
		{
			name: "Should not create reconcile request if ownerCluster not found",
			awsMachine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						clusterv1.ClusterLabelName: "aws-test-2",
					},
					Name: "aws-test-2",
				},
				Spec: clusterv1.MachineSpec{
					InfrastructureRef: corev1.ObjectReference{
						Kind:       "AWSMachine",
						Name:       "aws-machine-2",
						APIVersion: infrav1.GroupVersion.String(),
					}},
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "aws-test-2",
					OwnerReferences: []metav1.OwnerReference{
						{
							Name:       "capi-test-2",
							Kind:       "Cluster",
							APIVersion: clusterv1.GroupVersion.String(),
						},
					},
				},
			},
		},
		{
			name:         "Should not create reconcile request if owned Machines not found",
			ownerCluster: &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "capi-test-3"}},
			awsMachine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name: "aws-test-3",
				},
				Spec: clusterv1.MachineSpec{
					InfrastructureRef: corev1.ObjectReference{
						Kind:       "AWSMachine",
						Name:       "aws-machine-3",
						APIVersion: infrav1.GroupVersion.String(),
					}},
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "aws-test-3",
					OwnerReferences: []metav1.OwnerReference{
						{
							Name:       "capi-test-3",
							Kind:       "Cluster",
							APIVersion: clusterv1.GroupVersion.String(),
						},
					},
				},
			},
			requests: []reconcile.Request{},
		},
		{
			name:         "Should not create reconcile request if owned Machine type is not AWSMachine",
			ownerCluster: &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "capi-test-4"}},
			awsMachine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						clusterv1.ClusterLabelName: "capi-test-4",
					},
					Name:      "aws-test-4",
					Namespace: "default",
				},
				TypeMeta: metav1.TypeMeta{
					Kind: "Machine",
				},
				Spec: clusterv1.MachineSpec{
					InfrastructureRef: corev1.ObjectReference{
						Kind:       "Machine",
						Name:       "aws-machine-4",
						APIVersion: infrav1.GroupVersion.String(),
					}},
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "aws-test-4",
					OwnerReferences: []metav1.OwnerReference{
						{
							Name:       "capi-test-4",
							Kind:       "Cluster",
							APIVersion: clusterv1.GroupVersion.String(),
						},
					},
				},
			},
			requests: []reconcile.Request{},
		},
		{
			name:         "Should not create reconcile request if name for machine in infrastructure ref not found",
			ownerCluster: &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "capi-test-5"}},
			awsMachine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name: "aws-test-5",
					Labels: map[string]string{
						clusterv1.ClusterLabelName: "capi-test-5",
					},
				},
				Spec: clusterv1.MachineSpec{
					InfrastructureRef: corev1.ObjectReference{
						Kind:       "AWSMachine",
						APIVersion: infrav1.GroupVersion.String(),
					}},
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "aws-test-5",
					OwnerReferences: []metav1.OwnerReference{
						{
							Name:       "capi-test-5",
							Kind:       "Cluster",
							APIVersion: clusterv1.GroupVersion.String(),
						},
					},
				},
			},
			requests: []reconcile.Request{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			reconciler := &AWSMachineReconciler{
				Client: testEnv.Client,
				Log:    klogr.New(),
			}
			ns, err := testEnv.CreateNamespace(ctx, fmt.Sprintf("namespace-%s", util.RandomString(5)))
			g.Expect(err).To(BeNil())

			createObject(g, tc.ownerCluster, ns.Name)
			defer cleanupObject(g, tc.ownerCluster)

			createObject(g, tc.awsMachine, ns.Name)
			defer cleanupObject(g, tc.awsMachine)

			tc.awsCluster.Namespace = ns.Name
			defer t.Cleanup(func() {
				g.Expect(testEnv.Cleanup(ctx, tc.awsCluster, ns)).To(Succeed())
			})

			requests := reconciler.AWSClusterToAWSMachines(klogr.New())(tc.awsCluster)
			if tc.requests != nil {
				if len(tc.requests) > 0 {
					tc.requests[0].Namespace = ns.Name
				}
				g.Expect(requests).To(ConsistOf(tc.requests))
			} else {
				g.Expect(requests).To(BeNil())
			}
		})
	}
}

func TestAWSMachineReconciler_requeueAWSMachinesForUnpausedCluster(t *testing.T) {
	testCases := []struct {
		name         string
		ownerCluster *clusterv1.Cluster
		requests     []reconcile.Request
	}{
		{
			name:         "Should not create reconcile request for deleted clusters",
			ownerCluster: &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "capi-test-1", Namespace: "default", DeletionTimestamp: &metav1.Time{Time: time.Now()}}},
		},
		{
			name:         "Should create reconcile request successfully",
			ownerCluster: &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "capi-test-1", Namespace: "default"}},
			requests:     []reconcile.Request{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			reconciler := &AWSMachineReconciler{
				Client: testEnv.Client,
				Log:    klogr.New(),
			}
			requests := reconciler.requeueAWSMachinesForUnpausedCluster(klogr.New())(tc.ownerCluster)
			if tc.requests != nil {
				g.Expect(requests).To(ConsistOf(tc.requests))
			} else {
				g.Expect(requests).To(BeNil())
			}
		})
	}
}

func TestAWSMachineReconciler_indexAWSMachineByInstanceID(t *testing.T) {
	t.Run("Should not return instance id if cluster type is not AWSCluster", func(t *testing.T) {
		g := NewWithT(t)
		reconciler := &AWSMachineReconciler{
			Client: testEnv.Client,
			Log:    klogr.New(),
		}
		machine := &clusterv1.Machine{ObjectMeta: metav1.ObjectMeta{Name: "capi-test-1", Namespace: "default"}}
		requests := reconciler.indexAWSMachineByInstanceID(machine)
		g.Expect(requests).To(BeNil())
	})
	t.Run("Should return instance id successfully", func(t *testing.T) {
		g := NewWithT(t)
		reconciler := &AWSMachineReconciler{
			Client: testEnv.Client,
			Log:    klogr.New(),
		}
		awsMachine := &infrav1.AWSMachine{ObjectMeta: metav1.ObjectMeta{Name: "capi-test-1", Namespace: "default"}, Spec: infrav1.AWSMachineSpec{InstanceID: aws.String("12345")}}
		requests := reconciler.indexAWSMachineByInstanceID(awsMachine)
		g.Expect(requests).To(ConsistOf([]string{"12345"}))
	})
	t.Run("Should not return instance id if instance id is not present", func(t *testing.T) {
		g := NewWithT(t)
		reconciler := &AWSMachineReconciler{
			Client: testEnv.Client,
			Log:    klogr.New(),
		}
		awsMachine := &infrav1.AWSMachine{ObjectMeta: metav1.ObjectMeta{Name: "capi-test-1", Namespace: "default"}}
		requests := reconciler.indexAWSMachineByInstanceID(awsMachine)
		g.Expect(requests).To(BeNil())
	})
}

func TestAWSMachineReconciler_Reconcile(t *testing.T) {
	testCases := []struct {
		name         string
		awsMachine   *infrav1.AWSMachine
		ownerMachine *clusterv1.Machine
		ownerCluster *clusterv1.Cluster
		awsCluster   *infrav1.AWSCluster
		expectError  bool
		requeue      bool
	}{
		{
			name:        "Should Reconcile successfully if no AWSMachine found",
			expectError: false,
		},
		{
			name:        "Should Reconcile AWSMachine with requeue",
			awsMachine:  &infrav1.AWSMachine{ObjectMeta: metav1.ObjectMeta{Name: "aws-test-1"}, Spec: infrav1.AWSMachineSpec{InstanceType: "test"}},
			requeue:     true,
			expectError: false,
		},
		{
			name: "Should fail Reconcile with GetOwnerMachine failure",
			awsMachine: &infrav1.AWSMachine{ObjectMeta: metav1.ObjectMeta{Name: "aws-test-2",
				OwnerReferences: []metav1.OwnerReference{
					{
						APIVersion: clusterv1.GroupVersion.String(),
						Kind:       "Machine",
						Name:       "capi-test-machine",
						UID:        "1",
					}}},
				Spec: infrav1.AWSMachineSpec{InstanceType: "test"},
			},
			expectError: true,
		},
		{
			name: "Should not Reconcile if machine does not contain cluster label",
			awsMachine: &infrav1.AWSMachine{ObjectMeta: metav1.ObjectMeta{
				Name: "aws-test-3", Annotations: map[string]string{clusterv1.PausedAnnotation: ""}, OwnerReferences: []metav1.OwnerReference{
					{
						APIVersion: clusterv1.GroupVersion.String(),
						Kind:       "Machine",
						Name:       "capi-test-machine",
						UID:        "1",
					},
				}}, Spec: infrav1.AWSMachineSpec{InstanceType: "test"},
			},
			ownerMachine: &clusterv1.Machine{ObjectMeta: metav1.ObjectMeta{Name: "capi-test-machine"}},
			ownerCluster: &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "capi-test-1"}},
			expectError:  false,
		},
		{
			name: "Should not Reconcile if cluster is paused",
			awsMachine: &infrav1.AWSMachine{ObjectMeta: metav1.ObjectMeta{
				Name: "aws-test-4", Annotations: map[string]string{clusterv1.PausedAnnotation: ""}, OwnerReferences: []metav1.OwnerReference{
					{
						APIVersion: clusterv1.GroupVersion.String(),
						Kind:       "Machine",
						Name:       "capi-test-machine",
						UID:        "1",
					},
				}}, Spec: infrav1.AWSMachineSpec{InstanceType: "test"},
			},
			ownerMachine: &clusterv1.Machine{ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					clusterv1.ClusterLabelName: "capi-test-1",
				},
				Name: "capi-test-machine", Namespace: "default"}},
			ownerCluster: &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "capi-test-1"}},
			expectError:  false,
		},
		{
			name: "Should not Reconcile if AWSManagedControlPlane is not ready",
			awsMachine: &infrav1.AWSMachine{ObjectMeta: metav1.ObjectMeta{
				Name: "aws-test-5", OwnerReferences: []metav1.OwnerReference{
					{
						APIVersion: clusterv1.GroupVersion.String(),
						Kind:       "Machine",
						Name:       "capi-test-machine",
						UID:        "1",
					},
				}}, Spec: infrav1.AWSMachineSpec{InstanceType: "test"},
			},
			ownerMachine: &clusterv1.Machine{ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					clusterv1.ClusterLabelName: "capi-test-1",
				},
				Name: "capi-test-machine", Namespace: "default"}},
			ownerCluster: &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "capi-test-1"},
				Spec: clusterv1.ClusterSpec{
					ControlPlaneRef: &corev1.ObjectReference{Kind: AWSManagedControlPlaneRefKind},
				}},
			expectError: false,
		},
		{
			name: "Should not Reconcile if AWSCluster is not ready",
			awsMachine: &infrav1.AWSMachine{ObjectMeta: metav1.ObjectMeta{
				Name: "aws-test-5", OwnerReferences: []metav1.OwnerReference{
					{
						APIVersion: clusterv1.GroupVersion.String(),
						Kind:       "Machine",
						Name:       "capi-test-machine",
						UID:        "1",
					},
				}}, Spec: infrav1.AWSMachineSpec{InstanceType: "test"},
			},
			ownerMachine: &clusterv1.Machine{ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					clusterv1.ClusterLabelName: "capi-test-1",
				},
				Name: "capi-test-machine", Namespace: "default"}},
			ownerCluster: &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "capi-test-1"},
				Spec: clusterv1.ClusterSpec{
					InfrastructureRef: &corev1.ObjectReference{Name: "aws-test-5"},
				}},
			expectError: false,
		},
		{
			name: "Should fail to reconcile while fetching infra cluster",
			awsMachine: &infrav1.AWSMachine{ObjectMeta: metav1.ObjectMeta{
				Name: "aws-test-5", OwnerReferences: []metav1.OwnerReference{
					{
						APIVersion: clusterv1.GroupVersion.String(),
						Kind:       "Machine",
						Name:       "capi-test-machine",
						UID:        "1",
					},
				}}, Spec: infrav1.AWSMachineSpec{InstanceType: "test"},
			},
			ownerMachine: &clusterv1.Machine{ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					clusterv1.ClusterLabelName: "capi-test-1",
				},
				Name: "capi-test-machine", Namespace: "default"}},
			ownerCluster: &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "capi-test-1"},
				Spec: clusterv1.ClusterSpec{
					InfrastructureRef: &corev1.ObjectReference{Name: "aws-test-5"},
				}},
			awsCluster:  &infrav1.AWSCluster{ObjectMeta: metav1.ObjectMeta{Name: "aws-test-5"}},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			reconciler := &AWSMachineReconciler{
				Client: testEnv.Client,
			}
			ns, err := testEnv.CreateNamespace(ctx, fmt.Sprintf("namespace-%s", util.RandomString(5)))
			g.Expect(err).To(BeNil())
			defer func() {
				g.Expect(testEnv.Cleanup(ctx, ns)).To(Succeed())
			}()

			createObject(g, tc.ownerCluster, ns.Name)
			defer cleanupObject(g, tc.ownerCluster)

			createObject(g, tc.awsCluster, ns.Name)
			defer cleanupObject(g, tc.awsCluster)

			createObject(g, tc.ownerMachine, ns.Name)
			defer cleanupObject(g, tc.ownerMachine)

			createObject(g, tc.awsMachine, ns.Name)
			defer cleanupObject(g, tc.awsMachine)
			if tc.awsMachine != nil {
				g.Eventually(func() bool {
					machine := &infrav1.AWSMachine{}
					key := client.ObjectKey{
						Name:      tc.awsMachine.Name,
						Namespace: ns.Name,
					}
					err = testEnv.Get(ctx, key, machine)
					return err == nil
				}, 10*time.Second).Should(Equal(true))

				result, err := reconciler.Reconcile(ctx, ctrl.Request{
					NamespacedName: client.ObjectKey{
						Namespace: tc.awsMachine.Namespace,
						Name:      tc.awsMachine.Name,
					},
				})
				if tc.expectError {
					g.Expect(err).ToNot(BeNil())
				} else {
					g.Expect(err).To(BeNil())
				}
				if tc.requeue {
					g.Expect(result.RequeueAfter).To(BeZero())
				}
			} else {
				_, err = reconciler.Reconcile(ctx, ctrl.Request{
					NamespacedName: client.ObjectKey{
						Namespace: "default",
						Name:      "test",
					},
				})
				g.Expect(err).To(BeNil())
			}
		})
	}
}

func createObject(g *WithT, obj client.Object, namespace string) {
	if obj.DeepCopyObject() != nil {
		obj.SetNamespace(namespace)
		g.Expect(testEnv.Create(ctx, obj)).To(Succeed())
	}
}

func cleanupObject(g *WithT, obj client.Object) {
	if obj.DeepCopyObject() != nil {
		g.Expect(testEnv.Cleanup(ctx, obj)).To(Succeed())
	}
}
