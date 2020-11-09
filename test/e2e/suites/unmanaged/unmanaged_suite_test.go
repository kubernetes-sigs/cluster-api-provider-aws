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

package unmanaged

import (
	"context"
	"flag"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	cfn_bootstrap "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cloudformation/bootstrap"
	cloudformation "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cloudformation/service"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/credentials"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/yaml"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/test/framework/bootstrap"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
	"sigs.k8s.io/cluster-api/test/framework/kubernetesversions"

	"sigs.k8s.io/cluster-api-provider-aws/test/e2e/shared"
)

var (
	e2eCtx *shared.E2EContext
)

func init() {
	e2eCtx = shared.NewE2EContextWithFlags()
}

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

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "capa-e2e", []Reporter{framework.CreateJUnitReporterForProw(e2eCtx.ArtifactFolder)})
}

var _ = SynchronizedBeforeSuite(func() []byte {
	flag.Parse()
	Expect(e2eCtx.ConfigPath).To(BeAnExistingFile(), "Invalid test suite argument. configPath should be an existing file.")
	Expect(os.MkdirAll(e2eCtx.ArtifactFolder, 0o750)).To(Succeed(), "Invalid test suite argument. Can't create artifacts-folder %q", e2eCtx.ArtifactFolder)
	shared.Byf("Loading the e2e test configuration from %q", e2eCtx.ConfigPath)
	e2eCtx.E2EConfig = loadE2EConfig(e2eCtx.ConfigPath)
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
	e2eCtx.AWSSession = newAWSSession()
	if !e2eCtx.SkipCloudFormationCreation {
		createCloudFormationStack(e2eCtx.AWSSession, getBootstrapTemplate())
	}
	ensureNoServiceLinkedRoles(e2eCtx.AWSSession)
	ensureSSHKeyPair(e2eCtx.AWSSession, shared.DefaultSSHKeyPairName)
	e2eCtx.BootstrapAccessKey = newUserAccessKey(e2eCtx.AWSSession, getBootstrapTemplate().Spec.BootstrapUser.UserName)

	By("Initializing a runtime.Scheme with all the GVK relevant for this test")
	scheme := initScheme()

	// If using a version of Kubernetes from CI, override the image ID with a known good image
	if e2eCtx.UseCIArtifacts {
		e2eCtx.E2EConfig.Variables["IMAGE_ID"] = conformanceImageID()
	}

	shared.Byf("Creating a clusterctl local repository into %q", e2eCtx.ArtifactFolder)
	e2eCtx.ClusterctlConfigPath = createClusterctlLocalRepository(e2eCtx.E2EConfig, filepath.Join(e2eCtx.ArtifactFolder, "repository"))

	By("Setting up the bootstrap cluster")
	e2eCtx.BootstrapClusterProvider, e2eCtx.BootstrapClusterProxy = setupBootstrapCluster(e2eCtx.E2EConfig, scheme, e2eCtx.UseExistingCluster)

	setEnvVar("AWS_B64ENCODED_CREDENTIALS", encodeCredentials(e2eCtx.BootstrapAccessKey, getBootstrapTemplate().Spec.Region), true)

	By("Initializing the bootstrap cluster")
	initBootstrapCluster(e2eCtx.BootstrapClusterProxy, e2eCtx.E2EConfig, e2eCtx.ClusterctlConfigPath, e2eCtx.ArtifactFolder)

	conf := synchronizedBeforeTestSuiteConfig{
		ArtifactFolder:          e2eCtx.ArtifactFolder,
		ConfigPath:              e2eCtx.ConfigPath,
		ClusterctlConfigPath:    e2eCtx.ClusterctlConfigPath,
		KubeconfigPath:          e2eCtx.BootstrapClusterProxy.GetKubeconfigPath(),
		Region:                  getBootstrapTemplate().Spec.Region,
		E2EConfig:               *e2eCtx.E2EConfig,
		KubetestConfigFilePath:  e2eCtx.KubetestConfigFilePath,
		UseCIArtifacts:          e2eCtx.UseCIArtifacts,
		GinkgoNodes:             e2eCtx.GinkgoNodes,
		GinkgoSlowSpecThreshold: e2eCtx.GinkgoSlowSpecThreshold,
	}

	data, err := yaml.Marshal(conf)
	Expect(err).NotTo(HaveOccurred())
	return data

}, func(data []byte) {
	// Before each ParallelNode.
	conf := &synchronizedBeforeTestSuiteConfig{}
	err := yaml.UnmarshalStrict(data, conf)
	Expect(err).NotTo(HaveOccurred())
	e2eCtx.ArtifactFolder = conf.ArtifactFolder
	e2eCtx.ConfigPath = conf.ConfigPath
	e2eCtx.ClusterctlConfigPath = conf.ClusterctlConfigPath
	e2eCtx.BootstrapClusterProxy = framework.NewClusterProxy("bootstrap", conf.KubeconfigPath, initScheme())
	e2eCtx.E2EConfig = &conf.E2EConfig
	e2eCtx.KubetestConfigFilePath = conf.KubetestConfigFilePath
	e2eCtx.UseCIArtifacts = conf.UseCIArtifacts
	e2eCtx.GinkgoNodes = conf.GinkgoNodes
	e2eCtx.GinkgoSlowSpecThreshold = conf.GinkgoSlowSpecThreshold
	azs := getAvailabilityZones()
	setEnvVar(shared.AwsAvailabilityZone1, *azs[0].ZoneName, false)
	setEnvVar(shared.AwsAvailabilityZone2, *azs[1].ZoneName, false)
	setEnvVar("AWS_REGION", conf.Region, false)
	setEnvVar("AWS_SSH_KEY_NAME", shared.DefaultSSHKeyPairName, false)
	e2eCtx.AWSSession = newAWSSession()
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
					shared.DumpSpecResources(resourceCtx, e2eCtx, k)
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
					shared.DumpMachines(machineCtx, e2eCtx, k)
				}
			}
		}
	}()
})

