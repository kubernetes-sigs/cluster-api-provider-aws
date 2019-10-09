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
	"context"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apimachinerytypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/clientcmd"
	bootstrapv1 "sigs.k8s.io/cluster-api-bootstrap-provider-kubeadm/api/v1alpha2"
	kubeadmv1beta1 "sigs.k8s.io/cluster-api-bootstrap-provider-kubeadm/kubeadm/v1beta1"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha2"
	capiFlag "sigs.k8s.io/cluster-api/test/helpers/flag"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/kubeconfig"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	cniManifests  = capiFlag.DefineOrLookupStringFlag("cniManifests", "https://docs.projectcalico.org/v3.8/manifests/calico.yaml", "URL to CNI manifests to load")
	kubectlBinary = capiFlag.DefineOrLookupStringFlag("kubectlBinary", "kubectl", "path to the kubectl binary")
)

var _ = Describe("functional tests", func() {
	var (
		namespace               string
		clusterName             string
		awsClusterName          string
		cpMachinePrefix         string
		cpAWSMachinePrefix      string
		cpBootstrapConfigPrefix string
		testTmpDir              string
	)
	BeforeEach(func() {
		var err error
		testTmpDir, err = ioutil.TempDir(suiteTmpDir, "aws-test")
		Expect(err).NotTo(HaveOccurred())

		namespace = "test-" + util.RandomString(6)
		createNamespace(namespace)

		clusterName = "test-" + util.RandomString(6)
		awsClusterName = "test-infra-" + util.RandomString(6)
		cpMachinePrefix = "test-" + util.RandomString(6)
		cpAWSMachinePrefix = "test-infra-" + util.RandomString(6)
		cpBootstrapConfigPrefix = "test-boot-" + util.RandomString(6)
	})

	Describe("workload cluster lifecycle", func() {
		It("It should be creatable and deletable", func() {
			By("Creating an AWSCluster")
			makeAWSCluster(namespace, awsClusterName)

			By("Creating a Cluster")
			makeCluster(namespace, clusterName, awsClusterName)

			By("Ensuring Cluster Infrastructure Reports as Ready")
			waitForClusterInfrastructureReady(namespace, clusterName)

			By("Creating the initial Control Plane Machine")
			awsMachineName := cpAWSMachinePrefix + "-0"
			bootstrapConfigName := cpBootstrapConfigPrefix + "-0"
			machineName := cpMachinePrefix + "-0"
			createInitialControlPlaneMachine(namespace, clusterName, machineName, awsMachineName, bootstrapConfigName)

			By("Deploying CNI to created Cluster")
			deployCNI(testTmpDir, namespace, clusterName, *cniManifests)
			waitForMachineNodeReady(namespace, machineName)

			By("Creating the second Control Plane Machine")
			awsMachineName = cpAWSMachinePrefix + "-1"
			bootstrapConfigName = cpBootstrapConfigPrefix + "-1"
			machineName = cpMachinePrefix + "-1"
			createAdditionalControlPlaneMachine(namespace, clusterName, machineName, awsMachineName, bootstrapConfigName)

			By("Creating the third Control Plane Machine")
			awsMachineName = cpAWSMachinePrefix + "-2"
			bootstrapConfigName = cpBootstrapConfigPrefix + "-2"
			machineName = cpMachinePrefix + "-2"
			createAdditionalControlPlaneMachine(namespace, clusterName, machineName, awsMachineName, bootstrapConfigName)

			// TODO: Deploy a MachineDeployment
			// TODO: Scale MachineDeployment up
			// TODO: Scale MachineDeployment down

			By("Deleting the Cluster")
			deleteCluster(namespace, clusterName)

		})
	})
})

