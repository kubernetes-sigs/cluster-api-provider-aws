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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// SetupWebhookWithManager will setup the webhooks for the AWSManagedCluster.
func (r *AWSManagedCluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1alpha4-awsmanagedcluster,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsmanagedclusters,versions=v1alpha4,name=default.awsmanagedcluster.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1alpha4-awsmanagedcluster,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsmanagedclusters,versions=v1alpha4,name=validation.awsmanagedcluster.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1

var _ webhook.Defaulter = &AWSManagedCluster{}
var _ webhook.Validator = &AWSManagedCluster{}

// Default will set default values for the AWSManagedCluster.
func (r *AWSManagedCluster) Default() {
}

// ValidateCreate will do any extra validation when creating a AWSManagedCluster.
func (r *AWSManagedCluster) ValidateCreate() error {
	return nil
}

// ValidateUpdate will do any extra validation when updating a AWSManagedCluster.
func (r *AWSManagedCluster) ValidateUpdate(old runtime.Object) error {
	return nil
}

// ValidateDelete allows you to add any extra validation when deleting.
func (r *AWSManagedCluster) ValidateDelete() error {
	return nil
}
