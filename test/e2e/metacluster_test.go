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
	"flag"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/onsi/gomega/types"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"sigs.k8s.io/cluster-api-provider-aws/test/e2e/util/kind"
)

const (
	kindTimeout         = 5 * 60
	controllerNamespace = "aws-provider-system"
	controllerName      = "aws-provider-controller-manager"
)

const testClusterAndMachine = `apiVersion: "cluster.k8s.io/v1alpha1"
kind: Cluster
metadata:
  name: test1
spec:
  clusterNetwork:
    services:
      cidrBlocks: ["10.96.0.0/12"]
    pods:
      cidrBlocks: ["192.168.0.0/16"]
    serviceDomain: "cluster.local"
  providerSpec:
    value:
      apiVersion: "awsprovider/v1alpha1"
      kind: "AWSClusterProviderSpec"
      region: "us-east-1"
      sshKeyName: "default"
---
apiVersion: "cluster.k8s.io/v1alpha1"
kind: Machine
metadata:
  name: aws-controlplane-0
  labels:
    set: controlplane
spec:
  versions:
    kubelet: v1.13.0
    controlPlane: v1.13.0
  providerSpec:
    value:
      apiVersion: awsprovider/v1alpha1
      kind: AWSMachineProviderSpec
      instanceType: "t2.medium"
      iamInstanceProfile: "control-plane.cluster-api-provider-aws.sigs.k8s.io"
      keyName: "default"
`

var clusterCtlBinary = flag.String("clusterCtlBinary", "clusterctl", "Path to the clusterctl binary")

var _ = Describe("Metacluster", func() {
	var (
		cluster kind.Cluster
		client  kubernetes.Interface
	)
	BeforeEach(func() {
		cluster.Setup()
		client = cluster.KubeClient()
	}, kindTimeout)

	AfterEach(func() {
		// cluster.Teardown()
	})

	Describe("manager container", func() {
		It("Should be healthy", func() {
			sets := client.AppsV1().StatefulSets(controllerNamespace)

			Eventually(
				func() (*appsv1.StatefulSet, error) {
					return sets.Get(controllerName, metav1.GetOptions{})
				},
				2*time.Minute, 5*time.Second,
			).Should(haveReplicas(1))
		})

		It("Should create a cluster", func() {
			cmd := cluster.GetKubectlCommand(
				"create",
				"-f", "-",
			)
			cmd.Stdin = strings.NewReader(testClusterAndMachine)
			cluster.Run(cmd)
		})

	})
})

// haveReplicas matches a stateful set with i replicas
func haveReplicas(i int32) types.GomegaMatcher {
	return PointTo(
		MatchFields(IgnoreExtras, Fields{
			"Status": MatchFields(IgnoreExtras, Fields{
				"Replicas":      Equal(i),
				"ReadyReplicas": Equal(i),
			}),
		}),
	)
}
