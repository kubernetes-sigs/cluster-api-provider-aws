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
MANAGER_IMAGE_TAG ?= 0.0.2
MANAGER_IMAGE ?= $(STABLE_DOCKER_REPO)/$(MANAGER_IMAGE_NAME):$(MANAGER_IMAGE_TAG)
DEV_DOCKER_REPO ?= gcr.io/$(shell gcloud config get-value project)
DEV_MANAGER_IMAGE ?= $(DEV_DOCKER_REPO)/$(MANAGER_IMAGE_NAME):$(MANAGER_IMAGE_TAG)

DEPCACHEAGE ?= 24h # Enables caching for Dep
BAZEL_ARGS ?=

# Bazel variables
BAZEL_VERSION := $(shell command -v bazel 2> /dev/null)
DEP ?= bazel run dep

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

.PHONY: clusterctl
clusterctl: generate ## Build clusterctl binary.
	bazel build //cmd/clusterctl $(BAZEL_ARGS)

.PHONY: clusterawsadm
clusterawsadm: dep-ensure ## Build clusterawsadm binary.
	bazel build //cmd/clusterawsadm $(BAZEL_ARGS)

.PHONY: cluster-api-dev-helper
cluster-api-dev-helper: dep-ensure ## Build cluster-api-dev-helper binary
	bazel build //hack/cluster-api-dev-helper $(BAZEL_ARGS)

.PHONY: release-binaries
release-binaries: ## Build release binaries
	bazel build --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/clusterctl //cmd/clusterawsadm
	bazel build --platforms=@io_bazel_rules_go//go/toolchain:darwin_amd64 //cmd/clusterctl //cmd/clusterawsadm
	install bazel-bin/cmd/clusterawsadm/darwin_amd64_pure_stripped/clusterawsadm out/clusterawsadm-darwin-amd64
	install bazel-bin/cmd/clusterawsadm/linux_amd64_pure_stripped/clusterawsadm out/clusterawsadm-linux-amd64
	install bazel-bin/cmd/clusterctl/darwin_amd64_pure_stripped/clusterctl out/clusterctl-darwin-amd64
	install bazel-bin/cmd/clusterctl/linux_amd64_pure_stripped/clusterctl out/clusterctl-linux-amd64

.PHONY: test
test: generate ## Run tests
	bazel test --nosandbox_debug //pkg/... //cmd/... $(BAZEL_ARGS)

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

cmd/clusterctl/examples/aws/out/:
	./cmd/clusterctl/examples/aws/generate-yaml.sh

cmd/clusterctl/examples/aws/out/credentials: cmd/clusterctl/examples/aws/out/ ## Generate k8s secret for AWS credentials
	clusterawsadm alpha bootstrap generate-aws-default-profile > cmd/clusterctl/examples/aws/out/credentials

.PHONY: examples
examples: ## Generate example output
	$(MAKE) cmd/clusterctl/examples/aws/out/ IMAGE=${MANAGER_IMAGE}

.PHONY: examples-dev
examples-dev: ## Generate example output with developer image
	$(MAKE) cmd/clusterctl/examples/aws/out/ IMAGE=${DEV_MANAGER_IMAGE}

.PHONY: manifests
manifests: cmd/clusterctl/examples/aws/out/credentials ## Generate manifests for clusterctl
	go run vendor/sigs.k8s.io/controller-tools/cmd/controller-gen/main.go crd
	kustomize build config/default/ > cmd/clusterctl/examples/aws/out/provider-components.yaml
	echo "---" >> cmd/clusterctl/examples/aws/out/provider-components.yaml
	kustomize build vendor/sigs.k8s.io/cluster-api/config/default/ >> cmd/clusterctl/examples/aws/out/provider-components.yaml

.PHONY: manifests-dev
manifests-dev: ## Push development manifest
	MANAGER_IMAGE=$(DEV_MANAGER_IMAGE) $(MAKE) manifests

.PHONY: create-cluster
create-cluster: ## Create a Kubernetes cluster on AWS using examples
	clusterctl create cluster -v 4 --provider aws -m ./cmd/clusterctl/examples/aws/out/machines.yaml -c ./cmd/clusterctl/examples/aws/out/cluster.yaml -p ./cmd/clusterctl/examples/aws/out/provider-components.yaml -a ./cmd/clusterctl/examples/aws/out/addons.yaml

lint-full: dep-ensure ## Run slower linters to detect possible issues
	bazel run //:lint-full $(BAZEL_ARGS)

.PHONY: generate lint
ifneq ($(FASTBUILD),y)

## Define slow dependency targets here

reset-bazel: ## Deep cleaning for bazel
	bazel clean --expunge

generate: dep-ensure ## Run go generate
	GOPATH=$(shell go env GOPATH) bazel run //:generate $(BAZEL_ARGS)
	$(MAKE) dep-ensure

lint: dep-ensure ## Lint codebase
	@echo If you have genereated new mocks, run make copy-genmocks before linting
	bazel run //:lint $(BAZEL_ARGS)

else

## Add skips for slow depedency targets here

generate:
	@echo FASTBUILD is set: Skipping generate

lint:
	@echo FASTBUILD is set: Skipping lint

endif
