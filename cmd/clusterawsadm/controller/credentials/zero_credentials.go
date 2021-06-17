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

package credentials

// ZeroCredentialsInput defines the specs for zero credentials input.
type ZeroCredentialsInput struct {
	KubeconfigPath    string
	KubeconfigContext string
	Namespace         string
}

// ZeroCredentials zeroes out the CAPA controller bootstrap secret.
// RolloutControllers() must be called after any change to the controller bootstrap secret to take effect.
func ZeroCredentials(input ZeroCredentialsInput) error {
	return UpdateCredentials(UpdateCredentialsInput{
		KubeconfigPath:    input.KubeconfigPath,
		KubeconfigContext: input.KubeconfigContext,
		Credentials:       "Cg==",
		Namespace:         input.Namespace,
	})
}
