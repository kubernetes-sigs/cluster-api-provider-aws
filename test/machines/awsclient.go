package machines

import (
	"fmt"

	machineutils "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators/machine"
	awsclient "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/client"
	clusterv1alpha1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

type awsClientWrapper struct {
	client awsclient.Client
}

func (client *awsClientWrapper) GetRunningInstances(machine *clusterv1alpha1.Machine) ([]interface{}, error) {
	runningInstances, err := machineutils.GetRunningInstances(machine, client.client)
	if err != nil {
		return nil, err
	}

	var instances []interface{}
	for _, instance := range runningInstances {
		instances = append(instances, instance)
	}

	return instances, nil
}

func (client *awsClientWrapper) GetPublicDNSName(machine *clusterv1alpha1.Machine) (string, error) {
	runningInstances, err := machineutils.GetRunningInstances(machine, client.client)
	if err != nil {
		return "", err
	}

	if len(runningInstances) == 0 {
		return "", fmt.Errorf("no running machine instance found")
	}

	if *runningInstances[0].PublicDnsName == "" {
		return "", fmt.Errorf("machine instance public DNS name not set")
	}

	return *runningInstances[0].PublicDnsName, nil
}

func (client *awsClientWrapper) GetPrivateIP(machine *clusterv1alpha1.Machine) (string, error) {
	runningInstances, err := machineutils.GetRunningInstances(machine, client.client)
	if err != nil {
		return "", err
	}

	if len(runningInstances) == 0 {
		return "", fmt.Errorf("no running machine instance found")
	}

	if *runningInstances[0].PrivateIpAddress == "" {
		return "", fmt.Errorf("machine instance public DNS name not set")
	}

	return *runningInstances[0].PrivateIpAddress, nil
}
