package machine

import (
	"fmt"

	machinev1beta1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	awsclient "sigs.k8s.io/cluster-api-provider-aws/pkg/client"
)

// AwsClientWrapper implements CloudProviderClient for aws e2e framework
type AwsClientWrapper struct {
	client awsclient.Client
}

// NewAwsClientWrapper returns aws client implementation
func NewAwsClientWrapper(client awsclient.Client) *AwsClientWrapper {
	return &AwsClientWrapper{client: client}
}

// GetRunningInstances gets running instances (of a given cloud provider) managed by the machine object
func (client *AwsClientWrapper) GetRunningInstances(machine *machinev1beta1.Machine) ([]interface{}, error) {
	runningInstances, err := getRunningInstances(machine, client.client)
	if err != nil {
		return nil, err
	}

	var instances []interface{}
	for _, instance := range runningInstances {
		instances = append(instances, instance)
	}

	return instances, nil
}

// GetPublicDNSName gets running instance public DNS name
func (client *AwsClientWrapper) GetPublicDNSName(machine *machinev1beta1.Machine) (string, error) {
	instance, err := getRunningInstance(machine, client.client)
	if err != nil {
		return "", err
	}

	if *instance.PublicDnsName == "" {
		return "", fmt.Errorf("machine instance public DNS name not set")
	}

	return *instance.PublicDnsName, nil
}

// GetPrivateIP gets private IP
func (client *AwsClientWrapper) GetPrivateIP(machine *machinev1beta1.Machine) (string, error) {
	instance, err := getRunningInstance(machine, client.client)
	if err != nil {
		return "", err
	}

	if *instance.PrivateIpAddress == "" {
		return "", fmt.Errorf("machine instance public DNS name not set")
	}

	return *instance.PrivateIpAddress, nil
}

// GetSecurityGroups gets security groups
func (client *AwsClientWrapper) GetSecurityGroups(machine *machinev1beta1.Machine) ([]string, error) {
	instance, err := getRunningInstance(machine, client.client)
	if err != nil {
		return nil, err
	}
	var groups []string
	for _, groupIdentifier := range instance.SecurityGroups {
		if *groupIdentifier.GroupName != "" {
			groups = append(groups, *groupIdentifier.GroupName)
		}
	}
	return groups, nil
}

// GetIAMRole gets IAM role
func (client *AwsClientWrapper) GetIAMRole(machine *machinev1beta1.Machine) (string, error) {
	instance, err := getRunningInstance(machine, client.client)
	if err != nil {
		return "", err
	}
	if instance.IamInstanceProfile == nil {
		return "", err
	}
	return *instance.IamInstanceProfile.Id, nil
}

// GetTags gets tags
func (client *AwsClientWrapper) GetTags(machine *machinev1beta1.Machine) (map[string]string, error) {
	instance, err := getRunningInstance(machine, client.client)
	if err != nil {
		return nil, err
	}
	tags := make(map[string]string, len(instance.Tags))
	for _, tag := range instance.Tags {
		tags[*tag.Key] = *tag.Value
	}
	return tags, nil
}

// GetSubnet gets subnet
func (client *AwsClientWrapper) GetSubnet(machine *machinev1beta1.Machine) (string, error) {
	instance, err := getRunningInstance(machine, client.client)
	if err != nil {
		return "", err
	}
	if instance.SubnetId == nil {
		return "", err
	}
	return *instance.SubnetId, nil
}

// GetAvailabilityZone gets availability zone
func (client *AwsClientWrapper) GetAvailabilityZone(machine *machinev1beta1.Machine) (string, error) {
	instance, err := getRunningInstance(machine, client.client)
	if err != nil {
		return "", err
	}
	if instance.Placement == nil {
		return "", err
	}
	return *instance.Placement.AvailabilityZone, nil
}

// GetVolumes gets volumes attached to instance
func (client *AwsClientWrapper) GetVolumes(machine *machinev1beta1.Machine) (map[string]map[string]interface{}, error) {
	instance, err := getRunningInstance(machine, client.client)
	if err != nil {
		return nil, err
	}
	if instance.BlockDeviceMappings == nil {
		return nil, err
	}
	volumes := make(map[string]map[string]interface{}, len(instance.BlockDeviceMappings))
	for _, blockDeviceMapping := range instance.BlockDeviceMappings {
		volume, err := getVolume(client.client, *blockDeviceMapping.Ebs.VolumeId)
		if err != nil {
			return volumes, err
		}
		volumes[*blockDeviceMapping.DeviceName] = map[string]interface{}{
			"id":   *volume.VolumeId,
			"iops": *volume.Iops,
			"size": *volume.Size,
			"type": *volume.VolumeType,
		}
	}
	return volumes, nil
}
