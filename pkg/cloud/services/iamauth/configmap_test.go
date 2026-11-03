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
	"testing"

	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/yaml"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
)

var (
	existingNodeRoleMap = `
    - groups:
      - system:bootstrappers
      - system:nodes
      rolearn: arn:aws:iam::000000000000:role/KubernetesNode
      username: system:node:{{EC2PrivateDNSName}}
`

	existingUserMap = `
    - userarn: arn:aws:iam::000000000000:user/Alice
      username: alice
      groups:
      - system:masters
`
)

func TestAddRoleMappingCM(t *testing.T) {
	testCases := []struct {
		name                  string
		existingAuthConfigMap *corev1.ConfigMap
		roleToMap             ekscontrolplanev1.RoleMapping
		expectedRoleMaps      []ekscontrolplanev1.RoleMapping
		expectError           bool
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
			expectedRoleMaps: []ekscontrolplanev1.RoleMapping{
				{
					RoleARN: "arn:aws:iam::000000000000:role/KubernetesNode",
					KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
						UserName: "system:node:{{EC2PrivateDNSName}}",
						Groups:   []string{systemBootstrappersGroup, systemNodesGroup},
					},
				},
			},
			expectError: false,
		},
		{
			name: "existing mapping, add different mapping",
			roleToMap: ekscontrolplanev1.RoleMapping{
				RoleARN: "arn:aws:iam::000000000000:role/KubernetesAdmin",
				KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
					UserName: "admin:{{SessionName}}",
					Groups:   []string{"system:masters"},
				},
			},
			expectedRoleMaps: []ekscontrolplanev1.RoleMapping{
				{
					RoleARN: "arn:aws:iam::000000000000:role/KubernetesNode",
					KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
						UserName: "system:node:{{EC2PrivateDNSName}}",
						Groups:   []string{systemBootstrappersGroup, systemNodesGroup},
					},
				},
				{
					RoleARN: "arn:aws:iam::000000000000:role/KubernetesAdmin",
					KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
						UserName: "admin:{{SessionName}}",
						Groups:   []string{"system:masters"},
					},
				},
			},
			expectError:           false,
			existingAuthConfigMap: createFakeConfigMap(existingNodeRoleMap, ""),
		},
		{
			name: "existing mapping, add same mapping",
			roleToMap: ekscontrolplanev1.RoleMapping{
				RoleARN: "arn:aws:iam::000000000000:role/KubernetesNode",
				KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
					UserName: "system:node:{{EC2PrivateDNSName}}",
					Groups:   []string{systemBootstrappersGroup, systemNodesGroup},
				},
			},
			expectedRoleMaps: []ekscontrolplanev1.RoleMapping{
				{
					RoleARN: "arn:aws:iam::000000000000:role/KubernetesNode",
					KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
						UserName: "system:node:{{EC2PrivateDNSName}}",
						Groups:   []string{systemBootstrappersGroup, systemNodesGroup},
					},
				},
			},
			expectError:           false,
			existingAuthConfigMap: createFakeConfigMap(existingNodeRoleMap, ""),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			var client crclient.Client
			if tc.existingAuthConfigMap == nil {
				client = fake.NewClientBuilder().Build()
			} else {
				client = fake.NewClientBuilder().WithObjects(tc.existingAuthConfigMap).Build()
			}
			backend, err := NewBackend(BackendTypeConfigMap, client)
			g.Expect(err).To(BeNil())

			err = backend.MapRole(tc.roleToMap)
			if tc.expectError {
				g.Expect(err).ToNot(BeNil())
				return
			}

			g.Expect(err).To(BeNil())

			key := types.NamespacedName{
				Name:      "aws-auth",
				Namespace: "kube-system",
			}

			cm := &corev1.ConfigMap{}

			err = client.Get(context.TODO(), key, cm)
			g.Expect(err).To(BeNil())

			g.Expect(cm.Name).To(Equal("aws-auth"))
			g.Expect(cm.Namespace).To(Equal("kube-system"))
			g.Expect(cm.Data).ToNot(BeNil())

			actualRoleMappings, roleMappingsFound := cm.Data["mapRoles"]
			if len(tc.expectedRoleMaps) == 0 {
				g.Expect(roleMappingsFound).To(BeFalse())
			} else {
				roles := []ekscontrolplanev1.RoleMapping{}
				err := yaml.Unmarshal([]byte(actualRoleMappings), &roles)
				g.Expect(err).To(BeNil())
				g.Expect(len(roles)).To(Equal(len(tc.expectedRoleMaps)))
				//TODO: we may need to do a better match
				bothMatch := cmp.Equal(roles, tc.expectedRoleMaps)
				g.Expect(bothMatch).To(BeTrue())
			}

			_, userMappingsFound := cm.Data["mapUsers"]
			g.Expect(userMappingsFound).To(BeFalse())
		})
	}
}

