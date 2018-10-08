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

package machine

// should not need to import the ec2 sdk here
import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	client "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1"
)

// ec2Svc are the functions from the ec2 service, not the client, this actuator needs.
// This should never need to import the ec2 sdk.
type ec2Svc interface {
	CreateInstance(*clusterv1.Machine, *v1alpha1.AWSMachineProviderConfig, *v1alpha1.AWSClusterProviderStatus) (*v1alpha1.Instance, error)
	InstanceIfExists(*string) (*v1alpha1.Instance, error)
	TerminateInstance(string) error
	CreateOrGetMachine(*clusterv1.Machine, *v1alpha1.AWSMachineProviderStatus, *v1alpha1.AWSMachineProviderConfig, *v1alpha1.AWSClusterProviderStatus) (*v1alpha1.Instance, error)
	UpdateInstanceSecurityGroups(*string, []*string) error
	UpdateResourceTags(*string, map[string]string, map[string]string) error
}

// codec are the functions off the generated codec that this actuator uses.
type codec interface {
	DecodeFromProviderConfig(clusterv1.ProviderConfig, runtime.Object) error
	DecodeProviderStatus(*runtime.RawExtension, runtime.Object) error
	EncodeProviderStatus(runtime.Object) (*runtime.RawExtension, error)
}

// Actuator is responsible for performing machine reconciliation
type Actuator struct {
	codec codec

	// Services
	ec2            ec2Svc
	machinesGetter client.MachinesGetter
}

// ActuatorParams holds parameter information for Actuator
type ActuatorParams struct {
	// Codec is needed to work with the provider configs and statuses.
	Codec codec

	// Services

	// ClusterService is the interface to cluster-api.
	MachinesGetter client.MachinesGetter
	// EC2Service is the interface to ec2.
	EC2Service ec2Svc
}

// NewActuator returns an actuator.
func NewActuator(params ActuatorParams) (*Actuator, error) {
	return &Actuator{
		codec:          params.Codec,
		ec2:            params.EC2Service,
		machinesGetter: params.MachinesGetter,
	}, nil
}

// Create creates a machine and is invoked by the machine controller.
func (a *Actuator) Create(cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	glog.Infof("Creating machine %v for cluster %v", machine.Name, cluster.Name)
	status, err := a.machineProviderStatus(machine)
	if err != nil {
		return errors.Wrap(err, "failed to get machine provider status")
	}
	clusterStatus, err := a.clusterProviderStatus(cluster)
	if err != nil {
		return errors.Wrap(err, "failed to get cluster provider status")
	}
	config, err := a.machineProviderConfig(machine)
	if err != nil {
		return errors.Wrap(err, "failed to get machine config")
	}

	i, err := a.ec2.CreateOrGetMachine(machine, status, config, clusterStatus)
	if err != nil {
		return errors.Wrap(err, "failed to create or get machine")
	}

	status.InstanceID = &i.ID
	status.InstanceState = aws.String(string(i.State))

	if machine.Annotations == nil {
		machine.Annotations = map[string]string{}
	}
	machine.Annotations["cluster-api-provider-aws"] = "true"

	return a.updateStatus(machine, status)
}

// Delete deletes a machine and is invoked by the Machine Controller
func (a *Actuator) Delete(cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	glog.Infof("Deleting machine %v for cluster %v.", machine.Name, cluster.Name)

	status, err := a.machineProviderStatus(machine)
	if err != nil {
		return errors.Wrap(err, "failed to get machine provider status")
	}

	instance, err := a.ec2.InstanceIfExists(status.InstanceID)
	if err != nil {
		return errors.Wrap(err, "failed to get instance")
	}

	// The machine hasn't been created yet
	if instance == nil {
		return nil
	}

	// Check the instance state. If it's already shutting down or terminated,
	// do nothing. Otherwise attempt to delete it.
	// This decision is based on the ec2-instance-lifecycle graph at
	// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-lifecycle.html
	switch instance.State {
	case v1alpha1.InstanceStateShuttingDown, v1alpha1.InstanceStateTerminated:
		return nil
	default:
		err = a.ec2.TerminateInstance(aws.StringValue(status.InstanceID))
		if err != nil {
			return errors.Wrap(err, "failed to terminate instance")
		}
	}

	return nil
}

