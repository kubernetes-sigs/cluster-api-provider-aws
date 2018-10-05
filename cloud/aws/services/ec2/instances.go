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

	"github.com/golang/glog"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/certificates"
)

// InstanceIfExists returns the existing instance or nothing if it doesn't exist.
func (s *Service) InstanceIfExists(instanceID *string) (*v1alpha1.Instance, error) {
	glog.V(2).Infof("Looking for instance %q", *instanceID)
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{instanceID},
	}

	out, err := s.EC2.DescribeInstances(input)
	switch {
	case IsNotFound(err):
		return nil, nil
	case err != nil:
		return nil, errors.Errorf("failed to describe instances: %v", err)
	}

	if len(out.Reservations) > 0 && len(out.Reservations[0].Instances) > 0 {
		return fromSDKTypeToInstance(out.Reservations[0].Instances[0]), nil
	}

	return nil, nil
}

// CreateInstance runs an ec2 instance.
func (s *Service) CreateInstance(machine *clusterv1.Machine, config *v1alpha1.AWSMachineProviderConfig, clusterStatus *v1alpha1.AWSClusterProviderStatus) (*v1alpha1.Instance, error) {

	input := &v1alpha1.Instance{
		Type:       config.InstanceType,
		IAMProfile: config.IAMInstanceProfile,
	}

	// Pick image from the machine configuration, or use a default one.
	if config.AMI.ID != nil {
		input.ImageID = *config.AMI.ID
	} else {
		input.ImageID = s.defaultAMILookup(clusterStatus.Region)
	}

	// Pick subnet from the machine configuration, or default to the first private available.
	if config.Subnet != nil && config.Subnet.ID != nil {
		input.SubnetID = *config.Subnet.ID
	} else {
		sns := clusterStatus.Network.Subnets.FilterPrivate()
		if len(sns) == 0 {
			return nil, errors.New("failed to run instance, no subnets available")
		}
		input.SubnetID = sns[0].ID
	}

	// apply values based on the role of the machine
	if machine.ObjectMeta.Labels["set"] == "controlplane" {
		caCert, caKey, err := certificates.NewCertificateAuthority()
		if err != nil {
			return input, errors.Wrap(err, "Failed to generate a CA for the control plane")
		}
		clusterStatus.CACertificate = certificates.EncodeCertPEM(caCert)
		clusterStatus.CAPrivateKey = certificates.EncodePrivateKeyPEM(caKey)

		input.UserData = aws.String(initControlPlaneScript(clusterStatus.CACertificate, clusterStatus.CAPrivateKey))
		input.SecurityGroupIDs = append(input.SecurityGroupIDs, clusterStatus.Network.SecurityGroups[v1alpha1.SecurityGroupControlPlane].ID)
	}

	if machine.ObjectMeta.Labels["set"] == "node" {
		input.SecurityGroupIDs = append(input.SecurityGroupIDs, clusterStatus.Network.SecurityGroups[v1alpha1.SecurityGroupNode].ID)
	}

	// Pick SSH key, if any.
	if config.KeyName != "" {
		input.KeyName = aws.String(config.KeyName)
	}

	return s.runInstance(input)
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
		glog.V(2).Infof("Looking up instance %q", *status.InstanceID)
		instance, err := s.InstanceIfExists(status.InstanceID)

		// if there was no error, return the found instance
		if err == nil {
			return instance, nil
		}

		// if there was an error but it's not IsNotFound then it's a real error
		if !IsNotFound(err) {
			return instance, errors.Wrapf(err, "instance %q was not found", *status.InstanceID)
		}

		return instance, errors.Wrapf(err, "failed to look up instance %q", *status.InstanceID)
	}

	// otherwise let's create it
	return s.CreateInstance(machine, config, clusterStatus)
}

func (s *Service) runInstance(i *v1alpha1.Instance) (*v1alpha1.Instance, error) {
	input := &ec2.RunInstancesInput{
		InstanceType: aws.String(i.Type),
		SubnetId:     aws.String(i.SubnetID),
		ImageId:      aws.String(i.ImageID),
		KeyName:      i.KeyName,
		EbsOptimized: i.EBSOptimized,
		MaxCount:     aws.Int64(1),
		MinCount:     aws.Int64(1),
		UserData:     i.UserData,
	}

	if i.UserData != nil {
		input.UserData = aws.String(base64.StdEncoding.EncodeToString([]byte(*i.UserData)))
	}

	if len(i.SecurityGroupIDs) > 0 {
		input.SecurityGroupIds = aws.StringSlice(i.SecurityGroupIDs)
	}

	if i.IAMProfile != "" {
		input.IamInstanceProfile = &ec2.IamInstanceProfileSpecification{
			Name: aws.String(i.IAMProfile),
		}
	}

	if len(i.Tags) > 0 {
		spec := &ec2.TagSpecification{ResourceType: aws.String(ec2.ResourceTypeInstance)}
		for key, value := range i.Tags {
			spec.Tags = append(spec.Tags, &ec2.Tag{
				Key:   aws.String(key),
				Value: aws.String(value),
			})
		}

		input.TagSpecifications = append(input.TagSpecifications, spec)
	}

	out, err := s.EC2.RunInstances(input)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to run instance: %v", i)
	}

	if len(out.Instances) == 0 {
		return nil, errors.Errorf("no instance returned for reservation %v", out.GoString())
	}

	return fromSDKTypeToInstance(out.Instances[0]), nil
}

