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
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
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
	ctrl "sigs.k8s.io/controller-runtime"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	rosacontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"
)

// generateTestID creates a unique identifier for test resources.
func generateTestID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Intn(10000))
}

// makeRuntimeFactory returns a runtimeFactory that always returns the given pre-built runtime.
// This replaces the old pattern of setting reconciler.Runtime directly.
func makeRuntimeFactory(rt *rosacli.Runtime) func(context.Context, *scope.RosaRoleConfigScope) (*rosacli.Runtime, error) {
	return func(_ context.Context, _ *scope.RosaRoleConfigScope) (*rosacli.Runtime, error) {
		return rt, nil
	}
}

func TestROSARoleConfigReconcileCreate(t *testing.T) {
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
	// mock iam client to expect ListRoles call
	mockIamClient := rosaMocks.NewMockIamApiClient(mockCtrl)
	mockIamClient.EXPECT().ListRoles(gomock.Any(), gomock.Any()).Return(&iamv2.ListRolesOutput{
		Roles: []iamTypes.Role{},
	}, nil).AnyTimes()

	mockIamClient.EXPECT().ListOpenIDConnectProviders(gomock.Any(), gomock.Any()).Return(&iamv2.ListOpenIDConnectProvidersOutput{
		OpenIDConnectProviderList: []iamTypes.OpenIDConnectProviderListEntry{},
	}, nil).AnyTimes()

	// Mock GetRole calls - return role not found error to trigger role creation
	mockIamClient.EXPECT().GetRole(gomock.Any(), gomock.Any()).Return(nil, &iamTypes.NoSuchEntityException{
		Message: awsSdk.String("The role with name test-role does not exist."),
	}).AnyTimes()

	// Mock CreateRole calls for role creation
	mockIamClient.EXPECT().CreateRole(gomock.Any(), gomock.Any()).Return(&iamv2.CreateRoleOutput{
		Role: &iamTypes.Role{
			RoleName: awsSdk.String("test-role"),
			Arn:      awsSdk.String("arn:aws:iam::123456789012:role/test-role"),
		},
	}, nil).AnyTimes()

	providerARN := "test-oidc-id-created"
	mockIamClient.EXPECT().CreateOpenIDConnectProvider(gomock.Any(), gomock.Any(), gomock.Any()).Return(
		&iamv2.CreateOpenIDConnectProviderOutput{OpenIDConnectProviderArn: &providerARN}, nil).AnyTimes()

	// Mock AttachRolePolicy calls
	mockIamClient.EXPECT().AttachRolePolicy(gomock.Any(), gomock.Any()).Return(&iamv2.AttachRolePolicyOutput{}, nil).AnyTimes()

	// Mock CreatePolicy calls
	mockIamClient.EXPECT().CreatePolicy(gomock.Any(), gomock.Any()).Return(&iamv2.CreatePolicyOutput{
		Policy: &iamTypes.Policy{
			PolicyName: awsSdk.String("test-policy"),
			Arn:        awsSdk.String("arn:aws:iam::123456789012:policy/test-policy"),
		},
	}, nil).AnyTimes()

	// Mock GetPolicy calls - return success for AWS managed policies, not found for others
	mockIamClient.EXPECT().GetPolicy(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, input *iamv2.GetPolicyInput) (*iamv2.GetPolicyOutput, error) {
		switch *input.PolicyArn {
		case "arn:aws:iam::aws:policy/sts_hcp_installer_permission_policy":
			return &iamv2.GetPolicyOutput{
				Policy: &iamTypes.Policy{
					PolicyName: awsSdk.String("sts_hcp_installer_permission_policy"),
					Arn:        awsSdk.String("arn:aws:iam::aws:policy/sts_hcp_installer_permission_policy"),
				},
			}, nil
		case "arn:aws:iam::aws:policy/sts_hcp_support_permission_policy":
			return &iamv2.GetPolicyOutput{
				Policy: &iamTypes.Policy{
					PolicyName: awsSdk.String("sts_hcp_support_permission_policy"),
					Arn:        awsSdk.String("arn:aws:iam::aws:policy/sts_hcp_support_permission_policy"),
				},
			}, nil
		case "arn:aws:iam::aws:policy/sts_hcp_worker_permission_policy":
			return &iamv2.GetPolicyOutput{
				Policy: &iamTypes.Policy{
					PolicyName: awsSdk.String("sts_hcp_worker_permission_policy"),
					Arn:        awsSdk.String("arn:aws:iam::aws:policy/sts_hcp_worker_permission_policy"),
				},
			}, nil
		default:
			return nil, &iamTypes.NoSuchEntityException{
				Message: awsSdk.String("The policy does not exist."),
			}
		}
	}).AnyTimes()

	// Mock ListPolicies calls - return expected ROSA managed policies
	mockIamClient.EXPECT().ListPolicies(gomock.Any(), gomock.Any()).Return(&iamv2.ListPoliciesOutput{
		Policies: []iamTypes.Policy{
			{
				PolicyName: awsSdk.String("sts_hcp_installer_permission_policy"),
				Arn:        awsSdk.String("arn:aws:iam::aws:policy/sts_hcp_installer_permission_policy"),
			},
			{
				PolicyName: awsSdk.String("sts_hcp_support_permission_policy"),
				Arn:        awsSdk.String("arn:aws:iam::aws:policy/sts_hcp_support_permission_policy"),
			},
			{
				PolicyName: awsSdk.String("sts_hcp_worker_permission_policy"),
				Arn:        awsSdk.String("arn:aws:iam::aws:policy/sts_hcp_worker_permission_policy"),
			},
		},
	}, nil).AnyTimes()

	// mock sts - add common STS calls that might be needed during role creation
	mockSTSClient := rosaMocks.NewMockStsApiClient(mockCtrl)
	mockSTSClient.EXPECT().GetCallerIdentity(gomock.Any(), gomock.Any()).Return(&stsv2.GetCallerIdentityOutput{
		Arn:     awsSdk.String("fake"),
		Account: awsSdk.String("123"),
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

	rt := rosacli.NewRuntime()
	rt.OCMClient = ocmClient
	rt.AWSClient = awsClient
	rt.Creator = &aws.Creator{
		ARN:       "fake",
		AccountID: "123",
		IsSTS:     false,
	}
	// Mock OCM API calls using path-based routing
	apiServer.RouteToHandler("GET", "/api/clusters_mgmt/v1/aws_inquiries/sts_policies",
		func(w http.ResponseWriter, r *http.Request) {
			query := r.URL.Query().Get("search")
			if strings.Contains(query, "AccountRole") {
				// Return AccountRole policies
				ocmsdk.RespondWithJSON(http.StatusOK, `{
					"items": [
						{
							"id": "sts_hcp_installer_permission_policy",
							"arn": "arn:aws:iam::aws:policy/sts_hcp_installer_permission_policy",
							"type": "AccountRole"
						},
						{
							"id": "sts_hcp_support_permission_policy",
							"arn": "arn:aws:iam::aws:policy/sts_hcp_support_permission_policy",
							"type": "AccountRole"
						},
						{
							"id": "sts_hcp_worker_permission_policy",
							"arn": "arn:aws:iam::aws:policy/sts_hcp_worker_permission_policy",
							"type": "AccountRole"
						},
						{
							"id": "sts_hcp_instance_worker_permission_policy",
							"arn": "arn:aws:iam::aws:policy/sts_hcp_instance_worker_permission_policy",
							"type": "AccountRole"
						}
					]
				}`)(w, r)
			} else if strings.Contains(query, "OperatorRole") {
				// Return OperatorRole policies
				ocmsdk.RespondWithJSON(http.StatusOK, `{
					"items": [
						{
							"id": "openshift_hcp_ingress_policy",
							"arn": "arn:aws:iam::aws:policy/openshift_hcp_ingress_policy",
							"type": "OperatorRole"
						},
						{
							"id": "openshift_hcp_image_registry_policy",
							"arn": "arn:aws:iam::aws:policy/openshift_hcp_image_registry_policy",
							"type": "OperatorRole"
						},
						{
							"id": "openshift_hcp_storage_policy",
							"arn": "arn:aws:iam::aws:policy/openshift_hcp_storage_policy",
							"type": "OperatorRole"
						},
						{
							"id": "openshift_hcp_network_policy",
							"arn": "arn:aws:iam::aws:policy/openshift_hcp_network_policy",
							"type": "OperatorRole"
						},
						{
							"id": "openshift_hcp_kube_controller_policy",
							"arn": "arn:aws:iam::aws:policy/openshift_hcp_kube_controller_policy",
							"type": "OperatorRole"
						},
						{
							"id": "openshift_hcp_node_pool_policy",
							"arn": "arn:aws:iam::aws:policy/openshift_hcp_node_pool_policy",
							"type": "OperatorRole"
						},
						{
							"id": "openshift_hcp_control_plane_policy",
							"arn": "arn:aws:iam::aws:policy/openshift_hcp_control_plane_policy",
							"type": "OperatorRole"
						},
						{
							"id": "openshift_hcp_kms_policy",
							"arn": "arn:aws:iam::aws:policy/openshift_hcp_kms_policy",
							"type": "OperatorRole"
						}
					]
				}`)(w, r)
			} else {
				// Default response for other queries
				ocmsdk.RespondWithJSON(http.StatusOK, `{"items": []}`)(w, r)
			}
		})

	// Mock ocm API calls - first call gets tris response
	apiServer.AppendHandlers(
		ocmsdk.RespondWithJSON(
			http.StatusOK, "",
		),
	)
	// Mock GetOidcConfig call
	apiServer.AppendHandlers(
		ocmsdk.RespondWithJSON(
			http.StatusOK, `{"id": "test-oidc-id", "issuer_url": "https://test.oidc.url"}`,
		),
	)
	// Mock OIDC config creation calls - POST /api/clusters_mgmt/v1/oidc_configs
	apiServer.RouteToHandler("POST", "/api/clusters_mgmt/v1/oidc_configs",
		ocmsdk.RespondWithJSON(
			http.StatusCreated, `{"id": "test-oidc-id-created", "issuer_url": "https://test.oidc.url"}`,
		),
	)
	// Additional OIDC config call mock for GET requests
	apiServer.RouteToHandler("GET", "/api/clusters_mgmt/v1/oidc_configs/test-oidc-id-created",
		ocmsdk.RespondWithJSON(
			http.StatusOK, `{"id": "test-oidc-id-created", "issuer_url": "https://test.oidc.url"}`,
		),
	)

	// Mock GetAllCredRequests call
	apiServer.AppendHandlers(
		ocmsdk.RespondWithJSON(
			http.StatusOK, `[]`,
		),
	)
	// Mock HasAClusterUsingOperatorRolesPrefix call
	apiServer.AppendHandlers(
		ocmsdk.RespondWithJSON(
			http.StatusOK, `false`,
		),
	)
	// GET /api/clusters_mgmt/v1/products/rosa/technology_previews/hcp-zero-egress
	apiServer.AppendHandlers(
		ocmsdk.RespondWithJSON(
			http.StatusInternalServerError, "",
		),
	)

	// Create CRs with unique names to avoid conflicts
	ns, err := testEnv.CreateNamespace(ctx, fmt.Sprintf("test-namespace-%s", testID))
	rosaRoleConfig := &expinfrav1.ROSARoleConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:       fmt.Sprintf("test-rosa-role-%s", testID),
			Namespace:  ns.Name,
			Finalizers: []string{expinfrav1.RosaRoleConfigFinalizer},
		},
		Spec: expinfrav1.ROSARoleConfigSpec{
			AccountRoleConfig: expinfrav1.AccountRoleConfig{
				Prefix:  "test",
				Version: "4.15.0",
			},
			OperatorRoleConfig: expinfrav1.OperatorRoleConfig{
				Prefix: "test",
			},
			OidcProviderType: expinfrav1.Managed,
		},
	}
	g.Expect(err).ToNot(HaveOccurred())

	createObject(g, rosaRoleConfig, ns.Name)
	defer cleanupObject(g, rosaRoleConfig)

	// Setup the reconciler using runtimeFactory to inject mock clients per reconciliation.
	reconciler := &ROSARoleConfigReconciler{
		Client:         testEnv.Client,
		runtimeFactory: makeRuntimeFactory(rt),
	}
	req := ctrl.Request{}
	req.NamespacedName = types.NamespacedName{Name: rosaRoleConfig.Name, Namespace: rosaRoleConfig.Namespace}

	g.Eventually(func(g Gomega) {
		// Call the Reconcile function
		_, errReconcile := reconciler.Reconcile(ctx, req)

		// Assertions - expect the installer role empty error since AccountRolesRef is not populated yet
		g.Expect(errReconcile).ToNot(HaveOccurred())

		// Check the status of the ROSARoleConfig resource
		updatedRoleConfig := &expinfrav1.ROSARoleConfig{}
		err = reconciler.Client.Get(ctx, req.NamespacedName, updatedRoleConfig)
		g.Expect(err).ToNot(HaveOccurred())

		// We expect only oidcID to be set with first reconcile happen, Account roles and Operator roles should be empty
		g.Expect(updatedRoleConfig.Status.OIDCID).To(Equal("test-oidc-id-created"))
		g.Expect(updatedRoleConfig.Status.AccountRolesRef).To(Equal(expinfrav1.AccountRolesRef{}))
		g.Expect(updatedRoleConfig.Status.OperatorRolesRef).To(Equal(rosacontrolplanev1.AWSRolesRef{}))

		// Ready condition should be false.
		for _, condition := range updatedRoleConfig.Status.Conditions {
			if condition.Type == expinfrav1.RosaRoleConfigReadyCondition {
				g.Expect(condition.Status).To(Equal(corev1.ConditionFalse))
				break
			}
		}
	}).WithTimeout(30 * time.Second).Should(Succeed())
}

