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

package ssm

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/ssm"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
)

// Service holds a collection of interfaces.
// The interfaces are broken down like this to group functions together.
// One alternative is to have a large list of functions from the ec2 client.
type Service struct {
	scope     cloud.ClusterScoper
	SSMClient SSMAPI
}

// SSMAPI defines the interface for interacting with AWS SSM Parameter Store.
type SSMAPI interface {
	PutParameter(ctx context.Context, input *ssm.PutParameterInput) (*ssm.PutParameterOutput, error)
	DeleteParameter(ctx context.Context, input *ssm.DeleteParameterInput) (*ssm.DeleteParameterOutput, error)
	GetParameter(ctx context.Context, input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error)
	// Add more methods as needed
}

// SSMClientV2 is a concrete implementation of the SSMAPI interface using AWS SDK v2.
type SSMClientV2 struct {
	Client *ssm.Client
}

// PutParameter adds or overwrites a parameter in AWS SSM Parameter Store.
func (c *SSMClientV2) PutParameter(ctx context.Context, input *ssm.PutParameterInput) (*ssm.PutParameterOutput, error) {
	if c.Client == nil {
		return nil, errors.New("SSM client is not initialized")
	}
	return c.Client.PutParameter(ctx, input)
}

// DeleteParameter deletes a parameter from AWS SSM Parameter Store.
func (c *SSMClientV2) DeleteParameter(ctx context.Context, input *ssm.DeleteParameterInput) (*ssm.DeleteParameterOutput, error) {
	if c.Client == nil {
		return nil, errors.New("SSM client is not initialized")
	}
	return c.Client.DeleteParameter(ctx, input)
}

// GetParameter retrieves a parameter from AWS SSM Parameter Store.
func (c *SSMClientV2) GetParameter(ctx context.Context, input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	if c.Client == nil {
		return nil, errors.New("SSM client is not initialized")
	}
	return c.Client.GetParameter(ctx, input)
}

// Ensure SSMClientV2 satisfies the SSMAPI interface.
var _ SSMAPI = &SSMClientV2{}

// NewService creates a new Service for managing secrets in AWS SSM.
func NewService(secretsScope cloud.ClusterScoper) *Service {
	return &Service{
		scope: secretsScope,
		SSMClient: &SSMClientV2{
			Client: scope.NewSSMClient(secretsScope, secretsScope, secretsScope, secretsScope.InfraCluster()),
		},
	}
}
