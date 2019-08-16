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

package conversion

import (
	"fmt"
	"testing"

	capav1a2 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha2"
	capav1a1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
)

type asserter struct {
	*testing.T
}

func (a *asserter) stringEqual(expected, actual, name string) {
	if expected != actual {
		a.Errorf("expected %s to be %q, got %q", name, expected, actual)
	}
}

func (a *asserter) stringPtrEqual(expected, actual *string, name string) {
	if expected == nil {
		if actual != nil {
			a.Errorf("expected %s to be nil, but was %q", name, *actual)
		}
	} else if actual == nil {
		a.Errorf("expected %s to be %q, but was nil", name, *expected)
	} else {
		a.stringEqual(*expected, *actual, name)
	}
}

func (a *asserter) notEmpty(expected, name string) {
	if expected == "" {
		a.Errorf("expected %s to not be empty", name)
	}
}

func (a *asserter) stringArrayEqual(expected, actual []string, name string) {
	if len(expected) != len(actual) {
		a.Errorf("expected %s to have length %v, got %v", name, len(expected), len(actual))
	} else {
		for i, e := range expected {
			a.stringEqual(e, actual[i], fmt.Sprintf("%s[%d]", name, i))
		}
	}
}

func (a *asserter) awsRefEqual(expected *capav1a1.AWSResourceReference, actual *capav1a2.AWSResourceReference, name string) {
	a.stringPtrEqual(expected.ID, actual.ID, fmt.Sprintf("%s ID", name))

	a.stringPtrEqual(expected.ARN, actual.ARN, fmt.Sprintf("%s ID", name))

	if len(expected.Filters) != len(actual.Filters) {
		a.Errorf("%s filters should have length %d, but had length %d", name, len(expected.Filters), len(actual.Filters))
	} else {
		for i, exFilter := range expected.Filters {
			actFilter := actual.Filters[i]
			a.stringArrayEqual(exFilter.Values, actFilter.Values, fmt.Sprintf("%s filter[%d]", name, i))
		}
	}
}
