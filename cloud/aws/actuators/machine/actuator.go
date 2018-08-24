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
	"encoding/base64"
	"fmt"

	"github.com/golang/glog"
	log "github.com/sirupsen/logrus"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	providerconfigv1 "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	clusterclient "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset"
	client "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	awsclient "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/client"
	clustoplog "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/logging"
)

const (
	userDataSecretKey         = "userData"
	ec2InstanceIDNotFoundCode = "InvalidInstanceID.NotFound"
)

// Actuator is the AWS-specific actuator for the Cluster API machine controller
type Actuator struct {
	kubeClient       kubernetes.Interface
	clusterClient    clusterclient.Interface
	logger           *log.Entry
	awsClientBuilder awsclient.AwsClientBuilderFuncType
}

// ActuatorParams holds parameter information for Actuator
type ActuatorParams struct {
	ClusterClient client.ClusterInterface
}

// NewActuator returns a new AWS Actuator
func NewActuator(kubeClient kubernetes.Interface, clusterClient clusterclient.Interface, logger *log.Entry, awsClientBuilder awsclient.AwsClientBuilderFuncType) (*Actuator, error) {
	actuator := &Actuator{
		kubeClient:       kubeClient,
		clusterClient:    clusterClient,
		logger:           logger,
		awsClientBuilder: awsClientBuilder,
	}
	return actuator, nil
}

// Create runs a new EC2 instance
func (a *Actuator) Create(cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	mLog := clustoplog.WithMachine(a.logger, machine)
	mLog.Info("creating machine")
	instance, err := a.CreateMachine(cluster, machine)
	if err != nil {
		mLog.Errorf("error creating machine: %v", err)
		return err
	}

	// TODO(csrwng):
	// Part of the status that gets updated when the machine gets created is the PublicIP.
	// However, after a call to runInstance, most of the time a PublicIP is not yet allocated.
	// If we don't yet have complete status (ie. the instance state is Pending, instead of Running),
	// maybe we should return an error so the machine controller keeps retrying until we have complete status we can set.
	return a.updateStatus(machine, instance, mLog)
}

// removeStoppedMachine removes all instances of a specific machine that are in a stopped state.
func (a *Actuator) removeStoppedMachine(machine *clusterv1.Machine, client awsclient.Client, mLog log.FieldLogger) error {
	instances, err := GetStoppedInstances(machine, client)
	if err != nil {
		mLog.Errorf("Error getting stopped instances: %v", err)
		return fmt.Errorf("Error getting stopped instances: %v", err)
	}

	if len(instances) == 0 {
		mLog.Infof("no stopped instances found for machine %v", machine.Name)
		return nil
	}

	return TerminateInstances(client, instances, mLog)
}