func TestROSARoleConfigReconcileExist(t *testing.T) {
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
	// mock iam client to expect ListRoles call - return existing account roles and operator roles

	mockAWSClient := aws.NewMockClient(mockCtrl)
	mockAWSClient.EXPECT().HasManagedPolicies(gomock.Any()).Return(false, nil).AnyTimes()
	mockAWSClient.EXPECT().HasHostedCPPolicies(gomock.Any()).Return(true, nil).AnyTimes()

	// Return existing account roles by exact name (avoids ListAccountRoles → ListRoleTags throttling).
	mockAWSClient.EXPECT().GetRoleByName("test-HCP-ROSA-Installer-Role").Return(
		iamTypes.Role{Arn: awsSdk.String("arn:aws:iam::123456789012:role/test-HCP-ROSA-Installer-Role")}, nil).AnyTimes()
	mockAWSClient.EXPECT().GetRoleByName("test-HCP-ROSA-Support-Role").Return(
		iamTypes.Role{Arn: awsSdk.String("arn:aws:iam::123456789012:role/test-HCP-ROSA-Support-Role")}, nil).AnyTimes()
	mockAWSClient.EXPECT().GetRoleByName("test-HCP-ROSA-Worker-Role").Return(
		iamTypes.Role{Arn: awsSdk.String("arn:aws:iam::123456789012:role/test-HCP-ROSA-Worker-Role")}, nil).AnyTimes()

	// Return existing operator roles by exact name (avoids ListOperatorRoles → ListRoleTags throttling).
	mockAWSClient.EXPECT().GetRoleByName("test-openshift-ingress-operator-cloud-credentials").Return(
		iamTypes.Role{Arn: awsSdk.String("arn:aws:iam::123456789012:role/test-openshift-ingress-operator-cloud-credentials")}, nil).AnyTimes()
	mockAWSClient.EXPECT().GetRoleByName("test-openshift-image-registry-installer-cloud-credentials").Return(
		iamTypes.Role{Arn: awsSdk.String("arn:aws:iam::123456789012:role/test-openshift-image-registry-installer-cloud-credentials")}, nil).AnyTimes()
	mockAWSClient.EXPECT().GetRoleByName("test-openshift-cluster-csi-drivers-ebs-cloud-credentials").Return(
		iamTypes.Role{Arn: awsSdk.String("arn:aws:iam::123456789012:role/test-openshift-cluster-csi-drivers-ebs-cloud-credentials")}, nil).AnyTimes()
	mockAWSClient.EXPECT().GetRoleByName("test-openshift-cloud-network-config-controller-cloud-credentials").Return(
		iamTypes.Role{Arn: awsSdk.String("arn:aws:iam::123456789012:role/test-openshift-cloud-network-config-controller-cloud-credentials")}, nil).AnyTimes()
	mockAWSClient.EXPECT().GetRoleByName("test-kube-system-kube-controller-manager").Return(
		iamTypes.Role{Arn: awsSdk.String("arn:aws:iam::123456789012:role/test-kube-system-kube-controller-manager")}, nil).AnyTimes()
	mockAWSClient.EXPECT().GetRoleByName("test-kube-system-capa-controller-manager").Return(
		iamTypes.Role{Arn: awsSdk.String("arn:aws:iam::123456789012:role/test-kube-system-capa-controller-manager")}, nil).AnyTimes()
	mockAWSClient.EXPECT().GetRoleByName("test-kube-system-control-plane-operator").Return(
		iamTypes.Role{Arn: awsSdk.String("arn:aws:iam::123456789012:role/test-kube-system-control-plane-operator")}, nil).AnyTimes()
	mockAWSClient.EXPECT().GetRoleByName("test-kube-system-kms-provider").Return(
		iamTypes.Role{Arn: awsSdk.String("arn:aws:iam::123456789012:role/test-kube-system-kms-provider")}, nil).AnyTimes()

	// Return existing OIDC provider ARN by issuer URL (avoids ListOpenIDConnectProviderTags throttling).
	mockAWSClient.EXPECT().GetOpenIDConnectProviderByOidcEndpointUrl("https://test.existing.oidc.url").Return(
		"arn:aws:iam::123456789012:oidc-provider/test-existing-oidc-id", nil).AnyTimes()

	mockAWSClient.EXPECT().GetCreator().Return(&aws.Creator{
		ARN:       "arn:aws:iam::123456789012:user/test-user",
		AccountID: "123456789012",
		IsSTS:     false,
	}, nil).AnyTimes()

	rt := rosacli.NewRuntime()
	rt.OCMClient = ocmClient
	rt.AWSClient = mockAWSClient
	rt.Creator = &aws.Creator{
		ARN:       "arn:aws:iam::123456789012:user/test-user",
		AccountID: "123456789012",
		IsSTS:     false,
	}

	// mock ocm API calls - first call gets tris response
	apiServer.AppendHandlers(
		ocmsdk.RespondWithJSON(
			http.StatusOK, "",
		),
	)
	// Mock GetOidcConfig call - return existing OIDC config
	apiServer.AppendHandlers(
		ocmsdk.RespondWithJSON(
			http.StatusOK, `{"id": "test-existing-oidc-id", "issuer_url": "https://test.existing.oidc.url"}`,
		),
	)
	// Mock GetAllClusters call
	apiServer.AppendHandlers(
		ocmsdk.RespondWithJSON(
			http.StatusOK, `{"items": []}`,
		),
	)
	// Mock GetAllCredRequests call
	apiServer.AppendHandlers(
		ocmsdk.RespondWithJSON(
			http.StatusOK, `{}`,
		),
	)
	// Mock HasAClusterUsingOperatorRolesPrefix call
	apiServer.AppendHandlers(
		ocmsdk.RespondWithJSON(
			http.StatusOK, `false`,
		),
	)

	// Mock existing OIDC config GET request
	apiServer.RouteToHandler("GET", "/api/clusters_mgmt/v1/oidc_configs/test-existing-oidc-id",
		ocmsdk.RespondWithJSON(
			http.StatusOK, `{"id": "test-existing-oidc-id", "issuer_url": "https://test.existing.oidc.url"}`,
		),
	)

	// Create CRs with unique names to avoid conflicts
	ns, err := testEnv.CreateNamespace(ctx, fmt.Sprintf("test-namespace-existing-%s", testID))
	g.Expect(err).ToNot(HaveOccurred())

	rosaRoleConfig := &expinfrav1.ROSARoleConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:       fmt.Sprintf("test-rosarole-existing-%s", testID),
			Namespace:  ns.Name,
			Finalizers: []string{expinfrav1.RosaRoleConfigFinalizer},
		},
		Spec: expinfrav1.ROSARoleConfigSpec{
			AccountRoleConfig: expinfrav1.AccountRoleConfig{
				Prefix:  "test",
				Version: "4.15.0",
			},
			OperatorRoleConfig: expinfrav1.OperatorRoleConfig{
				Prefix: "test",
			},
			OidcProviderType: expinfrav1.Managed,
		},
		Status: expinfrav1.ROSARoleConfigStatus{
			OIDCID: "test-existing-oidc-id",
		},
	}

	createObject(g, rosaRoleConfig, ns.Name)
	defer cleanupObject(g, rosaRoleConfig)

	// Setup the reconciler using runtimeFactory to inject mock clients per reconciliation.
	reconciler := &ROSARoleConfigReconciler{
		Client:         testEnv.Client,
		runtimeFactory: makeRuntimeFactory(rt),
	}

	req := ctrl.Request{}
	req.NamespacedName = types.NamespacedName{Name: rosaRoleConfig.Name, Namespace: rosaRoleConfig.Namespace}

	g.Eventually(func(g Gomega) {
		// Call the Reconcile function
		_, errReconcile := reconciler.Reconcile(ctx, req)

		// Assertions - since all resources exist, reconciliation should succeed
		g.Expect(errReconcile).ToNot(HaveOccurred())

		// Check the status of the ROSARoleConfig resource
		updatedRoleConfig := &expinfrav1.ROSARoleConfig{}
		g.Expect(reconciler.Client.Get(ctx, req.NamespacedName, updatedRoleConfig)).ToNot(HaveOccurred())

		// Verify that all existing account roles are preserved
		g.Expect(updatedRoleConfig.Status.AccountRolesRef.InstallerRoleARN).To(Equal("arn:aws:iam::123456789012:role/test-HCP-ROSA-Installer-Role"))
		g.Expect(updatedRoleConfig.Status.AccountRolesRef.SupportRoleARN).To(Equal("arn:aws:iam::123456789012:role/test-HCP-ROSA-Support-Role"))
		g.Expect(updatedRoleConfig.Status.AccountRolesRef.WorkerRoleARN).To(Equal("arn:aws:iam::123456789012:role/test-HCP-ROSA-Worker-Role"))

		// Verify OIDC config is preserved
		g.Expect(updatedRoleConfig.Status.OIDCID).To(Equal("test-existing-oidc-id"))
		g.Expect(updatedRoleConfig.Status.OIDCProviderARN).To(Equal("arn:aws:iam::123456789012:oidc-provider/test-existing-oidc-id"))

		// Verify operator roles are populated with existing roles
		g.Expect(updatedRoleConfig.Status.OperatorRolesRef.IngressARN).To(Equal("arn:aws:iam::123456789012:role/test-openshift-ingress-operator-cloud-credentials"))
		g.Expect(updatedRoleConfig.Status.OperatorRolesRef.ImageRegistryARN).To(Equal("arn:aws:iam::123456789012:role/test-openshift-image-registry-installer-cloud-credentials"))
		g.Expect(updatedRoleConfig.Status.OperatorRolesRef.StorageARN).To(Equal("arn:aws:iam::123456789012:role/test-openshift-cluster-csi-drivers-ebs-cloud-credentials"))
		g.Expect(updatedRoleConfig.Status.OperatorRolesRef.NetworkARN).To(Equal("arn:aws:iam::123456789012:role/test-openshift-cloud-network-config-controller-cloud-credentials"))
		g.Expect(updatedRoleConfig.Status.OperatorRolesRef.KubeCloudControllerARN).To(Equal("arn:aws:iam::123456789012:role/test-kube-system-kube-controller-manager"))
		g.Expect(updatedRoleConfig.Status.OperatorRolesRef.NodePoolManagementARN).To(Equal("arn:aws:iam::123456789012:role/test-kube-system-capa-controller-manager"))
		g.Expect(updatedRoleConfig.Status.OperatorRolesRef.ControlPlaneOperatorARN).To(Equal("arn:aws:iam::123456789012:role/test-kube-system-control-plane-operator"))
		g.Expect(updatedRoleConfig.Status.OperatorRolesRef.KMSProviderARN).To(Equal("arn:aws:iam::123456789012:role/test-kube-system-kms-provider"))

		// Should have a condition indicating success - expect Ready condition to be True
		readyCondition := v1beta1conditions.Get(updatedRoleConfig, expinfrav1.RosaRoleConfigReadyCondition)
		g.Expect(readyCondition).ToNot(BeNil())
		g.Expect(readyCondition.Status).To(Equal(corev1.ConditionTrue))
		g.Expect(readyCondition.Reason).To(Equal(expinfrav1.RosaRoleConfigCreatedReason))
	}).WithTimeout(30 * time.Second).WithPolling(500 * time.Millisecond).Should(Succeed())
}

