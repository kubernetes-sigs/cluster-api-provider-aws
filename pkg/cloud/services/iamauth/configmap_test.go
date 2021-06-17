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
	"testing"

	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/yaml"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1alpha4"
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
					Groups:   []string{"system:bootstrappers", "system:nodes"},
				},
			},
			expectedRoleMaps: []ekscontrolplanev1.RoleMapping{
				{
					RoleARN: "arn:aws:iam::000000000000:role/KubernetesNode",
					KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
						UserName: "system:node:{{EC2PrivateDNSName}}",
						Groups:   []string{"system:bootstrappers", "system:nodes"},
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
						Groups:   []string{"system:bootstrappers", "system:nodes"},
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
					Groups:   []string{"system:bootstrappers", "system:nodes"},
				},
			},
			expectedRoleMaps: []ekscontrolplanev1.RoleMapping{
				{
					RoleARN: "arn:aws:iam::000000000000:role/KubernetesNode",
					KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
						UserName: "system:node:{{EC2PrivateDNSName}}",
						Groups:   []string{"system:bootstrappers", "system:nodes"},
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
				bothMatch := reflect.DeepEqual(roles, tc.expectedRoleMaps)
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
				bothMatch := reflect.DeepEqual(users, tc.expectedUsersMap)
				g.Expect(bothMatch).To(BeTrue())
			}

			_, roleMappingsFound := cm.Data["mapRoles"]
			g.Expect(roleMappingsFound).To(BeFalse())
		})
	}
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
