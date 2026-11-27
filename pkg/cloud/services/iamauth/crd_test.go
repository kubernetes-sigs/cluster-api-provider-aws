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
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	iamauthv1 "sigs.k8s.io/aws-iam-authenticator/pkg/mapper/crd/apis/iamauthenticator/v1alpha1"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
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
					Groups:   []string{systemBootstrappersGroup, systemNodesGroup},
				},
			},
			expectedRoleMapSpecs: []iamauthv1.IAMIdentityMappingSpec{
				{
					ARN:      "arn:aws:iam::000000000000:role/KubernetesNode",
					Username: "system:node:{{EC2PrivateDNSName}}",
					Groups:   []string{systemBootstrappersGroup, systemNodesGroup},
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
					Groups:   []string{systemBootstrappersGroup, systemNodesGroup},
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
					Groups:   []string{systemBootstrappersGroup, systemNodesGroup},
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
					Groups:   []string{systemBootstrappersGroup, systemNodesGroup},
				},
			},
			existingRoleMapping: createIAMAuthMapping("arn:aws:iam::000000000000:role/KubernetesNode", "system:node:{{EC2PrivateDNSName}}", []string{systemBootstrappersGroup, systemNodesGroup}),
			expectedRoleMapSpecs: []iamauthv1.IAMIdentityMappingSpec{
				{
					ARN:      "arn:aws:iam::000000000000:role/KubernetesNode",
					Username: "system:node:{{EC2PrivateDNSName}}",
					Groups:   []string{systemBootstrappersGroup, systemNodesGroup},
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
					Groups:   []string{systemBootstrappersGroup, systemNodesGroup},
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
					if cmp.Equal(actualMapping.Spec, expectedMappingSpec) {
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
					if cmp.Equal(actualMapping.Spec, expectedMappingSpec) {
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

// TestReconcileMappingsCRD exercises the desired-state semantics of
// crdBackend.ReconcileMappings. It must:
//   - delete stale CAPA-managed CRs when the desired set shrinks
//   - preserve out-of-band CRs (no capa-iamauth- GenerateName prefix)
//   - create missing entries
//   - be idempotent on repeated calls
//
// Scenarios:
//  1. Deletes stale CAPA-managed CR when a user is removed.
//  2. Preserves out-of-band CR (no CAPA GenerateName prefix).
//  3. Node role preserved when included in the desired set.
//  4. Idempotency: state unchanged on repeated calls.
func TestReconcileMappingsCRD(t *testing.T) {
	alice := ekscontrolplanev1.UserMapping{
		UserARN: "arn:aws:iam::000000000000:user/Alice",
		KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
			UserName: "alice",
			Groups:   []string{"system:masters"},
		},
	}
	bob := ekscontrolplanev1.UserMapping{
		UserARN: "arn:aws:iam::000000000000:user/Bob",
		KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
			UserName: "bob",
			Groups:   []string{"system:masters"},
		},
	}
	nodeRole := ekscontrolplanev1.RoleMapping{
		RoleARN: "arn:aws:iam::000000000000:role/KubernetesNode",
		KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
			UserName: "system:node:{{EC2PrivateDNSName}}",
			Groups:   []string{"system:bootstrappers", "system:nodes"},
		},
	}
	outOfBandCR := &iamauthv1.IAMIdentityMapping{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "operator-hand-crafted",
			Namespace: metav1.NamespaceSystem,
			// Note: no capa-iamauth- GenerateName prefix — this simulates an
			// entry a cluster operator created by hand.
		},
		Spec: iamauthv1.IAMIdentityMappingSpec{
			ARN:      "arn:aws:iam::000000000000:user/OpsAdmin",
			Username: "opsadmin",
			Groups:   []string{"system:masters"},
		},
	}
	capaBobCR := &iamauthv1.IAMIdentityMapping{
		ObjectMeta: metav1.ObjectMeta{
			Name:         "capa-iamauth-abc123",
			Namespace:    metav1.NamespaceSystem,
			GenerateName: capaGenerateName,
		},
		Spec: iamauthv1.IAMIdentityMappingSpec{
			ARN:      bob.UserARN,
			Username: bob.UserName,
			Groups:   bob.Groups,
		},
	}

	tests := []struct {
		name           string
		preexisting    []crclient.Object
		roles          []ekscontrolplanev1.RoleMapping
		users          []ekscontrolplanev1.UserMapping
		wantARNs       []string
		wantARNsAbsent []string
	}{
		{
			name:           "deletes stale CAPA-managed CR when user removed",
			preexisting:    []crclient.Object{capaBobCR},
			users:          []ekscontrolplanev1.UserMapping{alice},
			wantARNs:       []string{alice.UserARN},
			wantARNsAbsent: []string{bob.UserARN},
		},
		{
			name:        "preserves out-of-band CR (no CAPA GenerateName prefix)",
			preexisting: []crclient.Object{outOfBandCR},
			users:       []ekscontrolplanev1.UserMapping{alice},
			wantARNs:    []string{alice.UserARN, outOfBandCR.Spec.ARN},
		},
		{
			name:     "node role preserved when included in desired set",
			roles:    []ekscontrolplanev1.RoleMapping{nodeRole},
			users:    []ekscontrolplanev1.UserMapping{alice},
			wantARNs: []string{nodeRole.RoleARN, alice.UserARN},
		},
		{
			name:        "idempotent — existing matching CR is kept, no new CR created",
			preexisting: []crclient.Object{capaBobCR},
			users:       []ekscontrolplanev1.UserMapping{bob},
			wantARNs:    []string{bob.UserARN},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			scheme := runtime.NewScheme()
			_ = iamauthv1.AddToScheme(scheme)

			c := fake.NewClientBuilder().WithScheme(scheme).
				WithObjects(tc.preexisting...).Build()

			backend, err := NewBackend(BackendTypeCRD, c)
			g.Expect(err).To(BeNil())

			g.Expect(backend.ReconcileMappings(tc.roles, tc.users)).To(Succeed())

			list := &iamauthv1.IAMIdentityMappingList{}
			g.Expect(c.List(context.TODO(), list)).To(Succeed())

			gotARNs := make([]string, 0, len(list.Items))
			for _, cr := range list.Items {
				gotARNs = append(gotARNs, cr.Spec.ARN)
			}
			for _, arn := range tc.wantARNs {
				g.Expect(gotARNs).To(ContainElement(arn), "expected ARN %s to be present", arn)
			}
			for _, arn := range tc.wantARNsAbsent {
				g.Expect(gotARNs).ToNot(ContainElement(arn), "expected ARN %s to be absent", arn)
			}
		})
	}
}

