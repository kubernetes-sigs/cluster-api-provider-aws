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
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta1"
	"k8s.io/kubernetes/cmd/kubeadm/app/util"

	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
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
      {{.ClusterConfiguration | Indent 6}}
      ---
      {{.InitConfiguration | Indent 6}}
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
      {{.JoinConfiguration | Indent 6}}
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

	// TODO extract these since they contain values from above (not certs)
	ClusterConfiguration string
	InitConfiguration    string
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

	// TODO extract this
	JoinConfiguration string
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
func NewControlPlane(input *ControlPlaneInput, initConfiguration v1beta1.InitConfiguration) (string, error) {
	// Override critical variables
	// TODO(chuckha) add a warning if this is overwriting user input defined
	// in the configuration.
	initConfiguration.NodeRegistration.Name = "{{ ds.meta_data.hostname }}"
	initConfiguration.NodeRegistration.CRISocket = "/var/run/containerd/containerd.sock"

	if initConfiguration.NodeRegistration.KubeletExtraArgs == nil {
		initConfiguration.NodeRegistration.KubeletExtraArgs = map[string]string{}
	}
	initConfiguration.NodeRegistration.KubeletExtraArgs["cloud-provider"] = "aws"

	clusterConfiguration := v1beta1.ClusterConfiguration{
		APIServer: v1beta1.APIServer{
			CertSANs: []string{
				"{{ ds.meta_data.local_ipv4 }}",
				input.ELBAddress,
			},
			ControlPlaneComponent: v1beta1.ControlPlaneComponent{
				ExtraArgs: map[string]string{
					"cloud-provider": "aws",
				},
			},
		},
		ControlPlaneEndpoint: fmt.Sprintf("%s:%d", input.ELBAddress, 6443),
		ClusterName:          input.ClusterName,
		Networking: v1beta1.Networking{
			DNSDomain:     input.ServiceDomain,
			PodSubnet:     input.PodSubnet,
			ServiceSubnet: input.ServiceSubnet,
		},
		KubernetesVersion: input.KubernetesVersion,
	}
	initcfg, err := util.MarshalToYamlForCodecs(&initConfiguration, v1beta1.SchemeGroupVersion, v1alpha1.KubeadmCodecs)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal init configuration")
	}
	input.InitConfiguration = string(initcfg)

	clustercfg, err := util.MarshalToYamlForCodecs(&clusterConfiguration, v1beta1.SchemeGroupVersion, v1alpha1.KubeadmCodecs)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal cluster configuration")
	}
	input.ClusterConfiguration = string(clustercfg)
	input.Header = cloudConfigHeader
	if err := input.validateCertificates(); err != nil {
		return "", errors.Wrapf(err, "ControlPlaneInput is invalid")
	}

	fMap := map[string]interface{}{
		"Base64Encode": templateBase64Encode,
		"Indent":       templateYAMLIndent,
	}

	userData, err := generateWithFuncs("controlplane", controlPlaneCloudInit, funcMap(fMap), input)
	if err != nil {
		return "", errors.Wrapf(err, "failed to generate user data for new control plane machine")
	}

	return userData, err
}

// JoinControlPlane returns the user data string to be used on a new contrplplane instance.
func JoinControlPlane(input *ContolPlaneJoinInput, joinConfiguration v1beta1.JoinConfiguration) (string, error) {
	input.Header = cloudConfigHeader

	if err := input.validateCertificates(); err != nil {
		return "", errors.Wrapf(err, "ControlPlaneInput is invalid")
	}

	joinConfiguration.Discovery.BootstrapToken.Token = input.BootstrapToken
	joinConfiguration.Discovery.BootstrapToken.APIServerEndpoint = fmt.Sprintf("%s:%d", input.ELBAddress, 6443)
	joinConfiguration.Discovery.BootstrapToken.CACertHashes = append(joinConfiguration.Discovery.BootstrapToken.CACertHashes, input.CACertHash)
	joinConfiguration.NodeRegistration.Name = "{{ ds.meta_data.hostname }}"
	joinConfiguration.NodeRegistration.CRISocket = "/var/run/containerd/containerd.sock"
	joinConfiguration.NodeRegistration.KubeletExtraArgs["cloud-provider"] = "aws"
	joinConfiguration.ControlPlane.LocalAPIEndpoint.AdvertiseAddress = "{{ ds.meta_data.local_ipv4 }}"
	joinConfiguration.ControlPlane.LocalAPIEndpoint.BindPort = 6443
	joincfg, err := util.MarshalToYamlForCodecs(&joinConfiguration, v1beta1.SchemeGroupVersion, v1alpha1.KubeadmCodecs)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal cluster configuration")
	}
	input.JoinConfiguration = string(joincfg)

	fMap := map[string]interface{}{
		"Base64Encode": templateBase64Encode,
		"Indent":       templateYAMLIndent,
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

func templateYAMLIndent(i int, input string) string {
	split := strings.Split(input, "\n")
	ident := "\n" + strings.Repeat(" ", i)
	// Don't indent the first line, it's already indented in the template
	return strings.Join(split, ident)
}
