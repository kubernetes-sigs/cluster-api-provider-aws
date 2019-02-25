package resourceread

import (
	securityv1 "github.com/openshift/api/security/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	securityScheme = runtime.NewScheme()
	securityCodecs = serializer.NewCodecFactory(securityScheme)
)

func init() {
	if err := securityv1.AddToScheme(securityScheme); err != nil {
		panic(err)
	}
}

// ReadSecurityContextConstraintsV1OrDie reads clusterrolebinding object from bytes. Panics on error.
func ReadSecurityContextConstraintsV1OrDie(objBytes []byte) *securityv1.SecurityContextConstraints {
	requiredObj, err := runtime.Decode(securityCodecs.UniversalDecoder(securityv1.SchemeGroupVersion), objBytes)
	if err != nil {
		panic(err)
	}
	return requiredObj.(*securityv1.SecurityContextConstraints)
}
