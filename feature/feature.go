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

package feature

import (
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/component-base/featuregate"
)

const (
	// Every capa-specific feature gate should add method here following this template:
	//
	// // owner: @username
	// // alpha: v1.X
	// MyFeature featuregate.Feature = "MyFeature"

	// owner: @richardcase
	// alpha: v0.4
	EKS featuregate.Feature = "EKS"
)

func init() {
	runtime.Must(MutableGates.Add(defaultCAPAFeatureGates))
}

// defaultCAPAFeatureGates consists of all known capa-specific feature keys.
// To add a new feature, define a key for it above and add it here.
var defaultCAPAFeatureGates = map[featuregate.Feature]featuregate.FeatureSpec{
	// Every feature should be initiated here:
	EKS: {Default: false, PreRelease: featuregate.Alpha},
}
