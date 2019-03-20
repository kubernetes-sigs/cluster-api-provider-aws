DBG         ?= 0
#REGISTRY    ?= quay.io/openshift/
VERSION     ?= $(shell git describe --always --abbrev=7)
MUTABLE_TAG ?= latest
IMAGE        = $(REGISTRY)machine-api-operator

ifeq ($(DBG),1)
GOGCFLAGS ?= -gcflags=all="-N -l"
endif

.PHONY: all
all: check build test

NO_DOCKER ?= 0
ifeq ($(NO_DOCKER), 1)
  DOCKER_CMD =
  IMAGE_BUILD_CMD = imagebuilder
else
  DOCKER_CMD := docker run --rm -v "$(PWD)":/go/src/github.com/openshift/machine-api-operator:Z -w /go/src/github.com/openshift/machine-api-operator golang:1.10
  IMAGE_BUILD_CMD = docker build
endif

.PHONY: check
check: lint fmt vet verify-codegen test ## Run code validations

.PHONY: build
build: machine-api-operator nodelink-controller machine-healthcheck ## Build binaries

.PHONY: machine-api-operator
machine-api-operator:
	$(DOCKER_CMD) ./hack/go-build.sh machine-api-operator

.PHONY: nodelink-controller
nodelink-controller:
	$(DOCKER_CMD) ./hack/go-build.sh nodelink-controller

.PHONY: machine-healthcheck
machine-healthcheck:
	$(DOCKER_CMD) ./hack/go-build.sh machine-healthcheck

.PHONY: update-codegen
update-codegen:
	$(DOCKER_CMD) ./hack/update-codegen.sh

.PHONY: verify-codegen
verify-codegen:
	$(DOCKER_CMD) ./hack/verify-codegen.sh

.PHONY: build-integration
build-integration: ## Build integration test binary
	@echo -e "\033[32mBuilding integration test binary...\033[0m"
	mkdir -p bin
	$(DOCKER_CMD) go build $(GOGCFLAGS) -o bin/integration github.com/openshift/machine-api-operator/test/integration

.PHONY: test-e2e
test-e2e: ## Run openshift specific e2e test
	go test -timeout 60m \
		-v ./vendor/github.com/openshift/cluster-api-actuator-pkg/pkg/e2e \
		-kubeconfig $${KUBECONFIG:-~/.kube/config} \
		-machine-api-namespace $${NAMESPACE:-openshift-machine-api} \
		-ginkgo.v \
		-args -v 5 -logtostderr true

.PHONY: deploy-kubemark
deploy-kubemark:
	kustomize build config | kubectl apply -f -
	kustomize build | kubectl apply -f -
	kubectl apply -f config/kubemark-install-config.yaml

.PHONY: test
test: ## Run tests
	@echo -e "\033[32mTesting...\033[0m"
	$(DOCKER_CMD) go test ./...

.PHONY: image
image: ## Build docker image
	@echo -e "\033[32mBuilding image $(IMAGE):$(VERSION) and tagging also as $(IMAGE):$(MUTABLE_TAG)...\033[0m"
	$(IMAGE_BUILD_CMD) -t "$(IMAGE):$(VERSION)" -t "$(IMAGE):$(MUTABLE_TAG)" ./

.PHONY: push
push: ## Push image to docker registry
	@echo -e "\033[32mPushing images...\033[0m"
	docker push "$(IMAGE):$(VERSION)"
	docker push "$(IMAGE):$(MUTABLE_TAG)"

.PHONY: lint
lint: ## Go lint your code
	hack/go-lint.sh -min_confidence 0.3 $(go list -f '{{ .ImportPath }}' ./...)

.PHONY: fmt
fmt: ## Go fmt your code
	hack/go-fmt.sh .

.PHONY: vet
vet: ## Apply go vet to all go files
	hack/go-vet.sh ./...

.PHONY: help
help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
