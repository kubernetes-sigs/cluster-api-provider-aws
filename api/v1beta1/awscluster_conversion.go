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

package v1beta1

import (
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts the v1beta1 AWSCluster receiver to a v1beta1 AWSCluster.
func (src *AWSCluster) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.AWSCluster)

	if err := Convert_v1beta1_AWSCluster_To_v1beta2_AWSCluster(src, dst, nil); err != nil {
		return err
	}
	// Manually restore data.
	restored := &infrav1.AWSCluster{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}

	if restored.Spec.ControlPlaneLoadBalancer != nil {
		if dst.Spec.ControlPlaneLoadBalancer == nil {
			dst.Spec.ControlPlaneLoadBalancer = &infrav1.AWSLoadBalancerSpec{}
		}
		restoreControlPlaneLoadBalancer(restored.Spec.ControlPlaneLoadBalancer, dst.Spec.ControlPlaneLoadBalancer)
	}
	restoreControlPlaneLoadBalancerStatus(&restored.Status.Network.APIServerELB, &dst.Status.Network.APIServerELB)

	dst.Spec.S3Bucket = restored.Spec.S3Bucket
	if restored.Status.Bastion != nil {
		dst.Status.Bastion.InstanceMetadataOptions = restored.Status.Bastion.InstanceMetadataOptions
	}

	return nil
}

// restoreControlPlaneLoadBalancerStatus manually restores the control plane loadbalancer status data.
// Assumes restored and dst are non-nil.
func restoreControlPlaneLoadBalancerStatus(restored, dst *infrav1.LoadBalancer) {
	dst.ARN = restored.ARN
	dst.LoadBalancerType = restored.LoadBalancerType
	dst.ELBAttributes = restored.ELBAttributes
	dst.ELBListeners = restored.ELBListeners
}

// restoreControlPlaneLoadBalancer manually restores the control plane loadbalancer data.
// Assumes restored and dst are non-nil.
func restoreControlPlaneLoadBalancer(restored, dst *infrav1.AWSLoadBalancerSpec) {
	dst.Name = restored.Name
	dst.HealthCheckProtocol = restored.HealthCheckProtocol
	dst.LoadBalancerType = restored.LoadBalancerType
	dst.DisableHostsRewrite = restored.DisableHostsRewrite
	dst.PreserveClientIP = restored.PreserveClientIP
}

// ConvertFrom converts the v1beta1 AWSCluster receiver to a v1beta1 AWSCluster.
func (r *AWSCluster) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.AWSCluster)

	if err := Convert_v1beta2_AWSCluster_To_v1beta1_AWSCluster(src, r, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion.
	if err := utilconversion.MarshalData(src, r); err != nil {
		return err
	}

	return nil
}

// ConvertTo converts the v1beta1 AWSClusterList receiver to a v1beta2 AWSClusterList.
func (src *AWSClusterList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.AWSClusterList)

	return Convert_v1beta1_AWSClusterList_To_v1beta2_AWSClusterList(src, dst, nil)
}

// ConvertFrom converts the v1beta2 AWSClusterList receiver to a v1beta1 AWSClusterList.
func (r *AWSClusterList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.AWSClusterList)

	return Convert_v1beta2_AWSClusterList_To_v1beta1_AWSClusterList(src, r, nil)
}
