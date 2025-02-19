kind create cluster
clusterctl init --core cluster-api:v1.8.6 --bootstrap kubeadm:v1.8.6 --control-plane kubeadm:v1.8.6

make e2e-image # fails 

 > [toolchain 1/1] FROM docker.io/library/golang:1.22.6:
------
failed to load cache key: error getting credentials - err: docker-credential-dev-containers-fe22015d-c8a4-4559-80b9-ae4c6e9e73d1 resolves to executable in current directory (./docker-credential-dev-containers-fe22015d-c8a4-4559-80b9-ae4c6e9e73d1), out: ``
make: *** [Makefile:408: e2e-image] Error 1


RELEASE_TAG="e2e" make release-manifests
kubectl apply -f ./out/infrastructure.yaml
