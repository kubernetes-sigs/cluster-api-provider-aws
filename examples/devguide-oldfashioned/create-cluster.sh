kind create cluster
clusterctl init --core cluster-api:v1.8.6 --bootstrap kubeadm:v1.8.6 --control-plane kubeadm:v1.8.6

make e2e-image 

RELEASE_TAG="e2e" make release-manifests
kubectl apply -f ./out/infrastructure.yaml