func TestAddUserMappingCM(t *testing.T) {
	testCases := []struct {
		name                  string
		existingAuthConfigMap *corev1.ConfigMap
		userToMap             ekscontrolplanev1.UserMapping
		expectedUsersMap      []ekscontrolplanev1.UserMapping
		expectError           bool
	}{
		{
			name: "no existing user mappings, add user mapping",
			userToMap: ekscontrolplanev1.UserMapping{
				UserARN: "arn:aws:iam::000000000000:user/Alice",
				KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
					UserName: "alice",
					Groups:   []string{"system:masters"},
				},
			},
			expectedUsersMap: []ekscontrolplanev1.UserMapping{
				{
					UserARN: "arn:aws:iam::000000000000:user/Alice",
					KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
						UserName: "alice",
						Groups:   []string{"system:masters"},
					},
				},
			},
			expectError: false,
		},
		{
			name: "existing user mapping, add different user mapping",
			userToMap: ekscontrolplanev1.UserMapping{
				UserARN: "arn:aws:iam::000000000000:user/Bob",
				KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
					UserName: "bob",
					Groups:   []string{"system:masters"},
				},
			},
			expectedUsersMap: []ekscontrolplanev1.UserMapping{
				{
					UserARN: "arn:aws:iam::000000000000:user/Alice",
					KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
						UserName: "alice",
						Groups:   []string{"system:masters"},
					},
				},
				{
					UserARN: "arn:aws:iam::000000000000:user/Bob",
					KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
						UserName: "bob",
						Groups:   []string{"system:masters"},
					},
				},
			},
			expectError:           false,
			existingAuthConfigMap: createFakeConfigMap("", existingUserMap),
		},
		{
			name: "existing user mapping, add same user mapping",
			userToMap: ekscontrolplanev1.UserMapping{
				UserARN: "arn:aws:iam::000000000000:user/Alice",
				KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
					UserName: "alice",
					Groups:   []string{"system:masters"},
				},
			},
			expectedUsersMap: []ekscontrolplanev1.UserMapping{
				{
					UserARN: "arn:aws:iam::000000000000:user/Alice",
					KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
						UserName: "alice",
						Groups:   []string{"system:masters"},
					},
				},
			},
			expectError:           false,
			existingAuthConfigMap: createFakeConfigMap("", existingUserMap),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			var client crclient.Client
			if tc.existingAuthConfigMap == nil {
				client = fake.NewClientBuilder().Build()
			} else {
				client = fake.NewClientBuilder().WithObjects(tc.existingAuthConfigMap).Build()
			}
			backend, err := NewBackend(BackendTypeConfigMap, client)
			g.Expect(err).To(BeNil())

			err = backend.MapUser(tc.userToMap)
			if tc.expectError {
				g.Expect(err).ToNot(BeNil())
				return
			}

			g.Expect(err).To(BeNil())

			key := types.NamespacedName{
				Name:      "aws-auth",
				Namespace: "kube-system",
			}

			cm := &corev1.ConfigMap{}

			err = client.Get(context.TODO(), key, cm)
			g.Expect(err).To(BeNil())

			g.Expect(cm.Name).To(Equal("aws-auth"))
			g.Expect(cm.Namespace).To(Equal("kube-system"))
			g.Expect(cm.Data).ToNot(BeNil())

			actualUserMappings, userMappingsFound := cm.Data["mapUsers"]
			if len(tc.expectedUsersMap) == 0 {
				g.Expect(userMappingsFound).To(BeFalse())
			} else {
				users := []ekscontrolplanev1.UserMapping{}
				err := yaml.Unmarshal([]byte(actualUserMappings), &users)
				g.Expect(err).To(BeNil())
				g.Expect(len(users)).To(Equal(len(tc.expectedUsersMap)))
				//TODO: we may need to do a better match
				bothMatch := cmp.Equal(users, tc.expectedUsersMap)
				g.Expect(bothMatch).To(BeTrue())
			}

			_, roleMappingsFound := cm.Data["mapRoles"]
			g.Expect(roleMappingsFound).To(BeFalse())
		})
	}
}

