/*
Copyright 2026 The Kubernetes Authors.

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

package controllers

import (
	"testing"

	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
)

func TestShouldWaitForControlPlane(t *testing.T) {
	now := metav1.Now()

	tests := []struct {
		name                 string
		controlPlaneReady    bool
		controlPlaneDeleting bool
		objDeleting          bool
		want                 bool
	}{
		{
			name:              "control plane ready: no wait",
			controlPlaneReady: true,
			want:              false,
		},
		{
			name:                 "control plane ready, both deleting: no wait",
			controlPlaneReady:    true,
			controlPlaneDeleting: true,
			objDeleting:          true,
			want:                 false,
		},
		{
			name: "control plane not ready, nothing deleting: wait",
			want: true,
		},
		{
			name:        "control plane not ready, only obj deleting: wait",
			objDeleting: true,
			want:        true,
		},
		{
			name:                 "control plane not ready, only control plane deleting: wait so obj is not created",
			controlPlaneDeleting: true,
			want:                 true,
		},
		{
			name:                 "control plane not ready, both deleting: no wait so the deletions are not blocked",
			controlPlaneDeleting: true,
			objDeleting:          true,
			want:                 false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			controlPlane := &ekscontrolplanev1.AWSManagedControlPlane{
				Status: ekscontrolplanev1.AWSManagedControlPlaneStatus{Ready: tc.controlPlaneReady},
			}
			pool := &expinfrav1.AWSManagedMachinePool{}
			profile := &expinfrav1.AWSFargateProfile{}
			if tc.controlPlaneDeleting {
				controlPlane.DeletionTimestamp = &now
			}
			if tc.objDeleting {
				pool.DeletionTimestamp = &now
				profile.DeletionTimestamp = &now
			}

			g.Expect(shouldWaitForControlPlane(controlPlane, pool)).To(Equal(tc.want), "AWSManagedMachinePool")
			g.Expect(shouldWaitForControlPlane(controlPlane, profile)).To(Equal(tc.want), "AWSFargateProfile")
		})
	}
}
