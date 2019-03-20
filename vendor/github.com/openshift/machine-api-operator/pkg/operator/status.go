package operator

import (
	"fmt"
	"reflect"
	"strings"

	v1 "k8s.io/api/core/v1"

	"github.com/golang/glog"

	osconfigv1 "github.com/openshift/api/config/v1"
	cvoresourcemerge "github.com/openshift/cluster-version-operator/lib/resourcemerge"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// StatusReason is a MixedCaps string representing the reason for a
// status condition change.
type StatusReason string

// The default set of status change reasons.
const (
	ReasonEmpty      StatusReason = ""
	ReasonSyncing    StatusReason = "SyncingResources"
	ReasonSyncFailed StatusReason = "SyncingFailed"
)

const (
	clusterOperatorName = "machine-api"
)

// statusProgressing sets the Progressing condition to True, with the given
// reason and message, and sets both the Available and Failing conditions to
// False.
func (optr *Operator) statusProgressing() error {
	desiredVersions := optr.operandVersions
	currentVersions, err := optr.getCurrentVersions()
	if err != nil {
		glog.Errorf("Error getting operator current versions: %v", err)
		return err
	}
	var isProgressing osconfigv1.ConditionStatus

	co, err := optr.getOrCreateClusterOperator()
	if err != nil {
		glog.Errorf("Failed to get or create Cluster Operator: %v", err)
		return err
	}

	var message string
	if !reflect.DeepEqual(desiredVersions, currentVersions) {
		glog.V(2).Info("Syncing status: progressing")
		message = fmt.Sprintf("Progressing towards %s", optr.printOperandVersions())
		optr.eventRecorder.Eventf(co, v1.EventTypeNormal, "Status upgrade", message)
		isProgressing = osconfigv1.ConditionTrue
	} else {
		glog.V(2).Info("Syncing status: re-syncing")
		message = fmt.Sprintf("Running resync for %s", optr.printOperandVersions())
		isProgressing = osconfigv1.ConditionFalse
	}

	conds := []osconfigv1.ClusterOperatorStatusCondition{
		{
			Type:    osconfigv1.OperatorProgressing,
			Status:  isProgressing,
			Reason:  string(ReasonSyncing),
			Message: message,
		},
		{
			Type:   osconfigv1.OperatorAvailable,
			Status: osconfigv1.ConditionTrue,
		},
		{
			Type:   osconfigv1.OperatorFailing,
			Status: osconfigv1.ConditionFalse,
		},
	}

	return optr.syncStatus(co, conds)
}

// statusAvailable sets the Available condition to True, with the given reason
// and message, and sets both the Progressing and Failing conditions to False.
func (optr *Operator) statusAvailable() error {
	conds := []osconfigv1.ClusterOperatorStatusCondition{
		{
			Type:    osconfigv1.OperatorAvailable,
			Status:  osconfigv1.ConditionTrue,
			Reason:  string(ReasonEmpty),
			Message: fmt.Sprintf("Cluster Machine API Operator is available at %s", optr.printOperandVersions()),
		},
		{
			Type:   osconfigv1.OperatorProgressing,
			Status: osconfigv1.ConditionFalse,
		},

		{
			Type:   osconfigv1.OperatorFailing,
			Status: osconfigv1.ConditionFalse,
		},
	}

	co, err := optr.getOrCreateClusterOperator()
	if err != nil {
		return err
	}

	// 	important: we only write the version field if we report available at the present level
	co.Status.Versions = optr.operandVersions
	glog.V(2).Info("Syncing status: available")
	return optr.syncStatus(co, conds)
}

// statusFailing sets the Failing condition to True, with the given reason and
// message, and sets the Progressing condition to False, and the Available
// condition to True.  This indicates that the operator is present and may be
// partially functioning, but is in a degraded or failing state.
func (optr *Operator) statusFailing(error string) error {
	desiredVersions := optr.operandVersions
	currentVersions, err := optr.getCurrentVersions()
	if err != nil {
		glog.Errorf("Error getting current versions: %v", err)
		return err
	}

	var message string
	if !reflect.DeepEqual(desiredVersions, currentVersions) {
		message = fmt.Sprintf("Failed when progressing towards %s because %s", optr.printOperandVersions(), error)
	} else {
		message = fmt.Sprintf("Failed to resync for %s because %s", optr.printOperandVersions(), error)
	}

	conds := []osconfigv1.ClusterOperatorStatusCondition{
		{
			Type:    osconfigv1.OperatorFailing,
			Status:  osconfigv1.ConditionTrue,
			Reason:  string(ReasonSyncFailed),
			Message: message,
		},
		{
			Type:   osconfigv1.OperatorProgressing,
			Status: osconfigv1.ConditionFalse,
		},
		{
			Type:   osconfigv1.OperatorAvailable,
			Status: osconfigv1.ConditionTrue,
		},
	}

	co, err := optr.getOrCreateClusterOperator()
	if err != nil {
		return err
	}
	optr.eventRecorder.Eventf(co, v1.EventTypeWarning, "Status failing", error)
	glog.V(2).Info("Syncing status: failing")
	return optr.syncStatus(co, conds)
}

//syncStatus applies the new condition to the mao ClusterOperator object.
func (optr *Operator) syncStatus(co *osconfigv1.ClusterOperator, conds []osconfigv1.ClusterOperatorStatusCondition) error {
	for _, c := range conds {
		cvoresourcemerge.SetOperatorStatusCondition(&co.Status.Conditions, c)
	}

	_, err := optr.osClient.ConfigV1().ClusterOperators().UpdateStatus(co)
	return err
}

func (optr *Operator) getOrCreateClusterOperator() (*osconfigv1.ClusterOperator, error) {
	co, err := optr.osClient.ConfigV1().ClusterOperators().Get(clusterOperatorName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		// to report the status of all the managed components.
		// TODO we will report the version of the operands (so our machine api implementation version)
		// NOTE: related objects lets openshift/must-gather collect diagnostic content
		co = &osconfigv1.ClusterOperator{
			ObjectMeta: metav1.ObjectMeta{
				Name: clusterOperatorName,
			},
			Status: osconfigv1.ClusterOperatorStatus{
				Versions: optr.operandVersions,
				RelatedObjects: []osconfigv1.ObjectReference{
					{
						Group:    "",
						Resource: "namespaces",
						Name:     optr.namespace,
					},
				},
			},
		}

		glog.Infof("%s clusterOperator status does not exist, creating %v", clusterOperatorName, co)
		co, err := optr.osClient.ConfigV1().ClusterOperators().Create(co)
		if err != nil {
			return nil, err
		}
		return co, nil
	}
	if err != nil {
		return nil, err
	}
	return co, nil
}

func (optr *Operator) getCurrentVersions() ([]osconfigv1.OperandVersion, error) {
	co, err := optr.getOrCreateClusterOperator()
	if err != nil {
		return nil, err
	}
	return co.Status.Versions, nil
}

func (optr *Operator) printOperandVersions() string {
	versionsOutput := []string{}
	for _, operand := range optr.operandVersions {
		versionsOutput = append(versionsOutput, fmt.Sprintf("%s: %s", operand.Name, operand.Version))
	}
	return strings.Join(versionsOutput, ", ")
}
