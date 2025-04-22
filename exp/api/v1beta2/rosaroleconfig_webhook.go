package v1beta2

import (
	"context"

	runtime "k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// SetupWebhookWithManager will setup the webhooks for the ROSARoleConfig.
func (r *ROSARoleConfig) SetupWebhookWithManager(mgr ctrl.Manager) error {
	w := new(rosaRoleConfigWebhook)
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		WithValidator(w).
		WithDefaulter(w).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-rosaroleconfig,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=rosaroleconfigs,versions=v1beta2,name=validation.rosaroleconfig.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-rosaroleconfig,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=rosaroleconfigs,versions=v1beta2,name=default.rosaroleconfig.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

type rosaRoleConfigWebhook struct{}

var _ webhook.CustomDefaulter = &rosaRoleConfigWebhook{}
var _ webhook.CustomValidator = &rosaRoleConfigWebhook{}

// ValidateCreate implements admission.Validator.
func (r *rosaRoleConfigWebhook) ValidateCreate(ctx context.Context, obj runtime.Object) (warnings admission.Warnings, err error) {
	return nil, nil
}

// ValidateUpdate implements admission.Validator.
func (r *rosaRoleConfigWebhook) ValidateUpdate(ctx context.Context, old runtime.Object, updated runtime.Object) (warnings admission.Warnings, err error) {
	return nil, nil
}

// ValidateDelete implements admission.Validator.
func (r *rosaRoleConfigWebhook) ValidateDelete(ctx context.Context, obj runtime.Object) (warnings admission.Warnings, err error) {
	return nil, nil
}

// Default implements admission.Defaulter.
func (r *rosaRoleConfigWebhook) Default(ctx context.Context, obj runtime.Object) error {
	return nil
}
