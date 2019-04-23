/*
Copyright 2018 The Kubernetes Authors.

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

package machine

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	"sigs.k8s.io/cluster-api/pkg/controller/machine"
)

var (
	_ machine.Actuator = (*Actuator)(nil)
)

func contains(s []*clusterv1.Machine, e clusterv1.Machine) bool {
	exists := false
	for _, em := range s {
		if em.Name == e.Name && em.Namespace == e.Namespace {
			exists = true
			break
		}
	}
	return exists
}

func TestGetControlPlaneMachines(t *testing.T) {
	testCases := []struct {
		name        string
		input       *clusterv1.MachineList
		expectedOut []clusterv1.Machine
	}{
		{
			name: "0 machines",
			input: &clusterv1.MachineList{
				Items: []clusterv1.Machine{},
			},
			expectedOut: []clusterv1.Machine{},
		},
		{
			name: "only 2 controlplane machines",
			input: &clusterv1.MachineList{
				Items: []clusterv1.Machine{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "master-0",
							Namespace: "awesome-ns",
						},
						Spec: clusterv1.MachineSpec{
							Versions: clusterv1.MachineVersionInfo{
								Kubelet:      "v1.13.0",
								ControlPlane: "v1.13.0",
							},
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "master-1",
							Namespace: "awesome-ns",
						},
						Spec: clusterv1.MachineSpec{
							Versions: clusterv1.MachineVersionInfo{
								Kubelet:      "v1.13.0",
								ControlPlane: "v1.13.0",
							},
						},
					},
				},
			},
			expectedOut: []clusterv1.Machine{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "master-0",
						Namespace: "awesome-ns",
					},
					Spec: clusterv1.MachineSpec{
						Versions: clusterv1.MachineVersionInfo{
							Kubelet:      "v1.13.0",
							ControlPlane: "v1.13.0",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "master-1",
						Namespace: "awesome-ns",
					},
					Spec: clusterv1.MachineSpec{
						Versions: clusterv1.MachineVersionInfo{
							Kubelet:      "v1.13.0",
							ControlPlane: "v1.13.0",
						},
					},
				},
			},
		},
		{
			name: "2 controlplane machines and 2 worker machines",
			input: &clusterv1.MachineList{
				Items: []clusterv1.Machine{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "master-0",
							Namespace: "awesome-ns",
						},
						Spec: clusterv1.MachineSpec{
							Versions: clusterv1.MachineVersionInfo{
								Kubelet:      "v1.13.0",
								ControlPlane: "v1.13.0",
							},
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "master-1",
							Namespace: "awesome-ns",
						},
						Spec: clusterv1.MachineSpec{
							Versions: clusterv1.MachineVersionInfo{
								Kubelet:      "v1.13.0",
								ControlPlane: "v1.13.0",
							},
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "worker-0",
							Namespace: "awesome-ns",
						},
						Spec: clusterv1.MachineSpec{
							Versions: clusterv1.MachineVersionInfo{
								Kubelet: "v1.13.0",
							},
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "worker-1",
							Namespace: "awesome-ns",
						},
						Spec: clusterv1.MachineSpec{
							Versions: clusterv1.MachineVersionInfo{
								Kubelet: "v1.13.0",
							},
						},
					},
				},
			},
			expectedOut: []clusterv1.Machine{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "master-0",
						Namespace: "awesome-ns",
					},
					Spec: clusterv1.MachineSpec{
						Versions: clusterv1.MachineVersionInfo{
							Kubelet:      "v1.13.0",
							ControlPlane: "v1.13.0",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "master-1",
						Namespace: "awesome-ns",
					},
					Spec: clusterv1.MachineSpec{
						Versions: clusterv1.MachineVersionInfo{
							Kubelet:      "v1.13.0",
							ControlPlane: "v1.13.0",
						},
					},
				}},
		},
		{
			name: "only 2 worker machines",
			input: &clusterv1.MachineList{
				Items: []clusterv1.Machine{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "worker-0",
							Namespace: "awesome-ns",
						},
						Spec: clusterv1.MachineSpec{
							Versions: clusterv1.MachineVersionInfo{
								Kubelet: "v1.13.0",
							},
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "worker-1",
							Namespace: "awesome-ns",
						},
						Spec: clusterv1.MachineSpec{
							Versions: clusterv1.MachineVersionInfo{
								Kubelet: "v1.13.0",
							},
						},
					},
				},
			},
			expectedOut: []clusterv1.Machine{},
		},
	}

	for _, tc := range testCases {
		actual := GetControlPlaneMachines(tc.input)
		if len(actual) != len(tc.expectedOut) {
			t.Fatalf("[%s] Unexpected number of controlplane machines returned. Got: %d, Want: %d", tc.name, len(actual), len(tc.expectedOut))
		}
		if len(tc.expectedOut) > 1 {
			for _, em := range tc.expectedOut {
				if !contains(actual, em) {
					t.Fatalf("[%s] Expected controlplane machine %q in namespace %q not found", tc.name, em.Name, em.Namespace)
				}
			}
		}
	}
}

func TestMachineEqual(t *testing.T) {
	testCases := []struct {
		name          string
		inM1          clusterv1.Machine
		inM2          clusterv1.Machine
		expectedEqual bool
	}{
		{
			name: "machines are equal",
			inM1: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "machine1",
					Namespace: "my-awesome-ns",
				},
			},
			inM2: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "machine1",
					Namespace: "my-awesome-ns",
				},
			},
			expectedEqual: true,
		},
		{
			name: "machines are not equal: names are different",
			inM1: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "machine 1",
					Namespace: "my-awesome-ns",
				},
			},
			inM2: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "machine 2",
					Namespace: "my-awsesome-ns",
				},
			},
			expectedEqual: false,
		},
		{
			name: "machines are not equal: namespace are different",
			inM1: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "machine1",
					Namespace: "my-awesome-ns",
				},
			},
			inM2: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "machine1",
					Namespace: "your-awsesome-ns",
				},
			},
			expectedEqual: false,
		},
	}

	for _, tc := range testCases {
		actualEqual := machinesEqual(&tc.inM1, &tc.inM2)
		if tc.expectedEqual {
			if !actualEqual {
				t.Fatalf("[%s] Expected Machine1 [Name:%q, Namespace:%q], Equal Machine2 [Name:%q, Namespace:%q]",
					tc.name, tc.inM1.Name, tc.inM1.Namespace, tc.inM2.Name, tc.inM2.Namespace)
			}
		} else {
			if actualEqual {
				t.Fatalf("[%s] Expected Machine1 [Name:%q, Namespace:%q], NOT Equal Machine2 [Name:%q, Namespace:%q]",
					tc.name, tc.inM1.Name, tc.inM1.Namespace, tc.inM2.Name, tc.inM2.Namespace)
			}
		}
	}
}

func TestImmutableStateChange(t *testing.T) {
	testCases := []struct {
		name        string
		machineSpec v1alpha1.AWSMachineProviderSpec
		instance    v1alpha1.Instance
		// expected length of returned errors
		expected int
	}{
		{
			name: "instance type is unchanged",
			machineSpec: v1alpha1.AWSMachineProviderSpec{
				InstanceType: "t2.micro",
			},
			instance: v1alpha1.Instance{
				Type: "t2.micro",
			},
			expected: 0,
		},
		{
			name: "instance type is changed",
			machineSpec: v1alpha1.AWSMachineProviderSpec{
				InstanceType: "m5.large",
			},
			instance: v1alpha1.Instance{
				Type: "t2.micro",
			},
			expected: 1,
		},
		{
			name: "iam profile is unchanged",
			machineSpec: v1alpha1.AWSMachineProviderSpec{
				IAMInstanceProfile: "test-profile",
			},
			instance: v1alpha1.Instance{
				IAMProfile: "test-profile",
			},
			expected: 0,
		},
		{
			name: "iam profile is changed",
			machineSpec: v1alpha1.AWSMachineProviderSpec{
				IAMInstanceProfile: "test-profile-updated",
			},
			instance: v1alpha1.Instance{
				IAMProfile: "test-profile",
			},
			expected: 1,
		},
		{
			name: "keyname is unchanged",
			machineSpec: v1alpha1.AWSMachineProviderSpec{
				KeyName: "SSHKey",
			},
			instance: v1alpha1.Instance{
				KeyName: aws.String("SSHKey"),
			},
			expected: 0,
		},
		{
			name: "keyname is changed",
			machineSpec: v1alpha1.AWSMachineProviderSpec{
				KeyName: "SSHKey2",
			},
			instance: v1alpha1.Instance{
				KeyName: aws.String("SSHKey"),
			},
			expected: 1,
		},
		{
			name: "instance with public ip is unchanged",
			machineSpec: v1alpha1.AWSMachineProviderSpec{
				PublicIP: aws.Bool(true),
			},
			instance: v1alpha1.Instance{
				// This IP chosen from RFC5737 TEST-NET-1
				PublicIP: aws.String("192.0.2.1"),
			},
			expected: 0,
		},
		{
			name: "instance with public ip is changed",
			machineSpec: v1alpha1.AWSMachineProviderSpec{
				PublicIP: aws.Bool(false),
			},
			instance: v1alpha1.Instance{
				// This IP chosen from RFC5737 TEST-NET-1
				PublicIP: aws.String("192.0.2.1"),
			},
			expected: 1,
		},
		{
			name: "instance without public ip is unchanged",
			machineSpec: v1alpha1.AWSMachineProviderSpec{
				PublicIP: aws.Bool(false),
			},
			instance: v1alpha1.Instance{
				PublicIP: aws.String(""),
			},
			expected: 0,
		},
		{
			name: "instance without public ip is changed",
			machineSpec: v1alpha1.AWSMachineProviderSpec{
				PublicIP: aws.Bool(true),
			},
			instance: v1alpha1.Instance{
				PublicIP: aws.String(""),
			},
			expected: 1,
		},
		{
			name: "subnetid is unchanged",
			machineSpec: v1alpha1.AWSMachineProviderSpec{
				Subnet: &v1alpha1.AWSResourceReference{
					ID: aws.String("subnet-abcdef"),
				},
			},
			instance: v1alpha1.Instance{
				SubnetID: "subnet-abcdef",
			},
			expected: 0,
		},
		{
			name: "subnetid is changed",
			machineSpec: v1alpha1.AWSMachineProviderSpec{
				Subnet: &v1alpha1.AWSResourceReference{
					ID: aws.String("subnet-123456"),
				},
			},
			instance: v1alpha1.Instance{
				SubnetID: "subnet-abcdef",
			},
			expected: 1,
		},
		{
			name: "root device size is unchanged",
			machineSpec: v1alpha1.AWSMachineProviderSpec{
				RootDeviceSize: 12,
			},
			instance: v1alpha1.Instance{
				RootDeviceSize: 12,
			},
			expected: 0,
		},
		{
			name: "root device size is changed",
			machineSpec: v1alpha1.AWSMachineProviderSpec{
				RootDeviceSize: 12,
			},
			instance: v1alpha1.Instance{
				RootDeviceSize: 16,
			},
			expected: 1,
		},
		{
			name: "multiple immutable changes",
			machineSpec: v1alpha1.AWSMachineProviderSpec{
				IAMInstanceProfile: "test-profile-updated",
				PublicIP:           aws.Bool(false),
			},
			instance: v1alpha1.Instance{
				IAMProfile: "test-profile",
				// This IP chosen from RFC5737 TEST-NET-1
				PublicIP: aws.String("192.0.2.1"),
			},
			expected: 2,
		},
	}

	testActuator := NewActuator(ActuatorParams{})

	for _, tc := range testCases {
		changed := len(testActuator.isMachineOutdated(&tc.machineSpec, &tc.instance))

		if tc.expected != changed {
			t.Fatalf("[%s] Expected MachineSpec [%+v], NOT Equal Instance [%+v]",
				tc.name, tc.machineSpec, tc.instance)
		}
	}
}
