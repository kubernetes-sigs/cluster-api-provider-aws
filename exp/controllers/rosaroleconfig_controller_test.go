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
	"net/http"
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

	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
)

func TestROSARoleConfigReconcile(t *testing.T) {
	RegisterTestingT(t)
	g := NewWithT(t)

	ssoServer := ocmsdk.MakeTCPServer()
	apiServer := ocmsdk.MakeTCPServer()
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

	r := rosacli.NewRuntime()
	r.OCMClient = ocmClient
	r.AWSClient = awsClient
	r.Creator = &aws.Creator{
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

	// mock ocm API calls - first call gets tris response
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
	// Mock GetAllClusters call
	apiServer.AppendHandlers(
		ocmsdk.RespondWithJSON(
			http.StatusOK, `{"items": []}`,
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

	// prepare the role config

	// Create CRs
	ns, err := testEnv.CreateNamespace(ctx, "test-namespace")
	rosaRoleConfig := &expinfrav1.ROSARoleConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:       "test-rosa-role",
			Namespace:  ns.Name,
			Finalizers: []string{expinfrav1.RosaRoleConfigFinalizer},
		},
		Spec: expinfrav1.ROSARoleConfigSpec{
			AccountRoleConfig: expinfrav1.AccountRoleConfig{
				Prefix:  "test",
				Version: "4.15",
			},
			OperatorRoleConfig: expinfrav1.OperatorRoleConfig{
				Prefix: "test",
			},
		},
	}
	g.Expect(err).ToNot(HaveOccurred())

	createObject(g, rosaRoleConfig, ns.Name)
	defer cleanupObject(g, rosaRoleConfig)

	// Setup the reconciler with these mocks
	reconciler := &ROSARoleConfigReconciler{
		Client:  testEnv.Client,
		Runtime: r,
	}

	// Call the Reconcile function
	req := ctrl.Request{}
	req.NamespacedName = types.NamespacedName{Name: rosaRoleConfig.Name, Namespace: rosaRoleConfig.Namespace}
	_, errReconcile := reconciler.Reconcile(ctx, req)

	// Assertions - expect the installer role empty error since AccountRolesRef is not populated yet
	g.Expect(errReconcile).To(HaveOccurred())
	g.Expect(errReconcile.Error()).To(ContainSubstring("installer role is empty"))

	// Check the status of the ROSARoleConfig resource
	updatedRoleConfig := &expinfrav1.ROSARoleConfig{}
	err = reconciler.Client.Get(ctx, req.NamespacedName, updatedRoleConfig)
	g.Expect(err).ToNot(HaveOccurred())

	// Verify that account roles would be populated if the mocking was complete
	// (The actual population depends on the CreateAccountRoles implementation)
	// Verify conditions are set appropriately

	// Should have a condition indicating the failure at operator role creation
	hasFailureCondition := false
	for _, condition := range updatedRoleConfig.Status.Conditions {
		if condition.Type == expinfrav1.RosaRoleConfigReadyCondition &&
			condition.Status == corev1.ConditionFalse &&
			condition.Reason == expinfrav1.RosaRoleConfigReconciliationFailedReason {
			hasFailureCondition = true
			g.Expect(condition.Message).To(ContainSubstring("Failed to create Operator Roles"))
			break
		}
	}
	g.Expect(hasFailureCondition).To(BeTrue(), "Expected to find a failure condition for operator role creation")
}
