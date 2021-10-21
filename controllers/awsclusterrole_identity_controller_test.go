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

	"github.com/aws/aws-sdk-go/aws"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestAWSClusterRoleIdentityReconciler(t *testing.T) {
	g := NewWithT(t)
	ctx := context.Background()

	awsCluster := &infrav1.AWSCluster{
		TypeMeta:   metav1.TypeMeta{APIVersion: "infrastructure.cluster.x-k8s.io/v1beta1", Kind: awsClusterKind},
		ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "default"},
		Spec:       infrav1.AWSClusterSpec{IdentityRef: &infrav1.AWSIdentityReference{Name: "id-2", Kind: infrav1.ClusterRoleIdentityKind}}}

	ownerRef := []metav1.OwnerReference{
		{
			APIVersion:         "infrastructure.cluster.x-k8s.io/v1beta1",
			Kind:               awsClusterKind,
			Name:               awsCluster.Name,
			UID:                awsCluster.UID,
			BlockOwnerDeletion: aws.Bool(true),
		},
	}

	roleIdentity1 := infrav1.AWSClusterRoleIdentity{
		TypeMeta:   metav1.TypeMeta{Kind: string(infrav1.ClusterRoleIdentityKind)},
		ObjectMeta: metav1.ObjectMeta{Name: "id-1", OwnerReferences: ownerRef}}
	roleIdentity2 := infrav1.AWSClusterRoleIdentity{
		TypeMeta:   metav1.TypeMeta{Kind: string(infrav1.ClusterRoleIdentityKind)},
		ObjectMeta: metav1.ObjectMeta{Name: "id-2", OwnerReferences: ownerRef}}
	roleIdentity3 := infrav1.AWSClusterRoleIdentity{
		TypeMeta:   metav1.TypeMeta{Kind: string(infrav1.ClusterRoleIdentityKind)},
		ObjectMeta: metav1.ObjectMeta{Name: "id-3", OwnerReferences: ownerRef}}

	csClient := fake.NewClientBuilder().WithObjects(awsCluster, &roleIdentity1, &roleIdentity2, &roleIdentity3).Build()
	reconciler := &AWSClusterRoleIdentityReconciler{
		Client: csClient,
	}

	// Calling reconcile should not error and not requeue the request with insufficient set up
	result, err := reconciler.Reconcile(ctx, ctrl.Request{
		NamespacedName: client.ObjectKey{
			Namespace: awsCluster.Namespace,
			Name:      awsCluster.Name,
		},
	})
	g.Expect(err).To(BeNil())
	g.Expect(result).To(BeZero())

	unstructuredRoleIdentity1, err := getUnstructuredFromObjectReference(ctx, csClient, roleIdentity1.Kind, roleIdentity1.Name)
	g.Expect(err).To(BeNil())
	unstructuredRoleIdentity2, err := getUnstructuredFromObjectReference(ctx, csClient, roleIdentity2.Kind, roleIdentity2.Name)
	g.Expect(err).To(BeNil())
	unstructuredRoleIdentity3, err := getUnstructuredFromObjectReference(ctx, csClient, roleIdentity3.Kind, roleIdentity3.Name)
	g.Expect(err).To(BeNil())
	g.Expect(unstructuredRoleIdentity1.GetOwnerReferences()).NotTo(ConsistOf(ownerRef))
	g.Expect(unstructuredRoleIdentity2.GetOwnerReferences()).To(ConsistOf(ownerRef))
	g.Expect(unstructuredRoleIdentity3.GetOwnerReferences()).NotTo(ConsistOf(ownerRef))
}
