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
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	machinev1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	machinecontroller "github.com/openshift/machine-api-operator/pkg/controller/machine"
	mapierrors "github.com/openshift/machine-api-operator/pkg/controller/machine"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apimachineryerrors "k8s.io/apimachinery/pkg/api/errors"
	errorutil "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1beta1"
	providerconfigv1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1beta1"
	awsclient "sigs.k8s.io/cluster-api-provider-aws/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	userDataSecretKey         = "userData"
	ec2InstanceIDNotFoundCode = "InvalidInstanceID.NotFound"
	requeueAfterSeconds       = 20
	requeueAfterFatalSeconds  = 180
)

// Actuator is the AWS-specific actuator for the Cluster API machine controller
type Actuator struct {
	awsClientBuilder awsclient.AwsClientBuilderFuncType
	client           client.Client
	config           *rest.Config

	codec         *providerconfigv1.AWSProviderConfigCodec
	eventRecorder record.EventRecorder
}

// ActuatorParams holds parameter information for Actuator
type ActuatorParams struct {
	Client           client.Client
	Config           *rest.Config
	AwsClientBuilder awsclient.AwsClientBuilderFuncType
	Codec            *providerconfigv1.AWSProviderConfigCodec
	EventRecorder    record.EventRecorder
}

// NewActuator returns a new AWS Actuator
func NewActuator(params ActuatorParams) (*Actuator, error) {
	actuator := &Actuator{
		client:           params.Client,
		config:           params.Config,
		awsClientBuilder: params.AwsClientBuilder,
		codec:            params.Codec,
		eventRecorder:    params.EventRecorder,
	}
	return actuator, nil
}

const (
	createEventAction = "Create"
	updateEventAction = "Update"
	deleteEventAction = "Delete"
	noEventAction     = ""
)

// Set corresponding event based on error. It also returns the original error
// for convenience, so callers can do "return handleMachineError(...)".
func (a *Actuator) handleMachineError(machine *machinev1.Machine, err error, eventAction string) error {
	if eventAction != noEventAction {
		a.eventRecorder.Eventf(machine, corev1.EventTypeWarning, "Failed"+eventAction, "%v", err)
	}

	klog.Errorf("%s: Machine error: %v", machine.Name, err)
	return err
}

// Create runs a new EC2 instance
func (a *Actuator) Create(context context.Context, machine *machinev1.Machine) error {
	klog.Infof("%s: creating machine", machine.Name)

	machineToBePatched := client.MergeFrom(machine.DeepCopy())

	instance, err := a.CreateMachine(machine)
	if err != nil {
		klog.Errorf("%s: error creating machine: %v", machine.Name, err)
		conditionFailed := conditionFailed()
		conditionFailed.Message = err.Error()
		updateConditionError := a.setMachineProviderConditions(machine, conditionFailed)
		if updateConditionError != nil {
			klog.Errorf("%s: error updating machine conditions: %v", machine.Name, updateConditionError)
		}
		patchErr := a.patchMachine(context, machine, machineToBePatched)
		if patchErr != nil {
			return a.handleMachineError(machine, errors.Wrap(err, "failed to patch machine status"), createEventAction)
		}
		return err
	}

	err = a.setProviderID(machine, instance)
	if err != nil {
		return a.handleMachineError(machine, errors.Wrap(err, "failed to update machine object with providerID"), createEventAction)
	}

	err = a.setMachineCloudProviderSpecifics(machine, instance)
	if err != nil {
		return a.handleMachineError(machine, errors.Wrap(err, "failed to set machine cloud provider specifics"), createEventAction)
	}

	err = a.setStatus(machine, instance, conditionSuccess())
	if err != nil {
		return a.handleMachineError(machine, errors.Wrap(err, "failed to set machine status"), createEventAction)

	}

	err = a.patchMachine(context, machine, machineToBePatched)
	if err != nil {
		return a.handleMachineError(machine, errors.Wrap(err, "failed to patch machine"), createEventAction)
	}

	return a.requeueIfInstancePending(machine.Name, instance)
}

