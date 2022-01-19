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

package v1beta1

import (
	"testing"

	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/cluster-api-provider-aws/pkg/eks"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	utildefaulting "sigs.k8s.io/cluster-api/util/defaulting"
)

func TestAWSFargateProfileDefault(t *testing.T) {
	fargate := &AWSFargateProfile{ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "default"},
		Spec: FargateProfileSpec{
			ClusterName: "clustername",
		},
	}
	t.Run("for AWSFargateProfile", utildefaulting.DefaultValidateTest(fargate))
	fargate.Default()
	g := NewWithT(t)
	g.Expect(fargate.GetLabels()[clusterv1.ClusterLabelName]).To(BeEquivalentTo(fargate.Spec.ClusterName))
	name, err := eks.GenerateEKSName(fargate.Name, fargate.Namespace, maxProfileNameLength)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(fargate.Spec.ProfileName).To(BeEquivalentTo(name))
}

func TestAWSFargateProfileValidateRoleNameUpdate(t *testing.T) {
	g := NewWithT(t)

	before := &AWSFargateProfile{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: "default",
		},
		Spec: FargateProfileSpec{
			ClusterName: "clustername",
			ProfileName: "profilename",
		},
	}

	invalidRoleNameUpdate := before.DeepCopy()
	invalidRoleNameUpdate.Spec.RoleName = "invalid-role-name"

	defaultRoleNameUpdate := before.DeepCopy()
	defaultRoleNameUpdate.Spec.RoleName = DefaultEKSFargateRole

	validRoleNameUpdate := before.DeepCopy()
	validRoleNameUpdate.Spec.RoleName = "clustername-profilename_fargate"

	beforeWithDifferentRoleName := before.DeepCopy()
	beforeWithDifferentRoleName.Spec.RoleName = "different-role-name"

	tests := []struct {
		name           string
		expectErr      bool
		before         *AWSFargateProfile
		fargateProfile *AWSFargateProfile
	}{
		{
			name:           "update roleName should fail when existing roleName is empty and the new roleName is not the generated name",
			expectErr:      true,
			before:         before,
			fargateProfile: invalidRoleNameUpdate,
		},
		{
			name:           "update roleName should succeed when existing roleName is empty and the new roleName is the default name",
			expectErr:      false,
			before:         before,
			fargateProfile: defaultRoleNameUpdate,
		},
		{
			name:           "update roleName should succeed when existing roleName is empty and the new roleName is the generated name",
			expectErr:      false,
			before:         before,
			fargateProfile: validRoleNameUpdate,
		},
		{
			name:           "update roleName should fail when existing roleName is different from the new roleName",
			expectErr:      true,
			before:         beforeWithDifferentRoleName,
			fargateProfile: validRoleNameUpdate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fargateProfile.ValidateUpdate(tt.before.DeepCopy())
			if tt.expectErr {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(Succeed())
			}
		})
	}
}
