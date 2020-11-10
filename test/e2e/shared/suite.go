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
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"sigs.k8s.io/yaml"

	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
	"sigs.k8s.io/cluster-api/test/framework/kubernetesversions"
)

type synchronizedBeforeTestSuiteConfig struct {
	ArtifactFolder          string               `json:"artifactFolder,omitempty"`
	ConfigPath              string               `json:"configPath,omitempty"`
	ClusterctlConfigPath    string               `json:"clusterctlConfigPath,omitempty"`
	KubeconfigPath          string               `json:"kubeconfigPath,omitempty"`
	Region                  string               `json:"region,omitempty"`
	E2EConfig               clusterctl.E2EConfig `json:"e2eConfig,omitempty"`
	KubetestConfigFilePath  string               `json:"kubetestConfigFilePath,omitempty"`
	UseCIArtifacts          bool                 `json:"useCIArtifacts,omitempty"`
	GinkgoNodes             int                  `json:"ginkgoNodes,omitempty"`
	GinkgoSlowSpecThreshold int                  `json:"ginkgoSlowSpecThreshold,omitempty"`
}

func Node1BeforeSuite(e2eCtx *E2EContext) []byte {
	flag.Parse()
	Expect(e2eCtx.ConfigPath).To(BeAnExistingFile(), "Invalid test suite argument. configPath should be an existing file.")
	Expect(os.MkdirAll(e2eCtx.ArtifactFolder, 0o750)).To(Succeed(), "Invalid test suite argument. Can't create artifacts-folder %q", e2eCtx.ArtifactFolder)
	Byf("Loading the e2e test configuration from %q", e2eCtx.ConfigPath)
	e2eCtx.E2EConfig = LoadE2EConfig(e2eCtx.ConfigPath)
	sourceTemplate, err := ioutil.ReadFile("../../data/infrastructure-aws/cluster-template.yaml")
	Expect(err).NotTo(HaveOccurred())
	platformKustomization, err := ioutil.ReadFile("../../data/ci-artifacts-platform-kustomization.yaml")
	Expect(err).NotTo(HaveOccurred())
	ciTemplate, err := kubernetesversions.GenerateCIArtifactsInjectedTemplateForDebian(
		kubernetesversions.GenerateCIArtifactsInjectedTemplateForDebianInput{
			ArtifactsDirectory:    e2eCtx.ArtifactFolder,
			SourceTemplate:        sourceTemplate,
			PlatformKustomization: platformKustomization,
		},
	)
	clusterctlCITemplate := clusterctl.Files{
		SourcePath: ciTemplate,
		TargetName: "cluster-template-conformance-ci-artifacts.yaml",
	}
	providers := e2eCtx.E2EConfig.Providers
	for i, prov := range providers {
		if prov.Name != "aws" {
			continue
		}
		e2eCtx.E2EConfig.Providers[i].Files = append(e2eCtx.E2EConfig.Providers[i].Files, clusterctlCITemplate)
	}
	Expect(err).NotTo(HaveOccurred())
	e2eCtx.AWSSession = NewAWSSession()
	boostrapTemplate := getBootstrapTemplate(e2eCtx)
	if !e2eCtx.SkipCloudFormationCreation {
		createCloudFormationStack(e2eCtx.AWSSession, boostrapTemplate)
	}
	ensureNoServiceLinkedRoles(e2eCtx.AWSSession)
	ensureSSHKeyPair(e2eCtx.AWSSession, DefaultSSHKeyPairName)
	e2eCtx.BootstrapAccessKey = newUserAccessKey(e2eCtx.AWSSession, boostrapTemplate.Spec.BootstrapUser.UserName)

	By("Initializing a runtime.Scheme with all the GVK relevant for this test")
	scheme := e2eCtx.InitScheme()

	// If using a version of Kubernetes from CI, override the image ID with a known good image
	if e2eCtx.UseCIArtifacts {
		e2eCtx.E2EConfig.Variables["IMAGE_ID"] = conformanceImageID(e2eCtx)
	}

	Byf("Creating a clusterctl local repository into %q", e2eCtx.ArtifactFolder)
	e2eCtx.ClusterctlConfigPath = createClusterctlLocalRepository(e2eCtx.E2EConfig, filepath.Join(e2eCtx.ArtifactFolder, "repository"))

	By("Setting up the bootstrap cluster")
	e2eCtx.BootstrapClusterProvider, e2eCtx.BootstrapClusterProxy = setupBootstrapCluster(e2eCtx.E2EConfig, scheme, e2eCtx.UseExistingCluster)

	SetEnvVar("AWS_B64ENCODED_CREDENTIALS", encodeCredentials(e2eCtx.BootstrapAccessKey, boostrapTemplate.Spec.Region), true)

	By("Initializing the bootstrap cluster")
	initBootstrapCluster(e2eCtx)

	conf := synchronizedBeforeTestSuiteConfig{
		ArtifactFolder:          e2eCtx.ArtifactFolder,
		ConfigPath:              e2eCtx.ConfigPath,
		ClusterctlConfigPath:    e2eCtx.ClusterctlConfigPath,
		KubeconfigPath:          e2eCtx.BootstrapClusterProxy.GetKubeconfigPath(),
		Region:                  getBootstrapTemplate(e2eCtx).Spec.Region,
		E2EConfig:               *e2eCtx.E2EConfig,
		KubetestConfigFilePath:  e2eCtx.KubetestConfigFilePath,
		UseCIArtifacts:          e2eCtx.UseCIArtifacts,
		GinkgoNodes:             e2eCtx.GinkgoNodes,
		GinkgoSlowSpecThreshold: e2eCtx.GinkgoSlowSpecThreshold,
	}

	data, err := yaml.Marshal(conf)
	Expect(err).NotTo(HaveOccurred())
	return data
}

