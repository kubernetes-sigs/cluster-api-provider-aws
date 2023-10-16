/*
Copyright 2020 The Kubernetes Authors.

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

package remotemanifest

import "sigs.k8s.io/promo-tools/v4/internal/legacy/dockerregistry/schema"

// Fake is a fake remote manifest. It is fake in the sense that it
// will never fetch anything from any remote.
type Fake struct {
	manifests []schema.Manifest
}

// Fetch just returns the manifests that were set in NewFakeRemoteManifest.
func (remote *Fake) Fetch() ([]schema.Manifest, error) {
	return remote.manifests, nil
}

// NewFake creates a new Fake.
func NewFake(manifests []schema.Manifest) *Fake {
	remote := Fake{}

	remote.manifests = manifests

	return &remote
}
