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

package eks

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/eks/mock_eksiface"
	eksaddons "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/eks/addons"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
)

func TestUpdateAddonTagsWithLatestVersion(t *testing.T) {
	tests := []struct {
		name         string
		inputAddons  []*eksaddons.EKSAddon
		mockResponse *eks.DescribeAddonVersionsOutput
		mockError    error
		expected     []*eksaddons.EKSAddon
		expectError  bool
	}{
		{
			name: "replaces latest tag with highest semver",
			inputAddons: []*eksaddons.EKSAddon{
				{Name: aws.String("vpc-cni"), Version: aws.String("latest")},
			},
			mockResponse: &eks.DescribeAddonVersionsOutput{
				Addons: []ekstypes.AddonInfo{
					{
						AddonName: aws.String("vpc-cni"),
						AddonVersions: []ekstypes.AddonVersionInfo{
							{AddonVersion: aws.String("v1.12.0-eksbuild.1")},
							{AddonVersion: aws.String("v1.11.4-eksbuild.1")},
							{AddonVersion: aws.String("v1.12.5-eksbuild.2")},
							{AddonVersion: aws.String("v1.12.1-eksbuild.1")},
						},
					},
				},
			},
			expected: []*eksaddons.EKSAddon{
				{Name: aws.String("vpc-cni"), Version: aws.String("v1.12.5-eksbuild.2")},
			},
		},
		{
			name: "preserves specific version",
			inputAddons: []*eksaddons.EKSAddon{
				{Name: aws.String("vpc-cni"), Version: aws.String("v1.11.4-eksbuild.1")},
			},
			mockResponse: &eks.DescribeAddonVersionsOutput{
				Addons: []ekstypes.AddonInfo{
					{
						AddonName: aws.String("vpc-cni"),
						AddonVersions: []ekstypes.AddonVersionInfo{
							{AddonVersion: aws.String("v1.12.5-eksbuild.2")},
							{AddonVersion: aws.String("v1.11.4-eksbuild.1")},
						},
					},
				},
			},
			expected: []*eksaddons.EKSAddon{
				{Name: aws.String("vpc-cni"), Version: aws.String("v1.11.4-eksbuild.1")},
			},
		},
		{
			name: "handles multiple addons with latest",
			inputAddons: []*eksaddons.EKSAddon{
				{Name: aws.String("vpc-cni"), Version: aws.String("latest")},
				{Name: aws.String("coredns"), Version: aws.String("latest")},
				{Name: aws.String("kube-proxy"), Version: aws.String("v1.28.1-eksbuild.1")},
			},
			mockResponse: &eks.DescribeAddonVersionsOutput{
				Addons: []ekstypes.AddonInfo{
					{
						AddonName: aws.String("vpc-cni"),
						AddonVersions: []ekstypes.AddonVersionInfo{
							{AddonVersion: aws.String("v1.12.5-eksbuild.2")},
							{AddonVersion: aws.String("v1.12.0-eksbuild.1")},
						},
					},
					{
						AddonName: aws.String("coredns"),
						AddonVersions: []ekstypes.AddonVersionInfo{
							{AddonVersion: aws.String("v1.10.1-eksbuild.2")},
							{AddonVersion: aws.String("v1.10.1-eksbuild.1")},
							{AddonVersion: aws.String("v1.9.3-eksbuild.5")},
						},
					},
					{
						AddonName: aws.String("kube-proxy"),
						AddonVersions: []ekstypes.AddonVersionInfo{
							{AddonVersion: aws.String("v1.28.2-eksbuild.2")},
							{AddonVersion: aws.String("v1.28.1-eksbuild.1")},
						},
					},
				},
			},
			expected: []*eksaddons.EKSAddon{
				{Name: aws.String("vpc-cni"), Version: aws.String("v1.12.5-eksbuild.2")},
				{Name: aws.String("coredns"), Version: aws.String("v1.10.1-eksbuild.2")},
				{Name: aws.String("kube-proxy"), Version: aws.String("v1.28.1-eksbuild.1")},
			},
		},
		{
			name: "handles addon not in versions list",
			inputAddons: []*eksaddons.EKSAddon{
				{Name: aws.String("unknown-addon"), Version: aws.String("latest")},
			},
			mockResponse: &eks.DescribeAddonVersionsOutput{
				Addons: []ekstypes.AddonInfo{
					{
						AddonName: aws.String("vpc-cni"),
						AddonVersions: []ekstypes.AddonVersionInfo{
							{AddonVersion: aws.String("v1.12.5-eksbuild.2")},
						},
					},
				},
			},
			expected: []*eksaddons.EKSAddon{
				{Name: aws.String("unknown-addon"), Version: aws.String("latest")},
			},
		},
		{
			name: "correctly compares build numbers",
			inputAddons: []*eksaddons.EKSAddon{
				{Name: aws.String("vpc-cni"), Version: aws.String("latest")},
			},
			mockResponse: &eks.DescribeAddonVersionsOutput{
				Addons: []ekstypes.AddonInfo{
					{
						AddonName: aws.String("vpc-cni"),
						AddonVersions: []ekstypes.AddonVersionInfo{
							{AddonVersion: aws.String("v1.12.5-eksbuild.1")},
							{AddonVersion: aws.String("v1.12.5-eksbuild.10")},
							{AddonVersion: aws.String("v1.12.5-eksbuild.2")},
						},
					},
				},
			},
			expected: []*eksaddons.EKSAddon{
				{Name: aws.String("vpc-cni"), Version: aws.String("v1.12.5-eksbuild.10")},
			},
		},
		{
			name: "returns error when DescribeAddonVersions fails",
			inputAddons: []*eksaddons.EKSAddon{
				{Name: aws.String("vpc-cni"), Version: aws.String("latest")},
			},
			mockError:   errors.New("API error"),
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			mockControl := gomock.NewController(t)
			defer mockControl.Finish()

			eksMock := mock_eksiface.NewMockEKSAPI(mockControl)

			scheme := runtime.NewScheme()
			_ = infrav1.AddToScheme(scheme)
			_ = ekscontrolplanev1.AddToScheme(scheme)
			client := fake.NewClientBuilder().WithScheme(scheme).Build()

			scope, err := scope.NewManagedControlPlaneScope(scope.ManagedControlPlaneScopeParams{
				Client: client,
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "test-cluster",
					},
				},
				ControlPlane: &ekscontrolplanev1.AWSManagedControlPlane{
					Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
						Version: aws.String("1.34"),
					},
				},
			})
			g.Expect(err).To(BeNil())

			eksMock.EXPECT().
				DescribeAddonVersions(gomock.Any(), gomock.Any()).
				Return(tc.mockResponse, tc.mockError)

			s := NewService(scope)
			s.EKSClient = eksMock

			result, err := s.updateAddonTagsWithLatestVersion(context.TODO(), tc.inputAddons)

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
				return
			}

			g.Expect(err).To(BeNil())
			g.Expect(result).To(HaveLen(len(tc.expected)))

			for i := range result {
				g.Expect(*result[i].Name).To(Equal(*tc.expected[i].Name))
				g.Expect(*result[i].Version).To(Equal(*tc.expected[i].Version))
			}
		})
	}
}
