kind create cluster
clusterctl init --core cluster-api:v1.8.6 --bootstrap kubeadm:v1.8.6 --control-plane kubeadm:v1.8.6

make e2e-image 

RELEASE_TAG="e2e" make release-manifests

export IMAGE_TAG="gcr.io/k8s-staging-cluster-api/capa-manager:e2e"
kind load docker-image --name=capi-test $IMAGE_TAG
kubectl apply -f ./out/infrastructure-components.yaml
