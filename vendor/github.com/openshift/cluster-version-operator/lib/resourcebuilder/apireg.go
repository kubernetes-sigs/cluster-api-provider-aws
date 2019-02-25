package resourcebuilder

import (
	"github.com/openshift/cluster-version-operator/lib"
	"github.com/openshift/cluster-version-operator/lib/resourceapply"
	"github.com/openshift/cluster-version-operator/lib/resourceread"
	"k8s.io/client-go/rest"
	apiregclientv1 "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset/typed/apiregistration/v1"
)

type apiServiceBuilder struct {
	client   *apiregclientv1.ApiregistrationV1Client
	raw      []byte
	modifier MetaV1ObjectModifierFunc
}

func newAPIServiceBuilder(config *rest.Config, m lib.Manifest) Interface {
	return &apiServiceBuilder{
		client: apiregclientv1.NewForConfigOrDie(config),
		raw:    m.Raw,
	}
}

func (b *apiServiceBuilder) WithModifier(f MetaV1ObjectModifierFunc) Interface {
	b.modifier = f
	return b
}

func (b *apiServiceBuilder) Do() error {
	apiService := resourceread.ReadAPIServiceV1OrDie(b.raw)
	if b.modifier != nil {
		b.modifier(apiService)
	}
	_, _, err := resourceapply.ApplyAPIService(b.client, apiService)
	return err
}
