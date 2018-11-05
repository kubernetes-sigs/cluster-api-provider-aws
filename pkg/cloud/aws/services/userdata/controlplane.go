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
	controlPlaneBashScript = `{{.Header}}

mkdir -p /etc/kubernetes/pki

echo '{{.CACert}}' > /etc/kubernetes/pki/ca.crt
echo '{{.CAKey}}' > /etc/kubernetes/pki/ca.key

PRIVATE_IP=$(curl http://169.254.169.254/latest/meta-data/local-ipv4)
HOSTNAME="$(hostname -f 2>/dev/null || curl http://169.254.169.254/latest/meta-data/local-hostname)"

cat >/tmp/kubeadm.yaml <<EOF
---
apiVersion: kubeadm.k8s.io/v1alpha3
kind: ClusterConfiguration
apiServerCertSANs:
  - "$PRIVATE_IP"
  - "{{.ELBAddress}}"
apiServerExtraArgs:
  cloud-provider: aws
controlPlaneEndpoint: "{{.ELBAddress}}:6443"
clusterName: "{{.ClusterName}}"
networking:
  dnsDomain: "{{.ServiceDomain}}"
  podSubnet: "{{.PodSubnet}}"
  serviceSubnet: "{{.ServiceSubnet}}"
kubernetesVersion: "{{.KubernetesVersion}}"
---
apiVersion: kubeadm.k8s.io/v1alpha3
kind: InitConfiguration
nodeRegistration:
  name: ${HOSTNAME}
  criSocket: /var/run/containerd/containerd.sock
  kubeletExtraArgs:
    cloud-provider: aws
    allocate-node-cidrs: "false"
EOF

kubeadm init --config /tmp/kubeadm.yaml
`
)

// ControlPlaneInput defines the context to generate a controlplane instance user data.
type ControlPlaneInput struct {
	baseUserData

	CACert            string
	CAKey             string
	ELBAddress        string
	ClusterName       string
	PodSubnet         string
	ServiceDomain     string
	ServiceSubnet     string
	KubernetesVersion string
}

// NewControlPlane returns the user data string to be used on a controlplane instance.
func NewControlPlane(input *ControlPlaneInput) (string, error) {
	input.Header = defaultHeader
	return generate("controlplane", controlPlaneBashScript, input)
}
