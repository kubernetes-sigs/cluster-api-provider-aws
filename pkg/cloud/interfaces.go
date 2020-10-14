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
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/throttle"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
)

// Session represents an AWS session
type Session interface {
	Session() awsclient.ConfigProvider
	ServiceLimiter(string) *throttle.ServiceLimiter
}

// ScopeUsage is used to indicate which controller is using a scope
type ScopeUsage interface {
	// ControllerName returns the name of the controller that created the scope
	ControllerName() string
}

// ClusterObject represents a AWS cluster object
type ClusterObject interface {
	conditions.Setter
}

// ClusterScoper is the interface for a cluster scope
type ClusterScoper interface {
	logr.Logger
	Session
	ScopeUsage

	// Name returns the cluster name.
	Name() string
	// Namespace returns the cluster namespace.
	Namespace() string
	// Region returns the cluster region.
	Region() string

	// InfraCluster returns the AWS infrastructure cluster object.
	InfraCluster() ClusterObject

	// ListOptionsLabelSelector returns a ListOptions with a label selector for clusterName.
	ListOptionsLabelSelector() client.ListOption
	// APIServerPort returns the port to use when communicating with the API server.
	APIServerPort() int32
	// AdditionalTags returns any tags that you would like to attach to AWS resources. The returned value will never be nil.
	AdditionalTags() infrav1.Tags
	// SetFailureDomain sets the infrastructure provider failure domain key to the spec given as input.
	SetFailureDomain(id string, spec clusterv1.FailureDomainSpec)

	// PatchObject persists the cluster configuration and status.
	PatchObject() error
	// Close closes the current scope persisting the cluster configuration and status.
	Close() error
}
