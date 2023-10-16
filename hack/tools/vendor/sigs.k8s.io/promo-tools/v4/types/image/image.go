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

package image

// UserAgent header to be used in the requests
const UserAgent = "kpromo"

// Name can be just the bare name itself (e.g., "addon-builder" in
// "gcr.io/k8s-image-staging/addon-builder") or the prefix + name
// ("foo/bar/baz/quux" in "gcr.io/hello/foo/bar/baz/quux").
type Name string

// Registry is the leading part of an image name that includes the domain;
// it is everything that is not the actual image name itself. e.g.,
// "gcr.io/google-containers".
type Registry string

// Digest is a string that contains the SHA256 hash of a Docker container image.
type Digest string

// Tag is a Docker tag.
type Tag string
