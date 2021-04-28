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

package converters

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1alpha4"
)

// AddonSDKToAddonState is used to convert an AWS SDK Addon to a control plane AddonState
func AddonSDKToAddonState(eksAddon *eks.Addon) *ekscontrolplanev1.AddonState {
	addonState := &ekscontrolplanev1.AddonState{
		Name:                  aws.StringValue(eksAddon.AddonName),
		Version:               aws.StringValue(eksAddon.AddonVersion),
		ARN:                   aws.StringValue(eksAddon.AddonArn),
		CreatedAt:             metav1.NewTime(*eksAddon.CreatedAt),
		ModifiedAt:            metav1.NewTime(*eksAddon.ModifiedAt),
		Status:                eksAddon.Status,
		ServiceAccountRoleArn: eksAddon.ServiceAccountRoleArn,
		Issues:                []*ekscontrolplanev1.AddonIssue{},
	}
	if eksAddon.Health != nil {
		for _, issue := range eksAddon.Health.Issues {
			addonState.Issues = append(addonState.Issues, &ekscontrolplanev1.AddonIssue{
				Code:        issue.Code,
				Message:     issue.Message,
				ResourceIDs: issue.ResourceIds,
			})
		}
	}

	return addonState
}
