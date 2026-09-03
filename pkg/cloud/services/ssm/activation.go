/*
Copyright 2026 The Kubernetes Authors.

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
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/aws/smithy-go"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
)

const (
	// TagKeyNodeadmConfig is the tag key for the NodeadmConfig reference.
	// Uses the standard CAPA provider prefix.
	TagKeyNodeadmConfig = infrav1.NameAWSProviderPrefix + "nodeadmconfig"
	// TagKeyManaged is the tag key indicating the resource is managed by CAPA.
	// Uses the standard CAPA provider prefix.
	TagKeyManaged = infrav1.NameAWSProviderPrefix + "managed"
	// TagKeyMachine is the tag key for the Machine name that owns this activation.
	// Uses the cluster-api standard annotation key so the activation can be mapped
	// back to its owning Machine across the CAPI ecosystem.
	TagKeyMachine = clusterv1.MachineAnnotation
)

// HybridActivationParams contains parameters for creating an SSM hybrid activation.
type HybridActivationParams struct {
	// IAMRoleName is the IAM role name for hybrid nodes to assume.
	IAMRoleName string

	// RegistrationLimit is the max number of nodes that can use this activation.
	RegistrationLimit int32

	// ExpirationHours is the number of hours until expiration.
	ExpirationHours int32

	// DefaultInstanceName is the default name for registered instances.
	DefaultInstanceName string

	// Description for the activation.
	Description string

	// Tags to apply to the activation.
	Tags infrav1.Tags

	// ClusterName for tagging purposes.
	ClusterName string

	// Namespace of the NodeadmConfig.
	Namespace string

	// ConfigName is the name of the NodeadmConfig.
	ConfigName string

	// MachineName is the name of the Machine that owns this bootstrap config.
	// When set, a tag with key TagKeyMachine is added to the activation so the
	// activation can be mapped back to its owning Machine.
	MachineName string
}

// HybridActivationResult contains the result of creating an SSM activation.
type HybridActivationResult struct {
	// ActivationID is the unique identifier for the activation.
	ActivationID string

	// ActivationCode is the secret code for the activation.
	// This is only returned once at creation time.
	ActivationCode string

	// ExpirationTime is when the activation expires.
	ExpirationTime time.Time
}

// CreateHybridActivation creates an SSM activation for EKS hybrid nodes.
func (s *Service) CreateHybridActivation(ctx context.Context, params *HybridActivationParams) (*HybridActivationResult, error) {
	if params == nil {
		return nil, errors.New("params cannot be nil")
	}
	if params.IAMRoleName == "" {
		return nil, errors.New("IAMRoleName is required")
	}
	if params.RegistrationLimit <= 0 {
		return nil, errors.New("RegistrationLimit must be greater than 0")
	}
	if params.ExpirationHours <= 0 {
		return nil, errors.New("ExpirationHours must be greater than 0")
	}

	// Calculate expiration time
	expirationTime := time.Now().Add(time.Duration(params.ExpirationHours) * time.Hour)

	// Build tags - start with user-provided tags
	tags := make([]ssmtypes.Tag, 0, len(params.Tags)+3)
	for k, v := range params.Tags {
		tags = append(tags, ssmtypes.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}

	// Add standard CAPA tags
	if params.ClusterName != "" {
		tags = append(tags, ssmtypes.Tag{
			Key:   aws.String(infrav1.ClusterTagKey(params.ClusterName)),
			Value: aws.String(string(infrav1.ResourceLifecycleOwned)),
		})
	}
	if params.Namespace != "" && params.ConfigName != "" {
		tags = append(tags, ssmtypes.Tag{
			Key:   aws.String(TagKeyNodeadmConfig),
			Value: aws.String(fmt.Sprintf("%s/%s", params.Namespace, params.ConfigName)),
		})
	}
	if params.MachineName != "" {
		tags = append(tags, ssmtypes.Tag{
			Key:   aws.String(TagKeyMachine),
			Value: aws.String(params.MachineName),
		})
	}
	tags = append(tags, ssmtypes.Tag{
		Key:   aws.String(TagKeyManaged),
		Value: aws.String("true"),
	})

	input := &ssm.CreateActivationInput{
		IamRole:           aws.String(params.IAMRoleName),
		RegistrationLimit: aws.Int32(params.RegistrationLimit),
		ExpirationDate:    aws.Time(expirationTime),
		Tags:              tags,
	}

	// Set optional fields if provided
	if params.DefaultInstanceName != "" {
		input.DefaultInstanceName = aws.String(params.DefaultInstanceName)
	}
	if params.Description != "" {
		input.Description = aws.String(params.Description)
	}

	output, err := s.SSMClient.CreateActivation(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create SSM activation: %w", err)
	}

	return &HybridActivationResult{
		ActivationID:   aws.ToString(output.ActivationId),
		ActivationCode: aws.ToString(output.ActivationCode),
		ExpirationTime: expirationTime,
	}, nil
}

// DeleteHybridActivation deletes an SSM activation.
// It is idempotent - if the activation doesn't exist, no error is returned.
func (s *Service) DeleteHybridActivation(ctx context.Context, activationID string) error {
	if activationID == "" {
		return errors.New("activationID is required")
	}

	input := &ssm.DeleteActivationInput{
		ActivationId: aws.String(activationID),
	}

	_, err := s.SSMClient.DeleteActivation(ctx, input)
	if err != nil {
		// Check if activation doesn't exist (already deleted) - this is not an error
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			// InvalidActivation is returned when the activation ID is not found
			if apiErr.ErrorCode() == "InvalidActivation" {
				return nil
			}
		}
		return fmt.Errorf("failed to delete SSM activation: %w", err)
	}

	return nil
}
