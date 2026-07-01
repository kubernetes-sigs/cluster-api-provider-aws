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

package controllers

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	awsSdk "github.com/aws/aws-sdk-go-v2/aws"
	iamv2 "github.com/aws/aws-sdk-go-v2/service/iam"
	iamTypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	stsv2 "github.com/aws/aws-sdk-go-v2/service/sts"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	sdk "github.com/openshift-online/ocm-sdk-go"
	ocmlogging "github.com/openshift-online/ocm-sdk-go/logging"
	ocmsdk "github.com/openshift-online/ocm-sdk-go/testing"
	"github.com/openshift/rosa/pkg/aws"
	rosaMocks "github.com/openshift/rosa/pkg/aws/mocks"
	"github.com/openshift/rosa/pkg/ocm"
	rosacli "github.com/openshift/rosa/pkg/rosa"
	"github.com/sirupsen/logrus"
	"go.uber.org/mock/gomock"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"

	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	caparosa "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/rosa"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"
)

// getUniqueOrgID creates a unique organization ID for test isolation.
func getUniqueOrgID(testID string) string {
	return fmt.Sprintf("test-org-%s", testID)
}

// getUniqueAccountID creates a unique AWS account ID for test isolation.
func getUniqueAccountID(testID string) string {
	return "1" + testID[:11]
}

// getUniqueExternalID creates a unique external ID for test isolation.
func getUniqueExternalID(testID string) string {
	return getUniqueAccountID(testID)[4:12]
}

// common test infrastructure.
type testHarness struct {
	t             *testing.T
	g             *WithT
	testID        string
	orgID         string
	accountID     string
	externalID    string
	ssoServer     *ghttp.Server
	apiServer     *ghttp.Server
	ocmClient     *ocm.Client
	mockCtrl      *gomock.Controller
	testCtx       context.Context //nolint:containedctx // test helper struct
	ocmRoleConfig *expinfrav1.ROSAOCMRoleConfig
}

func newTestHarness(t *testing.T) *testHarness {
	t.Helper()
	RegisterTestingT(t)
	g := NewWithT(t)

	testID := generateTestID()
	orgID := getUniqueOrgID(testID)
	accountID := getUniqueAccountID(testID)
	externalID := getUniqueExternalID(testID)

	ssoServer := ocmsdk.MakeTCPServer()
	apiServer := ocmsdk.MakeTCPServer()
	apiServer.SetAllowUnhandledRequests(true)
	apiServer.SetUnhandledRequestStatusCode(http.StatusInternalServerError)

	accessToken := ocmsdk.MakeTokenString("Bearer", 15*time.Minute)
	ssoServer.AppendHandlers(ocmsdk.RespondWithAccessToken(accessToken))

	logger, err := ocmlogging.NewGoLoggerBuilder().Debug(false).Build()
	Expect(err).ToNot(HaveOccurred())
	connection, err := sdk.NewConnectionBuilder().
		Logger(logger).
		Tokens(accessToken).
		URL(apiServer.URL()).
		Build()
	Expect(err).To(BeNil())
	ocmClient := ocm.NewClientWithConnection(connection)

	mockCtrl := gomock.NewController(t)

	h := &testHarness{
		t:          t,
		g:          g,
		testID:     testID,
		orgID:      orgID,
		accountID:  accountID,
		externalID: externalID,
		ssoServer:  ssoServer,
		apiServer:  apiServer,
		ocmClient:  ocmClient,
		mockCtrl:   mockCtrl,
		testCtx:    context.TODO(),
	}

	t.Cleanup(func() {
		ssoServer.Close()
		apiServer.Close()
		mockCtrl.Finish()
	})

	return h
}

// setupOCMCurrentAccount adds the standard current_account handler.
func (h *testHarness) setupOCMCurrentAccount() {
	h.apiServer.RouteToHandler("GET", "/api/accounts_mgmt/v1/current_account",
		ocmsdk.RespondWithJSON(
			http.StatusOK, fmt.Sprintf(`{"id": %q, "organization": {"id": %q, "external_id": %q}}`,
				h.orgID, h.orgID, h.externalID),
		),
	)
}

// setupOCMLinkCheck adds a handler for checking if role is linked.
func (h *testHarness) setupOCMLinkCheck(linked bool) {
	h.apiServer.RouteToHandler("GET", fmt.Sprintf("/api/accounts_mgmt/v1/organizations/%s/labels/sts_ocm_role", h.orgID),
		func(w http.ResponseWriter, r *http.Request) {
			if linked {
				ocmsdk.RespondWithJSON(http.StatusOK,
					fmt.Sprintf(`{"value": "arn:aws:iam::%s:role/test-OCM-Role-%s"}`, h.accountID, h.externalID))(w, r)
			} else {
				ocmsdk.RespondWithJSON(http.StatusNotFound, `{}`)(w, r)
			}
		},
	)
}

// setupOCMLinkCheckProgressive adds a handler that returns NotFound first, then linked.
func (h *testHarness) setupOCMLinkCheckProgressive() {
	callCount := 0
	h.apiServer.RouteToHandler("GET", fmt.Sprintf("/api/accounts_mgmt/v1/organizations/%s/labels/sts_ocm_role", h.orgID),
		func(w http.ResponseWriter, r *http.Request) {
			callCount++
			if callCount == 1 {
				ocmsdk.RespondWithJSON(http.StatusNotFound, `{}`)(w, r)
			} else {
				ocmsdk.RespondWithJSON(http.StatusOK,
					fmt.Sprintf(`{"value": "arn:aws:iam::%s:role/test-OCM-Role-%s"}`, h.accountID, h.externalID))(w, r)
			}
		},
	)
}

// setupOCMLinkCheckDifferentRole adds a handler that returns a different role ARN (conflict scenario).
func (h *testHarness) setupOCMLinkCheckDifferentRole(differentRoleARN string) {
	h.apiServer.RouteToHandler("GET", fmt.Sprintf("/api/accounts_mgmt/v1/organizations/%s/labels/sts_ocm_role", h.orgID),
		ocmsdk.RespondWithJSON(http.StatusOK, fmt.Sprintf(`{"value": %q}`, differentRoleARN)),
	)
}

// setupOCMLinkPost adds a handler for linking a role.
func (h *testHarness) setupOCMLinkPost() {
	h.apiServer.RouteToHandler("POST", fmt.Sprintf("/api/accounts_mgmt/v1/organizations/%s/labels", h.orgID),
		ocmsdk.RespondWithJSON(
			http.StatusCreated,
			fmt.Sprintf(`{"key": "sts_ocm_role", "value": "arn:aws:iam::%s:role/test-OCM-Role-%s"}`, h.accountID, h.externalID),
		),
	)
}

