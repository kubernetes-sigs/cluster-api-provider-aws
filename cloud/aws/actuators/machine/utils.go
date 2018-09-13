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

package machine

import (
	"bytes"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	jsonserializer "k8s.io/apimachinery/pkg/runtime/serializer/json"

	awsclient "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/client"
	providerconfigv1 "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// SortInstances will examine the given slice of instances and return the current active instance for
// the machine, as well as a slice of all other instances which the caller may want to terminate. The
// active instance is calculated as the most recently launched instance.
// This function should only be called with running instances, not those which are stopped or
// terminated.
func SortInstances(instances []*ec2.Instance) (*ec2.Instance, []*ec2.Instance) {
	if len(instances) == 0 {
		return nil, []*ec2.Instance{}
	}
	var newestInstance *ec2.Instance
	inactiveInstances := make([]*ec2.Instance, 0, len(instances)-1)
	for _, i := range instances {
		if newestInstance == nil {
			newestInstance = i
			continue
		}
		tempInstance := chooseNewest(newestInstance, i)
		if *tempInstance.InstanceId != *newestInstance.InstanceId {
			inactiveInstances = append(inactiveInstances, newestInstance)
		} else {
			inactiveInstances = append(inactiveInstances, i)
		}
		newestInstance = tempInstance
	}
	return newestInstance, inactiveInstances
}

func chooseNewest(instance1, instance2 *ec2.Instance) *ec2.Instance {
	if instance1.LaunchTime == nil && instance2.LaunchTime == nil {
		// No idea what to do here, should not be possible, just return the first.
		return instance1
	}
	if instance1.LaunchTime != nil && instance2.LaunchTime == nil {
		return instance1
	}
	if instance1.LaunchTime == nil && instance2.LaunchTime != nil {
		return instance2
	}
	if (*instance1.LaunchTime).After(*instance2.LaunchTime) {
		return instance1
	}
	return instance2
}

// GetInstance returns the AWS instance for a given machine. If multiple instances match our machine,
// the most recently launched will be returned. If no instance exists, an error will be returned.
func GetInstance(machine *clusterv1.Machine, client awsclient.Client) (*ec2.Instance, error) {
	instances, err := GetRunningInstances(machine, client)
	if err != nil {
		return nil, err
	}
	if len(instances) == 0 {
		return nil, fmt.Errorf("no instance found for machine: %s", machine.Name)
	}

	instance, _ := SortInstances(instances)
	return instance, nil
}

// GetRunningInstances returns all running instances that have a tag matching our machine name,
// and cluster ID.
func GetRunningInstances(machine *clusterv1.Machine, client awsclient.Client) ([]*ec2.Instance, error) {
	runningInstanceStateFilter := []*string{aws.String(ec2.InstanceStateNameRunning), aws.String(ec2.InstanceStateNamePending)}
	return GetInstances(machine, client, runningInstanceStateFilter)
}

// GetStoppedInstances returns all stopped instances that have a tag matching our machine name,
// and cluster ID.
func GetStoppedInstances(machine *clusterv1.Machine, client awsclient.Client) ([]*ec2.Instance, error) {
	stoppedInstanceStateFilter := []*string{aws.String(ec2.InstanceStateNameStopped), aws.String(ec2.InstanceStateNameStopping)}
	return GetInstances(machine, client, stoppedInstanceStateFilter)
}

// GetInstances returns all instances that have a tag matching our machine name,
// and cluster ID.
func GetInstances(machine *clusterv1.Machine, client awsclient.Client, instanceStateFilter []*string) ([]*ec2.Instance, error) {

	machineName := machine.Name

	clusterID, ok := getClusterID(machine)
	if !ok {
		return []*ec2.Instance{}, fmt.Errorf("unable to get cluster ID for machine: %q", machine.Name)
	}

	requestFilters := []*ec2.Filter{
		{
			Name:   aws.String("tag:Name"),
			Values: []*string{&machineName},
		},
		{
			Name:   aws.String("tag:clusterid"),
			Values: []*string{&clusterID},
		},
	}

	if instanceStateFilter != nil {
		requestFilters = append(requestFilters, &ec2.Filter{
			Name:   aws.String("instance-state-name"),
			Values: instanceStateFilter,
		})
	}

	// Query instances with our machine's name, and in running/pending state.
	request := &ec2.DescribeInstancesInput{
		Filters: requestFilters,
	}

	result, err := client.DescribeInstances(request)
	if err != nil {
		return []*ec2.Instance{}, err
	}

	instances := make([]*ec2.Instance, 0, len(result.Reservations))
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			instances = append(instances, instance)
		}
	}

	return instances, nil
}

