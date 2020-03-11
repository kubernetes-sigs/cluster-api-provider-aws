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

# Ensure Make is run with bash shell as some syntax below is bash-specific
SHELL:=/usr/bin/env bash

.DEFAULT_GOAL:=help

# Use GOPROXY environment variable if set
GOPROXY := $(shell go env GOPROXY)
ifeq ($(GOPROXY),)
GOPROXY := https://proxy.golang.org
endif
export GOPROXY

# Active module mode, as we use go modules to manage dependencies
export GO111MODULE=on

# Directories.
TOOLS_DIR := hack/tools
TOOLS_BIN_DIR := $(TOOLS_DIR)/bin
BIN_DIR := bin
TEST_E2E_DIR := test/e2e

# Binaries.
CLUSTERCTL := $(BIN_DIR)/clusterctl
CONTROLLER_GEN := $(TOOLS_BIN_DIR)/controller-gen
ENVSUBST := $(TOOLS_BIN_DIR)/envsubst
GOLANGCI_LINT := $(TOOLS_BIN_DIR)/golangci-lint
MOCKGEN := $(TOOLS_BIN_DIR)/mockgen
CONVERSION_GEN := $(TOOLS_BIN_DIR)/conversion-gen
RELEASE_NOTES_BIN := bin/release-notes
RELEASE_NOTES := $(TOOLS_DIR)/$(RELEASE_NOTES_BIN)
GINKGO := $(abspath $(TOOLS_BIN_DIR)/ginkgo)

# Define Docker related variables. Releases should modify and double check these vars.
REGISTRY ?= gcr.io/$(shell gcloud config get-value project)
STAGING_REGISTRY := gcr.io/k8s-staging-cluster-api-aws
PROD_REGISTRY := us.gcr.io/k8s-artifacts-prod/cluster-api-aws
IMAGE_NAME ?= cluster-api-aws-controller
CONTROLLER_IMG ?= $(REGISTRY)/$(IMAGE_NAME)
TAG ?= dev
ARCH ?= amd64
ALL_ARCH = amd64 arm arm64 ppc64le s390x

# Allow overriding manifest generation destination directory
MANIFEST_ROOT ?= config
CRD_ROOT ?= $(MANIFEST_ROOT)/crd/bases
WEBHOOK_ROOT ?= $(MANIFEST_ROOT)/webhook
RBAC_ROOT ?= $(MANIFEST_ROOT)/rbac

# Allow overriding the imagePullPolicy
PULL_POLICY ?= Always

# Hosts running SELinux need :z added to volume mounts
SELINUX_ENABLED := $(shell cat /sys/fs/selinux/enforce 2> /dev/null || echo 0)

ifeq ($(SELINUX_ENABLED),1)
  DOCKER_VOL_OPTS?=:z
endif

# Set build time variables including version details
LDFLAGS := $(shell source ./hack/version.sh; version::ldflags)

GOLANG_VERSION := 1.13.8

# 'functional tests' as the ginkgo filter will run ALL tests ~ 2 hours @ 3 node concurrency.
E2E_FOCUS := "functional tests"
# Instead, you can run a quick smoke test, it should run fast (9 minutes)... 
# E2E_FOCUS := "Create cluster with name having"

## --------------------------------------
## Help
## --------------------------------------

help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

## --------------------------------------
## Testing
## --------------------------------------

.PHONY: test
test: ## Run tests
	source ./scripts/fetch_ext_bins.sh; fetch_tools; setup_envs; go test -v ./...

.PHONY: test-integration
test-integration: ## Run integration tests
	source ./scripts/fetch_ext_bins.sh; fetch_tools; setup_envs; go test -v -tags=integration ./test/integration/...

.PHONY: test-e2e
test-e2e: $(GINKGO) ## Run e2e tests
	PULL_POLICY=IfNotPresent $(MAKE) docker-build
	cd $(TEST_E2E_DIR); $(GINKGO) -nodes=2 -v -tags=e2e -focus=$(E2E_FOCUS) ./... -- -managerImage=$(CONTROLLER_IMG)-$(ARCH):$(TAG)

.PHONY: test-conformance
test-conformance: ## Run conformance test on workload cluster
	PULL_POLICY=IfNotPresent $(MAKE) docker-build
	cd $(TEST_E2E_DIR); go test -v -tags=e2e -timeout=4h . -args -ginkgo.v -ginkgo.focus "conformance tests" --managerImage $(CONTROLLER_IMG)-$(ARCH):$(TAG)

