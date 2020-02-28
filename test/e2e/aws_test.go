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
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apirand "k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	bootstrapv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1alpha3"
	kubeadmv1beta1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/types/v1beta1"
	capiFlag "sigs.k8s.io/cluster-api/test/helpers/flag"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/kubeconfig"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
)

var (
	cniManifests  = capiFlag.DefineOrLookupStringFlag("cniManifests", "https://docs.projectcalico.org/v3.12/manifests/calico.yaml", "URL to CNI manifests to load")
	kubectlBinary = capiFlag.DefineOrLookupStringFlag("kubectlBinary", "kubectl", "path to the kubectl binary")
)

const (
	scaleUpReplicas   = 5
	scaleDownReplicas = 0
)

var _ = Describe("functional tests", func() {
	var (
		namespace               string
		clusterName             string
		awsClusterName          string
		cpMachinePrefix         string
		cpAWSMachinePrefix      string
		cpBootstrapConfigPrefix string
		mdBootstrapConfig       string
		machineDeploymentName   string
		awsMachineTemplateName  string
		testTmpDir              string
		initialReplicas         int32
		cancelWatches           context.CancelFunc
	)
	BeforeEach(func() {
		var ctx context.Context
		ctx, cancelWatches = context.WithCancel(context.Background())

		var err error
		testTmpDir, err = ioutil.TempDir(suiteTmpDir, "aws-test")
		Expect(err).NotTo(HaveOccurred())

		namespace = "test-" + util.RandomString(6)
		createNamespace(namespace)

		go func() {
			defer GinkgoRecover()
			watchEvents(ctx, namespace)
		}()

		clusterName = "test-" + util.RandomString(6)
		awsClusterName = "test-infra-" + util.RandomString(6)
		cpMachinePrefix = "test-" + util.RandomString(6)
		cpAWSMachinePrefix = "test-infra-" + util.RandomString(6)
		cpBootstrapConfigPrefix = "test-boot-" + util.RandomString(6)
		mdBootstrapConfig = "test-boot-md" + util.RandomString(6)
		machineDeploymentName = "test-capa-md" + util.RandomString(6)
		awsMachineTemplateName = "test-infra-capa-mt" + util.RandomString(6)
		initialReplicas = 2
	})

	AfterEach(func() {
		defer cancelWatches()
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

			By("Waiting for workload cluster API server to be ready")
			waitForWorkerAPIServerReady(namespace, clusterName)

			By("Deploying CNI to created Cluster")
			deployCNI(testTmpDir, namespace, clusterName, *cniManifests)
			waitForMachineNodeReady(namespace, clusterName, machineName)

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

			By("Creating the MachineDeployment")
			createMachineDeployment(namespace, clusterName, machineDeploymentName, awsMachineTemplateName, mdBootstrapConfig, initialReplicas)

			By("Scale the MachineDeployment up")
			scaleMachineDeployment(namespace, machineDeploymentName, initialReplicas, scaleUpReplicas)

			By("Scale the MachineDeployment down")
			scaleMachineDeployment(namespace, machineDeploymentName, scaleUpReplicas, scaleDownReplicas)

			By("Scale the MachineDeployment up again")
			scaleMachineDeployment(namespace, machineDeploymentName, scaleDownReplicas, initialReplicas)

			By("Deleting the Cluster")
			deleteCluster(namespace, clusterName)
		})
	})
})

