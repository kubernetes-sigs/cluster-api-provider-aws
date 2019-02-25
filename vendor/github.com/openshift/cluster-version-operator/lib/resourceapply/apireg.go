package resourceapply

import (
	"github.com/openshift/cluster-version-operator/lib/resourcemerge"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiregv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	apiregclientv1 "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset/typed/apiregistration/v1"
	"k8s.io/utils/pointer"
)

func ApplyAPIService(client apiregclientv1.APIServicesGetter, required *apiregv1.APIService) (*apiregv1.APIService, bool, error) {
	existing, err := client.APIServices().Get(required.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		actual, err := client.APIServices().Create(required)
		return actual, true, err
	}
	if err != nil {
		return nil, false, err
	}

	modified := pointer.BoolPtr(false)
	resourcemerge.EnsureAPIService(modified, existing, *required)
	if !*modified {
		return existing, false, nil
	}

	actual, err := client.APIServices().Update(existing)
	return actual, true, err
}
