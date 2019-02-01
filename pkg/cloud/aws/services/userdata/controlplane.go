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

echo -n '{{.CACert}}' > /etc/kubernetes/pki/ca.crt
echo -n '{{.CAKey}}' > /etc/kubernetes/pki/ca.key
chmod 600 /etc/kubernetes/pki/ca.key

echo -n '{{.EtcdCACert}}' > /etc/kubernetes/pki/etcd/ca.crt
echo -n '{{.EtcdCAKey}}' > /etc/kubernetes/pki/etcd/ca.key
chmod 600 /etc/kubernetes/pki/etcd/ca.key

echo -n '{{.FrontProxyCACert}}' > /etc/kubernetes/pki/front-proxy-ca.crt
echo -n '{{.FrontProxyCAKey}}' > /etc/kubernetes/pki/front-proxy-ca.key
chmod 600 /etc/kubernetes/pki/front-proxy-ca.key

echo -n '{{.SaCert}}' > /etc/kubernetes/pki/sa.pub
echo -n '{{.SaKey}}' > /etc/kubernetes/pki/sa.key
chmod 600 /etc/kubernetes/pki/sa.key

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

mkdir -p /etc/kubernetes/pki/etcd

echo -n '{{.CACert}}' > /etc/kubernetes/pki/ca.crt
echo -n '{{.CAKey}}' > /etc/kubernetes/pki/ca.key
chmod 600 /etc/kubernetes/pki/ca.key

echo -n '{{.EtcdCACert}}' > /etc/kubernetes/pki/etcd/ca.crt
echo -n '{{.EtcdCAKey}}' > /etc/kubernetes/pki/etcd/ca.key
chmod 600 /etc/kubernetes/pki/etcd/ca.key

echo -n'{{.FrontProxyCACert}}' > /etc/kubernetes/pki/front-proxy-ca.crt
echo -n '{{.FrontProxyCAKey}}' > /etc/kubernetes/pki/front-proxy-ca.key
chmod 600 /etc/kubernetes/pki/front-proxy-ca.key

echo -n '{{.SaCert}}' > /etc/kubernetes/pki/sa.pub
echo -n '{{.SaKey}}' > /etc/kubernetes/pki/sa.key
chmod 600 /etc/kubernetes/pki/sa.key

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
	return cert != "" && key != ""
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

func (cpi *ControlPlaneInput) validateCertificates() error {
	if !isKeyPairValid(cpi.CACert, cpi.CAKey) {
		return errors.New("CA cert material in the ControlPlaneInput is missing cert/key")
	}

	if !isKeyPairValid(cpi.EtcdCACert, cpi.EtcdCAKey) {
		return errors.New("ETCD CA cert material in the ControlPlaneInput is  missing cert/key")
	}

	if !isKeyPairValid(cpi.FrontProxyCACert, cpi.FrontProxyCAKey) {
		return errors.New("FrontProxy CA cert material in ControlPlaneInput is  missing cert/key")
	}

	if !isKeyPairValid(cpi.SaCert, cpi.SaKey) {
		return errors.New("ServiceAccount cert material in ControlPlaneInput is  missing cert/key")
	}

	return nil
}

func (cpi *ContolPlaneJoinInput) validateCertificates() error {
	if !isKeyPairValid(cpi.CACert, cpi.CAKey) {
		return errors.New("CA cert material in the ContolPlaneJoinInput is  missing cert/key")
	}

	if !isKeyPairValid(cpi.EtcdCACert, cpi.EtcdCAKey) {
		return errors.New("ETCD cert material in the ContolPlaneJoinInput is  missing cert/key")
	}

	if !isKeyPairValid(cpi.FrontProxyCACert, cpi.FrontProxyCAKey) {
		return errors.New("FrontProxy cert material in ContolPlaneJoinInput is  missing cert/key")
	}

	if !isKeyPairValid(cpi.SaCert, cpi.SaKey) {
		return errors.New("ServiceAccount cert material in ContolPlaneJoinInput is  missing cert/key")
	}

	return nil
}

// NewControlPlane returns the user data string to be used on a controlplane instance.
func NewControlPlane(input *ControlPlaneInput) (string, error) {
	input.Header = defaultHeader
	if err := input.validateCertificates(); err != nil {
		return "", errors.Wrapf(err, "ControlPlaneInput is invalid")
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

	if err := input.validateCertificates(); err != nil {
		return "", errors.Wrapf(err, "ControlPlaneInput is invalid")
	}

	userData, err := generate("controlplane", controlPlaneJoinBashScript, input)
	if err != nil {
		return "", errors.Wrapf(err, "failed to generate user data for machine joining control plane")
	}
	return userData, err
}
