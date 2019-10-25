// +build e2e

/*
Copyright 2018 The Kubernetes Authors.

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

package e2e_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"testing"
	"text/template"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	awssts "github.com/aws/aws-sdk-go/service/sts"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	bootstrapv1 "sigs.k8s.io/cluster-api-bootstrap-provider-kubeadm/api/v1alpha2"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha2"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/cloudformation"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/sts"
	common "sigs.k8s.io/cluster-api/test/helpers/components"
	capiFlag "sigs.k8s.io/cluster-api/test/helpers/flag"
	"sigs.k8s.io/cluster-api/test/helpers/kind"
	"sigs.k8s.io/cluster-api/test/helpers/scheme"
	"sigs.k8s.io/cluster-api/util"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func TestE2e(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "e2e Suite")
}

const (
	CAPI_VERSION  = "v0.2.2"
	CABPK_VERSION = "v0.1.0"

	capiNamespace       = "capi-system"
	capiDeploymentName  = "capi-controller-manager"
	cabpkNamespace      = "cabpk-system"
	cabpkDeploymentName = "cabpk-controller-manager"
	capaNamespace       = "capa-system"
	capaDeploymentName  = "capa-controller-manager"
	setupTimeout        = 10 * 60
	stackName           = "cluster-api-provider-aws-sigs-k8s-io"
	keyPairName         = "cluster-api-provider-aws-sigs-k8s-io"
)

var (
	managerImage    = capiFlag.DefineOrLookupStringFlag("managerImage", "", "Docker image to load into the kind cluster for testing")
	cabpkComponents = capiFlag.DefineOrLookupStringFlag("cabpkComponents", "https://github.com/kubernetes-sigs/cluster-api-bootstrap-provider-kubeadm/releases/download/"+CABPK_VERSION+"/bootstrap-components.yaml", "URL to CAPI components to load")
	capaComponents  = capiFlag.DefineOrLookupStringFlag("capaComponents", "", "capa components to load")
	kustomizeBinary = capiFlag.DefineOrLookupStringFlag("kustomizeBinary", "kustomize", "path to the kustomize binary")

	kindCluster kind.Cluster
	kindClient  crclient.Client
	sess        client.ConfigProvider
	accountID   string
	accessKey   *iam.AccessKey
	suiteTmpDir string
	region      string
)

var _ = BeforeSuite(func() {
	fmt.Fprintf(GinkgoWriter, "Setting up kind cluster\n")

	var err error
	suiteTmpDir, err = ioutil.TempDir("", "capa-e2e-suite")
	Expect(err).NotTo(HaveOccurred())

	var ok bool
	region, ok = os.LookupEnv("AWS_REGION")
	fmt.Fprintf(GinkgoWriter, "Running in region: %s\n", region)
	if !ok {
		fmt.Fprintf(GinkgoWriter, "Environment variable AWS_REGION not found")
		Expect(ok).To(BeTrue())
	}

	sess = getSession()

	fmt.Fprintf(GinkgoWriter, "Creating AWS prerequisites\n")
	accountID = getAccountID(sess)
	createKeyPair(sess)
	createIAMRoles(sess, accountID)

	iamc := iam.New(sess)
	out, err := iamc.CreateAccessKey(&iam.CreateAccessKeyInput{UserName: aws.String("bootstrapper.cluster-api-provider-aws.sigs.k8s.io")})
	Expect(err).NotTo(HaveOccurred())
	Expect(out.AccessKey).NotTo(BeNil())
	accessKey = out.AccessKey

	kindCluster = kind.Cluster{
		Name: "capa-test-" + util.RandomString(6),
	}
	kindCluster.Setup()
	loadManagerImage(kindCluster)

	kindClient, err = crclient.New(kindCluster.RestConfig(), crclient.Options{Scheme: setupScheme()})
	Expect(err).NotTo(HaveOccurred())

	// Deploy the CAPI components
	common.DeployCAPIComponents(kindCluster)

	// Deploy the CABPK components
	applyManifests(kindCluster, cabpkComponents)

	// Deploy the CAPA components
	deployCAPAComponents(kindCluster)

	// Verify capi components are deployed
	common.WaitDeployment(kindClient, capiNamespace, capiDeploymentName)

	// Verify cabpk components are deployed
	common.WaitDeployment(kindClient, cabpkNamespace, cabpkDeploymentName)

	// Verify capa components are deployed
	common.WaitDeployment(kindClient, capaNamespace, capaDeploymentName)

	// Recreate kindClient so that it knows about the cluster api types
	kindClient, err = crclient.New(kindCluster.RestConfig(), crclient.Options{Scheme: setupScheme()})
	Expect(err).NotTo(HaveOccurred())
}, setupTimeout)

var _ = AfterSuite(func() {
	fmt.Fprintf(GinkgoWriter, "Tearing down kind cluster\n")
	retrieveAllLogs()
	kindCluster.Teardown()
	iamc := iam.New(sess)
	iamc.DeleteAccessKey(&iam.DeleteAccessKeyInput{UserName: accessKey.UserName, AccessKeyId: accessKey.AccessKeyId})
	os.RemoveAll(suiteTmpDir)
})

func retrieveAllLogs() {
	capiLogs := retrieveLogs(capiNamespace, capiDeploymentName)
	cabpkLogs := retrieveLogs(cabpkNamespace, cabpkDeploymentName)
	capaLogs := retrieveLogs(capaNamespace, capaDeploymentName)

	// If running in prow, output the logs to the artifacts path
	artifactPath, exists := os.LookupEnv("ARTIFACTS")
	if exists {
		ioutil.WriteFile(path.Join(artifactPath, "capi.log"), []byte(capiLogs), 0644)
		ioutil.WriteFile(path.Join(artifactPath, "cabpk.log"), []byte(cabpkLogs), 0644)
		ioutil.WriteFile(path.Join(artifactPath, "capa.log"), []byte(capaLogs), 0644)
		return
	}

	fmt.Fprintf(GinkgoWriter, "CAPI Logs:\n%s\n", capiLogs)
	fmt.Fprintf(GinkgoWriter, "CABPK Logs:\n%s\n", cabpkLogs)
	fmt.Fprintf(GinkgoWriter, "CAPA Logs:\n%s\n", capaLogs)
}

func retrieveLogs(namespace, deploymentName string) string {
	deployment := &appsv1.Deployment{}
	Expect(kindClient.Get(context.TODO(), crclient.ObjectKey{Namespace: namespace, Name: deploymentName}, deployment)).To(Succeed())

	pods := &corev1.PodList{}

	selector, err := metav1.LabelSelectorAsMap(deployment.Spec.Selector)
	Expect(err).NotTo(HaveOccurred())

	Expect(kindClient.List(context.TODO(), pods, crclient.InNamespace(namespace), crclient.MatchingLabels(selector))).To(Succeed())
	Expect(pods.Items).NotTo(BeEmpty())

	clientset, err := kubernetes.NewForConfig(kindCluster.RestConfig())
	Expect(err).NotTo(HaveOccurred())

	podLogs, err := clientset.CoreV1().Pods(namespace).GetLogs(pods.Items[0].Name, &corev1.PodLogOptions{Container: "manager"}).Stream()
	Expect(err).NotTo(HaveOccurred())
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	Expect(err).NotTo(HaveOccurred())

	return buf.String()
}

func getSession() client.ConfigProvider {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	Expect(err).NotTo(HaveOccurred())
	return sess
}

func getAccountID(prov client.ConfigProvider) string {
	stsSvc := sts.NewService(awssts.New(prov))
	accountID, err := stsSvc.AccountID()
	Expect(err).NotTo(HaveOccurred())
	return accountID
}

func createIAMRoles(prov client.ConfigProvider, accountID string) {
	cfnSvc := cloudformation.NewService(cfn.New(prov))
	Expect(
		cfnSvc.ReconcileBootstrapStack(stackName, accountID, "aws"),
	).To(Succeed())
}

func createKeyPair(prov client.ConfigProvider) {
	ec2c := ec2.New(prov)
	_, err := ec2c.CreateKeyPair(&ec2.CreateKeyPairInput{KeyName: aws.String(keyPairName)})
	if code, _ := awserrors.Code(err); code != "InvalidKeyPair.Duplicate" {
		Expect(err).NotTo(HaveOccurred())
	}
}

func loadManagerImage(kindCluster kind.Cluster) {
	if managerImage != nil && *managerImage != "" {
		kindCluster.LoadImage(*managerImage)
	}
}

func applyManifests(kindCluster kind.Cluster, manifests *string) {
	Expect(manifests).ToNot(BeNil())
	fmt.Fprintf(GinkgoWriter, "Applying manifests for %s\n", *manifests)
	Expect(*manifests).ToNot(BeEmpty())
	kindCluster.ApplyYAML(*manifests)
}

func deployCAPAComponents(kindCluster kind.Cluster) {
	if capaComponents != nil && *capaComponents != "" {
		applyManifests(kindCluster, capaComponents)
		return
	}

	fmt.Fprintf(GinkgoWriter, "Generating CAPA manifests\n")

	// Build the manifests using kustomize
	capaManifests, err := exec.Command(*kustomizeBinary, "build", "../../config/default").Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			fmt.Fprintf(GinkgoWriter, "Error: %s\n", string(exitError.Stderr))
		}
	}
	Expect(err).NotTo(HaveOccurred())

	// envsubst the credentials
	Expect(err).NotTo(HaveOccurred())
	b64credentials := generateB64Credentials()
	os.Setenv("AWS_B64ENCODED_CREDENTIALS", b64credentials)
	manifestsContent := os.ExpandEnv(string(capaManifests))

	// write out the manifests
	manifestFile := path.Join(suiteTmpDir, "infrastructure-components.yaml")
	Expect(ioutil.WriteFile(manifestFile, []byte(manifestsContent), 0644)).To(Succeed())

	// apply generated manifests
	applyManifests(kindCluster, &manifestFile)
}

const AWSCredentialsTemplate = `[default]
aws_access_key_id = {{ .AccessKeyID }}
aws_secret_access_key = {{ .SecretAccessKey }}
region = {{ .Region }}
`

type awsCredential struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
}

func generateB64Credentials() string {
	creds := awsCredential{
		Region:          region,
		AccessKeyID:     *accessKey.AccessKeyId,
		SecretAccessKey: *accessKey.SecretAccessKey,
	}

	tmpl, err := template.New("AWS Credentials").Parse(AWSCredentialsTemplate)
	Expect(err).NotTo(HaveOccurred())

	var profile bytes.Buffer
	Expect(tmpl.Execute(&profile, creds)).To(Succeed())

	encCreds := base64.StdEncoding.EncodeToString(profile.Bytes())
	return encCreds
}

func setupScheme() *runtime.Scheme {
	s := scheme.SetupScheme()
	Expect(bootstrapv1.AddToScheme(s)).To(Succeed())
	Expect(infrav1.AddToScheme(s)).To(Succeed())
	return s
}
