/*
Copyright 2019 The Kubernetes Authors.

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

package registry

import (
	"sigs.k8s.io/promo-tools/v4/internal/legacy/container"
	"sigs.k8s.io/promo-tools/v4/types/image"
)

// Various set manipulation operations. Some set operations are missing,
// because, we don't use them.

// ToSet converts a RegInvImage to a Set.
func (a RegInvImage) ToSet() container.Set {
	b := make(container.Set)
	for k, v := range a {
		b[k] = v
	}

	return b
}

func toRegistryInventory(a container.Set) RegInvImage {
	b := make(RegInvImage)
	for k, v := range a {
		// TODO: Why are we not checking errors here?
		//nolint:errcheck
		b[k.(image.Name)] = v.(DigestTags)
	}

	return b
}

// Minus is a set operation.
// TODO: ST1016: methods on the same type should have the same receiver name
func (a RegInvImage) Minus(b RegInvImage) RegInvImage {
	aSet := a.ToSet()
	bSet := b.ToSet()
	cSet := aSet.Minus(bSet)

	return toRegistryInventory(cSet)
}

// Union is a set operation.
func (a RegInvImage) Union(b RegInvImage) RegInvImage {
	aSet := a.ToSet()
	bSet := b.ToSet()
	cSet := aSet.Union(bSet)

	return toRegistryInventory(cSet)
}

// ToTagSet converts a TagSlice to a TagSet.
func (a TagSlice) ToTagSet() TagSet {
	b := make(TagSet)
	for _, t := range a {
		// The value doesn't matter.
		b[t] = nil
	}

	return b
}

// Minus is a set operation.
func (a TagSlice) Minus(b TagSlice) TagSet {
	aSet := a.ToTagSet()
	bSet := b.ToTagSet()
	cSet := aSet.Minus(bSet)

	return cSet
}

// Union is a set operation.
func (a TagSlice) Union(b TagSlice) TagSet {
	aSet := a.ToTagSet()
	bSet := b.ToTagSet()
	cSet := aSet.Union(bSet)

	return cSet
}

// Intersection is a set operation.
func (a TagSlice) Intersection(b TagSlice) TagSet {
	aSet := a.ToTagSet()
	bSet := b.ToTagSet()
	cSet := aSet.Intersection(bSet)

	return cSet
}

// ToSet converts a TagSet to a Set.
func (a TagSet) ToSet() container.Set {
	b := make(container.Set)
	for t := range a {
		// The value doesn't matter.
		b[t] = nil
	}

	return b
}

func setToTagSet(a container.Set) TagSet {
	b := make(TagSet)
	for k := range a {
		b[k.(image.Tag)] = nil
	}

	return b
}

// Minus is a set operation.
func (a TagSet) Minus(b TagSet) TagSet {
	aSet := a.ToSet()
	bSet := b.ToSet()
	cSet := aSet.Minus(bSet)

	return setToTagSet(cSet)
}

// Union is a set operation.
func (a TagSet) Union(b TagSet) TagSet {
	aSet := a.ToSet()
	bSet := b.ToSet()
	cSet := aSet.Union(bSet)

	return setToTagSet(cSet)
}

// Intersection is a set operation.
func (a TagSet) Intersection(b TagSet) TagSet {
	aSet := a.ToSet()
	bSet := b.ToSet()
	cSet := aSet.Intersection(bSet)

	return setToTagSet(cSet)
}