// CreateMachine starts a new AWS instance as described by the cluster and machine resources
func (a *Actuator) CreateMachine(cluster *clusterv1.Cluster, machine *clusterv1.Machine) (*ec2.Instance, error) {
	mLog := clustoplog.WithMachine(a.logger, machine)

	machineProviderConfig, err := MachineProviderConfigFromClusterAPIMachineSpec(&machine.Spec)
	if err != nil {
		mLog.Errorf("error decoding MachineProviderConfig: %v", err)
		return nil, err
	}

	credentialsSecretName := ""
	if machineProviderConfig.CredentialsSecret != nil {
		credentialsSecretName = machineProviderConfig.CredentialsSecret.Name
	}
	client, err := a.awsClientBuilder(a.kubeClient, credentialsSecretName, machine.Namespace, machineProviderConfig.Placement.Region)
	if err != nil {
		mLog.Errorf("unable to obtain AWS client: %v", err)
		return nil, fmt.Errorf("unable to obtain AWS client: %v", err)
	}

	// We explicitly do NOT want to remove stopped masters.
	if !MachineIsMaster(machine) {
		// Prevent having a lot of stopped nodes sitting around.
		err = a.removeStoppedMachine(machine, client, mLog)
		if err != nil {
			mLog.Errorf("unable to remove stopped machines: %v", err)
			return nil, fmt.Errorf("unable to remove stopped nodes: %v", err)
		}
	}

	// TODO(jchaloup): resolve ARN and/or Filters as well
	if machineProviderConfig.AMI.ID == nil {
		return nil, fmt.Errorf("machine does not have an AWS image set")
	}

	// Get AMI to use
	amiName := *machineProviderConfig.AMI.ID

	mLog.Debugf("Describing AMI %s", amiName)
	imageIds := []*string{aws.String(amiName)}
	describeImagesRequest := ec2.DescribeImagesInput{
		ImageIds: imageIds,
	}
	describeAMIResult, err := client.DescribeImages(&describeImagesRequest)
	if err != nil {
		mLog.Errorf("error describing AMI %s: %v", amiName, err)
		return nil, fmt.Errorf("error describing AMI %s: %v", amiName, err)
	}
	if len(describeAMIResult.Images) != 1 {
		mLog.Errorf("Unexpected number of images returned: %d", len(describeAMIResult.Images))
		return nil, fmt.Errorf("Unexpected number of images returned: %d", len(describeAMIResult.Images))
	}

	var securityGroupIds []*string
	for _, g := range machineProviderConfig.SecurityGroups {
		groupID := *g.ID
		securityGroupIds = append(securityGroupIds, &groupID)
	}

	// build list of networkInterfaces (just 1 for now)
	var networkInterfaces = []*ec2.InstanceNetworkInterfaceSpecification{
		{
			DeviceIndex:              aws.Int64(machineProviderConfig.DeviceIndex),
			AssociatePublicIpAddress: machineProviderConfig.PublicIP,
			SubnetId:                 machineProviderConfig.Subnet.ID,
			Groups:                   securityGroupIds,
		},
	}

	// Add tags to the created machine
	tagList := []*ec2.Tag{}
	for _, tag := range machineProviderConfig.Tags {
		tagList = append(tagList, &ec2.Tag{Key: aws.String(tag.Name), Value: aws.String(tag.Value)})
	}
	tagList = append(tagList, []*ec2.Tag{
		{Key: aws.String("clusterid"), Value: aws.String(cluster.Name)},
		{Key: aws.String("kubernetes.io/cluster/" + cluster.Name), Value: aws.String(cluster.Name)},
		{Key: aws.String("Name"), Value: aws.String(machine.Name)},
	}...)

	tagInstance := &ec2.TagSpecification{
		ResourceType: aws.String("instance"),
		Tags:         tagList,
	}
	tagVolume := &ec2.TagSpecification{
		ResourceType: aws.String("volume"),
		Tags:         []*ec2.Tag{{Key: aws.String("clusterid"), Value: aws.String(cluster.Name)}},
	}

	userData := []byte{}
	if machineProviderConfig.UserDataSecret != nil {
		userDataSecret, err := a.kubeClient.CoreV1().Secrets(machine.Namespace).Get(machineProviderConfig.UserDataSecret.Name, metav1.GetOptions{})
		if err != nil {
			mLog.Errorf("error getting user data secret %s: %v", machineProviderConfig.UserDataSecret.Name, err)
			return nil, err
		}
		if data, exists := userDataSecret.Data[userDataSecretKey]; exists {
			userData = data
		} else {
			glog.Warningf("Secret %v/%v does not have %q field set. Thus, no user data applied when creating an instance.", machine.Namespace, machineProviderConfig.UserDataSecret.Name, userDataSecretKey)
		}
	}

	userDataEnc := base64.StdEncoding.EncodeToString(userData)

	inputConfig := ec2.RunInstancesInput{
		ImageId:      describeAMIResult.Images[0].ImageId,
		InstanceType: aws.String(machineProviderConfig.InstanceType),
		// Only a single instance of the AWS instance allowed
		MinCount: aws.Int64(1),
		MaxCount: aws.Int64(1),
		KeyName:  machineProviderConfig.KeyName,
		IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
			Name: aws.String(*machineProviderConfig.IAMInstanceProfile.ID),
		},
		TagSpecifications: []*ec2.TagSpecification{tagInstance, tagVolume},
		NetworkInterfaces: networkInterfaces,
		UserData:          &userDataEnc,
	}

	runResult, err := client.RunInstances(&inputConfig)
	if err != nil {
		mLog.Errorf("error creating EC2 instance: %v", err)
		return nil, fmt.Errorf("error creating EC2 instance: %v", err)
	}

	if runResult == nil || len(runResult.Instances) != 1 {
		mLog.Errorf("unexpected reservation creating instances: %v", runResult)
		return nil, fmt.Errorf("unexpected reservation creating instance")
	}

	// TOOD(csrwng):
	// One issue we have right now with how this works, is that usually at the end of the RunInstances call,
	// the instance state is not yet 'Running'. Rather it is 'Pending'. Therefore the status
	// is not yet complete (like PublicIP). One possible fix would be to wait and poll here
	// until the instance is in the Running state. The other approach is to return an error
	// so that the machine is requeued and in the exists function return false if the status doesn't match.
	// That would require making the create re-entrant so we can just update the status.
	return runResult.Instances[0], nil
}

