// +build e2e

/*
Copyright 2020 The Kubernetes Authors.

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

package e2e_new

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	apimachinerytypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/test/e2e_new/repeaters"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	cabpkv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1alpha3"
	kubeadmv1beta1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/types/v1beta1"
	kcpv1 "sigs.k8s.io/cluster-api/controlplane/kubeadm/api/v1alpha3"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/yaml"
)

var _ = Describe("scale tests", func() {
	var (
		namespace *corev1.Namespace
		ctx       context.Context
		//cancel    context.CancelFunc
	)

	BeforeEach(func() {
		Expect(bootstrapClusterProxy).ToNot(BeNil(), "Invalid argument. BootstrapClusterProxy can't be nil")
		ctx = context.TODO()
		// Setup a Namespace where to host objects for this spec and create a watcher for the namespace events.
		namespace = setupSpecNamespace(ctx, "scale-tests", bootstrapClusterProxy, artifactFolder)
	})

	It("should complete a scale test", func() {
		scaleArtifactsDir := path.Join(artifactFolder, "scale-test")
		Expect(os.MkdirAll(scaleArtifactsDir, 0o750)).To(Succeed(), "Invalid test suite argument. Can't create scale-tests-folder %q", scaleArtifactsDir)

		By("Creating initial cluster")
		clusterName := fmt.Sprintf("cluster-%s", util.RandomString(6))
		workloadClusterTemplate := clusterctl.ConfigCluster(ctx, clusterctl.ConfigClusterInput{
			LogFolder:                filepath.Join(artifactFolder, "clusters", bootstrapClusterProxy.GetName()),
			ClusterctlConfigPath:     clusterctlConfigPath,
			KubeconfigPath:           bootstrapClusterProxy.GetKubeconfigPath(),
			InfrastructureProvider:   clusterctl.DefaultInfrastructureProvider,
			Flavor:                   clusterctl.DefaultFlavor,
			Namespace:                namespace.Name,
			ClusterName:              clusterName,
			KubernetesVersion:        e2eConfig.GetVariable(KubernetesVersion),
			ControlPlaneMachineCount: pointer.Int64Ptr(1),
			WorkerMachineCount:       pointer.Int64Ptr(0),
		})
		Expect(bootstrapClusterProxy.Apply(ctx, workloadClusterTemplate)).ShouldNot(HaveOccurred())

		By("Waiting for initial cluster to reach infrastructure ready")
		Eventually(func() bool {
			cluster := &clusterv1.Cluster{}
			if err := bootstrapClusterProxy.GetClient().Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: namespace.Name, Name: clusterName}, cluster); nil == err {
				if cluster.Status.InfrastructureReady {
					return true
				}
			}
			return false
		}, e2eConfig.GetIntervals("", "wait-cluster")...).Should(Equal(true))

		infracluster := &infrav1.AWSCluster{}
		err := bootstrapClusterProxy.GetClient().Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: namespace.Name, Name: clusterName}, infracluster)
		Expect(err).NotTo(HaveOccurred())

		networkSpec := infracluster.Spec.NetworkSpec.DeepCopy()
		clusterIndex := int64(-1)

		noClusters, err := strconv.ParseInt(e2eConfig.GetVariable("SCALE_TEST_NO_CLUSTERS"), 0, 0)
		Expect(err).NotTo(HaveOccurred())

		repeaters.AtOnce(uint64(noClusters), func(wg *sync.WaitGroup) {
			defer GinkgoRecover()
			i := atomic.AddInt64(&clusterIndex, 1)
			log := log.WithValues("cluster-index", i)
			log.Info("Generating cluster")
			name, clusterResources := generateCluster(namespace.Name, *networkSpec)
			log = log.WithValues("namespace", namespace.Name, "name", name)
			data, err := yaml.Marshal(clusterResources)
			if err != nil {
				log.Error(err, "error marshalling generated cluster")
				return
			}
			if err := ioutil.WriteFile(path.Join(scaleArtifactsDir, fmt.Sprintf("%d.%s.yaml", i, name)), data, 0o640); err != nil {
				log.Error(err, "error saving resource for debugging. continuing")
			}
			if err := bootstrapClusterProxy.Apply(ctx, data); err != nil {
				log.Error(err, "error applying cluster resources")
				return
			}

			log.Info("Waiting for the cluster infrastructure to be provisioned")
			cluster := framework.DiscoveryAndWaitForCluster(ctx, framework.DiscoveryAndWaitForClusterInput{
				Getter:    bootstrapClusterProxy.GetClient(),
				Namespace: namespace.Name,
				Name:      name,
			}, e2eConfig.GetIntervals("", "wait-cluster")...)

			log.Info("Waiting for control plane to be initialized")
			controlPlane := framework.DiscoveryAndWaitForControlPlaneInitialized(ctx, framework.DiscoveryAndWaitForControlPlaneInitializedInput{
				Lister:  bootstrapClusterProxy.GetClient(),
				Cluster: cluster,
			}, e2eConfig.GetIntervals("", "wait-cluster")...)

			log.Info("Installing a CNI plugin to the workload cluster")
			workloadCluster := bootstrapClusterProxy.GetWorkloadCluster(context.TODO(), cluster.Namespace, cluster.Name)

			cniYaml, err := ioutil.ReadFile(e2eConfig.GetVariable("CNI"))
			if err != nil {
				log.Error(err, "Unable to read CNI manifest")
			}

			if err := workloadCluster.Apply(context.TODO(), cniYaml); err != nil {
				log.Error(err, "Unable to apply CNI")
			}

			log.Info("Waiting for control plane to be ready")
			framework.WaitForControlPlaneAndMachinesReady(ctx, framework.WaitForControlPlaneAndMachinesReadyInput{
				GetLister:    bootstrapClusterProxy.GetClient(),
				Cluster:      cluster,
				ControlPlane: controlPlane,
			}, e2eConfig.GetIntervals("", "wait-control-plane")...)

			log.Info("Waiting for the worker machines to be provisioned")
			framework.DiscoveryAndWaitForMachineDeployments(ctx, framework.DiscoveryAndWaitForMachineDeploymentsInput{
				Lister:  bootstrapClusterProxy.GetClient(),
				Cluster: cluster,
			}, e2eConfig.GetIntervals("", "wait-worker-nodes")...)

		})

	})

	AfterEach(func() {
		// Dumps all the resources in the spec namespace, then cleanups the cluster object and the spec namespace itself.
		dumpSpecResourcesAndCleanup(ctx, "", bootstrapClusterProxy, artifactFolder, namespace, e2eConfig.GetIntervals, skipCleanup)
	})
})

func generateCluster(namespace string, networkSpec infrav1.NetworkSpec) (string, corev1.List) {
	name := fmt.Sprintf("cluster-%s", util.RandomString(6))
	cluster, infraCluster, kcp, kcpInfraTemplate := generateClusterResource(name, namespace, networkSpec)
	md, mdTemplate, mdKubeadmTemplate := generateMachineDeployment(cluster)
	return name, corev1.List{
		TypeMeta: metav1.TypeMeta{
			Kind:       "List",
			APIVersion: "v1",
		},
		Items: []runtime.RawExtension{
			{Object: &cluster},
			{Object: &infraCluster},
			{Object: &kcpInfraTemplate},
			{Object: &kcp},
			{Object: &md},
			{Object: &mdTemplate},
			{Object: &mdKubeadmTemplate},
		},
	}
}

func generateClusterResource(name, namespace string, networkSpec infrav1.NetworkSpec) (clusterv1.Cluster, infrav1.AWSCluster, kcpv1.KubeadmControlPlane, infrav1.AWSMachineTemplate) {
	infraCluster := infrav1.AWSCluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AWSCluster",
			APIVersion: infrav1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: infrav1.AWSClusterSpec{
			Region:      os.Getenv("AWS_REGION"),
			SSHKeyName:  pointer.StringPtr(os.Getenv("AWS_SSH_KEY_NAME")),
			NetworkSpec: networkSpec,
		},
	}
	kcp, kcpInfraTemplate := generateKcp(name, namespace)
	cluster := clusterv1.Cluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Cluster",
			APIVersion: clusterv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      infraCluster.GetName(),
			Namespace: infraCluster.GetNamespace(),
		},
		Spec: clusterv1.ClusterSpec{
			ClusterNetwork: &clusterv1.ClusterNetwork{
				Pods: &clusterv1.NetworkRanges{CIDRBlocks: []string{"192.168.0.0/16"}},
			},
			ControlPlaneRef: &corev1.ObjectReference{
				APIVersion: kcpv1.GroupVersion.String(),
				Kind:       "KubeadmControlPlane",
				Name:       kcp.GetName(),
				Namespace:  kcp.GetNamespace(),
			},
			InfrastructureRef: &corev1.ObjectReference{
				APIVersion: infrav1.GroupVersion.String(),
				Kind:       "AWSCluster",
				Name:       infraCluster.GetName(),
				Namespace:  infraCluster.GetNamespace(),
			},
		},
	}
	return cluster, infraCluster, kcp, kcpInfraTemplate
}

func generateInfraTemplate(name, namespace, instanceType string) infrav1.AWSMachineTemplate {
	return infrav1.AWSMachineTemplate{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AWSMachineTemplate",
			APIVersion: infrav1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: infrav1.AWSMachineTemplateSpec{
			Template: infrav1.AWSMachineTemplateResource{
				Spec: infrav1.AWSMachineSpec{
					InstanceType:       instanceType,
					IAMInstanceProfile: "control-plane.cluster-api-provider-aws.sigs.k8s.io",
					SSHKeyName:         pointer.StringPtr(e2eConfig.GetVariable("AWS_SSH_KEY_NAME")),
				},
			},
		},
	}
}

func generateKubeadmTemplate(name, namespace string) cabpkv1.KubeadmConfigTemplate {
	return cabpkv1.KubeadmConfigTemplate{
		TypeMeta: metav1.TypeMeta{
			Kind:       "KubeadmConfigTemplate",
			APIVersion: cabpkv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: cabpkv1.KubeadmConfigTemplateSpec{
			Template: cabpkv1.KubeadmConfigTemplateResource{
				Spec: cabpkv1.KubeadmConfigSpec{
					JoinConfiguration: &kubeadmv1beta1.JoinConfiguration{
						NodeRegistration: nodeRegistrationOptions(),
					},
				},
			},
		},
	}
}

func generateMachineDeployment(cluster clusterv1.Cluster) (clusterv1.MachineDeployment, infrav1.AWSMachineTemplate, cabpkv1.KubeadmConfigTemplate) {
	name := cluster.GetName() + "-md-0"
	infraTemplate := generateInfraTemplate(name, cluster.GetNamespace(), e2eConfig.GetVariable("AWS_NODE_MACHINE_TYPE"))
	kubeadmTemplate := generateKubeadmTemplate(name, cluster.GetNamespace())
	md := clusterv1.MachineDeployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "MachineDeployment",
			APIVersion: clusterv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: infraTemplate.GetNamespace(),
		},
		Spec: clusterv1.MachineDeploymentSpec{
			ClusterName: cluster.GetName(),
			Template: clusterv1.MachineTemplateSpec{
				Spec: clusterv1.MachineSpec{
					ClusterName: cluster.GetName(),
					Version:     pointer.StringPtr(e2eConfig.GetVariable("KUBERNETES_VERSION")),
					Bootstrap: clusterv1.Bootstrap{
						ConfigRef: &corev1.ObjectReference{
							APIVersion: cabpkv1.GroupVersion.String(),
							Kind:       "KubeadmConfigTemplate",
							Name:       kubeadmTemplate.GetName(),
							Namespace:  kubeadmTemplate.GetNamespace(),
						},
					},
					InfrastructureRef: corev1.ObjectReference{
						APIVersion: infrav1.GroupVersion.String(),
						Kind:       "AWSMachineTemplate",
						Name:       infraTemplate.GetName(),
						Namespace:  infraTemplate.GetNamespace(),
					},
				},
			},
		},
	}
	return md, infraTemplate, kubeadmTemplate
}

func cloudProviderArgs() map[string]string {
	return map[string]string{
		"cloud-provider": "aws",
	}
}

func nodeRegistrationOptions() kubeadmv1beta1.NodeRegistrationOptions {
	return kubeadmv1beta1.NodeRegistrationOptions{
		Name:             "{{ ds.meta_data.local_hostname }}",
		KubeletExtraArgs: cloudProviderArgs(),
	}
}

func generateKcp(name, namespace string) (kcpv1.KubeadmControlPlane, infrav1.AWSMachineTemplate) {
	infraTemplate := generateInfraTemplate(name+"-control-plane", namespace, e2eConfig.GetVariable("AWS_CONTROL_PLANE_MACHINE_TYPE"))
	kcp := kcpv1.KubeadmControlPlane{
		TypeMeta: metav1.TypeMeta{
			Kind:       "KubeadmControlPlane",
			APIVersion: kcpv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      infraTemplate.GetName(),
			Namespace: infraTemplate.GetNamespace(),
		},
		Spec: kcpv1.KubeadmControlPlaneSpec{
			Replicas: pointer.Int32Ptr(3),
			Version:  e2eConfig.GetVariable(KubernetesVersion),
			InfrastructureTemplate: corev1.ObjectReference{
				APIVersion: infrav1.GroupVersion.String(),
				Kind:       "AWSMachineTemplate",
				Name:       infraTemplate.GetName(),
				Namespace:  infraTemplate.GetNamespace(),
			},
			KubeadmConfigSpec: cabpkv1.KubeadmConfigSpec{
				ClusterConfiguration: &kubeadmv1beta1.ClusterConfiguration{
					APIServer: kubeadmv1beta1.APIServer{
						ControlPlaneComponent: kubeadmv1beta1.ControlPlaneComponent{
							ExtraArgs: cloudProviderArgs(),
						},
					},
				},
				InitConfiguration: &kubeadmv1beta1.InitConfiguration{
					NodeRegistration: nodeRegistrationOptions(),
				},
				JoinConfiguration: &kubeadmv1beta1.JoinConfiguration{
					NodeRegistration: nodeRegistrationOptions(),
				},
			},
		},
	}
	return kcp, infraTemplate
}
