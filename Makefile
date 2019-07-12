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
# A release does not need to define this
MANAGER_IMAGE_NAME ?= cluster-api-aws-controller
# A release should define this with the next version after 0.0.4
MANAGER_IMAGE_TAG ?= dev
# A release should define this with IfNotPresent
PULL_POLICY ?= Always

# Used in docker-* targets.
MANAGER_IMAGE ?= $(REGISTRY)/$(MANAGER_IMAGE_NAME):$(MANAGER_IMAGE_TAG)

## Image URL to use all building/pushing image targets
BAZEL_ARGS ?=
BAZEL_BUILD_ARGS := --define=REGISTRY=$(REGISTRY)\
 --define=PULL_POLICY=$(PULL_POLICY)\
 --define=MANAGER_IMAGE_NAME=$(MANAGER_IMAGE_NAME)\
 --define=MANAGER_IMAGE_TAG=$(MANAGER_IMAGE_TAG)\
 --host_force_python=PY2\
$(BAZEL_ARGS)

# Bazel variables
BAZEL_VERSION := $(shell command -v bazel 2> /dev/null)
ifndef BAZEL_VERSION
    $(error "Bazel is not available. \
		Installation instructions can be found at \
		https://docs.bazel.build/versions/master/install.html")
endif

.PHONY: all
all: verify-install test binaries

help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


## --------------------------------------
## Testing
## --------------------------------------

.PHONY: test
test: generate lint ## Run tests
	$(MAKE) test-go

.PHONY: test-go
test-go: ## Run tests
	go test -v -tags=integration ./pkg/... ./cmd/...

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
docker-build: generate lint-full ## Build the docker image for controller-manager
	docker build --pull . -t $(MANAGER_IMAGE)

.PHONY: docker-push
docker-push: docker-build ## Push the docker image
	docker push $(MANAGER_IMAGE)

## --------------------------------------
## Manifests
## --------------------------------------

.PHONY: manifests
manifests: cmd/clusterctl/examples/aws/provider-components-base.yaml ## Build example set of manifests from the current source
	./cmd/clusterctl/examples/aws/generate-yaml.sh

.PHONY: cmd/clusterctl/examples/aws/provider-components-base.yaml
cmd/clusterctl/examples/aws/provider-components-base.yaml:
	bazel build //cmd/clusterctl/examples/aws:provider-components-base $(BAZEL_BUILD_ARGS)
	install bazel-genfiles/cmd/clusterctl/examples/aws/provider-components-base.yaml cmd/clusterctl/examples/aws

## --------------------------------------
## Generate
## --------------------------------------

.PHONY: vendor
vendor: ## Runs go mod to ensure proper vendoring.
	./hack/update-vendor.sh
	$(MAKE) gazelle

.PHONY: gazelle
gazelle: ## Run Bazel Gazelle
	(which bazel && ./hack/update-bazel.sh) || true

.PHONY: generate
generate: ## Generate code
	$(MAKE) generate-go
	$(MAKE) generate-mocks
	$(MAKE) generate-manifests
	$(MAKE) gazelle

.PHONY: generate-go
generate-go: ## Runs go generate
	go generate ./pkg/... ./cmd/...

.PHONY: generate-mocks
generate-mocks: ## Generate mocks, CRDs and runs `go generate` through Bazel
	bazel build $(BAZEL_ARGS) //pkg/cloud/aws/services/mocks:mocks \
		//pkg/cloud/aws/services/ec2/mock_ec2iface:mocks \
		//pkg/cloud/aws/services/elb/mock_elbiface:mocks
	./hack/copy-bazel-mocks.sh

.PHONY: generate-manifests
generate-manifests: ## Generate manifests e.g. CRD, RBAC etc.
	go run vendor/sigs.k8s.io/controller-tools/cmd/controller-gen/main.go \
		paths=./pkg/apis/infrastructure/... \
		crd:trivialVersions=true \
		output:crd:dir=./config/crds
	go run vendor/sigs.k8s.io/controller-tools/cmd/controller-gen/main.go \
		paths=./pkg/controller/... \
		rbac:roleName=manager-role

## --------------------------------------
## Linting
## --------------------------------------

.PHONY: lint
lint: ## Lint codebase
	bazel run //:lint $(BAZEL_ARGS)

lint-full: ## Run slower linters to detect possible issues
	bazel run //:lint-full $(BAZEL_ARGS)

## --------------------------------------
## Binaries
## --------------------------------------

.PHONY: binaries
binaries: manager clusterawsadm clusterctl ## Builds and installs all binaries