## --------------------------------------
## Binaries
## --------------------------------------
$(GINKGO): $(TOOLS_DIR)/go.mod
	cd $(TOOLS_DIR) && go build -tags=tools -o $(BIN_DIR)/ginkgo github.com/onsi/ginkgo/ginkgo

.PHONY: binaries
binaries: manager clusterawsadm ## Builds and installs all binaries

.PHONY: manager
manager: ## Build manager binary.
	go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/manager .

.PHONY: clusterawsadm
clusterawsadm: ## Build clusterawsadm binary.
	go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/clusterawsadm ./cmd/clusterawsadm

## --------------------------------------
## Tooling Binaries
## --------------------------------------

$(CLUSTERCTL): go.mod ## Build clusterctl binary.
	go build -o $(BIN_DIR)/clusterctl sigs.k8s.io/cluster-api/cmd/clusterctl

$(CONTROLLER_GEN): $(TOOLS_DIR)/go.mod # Build controller-gen from tools folder.
	cd $(TOOLS_DIR); go build -tags=tools -o $(BIN_DIR)/controller-gen sigs.k8s.io/controller-tools/cmd/controller-gen

$(ENVSUBST): $(TOOLS_DIR)/go.mod # Build envsubst from tools folder.
	cd $(TOOLS_DIR); go build -tags=tools -o $(BIN_DIR)/envsubst github.com/a8m/envsubst/cmd/envsubst

$(GOLANGCI_LINT): $(TOOLS_DIR)/go.mod # Build golangci-lint from tools folder.
	cd $(TOOLS_DIR); go build -tags=tools -o $(BIN_DIR)/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint

$(MOCKGEN): $(TOOLS_DIR)/go.mod # Build mockgen from tools folder.
	cd $(TOOLS_DIR); go build -tags=tools -o $(BIN_DIR)/mockgen github.com/golang/mock/mockgen

$(CONVERSION_GEN): $(TOOLS_DIR)/go.mod
	cd $(TOOLS_DIR); go build -tags=tools -o $(BIN_DIR)/conversion-gen k8s.io/code-generator/cmd/conversion-gen

$(RELEASE_NOTES) : $(TOOLS_DIR)/go.mod
	cd $(TOOLS_DIR) && go build -tags tools -o $(BIN_DIR)/release-notes sigs.k8s.io/cluster-api/hack/tools/release

## --------------------------------------
## Linting
## --------------------------------------

.PHONY: lint
lint: $(GOLANGCI_LINT) ## Lint codebase
	$(GOLANGCI_LINT) run -v --fast=false

## --------------------------------------
## Generate
## --------------------------------------

.PHONY: modules
modules: ## Runs go mod to ensure proper vendoring.
	go mod tidy
	cd $(TOOLS_DIR); go mod tidy
	cd $(TEST_E2E_DIR); go mod tidy

.PHONY: generate
generate: ## Generate code
	$(MAKE) generate-go
	$(MAKE) generate-manifests

.PHONY: generate-go
generate-go: $(CONTROLLER_GEN) $(CONVERSION_GEN) $(MOCKGEN) ## Runs Go related generate targets
	$(CONTROLLER_GEN) \
		paths=./api/... \
		object:headerFile=./hack/boilerplate/boilerplate.generatego.txt

	$(CONVERSION_GEN) \
		--input-dirs=./api/v1alpha2 \
		--output-file-base=zz_generated.conversion \
		--go-header-file=./hack/boilerplate/boilerplate.generatego.txt
	go generate ./...

.PHONY: generate-manifests
generate-manifests: $(CONTROLLER_GEN) ## Generate manifests e.g. CRD, RBAC etc.
	$(CONTROLLER_GEN) \
		paths=./api/... \
		crd:crdVersions=v1 \
		output:crd:dir=$(CRD_ROOT) \
		output:webhook:dir=$(WEBHOOK_ROOT) \
		webhook
	$(CONTROLLER_GEN) \
		paths=./controllers/... \
		output:rbac:dir=$(RBAC_ROOT) \
		rbac:roleName=manager-role

## --------------------------------------
## Docker
## --------------------------------------

.PHONY: docker-build
docker-build: ## Build the docker image for controller-manager
	docker build --pull --build-arg ARCH=$(ARCH) --build-arg LDFLAGS="$(LDFLAGS)" . -t $(CONTROLLER_IMG)-$(ARCH):$(TAG)
	MANIFEST_IMG=$(CONTROLLER_IMG)-$(ARCH) MANIFEST_TAG=$(TAG) $(MAKE) set-manifest-image
	$(MAKE) set-manifest-pull-policy

.PHONY: docker-push
docker-push: ## Push the docker image
	docker push $(CONTROLLER_IMG)-$(ARCH):$(TAG)

