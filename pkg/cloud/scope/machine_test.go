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
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
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
				clusterv1.ClusterNameLabel: clusterName,
			},
			Name:      machineName,
			Namespace: "default",
		},
		Spec: clusterv1.MachineSpec{
			Bootstrap: clusterv1.Bootstrap{
				DataSecretName: pointer.String(machineName),
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
				clusterv1.ClusterNameLabel: clusterName,
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
				clusterv1.ClusterNameLabel: clusterName,
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

	initObjects := []client.Object{
		cluster, machine, secret, awsMachine, awsCluster,
	}

	client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(initObjects...).Build()
	return NewMachineScope(
		MachineScopeParams{
			Client:  client,
			Machine: machine,
			Cluster: cluster,
			InfraCluster: &ClusterScope{
				AWSCluster: awsCluster,
			},
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

func TestGetRawBootstrapDataWithFormat(t *testing.T) {
	t.Run("returns_empty_format_when_format_is_not_set_in_bootstrap_data", func(t *testing.T) {
		scope, err := setupMachineScope()
		if err != nil {
			t.Fatal(err)
		}

		_, format, err := scope.GetRawBootstrapDataWithFormat()
		if err != nil {
			t.Fatalf("Getting raw bootstrap data with format: %v", err)
		}

		if format != "" {
			t.Fatalf("Fromat should be empty when it's not defined in bootstrap data, got: %q", format)
		}
	})

	t.Run("returns_format_defined_in_bootstrap_data_when_available", func(t *testing.T) {
		scheme, err := setupScheme()
		if err != nil {
			t.Fatalf("Configuring schema: %v", err)
		}

		clusterName := "my-cluster"
		machineName := "my-machine-0"
		cluster := newCluster(clusterName)
		machine := newMachine(clusterName, machineName)
		awsMachine := newAWSMachine(clusterName, machineName)
		awsCluster := newAWSCluster(clusterName)

		expectedBootstrapDataFormat := "ignition"

		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					clusterv1.ClusterNameLabel: clusterName,
				},
				Name:      machineName,
				Namespace: "default",
			},
			Data: map[string][]byte{
				"value":  []byte("user data"),
				"format": []byte(expectedBootstrapDataFormat),
			},
		}

		initObjects := []client.Object{
			cluster, machine, secret, awsMachine, awsCluster,
		}

		client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(initObjects...).Build()

		machineScope, err := NewMachineScope(
			MachineScopeParams{
				Client:  client,
				Machine: machine,
				Cluster: cluster,
				InfraCluster: &ClusterScope{
					AWSCluster: awsCluster,
				},
				AWSMachine: awsMachine,
			},
		)
		if err != nil {
			t.Fatalf("Creating machine scope: %v", err)
		}

		_, format, err := machineScope.GetRawBootstrapDataWithFormat()
		if err != nil {
			t.Fatalf("Getting raw bootstrap data with format: %v", err)
		}

		if format != expectedBootstrapDataFormat {
			t.Fatalf("Unexpected bootstrap data format, expected %q, got %q", expectedBootstrapDataFormat, format)
		}
	})
}

func TestUseSecretsManagerTrue(t *testing.T) {
	scope, err := setupMachineScope()
	if err != nil {
		t.Fatal(err)
	}

	if !scope.UseSecretsManager("cloud-config") {
		t.Fatalf("UseSecretsManager should be true")
	}
}

func TestUseIgnition(t *testing.T) {
	t.Run("returns_true_when_given_bootstrap_data_format_is_ignition", func(t *testing.T) {
		scope, err := setupMachineScope()
		if err != nil {
			t.Fatal(err)
		}

		if !scope.UseIgnition("ignition") {
			t.Fatalf("UseIgnition should be true")
		}
	})

	// To retain backward compatibility, where KAPBK does not produce format field.
	t.Run("returns_false_when_given_bootstrap_data_format_is_empty", func(t *testing.T) {
		scope, err := setupMachineScope()
		if err != nil {
			t.Fatal(err)
		}

		if scope.UseIgnition("") {
			t.Fatalf("UseIgnition should be false")
		}
	})
}

func TestCompressUserData(t *testing.T) {
	// Ignition does not support compressed data in S3.
	t.Run("returns_false_when_bootstrap_data_is_in_ignition_format", func(t *testing.T) {
		scope, err := setupMachineScope()
		if err != nil {
			t.Fatal(err)
		}

		if scope.CompressUserData("ignition") {
			t.Fatalf("User data would be compressed despite Ignition format")
		}
	})
}

func TestGetSecretARNDefaultIsNil(t *testing.T) {
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
	if err != nil {
		t.Fatal(err)
	}

	scope.SetSecretPrefix(prefix)
	if val := scope.GetSecretPrefix(); val != prefix {
		t.Fatalf("prefix does not equal %s: %s", prefix, val)
	}
}

func TestSetProviderID(t *testing.T) {
	scope, err := setupMachineScope()
	if err != nil {
		t.Fatal(err)
	}

	scope.SetProviderID("test-id", "test-zone-1a")
	providerID := *scope.AWSMachine.Spec.ProviderID
	expectedProviderID := "aws:///test-zone-1a/test-id"
	if providerID != expectedProviderID {
		t.Fatalf("Expected providerID %s, got %s", expectedProviderID, providerID)
	}

	scope.SetProviderID("test-id", "")
	providerID = *scope.AWSMachine.Spec.ProviderID
	expectedProviderID = "aws:////test-id"
	if providerID != expectedProviderID {
		t.Fatalf("Expected providerID %s, got %s", expectedProviderID, providerID)
	}
}
