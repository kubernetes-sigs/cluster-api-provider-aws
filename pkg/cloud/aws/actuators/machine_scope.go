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

package actuators

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/tags"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	client "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1"
)

// MachineScopeParams defines the input parameters used to create a new MachineScope.
type MachineScopeParams struct {
	AWSClients
	Cluster    *clusterv1.Cluster
	Machine    *clusterv1.Machine
	Client     client.ClusterV1alpha1Interface
	KubeClient kubernetes.Interface
}

// NewMachineScope creates a new MachineScope from the supplied parameters.
// This is meant to be called for each machine actuator operation.
func NewMachineScope(params MachineScopeParams) (*MachineScope, error) {
	scope, err := NewScope(ScopeParams{AWSClients: params.AWSClients, Client: params.Client, Cluster: params.Cluster, Machine: params.Machine, KubeClient: params.KubeClient})
	if err != nil {
		return nil, err
	}
	machineConfig, err := v1alpha1.MachineConfigFromProviderSpec(params.Machine.Spec.ProviderSpec)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get machine config")
	}

	machineStatus, err := v1alpha1.MachineStatusFromProviderStatus(params.Machine.Status.ProviderStatus)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get machine provider status")
	}

	var machineClient client.MachineInterface
	if params.Client != nil {
		machineClient = params.Client.Machines(params.Machine.Namespace)
	}
	klog.Infof("machineStatus: %+v", machineStatus)

	return &MachineScope{
		Scope: scope,
		Machine:       params.Machine,
		MachineClient: machineClient,
		MachineConfig: machineConfig,
		MachineStatus: machineStatus,
	}, nil
}

// MachineScope defines a scope defined around a machine and its cluster.
type MachineScope struct {
	*Scope
	Machine       *clusterv1.Machine
	MachineClient client.MachineInterface
	MachineConfig *v1alpha1.AWSMachineProviderSpec
	MachineStatus *v1alpha1.AWSMachineProviderStatus
}

// Name returns the machine name.
func (m *MachineScope) Name() string {
	return m.Machine.Name
}

// ClusterName returns the cluster name.
func (m *MachineScope) ClusterName() string {
	return m.Scope.Name()
}

// ClusterID returns the clusterid from the labels.
func (m *MachineScope) ClusterID() string {
	return m.Machine.Labels[tags.NameAWSClusterIDLabel]
}

// Namespace returns the machine namespace.
func (m *MachineScope) Namespace() string {
	return m.Machine.Namespace
}

// Name returns the machine role from the labels.
func (m *MachineScope) Role() string {
	return m.Machine.Labels["set"]
}

// Region returns the machine region.
func (m *MachineScope) Region() string {
	return m.Scope.Region()
}

func (m *MachineScope) storeMachineSpec(machine *clusterv1.Machine) (*clusterv1.Machine, error) {
	ext, err := v1alpha1.EncodeMachineSpec(m.MachineConfig)
	if err != nil {
		return nil, err
	}

	latest, err := m.MachineClient.Get(m.Machine.Name, metav1.GetOptions{})
	if err != nil {
		klog.Errorf("failed to get machine from server: %v", err)
	}
	machine.Spec.ProviderSpec.Value = ext
	machine.ResourceVersion = latest.ResourceVersion
	return m.MachineClient.Update(machine)
}

func (m *MachineScope) storeMachineStatus(machine *clusterv1.Machine) (*clusterv1.Machine, error) {
	ext, err := v1alpha1.EncodeMachineStatus(m.MachineStatus)
	if err != nil {
		return nil, err
	}

	machine.Status.ProviderStatus = ext
	return m.MachineClient.UpdateStatus(machine)
}

func buildEC2Filters(inputFilters []v1alpha1.Filter) []*ec2.Filter {
	filters := make([]*ec2.Filter, len(inputFilters))
	for i, f := range inputFilters {
		values := make([]*string, len(f.Values))
		for j, v := range f.Values {
			values[j] = aws.String(v)
		}
		filters[i] = &ec2.Filter{
			Name:   aws.String(f.Name),
			Values: values,
		}
	}
	return filters
}

func (s *MachineScope) SecurityGroups() []string {
	var securityGroupIDs []string
	for _, g := range s.MachineConfig.AdditionalSecurityGroups {
		// ID has priority
		if g.ID != nil {
			securityGroupIDs = append(securityGroupIDs, *g.ID)
		} else if g.Filters != nil {
			klog.Info("Describing security groups based on filters")
			// Get groups based on filters
			describeSecurityGroupsRequest := ec2.DescribeSecurityGroupsInput{
				Filters: buildEC2Filters(g.Filters),
			}
			describeSecurityGroupsResult, err := s.AWSClients.EC2.DescribeSecurityGroups(&describeSecurityGroupsRequest)
			if err != nil {
				klog.Errorf("error describing security groups: %v", err)
				return nil
			}
			for _, g := range describeSecurityGroupsResult.SecurityGroups {
				groupID := *g.GroupId
				securityGroupIDs = append(securityGroupIDs, groupID)
			}
		}
	}

	if len(securityGroupIDs) == 0 {
		klog.Info("No security group found")
	}

	return securityGroupIDs

}

func (m *MachineScope) Close() {
	if m.MachineClient == nil {
		return
	}

	latestMachine, err := m.storeMachineSpec(m.Machine)
	if err != nil {
		klog.Errorf("[machinescope] failed to update machine %q in namespace %q: %v", m.Machine.Name, m.Machine.Namespace, err)
	}

	_, err = m.storeMachineStatus(latestMachine)
	if err != nil {
		klog.Errorf("[machinescope] failed to store provider status for machine %q in namespace %q: %v", m.Machine.Name, m.Machine.Namespace, err)
	}
}