// TestReconcileMappingsCRDIdempotent proves that calling ReconcileMappings
// twice with the same input does not create duplicate CRs. This complements
// the sub-case above by asserting the exact count after a second call.
func TestReconcileMappingsCRDIdempotent(t *testing.T) {
	g := NewGomegaWithT(t)

	alice := ekscontrolplanev1.UserMapping{
		UserARN: "arn:aws:iam::000000000000:user/Alice",
		KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
			UserName: "alice",
			Groups:   []string{"system:masters"},
		},
	}

	scheme := runtime.NewScheme()
	_ = iamauthv1.AddToScheme(scheme)
	c := fake.NewClientBuilder().WithScheme(scheme).Build()

	backend, err := NewBackend(BackendTypeCRD, c)
	g.Expect(err).To(BeNil())

	g.Expect(backend.ReconcileMappings(nil, []ekscontrolplanev1.UserMapping{alice})).To(Succeed())
	g.Expect(backend.ReconcileMappings(nil, []ekscontrolplanev1.UserMapping{alice})).To(Succeed())

	list := &iamauthv1.IAMIdentityMappingList{}
	g.Expect(c.List(context.TODO(), list)).To(Succeed())
	g.Expect(list.Items).To(HaveLen(1))
	g.Expect(list.Items[0].Spec.ARN).To(Equal(alice.UserARN))
}

// TestReconcileMappingsCRDInvalidARN proves the validate-before-mutate
// contract: an invalid input returns an error and the pre-existing CRs remain
// untouched.
func TestReconcileMappingsCRDInvalidARN(t *testing.T) {
	g := NewGomegaWithT(t)

	scheme := runtime.NewScheme()
	_ = iamauthv1.AddToScheme(scheme)

	existing := &iamauthv1.IAMIdentityMapping{
		ObjectMeta: metav1.ObjectMeta{
			Name:         "capa-iamauth-preexisting",
			Namespace:    metav1.NamespaceSystem,
			GenerateName: capaGenerateName,
		},
		Spec: iamauthv1.IAMIdentityMappingSpec{
			ARN:      "arn:aws:iam::000000000000:user/Alice",
			Username: "alice",
			Groups:   []string{"system:masters"},
		},
	}
	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(existing).Build()

	backend, err := NewBackend(BackendTypeCRD, c)
	g.Expect(err).To(BeNil())

	bogus := []ekscontrolplanev1.RoleMapping{{
		RoleARN: "not-a-valid-arn",
		KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
			UserName: "x",
			Groups:   []string{"g"},
		},
	}}
	err = backend.ReconcileMappings(bogus, nil)
	g.Expect(err).To(HaveOccurred())

	// Pre-existing CR still there — validation failure must not mutate state.
	list := &iamauthv1.IAMIdentityMappingList{}
	g.Expect(c.List(context.TODO(), list)).To(Succeed())
	g.Expect(list.Items).To(HaveLen(1))
	g.Expect(list.Items[0].Name).To(Equal("capa-iamauth-preexisting"))
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
