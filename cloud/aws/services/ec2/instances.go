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

package ec2

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

// instances is used to scope down/organize the ec2 client.
type instances interface {
	DescribeInstances(*ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error)
	RunInstances(*ec2.RunInstancesInput) (*ec2.Reservation, error)
}

// Instance is an internal representation of an AWS instance.
// This contains more data than the provider config struct tracked in the status.
type Instance struct {
	// State can be things like "running", "terminated", "stopped", etc.
	State string
	// ID is the AWS InstanceID.
	ID string
}

// InstanceIfExists returns the existing instance or nothing if it doesn't exist.
func (s *Service) InstanceIfExists(instanceID *string) (*Instance, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{instanceID},
	}
	out, err := s.Instances.DescribeInstances(input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe instances: %v", err)
	}
	if len(out.Reservations) > 0 && len(out.Reservations[0].Instances) > 0 {
		return &Instance{
			State: *out.Reservations[0].Instances[0].State.Name,
			ID:    *out.Reservations[0].Instances[0].InstanceId,
		}, nil
	}
	return nil, nil
}

// CreateInstance runs an ec2 instance.
func (s *Service) CreateInstance(machine *clusterv1.Machine) (*Instance, error) {
	input := &ec2.RunInstancesInput{}
	reservation, err := s.Instances.RunInstances(input)
	if err != nil {
		return nil, fmt.Errorf("failed to run instances: %v", err)
	}
	if len(reservation.Instances) <= 0 {
		return nil, errors.New("no instance was created after run was called")
	}
	return &Instance{
		State: *reservation.Instances[0].State.Name,
		ID:    *reservation.Instances[0].InstanceId,
	}, nil
}
