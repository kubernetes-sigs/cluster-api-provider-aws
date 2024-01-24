//go:build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1beta2

import (
	"k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/cluster-api/api/v1beta1"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSConfiguration) DeepCopyInto(out *AWSConfiguration) {
	*out = *in
	if in.PrivateLinkConfiguration != nil {
		in, out := &in.PrivateLinkConfiguration, &out.PrivateLinkConfiguration
		*out = new(PrivateLinkConfiguration)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSConfiguration.
func (in *AWSConfiguration) DeepCopy() *AWSConfiguration {
	if in == nil {
		return nil
	}
	out := new(AWSConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSRolesRef) DeepCopyInto(out *AWSRolesRef) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSRolesRef.
func (in *AWSRolesRef) DeepCopy() *AWSRolesRef {
	if in == nil {
		return nil
	}
	out := new(AWSRolesRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PrivateLinkConfiguration) DeepCopyInto(out *PrivateLinkConfiguration) {
	*out = *in
	if in.Principals != nil {
		in, out := &in.Principals, &out.Principals
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PrivateLinkConfiguration.
func (in *PrivateLinkConfiguration) DeepCopy() *PrivateLinkConfiguration {
	if in == nil {
		return nil
	}
	out := new(PrivateLinkConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ROSAControlPlane) DeepCopyInto(out *ROSAControlPlane) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ROSAControlPlane.
func (in *ROSAControlPlane) DeepCopy() *ROSAControlPlane {
	if in == nil {
		return nil
	}
	out := new(ROSAControlPlane)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ROSAControlPlane) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ROSAControlPlaneList) DeepCopyInto(out *ROSAControlPlaneList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ROSAControlPlane, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ROSAControlPlaneList.
func (in *ROSAControlPlaneList) DeepCopy() *ROSAControlPlaneList {
	if in == nil {
		return nil
	}
	out := new(ROSAControlPlaneList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ROSAControlPlaneList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RosaControlPlaneSpec) DeepCopyInto(out *RosaControlPlaneSpec) {
	*out = *in
	if in.Subnets != nil {
		in, out := &in.Subnets, &out.Subnets
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.AvailabilityZones != nil {
		in, out := &in.AvailabilityZones, &out.AvailabilityZones
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.MachineCIDR != nil {
		in, out := &in.MachineCIDR, &out.MachineCIDR
		*out = new(string)
		**out = **in
	}
	if in.Region != nil {
		in, out := &in.Region, &out.Region
		*out = new(string)
		**out = **in
	}
	if in.Version != nil {
		in, out := &in.Version, &out.Version
		*out = new(string)
		**out = **in
	}
	out.ControlPlaneEndpoint = in.ControlPlaneEndpoint
	out.RolesRef = in.RolesRef
	if in.OIDCID != nil {
		in, out := &in.OIDCID, &out.OIDCID
		*out = new(string)
		**out = **in
	}
	if in.AccountID != nil {
		in, out := &in.AccountID, &out.AccountID
		*out = new(string)
		**out = **in
	}
	if in.CreatorARN != nil {
		in, out := &in.CreatorARN, &out.CreatorARN
		*out = new(string)
		**out = **in
	}
	if in.InstallerRoleARN != nil {
		in, out := &in.InstallerRoleARN, &out.InstallerRoleARN
		*out = new(string)
		**out = **in
	}
	if in.SupportRoleARN != nil {
		in, out := &in.SupportRoleARN, &out.SupportRoleARN
		*out = new(string)
		**out = **in
	}
	if in.WorkerRoleARN != nil {
		in, out := &in.WorkerRoleARN, &out.WorkerRoleARN
		*out = new(string)
		**out = **in
	}
	if in.CredentialsSecretRef != nil {
		in, out := &in.CredentialsSecretRef, &out.CredentialsSecretRef
		*out = new(v1.LocalObjectReference)
		**out = **in
	}
	in.AWS.DeepCopyInto(&out.AWS)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RosaControlPlaneSpec.
func (in *RosaControlPlaneSpec) DeepCopy() *RosaControlPlaneSpec {
	if in == nil {
		return nil
	}
	out := new(RosaControlPlaneSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RosaControlPlaneStatus) DeepCopyInto(out *RosaControlPlaneStatus) {
	*out = *in
	if in.ExternalManagedControlPlane != nil {
		in, out := &in.ExternalManagedControlPlane, &out.ExternalManagedControlPlane
		*out = new(bool)
		**out = **in
	}
	if in.FailureMessage != nil {
		in, out := &in.FailureMessage, &out.FailureMessage
		*out = new(string)
		**out = **in
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make(v1beta1.Conditions, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RosaControlPlaneStatus.
func (in *RosaControlPlaneStatus) DeepCopy() *RosaControlPlaneStatus {
	if in == nil {
		return nil
	}
	out := new(RosaControlPlaneStatus)
	in.DeepCopyInto(out)
	return out
}
