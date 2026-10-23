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
	tests := []struct {
		controlPlaneVersion string
		machinePoolVersion  string
		supported           bool
	}{
		// RC control plane cases
		{"5.0.0-rc.0", "5.0.0-rc.0", true},  // pool matches RC control plane exactly
		{"5.0.0-rc.0", "5.0.0", false},      // GA pool must not exceed an RC control plane
		{"5.0.0-rc.0", "5.0.0-rc.1", false}, // later RC must not exceed the control plane RC
		{"5.0.0-rc.0", "5.0.3", false},      // patch ahead of the control plane
		{"5.0.0-rc.0", "4.22.9", true},      // previous major, admitted by cross-major window
		{"5.0.0-rc.0", "4.14.0", true},      // at the MinSupportedVersion floor — still admitted
		{"5.0.0-rc.0", "6.0.0", false},      // ahead of the control plane major
		// GA control plane cases
		{"5.0.0", "5.0.0-rc.0", true}, // RC pool is behind GA control plane — accepted
		{"5.0.0", "5.0.0", true},      // pool matches GA control plane exactly
		{"5.0.0", "5.0.1", false},     // patch ahead of GA control plane
		// 5.1 CP upgrade-path cases: 4.x pools admitted via cross-major window
		{"5.1.0", "5.1.0", true},   // pool matches 5.1 control plane
		{"5.1.0", "5.0.0", true},   // one minor behind, within skew
		{"5.1.0", "4.22.9", true},  // upgrade source: 4.x pool on 5.1 CP
		{"5.1.0", "4.23.0", true},  // upgrade target: 4.x pool on 5.1 CP
		{"5.1.0", "4.13.9", false}, // below MinSupportedVersion floor
		{"5.1.0", "5.1.1", false},  // patch ahead of control plane
		// 4.x RC pool against 4.x GA CP (different minor)
		{"4.22.0", "4.21.0-rc.0", true}, // RC pool within skew window of GA control plane
		{"4.22.0", "4.19.9", false},     // pool outside skew window
	}

	for _, tc := range tests {
		t.Run(tc.controlPlaneVersion+"/"+tc.machinePoolVersion, func(t *testing.T) {
			g := NewWithT(t)

			parsed, err := semver.Parse(tc.machinePoolVersion)
			g.Expect(err).ToNot(HaveOccurred())

			err = ValidateMachinePoolVersion(tc.controlPlaneVersion, parsed)
			if tc.supported {
				g.Expect(err).To(BeNil(), "expected version %s to be supported for control plane %s", tc.machinePoolVersion, tc.controlPlaneVersion)
			} else {
				g.Expect(err).To(HaveOccurred(), "expected version %s to be rejected for control plane %s", tc.machinePoolVersion, tc.controlPlaneVersion)
				g.Expect(err.Error()).To(ContainSubstring("is not supported"))
			}
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
			name:                "cross-major window for minor 0",
			controlPlaneVersion: "5.0.0",
			expectedMin:         "4.14.0",
			expectedMax:         "5.0.0",
		},
		{
			name:                "cross-major window for minor 1",
			controlPlaneVersion: "5.1.0",
			expectedMin:         "4.14.0",
			expectedMax:         "5.1.0",
		},
		{
			name:                "cross-major window for prerelease control plane",
			controlPlaneVersion: "5.0.0-rc.0",
			expectedMin:         "4.14.0",
			expectedMax:         "5.0.0-rc.0",
		},
		{
			name:                "preserves prerelease qualifier in the upper bound",
			controlPlaneVersion: "4.22.9-ec.3",
			expectedMin:         "4.20.0",
			expectedMax:         "4.22.9-ec.3",
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

			// A pool at the same version as the control plane must always be
			// accepted. This is the invariant that broke on 5.x: the underflowed
			// minimum put every same-version pool out of range.
			cpVersion, parseErr := semver.Parse(tc.controlPlaneVersion)
			g.Expect(parseErr).ToNot(HaveOccurred())
			g.Expect(ValidateMachinePoolVersion(tc.controlPlaneVersion, cpVersion)).To(Succeed(),
				"a pool at control plane version %s must be accepted by its own range [%s, %s]", tc.controlPlaneVersion, minVersion, maxVersion)
		})
	}
}
