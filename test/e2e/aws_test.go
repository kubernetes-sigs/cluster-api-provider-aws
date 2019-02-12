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
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/onsi/gomega/types"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	capa "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/test/e2e/util/kind"
	capi "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	clientset "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset"
)

const (
	kindTimeout      = 5 * 60
	namespace        = "capa-test"
	clusterName      = "capa-test-cluster"
	controlPlaneName = "capa-test-control-plane"

	kubeletVersion = "v1.13.0"
	instanceType   = "t2.medium"
	instanceRegion = "us-east-1"
)

var _ = Describe("AWS", func() {
	var (
		cluster kind.Cluster
		client  *clientset.Clientset
	)
	BeforeEach(func() {
		cluster.Setup()
		cfg := cluster.RestConfig()
		var err error
		client, err = clientset.NewForConfig(cfg)
		Expect(err).To(BeNil())
	}, kindTimeout)

	AfterEach(func() {
		cluster.Teardown()
	})

	Describe("control plane node", func() {
		It("Should be healthy", func() {
			_, err := cluster.KubeClient().CoreV1().Namespaces().Create(&corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			})
			Expect(err).To(BeNil())

			machineapi := client.ClusterV1alpha1().Machines(namespace)
			clusterapi := client.ClusterV1alpha1().Clusters(namespace)

			_, err = clusterapi.Create(makeCluster())
			Expect(err).To(BeNil())

			_, err = machineapi.Create(makeMachine())
			Expect(err).To(BeNil())

			Eventually(
				func() (*capa.AWSMachineProviderStatus, error) {
					machine, err := machineapi.Get(controlPlaneName, metav1.GetOptions{})
					if err != nil {
						return nil, err
					}
					return capa.MachineStatusFromProviderStatus(machine.Status.ProviderStatus)
				},
				10*time.Minute, 30*time.Second,
			).Should(beHealthy())
		})
	})
})

func beHealthy() types.GomegaMatcher {
	return PointTo(
		MatchFields(IgnoreExtras, Fields{
			"InstanceState": PointTo(Equal(capa.InstanceStateRunning)),
		}),
	)
}

func makeCluster() *capi.Cluster {
	clusterSpec, err := capa.EncodeClusterSpec(&capa.AWSClusterProviderSpec{
		Region:     instanceRegion,
		SSHKeyName: "default",
	})
	Expect(err).To(BeNil())

	return &capi.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name: clusterName,
		},
		Spec: capi.ClusterSpec{
			ClusterNetwork: capi.ClusterNetworkingConfig{
				Services: capi.NetworkRanges{
					CIDRBlocks: []string{"10.96.0.0/12"},
				},
				Pods: capi.NetworkRanges{
					CIDRBlocks: []string{"192.168.0.0/16"},
				},
				ServiceDomain: "cluster.local",
			},
			ProviderSpec: capi.ProviderSpec{
				Value: clusterSpec,
			},
		},
	}
}

func makeMachine() *capi.Machine {
	providerSpec, err := capa.EncodeMachineSpec(&capa.AWSMachineProviderSpec{
		InstanceType:       instanceType,
		IAMInstanceProfile: "control-plane.cluster-api-provider-aws.sigs.k8s.io",
		KeyName:            "default",
	})
	Expect(err).To(BeNil())

	return &capi.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name: controlPlaneName,
			Labels: map[string]string{
				"set":                         "controlplane",
				"cluster.k8s.io/cluster-name": clusterName,
			},
		},
		Spec: capi.MachineSpec{
			Versions: capi.MachineVersionInfo{
				Kubelet:      kubeletVersion,
				ControlPlane: kubeletVersion,
			},
			ProviderSpec: capi.ProviderSpec{
				Value: providerSpec,
			},
		},
	}
}