// UpdateInstanceSecurityGroups modifies the security groups of the given
// EC2 instance.
func (s *Service) UpdateInstanceSecurityGroups(instanceID *string, securityGroups []*string) error {
	input := &ec2.ModifyInstanceAttributeInput{
		InstanceId: instanceID,
		Groups:     securityGroups,
	}

	_, err := s.EC2.ModifyInstanceAttribute(input)
	if err != nil {
		return err
	}

	return nil
}

// UpdateResourceTags updates the tags for an instance.
// This will be called if there is anything to create (update) or delete.
// We may not always have to perform each action, so we check what we're
// receiving to avoid calling AWS if we don't need to.
func (s *Service) UpdateResourceTags(resourceID *string, create map[string]string, remove map[string]string) error {
	// If we have anything to create or update
	if len(create) > 0 {
		// Convert our create map into an array of *ec2.Tag
		createTagsInput := mapToTags(create)

		// Create the CreateTags input.
		input := &ec2.CreateTagsInput{
			Resources: []*string{resourceID},
			Tags:      createTagsInput,
		}

		// Create/Update tags in AWS.
		_, err := s.EC2.CreateTags(input)
		if err != nil {
			return err
		}
	}

	// If we have anything to remove
	if len(remove) > 0 {
		// Convert our remove map into an array of *ec2.Tag
		removeTagsInput := mapToTags(remove)

		// Create the DeleteTags input
		input := &ec2.DeleteTagsInput{
			Resources: []*string{resourceID},
			Tags:      removeTagsInput,
		}

		// Delete tags in AWS.
		_, err := s.EC2.DeleteTags(input)
		if err != nil {
			return err
		}
	}

	return nil
}

// fromSDKTypeToInstance takes a ec2.Instance and returns our v1.alpha1.Instance
// type. EC2 types are wrapped or converted to our own types here.
func fromSDKTypeToInstance(v *ec2.Instance) *v1alpha1.Instance {
	i := &v1alpha1.Instance{
		ID:           aws.StringValue(v.InstanceId),
		State:        v1alpha1.InstanceState(*v.State.Name),
		Type:         aws.StringValue(v.InstanceType),
		SubnetID:     aws.StringValue(v.SubnetId),
		ImageID:      aws.StringValue(v.ImageId),
		KeyName:      v.KeyName,
		PrivateIP:    v.PrivateIpAddress,
		PublicIP:     v.PublicIpAddress,
		ENASupport:   v.EnaSupport,
		EBSOptimized: v.EbsOptimized,
	}

	for _, sg := range v.SecurityGroups {
		i.SecurityGroupIDs = append(i.SecurityGroupIDs, *sg.GroupId)
	}

	// TODO: Handle returned IAM instance profile, since we are currently
	// using a string representing the name, but the InstanceProfile returned
	// from the sdk only returns ARN and ID.

	if len(v.Tags) > 0 {
		i.Tags = tagsToMap(v.Tags)
	}

	if len(v.SecurityGroups) > 0 {
		i.SecurityGroups = groupIdentifierToMap(v.SecurityGroups)
	}

	return i
}

// initControlPlaneScript returns the b64 encoded script to run on start up.
// The cert Must be CertPEM encoded and the key must be PrivateKeyPEM encoded
func initControlPlaneScript(caCert, caKey []byte) string {
	// The script must start with #!. If it goes on the next line Dedent will start the script with a \n.
	return fmt.Sprintf(`#!/usr/bin/env bash

mkdir -p /etc/kubernetes/pki

echo '%s' > /etc/kubernetes/pki/ca.crt
echo '%s' > /etc/kubernetes/pki/ca.key

cat >/tmp/kubeadm.yaml <<EOF
apiVersion: kubeadm.k8s.io/v1alpha3
kind: InitConfiguration
nodeRegistration:
  criSocket: /var/run/containerd/containerd.sock
EOF

kubeadm init --config /tmp/kubeadm.yaml

# Installation from https://docs.projectcalico.org/v3.2/getting-started/kubernetes/installation/calico
kubectl --kubeconfig /etc/kubernetes/admin.conf apply -f https://docs.projectcalico.org/v3.2/getting-started/kubernetes/installation/hosted/rbac-kdd.yaml
kubectl --kubeconfig /etc/kubernetes/admin.conf apply -f https://docs.projectcalico.org/v3.2/getting-started/kubernetes/installation/hosted/kubernetes-datastore/calico-networking/1.7/calico.yaml
`, caCert, caKey)
}
