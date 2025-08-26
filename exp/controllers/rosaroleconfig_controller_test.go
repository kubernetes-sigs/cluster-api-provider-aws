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
	"go.uber.org/mock/gomock"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"

	rosacontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api/util/conditions"
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
	mockIamClient.EXPECT().ListRoles(gomock.Any(), gomock.Any(), gomock.Any()).Return(&iamv2.ListRolesOutput{
		Roles: []iamTypes.Role{},
	}, nil).AnyTimes()
	mockIamClient.EXPECT().ListOpenIDConnectProviders(gomock.Any(), gomock.Any(), gomock.Any()).Return(&iamv2.ListOpenIDConnectProvidersOutput{
		OpenIDConnectProviderList: []iamTypes.OpenIDConnectProviderListEntry{},
	}, nil).AnyTimes()
	// Mock ListRoleTags calls
	mockIamClient.EXPECT().ListRoleTags(gomock.Any(), gomock.Any(), gomock.Any()).Return(&iamv2.ListRoleTagsOutput{
		Tags: []iamTypes.Tag{},
	}, nil).AnyTimes()

	// Mock ListAttachedRolePolicies calls - return empty policies since roles don't exist yet
	mockIamClient.EXPECT().ListAttachedRolePolicies(gomock.Any(), gomock.Any(), gomock.Any()).Return(&iamv2.ListAttachedRolePoliciesOutput{
		AttachedPolicies: []iamTypes.AttachedPolicy{},
	}, nil).AnyTimes()

	// Mock GetRole calls - return role not found error to trigger role creation
	mockIamClient.EXPECT().GetRole(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, &iamTypes.NoSuchEntityException{
		Message: awsSdk.String("The role with name test-role does not exist."),
	}).AnyTimes()
	// Mock CreateRole calls for role creation
	mockIamClient.EXPECT().CreateRole(gomock.Any(), gomock.Any(), gomock.Any()).Return(&iamv2.CreateRoleOutput{
		Role: &iamTypes.Role{
			RoleName: awsSdk.String("test-role"),
			Arn:      awsSdk.String("arn:aws:iam::123456789012:role/test-role"),
		},
	}, nil).AnyTimes()
	// Mock AttachRolePolicy calls
	mockIamClient.EXPECT().AttachRolePolicy(gomock.Any(), gomock.Any(), gomock.Any()).Return(&iamv2.AttachRolePolicyOutput{}, nil).AnyTimes()
	// Mock CreatePolicy calls
	mockIamClient.EXPECT().CreatePolicy(gomock.Any(), gomock.Any(), gomock.Any()).Return(&iamv2.CreatePolicyOutput{
		Policy: &iamTypes.Policy{
			PolicyName: awsSdk.String("test-policy"),
			Arn:        awsSdk.String("arn:aws:iam::123456789012:policy/test-policy"),
		},
	}, nil).AnyTimes()
	// Mock GetPolicy calls - return success for AWS managed policies, not found for others
	mockIamClient.EXPECT().GetPolicy(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, input *iamv2.GetPolicyInput, optFns ...func(*iamv2.Options)) (*iamv2.GetPolicyOutput, error) {
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
	mockIamClient.EXPECT().ListPolicies(gomock.Any(), gomock.Any(), gomock.Any()).Return(&iamv2.ListPoliciesOutput{
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
		Arn:     awsSdk.String("arn:aws:iam::123456789012:user/test-user"),
		Account: awsSdk.String("123456789012"),
		UserId:  awsSdk.String("test-user-id"),
	}, nil).AnyTimes()

	// Create mock AWS client and set up expectations
	mockAWSClient := aws.NewMockClient(mockCtrl)

	// Mock HasHostedCPPolicies to always return true and no error
	mockAWSClient.EXPECT().HasHostedCPPolicies(gomock.Any()).Return(true, nil).AnyTimes()

	// Mock other AWS client methods that might be called
	mockAWSClient.EXPECT().ListAccountRoles(gomock.Any()).Return([]aws.Role{}, nil).AnyTimes()
	mockAWSClient.EXPECT().ListOperatorRoles(gomock.Any(), gomock.Any(), gomock.Any()).Return(map[string][]aws.Role{}, nil).AnyTimes()
	mockAWSClient.EXPECT().GetCreator().Return(&aws.Creator{
		ARN:       "arn:aws:iam::123456789012:user/test-user",
		AccountID: "123456789012",
		IsSTS:     false,
	}, nil).AnyTimes()

	awsClient := mockAWSClient

	r := rosacli.NewRuntime()
	r.OCMClient = ocmClient
	r.AWSClient = awsClient
	r.Creator = &aws.Creator{
		ARN:       "arn:aws:iam::123456789012:user/test-user",
		AccountID: "123456789012",
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

	// Should have a condition indicating the failure at operator role creation
	failureCondition := conditions.Get(updatedRoleConfig, expinfrav1.RosaRoleConfigReadyCondition)
	g.Expect(failureCondition).ToNot(BeNil())
	g.Expect(failureCondition.Status).To(Equal(corev1.ConditionFalse))
	g.Expect(failureCondition.Reason).To(Equal(expinfrav1.RosaRoleConfigReconciliationFailedReason))
	g.Expect(failureCondition.Message).To(ContainSubstring("Failed to create Operator Roles"))
}

func TestROSARoleConfigReconcileWithExistingAccountRoles(t *testing.T) {
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
	// mock iam client to expect ListRoles call - return existing account roles
	mockIamClient := rosaMocks.NewMockIamApiClient(mockCtrl)
	mockIamClient.EXPECT().ListRoles(gomock.Any(), gomock.Any(), gomock.Any()).Return(&iamv2.ListRolesOutput{
		Roles: []iamTypes.Role{
			{
				RoleName: awsSdk.String("test-HCP-ROSA-Installer-Role"),
				Arn:      awsSdk.String("arn:aws:iam::123456789012:role/test-HCP-ROSA-Installer-Role"),
			},
			{
				RoleName: awsSdk.String("test-HCP-ROSA-Support-Role"),
				Arn:      awsSdk.String("arn:aws:iam::123456789012:role/test-HCP-ROSA-Support-Role"),
			},
			{
				RoleName: awsSdk.String("test-HCP-ROSA-Worker-Role"),
				Arn:      awsSdk.String("arn:aws:iam::123456789012:role/test-HCP-ROSA-Worker-Role"),
			},
		},
	}, nil).AnyTimes()

	mockIamClient.EXPECT().ListOpenIDConnectProviders(gomock.Any(), gomock.Any(), gomock.Any()).Return(&iamv2.ListOpenIDConnectProvidersOutput{
		OpenIDConnectProviderList: []iamTypes.OpenIDConnectProviderListEntry{
			{
				Arn: awsSdk.String("arn:aws:iam::123456789012:oidc-provider/test-oidc-id-created"),
			},
		},
	}, nil).AnyTimes()

	// Mock ListRoleTags calls
	mockIamClient.EXPECT().ListRoleTags(gomock.Any(), gomock.Any(), gomock.Any()).Return(&iamv2.ListRoleTagsOutput{
		Tags: []iamTypes.Tag{},
	}, nil).AnyTimes()

	// Mock ListOpenIDConnectProviderTags calls
	mockIamClient.EXPECT().ListOpenIDConnectProviderTags(gomock.Any(), gomock.Any(), gomock.Any()).Return(&iamv2.ListOpenIDConnectProviderTagsOutput{
		Tags: []iamTypes.Tag{},
	}, nil).AnyTimes()

	// Mock ListAttachedRolePolicies calls - return expected managed policies for installer role
	mockIamClient.EXPECT().ListAttachedRolePolicies(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, input *iamv2.ListAttachedRolePoliciesInput, optFns ...func(*iamv2.Options)) (*iamv2.ListAttachedRolePoliciesOutput, error) {
		roleName := *input.RoleName
		if strings.Contains(roleName, "Installer-Role") {
			return &iamv2.ListAttachedRolePoliciesOutput{
				AttachedPolicies: []iamTypes.AttachedPolicy{
					{
						PolicyName: awsSdk.String("sts_hcp_installer_permission_policy"),
						PolicyArn:  awsSdk.String("arn:aws:iam::aws:policy/sts_hcp_installer_permission_policy"),
					},
				},
			}, nil
		}
		// For other roles, return empty policies
		return &iamv2.ListAttachedRolePoliciesOutput{
			AttachedPolicies: []iamTypes.AttachedPolicy{},
		}, nil
	}).AnyTimes()

	// Mock CreateOpenIDConnectProvider calls for OIDC provider creation
	mockIamClient.EXPECT().CreateOpenIDConnectProvider(gomock.Any(), gomock.Any(), gomock.Any()).Return(&iamv2.CreateOpenIDConnectProviderOutput{
		OpenIDConnectProviderArn: awsSdk.String("arn:aws:iam::123456789012:oidc-provider/test-oidc-id-created"),
	}, nil).AnyTimes()

	// Mock GetRole calls - return existing roles for account roles
	mockIamClient.EXPECT().GetRole(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, input *iamv2.GetRoleInput, optFns ...func(*iamv2.Options)) (*iamv2.GetRoleOutput, error) {
		roleName := *input.RoleName
		switch {
		case strings.Contains(roleName, "Installer-Role"):
			return &iamv2.GetRoleOutput{
				Role: &iamTypes.Role{
					RoleName: awsSdk.String("test-HCP-ROSA-Installer-Role"),
					Arn:      awsSdk.String("arn:aws:iam::123456789012:role/test-HCP-ROSA-Installer-Role"),
				},
			}, nil
		case strings.Contains(roleName, "Support-Role"):
			return &iamv2.GetRoleOutput{
				Role: &iamTypes.Role{
					RoleName: awsSdk.String("test-HCP-ROSA-Support-Role"),
					Arn:      awsSdk.String("arn:aws:iam::123456789012:role/test-HCP-ROSA-Support-Role"),
				},
			}, nil
		case strings.Contains(roleName, "Worker-Role"):
			return &iamv2.GetRoleOutput{
				Role: &iamTypes.Role{
					RoleName: awsSdk.String("test-HCP-ROSA-Worker-Role"),
					Arn:      awsSdk.String("arn:aws:iam::123456789012:role/test-HCP-ROSA-Worker-Role"),
				},
			}, nil
		default:
			// For operator roles, return not found to trigger creation
			return nil, &iamTypes.NoSuchEntityException{
				Message: awsSdk.String("The role with name " + roleName + " does not exist."),
			}
		}
	}).AnyTimes()

	// Mock CreateRole calls for operator role creation (account roles already exist)
	mockIamClient.EXPECT().CreateRole(gomock.Any(), gomock.Any(), gomock.Any()).Return(&iamv2.CreateRoleOutput{
		Role: &iamTypes.Role{
			RoleName: awsSdk.String("test-operator-role"),
			Arn:      awsSdk.String("arn:aws:iam::123456789012:role/test-operator-role"),
		},
	}, nil).AnyTimes()

	// Mock AttachRolePolicy calls
	mockIamClient.EXPECT().AttachRolePolicy(gomock.Any(), gomock.Any(), gomock.Any()).Return(&iamv2.AttachRolePolicyOutput{}, nil).AnyTimes()

	// Mock CreatePolicy calls
	mockIamClient.EXPECT().CreatePolicy(gomock.Any(), gomock.Any(), gomock.Any()).Return(&iamv2.CreatePolicyOutput{
		Policy: &iamTypes.Policy{
			PolicyName: awsSdk.String("test-policy"),
			Arn:        awsSdk.String("arn:aws:iam::123456789012:policy/test-policy"),
		},
	}, nil).AnyTimes()

	// Mock GetPolicy calls - return success for AWS managed policies, not found for others
	mockIamClient.EXPECT().GetPolicy(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, input *iamv2.GetPolicyInput, optFns ...func(*iamv2.Options)) (*iamv2.GetPolicyOutput, error) {
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
		case "arn:aws:iam::123456789012:policy/test-HCP-ROSA-Installer-Policy":
			return &iamv2.GetPolicyOutput{
				Policy: &iamTypes.Policy{
					PolicyName: awsSdk.String("test-HCP-ROSA-Installer-Policy"),
					Arn:        awsSdk.String("arn:aws:iam::123456789012:policy/test-HCP-ROSA-Installer-Policy"),
				},
			}, nil
		case "arn:aws:iam::123456789012:role/test-HCP-ROSA-Installer-Role":
			return &iamv2.GetPolicyOutput{
				Policy: &iamTypes.Policy{
					PolicyName: awsSdk.String("test-HCP-ROSA-Installer-Policy"),
					Arn:        awsSdk.String("arn:aws:iam::123456789012:policy/test-HCP-ROSA-Installer-Policy"),
				},
			}, nil
		default:
			return nil, &iamTypes.NoSuchEntityException{
				Message: awsSdk.String("The policy does not exist."),
			}
		}
	}).AnyTimes()

	// Mock ListPolicies calls - return expected ROSA managed policies
	mockIamClient.EXPECT().ListPolicies(gomock.Any(), gomock.Any(), gomock.Any()).Return(&iamv2.ListPoliciesOutput{
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
		Arn:     awsSdk.String("arn:aws:iam::123456789012:user/test-user"),
		Account: awsSdk.String("123456789012"),
		UserId:  awsSdk.String("test-user-id"),
	}, nil).AnyTimes()

	// Create mock AWS client and set up expectations for second test
	mockAWSClient2 := aws.NewMockClient(mockCtrl)

	// Mock HasHostedCPPolicies to always return true and no error
	mockAWSClient2.EXPECT().HasHostedCPPolicies(gomock.Any()).Return(true, nil).AnyTimes()

	// Mock other AWS client methods that might be called in the second test
	mockAWSClient2.EXPECT().ListAccountRoles(gomock.Any()).Return([]aws.Role{
		{
			RoleName: "test-HCP-ROSA-Installer-Role",
			RoleARN:  "arn:aws:iam::123456789012:role/test-HCP-ROSA-Installer-Role",
		},
		{
			RoleName: "test-HCP-ROSA-Support-Role",
			RoleARN:  "arn:aws:iam::123456789012:role/test-HCP-ROSA-Support-Role",
		},
		{
			RoleName: "test-HCP-ROSA-Worker-Role",
			RoleARN:  "arn:aws:iam::123456789012:role/test-HCP-ROSA-Worker-Role",
		},
	}, nil).AnyTimes()
	mockAWSClient2.EXPECT().ListOperatorRoles(gomock.Any(), gomock.Any(), gomock.Any()).Return(map[string][]aws.Role{}, nil).AnyTimes()
	mockAWSClient2.EXPECT().ListOidcProviders(gomock.Any(), gomock.Any()).Return([]aws.OidcProviderOutput{
		{
			Arn: "arn:aws:iam::123456789012:oidc-provider/test-oidc-id-created",
		},
	}, nil).AnyTimes()
	mockAWSClient2.EXPECT().GetCreator().Return(&aws.Creator{
		ARN:       "arn:aws:iam::123456789012:user/test-user",
		AccountID: "123456789012",
		IsSTS:     false,
	}, nil).AnyTimes()

	awsClient := mockAWSClient2

	r := rosacli.NewRuntime()
	r.OCMClient = ocmClient
	r.AWSClient = awsClient
	r.Creator = &aws.Creator{
		ARN:       "arn:aws:iam::123456789012:user/test-user",
		AccountID: "123456789012",
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
						},
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

	// Create CRs
	ns, err := testEnv.CreateNamespace(ctx, "test-namespace-existing-roles")
	g.Expect(err).ToNot(HaveOccurred())

	rosaRoleConfig := &expinfrav1.ROSARoleConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:       "test-rosa-role-existing",
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
		// Status will be set separately after object creation
	}

	createObject(g, rosaRoleConfig, ns.Name)
	defer cleanupObject(g, rosaRoleConfig)

	// Update the status separately since status is a subresource in Kubernetes
	// First, get the created object from the API server
	err = testEnv.Client.Get(ctx, types.NamespacedName{Name: rosaRoleConfig.Name, Namespace: rosaRoleConfig.Namespace}, rosaRoleConfig)
	g.Expect(err).ToNot(HaveOccurred())

	// Set the status with pre-existing roles
	rosaRoleConfig.Status = expinfrav1.ROSARoleConfigStatus{
		// Pre-populate AccountRolesRef since account roles already exist
		AccountRolesRef: expinfrav1.AccountRolesRef{
			InstallerRoleARN: "arn:aws:iam::123456789012:role/test-HCP-ROSA-Installer-Role",
			SupportRoleARN:   "arn:aws:iam::123456789012:role/test-HCP-ROSA-Support-Role",
			WorkerRoleARN:    "arn:aws:iam::123456789012:role/test-HCP-ROSA-Worker-Role",
		},
		// OIDC config and provider will be created during reconciliation
		// Pre-populate operator roles since they already exist
		OperatorRolesRef: rosacontrolplanev1.AWSRolesRef{
			IngressARN:              "arn:aws:iam::123456789012:role/test-ingress-role",
			ImageRegistryARN:        "arn:aws:iam::123456789012:role/test-image-registry-role",
			StorageARN:              "arn:aws:iam::123456789012:role/test-storage-role",
			NetworkARN:              "arn:aws:iam::123456789012:role/test-network-role",
			KubeCloudControllerARN:  "arn:aws:iam::123456789012:role/test-kube-controller-role",
			NodePoolManagementARN:   "arn:aws:iam::123456789012:role/test-nodepool-role",
			ControlPlaneOperatorARN: "arn:aws:iam::123456789012:role/test-controlplane-role",
			KMSProviderARN:          "arn:aws:iam::123456789012:role/test-kms-role",
		},
	}

	// Update the status using the status subresource
	err = testEnv.Client.Status().Update(ctx, rosaRoleConfig)
	g.Expect(err).ToNot(HaveOccurred())

	// Setup the reconciler with these mocks
	reconciler := &ROSARoleConfigReconciler{
		Client:  testEnv.Client,
		Runtime: r,
	}

	// Call the Reconcile function
	req := ctrl.Request{}
	req.NamespacedName = types.NamespacedName{Name: rosaRoleConfig.Name, Namespace: rosaRoleConfig.Namespace}
	_, errReconcile := reconciler.Reconcile(ctx, req)

	// Assertions - since account roles, OIDC config, and operator roles already exist,
	// reconciliation should succeed
	g.Expect(errReconcile).ToNot(HaveOccurred())

	// Check the status of the ROSARoleConfig resource
	updatedRoleConfig := &expinfrav1.ROSARoleConfig{}
	err = reconciler.Client.Get(ctx, req.NamespacedName, updatedRoleConfig)
	g.Expect(err).ToNot(HaveOccurred())

	// Verify that all existing data is preserved
	g.Expect(updatedRoleConfig.Status.AccountRolesRef.InstallerRoleARN).ToNot(BeEmpty())
	g.Expect(updatedRoleConfig.Status.AccountRolesRef.InstallerRoleARN).To(Equal("arn:aws:iam::123456789012:role/test-HCP-ROSA-Installer-Role"))
	g.Expect(updatedRoleConfig.Status.AccountRolesRef.SupportRoleARN).To(Equal("arn:aws:iam::123456789012:role/test-HCP-ROSA-Support-Role"))
	g.Expect(updatedRoleConfig.Status.AccountRolesRef.WorkerRoleARN).To(Equal("arn:aws:iam::123456789012:role/test-HCP-ROSA-Worker-Role"))

	// Verify OIDC config was created during reconciliation
	g.Expect(updatedRoleConfig.Status.OIDCID).ToNot(BeEmpty())
	g.Expect(updatedRoleConfig.Status.OIDCID).To(Equal("test-oidc-id-created"))
	g.Expect(updatedRoleConfig.Status.OIDCProviderARN).ToNot(BeEmpty())
	// The provider ARN should contain the OIDC ID
	g.Expect(updatedRoleConfig.Status.OIDCProviderARN).To(ContainSubstring("test-oidc-id-created"))

	// Verify operator roles are preserved
	g.Expect(updatedRoleConfig.Status.OperatorRolesRef.IngressARN).To(Equal("arn:aws:iam::123456789012:role/test-ingress-role"))
	g.Expect(updatedRoleConfig.Status.OperatorRolesRef.ImageRegistryARN).To(Equal("arn:aws:iam::123456789012:role/test-image-registry-role"))

	// Should have a condition indicating success
	readyCondition := conditions.Get(updatedRoleConfig, expinfrav1.RosaRoleConfigReadyCondition)
	g.Expect(readyCondition).ToNot(BeNil())
	g.Expect(readyCondition.Status).To(Equal(corev1.ConditionTrue))
	g.Expect(readyCondition.Reason).To(Equal(expinfrav1.RosaRoleConfigCreatedReason))
}
