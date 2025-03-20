/*
Copyright 2025 The Kubernetes Authors.

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

package controllers

import (
	"slices"
	"testing"

	. "github.com/onsi/gomega"
)

func TestIsDisabled(t *testing.T) {
	g := NewWithT(t)
	defer resetDefaults()

	// Valid names default to false
	g.Expect(IsDisabled(Unmanaged)).To(BeFalse(), "unmanaged should default to false")

	// Invalid names are also false
	g.Expect(IsDisabled("eks")).To(BeFalse(), "invalid controller name eks should report disabled")
	g.Expect(IsDisabled("rosa")).To(BeFalse(), "invalid controller name rosa should report disabled")

	// Disable the known names
	disabledControllers = map[string]bool{
		Unmanaged: true,
	}

	// Valid names
	g.Expect(IsDisabled(Unmanaged)).To(BeTrue(), "unmanaged should have been disabled")

	// Invalid names are still false
	g.Expect(IsDisabled("eks")).To(BeFalse(), "invalid controller name eks should report disabled")
	g.Expect(IsDisabled("rosa")).To(BeFalse(), "invalid controller name rosa should report disabled")
}

func TestGetValidNames(t *testing.T) {
	g := NewWithT(t)
	defer resetDefaults()

	actual := GetValidNames()
	// Make sure we have a stable order for testing
	slices.Sort(actual)

	g.Expect(actual).To(Equal([]string{
		Unmanaged,
	}), "should only have 1 name")
}

func TestValidateNames(t *testing.T) {
	g := NewWithT(t)
	defer resetDefaults()

	// Valid set of names. Will mutate the map.
	err := ValidateNamesAndDisable([]string{Unmanaged})
	g.Expect(err).To(BeNil())
	g.Expect(disabledControllers[Unmanaged]).To(BeTrue(), "should disable valid name unmanaged")

	// TODO: This test should fail and require updating when EKS and ROSA controllers graduate.
	err = ValidateNamesAndDisable([]string{"eks", "rosa"})
	g.Expect(err.Error()).To(ContainSubstring("eks"), "should error on first invalid entry")
	g.Expect(disabledControllers[Unmanaged]).To(BeTrue(), "should not change existing key unmanaged")
	g.Expect(disabledControllers["eks"]).To(BeFalse(), "eks should not be in the default map")
}

// resetDefaults returns the disabledControllers map to expected default state.
func resetDefaults() {
	disabledControllers = map[string]bool{
		Unmanaged: false,
	}
}