// TerminateInstances terminates all provided instances with a single EC2 request.
func TerminateInstances(client awsclient.Client, instances []*ec2.Instance, mLog log.FieldLogger) error {
	instanceIDs := []*string{}
	// Cleanup all older instances:
	for _, instance := range instances {
		mLog.WithFields(log.Fields{
			"instanceID": *instance.InstanceId,
			"state":      *instance.State.Name,
			"launchTime": *instance.LaunchTime,
		}).Warn("cleaning up extraneous instance for machine")
		instanceIDs = append(instanceIDs, instance.InstanceId)
	}
	for _, instanceID := range instanceIDs {
		mLog.WithField("instanceID", *instanceID).Info("terminating instance")
	}

	terminateInstancesRequest := &ec2.TerminateInstancesInput{
		InstanceIds: instanceIDs,
	}
	_, err := client.TerminateInstances(terminateInstancesRequest)
	if err != nil {
		mLog.Errorf("error terminating instances: %v", err)
		return fmt.Errorf("error terminating instances: %v", err)
	}
	return nil
}

// ProviderConfigFromClusterAPIMachineSpec gets the machine provider config MachineSetSpec from the
// specified cluster-api MachineSpec.
func ProviderConfigFromClusterAPIMachineSpec(ms *clusterv1.MachineSpec) (*providerconfigv1.AWSMachineProviderConfig, error) {
	if ms.ProviderConfig.Value == nil {
		return nil, fmt.Errorf("no Value in ProviderConfig")
	}
	obj, gvk, err := providerconfigv1.Codecs.UniversalDecoder(providerconfigv1.SchemeGroupVersion).Decode([]byte(ms.ProviderConfig.Value.Raw), nil, nil)
	if err != nil {
		return nil, err
	}
	spec, ok := obj.(*providerconfigv1.AWSMachineProviderConfig)
	if !ok {
		return nil, fmt.Errorf("unexpected object when parsing machine provider config: %#v", gvk)
	}
	return spec, nil
}

// AWSMachineProviderStatusFromClusterAPIMachine gets the machine provider status from the specified machine.
func AWSMachineProviderStatusFromClusterAPIMachine(m *clusterv1.Machine) (*providerconfigv1.AWSMachineProviderStatus, error) {
	return AWSMachineProviderStatusFromMachineStatus(&m.Status)
}

// AWSMachineProviderStatusFromMachineStatus gets the machine provider status from the specified machine status.
func AWSMachineProviderStatusFromMachineStatus(s *clusterv1.MachineStatus) (*providerconfigv1.AWSMachineProviderStatus, error) {
	if s.ProviderStatus == nil {
		return &providerconfigv1.AWSMachineProviderStatus{}, nil
	}
	obj, gvk, err := providerconfigv1.Codecs.UniversalDecoder(providerconfigv1.SchemeGroupVersion).Decode([]byte(s.ProviderStatus.Raw), nil, nil)
	if err != nil {
		return nil, err
	}
	status, ok := obj.(*providerconfigv1.AWSMachineProviderStatus)
	if !ok {
		return nil, fmt.Errorf("unexpected object: %#v", gvk)
	}
	return status, nil
}

// EncodeAWSMachineProviderStatus encodes the machine status into RawExtension
func EncodeAWSMachineProviderStatus(awsStatus *providerconfigv1.AWSMachineProviderStatus) (*runtime.RawExtension, error) {
	awsStatus.TypeMeta = metav1.TypeMeta{
		APIVersion: providerconfigv1.SchemeGroupVersion.String(),
		Kind:       "AWSMachineProviderStatus",
	}
	serializer := jsonserializer.NewSerializer(jsonserializer.DefaultMetaFactory, providerconfigv1.Scheme, providerconfigv1.Scheme, false)
	var buffer bytes.Buffer
	err := serializer.Encode(awsStatus, &buffer)
	if err != nil {
		return nil, err
	}
	return &runtime.RawExtension{
		Raw: bytes.TrimSpace(buffer.Bytes()),
	}, nil
}

// IsMaster returns true if the machine is part of a cluster's control plane
func IsMaster(machine *clusterv1.Machine) bool {
	if machineType, exists := machine.ObjectMeta.Labels[providerconfigv1.MachineTypeLabel]; exists && machineType == "master" {
		return true
	}
	return false
}

// IsInfra returns true if the machine is part of a cluster's infra plane
func IsInfra(machine *clusterv1.Machine) bool {
	if machineRole, exists := machine.ObjectMeta.Labels[providerconfigv1.MachineRoleLabel]; exists && machineRole == "infra" {
		return true
	}
	return false
}

// StringPtrsEqual safely returns true if the value for each string pointer is equal, or both are nil.
func StringPtrsEqual(s1, s2 *string) bool {
	if s1 == s2 {
		return true
	}
	if s1 == nil || s2 == nil {
		return false
	}
	return *s1 == *s2
}
