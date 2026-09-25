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

package changelog

import (
	"bytes"
	"testing"
)

func TestWithFrontMatterAndParseFrontMatterRoundTrip(t *testing.T) {
	body := []byte("# Changelog\n\n- did a thing\n")
	withFM := WithFrontMatter("v1beta2", body)

	gotContract, gotBody := ParseFrontMatter(withFM)
	if gotContract != "v1beta2" {
		t.Errorf("contract = %q, want %q", gotContract, "v1beta2")
	}
	if !bytes.Equal(gotBody, body) {
		t.Errorf("body = %q, want %q", gotBody, body)
	}
}

func TestParseFrontMatterNoFrontMatter(t *testing.T) {
	body := []byte("# Changelog\n\n- did a thing\n")

	gotContract, gotBody := ParseFrontMatter(body)
	if gotContract != "" {
		t.Errorf("contract = %q, want empty", gotContract)
	}
	if !bytes.Equal(gotBody, body) {
		t.Errorf("body = %q, want %q", gotBody, body)
	}
}
