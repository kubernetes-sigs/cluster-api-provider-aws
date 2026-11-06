package rosa

import (
	"fmt"
	"strings"

	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	ocmerrors "github.com/openshift-online/ocm-sdk-go/errors"
	"github.com/openshift/rosa/pkg/ocm"
	errors "github.com/zgalor/weberr"
)

// IsNodePoolReady checks whether the NodePool is provisioned and all replicas are available.
// A NodePool with fixed replicas is ready once its current replicas match the requested replicas,
// including replicas: 0, which is ready when no replicas are running.
// If autoscaling is enabled, NodePool must have replicas >= autoscaling.MinReplica to be considered ready.
func IsNodePoolReady(nodePool *cmv1.NodePool) bool {
	if nodePool.Status().Message() != "" {
		return false
	}

	if nodePool.Replicas() == 0 && nodePool.Autoscaling() != nil {
		return nodePool.Status().CurrentReplicas() >= nodePool.Autoscaling().MinReplica()
	}

	return nodePool.Replicas() == nodePool.Status().CurrentReplicas()
}

// CreateVersionID creates a version ID based on the version, channelGroup, and channel.
// If channel is set, it takes priority over channelGroup (e.g., "fast-4.22" -> channelGroup "fast").
// If channelGroup is set and channel is not, channelGroup is used.
// If neither is set, the default channel group ("stable") behavior applies.
func CreateVersionID(version string, channelGroup string, channel string) string {
	if channel != "" {
		if idx := strings.Index(channel, "-"); idx > 0 {
			channelGroup = channel[:idx]
		}
	}
	if channelGroup == "" {
		channelGroup = ocm.DefaultChannelGroup
	}
	versionID := fmt.Sprintf("%s%s", ocm.VersionPrefix, version)
	if channelGroup != ocm.DefaultChannelGroup {
		versionID = fmt.Sprintf("%s-%s", versionID, channelGroup)
	}
	return versionID
}

func handleErr(res *ocmerrors.Error, err error) error {
	msg := res.Reason()
	if msg == "" {
		msg = err.Error()
	}
	// Hack to always display the correct terms and conditions message
	if res.Code() == "CLUSTERS-MGMT-451" {
		msg = "You must accept the Terms and Conditions in order to continue.\n" +
			"Go to https://www.redhat.com/wapps/tnc/ackrequired?site=ocm&event=register\n" +
			"Once you accept the terms, you will need to retry the action that was blocked."
	}
	errType := errors.ErrorType(res.Status()) //#nosec G115
	return errType.Set(errors.Errorf("%s", msg))
}