## --------------------------------------
## Docker — All ARCH
## --------------------------------------

.PHONY: docker-build-all ## Build all the architecture docker images
docker-build-all: $(addprefix docker-build-,$(ALL_ARCH))

docker-build-%:
	$(MAKE) ARCH=$* docker-build

.PHONY: docker-push-all ## Push all the architecture docker images
docker-push-all: $(addprefix docker-push-,$(ALL_ARCH))
	$(MAKE) docker-push-manifest

docker-push-%:
	$(MAKE) ARCH=$* docker-push

.PHONY: docker-push-manifest
docker-push-manifest: ## Push the fat manifest docker image.
	## Minimum docker version 18.06.0 is required for creating and pushing manifest images.
	docker manifest create --amend $(CONTROLLER_IMG):$(TAG) $(shell echo $(ALL_ARCH) | sed -e "s~[^ ]*~$(CONTROLLER_IMG)\-&:$(TAG)~g")
	@for arch in $(ALL_ARCH); do docker manifest annotate --arch $${arch} ${CONTROLLER_IMG}:${TAG} ${CONTROLLER_IMG}-$${arch}:${TAG}; done
	docker manifest push --purge ${CONTROLLER_IMG}:${TAG}
	MANIFEST_IMG=$(CONTROLLER_IMG) MANIFEST_TAG=$(TAG) $(MAKE) set-manifest-image
	$(MAKE) set-manifest-pull-policy

.PHONY: set-manifest-image
set-manifest-image:
	$(info Updating kustomize image patch file for manager resource)
	sed -i'' -e 's@image: .*@image: '"${MANIFEST_IMG}:$(MANIFEST_TAG)"'@' ./config/manager/manager_image_patch.yaml


.PHONY: set-manifest-pull-policy
set-manifest-pull-policy:
	$(info Updating kustomize pull policy file for manager resource)
	sed -i'' -e 's@imagePullPolicy: .*@imagePullPolicy: '"$(PULL_POLICY)"'@' ./config/manager/manager_pull_policy.yaml

## --------------------------------------
## Release
## --------------------------------------

RELEASE_TAG := $(shell git describe --abbrev=0 2>/dev/null)
RELEASE_DIR := out

$(RELEASE_DIR):
	mkdir -p $(RELEASE_DIR)/

.PHONY: release
release: clean-release  ## Builds and push container images using the latest git tag for the commit.
	@if [ -z "${RELEASE_TAG}" ]; then echo "RELEASE_TAG is not set"; exit 1; fi
	@if ! [ -z "$$(git status --porcelain)" ]; then echo "Your local git repository contains uncommitted changes, use git clean before proceeding."; exit 1; fi
	git checkout "${RELEASE_TAG}"
	# Set the manifest image to the production bucket.
	$(MAKE) set-manifest-image MANIFEST_IMG=$(PROD_REGISTRY)/$(IMAGE_NAME) MANIFEST_TAG=$(RELEASE_TAG)
	$(MAKE) set-manifest-pull-policy PULL_POLICY=IfNotPresent
	$(MAKE) release-manifests
	$(MAKE) release-binaries

.PHONY: release-manifests
release-manifests: $(RELEASE_DIR) ## Builds the manifests to publish with a release
	kustomize build config > $(RELEASE_DIR)/infrastructure-components.yaml

.PHONY: release-binaries
release-binaries: ## Builds the binaries to publish with a release
	RELEASE_BINARY=./cmd/clusterawsadm GOOS=linux GOARCH=amd64 $(MAKE) release-binary
	RELEASE_BINARY=./cmd/clusterawsadm GOOS=darwin GOARCH=amd64 $(MAKE) release-binary

.PHONY: release-binary
release-binary: $(RELEASE_DIR)
	docker run \
		--rm \
		-e CGO_ENABLED=0 \
		-e GOOS=$(GOOS) \
		-e GOARCH=$(GOARCH) \
		-v "$$(pwd):/workspace$(DOCKER_VOL_OPTS)" \
		-w /workspace \
		golang:$(GOLANG_VERSION) \
		go build -a -ldflags '$(LDFLAGS) -extldflags "-static"' \
		-o $(RELEASE_DIR)/$(notdir $(RELEASE_BINARY))-$(GOOS)-$(GOARCH) $(RELEASE_BINARY)

.PHONY: release-staging
release-staging: ## Builds and push container images to the staging bucket.
	REGISTRY=$(STAGING_REGISTRY) $(MAKE) docker-build-all docker-push-all release-alias-tag

RELEASE_ALIAS_TAG=$(PULL_BASE_REF)

