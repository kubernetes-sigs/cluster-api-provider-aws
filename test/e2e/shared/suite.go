//go:build e2e
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

// Package shared provides common utilities, setup and teardown for the e2e tests.
package shared

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/gofrs/flock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/yaml"

	"sigs.k8s.io/cluster-api-provider-aws/v2/test/helpers/kubernetesversions"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
)

type synchronizedBeforeTestSuiteConfig struct {
	ArtifactFolder           string               `json:"artifactFolder,omitempty"`
	ConfigPath               string               `json:"configPath,omitempty"`
	ClusterctlConfigPath     string               `json:"clusterctlConfigPath,omitempty"`
	KubeconfigPath           string               `json:"kubeconfigPath,omitempty"`
	Region                   string               `json:"region,omitempty"`
	E2EConfig                clusterctl.E2EConfig `json:"e2eConfig,omitempty"`
	BootstrapAccessKey       *iam.AccessKey       `json:"bootstrapAccessKey,omitempty"`
	KubetestConfigFilePath   string               `json:"kubetestConfigFilePath,omitempty"`
	UseCIArtifacts           bool                 `json:"useCIArtifacts,omitempty"`
	GinkgoNodes              int                  `json:"ginkgoNodes,omitempty"`
	GinkgoSlowSpecThreshold  int                  `json:"ginkgoSlowSpecThreshold,omitempty"`
	Base64EncodedCredentials string               `json:"base64EncodedCredentials,omitempty"`
}

