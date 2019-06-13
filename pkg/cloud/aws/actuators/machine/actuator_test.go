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
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/golang/mock/gomock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/filter"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/ec2/mock_ec2iface"
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
			name: "2 controlplane machines, 1 deleted",
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
							Name:      "master-2",
							Namespace: "awesome-ns",
							DeletionTimestamp: &metav1.Time{
								Time: time.Now(),
							},
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
			name:        "root device size is omitted",
			machineSpec: v1alpha1.AWSMachineProviderSpec{},
			instance: v1alpha1.Instance{
				// All instances have a root device size, even when we don't set one
				RootDeviceSize: 12,
			},
			expected: 0,
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

func getWorkerTestMachine() *clusterv1.Machine {
	return &clusterv1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "worker-0",
			Namespace: "awesome-ns",
			Labels: map[string]string{
				"set": "node",
			},
		},
		Spec: clusterv1.MachineSpec{
			Versions: clusterv1.MachineVersionInfo{
				Kubelet: "v1.13.0",
			},
		},
	}
}

func getControlplaneTestMachine() *clusterv1.Machine {
	return &clusterv1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "controlplane-0",
			Namespace: "awesome-ns",
			Labels: map[string]string{
				"set": "controlplane",
			},
		},
		Spec: clusterv1.MachineSpec{
			Versions: clusterv1.MachineVersionInfo{
				Kubelet: "v1.13.0",
			},
		},
	}
}

func getRandomTestMachine() *clusterv1.Machine {
	return &clusterv1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "random-0",
			Namespace: "awesome-ns",
			Labels: map[string]string{
				"set": "random",
			},
		},
	}
}

func getNoLabelTestMachine() *clusterv1.Machine {
	return &clusterv1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nolabel-0",
			Namespace: "awesome-ns",
		},
	}
}

func getTestMachineScope(m *clusterv1.Machine, t *testing.T, mockEC2 ec2iface.EC2API) *actuators.MachineScope {
	testCluster := &clusterv1.Cluster{}
	scope, err := actuators.NewMachineScope(actuators.MachineScopeParams{
		Machine: m,
		Cluster: testCluster,
		Client:  nil,
		Logger:  nil,
		AWSClients: actuators.AWSClients{
			EC2: mockEC2,
		},
	},
	)
	if err != nil {
		t.Fatalf("[getTestMachineScope] failed, Got %q; Want:nil", err)
	}
	return scope
}

func getWorkerMachineScope(t *testing.T, mockEC2 ec2iface.EC2API) *actuators.MachineScope {
	return getTestMachineScope(getWorkerTestMachine(), t, mockEC2)
}

func getControlplaneMachineScope(t *testing.T, mockEC2 ec2iface.EC2API) *actuators.MachineScope {
	return getTestMachineScope(getControlplaneTestMachine(), t, mockEC2)
}

func getRandoMachineScope(t *testing.T, mockEC2 ec2iface.EC2API) *actuators.MachineScope {
	return getTestMachineScope(getRandomTestMachine(), t, mockEC2)
}

func getInvalidMachineScope(t *testing.T) *actuators.MachineScope {
	is := getTestMachineScope(getControlplaneTestMachine(), t, nil)
	is.Cluster = nil
	return is
}

func getNoLabelMachineScope(t *testing.T, mockEC2 ec2iface.EC2API) *actuators.MachineScope {
	return getTestMachineScope(getNoLabelTestMachine(), t, mockEC2)
}

func getTestControlplaneMachines() []*clusterv1.Machine {
	return []*clusterv1.Machine{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "master-0",
				Namespace: "awesome-ns",
				Labels: map[string]string{
					"set": "controlplane",
				},
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
				Labels: map[string]string{
					"set": "controlplane",
				},
			},
			Spec: clusterv1.MachineSpec{
				Versions: clusterv1.MachineVersionInfo{
					Kubelet:      "v1.13.0",
					ControlPlane: "v1.13.0",
				},
			},
		},
	}
}

func getMockEC2APIDescribeInstancesNotFound(ne []string, mockCtrl *gomock.Controller) ec2iface.EC2API {
	mockEC2 := mock_ec2iface.NewMockEC2API(mockCtrl)
	for _, n := range ne {
		dmi := &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				filter.EC2.VPC(""),
				filter.EC2.ClusterOwned(""),
				filter.EC2.Name(n),
				filter.EC2.InstanceStates(ec2.InstanceStateNamePending, ec2.InstanceStateNameRunning),
			},
		}
		mockEC2.EXPECT().DescribeInstances(dmi).Return(
			nil,
			awserrors.NewNotFound(fmt.Errorf("mocked not found"))).Times(1)
	}

	return mockEC2
}

