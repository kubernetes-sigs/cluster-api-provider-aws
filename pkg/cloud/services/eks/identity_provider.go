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

package eks

import (
	"context"
	"fmt"
	"reflect"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/pkg/errors"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/converters"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/eks/identityprovider"
)

func (s *Service) reconcileIdentityProvider(ctx context.Context) error {
	s.scope.Info("reconciling oidc identity provider")

	clusterName := s.scope.KubernetesClusterName()
	currentByName, err := s.getAssociatedIdentityProviders(ctx, clusterName)
	if err != nil {
		return errors.Wrap(err, "unable to list associated identity providers")
	}

	desiredConfigs := s.scope.EffectiveOIDCIdentityProviderConfigs()
	if len(desiredConfigs) == 0 && len(currentByName) == 0 {
		s.scope.Info("no oidc provider config, skipping reconcile")
		return nil
	}

	desiredByName := make(map[string]*identityprovider.OidcIdentityProviderConfig, len(desiredConfigs))
	for i := range desiredConfigs {
		cfg := desiredConfigs[i]
		desiredByName[cfg.IdentityProviderConfigName] = converters.ConvertSDKToIdentityProvider(&cfg)
	}

	if len(desiredByName) == 0 && len(currentByName) == 0 {
		s.scope.Info("no identity provider required or installed, no action needed")
		return nil
	}

	s.scope.Debug("creating oidc provider plan", "desired", desiredByName, "current", currentByName)
	procedures, err := identityprovider.
		NewMultiProviderPlan(clusterName, currentByName, desiredByName, s.EKSClient, s.scope).
		Create(ctx)
	if err != nil {
		s.scope.Error(err, "failed creating eks identity provider plan")
		return fmt.Errorf("creating eks identity provider plan: %w", err)
	}

	if len(procedures) == 0 {
		return s.patchIdentityProviderStatus(ctx, clusterName, desiredByName)
	}

	s.scope.Debug("computed EKS identity provider plan", "numprocs", len(procedures))

	for _, procedure := range procedures {
		s.scope.Info("Executing identity provider procedure", "name", procedure.Name())
		if err := procedure.Do(ctx); err != nil {
			s.scope.Error(err, "failed executing identity provider procedure", "name", procedure.Name())
			return fmt.Errorf("%s: %w", procedure.Name(), err)
		}
	}

	return s.patchIdentityProviderStatus(ctx, clusterName, desiredByName)
}

func (s *Service) patchIdentityProviderStatus(
	ctx context.Context,
	clusterName string,
	desiredByName map[string]*identityprovider.OidcIdentityProviderConfig,
) error {
	latestByName, err := s.getAssociatedIdentityProviders(ctx, clusterName)
	if err != nil {
		return errors.Wrap(err, "getting associated identity providers")
	}

	statuses := make([]ekscontrolplanev1.IdentityProviderStatus, 0, len(desiredByName))
	for name := range desiredByName {
		latest := latestByName[name]
		if latest == nil {
			continue
		}
		statuses = append(statuses, ekscontrolplanev1.IdentityProviderStatus{
			IdentityProviderConfigName: name,
			ARN:                        latest.IdentityProviderConfigArn,
			Status:                     latest.Status,
		})
	}
	sort.Slice(statuses, func(i, j int) bool {
		return statuses[i].IdentityProviderConfigName < statuses[j].IdentityProviderConfigName
	})

	newLegacy := legacyIdentityProviderStatus(statuses)
	if identityProviderStatusesEqual(statuses, s.scope.ControlPlane.Status.IdentityProviderStatuses) &&
		reflect.DeepEqual(newLegacy, s.scope.ControlPlane.Status.IdentityProviderStatus) {
		return nil
	}

	s.scope.ControlPlane.Status.IdentityProviderStatuses = statuses
	s.scope.ControlPlane.Status.IdentityProviderStatus = newLegacy

	if err := s.scope.PatchObject(); err != nil {
		return errors.Wrap(err, "updating identity provider status")
	}

	return nil
}

func legacyIdentityProviderStatus(statuses []ekscontrolplanev1.IdentityProviderStatus) ekscontrolplanev1.IdentityProviderStatus {
	if len(statuses) == 1 {
		return statuses[0]
	}
	return ekscontrolplanev1.IdentityProviderStatus{}
}

func identityProviderStatusesEqual(a, b []ekscontrolplanev1.IdentityProviderStatus) bool {
	if len(a) != len(b) {
		return false
	}
	return reflect.DeepEqual(a, b)
}

func (s *Service) getAssociatedIdentityProviders(ctx context.Context, clusterName string) (map[string]*identityprovider.OidcIdentityProviderConfig, error) {
	list, err := s.EKSClient.ListIdentityProviderConfigs(ctx, &eks.ListIdentityProviderConfigsInput{
		ClusterName: aws.String(clusterName),
	})
	if err != nil {
		return nil, errors.Wrap(err, "listing identity provider configs")
	}

	providers := make(map[string]*identityprovider.OidcIdentityProviderConfig, len(list.IdentityProviderConfigs))
	for i := range list.IdentityProviderConfigs {
		idpRef := list.IdentityProviderConfigs[i]
		providerConfig, err := s.EKSClient.DescribeIdentityProviderConfig(ctx, &eks.DescribeIdentityProviderConfigInput{
			ClusterName:            aws.String(clusterName),
			IdentityProviderConfig: &idpRef,
		})
		if err != nil {
			return nil, errors.Wrap(err, "describing identity provider config")
		}

		config := providerConfig.IdentityProviderConfig.Oidc
		name := aws.ToString(config.IdentityProviderConfigName)
		providers[name] = &identityprovider.OidcIdentityProviderConfig{
			ClientID:                   aws.ToString(config.ClientId),
			GroupsClaim:                aws.ToString(config.GroupsClaim),
			GroupsPrefix:               aws.ToString(config.GroupsPrefix),
			IdentityProviderConfigArn:  aws.ToString(config.IdentityProviderConfigArn),
			IdentityProviderConfigName: name,
			IssuerURL:                  aws.ToString(config.IssuerUrl),
			RequiredClaims:             config.RequiredClaims,
			Status:                     string(config.Status),
			Tags:                       converters.MapPtrToMap(config.Tags),
			UsernameClaim:              aws.ToString(config.UsernameClaim),
			UsernamePrefix:             aws.ToString(config.UsernamePrefix),
		}
	}

	return providers, nil
}
