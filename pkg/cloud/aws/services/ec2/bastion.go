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
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"k8s.io/klog"
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
func (s *Service) ReconcileBastion(clusterName, keyName string, status *v1alpha1.AWSClusterProviderStatus) error {
	klog.V(2).Info("Reconciling bastion host")

	subnets := status.Network.Subnets
	if len(subnets.FilterPrivate()) == 0 {
		klog.V(2).Info("No private subnets available, skipping bastion host")
		return nil
	} else if len(subnets.FilterPublic()) == 0 {
		return errors.New("failed to reconcile bastion host, no public subnets are available")
	}

	if keyName == "" {
		keyName = defaultSSHKeyName
	}

	spec := s.getDefaultBastion(clusterName, status.Region, status.Network, keyName)

	// Describe bastion instance, if any.
	instance, err := s.describeBastionInstance(clusterName, status)
	if awserrors.IsNotFound(err) {
		instance, err = s.runInstance(spec)
		if err != nil {
			return err
		}

		klog.V(2).Infof("Created new bastion host: %+v", instance)

	} else if err != nil {
		return err
	}

	// TODO(vincepri): check for possible changes between the default spec and the instance.

	instance.DeepCopyInto(&status.Bastion)

	klog.V(2).Info("Reconcile bastion completed successfully")
	return nil
}

// DeleteBastion deletes the Bastion instance
func (s *Service) DeleteBastion(clusterName string, status *v1alpha1.AWSClusterProviderStatus) error {
	instance, err := s.describeBastionInstance(clusterName, status)
	if err != nil {
		if awserrors.IsNotFound(err) {
			klog.V(2).Info("bastion instance does not exist")
			return nil
		}
		return errors.Wrap(err, "unable to describe bastion instance")
	}

	err = s.TerminateInstanceAndWait(instance.ID)
	if err != nil {
		return errors.Wrap(err, "unable to delete bastion instance")
	}

	return nil
}

func (s *Service) describeBastionInstance(clusterName string, status *v1alpha1.AWSClusterProviderStatus) (*v1alpha1.Instance, error) {
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			filter.EC2.ProviderRole(tags.ValueBastionRole),
			filter.EC2.Cluster(clusterName),
			filter.EC2.InstanceStates(ec2.InstanceStateNamePending, ec2.InstanceStateNameRunning),
		},
	}

	out, err := s.EC2.DescribeInstances(input)
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

func (s *Service) getDefaultBastion(clusterName string, region string, network v1alpha1.Network, keyName string) *v1alpha1.Instance {
	name := fmt.Sprintf("%s-bastion", clusterName)
	userData, _ := userdata.NewBastion(&userdata.BastionInput{})

	i := &v1alpha1.Instance{
		Type:     "t2.micro",
		SubnetID: network.Subnets.FilterPublic()[0].ID,
		ImageID:  s.defaultBastionAMILookup(region),
		KeyName:  aws.String(keyName),
		UserData: aws.String(base64.StdEncoding.EncodeToString([]byte(userData))),
		SecurityGroupIDs: []string{
			network.SecurityGroups[v1alpha1.SecurityGroupBastion].ID,
		},
		Tags: tags.Build(tags.BuildParams{
			ClusterName: clusterName,
			Lifecycle:   tags.ResourceLifecycleOwned,
			Name:        aws.String(name),
			Role:        aws.String(tags.ValueBastionRole),
		}),
	}

	return i
}