// createOCMRoleConfig creates a test ROSAOCMRoleConfig resource.
func (h *testHarness) createOCMRoleConfig(name string, profile expinfrav1.ROSAOCMRoleProfile) {
	h.ocmRoleConfig = &expinfrav1.ROSAOCMRoleConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-%s", name, h.testID),
		},
		Spec: expinfrav1.ROSAOCMRoleConfigSpec{
			RolePrefix: "test",
			Profile:    profile,
		},
	}
	createObject(h.g, h.ocmRoleConfig, "")
	h.t.Cleanup(func() { cleanupObject(h.g, h.ocmRoleConfig) })
}

// createOCMRoleConfigWithFinalizer creates a test ROSAOCMRoleConfig resource with finalizer (for delete tests).
func (h *testHarness) createOCMRoleConfigWithFinalizer(name string, profile expinfrav1.ROSAOCMRoleProfile) {
	h.ocmRoleConfig = &expinfrav1.ROSAOCMRoleConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:       fmt.Sprintf("%s-%s", name, h.testID),
			Finalizers: []string{expinfrav1.ROSAOCMRoleConfigFinalizer},
		},
		Spec: expinfrav1.ROSAOCMRoleConfigSpec{
			RolePrefix: "test",
			Profile:    profile,
		},
	}
	createObject(h.g, h.ocmRoleConfig, "")
	h.t.Cleanup(func() { cleanupObject(h.g, h.ocmRoleConfig) })
}

// setOCMRoleConfigStatus sets the status on the ROSAOCMRoleConfig resource.
func (h *testHarness) setOCMRoleConfigStatus(roleARN, orgID string) {
	h.ocmRoleConfig.Status = expinfrav1.ROSAOCMRoleConfigStatus{
		RoleARN:        roleARN,
		OrganizationID: orgID,
	}
	err := testEnv.Client.Status().Update(h.testCtx, h.ocmRoleConfig)
	h.g.Expect(err).ToNot(HaveOccurred())
}

// setupOCMUnlinkDelete adds a handler for DELETE (unlink) endpoint with callback tracking.
func (h *testHarness) setupOCMUnlinkDelete(callbackPtr *bool) {
	h.apiServer.RouteToHandler("DELETE", fmt.Sprintf("/api/accounts_mgmt/v1/organizations/%s/labels/sts_ocm_role", h.orgID),
		func(w http.ResponseWriter, r *http.Request) {
			*callbackPtr = true
			ocmsdk.RespondWithJSON(http.StatusNoContent, `{}`)(w, r)
		},
	)
}

// buildReconciler creates a reconciler with the given AWS client.
func (h *testHarness) buildReconciler(awsClient aws.Client) *ROSAOCMRoleConfigReconciler {
	return &ROSAOCMRoleConfigReconciler{
		Client:   testEnv.Client,
		Recorder: record.NewFakeRecorder(32),
		runtimeFactory: func(ctx context.Context, scope *scope.ROSAOCMRoleConfigScope) (*rosacli.Runtime, error) {
			r := rosacli.NewRuntime()
			r.OCMClient = h.ocmClient
			r.AWSClient = awsClient
			r.Creator = &aws.Creator{
				AccountID: h.accountID,
				ARN:       fmt.Sprintf("arn:aws:iam::%s:user/test", h.accountID),
				Partition: "aws",
			}
			return r, nil
		},
	}
}

// buildAWSClientWithSDKMocks creates a real aws.Client with mocked AWS SDK clients
// The setupIAM and setupSTS callbacks receive the mock clients to configure expectations.
func (h *testHarness) buildAWSClientWithSDKMocks(
	setupIAM func(*rosaMocks.MockIamApiClient),
	setupSTS func(*rosaMocks.MockStsApiClient),
) aws.Client {
	mockIamClient := rosaMocks.NewMockIamApiClient(h.mockCtrl)
	mockSTSClient := rosaMocks.NewMockStsApiClient(h.mockCtrl)

	setupIAM(mockIamClient)
	setupSTS(mockSTSClient)

	return aws.New(
		awsSdk.Config{},
		aws.NewLoggerWrapper(logrus.New(), nil),
		mockIamClient,
		rosaMocks.NewMockEc2ApiClient(h.mockCtrl),
		rosaMocks.NewMockOrganizationsApiClient(h.mockCtrl),
		rosaMocks.NewMockS3ApiClient(h.mockCtrl),
		rosaMocks.NewMockSecretsManagerApiClient(h.mockCtrl),
		mockSTSClient,
		rosaMocks.NewMockCloudFormationApiClient(h.mockCtrl),
		rosaMocks.NewMockServiceQuotasApiClient(h.mockCtrl),
		rosaMocks.NewMockServiceQuotasApiClient(h.mockCtrl),
		&aws.AccessKey{},
		false,
	)
}

// setupOCMPoliciesHandler adds the OCM policies endpoint handler.
func (h *testHarness) setupOCMPoliciesHandler() {
	h.apiServer.RouteToHandler("GET", "/api/clusters_mgmt/v1/aws_inquiries/sts_policies",
		func(w http.ResponseWriter, r *http.Request) {
			query := r.URL.Query().Get("search")
			if query == "type='OCMRole'" {
				ocmsdk.RespondWithJSON(http.StatusOK, `{
                    "items": [
                        {
                            "id": "sts_ocm_permission_policy",
                            "arn": "arn:aws:iam::aws:policy/sts_ocm_permission_policy",
                            "type": "OCMRole",
                            "details": "{\"Version\": \"2012-10-17\", \"Statement\": [{\"Effect\": \"Allow\", \"Action\": \"*\", \"Resource\": \"*\"}]}"
                        },
                        {
                            "id": "sts_ocm_trust_policy",
                            "arn": "",
                            "type": "OCMRole",
                            "details": "{\"Version\": \"2012-10-17\", \"Statement\": [{\"Effect\": \"Allow\", \"Principal\": {\"AWS\": \"arn:aws:iam::710019948333:role/RH-Managed-OpenShift-Installer\"}, \"Action\": \"sts:AssumeRole\", \"Condition\": {\"StringEquals\": {\"sts:ExternalId\": \"${ocm_organization_id}\"}}}]}"
                        }
                    ]
                }`)(w, r)
			} else {
				ocmsdk.RespondWithJSON(http.StatusOK, `{"items": []}`)(w, r)
			}
		},
	)
}

// setupOCMPoliciesHandlerNoConsole adds the OCM policies endpoint handler for no-console tests.
func (h *testHarness) setupOCMPoliciesHandlerNoConsole() {
	h.apiServer.RouteToHandler("GET", "/api/clusters_mgmt/v1/aws_inquiries/sts_policies",
		func(w http.ResponseWriter, r *http.Request) {
			query := r.URL.Query().Get("search")
			if query == "policy_type = 'OCMRole'" || query == "type='OCMRole'" {
				ocmsdk.RespondWithJSON(http.StatusOK, `{
                    "items": [
                        {
                            "id": "sts_ocm_no_console_permission_policy",
                            "type": "OCMRole",
                            "details": "{\"Version\": \"2012-10-17\", \"Statement\": [{\"Effect\": \"Allow\", \"Action\": \"sts:AssumeRole\", \"Resource\": \"*\"}]}"
                        },
                        {
                            "id": "sts_ocm_trust_policy",
                            "type": "OCMRole",
                            "details": "{\"Version\": \"2012-10-17\", \"Statement\": []}"
                        }
                    ]
                }`)(w, r)
			} else {
				ocmsdk.RespondWithJSON(http.StatusOK, `{"items": []}`)(w, r)
			}
		},
	)
}

