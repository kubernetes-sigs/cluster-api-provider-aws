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

	"github.com/golang/glog"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	machinev1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/types"
	providerconfigv1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1beta1"
	awsclient "sigs.k8s.io/cluster-api-provider-aws/pkg/client"
)

// getRunningInstance returns the AWS instance for a given machine. If multiple instances match our machine,
// the most recently launched will be returned. If no instance exists, an error will be returned.
func getRunningInstance(machine *machinev1.Machine, client awsclient.Client) (*ec2.Instance, error) {
	instances, err := getRunningInstances(machine, client)
	if err != nil {
		return nil, err
	}
	if len(instances) == 0 {
		return nil, fmt.Errorf("no instance found for machine: %s", machine.Name)
	}

	sortInstances(instances)
	return instances[0], nil
}

// getRunningInstances returns all running instances that have a tag matching our machine name,
// and cluster ID.
func getRunningInstances(machine *machinev1.Machine, client awsclient.Client) ([]*ec2.Instance, error) {
	runningInstanceStateFilter := []*string{aws.String(ec2.InstanceStateNameRunning), aws.String(ec2.InstanceStateNamePending)}
	return getInstances(machine, client, runningInstanceStateFilter)
}

// getStoppedInstances returns all stopped instances that have a tag matching our machine name,
// and cluster ID.
func getStoppedInstances(machine *machinev1.Machine, client awsclient.Client) ([]*ec2.Instance, error) {
	stoppedInstanceStateFilter := []*string{aws.String(ec2.InstanceStateNameStopped), aws.String(ec2.InstanceStateNameStopping)}
	return getInstances(machine, client, stoppedInstanceStateFilter)
}

func getExistingInstances(machine *machinev1.Machine, client awsclient.Client) ([]*ec2.Instance, error) {
	return getInstances(machine, client, []*string{
		aws.String(ec2.InstanceStateNameRunning),
		aws.String(ec2.InstanceStateNamePending),
		aws.String(ec2.InstanceStateNameStopped),
		aws.String(ec2.InstanceStateNameStopping),
		aws.String(ec2.InstanceStateNameShuttingDown),
	})
}

// getInstances returns all instances that have a tag matching our machine name,
// and cluster ID.
func getInstances(machine *machinev1.Machine, client awsclient.Client, instanceStateFilter []*string) ([]*ec2.Instance, error) {

	clusterID, ok := getClusterID(machine)
	if !ok {
		return []*ec2.Instance{}, fmt.Errorf("unable to get cluster ID for machine: %q", machine.Name)
	}

	requestFilters := []*ec2.Filter{
		{
			Name:   awsTagFilter("Name"),
			Values: aws.StringSlice([]string{machine.Name}),
		},
		clusterFilter(clusterID),
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

func getVolume(client awsclient.Client, volumeID string) (*ec2.Volume, error) {
	request := &ec2.DescribeVolumesInput{
		VolumeIds: []*string{&volumeID},
	}
	result, err := client.DescribeVolumes(request)
	if err != nil {
		return &ec2.Volume{}, err
	}

	if len(result.Volumes) != 1 {
		return &ec2.Volume{}, fmt.Errorf("unable to get volume ID: %q", volumeID)
	}

	return result.Volumes[0], nil
}

// terminateInstances terminates all provided instances with a single EC2 request.
func terminateInstances(client awsclient.Client, instances []*ec2.Instance) error {
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
		glog.Errorf("Error terminating instances: %v", err)
		return fmt.Errorf("error terminating instances: %v", err)
	}
	return nil
}

// providerConfigFromMachine gets the machine provider config MachineSetSpec from the
// specified cluster-api MachineSpec.
func providerConfigFromMachine(machine *machinev1.Machine, codec *providerconfigv1.AWSProviderConfigCodec) (*providerconfigv1.AWSMachineProviderConfig, error) {
	if machine.Spec.ProviderSpec.Value == nil {
		return nil, fmt.Errorf("unable to find machine provider config: Spec.ProviderSpec.Value is not set")
	}

	var config providerconfigv1.AWSMachineProviderConfig
	if err := codec.DecodeProviderSpec(&machine.Spec.ProviderSpec, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// isMaster returns true if the machine is part of a cluster's control plane
func (a *Actuator) isMaster(machine *machinev1.Machine) (bool, error) {
	if machine.Status.NodeRef == nil {
		glog.Errorf("NodeRef not found in machine %s", machine.Name)
		return false, nil
	}
	node := &corev1.Node{}
	nodeKey := types.NamespacedName{
		Namespace: machine.Status.NodeRef.Namespace,
		Name:      machine.Status.NodeRef.Name,
	}

	err := a.client.Get(context.Background(), nodeKey, node)
	if err != nil {
		return false, fmt.Errorf("failed to get node from machine %s", machine.Name)
	}

	if _, exists := node.Labels["node-role.kubernetes.io/master"]; exists {
		return true, nil
	}
	return false, nil
}

// updateConditionCheck tests whether a condition should be updated from the
// old condition to the new condition. Returns true if the condition should
// be updated.
type updateConditionCheck func(oldReason, oldMessage, newReason, newMessage string) bool

// updateConditionAlways returns true. The condition will always be updated.
func updateConditionAlways(_, _, _, _ string) bool {
	return true
}

// updateConditionNever return false. The condition will never be updated,
// unless there is a change in the status of the condition.
func updateConditionNever(_, _, _, _ string) bool {
	return false
}

// updateConditionIfReasonOrMessageChange returns true if there is a change
// in the reason or the message of the condition.
func updateConditionIfReasonOrMessageChange(oldReason, oldMessage, newReason, newMessage string) bool {
	return oldReason != newReason ||
		oldMessage != newMessage
}

func shouldUpdateCondition(
	oldStatus corev1.ConditionStatus, oldReason, oldMessage string,
	newStatus corev1.ConditionStatus, newReason, newMessage string,
	updateConditionCheck updateConditionCheck,
) bool {
	if oldStatus != newStatus {
		return true
	}
	return updateConditionCheck(oldReason, oldMessage, newReason, newMessage)
}

// setAWSMachineProviderCondition sets the condition for the machine and
// returns the new slice of conditions.
// If the machine does not already have a condition with the specified type,
// a condition will be added to the slice if and only if the specified
// status is True.
// If the machine does already have a condition with the specified type,
// the condition will be updated if either of the following are true.
// 1) Requested status is different than existing status.
// 2) The updateConditionCheck function returns true.
func setAWSMachineProviderCondition(
	conditions []providerconfigv1.AWSMachineProviderCondition,
	conditionType providerconfigv1.AWSMachineProviderConditionType,
	status corev1.ConditionStatus,
	reason string,
	message string,
	updateConditionCheck updateConditionCheck,
) []providerconfigv1.AWSMachineProviderCondition {
	now := metav1.Now()
	existingCondition := findAWSMachineProviderCondition(conditions, conditionType)
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

// findAWSMachineProviderCondition finds in the machine the condition that has the
// specified condition type. If none exists, then returns nil.
func findAWSMachineProviderCondition(conditions []providerconfigv1.AWSMachineProviderCondition, conditionType providerconfigv1.AWSMachineProviderConditionType) *providerconfigv1.AWSMachineProviderCondition {
	for i, condition := range conditions {
		if condition.Type == conditionType {
			return &conditions[i]
		}
	}
	return nil
}
