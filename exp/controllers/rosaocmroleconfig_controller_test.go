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
	"fmt"
	"net/http"
	"testing"
	"time"

	awsSdk "github.com/aws/aws-sdk-go-v2/aws"
	iamv2 "github.com/aws/aws-sdk-go-v2/service/iam"
	iamTypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	stsv2 "github.com/aws/aws-sdk-go-v2/service/sts"
	. "github.com/onsi/gomega"
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
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/rosa"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"
)

func TestROSAOCMRoleConfigReconcileCreate(t *testing.T) {
	RegisterTestingT(t)
	g := NewWithT(t)

	// Generate unique test ID for resource isolation
	testID := generateTestID()

	ssoServer := ocmsdk.MakeTCPServer()
	apiServer := ocmsdk.MakeTCPServer()
	defer ssoServer.Close()
	defer apiServer.Close()
	apiServer.SetAllowUnhandledRequests(true)
	apiServer.SetUnhandledRequestStatusCode(http.StatusInternalServerError)
	ctx := context.TODO()

	// Create the token:
	accessToken := ocmsdk.MakeTokenString("Bearer", 15*time.Minute)

	// Prepare the server:
	ssoServer.AppendHandlers(
		ocmsdk.RespondWithAccessToken(accessToken),
	)
	logger, err := ocmlogging.NewGoLoggerBuilder().
		Debug(false).
		Build()
	Expect(err).ToNot(HaveOccurred())
	// Set up the connection with the fake config
	connection, err := sdk.NewConnectionBuilder().
		Logger(logger).
		Tokens(accessToken).
		URL(apiServer.URL()).
		Build()
	// Initialize client object
	Expect(err).To(BeNil())
	ocmClient := ocm.NewClientWithConnection(connection)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Mock IAM client
	mockIamClient := rosaMocks.NewMockIamApiClient(mockCtrl)

	// Mock GetRole calls - return role not found error to trigger role creation
	mockIamClient.EXPECT().GetRole(gomock.Any(), gomock.Any()).Return(nil, &iamTypes.NoSuchEntityException{
		Message: awsSdk.String("The role with name test-OCM-Role-123 does not exist."),
	}).AnyTimes()

	// Mock CreateRole calls for role creation
	mockIamClient.EXPECT().CreateRole(gomock.Any(), gomock.Any()).Return(&iamv2.CreateRoleOutput{
		Role: &iamTypes.Role{
			RoleName: awsSdk.String("test-OCM-Role-123"),
			Arn:      awsSdk.String("arn:aws:iam::123456789012:role/test-OCM-Role-123"),
		},
	}, nil).AnyTimes()

	// Mock AttachRolePolicy calls
	mockIamClient.EXPECT().AttachRolePolicy(gomock.Any(), gomock.Any()).Return(&iamv2.AttachRolePolicyOutput{}, nil).AnyTimes()

	// Mock CreatePolicy calls
	mockIamClient.EXPECT().CreatePolicy(gomock.Any(), gomock.Any()).Return(&iamv2.CreatePolicyOutput{
		Policy: &iamTypes.Policy{
			PolicyName: awsSdk.String("test-OCM-Role-123"),
			Arn:        awsSdk.String("arn:aws:iam::123456789012:policy/test-OCM-Role-123"),
		},
	}, nil).AnyTimes()

	// Mock GetPolicy calls - return not found for customer-managed policies
	mockIamClient.EXPECT().GetPolicy(gomock.Any(), gomock.Any()).Return(nil, &iamTypes.NoSuchEntityException{
		Message: awsSdk.String("The policy does not exist."),
	}).AnyTimes()

	// Mock tag operations
	mockIamClient.EXPECT().TagRole(gomock.Any(), gomock.Any()).Return(&iamv2.TagRoleOutput{}, nil).AnyTimes()

	// Mock STS client
	mockSTSClient := rosaMocks.NewMockStsApiClient(mockCtrl)
	mockSTSClient.EXPECT().GetCallerIdentity(gomock.Any(), gomock.Any()).Return(&stsv2.GetCallerIdentityOutput{
		Arn:     awsSdk.String("arn:aws:iam::123456789012:user/test-user"),
		Account: awsSdk.String("123456789012"),
		UserId:  awsSdk.String("test-user-id"),
	}, nil).AnyTimes()

	awsClient := aws.New(
		awsSdk.Config{},
		aws.NewLoggerWrapper(logrus.New(), nil),
		mockIamClient,
		rosaMocks.NewMockEc2ApiClient(mockCtrl),
		rosaMocks.NewMockOrganizationsApiClient(mockCtrl),
		rosaMocks.NewMockS3ApiClient(mockCtrl),
		rosaMocks.NewMockSecretsManagerApiClient(mockCtrl),
		mockSTSClient,
		rosaMocks.NewMockCloudFormationApiClient(mockCtrl),
		rosaMocks.NewMockServiceQuotasApiClient(mockCtrl),
		rosaMocks.NewMockServiceQuotasApiClient(mockCtrl),
		&aws.AccessKey{},
		false,
	)

	r := rosacli.NewRuntime()
	r.OCMClient = ocmClient
	r.AWSClient = awsClient
	r.Creator = &aws.Creator{
		ARN:       "arn:aws:iam::123456789012:user/test-user",
		AccountID: "123456789012",
		IsSTS:     false,
		Partition: "aws",
	}

	// Mock OCM API calls
	// Mock current organization call
	apiServer.RouteToHandler("GET", "/api/accounts_mgmt/v1/current_account",
		ocmsdk.RespondWithJSON(
			http.StatusOK, `{"id": "1jbXGbUAJ3eJSnZyDwfplNUMxFM", "organization": {"id": "1jbXGbUAJ3eJSnZyDwfplNUMxFM", "external_id": "13829321"}}`,
		),
	)

	// Mock OCM policies call
	apiServer.RouteToHandler("GET", "/api/clusters_mgmt/v1/aws_inquiries/sts_policies",
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

	// Mock checking if role is already linked (not linked initially, then linked after we link it)
	callCount := 0
	apiServer.RouteToHandler("GET", "/api/accounts_mgmt/v1/organizations/1jbXGbUAJ3eJSnZyDwfplNUMxFM/labels/sts_ocm_role",
		func(w http.ResponseWriter, r *http.Request) {
			callCount++
			if callCount == 1 {
				// First call: role not linked yet (404)
				ocmsdk.RespondWithJSON(http.StatusNotFound, `{}`)(w, r)
			} else {
				// Subsequent calls: role is linked
				ocmsdk.RespondWithJSON(http.StatusOK, `{"value": "arn:aws:iam::123456789012:role/test-OCM-Role-123"}`)(w, r)
			}
		},
	)

	// Mock link role to organization
	apiServer.RouteToHandler("POST", "/api/accounts_mgmt/v1/organizations/1jbXGbUAJ3eJSnZyDwfplNUMxFM/labels",
		ocmsdk.RespondWithJSON(
			http.StatusCreated, `{"key": "sts_ocm_role", "value": "arn:aws:iam::123456789012:role/test-OCM-Role-123"}`,
		),
	)

	// Create ROSAOCMRoleConfig CR
	ocmRoleConfig := &expinfrav1.ROSAOCMRoleConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("test-ocm-role-%s", testID),
		},
		Spec: expinfrav1.ROSAOCMRoleConfigSpec{
			RolePrefix: "test",
			Profile:    expinfrav1.ROSAOCMRoleProfileStandard,
		},
	}

	createObject(g, ocmRoleConfig, "")
	defer cleanupObject(g, ocmRoleConfig)

	// Setup the reconciler
	reconciler := &ROSAOCMRoleConfigReconciler{
		Client:   testEnv.Client,
		Runtime:  r,
		Recorder: record.NewFakeRecorder(32),
	}
	req := ctrl.Request{}
	req.NamespacedName = types.NamespacedName{Name: ocmRoleConfig.Name}

	g.Eventually(func(g Gomega) {
		// Call the Reconcile function
		_, errReconcile := reconciler.Reconcile(ctx, req)

		// Assertions
		g.Expect(errReconcile).ToNot(HaveOccurred())

		// Check the status of the ROSAOCMRoleConfig resource
		updatedRoleConfig := &expinfrav1.ROSAOCMRoleConfig{}
		err = reconciler.Client.Get(ctx, req.NamespacedName, updatedRoleConfig)
		g.Expect(err).ToNot(HaveOccurred())

		// Verify status fields are set
		g.Expect(updatedRoleConfig.Status.RoleARN).To(Equal("arn:aws:iam::123456789012:role/test-OCM-Role-123"))
		g.Expect(updatedRoleConfig.Status.OrganizationID).To(Equal("1jbXGbUAJ3eJSnZyDwfplNUMxFM"))

		// Ready condition should be true
		readyCondition := v1beta1conditions.Get(updatedRoleConfig, expinfrav1.ROSAOCMRoleConfigReadyCondition)
		g.Expect(readyCondition).ToNot(BeNil())
		g.Expect(readyCondition.Status).To(Equal(corev1.ConditionTrue))
	}).WithTimeout(30 * time.Second).Should(Succeed())
}

