/*
Copyright 2026 The Kubernetes Authors.

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

package eks

import (
	"testing"
)

func TestIsSymbolicLaunchTemplateVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    bool
	}{
		{name: "$Latest is symbolic", version: "$Latest", want: true},
		{name: "$Default is symbolic", version: "$Default", want: true},
		{name: "concrete version 1", version: "1", want: false},
		{name: "concrete version 42", version: "42", want: false},
		{name: "empty string", version: "", want: false},
		{name: "lowercase $latest is not symbolic", version: "$latest", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSymbolicLaunchTemplateVersion(tt.version); got != tt.want {
				t.Errorf("isSymbolicLaunchTemplateVersion(%q) = %v, want %v", tt.version, got, tt.want)
			}
		})
	}
}
