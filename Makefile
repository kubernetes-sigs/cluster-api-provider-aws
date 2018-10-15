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


# This Makefile uses https://suva.sh/posts/well-documented-makefiles/,
# Please document new targets with ## <help> at the end of the 
# declaration

REGISTRY ?= "gcr.io"
REPOSITORY ?= "cluster-api-provider-aws/cluster-api-aws-controller"
TAG ?= "latest"

IMG := $(REGISTRY)/$(REPOSITORY):$(TAG)

# Go environment flags.
GOFLAGS += -ldflags '-extldflags "-static"'
GOREBUILD :=
GOPATH := $(shell go env GOPATH)

.DEFAULT_GOAL:=help

all: test manager clusterctl clusterawsadm ## Run tests and build binaries

vendor: ## Run dep ensure to create the vendor directory
	dep version || go get -u github.com/golang/dep/cmd/dep
	dep ensure

depend-update: ## Update dependencies to the latest available
	dep version || go get -u github.com/golang/dep/cmd/dep
	dep ensure -update

generate: vendor ## Generate any necessary code
	GOPATH=${GOPATH} go generate ./pkg/... ./cmd/...

genmocks: vendor ## Generate mocks
	hack/generate-mocks.sh "github.com/aws/aws-sdk-go/service/ec2/ec2iface EC2API" "cloud/aws/services/ec2/mock_ec2iface/mock.go"
	hack/generate-mocks.sh "github.com/aws/aws-sdk-go/service/elb/elbiface ELBAPI" "cloud/aws/services/elb/mock_elbiface/mock.go"
	hack/generate-mocks.sh "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1 MachineInterface" "cloud/aws/actuators/machine/mock_machineiface/mock.go"
	hack/generate-mocks.sh "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1 ClusterInterface" "cloud/aws/actuators/cluster/mock_clusteriface/mock.go"
	hack/generate-mocks.sh "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services EC2Interface" "cloud/aws/services/mocks/ec2.go"
	hack/generate-mocks.sh "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services ELBInterface" "cloud/aws/services/mocks/elb.go"

manager: generate fmt vet ## Build the AWS Cluster API Provider binary
	CGO_ENABLED=0 go install $(GOFLAGS) $(GOREBUILD) sigs.k8s.io/cluster-api-provider-aws/cmd/manager

clusterctl: generate fmt vet ## Build the AWS version of the clusterctl tool
	CGO_ENABLED=0 go install $(GOFLAGS) $(GOREBUILD) sigs.k8s.io/cluster-api-provider-aws/clusterctl

clusterawsadm: vendor ## Build clusterawsadm binary.
	CGO_ENABLED=0 go install $(GOFLAGS) $(GOREBUILD) sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm

cluster-api-dev-helper: vendor ## Build the cluster-api-dev-helper binary
	CGO_ENABLED=0 go install $(GOFLAGS) sigs.k8s.io/cluster-api-provider-aws/hack/cluster-api-dev-helper

bazel-build: ## Build the binaries using Bazel
	bazel build //cmd/clusterawsadm
	bazel build //cmd/manager
	bazel build //clusterctl

bazel-push: ## Make & push the container image using Bazel
	bazel run //cmd/manager:push \
		--define registry=$(REGISTRY) \
		--define repository=$(REPOSITORY) \
		--define tag=$(TAG)

test: generate fmt vet ## Run tests tagged integration
	go test -v -tags=integration ./pkg/... ./cmd/...

fmt: ## Run go fmt
	go fmt ./pkg/... ./cmd/...

vet: ## Run go vet
	go vet ./pkg/... ./cmd/...

lint: ## Run golint
	golint || go get -u golang.org/x/lint/golint
	golint -set_exit_status ./cmd/... ./pkg/... ./clusterctl/...

docker-build: generate fmt vet ## Build the docker image
	docker build . -t ${IMG}

docker-push: ## Push the built Docker image
	docker push ${IMG}

examples = clusterctl/examples/aws/out/cluster.yaml clusterctl/examples/aws/out/machines.yaml clusterctl/examples/aws/out/provider-components.yaml
templates = clusterctl/examples/aws/cluster.yaml.template clusterctl/examples/aws/machines.yaml.template clusterctl/examples/aws/provider-components.yaml.template
example: $(examples) ## Generate example cluster configuration in ./clusterctl/examples/aws/out
$(examples) : envfile $(templates)
	source ./envfile && cd ./clusterctl/examples/aws && ./generate-yaml.sh

envfile: ## Create the envfile and exit if the envfile doesn't already exist
	cp -n envfile.example envfile
	echo "\033[0;31mPlease fill out your envfile!\033[0m"
	exit 1

clean: ## Delete generated manifests and kubeconfig
	rm -rf clusterctl/examples/aws/out
	rm -f kubeconfig

reset-repo-acls: ## Ensure image is publicly readable on GCR
	gsutil -m acl ch -r -u AllUsers:READ gs://artifacts.$${REPOSITORY%%/*}.appspot.com


help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@echo ""
	@echo "See docs/development.md for more information"