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

package services

import (
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
)

// EC2MachineInterface encapsulates the methods exposed to the machine
// actuator
type EC2MachineInterface interface {
	InstanceIfExists(id *string) (*infrav1.Instance, error)
	TerminateInstance(id string) error
	CreateInstance(scope *scope.MachineScope, userData []byte) (*infrav1.Instance, error)
	GetRunningInstanceByTags(scope *scope.MachineScope) (*infrav1.Instance, error)

	GetCoreSecurityGroups(machine *scope.MachineScope) ([]string, error)
	GetInstanceSecurityGroups(instanceID string) (map[string][]string, error)
	UpdateInstanceSecurityGroups(id string, securityGroups []string) error
	UpdateResourceTags(resourceID *string, create, remove map[string]string) error

	TerminateInstanceAndWait(instanceID string) error
	DetachSecurityGroupsFromNetworkInterface(groups []string, interfaceID string) error
}

// SecretsManagerInterface encapsulated the methods exposed to the
// machine actuator
type SecretsManagerInterface interface {
	Delete(m *scope.MachineScope) error
	Create(m *scope.MachineScope, data []byte) (string, int32, error)
}
