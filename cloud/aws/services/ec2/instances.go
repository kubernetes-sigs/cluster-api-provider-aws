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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
)

// InstanceIfExists returns the existing instance or nothing if it doesn't exist.
func (s *Service) InstanceIfExists(instanceID *string) (*v1alpha1.Instance, error) {
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
		ec2instance := out.Reservations[0].Instances[0]
		return fromSDKTypeToInstance(ec2instance), nil
	}

	return nil, nil
}

// CreateInstance runs an ec2 instance.
func (s *Service) CreateInstance(machine *clusterv1.Machine, config *v1alpha1.AWSMachineProviderConfig, clusterStatus *v1alpha1.AWSClusterProviderStatus) (*v1alpha1.Instance, error) {
	id := config.AMI.ID
	if id == nil {
		id = aws.String(s.defaultAMILookup(clusterStatus.Region))
	}

	subnets := clusterStatus.Network.Subnets.FilterPublic()
	if len(subnets) == 0 {
		return nil, errors.New("need at least one subnet but didn't find any")
	}

	input := &ec2.RunInstancesInput{
		ImageId:      id,
		InstanceType: aws.String(config.InstanceType),
		MaxCount:     aws.Int64(1),
		MinCount:     aws.Int64(1),
		SubnetId:     aws.String(subnets[0].ID),
	}

	if config.KeyName != "" {
		input.KeyName = aws.String(config.KeyName)
	}

	reservation, err := s.EC2.RunInstances(input)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to run instances")
	}

	if len(reservation.Instances) <= 0 {
		return nil, errors.New("no instance was created after run was called")
	}

	return fromSDKTypeToInstance(reservation.Instances[0]), nil
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

// CreateOrGetMachine will either return an existing instance or create and return an instance.
func (s *Service) CreateOrGetMachine(machine *clusterv1.Machine, status *v1alpha1.AWSMachineProviderStatus, config *v1alpha1.AWSMachineProviderConfig, clusterStatus *v1alpha1.AWSClusterProviderStatus) (*v1alpha1.Instance, error) {
	// instance id exists, try to get it
	if status.InstanceID != nil {
		instance, err := s.InstanceIfExists(status.InstanceID)

		// if there was no error, return the found instance
		if err == nil {
			return instance, err
		}

		// if there was an error but it's not IsNotFound then it's a real error
		if !IsNotFound(err) {
			return instance, err
		}
	}

	// otherwise let's create it
	return s.CreateInstance(machine, config, clusterStatus)
}

func fromSDKTypeToInstance(v *ec2.Instance) *v1alpha1.Instance {
	i := &v1alpha1.Instance{
		ID:           *v.InstanceId,
		State:        v1alpha1.InstanceState(*v.State.Name),
		Type:         *v.InstanceType,
		SubnetID:     *v.SubnetId,
		ImageID:      *v.ImageId,
		KeyName:      v.KeyName,
		PrivateIP:    v.PrivateIpAddress,
		PublicIP:     v.PublicIpAddress,
		ENASupport:   v.EnaSupport,
		EBSOptimized: v.EbsOptimized,
	}

	if v.IamInstanceProfile != nil && v.IamInstanceProfile.Id != nil {
		i.IamProfileID = v.IamInstanceProfile.Id
	}

	return i
}
