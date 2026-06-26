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
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/smithy-go"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/ssm/mock_ssmiface"
)

func TestCreateHybridActivation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name        string
		params      *HybridActivationParams
		expect      func(m *mock_ssmiface.MockSSMAPIMockRecorder)
		wantErr     bool
		errContains string
		validate    func(t *testing.T, result *HybridActivationResult)
	}{
		{
			name:        "nil params returns error",
			params:      nil,
			wantErr:     true,
			errContains: "params cannot be nil",
		},
		{
			name: "missing IAMRoleName returns error",
			params: &HybridActivationParams{
				RegistrationLimit: 1,
				ExpirationHours:   7,
			},
			wantErr:     true,
			errContains: "IAMRoleName is required",
		},
		{
			name: "zero RegistrationLimit returns error",
			params: &HybridActivationParams{
				IAMRoleName:       "TestRole",
				RegistrationLimit: 0,
				ExpirationHours:   7,
			},
			wantErr:     true,
			errContains: "RegistrationLimit must be greater than 0",
		},
		{
			name: "zero ExpirationHours returns error",
			params: &HybridActivationParams{
				IAMRoleName:       "TestRole",
				RegistrationLimit: 1,
				ExpirationHours:   0,
			},
			wantErr:     true,
			errContains: "ExpirationHours must be greater than 0",
		},
		{
			name: "successful creation with minimal params",
			params: &HybridActivationParams{
				IAMRoleName:       "TestRole",
				RegistrationLimit: 1,
				ExpirationHours:   7,
			},
			expect: func(m *mock_ssmiface.MockSSMAPIMockRecorder) {
				m.CreateActivation(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, input *ssm.CreateActivationInput, optFns ...func(*ssm.Options)) (*ssm.CreateActivationOutput, error) {
						// Validate input - IamRole should be passed directly
						if aws.ToString(input.IamRole) != "TestRole" {
							t.Errorf("expected IAM role name to be 'TestRole', got %q", aws.ToString(input.IamRole))
						}
						if aws.ToInt32(input.RegistrationLimit) != 1 {
							t.Errorf("expected RegistrationLimit to be 1")
						}
						// Verify managed tag is present
						hasManagedTag := false
						for _, tag := range input.Tags {
							if aws.ToString(tag.Key) == TagKeyManaged && aws.ToString(tag.Value) == "true" {
								hasManagedTag = true
								break
							}
						}
						if !hasManagedTag {
							t.Errorf("expected managed tag to be present")
						}
						return &ssm.CreateActivationOutput{
							ActivationId:   aws.String("test-activation-id"),
							ActivationCode: aws.String("test-activation-code"),
						}, nil
					},
				)
			},
			wantErr: false,
			validate: func(t *testing.T, result *HybridActivationResult) {
				if result.ActivationID != "test-activation-id" {
					t.Errorf("expected ActivationID to be 'test-activation-id', got %q", result.ActivationID)
				}
				if result.ActivationCode != "test-activation-code" {
					t.Errorf("expected ActivationCode to be 'test-activation-code', got %q", result.ActivationCode)
				}
				// Verify expiration time is approximately 7 hours from now
				expectedExpiration := time.Now().Add(7 * time.Hour)
				if result.ExpirationTime.Before(expectedExpiration.Add(-time.Minute)) ||
					result.ExpirationTime.After(expectedExpiration.Add(time.Minute)) {
					t.Errorf("expected ExpirationTime to be approximately 7 hours from now")
				}
			},
		},
		{
			name: "successful creation with all params",
			params: &HybridActivationParams{
				IAMRoleName:         "TestRole",
				RegistrationLimit:   10,
				ExpirationHours:     14,
				DefaultInstanceName: "hybrid-node",
				Description:         "Test activation for hybrid nodes",
				ClusterName:         "test-cluster",
				Namespace:           "default",
				ConfigName:          "test-config",
				MachineName:         "test-machine",
				Tags: infrav1.Tags{
					"Environment": "test",
					"Team":        "platform",
				},
			},
			expect: func(m *mock_ssmiface.MockSSMAPIMockRecorder) {
				m.CreateActivation(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, input *ssm.CreateActivationInput, optFns ...func(*ssm.Options)) (*ssm.CreateActivationOutput, error) {
						// Validate optional fields are set
						if aws.ToString(input.DefaultInstanceName) != "hybrid-node" {
							t.Errorf("expected DefaultInstanceName to be set")
						}
						if aws.ToString(input.Description) != "Test activation for hybrid nodes" {
							t.Errorf("expected Description to be set")
						}
						// Verify all tags are present
						tagMap := make(map[string]string)
						for _, tag := range input.Tags {
							tagMap[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
						}
						if tagMap[infrav1.ClusterTagKey("test-cluster")] != string(infrav1.ResourceLifecycleOwned) {
							t.Errorf("expected cluster tag to be set")
						}
						if tagMap[TagKeyNodeadmConfig] != "default/test-config" {
							t.Errorf("expected nodeadmconfig tag to be set")
						}
						if tagMap[TagKeyMachine] != "test-machine" {
							t.Errorf("expected machine tag to be set, got %q", tagMap[TagKeyMachine])
						}
						if tagMap["Environment"] != "test" {
							t.Errorf("expected Environment tag to be set")
						}
						if tagMap["Team"] != "platform" {
							t.Errorf("expected Team tag to be set")
						}
						return &ssm.CreateActivationOutput{
							ActivationId:   aws.String("full-test-activation-id"),
							ActivationCode: aws.String("full-test-activation-code"),
						}, nil
					},
				)
			},
			wantErr: false,
			validate: func(t *testing.T, result *HybridActivationResult) {
				if result.ActivationID != "full-test-activation-id" {
					t.Errorf("expected ActivationID to be 'full-test-activation-id', got %q", result.ActivationID)
				}
			},
		},
		{
			name: "AWS API error is propagated",
			params: &HybridActivationParams{
				IAMRoleName:       "TestRole",
				RegistrationLimit: 1,
				ExpirationHours:   7,
			},
			expect: func(m *mock_ssmiface.MockSSMAPIMockRecorder) {
				m.CreateActivation(gomock.Any(), gomock.Any()).Return(nil, &smithy.GenericAPIError{
					Code:    "ValidationException",
					Message: "Invalid IAM role ARN",
				})
			},
			wantErr:     true,
			errContains: "failed to create SSM activation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			scheme := runtime.NewScheme()
			_ = infrav1.AddToScheme(scheme)
			client := fake.NewClientBuilder().WithScheme(scheme).Build()

			clusterScope, err := getClusterScope(client)
			g.Expect(err).NotTo(HaveOccurred())

			ssmClientMock := mock_ssmiface.NewMockSSMAPI(mockCtrl)
			if tt.expect != nil {
				tt.expect(ssmClientMock.EXPECT())
			}

			s := NewService(clusterScope)
			s.SSMClient = ssmClientMock

			result, err := s.CreateHybridActivation(context.Background(), tt.params)

			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
				if tt.errContains != "" {
					g.Expect(err.Error()).To(ContainSubstring(tt.errContains))
				}
				return
			}

			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(result).NotTo(BeNil())

			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

