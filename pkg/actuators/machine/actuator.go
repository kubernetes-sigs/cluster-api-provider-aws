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
	"github.com/golang/glog"
	clusterv1 "github.com/openshift/cluster-api/pkg/apis/cluster/v1alpha1"
	machinev1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	clustererror "github.com/openshift/cluster-api/pkg/controller/error"
	machinecontroller "github.com/openshift/cluster-api/pkg/controller/machine"
	mapierrors "github.com/openshift/cluster-api/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apimachineryerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	errorutil "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog"
	providerconfigv1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1beta1"
	awsclient "sigs.k8s.io/cluster-api-provider-aws/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	userDataSecretKey         = "userData"
	ec2InstanceIDNotFoundCode = "InvalidInstanceID.NotFound"
	requeueAfterSeconds       = 20
	requeueAfterFatalSeconds  = 180

	// MachineCreationSucceeded indicates success for machine creation
	MachineCreationSucceeded = "MachineCreationSucceeded"

	// MachineCreationFailed indicates that machine creation failed
	MachineCreationFailed = "MachineCreationFailed"
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

	glog.Errorf("%s: Machine error: %v", machine.Name, err)
	return err
}

// Create runs a new EC2 instance
func (a *Actuator) Create(context context.Context, cluster *clusterv1.Cluster, machine *machinev1.Machine) error {
	glog.Infof("%s: creating machine", machine.Name)
	instance, err := a.CreateMachine(cluster, machine)
	if err != nil {
		glog.Errorf("%s: error creating machine: %v", machine.Name, err)
		updateConditionError := a.updateMachineProviderConditions(machine, providerconfigv1.MachineCreation, MachineCreationFailed, err.Error())
		if updateConditionError != nil {
			glog.Errorf("%s: error updating machine conditions: %v", machine.Name, updateConditionError)
		}
		return err
	}
	updatedMachine, err := a.updateProviderID(machine, instance)
	if err != nil {
		return fmt.Errorf("%s: failed to update machine object with providerID: %v", machine.Name, err)
	}
	machine = updatedMachine

	updatedMachine, err = a.setMachineCloudProviderSpecifics(machine, instance)
	if err != nil {
		return fmt.Errorf("%s: failed to set machine cloud provider specifics: %v", machine.Name, err)
	}

	return a.updateStatus(updatedMachine, instance)
}

func (a *Actuator) setMachineCloudProviderSpecifics(machine *machinev1.Machine, instance *ec2.Instance) (*machinev1.Machine, error) {
	if instance == nil {
		return machine, nil
	}
	machineCopy := machine.DeepCopy()

	if machineCopy.Labels == nil {
		machineCopy.Labels = make(map[string]string)
	}

	if machineCopy.Annotations == nil {
		machineCopy.Annotations = make(map[string]string)
	}

	// Reaching to machine provider config since the region is not directly
	// providing by *ec2.Instance object
	machineProviderConfig, err := providerConfigFromMachine(machine, a.codec)
	if err != nil {
		return nil, fmt.Errorf("error decoding MachineProviderConfig: %v", err)
	}

	machineCopy.Labels[machinecontroller.MachineRegionLabelName] = machineProviderConfig.Placement.Region

	if instance.Placement != nil {
		machineCopy.Labels[machinecontroller.MachineAZLabelName] = aws.StringValue(instance.Placement.AvailabilityZone)
	}

	if instance.InstanceType != nil {
		machineCopy.Labels[machinecontroller.MachineInstanceTypeLabelName] = aws.StringValue(instance.InstanceType)
	}

	if instance.State != nil && instance.State.Name != nil {
		machineCopy.Annotations[machinecontroller.MachineInstanceStateAnnotationName] = aws.StringValue(instance.State.Name)
	}

	if err := a.client.Update(context.Background(), machineCopy); err != nil {
		return nil, fmt.Errorf("%s: error updating machine spec: %v", machine.Name, err)
	}

	return machineCopy, nil
}

