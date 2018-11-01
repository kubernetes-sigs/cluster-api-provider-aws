// Copyright © 2018 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"github.com/aws/aws-sdk-go/aws/session"
	providerv1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

// Getter is a unified interfaces that includes all the getters.
type Getter interface {
	SDKSessionGetter
	EC2Getter
	ELBGetter
}

// SDKSessionGetter has a single method that returns an AWS session.
type SDKSessionGetter interface {
	Session(*providerv1.AWSClusterProviderConfig) *session.Session
}

// EC2Getter has a single method that returns an EC2 service interface.
type EC2Getter interface {
	EC2(*session.Session) EC2Interface
}

// ELBGetter has a single method that returns an ELB service interface.
type ELBGetter interface {
	ELB(*session.Session) ELBInterface
}

// EC2Interface encapsulates the methods exposed by the ec2 service.
type EC2Interface interface {
	EC2ClusterInterface
	EC2MachineInterface
}

// EC2ClusterInterface encapsulates the methods exposed to the cluster
// actuator
type EC2ClusterInterface interface {
	ReconcileNetwork(clusterName string, network *providerv1.Network) error
	ReconcileBastion(clusterName, keyName string, status *providerv1.AWSClusterProviderStatus) error
	DeleteNetwork(clusterName string, network *providerv1.Network) error
}

// EC2MachineInterface encapsulates the methods exposed to the machine
// actuator
type EC2MachineInterface interface {
	InstanceIfExists(instanceID *string) (*providerv1.Instance, error)
	TerminateInstance(instanceID string) error
	DeleteBastion(instanceID string, status *providerv1.AWSClusterProviderStatus) error
	CreateOrGetMachine(machine *clusterv1.Machine, status *providerv1.AWSMachineProviderStatus, config *providerv1.AWSMachineProviderConfig, clusterStatus *providerv1.AWSClusterProviderStatus, clusterConfig *providerv1.AWSClusterProviderConfig, cluster *clusterv1.Cluster, bootstrapToken string) (*providerv1.Instance, error)
	UpdateInstanceSecurityGroups(instanceID string, securityGroups []string) error
	UpdateResourceTags(resourceID *string, create map[string]string, remove map[string]string) error
}

// ELBInterface encapsulates the methods exposed by the elb service.
type ELBInterface interface {
	ReconcileLoadbalancers(clusterName string, network *providerv1.Network) error
	DeleteLoadbalancers(clusterName string, network *providerv1.Network) error
	RegisterInstanceWithAPIServerELB(clusterName string, instanceID string) error
	GetAPIServerDNSName(clusterName string) (string, error)
}