func deployCNI(tmpDir, namespace, clusterName, manifestPath string) {
	cluster := &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      clusterName,
		},
	}
	kubeConfigData, err := kubeconfig.FromSecret(kindClient, cluster)
	Expect(err).NotTo(HaveOccurred())
	kubeConfigPath := path.Join(tmpDir, clusterName+".kubeconfig")
	Expect(ioutil.WriteFile(kubeConfigPath, kubeConfigData, 0640)).To(Succeed())

	Expect(exec.Command(
		*kubectlBinary,
		"create",
		"--kubeconfig="+kubeConfigPath,
		"-f", manifestPath,
	).Run()).To(Succeed())
}

func waitForMachineNodeReady(namespace, name string) {
	machine := &clusterv1.Machine{}
	Expect(kindClient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: namespace, Name: name}, machine)).To(Succeed())

	nodeName := machine.Status.NodeRef.Name

	cluster := &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      machine.Labels[clusterv1.MachineClusterLabelName],
		},
	}
	kubeConfigData, err := kubeconfig.FromSecret(kindClient, cluster)
	Expect(err).NotTo(HaveOccurred())

	restConfig, err := clientcmd.RESTConfigFromKubeConfig(kubeConfigData)
	Expect(err).NotTo(HaveOccurred())

	nodeClient, err := crclient.New(restConfig, crclient.Options{})
	Expect(err).NotTo(HaveOccurred())

	Eventually(
		func() bool {
			node := &corev1.Node{}
			if err := nodeClient.Get(context.TODO(), apimachinerytypes.NamespacedName{Name: nodeName}, node); err != nil {
				fmt.Fprintf(GinkgoWriter, "Error retrieving node: %v", err)
				return false
			}

			conditionMap := make(map[corev1.NodeConditionType]*corev1.NodeCondition)
			for i, condition := range node.Status.Conditions {
				if condition.Type == corev1.NodeReady {
					conditionMap[condition.Type] = &node.Status.Conditions[i]
				}
			}

			if condition, ok := conditionMap[corev1.NodeReady]; ok {
				if condition.Status == corev1.ConditionTrue {
					return true
				}
			}
			return false
		},
		2*time.Minute, 15*time.Second,
	).Should(BeTrue())
}

func createAdditionalControlPlaneMachine(namespace, clusterName, machineName, awsMachineName, bootstrapConfigName string) {
	makeAWSMachine(namespace, awsMachineName)
	makeJoinBootstrapConfig(namespace, bootstrapConfigName)
	makeMachine(namespace, machineName, awsMachineName, bootstrapConfigName, clusterName)
	waitForMachineBootstrapReady(namespace, machineName)
	waitForAWSMachineRunning(namespace, awsMachineName)
	waitForAWSMachineReady(namespace, awsMachineName)
	waitForMachineNodeRef(namespace, machineName)
	waitForMachineNodeReady(namespace, machineName)
}

func createInitialControlPlaneMachine(namespace, clusterName, machineName, awsMachineName, bootstrapConfigName string) {
	makeAWSMachine(namespace, awsMachineName)
	makeInitBootstrapConfig(namespace, bootstrapConfigName)
	makeMachine(namespace, machineName, awsMachineName, bootstrapConfigName, clusterName)
	waitForMachineBootstrapReady(namespace, machineName)
	waitForAWSMachineRunning(namespace, awsMachineName)
	waitForAWSMachineReady(namespace, awsMachineName)
	waitForMachineNodeRef(namespace, machineName)
	waitForClusterControlPlaneInitialized(namespace, clusterName)
}

func deleteCluster(namespace, name string) {
	fmt.Fprintf(GinkgoWriter, "Deleting Cluster %s/%s\n", namespace, name)
	Expect(kindClient.Delete(context.TODO(), &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
	})).To(Succeed())

	Eventually(
		func() bool {
			cluster := &clusterv1.Cluster{}
			if err := kindClient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: namespace, Name: name}, cluster); err != nil {
				if apierrors.IsNotFound(err) {
					return true
				}
				return false
			}
			return false
		},
		20*time.Minute, 15*time.Second,
	).Should(BeTrue())
}