// updateProviderID adds providerID in the machine spec
func (a *Actuator) updateProviderID(machine *machinev1.Machine, instance *ec2.Instance) (*machinev1.Machine, error) {
	existingProviderID := machine.Spec.ProviderID
	machineCopy := machine.DeepCopy()
	if instance != nil {
		availabilityZone := ""
		if instance.Placement != nil {
			availabilityZone = aws.StringValue(instance.Placement.AvailabilityZone)
		}
		providerID := fmt.Sprintf("aws:///%s/%s", availabilityZone, aws.StringValue(instance.InstanceId))

		if existingProviderID != nil && *existingProviderID == providerID {
			glog.Infof("%s: ProviderID already set in the machine Spec with value:%s", machine.Name, *existingProviderID)
			return machine, nil
		}
		machineCopy.Spec.ProviderID = &providerID
		if err := a.client.Update(context.Background(), machineCopy); err != nil {
			return nil, fmt.Errorf("%s: error updating machine spec ProviderID: %v", machine.Name, err)
		}
		glog.Infof("%s: ProviderID updated at machine spec: %s", machine.Name, providerID)
	} else {
		machineCopy.Spec.ProviderID = nil
		if err := a.client.Update(context.Background(), machineCopy); err != nil {
			return nil, fmt.Errorf("%s: error updating ProviderID in machine spec: %v", machine.Name, err)
		}
		glog.Infof("%s: No instance found so clearing ProviderID field in the machine spec", machine.Name)
	}
	return machineCopy, nil
}

func (a *Actuator) updateMachineStatus(machine *machinev1.Machine, awsStatus *providerconfigv1.AWSMachineProviderStatus, networkAddresses []corev1.NodeAddress) error {
	awsStatusRaw, err := a.codec.EncodeProviderStatus(awsStatus)
	if err != nil {
		glog.Errorf("%s: error encoding AWS provider status: %v", machine.Name, err)
		return err
	}

	machineCopy := machine.DeepCopy()
	machineCopy.Status.ProviderStatus = awsStatusRaw
	if networkAddresses != nil {
		machineCopy.Status.Addresses = networkAddresses
	}

	oldAWSStatus := &providerconfigv1.AWSMachineProviderStatus{}
	if err := a.codec.DecodeProviderStatus(machine.Status.ProviderStatus, oldAWSStatus); err != nil {
		glog.Errorf("%s: error updating machine status: %v", machine.Name, err)
		return err
	}

	// TODO(vikasc): Revisit to compare complete machine status objects
	if !equality.Semantic.DeepEqual(awsStatus, oldAWSStatus) || !equality.Semantic.DeepEqual(machine.Status.Addresses, machineCopy.Status.Addresses) {
		glog.Infof("%s: machine status has changed, updating", machine.Name)
		time := metav1.Now()
		machineCopy.Status.LastUpdated = &time

		if err := a.client.Status().Update(context.Background(), machineCopy); err != nil {
			glog.Errorf("%s: error updating machine status: %v", machine.Name, err)
			return err
		}
	} else {
		glog.Infof("%s: status unchanged", machine.Name)
	}

	return nil
}

// updateMachineProviderConditions updates conditions set within machine provider status.
func (a *Actuator) updateMachineProviderConditions(machine *machinev1.Machine, conditionType providerconfigv1.AWSMachineProviderConditionType, reason string, msg string) error {

	glog.Infof("%s: updating machine conditions", machine.Name)

	awsStatus := &providerconfigv1.AWSMachineProviderStatus{}
	if err := a.codec.DecodeProviderStatus(machine.Status.ProviderStatus, awsStatus); err != nil {
		glog.Errorf("%s: error decoding machine provider status: %v", machine.Name, err)
		return err
	}

	awsStatus.Conditions = setAWSMachineProviderCondition(awsStatus.Conditions, conditionType, corev1.ConditionTrue, reason, msg, updateConditionIfReasonOrMessageChange)

	if err := a.updateMachineStatus(machine, awsStatus, nil); err != nil {
		return err
	}

	return nil
}

