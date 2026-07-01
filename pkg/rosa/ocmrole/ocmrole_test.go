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

package ocmrole

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	sdk "github.com/openshift-online/ocm-sdk-go"
	"github.com/openshift-online/ocm-sdk-go/logging"
	. "github.com/openshift-online/ocm-sdk-go/testing"
	rosaaws "github.com/openshift/rosa/pkg/aws"
	"github.com/openshift/rosa/pkg/config"
	"github.com/openshift/rosa/pkg/ocm"
	"github.com/openshift/rosa/pkg/reporter"
	rosacli "github.com/openshift/rosa/pkg/rosa"
	"go.uber.org/mock/gomock"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
)

var (
	r         *rosacli.Runtime
	ocmClient *ocm.Client
	awsClient *rosaaws.MockClient
	ctrl      *gomock.Controller
	testLog   *logger.Logger

	testAccountID = "111111111111"
)

func TestOCMRolePackage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "OCM Role Package Suite")
}

var _ = Describe("GetOrCreateOCMRole", func() {
	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		awsClient = rosaaws.NewMockClient(ctrl)
		testLog = logger.NewLogger(GinkgoLogr)

		logger, err := logging.NewGoLoggerBuilder().
			Debug(false).
			Build()
		Expect(err).To(BeNil())

		connection, err := sdk.NewConnectionBuilder().
			Logger(logger).
			Tokens("test-token").
			URL("http://fake.api").
			Build()
		Expect(err).To(BeNil())
		ocmClient = ocm.NewClientWithConnection(connection)

		r = &rosacli.Runtime{
			Reporter:  reporter.CreateReporter(),
			OCMClient: ocmClient,
			AWSClient: awsClient,
		}
		r.Creator = &rosaaws.Creator{
			AccountID: testAccountID,
			ARN:       fmt.Sprintf("arn:aws:iam::%s:user/test", testAccountID),
			Partition: "aws",
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("Input validation", func() {
		It("should fail when runtime is nil", func() {
			_, _, _, err := GetOrCreateOCMRole(nil, testLog, "test-prefix", ProfileStandard, "", "/")

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("runtime cannot be nil"))
		})

		It("should fail when AWS client is nil", func() {
			r.AWSClient = nil

			_, _, _, err := GetOrCreateOCMRole(r, testLog, "test-prefix", ProfileStandard, "", "/")

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("AWS client cannot be nil"))
		})

		It("should fail when creator is nil", func() {
			r.Creator = nil

			_, _, _, err := GetOrCreateOCMRole(r, testLog, "test-prefix", ProfileStandard, "", "/")

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("creator cannot be nil"))
		})

		It("should fail when reporter is nil", func() {
			r.Reporter = nil

			_, _, _, err := GetOrCreateOCMRole(r, testLog, "test-prefix", ProfileStandard, "", "/")

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("reporter cannot be nil"))
		})

		It("should fail when OCM client is nil", func() {
			r.OCMClient = nil

			_, _, _, err := GetOrCreateOCMRole(r, testLog, "test-prefix", ProfileStandard, "", "/")

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("OCM client cannot be nil"))
		})

		It("should fail when logger is nil", func() {
			_, _, _, err := GetOrCreateOCMRole(r, nil, "test-prefix", ProfileStandard, "", "/")

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("logger cannot be nil"))
		})

		It("should fail when prefix is empty", func() {
			_, _, _, err := GetOrCreateOCMRole(r, testLog, "", ProfileStandard, "", "/")

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("prefix cannot be empty"))
		})

		It("should fail when profile is invalid", func() {
			_, _, _, err := GetOrCreateOCMRole(r, testLog, "test-prefix", "InvalidProfile", "", "/")

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("profile must be one of"))
		})
	})

	Context("Role creation and idempotency", func() {
		var ssoServer, apiServer *ghttp.Server
		var testRuntime *rosacli.Runtime
		var tmpdir string

		BeforeEach(func() {
			var err error

			// Create temp directory for OCM config
			tmpdir, err = os.MkdirTemp("", ".ocm-config-*")
			Expect(err).ToNot(HaveOccurred())
			os.Setenv("OCM_CONFIG", tmpdir+"/ocm_config.json")

			// Create mock OCM servers
			ssoServer = MakeTCPServer()
			apiServer = MakeTCPServer()

			ssoServer.AppendHandlers(
				RespondWithAccessAndRefreshTokens(
					MakeTokenString("Bearer", 15*time.Minute),
					MakeTokenString("Refresh", 15*time.Minute),
				),
			)

			logger, err := logging.NewGoLoggerBuilder().Debug(false).Build()
			Expect(err).ToNot(HaveOccurred())

			connection, err := sdk.NewConnectionBuilder().
				Logger(logger).
				Tokens(MakeTokenString("Bearer", 15*time.Minute)).
				URL(apiServer.URL()).
				TokenURL(ssoServer.URL()).
				Build()
			Expect(err).ToNot(HaveOccurred())

			// Save OCM config so GetEnv() can read it
			config.Save(&config.Config{
				URL:          "https://api.openshift.com",
				AccessToken:  MakeTokenString("Bearer", 15*time.Minute),
				RefreshToken: MakeTokenString("Refresh", 15*time.Minute),
			})

			testRuntime = &rosacli.Runtime{
				Reporter:  reporter.CreateReporter(),
				OCMClient: ocm.NewClientWithConnection(connection),
				AWSClient: awsClient,
				Creator: &rosaaws.Creator{
					AccountID: testAccountID,
					ARN:       fmt.Sprintf("arn:aws:iam::%s:user/test", testAccountID),
					Partition: "aws",
				},
			}
		})

		AfterEach(func() {
			ssoServer.Close()
			apiServer.Close()
			os.Setenv("OCM_CONFIG", "")
			if tmpdir != "" {
				os.RemoveAll(tmpdir)
			}
		})

		It("should return existing role with created=false (idempotency)", func() {
			roleName := "test-prefix-OCM-Role-12345678"
			existingARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", testAccountID, roleName)
			expectedPolicyARN := fmt.Sprintf("arn:aws:iam::%s:policy/%s-Policy", testAccountID, roleName)

			// Mock OCM GetCurrentOrganization API
			apiServer.AppendHandlers(
				RespondWithJSON(http.StatusOK, `{
					"id": "test-org-id",
					"organization": {
						"id": "test-org-id",
						"external_id": "12345678"
					}
				}`),
			)

			// Mock OCM GetPolicies API
			apiServer.AppendHandlers(
				RespondWithJSON(http.StatusOK, `{
					"items": [
						{
							"id": "sts_ocm_permission_policy",
							"details": "{\"Version\":\"2012-10-17\",\"Statement\":[]}"
						}
					],
					"size": 1
				}`),
			)

			// Mock AWS: role exists as standard profile with policy already attached
			awsClient.EXPECT().CheckRoleExists(roleName).Return(true, existingARN, nil)
			awsClient.EXPECT().IsAdminRole(roleName).Return(false, nil)
			awsClient.EXPECT().GetRoleByName(roleName).Return(iamtypes.Role{}, nil)
			awsClient.EXPECT().ListAttachedRolePolicies(roleName).Return([]string{expectedPolicyARN}, nil)

			roleARN, orgID, created, err := GetOrCreateOCMRole(testRuntime, testLog, "test-prefix", ProfileStandard, "", "/")

			Expect(err).ToNot(HaveOccurred())
			Expect(created).To(BeFalse(), "should return created=false when role exists")
			Expect(roleARN).To(Equal(existingARN))
			Expect(orgID).To(Equal("test-org-id"))
		})

		It("should fail when no-console policy is missing (orphan prevention)", func() {
			roleName := "test-prefix-OCM-Role-12345678"

			// Mock OCM GetCurrentOrganization API
			apiServer.AppendHandlers(
				RespondWithJSON(http.StatusOK, `{
					"id": "test-org-id",
					"organization": {
						"id": "test-org-id",
						"external_id": "12345678"
					}
				}`),
			)

			awsClient.EXPECT().CheckRoleExists(roleName).Return(false, "", nil)

			// Mock OCM GetPolicies API - return policies WITHOUT no-console policy
			apiServer.AppendHandlers(
				RespondWithJSON(http.StatusOK, `{
					"items": [
						{
							"id": "sts_ocm_permission_policy",
							"details": "{\"Version\":\"2012-10-17\",\"Statement\":[]}"
						}
					]
				}`),
			)

			roleARN, orgID, created, err := GetOrCreateOCMRole(testRuntime, testLog, "test-prefix", ProfileNoConsole, "", "/")

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("no-console OCM role profile is not yet enabled"))
			Expect(created).To(BeFalse())
			Expect(roleARN).To(BeEmpty())
			Expect(orgID).To(BeEmpty())
		})

		It("should self-heal admin role when admin policy exists but tag is missing", func() {
			roleName := "test-prefix-OCM-Role-12345678"
			existingARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", testAccountID, roleName)
			standardPolicyARN := fmt.Sprintf("arn:aws:iam::%s:policy/%s-Policy", testAccountID, roleName)
			adminPolicyARN := fmt.Sprintf("arn:aws:iam::%s:policy/%s-Admin-Policy", testAccountID, roleName)

			// Mock OCM GetCurrentOrganization API
			apiServer.AppendHandlers(
				RespondWithJSON(http.StatusOK, `{
					"id": "test-org-id",
					"organization": {
						"id": "test-org-id",
						"external_id": "12345678"
					}
				}`),
			)

			// Mock OCM GetPolicies API
			apiServer.AppendHandlers(
				RespondWithJSON(http.StatusOK, `{
					"items": [
						{
							"id": "sts_ocm_permission_policy",
							"details": "{\"Version\":\"2012-10-17\",\"Statement\":[]}"
						},
						{
							"id": "sts_ocm_admin_permission_policy",
							"details": "{\"Version\":\"2012-10-17\",\"Statement\":[]}"
						}
					],
					"size": 2
				}`),
			)

			// Mock AWS: role exists, IsAdminRole returns false (tag missing)
			awsClient.EXPECT().CheckRoleExists(roleName).Return(true, existingARN, nil)
			awsClient.EXPECT().IsAdminRole(roleName).Return(false, nil)
			awsClient.EXPECT().GetRoleByName(roleName).Return(iamtypes.Role{}, nil)

			// Mock ListAttachedRolePolicies - both policies ARE attached
			awsClient.EXPECT().ListAttachedRolePolicies(roleName).Return([]string{standardPolicyARN, adminPolicyARN}, nil)

			// Expect self-healing: AddRoleTag should be called
			awsClient.EXPECT().AddRoleTag(roleName, "rosa_admin_role", "true").Return(nil)

			roleARN, orgID, created, err := GetOrCreateOCMRole(testRuntime, testLog, "test-prefix", ProfileAdmin, "", "/")

			Expect(err).ToNot(HaveOccurred())
			Expect(created).To(BeFalse())
			Expect(roleARN).To(Equal(existingARN))
			Expect(orgID).To(Equal("test-org-id"))
		})

		It("should self-heal standard role when policy is missing", func() {
			roleName := "test-prefix-OCM-Role-12345678"
			existingARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", testAccountID, roleName)
			standardPolicyARN := fmt.Sprintf("arn:aws:iam::%s:policy/%s-Policy", testAccountID, roleName)

			// Mock OCM GetCurrentOrganization API
			apiServer.AppendHandlers(
				RespondWithJSON(http.StatusOK, `{
					"id": "test-org-id",
					"organization": {
						"id": "test-org-id",
						"external_id": "12345678"
					}
				}`),
			)

			// Mock OCM GetPolicies API
			apiServer.AppendHandlers(
				RespondWithJSON(http.StatusOK, `{
					"items": [
						{
							"id": "sts_ocm_permission_policy",
							"details": "{\"Version\":\"2012-10-17\",\"Statement\":[]}"
						}
					],
					"size": 1
				}`),
			)

			// Mock AWS: role exists but policy is missing
			awsClient.EXPECT().CheckRoleExists(roleName).Return(true, existingARN, nil)
			awsClient.EXPECT().IsAdminRole(roleName).Return(false, nil)
			awsClient.EXPECT().GetRoleByName(roleName).Return(iamtypes.Role{}, nil)
			awsClient.EXPECT().ListAttachedRolePolicies(roleName).Return([]string{}, nil) // Empty - no policies attached

			// Expect self-healing: policy should be created and attached
			awsClient.EXPECT().EnsurePolicy(standardPolicyARN, gomock.Any(), "", gomock.Any(), "/").Return(standardPolicyARN, nil)
			awsClient.EXPECT().AttachRolePolicy(gomock.Any(), roleName, standardPolicyARN).Return(nil)

			roleARN, orgID, created, err := GetOrCreateOCMRole(testRuntime, testLog, "test-prefix", ProfileStandard, "", "/")

			Expect(err).ToNot(HaveOccurred())
			Expect(created).To(BeFalse())
			Expect(roleARN).To(Equal(existingARN))
			Expect(orgID).To(Equal("test-org-id"))
		})

		It("should self-heal admin role when standard policy is missing", func() {
			roleName := "test-prefix-OCM-Role-12345678"
			existingARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", testAccountID, roleName)
			standardPolicyARN := fmt.Sprintf("arn:aws:iam::%s:policy/%s-Policy", testAccountID, roleName)
			adminPolicyARN := fmt.Sprintf("arn:aws:iam::%s:policy/%s-Admin-Policy", testAccountID, roleName)

			// Mock OCM GetCurrentOrganization API
			apiServer.AppendHandlers(
				RespondWithJSON(http.StatusOK, `{
					"id": "test-org-id",
					"organization": {
						"id": "test-org-id",
						"external_id": "12345678"
					}
				}`),
			)

			// Mock OCM GetPolicies API
			apiServer.AppendHandlers(
				RespondWithJSON(http.StatusOK, `{
					"items": [
						{
							"id": "sts_ocm_permission_policy",
							"details": "{\"Version\":\"2012-10-17\",\"Statement\":[]}"
						},
						{
							"id": "sts_ocm_admin_permission_policy",
							"details": "{\"Version\":\"2012-10-17\",\"Statement\":[]}"
						}
					],
					"size": 2
				}`),
			)

			// Mock AWS: role exists with admin policy but missing standard policy
			awsClient.EXPECT().CheckRoleExists(roleName).Return(true, existingARN, nil)
			awsClient.EXPECT().IsAdminRole(roleName).Return(true, nil) // Has admin tag
			awsClient.EXPECT().GetRoleByName(roleName).Return(iamtypes.Role{}, nil)
			awsClient.EXPECT().ListAttachedRolePolicies(roleName).Return([]string{adminPolicyARN}, nil) // Only admin policy

			// Expect self-healing: standard policy should be created and attached
			awsClient.EXPECT().EnsurePolicy(standardPolicyARN, gomock.Any(), "", gomock.Any(), "/").Return(standardPolicyARN, nil)
			awsClient.EXPECT().AttachRolePolicy(gomock.Any(), roleName, standardPolicyARN).Return(nil)

			roleARN, orgID, created, err := GetOrCreateOCMRole(testRuntime, testLog, "test-prefix", ProfileAdmin, "", "/")

			Expect(err).ToNot(HaveOccurred())
			Expect(created).To(BeFalse())
			Expect(roleARN).To(Equal(existingARN))
			Expect(orgID).To(Equal("test-org-id"))
		})

		It("should self-heal admin role when admin policy is missing", func() {
			roleName := "test-prefix-OCM-Role-12345678"
			existingARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", testAccountID, roleName)
			standardPolicyARN := fmt.Sprintf("arn:aws:iam::%s:policy/%s-Policy", testAccountID, roleName)
			adminPolicyARN := fmt.Sprintf("arn:aws:iam::%s:policy/%s-Admin-Policy", testAccountID, roleName)

			// Mock OCM GetCurrentOrganization API
			apiServer.AppendHandlers(
				RespondWithJSON(http.StatusOK, `{
					"id": "test-org-id",
					"organization": {
						"id": "test-org-id",
						"external_id": "12345678"
					}
				}`),
			)

			// Mock OCM GetPolicies API
			apiServer.AppendHandlers(
				RespondWithJSON(http.StatusOK, `{
					"items": [
						{
							"id": "sts_ocm_permission_policy",
							"details": "{\"Version\":\"2012-10-17\",\"Statement\":[]}"
						},
						{
							"id": "sts_ocm_admin_permission_policy",
							"details": "{\"Version\":\"2012-10-17\",\"Statement\":[]}"
						}
					],
					"size": 2
				}`),
			)

			// Mock AWS: role exists with standard policy but missing admin policy
			awsClient.EXPECT().CheckRoleExists(roleName).Return(true, existingARN, nil)
			awsClient.EXPECT().IsAdminRole(roleName).Return(false, nil) // No admin tag
			awsClient.EXPECT().GetRoleByName(roleName).Return(iamtypes.Role{}, nil)
			awsClient.EXPECT().ListAttachedRolePolicies(roleName).Return([]string{standardPolicyARN}, nil) // Only standard policy

			// Expect self-healing: admin policy should be created and attached, tag added
			awsClient.EXPECT().EnsurePolicy(adminPolicyARN, gomock.Any(), "", gomock.Any(), "/").Return(adminPolicyARN, nil)
			awsClient.EXPECT().AttachRolePolicy(gomock.Any(), roleName, adminPolicyARN).Return(nil)
			awsClient.EXPECT().AddRoleTag(roleName, "rosa_admin_role", "true").Return(nil)

			roleARN, orgID, created, err := GetOrCreateOCMRole(testRuntime, testLog, "test-prefix", ProfileAdmin, "", "/")

			Expect(err).ToNot(HaveOccurred())
			Expect(created).To(BeFalse())
			Expect(roleARN).To(Equal(existingARN))
			Expect(orgID).To(Equal("test-org-id"))
		})

		It("should self-heal no-console role when policy is missing", func() {
			roleName := "test-prefix-OCM-Role-12345678"
			existingARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", testAccountID, roleName)
			noConsolePolicyARN := fmt.Sprintf("arn:aws:iam::%s:policy/%s-NoConsole-Policy", testAccountID, roleName)

			// Mock OCM GetCurrentOrganization API
			apiServer.AppendHandlers(
				RespondWithJSON(http.StatusOK, `{
					"id": "test-org-id",
					"organization": {
						"id": "test-org-id",
						"external_id": "12345678"
					}
				}`),
			)

			// Mock OCM GetPolicies API
			apiServer.AppendHandlers(
				RespondWithJSON(http.StatusOK, `{
					"items": [
						{
							"id": "sts_ocm_no_console_permission_policy",
							"details": "{\"Version\":\"2012-10-17\",\"Statement\":[]}"
						}
					],
					"size": 1
				}`),
			)

			// Mock AWS: role exists but no-console policy is missing
			awsClient.EXPECT().CheckRoleExists(roleName).Return(true, existingARN, nil)
			awsClient.EXPECT().IsAdminRole(roleName).Return(false, nil)
			awsClient.EXPECT().GetRoleByName(roleName).Return(iamtypes.Role{}, nil)
			awsClient.EXPECT().ListAttachedRolePolicies(roleName).Return([]string{}, nil) // No policies attached

			// Expect self-healing: no-console policy should be created and attached, tag added
			awsClient.EXPECT().EnsurePolicy(noConsolePolicyARN, gomock.Any(), "", gomock.Any(), "/").Return(noConsolePolicyARN, nil)
			awsClient.EXPECT().AttachRolePolicy(gomock.Any(), roleName, noConsolePolicyARN).Return(nil)
			awsClient.EXPECT().AddRoleTag(roleName, "rosa_no_console_role", "true").Return(nil)

			roleARN, orgID, created, err := GetOrCreateOCMRole(testRuntime, testLog, "test-prefix", ProfileNoConsole, "", "/")

			Expect(err).ToNot(HaveOccurred())
			Expect(created).To(BeFalse())
			Expect(roleARN).To(Equal(existingARN))
			Expect(orgID).To(Equal("test-org-id"))
		})
	})
})
