package v1beta2

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// SetupWebhookWithManager will setup the webhooks for the ROSARoleConfig.
func (r *ROSARoleConfig) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-rosaroleconfig,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=rosaroleconfigs,versions=v1beta2,name=validation.rosaroleconfig.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-rosaroleconfig,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=rosaroleconfigs,versions=v1beta2,name=default.rosaroleconfig.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

var _ webhook.Defaulter = &ROSARoleConfig{}
var _ webhook.Validator = &ROSARoleConfig{}

// ValidateCreate implements admission.Validator.
func (r *ROSARoleConfig) ValidateCreate() (warnings admission.Warnings, err error) {
	return nil, nil
}

// ValidateUpdate implements admission.Validator.
func (r *ROSARoleConfig) ValidateUpdate(old runtime.Object) (warnings admission.Warnings, err error) {
	return nil, nil
}

// ValidateDelete implements admission.Validator.
func (r *ROSARoleConfig) ValidateDelete() (warnings admission.Warnings, err error) {
	return nil, nil
}

// Default implements admission.Defaulter.
func (r *ROSARoleConfig) Default() {

}