// Node1BeforeSuite is the common setup down on the first ginkgo node before the test suite runs.
func Node1BeforeSuite(e2eCtx *E2EContext) []byte {
	flag.Parse()
	Expect(e2eCtx.Settings.ConfigPath).To(BeAnExistingFile(), "Invalid test suite argument. configPath should be an existing file.")
	Expect(os.MkdirAll(e2eCtx.Settings.ArtifactFolder, 0o750)).To(Succeed(), "Invalid test suite argument. Can't create artifacts-folder %q", e2eCtx.Settings.ArtifactFolder)
	By(fmt.Sprintf("Loading the e2e test configuration from %q", e2eCtx.Settings.ConfigPath))
	e2eCtx.E2EConfig = LoadE2EConfig(e2eCtx.Settings.ConfigPath)
	sourceTemplate, err := os.ReadFile(filepath.Join(e2eCtx.Settings.DataFolder, e2eCtx.Settings.SourceTemplate))
	Expect(err).NotTo(HaveOccurred())
	e2eCtx.StartOfSuite = time.Now()

	var clusterctlCITemplate clusterctl.Files
	if !e2eCtx.IsManaged {
		// Create CI manifest for upgrading to Kubernetes main test
		platformKustomization, err := os.ReadFile(filepath.Join(e2eCtx.Settings.DataFolder, "ci-artifacts-platform-kustomization-for-upgrade.yaml"))
		Expect(err).NotTo(HaveOccurred())
		sourceTemplateForUpgrade, err := os.ReadFile(filepath.Join(e2eCtx.Settings.DataFolder, "infrastructure-aws/withoutclusterclass/generated/cluster-template-upgrade-to-main.yaml"))
		Expect(err).NotTo(HaveOccurred())

		ciTemplateForUpgradePath, err := kubernetesversions.GenerateCIArtifactsInjectedTemplateForDebian(
			kubernetesversions.GenerateCIArtifactsInjectedTemplateForDebianInput{
				ArtifactsDirectory:    e2eCtx.Settings.ArtifactFolder,
				SourceTemplate:        sourceTemplateForUpgrade,
				PlatformKustomization: platformKustomization,
			},
		)
		Expect(err).NotTo(HaveOccurred())

		ciTemplateForUpgradeName := "cluster-template-upgrade-ci-artifacts.yaml"
		templateDir := path.Join(e2eCtx.Settings.ArtifactFolder, "templates")
		newTemplatePath := templateDir + "/" + ciTemplateForUpgradeName

		err = exec.Command("cp", ciTemplateForUpgradePath, newTemplatePath).Run() //nolint:gosec
		Expect(err).NotTo(HaveOccurred())

		clusterctlCITemplateForUpgrade := clusterctl.Files{
			SourcePath: newTemplatePath,
			TargetName: ciTemplateForUpgradeName,
		}

		// Create CI manifest for conformance test
		platformKustomization, err = os.ReadFile(filepath.Join(e2eCtx.Settings.DataFolder, "ci-artifacts-platform-kustomization.yaml"))
		Expect(err).NotTo(HaveOccurred())
		ciTemplatePath, err := kubernetesversions.GenerateCIArtifactsInjectedTemplateForDebian(
			kubernetesversions.GenerateCIArtifactsInjectedTemplateForDebianInput{
				ArtifactsDirectory:    e2eCtx.Settings.ArtifactFolder,
				SourceTemplate:        sourceTemplate,
				PlatformKustomization: platformKustomization,
			},
		)
		Expect(err).NotTo(HaveOccurred())

		clusterctlCITemplate = clusterctl.Files{
			SourcePath: ciTemplatePath,
			TargetName: "cluster-template-conformance-ci-artifacts.yaml",
		}

		providers := e2eCtx.E2EConfig.Providers
		for i, prov := range providers {
			if prov.Name != "aws" {
				continue
			}
			e2eCtx.E2EConfig.Providers[i].Files = append(
				e2eCtx.E2EConfig.Providers[i].Files,
				clusterctlCITemplate,
				clusterctlCITemplateForUpgrade,
			)
		}
	}

	Expect(err).NotTo(HaveOccurred())
	e2eCtx.AWSSession = NewAWSSession()

	logAccountDetails(e2eCtx.AWSSession)

	bootstrapTemplate := getBootstrapTemplate(e2eCtx)
	bootstrapTags := map[string]string{"capa-e2e-test": "true"}
	e2eCtx.CloudFormationTemplate = renderCustomCloudFormation(bootstrapTemplate)

	if !e2eCtx.Settings.SkipCloudFormationCreation {
		count := 0
		Eventually(func(gomega Gomega) bool {
			count++
			By(fmt.Sprintf("Trying to create CloudFormation stack... attempt %d", count))
			success := true
			if err := createCloudFormationStack(e2eCtx.AWSSession, bootstrapTemplate, bootstrapTags); err != nil {
				By(fmt.Sprintf("Failed to create CloudFormation stack in attempt %d: %s", count, err.Error()))
				deleteCloudFormationStack(e2eCtx.AWSSession, bootstrapTemplate)
				success = false
			}
			return success
		}, 45*time.Minute, 30*time.Second).Should(BeTrue(), "Should've eventually succeeded creating an AWS CloudFormation stack")
	}

	ensureStackTags(e2eCtx.AWSSession, bootstrapTemplate.Spec.StackName, bootstrapTags)
	ensureNoServiceLinkedRoles(e2eCtx.AWSSession)
	ensureSSHKeyPair(e2eCtx.AWSSession, DefaultSSHKeyPairName)
	e2eCtx.Environment.BootstrapAccessKey = newUserAccessKey(e2eCtx.AWSSession, bootstrapTemplate.Spec.BootstrapUser.UserName)
	e2eCtx.BootstrapUserAWSSession = NewAWSSessionWithKey(e2eCtx.Environment.BootstrapAccessKey)
	Expect(ensureTestImageUploaded(e2eCtx)).NotTo(HaveOccurred())

	// Image ID is needed when using a CI Kubernetes version. This is used in conformance test and upgrade to main test.
	if !e2eCtx.IsManaged {
		e2eCtx.E2EConfig.Variables["IMAGE_ID"] = conformanceImageID(e2eCtx)
	}

	By(fmt.Sprintf("Creating a clusterctl local repository into %q", e2eCtx.Settings.ArtifactFolder))
	e2eCtx.Environment.ClusterctlConfigPath = createClusterctlLocalRepository(e2eCtx, filepath.Join(e2eCtx.Settings.ArtifactFolder, "repository"))

	By("Setting up the bootstrap cluster")
	e2eCtx.Environment.BootstrapClusterProvider, e2eCtx.Environment.BootstrapClusterProxy = setupBootstrapCluster(e2eCtx.E2EConfig, e2eCtx.Environment.Scheme, e2eCtx.Settings.UseExistingCluster)

	base64EncodedCredentials := encodeCredentials(e2eCtx.Environment.BootstrapAccessKey, bootstrapTemplate.Spec.Region)
	SetEnvVar("AWS_B64ENCODED_CREDENTIALS", base64EncodedCredentials, true)

	if !e2eCtx.Settings.SkipQuotas {
		By("Writing AWS service quotas to a file for parallel tests")
		quotas, originalQuotas := EnsureServiceQuotas(e2eCtx.BootstrapUserAWSSession)
		WriteResourceQuotesToFile(ResourceQuotaFilePath, quotas)
		WriteResourceQuotesToFile(path.Join(e2eCtx.Settings.ArtifactFolder, "initial-resource-quotas.yaml"), quotas)
		WriteAWSResourceQuotesToFile(path.Join(e2eCtx.Settings.ArtifactFolder, "initial-aws-resource-quotas.yaml"), originalQuotas)
	}

	e2eCtx.Settings.InstanceVCPU, err = strconv.Atoi(e2eCtx.E2EConfig.GetVariable(InstanceVcpu))
	Expect(err).NotTo(HaveOccurred())

	By("Initializing the bootstrap cluster")
	initBootstrapCluster(e2eCtx)

	CreateAWSClusterControllerIdentity(e2eCtx.Environment.BootstrapClusterProxy.GetClient())

	if e2eCtx.IsManaged {
		By("Setting up AWS static credentials")
		SetupStaticCredentials(e2eCtx)
	}

	conf := synchronizedBeforeTestSuiteConfig{
		ArtifactFolder:           e2eCtx.Settings.ArtifactFolder,
		ConfigPath:               e2eCtx.Settings.ConfigPath,
		ClusterctlConfigPath:     e2eCtx.Environment.ClusterctlConfigPath,
		KubeconfigPath:           e2eCtx.Environment.BootstrapClusterProxy.GetKubeconfigPath(),
		Region:                   getBootstrapTemplate(e2eCtx).Spec.Region,
		E2EConfig:                *e2eCtx.E2EConfig,
		BootstrapAccessKey:       e2eCtx.Environment.BootstrapAccessKey,
		KubetestConfigFilePath:   e2eCtx.Settings.KubetestConfigFilePath,
		UseCIArtifacts:           e2eCtx.Settings.UseCIArtifacts,
		GinkgoNodes:              e2eCtx.Settings.GinkgoNodes,
		GinkgoSlowSpecThreshold:  e2eCtx.Settings.GinkgoSlowSpecThreshold,
		Base64EncodedCredentials: base64EncodedCredentials,
	}

	data, err := yaml.Marshal(conf)
	Expect(err).NotTo(HaveOccurred())
	return data
}