func getMockEC2APIDescribeInstancesFail(machineNames []string, mockCtrl *gomock.Controller) ec2iface.EC2API {
	mockEC2 := mock_ec2iface.NewMockEC2API(mockCtrl)
	for _, mn := range machineNames {
		dmi := &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				filter.EC2.VPC(""),
				filter.EC2.ClusterOwned(""),
				filter.EC2.Name(mn),
				filter.EC2.InstanceStates(ec2.InstanceStateNamePending, ec2.InstanceStateNameRunning),
			},
		}

		mockEC2.EXPECT().DescribeInstances(dmi).Return(
			nil,
			fmt.Errorf("mock API failure")).AnyTimes()
	}

	return mockEC2
}

func getMockEC2APIDescribeInstancesPass(machineNames []string, mockCtrl *gomock.Controller) ec2iface.EC2API {
	mockEC2 := mock_ec2iface.NewMockEC2API(mockCtrl)
	for _, mn := range machineNames {
		dmi := &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				filter.EC2.VPC(""),
				filter.EC2.ClusterOwned(""),
				filter.EC2.Name(mn),
				filter.EC2.InstanceStates(ec2.InstanceStateNamePending, ec2.InstanceStateNameRunning),
			},
		}

		mockEC2.EXPECT().DescribeInstances(dmi).Return(
			&ec2.DescribeInstancesOutput{
				Reservations: []*ec2.Reservation{
					{
						Instances: []*ec2.Instance{
							{
								InstanceId: aws.String(mn),
								State: &ec2.InstanceState{
									Code: aws.Int64(16),
									Name: aws.String("Running"),
								},
								InstanceType:     aws.String("t2.foo"),
								SubnetId:         aws.String("foo-subnet"),
								ImageId:          aws.String("foo"),
								KeyName:          aws.String("foo-key"),
								PrivateIpAddress: aws.String("1.2.3.4"),
								PublicIpAddress:  aws.String("5.6.7.8"),
								EnaSupport:       aws.Bool(true),
								EbsOptimized:     aws.Bool(true),
							},
						},
					},
				},
			},
			nil).AnyTimes()
	}

	return mockEC2
}

func getMockEC2APIDescribeInstancesNotFoundAndPass(machineNames []string, mockCtrl *gomock.Controller) ec2iface.EC2API {
	mockEC2 := mock_ec2iface.NewMockEC2API(mockCtrl)

	dmi := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			filter.EC2.VPC(""),
			filter.EC2.ClusterOwned(""),
			filter.EC2.Name(machineNames[0]),
			filter.EC2.InstanceStates(ec2.InstanceStateNamePending, ec2.InstanceStateNameRunning),
		},
	}
	mockEC2.EXPECT().DescribeInstances(dmi).Return(
		nil,
		awserrors.NewNotFound(fmt.Errorf("mocked not found"))).Times(1)

	dmi = &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			filter.EC2.VPC(""),
			filter.EC2.ClusterOwned(""),
			filter.EC2.Name(machineNames[1]),
			filter.EC2.InstanceStates(ec2.InstanceStateNamePending, ec2.InstanceStateNameRunning),
		},
	}
	mockEC2.EXPECT().DescribeInstances(dmi).Return(
		&ec2.DescribeInstancesOutput{
			Reservations: []*ec2.Reservation{
				{
					Instances: []*ec2.Instance{
						{
							InstanceId: aws.String(machineNames[1]),
							State: &ec2.InstanceState{
								Code: aws.Int64(16),
								Name: aws.String("Running"),
							},
							InstanceType:     aws.String("t2.foo"),
							SubnetId:         aws.String("foo-subnet"),
							ImageId:          aws.String("foo"),
							KeyName:          aws.String("foo-key"),
							PrivateIpAddress: aws.String("1.2.3.4"),
							PublicIpAddress:  aws.String("5.6.7.8"),
							EnaSupport:       aws.Bool(true),
							EbsOptimized:     aws.Bool(true),
						},
					},
				},
			},
		},
		nil).Times(1)

	return mockEC2
}