// CreateMachine starts a new AWS instance as described by the cluster and machine resources
func (a *Actuator) CreateMachine(cluster *clusterv1.Cluster, machine *machinev1.Machine) (*ec2.Instance, error) {
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
				errMsg := fmt.Sprintf("%s: unable to remove stopped machines: %v", machine.Name, err)
				glog.Errorf(errMsg)
				return nil, fmt.Errorf(errMsg)
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
func (a *Actuator) Delete(context context.Context, cluster *clusterv1.Cluster, machine *machinev1.Machine) error {
	glog.Infof("%s: deleting machine", machine.Name)
	if err := a.DeleteMachine(cluster, machine); err != nil {
		glog.Errorf("%s: error deleting machine: %v", machine.Name, err)
		return err
	}
	return nil
}

type glogLogger struct{}

func (gl *glogLogger) Log(v ...interface{}) {
	glog.Info(v...)
}

func (gl *glogLogger) Logf(format string, v ...interface{}) {
	glog.Infof(format, v...)
}

// DeleteMachine deletes an AWS instance
func (a *Actuator) DeleteMachine(cluster *clusterv1.Cluster, machine *machinev1.Machine) error {
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

	instances, err := getRunningInstances(machine, client)
	if err != nil {
		glog.Errorf("%s: error getting running instances: %v", machine.Name, err)
		return err
	}
	if len(instances) == 0 {
		glog.Warningf("%s: no instances found to delete for machine", machine.Name)
		return nil
	}

	terminatingInstances, err := terminateInstances(client, instances)
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
func (a *Actuator) Update(context context.Context, cluster *clusterv1.Cluster, machine *machinev1.Machine) error {
	glog.Infof("%s: updating machine", machine.Name)

	machineProviderConfig, err := providerConfigFromMachine(machine, a.codec)
	if err != nil {
		return a.handleMachineError(machine, mapierrors.InvalidMachineConfiguration("error decoding MachineProviderConfig: %v", err), updateEventAction)
	}

	region := machineProviderConfig.Placement.Region
	glog.Infof("%s: obtaining EC2 client for region", machine.Name)
	credentialsSecretName := ""
	if machineProviderConfig.CredentialsSecret != nil {
		credentialsSecretName = machineProviderConfig.CredentialsSecret.Name
	}
	client, err := a.awsClientBuilder(a.client, credentialsSecretName, machine.Namespace, region)
	if err != nil {
		return a.handleMachineError(machine, err, updateEventAction)
	}
	// Get all instances not terminated.
	existingInstances, err := a.getMachineInstances(cluster, machine)
	if err != nil {
		glog.Errorf("%s: error getting existing instances: %v", machine.Name, err)
		return err
	}
	existingLen := len(existingInstances)
	glog.Infof("%s: found %d existing instances for machine", machine.Name, existingLen)

	// Parent controller should prevent this from ever happening by calling Exists and then Create,
	// but instance could be deleted between the two calls.
	if existingLen == 0 {
		if machine.Spec.ProviderID != nil && (machine.Status.LastUpdated == nil || machine.Status.LastUpdated.Add(requeueAfterSeconds*time.Second).After(time.Now())) {
			glog.Infof("%s: Possible eventual-consistency discrepancy; returning an error to requeue", machine.Name)
			return &clustererror.RequeueAfterError{RequeueAfter: requeueAfterSeconds * time.Second}
		}

		glog.Warningf("%s: attempted to update machine but no instances found", machine.Name)

		a.handleMachineError(machine, mapierrors.UpdateMachine("no instance found, reason unknown"), updateEventAction)

		// Update status to clear out machine details.
		if err := a.updateStatus(machine, nil); err != nil {
			return err
		}
		// This is an unrecoverable error condition.  We should delay to
		// minimize unnecessary API calls.
		return &clustererror.RequeueAfterError{RequeueAfter: requeueAfterFatalSeconds * time.Second}
	}
	sortInstances(existingInstances)
	runningInstances := getRunningFromInstances(existingInstances)
	runningLen := len(runningInstances)
	var newestInstance *ec2.Instance
	if runningLen > 0 {
		// It would be very unusual to have more than one here, but it is
		// possible if someone manually provisions a machine with same tag name.
		glog.Infof("%s: found %d running instances for machine", machine.Name, runningLen)
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

	modMachine, err := a.setMachineCloudProviderSpecifics(machine, newestInstance)
	if err != nil {
		return fmt.Errorf("%s: failed to set machine cloud provider specifics: %v", machine.Name, err)
	}
	machine = modMachine

	updatedMachine, err := a.updateProviderID(machine, newestInstance)
	if err != nil {
		return fmt.Errorf("%s: failed to update machine object with providerID: %v", machine.Name, err)
	}
	machine = updatedMachine

	// We do not support making changes to pre-existing instances, just update status.
	return a.updateStatus(machine, newestInstance)
}

// Exists determines if the given machine currently exists.
// A machine which is not terminated is considered as existing.
func (a *Actuator) Exists(context context.Context, cluster *clusterv1.Cluster, machine *machinev1.Machine) (bool, error) {
	instance, err := a.Describe(cluster, machine)
	return instance != nil, err
}

// Describe provides information about machine's instance(s)
func (a *Actuator) Describe(cluster *clusterv1.Cluster, machine *machinev1.Machine) (*ec2.Instance, error) {
	glog.Infof("%s: Checking if machine exists", machine.Name)

	instances, err := a.getMachineInstances(cluster, machine)
	if err != nil {
		glog.Errorf("%s: Error getting existing instances: %v", machine.Name, err)
		return nil, err
	}
	if len(instances) == 0 {
		if machine.Spec.ProviderID != nil && (machine.Status.LastUpdated == nil || machine.Status.LastUpdated.Add(requeueAfterSeconds*time.Second).After(time.Now())) {
			glog.Infof("%s: Possible eventual-consistency discrepancy; returning an error to requeue", machine.Name)
			return nil, &clustererror.RequeueAfterError{RequeueAfter: requeueAfterSeconds * time.Second}
		}

		glog.Infof("%s: Instance does not exist", machine.Name)
		return nil, nil
	}

	return instances[0], nil
}

func (a *Actuator) getMachineInstances(cluster *clusterv1.Cluster, machine *machinev1.Machine) ([]*ec2.Instance, error) {
	machineProviderConfig, err := providerConfigFromMachine(machine, a.codec)
	if err != nil {
		glog.Errorf("%s: Error decoding MachineProviderConfig: %v", machine.Name, err)
		return nil, err
	}

	region := machineProviderConfig.Placement.Region
	credentialsSecretName := ""
	if machineProviderConfig.CredentialsSecret != nil {
		credentialsSecretName = machineProviderConfig.CredentialsSecret.Name
	}
	client, err := a.awsClientBuilder(a.client, credentialsSecretName, machine.Namespace, region)
	if err != nil {
		glog.Errorf("%s: Error getting EC2 client: %v", machine.Name, err)
		return nil, err
	}

	status := &providerconfigv1.AWSMachineProviderStatus{}
	err = a.codec.DecodeProviderStatus(machine.Status.ProviderStatus, status)

	// If the status was decoded successfully, and there is a non-empty instance
	// ID, search using that, otherwise fallback to filtering based on tags.
	if err == nil && status.InstanceID != nil && *status.InstanceID != "" {
		i, err := getExistingInstanceByID(*status.InstanceID, client)
		if err != nil {
			glog.Warningf("%s: Failed to find existing instance by id %s: %v",
				machine.Name, *status.InstanceID, err)
		} else {
			glog.Infof("%s: Found instance by id: %s", machine.Name, *status.InstanceID)
			return []*ec2.Instance{i}, nil
		}
	}

	return getExistingInstances(machine, client)
}

// updateLoadBalancers adds a given machine instance to the load balancers specified in its provider config
func (a *Actuator) updateLoadBalancers(client awsclient.Client, providerConfig *providerconfigv1.AWSMachineProviderConfig, instance *ec2.Instance, machineName string) error {
	if len(providerConfig.LoadBalancers) == 0 {
		glog.V(4).Infof("%s: Instance %q has no load balancers configured. Skipping", machineName, *instance.InstanceId)
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
			glog.Errorf("%s: Failed to register classic load balancers: %v", machineName, err)
			errs = append(errs, err)
		}
	}
	if len(networkLoadBalancerNames) > 0 {
		err = registerWithNetworkLoadBalancers(client, networkLoadBalancerNames, instance)
		if err != nil {
			glog.Errorf("%s: Failed to register network load balancers: %v", machineName, err)
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errorutil.NewAggregate(errs)
	}
	return nil
}

// updateStatus calculates the new machine status, checks if anything has changed, and updates if so.
func (a *Actuator) updateStatus(machine *machinev1.Machine, instance *ec2.Instance) error {

	glog.Infof("%s: Updating status", machine.Name)

	// Starting with a fresh status as we assume full control of it here.
	awsStatus := &providerconfigv1.AWSMachineProviderStatus{}
	if err := a.codec.DecodeProviderStatus(machine.Status.ProviderStatus, awsStatus); err != nil {
		glog.Errorf("%s: Error decoding machine provider status: %v", machine.Name, err)
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
		if instance.PublicIpAddress != nil {
			networkAddresses = append(networkAddresses, corev1.NodeAddress{
				Type:    corev1.NodeExternalIP,
				Address: *instance.PublicIpAddress,
			})
		}
		if instance.PrivateIpAddress != nil {
			networkAddresses = append(networkAddresses, corev1.NodeAddress{
				Type:    corev1.NodeInternalIP,
				Address: *instance.PrivateIpAddress,
			})
		}
		if instance.PublicDnsName != nil {
			networkAddresses = append(networkAddresses, corev1.NodeAddress{
				Type:    corev1.NodeExternalDNS,
				Address: *instance.PublicDnsName,
			})
		}
		if instance.PrivateDnsName != nil {
			networkAddresses = append(networkAddresses, corev1.NodeAddress{
				Type:    corev1.NodeInternalDNS,
				Address: *instance.PrivateDnsName,
			})
		}
	}
	glog.Infof("%s: finished calculating AWS status", machine.Name)

	awsStatus.Conditions = setAWSMachineProviderCondition(awsStatus.Conditions, providerconfigv1.MachineCreation, corev1.ConditionTrue, MachineCreationSucceeded, "machine successfully created", updateConditionIfReasonOrMessageChange)
	// TODO(jchaloup): do we really need to update tis?
	// origInstanceID := awsStatus.InstanceID
	// if !StringPtrsEqual(origInstanceID, awsStatus.InstanceID) {
	// 	mLog.Debug("AWS instance ID changed, clearing LastELBSync to trigger adding to ELBs")
	// 	awsStatus.LastELBSync = nil
	// }

	if err := a.updateMachineStatus(machine, awsStatus, networkAddresses); err != nil {
		return err
	}

	// If machine state is still pending, we will return an error to keep the controllers
	// attempting to update status until it hits a more permanent state. This will ensure
	// we get a public IP populated more quickly.
	if awsStatus.InstanceState != nil && *awsStatus.InstanceState == ec2.InstanceStateNamePending {
		glog.Infof("%s: Instance state still pending, returning an error to requeue", machine.Name)
		return &clustererror.RequeueAfterError{RequeueAfter: requeueAfterSeconds * time.Second}
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
