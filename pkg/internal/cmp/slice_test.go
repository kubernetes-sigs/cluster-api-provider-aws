/*
Copyright 2021 The Kubernetes Authors.

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

package cmp

import (
	"testing"

	. "github.com/onsi/gomega"
	"k8s.io/utils/pointer"
)

func TestCompareSlices(t *testing.T) {
	g := NewWithT(t)

	slice1 := []*string{pointer.String("foo"), pointer.String("bar")}
	slice2 := []*string{pointer.String("bar"), pointer.String("foo")}

	expected := Equals(slice1, slice2)
	g.Expect(expected, true)

	slice2 = append(slice2, pointer.String("test"))
	expected = Equals(slice1, slice2)
	g.Expect(expected, false)
}
