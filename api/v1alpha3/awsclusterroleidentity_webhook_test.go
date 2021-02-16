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

package v1alpha3

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAWSClusterRoleValidateCreate(t *testing.T) {
	tests := []struct {
		name      string
		identity  *AWSClusterRoleIdentity
		wantError bool
	}{
		{
			name: "do not allow nil sourceIdentityRef",
			identity: &AWSClusterRoleIdentity{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
			},
			wantError: true,
		},
		{
			name: "successfully create AWSClusterRoleIdentity",
			identity: &AWSClusterRoleIdentity{
				ObjectMeta: metav1.ObjectMeta{
					Name: "role",
				},
				Spec: AWSClusterRoleIdentitySpec{
					SourceIdentityRef: &AWSIdentityReference{
						Name: "another-role",
						Kind: ClusterRoleIdentityKind,
					},
				},
			},
			wantError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			identity := tt.identity.DeepCopy()
			identity.TypeMeta = metav1.TypeMeta{
				APIVersion: GroupVersion.String(),
				Kind:       "AWSClusterRoleIdentity",
			}
			ctx := context.TODO()
			if err := testEnv.Create(ctx, identity); (err != nil) != tt.wantError {
				t.Errorf("ValidateCreate() error = %v, wantErr %v", err, tt.wantError)
			}
			testEnv.Delete(ctx, identity)
		})
	}
}

func TestAWSClusterRoleValidateUpdate(t *testing.T) {
	roleIdentity := &AWSClusterRoleIdentity{
		TypeMeta: metav1.TypeMeta{
			APIVersion: GroupVersion.String(),
			Kind:       string(ClusterRoleIdentityKind),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "role",
		},
		Spec: AWSClusterRoleIdentitySpec{
			SourceIdentityRef: &AWSIdentityReference{
				Name: "another-role",
				Kind: ClusterRoleIdentityKind,
			},
		},
	}

	ctx := context.TODO()
	defer testEnv.Delete(ctx, roleIdentity)

	if err := testEnv.Create(ctx, roleIdentity); err != nil {
		t.Errorf("roleIdentity creation failed %v", err)
	}

	roleIdentity.Spec = AWSClusterRoleIdentitySpec{}
	if err := testEnv.Update(ctx, roleIdentity); err == nil {
		t.Errorf("roleIdentity is updated with nil sourceIdentityRef %v", err)
	}
}
