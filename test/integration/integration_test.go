// +build integration

/*
Copyright 2019 The Kubernetes Authors.

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
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	testutil "sigs.k8s.io/cluster-api-provider-aws/test/commons"
	"sigs.k8s.io/cluster-api/util"
)

var _ = Describe("integration tests", func() {
	var (
		namespace   string
		clusterName string
	)

	BeforeEach(func() {
		namespace = "test-namespace-" + util.RandomString(6)
		testutil.CreateNamespace(namespace, kindClient)
	})

	Describe("cluster name validation tests", func() {
		subDomainErrorMsg := "a DNS-1123 subdomain must consist of lower case alphanumeric characters, '-' or '.'," +
			" and must start and end with an alphanumeric character"

		It("cluster creation with name having capital letters should fail", func() {
			clusterName = "capa-cluster-CAPS"
			err := testutil.MakeAWSCluster(namespace, clusterName, KeyPairName, kindClient)
			Expect(err).NotTo(BeNil())
			Expect(strings.Contains(err.Error(), subDomainErrorMsg)).To(BeTrue())
			err = testutil.MakeCluster(namespace, clusterName, clusterName, kindClient)
			Expect(err).NotTo(BeNil())
			Expect(strings.Contains(err.Error(), subDomainErrorMsg)).To(BeTrue())
		})

		It("cluster creation with name having '_' should fail", func() {
			clusterName = "capa_cluster"
			err := testutil.MakeAWSCluster(namespace, clusterName, KeyPairName, kindClient)
			Expect(err).NotTo(BeNil())
			Expect(strings.Contains(err.Error(), subDomainErrorMsg)).To(BeTrue())
			err = testutil.MakeCluster(namespace, clusterName, clusterName, kindClient)
			Expect(err).NotTo(BeNil())
			Expect(strings.Contains(err.Error(), subDomainErrorMsg)).To(BeTrue())
		})

		It("cluster creation with name having only digits should succeed", func() {
			clusterName = "123456"
			Expect(testutil.MakeAWSCluster(namespace, clusterName, KeyPairName, kindClient)).ShouldNot(HaveOccurred())
			Expect(testutil.MakeCluster(namespace, clusterName, clusterName, kindClient)).ShouldNot(HaveOccurred())
			Expect(testutil.IsClusterProvisioningCompleted(clusterName, namespace, kindClient, 5)).To(BeTrue())
			testutil.DeleteCluster(namespace, clusterName, kindClient)
		})

		It("cluster creation with name having '.' should succeed", func() {
			clusterName = "capa.cluster"
			Expect(testutil.MakeAWSCluster(namespace, clusterName, KeyPairName, kindClient)).ShouldNot(HaveOccurred())
			Expect(testutil.MakeCluster(namespace, clusterName, clusterName, kindClient)).ShouldNot(HaveOccurred())
			Expect(testutil.IsClusterProvisioningCompleted(clusterName, namespace, kindClient, 5)).To(BeTrue())
			testutil.DeleteCluster(namespace, clusterName, kindClient)
		})

		It("cluster creation with name lenght more than 32 characters should succeed", func() {
			clusterName = "abcdefghijklmnopqrstuvwxyzabcdefgh"
			Expect(testutil.MakeAWSCluster(namespace, clusterName, KeyPairName, kindClient)).ShouldNot(HaveOccurred())
			Expect(testutil.MakeCluster(namespace, clusterName, clusterName, kindClient)).ShouldNot(HaveOccurred())
			Expect(testutil.IsClusterProvisioningCompleted(clusterName, namespace, kindClient, 5)).To(BeTrue())
			testutil.DeleteCluster(namespace, clusterName, kindClient)
		})
	})
})
