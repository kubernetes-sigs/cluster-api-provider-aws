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

package resource

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/config"
	rgapi "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	rgapitypes "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
)

// ListAWSResource fetches all AWS resources created by CAPA.
func ListAWSResource(region, clusterName string) (AWSResourceList, error) {
	var resourceList AWSResourceList
	ctx := context.TODO()

	sess, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return resourceList, err
	}

	resourceClient := rgapi.NewFromConfig(sess)
	input := &rgapi.GetResourcesInput{
		TagFilters: []rgapitypes.TagFilter{},
	}

	awsResourceTags := infrav1.Build(infrav1.BuildParams{
		ClusterName: clusterName,
		Lifecycle:   infrav1.ResourceLifecycleOwned,
	})

	for tagKey, tagValue := range awsResourceTags {
		tagFilter := rgapitypes.TagFilter{
			Key:    aws.String(tagKey),
			Values: []string{tagValue},
		}
		input.TagFilters = append(input.TagFilters, tagFilter)
	}

	output, err := resourceClient.GetResources(ctx, input)
	if err != nil {
		return resourceList, err
	}

	if len(output.ResourceTagMappingList) == 0 {
		fmt.Println("Could not find any AWS resource created by CAPA")
		return resourceList, nil
	}

	resourceList = AWSResourceList{
		ClusterName:  clusterName,
		AWSResources: []AWSResource{},
	}

	for _, eachResource := range output.ResourceTagMappingList {
		resourceARN, err := arn.Parse(*eachResource.ResourceARN)
		if err != nil {
			return resourceList, err
		}
		eachAWSResource := AWSResource{
			Partition: resourceARN.Partition,
			Service:   resourceARN.Service,
			Region:    resourceARN.Region,
			AccountID: resourceARN.AccountID,
			Resource:  resourceARN.Resource,
			ARN:       *eachResource.ResourceARN,
		}
		resourceList.AWSResources = append(resourceList.AWSResources, eachAWSResource)
	}

	return resourceList, nil
}
