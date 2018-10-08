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
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	providerv1 "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
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
	EC2(*session.Session, *aws.Config) EC2Interface
}

// ELBGetter has a single method that returns an ELB service interface.
type ELBGetter interface {
	ELB(*session.Session, *aws.Config) ELBInterface
}

// EC2Interface encapsulates the methods exposed by the ec2 service.
type EC2Interface interface {
	EC2ClusterInterface
	EC2MachineInterface
}

// EC2ClusterInterface encapsulates the methods exposed to the cluster
// actuator
type EC2ClusterInterface interface {
	ReconcileNetwork(ctx context.Context, clusterName string, network *providerv1.Network) error
	ReconcileBastion(ctx context.Context, clusterName, keyName string, status *providerv1.AWSClusterProviderStatus) error
	DeleteNetwork(ctx context.Context, clusterName string, network *providerv1.Network) error
}

// EC2MachineInterface encapsulates the methods exposed to the machine
// actuator
type EC2MachineInterface interface {
	InstanceIfExists(ctx context.Context, instanceID *string) (*providerv1.Instance, error)
	CreateInstance(ctx context.Context, machine *clusterv1.Machine, config *providerv1.AWSMachineProviderConfig, clusterStatus *providerv1.AWSClusterProviderStatus) (*providerv1.Instance, error)
	TerminateInstance(ctx context.Context, instanceID string) error
	DeleteBastion(ctx context.Context, instanceID string, status *providerv1.AWSClusterProviderStatus) error
	CreateOrGetMachine(ctx context.Context, machine *clusterv1.Machine, status *providerv1.AWSMachineProviderStatus, config *providerv1.AWSMachineProviderConfig, clusterStatus *providerv1.AWSClusterProviderStatus) (*providerv1.Instance, error)
	UpdateInstanceSecurityGroups(ctx context.Context, instanceID *string, securityGroups []*string) error
	UpdateResourceTags(ctx context.Context, resourceID *string, create map[string]string, remove map[string]string) error
}

// ELBInterface encapsulates the methods exposed by the elb service.
type ELBInterface interface {
	ReconcileLoadbalancers(ctx context.Context, clusterName string, network *providerv1.Network) error
}
