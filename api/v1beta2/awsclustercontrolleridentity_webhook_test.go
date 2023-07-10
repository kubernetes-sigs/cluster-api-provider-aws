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

func TestAWSClusterControllerIdentityCreationValidation(t *testing.T) {
	tests := []struct {
		name      string
		identity  *AWSClusterControllerIdentity
		wantError bool
	}{
		{
			name: "only allow AWSClusterControllerIdentity creation with name default",
			identity: &AWSClusterControllerIdentity{
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
				},
			},
			wantError: false,
		},
		{
			name: "do not allow AWSClusterControllerIdentity creation with name other than default",
			identity: &AWSClusterControllerIdentity{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			identity := tt.identity.DeepCopy()
			identity.TypeMeta = metav1.TypeMeta{
				APIVersion: GroupVersion.String(),
				Kind:       "AWSClusterControllerIdentity",
			}
			ctx := context.TODO()
			if err := testEnv.Create(ctx, identity); (err != nil) != tt.wantError {
				t.Errorf("ValidateCreate() error = %v, wantErr %v", err, tt.wantError)
			}
			testEnv.Delete(ctx, identity)
		})
	}
}

func TestAWSClusterControllerIdentityLabelSelectorAsSelectorValidation(t *testing.T) {
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
			identity := &AWSClusterControllerIdentity{
				TypeMeta: metav1.TypeMeta{
					APIVersion: GroupVersion.String(),
					Kind:       "AWSClusterControllerIdentity",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: AWSClusterControllerIdentityName,
				},
				Spec: AWSClusterControllerIdentitySpec{
					AWSClusterIdentitySpec: AWSClusterIdentitySpec{
						AllowedNamespaces: &AllowedNamespaces{
							Selector: metav1.LabelSelector{
								MatchLabels: tt.selectors,
							},
						},
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

func TestAWSClusterControllerValidateUpdate(t *testing.T) {
	controllerIdentity := &AWSClusterControllerIdentity{
		TypeMeta: metav1.TypeMeta{
			APIVersion: GroupVersion.String(),
			Kind:       "AWSClusterControllerIdentity",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: AWSClusterControllerIdentityName,
		},
		Spec: AWSClusterControllerIdentitySpec{
			AWSClusterIdentitySpec: AWSClusterIdentitySpec{
				AllowedNamespaces: &AllowedNamespaces{},
			},
		},
	}

	ctx := context.TODO()
	defer testEnv.Delete(ctx, controllerIdentity)

	if err := testEnv.Create(ctx, controllerIdentity); err != nil {
		t.Errorf("controllerIdentity creation failed %v", err)
	}

	tests := []struct {
		name      string
		identity  *AWSClusterControllerIdentity
		wantError bool
	}{
		{
			name: "do not allow any spec changes",
			identity: &AWSClusterControllerIdentity{
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
				},
			},
			wantError: true,
		},
		{
			name: "do not allow name change",
			identity: &AWSClusterControllerIdentity{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
			},
			wantError: true,
		},
		{
			name:      "no error when updating with same object",
			identity:  controllerIdentity,
			wantError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			identity := tt.identity.DeepCopy()
			identity.TypeMeta = metav1.TypeMeta{
				APIVersion: GroupVersion.String(),
				Kind:       "AWSClusterControllerIdentity",
			}
			ctx := context.TODO()
			if err := testEnv.Update(ctx, identity); (err != nil) != tt.wantError {
				t.Errorf("ValidateUpdate() error = %v, wantErr %v", err, tt.wantError)
			}
		})
	}
}

func TestAWSClusterControllerIdentityUpdateValidation(t *testing.T) {
	controllerIdentity := &AWSClusterControllerIdentity{
		TypeMeta: metav1.TypeMeta{
			APIVersion: GroupVersion.String(),
			Kind:       "AWSClusterControllerIdentityName",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: AWSClusterControllerIdentityName,
		},
		Spec: AWSClusterControllerIdentitySpec{
			AWSClusterIdentitySpec: AWSClusterIdentitySpec{
				AllowedNamespaces: &AllowedNamespaces{
					Selector: metav1.LabelSelector{
						MatchLabels: map[string]string{"foo": "bar"},
					},
				},
			},
		},
	}

	ctx := context.TODO()
	defer testEnv.Delete(ctx, controllerIdentity)

	if err := testEnv.Create(ctx, controllerIdentity); err != nil {
		t.Errorf("controllerIdentity creation failed %v", err)
	}

	tests := []struct {
		name      string
		identity  *AWSClusterControllerIdentity
		wantError bool
	}{
		{
			name:      "should not return error for valid selector",
			identity:  controllerIdentity,
			wantError: false,
		},
		{
			name: "should return error for invalid selector",
			identity: &AWSClusterControllerIdentity{
				ObjectMeta: metav1.ObjectMeta{
					Name: AWSClusterControllerIdentityName,
				},
				Spec: AWSClusterControllerIdentitySpec{
					AWSClusterIdentitySpec: AWSClusterIdentitySpec{
						AllowedNamespaces: &AllowedNamespaces{
							Selector: metav1.LabelSelector{
								MatchLabels: map[string]string{"-foo-123": "bar"},
							},
						},
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

func TestAWSClusterControllerIdentityDefault(t *testing.T) {
	g := NewWithT(t)
	tests := []struct {
		name                               string
		beforeAWSClusterControllerIdentity *AWSClusterControllerIdentity
		afterAWSClusterControllerIdentity  *AWSClusterControllerIdentity
	}{
		{
			name: "Set label while creating AWSClusterControllerIdentity if labels are undefined",
			beforeAWSClusterControllerIdentity: &AWSClusterControllerIdentity{
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
				},
			},
			afterAWSClusterControllerIdentity: &AWSClusterControllerIdentity{
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
					Labels: map[string]string{
						clusterv1.ClusterctlMoveHierarchyLabel: "",
					},
				},
			},
		},
		{
			name: "Not update any label while creating AWSClusterControllerIdentity if labels are already defined",
			beforeAWSClusterControllerIdentity: &AWSClusterControllerIdentity{
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
					Labels: map[string]string{
						clusterv1.ClusterctlMoveHierarchyLabel: "abc",
					},
				},
			},
			afterAWSClusterControllerIdentity: &AWSClusterControllerIdentity{
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
					Labels: map[string]string{
						clusterv1.ClusterctlMoveHierarchyLabel: "abc",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			awsClusterControllerIdentity := tt.beforeAWSClusterControllerIdentity.DeepCopy()
			g.Expect(testEnv.Create(ctx, awsClusterControllerIdentity)).To(Succeed())
			g.Expect(len(awsClusterControllerIdentity.ObjectMeta.Labels)).To(Not(Equal(0)))
			g.Expect(awsClusterControllerIdentity.ObjectMeta.Labels[clusterv1.ClusterctlMoveHierarchyLabel]).To(Equal(tt.afterAWSClusterControllerIdentity.ObjectMeta.Labels[clusterv1.ClusterctlMoveHierarchyLabel]))
			g.Expect(testEnv.Delete(ctx, awsClusterControllerIdentity)).To(Succeed())
		})
	}
}
