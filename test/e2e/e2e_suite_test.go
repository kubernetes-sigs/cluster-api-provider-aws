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
	"encoding/base64"
	"fmt"
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
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/session"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	awssts "github.com/aws/aws-sdk-go/service/sts"
	"k8s.io/apimachinery/pkg/runtime"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/cloudformation"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/sts"
	common "sigs.k8s.io/cluster-api/test/helpers/components"
	capiFlag "sigs.k8s.io/cluster-api/test/helpers/flag"
	"sigs.k8s.io/cluster-api/test/helpers/kind"
	"sigs.k8s.io/cluster-api/test/helpers/scheme"
	"sigs.k8s.io/cluster-api/util"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
	bootstrapv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1alpha2"
)

func TestE2e(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "e2e Suite")
}

const (
	capiNamespace       = "capi-system"
	capiDeploymentName  = "capi-controller-manager"
	capaNamespace       = "capa-system"
	capaDeploymentName  = "capa-controller-manager"
	setupTimeout        = 10 * 60
	stackName           = "cluster-api-provider-aws-sigs-k8s-io"
	keyPairName         = "cluster-api-provider-aws-sigs-k8s-io"
)

var (
	managerImage    = capiFlag.DefineOrLookupStringFlag("managerImage", "", "Docker image to load into the kind cluster for testing")
	capaComponents  = capiFlag.DefineOrLookupStringFlag("capaComponents", "", "capa components to load")
	kustomizeBinary = capiFlag.DefineOrLookupStringFlag("kustomizeBinary", "kustomize", "path to the kustomize binary")

	kindCluster kind.Cluster
	kindClient  crclient.Client
	sess        client.ConfigProvider
	accountID   string
	suiteTmpDir string
)

var _ = BeforeSuite(func() {
	fmt.Fprintf(GinkgoWriter, "Setting up kind cluster\n")

	var err error
	suiteTmpDir, err = ioutil.TempDir("", "capa-e2e-suite")
	Expect(err).NotTo(HaveOccurred())

	kindCluster = kind.Cluster{
		Name: "capa-test-" + util.RandomString(6),
	}
	kindCluster.Setup()
	loadManagerImage(kindCluster)

	// Deploy CertManager
	certmanagerYaml := "https://raw.githubusercontent.com/kubernetes-sigs/cluster-api/master/config/certmanager/cert-manager.yaml"
	applyManifests(kindCluster, &certmanagerYaml)

	kindClient, err = crclient.New(kindCluster.RestConfig(), crclient.Options{Scheme: setupScheme()})
	Expect(err).NotTo(HaveOccurred())
	common.WaitDeployment(kindClient, "cert-manager", "cert-manager-webhook")

	// Deploy the CAPI components
	// workaround since there isn't a v1alpha3 capi release yet
	deployCAPIComponents(kindCluster)

	// Deploy the CAPA components
	deployCAPAComponents(kindCluster)

	kindClient, err = crclient.New(kindCluster.RestConfig(), crclient.Options{Scheme: setupScheme()})
	Expect(err).NotTo(HaveOccurred())

	fmt.Fprintf(GinkgoWriter, "Creating AWS prerequisites\n")
	sess = getSession()
	accountID = getAccountID(sess)
	createKeyPair(sess)
	createIAMRoles(sess, accountID)

	// Verify capi components are deployed
	common.WaitDeployment(kindClient, capiNamespace, capiDeploymentName)

	// Verify capa components are deployed
	common.WaitDeployment(kindClient, capaNamespace, capaDeploymentName)
}, setupTimeout)

var _ = AfterSuite(func() {
	fmt.Fprintf(GinkgoWriter, "Tearing down kind cluster\n")
	kindCluster.Teardown()
	os.RemoveAll(suiteTmpDir)
})

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

func deployCAPIComponents(kindCluster kind.Cluster) {
	fmt.Fprintf(GinkgoWriter, "Generating CAPI manifests\n")

	// Build the manifests using kustomize
	capiManifests, err := exec.Command(*kustomizeBinary, "build", "https://github.com/kubernetes-sigs/cluster-api//config/default").Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			fmt.Fprintf(GinkgoWriter, "Error: %s\n", string(exitError.Stderr))
		}
	}
	Expect(err).NotTo(HaveOccurred())

	// write out the manifests
	manifestFile := path.Join(suiteTmpDir, "cluster-api-components.yaml")
	Expect(ioutil.WriteFile(manifestFile, capiManifests, 0644)).To(Succeed())

	// apply generated manifests
	applyManifests(kindCluster, &manifestFile)
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
	b64credentials, err := generateB64Credentials()
	Expect(err).NotTo(HaveOccurred())
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
{{if .SessionToken }}
aws_session_token = {{ .SessionToken }}
{{end}}
`

type awsCredential struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	Region          string
}

func generateB64Credentials() (string, error) {
	creds := awsCredential{}
	conf := aws.NewConfig()
	chain := defaults.CredChain(conf, defaults.Handlers())
	chainCreds, err := chain.Get()
	if err != nil {
		return "", err
	}

	// still needed as defaults.CredChain doesn't contain region
	region, ok := os.LookupEnv("AWS_REGION")
	if !ok {
		return "", fmt.Errorf("Environment variable AWS_REGION not found")
	}
	creds.Region = region

	creds.AccessKeyID = chainCreds.AccessKeyID
	creds.SecretAccessKey = chainCreds.SecretAccessKey
	creds.SessionToken = chainCreds.SessionToken

	tmpl, err := template.New("AWS Credentials").Parse(AWSCredentialsTemplate)
	if err != nil {
		return "", err
	}

	var profile bytes.Buffer
	err = tmpl.Execute(&profile, creds)
	if err != nil {
		return "", err
	}

	encCreds := base64.StdEncoding.EncodeToString(profile.Bytes())
	return encCreds, nil
}

func setupScheme() *runtime.Scheme {
	s := scheme.SetupScheme()
	Expect(bootstrapv1.AddToScheme(s)).To(Succeed())
	Expect(infrav1.AddToScheme(s)).To(Succeed())
	return s
}
