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
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

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
				clusterv1.ClusterLabelName: clusterName,
			},
			ClusterName: clusterName,
			Name:        machineName,
			Namespace:   "default",
		},
		Spec: clusterv1.MachineSpec{
			Bootstrap: clusterv1.Bootstrap{
				DataSecretName: pointer.StringPtr(machineName),
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
				clusterv1.ClusterLabelName: clusterName,
			},
			Name:      machineName,
			Namespace: "default",
		},
	}
}

func newBootstrapSecret(clusterName, machineName string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				clusterv1.ClusterLabelName: clusterName,
			},
			Name:      machineName,
			Namespace: "default",
		},
		Data: map[string][]byte{
			"value": []byte("user data"),
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
	secret := newBootstrapSecret(clusterName, "my-machine-0")
	awsMachine := newAWSMachine(clusterName, "my-machine-0")
	awsCluster := newAWSCluster(clusterName)

	initObjects := []runtime.Object{
		cluster, machine, secret, awsMachine, awsCluster,
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

func TestUseSecretsManagerTrue(t *testing.T) {
	scope, err := setupMachineScope()
	if err != nil {
		t.Fatal(err)
	}

	if !scope.UseSecretsManager() {
		t.Fatalf("UseSecretsManager should be true")
	}
}

func TestGetSecretARNDefaultIsNil(t *testing.T) {
	scope, err := setupMachineScope()
	if err != nil {
		t.Fatal(err)
	}

	if scope.GetSecretARN() != "" {
		t.Fatalf("GetSecretARN should be empty string")
	}
}

func TestSetSecretARN(t *testing.T) {
	secretARN := "secretARN"
	scope, err := setupMachineScope()
	if err != nil {
		t.Fatal(err)
	}

	scope.SetSecretARN(secretARN)
	val := scope.GetSecretARN()

	if val != secretARN {
		t.Fatalf("GetSecretARN does not equal %s: %s", secretARN, val)
	}
}
