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

import "github.com/pkg/errors"

const (
	controlPlaneBashScript = `{{.Header}}

set -eox

mkdir -p /etc/kubernetes/pki/etcd

echo '{{.CACert}}' > /etc/kubernetes/pki/ca.crt
echo '{{.CAKey}}' > /etc/kubernetes/pki/ca.key

echo '{{.EtcdCACert}}' > /etc/kubernetes/pki/etcd/ca.crt
echo '{{.EtcdCAKey}}' >/etc/kubernetes/pki/etcd/ca.key

echo '{{.FrontProxyCACert}}' > /etc/kubernetes/pki/front-proxy-ca.crt
echo '{{.FrontProxyCAKey}}' > /etc/kubernetes/pki/front-proxy-ca.key

echo '{{.SaCert}}' > /etc/kubernetes/pki/sa.pub
echo '{{.SaKey}}' > /etc/kubernetes/pki/sa.key

PRIVATE_IP=$(curl http://169.254.169.254/latest/meta-data/local-ipv4)
HOSTNAME="$(curl http://169.254.169.254/latest/meta-data/local-hostname)"

cat >/tmp/kubeadm.yaml <<EOF
---
apiVersion: kubeadm.k8s.io/v1beta1
kind: ClusterConfiguration
apiServer:
  certSANs:
    - "$PRIVATE_IP"
    - "{{.ELBAddress}}"
  extraArgs:
    cloud-provider: aws
controlPlaneEndpoint: "{{.ELBAddress}}:6443"
clusterName: "{{.ClusterName}}"
networking:
  dnsDomain: "{{.ServiceDomain}}"
  podSubnet: "{{.PodSubnet}}"
  serviceSubnet: "{{.ServiceSubnet}}"
kubernetesVersion: "{{.KubernetesVersion}}"
---
apiVersion: kubeadm.k8s.io/v1beta1
kind: InitConfiguration
nodeRegistration:
  name: ${HOSTNAME}
  criSocket: /var/run/containerd/containerd.sock
  kubeletExtraArgs:
    cloud-provider: aws
EOF

kubeadm init --config /tmp/kubeadm.yaml --v 10
`

	controlPlaneJoinBashScript = `{{.Header}}
    
set -eox

mkdir -p /etc/kubernetes/pki

echo '{{.CACert}}' > /etc/kubernetes/pki/ca.crt
echo '{{.CAKey}}' > /etc/kubernetes/pki/ca.key

echo '{{.EtcdCACert}}' > /etc/kubernetes/pki/etcd/ca.crt
echo '{{.EtcdCAKey}}' >/etc/kubernetes/pki/etcd/ca.key

echo '{{.FrontProxyCACert}}' > /etc/kubernetes/pki/front-proxy-ca.crt
echo '{{.FrontProxyCAKey}}' > /etc/kubernetes/pki/front-proxy-ca.key

echo '{{.SaCert}}' > /etc/kubernetes/pki/sa.pub
echo '{{.SaKey}}' > /etc/kubernetes/pki/sa.key

PRIVATE_IP=$(curl http://169.254.169.254/latest/meta-data/local-ipv4)
HOSTNAME="$(curl http://169.254.169.254/latest/meta-data/local-hostname)"

cat >/tmp/kubeadm-controlplane-join-config.yaml <<EOF
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
controlPlane:
  localAPIEndpoint:
    advertiseAddress: "${PRIVATE_IP}"
    bindPort: 6443
EOF

kubeadm join --config /tmp/kubeadm-controlplane-join-config.yaml --v 10
`
)

func isKeyPairValid(cert, key string) bool {
	if (cert == "" && key != "") ||
		(cert != "" && key == "") {
		return false
	}
	return true
}

// ControlPlaneInput defines the context to generate a controlplane instance user data.
type ControlPlaneInput struct {
	baseUserData

	CACert            string
	CAKey             string
	EtcdCACert        string
	EtcdCAKey         string
	FrontProxyCACert  string
	FrontProxyCAKey   string
	SaCert            string
	SaKey             string
	ELBAddress        string
	ClusterName       string
	PodSubnet         string
	ServiceDomain     string
	ServiceSubnet     string
	KubernetesVersion string
}

// ContolPlaneJoinInput defines context to generate controlplane instance user data for controlplane node join.
type ContolPlaneJoinInput struct {
	baseUserData

	CACertHash       string
	CACert           string
	CAKey            string
	EtcdCACert       string
	EtcdCAKey        string
	FrontProxyCACert string
	FrontProxyCAKey  string
	SaCert           string
	SaKey            string
	BootstrapToken   string
	ELBAddress       string
}

func (cpi *ControlPlaneInput) isValid() (bool, string) {
	// TODO: ashish-amarnath verify if ca cert and key is mandatory for kubeadm init
	if cpi.CACert == "" || cpi.CAKey == "" {
		return false, "CA cert material in the ControlPlaneInput is invalid"
	}

	if !isKeyPairValid(cpi.EtcdCACert, cpi.EtcdCAKey) {
		return false, "ETCD cert material in the ControlPlaneInput is invalid"
	}

	if !isKeyPairValid(cpi.FrontProxyCACert, cpi.FrontProxyCAKey) {
		return false, "FrontProxy cert material in ControlPlaneInput is invalid"
	}

	if !isKeyPairValid(cpi.SaCert, cpi.SaKey) {
		return false, "ServiceAccount cert material in ControlPlaneInput is invalid"
	}

	return true, ""
}

func (cpi *ContolPlaneJoinInput) isValid() (bool, string) {
	if !isKeyPairValid(cpi.CACert, cpi.CAKey) {
		return false, "CA cert material in the ContolPlaneJoinInput is invalid"
	}

	if !isKeyPairValid(cpi.EtcdCACert, cpi.EtcdCAKey) {
		return false, "ETCD cert material in the ContolPlaneJoinInput is invalid"
	}

	if !isKeyPairValid(cpi.FrontProxyCACert, cpi.FrontProxyCAKey) {
		return false, "FrontProxy cert material in ContolPlaneJoinInput is invalid"
	}

	if !isKeyPairValid(cpi.SaCert, cpi.SaKey) {
		return false, "ServiceAccount cert material in ContolPlaneJoinInput is invalid"
	}

	return true, ""
}

// NewControlPlane returns the user data string to be used on a controlplane instance.
func NewControlPlane(input *ControlPlaneInput) (string, error) {
	input.Header = defaultHeader
	valid, reason := input.isValid()
	if !valid {
		return "", errors.Errorf("ControlPlaneInput is invalid, Reason: %q", reason)
	}

	userData, err := generate("controlplane", controlPlaneBashScript, input)
	if err != nil {
		return "", errors.Wrapf(err, "failed to generate user data for new control plane machine")
	}

	return userData, err
}

// JoinControlPlane returns the user data string to be used on a new contrplplane instance.
func JoinControlPlane(input *ContolPlaneJoinInput) (string, error) {
	input.Header = defaultHeader

	valid, reason := input.isValid()
	if !valid {
		return "", errors.Errorf("ControlPlaneInput is invalid, Reason: %q", reason)
	}

	userData, err := generate("controlplane", controlPlaneJoinBashScript, input)
	if err != nil {
		return "", errors.Wrapf(err, "failed to generate user data for machine joining control plane")
	}
	return userData, err
}
