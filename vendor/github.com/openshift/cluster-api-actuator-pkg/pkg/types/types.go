package types

import clusterv1alpha1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

// CloudProviderClient interface to generalize queries over various cloud providers
type CloudProviderClient interface {
	// Get running instances (of a given cloud provider) managed by the machine object
	GetRunningInstances(machine *clusterv1alpha1.Machine) ([]interface{}, error)
	// Get running instance public DNS name
	GetPublicDNSName(machine *clusterv1alpha1.Machine) (string, error)
	// Get private IP
	GetPrivateIP(machine *clusterv1alpha1.Machine) (string, error)
}