func (a *Actuator) setMachineCloudProviderSpecifics(machine *machinev1.Machine, instance *ec2.Instance) error {
	if instance == nil {
		return nil
	}

	if machine.Labels == nil {
		machine.Labels = make(map[string]string)
	}

	if machine.Annotations == nil {
		machine.Annotations = make(map[string]string)
	}

	// Reaching to machine provider config since the region is not directly
	// providing by *ec2.Instance object
	machineProviderConfig, err := providerConfigFromMachine(machine, a.codec)
	if err != nil {
		return a.handleMachineError(machine, errors.Wrapf(err, "error decoding MachineProviderConfig"), noEventAction)
	}

	machine.Labels[machinecontroller.MachineRegionLabelName] = machineProviderConfig.Placement.Region

	if instance.Placement != nil {
		machine.Labels[machinecontroller.MachineAZLabelName] = aws.StringValue(instance.Placement.AvailabilityZone)
	}

	if instance.InstanceType != nil {
		machine.Labels[machinecontroller.MachineInstanceTypeLabelName] = aws.StringValue(instance.InstanceType)
	}

	if instance.State != nil && instance.State.Name != nil {
		machine.Annotations[machinecontroller.MachineInstanceStateAnnotationName] = aws.StringValue(instance.State.Name)
	}

	return nil
}

// setProviderID adds providerID in the machine spec
func (a *Actuator) setProviderID(machine *machinev1.Machine, instance *ec2.Instance) error {
	existingProviderID := machine.Spec.ProviderID
	if instance == nil {
		return nil
	}
	availabilityZone := ""
	if instance.Placement != nil {
		availabilityZone = aws.StringValue(instance.Placement.AvailabilityZone)
	}
	providerID := fmt.Sprintf("aws:///%s/%s", availabilityZone, aws.StringValue(instance.InstanceId))

	if existingProviderID != nil && *existingProviderID == providerID {
		klog.Infof("%s: ProviderID already set in the machine Spec with value:%s", machine.Name, *existingProviderID)
		return nil
	}
	machine.Spec.ProviderID = &providerID
	klog.Infof("%s: ProviderID set at machine spec: %s", machine.Name, providerID)
	return nil
}

func (a *Actuator) setMachineStatus(machine *machinev1.Machine, awsStatus *providerconfigv1.AWSMachineProviderStatus, networkAddresses []corev1.NodeAddress) error {
	awsStatusRaw, err := a.codec.EncodeProviderStatus(awsStatus)
	if err != nil {
		klog.Errorf("%s: error encoding AWS provider status: %v", machine.Name, err)
		return err
	}

	machine.Status.ProviderStatus = awsStatusRaw
	if networkAddresses != nil {
		machine.Status.Addresses = networkAddresses
	}

	return nil
}

// updateMachineProviderConditions updates conditions set within machine provider status.
func (a *Actuator) setMachineProviderConditions(machine *machinev1.Machine, condition providerconfigv1.AWSMachineProviderCondition) error {
	klog.Infof("%s: updating machine conditions", machine.Name)

	awsStatus := &providerconfigv1.AWSMachineProviderStatus{}
	if err := a.codec.DecodeProviderStatus(machine.Status.ProviderStatus, awsStatus); err != nil {
		klog.Errorf("%s: error decoding machine provider status: %v", machine.Name, err)
		return err
	}

	awsStatus.Conditions = setAWSMachineProviderCondition(condition, awsStatus.Conditions)

	if err := a.setMachineStatus(machine, awsStatus, nil); err != nil {
		return err
	}

	return nil
}