func TestROSAOCMRoleConfigReconcileCreate(t *testing.T) {
	h := newTestHarness(t)

	// Setup OCM mocks
	h.setupOCMCurrentAccount()
	h.setupOCMPoliciesHandler()
	h.setupOCMLinkCheckProgressive() // Not linked initially, then linked after creation
	h.setupOCMLinkPost()

	// Setup AWS mocks for role creation flow
	awsClient := h.buildAWSClientWithSDKMocks(
		func(mockIamClient *rosaMocks.MockIamApiClient) {
			// Role doesn't exist yet - trigger creation
			mockIamClient.EXPECT().GetRole(gomock.Any(), gomock.Any()).Return(nil, &iamTypes.NoSuchEntityException{
				Message: awsSdk.String("The role with name test-OCM-Role-123 does not exist."),
			}).AnyTimes()

			// Mock role creation
			mockIamClient.EXPECT().CreateRole(gomock.Any(), gomock.Any()).Return(&iamv2.CreateRoleOutput{
				Role: &iamTypes.Role{
					RoleName: awsSdk.String(fmt.Sprintf("test-OCM-Role-%s", h.externalID)),
					Arn:      awsSdk.String(fmt.Sprintf("arn:aws:iam::%s:role/test-OCM-Role-%s", h.accountID, h.externalID)),
				},
			}, nil).AnyTimes()

			// Mock policy operations
			mockIamClient.EXPECT().GetPolicy(gomock.Any(), gomock.Any()).Return(nil, &iamTypes.NoSuchEntityException{
				Message: awsSdk.String("The policy does not exist."),
			}).AnyTimes()

			mockIamClient.EXPECT().CreatePolicy(gomock.Any(), gomock.Any()).Return(&iamv2.CreatePolicyOutput{
				Policy: &iamTypes.Policy{
					PolicyName: awsSdk.String(fmt.Sprintf("test-OCM-Role-%s-Policy", h.externalID)),
					Arn:        awsSdk.String(fmt.Sprintf("arn:aws:iam::%s:policy/test-OCM-Role-%s-Policy", h.accountID, h.externalID)),
				},
			}, nil).AnyTimes()

			mockIamClient.EXPECT().AttachRolePolicy(gomock.Any(), gomock.Any()).Return(&iamv2.AttachRolePolicyOutput{}, nil).AnyTimes()
			mockIamClient.EXPECT().TagRole(gomock.Any(), gomock.Any()).Return(&iamv2.TagRoleOutput{}, nil).AnyTimes()
		},
		func(mockSTSClient *rosaMocks.MockStsApiClient) {
			mockSTSClient.EXPECT().GetCallerIdentity(gomock.Any(), gomock.Any()).Return(&stsv2.GetCallerIdentityOutput{
				Arn:     awsSdk.String(fmt.Sprintf("arn:aws:iam::%s:user/test-user", h.accountID)),
				Account: awsSdk.String(h.accountID),
				UserId:  awsSdk.String("test-user-id"),
			}, nil).AnyTimes()
		},
	)

	// Create test resource and reconciler
	h.createOCMRoleConfig("test-ocm-role", expinfrav1.ROSAOCMRoleProfileStandard)
	reconciler := h.buildReconciler(awsClient)

	// Test: Role creation should succeed
	req := ctrl.Request{NamespacedName: types.NamespacedName{Name: h.ocmRoleConfig.Name}}
	h.g.Eventually(func(g Gomega) {
		_, err := reconciler.Reconcile(h.testCtx, req)
		g.Expect(err).ToNot(HaveOccurred())

		updatedRoleConfig := &expinfrav1.ROSAOCMRoleConfig{}
		g.Expect(reconciler.Client.Get(h.testCtx, req.NamespacedName, updatedRoleConfig)).ToNot(HaveOccurred())

		g.Expect(updatedRoleConfig.Status.RoleARN).To(Equal(fmt.Sprintf("arn:aws:iam::%s:role/test-OCM-Role-%s", h.accountID, h.externalID)))
		g.Expect(updatedRoleConfig.Status.OrganizationID).To(Equal(h.orgID))

		readyCondition := v1beta1conditions.Get(updatedRoleConfig, expinfrav1.ROSAOCMRoleConfigReadyCondition)
		g.Expect(readyCondition).ToNot(BeNil())
		g.Expect(readyCondition.Status).To(Equal(corev1.ConditionTrue))
	}).WithTimeout(30 * time.Second).Should(Succeed())
}

func TestROSAOCMRoleConfigReconcileExist(t *testing.T) {
	h := newTestHarness(t)

	// Setup OCM mocks
	h.setupOCMCurrentAccount()
	h.setupOCMLinkCheck(true) // Role already linked

	// Setup AWS mocks - role already exists (standard profile)
	mockAWSClient := aws.NewMockClient(h.mockCtrl)
	mockAWSClient.EXPECT().CheckRoleExists(fmt.Sprintf("test-OCM-Role-%s", h.externalID)).Return(
		true,
		fmt.Sprintf("arn:aws:iam::%s:role/test-OCM-Role-%s", h.accountID, h.externalID),
		nil,
	).AnyTimes()
	mockAWSClient.EXPECT().IsAdminRole(fmt.Sprintf("test-OCM-Role-%s", h.externalID)).Return(false, nil).AnyTimes()
	mockAWSClient.EXPECT().GetRoleByName(fmt.Sprintf("test-OCM-Role-%s", h.externalID)).Return(iamTypes.Role{}, nil).AnyTimes()
	mockAWSClient.EXPECT().GetCreator().Return(&aws.Creator{
		ARN:       fmt.Sprintf("arn:aws:iam::%s:user/test-user", h.accountID),
		AccountID: h.accountID,
		IsSTS:     false,
		Partition: "aws",
	}, nil).AnyTimes()

	// Create test resource and reconciler
	h.createOCMRoleConfig("test-ocm-role-exist", expinfrav1.ROSAOCMRoleProfileStandard)
	reconciler := h.buildReconciler(mockAWSClient)

	// Test: Existing role should be preserved
	req := ctrl.Request{NamespacedName: types.NamespacedName{Name: h.ocmRoleConfig.Name}}
	h.g.Eventually(func(g Gomega) {
		_, err := reconciler.Reconcile(h.testCtx, req)
		g.Expect(err).ToNot(HaveOccurred())

		updatedRoleConfig := &expinfrav1.ROSAOCMRoleConfig{}
		g.Expect(reconciler.Client.Get(h.testCtx, req.NamespacedName, updatedRoleConfig)).ToNot(HaveOccurred())

		g.Expect(updatedRoleConfig.Status.RoleARN).To(Equal(fmt.Sprintf("arn:aws:iam::%s:role/test-OCM-Role-%s", h.accountID, h.externalID)))
		g.Expect(updatedRoleConfig.Status.OrganizationID).To(Equal(h.orgID))

		readyCondition := v1beta1conditions.Get(updatedRoleConfig, expinfrav1.ROSAOCMRoleConfigReadyCondition)
		g.Expect(readyCondition).ToNot(BeNil())
		g.Expect(readyCondition.Status).To(Equal(corev1.ConditionTrue))
	}).WithTimeout(30 * time.Second).WithPolling(500 * time.Millisecond).Should(Succeed())
}

