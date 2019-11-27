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

GO111MODULE = on
export GO111MODULE
GOFLAGS += -mod=vendor
export GOFLAGS
GOPROXY ?=
export GOPROXY

DBG ?= 0

ifeq ($(DBG),1)
GOGCFLAGS ?= -gcflags=all="-N -l"
endif

VERSION     ?= $(shell git describe --always --abbrev=7)
REPO_PATH   ?= sigs.k8s.io/cluster-api-provider-aws
LD_FLAGS    ?= -X $(REPO_PATH)/pkg/version.Raw=$(VERSION) -extldflags "-static"
MUTABLE_TAG ?= latest
IMAGE        = origin-aws-machine-controllers

.PHONY: all
all: generate build images check

NO_DOCKER ?= 0
ifeq ($(NO_DOCKER), 1)
  DOCKER_CMD =
  IMAGE_BUILD_CMD = imagebuilder
  CGO_ENABLED = 1
else
  DOCKER_CMD := docker run --rm -e CGO_ENABLED=1 -v "$(PWD)":/go/src/sigs.k8s.io/cluster-api-provider-aws:Z -w /go/src/sigs.k8s.io/cluster-api-provider-aws openshift/origin-release:golang-1.12
  IMAGE_BUILD_CMD = docker build
endif

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor
	go mod verify

.PHONY: generate
generate:
	$(DOCKER_CMD) go generate ./pkg/... ./cmd/...
	hack/goimports.sh .

.PHONY: test
test: unit

bin:
	@mkdir $@

.PHONY: build
build: ## build binaries
	$(DOCKER_CMD) go build $(GOGCFLAGS) -o "bin/machine-controller-manager" \
               -ldflags "$(LD_FLAGS)" "$(REPO_PATH)/cmd/manager"
	$(DOCKER_CMD) go build $(GOGCFLAGS) -o bin/manager -ldflags '-extldflags "-static"' \
               "$(REPO_PATH)/vendor/github.com/openshift/machine-api-operator/cmd/machineset"


aws-actuator:
	$(DOCKER_CMD) go build $(GOGCFLAGS) -o bin/aws-actuator sigs.k8s.io/cluster-api-provider-aws/cmd/aws-actuator

.PHONY: images
images: ## Create images
	$(IMAGE_BUILD_CMD) -t "$(IMAGE):$(VERSION)" -t "$(IMAGE):$(MUTABLE_TAG)" ./

.PHONY: push
push:
	docker push "$(IMAGE):$(VERSION)"
	docker push "$(IMAGE):$(MUTABLE_TAG)"

.PHONY: check
check: fmt vet lint test check-pkg ## Check your code

.PHONY: check-pkg
check-pkg:
	./hack/verify-actuator-pkg.sh

.PHONY: unit
unit: # Run unit test
	$(DOCKER_CMD) go test -race -cover ./cmd/... ./pkg/...

.PHONY: test-e2e
test-e2e: ## Run e2e tests
	hack/e2e.sh


.PHONY: lint
lint: ## Go lint your code
	hack/go-lint.sh -min_confidence 0.3 $$(go list -f '{{ .ImportPath }}' ./... | grep -v -e 'sigs.k8s.io/cluster-api-provider-aws/test' -e 'sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/client/mock')

.PHONY: fmt
fmt: ## Go fmt your code
	hack/go-fmt.sh .

.PHONY: goimports
goimports: ## Go fmt your code
	hack/goimports.sh .

.PHONY: vet
vet: ## Apply go vet to all go files
	hack/go-vet.sh ./...

.PHONY: help
help:
	@grep -E '^[a-zA-Z/0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
