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

package ec2

import (
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func setupScheme() (*runtime.Scheme, error) {
	scheme := runtime.NewScheme()
	if err := clusterv1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := corev1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := infrav1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := expinfrav1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	return scheme, nil
}

func newCluster() *clusterv1.Cluster {
	return &clusterv1.Cluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Cluster",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cluster-name",
			Namespace: "cluster-ns",
		},
		Spec: clusterv1.ClusterSpec{},
	}
}

func newAWSCluster() *infrav1.AWSCluster {
	return &infrav1.AWSCluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AWSCluster",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "aws-cluster-name",
			Namespace: "aws-cluster-ns",
		},
		Spec: infrav1.AWSClusterSpec{},
		Status: infrav1.AWSClusterStatus{
			Ready: true,
			Network: infrav1.NetworkStatus{
				SecurityGroups: map[infrav1.SecurityGroupRole]infrav1.SecurityGroup{},
			},
		},
	}
}

func setupClusterScope(g *WithT) *scope.ClusterScope {
	scheme, err := setupScheme()
	g.Expect(err).NotTo(HaveOccurred())

	awsCluster := newAWSCluster()
	client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(awsCluster).Build()

	cs, err := scope.NewClusterScope(scope.ClusterScopeParams{
		Client:     client,
		Cluster:    newCluster(),
		AWSCluster: awsCluster,
	})
	g.Expect(err).NotTo(HaveOccurred())
	return cs
}
