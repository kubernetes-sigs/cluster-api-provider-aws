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

func TestValidateTriggerCmd(t *testing.T) {
	cmd := validateTriggerCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--branch", "main", "--version", "v2.10.0"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	got := out.String()
	for _, want := range []string{
		"release_branch=release-2.10",
		"create_release_branch=true",
		"major=2",
		"minor=10",
		"patch=0",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("output %q missing %q", got, want)
		}
	}
}

func TestValidateTriggerCmdRejectsInvalid(t *testing.T) {
	cmd := validateTriggerCmd()
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	cmd.SetArgs([]string{"--branch", "main", "--version", "v2.10.1"})

	if err := cmd.Execute(); err == nil {
		t.Fatal("Execute() expected error for main with non-zero patch, got nil")
	}
}

func TestValidateTriggerCmdRequiresFlags(t *testing.T) {
	cmd := validateTriggerCmd()
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	cmd.SetArgs([]string{})

	if err := cmd.Execute(); err == nil {
		t.Fatal("Execute() expected error for missing required flags, got nil")
	}
}