// Using a SynchronizedAfterSuite for controlling how to delete resources shared across ParallelNodes (~ginkgo threads).
// The bootstrap cluster is shared across all the tests, so it should be deleted only after all ParallelNodes completes.
// The local clusterctl repository is preserved like everything else created into the artifact folder.
var _ = SynchronizedAfterSuite(func() {
	if e2eCtx.ResourceTickerDone != nil {
		e2eCtx.ResourceTickerDone <- true
	}
	if e2eCtx.MachineTickerDone != nil {
		e2eCtx.MachineTickerDone <- true
	}
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Minute)
	defer cancel()
	for k := range e2eCtx.Namespaces {
		shared.DumpSpecResourcesAndCleanup(ctx, "", k, e2eCtx)
		shared.DumpMachines(ctx, e2eCtx, k)
	}
}, func() {
	// After all ParallelNodes.
	By("Tearing down the management cluster")
	if !e2eCtx.SkipCleanup {
		tearDown(e2eCtx.BootstrapClusterProvider, e2eCtx.BootstrapClusterProxy)
		if !e2eCtx.SkipCloudFormationDeletion {
			deleteCloudFormationStack(e2eCtx.AWSSession, getBootstrapTemplate())
		}
	}
})

// initScheme creates a new GVK scheme
func initScheme() *runtime.Scheme {
	sc := runtime.NewScheme()
	framework.TryAddDefaultSchemes(sc)
	_ = v1alpha3.AddToScheme(sc)
	_ = clientgoscheme.AddToScheme(sc)
	return sc
}

