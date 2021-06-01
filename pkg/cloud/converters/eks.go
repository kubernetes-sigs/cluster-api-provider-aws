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
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1alpha3"
	infrav1exp "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1alpha3"
)

var (
	// ErrUnknowTaintEffect is an error when a unknown TaintEffect is used
	ErrUnknowTaintEffect = errors.New("uknown taint effect")
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

// TaintToSDK is used to a CAPA Taint to AWS SDK taint.
func TaintToSDK(taint infrav1exp.Taint) (*eks.Taint, error) {
	convertedEffect, err := TaintEffectToSDK(taint.Effect)
	if err != nil {
		return nil, fmt.Errorf("converting taint effect %s: %w", taint.Effect, err)
	}
	return &eks.Taint{
		Effect: aws.String(convertedEffect),
		Key:    aws.String(taint.Key),
		Value:  aws.String(taint.Value),
	}, nil

}

// TaintsToSDK is used to convert an array of CAPA Taints to AWS SDK taints.
func TaintsToSDK(taints infrav1exp.Taints) ([]*eks.Taint, error) {
	converted := []*eks.Taint{}

	for _, taint := range taints {
		convertedTaint, err := TaintToSDK(taint)
		if err != nil {
			return nil, fmt.Errorf("converting taint: %w", err)
		}
		converted = append(converted, convertedTaint)
	}

	return converted, nil
}

// TaintsFromSDK is used to convert an array of AWS SDK taints to CAPA Taints
func TaintsFromSDK(taints []*eks.Taint) (infrav1exp.Taints, error) {
	converted := infrav1exp.Taints{}
	for _, taint := range taints {
		convertedEffect, err := TaintEffectFromSDK(*taint.Effect)
		if err != nil {
			return nil, fmt.Errorf("converting taint effect %s: %w", *taint.Effect, err)
		}
		converted = append(converted, infrav1exp.Taint{
			Effect: convertedEffect,
			Key:    *taint.Key,
			Value:  *taint.Value,
		})
	}

	return converted, nil
}

// TaintEffectToSDK is used to convert a TaintEffect to the AWS SDK taint effect value
func TaintEffectToSDK(effect infrav1exp.TaintEffect) (string, error) {
	switch effect {
	case infrav1exp.TaintEffectNoExecute:
		return eks.TaintEffectNoExecute, nil
	case infrav1exp.TaintEffectPreferNoSchedule:
		return eks.TaintEffectPreferNoSchedule, nil
	case infrav1exp.TaintEffectNoSchedule:
		return eks.TaintEffectNoSchedule, nil
	default:
		return "", ErrUnknowTaintEffect
	}
}

// TaintEffectFromSDK is used to convert a AWS SDK taint effect value to a TaintEffect
func TaintEffectFromSDK(effect string) (infrav1exp.TaintEffect, error) {
	switch effect {
	case eks.TaintEffectNoExecute:
		return infrav1exp.TaintEffectNoExecute, nil
	case eks.TaintEffectPreferNoSchedule:
		return infrav1exp.TaintEffectPreferNoSchedule, nil
	case eks.TaintEffectNoSchedule:
		return infrav1exp.TaintEffectNoSchedule, nil
	default:
		return "", ErrUnknowTaintEffect
	}
}
