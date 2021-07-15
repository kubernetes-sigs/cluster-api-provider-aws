// +build e2e

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

package shared

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/gofrs/flock"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	cfn_bootstrap "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cloudformation/bootstrap"
	clusterctlv1 "sigs.k8s.io/cluster-api/cmd/clusterctl/api/v1alpha3"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/test/framework/bootstrap"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
)

// Option represents an option to use when creating a e2e context
type Option func(*E2EContext)

func NewE2EContext(options ...Option) *E2EContext {
	ctx := &E2EContext{
		IsManaged: false,
	}
	ctx.Environment.Scheme = DefaultScheme()
	ctx.Environment.Namespaces = map[*corev1.Namespace]context.CancelFunc{}
	//ctx.Lifecycle = DefaultGinkgoLifecycle()

	for _, opt := range options {
		opt(ctx)
	}

	return ctx
}

// E2EContext represents the context of the e2e test
type E2EContext struct {
	// Settings is the settings used for the test
	Settings Settings
	// E2EConfig to be used for this test, read from configPath.
	E2EConfig *clusterctl.E2EConfig
	// Environment represents the runtime enviroment
	Environment RuntimeEnvironment
	// AWSSession is the AWS session for the tests
	AWSSession client.ConfigProvider
	// BootstrapUserAWSSession is the AWS session for the bootstrap user
	BootstrapUserAWSSession client.ConfigProvider
	// IsManaged indicates that this is for the managed part of the provider
	IsManaged bool
	// CloudFormationTemplate is the rendered template created for the test
	CloudFormationTemplate *cloudformation.Template
	StartOfSuite           time.Time
}

// Settings represents the test settings
type Settings struct {
	// ConfigPath is the path to the e2e config file.
	ConfigPath string
	// useExistingCluster instructs the test to use the current cluster instead of creating a new one (default discovery rules apply).
	UseExistingCluster bool
	// ArtifactFolder is the folder to store e2e test artifacts.
	ArtifactFolder string
	// DataFolder is the root folder for the data required by the tests
	DataFolder string
	// SkipCleanup prevents cleanup of test resources e.g. for debug purposes.
	SkipCleanup bool
	// SkipCloudFormationCreation will skip the cloudformation execution - useful for debugging e2e tests
	SkipCloudFormationCreation bool
	// SkipCloudFormationDeletion prevents the deletion of the AWS CloudFormation stack
	SkipCloudFormationDeletion bool
	// number of ginkgo nodes to use for kubetest
	GinkgoNodes int
	// time in s before kubetest spec is marked as slow
	GinkgoSlowSpecThreshold int
	// kubetestConfigFilePath is the path to the kubetest configuration file
	KubetestConfigFilePath string
	// useCIArtifacts specifies whether or not to use the latest build from the main branch of the Kubernetes repository
	UseCIArtifacts bool
	// SourceTemplate specifies which source template to use
	SourceTemplate string
	// FileLock is the lock to be used to read the resource quotas file
	FileLock *flock.Flock
}

// RuntimeEnvironment represents the runtime environment of the test
type RuntimeEnvironment struct {
	// BootstrapClusterProvider manages provisioning of the the bootstrap cluster to be used for the e2e tests.
	// Please note that provisioning will be skipped if use-existing-cluster is provided.
	BootstrapClusterProvider bootstrap.ClusterProvider
	// BootstrapClusterProxy allows to interact with the bootstrap cluster to be used for the e2e tests.
	BootstrapClusterProxy framework.ClusterProxy
	// BootstrapTemplate is the clusterawsadm bootstrap template for this run
	BootstrapTemplate *cfn_bootstrap.Template
	// BootstrapAccessKey is the bootstrap user access key
	BootstrapAccessKey *iam.AccessKey
	// ResourceTicker for dumping resources
	ResourceTicker *time.Ticker
	// ResourceTickerDone to stop ticking
	ResourceTickerDone chan bool
	// MachineTicker for dumping resources
	MachineTicker *time.Ticker
	// MachineTickerDone to stop ticking
	MachineTickerDone chan bool
	// Namespaces holds the namespaces used in the tests
	Namespaces map[*corev1.Namespace]context.CancelFunc
	// ClusterctlConfigPath to be used for this test, created by generating a clusterctl local repository
	// with the providers specified in the configPath.
	ClusterctlConfigPath string
	// Scheme is the GVK scheme to use for the tests
	Scheme *runtime.Scheme
}

// InitSchemeFunc is a function that will create a scheme
type InitSchemeFunc func() *runtime.Scheme

// WithSchemeInit will set a different function to initalize the scheme
func WithSchemeInit(fn InitSchemeFunc) Option {
	return func(ctx *E2EContext) {
		ctx.Environment.Scheme = fn()
	}
}

// WithSchemeInit will set a different function to initalize the scheme
func WithManaged() Option {
	return func(ctx *E2EContext) {
		ctx.IsManaged = true
	}
}

func (c *E2EContext) InfrastructureProviders() []string {
	InfraProviders := []string{}
	for _, provider := range c.E2EConfig.Providers {
		if provider.Type == string(clusterctlv1.InfrastructureProviderType) {
			InfraProviders = append(InfraProviders, provider.Name)
		}
	}
	return InfraProviders
}

func (c *E2EContext) BootstrapProviders() []string {
	BootstrapProviders := []string{}
	for _, provider := range c.E2EConfig.Providers {
		if provider.Type == string(clusterctlv1.BootstrapProviderType) {
			BootstrapProviders = append(BootstrapProviders, provider.Name)
		}
	}
	return BootstrapProviders
}

func (c *E2EContext) ControlPlaneProviders() []string {
	ControlPlaneProviders := []string{}
	for _, provider := range c.E2EConfig.Providers {
		if provider.Type == string(clusterctlv1.ControlPlaneProviderType) {
			ControlPlaneProviders = append(ControlPlaneProviders, provider.Name)
		}
	}
	return ControlPlaneProviders
}