func TestROSAOCMRoleConfigSetUpRuntimeWithExpiredAWSCredentials(t *testing.T) {
	RegisterTestingT(t)
	g := NewWithT(t)
	testCtx := context.TODO()

	// Mock empty AWS credentials
	t.Setenv("AWS_ACCESS_KEY_ID", "")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "")
	t.Setenv("AWS_SESSION_TOKEN", "")
	t.Setenv("AWS_PROFILE", "")
	t.Setenv("AWS_SHARED_CREDENTIALS_FILE", "/dev/null")
	t.Setenv("AWS_CONFIG_FILE", "/dev/null")

	ocmRoleConfig := &expinfrav1.ROSAOCMRoleConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-role-config",
		},
		Spec: expinfrav1.ROSAOCMRoleConfigSpec{
			RolePrefix: "test",
			Profile:    expinfrav1.ROSAOCMRoleProfileStandard,
		},
	}
	createObject(g, ocmRoleConfig, "")
	defer cleanupObject(g, ocmRoleConfig)

	scope, err := scope.NewROSAOCMRoleConfigScope(scope.ROSAOCMRoleConfigScopeParams{
		Client:            testEnv.Client,
		ROSAOCMRoleConfig: ocmRoleConfig,
		ControllerName:    "rosaocmroleconfig",
	})
	g.Expect(err).ToNot(HaveOccurred())

	ssoServer := ocmsdk.MakeTCPServer()
	apiServer := ocmsdk.MakeTCPServer()
	defer ssoServer.Close()
	defer apiServer.Close()

	accessToken := ocmsdk.MakeTokenString("Bearer", 15*time.Minute)
	ssoServer.AppendHandlers(ocmsdk.RespondWithAccessToken(accessToken))

	logger, err := ocmlogging.NewGoLoggerBuilder().Debug(false).Build()
	Expect(err).ToNot(HaveOccurred())
	connection, err := sdk.NewConnectionBuilder().
		Logger(logger).
		Tokens(accessToken).
		URL(apiServer.URL()).
		Build()
	Expect(err).To(BeNil())

	ocmClient := ocm.NewClientWithConnection(connection)

	// Track NewOCMClient call count to verify retry behavior
	ocmClientCallCount := 0

	reconciler := &ROSAOCMRoleConfigReconciler{
		Client: testEnv.Client,
		NewOCMClient: func(ctx context.Context, scope caparosa.OCMSecretsRetriever) (caparosa.OCMClient, error) {
			ocmClientCallCount++
			return ocmClient, nil
		},
		runtimeFactory: nil, // Use default setUpRuntime (not the mock factory)
	}

	// ==================== FIRST RECONCILIATION ====================
	t.Log("First reconciliation: AWS credentials are missing/expired")

	_, err = reconciler.setUpRuntime(testCtx, scope)

	g.Expect(err).To(HaveOccurred(),
		"setUpRuntime should return error when AWS client creation fails")
	g.Expect(err.Error()).To(ContainSubstring("failed to create AWS client"),
		"Error should indicate AWS client failure")

	g.Expect(ocmClientCallCount).To(Equal(1),
		"NewOCMClient should be called once on first attempt")

	// ==================== SECOND RECONCILIATION ====================
	t.Log("Second reconciliation: Retry with still-missing AWS credentials")

	_, err = reconciler.setUpRuntime(testCtx, scope)

	g.Expect(err).To(HaveOccurred(),
		"setUpRuntime should still fail with missing AWS credentials")

	// NewOCMClient was called AGAIN (proves no early return)
	g.Expect(ocmClientCallCount).To(Equal(2),
		"NewOCMClient should be called twice - proves guard clause allows retry when Runtime is nil")
}

// TestSetUpRuntimeIdempotency verifies that setUpRuntime with runtimeFactory
// creates a new runtime per reconciliation (because identity may have changed).
func TestROSAOCMRoleConfigSetUpRuntimeIdempotency(t *testing.T) {
	h := newTestHarness(t)

	// Build AWS client with no expectations (not actually used in this test)
	awsClient := h.buildAWSClientWithSDKMocks(
		func(mockIamClient *rosaMocks.MockIamApiClient) {
		},
		func(mockSTSClient *rosaMocks.MockStsApiClient) {
		},
	)

	// Track factory calls to verify idempotency behavior
	callCount := 0
	reconciler := &ROSAOCMRoleConfigReconciler{
		runtimeFactory: func(ctx context.Context, scope *scope.ROSAOCMRoleConfigScope) (*rosacli.Runtime, error) {
			callCount++
			runtime := rosacli.NewRuntime()
			runtime.OCMClient = h.ocmClient
			runtime.AWSClient = awsClient
			runtime.Creator = &aws.Creator{
				ARN:       fmt.Sprintf("arn:aws:iam::%s:user/test", h.accountID),
				AccountID: h.accountID,
				IsSTS:     false,
			}
			return runtime, nil
		},
	}

	scope := &scope.ROSAOCMRoleConfigScope{}

	// Test: First call should invoke factory
	runtime1, err := reconciler.setUpRuntime(h.testCtx, scope)
	h.g.Expect(err).ToNot(HaveOccurred())
	h.g.Expect(callCount).To(Equal(1), "Factory should be called on first setUpRuntime")

	// Test: Second call should invoke factory again (new runtime per reconciliation)
	runtime2, err := reconciler.setUpRuntime(h.testCtx, scope)
	h.g.Expect(err).ToNot(HaveOccurred())
	h.g.Expect(callCount).To(Equal(2), "Factory should be called on second setUpRuntime")

	// Test: Should be different instances
	h.g.Expect(runtime1).ToNot(BeIdenticalTo(runtime2),
		"Runtime should be recreated each reconciliation to use correct identity")
}

