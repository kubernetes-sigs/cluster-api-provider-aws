#!/bin/bash

cat <<HEREDOC > /root/user-data.sh
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
yum install -y kubelet-1.13.1 kubeadm-1.13.1 kubectl-1.13.1 --disableexcludes=kubernetes

cat <<EOF > /etc/default/kubelet
KUBELET_KUBEADM_EXTRA_ARGS=--cgroup-driver=systemd
EOF

echo '1' > /proc/sys/net/bridge/bridge-nf-call-iptables

kubeadm init --apiserver-bind-port 8443 --token 2iqzqm.85bs0x6miyx1nm7l --apiserver-cert-extra-sans=\$(curl -s http://169.254.169.254/latest/meta-data/public-hostname) --apiserver-cert-extra-sans=\$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4) --pod-network-cidr=192.168.0.0/16 -v 6

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
HEREDOC

bash /root/user-data.sh > /root/user-data.logs
