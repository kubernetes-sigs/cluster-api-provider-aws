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

package identity

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	ststypes "github.com/aws/aws-sdk-go-v2/service/sts/types"
	"github.com/aws/smithy-go"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/sts/mock_stsiface"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
)

func TestAWSStaticPrincipalTypeProvider(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	secret := &corev1.Secret{
		Data: map[string][]byte{
			"AccessKeyID":     []byte("static-AccessKeyID"),
			"SecretAccessKey": []byte("static-SecretAccessKey"),
		},
	}

	staticProvider := NewAWSStaticPrincipalTypeProvider(&infrav1.AWSClusterStaticIdentity{}, secret)

	stsMock := mock_stsiface.NewMockSTSClient(mockCtrl)
	roleIdentity := &infrav1.AWSClusterRoleIdentity{
		Spec: infrav1.AWSClusterRoleIdentitySpec{
			AWSRoleSpec: infrav1.AWSRoleSpec{
				RoleArn:         "arn:*:iam::*:role/aws-role/firstroleprovider",
				SessionName:     "first-role-provider-session",
				DurationSeconds: 900,
			},
		},
	}

	expiresAt := time.Now()

	roleProvider := &AWSRolePrincipalTypeProvider{
		credentials:    nil,
		Principal:      roleIdentity,
		region:         "us-west-2",
		sourceProvider: staticProvider,
		stsClient:      stsMock,
	}

	roleIdentity2 := &infrav1.AWSClusterRoleIdentity{
		Spec: infrav1.AWSClusterRoleIdentitySpec{
			AWSRoleSpec: infrav1.AWSRoleSpec{
				RoleArn:         "arn:*:iam::*:role/aws-role/secondroleprovider",
				SessionName:     "second-role-provider-session",
				DurationSeconds: 900,
			},
		},
	}

	roleProvider2 := &AWSRolePrincipalTypeProvider{
		credentials:    nil,
		Principal:      roleIdentity2,
		region:         "us-west-2",
		sourceProvider: roleProvider,
		stsClient:      stsMock,
	}

	testCases := []struct {
		name      string
		provider  AWSPrincipalTypeProvider
		expect    func(m *mock_stsiface.MockSTSClientMockRecorder)
		expectErr bool
		value     aws.Credentials
	}{
		{
			name:      "Static provider successfully retrieves",
			provider:  staticProvider,
			expect:    func(m *mock_stsiface.MockSTSClientMockRecorder) {},
			expectErr: false,
			value: aws.Credentials{
				AccessKeyID:     "static-AccessKeyID",
				SecretAccessKey: "static-SecretAccessKey",
				Source:          "StaticCredentials",
				CanExpire:       false,
				Expires:         time.Time{},
			},
		},
		{
			name:     "Role provider with static provider source successfully retrieves",
			provider: roleProvider,
			expect: func(m *mock_stsiface.MockSTSClientMockRecorder) {
				m.AssumeRole(gomock.Any(), &sts.AssumeRoleInput{
					RoleArn:         aws.String(roleIdentity.Spec.RoleArn),
					RoleSessionName: aws.String(roleIdentity.Spec.SessionName),
					DurationSeconds: aws.Int32(roleIdentity.Spec.DurationSeconds),
				}).Return(&sts.AssumeRoleOutput{
					Credentials: &ststypes.Credentials{
						AccessKeyId:     aws.String("assumedAccessKeyId"),
						SecretAccessKey: aws.String("assumedSecretAccessKey"),
						SessionToken:    aws.String("assumedSessionToken"),
						Expiration:      aws.Time(expiresAt),
					},
				}, nil)
			},
			expectErr: false,
			value: aws.Credentials{
				AccessKeyID:     "assumedAccessKeyId",
				SecretAccessKey: "assumedSecretAccessKey",
				SessionToken:    "assumedSessionToken",
				Source:          "AssumeRoleProvider",
				CanExpire:       true,
				Expires:         expiresAt,
			},
		},
		{
			name:     "Role provider with role provider source successfully retrieves",
			provider: roleProvider2,
			expect: func(m *mock_stsiface.MockSTSClientMockRecorder) {
				m.AssumeRole(gomock.Any(), &sts.AssumeRoleInput{
					RoleArn:         aws.String(roleIdentity.Spec.RoleArn),
					RoleSessionName: aws.String(roleIdentity.Spec.SessionName),
					DurationSeconds: aws.Int32(roleIdentity.Spec.DurationSeconds),
				}).Return(&sts.AssumeRoleOutput{
					Credentials: &ststypes.Credentials{
						AccessKeyId:     aws.String("assumedAccessKeyId"),
						SecretAccessKey: aws.String("assumedSecretAccessKey"),
						SessionToken:    aws.String("assumedSessionToken"),
						Expiration:      aws.Time(time.Now().AddDate(+1, 0, 0)),
					},
				}, nil)

				m.AssumeRole(gomock.Any(), &sts.AssumeRoleInput{
					RoleArn:         aws.String(roleIdentity2.Spec.RoleArn),
					RoleSessionName: aws.String(roleIdentity2.Spec.SessionName),
					DurationSeconds: aws.Int32(roleIdentity2.Spec.DurationSeconds),
				}).Return(&sts.AssumeRoleOutput{
					Credentials: &ststypes.Credentials{
						AccessKeyId:     aws.String("assumedAccessKeyId2"),
						SecretAccessKey: aws.String("assumedSecretAccessKey2"),
						SessionToken:    aws.String("assumedSessionToken2"),
						Expiration:      aws.Time(expiresAt),
					},
				}, nil)
			},
			expectErr: false,
			value: aws.Credentials{
				AccessKeyID:     "assumedAccessKeyId2",
				SecretAccessKey: "assumedSecretAccessKey2",
				SessionToken:    "assumedSessionToken2",
				Source:          "AssumeRoleProvider",
				CanExpire:       true,
				Expires:         expiresAt,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			tc.expect(stsMock.EXPECT())
			value, err := tc.provider.Retrieve(context.TODO())
			if tc.expectErr {
				g.Expect(err).ToNot(BeNil())
				return
			}

			g.Expect(err).To(BeNil())

			if !cmp.Equal(tc.value, value) {
				t.Fatal("Did not get expected result")
			}
		})
	}
}