// TestROSAOCMRoleConfigRoleConflict tests when a different OCM role ARN is already linked
// for the same AWS account (controller should fail with clear error).
func TestROSAOCMRoleConfigRoleConflict(t *testing.T) {
	h := newTestHarness(t)

	// Setup OCM mocks - different role is already linked (conflict)
	h.setupOCMCurrentAccount()
	h.setupOCMPoliciesHandler()
	differentRoleARN := fmt.Sprintf("arn:aws:iam::%s:role/different-OCM-Role-%s", h.accountID, h.externalID)
	h.setupOCMLinkCheckDifferentRole(differentRoleARN)

	// Setup AWS mocks - role exists
	mockAWSClient := aws.NewMockClient(h.mockCtrl)
	mockAWSClient.EXPECT().CheckRoleExists(fmt.Sprintf("test-OCM-Role-%s", h.externalID)).Return(
		true,
		fmt.Sprintf("arn:aws:iam::%s:role/test-OCM-Role-%s", h.accountID, h.externalID),
		nil,
	).AnyTimes()
	mockAWSClient.EXPECT().IsAdminRole(fmt.Sprintf("test-OCM-Role-%s", h.externalID)).Return(false, nil).AnyTimes()
	mockAWSClient.EXPECT().GetRoleByName(fmt.Sprintf("test-OCM-Role-%s", h.externalID)).Return(iamTypes.Role{}, nil).AnyTimes()
	mockAWSClient.EXPECT().ListAttachedRolePolicies(fmt.Sprintf("test-OCM-Role-%s", h.externalID)).Return(
		[]string{fmt.Sprintf("arn:aws:iam::%s:policy/test-OCM-Role-%s-Policy", h.accountID, h.externalID)},
		nil,
	).AnyTimes()
	mockAWSClient.EXPECT().GetCreator().Return(&aws.Creator{
		ARN:       fmt.Sprintf("arn:aws:iam::%s:user/test-user", h.accountID),
		AccountID: h.accountID,
		IsSTS:     false,
		Partition: "aws",
	}, nil).AnyTimes()

	// Create test resource and reconciler
	h.createOCMRoleConfig("test-ocm-role-conflict", expinfrav1.ROSAOCMRoleProfileStandard)
	reconciler := h.buildReconciler(mockAWSClient)

	// Test: Should fail with role conflict error
	req := ctrl.Request{NamespacedName: types.NamespacedName{Name: h.ocmRoleConfig.Name}}
	h.g.Eventually(func(g Gomega) {
		_, err := reconciler.Reconcile(h.testCtx, req)

		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("Only one role can be linked per AWS account per organization"))

		updatedRoleConfig := &expinfrav1.ROSAOCMRoleConfig{}
		g.Expect(reconciler.Client.Get(h.testCtx, req.NamespacedName, updatedRoleConfig)).ToNot(HaveOccurred())

		readyCondition := v1beta1conditions.Get(updatedRoleConfig, expinfrav1.ROSAOCMRoleConfigReadyCondition)
		g.Expect(readyCondition).ToNot(BeNil())
		g.Expect(readyCondition.Status).To(Equal(corev1.ConditionFalse))
		g.Expect(readyCondition.Message).To(ContainSubstring("Only one role can be linked per AWS account per organization"))
	}).WithTimeout(30 * time.Second).Should(Succeed())
}

// TestROSAOCMRoleConfigNoConsoleTagPolicyMismatch tests that reconciliation self-heals when role has no-console tag but wrong policy.
func TestROSAOCMRoleConfigNoConsoleTagPolicyMismatch(t *testing.T) {
	h := newTestHarness(t)

	// Setup OCM mocks
	h.setupOCMCurrentAccount()
	h.setupOCMPoliciesHandlerNoConsole()
	h.setupOCMLinkCheck(false)
	h.setupOCMLinkPost()

	// Setup AWS mocks - role has no-console tag but wrong policy attached (mismatch scenario)
	awsClient := h.buildAWSClientWithSDKMocks(
		func(mockIamClient *rosaMocks.MockIamApiClient) {
			// Role exists with no-console tag
			mockIamClient.EXPECT().GetRole(gomock.Any(), gomock.Any()).Return(&iamv2.GetRoleOutput{
				Role: &iamTypes.Role{
					RoleName: awsSdk.String(fmt.Sprintf("test-OCM-Role-%s", h.externalID)),
					Arn:      awsSdk.String(fmt.Sprintf("arn:aws:iam::%s:role/test-OCM-Role-%s", h.accountID, h.externalID)),
					Tags: []iamTypes.Tag{
						{
							Key:   awsSdk.String("rosa_no_console_role"),
							Value: awsSdk.String("true"),
						},
					},
				},
			}, nil).AnyTimes()

			// Policy exists but doesn't match no-console expectations
			mockIamClient.EXPECT().GetPolicy(gomock.Any(), gomock.Any()).Return(&iamv2.GetPolicyOutput{
				Policy: &iamTypes.Policy{
					Arn: awsSdk.String(fmt.Sprintf("arn:aws:iam::%s:policy/test-OCM-Role-%s-Policy", h.accountID, h.externalID)),
				},
			}, nil).AnyTimes()

			// Wrong policy attached (not the no-console policy)
			mockIamClient.EXPECT().ListAttachedRolePolicies(gomock.Any(), gomock.Any()).Return(&iamv2.ListAttachedRolePoliciesOutput{
				AttachedPolicies: []iamTypes.AttachedPolicy{
					{
						PolicyArn:  awsSdk.String(fmt.Sprintf("arn:aws:iam::%s:policy/test-OCM-Role-%s-Policy", h.accountID, h.externalID)),
						PolicyName: awsSdk.String(fmt.Sprintf("test-OCM-Role-%s-Policy", h.externalID)),
					},
				},
			}, nil).AnyTimes()

			// Self-healing will create and attach no-console policy
			mockIamClient.EXPECT().CreatePolicy(gomock.Any(), gomock.Any()).Return(&iamv2.CreatePolicyOutput{
				Policy: &iamTypes.Policy{
					Arn: awsSdk.String(fmt.Sprintf("arn:aws:iam::%s:policy/test-OCM-Role-%s-NoConsole-Policy", h.accountID, h.externalID)),
				},
			}, nil).AnyTimes()
			mockIamClient.EXPECT().AttachRolePolicy(gomock.Any(), gomock.Any()).Return(&iamv2.AttachRolePolicyOutput{}, nil).AnyTimes()
			mockIamClient.EXPECT().TagRole(gomock.Any(), gomock.Any()).Return(&iamv2.TagRoleOutput{}, nil).AnyTimes()
		},
		func(mockSTSClient *rosaMocks.MockStsApiClient) {
			mockSTSClient.EXPECT().GetCallerIdentity(gomock.Any(), gomock.Any()).Return(&stsv2.GetCallerIdentityOutput{
				Arn:     awsSdk.String(fmt.Sprintf("arn:aws:iam::%s:user/test-user", h.accountID)),
				Account: awsSdk.String(h.accountID),
				UserId:  awsSdk.String("test-user-id"),
			}, nil).AnyTimes()
		},
	)

	// Create test resource and reconciler
	h.createOCMRoleConfig("test-ocm-role-noconsole", expinfrav1.ROSAOCMRoleProfileNoConsole)
	reconciler := h.buildReconciler(awsClient)

	// Test: Should self-heal by attaching no-console policy
	req := ctrl.Request{NamespacedName: types.NamespacedName{Name: h.ocmRoleConfig.Name}}
	h.g.Eventually(func(g Gomega) {
		_, err := reconciler.Reconcile(h.testCtx, req)

		g.Expect(err).ToNot(HaveOccurred(), "reconciliation should succeed after self-healing")

		updatedRoleConfig := &expinfrav1.ROSAOCMRoleConfig{}
		g.Expect(reconciler.Client.Get(h.testCtx, req.NamespacedName, updatedRoleConfig)).ToNot(HaveOccurred())

		readyCondition := v1beta1conditions.Get(updatedRoleConfig, expinfrav1.ROSAOCMRoleConfigReadyCondition)
		g.Expect(readyCondition).ToNot(BeNil())
		g.Expect(readyCondition.Status).To(Equal(corev1.ConditionTrue), "should be ready after self-healing policy")
	}).WithTimeout(5 * time.Second).Should(Succeed())
}

