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

package gc

import (
	"context"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/annotations"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/controllers/external"
)

const (
	testClusterName = "test-cluster"
)

func TestEnableGC(t *testing.T) {
	RegisterTestingT(t)

	testCases := []struct {
		name         string
		clusterName  string
		existingObjs []client.Object
		expectError  bool
	}{
		{
			name:         "no capi cluster",
			clusterName:  testClusterName,
			existingObjs: []client.Object{},
			expectError:  true,
		},
		{
			name:         "no infra cluster",
			clusterName:  testClusterName,
			existingObjs: newManagedCluster(testClusterName, true),
			expectError:  true,
		},
		{
			name:         "with managed control plane and no annotation",
			clusterName:  testClusterName,
			existingObjs: newManagedClusterWithAnnotations(testClusterName, nil),
			expectError:  false,
		},
		{
			name:         "with awscluster and no annotation",
			clusterName:  testClusterName,
			existingObjs: newUnManagedClusterWithAnnotations(testClusterName, nil),
			expectError:  false,
		},
		{
			name:         "with managed control plane and existing annotation",
			clusterName:  testClusterName,
			existingObjs: newManagedClusterWithAnnotations(testClusterName, map[string]string{infrav1.ExternalResourceGCAnnotation: "false"}),
			expectError:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			input := GCInput{
				ClusterName: tc.clusterName,
				Namespace:   "default",
			}

			fake := newFakeClient(scheme, tc.existingObjs...)
			ctx := context.TODO()

			proc, err := New(input, WithClient(fake))
			g.Expect(err).NotTo(HaveOccurred())

			resErr := proc.Enable(ctx)
			if tc.expectError {
				g.Expect(resErr).To(HaveOccurred())
				return
			}
			g.Expect(resErr).NotTo(HaveOccurred())

			cluster := tc.existingObjs[0].(*clusterv1.Cluster)
			ref := cluster.Spec.InfrastructureRef

			obj, err := external.Get(ctx, fake, ref)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(obj).NotTo(BeNil())

			annotationVal, hasAnnotation := annotations.Get(obj, infrav1.ExternalResourceGCAnnotation)
			g.Expect(hasAnnotation).To(BeTrue())
			g.Expect(annotationVal).To(Equal("true"))
		})
	}
}

func TestDisableGC(t *testing.T) {
	RegisterTestingT(t)

	testCases := []struct {
		name         string
		clusterName  string
		existingObjs []client.Object
		expectError  bool
	}{
		{
			name:         "no capi cluster",
			clusterName:  testClusterName,
			existingObjs: []client.Object{},
			expectError:  true,
		},
		{
			name:         "no infra cluster",
			clusterName:  testClusterName,
			existingObjs: newManagedCluster(testClusterName, true),
			expectError:  true,
		},
		{
			name:         "with managed control plane and with annotation",
			clusterName:  testClusterName,
			existingObjs: newManagedClusterWithAnnotations(testClusterName, map[string]string{infrav1.ExternalResourceGCAnnotation: "true"}),
			expectError:  false,
		},
		{
			name:         "with awscluster and with annotation",
			clusterName:  testClusterName,
			existingObjs: newUnManagedClusterWithAnnotations(testClusterName, map[string]string{infrav1.ExternalResourceGCAnnotation: "true"}),
			expectError:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			input := GCInput{
				ClusterName: tc.clusterName,
				Namespace:   "default",
			}

			fake := newFakeClient(scheme, tc.existingObjs...)
			ctx := context.TODO()

			proc, err := New(input, WithClient(fake))
			g.Expect(err).NotTo(HaveOccurred())

			resErr := proc.Disable(ctx)
			if tc.expectError {
				g.Expect(resErr).To(HaveOccurred())
				return
			}
			g.Expect(resErr).NotTo(HaveOccurred())

			cluster := tc.existingObjs[0].(*clusterv1.Cluster)
			ref := cluster.Spec.InfrastructureRef

			obj, err := external.Get(ctx, fake, ref)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(obj).NotTo(BeNil())

			annotationVal, hasAnnotation := annotations.Get(obj, infrav1.ExternalResourceGCAnnotation)
			g.Expect(hasAnnotation).To(BeTrue())
			g.Expect(annotationVal).To(Equal("false"))
		})
	}
}

