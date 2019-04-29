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

# If you update this file, please follow
# https://suva.sh/posts/well-documented-makefiles

.DEFAULT_GOAL:=help

# A release should define this with gcr.io/cluster-api-provider-aws
REGISTRY ?= gcr.io/$(shell gcloud config get-value project)

# A release should define this with IfNotPresent
PULL_POLICY ?= Always

# A release does not need to define this
MANAGER_IMAGE_NAME ?= cluster-api-aws-controller

# A release should define this with the next version after 0.0.4
MANAGER_IMAGE_TAG ?= dev

## Image URL to use all building/pushing image targets
DEPCACHEAGE ?= 24h # Enables caching for Dep
BAZEL_ARGS ?=

BAZEL_BUILD_ARGS := --define=REGISTRY=$(REGISTRY)\
 --define=PULL_POLICY=$(PULL_POLICY)\
 --define=MANAGER_IMAGE_NAME=$(MANAGER_IMAGE_NAME)\
 --define=MANAGER_IMAGE_TAG=$(MANAGER_IMAGE_TAG)\
$(BAZEL_ARGS)

# Bazel variables
BAZEL_VERSION := $(shell command -v bazel 2> /dev/null)
DEP ?= bazel run dep

# Determine the OS
HOSTOS := $(shell go env GOHOSTOS)
HOSTARCH := $(shell go env GOARCH)
BINARYPATHPATTERN :=${HOSTOS}_${HOSTARCH}_*

ifndef BAZEL_VERSION
    $(error "Bazel is not available. \
		Installation instructions can be found at \
		https://docs.bazel.build/versions/master/install.html")
endif

.PHONY: all
all: check-install test binaries

help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


## --------------------------------------
## Testing
## --------------------------------------

.PHONY: test
test: generate verify ## Run tests
	bazel test --nosandbox_debug //pkg/... //cmd/... $(BAZEL_ARGS)

.PHONY: integration
integration: generate verify ## Run integraion tests
	bazel test --define='gotags=integration' --test_output all //test/integration/...

JANITOR_ENABLED ?= 0
.PHONY: e2e
e2e: generate verify ## Run e2e tests
	JANITOR_ENABLED=$(JANITOR_ENABLED) ./hack/e2e.sh  $(BAZEL_BUILD_ARGS)

.PHONY: e2e-janitor
e2e-janitor:
	./hack/e2e-aws-janitor.sh

## --------------------------------------
## Docker
## --------------------------------------

.PHONY: docker-build
docker-build: generate ## Build the production docker image
	bazel run //cmd/manager:manager-image $(BAZEL_BUILD_ARGS)

.PHONY: docker-push
docker-push: ## Push production docker image
	bazel run //cmd/manager:manager-push $(BAZEL_BUILD_ARGS)

## --------------------------------------
## Cleanup / Verification
## --------------------------------------

.PHONY: clean
clean: ## Remove all generated files
	rm -rf cmd/clusterctl/examples/aws/out/
	rm -f kubeconfig
	rm -f minikube.kubeconfig
	rm -f bazel-*
	rm -rf out/
	rm -f cmd/clusterctl/examples/aws/provider-components-base.yaml

.PHONY: check-install
check-install: ## Checks that you've installed this repository correctly
	@./scripts/check-install.sh

.PHONY: verify
verify: ## Runs verification scripts to ensure correct execution
	./hack/verify_boilerplate.py

## --------------------------------------
## Manifests
## --------------------------------------

.PHONY: manifests
manifests: clusterawsadm cmd/clusterctl/examples/aws/provider-components-base.yaml ## Build example set of manifests from the current source
	./cmd/clusterctl/examples/aws/generate-yaml.sh

.PHONY: cmd/clusterctl/examples/aws/provider-components-base.yaml
cmd/clusterctl/examples/aws/provider-components-base.yaml:
	bazel build //cmd/clusterctl/examples/aws:provider-components-base $(BAZEL_BUILD_ARGS)
	install bazel-genfiles/cmd/clusterctl/examples/aws/provider-components-base.yaml cmd/clusterctl/examples/aws

## --------------------------------------
## Generate
## --------------------------------------

.PHONY: dep-ensure
dep-ensure: check-install ## Ensure dependencies are up to date
	@${DEP} ensure
	$(MAKE) gazelle

.PHONY: gazelle
gazelle: ## Run Bazel Gazelle
	bazel run //:gazelle $(BAZEL_ARGS)

.PHONY: generate
generate: ## Generate mocks, CRDs and runs `go generate` through Bazel
	GOPATH=$(shell go env GOPATH) bazel run //:generate $(BAZEL_ARGS)
	$(MAKE) dep-ensure
	bazel build $(BAZEL_ARGS) //pkg/cloud/aws/services/mocks:mocks \
		//pkg/cloud/aws/services/ec2/mock_ec2iface:mocks \
		//pkg/cloud/aws/services/elb/mock_elbiface:mocks
	./hack/copy-bazel-mocks.sh
	$(MAKE) generate-crds

