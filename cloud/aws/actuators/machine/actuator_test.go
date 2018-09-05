// Copyright © 2018 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package machine_test

import (
	"errors"
	"testing"

	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/actuators/machine"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
	ec2svc "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/ec2"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

type instanceError struct {
	instance *ec2svc.Instance
	err      error
}

type ec2 struct {
	createInstanceReturn   instanceError
	instanceIfExistsReturn instanceError
}

func (e *ec2) CreateInstance(machine *clusterv1.Machine) (*ec2svc.Instance, error) {
	return e.createInstanceReturn.instance, e.createInstanceReturn.err
}
func (e *ec2) InstanceIfExists(id *string) (*ec2svc.Instance, error) {
	return e.instanceIfExistsReturn.instance, e.instanceIfExistsReturn.err
}

type machines struct{}

func (m *machines) UpdateMachineStatus(machine *clusterv1.Machine) (*clusterv1.Machine, error) {
	return machine, nil
}

func TestCreate(t *testing.T) {
	testCases := []struct {
		name                   string
		createInstanceReturn   instanceError
		instanceIfExistsReturn instanceError

		expected error
	}{
		{
			name: "instance does not exist",
			createInstanceReturn: instanceError{
				instance: &ec2svc.Instance{
					State: "Running",
					ID:    "abc",
				},
			},
			instanceIfExistsReturn: instanceError{
				instance: nil,
				err:      nil,
			},
			expected: nil,
		},
		{
			name:                 "instance lookup fails",
			createInstanceReturn: instanceError{},
			instanceIfExistsReturn: instanceError{
				instance: nil,
				err:      errors.New("failed to create instance"),
			},
			expected: errors.New("failed to create instance"),
		},
	}

	// shared codec
	codec, err := v1alpha1.NewCodec()
	if err != nil {
		t.Fatalf("failed to create a codec: %v", err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ap := machine.ActuatorParams{
				Codec:           codec,
				MachinesService: &machines{},
				EC2Service: &ec2{
					createInstanceReturn:   tc.createInstanceReturn,
					instanceIfExistsReturn: tc.instanceIfExistsReturn,
				},
			}
			actuator, err := machine.NewActuator(ap)
			if err != nil {
				t.Fatalf("failed to create an actuator: %v", err)
			}

			if err := actuator.Create(&clusterv1.Cluster{}, &clusterv1.Machine{}); err != nil && err.Error() != tc.expected.Error() {
				t.Fatalf("expected error: %q but got error %q", err, tc.expected)
			}
		})
	}
}
