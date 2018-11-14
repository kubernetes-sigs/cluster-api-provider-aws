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
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"

	"github.com/golang/glog"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	providerconfigv1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1alpha1"
	awsclient "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/client"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/ghodss/yaml"
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

// GetRunningInstance returns the AWS instance for a given machine. If multiple instances match our machine,
// the most recently launched will be returned. If no instance exists, an error will be returned.
func GetRunningInstance(machine *clusterv1.Machine, client awsclient.Client) (*ec2.Instance, error) {
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
func TerminateInstances(client awsclient.Client, instances []*ec2.Instance) error {
	instanceIDs := []*string{}
	// Cleanup all older instances:
	for _, instance := range instances {
		glog.Infof("Cleaning up extraneous instance for machine: %v, state: %v, launchTime: %v", *instance.InstanceId, *instance.State.Name, *instance.LaunchTime)
		instanceIDs = append(instanceIDs, instance.InstanceId)
	}
	for _, instanceID := range instanceIDs {
		glog.Infof("Terminating %v instance", *instanceID)
	}

	terminateInstancesRequest := &ec2.TerminateInstancesInput{
		InstanceIds: instanceIDs,
	}
	_, err := client.TerminateInstances(terminateInstancesRequest)
	if err != nil {
		glog.Errorf("error terminating instances: %v", err)
		return fmt.Errorf("error terminating instances: %v", err)
	}
	return nil
}

// ProviderConfigFromMachine gets the machine provider config MachineSetSpec from the
// specified cluster-api MachineSpec.
func ProviderConfigFromMachine(machine *clusterv1.Machine) (*providerconfigv1.AWSMachineProviderConfig, error) {
	var config providerconfigv1.AWSMachineProviderConfig
	if err := yaml.Unmarshal(machine.Spec.ProviderConfig.Value.Raw, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// ProviderStatusFromMachine gets the machine provider status from the specified machine.
func ProviderStatusFromMachine(codec codec, m *clusterv1.Machine) (*providerconfigv1.AWSMachineProviderStatus, error) {
	status := &providerconfigv1.AWSMachineProviderStatus{}
	err := codec.DecodeProviderStatus(m.Status.ProviderStatus, status)
	return status, err
}

// EncodeProviderStatus encodes the machine status into RawExtension
func EncodeProviderStatus(codec codec, awsStatus *providerconfigv1.AWSMachineProviderStatus) (*runtime.RawExtension, error) {
	return codec.EncodeProviderStatus(awsStatus)
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

// UpdateConditionCheck tests whether a condition should be updated from the
// old condition to the new condition. Returns true if the condition should
// be updated.
type UpdateConditionCheck func(oldReason, oldMessage, newReason, newMessage string) bool

// UpdateConditionAlways returns true. The condition will always be updated.
func UpdateConditionAlways(_, _, _, _ string) bool {
	return true
}

// UpdateConditionNever return false. The condition will never be updated,
// unless there is a change in the status of the condition.
func UpdateConditionNever(_, _, _, _ string) bool {
	return false
}

// UpdateConditionIfReasonOrMessageChange returns true if there is a change
// in the reason or the message of the condition.
func UpdateConditionIfReasonOrMessageChange(oldReason, oldMessage, newReason, newMessage string) bool {
	return oldReason != newReason ||
		oldMessage != newMessage
}

func shouldUpdateCondition(
	oldStatus corev1.ConditionStatus, oldReason, oldMessage string,
	newStatus corev1.ConditionStatus, newReason, newMessage string,
	updateConditionCheck UpdateConditionCheck,
) bool {
	if oldStatus != newStatus {
		return true
	}
	return updateConditionCheck(oldReason, oldMessage, newReason, newMessage)
}

// SetAWSMachineProviderCondition sets the condition for the machine and
// returns the new slice of conditions.
// If the machine does not already have a condition with the specified type,
// a condition will be added to the slice if and only if the specified
// status is True.
// If the machine does already have a condition with the specified type,
// the condition will be updated if either of the following are true.
// 1) Requested status is different than existing status.
// 2) The updateConditionCheck function returns true.
func SetAWSMachineProviderCondition(
	conditions []providerconfigv1.AWSMachineProviderCondition,
	conditionType providerconfigv1.AWSMachineProviderConditionType,
	status corev1.ConditionStatus,
	reason string,
	message string,
	updateConditionCheck UpdateConditionCheck,
) []providerconfigv1.AWSMachineProviderCondition {
	now := metav1.Now()
	existingCondition := FindAWSMachineProviderCondition(conditions, conditionType)
	if existingCondition == nil {
		if status == corev1.ConditionTrue {
			conditions = append(
				conditions,
				providerconfigv1.AWSMachineProviderCondition{
					Type:               conditionType,
					Status:             status,
					Reason:             reason,
					Message:            message,
					LastTransitionTime: now,
					LastProbeTime:      now,
				},
			)
		}
	} else {
		if shouldUpdateCondition(
			existingCondition.Status, existingCondition.Reason, existingCondition.Message,
			status, reason, message,
			updateConditionCheck,
		) {
			if existingCondition.Status != status {
				existingCondition.LastTransitionTime = now
			}
			existingCondition.Status = status
			existingCondition.Reason = reason
			existingCondition.Message = message
			existingCondition.LastProbeTime = now
		}
	}
	return conditions
}

// FindAWSMachineProviderCondition finds in the machine the condition that has the
// specified condition type. If none exists, then returns nil.
func FindAWSMachineProviderCondition(conditions []providerconfigv1.AWSMachineProviderCondition, conditionType providerconfigv1.AWSMachineProviderConditionType) *providerconfigv1.AWSMachineProviderCondition {
	for i, condition := range conditions {
		if condition.Type == conditionType {
			return &conditions[i]
		}
	}
	return nil
}
