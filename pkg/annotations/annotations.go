/*
Copyright 2022 The Kubernetes Authors.

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

package annotations

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// ExternalResourceGCAnnotation is the name of an annotation that indicates if
	// external resources should be garbage collected for the cluster.
	ExternalResourceGCAnnotation = "aws.cluster.x-k8s.io/external-resource-gc"
)

// Set will set the value of an annotation on the supplied object. If there is no annotation it will be created.
func Set(obj metav1.Object, name, value string) {
	annotations := obj.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations[name] = value
	obj.SetAnnotations(annotations)
}

// Get will get the value of the supplied annotation.
func Get(obj metav1.Object, name string) (value string, found bool) {
	annotations := obj.GetAnnotations()
	if len(annotations) == 0 {
		return "", false
	}

	value, found = annotations[name]

	return
}

// Has returns true if the supplied object has the supplied annotation.
func Has(obj metav1.Object, name string) bool {
	annotations := obj.GetAnnotations()
	if len(annotations) == 0 {
		return false
	}

	_, found := annotations[name]

	return found
}

// // SetExternalResourceGC is a helper that will set the external resource garbage collection
// // annotation on the supplied object.
// func SetExternalResourceGC(obj metav1.Object, hasBeenGC bool) {
// 	Set(obj, ExternalResourceGCAnnotation, strconv.FormatBool(hasBeenGC))
// }

// // GetExternalResourceGC is a helper that will get the value of the external resource garbage collection
// // annotation from the supplied object. If the annotation is not found then its returns true (i.e. that
// // the external resource have been garbage collected).
// func GetExternalResourceGC(obj metav1.Object) (bool, error) {
// 	status, found := Get(obj, ExternalResourceGCAnnotation)
// 	if !found {
// 		// If not found then assume GC
// 		return true, nil
// 	}

// 	converted, err := strconv.ParseBool(status)
// 	if err != nil {
// 		return false, fmt.Errorf("converting annotation value to bool %s: %w", ExternalResourceGCAnnotation, err)
// 	}

// 	return converted, nil
// }
