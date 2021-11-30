package controllers

import (
	"sigs.k8s.io/controller-runtime/pkg/controller"

	"sigs.k8s.io/cluster-api-provider-aws/pkg/coalescing"
)

type (
	// Options are controller options extended.
	Options struct {
		controller.Options
		Cache *coalescing.ReconcileCache
	}
)
