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

package v1alpha4

import (
	"sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts the v1alpha4 AWSManagedControlPlane receiver to a v1beta1 AWSManagedControlPlane.
func (r *AWSManagedControlPlane) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.AWSManagedControlPlane)

	return Convert_v1alpha4_AWSManagedControlPlane_To_v1beta1_AWSManagedControlPlane(r, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSManagedControlPlane receiver to a v1alpha4 AWSManagedControlPlane.
func (r *AWSManagedControlPlane) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.AWSManagedControlPlane)

	return Convert_v1beta1_AWSManagedControlPlane_To_v1alpha4_AWSManagedControlPlane(src, r, nil)
}


// ConvertTo converts the v1alpha4 AWSManagedControlPlaneList receiver to a v1beta1 AWSManagedControlPlaneList.
func (r *AWSManagedControlPlaneList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.AWSManagedControlPlane)

	return Convert_v1alpha4_AWSManagedControlPlaneList_To_v1beta1_AWSManagedControlPlaneList(r, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSManagedControlPlaneList receiver to a v1alpha4 AWSManagedControlPlaneList.
func (r *AWSManagedControlPlaneList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.AWSManagedControlPlaneList)

	return Convert_v1beta1_AWSManagedControlPlaneListe_To_v1alpha4_AWSManagedControlPlaneList(src, r, nil)
}