// loadE2EConfig loads the e2econfig from the specified path
func loadE2EConfig(configPath string) *clusterctl.E2EConfig {
	config := clusterctl.LoadE2EConfig(context.TODO(), clusterctl.LoadE2EConfigInput{ConfigPath: configPath})
	Expect(config).ToNot(BeNil(), "Failed to load E2E config from %s", configPath)
	// Read CNI file and set CNI_RESOURCES environmental variable
	Expect(config.Variables).To(HaveKey(shared.CNIPath), "Missing %s variable in the config", shared.CNIPath)
	clusterctl.SetCNIEnvVar(config.GetVariable(shared.CNIPath), shared.CNIResources)
	return config
}

// createClusterctlLocalRepository generates a clusterctl repository.
// Must always be run after kubetest.NewConfiguration
func createClusterctlLocalRepository(config *clusterctl.E2EConfig, repositoryFolder string) string {
	clusterctlConfig := clusterctl.CreateRepository(context.TODO(), clusterctl.CreateRepositoryInput{
		E2EConfig:        config,
		RepositoryFolder: repositoryFolder,
	})
	Expect(clusterctlConfig).To(BeAnExistingFile(), "The clusterctl config file does not exists in the local repository %s", repositoryFolder)
	return clusterctlConfig
}

// setupBootstrapCluster installs Cluster API components via clusterctl
func setupBootstrapCluster(config *clusterctl.E2EConfig, scheme *runtime.Scheme, useExistingCluster bool) (bootstrap.ClusterProvider, framework.ClusterProxy) {
	var clusterProvider bootstrap.ClusterProvider
	kubeconfigPath := ""
	if !useExistingCluster {
		clusterProvider = bootstrap.CreateKindBootstrapClusterAndLoadImages(context.TODO(), bootstrap.CreateKindBootstrapClusterAndLoadImagesInput{
			Name:               config.ManagementClusterName,
			RequiresDockerSock: config.HasDockerProvider(),
			Images:             config.Images,
		})
		Expect(clusterProvider).ToNot(BeNil(), "Failed to create a bootstrap cluster")

		kubeconfigPath = clusterProvider.GetKubeconfigPath()
		Expect(kubeconfigPath).To(BeAnExistingFile(), "Failed to get the kubeconfig file for the bootstrap cluster")
	}

	clusterProxy := framework.NewClusterProxy("bootstrap", kubeconfigPath, scheme)
	Expect(clusterProxy).ToNot(BeNil(), "Failed to get a bootstrap cluster proxy")

	return clusterProvider, clusterProxy
}

// initBootstrapCluster uses kind to create a cluster
func initBootstrapCluster(bootstrapClusterProxy framework.ClusterProxy, config *clusterctl.E2EConfig, clusterctlConfig, artifactFolder string) {
	clusterctl.InitManagementClusterAndWatchControllerLogs(context.TODO(), clusterctl.InitManagementClusterAndWatchControllerLogsInput{
		ClusterProxy:            bootstrapClusterProxy,
		ClusterctlConfigPath:    clusterctlConfig,
		InfrastructureProviders: config.InfrastructureProviders(),
		LogFolder:               filepath.Join(e2eCtx.ArtifactFolder, "clusters", bootstrapClusterProxy.GetName()),
	}, config.GetIntervals(bootstrapClusterProxy.GetName(), "wait-controllers")...)
}

// tearDown the bootstrap kind cluster
func tearDown(bootstrapClusterProvider bootstrap.ClusterProvider, bootstrapClusterProxy framework.ClusterProxy) {
	if bootstrapClusterProxy != nil {
		bootstrapClusterProxy.Dispose(context.TODO())
	}
	if bootstrapClusterProvider != nil {
		bootstrapClusterProvider.Dispose(context.TODO())
	}
}