// TestReconcileMappingsCM exercises the desired-state semantics of
// configMapBackend.ReconcileMappings. Unlike the append-only MapRole / MapUser
// methods, ReconcileMappings must replace the aws-auth ConfigMap contents
// wholesale — including deleting entries that are absent from the input.
//
// Scenarios:
//  1. Empty desired set on a populated backend removes everything.
//  2. Add a user to an empty backend.
//  3. Replace {alice, bob} with {alice} — bob must be removed.
//  4. Node roles preserved when included in the desired set.
//  5. Idempotency: second call with the same input yields the same state.
//  6. Invalid ARN fails without mutating state.
func TestReconcileMappingsCM(t *testing.T) {
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

	tests := []struct {
		name        string
		existing    *corev1.ConfigMap
		roles       []ekscontrolplanev1.RoleMapping
		users       []ekscontrolplanev1.UserMapping
		wantRoles   []ekscontrolplanev1.RoleMapping
		wantUsers   []ekscontrolplanev1.UserMapping
		expectError bool
	}{
		{
			name:      "empty desired set on populated backend removes everything",
			existing:  cmWith(nil, []ekscontrolplanev1.UserMapping{alice, bob}),
			roles:     []ekscontrolplanev1.RoleMapping{},
			users:     []ekscontrolplanev1.UserMapping{},
			wantRoles: nil,
			wantUsers: nil,
		},
		{
			name:      "add user to empty backend",
			existing:  nil,
			roles:     nil,
			users:     []ekscontrolplanev1.UserMapping{alice},
			wantRoles: nil,
			wantUsers: []ekscontrolplanev1.UserMapping{alice},
		},
		{
			name:      "replace {alice,bob} with {alice} removes bob",
			existing:  cmWith(nil, []ekscontrolplanev1.UserMapping{alice, bob}),
			roles:     nil,
			users:     []ekscontrolplanev1.UserMapping{alice},
			wantRoles: nil,
			wantUsers: []ekscontrolplanev1.UserMapping{alice},
		},
		{
			name:      "node role preserved when included in desired set",
			existing:  nil,
			roles:     []ekscontrolplanev1.RoleMapping{nodeRole},
			users:     []ekscontrolplanev1.UserMapping{alice},
			wantRoles: []ekscontrolplanev1.RoleMapping{nodeRole},
			wantUsers: []ekscontrolplanev1.UserMapping{alice},
		},
		{
			name:      "idempotent — second call with same input produces same state",
			existing:  cmWith([]ekscontrolplanev1.RoleMapping{nodeRole}, []ekscontrolplanev1.UserMapping{alice}),
			roles:     []ekscontrolplanev1.RoleMapping{nodeRole},
			users:     []ekscontrolplanev1.UserMapping{alice},
			wantRoles: []ekscontrolplanev1.RoleMapping{nodeRole},
			wantUsers: []ekscontrolplanev1.UserMapping{alice},
		},
		{
			name:     "invalid ARN in role mapping fails without mutating state",
			existing: cmWith(nil, []ekscontrolplanev1.UserMapping{alice}),
			roles: []ekscontrolplanev1.RoleMapping{{
				RoleARN: "bogus",
				KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
					UserName: "x", Groups: []string{"g"},
				},
			}},
			users:       nil,
			expectError: true,
			// State on failure remains unchanged: alice user only.
			wantRoles: nil,
			wantUsers: []ekscontrolplanev1.UserMapping{alice},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			builder := fake.NewClientBuilder()
			if tc.existing != nil {
				builder = builder.WithObjects(tc.existing)
			}
			c := builder.Build()

			backend, err := NewBackend(BackendTypeConfigMap, c)
			g.Expect(err).To(BeNil())

			err = backend.ReconcileMappings(tc.roles, tc.users)
			if tc.expectError {
				g.Expect(err).ToNot(BeNil())
			} else {
				g.Expect(err).To(BeNil())
			}

			gotRoles, gotUsers := readAuthCM(t, c)
			// nil vs empty slice: normalize before compare.
			if len(tc.wantRoles) == 0 {
				g.Expect(gotRoles).To(BeEmpty())
			} else {
				g.Expect(gotRoles).To(Equal(tc.wantRoles))
			}
			if len(tc.wantUsers) == 0 {
				g.Expect(gotUsers).To(BeEmpty())
			} else {
				g.Expect(gotUsers).To(Equal(tc.wantUsers))
			}
		})
	}
}