// Update updates a machine and is invoked by the Machine Controller.
// If the Update attempts to mutate any immutable state, the method will error
// and no updates will be performed.
func (a *Actuator) Update(cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	glog.Infof("Updating machine %v for cluster %v.", machine.Name, cluster.Name)

	// Get the updated config. We're going to compare parts of this to the
	// current Tags and Security Groups that AWS is aware of.
	config, err := a.machineProviderConfig(machine)
	if err != nil {
		return errors.Wrap(err, "failed to get machine config")
	}

	// Get the new status from the provided machine object.
	// We need this in case any of it has changed, and we also require the
	// instanceID for various AWS queries.
	status, err := a.machineProviderStatus(machine)
	if err != nil {
		return errors.Wrap(err, "failed to get machine status")
	}

	// Get the current instance description from AWS.
	instanceDescription, err := a.ec2.InstanceIfExists(status.InstanceID)
	if err != nil {
		return errors.Wrap(err, "failed to get instance")
	}

	// We can now compare the various AWS state to the state we were passed.
	// We will check immutable state first, in order to fail quickly before
	// moving on to state that we can mutate.
	// TODO: Implement immutable state check.

	// Ensure that the security groups are correct.
	securityGroupsChanged, err := a.ensureSecurityGroups(
		machine,
		status.InstanceID,
		config.AdditionalSecurityGroups,
		instanceDescription.SecurityGroups,
	)
	if err != nil {
		return errors.Wrap(err, "failed to ensure security groups")
	}

	// Ensure that the tags are correct.
	tagsChanged, err := a.ensureTags(machine, status.InstanceID, config.AdditionalTags)
	if err != nil {
		return errors.Wrap(err, "failed to ensure tags")
	}

	// We need to update the machine since annotations may have changed.
	if securityGroupsChanged || tagsChanged {
		err = a.updateMachine(machine)
		if err != nil {
			return errors.Wrap(err, "failed to update machine")
		}
	}

	// Finally update the machine status.
	err = a.updateStatus(machine, status)
	if err != nil {
		return errors.Wrap(err, "failed to update machine status")
	}

	return nil
}

// Exists test for the existence of a machine and is invoked by the Machine Controller
func (a *Actuator) Exists(cluster *clusterv1.Cluster, machine *clusterv1.Machine) (bool, error) {
	glog.Infof("Checking if machine %v for cluster %v exists", machine.Name, cluster.Name)
	status, err := a.machineProviderStatus(machine)
	if err != nil {
		return false, err
	}

	// TODO worry about pointers. instance if exists returns *any* instance
	if status.InstanceID == nil {
		return false, nil
	}

	instance, err := a.ec2.InstanceIfExists(status.InstanceID)
	if err != nil {
		return false, err
	}
	if instance == nil {
		return false, nil
	}

	glog.Infof("Found an instance: %#v", instance)

	switch instance.State {
	case v1alpha1.InstanceStateRunning, v1alpha1.InstanceStatePending:
		return true, nil
	default:
		return false, nil
	}
}

func (a *Actuator) updateMachine(machine *clusterv1.Machine) error {
	machinesClient := a.machinesGetter.Machines(machine.Namespace)

	_, err := machinesClient.Update(machine)
	if err != nil {
		return fmt.Errorf("failed to update machine: %v", err)
	}

	return nil
}

func (a *Actuator) updateStatus(machine *clusterv1.Machine, status *v1alpha1.AWSMachineProviderStatus) error {
	machinesClient := a.machinesGetter.Machines(machine.Namespace)
	encodedProviderStatus, err := a.codec.EncodeProviderStatus(status)
	if err != nil {
		return fmt.Errorf("failed to encode machine status: %v", err)
	}
	if encodedProviderStatus != nil {
		machine.Status.ProviderStatus = encodedProviderStatus
		if _, err := machinesClient.UpdateStatus(machine); err != nil {
			return fmt.Errorf("failed to update machine status: %v", err)
		}
	}
	return nil
}
