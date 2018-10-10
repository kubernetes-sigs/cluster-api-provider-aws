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
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
	service "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services"
	ec2svc "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/ec2"
	elbsvc "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/elb"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	client "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1"
	controllerError "sigs.k8s.io/cluster-api/pkg/controller/error"
)

// codec are the functions off the generated codec that this actuator uses.
type codec interface {
	DecodeFromProviderConfig(clusterv1.ProviderConfig, runtime.Object) error
	DecodeProviderStatus(*runtime.RawExtension, runtime.Object) error
	EncodeProviderStatus(runtime.Object) (*runtime.RawExtension, error)
}

// Actuator is responsible for performing machine reconciliation.
type Actuator struct {
	codec codec

	machinesGetter client.MachinesGetter
	servicesGetter service.Getter
}

// ActuatorParams holds parameter information for Actuator.
type ActuatorParams struct {
	Codec          codec
	MachinesGetter client.MachinesGetter
	ServicesGetter service.Getter
}

// NewActuator returns an actuator.
func NewActuator(params ActuatorParams) (*Actuator, error) {
	res := &Actuator{
		codec:          params.Codec,
		machinesGetter: params.MachinesGetter,
		servicesGetter: params.ServicesGetter,
	}

	if res.servicesGetter == nil {
		res.servicesGetter = new(defaultServicesGetter)
	}

	return res, nil
}

func (a *Actuator) ec2(clusterConfig *v1alpha1.AWSClusterProviderConfig) service.EC2MachineInterface {
	return a.servicesGetter.EC2(a.servicesGetter.Session(clusterConfig))
}

func (a *Actuator) elb(clusterConfig *v1alpha1.AWSClusterProviderConfig) service.ELBInterface {
	return a.servicesGetter.ELB(a.servicesGetter.Session(clusterConfig))
}

func (a *Actuator) machineProviderConfig(machine *clusterv1.Machine) (*v1alpha1.AWSMachineProviderConfig, error) {
	machineProviderCfg := &v1alpha1.AWSMachineProviderConfig{}
	err := a.codec.DecodeFromProviderConfig(machine.Spec.ProviderConfig, machineProviderCfg)
	return machineProviderCfg, err
}

func (a *Actuator) machineProviderStatus(machine *clusterv1.Machine) (*v1alpha1.AWSMachineProviderStatus, error) {
	status := &v1alpha1.AWSMachineProviderStatus{}
	err := a.codec.DecodeProviderStatus(machine.Status.ProviderStatus, status)
	return status, err
}

func (a *Actuator) clusterProviderStatus(cluster *clusterv1.Cluster) (*v1alpha1.AWSClusterProviderStatus, error) {
	providerStatus := &v1alpha1.AWSClusterProviderStatus{}
	err := a.codec.DecodeProviderStatus(cluster.Status.ProviderStatus, providerStatus)
	return providerStatus, err
}

func (a *Actuator) clusterProviderConfig(cluster *clusterv1.Cluster) (*v1alpha1.AWSClusterProviderConfig, error) {
	providerConfig := &v1alpha1.AWSClusterProviderConfig{}
	err := a.codec.DecodeFromProviderConfig(cluster.Spec.ProviderConfig, providerConfig)
	return providerConfig, err
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

	clusterConfig, err := a.clusterProviderConfig(cluster)
	if err != nil {
		return errors.Wrap(err, "failed to get cluster provider config")
	}

	config, err := a.machineProviderConfig(machine)
	if err != nil {
		return errors.Wrap(err, "failed to get machine config")
	}

	defer func() {
		if err := a.updateStatus(machine, status); err != nil {
			glog.Errorf("failed to store provider status for machine %q: %v", machine.Name, err)
		}
	}()

	i, err := a.ec2(clusterConfig).CreateOrGetMachine(machine, status, config, clusterStatus, cluster)
	if err != nil {

		if ec2svc.IsFailedDependency(errors.Cause(err)) {
			glog.Errorf("network not ready to launch instances yet: %s", err)
			return &controllerError.RequeueAfterError{
				RequeueAfter: time.Minute,
			}
		}

		return errors.Wrap(err, "failed to create or get machine")
	}

	status.InstanceID = &i.ID
	status.InstanceState = aws.String(string(i.State))

	if machine.Annotations == nil {
		machine.Annotations = map[string]string{}
	}

	machine.Annotations["cluster-api-provider-aws"] = "true"

	if err := a.reconcileLBAttachment(cluster, machine, i); err != nil {
		return errors.Wrap(err, "failed to reconcile LB attachment")
	}

	return nil
}

