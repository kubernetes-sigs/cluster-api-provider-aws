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

package template

import (
	"bytes"
	"testing"
	texttemplate "text/template"
)

func TestFuncMap(t *testing.T) {
	tmpl, err := texttemplate.New("test").Funcs(FuncMap()).Parse(`{{ Indent 2 "a\nb" }}|{{ join . "," }}|{{ base64Encode "test" }}|{{ default "fallback" "" }}|{{ trimSpace " value " }}|{{ contains "value" "al" }}|{{ hasPrefix "value" "va" }}|{{ hasSuffix "value" "ue" }}|{{ replace "a-b" "-" "_" }}|{{ lower "VALUE" }}|{{ upper "value" }}`)
	if err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	if err := tmpl.Execute(&out, []string{"a", "b"}); err != nil {
		t.Fatal(err)
	}

	expected := "  a\n  b|a,b|dGVzdA==|fallback|value|true|true|true|a_b|value|VALUE"
	if out.String() != expected {
		t.Fatalf("unexpected output: got %q, want %q", out.String(), expected)
	}
}
