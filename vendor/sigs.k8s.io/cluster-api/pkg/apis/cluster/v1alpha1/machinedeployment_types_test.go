/*
Copyright 2018 The Kubernetes Authors.

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

package v1alpha1

import (
	"reflect"
	"testing"

	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func TestStorageMachineDeployment(t *testing.T) {
	key := types.NamespacedName{Name: "foo", Namespace: "default"}
	created := &MachineDeployment{ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "default"}}

	// Test Create
	fetched := &MachineDeployment{}
	if err := c.Create(context.TODO(), created); err != nil {
		t.Errorf("error creating machine deployment: %v", err)
	}

	if err := c.Get(context.TODO(), key, fetched); err != nil {
		t.Errorf("error getting machine deployment: %v", err)
	}
	if !reflect.DeepEqual(*fetched, *created) {
		t.Error("fetched value not what was created")
	}

	// Test Updating the Labels
	updated := fetched.DeepCopy()
	updated.Labels = map[string]string{"hello": "world"}
	if err := c.Update(context.TODO(), updated); err != nil {
		t.Errorf("error updating machine deployment: %v", err)
	}

	if err := c.Get(context.TODO(), key, fetched); err != nil {
		t.Errorf("error getting machine deployment: %v", err)
	}
	if !reflect.DeepEqual(*fetched, *updated) {
		t.Error("fetched value not what was updated")
	}

	// Test Delete
	if err := c.Delete(context.TODO(), fetched); err != nil {
		t.Errorf("error deleting machine deployment: %v", err)
	}
	if err := c.Get(context.TODO(), key, fetched); err == nil {
		t.Error("expected error getting machine deployment")
	}
}
