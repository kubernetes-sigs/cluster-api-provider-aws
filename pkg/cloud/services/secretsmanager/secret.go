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

package secretsmanager

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	apirand "k8s.io/apimachinery/pkg/util/rand"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/converters"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
)

const (
	parameterPrefix  = "aws.cluster.x-k8s.io"
	namespacesPrefix = "namespaces"
	clustersPrefix   = "clusters"
	instancesPrefix  = "instances"
)

// Create stores a secret in AWS Secrets Manager for a given machine
func (s *Service) Create(m *scope.MachineScope, data []byte) (string, error) {
	// Make sure to use the MachineScope here to get the merger of AWSCluster and AWSMachine tags
	additionalTags := m.AdditionalTags()
	// Set the cloud provider tag
	additionalTags[infrav1.ClusterAWSCloudProviderTagKey(s.scope.Name())] = string(infrav1.ResourceLifecycleOwned)

	tags := infrav1.Build(infrav1.BuildParams{
		ClusterName: s.scope.Name(),
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Name:        aws.String(m.Name()),
		Role:        aws.String(m.Role()),
		Additional:  additionalTags,
	})

	name := s.secretName(m)

	resp, err := s.scope.SecretsManager.CreateSecret(&secretsmanager.CreateSecretInput{
		Name:         aws.String(name),
		SecretBinary: data,
		Tags:         converters.MapToSecretsManagerTags(tags),
	})

	if err != nil {
		return "", err
	}

	return aws.StringValue(resp.ARN), nil
}

// Delete the secret belonging to a machine from AWS Secrets Manager
func (s *Service) Delete(m *scope.MachineScope) error {
	secretArn := m.AWSMachine.Spec.CloudInit.SecretARN
	if secretArn == "" {
		return nil
	}
	_, err := s.scope.SecretsManager.DeleteSecret(&secretsmanager.DeleteSecretInput{
		SecretId:                   aws.String(secretArn),
		ForceDeleteWithoutRecovery: aws.Bool(true),
	})

	if awserrors.IsNotFound(err) {
		return nil
	}

	return err
}

func (s *Service) secretName(m *scope.MachineScope) string {
	prefix := strings.Join(
		[]string{
			parameterPrefix,
			namespacesPrefix,
			s.scope.Namespace(),
			clustersPrefix,
			s.scope.Name(),
			instancesPrefix,
			m.Name(),
		},
		"/",
	) + "-"

	// apirand uses 27 runes, 27^54 is closest to 2^256 for 256-bits of entropy.
	randomStr := apirand.String(54)

	return prefix + randomStr
}
