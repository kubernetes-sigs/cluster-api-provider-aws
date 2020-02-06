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
	"net"
	"strings"

	"github.com/golang/glog"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	machinev1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/types"
	providerconfigv1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1beta1"
	awsclient "sigs.k8s.io/cluster-api-provider-aws/pkg/client"
)

// existingInstanceStates returns the list of states an EC2 instance can be in
// while being considered "existing", i.e. mostly anything but "Terminated".
func existingInstanceStates() []*string {
	return []*string{
		aws.String(ec2.InstanceStateNameRunning),
		aws.String(ec2.InstanceStateNamePending),
		aws.String(ec2.InstanceStateNameStopped),
		aws.String(ec2.InstanceStateNameStopping),
		aws.String(ec2.InstanceStateNameShuttingDown),
	}
}

// getRunningFromInstances returns all running instances from a list of instances.
func getRunningFromInstances(instances []*ec2.Instance) []*ec2.Instance {
	var runningInstances []*ec2.Instance
	for _, instance := range instances {
		if *instance.State.Name == ec2.InstanceStateNameRunning {
			runningInstances = append(runningInstances, instance)
		}
	}
	return runningInstances
}

// getStoppedInstances returns all stopped instances that have a tag matching our machine name,
// and cluster ID.
func getStoppedInstances(machine *machinev1.Machine, client awsclient.Client) ([]*ec2.Instance, error) {
	stoppedInstanceStateFilter := []*string{aws.String(ec2.InstanceStateNameStopped), aws.String(ec2.InstanceStateNameStopping)}
	return getInstances(machine, client, stoppedInstanceStateFilter)
}

// getExistingInstances returns all instances not terminated
func getExistingInstances(machine *machinev1.Machine, client awsclient.Client) ([]*ec2.Instance, error) {
	return getInstances(machine, client, existingInstanceStates())
}

func getExistingInstanceByID(id string, client awsclient.Client) (*ec2.Instance, error) {
	return getInstanceByID(id, client, existingInstanceStates())
}

func instanceHasAllowedState(instance *ec2.Instance, instanceStateFilter []*string) error {
	if instance.InstanceId == nil {
		return fmt.Errorf("instance has nil ID")
	}

	if instance.State == nil {
		return fmt.Errorf("instance %s has nil state", *instance.InstanceId)
	}

	if len(instanceStateFilter) == 0 {
		return nil
	}

	actualState := aws.StringValue(instance.State.Name)
	for _, allowedState := range instanceStateFilter {
		if aws.StringValue(allowedState) == actualState {
			return nil
		}
	}

	allowedStates := make([]string, 0, len(instanceStateFilter))
	for _, allowedState := range instanceStateFilter {
		allowedStates = append(allowedStates, aws.StringValue(allowedState))
	}
	return fmt.Errorf("instance %s state %q is not in %s", *instance.InstanceId, actualState, strings.Join(allowedStates, ", "))
}

// getInstanceByID returns the instance with the given ID if it exists.
func getInstanceByID(id string, client awsclient.Client, instanceStateFilter []*string) (*ec2.Instance, error) {
	if id == "" {
		return nil, fmt.Errorf("instance-id not specified")
	}

	request := &ec2.DescribeInstancesInput{
		InstanceIds: aws.StringSlice([]string{id}),
	}

	result, err := client.DescribeInstances(request)
	if err != nil {
		return nil, err
	}

	if len(result.Reservations) != 1 {
		return nil, fmt.Errorf("found %d reservations for instance-id %s", len(result.Reservations), id)
	}

	reservation := result.Reservations[0]

	if len(reservation.Instances) != 1 {
		return nil, fmt.Errorf("found %d instances for instance-id %s", len(reservation.Instances), id)
	}

	instance := reservation.Instances[0]

	return instance, instanceHasAllowedState(instance, instanceStateFilter)
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
			err := instanceHasAllowedState(instance, instanceStateFilter)
			if err != nil {
				glog.Errorf("Excluding instance matching %s: %v", machine.Name, err)
			} else {
				instances = append(instances, instance)
			}
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
func terminateInstances(client awsclient.Client, instances []*ec2.Instance) ([]*ec2.InstanceStateChange, error) {
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
	output, err := client.TerminateInstances(terminateInstancesRequest)
	if err != nil {
		glog.Errorf("Error terminating instances: %v", err)
		return nil, fmt.Errorf("error terminating instances: %v", err)
	}

	if output == nil {
		return nil, nil
	}

	return output.TerminatingInstances, nil
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

// extractNodeAddresses maps the instance information from EC2 to an array of NodeAddresses
func extractNodeAddresses(instance *ec2.Instance) ([]corev1.NodeAddress, error) {
	// Not clear if the order matters here, but we might as well indicate a sensible preference order

	if instance == nil {
		return nil, fmt.Errorf("nil instance passed to extractNodeAddresses")
	}

	addresses := []corev1.NodeAddress{}

	// handle internal network interfaces
	for _, networkInterface := range instance.NetworkInterfaces {
		// skip network interfaces that are not currently in use
		if aws.StringValue(networkInterface.Status) != ec2.NetworkInterfaceStatusInUse {
			continue
		}

		// Treating IPv6 addresses as type NodeInternalIP to match what the KNI
		// patch to the AWS cloud-provider code is doing:
		//
		// https://github.com/openshift-kni/origin/commit/7db21c1e26a344e25ae1b825d4f21e7bef5c3650
		for _, ipv6Address := range networkInterface.Ipv6Addresses {
			if addr := aws.StringValue(ipv6Address.Ipv6Address); addr != "" {
				ip := net.ParseIP(addr)
				if ip == nil {
					return nil, fmt.Errorf("EC2 instance had invalid IPv6 address: %s (%q)", aws.StringValue(instance.InstanceId), addr)
				}
				addresses = append(addresses, corev1.NodeAddress{Type: corev1.NodeInternalIP, Address: ip.String()})
			}
		}

		for _, internalIP := range networkInterface.PrivateIpAddresses {
			if ipAddress := aws.StringValue(internalIP.PrivateIpAddress); ipAddress != "" {
				ip := net.ParseIP(ipAddress)
				if ip == nil {
					return nil, fmt.Errorf("EC2 instance had invalid private address: %s (%q)", aws.StringValue(instance.InstanceId), ipAddress)
				}
				addresses = append(addresses, corev1.NodeAddress{Type: corev1.NodeInternalIP, Address: ip.String()})
			}
		}
	}

	// TODO: Other IP addresses (multiple ips)?
	publicIPAddress := aws.StringValue(instance.PublicIpAddress)
	if publicIPAddress != "" {
		ip := net.ParseIP(publicIPAddress)
		if ip == nil {
			return nil, fmt.Errorf("EC2 instance had invalid public address: %s (%s)", aws.StringValue(instance.InstanceId), publicIPAddress)
		}
		addresses = append(addresses, corev1.NodeAddress{Type: corev1.NodeExternalIP, Address: ip.String()})
	}

	privateDNSName := aws.StringValue(instance.PrivateDnsName)
	if privateDNSName != "" {
		addresses = append(addresses, corev1.NodeAddress{Type: corev1.NodeInternalDNS, Address: privateDNSName})
		addresses = append(addresses, corev1.NodeAddress{Type: corev1.NodeHostName, Address: privateDNSName})
	}

	publicDNSName := aws.StringValue(instance.PublicDnsName)
	if publicDNSName != "" {
		addresses = append(addresses, corev1.NodeAddress{Type: corev1.NodeExternalDNS, Address: publicDNSName})
	}

	return addresses, nil
}