func waitForMachineNodeRef(namespace, name string) {
	Eventually(
		func() *corev1.ObjectReference {
			machine := &clusterv1.Machine{}
			if err := kindClient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: namespace, Name: name}, machine); err != nil {
				return nil
			}
			return machine.Status.NodeRef

		},
		5*time.Minute, 15*time.Second,
	).ShouldNot(BeNil())
}

func waitForClusterControlPlaneInitialized(namespace, name string) {
	fmt.Fprintf(GinkgoWriter, "Ensuring control plane initialized for cluster %s/%s\n", namespace, name)
	Eventually(
		func() (bool, error) {
			cluster := &clusterv1.Cluster{}
			if err := kindClient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: namespace, Name: name}, cluster); err != nil {
				return false, err
			}
			return cluster.Status.ControlPlaneInitialized, nil
		},
		10*time.Minute, 15*time.Second,
	).Should(BeTrue())
}

func waitForAWSMachineRunning(namespace, name string) {
	Eventually(
		func() (bool, error) {
			awsMachine := &infrav1.AWSMachine{}
			if err := kindClient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: namespace, Name: name}, awsMachine); err != nil {
				return false, err
			}
			if awsMachine.Status.InstanceState == nil {
				return false, nil
			}
			return *awsMachine.Status.InstanceState == infrav1.InstanceStateRunning, nil
		},
		5*time.Minute, 15*time.Second,
	).Should(BeTrue())
}

func waitForClusterInfrastructureReady(namespace, name string) {
	fmt.Fprintf(GinkgoWriter, "Ensuring infrastructure is ready for cluster %s/%s\n", namespace, name)
	Eventually(
		func() (bool, error) {
			cluster := &clusterv1.Cluster{}
			if err := kindClient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: namespace, Name: name}, cluster); err != nil {
				return false, err
			}
			return cluster.Status.InfrastructureReady, nil
		},
		10*time.Minute, 15*time.Second,
	).Should(BeTrue())
}

func waitForMachineBootstrapReady(namespace, name string) {
	fmt.Fprintf(GinkgoWriter, "Ensuring Machine %s/%s has bootstrapReady\n", namespace, name)
	Eventually(
		func() (bool, error) {
			machine := &clusterv1.Machine{}
			if err := kindClient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: namespace, Name: name}, machine); err != nil {
				return false, err
			}
			return machine.Status.BootstrapReady, nil
		},
		2*time.Minute, 15*time.Second,
	).Should(BeTrue())
}

func waitForAWSMachineReady(namespace, name string) {
	Eventually(
		func() (bool, error) {
			awsMachine := &infrav1.AWSMachine{}
			if err := kindClient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: namespace, Name: name}, awsMachine); err != nil {
				return false, err
			}
			return awsMachine.Status.Ready, nil
		},
		2*time.Minute, 15*time.Second,
	).Should(BeTrue())
}

func makeMachine(namespace, name, awsMachineName, bootstrapConfigName, clusterName string) {
	fmt.Fprintf(GinkgoWriter, "Creating Machine %s/%s\n", namespace, name)
	k8s_version := "v1.15.3"
	machine := &clusterv1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				clusterv1.MachineClusterLabelName:      clusterName,
				clusterv1.MachineControlPlaneLabelName: "true",
			},
		},
		Spec: clusterv1.MachineSpec{
			Bootstrap: clusterv1.Bootstrap{
				ConfigRef: &corev1.ObjectReference{
					Kind:       "KubeadmConfig",
					APIVersion: bootstrapv1.GroupVersion.String(),
					Name:       bootstrapConfigName,
					Namespace:  namespace,
				},
			},
			InfrastructureRef: corev1.ObjectReference{
				Kind:       "AWSMachine",
				APIVersion: infrav1.GroupVersion.String(),
				Name:       awsMachineName,
				Namespace:  namespace,
			},
			Version: &k8s_version,
		},
	}
	Expect(kindClient.Create(context.TODO(), machine)).To(Succeed())
}

