package resourceapply

import (
	"github.com/openshift/cluster-version-operator/lib/resourcemerge"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextclientv1beta1 "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	apiextlistersv1beta1 "k8s.io/apiextensions-apiserver/pkg/client/listers/apiextensions/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

func ApplyCustomResourceDefinition(client apiextclientv1beta1.CustomResourceDefinitionsGetter, required *apiextv1beta1.CustomResourceDefinition) (*apiextv1beta1.CustomResourceDefinition, bool, error) {
	existing, err := client.CustomResourceDefinitions().Get(required.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		actual, err := client.CustomResourceDefinitions().Create(required)
		return actual, true, err
	}
	if err != nil {
		return nil, false, err
	}

	modified := pointer.BoolPtr(false)
	resourcemerge.EnsureCustomResourceDefinition(modified, existing, *required)
	if !*modified {
		return existing, false, nil
	}

	actual, err := client.CustomResourceDefinitions().Update(existing)
	return actual, true, err
}

func ApplyCustomResourceDefinitionFromCache(lister apiextlistersv1beta1.CustomResourceDefinitionLister, client apiextclientv1beta1.CustomResourceDefinitionsGetter, required *apiextv1beta1.CustomResourceDefinition) (*apiextv1beta1.CustomResourceDefinition, bool, error) {
	existing, err := lister.Get(required.Name)
	if apierrors.IsNotFound(err) {
		actual, err := client.CustomResourceDefinitions().Create(required)
		return actual, true, err
	}
	if err != nil {
		return nil, false, err
	}

	existing = existing.DeepCopy()
	modified := pointer.BoolPtr(false)
	resourcemerge.EnsureCustomResourceDefinition(modified, existing, *required)
	if !*modified {
		return existing, false, nil
	}

	actual, err := client.CustomResourceDefinitions().Update(existing)
	return actual, true, err
}
