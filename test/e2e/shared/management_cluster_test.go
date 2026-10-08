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

func TestValidateManagementClusterFlags(t *testing.T) {
	tests := []struct {
		name        string
		useExisting bool
		provision   bool
		teardown    bool
		wantErr     bool
	}{
		{name: "kind management cluster"},
		{name: "existing management cluster", useExisting: true},
		{name: "self-hosted management cluster", provision: true},
		{
			name:        "existing and self-hosted management clusters",
			useExisting: true,
			provision:   true,
			wantErr:     true,
		},
		{
			name:      "provision and teardown",
			provision: true,
			teardown:  true,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateManagementClusterFlags(tt.useExisting, tt.provision, tt.teardown)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateManagementClusterFlags() error = %v, wantErr %t", err, tt.wantErr)
			}
		})
	}
}

func TestManagementClusterLifecycleOnly(t *testing.T) {
	tests := []struct {
		name      string
		provision bool
		teardown  bool
		want      bool
	}{
		{name: "zero settings", want: false},
		{name: "provision only", provision: true, want: true},
		{name: "teardown only", teardown: true, want: true},
		{name: "both provision and teardown", provision: true, teardown: true, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := managementClusterLifecycleOnly(tt.provision, tt.teardown)
			if got != tt.want {
				t.Fatalf("managementClusterLifecycleOnly() = %v, want %v", got, tt.want)
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

func TestRemoveFileIfExists(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() string // returns file path to test
		wantErr bool
	}{
		{
			name: "remove existing file",
			setup: func() string {
				tmpDir := t.TempDir()
				filePath := filepath.Join(tmpDir, "testfile")
				if err := os.WriteFile(filePath, []byte("test"), 0o600); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				return filePath
			},
			wantErr: false,
		},
		{
			name: "missing path returns nil",
			setup: func() string {
				return filepath.Join(t.TempDir(), "nonexistent.file")
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			err := removeFileIfExists(path)
			if (err != nil) != tt.wantErr {
				t.Fatalf("removeFileIfExists() error = %v, wantErr %t", err, tt.wantErr)
			}
			// Verify file is actually gone
			if _, err := os.Stat(path); err == nil {
				t.Fatalf("file still exists after removeFileIfExists()")
			}
		})
	}
}
