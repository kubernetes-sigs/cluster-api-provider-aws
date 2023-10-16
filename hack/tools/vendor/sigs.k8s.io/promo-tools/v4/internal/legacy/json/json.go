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

package json

import (
	"io"

	yaml "gopkg.in/yaml.v2"
)

// Object is a JSON object.
type Object map[string]interface{}

// Objects is a slice of Object values.
type Objects []Object

// Consume decodes JSON from a given io.Reader handle.
func Consume(h io.Reader) (Objects, error) {
	// Generic-looking type to hold whatever JSON objects we get as a stream
	// (that's why it's a slice, not just a plain map).
	var m Objects
	decoder := yaml.NewDecoder(h)
	for {
		err := decoder.Decode(&m)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}
	return m, nil
}
