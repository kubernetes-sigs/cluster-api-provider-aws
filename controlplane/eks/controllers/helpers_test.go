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
package controllers

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta1")

func getAWSManagedControlPlaneScope(cluster *clusterv1.Cluster, awsManagedControlPlane *ekscontrolplanev1.AWSManagedControlPlane) *scope.ManagedControlPlaneScope {
	scope, err := scope.NewManagedControlPlaneScope(
		scope.ManagedControlPlaneScopeParams{
			Client:                    testEnv.Client,
			Cluster:                   cluster,
			ControlPlane:              awsManagedControlPlane,
			EnableIAM:                 true,
			MaxWaitActiveUpdateDelete: maxActiveUpdateDeleteWait,
		},
	)
	utilruntime.Must(err)
	return scope
}

func getManagedClusterObjects(name, namespace string) (clusterv1.Cluster, infrav1.AWSManagedCluster, ekscontrolplanev1.AWSManagedControlPlane) {
	cluster := clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			UID:       "1",
		},
		Spec: clusterv1.ClusterSpec{
			ControlPlaneRef: &corev1.ObjectReference{
				APIVersion: ekscontrolplanev1.GroupVersion.String(),
				Name:       name,
				Kind:       "AWSManagedControlPlane",
				Namespace:  namespace,
			},
			InfrastructureRef: &corev1.ObjectReference{
				APIVersion: infrav1.GroupVersion.String(),
				Name:       name,
				Kind:       "AWSManagedCluster",
				Namespace:  namespace,
			},
		},
	}
	awsManagedCluster := infrav1.AWSManagedCluster{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
	}
	awsManagedControlPlane := ekscontrolplanev1.AWSManagedControlPlane{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: clusterv1.GroupVersion.String(),
					Kind:       "Cluster",
					Name:       cluster.Name,
					UID:        "1",
				},
			},
		},
		Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
			EKSClusterName: name,
			Region:         "us-east-1",
			NetworkSpec: infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID:        "vpc-exists",
					CidrBlock: "10.0.0.0/8",
				},
				Subnets: infrav1.Subnets{
					{
						ID:               "subnet-1",
						AvailabilityZone: "us-east-1a",
						CidrBlock:        "10.0.10.0/24",
						IsPublic:         false,
					},
					{
						ID:               "subnet-2",
						AvailabilityZone: "us-east-1b",
						CidrBlock:        "10.0.11.0/24",
						IsPublic:         true,
					},
					{
						ID:               "subnet-3",
						AvailabilityZone: "us-east-1c",
						CidrBlock:        "10.0.12.0/24",
						IsPublic:         true,
					},
				},
				SecurityGroupOverrides: map[infrav1.SecurityGroupRole]string{},
			},
			Bastion: infrav1.Bastion{Enabled: true},
		},
	}
	return cluster, awsManagedCluster, awsManagedControlPlane
}

func getManagedControlPlaneScope(cp ekscontrolplanev1.AWSManagedControlPlane) (*scope.ManagedControlPlaneScope, error) {
	scheme := runtime.NewScheme()
	_ = ekscontrolplanev1.AddToScheme(scheme)
	_ = infrav1.AddToScheme(scheme)
	client := fake.NewClientBuilder().WithScheme(scheme).Build()

	return scope.NewManagedControlPlaneScope(
		scope.ManagedControlPlaneScopeParams{
			Client: client,
			Cluster: &clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-cluster",
				},
			},
			ControlPlane: &cp,
		},
	)
}
