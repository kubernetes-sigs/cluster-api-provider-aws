/*
Copyright (c) 2019 Red Hat, Inc.

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
	"regexp"
	"text/template"

	"github.com/onsi/gomega/types"

	. "github.com/onsi/gomega" // nolint
)

// EvaluateTemplate generates a string from the given templlate source and name value pairs. For
// example the following code:
//
//	EvaluateTemplate(
//		`{
//			"access_token": "{{ .AccessToken }}"
//			"refresh_token": "{{ .RefreshToken }}"
//		}`,
//		"AccessToken", "myaccesstoken",
//		"RefreshToken", "myrefreshtoken",
//	)
//
// Will generate the following string:
//
//	{
//		"access_token": "myaccesstoken"
//		"access_token": "myrefreshtoken"
//	}
//
// To simplify embeding of the templates in Go source the function also removes the leading tabs
// from the generated text.
func EvaluateTemplate(source string, args ...interface{}) string {
	// Check that there is an even number of args, and that the first of each pair is a string:
	count := len(args)
	Expect(count%2).To(
		Equal(0),
		"Template '%s' should have an even number of arguments, but it has %d",
		source, count,
	)
	for i := 0; i < count; i = i + 2 {
		name := args[i]
		_, ok := name.(string)
		Expect(ok).To(
			BeTrue(),
			"Argument %d of template '%s' is a key, so it should be a string, "+
				"but its type is %T",
			i, source, name,
		)
	}

	// Put the variables in the map that will be passed as the data object for the execution of
	// the template:
	data := make(map[string]interface{})
	for i := 0; i < count; i = i + 2 {
		name := args[i].(string)
		value := args[i+1]
		data[name] = value
	}

	// Parse the template:
	tmpl, err := template.New("").Parse(source)
	Expect(err).ToNot(
		HaveOccurred(),
		"Can't parse template '%s': %v",
		source, err,
	)

	// Execute the template:
	buffer := new(bytes.Buffer)
	err = tmpl.Execute(buffer, data)
	Expect(err).ToNot(
		HaveOccurred(),
		"Can't execute template '%s': %v",
		source, err,
	)
	result := buffer.String()

	// Remove the leading tabs:
	result = RemoveLeadingTabs(result)

	return result
}

// MatchJSONTemplate succeeds if actual is a string or stringer of JSON that matches the result of
// evaluating the given template with the given arguments.
func MatchJSONTemplate(template string, args ...interface{}) types.GomegaMatcher {
	return MatchJSON(EvaluateTemplate(template, args...))
}

// RemoveLeadingTabs removes the leading tabs from the lines of the given string.
func RemoveLeadingTabs(s string) string {
	return leadingTabsRE.ReplaceAllString(s, "")
}

// leadingTabsRE is the regular expression used to remove leading tabs from strings generated with
// the EvaluateTemplate function.
var leadingTabsRE = regexp.MustCompile(`(?m)^\t*`)
