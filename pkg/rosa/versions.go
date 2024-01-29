package rosa

import (
	"fmt"
	"time"

	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

// IsVersionSupported checks whether the input version is supported for ROSA clusters.
func (c *RosaClient) IsVersionSupported(versionID string) (bool, error) {
	filter := fmt.Sprintf("raw_id='%s' AND channel_group = '%s'", versionID, "stable")
	response, err := c.ocm.ClustersMgmt().V1().
		Versions().
		List().
		Search(filter).
		Page(1).Size(1).
		Parameter("product", "hcp").
		Send()
	if err != nil {
		return false, handleErr(response.Error(), err)
	}
	if response.Total() == 0 {
		return false, nil
	}

	version := response.Items().Get(0)
	return version.ROSAEnabled() && version.HostedControlPlaneEnabled(), nil
}

// CheckExistingScheduledUpgrade checks and returns the current upgrade schedule if any.
func (c *RosaClient) CheckExistingScheduledUpgrade(cluster *cmv1.Cluster) (*cmv1.ControlPlaneUpgradePolicy, error) {
	upgradePolicies, err := c.getControlPlaneUpgradePolicies(cluster.ID())
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
func (c *RosaClient) ScheduleControlPlaneUpgrade(cluster *cmv1.Cluster, version string, nextRun time.Time) (*cmv1.ControlPlaneUpgradePolicy, error) {
	// earliestNextRun is set to at least 5 min from now by the OCM API.
	// we set it to 6 min here to account for latencty.
	earliestNextRun := time.Now().Add(time.Minute * 6)
	if nextRun.Before(earliestNextRun) {
		nextRun = earliestNextRun
	}

	upgradePolicy, err := cmv1.NewControlPlaneUpgradePolicy().
		UpgradeType(cmv1.UpgradeTypeControlPlane).
		ScheduleType(cmv1.ScheduleTypeManual).
		Version(version).
		NextRun(nextRun).
		Build()
	if err != nil {
		return nil, err
	}

	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(cluster.ID()).
		ControlPlane().
		UpgradePolicies().
		Add().Body(upgradePolicy).
		Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}

	return response.Body(), nil
}

func (c *RosaClient) getControlPlaneUpgradePolicies(clusterID string) (controlPlaneUpgradePolicies []*cmv1.ControlPlaneUpgradePolicy, err error) {
	collection := c.ocm.ClustersMgmt().V1().
		Clusters().
		Cluster(clusterID).
		ControlPlane().
		UpgradePolicies()
	page := 1
	size := 100
	for {
		response, err := collection.List().
			Page(page).
			Size(size).
			Send()
		if err != nil {
			return nil, handleErr(response.Error(), err)
		}
		controlPlaneUpgradePolicies = append(controlPlaneUpgradePolicies, response.Items().Slice()...)
		if response.Size() < size {
			break
		}
		page++
	}
	return
}