func TestROSAOCMRoleConfigReconcileExist(t *testing.T) {
	RegisterTestingT(t)
	g := NewWithT(t)

	// Generate unique test ID for resource isolation
	testID := generateTestID()

	ssoServer := ocmsdk.MakeTCPServer()
	apiServer := ocmsdk.MakeTCPServer()
	defer ssoServer.Close()
	defer apiServer.Close()
	apiServer.SetAllowUnhandledRequests(true)
	apiServer.SetUnhandledRequestStatusCode(http.StatusInternalServerError)
	ctx := context.TODO()

	// Create the token:
	accessToken := ocmsdk.MakeTokenString("Bearer", 15*time.Minute)

	// Prepare the server:
	ssoServer.AppendHandlers(
		ocmsdk.RespondWithAccessToken(accessToken),
	)
	logger, err := ocmlogging.NewGoLoggerBuilder().
		Debug(false).
		Build()
	Expect(err).ToNot(HaveOccurred())
	// Set up the connection with the fake config
	connection, err := sdk.NewConnectionBuilder().
		Logger(logger).
		Tokens(accessToken).
		URL(apiServer.URL()).
		Build()
	// Initialize client object
	Expect(err).To(BeNil())
	ocmClient := ocm.NewClientWithConnection(connection)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Use mock AWS client (ROSA client, not IAM client)
	mockAWSClient := aws.NewMockClient(mockCtrl)

	// Mock CheckRoleExists to return existing role (roleName includes externalID: test-OCM-Role-13829321)
	mockAWSClient.EXPECT().CheckRoleExists("test-OCM-Role-13829321").Return(
		true,
		"arn:aws:iam::123456789012:role/test-OCM-Role-13829321",
		nil,
	).AnyTimes()

	// Mock IsAdminRole to return false (standard role)
	mockAWSClient.EXPECT().IsAdminRole("test-OCM-Role-13829321").Return(false, nil).AnyTimes()

	// Mock IsNoConsoleRole to return false (standard role)
	mockAWSClient.EXPECT().IsNoConsoleRole("test-OCM-Role-13829321").Return(false, nil).AnyTimes()

	mockAWSClient.EXPECT().GetCreator().Return(&aws.Creator{
		ARN:       "arn:aws:iam::123456789012:user/test-user",
		AccountID: "123456789012",
		IsSTS:     false,
		Partition: "aws",
	}, nil).AnyTimes()

	awsClient := mockAWSClient

	r := rosacli.NewRuntime()
	r.OCMClient = ocmClient
	r.AWSClient = awsClient
	r.Creator = &aws.Creator{
		ARN:       "arn:aws:iam::123456789012:user/test-user",
		AccountID: "123456789012",
		IsSTS:     false,
		Partition: "aws",
	}

	// Mock OCM API calls
	apiServer.RouteToHandler("GET", "/api/accounts_mgmt/v1/current_account",
		ocmsdk.RespondWithJSON(
			http.StatusOK, `{"id": "1jbXGbUAJ3eJSnZyDwfplNUMxFM", "organization": {"id": "1jbXGbUAJ3eJSnZyDwfplNUMxFM", "external_id": "13829321"}}`,
		),
	)

	// Mock checking if role is already linked (already linked)
	apiServer.RouteToHandler("GET", "/api/accounts_mgmt/v1/organizations/1jbXGbUAJ3eJSnZyDwfplNUMxFM/labels/sts_ocm_role",
		ocmsdk.RespondWithJSON(
			http.StatusOK, `{"value": "arn:aws:iam::123456789012:role/test-OCM-Role-13829321"}`,
		),
	)

	// Create ROSAOCMRoleConfig CR
	ocmRoleConfig := &expinfrav1.ROSAOCMRoleConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("test-ocm-role-exist-%s", testID),
		},
		Spec: expinfrav1.ROSAOCMRoleConfigSpec{
			RolePrefix: "test",
			Profile:    expinfrav1.ROSAOCMRoleProfileStandard,
		},
	}

	createObject(g, ocmRoleConfig, "")
	defer cleanupObject(g, ocmRoleConfig)

	// Setup the reconciler
	reconciler := &ROSAOCMRoleConfigReconciler{
		Client:   testEnv.Client,
		Runtime:  r,
		Recorder: record.NewFakeRecorder(32),
	}

	req := ctrl.Request{}
	req.NamespacedName = types.NamespacedName{Name: ocmRoleConfig.Name}

	g.Eventually(func(g Gomega) {
		// Call the Reconcile function
		_, errReconcile := reconciler.Reconcile(ctx, req)

		// Assertions - since role already exists and is linked, reconciliation should succeed
		g.Expect(errReconcile).ToNot(HaveOccurred())

		// Check the status of the ROSAOCMRoleConfig resource
		updatedRoleConfig := &expinfrav1.ROSAOCMRoleConfig{}
		g.Expect(reconciler.Client.Get(ctx, req.NamespacedName, updatedRoleConfig)).ToNot(HaveOccurred())

		// Verify that existing role is preserved
		g.Expect(updatedRoleConfig.Status.RoleARN).To(Equal("arn:aws:iam::123456789012:role/test-OCM-Role-13829321"))
		g.Expect(updatedRoleConfig.Status.OrganizationID).To(Equal("1jbXGbUAJ3eJSnZyDwfplNUMxFM"))

		// Should have a condition indicating success - expect Ready condition to be True
		readyCondition := v1beta1conditions.Get(updatedRoleConfig, expinfrav1.ROSAOCMRoleConfigReadyCondition)
		g.Expect(readyCondition).ToNot(BeNil())
		g.Expect(readyCondition.Status).To(Equal(corev1.ConditionTrue))
	}).WithTimeout(30 * time.Second).WithPolling(500 * time.Millisecond).Should(Succeed())
}