// TestROSARoleConfigSetUpRuntimePerReconciliation verifies that setUpRuntime always invokes
// runtimeFactory on every call — i.e. no singleton caching across reconciliations.
// This is the regression test for the bug where the first reconciliation's AWSClient was
// reused for all subsequent ones regardless of identityRef.
//
// The test calls setUpRuntime directly with two distinct stub scopes so it doesn't depend
// on the reconciler's Get/scope-creation path, avoiding any environment-specific flakiness.
func TestROSARoleConfigSetUpRuntimePerReconciliation(t *testing.T) {
	g := NewWithT(t)
	ctx := context.TODO()

	callCount := 0
	var capturedScopes []*scope.RosaRoleConfigScope

	reconciler := &ROSARoleConfigReconciler{
		runtimeFactory: func(_ context.Context, s *scope.RosaRoleConfigScope) (*rosacli.Runtime, error) {
			callCount++
			capturedScopes = append(capturedScopes, s)
			return rosacli.NewRuntime(), nil
		},
	}

	scope1 := &scope.RosaRoleConfigScope{}
	scope2 := &scope.RosaRoleConfigScope{}

	_, err := reconciler.setUpRuntime(ctx, scope1)
	g.Expect(err).ToNot(HaveOccurred())

	_, err = reconciler.setUpRuntime(ctx, scope2)
	g.Expect(err).ToNot(HaveOccurred())

	// runtimeFactory must be called once per setUpRuntime call — no early-return caching.
	g.Expect(callCount).To(Equal(2), "runtimeFactory must be called for every reconciliation")

	// Each call received the scope it was given, not a shared/cached one.
	g.Expect(capturedScopes).To(HaveLen(2))
	g.Expect(capturedScopes[0]).To(BeIdenticalTo(scope1))
	g.Expect(capturedScopes[1]).To(BeIdenticalTo(scope2))
}

