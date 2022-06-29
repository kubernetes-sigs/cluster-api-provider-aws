/*
Copyright 2022 The Kubernetes Authors.

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

package workload

import (
	"context"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func TestReconcileDelete(t *testing.T) {
	testCases := []struct {
		name               string
		cluster            *clusterv1.Cluster
		infraCluster       client.Object
		objs               []client.Object
		expectErr          bool
		expectRequeue      bool
		expectRequeueAfter bool
		expectDeleted      []string
	}{
		{
			name:         "eks with services of type load balancer",
			cluster:      createEKSCluster(),
			infraCluster: createManagedControlPlane(),
			objs: []client.Object{
				createService("svc1", corev1.ServiceTypeLoadBalancer, false),
			},
			expectErr:          false,
			expectRequeue:      false,
			expectRequeueAfter: true,
			expectDeleted:      []string{"svc1"},
		},
		{
			name:         "awscluster with services of type load balancer",
			cluster:      createUnmanagedCluster(),
			infraCluster: createAWSCluser(),
			objs: []client.Object{
				createService("svc1", corev1.ServiceTypeLoadBalancer, false),
			},
			expectErr:          false,
			expectRequeue:      false,
			expectRequeueAfter: true,
			expectDeleted:      []string{"svc1"},
		},
		{
			name:         "eks with no services of type load balancer",
			cluster:      createEKSCluster(),
			infraCluster: createManagedControlPlane(),
			objs: []client.Object{
				createService("svc1", corev1.ServiceTypeClusterIP, false),
			},
			expectErr:          false,
			expectRequeue:      false,
			expectRequeueAfter: false,
			expectDeleted:      []string{},
		},
		{
			name:               "eks with no services",
			cluster:            createEKSCluster(),
			infraCluster:       createManagedControlPlane(),
			objs:               []client.Object{},
			expectErr:          false,
			expectRequeue:      false,
			expectRequeueAfter: false,
			expectDeleted:      []string{},
		},
		{
			name:         "eks with services of type load balancer",
			cluster:      createEKSCluster(),
			infraCluster: createManagedControlPlane(),
			objs: []client.Object{
				createService("svc1", corev1.ServiceTypeLoadBalancer, true),
			},
			expectErr:          false,
			expectRequeue:      false,
			expectRequeueAfter: true,
			expectDeleted:      []string{},
		},
	}

	scheme := runtime.NewScheme()
	_ = corev1.AddToScheme(scheme)
	_ = ekscontrolplanev1.AddToScheme(scheme)
	_ = infrav1.AddToScheme(scheme)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			ctx := context.TODO()
			client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(tc.objs...).Build()

			s, err := scope.NewExternalResourceGCScope(scope.ExternalResourceGCScopeParams{
				Cluster:      tc.cluster,
				InfraCluster: tc.infraCluster,
				Client:       client,
			}, scope.WithExternalResourceGCScopeTenantClient(client))
			g.Expect(err).To(BeNil())
			g.Expect(s).NotTo(BeNil())

			wkSvc := NewService(s)
			res, err := wkSvc.ReconcileDelete(ctx)

			if tc.expectErr {
				g.Expect(err).NotTo(BeNil())
				return
			}

			g.Expect(err).To(BeNil())
			g.Expect(res).NotTo(BeNil())
			if tc.expectRequeueAfter {
				g.Expect(res.RequeueAfter).NotTo(BeNil())
				g.Expect(res.RequeueAfter).NotTo(BeZero())
			}
			if tc.expectRequeue {
				g.Expect(res.Requeue).To(BeTrue())
			}
			for _, svcDeleted := range tc.expectDeleted {
				svc := &corev1.Service{}
				getErr := client.Get(ctx, types.NamespacedName{
					Name:      svcDeleted,
					Namespace: "default",
				}, svc)
				g.Expect(apierrors.IsNotFound(getErr)).To(BeTrue(), "expecting the service to be deleted")
			}
		})
	}
}

func createEKSCluster() *clusterv1.Cluster {
	return &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cluster1",
			Namespace: "default",
		},
		Spec: clusterv1.ClusterSpec{
			InfrastructureRef: &corev1.ObjectReference{
				Kind:       "AWSManagedControlPlane",
				APIVersion: ekscontrolplanev1.GroupVersion.String(),
				Name:       "cp1",
				Namespace:  "default",
			},
		},
	}
}

func createManagedControlPlane() *ekscontrolplanev1.AWSManagedControlPlane {
	return &ekscontrolplanev1.AWSManagedControlPlane{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AWSManagedControlPlane",
			APIVersion: ekscontrolplanev1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cp1",
			Namespace: "default",
		},
		Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{},
	}
}

func createAWSCluser() *infrav1.AWSCluster {
	return &infrav1.AWSCluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AWSCluster",
			APIVersion: infrav1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cluster1",
			Namespace: "default",
		},
		Spec: infrav1.AWSClusterSpec{},
	}
}

func createUnmanagedCluster() *clusterv1.Cluster {
	return &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cluster1",
			Namespace: "default",
		},
		Spec: clusterv1.ClusterSpec{
			InfrastructureRef: &corev1.ObjectReference{
				Kind:       "AWSCluster",
				APIVersion: infrav1.GroupVersion.String(),
				Name:       "cluster1",
				Namespace:  "default",
			},
		},
	}
}

func createService(name string, svcType corev1.ServiceType, deleting bool) *corev1.Service {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Spec: corev1.ServiceSpec{
			Type: svcType,
		},
	}

	if deleting {
		svc.ObjectMeta.DeletionTimestamp = &metav1.Time{
			Time: time.Now(),
		}
	}

	return svc
}
