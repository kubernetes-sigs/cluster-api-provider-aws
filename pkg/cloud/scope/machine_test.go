/*
Copyright 2020 The Kubernetes Authors.

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

package scope

import (
	"encoding/base64"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha2"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

const clusterLabelName = "cluster.x-k8s.io/cluster-name"

func setupScheme() (*runtime.Scheme, error) {
	scheme := runtime.NewScheme()
	if err := infrav1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := clusterv1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := corev1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	return scheme, nil
}

func newMachine(clusterName, machineName string) *clusterv1.Machine {
	return &clusterv1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				clusterLabelName: clusterName,
			},
			ClusterName: clusterName,
			Name:        machineName,
			Namespace:   "default",
		},
		Spec: clusterv1.MachineSpec{
			Bootstrap: clusterv1.Bootstrap{
				Data: pointer.StringPtr(base64.StdEncoding.EncodeToString([]byte("some base 64 data"))),
			},
		},
	}
}

func newCluster(name string) *clusterv1.Cluster {
	return &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
	}
}

func newAWSCluster(name string) *infrav1.AWSCluster {
	return &infrav1.AWSCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
	}
}

func newAWSMachine(clusterName, machineName string) *infrav1.AWSMachine {
	return &infrav1.AWSMachine{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				clusterLabelName: clusterName,
			},
			Name:      machineName,
			Namespace: "default",
		},
	}
}

func setupMachineScope() (*MachineScope, error) {
	scheme, err := setupScheme()
	if err != nil {
		return nil, err
	}
	clusterName := "my-cluster"
	cluster := newCluster(clusterName)
	machine := newMachine(clusterName, "my-machine-0")
	awsMachine := newAWSMachine(clusterName, "my-machine-0")
	awsCluster := newAWSCluster(clusterName)

	initObjects := []runtime.Object{
		cluster, machine, awsMachine, awsCluster,
	}

	client := fake.NewFakeClientWithScheme(scheme, initObjects...)
	return NewMachineScope(
		MachineScopeParams{
			Client:     client,
			Machine:    machine,
			Cluster:    cluster,
			AWSCluster: awsCluster,
			AWSMachine: awsMachine,
		},
	)
}

func TestGetBootstrapDataIsBase64Encoded(t *testing.T) {
	scope, err := setupMachineScope()
	if err != nil {
		t.Fatal(err)
	}

	userdata, err := scope.GetBootstrapData()
	if err != nil {
		t.Fatal(err)
	}
	_, err = base64.StdEncoding.DecodeString(userdata)
	if err != nil {
		t.Fatalf("GetBootstrapData isn't base 64 encoded: %+v", err)
	}
}

func TestGetRawBootstrapDataIsNotBase64Encoded(t *testing.T) {
	scope, err := setupMachineScope()
	if err != nil {
		t.Fatal(err)
	}

	userdata, err := scope.GetRawBootstrapData()
	if err != nil {
		t.Fatal(err)
	}
	_, err = base64.StdEncoding.DecodeString(string(userdata))
	if err == nil {
		t.Fatalf("GetBootstrapData is base 64 encoded: %+v", err)
	}
}

func TestUseSecretsManagerFalse(t *testing.T) {
	scope, err := setupMachineScope()
	if err != nil {
		t.Fatal(err)
	}

	if scope.UseSecretsManager() {
		t.Fatalf("UseSecretsManager should be false by default")
	}
}

func TestGetSecretPrefixDefaultIsNil(t *testing.T) {
	scope, err := setupMachineScope()
	if err != nil {
		t.Fatal(err)
	}

	if scope.GetSecretPrefix() != "" {
		t.Fatalf("GetSecretPrefix should be empty string")
	}
}

func TestSetSecretARN(t *testing.T) {
	prefix := "secret"
	scope, err := setupMachineScope()
	scope.AWSMachine.Spec.CloudInit = &infrav1.CloudInit{
		EnableSecureSecretsManager: true,
	}
	if err != nil {
		t.Fatal(err)
	}

	scope.SetSecretPrefix(prefix)
	val := scope.GetSecretPrefix()
	if val != prefix {
		t.Fatalf("prefix does not equal %s: %s", prefix, val)
	}
}
