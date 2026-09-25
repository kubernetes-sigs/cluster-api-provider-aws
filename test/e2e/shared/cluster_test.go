//go:build e2e

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

package shared

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateManagementClusterSettings(t *testing.T) {
	tests := []struct {
		name     string
		settings Settings
		wantErr  bool
	}{
		{name: "kind management cluster"},
		{name: "existing management cluster", settings: Settings{UseExistingCluster: true}},
		{name: "self-hosted management cluster", settings: Settings{ProvisionSelfHostedManagementCluster: true}},
		{
			name: "existing and self-hosted management clusters",
			settings: Settings{
				UseExistingCluster:                   true,
				ProvisionSelfHostedManagementCluster: true,
			},
			wantErr: true,
		},
		{
			name: "provision and teardown",
			settings: Settings{
				ProvisionSelfHostedManagementCluster: true,
				TeardownSelfHostedManagementCluster:  true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateManagementClusterSettings(tt.settings)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateManagementClusterSettings() error = %v, wantErr %t", err, tt.wantErr)
			}
		})
	}
}

func TestPersistSelfHostedManagementClusterKubeconfig(t *testing.T) {
	artifactFolder := t.TempDir()
	sourcePath := filepath.Join(t.TempDir(), "source.kubeconfig")
	if err := os.WriteFile(sourcePath, []byte("apiVersion: v1\n"), 0o600); err != nil {
		t.Fatalf("failed to create source kubeconfig: %v", err)
	}

	gotPath, err := persistSelfHostedManagementClusterKubeconfig(sourcePath, artifactFolder)
	if err != nil {
		t.Fatalf("persistSelfHostedManagementClusterKubeconfig() error = %v", err)
	}
	wantPath := filepath.Join(artifactFolder, selfHostedManagementClusterKubeconfigFile)
	if gotPath != wantPath {
		t.Fatalf("persistSelfHostedManagementClusterKubeconfig() path = %q, want %q", gotPath, wantPath)
	}
	contents, err := os.ReadFile(gotPath) //nolint:gosec // The test controls this path.
	if err != nil {
		t.Fatalf("failed to read persisted kubeconfig: %v", err)
	}
	if string(contents) != "apiVersion: v1\n" {
		t.Fatalf("persisted kubeconfig contents = %q", contents)
	}
}