func TestConfigureGC(t *testing.T) {
	RegisterTestingT(t)

	testCases := []struct {
		name         string
		clusterName  string
		gcTasks      []string
		existingObjs []client.Object
		expectError  bool
	}{
		{
			name:         "no capi cluster",
			clusterName:  testClusterName,
			existingObjs: []client.Object{},
			expectError:  true,
		},
		{
			name:         "no infra cluster",
			clusterName:  testClusterName,
			existingObjs: newManagedCluster(testClusterName, true),
			expectError:  true,
		},
		{
			name:         "with managed control plane and no annotation",
			clusterName:  testClusterName,
			existingObjs: newManagedCluster(testClusterName, false),
			gcTasks:      []string{"load-balancer", "target-group"},
			expectError:  false,
		},
		{
			name:         "with awscluster and no annotation",
			clusterName:  testClusterName,
			existingObjs: newUnManagedCluster(testClusterName, false),
			gcTasks:      []string{"load-balancer", "security-group"},
			expectError:  false,
		},
		{
			name:         "with managed control plane and with annotation",
			clusterName:  testClusterName,
			existingObjs: newManagedClusterWithAnnotations(testClusterName, map[string]string{infrav1.ExternalResourceGCTasksAnnotation: "security-group"}),
			gcTasks:      []string{"load-balancer", "target-group"},
			expectError:  false,
		},
		{
			name:         "with awscluster and with annotation",
			clusterName:  testClusterName,
			existingObjs: newUnManagedClusterWithAnnotations(testClusterName, map[string]string{infrav1.ExternalResourceGCTasksAnnotation: "security-group"}),
			gcTasks:      []string{"load-balancer", "target-group"},
			expectError:  false,
		},
		{
			name:         "with awscluster and invalid gc tasks",
			clusterName:  testClusterName,
			existingObjs: newUnManagedCluster(testClusterName, false),
			gcTasks:      []string{"load-balancer", "INVALID"},
			expectError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			input := GCInput{
				ClusterName: tc.clusterName,
				Namespace:   "default",
			}

			fake := newFakeClient(scheme, tc.existingObjs...)
			ctx := context.TODO()

			proc, err := New(input, WithClient(fake))
			g.Expect(err).NotTo(HaveOccurred())

			resErr := proc.Configure(ctx, tc.gcTasks)
			if tc.expectError {
				g.Expect(resErr).To(HaveOccurred())
				return
			}
			g.Expect(resErr).NotTo(HaveOccurred())

			cluster := tc.existingObjs[0].(*clusterv1.Cluster)
			ref := cluster.Spec.InfrastructureRef

			obj, err := external.Get(ctx, fake, ref)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(obj).NotTo(BeNil())

			expected := strings.Join(tc.gcTasks, ",")
			annotationVal, hasAnnotation := annotations.Get(obj, infrav1.ExternalResourceGCTasksAnnotation)

			if expected != "" {
				g.Expect(hasAnnotation).To(BeTrue())
				g.Expect(annotationVal).To(Equal(expected))
			} else {
				g.Expect(hasAnnotation).To(BeFalse())
			}
		})
	}
}

func newFakeClient(scheme *runtime.Scheme, objs ...client.Object) client.Client {
	return fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build()
}

func newManagedCluster(name string, excludeInfra bool) []client.Object {
	objs := []client.Object{
		&clusterv1.Cluster{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Cluster",
				APIVersion: clusterv1.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: "default",
			},
			Spec: clusterv1.ClusterSpec{
				InfrastructureRef: &corev1.ObjectReference{
					Name:       name,
					Namespace:  "default",
					Kind:       "AWSManagedControlPlane",
					APIVersion: ekscontrolplanev1.GroupVersion.String(),
				},
			},
		},
	}

	if !excludeInfra {
		objs = append(objs, &ekscontrolplanev1.AWSManagedControlPlane{
			TypeMeta: metav1.TypeMeta{
				Kind:       "AWSManagedControlPlane",
				APIVersion: ekscontrolplanev1.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: "default",
			},
		})
	}

	return objs
}

func newManagedClusterWithAnnotations(name string, annotations map[string]string) []client.Object {
	objs := newManagedCluster(name, false)

	mcp := objs[1].(*ekscontrolplanev1.AWSManagedControlPlane)
	mcp.ObjectMeta.Annotations = annotations

	return objs
}

func newUnManagedCluster(name string, excludeInfra bool) []client.Object {
	objs := []client.Object{
		&clusterv1.Cluster{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Cluster",
				APIVersion: clusterv1.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: "default",
			},
			Spec: clusterv1.ClusterSpec{
				InfrastructureRef: &corev1.ObjectReference{
					Name:       name,
					Namespace:  "default",
					Kind:       "AWSCluster",
					APIVersion: infrav1.GroupVersion.String(),
				},
			},
		},
	}

	if !excludeInfra {
		objs = append(objs, &infrav1.AWSCluster{
			TypeMeta: metav1.TypeMeta{
				Kind:       "AWSCluster",
				APIVersion: infrav1.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: "default",
			},
		})
	}

	return objs
}

func newUnManagedClusterWithAnnotations(name string, annotations map[string]string) []client.Object {
	objs := newUnManagedCluster(name, false)

	awsc := objs[1].(*infrav1.AWSCluster)
	awsc.ObjectMeta.Annotations = annotations

	return objs
}
