package resourcebuilder

import (
	securityclientv1 "github.com/openshift/client-go/security/clientset/versioned/typed/security/v1"
	"github.com/openshift/cluster-version-operator/lib"
	"github.com/openshift/cluster-version-operator/lib/resourceapply"
	"github.com/openshift/cluster-version-operator/lib/resourceread"
	"k8s.io/client-go/rest"
)

type securityBuilder struct {
	client   *securityclientv1.SecurityV1Client
	raw      []byte
	modifier MetaV1ObjectModifierFunc
}

func newSecurityBuilder(config *rest.Config, m lib.Manifest) Interface {
	return &securityBuilder{
		client: securityclientv1.NewForConfigOrDie(withProtobuf(config)),
		raw:    m.Raw,
	}
}

func (b *securityBuilder) WithModifier(f MetaV1ObjectModifierFunc) Interface {
	b.modifier = f
	return b
}

func (b *securityBuilder) Do() error {
	scc := resourceread.ReadSecurityContextConstraintsV1OrDie(b.raw)
	if b.modifier != nil {
		b.modifier(scc)
	}
	_, _, err := resourceapply.ApplySecurityContextConstraints(b.client, scc)
	return err
}
