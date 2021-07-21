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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"k8s.io/klog/v2/klogr"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha4"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/eks/mock_eksiface"
)

func TestEKSAddonPlan(t *testing.T) {
	clusterName := "default.cluster"
	identityProviderARN := "aws:mock:provider:arn"
	idnetityProviderName := "IdentityProviderConfigName"
	log := klogr.New()

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
				m.AssociateIdentityProviderConfigWithContext(gomock.Eq(context.TODO()), gomock.Eq(&eks.AssociateIdentityProviderConfigInput{
					ClusterName: aws.String(clusterName),
					Oidc:        createDesiredIdentityProviderRequest(aws.String(idnetityProviderName)),
					Tags:        aws.StringMap(createTags()),
				}))
			},
			desiredIdentityProvider: createDesiredIdentityProvider(idnetityProviderName, createTags()),
			expectCreateError:       false,
			expectDoError:           false,
		},
		{
			name:                    "1 installed and 1 desired - both same and installed active",
			currentIdentityProvider: createCurrentIdentityProvider(idnetityProviderName, identityProviderARN, eks.ConfigStatusActive, createTags()),
			desiredIdentityProvider: createDesiredIdentityProvider(idnetityProviderName, createTags()),
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {

			},
			expectCreateError: false,
			expectDoError:     false,
		},
		{
			name:                    "1 installed and 1 desired - both same and installed is creating",
			currentIdentityProvider: createCurrentIdentityProvider(idnetityProviderName, identityProviderARN, eks.ConfigStatusCreating, createTags()),
			desiredIdentityProvider: createDesiredIdentityProvider(idnetityProviderName, createTags()),
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.DescribeIdentityProviderConfigWithContext(gomock.Eq(context.TODO()),
					gomock.Eq(&eks.DescribeIdentityProviderConfigInput{
						ClusterName: aws.String(clusterName),
						IdentityProviderConfig: &eks.IdentityProviderConfig{
							Name: aws.String("IdentityProviderConfigName"),
							Type: oidcType,
						},
					})).
					Return(&eks.DescribeIdentityProviderConfigOutput{
						IdentityProviderConfig: &eks.IdentityProviderConfigResponse{
							Oidc: &eks.OidcIdentityProviderConfig{
								ClusterName:                aws.String(clusterName),
								IdentityProviderConfigArn:  aws.String(identityProviderARN),
								IdentityProviderConfigName: aws.String("IdentityProviderConfigName"),
								Status:                     aws.String(eks.ConfigStatusActive),
								Tags:                       aws.StringMap(createTags()),
							},
						},
					}, nil)
			},
			expectCreateError: false,
			expectDoError:     false,
		},

		{
			name:                    "1 installed and 1 desired - both same and installed is active",
			currentIdentityProvider: createCurrentIdentityProvider(idnetityProviderName, identityProviderARN, eks.ConfigStatusActive, createTags()),
			desiredIdentityProvider: createDesiredIdentityProvider(idnetityProviderName, createTags()),
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {

			},
			expectCreateError: false,
			expectDoError:     false,
		},
		{
			name:                    "1 installed and 1 desired - both same and installed is active, and tags added",
			currentIdentityProvider: createCurrentIdentityProvider(idnetityProviderName, identityProviderARN, eks.ConfigStatusActive, createTags()),
			desiredIdentityProvider: createDesiredIdentityProvider(idnetityProviderName, changeTags(createTags())),
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.TagResource(gomock.Eq(&eks.TagResourceInput{
					ResourceArn: aws.String(identityProviderARN),
					Tags:        aws.StringMap(changeTags(createTags())),
				}))
			},
			expectCreateError: false,
			expectDoError:     false,
		},
		{
			name:                    "1 installed and 1 desired - both same and installed is active, and tags removed",
			currentIdentityProvider: createCurrentIdentityProvider(idnetityProviderName, identityProviderARN, eks.ConfigStatusActive, createTags()),
			desiredIdentityProvider: createDesiredIdentityProvider(idnetityProviderName, nil),
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.UntagResource(gomock.Eq(&eks.UntagResourceInput{
					ResourceArn: aws.String(identityProviderARN),
					TagKeys:     []*string{aws.String("key1")},
				}))
			},
			expectCreateError: false,
			expectDoError:     false,
		},

		{
			name:                    "1 installed and 0 desired - installed provider is removed",
			currentIdentityProvider: createCurrentIdentityProvider(idnetityProviderName, identityProviderARN, eks.ConfigStatusActive, createTags()),
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.DisassociateIdentityProviderConfigWithContext(gomock.Eq(context.TODO()),
					gomock.Eq(&eks.DisassociateIdentityProviderConfigInput{
						ClusterName: aws.String(clusterName),
						IdentityProviderConfig: &eks.IdentityProviderConfig{
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
			currentIdentityProvider: createCurrentIdentityProvider(idnetityProviderName, identityProviderARN, eks.ConfigStatusActive, createTags()),
			desiredIdentityProvider: createDesiredIdentityProviderWithDifferentClientID(idnetityProviderName, createTags()),
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.DisassociateIdentityProviderConfigWithContext(gomock.Eq(context.TODO()),
					gomock.Eq(&eks.DisassociateIdentityProviderConfigInput{
						ClusterName: aws.String(clusterName),
						IdentityProviderConfig: &eks.IdentityProviderConfig{
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
		Tags:                       tags,
	}
}

func createCurrentIdentityProvider(name string, arn, status string, tags infrav1.Tags) *OidcIdentityProviderConfig {
	config := createDesiredIdentityProvider(name, tags)
	config.IdentityProviderConfigArn = aws.String(arn)
	config.Status = aws.String(status)

	return config
}

func changeTags(original infrav1.Tags) infrav1.Tags {
	original["key2"] = "value2"
	return original
}

func createDesiredIdentityProviderRequest(name *string) *eks.OidcIdentityProviderConfigRequest {
	return &eks.OidcIdentityProviderConfigRequest{
		ClientId:                   aws.String("clientId"),
		IdentityProviderConfigName: name,
		IssuerUrl:                  aws.String("http://IssuerURL.com"),
	}
}

func createDesiredIdentityProviderWithDifferentClientID(name string, tags infrav1.Tags) *OidcIdentityProviderConfig {
	p := createDesiredIdentityProvider(name, tags)
	p.ClientID = "clientId2"
	return p
}
