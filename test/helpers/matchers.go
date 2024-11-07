/*
Copyright 2024 The Kubernetes Authors.

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

package helpers

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/golang/mock/gomock"
)

// PartialMatchCreateTargetGroupInput matches a partial CreateTargetGroupInput struct based on fuzzy matching rules.
func PartialMatchCreateTargetGroupInput(t *testing.T, i *elbv2.CreateTargetGroupInput) gomock.Matcher {
	t.Helper()
	return &createTargetGroupInputPartialMatcher{
		in: i,
		t:  t,
	}
}

// createTargetGroupInputPartialMatcher conforms to the gomock.Matcher interface in order to implement a match against a partial
// CreateTargetGroupInput expected value.
// In particular, the TargetGroupName expected value is used as a prefix, in order to support generated names.
type createTargetGroupInputPartialMatcher struct {
	in *elbv2.CreateTargetGroupInput
	t  *testing.T
}

func (m *createTargetGroupInputPartialMatcher) Matches(x interface{}) bool {
	actual, ok := x.(*elbv2.CreateTargetGroupInput)
	if !ok {
		return false
	}

	// Check for a perfect match across all fields first.
	eq := gomock.Eq(m.in).Matches(actual)

	if !eq && (actual.Name != nil && m.in.Name != nil) {
		// If the actual name is prefixed with the expected value, then it matches
		if (*actual.Name != *m.in.Name) && strings.HasPrefix(*actual.Name, *m.in.Name) {
			return true
		}
	}

	return eq
}

func (m *createTargetGroupInputPartialMatcher) String() string {
	return fmt.Sprintf("%v (%T)", m.in, m.in)
}
