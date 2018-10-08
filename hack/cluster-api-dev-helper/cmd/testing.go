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

package cmd

import (
	"fmt"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

func defineTestingCmd(parent *cobra.Command) {
	newCmd := &cobra.Command{
		Use:   "testing",
		Short: "Test Cluster API",
		Long:  `Commands to help test Cluster API`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.Help())
		},
	}

	defineTestingStartCmd(newCmd)
	defineTestingLogsCmd(newCmd)
	defineTestingDestroyAllCmd(newCmd)
	defineTestingDeleteClustersCmd(newCmd)
	defineTestingDestroyClustersCmd(newCmd)
	defineTestingDeleteMachinesCmd(newCmd)
	defineTestingDestroyMachinesCmd(newCmd)
	defineTestingClusterLogs(newCmd)
	defineTestingMachineLogs(newCmd)
	defineTestingRestartControllerCmd(newCmd)
	defineTestingApplyControllerManifestsCmd(newCmd)

	parent.AddCommand(newCmd)
}

func defineTestingStartCmd(parent *cobra.Command) {
	newCmd := &cobra.Command{
		Use:   "start-with-existing-context",
		Short: "Start tests with an existing kubeconfig context",
		Long:  `Start tests with an existing kubeconfig context. Expects to be run in the repository using Go 1.11`,
		Run: func(cmd *cobra.Command, args []string) {
			var wg sync.WaitGroup
			a := runShell(`go run ./clusterctl create cluster -v2  \
				-m ./clusterctl/examples/aws/out/machines.yaml \
				-c ./clusterctl/examples/aws/out/cluster.yaml \
				-p ./clusterctl/examples/aws/out/provider-components.yaml \
				--provider aws \
				--existing-bootstrap-cluster-kubeconfig ~/.kube/config`)
			wg.Add(2)
			go clusterLogs(&wg)
			go machineLogs(&wg)
			a.Wait()
		},
	}

	parent.AddCommand(newCmd)
}

func defineTestingLogsCmd(parent *cobra.Command) {
	newCmd := &cobra.Command{
		Use:   "controller-logs",
		Short: "Start tailing controller logs",
		Long:  `Start tailing controller logs`,
		Run: func(cmd *cobra.Command, args []string) {
			var wg sync.WaitGroup
			wg.Add(2)
			go clusterLogs(&wg)
			go machineLogs(&wg)
			wg.Wait()
		},
	}
	parent.AddCommand(newCmd)
}

func defineTestingDestroyAllCmd(parent *cobra.Command) {
	newCmd := &cobra.Command{
		Use:   "destroy-all",
		Short: "Destroys Cluster API without finalizers",
		Long:  `Destroys Cluster API without finalizers`,
		Run: func(cmd *cobra.Command, args []string) {
			destroyMachines()
			destroyClusters()
			destroyControlPlane()
		},
	}
	parent.AddCommand(newCmd)
}

func defineTestingDeleteClustersCmd(parent *cobra.Command) {
	newCmd := &cobra.Command{
		Use:   "delete-clusters",
		Short: "Delete all clusters and tail logs",
		Long:  `Delete all clusters and tail logs`,
		Run: func(cmd *cobra.Command, args []string) {
			var wg sync.WaitGroup
			wg.Add(1)
			go clusterLogs(&wg)
			runShell("kubectl delete clusters --all --wait=true")
			wg.Wait()
		},
	}
	parent.AddCommand(newCmd)
}

func defineTestingDestroyClustersCmd(parent *cobra.Command) {
	newCmd := &cobra.Command{
		Use:   "destroy-clusters",
		Short: "Forcefully destroy clusters",
		Long:  `Forcefully destroy clusters`,
		Run: func(cmd *cobra.Command, args []string) {
			destroyClusters()
		},
	}
	parent.AddCommand(newCmd)
}

