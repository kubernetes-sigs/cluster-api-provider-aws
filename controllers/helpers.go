/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"strings"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clusterutil "sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

func PausedPredicates(logger logr.Logger) predicate.Funcs {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			return processIfUnpaused(logger.WithValues("predicate", "updateEvent"), e.ObjectNew, e.MetaNew)
		},
		CreateFunc: func(e event.CreateEvent) bool {
			return processIfUnpaused(logger.WithValues("predicate", "createEvent"), e.Object, e.Meta)
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return processIfUnpaused(logger.WithValues("predicate", "deleteEvent"), e.Object, e.Meta)
		},
		GenericFunc: func(e event.GenericEvent) bool {
			return processIfUnpaused(logger.WithValues("predicate", "genericEvent"), e.Object, e.Meta)
		},
	}
}

func processIfUnpaused(logger logr.Logger, obj runtime.Object, meta metav1.Object) bool {
	kind := strings.ToLower(obj.GetObjectKind().GroupVersionKind().Kind)
	log := logger.WithValues("namespace", meta.GetNamespace(), kind, meta.GetName())
	if clusterutil.HasPausedAnnotation(meta) {
		log.V(4).Info("Resource is paused, will not attempt to map resource")
		return false
	}
	log.V(4).Info("Resource is not paused, will attempt to map resource")
	return true
}
