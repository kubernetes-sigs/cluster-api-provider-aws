/*
Copyright 2020 The Kubernetes Authors.

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

package mime

import (
	"bytes"
	"net/mail"
	"testing"
)

func TestGenerateInitDocument(t *testing.T) {
	secretARN := "secretARN"
	doc, _ := GenerateInitDocument(secretARN, 1, "eu-west-1", "localhost", "abc123")

	_, err := mail.ReadMessage(bytes.NewBuffer(doc))
	if err != nil {
		t.Fatalf("Cannot parse MIME doc: %+v\n%s", err, string(doc))
	}
}