// TestROSARoleConfigSetUpRuntimeError verifies that when runtimeFactory fails, the error
// propagates and no stale runtime leaks into subsequent reconciliations.
func TestROSARoleConfigSetUpRuntimeError(t *testing.T) {
	RegisterTestingT(t)
	g := NewWithT(t)
	ctx := context.TODO()

	callCount := 0
	reconciler := &ROSARoleConfigReconciler{
		Client: testEnv.Client,
		runtimeFactory: func(c context.Context, s *scope.RosaRoleConfigScope) (*rosacli.Runtime, error) {
			callCount++
			return nil, fmt.Errorf("simulated AWS credential failure (call %d)", callCount)
		},
	}

	ns, err := testEnv.CreateNamespace(ctx, fmt.Sprintf("test-ns-err-%s", generateTestID()))
	g.Expect(err).ToNot(HaveOccurred())

	rc := &expinfrav1.ROSARoleConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("rc-err-%s", generateTestID()),
			Namespace: ns.Name,
		},
		Spec: expinfrav1.ROSARoleConfigSpec{
			AccountRoleConfig:  expinfrav1.AccountRoleConfig{Prefix: "err", Version: "4.15.0"},
			OperatorRoleConfig: expinfrav1.OperatorRoleConfig{Prefix: "err"},
			OidcProviderType:   expinfrav1.Managed,
		},
	}
	createObject(g, rc, ns.Name)
	defer cleanupObject(g, rc)

	req := ctrl.Request{NamespacedName: types.NamespacedName{Name: rc.Name, Namespace: rc.Namespace}}

	// First reconciliation: runtimeFactory fails.
	_, err1 := reconciler.Reconcile(ctx, req)
	g.Expect(err1).To(HaveOccurred())
	g.Expect(err1.Error()).To(ContainSubstring("failed to set up runtime"))
	g.Expect(callCount).To(Equal(1))

	// Second reconciliation: runtimeFactory is called again (no early-return / stale runtime).
	_, err2 := reconciler.Reconcile(ctx, req)
	g.Expect(err2).To(HaveOccurred())
	g.Expect(callCount).To(Equal(2), "runtimeFactory must be called on every reconciliation, even after a previous failure")
}