.PHONY: manager
manager: ## Build manager binary.
	go build -o bin/manager sigs.k8s.io/cluster-api-provider-aws/cmd/manager

.PHONY: clusterctl
clusterctl: ## Build clusterctl binary.
	go build -o bin/clusterctl sigs.k8s.io/cluster-api-provider-aws/cmd/clusterctl

.PHONY: clusterawsadm
clusterawsadm: ## Build clusterawsadm binary.
	go build -o bin/clusterawsadm sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm

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

.PHONY: create-cluster
create-cluster: binaries ## Create a development Kubernetes cluster on AWS using examples
	bin/clusterctl create cluster -v 4 \
	--provider aws \
	--bootstrap-type kind \
	-m ./cmd/clusterctl/examples/aws/out/machines.yaml \
	-c ./cmd/clusterctl/examples/aws/out/cluster.yaml \
	-p ./cmd/clusterctl/examples/aws/out/provider-components.yaml \
	-a ./cmd/clusterctl/examples/aws/out/addons.yaml

.PHONY: create-cluster-ha
create-cluster-ha: binaries ## Create a development Kubernetes cluster on AWS using HA examples
	bin/clusterctl create cluster -v 4 \
	--provider aws \
	--bootstrap-type kind \
	-m ./cmd/clusterctl/examples/aws/out/machines-ha.yaml \
	-c ./cmd/clusterctl/examples/aws/out/cluster.yaml \
	-p ./cmd/clusterctl/examples/aws/out/provider-components.yaml \
	-a ./cmd/clusterctl/examples/aws/out/addons.yaml

.PHONY: create-cluster-management
create-cluster-management: ## Create a development Kubernetes cluster on AWS in a KIND management cluster.
	kind create cluster --name=clusterapi
	# Apply provider-components.
	kubectl \
		--kubeconfig=$$(kind get kubeconfig-path --name="clusterapi") \
		create -f cmd/clusterctl/examples/aws/out/provider-components.yaml
	# Create Cluster.
	kubectl \
		--kubeconfig=$$(kind get kubeconfig-path --name="clusterapi") \
		create -f cmd/clusterctl/examples/aws/out/cluster.yaml
	# Create control plane machine.
	kubectl \
		--kubeconfig=$$(kind get kubeconfig-path --name="clusterapi") \
		create -f cmd/clusterctl/examples/aws/out/controlplane-machine.yaml
	# Get KubeConfig using clusterctl.
	bin/clusterctl alpha phases get-kubeconfig -v=3 \
		--kubeconfig=$$(kind get kubeconfig-path --name="clusterapi") \
		--provider=aws \
		--cluster-name=test1
	# Apply addons on the target cluster, waiting for the control-plane to become available.
	bin/clusterctl alpha phases apply-addons -v=3 \
		--kubeconfig=./kubeconfig \
		-a cmd/clusterctl/examples/aws/out/addons.yaml
	# Create a worker node with MachineDeployment.
	kubectl \
		--kubeconfig=$$(kind get kubeconfig-path --name="clusterapi") \
		create -f cmd/clusterctl/examples/aws/out/machine-deployment.yaml

.PHONY: delete-cluster
delete-cluster: binaries ## Deletes the development Kubernetes Cluster "test1"
	bin/clusterctl delete cluster -v 4 \
	--bootstrap-type kind \
	--cluster test1 \
	--kubeconfig ./kubeconfig \
	-p ./cmd/clusterctl/examples/aws/out/provider-components.yaml \

kind-reset: ## Destroys the "clusterapi" kind cluster.
	kind delete cluster --name=clusterapi || true

## --------------------------------------
## Cleanup / Verification
## --------------------------------------

.PHONY: clean
clean: ## Remove all generated files
	$(MAKE) clean-bazel
	$(MAKE) clean-bin
	$(MAKE) clean-temporary

.PHONY: clean-bazel
clean-bazel: ## Remove all generated bazel symlinks
	bazel clean

.PHONY: clean-bin
clean-bin: ## Remove all generated binaries
	rm -rf bin

.PHONY: clean-temporary
clean-temporary: ## Remove all temporary files and folders
	rm -f minikube.kubeconfig
	rm -f kubeconfig
	rm -rf out/
	rm -rf cmd/clusterctl/examples/aws/out/
	rm -f cmd/clusterctl/examples/aws/provider-components-base.yaml

.PHONY: verify
verify: ## Runs verification scripts to ensure correct execution
	./hack/verify-boilerplate.sh
	./hack/verify-bazel.sh

.PHONY: verify-install
verify-install: ## Checks that you've installed this repository correctly
	./hack/verify-install.sh
