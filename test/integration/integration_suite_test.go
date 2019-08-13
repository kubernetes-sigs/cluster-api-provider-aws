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
	"context"
	"flag"
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	appsv1 "k8s.io/api/apps/v1"
	apimachinerytypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/cluster-api/pkg/util"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/cluster-api-provider-aws/test/e2e/util/kind"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

const (
	setupTimeout          = 5 * 60
	capaProviderNamespace = "aws-provider-system"
	capaStatefulSetName   = "aws-provider-controller-manager"
)

var (
	providerComponentsYAML = flag.String("providerComponentsYAML", "", "path to the provider components YAML for the cluster API")
	managerImageTar        = flag.String("managerImageTar", "", "a script to load the manager Docker image into Docker")

	kindCluster kind.Cluster
	kindClient  crclient.Client
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
	var err error
	kindClient, err = crclient.New(cfg, crclient.Options{})
	Expect(err).To(BeNil())

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
}, setupTimeout)

var _ = AfterSuite(func() {
	fmt.Fprintf(GinkgoWriter, "Tearing down kind cluster\n")
	kindCluster.Teardown()
})

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