func watchEvents(ctx context.Context, namespace string) {
	logFile := path.Join(artifactPath, "resources", namespace, "events.log")
	fmt.Fprintf(GinkgoWriter, "Creating directory: %s\n", filepath.Dir(logFile))
	Expect(os.MkdirAll(filepath.Dir(logFile), 0755)).To(Succeed())

	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	Expect(err).NotTo(HaveOccurred())
	defer f.Close()

	informerFactory := informers.NewSharedInformerFactoryWithOptions(
		clientSet,
		10*time.Minute,
		informers.WithNamespace(namespace),
	)
	eventInformer := informerFactory.Core().V1().Events().Informer()
	eventInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			e := obj.(*corev1.Event)
			f.WriteString(fmt.Sprintf("[New Event] %s/%s\n\tresource: %s/%s/%s\n\treason: %s\n\tmessage: %s\n\tfull: %#v\n",
				e.Namespace, e.Name, e.InvolvedObject.APIVersion, e.InvolvedObject.Kind, e.InvolvedObject.Name, e.Reason, e.Message, e))
		},
		UpdateFunc: func(_, obj interface{}) {
			e := obj.(*corev1.Event)
			f.WriteString(fmt.Sprintf("[Updated Event] %s/%s\n\tresource: %s/%s/%s\n\treason: %s\n\tmessage: %s\n\tfull: %#v\n",
				e.Namespace, e.Name, e.InvolvedObject.APIVersion, e.InvolvedObject.Kind, e.InvolvedObject.Name, e.Reason, e.Message, e))
		},
		DeleteFunc: func(obj interface{}) {},
	})

	stopInformer := make(chan struct{})
	defer close(stopInformer)
	informerFactory.Start(stopInformer)
	<-ctx.Done()
	stopInformer <- struct{}{}
}

func scaleMachineDeployment(namespace, machineDeployment string, replicasCurrent int32, replicasProposed int32) {
	fmt.Fprintf(GinkgoWriter, "Scaling MachineDeployment from %d to %d\n", replicasCurrent, replicasProposed)
	deployment := &clusterv1.MachineDeployment{}
	Expect(kindClient.Get(context.TODO(), crclient.ObjectKey{Namespace: namespace, Name: machineDeployment}, deployment))
	deployment.Spec.Replicas = &replicasProposed
	Expect(kindClient.Update(context.TODO(), deployment)).NotTo(HaveOccurred())
	waitForMachinesCountMatch(namespace, machineDeployment, replicasCurrent, replicasProposed)
	waitForMachineDeploymentRunning(namespace, machineDeployment)
}

func createMachineDeployment(namespace, clusterName, machineDeploymentName, awsMachineTemplateName, bootstrapConfigName string, replicas int32) {
	fmt.Fprintf(GinkgoWriter, "Creating MachineDeployment in namespace %s with name %s\n", namespace, machineDeploymentName)
	makeAWSMachineTemplate(namespace, awsMachineTemplateName)
	makeJoinBootstrapConfigTemplate(namespace, bootstrapConfigName)
	makeMachineDeployment(namespace, machineDeploymentName, awsMachineTemplateName, bootstrapConfigName, clusterName, replicas)
	waitForMachinesCountMatch(namespace, machineDeploymentName, replicas, replicas)
	waitForMachineDeploymentRunning(namespace, machineDeploymentName)
}

func waitForMachinesCountMatch(namespace, machineDeploymentName string, replicasCurrent int32, replicasProposed int32) {
	fmt.Fprintf(GinkgoWriter, "Ensuring Machine count matched to %d \n", replicasProposed)
	machineDeployment := &clusterv1.MachineDeployment{}
	if err := kindClient.Get(context.TODO(), crclient.ObjectKey{Namespace: namespace, Name: machineDeploymentName}, machineDeployment); err != nil {
		return
	}
	scalingDown := replicasCurrent > replicasProposed
	Eventually(
		func() (int32, error) {
			machineList := &clusterv1.MachineList{}
			selector, err := metav1.LabelSelectorAsMap(&machineDeployment.Spec.Selector)
			if err != nil {
				return -1, err
			}
			if err := kindClient.List(context.TODO(), machineList, crclient.InNamespace(namespace), crclient.MatchingLabels(selector)); err != nil {
				return -1, err
			} else {
				var runningMachines int32
				var deletingMachines int32
				for _, item := range machineList.Items {
					if string(clusterv1.MachinePhaseRunning) == item.Status.Phase {
						runningMachines += 1
					}
					if string(clusterv1.MachinePhaseDeleting) == item.Status.Phase {
						deletingMachines += 1
					}
				}
				if scalingDown {
					if deletingMachines > 0 {
						return -1, nil
					}
				}
				return runningMachines, nil
			}

		},
		20*time.Minute, 30*time.Second,
	).Should(Equal(replicasProposed))
}

