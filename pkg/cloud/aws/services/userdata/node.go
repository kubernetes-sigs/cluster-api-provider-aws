// Copyright © 2018 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package userdata

const (
	nodeBashScript = `{{.Header}}

HOSTNAME="$(curl http://169.254.169.254/latest/meta-data/local-hostname)"

cat >/tmp/kubeadm-node.yaml <<EOF
---
apiVersion: kubeadm.k8s.io/v1beta1
kind: JoinConfiguration
discovery:
  bootstrapToken:
    token: "{{.BootstrapToken}}"
    apiServerEndpoint: "{{.ELBAddress}}:6443"
    caCertHashes:
      - "{{.CACertHash}}"
nodeRegistration:
  name: "${HOSTNAME}"
  criSocket: /var/run/containerd/containerd.sock
  kubeletExtraArgs:
    cloud-provider: aws
EOF

kubeadm join --config /tmp/kubeadm-node.yaml
`
)

// NodeInput defines the context to generate a node user data.
type NodeInput struct {
	baseUserData

	CACertHash     string
	BootstrapToken string
	ELBAddress     string
}

// NewNode returns the user data string to be used on a node instance.
func NewNode(input *NodeInput) (string, error) {
	input.Header = defaultHeader
	return generate("node", nodeBashScript, input)
}