// cmWith builds a fake aws-auth ConfigMap seeded with the given mappings.
// Empty / nil slices produce a ConfigMap without the corresponding YAML key,
// mirroring what saveAuthConfig produces at runtime.
func cmWith(roles []ekscontrolplanev1.RoleMapping, users []ekscontrolplanev1.UserMapping) *corev1.ConfigMap {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: configMapNS,
			UID:       "1234567",
		},
		Data: map[string]string{},
	}
	if len(roles) > 0 {
		b, err := yaml.Marshal(roles)
		if err != nil {
			panic(err)
		}
		cm.Data[roleKey] = string(b)
	}
	if len(users) > 0 {
		b, err := yaml.Marshal(users)
		if err != nil {
			panic(err)
		}
		cm.Data[usersKey] = string(b)
	}
	return cm
}

// readAuthCM reads the aws-auth ConfigMap from the fake client and returns the
// deserialized role and user mapping slices. Returns empty slices if the
// ConfigMap does not exist or the keys are absent.
func readAuthCM(t *testing.T, c crclient.Client) ([]ekscontrolplanev1.RoleMapping, []ekscontrolplanev1.UserMapping) {
	t.Helper()
	cm := &corev1.ConfigMap{}
	err := c.Get(context.TODO(), types.NamespacedName{Name: configMapName, Namespace: configMapNS}, cm)
	if err != nil {
		// Missing ConfigMap → nothing configured.
		return nil, nil
	}
	var roles []ekscontrolplanev1.RoleMapping
	if v, ok := cm.Data[roleKey]; ok {
		if err := yaml.Unmarshal([]byte(v), &roles); err != nil {
			t.Fatalf("unmarshalling roles: %v", err)
		}
	}
	var users []ekscontrolplanev1.UserMapping
	if v, ok := cm.Data[usersKey]; ok {
		if err := yaml.Unmarshal([]byte(v), &users); err != nil {
			t.Fatalf("unmarshalling users: %v", err)
		}
	}
	return roles, users
}

func createFakeConfigMap(roleMappings string, userMappings string) *corev1.ConfigMap {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "aws-auth",
			Namespace: "kube-system",
			UID:       "1234567",
		},
		Data: make(map[string]string),
	}

	if roleMappings != "" {
		cm.Data["mapRoles"] = roleMappings
	}

	if userMappings != "" {
		cm.Data["mapUsers"] = userMappings
	}

	return cm
}
