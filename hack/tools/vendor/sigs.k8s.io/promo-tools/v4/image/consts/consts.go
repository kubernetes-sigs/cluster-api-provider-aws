/*
Copyright 2023 The Kubernetes Authors.

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

package consts

const (
	// Production registry root URL
	ProdRegistry = "registry.k8s.io"

	// Staging repository root URL prefix
	StagingRepoPrefix = "gcr.io/k8s-staging-"

	// The suffix of the default image repository to promote images from
	// i.e., gcr.io/<staging-prefix>-<staging-suffix>
	// e.g., gcr.io/k8s-staging-foo
	StagingRepoSuffix = "kubernetes"
)
