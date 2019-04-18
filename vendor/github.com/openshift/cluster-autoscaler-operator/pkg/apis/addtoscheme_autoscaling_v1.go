package apis

import (
	v1 "github.com/openshift/cluster-autoscaler-operator/pkg/apis/autoscaling/v1"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes, v1.SchemeBuilder.AddToScheme)
}
