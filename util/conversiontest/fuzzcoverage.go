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

package conversiontest

import (
	"reflect"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
)

// RequireFuzzCoverage fails tb if any Convertible kind in gv lacks a fuzz test in cases.
// knownGaps maps kind names to the reason they are excluded; entries that become covered
// or whose kind no longer exists also fail, so the allowlist cannot go stale.
func RequireFuzzCoverage(tb testing.TB, scheme *runtime.Scheme, gv schema.GroupVersion,
	cases []utilconversion.FuzzTestFuncInput, knownGaps map[string]string) {
	tb.Helper()
	// Collect all Convertible kinds in the group version
	knownTypes := scheme.KnownTypes(gv)
	convertibleKinds := make(map[string]reflect.Type)

	for kind, typ := range knownTypes {
		// Skip List types by convention; individual kinds are checked
		if strings.HasSuffix(kind, "List") {
			continue
		}

		// Check if implements Convertible
		v := reflect.New(typ).Interface()
		if _, ok := v.(conversion.Convertible); ok {
			convertibleKinds[kind] = typ
		}
	}

	// Collect covered kinds from fuzz cases
	covered := make(map[string]bool)
	for _, c := range cases {
		if c.Spoke != nil {
			spokeType := reflect.TypeOf(c.Spoke).Elem()
			covered[spokeType.Name()] = true
		}
	}

	// Check for missing coverage
	for kind := range convertibleKinds {
		if _, inGaps := knownGaps[kind]; inGaps {
			// Kind exists and is in known gaps, but is it covered?
			if covered[kind] {
				tb.Errorf("kind %s is covered but still in knownGaps; remove the entry", kind)
			}
		} else {
			// Kind is not in known gaps, so it must be covered
			if !covered[kind] {
				tb.Errorf("kind %s is not covered by any FuzzTestFuncInput; add fuzz test or document in knownGaps", kind)
			}
		}
	}

	// Check for stale knownGaps entries
	for kind := range knownGaps {
		if _, exists := convertibleKinds[kind]; !exists {
			tb.Errorf("knownGaps entry %s no longer exists in scheme; remove it", kind)
		}
	}
}
