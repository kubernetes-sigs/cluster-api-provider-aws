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

import "testing"

func TestCmdRegistersAllSubcommands(t *testing.T) {
	want := []string{"validate-trigger", "notes-pr", "body", "previous-version", "verify-metadata"}

	got := map[string]bool{}
	for _, c := range Cmd().Commands() {
		got[c.Name()] = true
	}

	for _, name := range want {
		if !got[name] {
			t.Errorf("Cmd() did not register subcommand %q", name)
		}
	}
	if len(got) != len(want) {
		t.Errorf("Cmd() registered %d subcommands, want %d (got: %v)", len(got), len(want), got)
	}
}