// TestROSAOCMRoleConfigAdminProfile tests reconciling an existing role with Admin profile.
func TestROSAOCMRoleConfigAdminProfile(t *testing.T) {
	h := newTestHarness(t)

	// Setup OCM mocks
	h.setupOCMCurrentAccount()
	h.setupOCMLinkCheck(true) // Role already linked

	// Setup AWS mocks
	mockAWSClient := aws.NewMockClient(h.mockCtrl)
	mockAWSClient.EXPECT().CheckRoleExists(fmt.Sprintf("test-OCM-Role-%s", h.externalID)).Return(
		true,
		fmt.Sprintf("arn:aws:iam::%s:role/test-OCM-Role-%s", h.accountID, h.externalID),
		nil,
	).AnyTimes()
	mockAWSClient.EXPECT().IsAdminRole(fmt.Sprintf("test-OCM-Role-%s", h.externalID)).Return(true, nil).AnyTimes()
	mockAWSClient.EXPECT().GetRoleByName(gomock.Any()).Return(iamTypes.Role{}, nil).AnyTimes()
	mockAWSClient.EXPECT().GetCreator().Return(&aws.Creator{
		ARN:       fmt.Sprintf("arn:aws:iam::%s:user/test-user", h.accountID),
		AccountID: h.accountID,
		IsSTS:     false,
		Partition: "aws",
	}, nil).AnyTimes()

	// Create test resource and reconciler
	h.createOCMRoleConfig("test-ocm-role-admin", expinfrav1.ROSAOCMRoleProfileAdmin)
	reconciler := h.buildReconciler(mockAWSClient)

	// Test: Reconcile should succeed and set Ready condition
	req := ctrl.Request{NamespacedName: types.NamespacedName{Name: h.ocmRoleConfig.Name}}
	h.g.Eventually(func(g Gomega) {
		_, err := reconciler.Reconcile(h.testCtx, req)
		g.Expect(err).ToNot(HaveOccurred())

		updatedRoleConfig := &expinfrav1.ROSAOCMRoleConfig{}
		g.Expect(reconciler.Client.Get(h.testCtx, req.NamespacedName, updatedRoleConfig)).ToNot(HaveOccurred())

		g.Expect(updatedRoleConfig.Status.RoleARN).To(Equal(fmt.Sprintf("arn:aws:iam::%s:role/test-OCM-Role-%s", h.accountID, h.externalID)))
		g.Expect(updatedRoleConfig.Status.OrganizationID).To(Equal(h.orgID))

		readyCondition := v1beta1conditions.Get(updatedRoleConfig, expinfrav1.ROSAOCMRoleConfigReadyCondition)
		g.Expect(readyCondition).ToNot(BeNil())
		g.Expect(readyCondition.Status).To(Equal(corev1.ConditionTrue))
	}).WithTimeout(30 * time.Second).Should(Succeed())
}

// TestROSAOCMRoleConfigAdminTagSelfHealing tests that the admin tag is added when role has admin policy but missing tag.
func TestROSAOCMRoleConfigAdminTagSelfHealing(t *testing.T) {
	h := newTestHarness(t)

	// Setup OCM mocks
	h.setupOCMCurrentAccount()
	h.setupOCMLinkCheck(true) // Role already linked

	// Setup AWS mocks - role exists with admin policy but NO admin tag (self-healing scenario)
	mockAWSClient := aws.NewMockClient(h.mockCtrl)
	mockAWSClient.EXPECT().CheckRoleExists(fmt.Sprintf("test-OCM-Role-%s", h.externalID)).Return(
		true,
		fmt.Sprintf("arn:aws:iam::%s:role/test-OCM-Role-%s", h.accountID, h.externalID),
		nil,
	).AnyTimes()
	mockAWSClient.EXPECT().IsAdminRole(fmt.Sprintf("test-OCM-Role-%s", h.externalID)).Return(true, nil).AnyTimes()
	mockAWSClient.EXPECT().GetRoleByName(fmt.Sprintf("test-OCM-Role-%s", h.externalID)).Return(iamTypes.Role{
		RoleName: awsSdk.String(fmt.Sprintf("test-OCM-Role-%s", h.externalID)),
		Arn:      awsSdk.String(fmt.Sprintf("arn:aws:iam::%s:role/test-OCM-Role-%s", h.accountID, h.externalID)),
		Tags:     []iamTypes.Tag{}, // No tags - missing admin tag!
	}, nil).AnyTimes()
	mockAWSClient.EXPECT().GetCreator().Return(&aws.Creator{
		ARN:       fmt.Sprintf("arn:aws:iam::%s:user/test-user", h.accountID),
		AccountID: h.accountID,
		IsSTS:     false,
		Partition: "aws",
	}, nil).AnyTimes()

	// Create test resource and reconciler
	h.createOCMRoleConfig("test-ocm-role-admin-selfheal", expinfrav1.ROSAOCMRoleProfileAdmin)
	reconciler := h.buildReconciler(mockAWSClient)

	// Test: Self-healing should succeed - controller detects admin policy and proceeds
	req := ctrl.Request{NamespacedName: types.NamespacedName{Name: h.ocmRoleConfig.Name}}
	h.g.Eventually(func(g Gomega) {
		_, err := reconciler.Reconcile(h.testCtx, req)
		g.Expect(err).ToNot(HaveOccurred())

		updatedRoleConfig := &expinfrav1.ROSAOCMRoleConfig{}
		g.Expect(reconciler.Client.Get(h.testCtx, req.NamespacedName, updatedRoleConfig)).ToNot(HaveOccurred())

		g.Expect(updatedRoleConfig.Status.RoleARN).To(Equal(fmt.Sprintf("arn:aws:iam::%s:role/test-OCM-Role-%s", h.accountID, h.externalID)))
		g.Expect(updatedRoleConfig.Status.OrganizationID).To(Equal(h.orgID))

		readyCondition := v1beta1conditions.Get(updatedRoleConfig, expinfrav1.ROSAOCMRoleConfigReadyCondition)
		g.Expect(readyCondition).ToNot(BeNil())
		g.Expect(readyCondition.Status).To(Equal(corev1.ConditionTrue))
	}).WithTimeout(30 * time.Second).Should(Succeed())
}

