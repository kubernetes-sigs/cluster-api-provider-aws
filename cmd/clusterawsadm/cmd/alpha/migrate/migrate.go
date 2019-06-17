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

package migrate

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awstags "github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/spf13/cobra"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
)

var (
	// See: https://docs.aws.amazon.com/sdk-for-go/api/service/resourcegroupstaggingapi/#TagResourcesInput
	maxARNs           = 20
	supportedVersions = []string{"0.3.0"}

	clusterName string
)

// MigrateCmd is the command for migrating AWS resources to be compatible
// with specific CAPA versions
func MigrateCmd() *cobra.Command { // nolint
	newCmd := &cobra.Command{
		Use:   "migrate [target version]",
		Short: "migrate between CAPA versions",
		Long: fmt.Sprintf(`Migrate AWS resources between incompatible versions of Cluster API Provider AWS.
Supported versions: %v`, supportedVersions),
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				fmt.Printf("Error: requires target version as an argument. Supported versions: %v\n\n", supportedVersions)
				if err := cmd.Help(); err != nil {
					return err
				}
				os.Exit(200)
			}
			if !isValidVersion(args[0]) {
				fmt.Printf("Error: unsupported migration target. Supported versions: %v\n\n", supportedVersions)
				cmd.Help()
				os.Exit(201)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			sess, err := session.NewSessionWithOptions(session.Options{
				SharedConfigState: session.SharedConfigEnable,
			})
			if err != nil {
				fmt.Printf("Error: %v", err)
				return nil
			}

			tagsSvc := awstags.New(sess)

			resources, err := getResourcesByCluster(tagsSvc, clusterName)
			if err != nil {
				return err
			}

			fmt.Printf("Found %v resources owned by cluster %q.\n", len(resources), clusterName)
			fmt.Printf("Applying new tags to cluster %q.\n", clusterName)

			err = applyNewTags(tagsSvc, resources, clusterName)
			if err != nil {
				return err
			}

			fmt.Printf("Removing legacy tags from cluster %q\n", clusterName)

			return removeOldTags(tagsSvc, resources, clusterName)
		},
	}

	newCmd.Flags().StringVarP(&clusterName, "clusterName", "n", "", "name of existing Cluster object")
	newCmd.MarkFlagRequired("clusterName")

	return newCmd
}

func getResourcesByCluster(svc *awstags.ResourceGroupsTaggingAPI, name string) ([]*string, error) {
	input := &awstags.GetResourcesInput{
		TagFilters: []*awstags.TagFilter{
			{
				Key:    aws.String(fmt.Sprintf("kubernetes.io/cluster/%s", name)),
				Values: []*string{aws.String("owned")},
			},
			{
				Key:    aws.String("sigs.k8s.io/cluster-api-provider-aws/managed"),
				Values: []*string{aws.String("true")},
			},
		},
	}

	out, err := svc.GetResources(input)
	if err != nil {
		return nil, err
	}

	arns := make([]*string, 0, len(out.ResourceTagMappingList))
	for _, resource := range out.ResourceTagMappingList {
		arns = append(arns, resource.ResourceARN)
	}

	return arns, nil
}

func applyNewTags(svc *awstags.ResourceGroupsTaggingAPI, arns []*string, name string) error {

	for i := 0; i <= (len(arns) / maxARNs); i++ {
		end := (i + 1) * maxARNs
		if end > len(arns) {
			end = len(arns)
		}

		input := &awstags.TagResourcesInput{
			ResourceARNList: arns[i*maxARNs : end],
			Tags: map[string]*string{
				v1alpha1.ClusterTagKey(name): aws.String("owned"),
			},
		}

		_, err := svc.TagResources(input)
		if err != nil {
			return err
		}
	}

	return nil
}

func removeOldTags(svc *awstags.ResourceGroupsTaggingAPI, arns []*string, name string) error {
	for i := 0; i <= (len(arns) / maxARNs); i++ {
		end := (i + 1) * maxARNs
		if end > len(arns) {
			end = len(arns)
		}

		managedInput := &awstags.UntagResourcesInput{
			ResourceARNList: arns[i*maxARNs : end],
			TagKeys: []*string{
				aws.String("sigs.k8s.io/cluster-api-provider-aws/managed"),
			},
		}

		_, err := svc.UntagResources(managedInput)
		if err != nil {
			return err
		}
	}

	var filteredARNs []*string

	for _, v := range arns {
		// instances should have both ownership tags, so filter those out
		// TODO(rudoi): is there a better way to filter?
		if !strings.Contains(aws.StringValue(v), "instance") {
			filteredARNs = append(filteredARNs, v)
		}
	}

	for i := 0; i <= (len(filteredARNs) / maxARNs); i++ {
		end := (i + 1) * maxARNs
		if i == (len(filteredARNs) / maxARNs) {
			end = len(filteredARNs)
		}

		ownedInput := &awstags.UntagResourcesInput{
			ResourceARNList: filteredARNs[i*maxARNs : end],
			TagKeys: []*string{
				aws.String(fmt.Sprintf("kubernetes.io/cluster/%s", name)),
			},
		}

		_, err := svc.UntagResources(ownedInput)
		if err != nil {
			return err
		}
	}

	return nil
}

func isValidVersion(s string) bool {
	for _, v := range supportedVersions {
		if s == v {
			return true
		}
	}
	return false
}