func TestIsNodeJoin(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name                      string
		inputScope                *actuators.MachineScope
		inputControlplaneMachines []*clusterv1.Machine
		expectedIsNodeJoin        bool
		expectedError             error
	}{
		{
			name:                      "should always join worker machines",
			inputScope:                getWorkerMachineScope(t, nil),
			inputControlplaneMachines: nil,
			expectedIsNodeJoin:        true,
			expectedError:             nil,
		},
		{
			name:                      "should return error for random machines joining",
			inputScope:                getRandoMachineScope(t, nil),
			inputControlplaneMachines: nil,
			expectedIsNodeJoin:        false,
			expectedError:             fmt.Errorf("Unknown value %q for label `set` on machine %q", "random", "random-0"),
		},
		{
			name:                      "should return error when called with invalid scope",
			inputScope:                getInvalidMachineScope(t),
			inputControlplaneMachines: getTestControlplaneMachines(),
			expectedIsNodeJoin:        false,
			expectedError:             fmt.Errorf("failed to create machine scope for awesome-ns/master-0: failed to generate new scope from nil cluster"),
		},
		{
			name:                      "should return error for no `set` label",
			inputScope:                getNoLabelMachineScope(t, nil),
			inputControlplaneMachines: nil,
			expectedIsNodeJoin:        false,
			expectedError:             fmt.Errorf("Unknown value %q for label `set` on machine %q", "", "nolabel-0"),
		},
		{
			name:                      "should not join first controlplane machine",
			inputScope:                getControlplaneMachineScope(t, nil),
			inputControlplaneMachines: nil,
			expectedIsNodeJoin:        false,
			expectedError:             nil,
		},
		{
			name: "should not join controlplane machine when no controlplane machines exist",
			inputScope: getControlplaneMachineScope(
				t,
				getMockEC2APIDescribeInstancesNotFound([]string{"master-0", "master-1"}, mockCtrl),
			),
			inputControlplaneMachines: getTestControlplaneMachines(),
			expectedIsNodeJoin:        false,
			expectedError:             nil,
		},
		{
			name: "should return error when unable to verify controlplane machine existence",
			inputScope: getControlplaneMachineScope(
				t,
				getMockEC2APIDescribeInstancesFail([]string{"master-0", "master-1"}, mockCtrl),
			),
			inputControlplaneMachines: getTestControlplaneMachines(),
			expectedIsNodeJoin:        false,
			expectedError: fmt.Errorf(
				`failed to verify existence of machine awesome-ns/master-0: failed to lookup machine "master-0": failed to describe instances by tags: mock API failure`),
		},
		{
			name: "should join controlplane machine when other controlplane machine exists",
			inputScope: getControlplaneMachineScope(
				t,
				getMockEC2APIDescribeInstancesPass([]string{"master-0", "master-1"}, mockCtrl),
			),
			inputControlplaneMachines: getTestControlplaneMachines(),
			expectedIsNodeJoin:        true,
			expectedError:             nil,
		},
		{
			name: "should join controlplane machine when first controlplane machine doesn't exist but second controlplane machine exists",
			inputScope: getControlplaneMachineScope(
				t,
				getMockEC2APIDescribeInstancesNotFoundAndPass([]string{"master-0", "master-1"}, mockCtrl),
			),
			inputControlplaneMachines: getTestControlplaneMachines(),
			expectedIsNodeJoin:        true,
			expectedError:             nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testActuator := NewActuator(ActuatorParams{})
			actualIsNodeJoin, actualError := testActuator.isNodeJoin(tc.inputScope, tc.inputControlplaneMachines)

			if tc.expectedIsNodeJoin != actualIsNodeJoin {
				t.Fatalf("isNodeJoin failed, [%s], Got: %t, Want: %t", tc.name, actualIsNodeJoin, tc.expectedIsNodeJoin)
			}

			if tc.expectedError == nil && actualError != nil {
				t.Fatalf("isNodeJoin failed, [%s], GotError: %q, WantError: nil", tc.name, actualError)
			}

			if tc.expectedError != nil && actualError == nil {
				t.Fatalf("isNodeJoin failed, [%s], GotError: nil, WantError: %q", tc.name, tc.expectedError)
			}

			if tc.expectedError != nil && tc.expectedError.Error() != actualError.Error() {
				t.Fatalf("isNodeJoin Failed, [%s], GotError: %q, WantError: %q", tc.name, actualError, tc.expectedError)
			}

		})
	}
}
