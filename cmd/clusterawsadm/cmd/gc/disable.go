/*
Copyright 2022 The Kubernetes Authors.

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

func newDisableCmd() *cobra.Command {
	var (
		clusterName       string
		namespace         string
		kubeConfig        string
		kubeConfigDefault string
	)

	if home := homedir.HomeDir(); home != "" {
		kubeConfigDefault = filepath.Join(home, ".kube", "config")
	}

	newCmd := &cobra.Command{
		Use:   "disable",
		Short: "Mark a cluster as NOT requiring external resource garbage collection",
		Long: templates.LongDesc(`
			This command will mark the given cluster as not requiring external
			resource garbage collection (i.e. deleting) when the cluster is
			requested to be deleted.
		`),
		Example: templates.Examples(`
			# Disable GC for a cluster using existing k8s context
			clusterawsadm gc disable --cluster-name=test-cluster

			# Disable GC for a cluster using kubeconfig
			clusterawsadm gc disable --cluster-name=test-cluster --kubeconfig=test.kubeconfig
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

			if err := proc.Disable(cmd.Context()); err != nil {
				return fmt.Errorf("disabling garbage collection: %w", err)
			}
			fmt.Printf("Disabled garbage collection for cluster %s/%s\n", namespace, clusterName)

			return nil
		},
	}

	newCmd.Flags().StringVar(&clusterName, "cluster-name", "", "The name of the CAPA cluster")
	newCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "The namespace for the cluster definition")
	newCmd.Flags().StringVar(&kubeConfig, "kubeconfig", kubeConfigDefault, "Path to the kubeconfig file to use")

	newCmd.MarkFlagRequired("cluster-name") //nolint: errcheck

	return newCmd
}
