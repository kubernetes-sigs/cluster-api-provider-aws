package cvo

import (
	"bytes"
	"fmt"

	"github.com/golang/glog"

	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/diff"
	"k8s.io/apimachinery/pkg/util/validation/field"

	configv1 "github.com/openshift/api/config/v1"
	configclientv1 "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"

	"github.com/openshift/cluster-version-operator/lib/resourcemerge"
	"github.com/openshift/cluster-version-operator/pkg/payload"
)

func mergeEqualVersions(current *configv1.UpdateHistory, desired configv1.Update) bool {
	if len(desired.Image) > 0 && desired.Image == current.Image {
		if len(current.Version) == 0 || desired.Version == current.Version {
			current.Version = desired.Version
			return true
		}
	}
	if len(desired.Version) > 0 && desired.Version == current.Version {
		if len(current.Image) == 0 || desired.Image == current.Image {
			current.Image = desired.Image
			return true
		}
	}
	return false
}

func mergeOperatorHistory(config *configv1.ClusterVersion, desired configv1.Update, now metav1.Time, completed bool) {
	// if we have no image, we cannot reproduce the update later and so it cannot be part of the history
	if len(desired.Image) == 0 {
		// make the array empty
		if config.Status.History == nil {
			config.Status.History = []configv1.UpdateHistory{}
		}
		return
	}

	if len(config.Status.History) == 0 {
		config.Status.History = append(config.Status.History, configv1.UpdateHistory{
			Version: desired.Version,
			Image:   desired.Image,

			State:       configv1.PartialUpdate,
			StartedTime: now,
		})
	}

	last := &config.Status.History[0]

	if !mergeEqualVersions(last, desired) {
		last.CompletionTime = &now
		config.Status.History = append([]configv1.UpdateHistory{
			{
				Version: desired.Version,
				Image:   desired.Image,

				State:       configv1.PartialUpdate,
				StartedTime: now,
			},
		}, config.Status.History...)
		last = &config.Status.History[0]
	}

	if len(config.Status.History) > 10 {
		config.Status.History = config.Status.History[:10]
	}

	if completed {
		last.State = configv1.CompletedUpdate
		if last.CompletionTime == nil {
			last.CompletionTime = &now
		}
	}
	if len(last.State) == 0 {
		last.State = configv1.PartialUpdate
	}

	config.Status.Desired = desired
}

// ClusterVersionInvalid indicates that the cluster version has an error that prevents the server from
// taking action. The cluster version operator will only reconcile the current state as long as this
// condition is set.
const ClusterVersionInvalid configv1.ClusterStatusConditionType = "Invalid"

