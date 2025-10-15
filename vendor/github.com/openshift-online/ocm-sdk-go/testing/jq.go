/*
Copyright (c) 2021 Red Hat, Inc.

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

package testing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/itchyny/gojq"
	"github.com/onsi/gomega/types"

	. "github.com/onsi/gomega" //nolint
)

// JQ runs the given `jq` filter on the given object and returns the list of results. The returned
// slice will never be nil; if there are no results it will be empty.
func JQ(filter string, input interface{}) (results []interface{}, err error) {
	query, err := gojq.Parse(filter)
	if err != nil {
		return
	}
	iterator := query.Run(input)
	for {
		result, ok := iterator.Next()
		if !ok {
			break
		}
		results = append(results, result)
	}
	return
}

// MatchJQ creates a matcher that checks that the all the results of applying a `jq` filter to the
// actual value is the given expected value.
func MatchJQ(filter string, expected interface{}) types.GomegaMatcher {
	return &jqMatcher{
		filter:   filter,
		expected: expected,
	}
}

type jqMatcher struct {
	filter   string
	expected interface{}
	results  []interface{}
}

func (m *jqMatcher) Match(actual interface{}) (success bool, err error) {
	// Run the query:
	m.results, err = JQ(m.filter, actual)
	if err != nil {
		return
	}

	// Check that there is at least one result:
	if len(m.results) == 0 {
		return
	}

	// We consider the match sucessful if all the results returned by the JQ filter are exactly
	// equal to the expected value.
	success = true
	for _, result := range m.results {
		if !reflect.DeepEqual(result, m.expected) {
			success = false
			break
		}
	}
	return
}

func (m *jqMatcher) FailureMessage(actual interface{}) string {
	return fmt.Sprintf(
		"Expected all results of running JQ filter\n\t%s\n"+
			"on input\n\t%s\n"+
			"to be\n\t%s\n"+
			"but at list one of the following results isn't\n\t%s\n",
		m.filter, m.pretty(actual), m.pretty(m.expected), m.pretty(m.results),
	)
}

func (m *jqMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf(
		"Expected results of running JQ filter\n\t%s\n"+
			"on\n\t%s\n"+
			"to not be\n\t%s\n",
		m.filter, m.pretty(actual), m.pretty(m.expected),
	)
}

func (m *jqMatcher) pretty(object interface{}) string {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetIndent("\t", "  ")
	err := encoder.Encode(object)
	if err != nil {
		return fmt.Sprintf("\t%v", object)
	}
	return strings.TrimRight(buffer.String(), "\n")
}

// VerifyJQ verifies that the result of applying the given `jq` filter to the request body matches
// the given expected value.
func VerifyJQ(filter string, expected interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Read the body completely:
		body, err := io.ReadAll(r.Body)
		Expect(err).ToNot(HaveOccurred())
		err = r.Body.Close()
		Expect(err).ToNot(HaveOccurred())

		// Replace the body with a buffer so that other calls to this same method, or to
		// other verification methods will also be able to work with it.
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		// Parse the body as JSON and verify that it matches the filter:
		var data interface{}
		err = json.Unmarshal(body, &data)
		Expect(err).ToNot(HaveOccurred())
		Expect(data).To(MatchJQ(filter, expected))
	}
}