// Delete deletes a machine and updates its finalizer
func (a *Actuator) Delete(cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	mLog := clustoplog.WithMachine(a.logger, machine)
	mLog.Info("deleting machine")
	if err := a.DeleteMachine(cluster, machine); err != nil {
		mLog.Errorf("error deleting machine: %v", err)
		return err
	}
	return nil
}

// DeleteMachine deletes an AWS instance
func (a *Actuator) DeleteMachine(cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	mLog := clustoplog.WithMachine(a.logger, machine)

	machineProviderConfig, err := MachineProviderConfigFromClusterAPIMachineSpec(&machine.Spec)
	if err != nil {
		mLog.Errorf("error decoding MachineProviderConfig: %v", err)
		return err
	}

	region := machineProviderConfig.Placement.Region
	credentialsSecretName := ""
	if machineProviderConfig.CredentialsSecret != nil {
		credentialsSecretName = machineProviderConfig.CredentialsSecret.Name
	}
	client, err := a.awsClientBuilder(a.kubeClient, credentialsSecretName, machine.Namespace, region)
	if err != nil {
		mLog.Errorf("error getting EC2 client: %v", err)
		return fmt.Errorf("error getting EC2 client: %v", err)
	}

	instances, err := GetRunningInstances(machine, client)
	if err != nil {
		mLog.Errorf("error getting running instances: %v", err)
		return err
	}
	if len(instances) == 0 {
		mLog.Warnf("no instances found to delete for machine")
		return nil
	}

	return TerminateInstances(client, instances, mLog)
}

// Update attempts to sync machine state with an existing instance. Today this just updates status
// for details that may have changed. (IPs and hostnames) We do not currently support making any
// changes to actual machines in AWS. Instead these will be replaced via MachineDeployments.
func (a *Actuator) Update(cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	mLog := clustoplog.WithMachine(a.logger, machine)
	mLog.Debugf("updating machine")

	machineProviderConfig, err := MachineProviderConfigFromClusterAPIMachineSpec(&machine.Spec)
	if err != nil {
		mLog.Errorf("error decoding MachineProviderConfig: %v", err)
		return err
	}

	region := machineProviderConfig.Placement.Region
	mLog.WithField("region", region).Debugf("obtaining EC2 client for region")
	credentialsSecretName := ""
	if machineProviderConfig.CredentialsSecret != nil {
		credentialsSecretName = machineProviderConfig.CredentialsSecret.Name
	}
	client, err := a.awsClientBuilder(a.kubeClient, credentialsSecretName, machine.Namespace, region)
	if err != nil {
		mLog.Errorf("error getting EC2 client: %v", err)
		return fmt.Errorf("unable to obtain EC2 client: %v", err)
	}

	instances, err := GetRunningInstances(machine, client)
	if err != nil {
		mLog.Errorf("error getting running instances: %v", err)
		return err
	}
	mLog.Debugf("found %d instances for machine", len(instances))

	// Parent controller should prevent this from ever happening by calling Exists and then Create,
	// but instance could be deleted between the two calls.
	if len(instances) == 0 {
		mLog.Warnf("attempted to update machine but no instances found")
		// Update status to clear out machine details.
		err := a.updateStatus(machine, nil, mLog)
		if err != nil {
			return err
		}
		mLog.Errorf("attempted to update machine but no instances found")
		return fmt.Errorf("attempted to update machine but no instances found")
	}
	newestInstance, terminateInstances := SortInstances(instances)

	// In very unusual circumstances, there could be more than one machine running matching this
	// machine name and cluster ID. In this scenario we will keep the newest, and delete all others.
	mLog = mLog.WithField("instanceID", *newestInstance.InstanceId)
	mLog.Debug("instance found")

	if len(instances) > 1 {
		err = TerminateInstances(client, terminateInstances, mLog)
		if err != nil {
			return err
		}

	}

	// We do not support making changes to pre-existing instances, just update status.
	return a.updateStatus(machine, newestInstance, mLog)
}

