package util

import (
	configv1 "github.com/openshift/api/config/v1"
	cvorm "github.com/openshift/cluster-version-operator/lib/resourcemerge"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Common Kubernetes object annotations.
const (
	ReleaseVersionAnnotation = "release.openshift.io/version"
	CriticalPodAnnotation    = "scheduler.alpha.kubernetes.io/critical-pod"
)

// FilterString removes any instances of the needle from haystack.  It
// returns a new slice with all instances of needle removed, and a
// count of the number instances encountered.
func FilterString(haystack []string, needle string) ([]string, int) {
	var newSlice []string
	found := 0

	for _, x := range haystack {
		if x != needle {
			newSlice = append(newSlice, x)
		} else {
			found++
		}
	}

	return newSlice, found
}

// ReleaseVersionMatches checks whether a Kubernetes object has an OpenShift
// release version annotation that matches the given version.
func ReleaseVersionMatches(obj metav1.Object, version string) bool {
	annotations := obj.GetAnnotations()

	value, found := annotations[ReleaseVersionAnnotation]
	if !found || value != version {
		return false
	}

	return true
}

// DeploymentUpdated checks whether a Kubernetes deployment object's replicas
// are fully updated and available.
func DeploymentUpdated(dep *appsv1.Deployment) bool {
	if dep.Status.ObservedGeneration < dep.Generation {
		return false
	}

	if dep.Status.UpdatedReplicas != dep.Status.Replicas {
		return false
	}

	if dep.Status.AvailableReplicas == 0 {
		return false
	}

	return true
}

// ResetProgressingTime finds the Progressing condition in the given slice, or
// creates a default one if none is found, and sets the LastTransitionTime to
// the current time.
func ResetProgressingTime(conds *[]configv1.ClusterOperatorStatusCondition) {
	prog := cvorm.FindOperatorStatusCondition(*conds, configv1.OperatorProgressing)

	// If the Progressing condition wasn't found, set a default one.
	if prog == nil {
		prog = &configv1.ClusterOperatorStatusCondition{
			Type:   configv1.OperatorProgressing,
			Status: configv1.ConditionFalse,
		}
	}

	prog.LastTransitionTime = metav1.Now()

	cvorm.SetOperatorStatusCondition(conds, *prog)
}
