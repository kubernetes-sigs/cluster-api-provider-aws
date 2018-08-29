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
	"testing"

	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/actuators/machine"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
	ec2svc "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/ec2"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

type ec2 struct{}

func (e *ec2) CreateInstance(machine *clusterv1.Machine) (*ec2svc.Instance, error) {
	return &ec2svc.Instance{
		ID: "abc",
	}, nil
}
func (e *ec2) InstanceIfExists(id *string) (*ec2svc.Instance, error) {
	if id == nil {
		return nil, nil
	}
	if *id == "abc" {
		return &ec2svc.Instance{
			ID: "abc",
		}, nil
	}
	return nil, nil
}

type machines struct{}

func (m *machines) UpdateMachineStatus(machine *clusterv1.Machine) (*clusterv1.Machine, error) {
	return machine, nil
}

func TestCreate(t *testing.T) {
	codec, err := v1alpha1.NewCodec()
	if err != nil {
		t.Fatalf("failed to create a codec: %v", err)
	}
	ap := machine.ActuatorParams{
		Codec:           codec,
		MachinesService: &machines{},
		EC2Service:      &ec2{},
	}
	actuator, err := machine.NewActuator(ap)
	if err != nil {
		t.Fatalf("failed to create an actuator: %v", err)
	}

	if err := actuator.Create(&clusterv1.Cluster{}, &clusterv1.Machine{}); err != nil {
		t.Fatalf("failed to create machine: %v", err)
	}
}
