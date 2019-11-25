// +build integration

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

package integration_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/iam"
	testutil "sigs.k8s.io/cluster-api-provider-aws/test/commons"
	common "sigs.k8s.io/cluster-api/test/helpers/components"
	capiFlag "sigs.k8s.io/cluster-api/test/helpers/flag"
	"sigs.k8s.io/cluster-api/test/helpers/kind"
	"sigs.k8s.io/cluster-api/util"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)

	// If running in prow, output the junit files to the artifacts path
	junitPath := fmt.Sprintf("junit.integration_suite.%d.xml", config.GinkgoConfig.ParallelNode)
	artifactPath, exists := os.LookupEnv("ARTIFACTS")
	if exists {
		junitPath = path.Join(artifactPath, junitPath)
	}
	junitReporter := reporters.NewJUnitReporter(junitPath)
	RunSpecsWithDefaultAndCustomReporters(t, "integration Suite", []Reporter{junitReporter})
}

const (
	CABPK_VERSION = "v0.1.0"
	setupTimeout  = 10 * 60
	stackName     = "cluster-api-provider-aws-sigs-k8s-io"
	KeyPairName   = "cluster-api-provider-aws-sigs-k8s-io"
)

var (
	managerImage    = capiFlag.DefineOrLookupStringFlag("managerImage", "", "Docker image to load into the kind cluster for testing")
	cabpkComponents = capiFlag.DefineOrLookupStringFlag("cabpkComponents", "https://github.com/kubernetes-sigs/cluster-api-bootstrap-provider-kubeadm/releases/download/"+CABPK_VERSION+"/bootstrap-components.yaml", "URL to CAPI components to load")

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
	suiteTmpDir, err = ioutil.TempDir("", "capa-integration-suite")
	Expect(err).NotTo(HaveOccurred())

	var ok bool
	region, ok = os.LookupEnv("AWS_REGION")
	testutil.Region = &region
	fmt.Fprintf(GinkgoWriter, "Running in region: %s\n", region)
	if !ok {
		fmt.Fprintf(GinkgoWriter, "Environment variable AWS_REGION not found")
		Expect(ok).To(BeTrue())
	}

	sess = testutil.GetSession()

	fmt.Fprintf(GinkgoWriter, "Creating AWS prerequisites\n")
	accountID = testutil.GetAccountID(sess)
	testutil.CreateKeyPair(sess, KeyPairName)
	testutil.CreateIAMRoles(sess, accountID, stackName)

	iamc := iam.New(sess)
	out, err := iamc.CreateAccessKey(&iam.CreateAccessKeyInput{UserName: aws.String("bootstrapper.cluster-api-provider-aws.sigs.k8s.io")})
	Expect(err).NotTo(HaveOccurred())
	Expect(out.AccessKey).NotTo(BeNil())
	accessKey = out.AccessKey

	kindCluster = kind.Cluster{
		Name: "capa-test-" + util.RandomString(6),
	}
	kindCluster.Setup()
	testutil.LoadManagerImage(kindCluster, managerImage)

	kindClient, err = crclient.New(kindCluster.RestConfig(), crclient.Options{Scheme: testutil.SetupScheme()})
	Expect(err).NotTo(HaveOccurred())

	// Deploy the CAPI components
	common.DeployCAPIComponents(kindCluster)

	// Deploy the CABPK components
	testutil.ApplyManifests(kindCluster, cabpkComponents)

	// Deploy the CAPA components
	testutil.DeployCAPAComponents(kindCluster, suiteTmpDir, accessKey)

	// Verify capi components are deployed
	common.WaitDeployment(kindClient, testutil.CapiNamespace, testutil.CapiDeploymentName)

	// Verify cabpk components are deployed
	common.WaitDeployment(kindClient, testutil.CabpkNamespace, testutil.CabpkDeploymentName)

	// Verify capa components are deployed
	common.WaitDeployment(kindClient, testutil.CapaNamespace, testutil.CapaDeploymentName)

	// Recreate kindClient so that it knows about the cluster api types
	kindClient, err = crclient.New(kindCluster.RestConfig(), crclient.Options{Scheme: testutil.SetupScheme()})
	Expect(err).NotTo(HaveOccurred())
}, setupTimeout)

var _ = AfterSuite(func() {
	fmt.Fprintf(GinkgoWriter, "Tearing down kind cluster\n")
	testutil.RetrieveAllLogs(kindCluster, kindClient)
	kindCluster.Teardown()
	iamc := iam.New(sess)
	iamc.DeleteAccessKey(&iam.DeleteAccessKeyInput{UserName: accessKey.UserName, AccessKeyId: accessKey.AccessKeyId})
	testutil.DeleteIAMRoles(sess, stackName)
	os.RemoveAll(suiteTmpDir)
})
