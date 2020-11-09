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
	corev1 "k8s.io/api/core/v1"

	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/test/framework/bootstrap"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"

	cfn_bootstrap "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cloudformation/bootstrap"
)

func NewE2EContextWithFlags() *E2EContext {
	ctx := &E2EContext{
		Namespaces: map[*corev1.Namespace]context.CancelFunc{},
	}
	CreateDefaultFlags(ctx)

	return ctx
}

type E2EContext struct {
	// Flags

	// ConfigPath is the path to the e2e config file.
	ConfigPath string
	// useExistingCluster instructs the test to use the current cluster instead of creating a new one (default discovery rules apply).
	UseExistingCluster bool
	// ArtifactFolder is the folder to store e2e test artifacts.
	ArtifactFolder string
	// DataFolder is the root folder for the data required by the tests
	//DataFolder string
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

	//Globals

	// e2eConfig to be used for this test, read from configPath.
	E2EConfig *clusterctl.E2EConfig
	// clusterctlConfigPath to be used for this test, created by generating a clusterctl local repository
	// with the providers specified in the configPath.
	ClusterctlConfigPath string
	// bootstrapClusterProvider manages provisioning of the the bootstrap cluster to be used for the e2e tests.
	// Please note that provisioning will be skipped if use-existing-cluster is provided.
	BootstrapClusterProvider bootstrap.ClusterProvider
	// bootstrapClusterProxy allows to interact with the bootstrap cluster to be used for the e2e tests.
	BootstrapClusterProxy framework.ClusterProxy
	// bootstrapTemplate is the clusterawsadm bootstrap template for this run
	BootstrapTemplate *cfn_bootstrap.Template
	// bootstrapAccessKey
	BootstrapAccessKey *iam.AccessKey
	// ticker for dumping resources
	ResourceTicker *time.Ticker
	// tickerDone to stop ticking
	ResourceTickerDone chan bool
	// ticker for dumping resources
	MachineTicker *time.Ticker
	// tickerDone to stop ticking
	MachineTickerDone chan bool

	Namespaces map[*corev1.Namespace]context.CancelFunc

	AWSSession client.ConfigProvider
}