func TestAWSRolePrincipalTypeProvider_ClearsCredentialCacheOnInvalidClientTokenId(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	g := NewWithT(t)

	secret := &corev1.Secret{
		Data: map[string][]byte{
			"AccessKeyID":     []byte("test-key"),
			"SecretAccessKey": []byte("test-secret"),
		},
	}
	staticProvider := NewAWSStaticPrincipalTypeProvider(&infrav1.AWSClusterStaticIdentity{}, secret)

	stsMock := mock_stsiface.NewMockSTSClient(mockCtrl)
	roleIdentity := &infrav1.AWSClusterRoleIdentity{
		Spec: infrav1.AWSClusterRoleIdentitySpec{
			AWSRoleSpec: infrav1.AWSRoleSpec{
				RoleArn:         "arn:aws:iam::123456789012:role/test-role",
				SessionName:     "test-session",
				DurationSeconds: 900,
			},
		},
	}

	testLog := logger.NewLogger(klog.NewKlogr())
	roleProvider := &AWSRolePrincipalTypeProvider{
		credentials:    nil,
		Principal:      roleIdentity,
		region:         "us-east-1",
		sourceProvider: staticProvider,
		stsClient:      stsMock,
		log:            testLog.WithName("test"),
	}

	// First call: return InvalidClientTokenId error
	invalidTokenErr := &smithy.GenericAPIError{
		Code:    "InvalidClientTokenId",
		Message: "The security token included in the request is invalid",
	}
	stsMock.EXPECT().AssumeRole(gomock.Any(), gomock.Any()).Return(nil, invalidTokenErr)

	_, err := roleProvider.Retrieve(context.TODO())
	g.Expect(err).To(HaveOccurred())
	g.Expect(err.Error()).To(ContainSubstring("InvalidClientTokenId"))
	// Credentials cache should be cleared
	g.Expect(roleProvider.credentials).To(BeNil())

	// Second call: succeed — proves the cache was cleared and a fresh attempt is made
	stsMock.EXPECT().AssumeRole(gomock.Any(), gomock.Any()).Return(&sts.AssumeRoleOutput{
		Credentials: &ststypes.Credentials{
			AccessKeyId:     aws.String("new-key"),
			SecretAccessKey: aws.String("new-secret"),
			SessionToken:    aws.String("new-token"),
			Expiration:      aws.Time(time.Now().Add(1 * time.Hour)),
		},
	}, nil)

	creds, err := roleProvider.Retrieve(context.TODO())
	g.Expect(err).To(BeNil())
	g.Expect(creds.AccessKeyID).To(Equal("new-key"))
}

func TestIsInvalidClientTokenIDError(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		expect bool
	}{
		{
			name:   "nil error",
			err:    nil,
			expect: false,
		},
		{
			name: "typed smithy APIError with InvalidClientTokenId",
			err: &smithy.GenericAPIError{
				Code:    "InvalidClientTokenId",
				Message: "The security token included in the request is invalid",
			},
			expect: true,
		},
		{
			name:   "error containing InvalidClientTokenId in message",
			err:    fmt.Errorf("operation failed: InvalidClientTokenId: token not valid"),
			expect: true,
		},
		{
			name:   "unrelated error",
			err:    fmt.Errorf("connection timeout"),
			expect: false,
		},
		{
			name: "different API error code",
			err: &smithy.GenericAPIError{
				Code:    "AccessDenied",
				Message: "not authorized",
			},
			expect: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			g.Expect(IsInvalidClientTokenIDError(tt.err)).To(Equal(tt.expect))
		})
	}
}
