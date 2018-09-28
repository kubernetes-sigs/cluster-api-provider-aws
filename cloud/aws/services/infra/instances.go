package infra

import (
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/ec2"
)

type Infra struct {
	EC2Svc *ec2.Service
}

func (i *Infra) CreateOrGetMachine(machine *clusterv1.Machine, status *v1alpha1.AWSMachineProviderStatus, config *v1alpha1.AWSMachineProviderConfig, cluster *clusterv1.Cluster, clusterStatus *v1alpha1.AWSClusterProviderStatus, clusterConfig *v1alpha1.AWSClusterProviderConfig) (*ec2.Instance, error) {
	// instance id exists, try to get it
	if status.InstanceID != nil {
		instance, err := i.EC2Svc.InstanceIfExists(status.InstanceID)

		// if there was no error, return the found instance
		if err == nil {
			return instance, err
		}

		// if there was an error but it's not IsNotFound then it's a real error
		if !ec2.IsNotFound(err) {
			return instance, err
		}
	}

	// otherwise let's create it
	return i.EC2Svc.CreateInstance(machine, config, cluster, clusterConfig, clusterStatus)
}
