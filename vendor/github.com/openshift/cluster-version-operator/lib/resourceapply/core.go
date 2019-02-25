package resourceapply

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	coreclientv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/utils/pointer"

	"github.com/openshift/cluster-version-operator/lib/resourcemerge"
)

// ApplyNamespace merges objectmeta, does not worry about anything else
func ApplyNamespace(client coreclientv1.NamespacesGetter, required *corev1.Namespace) (*corev1.Namespace, bool, error) {
	existing, err := client.Namespaces().Get(required.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		actual, err := client.Namespaces().Create(required)
		return actual, true, err
	}
	if err != nil {
		return nil, false, err
	}

	modified := pointer.BoolPtr(false)
	resourcemerge.EnsureObjectMeta(modified, &existing.ObjectMeta, required.ObjectMeta)
	if !*modified {
		return existing, false, nil
	}

	actual, err := client.Namespaces().Update(existing)
	return actual, true, err
}

// ApplyService merges objectmeta and requires
// TODO, since this cannot determine whether changes are due to legitimate actors (api server) or illegitimate ones (users), we cannot update
// TODO I've special cased the selector for now
func ApplyService(client coreclientv1.ServicesGetter, required *corev1.Service) (*corev1.Service, bool, error) {
	existing, err := client.Services(required.Namespace).Get(required.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		actual, err := client.Services(required.Namespace).Create(required)
		return actual, true, err
	}
	if err != nil {
		return nil, false, err
	}

	modified := pointer.BoolPtr(false)
	resourcemerge.EnsureObjectMeta(modified, &existing.ObjectMeta, required.ObjectMeta)
	selectorSame := equality.Semantic.DeepEqual(existing.Spec.Selector, required.Spec.Selector)
	typeSame := equality.Semantic.DeepEqual(existing.Spec.Type, required.Spec.Type)
	if selectorSame && typeSame && !*modified {
		return nil, false, nil
	}
	existing.Spec.Selector = required.Spec.Selector
	existing.Spec.Type = required.Spec.Type // if this is different, the update will fail.  Status will indicate it.

	actual, err := client.Services(required.Namespace).Update(existing)
	return actual, true, err
}

// ApplyServiceAccount applies the required serviceaccount to the cluster.
func ApplyServiceAccount(client coreclientv1.ServiceAccountsGetter, required *corev1.ServiceAccount) (*corev1.ServiceAccount, bool, error) {
	existing, err := client.ServiceAccounts(required.Namespace).Get(required.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		actual, err := client.ServiceAccounts(required.Namespace).Create(required)
		return actual, true, err
	}
	if err != nil {
		return nil, false, err
	}

	modified := pointer.BoolPtr(false)
	resourcemerge.EnsureObjectMeta(modified, &existing.ObjectMeta, required.ObjectMeta)
	if !*modified {
		return existing, false, nil
	}

	actual, err := client.ServiceAccounts(required.Namespace).Update(existing)
	return actual, true, err
}

// ApplyConfigMap applies the required serviceaccount to the cluster.
func ApplyConfigMap(client coreclientv1.ConfigMapsGetter, required *corev1.ConfigMap) (*corev1.ConfigMap, bool, error) {
	existing, err := client.ConfigMaps(required.Namespace).Get(required.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		actual, err := client.ConfigMaps(required.Namespace).Create(required)
		return actual, true, err
	}
	if err != nil {
		return nil, false, err
	}

	modified := pointer.BoolPtr(false)
	resourcemerge.EnsureConfigMap(modified, existing, *required)
	if !*modified {
		return existing, false, nil
	}

	actual, err := client.ConfigMaps(required.Namespace).Update(existing)
	return actual, true, err
}