func makeJoinBootstrapConfig(namespace, name string) {
	fmt.Fprintf(GinkgoWriter, "Creating Join KubeadmConfig %s/%s\n", namespace, name)
	config := &bootstrapv1.KubeadmConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: bootstrapv1.KubeadmConfigSpec{
			JoinConfiguration: &kubeadmv1beta1.JoinConfiguration{
				ControlPlane: &kubeadmv1beta1.JoinControlPlane{},
				NodeRegistration: kubeadmv1beta1.NodeRegistrationOptions{
					Name:             "{{ ds.meta_data.hostname }}",
					KubeletExtraArgs: map[string]string{"cloud-provider": "aws"},
				},
			},
		},
	}
	Expect(kindClient.Create(context.TODO(), config)).To(Succeed())
}

func makeInitBootstrapConfig(namespace, name string) {
	fmt.Fprintf(GinkgoWriter, "Creating Init KubeadmConfig %s/%s\n", namespace, name)
	config := &bootstrapv1.KubeadmConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: bootstrapv1.KubeadmConfigSpec{
			ClusterConfiguration: &kubeadmv1beta1.ClusterConfiguration{
				APIServer: kubeadmv1beta1.APIServer{
					ControlPlaneComponent: kubeadmv1beta1.ControlPlaneComponent{
						ExtraArgs: map[string]string{"cloud-provider": "aws"},
					},
				},
				ControllerManager: kubeadmv1beta1.ControlPlaneComponent{
					ExtraArgs: map[string]string{"cloud-provider": "aws"},
				},
			},
			InitConfiguration: &kubeadmv1beta1.InitConfiguration{
				NodeRegistration: kubeadmv1beta1.NodeRegistrationOptions{
					Name:             "{{ ds.meta_data.hostname }}",
					KubeletExtraArgs: map[string]string{"cloud-provider": "aws"},
				},
			},
		},
	}
	Expect(kindClient.Create(context.TODO(), config)).To(Succeed())
}

func makeAWSMachine(namespace, name string) {
	fmt.Fprintf(GinkgoWriter, "Creating AWSMachine %s/%s\n", namespace, name)
	awsMachine := &infrav1.AWSMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: infrav1.AWSMachineSpec{
			InstanceType:       "t3.large",
			IAMInstanceProfile: "control-plane.cluster-api-provider-aws.sigs.k8s.io",
			SSHKeyName:         keyPairName,
		},
	}
	Expect(kindClient.Create(context.TODO(), awsMachine)).To(Succeed())
}

func makeCluster(namespace, name, awsClusterName string) {
	fmt.Fprintf(GinkgoWriter, "Creating Cluster %s/%s\n", namespace, name)
	cluster := &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: clusterv1.ClusterSpec{
			ClusterNetwork: &clusterv1.ClusterNetwork{
				Pods: &clusterv1.NetworkRanges{
					CIDRBlocks: []string{"192.168.0.0/16"},
				},
			},
			InfrastructureRef: &corev1.ObjectReference{
				Kind:       "AWSCluster",
				APIVersion: infrav1.GroupVersion.String(),
				Name:       awsClusterName,
				Namespace:  namespace,
			},
		},
	}
	Expect(kindClient.Create(context.TODO(), cluster)).To(Succeed())
}

func makeAWSCluster(namespace, name string) {
	fmt.Fprintf(GinkgoWriter, "Creating AWSCluster %s/%s\n", namespace, name)
	awsCluster := &infrav1.AWSCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: infrav1.AWSClusterSpec{
			Region:     "us-east-1",
			SSHKeyName: keyPairName,
		},
	}
	Expect(kindClient.Create(context.TODO(), awsCluster)).To(Succeed())
}

func createNamespace(namespace string) {
	fmt.Fprintf(GinkgoWriter, "creating namespace %q\n", namespace)
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
	Expect(kindClient.Create(context.TODO(), ns)).To(Succeed())
}
