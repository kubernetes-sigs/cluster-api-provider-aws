/*
Copyright The Kubernetes Authors.

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

package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	sdk "github.com/openshift-online/ocm-sdk-go"
	"github.com/openshift-online/ocm-sdk-go/logging"
	. "github.com/openshift-online/ocm-sdk-go/testing"
	rosaocm "github.com/openshift/rosa/pkg/ocm"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/s3/mock_stsiface"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/rosa"
	"sigs.k8s.io/cluster-api/util/conditions"
)

// createFakeOCMClient creates a fake OCM client that doesn't require real network connections
func createFakeOCMClient() (*rosaocm.Client, *ghttp.Server, *ghttp.Server, error) {
	// Create fake servers
	ssoServer := MakeTCPServer()
	apiServer := MakeTCPServer()
	apiServer.SetAllowUnhandledRequests(true)
	apiServer.SetUnhandledRequestStatusCode(http.StatusInternalServerError)

	// Create a fake access token
	accessToken := "faketoken"

	// Prepare the SSO server to respond with the access token
	ssoServer.AppendHandlers(
		ghttp.RespondWith(
			200,
			accessToken,
			http.Header{
				"Content-Type": []string{
					"application/json",
				},
			},
		),
	)

	// Create a logger for the OCM SDK
	logger, err := logging.NewGoLoggerBuilder().
		Debug(false).
		Build()
	if err != nil {
		return nil, nil, nil, err
	}

	// Create the fake connection
	connection, err := sdk.NewConnectionBuilder().
		Logger(logger).
		Tokens(accessToken).
		URL(apiServer.URL()).
		Build()
	if err != nil {
		return nil, nil, nil, err
	}

	// Create the OCM client with the fake connection
	ocmClient := rosaocm.NewClientWithConnection(connection)

	return ocmClient, ssoServer, apiServer, nil
}

// setupAPIServerHandlers adds the OCM API handlers to the server
func setupAPIServerHandlers(apiServer *ghttp.Server) {
	// Handler for GetPolicies - responds to /api/clusters_mgmt/v1/aws_sts_policies
	mockPoliciesResponse := `{
		"kind": "AWSSTSPolicyList",
		"page": 1,
		"size": 2,
		"total": 2,
		"items": [
			{
				"kind": "AWSSTSPolicy",
				"id": "account-role-policy",
				"type": "json",
				"details": "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Effect\":\"Allow\",\"Action\":\"*\",\"Resource\":\"*\"}]}"
			},
			{
				"kind": "AWSSTSPolicy", 
				"id": "operator-role-policy",
				"type": "json",
				"details": "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Effect\":\"Allow\",\"Action\":\"*\",\"Resource\":\"*\"}]}"
			}
		]
	}`

	// Handler for GetOidcConfig - responds to /api/clusters_mgmt/v1/oidc_configs/{id}
	mockOidcConfigResponse := `{
		"kind": "OidcConfig",
		"id": "test-oidc-config-id",
		"issuer_url": "https://test-oidc.s3.amazonaws.com",
		"managed": true,
		"secret_arn": "arn:aws:secretsmanager:us-east-1:123456789012:secret:test-secret"
	}`

	// Add handlers to the API server
	apiServer.AppendHandlers(
		// Handle GetPolicies requests
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/api/clusters_mgmt/v1/aws_sts_policies"),
			ghttp.RespondWith(http.StatusOK, mockPoliciesResponse, http.Header{
				"Content-Type": []string{"application/json"},
			}),
		),
		// Handle GetOidcConfig requests
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/api/clusters_mgmt/v1/oidc_configs/test-oidc-config-id"),
			ghttp.RespondWith(http.StatusOK, mockOidcConfigResponse, http.Header{
				"Content-Type": []string{"application/json"},
			}),
		),
		// Handle OIDC config creation
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("POST", "/api/clusters_mgmt/v1/oidc_configs"),
			ghttp.RespondWith(http.StatusCreated, mockOidcConfigResponse, http.Header{
				"Content-Type": []string{"application/json"},
			}),
		),
		// Handle clusters list requests (for deletion checks)
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/api/clusters_mgmt/v1/clusters"),
			ghttp.RespondWith(http.StatusOK, `{
				"kind": "ClusterList",
				"page": 1,
				"size": 0,
				"total": 0,
				"items": []
			}`, http.Header{
				"Content-Type": []string{"application/json"},
			}),
		),
		// Handle credential requests
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/api/clusters_mgmt/v1/aws_sts_operators"),
			ghttp.RespondWith(http.StatusOK, `{
				"kind": "STSOperatorList",
				"page": 1,
				"size": 0,
				"total": 0,
				"items": []
			}`, http.Header{
				"Content-Type": []string{"application/json"},
			}),
		),
	)
}

func TestROSARoleConfigReconciler_Reconcile(t *testing.T) {
	g := NewWithT(t)
	ns, err := testEnv.CreateNamespace(ctx, "test-namespace")
	g.Expect(err).ToNot(HaveOccurred())

	// Create test secret for OCM credentials
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "rosa-secret",
			Namespace: ns.Name,
		},
		Data: map[string][]byte{
			"ocmToken":  []byte("test-ocm-token"),
			"ocmApiUrl": []byte("https://api.openshift.com"),
		},
	}

	// Create test identity
	identity := &infrav1.AWSClusterControllerIdentity{
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
		},
		Spec: infrav1.AWSClusterControllerIdentitySpec{
			AWSClusterIdentitySpec: infrav1.AWSClusterIdentitySpec{
				AllowedNamespaces: &infrav1.AllowedNamespaces{},
			},
		},
	}
	identity.SetGroupVersionKind(infrav1.GroupVersion.WithKind("AWSClusterControllerIdentity"))

	// Create objects in the test environment
	g.Expect(testEnv.Create(ctx, secret)).To(Succeed())
	defer func() { g.Expect(testEnv.Delete(ctx, secret)).To(Succeed()) }()
	g.Expect(testEnv.Create(ctx, identity)).To(Succeed())
	defer func() { g.Expect(testEnv.Delete(ctx, identity)).To(Succeed()) }()

	tests := []struct {
		name               string
		rosaRoleConfig     *expinfrav1.ROSARoleConfig
		expectError        bool
		expectRequeue      bool
		expectedConditions func(*expinfrav1.ROSARoleConfig)
	}{
		{
			name: "successful reconciliation - happy path with API server handlers",
			rosaRoleConfig: &expinfrav1.ROSARoleConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-rosa-role-config",
					Namespace: ns.Name,
				},
				Spec: expinfrav1.ROSARoleConfigSpec{
					Region: "us-east-1",
					AccountRoleConfig: expinfrav1.AccountRoleConfig{
						Prefix:  "test",
						Version: "4.15.0",
					},
					OperatorRoleConfig: expinfrav1.OperatorRoleConfig{
						Prefix: "test",
					},
					OIDCConfig: expinfrav1.OIDCConfig{
						ManagedOIDC: true,
						Prefix:      "test",
						Region:      "us-east-1",
					},
					IdentityRef: &infrav1.AWSIdentityReference{
						Name: identity.Name,
						Kind: infrav1.ControllerIdentityKind,
					},
					CredentialsSecretRef: &corev1.LocalObjectReference{
						Name: secret.Name,
					},
				},
			},
			expectError:   false, // Now we expect success with API server handlers
			expectRequeue: false,
			expectedConditions: func(roleConfig *expinfrav1.ROSARoleConfig) {
				// Expect finalizer to be added
				g.Expect(roleConfig.GetFinalizers()).To(ContainElement(expinfrav1.RosaRoleConfigFinalizer))
				// Should have ready condition set (API server will provide proper responses)
				readyCondition := conditions.Get(roleConfig, expinfrav1.RosaRoleConfigReadyCondition)
				g.Expect(readyCondition).ToNot(BeNil())
				// The condition status will depend on how far the reconciliation gets with mocked AWS calls
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := NewWithT(t)

			// Create fake OCM client with API server handlers
			ocmClient, ssoServer, apiServer, err := createFakeOCMClient()
			g.Expect(err).ToNot(HaveOccurred())
			defer ssoServer.Close()
			defer apiServer.Close()
			defer ocmClient.Close()

			// Setup API server handlers after client creation
			setupAPIServerHandlers(apiServer)

			g.Expect(testEnv.Create(ctx, test.rosaRoleConfig)).To(Succeed())
			defer func() { g.Expect(testEnv.Delete(ctx, test.rosaRoleConfig)).To(Succeed()) }()

			// Setup mocks
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			stsMock := mock_stsiface.NewMockSTSAPI(mockCtrl)
			stsMock.EXPECT().GetCallerIdentity(gomock.Any()).Return(nil, nil).AnyTimes()

			// Create reconciler with a custom NewOCMClient function that returns our fake client
			reconciler := &ROSARoleConfigReconciler{
				Client:           testEnv.Client,
				WatchFilterValue: "",
				Endpoints:        []scope.ServiceEndpoint{},
				NewStsClient: func(cloud.ScopeUsage, cloud.Session, logger.Wrapper, runtime.Object) stsiface.STSAPI {
					return stsMock
				},
				// Use a custom NewOCMClient that returns our fake client with API server handlers
				NewOCMClient: func(ctx context.Context, scope rosa.OCMSecretsRetriever) (rosa.OCMClient, error) {
					// Return a wrapper around our fake OCM client that has real API server responses
					return rosa.NewWrappedOCMClientFromOCMClient(ctx, ocmClient)
				},
			}

			// Create reconcile request
			req := ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      test.rosaRoleConfig.Name,
					Namespace: test.rosaRoleConfig.Namespace,
				},
			}

			// Perform reconciliation
			result, err := reconciler.Reconcile(ctx, req)

			// Assertions
			if test.expectError {
				g.Expect(err).To(HaveOccurred())
				if test.expectRequeue {
					g.Expect(result.RequeueAfter).To(Equal(time.Second * 60))
				}
			} else {
				g.Expect(err).ToNot(HaveOccurred())
				if test.expectRequeue {
					g.Expect(result.RequeueAfter).To(BeNumerically(">", 0))
				} else {
					g.Expect(result).To(Equal(ctrl.Result{}))
				}
			}

			// For the successful object creation case, verify conditions and finalizers
			updatedRoleConfig := &expinfrav1.ROSARoleConfig{}
			key := client.ObjectKey{Name: test.rosaRoleConfig.Name, Namespace: test.rosaRoleConfig.Namespace}
			g.Expect(testEnv.Get(ctx, key, updatedRoleConfig)).To(Succeed())

			// Apply expected conditions validation
			test.expectedConditions(updatedRoleConfig)
		})
	}
}

// Simple test to verify API server handlers work correctly
func TestAPIServerHandlers(t *testing.T) {
	g := NewWithT(t)

	// Create fake OCM client with API server handlers
	ocmClient, ssoServer, apiServer, err := createFakeOCMClient()
	g.Expect(err).ToNot(HaveOccurred())
	defer ssoServer.Close()
	defer apiServer.Close()
	defer ocmClient.Close()

	// Setup API server handlers after client creation
	setupAPIServerHandlers(apiServer)

	t.Run("GetOidcConfig returns correct response", func(t *testing.T) {
		g := NewWithT(t)

		// Test GetOidcConfig directly
		oidcConfig, err := ocmClient.GetOidcConfig("test-oidc-config-id")
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(oidcConfig).ToNot(BeNil())
		g.Expect(oidcConfig.ID()).To(Equal("test-oidc-config-id"))
		g.Expect(oidcConfig.IssuerUrl()).To(Equal("https://test-oidc.s3.amazonaws.com"))
		g.Expect(oidcConfig.Managed()).To(BeTrue())
	})

	t.Run("GetPolicies returns correct response", func(t *testing.T) {
		g := NewWithT(t)

		// Test GetPolicies directly
		policies, err := ocmClient.GetPolicies("AccountRole")
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(policies).ToNot(BeEmpty())
		g.Expect(len(policies)).To(BeNumerically(">=", 1))

		// Check that we have the expected policies
		var foundAccountPolicy bool
		for _, policy := range policies {
			if policy.ID() == "account-role-policy" {
				foundAccountPolicy = true
				g.Expect(policy.Type()).To(Equal("json"))
			}
		}
		g.Expect(foundAccountPolicy).To(BeTrue(), "Should find account-role-policy")
	})

	t.Run("API server receives correct requests", func(t *testing.T) {
		g := NewWithT(t)

		// Make a direct HTTP request to verify the server is responding
		client := &http.Client{}

		// Test OIDC config endpoint
		fmt.Println(apiServer.URL())
		req, err := http.NewRequest("GET", apiServer.URL()+"/api/clusters_mgmt/v1/oidc_configs/test-oidc-config-id", nil)
		g.Expect(err).ToNot(HaveOccurred())
		req.Header.Set("Authorization", "Bearer faketoken")

		resp, err := client.Do(req)
		g.Expect(err).ToNot(HaveOccurred())
		defer resp.Body.Close()

		fmt.Printf("Response status: %d\n", resp.StatusCode)

		g.Expect(resp.StatusCode).To(Equal(http.StatusOK))

		// Read and verify response body
		body, err := io.ReadAll(resp.Body)
		g.Expect(err).ToNot(HaveOccurred())

		fmt.Printf("Response body: %s\n", string(body))

		var oidcResponse map[string]interface{}
		err = json.Unmarshal(body, &oidcResponse)
		g.Expect(err).ToNot(HaveOccurred())

		g.Expect(oidcResponse["id"]).To(Equal("test-oidc-config-id"))
		g.Expect(oidcResponse["issuer_url"]).To(Equal("https://test-oidc.s3.amazonaws.com"))
		g.Expect(oidcResponse["managed"]).To(BeTrue())
	})
}
