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
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/cluster-api/api/v1beta1"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DiskSetup) DeepCopyInto(out *DiskSetup) {
	*out = *in
	if in.Partitions != nil {
		in, out := &in.Partitions, &out.Partitions
		*out = make([]Partition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Filesystems != nil {
		in, out := &in.Filesystems, &out.Filesystems
		*out = make([]Filesystem, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DiskSetup.
func (in *DiskSetup) DeepCopy() *DiskSetup {
	if in == nil {
		return nil
	}
	out := new(DiskSetup)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EKSConfig) DeepCopyInto(out *EKSConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EKSConfig.
func (in *EKSConfig) DeepCopy() *EKSConfig {
	if in == nil {
		return nil
	}
	out := new(EKSConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EKSConfig) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EKSConfigList) DeepCopyInto(out *EKSConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]EKSConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EKSConfigList.
func (in *EKSConfigList) DeepCopy() *EKSConfigList {
	if in == nil {
		return nil
	}
	out := new(EKSConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EKSConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EKSConfigSpec) DeepCopyInto(out *EKSConfigSpec) {
	*out = *in
	if in.KubeletExtraArgs != nil {
		in, out := &in.KubeletExtraArgs, &out.KubeletExtraArgs
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.ContainerRuntime != nil {
		in, out := &in.ContainerRuntime, &out.ContainerRuntime
		*out = new(string)
		**out = **in
	}
	if in.DNSClusterIP != nil {
		in, out := &in.DNSClusterIP, &out.DNSClusterIP
		*out = new(string)
		**out = **in
	}
	if in.DockerConfigJSON != nil {
		in, out := &in.DockerConfigJSON, &out.DockerConfigJSON
		*out = new(string)
		**out = **in
	}
	if in.APIRetryAttempts != nil {
		in, out := &in.APIRetryAttempts, &out.APIRetryAttempts
		*out = new(int)
		**out = **in
	}
	if in.PauseContainer != nil {
		in, out := &in.PauseContainer, &out.PauseContainer
		*out = new(PauseContainer)
		**out = **in
	}
	if in.UseMaxPods != nil {
		in, out := &in.UseMaxPods, &out.UseMaxPods
		*out = new(bool)
		**out = **in
	}
	if in.ServiceIPV6Cidr != nil {
		in, out := &in.ServiceIPV6Cidr, &out.ServiceIPV6Cidr
		*out = new(string)
		**out = **in
	}
	if in.PreBootstrapCommands != nil {
		in, out := &in.PreBootstrapCommands, &out.PreBootstrapCommands
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.PostBootstrapCommands != nil {
		in, out := &in.PostBootstrapCommands, &out.PostBootstrapCommands
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.BootstrapCommandOverride != nil {
		in, out := &in.BootstrapCommandOverride, &out.BootstrapCommandOverride
		*out = new(string)
		**out = **in
	}
	if in.Files != nil {
		in, out := &in.Files, &out.Files
		*out = make([]File, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.DiskSetup != nil {
		in, out := &in.DiskSetup, &out.DiskSetup
		*out = new(DiskSetup)
		(*in).DeepCopyInto(*out)
	}
	if in.Mounts != nil {
		in, out := &in.Mounts, &out.Mounts
		*out = make([]MountPoints, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = make(MountPoints, len(*in))
				copy(*out, *in)
			}
		}
	}
	if in.Users != nil {
		in, out := &in.Users, &out.Users
		*out = make([]User, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.NTP != nil {
		in, out := &in.NTP, &out.NTP
		*out = new(NTP)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EKSConfigSpec.
func (in *EKSConfigSpec) DeepCopy() *EKSConfigSpec {
	if in == nil {
		return nil
	}
	out := new(EKSConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EKSConfigStatus) DeepCopyInto(out *EKSConfigStatus) {
	*out = *in
	if in.DataSecretName != nil {
		in, out := &in.DataSecretName, &out.DataSecretName
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
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EKSConfigStatus.
func (in *EKSConfigStatus) DeepCopy() *EKSConfigStatus {
	if in == nil {
		return nil
	}
	out := new(EKSConfigStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EKSConfigTemplate) DeepCopyInto(out *EKSConfigTemplate) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EKSConfigTemplate.
func (in *EKSConfigTemplate) DeepCopy() *EKSConfigTemplate {
	if in == nil {
		return nil
	}
	out := new(EKSConfigTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EKSConfigTemplate) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EKSConfigTemplateList) DeepCopyInto(out *EKSConfigTemplateList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]EKSConfigTemplate, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EKSConfigTemplateList.
func (in *EKSConfigTemplateList) DeepCopy() *EKSConfigTemplateList {
	if in == nil {
		return nil
	}
	out := new(EKSConfigTemplateList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EKSConfigTemplateList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EKSConfigTemplateResource) DeepCopyInto(out *EKSConfigTemplateResource) {
	*out = *in
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EKSConfigTemplateResource.
func (in *EKSConfigTemplateResource) DeepCopy() *EKSConfigTemplateResource {
	if in == nil {
		return nil
	}
	out := new(EKSConfigTemplateResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EKSConfigTemplateSpec) DeepCopyInto(out *EKSConfigTemplateSpec) {
	*out = *in
	in.Template.DeepCopyInto(&out.Template)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EKSConfigTemplateSpec.
func (in *EKSConfigTemplateSpec) DeepCopy() *EKSConfigTemplateSpec {
	if in == nil {
		return nil
	}
	out := new(EKSConfigTemplateSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *File) DeepCopyInto(out *File) {
	*out = *in
	if in.ContentFrom != nil {
		in, out := &in.ContentFrom, &out.ContentFrom
		*out = new(FileSource)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new File.
func (in *File) DeepCopy() *File {
	if in == nil {
		return nil
	}
	out := new(File)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FileSource) DeepCopyInto(out *FileSource) {
	*out = *in
	out.Secret = in.Secret
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FileSource.
func (in *FileSource) DeepCopy() *FileSource {
	if in == nil {
		return nil
	}
	out := new(FileSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Filesystem) DeepCopyInto(out *Filesystem) {
	*out = *in
	if in.Partition != nil {
		in, out := &in.Partition, &out.Partition
		*out = new(string)
		**out = **in
	}
	if in.Overwrite != nil {
		in, out := &in.Overwrite, &out.Overwrite
		*out = new(bool)
		**out = **in
	}
	if in.ExtraOpts != nil {
		in, out := &in.ExtraOpts, &out.ExtraOpts
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Filesystem.
func (in *Filesystem) DeepCopy() *Filesystem {
	if in == nil {
		return nil
	}
	out := new(Filesystem)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in MountPoints) DeepCopyInto(out *MountPoints) {
	{
		in := &in
		*out = make(MountPoints, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MountPoints.
func (in MountPoints) DeepCopy() MountPoints {
	if in == nil {
		return nil
	}
	out := new(MountPoints)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NTP) DeepCopyInto(out *NTP) {
	*out = *in
	if in.Servers != nil {
		in, out := &in.Servers, &out.Servers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Enabled != nil {
		in, out := &in.Enabled, &out.Enabled
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NTP.
func (in *NTP) DeepCopy() *NTP {
	if in == nil {
		return nil
	}
	out := new(NTP)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Partition) DeepCopyInto(out *Partition) {
	*out = *in
	if in.Overwrite != nil {
		in, out := &in.Overwrite, &out.Overwrite
		*out = new(bool)
		**out = **in
	}
	if in.TableType != nil {
		in, out := &in.TableType, &out.TableType
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Partition.
func (in *Partition) DeepCopy() *Partition {
	if in == nil {
		return nil
	}
	out := new(Partition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PasswdSource) DeepCopyInto(out *PasswdSource) {
	*out = *in
	out.Secret = in.Secret
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PasswdSource.
func (in *PasswdSource) DeepCopy() *PasswdSource {
	if in == nil {
		return nil
	}
	out := new(PasswdSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PauseContainer) DeepCopyInto(out *PauseContainer) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PauseContainer.
func (in *PauseContainer) DeepCopy() *PauseContainer {
	if in == nil {
		return nil
	}
	out := new(PauseContainer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretFileSource) DeepCopyInto(out *SecretFileSource) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretFileSource.
func (in *SecretFileSource) DeepCopy() *SecretFileSource {
	if in == nil {
		return nil
	}
	out := new(SecretFileSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretPasswdSource) DeepCopyInto(out *SecretPasswdSource) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretPasswdSource.
func (in *SecretPasswdSource) DeepCopy() *SecretPasswdSource {
	if in == nil {
		return nil
	}
	out := new(SecretPasswdSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *User) DeepCopyInto(out *User) {
	*out = *in
	if in.Gecos != nil {
		in, out := &in.Gecos, &out.Gecos
		*out = new(string)
		**out = **in
	}
	if in.Groups != nil {
		in, out := &in.Groups, &out.Groups
		*out = new(string)
		**out = **in
	}
	if in.HomeDir != nil {
		in, out := &in.HomeDir, &out.HomeDir
		*out = new(string)
		**out = **in
	}
	if in.Inactive != nil {
		in, out := &in.Inactive, &out.Inactive
		*out = new(bool)
		**out = **in
	}
	if in.Shell != nil {
		in, out := &in.Shell, &out.Shell
		*out = new(string)
		**out = **in
	}
	if in.Passwd != nil {
		in, out := &in.Passwd, &out.Passwd
		*out = new(string)
		**out = **in
	}
	if in.PasswdFrom != nil {
		in, out := &in.PasswdFrom, &out.PasswdFrom
		*out = new(PasswdSource)
		**out = **in
	}
	if in.PrimaryGroup != nil {
		in, out := &in.PrimaryGroup, &out.PrimaryGroup
		*out = new(string)
		**out = **in
	}
	if in.LockPassword != nil {
		in, out := &in.LockPassword, &out.LockPassword
		*out = new(bool)
		**out = **in
	}
	if in.Sudo != nil {
		in, out := &in.Sudo, &out.Sudo
		*out = new(string)
		**out = **in
	}
	if in.SSHAuthorizedKeys != nil {
		in, out := &in.SSHAuthorizedKeys, &out.SSHAuthorizedKeys
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new User.
func (in *User) DeepCopy() *User {
	if in == nil {
		return nil
	}
	out := new(User)
	in.DeepCopyInto(out)
	return out
}
