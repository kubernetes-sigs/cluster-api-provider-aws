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
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	providerconfigv1 "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
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

	// TODO: Get default AMI using a lookup/filter based on tags added with image baking utils
	defaultAMIID = "ami-0de61b6929e8f091c"

	defaultInstanceType = ec2.InstanceTypeT3Medium

	defaultKeyName = "default"
)

// Instance is an internal representation of an AWS instance.
// This contains more data than the provider config struct tracked in the status.
type Instance struct {
	// State can be things like "running", "terminated", "stopped", etc.
	State string
	// ID is the AWS InstanceID.
	ID string
	// SubnetID is the AWS SubnetID that the interface exists in
	SubnetID string
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
			State:    aws.StringValue(out.Reservations[0].Instances[0].State.Name),
			ID:       aws.StringValue(out.Reservations[0].Instances[0].InstanceId),
			SubnetID: aws.StringValue(out.Reservations[0].Instances[0].SubnetId),
		}, nil
	}

	return nil, nil
}

// CreateInstance runs an ec2 instance.
func (s *Service) CreateInstance(machine *clusterv1.Machine, providerConfig *providerconfigv1.AWSMachineProviderConfig, cluster *clusterv1.Cluster, clusterProviderConfig *providerconfigv1.AWSClusterProviderConfig, clusterStatus *providerconfigv1.AWSClusterProviderStatus) (*Instance, error) {
	// TODO: attempt to infer image id from cluster
	id := providerConfig.AMI.ID
	if id == nil {
		id = aws.String(defaultAMIID)
	}

	// TODO: attempt to infer instance type from cluster
	instanceType := providerConfig.InstanceType
	if instanceType == "" {
		instanceType = defaultInstanceType
	}

	// TODO: attempt to infer KeyName from cluster
	keyName := providerConfig.KeyName
	if keyName == "" {
		keyName = defaultKeyName
	}

	var subnetID *string
	if providerConfig.Subnet == nil || providerConfig.Subnet.ID == nil {
		subnets := clusterStatus.Network.Subnets.FilterPrivate()
		if len(subnets) == 0 {
			return nil, errors.New("need at least one subnet but didn't find any")
		}
		subnetID = aws.String(subnets[0].ID)
	} else {
		subnetID = providerConfig.Subnet.ID
	}

	// TODO: handle common cluster config for security groups
	// TODO: better error handling
	var sgIds []*string
	for _, sg := range providerConfig.AdditionalSecurityGroups {
		sgIds = append(sgIds, sg.ID)
	}

	// TODO: handle common cluster config for tags
	additionalTags := providerConfig.AdditionalTags
	if additionalTags == nil {
		additionalTags = make(map[string]string)
	}
	if _, ok := additionalTags["Name"]; !ok {
		additionalTags["Name"] = machine.Name
	}
	tags := tagMapToEC2Tags(s.buildTags(cluster.Name, ResourceLifecycleOwned, additionalTags))

	var userDataText string

	// TODO: ensure cluster config and additional machine config are plumbed through
	if providerConfig.NodeRole == "controlplane" {
		// TODO: remove hard coded bootstrap token
		userDataText = `#cloud-config
kubeadm:
  operation: init
  config: /run/kubeadm/kubeadm.config
write_files:
- path: /run/kubeadm/kubeadm.config
  content: |
    apiVersion: kubeadm.k8s.io/v1alpha3
    kind: InitConfiguration
`
		userDataText += fmt.Sprintf("    apiServerCertSANs: [\"%s\"]\n", clusterStatus.Network.APIServerLoadBalancer.DNSName)
		userDataText += fmt.Sprintf("    controlPlaneEndpoint: \"%s:443\"\n", clusterStatus.Network.APIServerLoadBalancer.DNSName)
		userDataText += fmt.Sprintf("    clusterName: %s\n", cluster.Name)
		userDataText += "    networking:\n"
		userDataText += fmt.Sprintf("      dnsDomain: %s\n", cluster.Spec.ClusterNetwork.ServiceDomain)
		userDataText += fmt.Sprintf("      podSubnet: %s\n", cluster.Spec.ClusterNetwork.Pods.CIDRBlocks[0])
		userDataText += fmt.Sprintf("      serviceSubnet: %s\n", cluster.Spec.ClusterNetwork.Services.CIDRBlocks[0])

		// Set kubernetes version
		userDataText += fmt.Sprintf("    kubernetesVersion: %s\n", machine.Spec.Versions.ControlPlane)
		userDataText += `
    nodeRegistration:
      criSocket: /run/containerd/containerd.sock
    bootstrapTokens:
    - groups:
      - system:bootstrappers:kubeadm:default-node-token
      token: abcdef.0123456789abcdef
      ttl: 8760h0m0s
      usages:
      - signing
      - authentication
`
	} else {
		// TODO: remove hard coded bootstrap token
		userDataText = `#cloud-config
kubeadm:
  operation: join
  config: /run/kubeadm/kubeadm.config
write_files:
- path: /run/kubeadm/kubeadm.config
  content: |
    apiVersion: kubeadm.k8s.io/v1alpha3
    kind: JoinConfiguration
    nodeRegistration:
      criSocket: /run/containerd/containerd.sock
	token: abcdef.0123456789abcdef
    discoveryTokenUnsafeSkipCAVerification: true
`
		userDataText += "    discoveryTokenAPIServers:\n"
		userDataText += fmt.Sprintf("    - %s:443\n", clusterStatus.Network.APIServerLoadBalancer.DNSName)
		userDataText += fmt.Sprintf("    clusterName: %s\n", cluster.Name)
	}

	userData := base64.StdEncoding.EncodeToString([]byte(userDataText))

	var iamInstanceProfileSpec *ec2.IamInstanceProfileSpecification
	if providerConfig.IAMInstanceProfile != nil {
		iamInstanceProfileSpec = &ec2.IamInstanceProfileSpecification{
			Arn:  providerConfig.IAMInstanceProfile.ARN,
			Name: providerConfig.IAMInstanceProfile.Name,
		}
	}

	tagSpecifications := []*ec2.TagSpecification{}
	if len(tags) > 0 {
		tagSpecifications = append(tagSpecifications, &ec2.TagSpecification{
			ResourceType: aws.String("instance"),
			Tags:         tags,
		})
		tagSpecifications = append(tagSpecifications, &ec2.TagSpecification{
			ResourceType: aws.String("volume"),
			Tags:         tags,
		})
	}

	// TODO: handle providerConfig.PublicIP
	input := &ec2.RunInstancesInput{
		// TODO: using machine.UID here will likely cause issues for update workflows that involve
		// deleting and re-creating the instance.
		//ClientToken:        aws.String(string(machine.UID)),
		IamInstanceProfile: iamInstanceProfileSpec,
		ImageId:            id,
		InstanceType:       aws.String(instanceType),
		KeyName:            aws.String(keyName),
		MaxCount:           aws.Int64(1),
		MinCount:           aws.Int64(1),
		SecurityGroupIds:   sgIds,
		SubnetId:           subnetID,
		TagSpecifications:  tagSpecifications,
		UserData:           aws.String(userData),
	}

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
