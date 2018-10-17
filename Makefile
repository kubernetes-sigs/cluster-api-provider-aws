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
IMG ?= gcr.io/cluster-api-provider-aws/aws-cluster-api-controller:latest

# Go environment flags.
GOFLAGS += -ldflags '-extldflags "-static"'
GOREBUILD :=
GOPATH := $(shell go env GOPATH)

all: check-install test manager clusterctl clusterawsadm

.PHONY: vendor

# Dependency managemnt.
vendor:
	dep version || go get -u github.com/golang/dep/cmd/dep
	dep ensure

check-install:
	./scripts/check-install.sh

depend-update:
	dep version || go get -u github.com/golang/dep/cmd/dep
	dep ensure -update

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
	CGO_ENABLED=0 go install $(GOFLAGS) $(GOREBUILD) sigs.k8s.io/cluster-api-provider-aws/cmd/manager

# Build clusterctl binary.
clusterctl: check-install generate fmt vet
	CGO_ENABLED=0 go install $(GOFLAGS) $(GOREBUILD) sigs.k8s.io/cluster-api-provider-aws/clusterctl

# Build clusterawsadm binary.
clusterawsadm: check-install vendor
	CGO_ENABLED=0 go install $(GOFLAGS) $(GOREBUILD) sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm

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
	rm -rf clusterctl/examples/aws/out

manifests:
	go run vendor/sigs.k8s.io/controller-tools/cmd/controller-gen/main.go all
