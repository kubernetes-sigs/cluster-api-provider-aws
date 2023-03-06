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
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/sts/mock_stsiface"
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

	var staticProvider AWSPrincipalTypeProvider = NewAWSStaticPrincipalTypeProvider(&infrav1.AWSClusterStaticIdentity{}, secret)

	stsMock := mock_stsiface.NewMockSTSAPI(mockCtrl)
	roleIdentity := &infrav1.AWSClusterRoleIdentity{
		Spec: infrav1.AWSClusterRoleIdentitySpec{
			AWSRoleSpec: infrav1.AWSRoleSpec{
				RoleArn:         "arn:*:iam::*:role/aws-role/firstroleprovider",
				SessionName:     "first-role-provider-session",
				DurationSeconds: 900,
			},
		},
	}

	var roleProvider AWSPrincipalTypeProvider = &AWSRolePrincipalTypeProvider{
		credentials:    nil,
		Principal:      roleIdentity,
		sourceProvider: &staticProvider,
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

	var roleProvider2 AWSPrincipalTypeProvider = &AWSRolePrincipalTypeProvider{
		credentials:    nil,
		Principal:      roleIdentity2,
		sourceProvider: &roleProvider,
		stsClient:      stsMock,
	}

	testCases := []struct {
		name      string
		provider  AWSPrincipalTypeProvider
		expect    func(m *mock_stsiface.MockSTSAPIMockRecorder)
		expectErr bool
		value     credentials.Value
	}{
		{
			name:      "Static provider successfully retrieves",
			provider:  staticProvider,
			expect:    func(m *mock_stsiface.MockSTSAPIMockRecorder) {},
			expectErr: false,
			value: credentials.Value{
				AccessKeyID:     "static-AccessKeyID",
				SecretAccessKey: "static-SecretAccessKey",
				ProviderName:    "StaticProvider",
			},
		},
		{
			name:     "Role provider with static provider source successfully retrieves",
			provider: roleProvider,
			expect: func(m *mock_stsiface.MockSTSAPIMockRecorder) {
				m.AssumeRoleWithContext(gomock.Any(), &sts.AssumeRoleInput{
					RoleArn:         aws.String(roleIdentity.Spec.RoleArn),
					RoleSessionName: aws.String(roleIdentity.Spec.SessionName),
					DurationSeconds: pointer.Int64(int64(roleIdentity.Spec.DurationSeconds)),
				}).Return(&sts.AssumeRoleOutput{
					Credentials: &sts.Credentials{
						AccessKeyId:     aws.String("assumedAccessKeyId"),
						SecretAccessKey: aws.String("assumedSecretAccessKey"),
						SessionToken:    aws.String("assumedSessionToken"),
						Expiration:      aws.Time(time.Now()),
					},
				}, nil)
			},
			expectErr: false,
			value: credentials.Value{
				AccessKeyID:     "assumedAccessKeyId",
				SecretAccessKey: "assumedSecretAccessKey",
				SessionToken:    "assumedSessionToken",
				ProviderName:    "AssumeRoleProvider",
			},
		},
		{
			name:     "Role provider with role provider source successfully retrieves",
			provider: roleProvider2,
			expect: func(m *mock_stsiface.MockSTSAPIMockRecorder) {
				m.AssumeRoleWithContext(gomock.Any(), &sts.AssumeRoleInput{
					RoleArn:         aws.String(roleIdentity.Spec.RoleArn),
					RoleSessionName: aws.String(roleIdentity.Spec.SessionName),
					DurationSeconds: pointer.Int64(int64(roleIdentity.Spec.DurationSeconds)),
				}).Return(&sts.AssumeRoleOutput{
					Credentials: &sts.Credentials{
						AccessKeyId:     aws.String("assumedAccessKeyId"),
						SecretAccessKey: aws.String("assumedSecretAccessKey"),
						SessionToken:    aws.String("assumedSessionToken"),
						Expiration:      aws.Time(time.Now().AddDate(+1, 0, 0)),
					},
				}, nil)

				m.AssumeRoleWithContext(gomock.Any(), &sts.AssumeRoleInput{
					RoleArn:         aws.String(roleIdentity2.Spec.RoleArn),
					RoleSessionName: aws.String(roleIdentity2.Spec.SessionName),
					DurationSeconds: pointer.Int64(int64(roleIdentity2.Spec.DurationSeconds)),
				}).Return(&sts.AssumeRoleOutput{
					Credentials: &sts.Credentials{
						AccessKeyId:     aws.String("assumedAccessKeyId2"),
						SecretAccessKey: aws.String("assumedSecretAccessKey2"),
						SessionToken:    aws.String("assumedSessionToken2"),
						Expiration:      aws.Time(time.Now()),
					},
				}, nil)
			},
			expectErr: false,
			value: credentials.Value{
				AccessKeyID:     "assumedAccessKeyId2",
				SecretAccessKey: "assumedSecretAccessKey2",
				SessionToken:    "assumedSessionToken2",
				ProviderName:    "AssumeRoleProvider",
			},
		},
		{
			name:     "Role provider with role provider source fails to retrieve when the source's source cannot assume source",
			provider: roleProvider2,
			expect: func(m *mock_stsiface.MockSTSAPIMockRecorder) {
				roleProvider.(*AWSRolePrincipalTypeProvider).credentials.Expire()
				roleProvider2.(*AWSRolePrincipalTypeProvider).credentials.Expire()
				// AssumeRoleWithContext() call is not needed for roleIdentity as it has unexpired credentials
				m.AssumeRoleWithContext(gomock.Any(), &sts.AssumeRoleInput{
					RoleArn:         aws.String(roleIdentity.Spec.RoleArn),
					RoleSessionName: aws.String(roleIdentity.Spec.SessionName),
					DurationSeconds: pointer.Int64(int64(roleIdentity.Spec.DurationSeconds)),
				}).Return(&sts.AssumeRoleOutput{}, errors.New("Not authorized to assume role"))
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			tc.expect(stsMock.EXPECT())
			value, err := tc.provider.Retrieve()
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