func (a *Actuator) reconcileLBAttachment(c *clusterv1.Cluster, m *clusterv1.Machine, i *v1alpha1.Instance) error {
	clusterConfig, err := a.clusterProviderConfig(c)
	if err != nil {
		return errors.Wrap(err, "failed to get cluster provider config")
	}

	if m.ObjectMeta.Labels["set"] == "controlplane" {
		if err := a.elb(clusterConfig).RegisterInstanceWithAPIServerELB(c.Name, i.ID); err != nil {
			return errors.Wrapf(err, "could not register control plane instance %q with load balancer", i.ID)
		}
	}

	return nil
}

// Delete deletes a machine and is invoked by the Machine Controller
func (a *Actuator) Delete(cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	glog.Infof("Deleting machine %v for cluster %v.", machine.Name, cluster.Name)

	status, err := a.machineProviderStatus(machine)
	if err != nil {
		return errors.Wrap(err, "failed to get machine provider status")
	}

	clusterConfig, err := a.clusterProviderConfig(cluster)
	if err != nil {
		return errors.Wrap(err, "failed to get cluster provider config")
	}

	if status.InstanceID == nil {
		// Instance was never created
		return nil
	}

	ec2svc := a.ec2(clusterConfig)

	instance, err := ec2svc.InstanceIfExists(status.InstanceID)
	if err != nil {
		return errors.Wrap(err, "failed to get instance")
	}

	if instance == nil {
		// The machine hasn't been created yet
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
		if err := ec2svc.TerminateInstance(aws.StringValue(status.InstanceID)); err != nil {
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

	clusterConfig, err := a.clusterProviderConfig(cluster)
	if err != nil {
		return errors.Wrap(err, "failed to get cluster provider config")
	}

	ec2svc := a.ec2(clusterConfig)

	// Get the current instance description from AWS.
	instanceDescription, err := ec2svc.InstanceIfExists(status.InstanceID)
	if err != nil {
		return errors.Wrap(err, "failed to get instance")
	}

	// We can now compare the various AWS state to the state we were passed.
	// We will check immutable state first, in order to fail quickly before
	// moving on to state that we can mutate.
	// TODO: Implement immutable state check.

	// Ensure that the security groups are correct.
	securityGroupsChanged, err := a.ensureSecurityGroups(
		ec2svc,
		machine,
		*status.InstanceID,
		config.AdditionalSecurityGroups,
		instanceDescription.SecurityGroupIDs,
	)
	if err != nil {
		return errors.Wrap(err, "failed to ensure security groups")
	}

	// Ensure that the tags are correct.
	tagsChanged, err := a.ensureTags(ec2svc, machine, status.InstanceID, config.AdditionalTags)
	if err != nil {
		return errors.Wrap(err, "failed to ensure tags")
	}

	// We need to update the machine since annotations may have changed.
	if securityGroupsChanged || tagsChanged {
		if err := a.updateMachine(machine); err != nil {
			return errors.Wrap(err, "failed to update machine")
		}
	}

	// Finally update the machine status.
	if err := a.updateStatus(machine, status); err != nil {
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

	clusterConfig, err := a.clusterProviderConfig(cluster)
	if err != nil {
		return false, errors.Wrap(err, "failed to get cluster provider config")
	}

	// TODO worry about pointers. instance if exists returns *any* instance
	if status.InstanceID == nil {
		return false, nil
	}

	instance, err := a.ec2(clusterConfig).InstanceIfExists(status.InstanceID)
	if err != nil {
		return false, err
	}

	if instance == nil {
		return false, nil
	}

	glog.Infof("Found an instance: %v", instance)

	switch instance.State {
	case v1alpha1.InstanceStateRunning:
		glog.Infof("Machine %v is running", status.InstanceID)
	case v1alpha1.InstanceStatePending:
		glog.Infof("Machine %v is pending", status.InstanceID)
	default:
		return false, nil
	}

	if err := a.reconcileLBAttachment(cluster, machine, instance); err != nil {
		return true, err
	}

	return true, nil
}

func (a *Actuator) updateMachine(machine *clusterv1.Machine) error {
	machinesClient := a.machinesGetter.Machines(machine.Namespace)

	if _, err := machinesClient.Update(machine); err != nil {
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

type defaultServicesGetter struct{}

func (d *defaultServicesGetter) Session(clusterConfig *v1alpha1.AWSClusterProviderConfig) *session.Session {
	return session.Must(session.NewSession(aws.NewConfig().WithRegion(clusterConfig.Region)))
}

func (d *defaultServicesGetter) EC2(session *session.Session) service.EC2Interface {
	return ec2svc.NewService(ec2.New(session))
}

func (d *defaultServicesGetter) ELB(session *session.Session) service.ELBInterface {
	return elbsvc.NewService(elb.New(session))
}
