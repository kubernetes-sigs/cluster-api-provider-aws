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
			a := runShell(`bazel-bin/cmd/clusterctl/linux_amd64_stripped/clusterctl -v4  \
			  create cluster \
				-m ./cmd/clusterctl/examples/aws/out/machines.yaml \
				-c ./cmd/clusterctl/examples/aws/out/cluster.yaml \
				-p ./cmd/clusterctl/examples/aws/out/provider-components.yaml \
				-a ./cmd/clusterctl/examples/aws/out/addons.yaml \
				--provider aws \
				--existing-bootstrap-cluster-kubeconfig ${HOME}/.kube/config`)
			wg.Add(1)
			go controllerLogs(&wg)
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
			wg.Add(1)
			go controllerLogs(&wg)
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
			go controllerLogs(&wg)
			runShell("kubectl delete clusters -n aws-provider-system --all --wait=true")
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
			go controllerLogs(&wg)
			runShell("kubectl delete machines -n aws-provider-system --all --wait=true")
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
			runShellWithWait("kubectl apply -f ./out/provider-components.yaml")
		},
	}
	parent.AddCommand(newCmd)
}

func destroyMachines() {
	runShellWithWait("kubectl get machine -o name -n aws-provider-system | xargs kubectl patch  -n aws-provider-system -p '{\"metadata\":{\"finalizers\":null}}' --type=merge")
	runShellWithWait("kubectl delete machines -n aws-provider-system --force=true --grace-period 0 --all --wait=true")
}

func destroyClusters() {
	runShellWithWait("kubectl patch cluster test1  -n aws-provider-system -p '{\"metadata\":{\"finalizers\":null}}' --type=merge")
	runShellWithWait("kubectl delete clusters  -n aws-provider-system --force=true --grace-period 0 --all --wait=true")
}

func destroyControlPlane() {
	runShellWithWait("kubectl delete statefulset aws-provider-controller-manager  -n aws-provider-system --force=true --grace-period 0 --wait=true")
	runShellWithWait("kubectl delete crds  -n aws-provider-system --all --force=true --grace-period 0 --wait=true")
}

func controllerLogs(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		b := runShell("kubectl get po -o name -n aws-provider-system | grep aws-provider-controller-manager | xargs kubectl logs -n aws-provider-system -c manager -f")
		b.Wait()
		time.Sleep(5 * 1000 * 1000 * 1000)
	}
}

func restartControllers() {
	runShellWithWait("kubectl get po -o name -n aws-provider-system | grep aws-provider-controller-manager | xargs kubectl -n aws-provider-system delete")
}
