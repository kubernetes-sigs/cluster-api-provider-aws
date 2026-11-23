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

package v1beta2

import (
	"errors"
)

// ErrConflictingOIDCIdentityProviderFields is returned when both OIDC identity provider spec fields are set.
var ErrConflictingOIDCIdentityProviderFields = errors.New("oidcIdentityProviderConfig and oidcIdentityProviderConfigs cannot both be set")

// HasSingularOIDCIdentityProviderConfig reports whether the deprecated singular OIDC identity provider field is set.
func (s *AWSManagedControlPlaneSpec) HasSingularOIDCIdentityProviderConfig() bool {
	return s.OIDCIdentityProviderConfig != nil
}

// SingularOIDCIdentityProviderConfig returns the deprecated singular OIDC identity provider config, if set.
func (s *AWSManagedControlPlaneSpec) SingularOIDCIdentityProviderConfig() *OIDCIdentityProviderConfig {
	return s.OIDCIdentityProviderConfig
}

// SetSingularOIDCIdentityProviderConfig sets the deprecated singular OIDC identity provider config.
func (s *AWSManagedControlPlaneSpec) SetSingularOIDCIdentityProviderConfig(cfg *OIDCIdentityProviderConfig) {
	s.OIDCIdentityProviderConfig = cfg
}

// WithSingularOIDCIdentityProviderConfig returns a copy of the spec with the deprecated singular OIDC field set.
func (s AWSManagedControlPlaneSpec) WithSingularOIDCIdentityProviderConfig(cfg *OIDCIdentityProviderConfig) AWSManagedControlPlaneSpec {
	s.OIDCIdentityProviderConfig = cfg
	return s
}

// EffectiveOIDCIdentityProviderConfigs returns the OIDC identity provider configs declared on the control plane.
func (s *AWSManagedControlPlaneSpec) EffectiveOIDCIdentityProviderConfigs() ([]OIDCIdentityProviderConfig, error) {
	if s.HasSingularOIDCIdentityProviderConfig() && len(s.OIDCIdentityProviderConfigs) > 0 {
		return nil, ErrConflictingOIDCIdentityProviderFields
	}
	if len(s.OIDCIdentityProviderConfigs) > 0 {
		return s.OIDCIdentityProviderConfigs, nil
	}
	if cfg := s.SingularOIDCIdentityProviderConfig(); cfg != nil {
		return []OIDCIdentityProviderConfig{*cfg}, nil
	}
	return nil, nil
}
