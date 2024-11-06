/*
Copyright 2023 The Kubernetes Authors.

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

package gc

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"k8s.io/client-go/util/homedir"
	"k8s.io/kubectl/pkg/util/templates"

	gcproc "sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/gc"
)

func newConfigureCmd() *cobra.Command {
	var (
		clusterName       string
		namespace         string
		kubeConfig        string
		kubeConfigDefault string
		gcTasks           []string
	)

	if home := homedir.HomeDir(); home != "" {
		kubeConfigDefault = filepath.Join(home, ".kube", "config")
	}

	newCmd := &cobra.Command{
		Use:   "configure",
		Short: "Specify what cleanup tasks will be executed on a given cluster",
		Long: templates.LongDesc(`
			This command will set what cleanup tasks to execute on the given cluster
			during garbage collection (i.e. deleting) when the cluster is
			requested to be deleted. Supported values: load-balancer, security-group, target-group.
		`),
		Example: templates.Examples(`
			# Configure GC for a cluster to delete only load balancers and security groups using existing k8s context
			clusterawsadm gc configure --cluster-name=test-cluster --gc-task load-balancer --gc-task security-group

			# Reset GC configuration for a cluster using kubeconfig
			clusterawsadm gc configure --cluster-name=test-cluster --kubeconfig=test.kubeconfig
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			proc, err := gcproc.New(gcproc.GCInput{
				ClusterName:    clusterName,
				Namespace:      namespace,
				KubeconfigPath: kubeConfig,
			})
			if err != nil {
				return fmt.Errorf("creating command processor: %w", err)
			}

			if err := proc.Configure(cmd.Context(), gcTasks); err != nil {
				return fmt.Errorf("configuring garbage collection: %w", err)
			}
			fmt.Printf("Configuring garbage collection for cluster %s/%s\n", namespace, clusterName)

			return nil
		},
	}

	newCmd.Flags().StringVar(&clusterName, "cluster-name", "", "The name of the CAPA cluster")
	newCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "The namespace for the cluster definition")
	newCmd.Flags().StringVar(&kubeConfig, "kubeconfig", kubeConfigDefault, "Path to the kubeconfig file to use")
	newCmd.Flags().StringSliceVar(&gcTasks, "gc-task", []string{}, "Garbage collection tasks to execute during cluster deletion")

	newCmd.MarkFlagRequired("cluster-name") //nolint: errcheck

	return newCmd
}
