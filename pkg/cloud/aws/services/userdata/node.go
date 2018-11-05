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

certificate=$(echo '{{.CACert}}' | base64 -w0)

cat >/tmp/cluster-info.yaml <<EOF
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: ${certificate}
    server: https://{{.ELBAddress}}:6443
  name: ""
contexts: []
current-context: ""
kind: Config
preferences: {}
users: []
EOF

HOSTNAME="$(hostname -f 2>/dev/null || curl http://169.254.169.254/latest/meta-data/local-hostname)"

cat >/tmp/kubeadm-node.yaml <<EOF
---
apiVersion: kubeadm.k8s.io/v1alpha3
kind: NodeConfiguration
token: {{.BootstrapToken}}
discoveryTokenAPIServers: 
- "{{.ELBAddress}}:6443"
discoveryFile: /tmp/cluster-info.yaml
nodeRegistration:
  name: ${HOSTNAME}
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

	CACert         string
	BootstrapToken string
	ELBAddress     string
}

// NewNode returns the user data string to be used on a node instance.
func NewNode(input *NodeInput) (string, error) {
	input.Header = defaultHeader
	return generate("node", nodeBashScript, input)
}
