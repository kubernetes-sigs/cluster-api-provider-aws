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

package addons

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/spf13/cobra"
	"k8s.io/utils/ptr"

	cmdout "sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/printers"
)

func listAvailableCmd() *cobra.Command {
	clusterName := ""
	region := ""
	outputPrinter := ""

	newCmd := &cobra.Command{
		Use:   "list-available",
		Short: "List available EKS addons",
		Long:  "Lists the addons that are available for use with an EKS cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listAvailableAddons(&region, &clusterName, &outputPrinter)
		},
	}

	newCmd.Flags().StringVarP(&region, "region", "r", "", "The AWS region containing the EKS cluster")
	newCmd.Flags().StringVarP(&clusterName, "cluster-name", "n", "", "The name of the cluster to get the list of available addons for")
	newCmd.Flags().StringVarP(&outputPrinter, "output", "o", "table", "The output format of the results. Possible values: table,json,yaml")
	newCmd.MarkFlagRequired("cluster-name") //nolint: errcheck

	return newCmd
}

func listAvailableAddons(region, clusterName, printerType *string) error {
	ctx := context.TODO()

	optFns := []func(*config.LoadOptions) error{
		config.WithRegion(ptr.Deref(region, "")),
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), optFns...)

	if err != nil {
		return err
	}

	eksClient := eks.NewFromConfig(cfg)

	input := &eks.ListAddonsInput{
		ClusterName: clusterName,
	}
	output, err := eksClient.ListAddons(ctx, input)
	if err != nil {
		return fmt.Errorf("list addons: %w", err)
	}

	if len(output.Addons) == 0 {
		fmt.Println("No EKS addons found")
		return nil
	}

	addonsList := availableAddonsList{
		Cluster: *clusterName,
		Addons:  []availableAddon{},
	}
	for _, addon := range output.Addons {
		describeInput := &eks.DescribeAddonVersionsInput{
			AddonName: aws.String(addon),
		}
		describeOutput, err := eksClient.DescribeAddonVersions(ctx, describeInput)
		if err != nil {
			return fmt.Errorf("describing addon versions %s: %w", addon, err)
		}

		for _, info := range describeOutput.Addons {
			for _, version := range info.AddonVersions {
				newAddon := availableAddon{
					Name:            *info.AddonName,
					Type:            *info.Type,
					Version:         *version.AddonVersion,
					Architecture:    []string{},
					Compatibilities: []compatibility{},
				}
				newAddon.Architecture = append(newAddon.Architecture, version.Architecture...)

				for _, compat := range version.Compatibilities {
					compatibility := compatibility{
						ClusterVersion:   *compat.ClusterVersion,
						DefaultVersion:   compat.DefaultVersion,
						PlatformVersions: []string{},
					}
					compatibility.PlatformVersions = append(compatibility.PlatformVersions, compat.PlatformVersions...)
					newAddon.Compatibilities = append(newAddon.Compatibilities, compatibility)
				}
				addonsList.Addons = append(addonsList.Addons, newAddon)
			}
		}
	}

	outputPrinter, err := cmdout.New(*printerType, os.Stderr)
	if err != nil {
		return fmt.Errorf("failed creating output printer: %w", err)
	}

	if *printerType == string(cmdout.PrinterTypeTable) {
		table := addonsList.ToTable()
		fmt.Printf("Available addons for cluster %s:", *clusterName)
		outputPrinter.Print(table)
	} else {
		outputPrinter.Print(addonsList)
	}

	return nil
}