// syncStatus calculates the new status of the ClusterVersion based on the current sync state and any
// validation errors found. We allow the caller to pass the original object to avoid DeepCopying twice.
func (optr *Operator) syncStatus(original, config *configv1.ClusterVersion, status *SyncWorkerStatus, validationErrs field.ErrorList) error {
	glog.V(5).Infof("Synchronizing errs=%#v status=%#v", validationErrs, status)

	// update the config with the latest available updates
	if updated := optr.getAvailableUpdates().NeedsUpdate(config); updated != nil {
		config = updated
	} else if original == nil || original == config {
		original = config.DeepCopy()
	}

	config.Status.ObservedGeneration = status.Generation
	if len(status.VersionHash) > 0 {
		config.Status.VersionHash = status.VersionHash
	}

	now := metav1.Now()
	version := versionString(status.Actual)

	// update validation errors
	var reason string
	if len(validationErrs) > 0 {
		buf := &bytes.Buffer{}
		if len(validationErrs) == 1 {
			fmt.Fprintf(buf, "The cluster version is invalid: %s", validationErrs[0].Error())
		} else {
			fmt.Fprintf(buf, "The cluster version is invalid:\n")
			for _, err := range validationErrs {
				fmt.Fprintf(buf, "* %s\n", err.Error())
			}
		}
		reason = "InvalidClusterVersion"

		resourcemerge.SetOperatorStatusCondition(&config.Status.Conditions, configv1.ClusterOperatorStatusCondition{
			Type:               ClusterVersionInvalid,
			Status:             configv1.ConditionTrue,
			Reason:             reason,
			Message:            buf.String(),
			LastTransitionTime: now,
		})
	} else {
		resourcemerge.RemoveOperatorStatusCondition(&config.Status.Conditions, ClusterVersionInvalid)
	}

	// set the available condition
	if status.Completed > 0 {
		resourcemerge.SetOperatorStatusCondition(&config.Status.Conditions, configv1.ClusterOperatorStatusCondition{
			Type:    configv1.OperatorAvailable,
			Status:  configv1.ConditionTrue,
			Message: fmt.Sprintf("Done applying %s", version),

			LastTransitionTime: now,
		})
	}
	// default the available condition if not set
	if resourcemerge.FindOperatorStatusCondition(config.Status.Conditions, configv1.OperatorAvailable) == nil {
		resourcemerge.SetOperatorStatusCondition(&config.Status.Conditions, configv1.ClusterOperatorStatusCondition{
			Type:               configv1.OperatorAvailable,
			Status:             configv1.ConditionFalse,
			LastTransitionTime: now,
		})
	}

	if err := status.Failure; err != nil {
		var reason string
		msg := "an error occurred"
		if uErr, ok := err.(*payload.UpdateError); ok {
			reason = uErr.Reason
			msg = payload.SummaryForReason(reason, uErr.Name)
		}

		// set the failing condition
		resourcemerge.SetOperatorStatusCondition(&config.Status.Conditions, configv1.ClusterOperatorStatusCondition{
			Type:               configv1.OperatorFailing,
			Status:             configv1.ConditionTrue,
			Reason:             reason,
			Message:            err.Error(),
			LastTransitionTime: now,
		})

		// update progressing
		if status.Reconciling {
			resourcemerge.SetOperatorStatusCondition(&config.Status.Conditions, configv1.ClusterOperatorStatusCondition{
				Type:               configv1.OperatorProgressing,
				Status:             configv1.ConditionFalse,
				Reason:             reason,
				Message:            fmt.Sprintf("Error while reconciling %s: %s", version, msg),
				LastTransitionTime: now,
			})
		} else {
			resourcemerge.SetOperatorStatusCondition(&config.Status.Conditions, configv1.ClusterOperatorStatusCondition{
				Type:               configv1.OperatorProgressing,
				Status:             configv1.ConditionTrue,
				Reason:             reason,
				Message:            fmt.Sprintf("Unable to apply %s: %s", version, msg),
				LastTransitionTime: now,
			})
		}

	} else {
		// clear the failure condition
		resourcemerge.SetOperatorStatusCondition(&config.Status.Conditions, configv1.ClusterOperatorStatusCondition{Type: configv1.OperatorFailing, Status: configv1.ConditionFalse, LastTransitionTime: now})

		// update progressing
		if status.Reconciling {
			message := fmt.Sprintf("Cluster version is %s", version)
			if len(validationErrs) > 0 {
				message = fmt.Sprintf("Stopped at %s: the cluster version is invalid", version)
			}
			resourcemerge.SetOperatorStatusCondition(&config.Status.Conditions, configv1.ClusterOperatorStatusCondition{
				Type:               configv1.OperatorProgressing,
				Status:             configv1.ConditionFalse,
				Reason:             reason,
				Message:            message,
				LastTransitionTime: now,
			})
		} else {
			var message string
			switch {
			case len(validationErrs) > 0:
				message = fmt.Sprintf("Reconciling %s: the cluster version is invalid", version)
			case status.Fraction > 0:
				message = fmt.Sprintf("Working towards %s: %.0f%% complete", version, status.Fraction*100)
			case status.Step == "RetrievePayload":
				if len(reason) == 0 {
					reason = "DownloadingUpdate"
				}
				message = fmt.Sprintf("Working towards %s: downloading update", version)
			default:
				message = fmt.Sprintf("Working towards %s", version)
			}
			resourcemerge.SetOperatorStatusCondition(&config.Status.Conditions, configv1.ClusterOperatorStatusCondition{
				Type:               configv1.OperatorProgressing,
				Status:             configv1.ConditionTrue,
				Reason:             reason,
				Message:            message,
				LastTransitionTime: now,
			})
		}
	}

	// default retrieved updates if it is not set
	if resourcemerge.FindOperatorStatusCondition(config.Status.Conditions, configv1.RetrievedUpdates) == nil {
		resourcemerge.SetOperatorStatusCondition(&config.Status.Conditions, configv1.ClusterOperatorStatusCondition{
			Type:               configv1.RetrievedUpdates,
			Status:             configv1.ConditionFalse,
			LastTransitionTime: now,
		})
	}

	mergeOperatorHistory(config, status.Actual, now, status.Completed > 0)

	if glog.V(6) {
		glog.Infof("Apply config: %s", diff.ObjectReflectDiff(original, config))
	}
	updated, err := applyClusterVersionStatus(optr.client.ConfigV1(), config, original)
	optr.rememberLastUpdate(updated)
	return err
}