func TestROSAOCMRoleConfigSetUpRuntimeWithExpiredAWSCredentials(t *testing.T) {
	RegisterTestingT(t)
	g := NewWithT(t)
	ctx := context.TODO()

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
		Client:         testEnv.Client,
		ROSAOCMRoleConfig:  ocmRoleConfig,
		ControllerName: "rosaocmroleconfig",
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
		NewOCMClient: func(ctx context.Context, scope rosa.OCMSecretsRetriever) (rosa.OCMClient, error) {
			ocmClientCallCount++
			return ocmClient, nil
		},
		Runtime: nil, // Start with nil Runtime
	}

	// ==================== FIRST RECONCILIATION ====================
	t.Log("First reconciliation: AWS credentials are missing/expired")

	err = reconciler.setUpRuntime(ctx, scope)

	g.Expect(err).To(HaveOccurred(),
		"setUpRuntime should return error when AWS client creation fails")
	g.Expect(err.Error()).To(ContainSubstring("failed to create aws client"),
		"Error should indicate AWS client failure")

	g.Expect(reconciler.Runtime).To(BeNil(),
		"Runtime MUST be nil after failed initialization (atomic assignment fix)")

	g.Expect(ocmClientCallCount).To(Equal(1),
		"NewOCMClient should be called once on first attempt")

	// ==================== SECOND RECONCILIATION ====================
	t.Log("Second reconciliation: Retry with still-missing AWS credentials")

	err = reconciler.setUpRuntime(ctx, scope)

	g.Expect(err).To(HaveOccurred(),
		"setUpRuntime should still fail with missing AWS credentials")

	g.Expect(reconciler.Runtime).To(BeNil(),
		"Runtime should still be nil after second failed attempt")

	// NewOCMClient was called AGAIN (proves no early return)
	g.Expect(ocmClientCallCount).To(Equal(2),
		"NewOCMClient should be called twice - proves guard clause allows retry when Runtime is nil")
}

