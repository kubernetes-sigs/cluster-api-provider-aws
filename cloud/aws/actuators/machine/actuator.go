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

	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
	ec2svc "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/ec2"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

// machinesSvc are the functions of the cluster-api that this actuator needs.
type machinesSvc interface {
	UpdateMachineStatus(*clusterv1.Machine) (*clusterv1.Machine, error)
}

// ec2Svc are the functions from the ec2 service, not the client, this actuator needs.
// This should never need to import the ec2 sdk.
type ec2Svc interface {
	CreateInstance(*clusterv1.Machine) (*ec2svc.Instance, error)
	InstanceIfExists(*string) (*ec2svc.Instance, error)
	TerminateInstance(*string) error
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
	ec2      ec2Svc
	machines machinesSvc
}

// ActuatorParams holds parameter information for Actuator
type ActuatorParams struct {
	// Codec is needed to work with the provider configs and statuses.
	Codec codec

	// Services

	// ClusterService is the interface to cluster-api.
	MachinesService machinesSvc
	// EC2Service is the interface to ec2.
	EC2Service ec2Svc
}

// NewActuator returns an actuator.
func NewActuator(params ActuatorParams) (*Actuator, error) {
	return &Actuator{
		codec:    params.Codec,
		ec2:      params.EC2Service,
		machines: params.MachinesService,
	}, nil
}

// Create creates a machine and is invoked by the machine controller.
func (a *Actuator) Create(cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	// will need this machine config in a bit
	_, err := a.machineProviderConfig(machine.Spec.ProviderConfig)
	if err != nil {
		glog.Errorf("Failed to decode the machine provider config: %v", err)
		return err
	}

	// Get the machine status
	status, err := a.machineProviderStatus(machine)
	if err != nil {
		return err
	}

	// does the instance exist with a valid status? we're good
	// otherwise create it and move on.
	_, err = a.ec2.InstanceIfExists(status.InstanceID)
	if err != nil {
		return err
	}

	i, err := a.ec2.CreateInstance(machine)
	if err != nil {
		return err
	}

	status.InstanceID = &i.ID
	status.InstanceState = &i.State
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

	// Check the instance state. If it's already shutting down or terminated,
	// do nothing. Otherwise attempt to delete it.
	// This decision is based on the ec2-instance-lifecycle graph at
	// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-lifecycle.html
	switch instance.State {
	case ec2svc.InstanceStateShuttingDown, ec2svc.InstanceStateTerminated:
		return nil
	default:
		err = a.ec2.TerminateInstance(status.InstanceID)
		if err != nil {
			return errors.Wrap(err, "failed to terminate instance")
		}
	}

	return nil
}

// Update updates a machine and is invoked by the Machine Controller
func (a *Actuator) Update(cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	glog.Infof("Updating machine %v for cluster %v.", machine.Name, cluster.Name)
	return fmt.Errorf("TODO: Not yet implemented")
}

// Exists test for the existence of a machine and is invoked by the Machine Controller
func (a *Actuator) Exists(cluster *clusterv1.Cluster, machine *clusterv1.Machine) (bool, error) {
	glog.Info("Checking if machine %v for cluster %v exists.", machine.Name, cluster.Name)
	status, err := a.machineProviderStatus(machine)
	if err != nil {
		return false, err
	}

	instance, err := a.ec2.InstanceIfExists(status.InstanceID)
	if err != nil {
		return false, err
	}
	if instance == nil {
		return false, nil
	}
	// TODO update status here
	switch instance.State {
	case ec2svc.InstanceStateRunning, ec2svc.InstanceStatePending:
		return true, nil
	default:
		return false, nil
	}
}

func (a *Actuator) machineProviderConfig(providerConfig clusterv1.ProviderConfig) (*v1alpha1.AWSMachineProviderConfig, error) {
	machineProviderCfg := &v1alpha1.AWSMachineProviderConfig{}
	err := a.codec.DecodeFromProviderConfig(providerConfig, machineProviderCfg)
	return machineProviderCfg, err
}

func (a *Actuator) machineProviderStatus(machine *clusterv1.Machine) (*v1alpha1.AWSMachineProviderStatus, error) {
	status := &v1alpha1.AWSMachineProviderStatus{}
	err := a.codec.DecodeProviderStatus(machine.Status.ProviderStatus, status)
	return status, err
}

func (a *Actuator) updateStatus(machine *clusterv1.Machine, status *v1alpha1.AWSMachineProviderStatus) error {
	encodedProviderStatus, err := a.codec.EncodeProviderStatus(status)
	if err != nil {
		return fmt.Errorf("failed to encode machine status: %v", err)
	}
	if encodedProviderStatus != nil {
		machine.Status.ProviderStatus = encodedProviderStatus
		if _, err := a.machines.UpdateMachineStatus(machine); err != nil {
			return fmt.Errorf("failed to update machine status: %v", err)
		}
	}
	return nil
}