// CreateMachine starts a new AWS instance as described by the cluster and machine resources
func (a *Actuator) CreateMachine(machine *machinev1.Machine) (*ec2.Instance, error) {
	machineProviderConfig, err := providerConfigFromMachine(machine, a.codec)
	if err != nil {
		return nil, a.handleMachineError(machine, mapierrors.InvalidMachineConfiguration("error decoding MachineProviderConfig: %v", err), createEventAction)
	}

	credentialsSecretName := ""
	if machineProviderConfig.CredentialsSecret != nil {
		credentialsSecretName = machineProviderConfig.CredentialsSecret.Name
	}
	awsClient, err := a.awsClientBuilder(a.client, credentialsSecretName, machine.Namespace, machineProviderConfig.Placement.Region)
	if err != nil {
		return nil, a.handleMachineError(machine, err, createEventAction)
	}

	// We explicitly do NOT want to remove stopped masters.
	isMaster, err := a.isMaster(machine)
	// Unable to determine if a machine is a master machine.
	// Yet, it's only used to delete stopped machines that are not masters.
	// So we can safely continue to create a new machine since in the worst case
	// we just don't delete any stopped machine.
	if err != nil {
		klog.Errorf("%s: Error determining if machine is master: %v", machine.Name, err)
	} else {
		if !isMaster {
			// Prevent having a lot of stopped nodes sitting around.
			err = removeStoppedMachine(machine, awsClient)
			if err != nil {
				return nil, a.handleMachineError(machine, errors.Wrapf(err, "unable to remove stopped machines"), createEventAction)
			}
		}
	}

	userData, err := a.getUserData(machine, machineProviderConfig)
	if err != nil {
		return nil, err
	}

	instance, err := launchInstance(machine, machineProviderConfig, userData, awsClient)
	if err != nil {
		return nil, a.handleMachineError(machine, err, createEventAction)
	}

	err = a.updateLoadBalancers(awsClient, machineProviderConfig, instance, machine.Name)
	if err != nil {
		return nil, a.handleMachineError(machine, err, createEventAction)
	}

	a.eventRecorder.Eventf(machine, corev1.EventTypeNormal, "Created", "Created Machine %v", machine.Name)
	return instance, nil
}

func (a *Actuator) getUserData(machine *machinev1.Machine, machineProviderConfig *providerconfigv1.AWSMachineProviderConfig) ([]byte, error) {
	if machineProviderConfig.UserDataSecret == nil {
		return []byte{}, nil
	}

	var userDataSecret corev1.Secret
	err := a.client.Get(context.Background(), client.ObjectKey{Namespace: machine.Namespace, Name: machineProviderConfig.UserDataSecret.Name}, &userDataSecret)
	if err != nil {
		if apimachineryerrors.IsNotFound(err) {
			return nil, a.handleMachineError(machine, mapierrors.InvalidMachineConfiguration("user data secret %s: %v not found", machineProviderConfig.UserDataSecret.Name, err), createEventAction)
		}
		return nil, a.handleMachineError(machine, mapierrors.CreateMachine("error getting user data secret %s: %v", machineProviderConfig.UserDataSecret.Name, err), createEventAction)
	}
	userData, exists := userDataSecret.Data[userDataSecretKey]
	if !exists {
		return nil, a.handleMachineError(machine, mapierrors.InvalidMachineConfiguration("%s: Secret %v/%v does not have %q field set. Thus, no user data applied when creating an instance.", machine.Name, machine.Namespace, machineProviderConfig.UserDataSecret.Name, userDataSecretKey), createEventAction)
	}
	return userData, nil
}

// Delete deletes a machine and updates its finalizer
func (a *Actuator) Delete(context context.Context, machine *machinev1.Machine) error {
	klog.Infof("%s: deleting machine", machine.Name)
	if err := a.DeleteMachine(machine); err != nil {
		klog.Errorf("%s: error deleting machine: %v", machine.Name, err)
		return err
	}
	return nil
}

// DeleteMachine deletes an AWS instance
func (a *Actuator) DeleteMachine(machine *machinev1.Machine) error {
	machineProviderConfig, err := providerConfigFromMachine(machine, a.codec)
	if err != nil {
		return a.handleMachineError(machine, mapierrors.InvalidMachineConfiguration("error decoding MachineProviderConfig: %v", err), deleteEventAction)
	}

	region := machineProviderConfig.Placement.Region
	credentialsSecretName := ""
	if machineProviderConfig.CredentialsSecret != nil {
		credentialsSecretName = machineProviderConfig.CredentialsSecret.Name
	}
	client, err := a.awsClientBuilder(a.client, credentialsSecretName, machine.Namespace, region)
	if err != nil {
		return a.handleMachineError(machine, err, deleteEventAction)
	}

	// Get all instances not terminated.
	existingInstances, err := a.getMachineInstances(machine)
	if err != nil {
		klog.Errorf("%s: error getting existing instances: %v", machine.Name, err)
		return err
	}
	existingLen := len(existingInstances)
	klog.Infof("%s: found %d existing instances for machine", machine.Name, existingLen)
	if existingLen == 0 {
		klog.Warningf("%s: no instances found to delete for machine", machine.Name)
		return nil
	}

	terminatingInstances, err := terminateInstances(client, existingInstances)
	if err != nil {
		return a.handleMachineError(machine, mapierrors.DeleteMachine(err.Error()), noEventAction)
	}

	if len(terminatingInstances) == 1 {
		if terminatingInstances[0] != nil && terminatingInstances[0].CurrentState != nil && terminatingInstances[0].CurrentState.Name != nil {
			machineCopy := machine.DeepCopy()
			machineCopy.Annotations[machinecontroller.MachineInstanceStateAnnotationName] = aws.StringValue(terminatingInstances[0].CurrentState.Name)
			a.client.Update(context.Background(), machineCopy)
		}
	}

	a.eventRecorder.Eventf(machine, corev1.EventTypeNormal, "Deleted", "Deleted machine %v", machine.Name)

	return nil
}

