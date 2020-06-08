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

package ec2

import (
	"encoding/base64"
	"fmt"
	"sigs.k8s.io/cluster-api/util/conditions"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/filter"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/userdata"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/record"
)

const (
	defaultSSHKeyName = "default"
)

// ReconcileBastion ensures a bastion is created for the cluster
func (s *Service) ReconcileBastion() error {
	if s.scope.VPC().IsUnmanaged(s.scope.Name()) {
		s.scope.V(4).Info("Skipping bastion reconcile in unmanaged mode")
		return nil
	}

	if !s.scope.AWSCluster.Spec.Bastion.Enabled {
		_, err := s.describeBastionInstance()
		if err != nil {
			if awserrors.IsNotFound(err) {
				return nil
			}
			return err
		}
		return s.DeleteBastion()
	}

	s.scope.V(2).Info("Reconciling bastion host")

	subnets := s.scope.Subnets()
	if len(subnets.FilterPrivate()) == 0 {
		s.scope.V(2).Info("No private subnets available, skipping bastion host")
		return nil
	} else if len(subnets.FilterPublic()) == 0 {
		return errors.New("failed to reconcile bastion host, no public subnets are available")
	}

	spec := s.getDefaultBastion()

	// Describe bastion instance, if any.
	instance, err := s.describeBastionInstance()
	if awserrors.IsNotFound(err) {
		instance, err = s.runInstance("bastion", spec)
		if err != nil {
			record.Warnf(s.scope.AWSCluster, "FailedCreateBastion", "Failed to create bastion instance: %v", err)
			return err
		}

		record.Eventf(s.scope.AWSCluster, "SuccessfulCreateBastion", "Created bastion instance %q", instance.ID)
		s.scope.V(2).Info("Created new bastion host", "instance", instance)

	} else if err != nil {
		return err
	}

	// TODO(vincepri): check for possible changes between the default spec and the instance.

	s.scope.AWSCluster.Status.Bastion = instance.DeepCopy()
	conditions.MarkTrue(s.scope.AWSCluster, infrav1.BastionHostReadyCondition)
	s.scope.V(2).Info("Reconcile bastion completed successfully")

	return nil
}

// DeleteBastion deletes the Bastion instance
func (s *Service) DeleteBastion() error {
	if s.scope.VPC().IsUnmanaged(s.scope.Name()) {
		s.scope.V(4).Info("Skipping bastion deletion in unmanaged mode")
		return nil
	}

	instance, err := s.describeBastionInstance()
	if err != nil {
		if awserrors.IsNotFound(err) {
			s.scope.V(4).Info("bastion instance does not exist")
			return nil
		}
		return errors.Wrap(err, "unable to describe bastion instance")
	}

	if err := s.TerminateInstanceAndWait(instance.ID); err != nil {
		record.Warnf(s.scope.AWSCluster, "FailedTerminateBastion", "Failed to terminate bastion instance %q: %v", instance.ID, err)
		return errors.Wrap(err, "unable to delete bastion instance")
	}
	record.Eventf(s.scope.AWSCluster, "SuccessfulTerminateBastion", "Terminated bastion instance %q", instance.ID)

	return nil
}

func (s *Service) describeBastionInstance() (*infrav1.Instance, error) {
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			filter.EC2.ProviderRole(infrav1.BastionRoleTagValue),
			filter.EC2.Cluster(s.scope.Name()),
			filter.EC2.InstanceStates(
				ec2.InstanceStateNamePending,
				ec2.InstanceStateNameRunning,
				ec2.InstanceStateNameStopping,
				ec2.InstanceStateNameStopped,
			),
		},
	}

	out, err := s.scope.EC2.DescribeInstances(input)
	if err != nil {
		record.Eventf(s.scope.AWSCluster, "FailedDescribeBastionHost", "Failed to describe bastion host: %v", err)
		return nil, errors.Wrap(err, "failed to describe bastion host")
	}

	// TODO: properly handle multiple bastions found rather than just returning
	// the first non-terminated.
	for _, res := range out.Reservations {
		for _, instance := range res.Instances {
			if aws.StringValue(instance.State.Name) != ec2.InstanceStateNameTerminated {
				return s.SDKToInstance(instance)
			}
		}
	}

	return nil, awserrors.NewNotFound(errors.New("bastion host not found"))
}

func (s *Service) getDefaultBastion() *infrav1.Instance {
	name := fmt.Sprintf("%s-bastion", s.scope.Name())
	userData, _ := userdata.NewBastion(&userdata.BastionInput{})

	// If SSHKeyName WAS NOT provided, use the defaultSSHKeyName
	keyName := s.scope.AWSCluster.Spec.SSHKeyName
	if keyName == nil {
		keyName = aws.String(defaultSSHKeyName)
	}

	i := &infrav1.Instance{
		Type:       "t2.micro",
		SubnetID:   s.scope.Subnets().FilterPublic()[0].ID,
		ImageID:    s.defaultBastionAMILookup(s.scope.AWSCluster.Spec.Region),
		SSHKeyName: keyName,
		UserData:   aws.String(base64.StdEncoding.EncodeToString([]byte(userData))),
		SecurityGroupIDs: []string{
			s.scope.Network().SecurityGroups[infrav1.SecurityGroupBastion].ID,
		},
		Tags: infrav1.Build(infrav1.BuildParams{
			ClusterName: s.scope.Name(),
			Lifecycle:   infrav1.ResourceLifecycleOwned,
			Name:        aws.String(name),
			Role:        aws.String(infrav1.BastionRoleTagValue),
			Additional:  s.scope.AdditionalTags(),
		}),
	}

	return i
}
