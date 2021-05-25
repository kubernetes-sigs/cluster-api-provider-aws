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
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"

	cmdout "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/printers"
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
	cfg := aws.Config{}
	if *region != "" {
		cfg.Region = region
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            cfg,
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	eksClient := eks.New(sess)

	input := &eks.ListAddonsInput{
		ClusterName: clusterName,
	}
	output, err := eksClient.ListAddons(input)
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
			AddonName: addon,
		}
		describeOutput, err := eksClient.DescribeAddonVersions(describeInput)
		if err != nil {
			return fmt.Errorf("describing addon versions %s: %w", *addon, err)
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
				for _, architecture := range version.Architecture {
					newAddon.Architecture = append(newAddon.Architecture, *architecture)
				}
				for _, compat := range version.Compatibilities {
					compatibility := compatibility{
						ClusterVersion:   *compat.ClusterVersion,
						DefaultVersion:   *compat.DefaultVersion,
						PlatformVersions: []string{},
					}
					for _, platformVersion := range compat.PlatformVersions {
						compatibility.PlatformVersions = append(compatibility.PlatformVersions, *platformVersion)
					}
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