func TestROSARoleConfigReconcileDelete(t *testing.T) {
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
	mockAWSClient := aws.NewMockClient(mockCtrl)
	mockAWSClient.EXPECT().HasManagedPolicies(gomock.Any()).Return(false, nil).AnyTimes()
	mockAWSClient.EXPECT().HasHostedCPPolicies(gomock.Any()).Return(true, nil).AnyTimes()
	mockAWSClient.EXPECT().GetOperatorRolesFromAccountByPrefix(gomock.Any(), gomock.Any()).Return([]string{
		"test-openshift-ingress-operator-cloud-credentials",
		"test-openshift-image-registry-installer-cloud-credentials",
		"test-openshift-cluster-csi-drivers-ebs-cloud-credentials",
		"test-openshift-cloud-network-config-controller-cloud-credentials",
		"test-kube-system-kube-controller-manager",
		"test-kube-system-capa-controller-manager",
		"test-kube-system-control-plane-operator",
		"test-kube-system-kms-provider",
	}, nil).AnyTimes()

	// Delete operator roles (called individually for each role)
	mockAWSClient.EXPECT().DeleteOperatorRole(gomock.Any(), gomock.Any(), true).Return(map[string]bool{}, nil).AnyTimes()

	// Mock OIDC provider deletion
	mockAWSClient.EXPECT().DeleteOpenIDConnectProvider("arn:aws:iam::123456789012:oidc-provider/test-existing-oidc-id").Return(nil).AnyTimes()

	// Delete account roles (called individually for each role)
	mockAWSClient.EXPECT().DeleteAccountRole(gomock.Any(), gomock.Any(), true, false).Return(nil).AnyTimes()

	mockAWSClient.EXPECT().GetCreator().Return(&aws.Creator{
		ARN:       "arn:aws:iam::123456789012:user/test-user",
		AccountID: "123456789012",
		IsSTS:     false,
	}, nil).AnyTimes()

	rt := rosacli.NewRuntime()
	rt.OCMClient = ocmClient
	rt.AWSClient = mockAWSClient
	rt.Creator = &aws.Creator{
		ARN:       "arn:aws:iam::123456789012:user/test-user",
		AccountID: "123456789012",
		IsSTS:     false,
	}

	// Mock OCM API calls
	apiServer.AppendHandlers(
		ocmsdk.RespondWithJSON(
			http.StatusOK, "",
		),
	)
	// Mock GetOidcConfig call - return existing OIDC config
	apiServer.AppendHandlers(
		ocmsdk.RespondWithJSON(
			http.StatusOK, `{"id": "test-existing-oidc-id", "issuer_url": "https://test.existing.oidc.url"}`,
		),
	)
	// Mock GetAllClusters call
	apiServer.AppendHandlers(
		ocmsdk.RespondWithJSON(
			http.StatusOK, `{"items": []}`,
		),
	)
	// Mock GetAllCredRequests call
	apiServer.AppendHandlers(
		ocmsdk.RespondWithJSON(
			http.StatusOK, `{}`,
		),
	)
	// Mock HasAClusterUsingOperatorRolesPrefix call
	apiServer.AppendHandlers(
		ocmsdk.RespondWithJSON(
			http.StatusOK, `false`,
		),
	)

	// Mock existing OIDC config GET request
	apiServer.RouteToHandler("GET", "/api/clusters_mgmt/v1/oidc_configs/test-existing-oidc-id",
		ocmsdk.RespondWithJSON(
			http.StatusOK, `{"id": "test-existing-oidc-id", "issuer_url": "https://test.existing.oidc.url"}`,
		),
	)

	// Mock OIDC config deletion
	apiServer.RouteToHandler("DELETE", "/api/clusters_mgmt/v1/oidc_configs/test-existing-oidc-id",
		ocmsdk.RespondWithJSON(
			http.StatusOK, `{}`,
		),
	)

	// Create CRs with unique names to avoid conflicts
	ns, err := testEnv.CreateNamespace(ctx, fmt.Sprintf("test-namespace-delete-%s", testID))
	g.Expect(err).ToNot(HaveOccurred())

	// Create ROSARoleConfig with populated status (simulating existing resources)
	rosaRoleConfig := &expinfrav1.ROSARoleConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:       fmt.Sprintf("test-rosa-role-delete-%s", testID),
			Namespace:  ns.Name,
			Finalizers: []string{expinfrav1.RosaRoleConfigFinalizer},
			// Set deletion timestamp to simulate deletion request
			DeletionTimestamp: &metav1.Time{Time: time.Now()},
		},
		Spec: expinfrav1.ROSARoleConfigSpec{
			AccountRoleConfig: expinfrav1.AccountRoleConfig{
				Prefix:  "test",
				Version: "4.15.0",
			},
			OperatorRoleConfig: expinfrav1.OperatorRoleConfig{
				Prefix: "test",
			},
			OidcProviderType: expinfrav1.Managed,
		},
		Status: expinfrav1.ROSARoleConfigStatus{
			OIDCID:          "test-existing-oidc-id",
			OIDCProviderARN: "arn:aws:iam::123456789012:oidc-provider/test-existing-oidc-id",
			AccountRolesRef: expinfrav1.AccountRolesRef{
				InstallerRoleARN: "arn:aws:iam::123456789012:role/test-HCP-ROSA-Installer-Role",
				SupportRoleARN:   "arn:aws:iam::123456789012:role/test-HCP-ROSA-Support-Role",
				WorkerRoleARN:    "arn:aws:iam::123456789012:role/test-HCP-ROSA-Worker-Role",
			},
			OperatorRolesRef: rosacontrolplanev1.AWSRolesRef{
				IngressARN:              "arn:aws:iam::123456789012:role/test-openshift-ingress-operator-cloud-credentials",
				ImageRegistryARN:        "arn:aws:iam::123456789012:role/test-openshift-image-registry-installer-cloud-credentials",
				StorageARN:              "arn:aws:iam::123456789012:role/test-openshift-cluster-csi-drivers-ebs-cloud-credentials",
				NetworkARN:              "arn:aws:iam::123456789012:role/test-openshift-cloud-network-config-controller-cloud-credentials",
				KubeCloudControllerARN:  "arn:aws:iam::123456789012:role/test-kube-system-kube-controller-manager",
				NodePoolManagementARN:   "arn:aws:iam::123456789012:role/test-kube-system-capa-controller-manager",
				ControlPlaneOperatorARN: "arn:aws:iam::123456789012:role/test-kube-system-control-plane-operator",
				KMSProviderARN:          "arn:aws:iam::123456789012:role/test-kube-system-kms-provider",
			},
		},
	}

	createObject(g, rosaRoleConfig, ns.Name)
	defer cleanupObject(g, rosaRoleConfig)

	// Setup the reconciler using runtimeFactory to inject mock clients per reconciliation.
	reconciler := &ROSARoleConfigReconciler{
		Client:         testEnv.Client,
		runtimeFactory: makeRuntimeFactory(rt),
	}

	// Call the Reconcile function
	req := ctrl.Request{}
	req.NamespacedName = types.NamespacedName{Name: rosaRoleConfig.Name, Namespace: rosaRoleConfig.Namespace}

	err = reconciler.Client.Delete(ctx, rosaRoleConfig)
	g.Expect(err).ToNot(HaveOccurred())

	// Sleep to ensure the status is updated
	time.Sleep(100 * time.Millisecond)

	_, errReconcile := reconciler.Reconcile(ctx, req)

	// Assertions - deletion should succeed
	g.Expect(errReconcile).ToNot(HaveOccurred())

	// Sleep to ensure the status is updated
	time.Sleep(100 * time.Millisecond)

	deletedRoleConfig := &expinfrav1.ROSARoleConfig{}

	// Verify the resource has been deleted (finalizers removed)
	err = reconciler.Client.Get(ctx, req.NamespacedName, deletedRoleConfig)

	// The object should either be not found (fully deleted) or have no finalizers
	if err == nil {
		// If object still exists, verify finalizers are removed
		g.Expect(deletedRoleConfig.Finalizers).To(BeEmpty(), "Finalizers should be removed after successful deletion")
	}
}