.PHONY: release-alias-tag
release-alias-tag: # Adds the tag to the last build tag.
	gcloud container images add-tag $(CONTROLLER_IMG):$(TAG) $(CONTROLLER_IMG):$(RELEASE_ALIAS_TAG)

.PHONY: release-notes
release-notes: $(RELEASE_NOTES)
	$(RELEASE_NOTES) $(ARGS)

## --------------------------------------
## Development
## --------------------------------------

# This is used in the get-kubeconfig call below in the create-cluster target. It may be overridden by the
# e2e-conformance.sh script, which is why we need it as a variable here.
CLUSTER_NAME ?= test1

# NOTE: do not add 'generate-exmaples' as a prerequisite of this target. It will break e2e conformance testing.
.PHONY: create-cluster
create-cluster: $(CLUSTERCTL) $(ENVSUBST) ## Create a development Kubernetes cluster on AWS in a KIND management cluster.
	kind create cluster --name=clusterapi
	@if [ ! -z "${LOAD_IMAGE}" ]; then \
		echo "loading ${LOAD_IMAGE} into kind cluster ..." && \
		kind --name="clusterapi" load docker-image "${LOAD_IMAGE}"; \
	fi
	# Install cert manager and wait for availability
	kubectl create -f https://github.com/jetstack/cert-manager/releases/download/v0.11.1/cert-manager.yaml
	kubectl wait --for=condition=Available --timeout=5m apiservice v1beta1.webhook.cert-manager.io

	# Deploy CAPI
	kubectl apply -f https://github.com/kubernetes-sigs/cluster-api/releases/download/v0.3.0/cluster-api-components.yaml

	# Deploy CAPA
	kustomize build config | $(ENVSUBST) | kubectl apply -f -

	# Wait for CAPI pods
	kubectl wait --for=condition=Ready --timeout=5m -n capi-system pod -l cluster.x-k8s.io/provider=cluster-api
	kubectl wait --for=condition=Ready --timeout=5m -n capi-kubeadm-bootstrap-system pod -l cluster.x-k8s.io/provider=bootstrap-kubeadm
	kubectl wait --for=condition=Ready --timeout=5m -n capi-kubeadm-control-plane-system pod -l cluster.x-k8s.io/provider=control-plane-kubeadm

	# Wait for CAPA pods
	kubectl wait --for=condition=Ready --timeout=5m -n capa-system pod -l cluster.x-k8s.io/provider=infrastructure-aws

	# Create Cluster.
	sleep 10
	kustomize build templates | $(ENVSUBST) | kubectl apply -f -

	# Wait for the kubeconfig to become available.
	timeout 300 bash -c "while ! kubectl get secrets | grep $(CLUSTER_NAME)-kubeconfig; do sleep 1; done"
	# Get kubeconfig and store it locally.
	kubectl get secrets $(CLUSTER_NAME)-kubeconfig -o json | jq -r .data.value | base64 --decode > ./kubeconfig
	timeout 300 bash -c "while ! kubectl --kubeconfig=./kubeconfig get nodes | grep master; do sleep 1; done"

	# Deploy calico
	kubectl --kubeconfig=./kubeconfig apply -f https://docs.projectcalico.org/manifests/calico.yaml

.PHONY: kind-reset
kind-reset: ## Destroys the "clusterapi" kind cluster.
	kind delete cluster --name=clusterapi || true

## --------------------------------------
## Cleanup / Verification
## --------------------------------------

.PHONY: clean
clean: ## Remove all generated files
	$(MAKE) clean-bin
	$(MAKE) clean-temporary

.PHONY: clean-bin
clean-bin: ## Remove all generated binaries
	rm -rf bin
	rm -rf hack/tools/bin

.PHONY: clean-temporary
clean-temporary: ## Remove all temporary files and folders
	rm -f minikube.kubeconfig
	rm -f kubeconfig

.PHONY: clean-release
clean-release: ## Remove the release folder
	rm -rf $(RELEASE_DIR)

.PHONY: verify
verify: verify-boilerplate verify-modules verify-gen

.PHONY: verify-boilerplate
verify-boilerplate:
	./hack/verify-boilerplate.sh

.PHONY: verify-modules
verify-modules: modules
	@if !(git diff --quiet HEAD -- go.sum go.mod hack/tools/go.mod hack/tools/go.sum); then \
		git diff; \
		echo "go module files are out of date"; exit 1; \
	fi

verify-gen: generate
	@if !(git diff --quiet HEAD); then \
		git diff; \
		echo "generated files are out of date, run make generate"; exit 1; \
	fi
