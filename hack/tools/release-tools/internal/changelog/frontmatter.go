/*
Copyright 2026 The Kubernetes Authors.

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

package changelog

import (
	"fmt"
	"regexp"
	"strings"
)

var frontMatterRE = regexp.MustCompile(`(?s)\A---\r?\n(.*?)\r?\n---\r?\n(.*)\z`)
var contractFieldRE = regexp.MustCompile(`(?m)^contract:\s*(\S+)\s*$`)

// WithFrontMatter prepends a "contract: <value>" front-matter block to body.
func WithFrontMatter(contract string, body []byte) []byte {
	return []byte(fmt.Sprintf("---\ncontract: %s\n---\n%s", contract, string(body)))
}

// ParseFrontMatter splits a CHANGELOG file into its "contract" front-matter
// value and the remaining body. If content has no front-matter block, it
// returns an empty contract and the content unchanged.
func ParseFrontMatter(content []byte) (contract string, body []byte) {
	m := frontMatterRE.FindSubmatch(content)
	if m == nil {
		return "", content
	}
	fm := m[1]
	body = m[2]
	if cm := contractFieldRE.FindSubmatch(fm); cm != nil {
		contract = strings.TrimSpace(string(cm[1]))
	}
	return contract, body
}
