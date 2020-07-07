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

package e2e_new

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	cfn_bootstrap "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cloudformation/bootstrap"
	cloudformation "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cloudformation/service"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/credentials"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/yaml"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/test/framework/bootstrap"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
)

// Test suite flags
var (
	// configPath is the path to the e2e config file.
	configPath string

	// useExistingCluster instructs the test to use the current cluster instead of creating a new one (default discovery rules apply).
	useExistingCluster bool

	// artifactFolder is the folder to store e2e test artifacts.
	artifactFolder string

	// skipCleanup prevents cleanup of test resources e.g. for debug purposes.
	skipCleanup bool

	// skipCloudFormationDeletion prevents the deletion of the AWS CloudFormation stack
	skipCloudFormationDeletion bool
)

// Test suite global vars
var (
	// e2eConfig to be used for this test, read from configPath.
	e2eConfig *clusterctl.E2EConfig

	// clusterctlConfigPath to be used for this test, created by generating a clusterctl local repository
	// with the providers specified in the configPath.
	clusterctlConfigPath string

	// bootstrapClusterProvider manages provisioning of the the bootstrap cluster to be used for the e2e tests.
	// Please note that provisioning will be skipped if use-existing-cluster is provided.
	bootstrapClusterProvider bootstrap.ClusterProvider

	// bootstrapClusterProxy allows to interact with the bootstrap cluster to be used for the e2e tests.
	bootstrapClusterProxy framework.ClusterProxy

	// bootstrapTemplate is the clusterawsadm bootstrap template for this run
	bootstrapTemplate *cfn_bootstrap.Template

	// awsSession is a holder for AWS credentials
	awsSession client.ConfigProvider

	// bootstrapAccessKey
	bootstrapAccessKey *iam.AccessKey

	defaultSSHKeyPairName = "cluster-api-provider-aws-sigs-k8s-io"
)

const (
	KubernetesVersion = "KUBERNETES_VERSION"
	CNIPath           = "CNI"
)

func init() {
	flag.StringVar(&configPath, "config-path", "", "path to the e2e config file")
	flag.StringVar(&artifactFolder, "artifacts-folder", "", "folder where e2e test artifact should be stored")
	flag.BoolVar(&useExistingCluster, "use-existing-cluster", false, "if true, the test uses the current cluster instead of creating a new one (default discovery rules apply)")
	flag.BoolVar(&skipCleanup, "skip-cleanup", false, "if true, the resource cleanup after tests will be skipped")
	flag.BoolVar(&skipCloudFormationDeletion, "skip-cloudformation-deletion", false, "if true, an AWS CloudFormation stack will not be deleted")
}

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	junitPath := path.Join(artifactFolder, fmt.Sprintf("junit.e2e_suite.%d.xml", config.GinkgoConfig.ParallelNode))
	junitReporter := reporters.NewJUnitReporter(junitPath)

	RunSpecsWithDefaultAndCustomReporters(t, "capa-e2e", []Reporter{junitReporter})
}

