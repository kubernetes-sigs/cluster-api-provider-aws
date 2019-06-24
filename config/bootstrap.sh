#!/bin/bash

cat <<HEREDOC > /root/user-data.sh
#!/bin/bash

################################################
######## Install packages and binaries
################################################

cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://packages.cloud.google.com/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://packages.cloud.google.com/yum/doc/yum-key.gpg https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg
exclude=kube*
EOF

setenforce 0
yum install -y docker
systemctl enable docker
systemctl start docker
yum install -y kubelet-1.13.1 kubeadm-1.13.1 kubectl-1.13.1 kubernetes-cni-0.6.0-0 --disableexcludes=kubernetes

cat <<EOF > /etc/default/kubelet
KUBELET_KUBEADM_EXTRA_ARGS=--cgroup-driver=systemd
EOF

echo '1' > /proc/sys/net/bridge/bridge-nf-call-iptables

curl -s https://api.github.com/repos/kubernetes-sigs/kustomize/releases/latest |\
  grep browser_download |\
  grep linux |\
  cut -d '"' -f 4 |\
  xargs curl -O -L
chmod u+x kustomize_*_linux_amd64
sudo mv kustomize_*_linux_amd64 /usr/bin/kustomize

sudo yum install -y git

################################################
######## Deploy kubernetes master
################################################

kubeadm init --apiserver-bind-port 8443 --token 2iqzqm.85bs0x6miyx1nm7l --apiserver-cert-extra-sans=$(curl icanhazip.com) --pod-network-cidr=192.168.0.0/16 -v 6

# Enable networking by default.
kubectl apply -f https://raw.githubusercontent.com/cloudnativelabs/kube-router/master/daemonset/kubeadm-kuberouter.yaml --kubeconfig /etc/kubernetes/admin.conf

# Binaries expected under /opt/cni/bin are actually under /usr/libexec/cni
if [[ ! -e /opt/cni/bin ]]; then
  mkdir -p /opt/cni/bin
  cp /usr/libexec/cni/bridge /opt/cni/bin
  cp /usr/libexec/cni/loopback /opt/cni/bin
  cp /usr/libexec/cni/host-local /opt/cni/bin
fi

mkdir -p /root/.kube
cp -i /etc/kubernetes/admin.conf /root/.kube/config
chown $(id -u):$(id -g) /root/.kube/config

################################################
######## Deploy machine-api plane
################################################

git clone https://github.com/openshift/cluster-api-provider-aws.git
cd cluster-api-provider-aws

cat <<EOF > secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: aws-credentials-secret
  namespace: default
type: Opaque
data:
  aws_access_key_id: FILLIN
  aws_secret_access_key: FILLIN
EOF

sudo kubectl apply -f secret.yaml

kustomize build config | sudo kubectl apply -f -

kubectl apply -f config/master-user-data-secret.yaml
kubectl apply -f config/master-machine.yaml

################################################
######## generate worker machineset user data
################################################

cat <<WORKERSET > /root/workerset-user-data.sh
#!/bin/bash

cat <<WORKERHEREDOC > /root/workerset-user-data.sh
#!/bin/bash

cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://packages.cloud.google.com/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://packages.cloud.google.com/yum/doc/yum-key.gpg https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg
exclude=kube*
EOF
setenforce 0
yum install -y docker
systemctl enable docker
systemctl start docker
yum install -y kubelet-1.13.1 kubeadm-1.13.1 kubernetes-cni-0.6.0-0 --disableexcludes=kubernetes

cat <<EOF > /etc/default/kubelet
KUBELET_KUBEADM_EXTRA_ARGS=--cgroup-driver=systemd
EOF

echo '1' > /proc/sys/net/bridge/bridge-nf-call-iptables

kubeadm join $(curl icanhazip.com):8443 --token 2iqzqm.85bs0x6miyx1nm7l --discovery-token-unsafe-skip-ca-verification
WORKERHEREDOC

bash /root/workerset-user-data.sh 2>&1 > /root/workerset-user-data.logs

WORKERSET

################################################
######## deploy worker user data and machineset
################################################

# NOTE: The secret is rendered twice, the first time when it's run during bootstrapping.
#       During bootstrapping, /root/workerset-user-data.sh does not exist yet.
#       So \$ needs to be used so the command is executed the second time
#       the script is executed.

cat <<EOF > /root/worker-secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: worker-user-data-secret
  namespace: default
type: Opaque
data:
  userData: \$(cat /root/workerset-user-data.sh | base64 --w=0)
EOF

sudo kubectl apply -f /root/worker-secret.yaml

sudo kubectl apply -f config/worker-machineset.yaml
HEREDOC

bash /root/user-data.sh 2>&1 > /root/user-data.logs
