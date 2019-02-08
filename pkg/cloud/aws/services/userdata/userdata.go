/*
Copyright 2018 The Kubernetes Authors.

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
	"text/template"

	"github.com/pkg/errors"
)

const (
	cloudConfigHeader = `## template: jinja
#cloud-config
`
	defaultHeader = `#!/usr/bin/env bash

# Copyright 2018 by the contributors
#
#    Licensed under the Apache License, Version 2.0 (the "License");
#    you may not use this file except in compliance with the License.
#    You may obtain a copy of the License at
#
#        http://www.apache.org/licenses/LICENSE-2.0
#
#    Unless required by applicable law or agreed to in writing, software
#    distributed under the License is distributed on an "AS IS" BASIS,
#    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#    See the License for the specific language governing permissions and
#    limitations under the License.

set -o verbose
set -o errexit
set -o nounset
set -o pipefail
`
)

type baseUserData struct {
	Header string
}

func generate(kind string, tpl string, data interface{}) (string, error) {
	t, err := template.New(kind).Parse(tpl)
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse %s template", kind)
	}

	var out bytes.Buffer
	if err := t.Execute(&out, data); err != nil {
		return "", errors.Wrapf(err, "failed to generate %s template", kind)
	}

	return out.String(), nil
}

func generateWithFuncs(kind string, tpl string, funcsMap template.FuncMap, data interface{}) (string, error) {
	t, err := template.New(kind).Funcs(funcsMap).Parse(tpl)
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse %s template", kind)
	}

	var out bytes.Buffer
	if err := t.Execute(&out, data); err != nil {
		return "", errors.Wrapf(err, "failed to generate %s template", kind)
	}

	return out.String(), nil
}

func funcMap(funcs map[string]interface{}) template.FuncMap {
	funcMap := template.FuncMap{}
	for name, function := range funcs {
		funcMap[name] = function
	}

	return funcMap
}