func waitForMachineDeploymentRunning(namespace, machineDeploymentName string) {
	Eventually(
		func() (bool, error) {
			machineDeployment := &clusterv1.MachineDeployment{}
			if err := kindClient.Get(context.TODO(), crclient.ObjectKey{Namespace: namespace, Name: machineDeploymentName}, machineDeployment); err != nil {
				return false, err
			}
			return *machineDeployment.Spec.Replicas == machineDeployment.Status.ReadyReplicas, nil
		},
		5*time.Minute, 15*time.Second,
	).Should(BeTrue())
}

func makeMachineDeployment(namespace, mdName, awsMachineTemplateName, bootstrapConfigName, clusterName string, replicas int32) {
	fmt.Fprintf(GinkgoWriter, "Creating MachineDeployment %s/%s\n", namespace, mdName)
	k8s_version := "v1.15.3"
	machineDeployment := &clusterv1.MachineDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mdName,
			Namespace: namespace,
			Labels: map[string]string{
				"cluster.x-k8s.io/cluster-name": clusterName,
				"nodepool":                      mdName,
			},
		},
		Spec: clusterv1.MachineDeploymentSpec{
			Replicas: &replicas,
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"cluster.x-k8s.io/cluster-name": clusterName,
					"nodepool":                      mdName,
				},
			},
			ClusterName: clusterName,
			Template: clusterv1.MachineTemplateSpec{
				ObjectMeta: clusterv1.ObjectMeta{
					Labels: map[string]string{
						"cluster.x-k8s.io/cluster-name": clusterName,
						"nodepool":                      mdName,
					},
				},
				Spec: clusterv1.MachineSpec{
					ClusterName: clusterName,
					Bootstrap: clusterv1.Bootstrap{
						ConfigRef: &corev1.ObjectReference{
							Kind:       "KubeadmConfigTemplate",
							APIVersion: bootstrapv1.GroupVersion.String(),
							Name:       bootstrapConfigName,
							Namespace:  namespace,
						},
					},
					InfrastructureRef: corev1.ObjectReference{
						Kind:       "AWSMachineTemplate",
						APIVersion: infrav1.GroupVersion.String(),
						Name:       awsMachineTemplateName,
						Namespace:  namespace,
					},
					Version: &k8s_version,
				},
			},
		},
	}
	Expect(kindClient.Create(context.TODO(), machineDeployment)).NotTo(HaveOccurred())
}

func makeJoinBootstrapConfigTemplate(namespace, name string) {
	fmt.Fprintf(GinkgoWriter, "Creating Join KubeadmConfigTemplate %s/%s\n", namespace, name)
	configTemplate := &bootstrapv1.KubeadmConfigTemplate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: bootstrapv1.KubeadmConfigTemplateSpec{
			Template: bootstrapv1.KubeadmConfigTemplateResource{
				Spec: bootstrapv1.KubeadmConfigSpec{
					JoinConfiguration: &kubeadmv1beta1.JoinConfiguration{
						NodeRegistration: kubeadmv1beta1.NodeRegistrationOptions{
							Name:             "{{ ds.meta_data.local_hostname }}",
							KubeletExtraArgs: map[string]string{"cloud-provider": "aws"},
						},
					},
				},
			},
		},
	}
	Expect(kindClient.Create(context.TODO(), configTemplate)).NotTo(HaveOccurred())
}

func makeAWSMachineTemplate(namespace, name string) {
	fmt.Fprintf(GinkgoWriter, "Creating AWSMachineTemplate %s/%s\n", namespace, name)
	awsMachine := &infrav1.AWSMachineTemplate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: infrav1.AWSMachineTemplateSpec{
			Template: infrav1.AWSMachineTemplateResource{
				Spec: infrav1.AWSMachineSpec{
					InstanceType:       "t3.large",
					IAMInstanceProfile: "nodes.cluster-api-provider-aws.sigs.k8s.io",
					SSHKeyName:         keyPairName,
				},
			},
		},
	}
	Expect(kindClient.Create(context.TODO(), awsMachine)).NotTo(HaveOccurred())
}

func deployCNI(tmpDir, namespace, clusterName, manifestPath string) {
	kubeConfigData, err := kubeconfig.FromSecret(context.TODO(), kindClient, crclient.ObjectKey{Name: clusterName, Namespace: namespace})
	Expect(err).NotTo(HaveOccurred())
	kubeConfigPath := path.Join(tmpDir, clusterName+".kubeconfig")
	Expect(ioutil.WriteFile(kubeConfigPath, kubeConfigData, 0640)).To(Succeed())

	Expect(exec.Command(
		*kubectlBinary,
		"apply",
		"--kubeconfig",
		kubeConfigPath,
		"-f",
		manifestPath,
	).Run()).To(Succeed())
}

