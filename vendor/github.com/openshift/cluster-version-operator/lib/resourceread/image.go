package resourceread

import (
	imagev1 "github.com/openshift/api/image/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	imageScheme = runtime.NewScheme()
	imageCodecs = serializer.NewCodecFactory(imageScheme)
)

func init() {
	if err := imagev1.AddToScheme(imageScheme); err != nil {
		panic(err)
	}
}

// ReadImageStreamV1 reads imagestream object from bytes or reports an error.
func ReadImageStreamV1(objBytes []byte) (*imagev1.ImageStream, error) {
	requiredObj, err := runtime.Decode(imageCodecs.UniversalDecoder(imagev1.SchemeGroupVersion), objBytes)
	if err != nil {
		return nil, err
	}
	return requiredObj.(*imagev1.ImageStream), nil
}