// AllNodesBeforeSuite is the common setup down on each ginkgo parallel node before the test suite runs.
func AllNodesBeforeSuite(e2eCtx *E2EContext, data []byte) {
	conf := &synchronizedBeforeTestSuiteConfig{}
	err := yaml.UnmarshalStrict(data, conf)
	Expect(err).NotTo(HaveOccurred())
	e2eCtx.Settings.ArtifactFolder = conf.ArtifactFolder
	e2eCtx.Settings.ConfigPath = conf.ConfigPath
	e2eCtx.Environment.ClusterctlConfigPath = conf.ClusterctlConfigPath
	e2eCtx.Environment.BootstrapClusterProxy = framework.NewClusterProxy("bootstrap", conf.KubeconfigPath, e2eCtx.Environment.Scheme)
	e2eCtx.E2EConfig = &conf.E2EConfig
	e2eCtx.BootstrapUserAWSSession = NewAWSSessionWithKey(conf.BootstrapAccessKey)
	e2eCtx.Settings.FileLock = flock.New(ResourceQuotaFilePath)
	e2eCtx.Settings.KubetestConfigFilePath = conf.KubetestConfigFilePath
	e2eCtx.Settings.UseCIArtifacts = conf.UseCIArtifacts
	e2eCtx.Settings.GinkgoNodes = conf.GinkgoNodes
	e2eCtx.Settings.GinkgoSlowSpecThreshold = conf.GinkgoSlowSpecThreshold
	e2eCtx.AWSSession = NewAWSSession()
	azs := GetAvailabilityZones(e2eCtx.AWSSession)
	SetEnvVar(AwsAvailabilityZone1, *azs[0].ZoneName, false)
	SetEnvVar(AwsAvailabilityZone2, *azs[1].ZoneName, false)
	SetEnvVar("AWS_REGION", conf.Region, false)
	SetEnvVar("AWS_SSH_KEY_NAME", DefaultSSHKeyPairName, false)
	SetEnvVar("AWS_B64ENCODED_CREDENTIALS", conf.Base64EncodedCredentials, true)
	e2eCtx.Environment.ResourceTicker = time.NewTicker(time.Second * 5)
	e2eCtx.Environment.ResourceTickerDone = make(chan bool)
	// Get EC2 logs every minute
	e2eCtx.Environment.MachineTicker = time.NewTicker(time.Second * 60)
	e2eCtx.Environment.MachineTickerDone = make(chan bool)
	resourceCtx, resourceCancel := context.WithCancel(context.Background())
	machineCtx, machineCancel := context.WithCancel(context.Background())
	// Dump resources every 5 seconds
	go func() {
		defer GinkgoRecover()
		for {
			select {
			case <-e2eCtx.Environment.ResourceTickerDone:
				resourceCancel()
				return
			case <-e2eCtx.Environment.ResourceTicker.C:
				for k := range e2eCtx.Environment.Namespaces {
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
			case <-e2eCtx.Environment.MachineTickerDone:
				machineCancel()
				return
			case <-e2eCtx.Environment.MachineTicker.C:
				for k := range e2eCtx.Environment.Namespaces {
					DumpMachines(machineCtx, e2eCtx, k)
				}
			}
		}
	}()
}

// Node1AfterSuite is cleanup that runs on the first ginkgo node after the test suite finishes.
func Node1AfterSuite(e2eCtx *E2EContext) {
	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Minute)
	DumpEKSClusters(ctx, e2eCtx)
	DumpCloudTrailEvents(e2eCtx)

	if e2eCtx.IsManaged {
		By("Deleting AWS static credentials")
		CleanupStaticCredentials(ctx, e2eCtx)
	}

	defer cancel()
	By("Tearing down the management cluster")
	if !e2eCtx.Settings.SkipCleanup {
		tearDown(e2eCtx.Environment.BootstrapClusterProvider, e2eCtx.Environment.BootstrapClusterProxy)
		if !e2eCtx.Settings.SkipCloudFormationDeletion {
			deleteCloudFormationStack(e2eCtx.AWSSession, getBootstrapTemplate(e2eCtx))
		}
	}
}

// AllNodesAfterSuite is cleanup that runs on all ginkgo parallel nodes after the test suite finishes.
func AllNodesAfterSuite(e2eCtx *E2EContext) {
	if e2eCtx.Environment.ResourceTickerDone != nil {
		e2eCtx.Environment.ResourceTickerDone <- true
	}
	if e2eCtx.Environment.MachineTickerDone != nil {
		e2eCtx.Environment.MachineTickerDone <- true
	}
	ctx, cancel := context.WithTimeout(context.TODO(), 45*time.Minute)
	defer cancel()
	for k := range e2eCtx.Environment.Namespaces {
		DumpSpecResourcesAndCleanup(ctx, "", k, e2eCtx)
		DumpMachines(ctx, e2eCtx, k)
	}
}
