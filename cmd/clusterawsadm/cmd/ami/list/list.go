/*
Copyright 2021 The Kubernetes Authors.

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

// Package list provides a way to list AMIs from the default AWS account where AMIs are stored.
package list

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/ami"
	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/cmd/flags"
	cmdout "sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/printers"
)

var (
	kubernetesVersion string
	opSystem          string
	outputPrinter     string
	ownerID           string
)

// ListAMICmd is a CLI command that will list AMIs from the default AWS account where AMIs are stored.
func ListAMICmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "list",
		Short: "List AMIs from the default AWS account where AMIs are stored",
		Long: templates.LongDesc(`
			List AMIs based on Kubernetes version, OS, region. If no arguments are provided,
			it will print all AMIs in all regions, OS types for the supported Kubernetes versions.
            Supported Kubernetes versions start from the latest stable version and goes 2 release back:
			if the latest stable release is v1.20.4- v1.19.x and v1.18.x are supported.
			Note: First release of each version will be skipped, e.g., v1.21.0
			To list AMIs of unsupported Kubernetes versions, --kubernetes-version flag needs to be provided.
		`),
		Example: templates.Examples(`
		# List AMIs from the default AWS account where AMIs are stored.
		# Available os options: centos-7, ubuntu-24.04, ubuntu-22.04, amazon-2, flatcar-stable
		clusterawsadm ami list --kubernetes-version=v1.18.12 --os=ubuntu-20.04  --region=us-west-2
		# To list all supported AMIs in all supported Kubernetes versions, regions, and linux distributions:
		clusterawsadm ami list
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			region, _ := flags.GetRegion(cmd)

			printer, err := cmdout.New(outputPrinter, os.Stdout)
			if err != nil {
				return fmt.Errorf("failed creating output printer: %w", err)
			}

			listByVersion, err := ami.List(ami.ListInput{
				Region:            region,
				KubernetesVersion: kubernetesVersion,
				OperatingSystem:   opSystem,
				OwnerID:           ownerID,
			})
			if err != nil {
				return err
			}
			if len(listByVersion.Items) == 0 {
				fmt.Println("No AMIs found")
				return nil
			}

			if outputPrinter == string(cmdout.PrinterTypeTable) {
				table := listByVersion.ToTable()
				printer.Print(table)
			} else {
				printer.Print(listByVersion)
			}

			return nil
		},
	}

	flags.AddRegionFlag(newCmd)
	addOsFlag(newCmd)
	addKubernetesVersionFlag(newCmd)
	addOutputFlag(newCmd)
	addOwnerIDFlag(newCmd)
	return newCmd
}

func addOsFlag(c *cobra.Command) {
	c.Flags().StringVar(&opSystem, "os", "", "Operating system of the AMI to be listed")
}

func addKubernetesVersionFlag(c *cobra.Command) {
	c.Flags().StringVar(&kubernetesVersion, "kubernetes-version", "", "Kubernetes version of the AMI to be copied")
}

func addOutputFlag(c *cobra.Command) {
	c.Flags().StringVarP(&outputPrinter, "output", "o", "table", "The output format of the results. Possible values: table,json,yaml")
}

func addOwnerIDFlag(c *cobra.Command) {
	c.Flags().StringVarP(&ownerID, "owner-id", "", "", "The owner ID of the AWS account to be used for listing AMIs")
}