// TestROSAOCMRoleConfigDeleteNoRoleLinked tests deletion when no role is linked.
func TestROSAOCMRoleConfigDeleteNoRoleLinked(t *testing.T) {
	h := newTestHarness(t)

	// Setup OCM mocks
	h.setupOCMCurrentAccount()
	h.setupOCMLinkCheck(false) // No role linked
	unlinkCalled := false
	h.setupOCMUnlinkDelete(&unlinkCalled)

	// Setup AWS mocks
	mockAWSClient := aws.NewMockClient(h.mockCtrl)
	mockAWSClient.EXPECT().CheckRoleExists(fmt.Sprintf("test-OCM-Role-%s", h.externalID)).Return(
		true,
		fmt.Sprintf("arn:aws:iam::%s:role/test-OCM-Role-%s", h.accountID, h.externalID),
		nil,
	).AnyTimes()
	mockAWSClient.EXPECT().IsAdminRole(gomock.Any()).Return(false, nil).AnyTimes()
	mockAWSClient.EXPECT().GetRoleByName(gomock.Any()).Return(iamTypes.Role{}, nil).AnyTimes()
	mockAWSClient.EXPECT().DeleteOCMRole(fmt.Sprintf("test-OCM-Role-%s", h.externalID), false).Return(nil).Times(1)
	mockAWSClient.EXPECT().GetCreator().Return(&aws.Creator{
		ARN:       fmt.Sprintf("arn:aws:iam::%s:user/test-user", h.accountID),
		AccountID: h.accountID,
		IsSTS:     false,
		Partition: "aws",
	}, nil).AnyTimes()

	// Create test resource with finalizer and status
	h.createOCMRoleConfigWithFinalizer("test-ocm-role-delete-nolink", expinfrav1.ROSAOCMRoleProfileStandard)
	h.setOCMRoleConfigStatus(
		fmt.Sprintf("arn:aws:iam::%s:role/test-OCM-Role-%s", h.accountID, h.externalID),
		h.orgID,
	)
	reconciler := h.buildReconciler(mockAWSClient)

	// Delete the CR
	req := ctrl.Request{NamespacedName: types.NamespacedName{Name: h.ocmRoleConfig.Name}}
	err := reconciler.Client.Delete(h.testCtx, h.ocmRoleConfig)
	h.g.Expect(err).ToNot(HaveOccurred())

	// Test: Delete should succeed without calling unlink (no role linked)
	h.g.Eventually(func(g Gomega) {
		_, err := reconciler.Reconcile(h.testCtx, req)
		g.Expect(err).ToNot(HaveOccurred())

		deletedRoleConfig := &expinfrav1.ROSAOCMRoleConfig{}
		err = reconciler.Client.Get(h.testCtx, req.NamespacedName, deletedRoleConfig)
		if err == nil {
			g.Expect(deletedRoleConfig.Finalizers).To(BeEmpty(), "Finalizers should be removed after successful deletion")
		}
	}).WithTimeout(30 * time.Second).Should(Succeed())

	h.g.Expect(unlinkCalled).To(BeFalse(), "Should skip unlink when no role is linked")
}

// TestROSAOCMRoleConfigDeleteLinkedRole tests deletion when role exists and is linked.
func TestROSAOCMRoleConfigDeleteLinkedRole(t *testing.T) {
	h := newTestHarness(t)

	// Setup OCM mocks
	h.setupOCMCurrentAccount()
	h.setupOCMLinkCheck(true) // Role is linked
	unlinkCalled := false
	h.setupOCMUnlinkDelete(&unlinkCalled)

	// Setup AWS mocks
	mockAWSClient := aws.NewMockClient(h.mockCtrl)
	mockAWSClient.EXPECT().CheckRoleExists(fmt.Sprintf("test-OCM-Role-%s", h.externalID)).Return(
		true,
		fmt.Sprintf("arn:aws:iam::%s:role/test-OCM-Role-%s", h.accountID, h.externalID),
		nil,
	).AnyTimes()
	mockAWSClient.EXPECT().IsAdminRole(gomock.Any()).Return(false, nil).AnyTimes()
	mockAWSClient.EXPECT().GetRoleByName(gomock.Any()).Return(iamTypes.Role{}, nil).AnyTimes()
	mockAWSClient.EXPECT().DeleteOCMRole(fmt.Sprintf("test-OCM-Role-%s", h.externalID), false).Return(nil).Times(1)
	mockAWSClient.EXPECT().GetCreator().Return(&aws.Creator{
		ARN:       fmt.Sprintf("arn:aws:iam::%s:user/test-user", h.accountID),
		AccountID: h.accountID,
		IsSTS:     false,
		Partition: "aws",
	}, nil).AnyTimes()

	// Create test resource with finalizer and status
	h.createOCMRoleConfigWithFinalizer("test-ocm-role-delete", expinfrav1.ROSAOCMRoleProfileStandard)
	h.setOCMRoleConfigStatus(
		fmt.Sprintf("arn:aws:iam::%s:role/test-OCM-Role-%s", h.accountID, h.externalID),
		h.orgID,
	)
	reconciler := h.buildReconciler(mockAWSClient)

	// Delete the CR
	req := ctrl.Request{NamespacedName: types.NamespacedName{Name: h.ocmRoleConfig.Name}}
	err := reconciler.Client.Delete(h.testCtx, h.ocmRoleConfig)
	h.g.Expect(err).ToNot(HaveOccurred())

	// Test: Delete should succeed and call unlink (role is linked)
	h.g.Eventually(func(g Gomega) {
		_, err := reconciler.Reconcile(h.testCtx, req)
		g.Expect(err).ToNot(HaveOccurred())

		deletedRoleConfig := &expinfrav1.ROSAOCMRoleConfig{}
		err = reconciler.Client.Get(h.testCtx, req.NamespacedName, deletedRoleConfig)
		if err == nil {
			g.Expect(deletedRoleConfig.Finalizers).To(BeEmpty(), "Finalizers should be removed after successful deletion")
		}
	}).WithTimeout(30 * time.Second).Should(Succeed())

	h.g.Expect(unlinkCalled).To(BeTrue(), "OCM role should be unlinked")
}

