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

ifndef FASTBUILD
.PHONY: vendor
endif

all: test manager clusterctl clusterawsadm

# Dependency management.
vendor:
	dep version || go get -u github.com/golang/dep/cmd/dep
	dep ensure

depend-update:
	dep version || go get -u github.com/golang/dep/cmd/dep
	dep ensure -update

# Generate code.
generate: vendor
	GOPATH=${GOPATH} go generate ./pkg/... ./cmd/...

genmocks: vendor
	hack/generate-mocks.sh "github.com/aws/aws-sdk-go/service/ec2/ec2iface EC2API" "cloud/aws/services/ec2/mock_ec2iface/mock.go"
	hack/generate-mocks.sh "github.com/aws/aws-sdk-go/service/elb/elbiface ELBAPI" "cloud/aws/services/elb/mock_elbiface/mock.go"
	hack/generate-mocks.sh "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1 MachineInterface" "cloud/aws/actuators/machine/mock_machineiface/mock.go"
	hack/generate-mocks.sh "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1 ClusterInterface" "cloud/aws/actuators/cluster/mock_clusteriface/mock.go"
	hack/generate-mocks.sh "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services EC2Interface" "cloud/aws/services/mocks/ec2.go"
	hack/generate-mocks.sh "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services ELBInterface" "cloud/aws/services/mocks/elb.go"

# Build manager binary.
manager: generate
	CGO_ENABLED=0 go install $(GOFLAGS) sigs.k8s.io/cluster-api-provider-aws/cmd/manager

# Build clusterctl binary.
clusterctl: generate
	CGO_ENABLED=0 go install $(GOFLAGS) sigs.k8s.io/cluster-api-provider-aws/clusterctl

# Build clusterawsadm binary.
clusterawsadm: vendor
	CGO_ENABLED=0 go install $(GOFLAGS) sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm

# Build cluster-api-dev-helper binary.
cluster-api-dev-helper: vendor
	CGO_ENABLED=0 go install $(GOFLAGS) sigs.k8s.io/cluster-api-provider-aws/hack/cluster-api-dev-helper

# Run tests
test: generate lint
	go test -v -tags=integration ./pkg/... ./cmd/...

# Run go fmt against code.
lint:
	golangci-lint > /dev/null || scripts/install-golint.sh -b $(GOPATH)/bin v1.10.2
	golangci-lint run

# Build the docker image
docker-build: generate
	docker build . -t ${IMG}

# Push the docker image
docker-push:
	docker push ${IMG}

# Example.
examples = clusterctl/examples/aws/out/cluster.yaml clusterctl/examples/aws/out/machines.yaml clusterctl/examples/aws/out/provider-components.yaml
templates = clusterctl/examples/aws/cluster.yaml.template clusterctl/examples/aws/machines.yaml.template clusterctl/examples/aws/provider-components.yaml.template
example: $(examples)
$(examples) : envfile $(templates)
	source ./envfile && cd ./clusterctl/examples/aws && ./generate-yaml.sh

envfile:
	# create the envfile and exit if the envfile doesn't already exist
	cp -n envfile.example envfile
	echo "\033[0;31mPlease fill out your envfile!\033[0m"
	exit 1

# Cleanup
clean:
	rm -rf clusterctl/examples/aws/out
