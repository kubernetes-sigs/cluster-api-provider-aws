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

const (
	controlPlaneBashScript = `{{.Header}}

mkdir -p /etc/kubernetes/pki

echo '{{.CACert}}' > /etc/kubernetes/pki/ca.crt
echo '{{.CAKey}}' > /etc/kubernetes/pki/ca.key

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

kubeadm init --config /tmp/kubeadm.yaml

kubectl -n kube-system --kubeconfig /etc/kubernetes/admin.conf \
create secret tls kubeadm-certs-ca \
--key /etc/kubernetes/pki/ca.key \
--cert /etc/kubernetes/pki/ca.crt

kubectl -n kube-system --kubeconfig /etc/kubernetes/admin.conf \
create secret tls kubeadm-certs-etcd-ca \
--key /etc/kubernetes/pki/etcd/ca.key \
--cert /etc/kubernetes/pki/etcd/ca.crt

kubectl -n kube-system --kubeconfig /etc/kubernetes/admin.conf \
create secret tls kubeadm-certs-front-proxy \
--key /etc/kubernetes/pki/front-proxy-ca.key \
--cert /etc/kubernetes/pki/front-proxy-ca.crt

# service account keys are different
tar -cvzf /etc/kubernetes/pki/sa-certs.tar.gz /etc/kubernetes/pki/sa.*
kubectl -n kube-system --kubeconfig /etc/kubernetes/admin.conf create secret generic kubeadm-sa-certs --from-file=/etc/kubernetes/pki/sa-certs.tar.gz
`

	controlPlaneJoinBashScript = `{{.Header}}

mkdir -p /etc/kubernetes/pki

echo '{{.CACert}}' > /etc/kubernetes/pki/ca.crt
echo '{{.CAKey}}' > /etc/kubernetes/pki/ca.key

echo '{{.KubeConfig}}' > /etc/kubernetes/admin.conf

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

kubectl --kubeconfig /etc/kubernetes/admin.conf -n kube-system \
get secret kubeadm-certs-ca -ojson | jq '.data."tls.crt"' -r | base64 --decode > /etc/kubernetes/pki/ca.crt
kubectl --kubeconfig /etc/kubernetes/admin.conf -n kube-system \
get secret kubeadm-certs-ca -ojson | jq '.data."tls.key"' -r | base64 --decode > /etc/kubernetes/pki/ca.key

mkdir -p /etc/kubernetes/pki/etcd
kubectl --kubeconfig /etc/kubernetes/admin.conf -n kube-system \
get secret kubeadm-certs-etcd-ca -ojson | jq '.data."tls.crt"' -r | base64 --decode > /etc/kubernetes/pki/etcd/ca.crt
kubectl --kubeconfig /etc/kubernetes/admin.conf -n kube-system \
get secret kubeadm-certs-etcd-ca -ojson | jq '.data."tls.key"' -r | base64 --decode > /etc/kubernetes/pki/etcd/ca.key

kubectl --kubeconfig /etc/kubernetes/admin.conf -n kube-system \
get secret kubeadm-certs-front-proxy -ojson | jq '.data."tls.crt"' -r | base64 --decode > /etc/kubernetes/pki/front-proxy-ca.crt
kubectl --kubeconfig /etc/kubernetes/admin.conf -n kube-system \
get secret kubeadm-certs-front-proxy -ojson | jq '.data."tls.key"' -r | base64 --decode > /etc/kubernetes/pki/front-proxy-ca.key


kubectl --kubeconfig /etc/kubernetes/admin.conf -n kube-system get secrets kubeadm-sa-certs -ojson | jq '.data."sa-certs.tar.gz"' -r | base64 --decode > /etc/kubernetes/pki/sa-certs.tar.gz
cd / 
tar -xvf /etc/kubernetes/pki/sa-certs.tar.gz

kubeadm join --config /tmp/kubeadm-controlplane-join-config.yaml --v 10
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

// TODO ashish-amarnath: incomplete
// https://github.com/kubernetes/kubernetes/blob/8d9ac261c4b49759179856d0a9db3ad4dc09e575/cmd/kubeadm/app/apis/kubeadm/types.go#L315:6

// ContolPlaneJoinInput defines context to generate controlplane instance user data for controlplane node join.
type ContolPlaneJoinInput struct {
	baseUserData

	CACertHash     string
	CACert         string
	CAKey          string
	BootstrapToken string
	ELBAddress     string
	KubeConfig     string
}

// NewControlPlane returns the user data string to be used on a controlplane instance.
func NewControlPlane(input *ControlPlaneInput) (string, error) {
	input.Header = defaultHeader
	return generate("controlplane", controlPlaneBashScript, input)
}

// JoinControlPlane returns the user data string to be used on a new contrplplane instance.
func JoinControlPlane(input *ContolPlaneJoinInput) (string, error) {
	input.Header = defaultHeader
	return generate("controlplane", controlPlaneJoinBashScript, input)
}
