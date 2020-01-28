// +build e2e

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

package e2e_test

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"text/template"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	"github.com/vmware-tanzu/sonobuoy/pkg/client"
	sonodynamic "github.com/vmware-tanzu/sonobuoy/pkg/dynamic"
	"golang.org/x/sync/errgroup"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/kubeconfig"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("conformance tests", func() {
	var (
		setup testSetup
	)

	BeforeEach(func() {
		var err error
		setup = testSetup{}
		setup.testTmpDir, err = ioutil.TempDir(suiteTmpDir, "conformance-test")
		Expect(err).NotTo(HaveOccurred())

		setup.namespace = "conformance-" + util.RandomString(6)
		createNamespace(setup.namespace)

		setup.clusterName = "conformance-" + util.RandomString(6)
		setup.awsClusterName = "conformance-infra-" + util.RandomString(6)
		setup.cpMachinePrefix = "conformance-" + util.RandomString(6)
		setup.cpAWSMachinePrefix = "conformance-infra-" + util.RandomString(6)
		setup.cpBootstrapConfigPrefix = "conformance-boot-" + util.RandomString(6)
		setup.mdBootstrapConfig = "conformance-boot-md" + util.RandomString(6)
		setup.machineDeploymentName = "conformance-capa-md" + util.RandomString(6)
		setup.awsMachineTemplateName = "conformance-infra-capa-mt" + util.RandomString(6)
		setup.initialReplicas = 2
		setup.instanceType = "t3.large"
		setup.multipleAZ = false
	})

	Describe("conformance on workload cluster", func() {
		It("It should pass k8s certified-conformance tests", func() {
			By("Creating a cluster with single control plane")
			makeSingleControlPlaneCluster(setup)

			By("Deploying a MachineDeployment")
			createMachineDeployment(setup)

			By("Running conformance on the workload cluster")
			err := runConformance(setup.testTmpDir, setup.namespace, setup.clusterName)
			Expect(err).NotTo(HaveOccurred())

			By("Deleting the Cluster")
			deleteCluster(setup.namespace, setup.clusterName)
		})
	})
})

type sonobuoyConfig struct {
	SonobuoyVersion string
	K8sVersion      string
}

func runConformance(tmpDir, namespace, clusterName string) error {
	cluster := crclient.ObjectKey{
		Namespace: namespace,
		Name:      clusterName,
	}
	kubeConfigData, err := kubeconfig.FromSecret(context.TODO(), kindClient, cluster)
	if err != nil {
		return errors.Wrap(err, "couldn't get kubeconfig of workload cluster")
	}

	restConfig, err := clientcmd.RESTConfigFromKubeConfig(kubeConfigData)
	if err != nil {
		return errors.Wrap(err, "couldn't get rest config from kubeconfig")
	}

	sonobuoyKubeCli, err := sonodynamic.NewAPIHelperFromRESTConfig(restConfig)
	if err != nil {
		return errors.Wrap(err, "couldn't create sonobuoy client")
	}

	sonobuoyClient, err := client.NewSonobuoyClient(restConfig, sonobuoyKubeCli)
	if err != nil {
		return errors.Wrap(err, "couldn't create sonobuoy client")
	}

	k8sClient, err := sonobuoyClient.Client()
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve k8s client")
	}

	apiVersion, err := k8sClient.Discovery().ServerVersion()
	if err != nil {
		return errors.Wrap(err, "couldn't get workload cluster's k8s version")
	}

	tmpl, err := template.ParseFiles("sonobuoy-config-tmpl.yaml")
	if err != nil {
		return errors.Wrap(err, "couldn't parse sonobuoy config template")
	}

	sonobuoyConfigPath := path.Join(tmpDir, clusterName+"-sonobuoy-config.yaml")
	fileP, err := os.Create(sonobuoyConfigPath)
	if err != nil {
		return errors.Wrap(err, "couldn't create sonobuoy config file")
	}

	sbConfig := sonobuoyConfig{SonobuoyVersion: *sonobuoyVersion, K8sVersion: apiVersion.GitVersion}
	err = tmpl.Execute(fileP, sbConfig)
	if err != nil {
		return errors.Wrap(err, "couldn't execute template")
	}

	runConfig := &client.RunConfig{
		Wait:    time.Duration(4) * time.Hour,
		GenFile: sonobuoyConfigPath,
	}
	if err := sonobuoyClient.Run(runConfig); err != nil {
		return errors.Wrap(err, "error attempting to run sonobuoy")
	}
	reader, ec, err := sonobuoyClient.RetrieveResults(&client.RetrieveConfig{Namespace: "sonobuoy"})
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve sonobuoy results")
	}
	var fileName string
	eg := &errgroup.Group{}
	eg.Go(func() error { return <-ec })
	eg.Go(func() error {
		filesCreated, err := client.UntarAll(reader, ".", "")
		if err != nil {
			return errors.Wrap(err, "couldn't untar sonobuoy results")
		}
		fmt.Fprintf(GinkgoWriter, "Files created by sonobuoy: ")
		for _, fileName = range filesCreated {
			fmt.Fprintf(GinkgoWriter, "%s\n", fileName)
		}
		return nil
	})

	err = eg.Wait()
	if err != nil {
		return errors.Wrap(err, "error retrieving results")
	}

	_, err = exec.Command("tar", "-C", "tmp/sonobuoy", "-xf", fileName).Output()
	if err != nil {
		fmt.Fprintf(GinkgoWriter, "untar %s failed\n", fileName)
	} else {
		src, err := os.Open("tmp/sonobuoy/plugins/e2e/results/global/junit_01.xml")
		defer src.Close()
		if err == nil {
			dst, err := os.Create(path.Join(artifactPath, "junit.k8s_conf.xml"))
			defer dst.Close()
			if err == nil {
				_, err = io.Copy(dst, src)
				if err != nil {
					fmt.Fprintf(GinkgoWriter, "couldn't fetch junit.k8s_conf.xml %v", err)
				}
			}
		}
	}

	err = sonobuoyClient.Delete(&client.DeleteConfig{
		Namespace:  "sonobuoy",
		EnableRBAC: true,
		DeleteAll:  true,
		Wait:       15 * time.Minute,
	})
	if err != nil {
		return errors.Wrap(err, "error deleting sonobuoy")
	}
	return nil
}