// Update attempts to sync machine state with an existing instance. Today this just updates status
// for details that may have changed. (IPs and hostnames) We do not currently support making any
// changes to actual machines in AWS. Instead these will be replaced via MachineDeployments.
func (a *Actuator) Update(context context.Context, machine *machinev1.Machine) error {
	klog.Infof("%s: updating machine", machine.Name)

	machineToBePatched := client.MergeFrom(machine.DeepCopy())

	machineProviderConfig, err := providerConfigFromMachine(machine, a.codec)
	if err != nil {
		return a.handleMachineError(machine, mapierrors.InvalidMachineConfiguration("error decoding MachineProviderConfig: %v", err), updateEventAction)
	}

	region := machineProviderConfig.Placement.Region
	klog.Infof("%s: obtaining EC2 client for region", machine.Name)
	credentialsSecretName := ""
	if machineProviderConfig.CredentialsSecret != nil {
		credentialsSecretName = machineProviderConfig.CredentialsSecret.Name
	}
	client, err := a.awsClientBuilder(a.client, credentialsSecretName, machine.Namespace, region)
	if err != nil {
		return a.handleMachineError(machine, err, updateEventAction)
	}
	// Get all instances not terminated.
	existingInstances, err := a.getMachineInstances(machine)
	if err != nil {
		klog.Errorf("%s: error getting existing instances: %v", machine.Name, err)
		return err
	}
	existingLen := len(existingInstances)
	klog.Infof("%s: found %d existing instances for machine", machine.Name, existingLen)

	// Parent controller should prevent this from ever happening by calling Exists and then Create,
	// but instance could be deleted between the two calls.
	if existingLen == 0 {
		if machine.Spec.ProviderID != nil && *machine.Spec.ProviderID != "" && (machine.Status.LastUpdated == nil || machine.Status.LastUpdated.Add(requeueAfterSeconds*time.Second).After(time.Now())) {
			klog.Infof("%s: Possible eventual-consistency discrepancy; returning an error to requeue", machine.Name)
			return &machinecontroller.RequeueAfterError{RequeueAfter: requeueAfterSeconds * time.Second}
		}

		klog.Warningf("%s: attempted to update machine but no instances found", machine.Name)

		a.handleMachineError(machine, mapierrors.UpdateMachine("no instance found, reason unknown"), updateEventAction)

		// Update status to clear out machine details.
		if err := a.setStatus(machine, nil, conditionSuccess()); err != nil {
			return err
		}
		// This is an unrecoverable error condition.  We should delay to
		// minimize unnecessary API calls.
		return &machinecontroller.RequeueAfterError{RequeueAfter: requeueAfterFatalSeconds * time.Second}
	}
	sortInstances(existingInstances)
	runningInstances := getRunningFromInstances(existingInstances)
	runningLen := len(runningInstances)
	var newestInstance *ec2.Instance
	if runningLen > 0 {
		// It would be very unusual to have more than one here, but it is
		// possible if someone manually provisions a machine with same tag name.
		klog.Infof("%s: found %d running instances for machine", machine.Name, runningLen)
		newestInstance = runningInstances[0]

		err = a.updateLoadBalancers(client, machineProviderConfig, newestInstance, machine.Name)
		if err != nil {
			a.handleMachineError(machine, mapierrors.CreateMachine("Error updating load balancers: %v", err), updateEventAction)
			return err
		}
	} else {
		// Didn't find any running instances, just newest existing one.
		// In most cases, there should only be one existing Instance.
		newestInstance = existingInstances[0]
	}

	a.eventRecorder.Eventf(machine, corev1.EventTypeNormal, "Updated", "Updated machine %v", machine.Name)

	err = a.setMachineCloudProviderSpecifics(machine, newestInstance)
	if err != nil {
		return a.handleMachineError(machine, errors.Wrap(err, "failed to set machine cloud provider specifics"), updateEventAction)
	}

	err = a.setProviderID(machine, newestInstance)
	if err != nil {
		return a.handleMachineError(machine, errors.Wrap(err, "failed to update machine object with providerID"), updateEventAction)
	}

	// We do not support making changes to pre-existing instances, just update status.
	err = a.setStatus(machine, newestInstance, conditionSuccess())
	if err != nil {
		return a.handleMachineError(machine, errors.Wrap(err, "failed to set machine status"), updateEventAction)
	}

	err = a.patchMachine(context, machine, machineToBePatched)
	if err != nil {
		return a.handleMachineError(machine, errors.Wrap(err, "failed to patch machine"), updateEventAction)
	}

	return a.requeueIfInstancePending(machine.Name, newestInstance)
}