var _ = SynchronizedBeforeSuite(func() []byte {
	Expect(configPath).To(BeAnExistingFile(), "Invalid test suite argument. configPath should be an existing file.")
	Expect(os.MkdirAll(artifactFolder, 0o755)).To(Succeed(), "Invalid test suite argument. Can't create artifacts-folder %q", artifactFolder)

	Byf("Loading the e2e test configuration from %q", configPath)
	e2eConfig = loadE2EConfig(configPath)

	awsSession = newAWSSession()
	createCloudFormationStack(awsSession, getBootstrapTemplate())
	ensureNoServiceLinkedRoles(awsSession)
	ensureSSHKeyPair(awsSession, defaultSSHKeyPairName)
	bootstrapAccessKey = newUserAccessKey(awsSession, getBootstrapTemplate().Spec.BootstrapUser.UserName)

	By("Initializing a runtime.Scheme with all the GVK relevant for this test")
	scheme := initScheme()

	Byf("Creating a clusterctl local repository into %q", artifactFolder)
	clusterctlConfigPath = createClusterctlLocalRepository(e2eConfig, filepath.Join(artifactFolder, "repository"))

	By("Setting up the bootstrap cluster")
	bootstrapClusterProvider, bootstrapClusterProxy = setupBootstrapCluster(e2eConfig, scheme, useExistingCluster)

	setEnvVar("AWS_B64ENCODED_CREDENTIALS", encodeCredentials(bootstrapAccessKey, getBootstrapTemplate().Spec.Region), true)

	By("Initializing the bootstrap cluster")
	initBootstrapCluster(bootstrapClusterProxy, e2eConfig, clusterctlConfigPath, artifactFolder)
	return []byte(
		strings.Join([]string{
			artifactFolder,
			configPath,
			clusterctlConfigPath,
			bootstrapClusterProxy.GetKubeconfigPath(),
			getBootstrapTemplate().Spec.Region,
		}, ","),
	)
}, func(data []byte) {
	// Before each ParallelNode.

	parts := strings.Split(string(data), ",")
	Expect(parts).To(HaveLen(5))

	artifactFolder = parts[0]
	configPath = parts[1]
	clusterctlConfigPath = parts[2]
	kubeconfigPath := parts[3]
	setEnvVar("AWS_REGION", parts[4], false)
	setEnvVar("AWS_SSH_KEY_NAME", defaultSSHKeyPairName, false)
	awsSession = newAWSSession()
	e2eConfig = loadE2EConfig(configPath)
	bootstrapClusterProxy = framework.NewClusterProxy("bootstrap", kubeconfigPath, initScheme())
})

// Using a SynchronizedAfterSuite for controlling how to delete resources shared across ParallelNodes (~ginkgo threads).
// The bootstrap cluster is shared across all the tests, so it should be deleted only after all ParallelNodes completes.
// The local clusterctl repository is preserved like everything else created into the artifact folder.
var _ = SynchronizedAfterSuite(func() {
	// After each ParallelNode.
}, func() {
	// After all ParallelNodes.

	By("Tearing down the management cluster")
	if !skipCleanup {
		tearDown(bootstrapClusterProvider, bootstrapClusterProxy)
		if !skipCloudFormationDeletion {
			awsSession = newAWSSession()
			deleteCloudFormationStack(awsSession, getBootstrapTemplate())
		}
	}
})

func initScheme() *runtime.Scheme {
	sc := runtime.NewScheme()
	framework.TryAddDefaultSchemes(sc)
	_ = v1alpha3.AddToScheme(sc)
	return sc
}

func loadE2EConfig(configPath string) *clusterctl.E2EConfig {
	config := clusterctl.LoadE2EConfig(context.TODO(), clusterctl.LoadE2EConfigInput{ConfigPath: configPath})
	Expect(config).ToNot(BeNil(), "Failed to load E2E config from %s", configPath)
	return config
}

func createClusterctlLocalRepository(config *clusterctl.E2EConfig, repositoryFolder string) string {
	clusterctlConfig := clusterctl.CreateRepository(context.TODO(), clusterctl.CreateRepositoryInput{
		E2EConfig:        config,
		RepositoryFolder: repositoryFolder,
	})
	Expect(clusterctlConfig).To(BeAnExistingFile(), "The clusterctl config file does not exists in the local repository %s", repositoryFolder)
	return clusterctlConfig
}

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

func initBootstrapCluster(bootstrapClusterProxy framework.ClusterProxy, config *clusterctl.E2EConfig, clusterctlConfig, artifactFolder string) {
	clusterctl.InitManagementClusterAndWatchControllerLogs(context.TODO(), clusterctl.InitManagementClusterAndWatchControllerLogsInput{
		ClusterProxy:            bootstrapClusterProxy,
		ClusterctlConfigPath:    clusterctlConfig,
		InfrastructureProviders: config.InfrastructureProviders(),
		LogFolder:               filepath.Join(artifactFolder, "clusters", bootstrapClusterProxy.GetName()),
	}, config.GetIntervals(bootstrapClusterProxy.GetName(), "wait-controllers")...)
}

func tearDown(bootstrapClusterProvider bootstrap.ClusterProvider, bootstrapClusterProxy framework.ClusterProxy) {
	if bootstrapClusterProxy != nil {
		bootstrapClusterProxy.Dispose(context.TODO())
	}
	if bootstrapClusterProvider != nil {
		bootstrapClusterProvider.Dispose(context.TODO())
	}
}