func waitForWorkerAPIServerReady(namespace, clusterName string) {
	kubeConfigData, err := kubeconfig.FromSecret(context.TODO(), kindClient, crclient.ObjectKey{Name: clusterName, Namespace: namespace})
	Expect(err).NotTo(HaveOccurred())

	restConfig, err := clientcmd.RESTConfigFromKubeConfig(kubeConfigData)
	Expect(err).NotTo(HaveOccurred())

	nodeClient, err := crclient.New(restConfig, crclient.Options{})
	Expect(err).NotTo(HaveOccurred())

	Eventually(
		func() bool {
			nodes := &corev1.NodeList{}
			if err := nodeClient.List(context.TODO(), nodes); err != nil {
				fmt.Fprintf(GinkgoWriter, "Error retrieving nodes: %v", err)
				return false
			}
			return true
		},
		5*time.Minute, 15*time.Second,
	).Should(BeTrue())
}

func waitForMachineNodeReady(namespace, clusterName, name string) {
	machine := &clusterv1.Machine{}
	Expect(kindClient.Get(context.TODO(), crclient.ObjectKey{Namespace: namespace, Name: name}, machine)).To(Succeed())

	nodeName := machine.Status.NodeRef.Name

	kubeConfigData, err := kubeconfig.FromSecret(context.TODO(), kindClient, crclient.ObjectKey{Name: clusterName, Namespace: namespace})
	Expect(err).NotTo(HaveOccurred())

	restConfig, err := clientcmd.RESTConfigFromKubeConfig(kubeConfigData)
	Expect(err).NotTo(HaveOccurred())

	nodeClient, err := crclient.New(restConfig, crclient.Options{})
	Expect(err).NotTo(HaveOccurred())

	Eventually(
		func() bool {
			node := &corev1.Node{}
			if err := nodeClient.Get(context.TODO(), crclient.ObjectKey{Name: nodeName}, node); err != nil {
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
	waitForMachineNodeReady(namespace, clusterName, machineName)
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
			if err := kindClient.Get(context.TODO(), crclient.ObjectKey{Namespace: namespace, Name: name}, cluster); err != nil {
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
			if err := kindClient.Get(context.TODO(), crclient.ObjectKey{Namespace: namespace, Name: name}, machine); err != nil {
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
			if err := kindClient.Get(context.TODO(), crclient.ObjectKey{Namespace: namespace, Name: name}, cluster); err != nil {
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
			if err := kindClient.Get(context.TODO(), crclient.ObjectKey{Namespace: namespace, Name: name}, awsMachine); err != nil {
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
			if err := kindClient.Get(context.TODO(), crclient.ObjectKey{Namespace: namespace, Name: name}, cluster); err != nil {
				return false, err
			}
			return cluster.Status.InfrastructureReady, nil
		},
		15*time.Minute, 15*time.Second,
	).Should(BeTrue())
}

func waitForMachineBootstrapReady(namespace, name string) {
	fmt.Fprintf(GinkgoWriter, "Ensuring Machine %s/%s has bootstrapReady\n", namespace, name)
	Eventually(
		func() (bool, error) {
			machine := &clusterv1.Machine{}
			if err := kindClient.Get(context.TODO(), crclient.ObjectKey{Namespace: namespace, Name: name}, machine); err != nil {
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
			if err := kindClient.Get(context.TODO(), crclient.ObjectKey{Namespace: namespace, Name: name}, awsMachine); err != nil {
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
				clusterv1.MachineControlPlaneLabelName: "true",
			},
		},
		Spec: clusterv1.MachineSpec{
			ClusterName: clusterName,
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
					Name:             "{{ ds.meta_data.local_hostname }}",
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
					Name:             "{{ ds.meta_data.local_hostname }}",
					KubeletExtraArgs: map[string]string{"cloud-provider": "aws"},
				},
			},
			Files: []bootstrapv1.File{
				{
					Path:    "/tmp/userdata-length-test.txt",
					Content: apirand.String(20000),
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
			Region:     region,
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
