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
	"context"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
	"go.opencensus.io/trace"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/events"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/instrumentation"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/ec2"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/elb"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	client "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1"
	controllerError "sigs.k8s.io/cluster-api/pkg/controller/error"
	"time"
)

// ec2Svc are the functions from the ec2 service, not the client, this actuator needs.
// This should never need to import the ec2 sdk.
type ec2Svc interface {
	InstanceIfExists(context.Context, string) (*v1alpha1.Instance, error)
	TerminateInstance(context.Context, string) error
	ReconcileInstance(context.Context, *clusterv1.Machine, *v1alpha1.AWSMachineProviderStatus, *v1alpha1.AWSMachineProviderConfig, *v1alpha1.AWSClusterProviderStatus) (*v1alpha1.Instance, error)
	UpdateInstanceSecurityGroups(context.Context, string, []string) error
	UpdateResourceTags(context.Context, string, map[string]string, map[string]string) error
}

type elbSvc interface {
	RegisterInstanceWithClassicELB(context.Context, string, string) error
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
	elb            elbSvc
	machinesGetter client.MachinesGetter

	i instruments
}

type instruments struct {
	createCount *stats.Int64Measure
	deleteCount *stats.Int64Measure
	updateCount *stats.Int64Measure
	existsCount *stats.Int64Measure
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
	// ELBService is the interface to load balancing
	ELBService elbSvc
}

// NewActuator returns an actuator.
func NewActuator(params ActuatorParams) (*Actuator, error) {

	counters := instruments{
		createCount: instrumentation.NewCounter("number of create operations", "actuators", "machine", "Create"),
		deleteCount: instrumentation.NewCounter("number of delete operations", "actuators", "machine", "Delete"),
		updateCount: instrumentation.NewCounter("number of update operations", "actuators", "machine", "Update"),
		existsCount: instrumentation.NewCounter("number of update operations", "actuators", "machine", "Exists"),
	}

	return &Actuator{
		codec:          params.Codec,
		ec2:            params.EC2Service,
		elb:            params.ELBService,
		machinesGetter: params.MachinesGetter,
		i:              counters,
	}, nil
}

// Create creates a machine and is invoked by the machine controller.
func (a *Actuator) Create(cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	ctx := context.Background()
	ctx, span := trace.StartSpan(
		ctx, instrumentation.MethodName("actuators", "machine", "Create"),
	)
	defer span.End()

	stats.Record(ctx, a.i.createCount.M(1))

	rec, _ := events.NewStdObjRecorder(machine)

	rec.Event(events.Warning, "starting machine reconciliation", "object created")

	status, err := a.machineProviderStatus(machine)
	if err != nil {
		rec.Event(events.Warning, "retreiving machine provider status", "no provider status")
		return err
	}
	clusterStatus, err := a.clusterProviderStatus(cluster)
	if err != nil {
		rec.Event(events.Warning, "retreiving cluster provider status", "no provider status")
		return err
	}
	config, err := a.machineProviderConfig(machine)
	if err != nil {
		rec.Event(events.Warning, "retreiving machine provider configuration", "no machine configuration")
		return err
	}

	// Get a cluster api client for the namespace of the cluster.
	clusterClient := a.machinesGetter.Machines(machine.Namespace)

	storeUpdate := func() {
		if err := a.storeProviderStatus(clusterClient, machine, status); err != nil {
			rec.Event(events.Warning, "updating provider status", "update failed")
		}
		rec.Info(events.Normal, "updating provider status", "update succeeded")
	}

	defer storeUpdate()

	i, err := a.ec2.ReconcileInstance(ctx, machine, status, config, clusterStatus)

	if i != nil && i.ID != "" && i.State != "" {
		rec.Infof(events.Normal, "instance IDs received: %q", i.ID)
		state := string(i.State)
		status.InstanceID = &i.ID
		status.InstanceState = &state
		storeUpdate()
	}

	if err != nil {

		if ec2.IsFailedDependency(errors.Cause(err)) {
			rec.Event(events.Warning, "instance reconciliation", "dependencies not available yet")
			duration, _ := time.ParseDuration("60s")
			return &controllerError.RequeueAfterError{
				RequeueAfter: duration,
			}
		}
		rec.Event(events.Warning, "instance reconciliation", "could not create or find instance")
		return err
	}

	if machine.Annotations == nil {
		machine.Annotations = map[string]string{}
	}
	machine.Annotations["cluster-api-provider-aws"] = "true"

	if machine.ObjectMeta.Labels["set"] == "controlplane" {
		if err := a.elb.RegisterInstanceWithClassicELB(ctx, i.ID, elb.TagValueAPIServerRole); err != nil {
			rec.Event(events.Warning, "ELB registration", "not registered")
			return err
		}
	}
	rec.Event(events.Normal, "instance reconciliation", "machine created")
	return nil
}

// Delete deletes a machine and is invoked by the Machine Controller
func (a *Actuator) Delete(cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	ctx := context.Background()
	ctx, span := trace.StartSpan(
		ctx, instrumentation.MethodName("actuators", "machine", "Delete"),
	)
	defer span.End()

	stats.Record(ctx, a.i.deleteCount.M(1))

	rec, _ := events.NewStdObjRecorder(machine)
	rec.Event(events.Warning, "object deleted", "starting machine deletion")

	status, err := a.machineProviderStatus(machine)
	if err != nil {
		rec.Event(events.Warning, "retreiving machine provider status", "no provider status")
		return err
	}

	if status.InstanceID == nil {
		// Instance was never created
		return nil
	}

	instance, err := a.ec2.InstanceIfExists(ctx, *status.InstanceID)
	if err != nil {
		rec.Event(events.Warning, "retreiving EC2 instance for deletion", "failed to get instance")
		return err
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
		rec.Event(events.Normal, "EC2 instance terminated", "object deleted")
		return nil
	default:
		err = a.ec2.TerminateInstance(ctx, *status.InstanceID)
		if err != nil {
			rec.Event(events.Warning, "instance termination request", "failed to terminate instance")
			return err
		}
	}

	return nil
}