// Exists determines if the given machine currently exists.
// A machine which is not terminated is considered as existing.
func (a *Actuator) Exists(context context.Context, machine *machinev1.Machine) (bool, error) {
	instance, err := a.Describe(machine)
	return instance != nil, err
}

// Describe provides information about machine's instance(s)
func (a *Actuator) Describe(machine *machinev1.Machine) (*ec2.Instance, error) {
	klog.Infof("%s: Checking if machine exists", machine.Name)

	instances, err := a.getMachineInstances(machine)
	if err != nil {
		klog.Errorf("%s: Error getting existing instances: %v", machine.Name, err)
		return nil, err
	}
	if len(instances) == 0 {
		if machine.Spec.ProviderID != nil && *machine.Spec.ProviderID != "" && (machine.Status.LastUpdated == nil || machine.Status.LastUpdated.Add(requeueAfterSeconds*time.Second).After(time.Now())) {
			klog.Infof("%s: Possible eventual-consistency discrepancy; returning an error to requeue", machine.Name)
			return nil, &machinecontroller.RequeueAfterError{RequeueAfter: requeueAfterSeconds * time.Second}
		}

		klog.Infof("%s: Instance does not exist", machine.Name)
		return nil, nil
	}

	return instances[0], nil
}

func (a *Actuator) getMachineInstances(machine *machinev1.Machine) ([]*ec2.Instance, error) {
	machineProviderConfig, err := providerConfigFromMachine(machine, a.codec)
	if err != nil {
		klog.Errorf("%s: Error decoding MachineProviderConfig: %v", machine.Name, err)
		return nil, err
	}

	region := machineProviderConfig.Placement.Region
	credentialsSecretName := ""
	if machineProviderConfig.CredentialsSecret != nil {
		credentialsSecretName = machineProviderConfig.CredentialsSecret.Name
	}
	client, err := a.awsClientBuilder(a.client, credentialsSecretName, machine.Namespace, region)
	if err != nil {
		klog.Errorf("%s: Error getting EC2 client: %v", machine.Name, err)
		return nil, err
	}

	status := &providerconfigv1.AWSMachineProviderStatus{}
	err = a.codec.DecodeProviderStatus(machine.Status.ProviderStatus, status)

	// If the status was decoded successfully, and there is a non-empty instance
	// ID, search using that, otherwise fallback to filtering based on tags.
	if err == nil && status.InstanceID != nil && *status.InstanceID != "" {
		i, err := getExistingInstanceByID(*status.InstanceID, client)
		if err != nil {
			klog.Warningf("%s: Failed to find existing instance by id %s: %v",
				machine.Name, *status.InstanceID, err)
		} else {
			klog.Infof("%s: Found instance by id: %s", machine.Name, *status.InstanceID)
			return []*ec2.Instance{i}, nil
		}
	}

	return getExistingInstances(machine, client)
}

