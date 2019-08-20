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
	"flag"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	awssts "github.com/aws/aws-sdk-go/service/sts"
	appsv1 "k8s.io/api/apps/v1"
	apimachinerytypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	capi "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	"sigs.k8s.io/cluster-api/pkg/util"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/cloudformation"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/sts"
	"sigs.k8s.io/cluster-api-provider-aws/test/e2e/util/kind"
)

func TestE2e(t *testing.T) {
	if err := initRegion(); err != nil {
		t.Fatal(err)
	}
	RegisterFailHandler(Fail)
	RunSpecs(t, "e2e Suite")
}

const (
	setupTimeout          = 10 * 60
	capaProviderNamespace = "aws-provider-system"
	capaStatefulSetName   = "aws-provider-controller-manager"
	stackName             = "cluster-api-provider-aws-sigs-k8s-io"
	keyPairName           = "cluster-api-provider-aws-sigs-k8s-io"
)

var (
	credFile               = flag.String("credFile", "", "path to an AWS credentials file")
	regionFile             = flag.String("regionFile", "", "The path to a text file containing the AWS region")
	providerComponentsYAML = flag.String("providerComponentsYAML", "", "path to the provider components YAML for the cluster API")
	managerImageTar        = flag.String("managerImageTar", "", "a script to load the manager Docker image into Docker")

	kindCluster kind.Cluster
	kindClient  crclient.Client
	sess        client.ConfigProvider
	accountID   string
	region      string
)

var _ = BeforeSuite(func() {
	fmt.Fprintf(GinkgoWriter, "Setting up kind cluster\n")
	kindCluster = kind.Cluster{
		Name: "capa-test-" + util.RandomString(6),
	}
	kindCluster.Setup()
	loadManagerImage(kindCluster)

	fmt.Fprintf(GinkgoWriter, "Applying Provider Components to the kind cluster\n")
	applyProviderComponents(kindCluster)
	cfg := kindCluster.RestConfig()

	Expect(capi.SchemeBuilder.AddToScheme(scheme.Scheme)).NotTo(HaveOccurred())

	var err error
	kindClient, err = crclient.New(cfg, crclient.Options{Scheme: scheme.Scheme})
	Expect(err).To(BeNil())

	fmt.Fprintf(GinkgoWriter, "Creating AWS prerequisites\n")
	sess = getSession()
	accountID = getAccountID(sess)
	createKeyPair(sess)
	createIAMRoles(sess, accountID)

	fmt.Fprintf(GinkgoWriter, "Ensuring ProviderComponents are deployed\n")
	Eventually(
		func() (int32, error) {
			statefulSet := &appsv1.StatefulSet{}
			if err := kindClient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: capaProviderNamespace, Name: capaStatefulSetName}, statefulSet); err != nil {
				return 0, err
			}
			return statefulSet.Status.ReadyReplicas, nil
		}, 5*time.Minute, 15*time.Second,
	).ShouldNot(BeZero())

	fmt.Fprintf(GinkgoWriter, "Running in AWS region: %s\n", region)
}, setupTimeout)

var _ = AfterSuite(func() {
	fmt.Fprintf(GinkgoWriter, "Tearing down kind cluster\n")
	kindCluster.Teardown()
})

func initRegion() error {
	if regionFile != nil && *regionFile != "" {
		data, err := ioutil.ReadFile(*regionFile)
		if err != nil {
			return fmt.Errorf("error reading AWS region file: %v", err)
		}
		region = string(bytes.TrimSpace(data))
		return nil
	}

	region = "us-east-1"
	return nil
}

func getSession() client.ConfigProvider {
	if credFile != nil && *credFile != "" {
		creds := credentials.NewCredentials(&credentials.SharedCredentialsProvider{
			Filename: *credFile,
		})
		sess, err := session.NewSession(aws.NewConfig().WithCredentials(creds).WithRegion(region))
		Expect(err).To(BeNil())
		return sess
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	Expect(err).To(BeNil())
	return sess
}

func getAccountID(prov client.ConfigProvider) string {
	stsSvc := sts.NewService(awssts.New(prov))
	accountID, err := stsSvc.AccountID()
	Expect(err).To(BeNil())
	return accountID
}

func createIAMRoles(prov client.ConfigProvider, accountID string) {
	cfnSvc := cloudformation.NewService(cfn.New(prov))
	Expect(
		cfnSvc.ReconcileBootstrapStack(stackName, accountID),
	).To(Succeed())
}

func createKeyPair(prov client.ConfigProvider) {
	ec2c := ec2.New(prov)
	_, err := ec2c.CreateKeyPair(&ec2.CreateKeyPairInput{KeyName: aws.String(keyPairName)})
	if code, _ := awserrors.Code(err); code != "InvalidKeyPair.Duplicate" {
		Expect(err).To(BeNil())
	}
}

func loadManagerImage(kindCluster kind.Cluster) {
	if managerImageTar != nil && *managerImageTar != "" {
		kindCluster.LoadImageArchive(*managerImageTar)
	}
}

func applyProviderComponents(kindCluster kind.Cluster) {
	Expect(providerComponentsYAML).ToNot(BeNil())
	Expect(*providerComponentsYAML).ToNot(BeEmpty())
	kindCluster.ApplyYAML(*providerComponentsYAML)
}
