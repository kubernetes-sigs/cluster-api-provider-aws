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

package shared

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func SetupSpecNamespace(ctx context.Context, specName string, e2eCtx *E2EContext) *corev1.Namespace {
	Byf("Creating a namespace for hosting the %q test spec", specName)
	namespace, cancelWatches := framework.CreateNamespaceAndWatchEvents(ctx, framework.CreateNamespaceAndWatchEventsInput{
		Creator:   e2eCtx.BootstrapClusterProxy.GetClient(),
		ClientSet: e2eCtx.BootstrapClusterProxy.GetClientSet(),
		Name:      fmt.Sprintf("%s-%s", specName, util.RandomString(6)),
		LogFolder: filepath.Join(e2eCtx.ArtifactFolder, "clusters", e2eCtx.BootstrapClusterProxy.GetName()),
	})

	e2eCtx.Namespaces[namespace] = cancelWatches

	return namespace
}

// (ctx, "", e2eCtx.BootstrapClusterProxy, e2eCtx.ArtifactFolder, k, e2eCtx.E2EConfig.GetIntervals, e2eCtx.SkipCleanup)
func DumpSpecResourcesAndCleanup(ctx context.Context, specName string, namespace *corev1.Namespace, e2eCtx *E2EContext) {
	Byf("Dumping all the Cluster API resources in the %q namespace", namespace.Name)
	// Dump all Cluster API related resources to artifacts before deleting them.
	cancelWatches := e2eCtx.Namespaces[namespace]
	DumpSpecResources(ctx, e2eCtx, namespace)
	Byf("Dumping all EC2 instances in the %q namespace", namespace.Name)
	DumpMachines(ctx, e2eCtx, namespace)
	if !e2eCtx.SkipCleanup {
		framework.DeleteAllClustersAndWait(ctx, framework.DeleteAllClustersAndWaitInput{
			Client:    e2eCtx.BootstrapClusterProxy.GetClient(),
			Namespace: namespace.Name,
		}, e2eCtx.E2EConfig.GetIntervals(specName, "wait-delete-cluster")...)

		Byf("Deleting namespace used for hosting the %q test spec", specName)
		framework.DeleteNamespace(ctx, framework.DeleteNamespaceInput{
			Deleter: e2eCtx.BootstrapClusterProxy.GetClient(),
			Name:    namespace.Name,
		})
	}
	cancelWatches()
	delete(e2eCtx.Namespaces, namespace)
}

func DumpMachines(ctx context.Context, e2eCtx *E2EContext, namespace *corev1.Namespace) {
	machines := machinesForSpec(ctx, e2eCtx.BootstrapClusterProxy, namespace)
	instances, err := allMachines(ctx, e2eCtx)
	if err != nil {
		return
	}
	instanceID := ""
	for _, m := range machines.Items {
		for _, i := range instances {
			if i.name == m.Name {
				instanceID = i.instanceID
				break
			}
		}
		if instanceID == "" {
			return
		}
		DumpMachine(ctx, e2eCtx, m, instanceID)
	}
}

func machinesForSpec(ctx context.Context, clusterProxy framework.ClusterProxy, namespace *corev1.Namespace) *infrav1.AWSMachineList {
	lister := clusterProxy.GetClient()
	list := new(infrav1.AWSMachineList)
	if err := lister.List(ctx, list, client.InNamespace(namespace.GetName())); err != nil {
		fmt.Fprintln(GinkgoWriter, "couldn't find machines")
		return nil
	}
	return list
}

func DumpMachine(ctx context.Context, e2eCtx *E2EContext, machine infrav1.AWSMachine, instanceID string) {
	logPath := filepath.Join(e2eCtx.ArtifactFolder, "clusters", e2eCtx.BootstrapClusterProxy.GetName())
	machineLogBase := path.Join(logPath, "instances", machine.Namespace, machine.Name)
	metaLog := path.Join(machineLogBase, "instance.log")
	if err := os.MkdirAll(filepath.Dir(metaLog), 0750); err != nil {
		fmt.Fprintf(GinkgoWriter, "couldn't create directory for file: path=%s, err=%s", metaLog, err)
	}
	f, err := os.OpenFile(metaLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	fmt.Fprintf(f, "instance found: instance-id=%q\n", instanceID)
	commandsForMachine(
		ctx,
		e2eCtx,
		f,
		instanceID,
		[]command{
			{
				title: "systemd",
				cmd:   "journalctl --no-pager --output=short-precise",
			},
			{
				title: "kern",
				cmd:   "journalctl --no-pager --output=short-precise -k",
			},
			{
				title: "containerd-info",
				cmd:   "crictl info",
			},
			{
				title: "cloud-final",
				cmd:   "journalctl --no-pager -u cloud-final",
			},
			{
				title: "kubelet",
				cmd:   "journalctl --no-pager -u kubelet.service",
			},
			{
				title: "containerd",
				cmd:   "journalctl --no-pager -u containerd.service",
			},
		},
	)
}

func DumpSpecResources(ctx context.Context, e2eCtx *E2EContext, namespace *corev1.Namespace) {
	framework.DumpAllResources(ctx, framework.DumpAllResourcesInput{
		Lister:    e2eCtx.BootstrapClusterProxy.GetClient(),
		Namespace: namespace.Name,
		LogPath:   filepath.Join(e2eCtx.ArtifactFolder, "clusters", e2eCtx.BootstrapClusterProxy.GetName(), "resources"),
	})
}

func Byf(format string, a ...interface{}) {
	By(fmt.Sprintf(format, a...))
}

// LoadE2EConfig loads the e2econfig from the specified path
func LoadE2EConfig(configPath string) *clusterctl.E2EConfig {
	config := clusterctl.LoadE2EConfig(context.TODO(), clusterctl.LoadE2EConfigInput{ConfigPath: configPath})
	Expect(config).ToNot(BeNil(), "Failed to load E2E config from %s", configPath)
	// Read CNI file and set CNI_RESOURCES environmental variable
	Expect(config.Variables).To(HaveKey(CNIPath), "Missing %s variable in the config", CNIPath)
	clusterctl.SetCNIEnvVar(config.GetVariable(CNIPath), CNIResources)
	return config
}

// SetEnvVar sets an environment variable in the process. If marked private,
// the value is not printed.
func SetEnvVar(key, value string, private bool) {
	printableValue := "*******"
	if !private {
		printableValue = value
	}

	Byf("Setting environment variable: key=%s, value=%s", key, printableValue)
	os.Setenv(key, value)
}
