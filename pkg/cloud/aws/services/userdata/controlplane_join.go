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
	controlPlaneJoinCloudInit = `{{.Header}}
{{template "files" .WriteFiles}}
-   path: /tmp/kubeadm-controlplane-join-config.yaml
    owner: root:root
    permissions: '0640'
    content: |
{{.JoinConfiguration | Indent 6}}
kubeadm:
  operation: join
  config: /tmp/kubeadm-controlplane-join-config.yaml
`
)

// ControlPlaneJoinInput defines context to generate controlplane instance user data for controlplane node join.
type ControlPlaneJoinInput struct {
	baseUserData
	Certificates

	BootstrapToken    string
	ELBAddress        string
	JoinConfiguration string
}

// NewJoinControlPlane returns the user data string to be used on a new contrplplane instance.
func NewJoinControlPlane(input *ControlPlaneJoinInput) (string, error) {
	input.Header = cloudConfigHeader
	if err := input.Certificates.validate(); err != nil {
		return "", errors.Wrapf(err, "ControlPlaneInput is invalid")
	}

	input.WriteFiles = certificatesToFiles(input.Certificates)
	userData, err := generate("JoinControlplane", controlPlaneJoinCloudInit, input)
	if err != nil {
		return "", errors.Wrapf(err, "failed to generate user data for machine joining control plane")
	}

	return userData, err
}
