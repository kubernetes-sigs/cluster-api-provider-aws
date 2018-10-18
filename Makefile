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
IMG ?= gcr.io/cluster-api-provider-aws/cluster-api-aws-controller:latest

# Go environment flags.
GOFLAGS += -ldflags '-extldflags "-static"'
GOPATH := $(shell go env GOPATH)

all: check-install test manager clusterctl clusterawsadm

ifndef FASTBUILD
# If FASTBUILD isn't defined, always run dep ensure
.PHONY: vendor
endif

# Dependency managemnt.
vendor:
	${GOPATH}/bin/dep version || go get -u github.com/golang/dep/cmd/dep
	${GOPATH}/bin/dep ensure

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
	CGO_ENABLED=0 go install $(GOFLAGS) sigs.k8s.io/cluster-api-provider-aws/clusterctl

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
	golint -set_exit_status ./cmd/... ./pkg/... ./clusterctl/...

# Build the docker image
docker-build: generate fmt vet
	docker build . -t ${IMG}

# Push the docker image
docker-push:
	docker push ${IMG}

# Cleanup
clean:
	rm -rf out
	rm -f kubeconfig
	rm -f minikube.kubeconfig

manifests:
	go run vendor/sigs.k8s.io/controller-tools/cmd/controller-gen/main.go rbac --name aws-manager
	go run vendor/sigs.k8s.io/controller-tools/cmd/controller-gen/main.go crd

copy-creds-to-minikube: clusterawsadm
	minikube ssh 'mkdir -p .aws'
	echo "Hit ctrl+c next to work around minikube"
	source ./envfile && clusterawsadm alpha bootstrap generate-aws-default-profile | minikube ssh 'cat > .aws/credentials'
