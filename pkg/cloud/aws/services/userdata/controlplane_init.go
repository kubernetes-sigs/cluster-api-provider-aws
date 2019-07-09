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
	"github.com/pkg/errors"
)

const (
	controlPlaneCloudInit = `{{.Header}}
{{template "files" .WriteFiles}}
-   path: /tmp/kubeadm.yaml
    owner: root:root
    permissions: '0640'
    content: |
      ---
{{.ClusterConfiguration | Indent 6}}
      ---
{{.InitConfiguration | Indent 6}}
kubeadm:
  operation: init
  config: /tmp/kubeadm.yaml
`
)

// ControlPlaneInput defines the context to generate a controlplane instance user data.
type ControlPlaneInput struct {
	baseUserData
	Certificates

	ClusterConfiguration string
	InitConfiguration    string
}

// NewInitControlPlane returns the user data string to be used on a controlplane instance.
func NewInitControlPlane(input *ControlPlaneInput) (string, error) {
	input.Header = cloudConfigHeader
	if err := input.Certificates.validate(); err != nil {
		return "", errors.Wrapf(err, "ControlPlaneInput is invalid")
	}

	input.WriteFiles = certificatesToFiles(input.Certificates)
	userData, err := generate("InitControlplane", controlPlaneCloudInit, input)
	if err != nil {
		return "", errors.Wrapf(err, "failed to generate user data for new control plane machine")
	}

	return userData, err
}
