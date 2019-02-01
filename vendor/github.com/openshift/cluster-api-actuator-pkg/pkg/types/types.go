package types

import machinev1beta1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"

// CloudProviderClient interface to generalize queries over various cloud providers
type CloudProviderClient interface {
	// Get running instances (of a given cloud provider) managed by the machine object
	GetRunningInstances(machine *machinev1beta1.Machine) ([]interface{}, error)
	// Get running instance public DNS name
	GetPublicDNSName(machine *machinev1beta1.Machine) (string, error)
	// Get private IP
	GetPrivateIP(machine *machinev1beta1.Machine) (string, error)
}
