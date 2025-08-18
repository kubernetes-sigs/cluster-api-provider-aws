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

func listInstalledCmd() *cobra.Command {
	clusterName := ""
	region := ""
	outputPrinter := ""

	newCmd := &cobra.Command{
		Use:   "list-installed",
		Short: "List installed EKS addons",
		Long:  "Lists the addons that are installed for an EKS cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listInstalledAddons(&region, &clusterName, &outputPrinter)
		},
	}

	newCmd.Flags().StringVarP(&region, "region", "r", "", "The AWS region containing the EKS cluster")
	newCmd.Flags().StringVarP(&clusterName, "cluster-name", "n", "", "The name of the cluster to get the list of installed addons for")
	newCmd.Flags().StringVarP(&outputPrinter, "output", "o", "table", "The output format of the results. Possible values: table,json,yaml")
	newCmd.MarkFlagRequired("cluster-name") //nolint: errcheck

	return newCmd
}

func listInstalledAddons(region, clusterName, printerType *string) error {
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

	addonsList := installedAddonsList{
		Cluster: *clusterName,
		Addons:  []installedAddon{},
	}
	for _, addon := range output.Addons {
		describeInput := &eks.DescribeAddonInput{
			AddonName:   aws.String(addon),
			ClusterName: clusterName,
		}
		describeOutput, err := eksClient.DescribeAddon(ctx, describeInput)
		if err != nil {
			return fmt.Errorf("describing addon %s: %w", addon, err)
		}

		if describeOutput.Addon == nil {
			continue
		}

		installedAddon := installedAddon{
			Name:         *describeOutput.Addon.AddonName,
			Version:      *describeOutput.Addon.AddonVersion,
			AddonARN:     *describeOutput.Addon.AddonArn,
			RoleARN:      describeOutput.Addon.ServiceAccountRoleArn,
			Status:       string(describeOutput.Addon.Status),
			CreatedAt:    *describeOutput.Addon.CreatedAt,
			ModifiedAt:   *describeOutput.Addon.ModifiedAt,
			Tags:         describeOutput.Addon.Tags,
			HealthIssues: []issue{},
		}
		for _, addonIssue := range describeOutput.Addon.Health.Issues {
			newIssue := issue{
				Code:        string(addonIssue.Code),
				Message:     *addonIssue.Message,
				ResourceIDs: []string{},
			}
			newIssue.ResourceIDs = append(newIssue.ResourceIDs, addonIssue.ResourceIds...)
			installedAddon.HealthIssues = append(installedAddon.HealthIssues, newIssue)
		}

		addonsList.Addons = append(addonsList.Addons, installedAddon)
	}

	outputPrinter, err := cmdout.New(*printerType, os.Stderr)
	if err != nil {
		return fmt.Errorf("failed creating output printer: %w", err)
	}

	if *printerType == string(cmdout.PrinterTypeTable) {
		table := addonsList.ToTable()
		fmt.Printf("Installed addons for cluster %s:", *clusterName)
		outputPrinter.Print(table)
	} else {
		outputPrinter.Print(addonsList)
	}

	return nil
}