func defineTestingDeleteMachinesCmd(parent *cobra.Command) {
	newCmd := &cobra.Command{
		Use:   "delete-machines",
		Short: "Delete all machines and tail logs",
		Long:  `Delete all machines and tail logs`,
		Run: func(cmd *cobra.Command, args []string) {
			var wg sync.WaitGroup
			wg.Add(1)
			go machineLogs(&wg)
			runShell("kubectl delete machines --all --wait=true")
			wg.Wait()
		},
	}
	parent.AddCommand(newCmd)
}

func defineTestingDestroyMachinesCmd(parent *cobra.Command) {
	newCmd := &cobra.Command{
		Use:   "destroy-machines",
		Short: "Forcefully destroy machines",
		Long:  `Forcefully destroy machines`,
		Run: func(cmd *cobra.Command, args []string) {
			destroyMachines()
		},
	}
	parent.AddCommand(newCmd)
}

func defineTestingClusterLogs(parent *cobra.Command) {
	newCmd := &cobra.Command{
		Use:   "cluster-controller-logs",
		Short: "Tail cluster controller logs",
		Long:  `Tail cluster controller logs`,
		Run: func(cmd *cobra.Command, args []string) {
			var wg sync.WaitGroup
			wg.Add(1)
			go clusterLogs(&wg)
			wg.Wait()
		},
	}
	parent.AddCommand(newCmd)
}

func defineTestingMachineLogs(parent *cobra.Command) {
	newCmd := &cobra.Command{
		Use:   "machine-controller-logs",
		Short: "Tail cluster controller logs",
		Long:  `Tail cluster controller logs`,
		Run: func(cmd *cobra.Command, args []string) {
			var wg sync.WaitGroup
			wg.Add(1)
			go machineLogs(&wg)
			wg.Wait()
		},
	}
	parent.AddCommand(newCmd)
}

func defineTestingRestartControllerCmd(parent *cobra.Command) {
	newCmd := &cobra.Command{
		Use:   "restart-controllers",
		Short: "Restart controllers",
		Long:  `Restart controllers`,
		Run: func(cmd *cobra.Command, args []string) {
			restartControllers()
		},
	}
	parent.AddCommand(newCmd)
}

func defineTestingApplyControllerManifestsCmd(parent *cobra.Command) {
	newCmd := &cobra.Command{
		Use:   "apply-controller-manifests",
		Short: "Apply controller manifests",
		Long:  "Apply controller manifests",
		Run: func(cmd *cobra.Command, args []string) {
			runShellWithWait("kubectl apply -f clusterctl/examples/aws/out/provider-components.yaml")
		},
	}
	parent.AddCommand(newCmd)
}

func destroyMachines() {
	runShellWithWait("kubectl get machine -o name | xargs kubectl patch -p '{\"metadata\":{\"finalizers\":null}}'")
	runShellWithWait("kubectl delete machines --force=true --grace-period 0 --all --wait=true")
}

func destroyClusters() {
	runShellWithWait("kubectl patch cluster test1 -p '{\"metadata\":{\"finalizers\":null}}'")
	runShellWithWait("kubectl delete clusters --force=true --grace-period 0 --all --wait=true")
}

func destroyControlPlane() {
	runShellWithWait("kubectl delete deployment clusterapi-apiserver --force=true --grace-period 0 --wait=true")
	runShellWithWait("kubectl delete deployment clusterapi-controllers --force=true --grace-period 0 --wait=true")
	runShellWithWait("kubectl delete statefulsets etcd-clusterapi --force=true --grace-period 0 --wait=true")
}

func clusterLogs(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		b := runShell("kubectl get po -o name | grep clusterapi-controllers | xargs kubectl logs -c aws-cluster-controller -f")
		b.Wait()
		time.Sleep(5 * 1000 * 1000 * 1000)
	}
}

func restartControllers() {
	runShellWithWait("kubectl get po -o name | grep clusterapi-controllers | xargs kubectl delete")
}

func machineLogs(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		c := runShell("kubectl get po -o name | grep clusterapi-controllers | xargs kubectl logs -c aws-machine-controller -f")
		c.Wait()
		time.Sleep(5 * 1000 * 1000 * 1000)
	}
}
