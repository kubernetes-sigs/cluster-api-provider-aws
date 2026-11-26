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

package controllers

import (
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
)

// shouldWaitForControlPlane reports whether reconciliation of obj (a managed machine pool or a
// fargate profile) must wait for the EKS control plane to become ready. It does not wait when obj
// and the control plane are both being deleted: a control plane in deletion never becomes ready,
// and its deletion waits for obj to be gone, so waiting would block both deletions forever.
// Deleting obj only needs its spec. This happens when the objects are restored without status
// (for example from a backup) and the control plane is deleted before it is reconciled.
func shouldWaitForControlPlane(controlPlane *ekscontrolplanev1.AWSManagedControlPlane, obj metav1.Object) bool {
	if controlPlane.Status.Ready {
		return false
	}
	return obj.GetDeletionTimestamp().IsZero() || controlPlane.GetDeletionTimestamp().IsZero()
}

// validateEnclaveEdgeZones returns an error if enclaveOptions is enabled and any
// subnet ID or availability zone resolves to a Local Zone or Wavelength Zone.
// AWS does not support Nitro Enclaves in edge zones.
func validateEnclaveEdgeZones(
	enclaveOptions *infrav1.EnclaveOptions,
	clusterSubnets infrav1.Subnets,
	subnetIDs []string,
	azs []string,
) error {
	if enclaveOptions == nil || !ptr.Deref(enclaveOptions.Enabled, false) {
		return nil
	}

	byID := make(map[string]infrav1.SubnetSpec, len(clusterSubnets))
	byAZ := make(map[string]infrav1.SubnetSpec, len(clusterSubnets))
	for _, s := range clusterSubnets {
		byID[s.GetResourceID()] = s
		byAZ[s.AvailabilityZone] = s
	}

	for _, id := range subnetIDs {
		if s, ok := byID[id]; ok && s.IsEdge() {
			return errors.New("Nitro Enclaves are not supported in Local Zones or Wavelength Zones")
		}
	}

	for _, az := range azs {
		if s, ok := byAZ[az]; ok && s.IsEdge() {
			return errors.New("Nitro Enclaves are not supported in Local Zones or Wavelength Zones")
		}
	}

	return nil
}
