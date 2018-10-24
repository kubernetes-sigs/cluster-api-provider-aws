# Copyright 2018 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Image URL to use all building/pushing image targets
MANAGER_IMAGE ?= gcr.io/cluster-api-provider-aws/cluster-api-aws-controller:latest
DEV_MANAGER_IMAGE ?= gcr.io/$(shell gcloud config get-value project)/cluster-api-aws-controller:0.0.1

# Go environment flags.
GOFLAGS += -ldflags '-extldflags "-static"'
GOPATH := $(shell go env GOPATH)

all: check-install test manager clusterctl clusterawsadm

ifdef REVENDOR
.PHONY: vendor
endif

# Dependency managemnt.
vendor:
	${GOPATH}/bin/dep version || go get -u github.com/golang/dep/cmd/dep
	${GOPATH}/bin/dep ensure
	bazel run //:gazelle

check-install:
	./scripts/check-install.sh

# Generate code.
generate: vendor
	GOPATH=${GOPATH} go generate ./pkg/... ./cmd/...

genmocks: vendor
	hack/generate-mocks.sh "github.com/aws/aws-sdk-go/service/ec2/ec2iface EC2API" "pkg/cloud/aws/services/ec2/mock_ec2iface/mock.go"
	hack/generate-mocks.sh "github.com/aws/aws-sdk-go/service/elb/elbiface ELBAPI" "pkg/cloud/aws/services/elb/mock_elbiface/mock.go"
	hack/generate-mocks.sh "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1 MachineInterface" "pkg/cloud/aws/actuators/machine/mock_machineiface/mock.go"
	hack/generate-mocks.sh "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1 ClusterInterface" "pkg/cloud/aws/actuators/cluster/mock_clusteriface/mock.go"
	hack/generate-mocks.sh "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services EC2Interface" "pkg/cloud/aws/services/mocks/ec2.go"
	hack/generate-mocks.sh "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services ELBInterface" "pkg/cloud/aws/services/mocks/elb.go"

# Build manager binary.
manager: generate fmt vet
	CGO_ENABLED=0 go install $(GOFLAGS) sigs.k8s.io/cluster-api-provider-aws/cmd/manager

# Build clusterctl binary.
clusterctl: check-install generate fmt vet
	CGO_ENABLED=0 go install $(GOFLAGS) sigs.k8s.io/cluster-api-provider-aws/cmd/clusterctl

# Build clusterawsadm binary.
clusterawsadm: check-install vendor
	CGO_ENABLED=0 go install $(GOFLAGS) sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm

# Build cluster-api-dev-helper binary.
cluster-api-dev-helper: check-install vendor
	CGO_ENABLED=0 go install $(GOFLAGS) sigs.k8s.io/cluster-api-provider-aws/hack/cluster-api-dev-helper

# Run tests
test: generate fmt vet
	go test -v -tags=integration ./pkg/... ./cmd/...

# Run go fmt against code.
fmt:
	go fmt ./pkg/... ./cmd/...

# Run go vet against code.
vet:
	go vet ./pkg/... ./cmd/...

# Run go lint.
lint:
	golint || go get -u golang.org/x/lint/golint
	golint -set_exit_status ./cmd/... ./pkg/...

# Build the docker image
docker-build: generate fmt vet
	docker build . -t ${MANAGER_IMAGE}

# Push the docker image
docker-push:
	docker push ${MANAGER_IMAGE}

# Cleanup
clean:
	rm -rf cmd/clusterctl/examples/aws/out/
	rm -f kubeconfig
	rm -f minikube.kubeconfig

# Manifests.
cmd/clusterctl/examples/aws/out/:
	MANAGER_IMAGE=${MANAGER_IMAGE} ./cmd/clusterctl/examples/aws/generate-yaml.sh

cmd/clusterctl/examples/aws/out/credentials: cmd/clusterctl/examples/aws/out/ clusterawsadm
	clusterawsadm alpha bootstrap generate-aws-default-profile > cmd/clusterctl/examples/aws/out/credentials

manifests: cmd/clusterctl/examples/aws/out/credentials
	go run vendor/sigs.k8s.io/controller-tools/cmd/controller-gen/main.go crd
	kustomize build config/default/ > cmd/clusterctl/examples/aws/out/provider-components.yaml
	echo "---" >> cmd/clusterctl/examples/aws/out/provider-components.yaml
	kustomize build vendor/sigs.k8s.io/cluster-api/config/default/ >> cmd/clusterctl/examples/aws/out/provider-components.yaml

# Create cluster.
create-cluster:
	clusterctl create cluster -v3 --provider aws -m ./cmd/clusterctl/examples/aws/out/machines.yaml -c ./cmd/clusterctl/examples/aws/out/cluster.yaml -p ./cmd/clusterctl/examples/aws/out/provider-components.yaml

# Development.
dev-images:
	MANAGER_IMAGE=$(DEV_MANAGER_IMAGE) $(MAKE) docker-build docker-push

dev-manifests:
	MANAGER_IMAGE=$(DEV_MANAGER_IMAGE) $(MAKE) manifests