// fakeSTSAssumeRoleResponse returns valid STS AssumeRole XML for the fake STS httptest server.
func fakeSTSAssumeRoleResponse(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/xml")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `<AssumeRoleResponse xmlns="https://sts.amazonaws.com/doc/2011-06-15/">
  <AssumeRoleResult>
    <Credentials>
      <AccessKeyId>ASIAIOSFODNN7EXAMPLE</AccessKeyId>
      <SecretAccessKey>wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY</SecretAccessKey>
      <SessionToken>AQoXnyc4MCrrlandlJKwBQ==</SessionToken>
      <Expiration>2030-01-01T00:00:00Z</Expiration>
    </Credentials>
    <AssumedRoleUser>
      <Arn>arn:aws:sts::123456789012:assumed-role/fake-rosa-role/test</Arn>
      <AssumedRoleId>ARO123EXAMPLE123:test</AssumedRoleId>
    </AssumedRoleUser>
  </AssumeRoleResult>
  <ResponseMetadata>
    <RequestId>12345678-1234-1234-1234-123456789012</RequestId>
  </ResponseMetadata>
</AssumeRoleResponse>`)
}

// TestROSARoleConfigReconcilerWithRoleIdentity verifies that ROSARoleConfigReconciler can
// create its scope (resolving credentials) when the ROSARoleConfig's IdentityRef points
// to an AWSClusterRoleIdentity backed by a fake IAM role.
//
// The fake STS server satisfies the AssumeRole call that happens during scope creation.
// The runtimeFactory is used as a canary: if it is called, scope creation (including
// credential resolution via AWSClusterRoleIdentity) succeeded.
func TestROSARoleConfigReconcilerWithRoleIdentity(t *testing.T) {
	RegisterTestingT(t)
	g := NewWithT(t)
	ctx := context.TODO()

	// Start a local STS server that returns fake AssumeRole credentials.
	stsServer := httptest.NewServer(http.HandlerFunc(fakeSTSAssumeRoleResponse))
	defer stsServer.Close()

	// Point the AWS SDK at the fake STS server and supply dummy source credentials so
	// config.LoadDefaultConfig has something to work with when building the STS client.
	t.Setenv("AWS_ENDPOINT_URL_STS", stsServer.URL)
	t.Setenv("AWS_ACCESS_KEY_ID", "fake-access-key-id")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "fake-secret-access-key")
	t.Setenv("AWS_REGION", "us-east-1")

	testID := generateTestID()

	ns, err := testEnv.CreateNamespace(ctx, fmt.Sprintf("test-ns-roleident-%s", testID))
	g.Expect(err).ToNot(HaveOccurred())

	// AWSClusterControllerIdentity is the source for the role identity.
	// The webhook requires sourceIdentityRef to be non-nil, so we use the controller
	// identity as a source — it returns no providers of its own, causing the role
	// provider to fall back to the env-var credentials when calling STS.
	controllerIdentity := &infrav1.AWSClusterControllerIdentity{
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
		},
		Spec: infrav1.AWSClusterControllerIdentitySpec{
			AWSClusterIdentitySpec: infrav1.AWSClusterIdentitySpec{
				AllowedNamespaces: &infrav1.AllowedNamespaces{},
			},
		},
	}
	createObject(g, controllerIdentity, ns.Name)
	defer cleanupObject(g, controllerIdentity)

	// AWSClusterRoleIdentity with a fake IAM role ARN — AllowedNamespaces is empty,
	// meaning every namespace may use this identity.
	roleIdentity := &infrav1.AWSClusterRoleIdentity{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("fake-role-identity-%s", testID),
		},
		Spec: infrav1.AWSClusterRoleIdentitySpec{
			AWSRoleSpec: infrav1.AWSRoleSpec{
				RoleArn:     fmt.Sprintf("arn:aws:iam::123456789012:role/fake-rosa-role-%s", testID),
				SessionName: "test-session",
			},
			AWSClusterIdentitySpec: infrav1.AWSClusterIdentitySpec{
				AllowedNamespaces: &infrav1.AllowedNamespaces{},
			},
			SourceIdentityRef: &infrav1.AWSIdentityReference{
				Name: controllerIdentity.Name,
				Kind: infrav1.ControllerIdentityKind,
			},
		},
	}
	roleIdentity.SetGroupVersionKind(infrav1.GroupVersion.WithKind("AWSClusterRoleIdentity"))
	createObject(g, roleIdentity, ns.Name)
	defer cleanupObject(g, roleIdentity)

	rosaRoleConfig := &expinfrav1.ROSARoleConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:       fmt.Sprintf("test-role-config-%s", testID),
			Namespace:  ns.Name,
			Finalizers: []string{expinfrav1.RosaRoleConfigFinalizer},
		},
		Spec: expinfrav1.ROSARoleConfigSpec{
			AccountRoleConfig: expinfrav1.AccountRoleConfig{
				Prefix:  "test",
				Version: "4.15.0",
			},
			OperatorRoleConfig: expinfrav1.OperatorRoleConfig{
				Prefix: "test",
			},
			OidcProviderType: expinfrav1.Managed,
			IdentityRef: &infrav1.AWSIdentityReference{
				Name: roleIdentity.Name,
				Kind: infrav1.ClusterRoleIdentityKind,
			},
		},
	}
	createObject(g, rosaRoleConfig, ns.Name)
	defer cleanupObject(g, rosaRoleConfig)

	runtimeCalled := false
	reconciler := &ROSARoleConfigReconciler{
		Client: testEnv.Client,
		// runtimeFactory acts as a canary: if it is invoked, scope creation (and
		// therefore AWSClusterRoleIdentity credential resolution) succeeded.
		runtimeFactory: func(_ context.Context, _ *scope.RosaRoleConfigScope) (*rosacli.Runtime, error) {
			runtimeCalled = true
			return nil, fmt.Errorf("sentinel: runtimeFactory reached after role-identity scope creation")
		},
	}

	req := ctrl.Request{
		NamespacedName: types.NamespacedName{Name: rosaRoleConfig.Name, Namespace: rosaRoleConfig.Namespace},
	}

	g.Eventually(func(g Gomega) {
		_, errReconcile := reconciler.Reconcile(ctx, req)

		// The reconciler must have advanced past scope creation and reached runtimeFactory.
		// A "failed to create rosaroleconfig scope" error here would mean credential
		// resolution via AWSClusterRoleIdentity failed.
		g.Expect(errReconcile).To(HaveOccurred())
		g.Expect(errReconcile.Error()).To(ContainSubstring("failed to set up runtime"))
		g.Expect(errReconcile.Error()).NotTo(ContainSubstring("failed to create rosaroleconfig scope"))
		g.Expect(runtimeCalled).To(BeTrue())
	}).WithTimeout(30 * time.Second).WithPolling(500 * time.Millisecond).Should(Succeed())
}

// TestROSARoleConfigReconcilerWithRoleIdentityNamespaceNotAllowed verifies that when the
// AWSClusterRoleIdentity's AllowedNamespaces does not include the ROSARoleConfig's namespace,
// scope creation fails with a credential error and the runtimeFactory is never called.
func TestROSARoleConfigReconcilerWithRoleIdentityNamespaceNotAllowed(t *testing.T) {
	RegisterTestingT(t)
	g := NewWithT(t)
	ctx := context.TODO()

	stsServer := httptest.NewServer(http.HandlerFunc(fakeSTSAssumeRoleResponse))
	defer stsServer.Close()

	t.Setenv("AWS_ENDPOINT_URL_STS", stsServer.URL)
	t.Setenv("AWS_ACCESS_KEY_ID", "fake-access-key-id")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "fake-secret-access-key")
	t.Setenv("AWS_REGION", "us-east-1")

	testID := generateTestID()

	ns, err := testEnv.CreateNamespace(ctx, fmt.Sprintf("test-ns-roleident-denied-%s", testID))
	g.Expect(err).ToNot(HaveOccurred())

	controllerIdentity := &infrav1.AWSClusterControllerIdentity{
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
		},
		Spec: infrav1.AWSClusterControllerIdentitySpec{
			AWSClusterIdentitySpec: infrav1.AWSClusterIdentitySpec{
				AllowedNamespaces: &infrav1.AllowedNamespaces{},
			},
		},
	}
	createObject(g, controllerIdentity, ns.Name)
	defer cleanupObject(g, controllerIdentity)

	// AllowedNamespaces lists only "other-namespace", so ns.Name is not permitted.
	roleIdentity := &infrav1.AWSClusterRoleIdentity{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("restricted-role-identity-%s", testID),
		},
		Spec: infrav1.AWSClusterRoleIdentitySpec{
			AWSRoleSpec: infrav1.AWSRoleSpec{
				RoleArn:     fmt.Sprintf("arn:aws:iam::123456789012:role/restricted-rosa-role-%s", testID),
				SessionName: "test-session",
			},
			AWSClusterIdentitySpec: infrav1.AWSClusterIdentitySpec{
				AllowedNamespaces: &infrav1.AllowedNamespaces{
					NamespaceList: []string{"other-namespace"},
				},
			},
			SourceIdentityRef: &infrav1.AWSIdentityReference{
				Name: controllerIdentity.Name,
				Kind: infrav1.ControllerIdentityKind,
			},
		},
	}
	roleIdentity.SetGroupVersionKind(infrav1.GroupVersion.WithKind("AWSClusterRoleIdentity"))
	createObject(g, roleIdentity, ns.Name)
	defer cleanupObject(g, roleIdentity)

	rosaRoleConfig := &expinfrav1.ROSARoleConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:       fmt.Sprintf("test-role-config-denied-%s", testID),
			Namespace:  ns.Name,
			Finalizers: []string{expinfrav1.RosaRoleConfigFinalizer},
		},
		Spec: expinfrav1.ROSARoleConfigSpec{
			AccountRoleConfig: expinfrav1.AccountRoleConfig{
				Prefix:  "test",
				Version: "4.15.0",
			},
			OperatorRoleConfig: expinfrav1.OperatorRoleConfig{Prefix: "test"},
			OidcProviderType:   expinfrav1.Managed,
			IdentityRef: &infrav1.AWSIdentityReference{
				Name: roleIdentity.Name,
				Kind: infrav1.ClusterRoleIdentityKind,
			},
		},
	}
	createObject(g, rosaRoleConfig, ns.Name)
	defer cleanupObject(g, rosaRoleConfig)

	runtimeCalled := false
	reconciler := &ROSARoleConfigReconciler{
		Client: testEnv.Client,
		runtimeFactory: func(_ context.Context, _ *scope.RosaRoleConfigScope) (*rosacli.Runtime, error) {
			runtimeCalled = true
			return nil, nil
		},
	}

	req := ctrl.Request{
		NamespacedName: types.NamespacedName{Name: rosaRoleConfig.Name, Namespace: rosaRoleConfig.Namespace},
	}

	_, errReconcile := reconciler.Reconcile(ctx, req)

	// Scope creation must fail because the namespace is not in AllowedNamespaces.
	g.Expect(errReconcile).To(HaveOccurred())
	g.Expect(errReconcile.Error()).To(ContainSubstring("failed to create rosaroleconfig scope"))
	g.Expect(runtimeCalled).To(BeFalse())
}
