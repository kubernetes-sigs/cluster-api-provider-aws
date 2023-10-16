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

package stream

import (
	"bytes"
	"io"
	"strings"
)

// Fake is a predefined stream (set with Bytes).
type Fake struct {
	stream io.Reader
	Bytes  []byte
}

// Produce a fake stream on stdout (and an empty stderr). Unlike a real stream,
// this does not call a subprocess --- instead it just provides a predefined
// stream (Bytes) to create an io.Reader for stdout. The stderr stream is empty.
func (producer *Fake) Produce() (stream, blankStderr io.Reader, err error) {
	producer.stream = bytes.NewReader(producer.Bytes)
	blankStderr = strings.NewReader("")
	return producer.stream, blankStderr, nil
}

// Close does nothing, as there is no actual subprocess to close.
func (producer *Fake) Close() error {
	return nil
}
