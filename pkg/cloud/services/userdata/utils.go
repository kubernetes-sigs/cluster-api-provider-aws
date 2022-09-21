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

package userdata

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

var defaultTemplateFuncMap = template.FuncMap{
	"Base64Encode": templateBase64Encode,
	"Indent":       templateYAMLIndent,
}

func templateBase64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func templateYAMLIndent(i int, input string) string {
	split := strings.Split(input, "\n")
	ident := "\n" + strings.Repeat(" ", i)
	return strings.Repeat(" ", i) + strings.Join(split, ident)
}

// GzipBytes will gzip a byte array.
func GzipBytes(dat []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(dat); err != nil {
		return []byte{}, errors.Wrap(err, "failed to gzip bytes")
	}

	if err := gz.Close(); err != nil {
		return []byte{}, errors.Wrap(err, "failed to gzip bytes")
	}

	return buf.Bytes(), nil
}

// ComputeHash returns the SHA256 hash of the user data byte array.
func ComputeHash(dat []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(dat))
}
