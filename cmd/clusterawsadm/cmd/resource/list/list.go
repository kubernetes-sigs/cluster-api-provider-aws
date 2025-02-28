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

// Package list provides the list command for the resource package.
package list

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/cmd/flags"
	cmdout "sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/printers"
	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/resource"
)

// ListAWSResourceCmd is the root cmd to list AWS resources created by CAPA.
func ListAWSResourceCmd() *cobra.Command {
	outputPrinterType := ""
	clusterName := ""
	region := ""
	newCmd := &cobra.Command{
		Use:   "list",
		Short: "List all AWS resources created by CAPA",
		Long: templates.LongDesc(`
			List AWS resources directly created by CAPA based on region and cluster-name. There are some indirect resources like Cloudwatch alarms, rules, etc
			which are not directly created by CAPA, so those resources are not listed here.
			If region and cluster-name are not set, then it will throw an error.
		`),
		Example: templates.Examples(`
		# List AWS resources directly created by CAPA in given region and clustername
		clusterawsadm resource list --region=us-east-1 --cluster-name=test-cluster
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			region, err = flags.GetRegion(cmd)
			if err != nil {
				return err
			}

			fmt.Fprintf(os.Stdout, "Attempting to fetch resources created by CAPA for cluster:%s present in %s\n\n", clusterName, region)
			resourceList, err := resource.ListAWSResource(&region, &clusterName)
			if err != nil || len(resourceList.AWSResources) == 0 {
				return err
			}

			outputPrinter, err := cmdout.New(outputPrinterType, os.Stdout)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed creating output printer: %s\n", err.Error())
				return err
			}

			fmt.Fprintln(os.Stdout, "Following resources are found: ")
			if outputPrinterType == string(cmdout.PrinterTypeTable) {
				table := resourceList.ToTable()
				outputPrinter.Print(table)
			} else {
				outputPrinter.Print(resourceList)
			}
			return nil
		},
	}

	newCmd.Flags().StringVarP(&region, "region", "r", "", "The AWS region where resources are created by CAPA")
	newCmd.Flags().StringVarP(&clusterName, "cluster-name", "n", "", "The name of the cluster where AWS resources created by CAPA")
	newCmd.Flags().StringVarP(&outputPrinterType, "output", "o", "table", "The output format of the results. Possible values: table, json, yaml")
	newCmd.MarkFlagRequired("cluster-name") //nolint: errcheck
	return newCmd
}
