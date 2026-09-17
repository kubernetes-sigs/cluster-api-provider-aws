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
	"strings"
	"testing"
)

func TestPreviousVersionCmd(t *testing.T) {
	cmd := previousVersionCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetIn(strings.NewReader("v2.11.0\nv2.12.0\nv2.13.0\nv2.13.2\n"))
	cmd.SetArgs([]string{"--version", "v2.14.0"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	want := "v2.13.2\n"
	if out.String() != want {
		t.Errorf("output = %q, want %q", out.String(), want)
	}
}

func TestPreviousVersionCmdNoCandidate(t *testing.T) {
	cmd := previousVersionCmd()
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	cmd.SetIn(strings.NewReader("v2.9.0\n"))
	cmd.SetArgs([]string{"--version", "v1.0.0"})

	if err := cmd.Execute(); err == nil {
		t.Fatal("Execute() expected error when no previous version exists, got nil")
	}
}

func TestPreviousVersionCmdRequiresVersion(t *testing.T) {
	cmd := previousVersionCmd()
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	cmd.SetIn(strings.NewReader(""))
	cmd.SetArgs([]string{})

	if err := cmd.Execute(); err == nil {
		t.Fatal("Execute() expected error for missing --version, got nil")
	}
}