func AllNodesBeforeSuite(e2eCtx *E2EContext, data []byte) {
	// Before each ParallelNode.
	conf := &synchronizedBeforeTestSuiteConfig{}
	err := yaml.UnmarshalStrict(data, conf)
	Expect(err).NotTo(HaveOccurred())
	e2eCtx.ArtifactFolder = conf.ArtifactFolder
	e2eCtx.ConfigPath = conf.ConfigPath
	e2eCtx.ClusterctlConfigPath = conf.ClusterctlConfigPath
	e2eCtx.BootstrapClusterProxy = framework.NewClusterProxy("bootstrap", conf.KubeconfigPath, e2eCtx.InitScheme())
	e2eCtx.E2EConfig = &conf.E2EConfig
	e2eCtx.KubetestConfigFilePath = conf.KubetestConfigFilePath
	e2eCtx.UseCIArtifacts = conf.UseCIArtifacts
	e2eCtx.GinkgoNodes = conf.GinkgoNodes
	e2eCtx.GinkgoSlowSpecThreshold = conf.GinkgoSlowSpecThreshold
	azs := GetAvailabilityZones(GetSession())
	SetEnvVar(AwsAvailabilityZone1, *azs[0].ZoneName, false)
	SetEnvVar(AwsAvailabilityZone2, *azs[1].ZoneName, false)
	SetEnvVar("AWS_REGION", conf.Region, false)
	SetEnvVar("AWS_SSH_KEY_NAME", DefaultSSHKeyPairName, false)
	e2eCtx.AWSSession = NewAWSSession()
	e2eCtx.ResourceTicker = time.NewTicker(time.Second * 5)
	e2eCtx.ResourceTickerDone = make(chan bool)
	// Get EC2 logs every minute
	e2eCtx.MachineTicker = time.NewTicker(time.Second * 60)
	e2eCtx.MachineTickerDone = make(chan bool)
	resourceCtx, resourceCancel := context.WithCancel(context.Background())
	machineCtx, machineCancel := context.WithCancel(context.Background())

	// Dump resources every 5 seconds
	go func() {
		defer GinkgoRecover()
		for {
			select {
			case <-e2eCtx.ResourceTickerDone:
				resourceCancel()
				return
			case <-e2eCtx.ResourceTicker.C:
				for k := range e2eCtx.Namespaces {
					DumpSpecResources(resourceCtx, e2eCtx, k)
				}
			}
		}
	}()

	// Dump machine logs every 60 seconds
	go func() {
		defer GinkgoRecover()
		for {
			select {
			case <-e2eCtx.MachineTickerDone:
				machineCancel()
				return
			case <-e2eCtx.MachineTicker.C:
				for k := range e2eCtx.Namespaces {
					DumpMachines(machineCtx, e2eCtx, k)
				}
			}
		}
	}()
}

func Node1AfterSuite(e2eCtx *E2EContext) {
	if e2eCtx.ResourceTickerDone != nil {
		e2eCtx.ResourceTickerDone <- true
	}
	if e2eCtx.MachineTickerDone != nil {
		e2eCtx.MachineTickerDone <- true
	}
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Minute)
	defer cancel()
	for k := range e2eCtx.Namespaces {
		DumpSpecResourcesAndCleanup(ctx, "", k, e2eCtx)
		DumpMachines(ctx, e2eCtx, k)
	}
}

func AllNodesAfterSuite(e2eCtx *E2EContext) {
	By("Tearing down the management cluster")
	if !e2eCtx.SkipCleanup {
		tearDown(e2eCtx.BootstrapClusterProvider, e2eCtx.BootstrapClusterProxy)
		if !e2eCtx.SkipCloudFormationDeletion {
			deleteCloudFormationStack(e2eCtx.AWSSession, getBootstrapTemplate(e2eCtx))
		}
	}
}
