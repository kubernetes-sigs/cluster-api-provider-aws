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

package identityprovider

import (
	"context"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/eks"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/planner"
)

// NewMultiProviderPlan creates a plan that reconciles multiple EKS OIDC identity providers.
func NewMultiProviderPlan(
	clusterName string,
	currentByName map[string]*OidcIdentityProviderConfig,
	desiredByName map[string]*OidcIdentityProviderConfig,
	client eks.Client,
	log logger.Wrapper,
) planner.Plan {
	return &multiProviderPlan{
		clusterName:   clusterName,
		currentByName: currentByName,
		desiredByName: desiredByName,
		eksClient:     client,
		log:           log,
	}
}

type multiProviderPlan struct {
	clusterName   string
	currentByName map[string]*OidcIdentityProviderConfig
	desiredByName map[string]*OidcIdentityProviderConfig
	eksClient     eks.Client
	log           logger.Wrapper
}

func (p *multiProviderPlan) Create(ctx context.Context) ([]planner.Procedure, error) {
	procedures := []planner.Procedure{}
	names := unionProviderNames(p.currentByName, p.desiredByName)

	for _, name := range names {
		current := p.currentByName[name]
		desired := p.desiredByName[name]

		providerPlan := NewPlan(p.clusterName, current, desired, p.eksClient, p.log)
		providerProcedures, err := providerPlan.Create(ctx)
		if err != nil {
			return nil, err
		}
		procedures = append(procedures, providerProcedures...)
	}

	return procedures, nil
}

func unionProviderNames(currentByName, desiredByName map[string]*OidcIdentityProviderConfig) []string {
	seen := map[string]struct{}{}
	names := []string{}

	for name := range currentByName {
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		names = append(names, name)
	}

	for name := range desiredByName {
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		names = append(names, name)
	}

	return names
}
