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

package identityprovider

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"k8s.io/klog/v2"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/eks/mock_eksiface"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
)

func TestEKSAddonPlan(t *testing.T) {
	clusterName := "default.cluster"
	identityProviderARN := "aws:mock:provider:arn"
	idnetityProviderName := "IdentityProviderConfigName"
	log := logger.NewLogger(klog.Background())

	testCases := []struct {
		name                    string
		currentIdentityProvider *OidcIdentityProviderConfig
		desiredIdentityProvider *OidcIdentityProviderConfig
		expect                  func(m *mock_eksiface.MockEKSAPIMockRecorder)
		expectCreateError       bool
		expectDoError           bool
	}{
		{
			name: "no desired and no installed",
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				// Do nothing
			},
			expectCreateError: false,
			expectDoError:     false,
		},
		{
			name: "no installed and 1 desired",
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.AssociateIdentityProviderConfig(gomock.Eq(context.TODO()), gomock.Eq(&eks.AssociateIdentityProviderConfigInput{
					ClusterName: aws.String(clusterName),
					Oidc:        createDesiredIdentityProviderRequest(aws.String(idnetityProviderName)),
					Tags:        createTags(),
				}))
			},
			desiredIdentityProvider: createDesiredIdentityProvider(idnetityProviderName, createTags()),
			expectCreateError:       false,
			expectDoError:           false,
		},
		{
			name:                    "1 installed and 1 desired - both same and installed active",
			currentIdentityProvider: createCurrentIdentityProvider(idnetityProviderName, identityProviderARN, string(ekstypes.ConfigStatusActive), createTags()),
			desiredIdentityProvider: createDesiredIdentityProvider(idnetityProviderName, createTags()),
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {

			},
			expectCreateError: false,
			expectDoError:     false,
		},
		{
			name:                    "1 installed and 1 desired - both same and installed is creating",
			currentIdentityProvider: createCurrentIdentityProvider(idnetityProviderName, identityProviderARN, string(ekstypes.ConfigStatusCreating), createTags()),
			desiredIdentityProvider: createDesiredIdentityProvider(idnetityProviderName, createTags()),
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.DescribeIdentityProviderConfig(gomock.Eq(context.TODO()),
					gomock.Eq(&eks.DescribeIdentityProviderConfigInput{
						ClusterName: aws.String(clusterName),
						IdentityProviderConfig: &ekstypes.IdentityProviderConfig{
							Name: aws.String("IdentityProviderConfigName"),
							Type: oidcType,
						},
					})).
					Return(&eks.DescribeIdentityProviderConfigOutput{
						IdentityProviderConfig: &ekstypes.IdentityProviderConfigResponse{
							Oidc: &ekstypes.OidcIdentityProviderConfig{
								ClusterName:                aws.String(clusterName),
								IdentityProviderConfigArn:  aws.String(identityProviderARN),
								IdentityProviderConfigName: aws.String("IdentityProviderConfigName"),
								Status:                     ekstypes.ConfigStatusActive,
								Tags:                       createTags(),
							},
						},
					}, nil)
			},
			expectCreateError: false,
			expectDoError:     false,
		},

		{
			name:                    "1 installed and 1 desired - both same and installed is active",
			currentIdentityProvider: createCurrentIdentityProvider(idnetityProviderName, identityProviderARN, string(ekstypes.ConfigStatusActive), createTags()),
			desiredIdentityProvider: createDesiredIdentityProvider(idnetityProviderName, createTags()),
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {

			},
			expectCreateError: false,
			expectDoError:     false,
		},
		{
			name:                    "1 installed and 1 desired - both same and installed is active, and tags added",
			currentIdentityProvider: createCurrentIdentityProvider(idnetityProviderName, identityProviderARN, string(ekstypes.ConfigStatusActive), createTags()),
			desiredIdentityProvider: createDesiredIdentityProvider(idnetityProviderName, changeTags(createTags())),
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.TagResource(gomock.Eq(context.TODO()), gomock.Eq(&eks.TagResourceInput{
					ResourceArn: aws.String(identityProviderARN),
					Tags:        changeTags(createTags()),
				}))
			},
			expectCreateError: false,
			expectDoError:     false,
		},
		{
			name:                    "1 installed and 1 desired - both same and installed is active, and tags removed",
			currentIdentityProvider: createCurrentIdentityProvider(idnetityProviderName, identityProviderARN, string(ekstypes.ConfigStatusActive), createTags()),
			desiredIdentityProvider: createDesiredIdentityProvider(idnetityProviderName, nil),
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.UntagResource(gomock.Eq(context.TODO()), gomock.Eq(&eks.UntagResourceInput{
					ResourceArn: aws.String(identityProviderARN),
					TagKeys:     []string{"key1"},
				}))
			},
			expectCreateError: false,
			expectDoError:     false,
		},

		{
			name:                    "1 installed and 0 desired - installed provider is removed",
			currentIdentityProvider: createCurrentIdentityProvider(idnetityProviderName, identityProviderARN, string(ekstypes.ConfigStatusActive), createTags()),
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.DisassociateIdentityProviderConfig(gomock.Eq(context.TODO()),
					gomock.Eq(&eks.DisassociateIdentityProviderConfigInput{
						ClusterName: aws.String(clusterName),
						IdentityProviderConfig: &ekstypes.IdentityProviderConfig{
							Name: aws.String("IdentityProviderConfigName"),
							Type: oidcType,
						},
					}))
			},
			expectCreateError: false,
			expectDoError:     false,
		},

		{
			name:                    "1 installed and desired client id changed - installed provider is removed",
			currentIdentityProvider: createCurrentIdentityProvider(idnetityProviderName, identityProviderARN, string(ekstypes.ConfigStatusActive), createTags()),
			desiredIdentityProvider: createDesiredIdentityProviderWithDifferentClientID(idnetityProviderName, createTags()),
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.DisassociateIdentityProviderConfig(gomock.Eq(context.TODO()),
					gomock.Eq(&eks.DisassociateIdentityProviderConfigInput{
						ClusterName: aws.String(clusterName),
						IdentityProviderConfig: &ekstypes.IdentityProviderConfig{
							Name: aws.String("IdentityProviderConfigName"),
							Type: oidcType,
						},
					}))
			},
			expectCreateError: false,
			expectDoError:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			mockControl := gomock.NewController(t)
			defer mockControl.Finish()

			eksMock := mock_eksiface.NewMockEKSAPI(mockControl)
			tc.expect(eksMock.EXPECT())

			ctx := context.TODO()

			planner := NewPlan(clusterName, tc.currentIdentityProvider, tc.desiredIdentityProvider, eksMock, log)
			procedures, err := planner.Create(ctx)
			if tc.expectCreateError {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).To(BeNil())
			g.Expect(procedures).NotTo(BeNil())

			for _, proc := range procedures {
				procErr := proc.Do(ctx)
				if tc.expectDoError {
					g.Expect(procErr).To(HaveOccurred())
					return
				}
				g.Expect(procErr).To(BeNil())
			}
		})
	}
}

