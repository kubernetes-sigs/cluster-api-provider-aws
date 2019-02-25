package payload

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/openshift/cluster-version-operator/lib"
)

var (
	metricPayloadErrors = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cluster_operator_payload_errors",
		Help: "Report the number of errors encountered applying the payload.",
	}, []string{"version"})
)

func init() {
	prometheus.MustRegister(
		metricPayloadErrors,
	)
}

// ResourceBuilder abstracts how a manifest is created on the server. Introduced for testing.
type ResourceBuilder interface {
	Apply(*lib.Manifest) error
}

type Task struct {
	Index    int
	Total    int
	Manifest *lib.Manifest
	Requeued int
	Backoff  wait.Backoff
}

func (st *Task) String() string {
	ns := st.Manifest.Object().GetNamespace()
	if len(ns) == 0 {
		return fmt.Sprintf("%s %q (%d of %d)", strings.ToLower(st.Manifest.GVK.Kind), st.Manifest.Object().GetName(), st.Index, st.Total)
	}
	return fmt.Sprintf("%s \"%s/%s\" (%d of %d)", strings.ToLower(st.Manifest.GVK.Kind), ns, st.Manifest.Object().GetName(), st.Index, st.Total)
}

func (st *Task) Run(version string, builder ResourceBuilder) error {
	var lastErr error
	if err := wait.ExponentialBackoff(st.Backoff, func() (bool, error) {
		// run builder for the manifest
		if err := builder.Apply(st.Manifest); err != nil {
			utilruntime.HandleError(errors.Wrapf(err, "error running apply for %s", st))
			lastErr = err
			metricPayloadErrors.WithLabelValues(version).Inc()
			if !shouldRequeueApplyOnErr(err) {
				return false, err
			}
			return false, nil
		}
		return true, nil
	}); err != nil {
		if uerr, ok := lastErr.(*UpdateError); ok {
			return uerr
		}
		reason, cause := reasonForPayloadSyncError(lastErr)
		if len(cause) > 0 {
			cause = ": " + cause
		}
		return &UpdateError{
			Nested:  lastErr,
			Reason:  reason,
			Message: fmt.Sprintf("Could not update %s%s", st, cause),
		}
	}
	return nil
}

func shouldRequeueApplyOnErr(err error) bool {
	if apierrors.IsInvalid(err) {
		return false
	}
	return true
}

// UpdateError is a wrapper for errors that occur during a payload sync.
type UpdateError struct {
	Nested  error
	Reason  string
	Message string
	Name    string
}

func (e *UpdateError) Error() string {
	return e.Message
}

func (e *UpdateError) Cause() error {
	return e.Nested
}

// reasonForUpdateError provides a succint explanation of a known error type for use in a human readable
// message during update. Since all objects in the image should be successfully applied, messages
// should direct the reader (likely a cluster administrator) to a possible cause in their own config.
func reasonForPayloadSyncError(err error) (string, string) {
	err = errors.Cause(err)
	switch {
	case apierrors.IsNotFound(err), apierrors.IsAlreadyExists(err):
		return "UpdatePayloadResourceNotFound", "resource may have been deleted"
	case apierrors.IsConflict(err):
		return "UpdatePayloadResourceConflict", "someone else is updating this resource"
	case apierrors.IsTimeout(err), apierrors.IsServiceUnavailable(err), apierrors.IsUnexpectedServerError(err):
		return "UpdatePayloadClusterDown", "the server is down or not responding"
	case apierrors.IsInternalError(err):
		return "UpdatePayloadClusterError", "the server is reporting an internal error"
	case apierrors.IsInvalid(err):
		return "UpdatePayloadResourceInvalid", "the object is invalid, possibly due to local cluster configuration"
	case apierrors.IsUnauthorized(err):
		return "UpdatePayloadClusterUnauthorized", "could not authenticate to the server"
	case apierrors.IsForbidden(err):
		return "UpdatePayloadResourceForbidden", "the server has forbidden updates to this resource"
	case apierrors.IsServerTimeout(err), apierrors.IsTooManyRequests(err):
		return "UpdatePayloadClusterOverloaded", "the server is overloaded and is not accepting updates"
	case meta.IsNoMatchError(err):
		return "UpdatePayloadResourceTypeMissing", "the server does not recognize this resource, check extension API servers"
	default:
		return "UpdatePayloadFailed", ""
	}
}

func SummaryForReason(reason, name string) string {
	switch reason {

	// likely temporary errors
	case "UpdatePayloadResourceNotFound", "UpdatePayloadResourceConflict":
		return "some resources could not be updated"
	case "UpdatePayloadClusterDown":
		return "the control plane is down or not responding"
	case "UpdatePayloadClusterError":
		return "the control plane is reporting an internal error"
	case "UpdatePayloadClusterOverloaded":
		return "the control plane is overloaded and is not accepting updates"
	case "UpdatePayloadClusterUnauthorized":
		return "could not authenticate to the server"
	case "UpdatePayloadRetrievalFailed":
		return "could not download the update"

	// likely a policy or other configuration error due to end user action
	case "UpdatePayloadResourceForbidden":
		return "the server is rejecting updates"

	// the image may not be correct, or the cluster may be in an unexpected
	// state
	case "UpdatePayloadResourceTypeMissing":
		return "a required extension is not available to update"
	case "UpdatePayloadResourceInvalid":
		return "some cluster configuration is invalid"
	case "UpdatePayloadIntegrity":
		return "the contents of the update are invalid"

	case "ClusterOperatorFailing":
		if len(name) > 0 {
			return fmt.Sprintf("the cluster operator %s is failing", name)
		}
		return "a cluster operator is failing"
	case "ClusterOperatorNotAvailable":
		if len(name) > 0 {
			return fmt.Sprintf("the cluster operator %s has not yet successfully rolled out", name)
		}
		return "a cluster operator has not yet rolled out"
	}

	if strings.HasPrefix(reason, "UpdatePayload") {
		return "the update could not be applied"
	}
	return "an unknown error has occurred"
}
