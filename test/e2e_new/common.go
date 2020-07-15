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

	. "github.com/onsi/ginkgo"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	corev1 "k8s.io/api/core/v1"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

var (
	namespaces map[*corev1.Namespace]context.CancelFunc
)

func setupSpecNamespace(ctx context.Context, specName string, clusterProxy framework.ClusterProxy, artifactFolder string) *corev1.Namespace {
	Byf("Creating a namespace for hosting the %q test spec", specName)
	namespace, cancelWatches := framework.CreateNamespaceAndWatchEvents(ctx, framework.CreateNamespaceAndWatchEventsInput{
		Creator:   clusterProxy.GetClient(),
		ClientSet: clusterProxy.GetClientSet(),
		Name:      fmt.Sprintf("%s-%s", specName, util.RandomString(6)),
		LogFolder: filepath.Join(artifactFolder, "clusters", clusterProxy.GetName()),
	})

	namespaces[namespace] = cancelWatches

	return namespace
}

func dumpSpecResourcesAndCleanup(ctx context.Context, specName string, clusterProxy framework.ClusterProxy, artifactFolder string, namespace *corev1.Namespace, intervalsGetter func(spec, key string) []interface{}, skipCleanup bool) {
	Byf("Dumping all the Cluster API resources in the %q namespace", namespace.Name)
	// Dump all Cluster API related resources to artifacts before deleting them.
	cancelWatches := namespaces[namespace]
	dumpSpecResources(ctx, clusterProxy, artifactFolder, namespace)
	Byf("Dumping all EC2 instances in the %q namespace", namespace.Name)
	dumpMachines(ctx, clusterProxy, namespace, artifactFolder)
	if !skipCleanup {
		framework.DeleteAllClustersAndWait(ctx, framework.DeleteAllClustersAndWaitInput{
			Client:    clusterProxy.GetClient(),
			Namespace: namespace.Name,
		}, intervalsGetter(specName, "wait-delete-cluster")...)

		Byf("Deleting namespace used for hosting the %q test spec", specName)
		framework.DeleteNamespace(ctx, framework.DeleteNamespaceInput{
			Deleter: clusterProxy.GetClient(),
			Name:    namespace.Name,
		})
	}
	cancelWatches()
	delete(namespaces, namespace)
}

func dumpMachines(ctx context.Context, clusterProxy framework.ClusterProxy, namespace *corev1.Namespace, artifactFolder string) {
	machines := machinesForSpec(ctx, clusterProxy, namespace)
	instances, err := allMachines(ctx)
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
		dumpMachine(ctx, m, instanceID, filepath.Join(artifactFolder, "clusters", clusterProxy.GetName()))
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

func dumpMachine(ctx context.Context, machine infrav1.AWSMachine, instanceID string, logPath string) {
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

func dumpSpecResources(ctx context.Context, clusterProxy framework.ClusterProxy, artifactFolder string, namespace *corev1.Namespace) {
	framework.DumpAllResources(ctx, framework.DumpAllResourcesInput{
		Lister:    clusterProxy.GetClient(),
		Namespace: namespace.Name,
		LogPath:   filepath.Join(artifactFolder, "clusters", clusterProxy.GetName(), "resources"),
	})
}

func Byf(format string, a ...interface{}) {
	By(fmt.Sprintf(format, a...))
}

func dumpServiceQuotas(ctx context.Context, artifactFolder string) error {
	if err := dumpServiceQuota(ctx, artifactFolder, "ec2"); err != nil {
		return err
	}
	if err := dumpServiceQuota(ctx, artifactFolder, "elasticloadbalancing"); err != nil {
		return err
	}
	return nil
}

func dumpServiceQuota(ctx context.Context, artifactFolder string, service string) error {
	s := servicequotas.New(awsSession)
	resp, err := s.ListServiceQuotasWithContext(ctx, &servicequotas.ListServiceQuotasInput{
		ServiceCode: aws.String(service),
	})
	if err != nil {
		return err
	}
	quotas := resp.Quotas
	for i := range quotas {
		quotas[i].QuotaArn = aws.String("")
	}
	data, err := yaml.Marshal(resp.Quotas)
	if err != nil {
		return err
	}
	filename := path.Join(artifactFolder, "service-quotas", service+".yaml")
	if err := os.MkdirAll(path.Dir(filename), 0o750); err != nil {
		return err
	}
	ioutil.WriteFile(filename, data, 0o640)
	return nil
}
