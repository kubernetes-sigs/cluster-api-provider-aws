package rosa

import (
	"fmt"
	"time"

	"github.com/blang/semver"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	"github.com/openshift/rosa/pkg/ocm"
)

// MinSupportedVersion is the minimum supported version for ROSA.
var MinSupportedVersion = semver.MustParse("4.14.0")

// CheckExistingScheduledUpgrade checks and returns the current upgrade schedule if any.
func CheckExistingScheduledUpgrade(client OCMClient, cluster *cmv1.Cluster) (*cmv1.ControlPlaneUpgradePolicy, error) {
	upgradePolicies, err := client.GetControlPlaneUpgradePolicies(cluster.ID())
	if err != nil {
		return nil, err
	}
	for _, upgradePolicy := range upgradePolicies {
		if upgradePolicy.UpgradeType() == cmv1.UpgradeTypeControlPlane {
			return upgradePolicy, nil
		}
	}
	return nil, nil
}

// ScheduleControlPlaneUpgrade schedules a new control plane upgrade to the specified version at the specified time.
func ScheduleControlPlaneUpgrade(client OCMClient, cluster *cmv1.Cluster, version string, nextRun time.Time, ack bool) (*cmv1.ControlPlaneUpgradePolicy, error) {
	// earliestNextRun is set to at least 5 min from now by the OCM API.
	// Set our next run request to something slightly longer than 5min to make sure we account for the latency between when we send this
	// request and when the server processes it.
	earliestNextRun := time.Now().Add(time.Minute * 6)
	if nextRun.Before(earliestNextRun) {
		nextRun = earliestNextRun
	}

	upgradePolicy, err := cmv1.NewControlPlaneUpgradePolicy().
		UpgradeType(cmv1.UpgradeTypeControlPlane).
		ScheduleType(cmv1.ScheduleTypeManual).
		Version(version).
		NextRun(nextRun).
		EnableMinorVersionUpgrades(true).
		Build()
	if err != nil {
		return nil, err
	}

	versionGates, err := client.GetMissingGateAgreementsHypershift(cluster.ID(), upgradePolicy)
	if err != nil {
		return nil, err
	}

	if !ack && len(versionGates) > 0 {
		errMess := "version gate acknowledgment required"
		for id := range versionGates {
			errMess = fmt.Sprintf(errMess+"\nid:%s\n %s\n %s\n %s\n", versionGates[id].ID(), versionGates[id].Description(), versionGates[id].DocumentationURL(), versionGates[id].WarningMessage())
		}

		return nil, fmt.Errorf("%s", errMess)
	}

	for id := range versionGates {
		if err = client.AckVersionGate(cluster.ID(), versionGates[id].ID()); err != nil {
			return nil, err
		}
	}

	return client.ScheduleHypershiftControlPlaneUpgrade(cluster.ID(), upgradePolicy)
}

// ScheduleNodePoolUpgrade schedules a new nodePool upgrade to the specified version at the specified time.
func ScheduleNodePoolUpgrade(client OCMClient, clusterID string, nodePool *cmv1.NodePool, version string, nextRun time.Time) (*cmv1.NodePoolUpgradePolicy, error) {
	// earliestNextRun is set to at least 5 min from now by the OCM API.
	// Set our next run request to something slightly longer than 5min to make sure we account for the latency between when we send this
	// request and when the server processes it.
	earliestNextRun := time.Now().Add(time.Minute * 6)
	if nextRun.Before(earliestNextRun) {
		nextRun = earliestNextRun
	}

	upgradePolicy, err := cmv1.NewNodePoolUpgradePolicy().
		UpgradeType(cmv1.UpgradeTypeNodePool).
		NodePoolID(nodePool.ID()).
		ScheduleType(cmv1.ScheduleTypeManual).
		Version(version).
		NextRun(nextRun).
		EnableMinorVersionUpgrades(true).
		Build()
	if err != nil {
		return nil, err
	}

	scheduledUpgrade, err := client.ScheduleNodePoolUpgrade(clusterID, nodePool.ID(), upgradePolicy)
	if err != nil {
		return nil, fmt.Errorf("failed to schedule nodePool upgrade to version %s: %w", version, err)
	}

	return scheduledUpgrade, nil
}

// machinepools can be created with a minimal of two minor versions from the control plane.
const minorVersionsAllowedDeviation = 2

// CoreVersion strips any prerelease and build qualifiers, leaving major.minor.patch.
// The machine pool skew policy is expressed purely in terms of minor versions, and
// semver sorts a prerelease below its own release ("5.0.0-rc.0" < "5.0.0"), so the
// qualifiers have to be dropped before any range comparison -- otherwise a
// prerelease control plane falls outside the range derived from itself.
func CoreVersion(version semver.Version) semver.Version {
	return semver.Version{Major: version.Major, Minor: version.Minor, Patch: version.Patch}
}

// MachinePoolSupportedVersionsRange returns the supported range of versions
// for a machine pool based on the control plane version. Both bounds are core
// versions; compare candidates with CoreVersion.
func MachinePoolSupportedVersionsRange(controlPlaneVersion string) (*semver.Version, *semver.Version, error) {
	parsed, err := semver.Parse(controlPlaneVersion)
	if err != nil {
		return nil, nil, err
	}
	maxVersion := CoreVersion(parsed)

	// Minor is a uint64, so subtracting the allowed deviation from a minor
	// version below it wraps around instead of clamping -- max() cannot undo
	// that, because the value is already enormous by the time it is applied.
	// Clamp the subtrahend instead so x.0 and x.1 control planes floor at x.0.
	minVersion := semver.Version{
		Major: maxVersion.Major,
		Minor: maxVersion.Minor - min(maxVersion.Minor, minorVersionsAllowedDeviation),
		Patch: 0,
	}

	if minVersion.LT(MinSupportedVersion) {
		minVersion = MinSupportedVersion
	}

	return &minVersion, &maxVersion, nil
}

// RawVersionID returns the rawID from the provided OCM version object.
func RawVersionID(version *cmv1.Version) string {
	rawID := version.RawID()
	if rawID != "" {
		return rawID
	}

	return ocm.GetRawVersionId(version.ID())
}
