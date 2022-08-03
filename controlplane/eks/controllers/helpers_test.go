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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func getAWSManagedControlPlane(name, namespace string) ekscontrolplanev1.AWSManagedControlPlane {
	return ekscontrolplanev1.AWSManagedControlPlane{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
			Region: "us-east-1",
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
						AvailabilityZone: "us-east-1c",
						CidrBlock:        "10.0.11.0/24",
						IsPublic:         true,
					},
				},
				SecurityGroupOverrides: map[infrav1.SecurityGroupRole]string{},
			},
			Bastion: infrav1.Bastion{Enabled: true},
		},
	}
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
