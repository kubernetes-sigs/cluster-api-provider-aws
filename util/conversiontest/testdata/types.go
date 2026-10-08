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

package testdata

// Spoke types (v1beta1-like)

type LayoutIdentical struct {
	Field int
}

type LayoutDifferent struct {
	Field int
	Extra string
}

type Spoke struct {
	Kind1 InnerType
}

type InnerType struct {
	Field string
}

type SpokeWithBypass struct {
	Inner InnerType
}

type CopyOnlyType struct {
	Field string
}

// Hub types (v1beta2-like)

type Hub struct {
	Field int
}

type HubInner struct {
	Field string
}

type HubSpoke struct {
	Kind1 HubInner
}

type HubWithBypass struct {
	Inner HubInner
}

type HubCopyOnly struct {
	Field string
}