.PHONY: generate-crds
generate-crds:
	bazel build //config
	cp -R bazel-genfiles/config/crds/* config/crds/
	cp -R bazel-genfiles/config/rbac/* config/rbac/

## --------------------------------------
## Linting
## --------------------------------------

.PHONY: lint
lint: dep-ensure ## Lint codebase
	bazel run //:lint $(BAZEL_ARGS)

lint-full: dep-ensure ## Run slower linters to detect possible issues
	bazel run //:lint-full $(BAZEL_ARGS)

## --------------------------------------
## Binaries
## --------------------------------------

.PHONY: binaries
binaries: generate manager clusterawsadm clusterctl ## Builds and installs all binaries

.PHONY: manager
manager: ## Build manager binary.
	bazel build //cmd/manager $(BAZEL_ARGS)
	install bazel-bin/cmd/manager/${BINARYPATHPATTERN}/manager $(shell go env GOPATH)/bin/aws-manager

.PHONY: clusterctl
clusterctl: ## Build clusterctl binary.
	bazel build --workspace_status_command=./hack/print-workspace-status.sh //cmd/clusterctl $(BAZEL_ARGS)
	install bazel-bin/cmd/clusterctl/${BINARYPATHPATTERN}/clusterctl $(shell go env GOPATH)/bin/clusterctl

.PHONY: clusterawsadm
clusterawsadm: ## Build clusterawsadm binary.
	bazel build --workspace_status_command=./hack/print-workspace-status.sh //cmd/clusterawsadm $(BAZEL_ARGS)
	install bazel-bin/cmd/clusterawsadm/${BINARYPATHPATTERN}/clusterawsadm $(shell go env GOPATH)/bin/clusterawsadm

## --------------------------------------
## Release
## --------------------------------------

.PHONY: release-artifacts
release-artifacts: ## Build release artifacts
	bazel build --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/clusterctl //cmd/clusterawsadm
	bazel build --platforms=@io_bazel_rules_go//go/toolchain:darwin_amd64 //cmd/clusterctl //cmd/clusterawsadm
	bazel build //cmd/clusterctl/examples/aws $(BAZEL_BUILD_ARGS)
	mkdir -p out
	install bazel-bin/cmd/clusterawsadm/darwin_amd64_pure_stripped/clusterawsadm out/clusterawsadm-darwin-amd64
	install bazel-bin/cmd/clusterawsadm/linux_amd64_pure_stripped/clusterawsadm out/clusterawsadm-linux-amd64
	install bazel-bin/cmd/clusterctl/darwin_amd64_pure_stripped/clusterctl out/clusterctl-darwin-amd64
	install bazel-bin/cmd/clusterctl/linux_amd64_pure_stripped/clusterctl out/clusterctl-linux-amd64
	install bazel-bin/cmd/clusterctl/examples/aws/aws.tar out/cluster-api-provider-aws-examples.tar

## --------------------------------------
## Define local development targets here
## --------------------------------------

.PHONY: binaries-dev
binaries-dev: ## Builds and installs all development binaries using go get
	go get -v ./...

.PHONY: create-cluster
create-cluster: binaries-dev ## Create a development Kubernetes cluster on AWS using examples
	clusterctl create cluster -v 4 \
	--provider aws \
	--bootstrap-type kind \
	-m ./cmd/clusterctl/examples/aws/out/machines.yaml \
	-c ./cmd/clusterctl/examples/aws/out/cluster.yaml \
	-p ./cmd/clusterctl/examples/aws/out/provider-components.yaml \
	-a ./cmd/clusterctl/examples/aws/out/addons.yaml

.PHONY: create-cluster-ha
create-cluste-ha: binaries-dev ## Create a development Kubernetes cluster on AWS using HA examples
	clusterctl create cluster -v 4 \
	--provider aws \
	--bootstrap-type kind \
	-m ./cmd/clusterctl/examples/aws/out/machines-ha.yaml \
	-c ./cmd/clusterctl/examples/aws/out/cluster.yaml \
	-p ./cmd/clusterctl/examples/aws/out/provider-components.yaml \
	-a ./cmd/clusterctl/examples/aws/out/addons.yaml

.PHONY: delete-cluster
delete-cluster: binaries-dev ## Deletes the development Kubernetes Cluster "test1"
	clusterctl delete cluster -v 4 \
	--bootstrap-type kind \
	--cluster test1 \
	--kubeconfig ./kubeconfig \
	-p ./cmd/clusterctl/examples/aws/out/provider-components.yaml \

kind-reset: ## Destroys the "clusterapi" kind cluster.
	kind delete cluster --name=clusterapi || true

.PHONY: reset-bazel
reset-bazel: ## Deep cleaning for bazel
	bazel clean --expunge
