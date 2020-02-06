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
	"bytes"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
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

	// we set the max secret size to well below the 10240 byte limit, because this is limit after base64 encoding,
	// but the aws sdk handles encoding for us, so we can't send a full 10240.
	maxSecretSizeBytes = 7000
)

// Create stores data in AWS Secrets Manager for a given machine, chunking at 10kb per secret. The prefix of the secret
// ARN and the number of chunks are returned.
func (s *Service) Create(m *scope.MachineScope, data []byte) (string, int32, error) {
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

	chunks := len(data) / maxSecretSizeBytes
	remainder := chunks % maxSecretSizeBytes
	if remainder != 0 {
		chunks++
	}

	buf := bytes.NewBuffer(data)
	prefix := s.secretNamePrefix(m)

	for i := 0; i < chunks; i++ {
		name := fmt.Sprintf("%s-%d", prefix, i)

		chunk := buf.Next(maxSecretSizeBytes)

		_, err := s.scope.SecretsManager.CreateSecret(&secretsmanager.CreateSecretInput{
			Name:         aws.String(name),
			SecretBinary: chunk,
			Tags:         converters.MapToSecretsManagerTags(tags),
		})

		if err != nil {
			return "", 0, err
		}
	}

	return prefix, int32(chunks), nil
}

// Delete the secret belonging to a machine from AWS Secrets Manager
func (s *Service) Delete(m *scope.MachineScope) error {
	var errors []error

	for i := int32(0); i < m.GetSecretCount(); i++ {
		_, err := s.scope.SecretsManager.DeleteSecret(&secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String(fmt.Sprintf("%s-%d", m.GetSecretPrefix(), i)),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})

		if awserrors.IsNotFound(err) {
			continue
		}
		if err != nil {
			errors = append(errors, err)
		}
	}

	return kerrors.NewAggregate(errors)
}

func (s *Service) secretNamePrefix(m *scope.MachineScope) string {
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