func newBootstrapTemplate() *cfn_bootstrap.Template {
	By("Creating a bootstrap AWSIAMConfiguration")
	t := cfn_bootstrap.NewTemplate()
	t.Spec.BootstrapUser.Enable = true
	region, err := credentials.ResolveRegion("")
	Expect(err).NotTo(HaveOccurred())
	t.Spec.Region = region
	str, err := yaml.Marshal(t.Spec)
	Expect(err).NotTo(HaveOccurred())
	fmt.Printf("AWSIAMConfiguration created:\n %s\n", string(str))
	return &t
}

func getBootstrapTemplate() *cfn_bootstrap.Template {
	if bootstrapTemplate == nil {
		bootstrapTemplate = newBootstrapTemplate()
	}
	return bootstrapTemplate
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

func ensureNoServiceLinkedRoles(prov client.ConfigProvider) {
	Byf("Deleting AWS IAM Service Linked Role: role-name=AWSServiceRoleForElasticLoadBalancing")
	iamSvc := iam.New(prov)
	_, err := iamSvc.DeleteServiceLinkedRole(&iam.DeleteServiceLinkedRoleInput{
		RoleName: aws.String("AWSServiceRoleForElasticLoadBalancing"),
	})
	if code, _ := awserrors.Code(err); code != iam.ErrCodeNoSuchEntityException {
		Expect(err).NotTo(HaveOccurred())
	}
}

func createCloudFormationStack(prov client.ConfigProvider, t *cfn_bootstrap.Template) {
	Byf("Creating AWS CloudFormation stack for AWS IAM resources: stack-name=%s", t.Spec.StackName)
	cfnSvc := cloudformation.NewService(cfn.New(prov))
	cfnTemplate := t.RenderCloudFormation()
	Expect(
		cfnSvc.ReconcileBootstrapStack(t.Spec.StackName, *cfnTemplate),
	).To(Succeed())
}

func deleteCloudFormationStack(prov client.ConfigProvider, t *cfn_bootstrap.Template) {
	Byf("Deleting %s CloudFormation stack", t.Spec.StackName)
	cfnSvc := cloudformation.NewService(cfn.New(prov))
	Expect(
		cfnSvc.DeleteStack(t.Spec.StackName),
	).To(Succeed())
}

func ensureSSHKeyPair(prov client.ConfigProvider, keyPairName string) {
	Byf("Ensuring presence of SSH key in EC2: key-name=%s", keyPairName)
	ec2c := ec2.New(prov)
	_, err := ec2c.CreateKeyPair(&ec2.CreateKeyPairInput{KeyName: aws.String(keyPairName)})
	if code, _ := awserrors.Code(err); code != "InvalidKeyPair.Duplicate" {
		Expect(err).NotTo(HaveOccurred())
	}
}

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

func setEnvVar(key, value string, private bool) {
	printableValue := "*******"
	if !private {
		printableValue = value
	}

	Byf("Setting environment variable: key=%s, value=%s", key, printableValue)
	os.Setenv(key, value)
}

func newUserAccessKey(prov client.ConfigProvider, userName string) *iam.AccessKey {
	iamSvc := iam.New(prov)
	keyOuts, err := iamSvc.ListAccessKeys(&iam.ListAccessKeysInput{
		UserName: aws.String(userName),
	})
	for i := range keyOuts.AccessKeyMetadata {
		Byf("Deleting an existing access key: user-name=%s", userName)
		_, err := iamSvc.DeleteAccessKey(&iam.DeleteAccessKeyInput{
			UserName:    aws.String(userName),
			AccessKeyId: keyOuts.AccessKeyMetadata[i].AccessKeyId,
		})
		Expect(err).NotTo(HaveOccurred())
	}
	Byf("Creating an access key: user-name=%s", userName)
	out, err := iamSvc.CreateAccessKey(&iam.CreateAccessKeyInput{UserName: aws.String(userName)})
	Expect(err).NotTo(HaveOccurred())
	Expect(out.AccessKey).ToNot(BeNil())

	return &iam.AccessKey{
		AccessKeyId:     out.AccessKey.AccessKeyId,
		SecretAccessKey: out.AccessKey.SecretAccessKey,
	}
}
