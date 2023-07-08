/*
Copyright 2022 The Kubernetes Authors.

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

package registry

import (
	"fmt"

	"sigs.k8s.io/promo-tools/v3/internal/legacy/gcloud"
	"sigs.k8s.io/promo-tools/v3/types/image"
)

// Context holds information about a registry, to be written in a
// manifest file.
type Context struct {
	Name           image.Registry `yaml:"name,omitempty"`
	ServiceAccount string         `yaml:"service-account,omitempty"`
	Token          gcloud.Token   `yaml:"-"`
	Src            bool           `yaml:"src,omitempty"`
}

// GetSrcRegistry gets the source registry.
func GetSrcRegistry(rcs []Context) (*Context, error) {
	for _, registry := range rcs {
		registry := registry
		if registry.Src {
			return &registry, nil
		}
	}

	return nil, fmt.Errorf("could not find source registry")
}