// updateLoadBalancers adds a given machine instance to the load balancers specified in its provider config
func (a *Actuator) updateLoadBalancers(client awsclient.Client, providerConfig *providerconfigv1.AWSMachineProviderConfig, instance *ec2.Instance, machineName string) error {
	if len(providerConfig.LoadBalancers) == 0 {
		klog.V(4).Infof("%s: Instance %q has no load balancers configured. Skipping", machineName, *instance.InstanceId)
		return nil
	}
	errs := []error{}
	classicLoadBalancerNames := []string{}
	networkLoadBalancerNames := []string{}
	for _, loadBalancerRef := range providerConfig.LoadBalancers {
		switch loadBalancerRef.Type {
		case providerconfigv1.NetworkLoadBalancerType:
			networkLoadBalancerNames = append(networkLoadBalancerNames, loadBalancerRef.Name)
		case providerconfigv1.ClassicLoadBalancerType:
			classicLoadBalancerNames = append(classicLoadBalancerNames, loadBalancerRef.Name)
		}
	}

	var err error
	if len(classicLoadBalancerNames) > 0 {
		err := registerWithClassicLoadBalancers(client, classicLoadBalancerNames, instance)
		if err != nil {
			klog.Errorf("%s: Failed to register classic load balancers: %v", machineName, err)
			errs = append(errs, err)
		}
	}
	if len(networkLoadBalancerNames) > 0 {
		err = registerWithNetworkLoadBalancers(client, networkLoadBalancerNames, instance)
		if err != nil {
			klog.Errorf("%s: Failed to register network load balancers: %v", machineName, err)
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errorutil.NewAggregate(errs)
	}
	return nil
}

// setStatus calculates the new machine status, checks if anything has changed, and updates if so.
func (a *Actuator) setStatus(machine *machinev1.Machine, instance *ec2.Instance, condition v1beta1.AWSMachineProviderCondition) error {
	klog.Infof("%s: Updating status", machine.Name)

	// Starting with a fresh status as we assume full control of it here.
	awsStatus := &providerconfigv1.AWSMachineProviderStatus{}
	if err := a.codec.DecodeProviderStatus(machine.Status.ProviderStatus, awsStatus); err != nil {
		klog.Errorf("%s: Error decoding machine provider status: %v", machine.Name, err)
		return err
	}

	// Save this, we need to check if it changed later.
	networkAddresses := []corev1.NodeAddress{}

	// Instance may have existed but been deleted outside our control, clear it's status if so:
	if instance == nil {
		awsStatus.InstanceID = nil
		awsStatus.InstanceState = nil
	} else {
		awsStatus.InstanceID = instance.InstanceId
		awsStatus.InstanceState = instance.State.Name

		addresses, err := extractNodeAddresses(instance)
		if err != nil {
			klog.Errorf("%s: Error extracting instance IP addresses: %v", machine.Name, err)
			return err
		}

		networkAddresses = append(networkAddresses, addresses...)
	}
	klog.Infof("%s: finished calculating AWS status", machine.Name)

	awsStatus.Conditions = setAWSMachineProviderCondition(condition, awsStatus.Conditions)
	if err := a.setMachineStatus(machine, awsStatus, networkAddresses); err != nil {
		return err
	}

	return nil
}

func getClusterID(machine *machinev1.Machine) (string, bool) {
	clusterID, ok := machine.Labels[providerconfigv1.ClusterIDLabel]
	// NOTE: This block can be removed after the label renaming transition to machine.openshift.io
	if !ok {
		clusterID, ok = machine.Labels["sigs.k8s.io/cluster-api-cluster"]
	}
	return clusterID, ok
}

func (a *Actuator) patchMachine(ctx context.Context, machine *machinev1.Machine, machineToBePatched client.Patch) error {
	statusCopy := *machine.Status.DeepCopy()

	// Patch machine
	if err := a.client.Patch(ctx, machine, machineToBePatched); err != nil {
		klog.Errorf("Failed to update machine %q: %v", machine.GetName(), err)
		return err
	}

	machine.Status = statusCopy

	//Patch status
	if err := a.client.Status().Patch(ctx, machine, machineToBePatched); err != nil {
		klog.Errorf("Failed to update machine status %q: %v", machine.GetName(), err)
		return err
	}
	return nil
}

func (a *Actuator) requeueIfInstancePending(machineName string, instance *ec2.Instance) error {
	// If machine state is still pending, we will return an error to keep the controllers
	// attempting to update status until it hits a more permanent state. This will ensure
	// we get a public IP populated more quickly.
	if instance.State != nil && *instance.State.Name == ec2.InstanceStateNamePending {
		klog.Infof("%s: Instance state still pending, returning an error to requeue", machineName)
		return &machinecontroller.RequeueAfterError{RequeueAfter: requeueAfterSeconds * time.Second}
	}

	return nil
}