// newBootstrapTemplate generates a clusterawsadm configuration, and prints it
// and the resultant cloudformation template to the artifacts directory
func newBootstrapTemplate() *cfn_bootstrap.Template {
	By("Creating a bootstrap AWSIAMConfiguration")
	t := cfn_bootstrap.NewTemplate()
	t.Spec.BootstrapUser.Enable = true
	t.Spec.SecureSecretsBackends = []v1alpha3.SecretBackend{
		v1alpha3.SecretBackendSecretsManager,
		v1alpha3.SecretBackendSSMParameterStore,
	}
	region, err := credentials.ResolveRegion("")
	Expect(err).NotTo(HaveOccurred())
	t.Spec.Region = region
	str, err := yaml.Marshal(t.Spec)
	Expect(err).NotTo(HaveOccurred())
	Expect(ioutil.WriteFile(path.Join(e2eCtx.ArtifactFolder, "awsiamconfiguration.yaml"), str, 0644)).To(Succeed())
	cfnData, err := t.RenderCloudFormation().YAML()
	Expect(err).NotTo(HaveOccurred())
	Expect(ioutil.WriteFile(path.Join(e2eCtx.ArtifactFolder, "cloudformation.yaml"), cfnData, 0644)).To(Succeed())
	return &t
}

// getBootstrapTemplate gets or generates a new bootstrap template
func getBootstrapTemplate() *cfn_bootstrap.Template {
	if e2eCtx.BootstrapTemplate == nil {
		e2eCtx.BootstrapTemplate = newBootstrapTemplate()
	}
	return e2eCtx.BootstrapTemplate
}

func newAWSSession() client.ConfigProvider {
	By("Getting an AWS IAM session")
	region, err := credentials.ResolveRegion("")
	Expect(err).NotTo(HaveOccurred())
	config := aws.NewConfig().WithCredentialsChainVerboseErrors(true).WithRegion(region)
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            *config,
	})
	Expect(err).NotTo(HaveOccurred())
	_, err = sess.Config.Credentials.Get()
	Expect(err).NotTo(HaveOccurred())
	return sess
}

// ensureNoServiceLinkedRoles removes an auto-created IAM role, and tests
// the controller's IAM permissions to use ELB and Spot instances successfully
func ensureNoServiceLinkedRoles(prov client.ConfigProvider) {
	shared.Byf("Deleting AWS IAM Service Linked Role: role-name=AWSServiceRoleForElasticLoadBalancing")
	iamSvc := iam.New(prov)
	_, err := iamSvc.DeleteServiceLinkedRole(&iam.DeleteServiceLinkedRoleInput{
		RoleName: aws.String("AWSServiceRoleForElasticLoadBalancing"),
	})
	if code, _ := awserrors.Code(err); code != iam.ErrCodeNoSuchEntityException {
		Expect(err).NotTo(HaveOccurred())
	}

	shared.Byf("Deleting AWS IAM Service Linked Role: role-name=AWSServiceRoleForEC2Spot")
	_, err = iamSvc.DeleteServiceLinkedRole(&iam.DeleteServiceLinkedRoleInput{
		RoleName: aws.String("AWSServiceRoleForEC2Spot"),
	})
	if code, _ := awserrors.Code(err); code != iam.ErrCodeNoSuchEntityException {
		Expect(err).NotTo(HaveOccurred())
	}
}

// createCloudFormationStack ensures the cloudformation stack is up to date
func createCloudFormationStack(prov client.ConfigProvider, t *cfn_bootstrap.Template) {
	shared.Byf("Creating AWS CloudFormation stack for AWS IAM resources: stack-name=%s", t.Spec.StackName)
	cfnSvc := cloudformation.NewService(cfn.New(prov))
	cfnTemplate := t.RenderCloudFormation()
	Expect(
		cfnSvc.ReconcileBootstrapStack(t.Spec.StackName, *cfnTemplate),
	).To(Succeed())
}

// deleteCloudFormationStack removes the provisioned clusterawsadm stack
func deleteCloudFormationStack(prov client.ConfigProvider, t *cfn_bootstrap.Template) {
	shared.Byf("Deleting %s CloudFormation stack", t.Spec.StackName)
	cfnSvc := cloudformation.NewService(cfn.New(prov))
	Expect(
		cfnSvc.DeleteStack(t.Spec.StackName),
	).To(Succeed())
}