func TestDeleteHybridActivation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name         string
		activationID string
		expect       func(m *mock_ssmiface.MockSSMAPIMockRecorder)
		wantErr      bool
		errContains  string
	}{
		{
			name:         "empty activationID returns error",
			activationID: "",
			wantErr:      true,
			errContains:  "activationID is required",
		},
		{
			name:         "successful deletion",
			activationID: "test-activation-id",
			expect: func(m *mock_ssmiface.MockSSMAPIMockRecorder) {
				m.DeleteActivation(gomock.Any(), gomock.Eq(&ssm.DeleteActivationInput{
					ActivationId: aws.String("test-activation-id"),
				})).Return(&ssm.DeleteActivationOutput{}, nil)
			},
			wantErr: false,
		},
		{
			name:         "activation not found is not an error (idempotent)",
			activationID: "nonexistent-activation",
			expect: func(m *mock_ssmiface.MockSSMAPIMockRecorder) {
				m.DeleteActivation(gomock.Any(), gomock.Any()).Return(nil, &mockAPIError{
					Code:    "InvalidActivation",
					Message: "Activation not found",
				})
			},
			wantErr: false,
		},
		{
			name:         "other AWS errors are propagated",
			activationID: "test-activation-id",
			expect: func(m *mock_ssmiface.MockSSMAPIMockRecorder) {
				m.DeleteActivation(gomock.Any(), gomock.Any()).Return(nil, &smithy.GenericAPIError{
					Code:    "AccessDenied",
					Message: "Access denied",
				})
			},
			wantErr:     true,
			errContains: "failed to delete SSM activation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			scheme := runtime.NewScheme()
			_ = infrav1.AddToScheme(scheme)
			client := fake.NewClientBuilder().WithScheme(scheme).Build()

			clusterScope, err := getClusterScope(client)
			g.Expect(err).NotTo(HaveOccurred())

			ssmClientMock := mock_ssmiface.NewMockSSMAPI(mockCtrl)
			if tt.expect != nil {
				tt.expect(ssmClientMock.EXPECT())
			}

			s := NewService(clusterScope)
			s.SSMClient = ssmClientMock

			err = s.DeleteHybridActivation(context.Background(), tt.activationID)

			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
				if tt.errContains != "" {
					g.Expect(err.Error()).To(ContainSubstring(tt.errContains))
				}
				return
			}

			g.Expect(err).NotTo(HaveOccurred())
		})
	}
}
