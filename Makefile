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
REGISTRY_DEV ?= gcr.io/$(shell gcloud config get-value project)
DEPCACHEAGE ?= 24h # Enables caching for Dep
BAZEL_ARGS ?=
BAZEL_DOCKER_ARGS := --define=REGISTRY_DEV=$(REGISTRY_DEV) $(BAZEL_ARGS)

# Bazel variables
BAZEL_VERSION := $(shell command -v bazel 2> /dev/null)
DEP ?= bazel run dep

# determine the OS
HOSTOS := $(shell go env GOHOSTOS)
HOSTARCH := $(shell go env GOARCH)
GOPATH := $(shell go env GOPATH)
BINARYPATHPATTERN :=${HOSTOS}_${HOSTARCH}_*

export GOPATH

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

.PHONY: release-binaries
release-binaries: ## Build release binaries
	bazel build --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64,@io_bazel_rules_go//go/toolchain:darwin_amd64 --build_tag_filters=cli //...
	mkdir -p out
	install bazel-bin/cmd/clusterawsadm/darwin_amd64_pure_stripped/clusterawsadm out/clusterawsadm-darwin-amd64
	install bazel-bin/cmd/clusterawsadm/linux_amd64_pure_stripped/clusterawsadm out/clusterawsadm-linux-amd64
	install bazel-bin/cmd/clusterctl/darwin_amd64_pure_stripped/clusterctl out/clusterctl-darwin-amd64
	install bazel-bin/cmd/clusterctl/linux_amd64_pure_stripped/clusterctl out/clusterctl-linux-amd64

.PHONY: test verify
test: generate verify ## Run tests
	bazel test  //pkg/... //cmd/... $(BAZEL_ARGS)

verify:
	./hack/verify_boilerplate.py

.PHONY: docker-build
docker-build: generate provider-components ## Build the production docker image
	bazel run //cmd/manager:manager-image $(BAZEL_DOCKER_ARGS)

.PHONY: docker-push
docker-push: generate provider-components ## Push production docker image
	bazel run //cmd/manager:manager-push $(BAZEL_DOCKER_ARGS)

.PHONY: docker-push-dev
docker-push-dev: generate provider-components-dev ## Push development image
	bazel run //cmd/manager:manager-push-dev $(BAZEL_DOCKER_ARGS)

.PHONY: provider-components-dev
provider-components-dev: provider-components ## Generate Kustomize version patch for the Docker image
	bazel build //cmd/manager:manager-version-patch-dev $(BAZEL_DOCKER_ARGS)
	install bazel-genfiles/cmd/manager/manager-version-patch-dev.yaml cmd/manager

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

.PHONY: cli-dev
cli-dev: ## Builds and installs the binaries on the local GOPATH
	bazel build --build_tag_filters=cli //... $(BAZEL_ARGS)
	install bazel-bin/cmd/clusterawsadm/${BINARYPATHPATTERN}/clusterawsadm $(GOPATH)/bin/clusterawsadm
	install bazel-bin/cmd/clusterctl/${BINARYPATHPATTERN}/clusterctl $(GOPATH)/bin/clusterctl


.PHONY: lint-full
lint-full: dep-ensure ## Run slower linters to detect possible issues
	bazel run //:lint-full $(BAZEL_ARGS)

kind-reset: ## Destroys the "clusterapi" kind cluster.
	bazel run //:kind-reset $(BAZEL_ARGS)

.PHONY: config
config: ## Create Kubernetes API components
	bazel build //config
	cp -Rf bazel-genfiles/config/* config/

ifneq ($(FASTBUILD),y)

## Define slow dependency targets here

.PHONY: generate
generate: gazelle dep-ensure ## Run go generate
	bazel build --build_tag_filters=generated //... $(BAZEL_ARGS)

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
