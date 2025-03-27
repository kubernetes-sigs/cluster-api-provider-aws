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

package v1 // github.com/openshift-online/ocm-sdk-go/authorizations/v1

// SelfFeatureReviewRequestBuilder contains the data and logic needed to build 'self_feature_review_request' objects.
//
// Representation of a feature review performed against oneself
type SelfFeatureReviewRequestBuilder struct {
	bitmap_ uint32
	feature string
}

// NewSelfFeatureReviewRequest creates a new builder of 'self_feature_review_request' objects.
func NewSelfFeatureReviewRequest() *SelfFeatureReviewRequestBuilder {
	return &SelfFeatureReviewRequestBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *SelfFeatureReviewRequestBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Feature sets the value of the 'feature' attribute to the given value.
func (b *SelfFeatureReviewRequestBuilder) Feature(value string) *SelfFeatureReviewRequestBuilder {
	b.feature = value
	b.bitmap_ |= 1
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *SelfFeatureReviewRequestBuilder) Copy(object *SelfFeatureReviewRequest) *SelfFeatureReviewRequestBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.feature = object.feature
	return b
}

// Build creates a 'self_feature_review_request' object using the configuration stored in the builder.
func (b *SelfFeatureReviewRequestBuilder) Build() (object *SelfFeatureReviewRequest, err error) {
	object = new(SelfFeatureReviewRequest)
	object.bitmap_ = b.bitmap_
	object.feature = b.feature
	return
}
