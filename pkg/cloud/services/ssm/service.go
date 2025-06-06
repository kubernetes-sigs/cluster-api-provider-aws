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
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
)

// SSMClient defines the interface for the SSM client methods used by the service.
// It has been updated to include all methods required by the secret management logic.
type SSMClient interface {
	GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
	PutParameter(ctx context.Context, params *ssm.PutParameterInput, optFns ...func(*ssm.Options)) (*ssm.PutParameterOutput, error)
	DeleteParameter(ctx context.Context, params *ssm.DeleteParameterInput, optFns ...func(*ssm.Options)) (*ssm.DeleteParameterOutput, error)
}

// Service holds a collection of interfaces.
type Service struct {
	scope     cloud.ClusterScoper
	SSMClient SSMClient
}

// NewService returns a new service for the AWS SSM service.
func NewService(secretsScope cloud.ClusterScoper) *Service {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(secretsScope.Region()))
	if err != nil {
		panic(fmt.Sprintf("failed to load AWS config for SSM service: %v", err))
	}

	return &Service{
		scope:     secretsScope,
		SSMClient: ssm.NewFromConfig(cfg),
	}
}
