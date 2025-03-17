/*
Copyright (c) 2020 Red Hat, Inc.

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

// IMPORTANT: This file has been generated automatically, refrain from modifying it manually as all
// your changes will be lost when the file is generated again.

package v1 // github.com/openshift-online/ocm-sdk-go/addonsmgmt/v1

// AddonSecretPropagation represents the values of the 'addon_secret_propagation' type.
//
// Representation of an addon secret propagation
type AddonSecretPropagation struct {
	bitmap_           uint32
	id                string
	destinationSecret string
	sourceSecret      string
	enabled           bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AddonSecretPropagation) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// ID returns the value of the 'ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// ID of the secret propagation
func (o *AddonSecretPropagation) ID() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the value of the 'ID' attribute and
// a flag indicating if the attribute has a value.
//
// ID of the secret propagation
func (o *AddonSecretPropagation) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.id
	}
	return
}

// DestinationSecret returns the value of the 'destination_secret' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// DestinationSecret is location of the secret to be added
func (o *AddonSecretPropagation) DestinationSecret() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.destinationSecret
	}
	return ""
}

// GetDestinationSecret returns the value of the 'destination_secret' attribute and
// a flag indicating if the attribute has a value.
//
// DestinationSecret is location of the secret to be added
func (o *AddonSecretPropagation) GetDestinationSecret() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.destinationSecret
	}
	return
}

// Enabled returns the value of the 'enabled' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates is this secret propagation is enabled for the addon
func (o *AddonSecretPropagation) Enabled() bool {
	if o != nil && o.bitmap_&4 != 0 {
		return o.enabled
	}
	return false
}

// GetEnabled returns the value of the 'enabled' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates is this secret propagation is enabled for the addon
func (o *AddonSecretPropagation) GetEnabled() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.enabled
	}
	return
}

// SourceSecret returns the value of the 'source_secret' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// SourceSecret is location of the source secret
func (o *AddonSecretPropagation) SourceSecret() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.sourceSecret
	}
	return ""
}

// GetSourceSecret returns the value of the 'source_secret' attribute and
// a flag indicating if the attribute has a value.
//
// SourceSecret is location of the source secret
func (o *AddonSecretPropagation) GetSourceSecret() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.sourceSecret
	}
	return
}

// AddonSecretPropagationListKind is the name of the type used to represent list of objects of
// type 'addon_secret_propagation'.
const AddonSecretPropagationListKind = "AddonSecretPropagationList"

// AddonSecretPropagationListLinkKind is the name of the type used to represent links to list
// of objects of type 'addon_secret_propagation'.
const AddonSecretPropagationListLinkKind = "AddonSecretPropagationListLink"

// AddonSecretPropagationNilKind is the name of the type used to nil lists of objects of
// type 'addon_secret_propagation'.
const AddonSecretPropagationListNilKind = "AddonSecretPropagationListNil"

// AddonSecretPropagationList is a list of values of the 'addon_secret_propagation' type.
type AddonSecretPropagationList struct {
	href  string
	link  bool
	items []*AddonSecretPropagation
}

// Len returns the length of the list.
func (l *AddonSecretPropagationList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *AddonSecretPropagationList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AddonSecretPropagationList) Get(i int) *AddonSecretPropagation {
	if l == nil || i < 0 || i >= len(l.items) {
		return nil
	}
	return l.items[i]
}

// Slice returns an slice containing the items of the list. The returned slice is a
// copy of the one used internally, so it can be modified without affecting the
// internal representation.
//
// If you don't need to modify the returned slice consider using the Each or Range
// functions, as they don't need to allocate a new slice.
func (l *AddonSecretPropagationList) Slice() []*AddonSecretPropagation {
	var slice []*AddonSecretPropagation
	if l == nil {
		slice = make([]*AddonSecretPropagation, 0)
	} else {
		slice = make([]*AddonSecretPropagation, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AddonSecretPropagationList) Each(f func(item *AddonSecretPropagation) bool) {
	if l == nil {
		return
	}
	for _, item := range l.items {
		if !f(item) {
			break
		}
	}
}

// Range runs the given function for each index and item of the list, in order. If
// the function returns false the iteration stops, otherwise it continues till all
// the elements of the list have been processed.
func (l *AddonSecretPropagationList) Range(f func(index int, item *AddonSecretPropagation) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
