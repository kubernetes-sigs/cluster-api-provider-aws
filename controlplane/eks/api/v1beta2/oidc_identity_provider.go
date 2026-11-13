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

// EffectiveOIDCIdentityProviderConfigs returns the OIDC identity provider configs declared on the control plane.
// When oidcIdentityProviderConfigs is set it takes precedence over the singular oidcIdentityProviderConfig field.
func (s *AWSManagedControlPlaneSpec) EffectiveOIDCIdentityProviderConfigs() []OIDCIdentityProviderConfig {
	if len(s.OIDCIdentityProviderConfigs) > 0 {
		return s.OIDCIdentityProviderConfigs
	}
	if s.OIDCIdentityProviderConfig != nil {
		return []OIDCIdentityProviderConfig{*s.OIDCIdentityProviderConfig}
	}
	return nil
}
