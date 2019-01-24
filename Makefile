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

FASTBUILD ?= n ## Set FASTBUILD=y (case-sensitive) to skip some slow tasks

## Image URL to use all building/pushing image targets
STABLE_DOCKER_REPO ?= gcr.io/cluster-api-provider-aws
MANAGER_IMAGE_NAME ?= cluster-api-aws-controller
MANAGER_IMAGE_TAG ?= 0.0.3
MANAGER_IMAGE ?= $(STABLE_DOCKER_REPO)/$(MANAGER_IMAGE_NAME):$(MANAGER_IMAGE_TAG)
DEV_DOCKER_REPO ?= gcr.io/$(shell gcloud config get-value project)
DEV_MANAGER_IMAGE ?= $(DEV_DOCKER_REPO)/$(MANAGER_IMAGE_NAME):$(MANAGER_IMAGE_TAG)

DEPCACHEAGE ?= 24h # Enables caching for Dep
BAZEL_ARGS ?=

# Bazel variables
BAZEL_VERSION := $(shell command -v bazel 2> /dev/null)
DEP ?= bazel run dep

# determine the OS
HOSTOS := $(shell go env GOHOSTOS)
HOSTARCH := $(shell go env GOARCH)
BINARYPATHPATTERN :=${HOSTOS}_${HOSTARCH}_*

ifndef BAZEL_VERSION
    $(error "Bazel is not available. \
		Installation instructions can be found at \
		https://docs.bazel.build/versions/master/install.html")
endif

.PHONY: all
all: check-install test manager clusterctl clusterawsadm

help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: dep-ensure
dep-ensure: check-install ## Ensure dependencies are up to date
	@echo Checking status of dependencies
	@${DEP} status 2>&1 > /dev/null || make dep-install
	@echo Finished verifying dependencies

.PHONY: dep-install
dep-install: ## Force install go dependencies
	${DEP} ensure
	bazel run //:gazelle $(BAZEL_ARGS)

.PHONY: gazelle
gazelle: ## Run Bazel Gazelle
	bazel run //:gazelle $(BAZEL_ARGS)

.PHONY: check-install
check-install: ## Checks that you've installed this repository correctly
	@./scripts/check-install.sh

.PHONY: manager
manager: generate  ## Build manager binary.
	bazel build //cmd/manager $(BAZEL_ARGS)
	install bazel-bin/cmd/manager/${BINARYPATHPATTERN}/manager $(shell go env GOPATH)/bin/aws-manager

.PHONY: clusterctl
clusterctl: generate ## Build clusterctl binary.
	bazel build --workspace_status_command=./hack/print-workspace-status.sh //cmd/clusterctl $(BAZEL_ARGS)
	install bazel-bin/cmd/clusterctl/${BINARYPATHPATTERN}/clusterctl $(shell go env GOPATH)/bin/clusterctl

.PHONY: clusterawsadm
clusterawsadm: dep-ensure ## Build clusterawsadm binary.
	bazel build --workspace_status_command=./hack/print-workspace-status.sh //cmd/clusterawsadm $(BAZEL_ARGS)
	install bazel-bin/cmd/clusterawsadm/${BINARYPATHPATTERN}/clusterawsadm $(shell go env GOPATH)/bin/clusterawsadm

.PHONY: release-binaries
release-binaries: ## Build release binaries
	bazel build --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/clusterctl //cmd/clusterawsadm
	bazel build --platforms=@io_bazel_rules_go//go/toolchain:darwin_amd64 //cmd/clusterctl //cmd/clusterawsadm
	mkdir -p out
	install bazel-bin/cmd/clusterawsadm/darwin_amd64_pure_stripped/clusterawsadm out/clusterawsadm-darwin-amd64
	install bazel-bin/cmd/clusterawsadm/linux_amd64_pure_stripped/clusterawsadm out/clusterawsadm-linux-amd64
	install bazel-bin/cmd/clusterctl/darwin_amd64_pure_stripped/clusterctl out/clusterctl-darwin-amd64
	install bazel-bin/cmd/clusterctl/linux_amd64_pure_stripped/clusterctl out/clusterctl-linux-amd64

.PHONY: test verify
test: generate verify ## Run tests
	bazel test --nosandbox_debug //pkg/... //cmd/... $(BAZEL_ARGS)

verify:
	./hack/verify_boilerplate.py

