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
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/exp/api/v1beta1"
)

func setupClusterScope(cl client.Client) (*scope.ClusterScope, error) {
	return scope.NewClusterScope(scope.ClusterScopeParams{
		Client:     cl,
		Cluster:    newCluster(),
		AWSCluster: newAWSCluster(),
	})
}

func setupNewManagedControlPlaneScope(cl client.Client) (*scope.ManagedControlPlaneScope, error) {
	return scope.NewManagedControlPlaneScope(scope.ManagedControlPlaneScopeParams{
		Client:       cl,
		Cluster:      newCluster(),
		ControlPlane: newAWSManagedControlPlane(),
	})
}

func setupMachinePoolScope(cl client.Client, ec2Scope scope.EC2Scope) (*scope.MachinePoolScope, error) {
	return scope.NewMachinePoolScope(scope.MachinePoolScopeParams{
		Client:         cl,
		InfraCluster:   ec2Scope,
		Cluster:        newCluster(),
		MachinePool:    newMachinePool(),
		AWSMachinePool: newAWSMachinePool(),
	})
}

func defaultEC2Tags(name, clusterName string) []*ec2.Tag {
	return []*ec2.Tag{
		{
			Key:   aws.String("Name"),
			Value: aws.String(name),
		},
		{
			Key:   aws.String(infrav1.ClusterAWSCloudProviderTagKey(clusterName)),
			Value: aws.String("owned"),
		},
		{
			Key:   aws.String(infrav1.ClusterTagKey(clusterName)),
			Value: aws.String("owned"),
		},
		{
			Key:   aws.String(infrav1.NameAWSClusterAPIRole),
			Value: aws.String("node"),
		},
	}
}

func newAWSMachinePool() *expinfrav1.AWSMachinePool {
	return &expinfrav1.AWSMachinePool{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AWSMachinePool",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "aws-mp-name",
			Namespace: "aws-mp-ns",
		},
		Spec: expinfrav1.AWSMachinePoolSpec{
			AvailabilityZones: []string{"us-east-1"},
			AdditionalTags:    infrav1.Tags{},
			AWSLaunchTemplate: expinfrav1.AWSLaunchTemplate{
				Name:               "aws-launch-template",
				IamInstanceProfile: "instance-profile",
				AMI:                infrav1.AMIReference{},
				InstanceType:       "t3.large",
				SSHKeyName:         aws.String("default"),
				SpotMarketOptions:  &infrav1.SpotMarketOptions{MaxPrice: aws.String("0.9")},
			},
		},
		Status: expinfrav1.AWSMachinePoolStatus{
			LaunchTemplateID: "launch-template-id",
		},
	}
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
		Spec: infrav1.AWSClusterSpec{
			ImageLookupFormat: "img-lookup-format",
			ImageLookupBaseOS: "img-lookup-os",
			ImageLookupOrg:    "img-lookup-org",
		},
		Status: infrav1.AWSClusterStatus{
			Ready: true,
			Network: infrav1.NetworkStatus{
				SecurityGroups: map[infrav1.SecurityGroupRole]infrav1.SecurityGroup{
					infrav1.SecurityGroupNode: {ID: "nodeSG"},
					infrav1.SecurityGroupLB:   {ID: "lbSG"},
				},
			},
		},
	}
}

func newAWSManagedControlPlane() *ekscontrolplanev1.AWSManagedControlPlane {
	return &ekscontrolplanev1.AWSManagedControlPlane{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AWSManagedControlPlane",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "aws-cluster-name",
			Namespace: "aws-cluster-ns",
		},
	}
}

func newMachinePool() *v1beta1.MachinePool {
	return &v1beta1.MachinePool{
		TypeMeta: metav1.TypeMeta{
			Kind:       "MachinePool",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "mp",
		},
		Spec: v1beta1.MachinePoolSpec{
			Template: clusterv1.MachineTemplateSpec{
				Spec: clusterv1.MachineSpec{
					Version: pointer.String("v1.23.3"),
				},
			},
		},
	}
}

func sortTags(a []*ec2.Tag) {
	sort.Slice(a, func(i, j int) bool {
		return *(a[i].Key) < *(a[j].Key)
	})
}

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
	if err := ekscontrolplanev1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := v1beta1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	return scheme, nil
}
