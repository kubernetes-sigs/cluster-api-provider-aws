package v1beta2

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"sigs.k8s.io/cluster-api-provider-aws/v2/feature"
)

// Validate will validate the spec fields.
func (s *AWSClusterSpec) Validate() []*field.Error {
	var errs field.ErrorList

	// Check the feature gate is enabled for OIDC Provider.
	if s.AssociateOIDCProvider && !feature.Gates.Enabled(feature.OIDCProviderSupport) {
		errs = append(errs,
			field.Forbidden(field.NewPath("spec", "associateOIDCProvider"),
				"can be enabled only if the OIDCProviderSupport feature gate is enabled"),
		)
		return errs
	}

	return errs
}
