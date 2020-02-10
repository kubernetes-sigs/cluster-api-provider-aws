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
	"fmt"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/uuid"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha2"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/converters"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
)

const (
	entryPrefix = "aws.cluster.x-k8s.io"

	// we set the max secret size to well below the 10240 byte limit, because this is limit after base64 encoding,
	// but the aws sdk handles encoding for us, so we can't send a full 10240.
	maxSecretSizeBytes = 7000
)

// Create stores data in AWS Secrets Manager for a given machine, chunking at 10kb per secret. The prefix of the secret
// ARN and the number of chunks are returned.
func (s *Service) Create(m *scope.MachineScope, data []byte) (string, int32, error) {
	// Build the tags to apply to the secret.
	additionalTags := m.AdditionalTags()
	additionalTags[infrav1.ClusterAWSCloudProviderTagKey(s.scope.Name())] = string(infrav1.ResourceLifecycleOwned)
	tags := infrav1.Build(infrav1.BuildParams{
		ClusterName: s.scope.Name(),
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Name:        aws.String(m.Name()),
		Role:        aws.String(m.Role()),
		Additional:  additionalTags,
	})

	// Build the prefix.
	prefix := path.Join(entryPrefix, string(uuid.NewUUID()))

	// Split the data into chunks and create the secrets on demand.
	var err error
	chunks := int32(0)
	splitBytes(data, maxSecretSizeBytes, func(chunk []byte) {
		name := fmt.Sprintf("%s-%d", prefix, chunks)
		_, callErr := s.scope.SecretsManager.CreateSecret(&secretsmanager.CreateSecretInput{
			Name:         aws.String(name),
			SecretBinary: chunk,
			Tags:         converters.MapToSecretsManagerTags(tags),
		})
		if callErr != nil {
			err = kerrors.NewAggregate([]error{callErr})
			return
		}
		chunks++
	})

	return prefix, chunks, err
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
