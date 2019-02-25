package resourceapply

import (
	"github.com/openshift/cluster-version-operator/lib/resourcemerge"
	batchv1 "k8s.io/api/batch/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	batchclientv1 "k8s.io/client-go/kubernetes/typed/batch/v1"
	"k8s.io/utils/pointer"
)

// ApplyJob applies the required Job to the cluster.
func ApplyJob(client batchclientv1.JobsGetter, required *batchv1.Job) (*batchv1.Job, bool, error) {
	existing, err := client.Jobs(required.Namespace).Get(required.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		actual, err := client.Jobs(required.Namespace).Create(required)
		return actual, true, err
	}
	if err != nil {
		return nil, false, err
	}

	modified := pointer.BoolPtr(false)
	resourcemerge.EnsureJob(modified, existing, *required)
	if !*modified {
		return existing, false, nil
	}

	actual, err := client.Jobs(required.Namespace).Update(existing)
	return actual, true, err
}
