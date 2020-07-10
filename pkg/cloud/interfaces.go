/*
Copyright 2020 The Kubernetes Authors.

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

package cloud

import (
	"github.com/go-logr/logr"

	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/controller-runtime/pkg/client"

	awsclient "github.com/aws/aws-sdk-go/aws/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
)

// Session represents an AWS session
type Session interface {
	Session() awsclient.ConfigProvider
}

// ClusterObject represents a AWS cluster object
type ClusterObject interface {
	conditions.Setter
}

// ClusterScoper is the interface for a cluster scopde
type ClusterScoper interface {
	logr.Logger
	Session

	Name() string
	Namespace() string
	Region() string

	InfraCluster() ClusterObject

	Network() *infrav1.Network
	VPC() *infrav1.VPCSpec
	Subnets() infrav1.Subnets
	SetSubnets(subnets infrav1.Subnets)
	CNIIngressRules() infrav1.CNIIngressRules
	SecurityGroups() map[infrav1.SecurityGroupRole]infrav1.SecurityGroup
	ListOptionsLabelSelector() client.ListOption
	PatchObject() error
	Close() error
	APIServerPort() int32
	AdditionalTags() infrav1.Tags
	SetFailureDomain(id string, spec clusterv1.FailureDomainSpec)
	Bastion() *infrav1.Bastion
	SetBastionInstance(instance *infrav1.Instance)
	SSHKeyName() *string
}