// ensureSSHKeyPair ensures A SSH key is present under the name
func ensureSSHKeyPair(prov client.ConfigProvider, keyPairName string) {
	shared.Byf("Ensuring presence of SSH key in EC2: key-name=%s", keyPairName)
	ec2c := ec2.New(prov)
	_, err := ec2c.CreateKeyPair(&ec2.CreateKeyPairInput{KeyName: aws.String(keyPairName)})
	if code, _ := awserrors.Code(err); code != "InvalidKeyPair.Duplicate" {
		Expect(err).NotTo(HaveOccurred())
	}
}

// encodeCredentials leverages clusterawsadm to encode AWS credentials
func encodeCredentials(accessKey *iam.AccessKey, region string) string {
	creds := credentials.AWSCredentials{
		Region:          region,
		AccessKeyID:     *accessKey.AccessKeyId,
		SecretAccessKey: *accessKey.SecretAccessKey,
	}
	encCreds, err := creds.RenderBase64EncodedAWSDefaultProfile()
	Expect(err).NotTo(HaveOccurred())
	return encCreds
}

// setEnvVar sets an environment variable in the process. If marked private,
// the value is not printed.
func setEnvVar(key, value string, private bool) {
	printableValue := "*******"
	if !private {
		printableValue = value
	}

	shared.Byf("Setting environment variable: key=%s, value=%s", key, printableValue)
	os.Setenv(key, value)
}

// newUserAccessKey generates a new AWS Access Key pair based off of the
// bootstrap user. This tests that the CloudFormation policy is correct.
func newUserAccessKey(prov client.ConfigProvider, userName string) *iam.AccessKey {
	iamSvc := iam.New(prov)
	keyOuts, _ := iamSvc.ListAccessKeys(&iam.ListAccessKeysInput{
		UserName: aws.String(userName),
	})
	for i := range keyOuts.AccessKeyMetadata {
		shared.Byf("Deleting an existing access key: user-name=%s", userName)
		_, err := iamSvc.DeleteAccessKey(&iam.DeleteAccessKeyInput{
			UserName:    aws.String(userName),
			AccessKeyId: keyOuts.AccessKeyMetadata[i].AccessKeyId,
		})
		Expect(err).NotTo(HaveOccurred())
	}
	shared.Byf("Creating an access key: user-name=%s", userName)
	out, err := iamSvc.CreateAccessKey(&iam.CreateAccessKeyInput{UserName: aws.String(userName)})
	Expect(err).NotTo(HaveOccurred())
	Expect(out.AccessKey).ToNot(BeNil())

	return &iam.AccessKey{
		AccessKeyId:     out.AccessKey.AccessKeyId,
		SecretAccessKey: out.AccessKey.SecretAccessKey,
	}
}

// conformanceImageID looks up a specific image for a given
// Kubernetes version in the e2econfig
func conformanceImageID() string {
	ver := e2eCtx.E2EConfig.GetVariable("CONFORMANCE_CI_ARTIFACTS_KUBERNETES_VERSION")
	strippedVer := strings.Replace(ver, "v", "", 1)
	amiName := shared.AMIPrefix + strippedVer + "*"

	shared.Byf("Searching for AMI: name=%s", amiName)
	ec2Svc := ec2.New(e2eCtx.AWSSession)
	filters := []*ec2.Filter{
		{
			Name:   aws.String("name"),
			Values: []*string{aws.String(amiName)},
		},
	}
	filters = append(filters, &ec2.Filter{
		Name:   aws.String("owner-id"),
		Values: []*string{aws.String(shared.DefaultImageLookupOrg)},
	})
	resp, err := ec2Svc.DescribeImages(&ec2.DescribeImagesInput{
		Filters: filters,
	})
	Expect(err).NotTo(HaveOccurred())
	Expect(len(resp.Images)).To(Not(BeZero()))
	imageID := aws.StringValue(resp.Images[0].ImageId)
	shared.Byf("Using AMI: image-id=%s", imageID)
	return imageID
}
