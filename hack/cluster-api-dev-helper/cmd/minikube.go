// Copyright © 2018 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func defineMinikubeCmd(parent *cobra.Command) {
	newCmd := &cobra.Command{
		Use:   "minikube",
		Short: "minikube helper commands",
		Long:  `minikube helper commands`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.Help())
		},
	}

	defineMinikubeFixLibVirtNetworkCmd(newCmd)
	defineMinikubeClean(newCmd)
	defineMinikubeKVMDestroy(newCmd)
	defineMinikubeSetup(newCmd)
	defineMinikubeReset(newCmd)

	parent.AddCommand(newCmd)
}

func defineMinikubeFixLibVirtNetworkCmd(parent *cobra.Command) {
	newCmd := &cobra.Command{
		Use:   "fix-libvirt-network",
		Short: "Make the KVM minikube network permanent using LibVirt",
		Long:  `Make the KVM minikube network permanent using LibVirt. Requires virsh`,
		Run: func(cmd *cobra.Command, args []string) {
			runCommandWithWait("virsh", "-c", "qemu:///system", "net-autostart", "minikube-net")
		},
	}
	parent.AddCommand(newCmd)
}

func defineMinikubeClean(parent *cobra.Command) {
	newCmd := &cobra.Command{
		Use:   "clean",
		Short: "Clean up .minikube directory",
		Long:  `Clean up .minikube directory, deleting machines`,
		Run: func(cmd *cobra.Command, args []string) {
			minikubeClean()
		},
	}
	parent.AddCommand(newCmd)
}

func defineMinikubeKVMDestroy(parent *cobra.Command) {
	newCmd := &cobra.Command{
		Use:   "kvm-destroy",
		Short: "Destroy the KVM VM backing minikube",
		Long:  `Destroy the KVM VM backing minikube`,
		Run: func(cmd *cobra.Command, args []string) {
			runCommandWithWait("virsh", "-c", "qemu:///system", "destroy", "minikube")
			runCommandWithWait("virsh", "-c", "qemu:///system", "undefine", "--remove-all-storage", "--delete-snapshots", "minikube")
			runCommandWithWait("virsh", "-c", "qemu:///system", "net-undefine", "minikube-net")
			minikubeClean()
		},
	}
	parent.AddCommand(newCmd)
}

func defineMinikubeSetup(parent *cobra.Command) {
	newCmd := &cobra.Command{
		Use:   "setup",
		Short: "Set up minikube for use with Cluster API",
		Long:  `Set up minikube for use with Cluster API`,
		Run: func(cmd *cobra.Command, args []string) {
			runCommandWithWait("minikube", "config", "set", "bootstrapper", "kubeadm")
			runCommandWithWait("minikube", "config", "set", "kubernetes-version", "v1.9.4")
		},
	}
	parent.AddCommand(newCmd)
}

func defineMinikubeReset(parent *cobra.Command) {
	newCmd := &cobra.Command{
		Use:   "reset",
		Short: "Remove Cluster API specific minikube configuration",
		Long:  `Remove Cluster API specific minikube configuration`,
		Run: func(cmd *cobra.Command, args []string) {
			runCommandWithWait("minikube", "config", "unset", "kubernetes-version")
			runCommandWithWait("minikube", "config", "unset", "bootstrapper")
		},
	}
	parent.AddCommand(newCmd)
}

func minikubeClean() {
	runCommandWithWait("sh", "-c", "rm -rf $HOME/.minikube/machines")
	runCommandWithWait("sh", "-c", "rm -rf $HOME/.minikube/profiles")
	runCommandWithWait("sh", "-c", "rm -rf $HOME/.minikube/certs")
	runCommandWithWait("sh", "-c", "rm -rf $HOME/.minikube/*.pem")
}
