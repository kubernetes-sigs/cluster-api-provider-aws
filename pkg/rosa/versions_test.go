/*
Copyright The Kubernetes Authors.

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

package rosa

import (
	"testing"

	"github.com/blang/semver"
	. "github.com/onsi/gomega"
)

// TestMachinePoolVersionRangeMembership mirrors the comparison validateMachinePoolSpec
// performs, so that widening the range to admit prerelease control planes cannot
// silently start admitting genuinely unsupported machine pool versions.
func TestMachinePoolVersionRangeMembership(t *testing.T) {
	const controlPlaneVersion = "5.0.0-rc.0"

	tests := []struct {
		machinePoolVersion string
		supported          bool
	}{
		{"5.0.0-rc.0", true}, // the real-world case: pool matches the control plane
		{"5.0.0", true},      // GA pool under an RC control plane of the same version
		{"5.0.3", false},     // patch ahead of the control plane
		{"4.22.9", false},    // previous major, outside the same-major floor
		{"6.0.0", false},     // ahead of the control plane
	}

	for _, tc := range tests {
		t.Run(tc.machinePoolVersion, func(t *testing.T) {
			g := NewWithT(t)

			minVersion, maxVersion, err := MachinePoolSupportedVersionsRange(controlPlaneVersion)
			g.Expect(err).ToNot(HaveOccurred())

			parsed, err := semver.Parse(tc.machinePoolVersion)
			g.Expect(err).ToNot(HaveOccurred())

			core := CoreVersion(parsed)
			inRange := !core.GT(*maxVersion) && !core.LT(*minVersion)
			g.Expect(inRange).To(Equal(tc.supported),
				"machine pool %s against control plane %s (range >= %s, <= %s)",
				tc.machinePoolVersion, controlPlaneVersion, minVersion, maxVersion)
		})
	}
}

func TestMachinePoolSupportedVersionsRange(t *testing.T) {
	tests := []struct {
		name                string
		controlPlaneVersion string
		expectedMin         string
		expectedMax         string
		expectErr           bool
	}{
		{
			name:                "subtracts two minor versions",
			controlPlaneVersion: "4.22.9",
			expectedMin:         "4.20.0",
			expectedMax:         "4.22.9",
		},
		{
			name:                "clamps to the floor when minor is 0",
			controlPlaneVersion: "5.0.0",
			expectedMin:         "5.0.0",
			expectedMax:         "5.0.0",
		},
		{
			name:                "clamps to the floor when minor is 1",
			controlPlaneVersion: "5.1.0",
			expectedMin:         "5.0.0",
			expectedMax:         "5.1.0",
		},
		{
			name:                "clamps to the floor for a prerelease control plane",
			controlPlaneVersion: "5.0.0-rc.0",
			expectedMin:         "5.0.0",
			expectedMax:         "5.0.0",
		},
		{
			name:                "drops prerelease qualifiers from the bounds",
			controlPlaneVersion: "4.22.9-ec.3",
			expectedMin:         "4.20.0",
			expectedMax:         "4.22.9",
		},
		{
			name:                "does not go below the minimum supported version",
			controlPlaneVersion: "4.15.0",
			expectedMin:         MinSupportedVersion.String(),
			expectedMax:         "4.15.0",
		},
		{
			name:                "rejects an unparsable version",
			controlPlaneVersion: "not-a-version",
			expectErr:           true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			minVersion, maxVersion, err := MachinePoolSupportedVersionsRange(tc.controlPlaneVersion)
			if tc.expectErr {
				g.Expect(err).To(HaveOccurred())
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(minVersion.String()).To(Equal(tc.expectedMin))
			g.Expect(maxVersion.String()).To(Equal(tc.expectedMax))

			// The range must be usable: a control plane is always allowed to run
			// machine pools at its own version. This is what broke on 5.x -- the
			// underflowed minimum put every pool out of range, and the prerelease
			// qualifier on an RC control plane did the same thing again.
			cpVersion, parseErr := semver.Parse(tc.controlPlaneVersion)
			g.Expect(parseErr).ToNot(HaveOccurred())
			core := CoreVersion(cpVersion)
			g.Expect(core.LT(*minVersion)).To(BeFalse(),
				"control plane version %s must not fall below the minimum %s", cpVersion, minVersion)
			g.Expect(core.GT(*maxVersion)).To(BeFalse(),
				"control plane version %s must not exceed the maximum %s", cpVersion, maxVersion)
		})
	}
}