// TestROSAOCMRoleConfigDeleteDifferentRoleLinked tests deletion when a different role is linked.
func TestROSAOCMRoleConfigDeleteDifferentRoleLinked(t *testing.T) {
	h := newTestHarness(t)

	// Setup OCM mocks - different role is linked
	h.setupOCMCurrentAccount()
	differentRoleARN := fmt.Sprintf("arn:aws:iam::%s:role/different-OCM-Role-%s", h.accountID, h.externalID)
	h.setupOCMLinkCheckDifferentRole(differentRoleARN)
	unlinkCalled := false
	h.setupOCMUnlinkDelete(&unlinkCalled)

	// Setup AWS mocks
	mockAWSClient := aws.NewMockClient(h.mockCtrl)
	mockAWSClient.EXPECT().CheckRoleExists(fmt.Sprintf("test-OCM-Role-%s", h.externalID)).Return(
		true,
		fmt.Sprintf("arn:aws:iam::%s:role/test-OCM-Role-%s", h.accountID, h.externalID),
		nil,
	).AnyTimes()
	mockAWSClient.EXPECT().IsAdminRole(gomock.Any()).Return(false, nil).AnyTimes()
	mockAWSClient.EXPECT().GetRoleByName(gomock.Any()).Return(iamTypes.Role{}, nil).AnyTimes()
	mockAWSClient.EXPECT().DeleteOCMRole(fmt.Sprintf("test-OCM-Role-%s", h.externalID), false).Return(nil).Times(1)
	mockAWSClient.EXPECT().GetCreator().Return(&aws.Creator{
		ARN:       fmt.Sprintf("arn:aws:iam::%s:user/test-user", h.accountID),
		AccountID: h.accountID,
		IsSTS:     false,
		Partition: "aws",
	}, nil).AnyTimes()

	// Create test resource with finalizer and status
	h.createOCMRoleConfigWithFinalizer("test-ocm-role-delete-diff", expinfrav1.ROSAOCMRoleProfileStandard)
	h.setOCMRoleConfigStatus(
		fmt.Sprintf("arn:aws:iam::%s:role/test-OCM-Role-%s", h.accountID, h.externalID),
		h.orgID,
	)
	reconciler := h.buildReconciler(mockAWSClient)

	// Delete the CR
	req := ctrl.Request{NamespacedName: types.NamespacedName{Name: h.ocmRoleConfig.Name}}
	err := reconciler.Client.Delete(h.testCtx, h.ocmRoleConfig)
	h.g.Expect(err).ToNot(HaveOccurred())

	// Test: Delete should succeed without calling unlink (different role is linked)
	h.g.Eventually(func(g Gomega) {
		_, err := reconciler.Reconcile(h.testCtx, req)
		g.Expect(err).ToNot(HaveOccurred())

		deletedRoleConfig := &expinfrav1.ROSAOCMRoleConfig{}
		err = reconciler.Client.Get(h.testCtx, req.NamespacedName, deletedRoleConfig)
		if err == nil {
			g.Expect(deletedRoleConfig.Finalizers).To(BeEmpty(), "Finalizers should be removed after successful deletion")
		}
	}).WithTimeout(30 * time.Second).Should(Succeed())

	h.g.Expect(unlinkCalled).To(BeFalse(), "Should skip unlink when different role is linked")
}

// TestROSAOCMRoleConfigDeleteWithRetainPolicy tests deletion when DeletionPolicy is Retain.
func TestROSAOCMRoleConfigDeleteWithRetainPolicy(t *testing.T) {
	h := newTestHarness(t)

	h.setupOCMCurrentAccount()

	// Setup OCM endpoints - these should not be called with Retain policy
	unlinkCalled := false
	h.setupOCMUnlinkDelete(&unlinkCalled)

	// Setup AWS mocks - allow normal reconcile calls before deletion
	mockAWSClient := aws.NewMockClient(h.mockCtrl)
	mockAWSClient.EXPECT().GetCreator().Return(&aws.Creator{
		ARN:       fmt.Sprintf("arn:aws:iam::%s:user/test-user", h.accountID),
		AccountID: h.accountID,
		IsSTS:     false,
		Partition: "aws",
	}, nil).AnyTimes()

	// CheckRoleExists may be called during normal reconcile (before deletion timestamp is set)
	// but should indicate role already exists to avoid creation attempts
	mockAWSClient.EXPECT().CheckRoleExists(fmt.Sprintf("test-OCM-Role-%s", h.externalID)).Return(
		true,
		fmt.Sprintf("arn:aws:iam::%s:role/test-OCM-Role-%s", h.accountID, h.externalID),
		nil,
	).AnyTimes()

	mockAWSClient.EXPECT().IsAdminRole(gomock.Any()).Return(false, nil).AnyTimes()
	mockAWSClient.EXPECT().GetRoleByName(gomock.Any()).Return(iamTypes.Role{}, nil).AnyTimes()

	// DeleteOCMRole should never be called (DeletionPolicy=Retain)
	deleteRoleCalled := false
	mockAWSClient.EXPECT().DeleteOCMRole(gomock.Any(), gomock.Any()).DoAndReturn(func(name string, admin bool) error {
		deleteRoleCalled = true
		return nil
	}).Times(0)

	// Create test resource with finalizer
	h.createOCMRoleConfigWithFinalizer("test-ocm-role-retain", expinfrav1.ROSAOCMRoleProfileStandard)

	// Set DeletionPolicy to Retain and update in API server
	h.ocmRoleConfig.Spec.DeletionPolicy = expinfrav1.ROSAOCMRoleDeletionPolicyRetain
	err := testEnv.Client.Update(h.testCtx, h.ocmRoleConfig)
	h.g.Expect(err).ToNot(HaveOccurred())

	h.setOCMRoleConfigStatus(
		fmt.Sprintf("arn:aws:iam::%s:role/test-OCM-Role-%s", h.accountID, h.externalID),
		h.orgID,
	)

	reconciler := h.buildReconciler(mockAWSClient)

	req := ctrl.Request{}
	req.NamespacedName = types.NamespacedName{Name: h.ocmRoleConfig.Name}

	// Delete the CR
	err = reconciler.Client.Delete(h.testCtx, h.ocmRoleConfig)
	h.g.Expect(err).ToNot(HaveOccurred())

	h.g.Eventually(func(g Gomega) {
		_, errReconcile := reconciler.Reconcile(h.testCtx, req)
		g.Expect(errReconcile).ToNot(HaveOccurred())

		// Verify the resource is eventually deleted or finalizers removed
		deletedRoleConfig := &expinfrav1.ROSAOCMRoleConfig{}
		err := reconciler.Client.Get(h.testCtx, req.NamespacedName, deletedRoleConfig)
		if err == nil {
			g.Expect(deletedRoleConfig.Finalizers).To(BeEmpty(), "Finalizers should be removed after successful deletion")
		}
	}).WithTimeout(30 * time.Second).Should(Succeed())

	// Verify unlink was not called (DeletionPolicy=Retain)
	h.g.Expect(unlinkCalled).To(BeFalse(), "Should skip unlink when DeletionPolicy is Retain")
	h.g.Expect(deleteRoleCalled).To(BeFalse(), "Should skip IAM role deletion when DeletionPolicy is Retain")
}
