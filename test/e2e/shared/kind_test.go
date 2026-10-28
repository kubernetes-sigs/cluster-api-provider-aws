//go:build e2e
// +build e2e

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
	"strings"
	"testing"

	. "github.com/onsi/gomega"

	"sigs.k8s.io/cluster-api/test/framework/bootstrap"
)

func TestKindNodeImageForVersion(t *testing.T) {
	// The framework's default node image is always pinned by digest in its kind mapping.
	defaultVersion, _, _ := strings.Cut(bootstrap.DefaultNodeImageVersion, "@")

	tests := []struct {
		name           string
		version        string
		wantImage      string
		wantNeedsBuild bool
		wantErr        bool
	}{
		{
			name:           "version published by kind uses the pinned image",
			version:        defaultVersion,
			wantImage:      bootstrap.DefaultNodeImageRepository + ":" + bootstrap.DefaultNodeImageVersion,
			wantNeedsBuild: false,
		},
		{
			name:           "version not published by kind is built",
			version:        "v1.99.0",
			wantImage:      "kindest/node:v1.99.0",
			wantNeedsBuild: true,
		},
		{
			name:    "invalid version",
			version: "not-a-version",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			image, needsBuild, err := kindNodeImageForVersion(tt.version)
			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(image).To(Equal(tt.wantImage))
			g.Expect(needsBuild).To(Equal(tt.wantNeedsBuild))
		})
	}
}
