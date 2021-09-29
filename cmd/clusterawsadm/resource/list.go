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
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/session"
	tagapi "github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha4"
)

// ListAWSResource fetches all AWS resources created by CAPA.
func ListAWSResource(region, clusterName *string) (AWSResourceList, error) {
	var resourceList AWSResourceList
	cfg := aws.Config{}
	if *region != "" {
		cfg.Region = region
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            cfg,
	})
	if err != nil {
		return resourceList, err
	}

	resourceClient := tagapi.New(sess)
	input := &tagapi.GetResourcesInput{
		TagFilters: []*tagapi.TagFilter{},
	}

	awsResourceTags := infrav1.Build(infrav1.BuildParams{
		ClusterName: *clusterName,
		Lifecycle:   infrav1.ResourceLifecycleOwned,
	})

	for tagKey, tagValue := range awsResourceTags {
		tagFilter := &tagapi.TagFilter{}
		tagFilter.SetKey(tagKey)
		tagFilter.SetValues([]*string{aws.String(tagValue)})
		input.TagFilters = append(input.TagFilters, tagFilter)
	}

	output, err := resourceClient.GetResources(input)
	if err != nil {
		return resourceList, err
	}

	if len(output.ResourceTagMappingList) == 0 {
		fmt.Println("Could not find any AWS resource created by CAPA")
		return resourceList, nil
	}

	resourceList = AWSResourceList{
		ClusterName:  *clusterName,
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
