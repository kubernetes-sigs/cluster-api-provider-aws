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
	"flag"
	"fmt"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	awssts "github.com/aws/aws-sdk-go/service/sts"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/cloudformation"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/sts"
	"sigs.k8s.io/cluster-api-provider-aws/test/e2e/util/kind"
	"sigs.k8s.io/cluster-api/util"
)

const (
	kindTimeout = 5 * 60
	stackName   = "cluster-api-provider-aws-sigs-k8s-io"
	keyPairName = "cluster-api-provider-aws-sigs-k8s-io"
)

var (
	credFile   = flag.String("credFile", "", "path to an AWS credentials file")
	regionFile = flag.String("regionFile", "", "The path to a text file containing the AWS region")
	region     string
)

func initRegion() error {
	data, err := ioutil.ReadFile(*regionFile)
	if err != nil {
		return fmt.Errorf("error reading AWS region file: %v", err)
	}
	region = string(bytes.TrimSpace(data))
	return nil
}

var _ = Describe("AWS", func() {
	var (
		kindCluster kind.Cluster
	)
	BeforeEach(func() {
		fmt.Fprintf(GinkgoWriter, "running in AWS region: %s\n", region)
		kindCluster = kind.Cluster{
			Name: "capa-test-" + util.RandomString(6),
		}
		kindCluster.Setup()
	}, kindTimeout)

	AfterEach(func() {
		kindCluster.Teardown()
	})

	Describe("control plane node", func() {
		It("should be running", func() {
			sess := getSession()
			accountID := getAccountID(sess)

			createKeyPair(sess)
			createIAMRoles(sess, accountID)

			namespace := "test-" + util.RandomString(6)
			createNamespace(kindCluster.KubeClient(), namespace)
		})
	})
})

func createNamespace(client kubernetes.Interface, namespace string) {
	_, err := client.CoreV1().Namespaces().Create(&corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	})
	Expect(err).To(BeNil())
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

func getSession() client.ConfigProvider {
	creds := credentials.NewCredentials(&credentials.SharedCredentialsProvider{
		Filename: *credFile,
	})
	sess, err := session.NewSession(aws.NewConfig().WithCredentials(creds).WithRegion(region))
	Expect(err).To(BeNil())
	return sess
}
