package rosa

import (
	"fmt"

	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	ocmerrors "github.com/openshift-online/ocm-sdk-go/errors"
	errors "github.com/zgalor/weberr"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
)

// IsNodePoolReady checkes whether the nodepool is provisoned and all replicas are available.
// If autosacling is enabled, NodePool must have replicas >= autosacling.MinReplica to be considered ready.
func IsNodePoolReady(nodePool *cmv1.NodePool) bool {
	if nodePool.Status().Message() != "" {
		return false
	}

	if nodePool.Replicas() != 0 {
		return nodePool.Replicas() == nodePool.Status().CurrentReplicas()
	}

	if nodePool.Autoscaling() != nil {
		return nodePool.Status().CurrentReplicas() >= nodePool.Autoscaling().MinReplica()
	}

	return false
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

// Reporter is a helper struct used by rosa CLI runtime to report errors.
// We need it to limit number of reports by rosa CLI runtime.
type Reporter struct {
	scope scope.RosaRoleConfigScope
}

// Debugf prints a debug message with the given format and arguments.
func (r *Reporter) Debugf(format string, args ...interface{}) {
}

// Infof prints an informative message with the given format and arguments.
func (r *Reporter) Infof(format string, args ...interface{}) {
}

// Warnf prints an warning message with the given format and arguments.
func (r *Reporter) Warnf(format string, args ...interface{}) {
}

// Errorf prints an error message with the given format and arguments. It also return an error
// containing the same information, which will be usually discarded, except when the caller needs to
// report the error and also return it.
func (r *Reporter) Errorf(format string, args ...interface{}) error {
	r.scope.Error(fmt.Errorf(format, args...), format, args...)
	return nil
}

// IsTerminal indicates that the reporter is terminal.
func (r *Reporter) IsTerminal() bool {
	return true
}