// Exists determines if the given machine currently exists. For AWS we query for instances in
// running state, with a matching name tag, to determine a match.
func (a *Actuator) Exists(cluster *clusterv1.Cluster, machine *clusterv1.Machine) (bool, error) {
	mLog := clustoplog.WithMachine(a.logger, machine)
	mLog.Debugf("checking if machine exists")

	machineProviderConfig, err := MachineProviderConfigFromClusterAPIMachineSpec(&machine.Spec)
	if err != nil {
		mLog.Errorf("error decoding MachineProviderConfig: %v", err)
		return false, err
	}

	region := machineProviderConfig.Placement.Region
	credentialsSecretName := ""
	if machineProviderConfig.CredentialsSecret != nil {
		credentialsSecretName = machineProviderConfig.CredentialsSecret.Name
	}
	client, err := a.awsClientBuilder(a.kubeClient, credentialsSecretName, machine.Namespace, region)
	if err != nil {
		mLog.Errorf("error getting EC2 client: %v", err)
		return false, fmt.Errorf("error getting EC2 client: %v", err)
	}

	instances, err := GetRunningInstances(machine, client)
	if err != nil {
		mLog.Errorf("error getting running instances: %v", err)
		return false, err
	}
	if len(instances) == 0 {
		mLog.Debug("instance does not exist")
		return false, nil
	}

	// If more than one result was returned, it will be handled in Update.
	mLog.Debugf("instance exists as %q", *instances[0].InstanceId)
	return true, nil
}

// updateStatus calculates the new machine status, checks if anything has changed, and updates if so.
func (a *Actuator) updateStatus(machine *clusterv1.Machine, instance *ec2.Instance, mLog log.FieldLogger) error {

	mLog.Debug("updating status")

	// Starting with a fresh status as we assume full control of it here.
	awsStatus, err := AWSMachineProviderStatusFromClusterAPIMachine(machine)
	if err != nil {
		mLog.Errorf("error decoding machine provider status: %v", err)
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
	mLog.Debug("finished calculating AWS status")

	// TODO(jchaloup): do we really need to update tis?
	// origInstanceID := awsStatus.InstanceID
	// if !StringPtrsEqual(origInstanceID, awsStatus.InstanceID) {
	// 	mLog.Debug("AWS instance ID changed, clearing LastELBSync to trigger adding to ELBs")
	// 	awsStatus.LastELBSync = nil
	// }

	awsStatusRaw, err := EncodeAWSMachineProviderStatus(awsStatus)
	if err != nil {
		mLog.Errorf("error encoding AWS provider status: %v", err)
		return err
	}

	machineCopy := machine.DeepCopy()
	machineCopy.Status.ProviderStatus = awsStatusRaw
	machineCopy.Status.Addresses = networkAddresses

	if !equality.Semantic.DeepEqual(machine.Status, machineCopy.Status) {
		mLog.Info("machine status has changed, updating")
		machineCopy.Status.LastUpdated = metav1.Now()

		_, err := a.clusterClient.ClusterV1alpha1().Machines(machineCopy.Namespace).UpdateStatus(machineCopy)
		if err != nil {
			mLog.Errorf("error updating machine status: %v", err)
			return err
		}
	} else {
		mLog.Debug("status unchanged")
	}
	return nil
}

func getClusterID(machine *clusterv1.Machine) (string, bool) {
	clusterID, ok := machine.Labels[providerconfigv1.ClusterNameLabel]
	return clusterID, ok
}