// Update updates a machine and is invoked by the Machine Controller.
// If the Update attempts to mutate any immutable state, the method will error
// and no updates will be performed.
func (a *Actuator) Update(cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	ctx := context.Background()
	ctx, span := trace.StartSpan(
		ctx, instrumentation.MethodName("actuators", "machine", "Update"),
	)
	defer span.End()

	stats.Record(ctx, a.i.updateCount.M(1))

	rec, _ := events.NewStdObjRecorder(machine)

	rec.Event(events.Warning, "starting machine update", "object updated")

	// Get the updated config. We're going to compare parts of this to the
	// current Tags and Security Groups that AWS is aware of.
	config, err := a.machineProviderConfig(machine)
	if err != nil {
		rec.Event(events.Warning, "retreiving machine provider config", "no provider config")
		return err
	}

	// Get the new status from the provided machine object.
	// We need this in case any of it has changed, and we also require the
	// instanceID for various AWS queries.
	status, err := a.machineProviderStatus(machine)
	if err != nil {
		rec.Event(events.Warning, "retreiving machine provider status", "no provider status")
		return err
	}

	// Get a cluster api client for the namespace of the cluster.
	clusterClient := a.machinesGetter.Machines(machine.Namespace)

	defer func() {
		if err := a.storeProviderStatus(clusterClient, machine, status); err != nil {
			rec.Event(events.Warning, "updating provider status", "update failed")
		}
	}()

	// Get the current instance description from AWS.
	instanceDescription, err := a.ec2.InstanceIfExists(ctx, *status.InstanceID)
	if err != nil {
		rec.Event(events.Warning, "instance reconciliation", "could not create or find instance")
		return err
	}

	// We can now compare the various AWS state to the state we were passed.
	// We will check immutable state first, in order to fail quickly before
	// moving on to state that we can mutate.
	// TODO: Implement immutable state check.

	// Ensure that the security groups are correct.
	securityGroupsChanged, err := a.ensureSecurityGroups(
		ctx,
		machine,
		status.InstanceID,
		config.AdditionalSecurityGroups,
		instanceDescription.SecurityGroups,
	)
	if err != nil {
		rec.Event(events.Warning, "security group reconciliation", "failed to ensure security groups")
		return err
	}

	// Ensure that the tags are correct.
	tagsChanged, err := a.ensureTags(ctx, machine, status.InstanceID, config.AdditionalTags)
	if err != nil {
		rec.Event(events.Warning, "tag reconciliation", "failed to ensure tags")
		return err
	}

	// We need to update the machine since annotations may have changed.
	if securityGroupsChanged || tagsChanged {
		err = a.updateMachine(ctx, machine)
		if err != nil {
			rec.Event(events.Warning, "updating machine object", "failed to update machine tags")
			return err
		}
	}

	return nil
}

// Exists test for the existence of a machine and is invoked by the Machine Controller
func (a *Actuator) Exists(cluster *clusterv1.Cluster, machine *clusterv1.Machine) (bool, error) {
	ctx := context.Background()
	ctx, span := trace.StartSpan(
		ctx, instrumentation.MethodName("actuators", "machine", "Exists"),
	)
	defer span.End()

	stats.Record(ctx, a.i.existsCount.M(1))

	rec, _ := events.NewStdObjRecorder(machine)
	rec.Info(events.Normal, "machine existence check", "Checking if machine exists")

	status, err := a.machineProviderStatus(machine)
	if err != nil {
		return false, err
	}

	// TODO worry about pointers. instance if exists returns *any* instance
	if status.InstanceID == nil {
		rec.Event(events.Warning, "checking for matching EC2 instance", "machine status Instance ID is nil")
		return false, nil
	}

	instance, err := a.ec2.InstanceIfExists(ctx, *status.InstanceID)
	if err != nil {
		rec.Event(events.Warning, "checking for matching EC2 instance", err.Error())
		return false, err
	}
	if instance == nil {
		rec.Event(events.Warning, "checking for matching EC2 instance", "no instance found")
		return false, nil
	}

	rec.Infof(events.Normal, "checking for matching EC2 instance", "found an instance: %#v", instance)

	switch instance.State {
	case v1alpha1.InstanceStateRunning, v1alpha1.InstanceStatePending:
		rec.Infof(events.Normal, "checking EC2 instance state", "instance in running state: %#v", instance)
		return true, nil
	default:
		rec.Event(events.Warning, "checking EC2 instance state", "no running EC2 instances found")
		return false, nil
	}
}

func (a *Actuator) updateMachine(ctx context.Context, machine *clusterv1.Machine) error {

	ctx, span := trace.StartSpan(ctx, "actuator.UpdateMachine")
	defer span.End()
	rec, _ := events.NewStdObjRecorder(machine)
	machinesClient := a.machinesGetter.Machines(machine.Namespace)

	_, err := machinesClient.Update(machine)
	if err != nil {
		rec.Error(events.Failure, "updating machine object", err.Error())
		return err
	}

	return nil
}
