/*
Copyright 2022 The Kubernetes Authors.

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

package v1beta2

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	clusterv1 "sigs.k8s.io/cluster-api/cmd/clusterctl/api/v1alpha3"
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

func TestCreateAWSClusterRoleIdentityLabelSelectorAsSelectorValidation(t *testing.T) {
	tests := []struct {
		name      string
		selectors map[string]string
		wantError bool
	}{
		{
			name:      "should not return error for valid selector",
			selectors: map[string]string{"foo": "bar"},
			wantError: false,
		},
		{
			name:      "should return error for invalid selector",
			selectors: map[string]string{"-123-foo": "bar"},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			identity := &AWSClusterRoleIdentity{
				ObjectMeta: metav1.ObjectMeta{
					Name: "role",
				},
				Spec: AWSClusterRoleIdentitySpec{
					AWSClusterIdentitySpec: AWSClusterIdentitySpec{
						AllowedNamespaces: &AllowedNamespaces{
							Selector: metav1.LabelSelector{
								MatchLabels: tt.selectors,
							},
						},
					},
					SourceIdentityRef: &AWSIdentityReference{
						Name: "another-role",
						Kind: ClusterRoleIdentityKind,
					},
				},
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

func TestAWSClusterRoleIdentityUpdateValidation(t *testing.T) {
	roleIdentity := &AWSClusterRoleIdentity{
		TypeMeta: metav1.TypeMeta{
			APIVersion: GroupVersion.String(),
			Kind:       string(ClusterRoleIdentityKind),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "role",
		},
		Spec: AWSClusterRoleIdentitySpec{
			AWSClusterIdentitySpec: AWSClusterIdentitySpec{
				AllowedNamespaces: &AllowedNamespaces{
					Selector: metav1.LabelSelector{
						MatchLabels: map[string]string{"foo": "bar"},
					},
				},
			},
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

	tests := []struct {
		name      string
		identity  *AWSClusterRoleIdentity
		wantError bool
	}{
		{
			name:      "should not return error for valid selector",
			identity:  roleIdentity,
			wantError: false,
		},
		{
			name: "should return error for invalid selector",
			identity: &AWSClusterRoleIdentity{
				ObjectMeta: metav1.ObjectMeta{
					Name: "role",
				},
				Spec: AWSClusterRoleIdentitySpec{
					AWSClusterIdentitySpec: AWSClusterIdentitySpec{
						AllowedNamespaces: &AllowedNamespaces{
							Selector: metav1.LabelSelector{
								MatchLabels: map[string]string{"-foo-123": "bar"},
							},
						},
					},
					SourceIdentityRef: &AWSIdentityReference{
						Name: "another-role",
						Kind: ClusterRoleIdentityKind,
					},
				},
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			identity := tt.identity.DeepCopy()
			ctx := context.TODO()
			if err := testEnv.Update(ctx, identity); (err != nil) != tt.wantError {
				t.Errorf("ValidateUpdate() error = %v, wantErr %v", err, tt.wantError)
			}
		})
	}
}

func TestAWSClusterRoleIdentityDefault(t *testing.T) {
	g := NewWithT(t)
	tests := []struct {
		name                         string
		beforeAWSClusterRoleIdentity *AWSClusterRoleIdentity
		afterAWSClusterRoleIdentity  *AWSClusterRoleIdentity
	}{
		{
			name: "Set label while creating AWSClusterRoleIdentity if labels are undefined",
			beforeAWSClusterRoleIdentity: &AWSClusterRoleIdentity{
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
				},
				Spec: AWSClusterRoleIdentitySpec{
					SourceIdentityRef: &AWSIdentityReference{
						Name: "another-role",
						Kind: ClusterRoleIdentityKind,
					},
				},
			},
			afterAWSClusterRoleIdentity: &AWSClusterRoleIdentity{
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
					Labels: map[string]string{
						clusterv1.ClusterctlMoveHierarchyLabel: "",
					},
				},
				Spec: AWSClusterRoleIdentitySpec{
					SourceIdentityRef: &AWSIdentityReference{
						Name: "another-role",
						Kind: ClusterRoleIdentityKind,
					},
				},
			},
		},
		{
			name: "Not update any label while creating AWSClusterRoleIdentity if labels are already defined",
			beforeAWSClusterRoleIdentity: &AWSClusterRoleIdentity{
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
					Labels: map[string]string{
						clusterv1.ClusterctlMoveHierarchyLabel: "abc",
					},
				},
				Spec: AWSClusterRoleIdentitySpec{
					SourceIdentityRef: &AWSIdentityReference{
						Name: "another-role",
						Kind: ClusterRoleIdentityKind,
					},
				},
			},
			afterAWSClusterRoleIdentity: &AWSClusterRoleIdentity{
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
					Labels: map[string]string{
						clusterv1.ClusterctlMoveHierarchyLabel: "abc",
					},
				},
				Spec: AWSClusterRoleIdentitySpec{
					SourceIdentityRef: &AWSIdentityReference{
						Name: "another-role",
						Kind: ClusterRoleIdentityKind,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			awsClusterRoleIdentity := tt.beforeAWSClusterRoleIdentity.DeepCopy()
			g.Expect(testEnv.Create(ctx, awsClusterRoleIdentity)).To(Succeed())
			g.Expect(len(awsClusterRoleIdentity.ObjectMeta.Labels)).To(Not(Equal(0)))
			g.Expect(awsClusterRoleIdentity.ObjectMeta.Labels[clusterv1.ClusterctlMoveHierarchyLabel]).To(Equal(tt.afterAWSClusterRoleIdentity.ObjectMeta.Labels[clusterv1.ClusterctlMoveHierarchyLabel]))
			g.Expect(testEnv.Delete(ctx, awsClusterRoleIdentity)).To(Succeed())
		})
	}
}
