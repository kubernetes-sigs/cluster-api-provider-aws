/*
Copyright 2020 The Kubernetes Authors.

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

package iamauth

import (
	"context"
	"reflect"
	"strings"
	"testing"

	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	iamauthv1 "sigs.k8s.io/aws-iam-authenticator/pkg/mapper/crd/apis/iamauthenticator/v1alpha1"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1alpha4"
)

func TestAddRoleMappingCRD(t *testing.T) {
	testCases := []struct {
		name                 string
		existingRoleMapping  *iamauthv1.IAMIdentityMapping
		roleToMap            ekscontrolplanev1.RoleMapping
		expectedRoleMapSpecs []iamauthv1.IAMIdentityMappingSpec
		expectError          bool
	}{
		{
			name: "no existing mappings, add role mapping",
			roleToMap: ekscontrolplanev1.RoleMapping{
				RoleARN: "arn:aws:iam::000000000000:role/KubernetesNode",
				KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
					UserName: "system:node:{{EC2PrivateDNSName}}",
					Groups:   []string{"system:bootstrappers", "system:nodes"},
				},
			},
			expectedRoleMapSpecs: []iamauthv1.IAMIdentityMappingSpec{
				{
					ARN:      "arn:aws:iam::000000000000:role/KubernetesNode",
					Username: "system:node:{{EC2PrivateDNSName}}",
					Groups:   []string{"system:bootstrappers", "system:nodes"},
				},
			},
			expectError: false,
		},
		{
			name: "existing mapping, add different role mapping",
			roleToMap: ekscontrolplanev1.RoleMapping{
				RoleARN: "arn:aws:iam::000000000000:role/KubernetesNode",
				KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
					UserName: "system:node:{{EC2PrivateDNSName}}",
					Groups:   []string{"system:bootstrappers", "system:nodes"},
				},
			},
			existingRoleMapping: createIAMAuthMapping("arn:aws:iam::000000000000:role/KubernetesAdmin", "admin:{{SessionName}}", []string{"system:masters"}),
			expectedRoleMapSpecs: []iamauthv1.IAMIdentityMappingSpec{
				{
					ARN:      "arn:aws:iam::000000000000:role/KubernetesAdmin",
					Username: "admin:{{SessionName}}",
					Groups:   []string{"system:masters"},
				},
				{
					ARN:      "arn:aws:iam::000000000000:role/KubernetesNode",
					Username: "system:node:{{EC2PrivateDNSName}}",
					Groups:   []string{"system:bootstrappers", "system:nodes"},
				},
			},
			expectError: false,
		},
		{
			name: "existing mapping, add same role mapping",
			roleToMap: ekscontrolplanev1.RoleMapping{
				RoleARN: "arn:aws:iam::000000000000:role/KubernetesNode",
				KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
					UserName: "system:node:{{EC2PrivateDNSName}}",
					Groups:   []string{"system:bootstrappers", "system:nodes"},
				},
			},
			existingRoleMapping: createIAMAuthMapping("arn:aws:iam::000000000000:role/KubernetesNode", "system:node:{{EC2PrivateDNSName}}", []string{"system:bootstrappers", "system:nodes"}),
			expectedRoleMapSpecs: []iamauthv1.IAMIdentityMappingSpec{
				{
					ARN:      "arn:aws:iam::000000000000:role/KubernetesNode",
					Username: "system:node:{{EC2PrivateDNSName}}",
					Groups:   []string{"system:bootstrappers", "system:nodes"},
				},
			},
			expectError: false,
		},
		{
			name: "no existing mapping, add role with not role ARN",
			roleToMap: ekscontrolplanev1.RoleMapping{
				RoleARN: "arn:aws:iam::000000000000:user/Alice",
				KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
					UserName: "system:node:{{EC2PrivateDNSName}}",
					Groups:   []string{"system:bootstrappers", "system:nodes"},
				},
			},
			expectedRoleMapSpecs: []iamauthv1.IAMIdentityMappingSpec{},
			expectError:          true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			scheme := runtime.NewScheme()
			iamauthv1.AddToScheme(scheme)

			var client crclient.Client
			if tc.existingRoleMapping == nil {
				client = fake.NewClientBuilder().WithScheme(scheme).Build()
			} else {
				client = fake.NewClientBuilder().WithScheme(scheme).WithObjects(tc.existingRoleMapping).Build()
			}
			backend, err := NewBackend(BackendTypeCRD, client)
			g.Expect(err).To(BeNil())

			err = backend.MapRole(tc.roleToMap)
			if tc.expectError {
				g.Expect(err).ToNot(BeNil())
				return
			}

			g.Expect(err).To(BeNil())

			mappings := &iamauthv1.IAMIdentityMappingList{}
			err = client.List(context.TODO(), mappings)
			g.Expect(err).To(BeNil())

			g.Expect(len(mappings.Items)).To(Equal(len(tc.expectedRoleMapSpecs)))

			for _, actualMapping := range mappings.Items {
				found := false
				for _, expectedMappingSpec := range tc.expectedRoleMapSpecs {
					if reflect.DeepEqual(actualMapping.Spec, expectedMappingSpec) {
						found = true
					}
				}
				g.Expect(found).To(BeTrue())
				g.Expect(actualMapping.Namespace).To(Equal("kube-system"))
				g.Expect(strings.HasPrefix(actualMapping.Name, "capa-iamauth-")).To(BeTrue())
			}
		})
	}
}
func TestAddUserMappingCRD(t *testing.T) {
	testCases := []struct {
		name                 string
		existingUserMapping  *iamauthv1.IAMIdentityMapping
		userToMap            ekscontrolplanev1.UserMapping
		expectedUserMapSpecs []iamauthv1.IAMIdentityMappingSpec
		expectError          bool
	}{
		{
			name: "no existing mappings, add user mapping",
			userToMap: ekscontrolplanev1.UserMapping{
				UserARN: "arn:aws:iam::000000000000:user/Alice",
				KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
					UserName: "alice",
					Groups:   []string{"system:masters"},
				},
			},
			expectedUserMapSpecs: []iamauthv1.IAMIdentityMappingSpec{
				{
					ARN:      "arn:aws:iam::000000000000:user/Alice",
					Username: "alice",
					Groups:   []string{"system:masters"},
				},
			},
			expectError: false,
		},
		{
			name: "existing mapping, add different user mapping",
			userToMap: ekscontrolplanev1.UserMapping{
				UserARN: "arn:aws:iam::000000000000:user/Alice",
				KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
					UserName: "alice",
					Groups:   []string{"system:masters"},
				},
			},
			existingUserMapping: createIAMAuthMapping("arn:aws:iam::000000000000:user/Bob", "bob", []string{"system:masters"}),
			expectedUserMapSpecs: []iamauthv1.IAMIdentityMappingSpec{
				{
					ARN:      "arn:aws:iam::000000000000:user/Bob",
					Username: "bob",
					Groups:   []string{"system:masters"},
				},
				{
					ARN:      "arn:aws:iam::000000000000:user/Alice",
					Username: "alice",
					Groups:   []string{"system:masters"},
				},
			},
			expectError: false,
		},
		{
			name: "existing mapping, add same user mapping",
			userToMap: ekscontrolplanev1.UserMapping{
				UserARN: "arn:aws:iam::000000000000:user/Alice",
				KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
					UserName: "alice",
					Groups:   []string{"system:masters"},
				},
			},
			existingUserMapping: createIAMAuthMapping("arn:aws:iam::000000000000:user/Alice", "alice", []string{"system:masters"}),
			expectedUserMapSpecs: []iamauthv1.IAMIdentityMappingSpec{
				{
					ARN:      "arn:aws:iam::000000000000:user/Alice",
					Username: "alice",
					Groups:   []string{"system:masters"},
				},
			},
			expectError: false,
		},
		{
			name: "no existing mapping, add role with not role ARN",
			userToMap: ekscontrolplanev1.UserMapping{
				UserARN: "arn:aws:iam::000000000000:role/KubernetesNode",
				KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
					UserName: "system:node:{{EC2PrivateDNSName}}",
					Groups:   []string{"system:masters"},
				},
			},
			expectedUserMapSpecs: []iamauthv1.IAMIdentityMappingSpec{},
			expectError:          true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			scheme := runtime.NewScheme()
			iamauthv1.AddToScheme(scheme)

			var client crclient.Client
			if tc.existingUserMapping == nil {
				client = fake.NewClientBuilder().WithScheme(scheme).Build()
			} else {
				client = fake.NewClientBuilder().WithScheme(scheme).WithObjects(tc.existingUserMapping).Build()
			}
			backend, err := NewBackend(BackendTypeCRD, client)
			g.Expect(err).To(BeNil())

			err = backend.MapUser(tc.userToMap)
			if tc.expectError {
				g.Expect(err).ToNot(BeNil())
				return
			}

			g.Expect(err).To(BeNil())

			mappings := &iamauthv1.IAMIdentityMappingList{}
			err = client.List(context.TODO(), mappings)
			g.Expect(err).To(BeNil())

			g.Expect(len(mappings.Items)).To(Equal(len(tc.expectedUserMapSpecs)))

			for _, actualMapping := range mappings.Items {
				found := false
				for _, expectedMappingSpec := range tc.expectedUserMapSpecs {
					if reflect.DeepEqual(actualMapping.Spec, expectedMappingSpec) {
						found = true
					}
				}
				g.Expect(found).To(BeTrue())
				g.Expect(actualMapping.Namespace).To(Equal("kube-system"))
				g.Expect(strings.HasPrefix(actualMapping.Name, "capa-iamauth-")).To(BeTrue())
			}
		})
	}
}

func createIAMAuthMapping(arn string, username string, groups []string) *iamauthv1.IAMIdentityMapping {
	return &iamauthv1.IAMIdentityMapping{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "capa-iamauth-abcd1234",
			Namespace: "kube-system",
			UID:       "1234567890",
		},
		Spec: iamauthv1.IAMIdentityMappingSpec{
			ARN:      arn,
			Username: username,
			Groups:   groups,
		},
	}
}
