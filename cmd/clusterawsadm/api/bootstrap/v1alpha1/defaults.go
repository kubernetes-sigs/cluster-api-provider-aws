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
)

const (
	// DefaultNameSuffix is the default suffix appended to all AWS IAM roles created by clusterawsadm.
	DefaultNameSuffix = ".cluster-api-provider-aws.sigs.k8s.io"
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
		obj.NameSuffix = utilpointer.StringPtr(DefaultNameSuffix)
	}
	if obj.StackName == "" {
		obj.StackName = DefaultStackName
	}
}

// SetDefaults_AWSIAMConfiguration is used by defaulter-gen
func SetDefaults_AWSIAMConfiguration(obj *AWSIAMConfiguration) { //nolint:golint,stylecheck
	obj.APIVersion = SchemeGroupVersion.String()
	obj.Kind = "AWSIAMConfiguration"
	if obj.Spec.NameSuffix == nil {
		obj.Spec.NameSuffix = utilpointer.StringPtr(DefaultNameSuffix)
	}
	if obj.Spec.StackName == "" {
		obj.Spec.StackName = DefaultStackName
	}
}
