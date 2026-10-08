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
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	ststypes "github.com/aws/aws-sdk-go-v2/service/sts/types"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"

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

	staticProvider := NewAWSStaticPrincipalTypeProvider(&infrav1.AWSClusterStaticIdentity{}, secret)

	stsMock := mock_stsiface.NewMockSTSClient(mockCtrl)
	roleIdentity := &infrav1.AWSClusterRoleIdentity{
		Spec: infrav1.AWSClusterRoleIdentitySpec{
			AWSRoleSpec: infrav1.AWSRoleSpec{
				RoleArn:         "arn:*:iam::*:role/aws-role/firstroleprovider",
				SessionName:     "first-role-provider-session",
				DurationSeconds: 900,
			},
		},
	}

	expiresAt := time.Now()

	roleProvider := &AWSRolePrincipalTypeProvider{
		credentials:    nil,
		Principal:      roleIdentity,
		region:         "us-west-2",
		sourceProvider: staticProvider,
		stsClient:      stsMock,
	}

	testCases := []struct {
		name      string
		provider  AWSPrincipalTypeProvider
		expect    func(m *mock_stsiface.MockSTSClientMockRecorder)
		expectErr bool
		value     aws.Credentials
	}{
		{
			name:      "Static provider successfully retrieves",
			provider:  staticProvider,
			expect:    func(m *mock_stsiface.MockSTSClientMockRecorder) {},
			expectErr: false,
			value: aws.Credentials{
				AccessKeyID:     "static-AccessKeyID",
				SecretAccessKey: "static-SecretAccessKey",
				Source:          "StaticCredentials",
				CanExpire:       false,
				Expires:         time.Time{},
			},
		},
		{
			name:     "Role provider with static provider source successfully retrieves",
			provider: roleProvider,
			expect: func(m *mock_stsiface.MockSTSClientMockRecorder) {
				m.AssumeRole(gomock.Any(), &sts.AssumeRoleInput{
					RoleArn:         aws.String(roleIdentity.Spec.RoleArn),
					RoleSessionName: aws.String(roleIdentity.Spec.SessionName),
					DurationSeconds: aws.Int32(roleIdentity.Spec.DurationSeconds),
				}).Return(&sts.AssumeRoleOutput{
					Credentials: &ststypes.Credentials{
						AccessKeyId:     aws.String("assumedAccessKeyId"),
						SecretAccessKey: aws.String("assumedSecretAccessKey"),
						SessionToken:    aws.String("assumedSessionToken"),
						Expiration:      aws.Time(expiresAt),
					},
				}, nil)
			},
			expectErr: false,
			value: aws.Credentials{
				AccessKeyID:     "assumedAccessKeyId",
				SecretAccessKey: "assumedSecretAccessKey",
				SessionToken:    "assumedSessionToken",
				Source:          "AssumeRoleProvider",
				CanExpire:       true,
				Expires:         expiresAt,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			tc.expect(stsMock.EXPECT())
			value, err := tc.provider.Retrieve(context.TODO())
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

type rotatingCredentialsProvider struct {
	calls int
}

func (p *rotatingCredentialsProvider) Retrieve(context.Context) (aws.Credentials, error) {
	p.calls++
	if p.calls == 1 {
		return aws.Credentials{
			AccessKeyID:     "root-old",
			SecretAccessKey: "root-secret-old",
			SessionToken:    "root-token-old",
			Source:          "RotatingCredentialsProvider",
			CanExpire:       true,
			Expires:         time.Now().Add(-time.Minute),
		}, nil
	}

	return aws.Credentials{
		AccessKeyID:     "root-new",
		SecretAccessKey: "root-secret-new",
		SessionToken:    "root-token-new",
		Source:          "RotatingCredentialsProvider",
		CanExpire:       true,
		Expires:         time.Now().Add(time.Hour),
	}, nil
}

func (*rotatingCredentialsProvider) Hash() (string, error) {
	return "rotating-root", nil
}

func (*rotatingCredentialsProvider) Name() string {
	return "rotating-root"
}

func TestAWSRolePrincipalTypeProviderRefreshesNestedSourceProviders(t *testing.T) {
	rootProvider := &rotatingCredentialsProvider{}
	requests := make([]string, 0, 4)

	stsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		params, err := url.ParseQuery(string(body))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		authorization := r.Header.Get("Authorization")
		credentialStart := strings.Index(authorization, "Credential=")
		if credentialStart < 0 {
			http.Error(w, "missing access key", http.StatusBadRequest)
			return
		}
		credentialStart += len("Credential=")
		credentialEnd := strings.Index(authorization[credentialStart:], "/")
		if credentialEnd < 0 {
			http.Error(w, "missing access key", http.StatusBadRequest)
			return
		}
		accessKeyID := authorization[credentialStart : credentialStart+credentialEnd]
		roleARN := params.Get("RoleArn")
		requests = append(requests, roleARN+"="+accessKeyID)

		call := len(requests)
		accessKeySuffix := "old"
		expiration := time.Now().Add(-time.Minute).UTC().Format(time.RFC3339)
		if call > 2 {
			accessKeySuffix = "new"
			expiration = time.Now().Add(time.Hour).UTC().Format(time.RFC3339)
		}
		roleName := params.Get("RoleArn")
		if separator := strings.LastIndex(roleName, "/"); separator >= 0 {
			roleName = roleName[separator+1:]
		}

		w.Header().Set("Content-Type", "text/xml")
		fmt.Fprintf(w, `<AssumeRoleResponse xmlns="https://sts.amazonaws.com/doc/2011-06-15/">
  <AssumeRoleResult>
    <Credentials>
      <AccessKeyId>%s-%s</AccessKeyId>
      <SecretAccessKey>secret-%s</SecretAccessKey>
      <SessionToken>token-%s</SessionToken>
      <Expiration>%s</Expiration>
    </Credentials>
  </AssumeRoleResult>
</AssumeRoleResponse>`, roleName, accessKeySuffix, accessKeySuffix, accessKeySuffix, expiration)
	}))
	defer stsServer.Close()

	t.Setenv("AWS_ENDPOINT_URL_STS", stsServer.URL)
	t.Setenv("AWS_REGION", "us-east-1")

	innerProvider := &AWSRolePrincipalTypeProvider{
		Principal: &infrav1.AWSClusterRoleIdentity{Spec: infrav1.AWSClusterRoleIdentitySpec{AWSRoleSpec: infrav1.AWSRoleSpec{
			RoleArn:         "arn:aws:iam::123456789012:role/inner",
			SessionName:     "inner-session",
			DurationSeconds: 900,
		}}},
		region:         "us-east-1",
		sourceProvider: rootProvider,
	}
	outerProvider := &AWSRolePrincipalTypeProvider{
		Principal: &infrav1.AWSClusterRoleIdentity{Spec: infrav1.AWSClusterRoleIdentitySpec{AWSRoleSpec: infrav1.AWSRoleSpec{
			RoleArn:         "arn:aws:iam::123456789012:role/outer",
			SessionName:     "outer-session",
			DurationSeconds: 900,
		}}},
		region:         "us-east-1",
		sourceProvider: innerProvider,
	}

	first, err := outerProvider.Retrieve(context.Background())
	if err != nil {
		t.Fatalf("first retrieve failed: %v", err)
	}
	if first.AccessKeyID != "outer-old" {
		t.Fatalf("unexpected first credentials: %#v", first)
	}

	second, err := outerProvider.Retrieve(context.Background())
	if err != nil {
		t.Fatalf("second retrieve failed: %v", err)
	}
	if second.AccessKeyID != "outer-new" {
		t.Fatalf("unexpected refreshed credentials: %#v", second)
	}
	if rootProvider.calls != 2 {
		t.Fatalf("expected root provider to be called twice, got %d", rootProvider.calls)
	}

	expectedRequests := []string{
		"arn:aws:iam::123456789012:role/inner=root-old",
		"arn:aws:iam::123456789012:role/outer=inner-old",
		"arn:aws:iam::123456789012:role/inner=root-new",
		"arn:aws:iam::123456789012:role/outer=inner-new",
	}
	if diff := cmp.Diff(expectedRequests, requests); diff != "" {
		t.Fatalf("unexpected AssumeRole requests (-want +got):\n%s", diff)
	}
}