.PHONY: copy-genmocks
copy-genmocks: ## Copies generated mocks into the repository
	cp -Rf bazel-genfiles/pkg/* pkg/

BAZEL_DOCKER_ARGS_COMMON := --define=MANAGER_IMAGE_NAME=$(MANAGER_IMAGE_NAME) --define=MANAGER_IMAGE_TAG=$(MANAGER_IMAGE_TAG) $(BAZEL_ARGS)
BAZEL_DOCKER_ARGS := --define=DOCKER_REPO=$(STABLE_DOCKER_REPO) $(BAZEL_DOCKER_ARGS_COMMON)
BAZEL_DOCKER_ARGS_DEV := --define=DOCKER_REPO=$(DEV_DOCKER_REPO) $(BAZEL_DOCKER_ARGS_COMMON)

.PHONY: docker-build
docker-build: generate ## Build the production docker image
	bazel run //cmd/manager:manager-image $(BAZEL_DOCKER_ARGS)

.PHONY: docker-build-dev
docker-build-dev: generate ## Build the development docker image
	bazel run //cmd/manager:manager-image $(BAZEL_DOCKER_ARGS_DEV)

.PHONY: docker-push
docker-push: generate ## Push production docker image
	bazel run //cmd/manager:manager-push $(BAZEL_DOCKER_ARGS)

.PHONY: docker-push-dev
docker-push-dev: generate ## Push development image
	bazel run //cmd/manager:manager-push $(BAZEL_DOCKER_ARGS_DEV)

.PHONY: clean
clean: ## Remove all generated files
	rm -rf cmd/clusterctl/examples/aws/out/
	rm -f kubeconfig
	rm -f minikube.kubeconfig
	rm -f bazel-*
	rm -rf out/

.PHONY: reset-bazel
reset-bazel: ## Deep cleaning for bazel
	bazel clean --expunge

cmd/clusterctl/examples/aws/out:
	./cmd/clusterctl/examples/aws/generate-yaml.sh

cmd/clusterctl/examples/aws/out/credentials: cmd/clusterctl/examples/aws/out ## Generate k8s secret for AWS credentials
	clusterawsadm alpha bootstrap generate-aws-default-profile > cmd/clusterctl/examples/aws/out/credentials

.PHONY: examples
examples: ## Generate example output
	$(MAKE) cmd/clusterctl/examples/aws/out MANAGER_IMAGE=${MANAGER_IMAGE}

.PHONY: examples-dev
examples-dev: ## Generate example output with developer image
	$(MAKE) cmd/clusterctl/examples/aws/out MANAGER_IMAGE=${DEV_MANAGER_IMAGE}

.PHONY: manifests
manifests: cmd/clusterctl/examples/aws/out/credentials ## Generate manifests for clusterctl
	go run vendor/sigs.k8s.io/controller-tools/cmd/controller-gen/main.go crd
	kustomize build config/default/ > cmd/clusterctl/examples/aws/out/provider-components.yaml
	echo "---" >> cmd/clusterctl/examples/aws/out/provider-components.yaml
	kustomize build vendor/sigs.k8s.io/cluster-api/config/default/ >> cmd/clusterctl/examples/aws/out/provider-components.yaml

.PHONY: manifests-dev
manifests-dev: dep-ensure dep-install binaries-dev ## Builds development manifests
	MANAGER_IMAGE=$(DEV_MANAGER_IMAGE) MANAGER_IMAGE_PULL_POLICY="Always" $(MAKE) manifests

# TODO(vincepri): This should move to rebuild Bazel binaries once every
# make target uses Bazel bins to run operations.
.PHONY: binaries-dev
binaries-dev: ## Builds and installs the binaries on the local GOPATH
	go get -v ./...
	go install -v ./...

.PHONY: create-cluster
create-cluster: ## Create a Kubernetes cluster on AWS using examples
	clusterctl create cluster -v 3 \
	--provider aws \
	--bootstrap-type kind \
	-m ./cmd/clusterctl/examples/aws/out/machines.yaml \
	-c ./cmd/clusterctl/examples/aws/out/cluster.yaml \
	-p ./cmd/clusterctl/examples/aws/out/provider-components.yaml \
	-a ./cmd/clusterctl/examples/aws/out/addons.yaml

lint-full: dep-ensure ## Run slower linters to detect possible issues
	bazel run //:lint-full $(BAZEL_ARGS)

## Define kind dependencies here.

kind-reset: ## Destroys the "clusterapi" kind cluster.
	kind delete cluster --name=clusterapi || true

ifneq ($(FASTBUILD),y)

## Define slow dependency targets here

.PHONY: generate
generate: gazelle dep-ensure ## Run go generate
	GOPATH=$(shell go env GOPATH) bazel run //:generate $(BAZEL_ARGS)
	$(MAKE) dep-ensure
	bazel build $(BAZEL_ARGS) //pkg/cloud/aws/services/mocks:go_mock_interfaces \
		//pkg/cloud/aws/services/ec2/mock_ec2iface:go_default_library \
		//pkg/cloud/aws/services/elb/mock_elbiface:go_default_library
	cp -Rf bazel-genfiles/pkg/* pkg/

.PHONY: lint
lint: dep-ensure ## Lint codebase
	@echo If you have genereated new mocks, run make copy-genmocks before linting
	bazel run //:lint $(BAZEL_ARGS)

else

## Add skips for slow depedency targets here

.PHONY: generate
generate:
	@echo FASTBUILD is set: Skipping generate

.PHONY: lint
lint:
	@echo FASTBUILD is set: Skipping lint

endif
