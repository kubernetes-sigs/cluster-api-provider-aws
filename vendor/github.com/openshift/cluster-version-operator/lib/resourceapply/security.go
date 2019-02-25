package resourceapply

import (
	securityv1 "github.com/openshift/api/security/v1"
	securityclientv1 "github.com/openshift/client-go/security/clientset/versioned/typed/security/v1"
	"github.com/openshift/cluster-version-operator/lib/resourcemerge"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

// ApplySecurityContextConstraints applies the required SecurityContextConstraints to the cluster.
func ApplySecurityContextConstraints(client securityclientv1.SecurityContextConstraintsGetter, required *securityv1.SecurityContextConstraints) (*securityv1.SecurityContextConstraints, bool, error) {
	existing, err := client.SecurityContextConstraints().Get(required.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		actual, err := client.SecurityContextConstraints().Create(required)
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

	actual, err := client.SecurityContextConstraints().Update(existing)
	return actual, true, err
}
