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

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
	utilpointer "k8s.io/utils/pointer"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
)

const (
	// DefaultBootstrapUserName is the default bootstrap user name.
	DefaultBootstrapUserName = "bootstrapper.cluster-api-provider-aws.sigs.k8s.io"
	// DefaultStackName is the default CloudFormation stack name.
	DefaultStackName = "cluster-api-provider-aws-sigs-k8s-io"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

// SetDefaults_BootstrapUser is used by defaulter-gen
func SetDefaults_BootstrapUser(obj *BootstrapUser) { //nolint:golint,stylecheck
	obj.UserName = DefaultBootstrapUserName
}

// SetDefaults_AWSIAMConfigurationSpec is used by defaulter-gen
func SetDefaults_AWSIAMConfigurationSpec(obj *AWSIAMConfigurationSpec) { //nolint:golint,stylecheck
	if obj.NameSuffix == nil {
		obj.NameSuffix = utilpointer.StringPtr(infrav1.DefaultNameSuffix)
	}
	if obj.StackName == "" {
		obj.StackName = DefaultStackName
	}
	if obj.EKS == nil {
		obj.EKS = &EKSConfig{
			Enable:               false,
			AllowIAMRoleCreation: false,
			DefaultControlPlaneRole: AWSIAMRoleSpec{
				Disable: true,
			},
		}
	} else if obj.EKS.Enable {
		obj.Nodes.EC2ContainerRegistryReadOnly = true
	}
	if obj.EKS.ManagedMachinePool == nil {
		obj.EKS.ManagedMachinePool = &AWSIAMRoleSpec{
			Disable: true,
		}
	}
	if len(obj.SecureSecretsBackends) == 0 {
		obj.SecureSecretsBackends = []infrav1.SecretBackend{
			infrav1.SecretBackendSecretsManager,
		}
	}
}

// SetDefaults_AWSIAMConfiguration is used by defaulter-gen
func SetDefaults_AWSIAMConfiguration(obj *AWSIAMConfiguration) { //nolint:golint,stylecheck
	obj.APIVersion = SchemeGroupVersion.String()
	obj.Kind = "AWSIAMConfiguration"
	if obj.Spec.NameSuffix == nil {
		obj.Spec.NameSuffix = utilpointer.StringPtr(infrav1.DefaultNameSuffix)
	}
	if obj.Spec.StackName == "" {
		obj.Spec.StackName = DefaultStackName
	}
}