// syncDegradedStatus handles generic errors in the cluster version. It tries to preserve
// all status fields that it can by using the provided config or loading the latest version
// from the cache (instead of clearing the status).
// if ierr is nil, return nil
// if ierr is not nil, update OperatorStatus as Failing and return ierr
func (optr *Operator) syncFailingStatus(config *configv1.ClusterVersion, ierr error) error {
	if ierr == nil {
		return nil
	}

	// try to reuse the most recent status if available
	if config == nil {
		config, _ = optr.cvLister.Get(optr.name)
	}
	if config == nil {
		config = &configv1.ClusterVersion{
			ObjectMeta: metav1.ObjectMeta{
				Name: optr.name,
			},
		}
	}

	original := config.DeepCopy()

	now := metav1.Now()
	msg := fmt.Sprintf("Error ensuring the cluster version is up to date: %v", ierr)

	// clear the available condition
	resourcemerge.SetOperatorStatusCondition(&config.Status.Conditions, configv1.ClusterOperatorStatusCondition{Type: configv1.OperatorAvailable, Status: configv1.ConditionFalse, LastTransitionTime: now})

	// reset the failing message
	resourcemerge.SetOperatorStatusCondition(&config.Status.Conditions, configv1.ClusterOperatorStatusCondition{
		Type:               configv1.OperatorFailing,
		Status:             configv1.ConditionTrue,
		Message:            ierr.Error(),
		LastTransitionTime: now,
	})

	// preserve the status of the existing progressing condition
	progressingStatus := configv1.ConditionFalse
	if resourcemerge.IsOperatorStatusConditionTrue(config.Status.Conditions, configv1.OperatorProgressing) {
		progressingStatus = configv1.ConditionTrue
	}
	resourcemerge.SetOperatorStatusCondition(&config.Status.Conditions, configv1.ClusterOperatorStatusCondition{
		Type:               configv1.OperatorProgressing,
		Status:             progressingStatus,
		Message:            msg,
		LastTransitionTime: now,
	})

	mergeOperatorHistory(config, optr.currentVersion(), now, false)

	updated, err := applyClusterVersionStatus(optr.client.ConfigV1(), config, original)
	optr.rememberLastUpdate(updated)
	if err != nil {
		return err
	}
	return ierr
}

// applyClusterVersionStatus attempts to overwrite the status subresource of required. If
// original is provided it is compared to required and no update will be made if the
// object does not change. The method will retry a conflict by retrieving the latest live
// version and updating the metadata of required. required is modified if the object on
// the server is newer.
func applyClusterVersionStatus(client configclientv1.ClusterVersionsGetter, required, original *configv1.ClusterVersion) (*configv1.ClusterVersion, error) {
	if original != nil && equality.Semantic.DeepEqual(&original.Status, &required.Status) {
		return required, nil
	}
	actual, err := client.ClusterVersions().UpdateStatus(required)
	if apierrors.IsConflict(err) {
		existing, cErr := client.ClusterVersions().Get(required.Name, metav1.GetOptions{})
		if err != nil {
			return nil, cErr
		}
		if existing.UID != required.UID {
			return nil, fmt.Errorf("cluster version was deleted and recreated, cannot update status")
		}
		if equality.Semantic.DeepEqual(&existing.Status, &required.Status) {
			return existing, nil
		}
		required.ObjectMeta = existing.ObjectMeta
		actual, err = client.ClusterVersions().UpdateStatus(required)
	}
	if err != nil {
		return nil, err
	}
	required.ObjectMeta = actual.ObjectMeta
	return actual, nil
}
