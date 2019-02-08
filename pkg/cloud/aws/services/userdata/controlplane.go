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
	"encoding/base64"

	"github.com/pkg/errors"
)

const (
	controlPlaneCloudInit = `{{.Header}}
write_files:
-   path: /etc/kubernetes/pki/ca.crt
    encoding: "base64"
    owner: root:root
    permissions: '0640'
    content: |
      {{.CACert | Base64Encode}}

-   path: /etc/kubernetes/pki/ca.key
    encoding: "base64"
    owner: root:root
    permissions: '0600'
    content: |
      {{.CAKey | Base64Encode}}

-   path: /etc/kubernetes/pki/etcd/ca.crt
    encoding: "base64"
    owner: root:root
    permissions: '0640'
    content: |
      {{.EtcdCACert | Base64Encode}}

-   path: /etc/kubernetes/pki/etcd/ca.key
    encoding: "base64"
    owner: root:root
    permissions: '0600'
    content: |
      {{.EtcdCAKey | Base64Encode}}

-   path: /etc/kubernetes/pki/front-proxy-ca.crt
    encoding: "base64"
    owner: root:root
    permissions: '0640'
    content: |
      {{.FrontProxyCACert | Base64Encode}}

-   path: /etc/kubernetes/pki/front-proxy-ca.key
    encoding: "base64"
    owner: root:root
    permissions: '0600'
    content: |
      {{.FrontProxyCAKey | Base64Encode}}

-   path: /etc/kubernetes/pki/sa.pub
    encoding: "base64"
    owner: root:root
    permissions: '0640'
    content: |
      {{.SaCert | Base64Encode}}

-   path: /etc/kubernetes/pki/sa.key
    encoding: "base64"
    owner: root:root
    permissions: '0600'
    content: |
      {{.SaKey | Base64Encode}}

-   path: /tmp/kubeadm.yaml
    owner: root:root
    permissions: '0640'
    content: |
      ---
      apiVersion: kubeadm.k8s.io/v1beta1
      kind: ClusterConfiguration
      apiServer:
        certSANs:
          - {{ "{{ ds.meta_data.local_ipv4 }}" }}
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
        name: {{ "{{ ds.meta_data.hostname }}" }}
        criSocket: /var/run/containerd/containerd.sock
        kubeletExtraArgs:
          cloud-provider: aws
kubeadm:
  operation: init
  config: /tmp/kubeadm.yaml

`

	controlPlaneJoinCloudInit = `{{.Header}}
write_files:
-   path: /etc/kubernetes/pki/ca.crt
    encoding: "base64"
    owner: root:root
    permissions: '0640'
    content: |
      {{.CACert | Base64Encode}}

-   path: /etc/kubernetes/pki/ca.key
    encoding: "base64"
    owner: root:root
    permissions: '0600'
    content: |
      {{.CAKey | Base64Encode}}

-   path: /etc/kubernetes/pki/etcd/ca.crt
    encoding: "base64"
    owner: root:root
    permissions: '0640'
    content: |
      {{.EtcdCACert | Base64Encode}}

-   path: /etc/kubernetes/pki/etcd/ca.key
    encoding: "base64"
    owner: root:root
    permissions: '0600'
    content: |
      {{.EtcdCAKey | Base64Encode}}

-   path: /etc/kubernetes/pki/front-proxy-ca.crt
    encoding: "base64"
    owner: root:root
    permissions: '0640'
    content: |
      {{.FrontProxyCACert | Base64Encode}}

-   path: /etc/kubernetes/pki/front-proxy-ca.key
    encoding: "base64"
    owner: root:root
    permissions: '0600'
    content: |
      {{.FrontProxyCAKey | Base64Encode}}

-   path: /etc/kubernetes/pki/sa.pub
    encoding: "base64"
    owner: root:root
    permissions: '0640'
    content: |
      {{.SaCert | Base64Encode}}

-   path: /etc/kubernetes/pki/sa.key
    encoding: "base64"
    owner: root:root
    permissions: '0600'
    content: |
      {{.SaKey | Base64Encode}}

-   path: /tmp/kubeadm-controlplane-join-config.yaml
    owner: root:root
    permissions: '0640'
    content: |
      apiVersion: kubeadm.k8s.io/v1beta1
      kind: JoinConfiguration
      discovery:
        bootstrapToken:
          token: "{{.BootstrapToken}}"
          apiServerEndpoint: "{{.ELBAddress}}:6443"
          caCertHashes:
            - "{{.CACertHash}}"
      nodeRegistration:
        name: {{ "{{ ds.meta_data.hostname }}" }}
        criSocket: /var/run/containerd/containerd.sock
        kubeletExtraArgs:
          cloud-provider: aws
      controlPlane:
        localAPIEndpoint:
          advertiseAddress: {{ "{{ ds.meta_data.local_ipv4 }}" }}
          bindPort: 6443
kubeadm:
  operation: join
  config: /tmp/kubeadm-controlplane-join-config.yaml
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
	input.Header = cloudConfigHeader
	if err := input.validateCertificates(); err != nil {
		return "", errors.Wrapf(err, "ControlPlaneInput is invalid")
	}

	fMap := map[string]interface{}{
		"Base64Encode": templateBase64Encode,
	}

	userData, err := generateWithFuncs("controlplane", controlPlaneCloudInit, funcMap(fMap), input)
	if err != nil {
		return "", errors.Wrapf(err, "failed to generate user data for new control plane machine")
	}

	return userData, err
}

// JoinControlPlane returns the user data string to be used on a new contrplplane instance.
func JoinControlPlane(input *ContolPlaneJoinInput) (string, error) {
	input.Header = cloudConfigHeader

	if err := input.validateCertificates(); err != nil {
		return "", errors.Wrapf(err, "ControlPlaneInput is invalid")
	}

	fMap := map[string]interface{}{
		"Base64Encode": templateBase64Encode,
	}

	userData, err := generateWithFuncs("controlplane", controlPlaneJoinCloudInit, funcMap(fMap), input)
	if err != nil {
		return "", errors.Wrapf(err, "failed to generate user data for machine joining control plane")
	}
	return userData, err
}

func templateBase64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}