// TestSetUpRuntimeIdempotency verifies that setUpRuntime returns early
// when Runtime is already fully initialized.
func TestROSAOCMRoleConfigSetUpRuntimeIdempotency(t *testing.T) {
	RegisterTestingT(t)
	g := NewWithT(t)
	ctx := context.TODO()

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

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockIamClient := rosaMocks.NewMockIamApiClient(mockCtrl)
	mockSTSClient := rosaMocks.NewMockStsApiClient(mockCtrl)

	awsClient := aws.New(
		awsSdk.Config{},
		aws.NewLoggerWrapper(logrus.New(), nil),
		mockIamClient,
		rosaMocks.NewMockEc2ApiClient(mockCtrl),
		rosaMocks.NewMockOrganizationsApiClient(mockCtrl),
		rosaMocks.NewMockS3ApiClient(mockCtrl),
		rosaMocks.NewMockSecretsManagerApiClient(mockCtrl),
		mockSTSClient,
		rosaMocks.NewMockCloudFormationApiClient(mockCtrl),
		rosaMocks.NewMockServiceQuotasApiClient(mockCtrl),
		rosaMocks.NewMockServiceQuotasApiClient(mockCtrl),
		&aws.AccessKey{},
		false,
	)

	runtime := rosacli.NewRuntime()
	runtime.OCMClient = ocmClient
	runtime.AWSClient = awsClient
	runtime.Creator = &aws.Creator{
		ARN:       "arn:aws:iam::123456789012:user/test",
		AccountID: "123456789012",
		IsSTS:     false,
	}

	callCount := 0
	reconciler := &ROSAOCMRoleConfigReconciler{
		Runtime: runtime, // already initialized
		NewOCMClient: func(ctx context.Context, scope rosa.OCMSecretsRetriever) (rosa.OCMClient, error) {
			callCount++
			return ocmClient, nil
		},
	}

	scope := &scope.ROSAOCMRoleConfigScope{}

	// Call setUpRuntime - should return early
	err = reconciler.setUpRuntime(ctx, scope)
	g.Expect(err).ToNot(HaveOccurred())

	// Runtime should be unchanged (same instance)
	g.Expect(reconciler.Runtime).To(BeIdenticalTo(runtime),
		"Runtime should not be recreated when already initialized")

	// NewOCMClient should NOT be called (early return)
	g.Expect(callCount).To(Equal(0),
		"NewOCMClient should not be called when Runtime already exists")

	// Call again - should still return early
	err = reconciler.setUpRuntime(ctx, scope)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(callCount).To(Equal(0), "Still should not call NewOCMClient")
}
