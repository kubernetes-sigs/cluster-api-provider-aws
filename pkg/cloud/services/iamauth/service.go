/*
Copyright 2020 The Kubernetes Authors.

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

package iamauth

import (
	"github.com/aws/aws-sdk-go/service/sts/stsiface"

	"sigs.k8s.io/controller-runtime/pkg/client"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
)

// Scope is a scope for use with the iamauth reconciling service
type Scope interface {
	cloud.ClusterScoper

	// IAMAuthConfig returns the IAM authenticator config
	IAMAuthConfig() *ekscontrolplanev1.IAMAuthenticatorConfig
}

type Service struct {
	scope     Scope
	backend   BackendType
	client    client.Client
	STSClient stsiface.STSAPI
}

func NewService(iamScope Scope, backend BackendType, client client.Client) *Service {
	return &Service{
		scope:     iamScope,
		backend:   backend,
		client:    client,
		STSClient: scope.NewSTSClient(iamScope, iamScope, iamScope, iamScope.InfraCluster()),
	}
}
