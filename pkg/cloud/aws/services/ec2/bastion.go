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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/converters"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/filter"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/userdata"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/tags"
)

const (
	defaultSSHKeyName = "default"
)

// ReconcileBastion ensures a bastion is created for the cluster
func (s *Service) ReconcileBastion() error {
	if s.scope.VPC().IsProvided() {
		s.scope.V(4).Info("Skipping bastion reconcile in unmanaged mode")
		return nil
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
			return err
		}

		s.scope.V(2).Info("Created new bastion host", "instance", instance)

	} else if err != nil {
		return err
	}

	// TODO(vincepri): check for possible changes between the default spec and the instance.

	instance.DeepCopyInto(&s.scope.ClusterStatus.Bastion)
	s.scope.V(2).Info("Reconcile bastion completed successfully")
	return nil
}

// DeleteBastion deletes the Bastion instance
func (s *Service) DeleteBastion() error {
	if s.scope.VPC().IsProvided() {
		s.scope.V(4).Info("Skipping bastion deletion in unmanaged mode")
		return nil
	}

	instance, err := s.describeBastionInstance()
	if err != nil {
		if awserrors.IsNotFound(err) {
			s.scope.V(2).Info("bastion instance does not exist", "instance-id", instance.ID)
			return nil
		}
		return errors.Wrap(err, "unable to describe bastion instance")
	}

	if err := s.TerminateInstanceAndWait(instance.ID); err != nil {
		return errors.Wrap(err, "unable to delete bastion instance")
	}

	return nil
}

func (s *Service) describeBastionInstance() (*v1alpha1.Instance, error) {
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			filter.EC2.ProviderRole(tags.ValueBastionRole),
			filter.EC2.Cluster(s.scope.Name()),
			filter.EC2.InstanceStates(ec2.InstanceStateNamePending, ec2.InstanceStateNameRunning),
		},
	}

	out, err := s.scope.EC2.DescribeInstances(input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to describe bastion host")
	}

	// TODO: properly handle multiple bastions found rather than just returning
	// the first non-terminated.
	for _, res := range out.Reservations {
		for _, instance := range res.Instances {
			if aws.StringValue(instance.State.Name) != ec2.InstanceStateNameTerminated {
				return converters.SDKToInstance(instance), nil
			}
		}
	}

	return nil, awserrors.NewNotFound(errors.New("bastion host not found"))
}

func (s *Service) getDefaultBastion() *v1alpha1.Instance {
	name := fmt.Sprintf("%s-bastion", s.scope.Name())
	userData, _ := userdata.NewBastion(&userdata.BastionInput{})

	keyName := defaultSSHKeyName
	if s.scope.ClusterConfig.SSHKeyName != "" {
		keyName = s.scope.ClusterConfig.SSHKeyName
	}

	i := &v1alpha1.Instance{
		Type:     "t2.micro",
		SubnetID: s.scope.Subnets().FilterPublic()[0].ID,
		ImageID:  s.defaultBastionAMILookup(s.scope.ClusterConfig.Region),
		KeyName:  aws.String(keyName),
		UserData: aws.String(base64.StdEncoding.EncodeToString([]byte(userData))),
		SecurityGroupIDs: []string{
			s.scope.Network().SecurityGroups[v1alpha1.SecurityGroupBastion].ID,
		},
		Tags: tags.Build(tags.BuildParams{
			ClusterName: s.scope.Name(),
			Lifecycle:   tags.ResourceLifecycleOwned,
			Name:        aws.String(name),
			Role:        aws.String(tags.ValueBastionRole),
		}),
	}

	return i
}
