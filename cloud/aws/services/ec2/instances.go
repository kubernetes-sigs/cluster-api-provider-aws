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
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

const (
	// InstanceStateShuttingDown indicates the instance is shutting-down
	InstanceStateShuttingDown = ec2.InstanceStateNameShuttingDown

	// InstanceStateTerminated indicates the instance has been terminated
	InstanceStateTerminated = ec2.InstanceStateNameTerminated

	// InstanceStateRunning indicates the instance is running
	InstanceStateRunning = ec2.InstanceStateNameRunning

	// InstanceStatePending indicates the instance is pending
	InstanceStatePending = ec2.InstanceStateNamePending
)

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
	out, err := s.EC2.DescribeInstances(input)

	switch {
	case IsNotFound(err):
		return nil, nil
	case err != nil:
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

	reservation, err := s.EC2.RunInstances(input)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to run instances")
	}

	if len(reservation.Instances) <= 0 {
		return nil, errors.New("no instance was created after run was called")
	}

	return &Instance{
		State: *reservation.Instances[0].State.Name,
		ID:    *reservation.Instances[0].InstanceId,
	}, nil
}

// TerminateInstance terminates an EC2 instance.
// Returns nil on success, error in all other cases.
func (s *Service) TerminateInstance(instanceID *string) error {
	input := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{
			instanceID,
		},
	}

	_, err := s.EC2.TerminateInstances(input)
	if err != nil {
		return err
	}

	return nil
}