func createTags() infrav1.Tags {
	tags := infrav1.Tags{}
	tags["key1"] = "value1"
	return tags
}

func createDesiredIdentityProvider(name string, tags infrav1.Tags) *OidcIdentityProviderConfig {
	return &OidcIdentityProviderConfig{
		ClientID:                   "clientId",
		IdentityProviderConfigName: name,
		IssuerURL:                  "http://IssuerURL.com",
		RequiredClaims:             make(map[string]string),
		Tags:                       tags,
	}
}

func createCurrentIdentityProvider(name string, arn, status string, tags infrav1.Tags) *OidcIdentityProviderConfig {
	config := createDesiredIdentityProvider(name, tags)
	config.IdentityProviderConfigArn = arn
	config.Status = status

	return config
}

func changeTags(original infrav1.Tags) infrav1.Tags {
	original["key2"] = "value2"
	return original
}

func createDesiredIdentityProviderRequest(name *string) *ekstypes.OidcIdentityProviderConfigRequest {
	return &ekstypes.OidcIdentityProviderConfigRequest{
		ClientId:                   aws.String("clientId"),
		IdentityProviderConfigName: name,
		IssuerUrl:                  aws.String("http://IssuerURL.com"),
		RequiredClaims:             make(map[string]string),
		GroupsClaim:                aws.String(""),
		GroupsPrefix:               aws.String(""),
		UsernameClaim:              aws.String(""),
		UsernamePrefix:             aws.String(""),
	}
}

func createDesiredIdentityProviderWithDifferentClientID(name string, tags infrav1.Tags) *OidcIdentityProviderConfig {
	p := createDesiredIdentityProvider(name, tags)
	p.ClientID = "clientId2"
	return p
}
