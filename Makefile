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
  DOCKER_CMD := docker run --rm -e CGO_ENABLED=1 -v "$(PWD)":/go/src/sigs.k8s.io/cluster-api-provider-aws:Z -w /go/src/sigs.k8s.io/cluster-api-provider-aws openshift/origin-release:golang-1.10
  IMAGE_BUILD_CMD = docker build
endif

.PHONY: depend
depend:
	dep version || go get -u github.com/golang/dep/cmd/dep
	dep ensure

.PHONY: vendor
vendor:
	dep version || go get -u github.com/golang/dep/cmd/dep
	dep ensure -v

.PHONY: generate
generate:
	go install $(GOGCFLAGS) -ldflags '-extldflags "-static"' sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/golang/mock/mockgen
	go generate ./pkg/... ./cmd/...

.PHONY: test
test: unit

bin:
	@mkdir $@

.PHONY: build
build: ## build binaries
	$(DOCKER_CMD) go build $(GOGCFLAGS) -o "bin/machine-controller-manager" \
               -ldflags "$(LD_FLAGS)" "$(REPO_PATH)/cmd/manager"
	$(DOCKER_CMD) go build $(GOGCFLAGS) -o bin/manager -ldflags '-extldflags "-static"' \
               "$(REPO_PATH)/vendor/github.com/openshift/cluster-api/cmd/manager"

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

.PHONY: integration
integration: ## Run integration test
	$(DOCKER_CMD) go test -v sigs.k8s.io/cluster-api-provider-aws/test/integration

.PHONY: build-e2e
build-e2e:
	go test -c -o bin/e2e.test sigs.k8s.io/cluster-api-provider-aws/test/machines

.PHONY: k8s-e2e
k8s-e2e: ## Run k8s specific e2e test
	# KUBECONFIG and SSH_PK dirs needs to be mounted inside a container if tests are run in containers
	go test -timeout 30m \
		-v sigs.k8s.io/cluster-api-provider-aws/test/machines \
		-kubeconfig $${KUBECONFIG:-~/.kube/config} \
		-ssh-key $${SSH_PK:-~/.ssh/id_rsa} \
		-machine-controller-image $${ACTUATOR_IMAGE:-gcr.io/k8s-cluster-api/aws-machine-controller:0.0.1} \
		-machine-manager-image $${ACTUATOR_IMAGE:-gcr.io/k8s-cluster-api/aws-machine-controller:0.0.1} \
		-nodelink-controller-image $$(docker run registry.svc.ci.openshift.org/openshift/origin-release:v4.0 image machine-api-operator) \
		-cluster-id $${ENVIRONMENT_ID:-""} \
		-ginkgo.v \
		-args -v 5 -logtostderr true

.PHONY: test-e2e
test-e2e: ## Run e2e tests
	hack/e2e.sh


.PHONY: lint
lint: ## Go lint your code
	hack/go-lint.sh -min_confidence 0.3 $$(go list -f '{{ .ImportPath }}' ./... | grep -v -e 'sigs.k8s.io/cluster-api-provider-aws/test' -e 'sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/client/mock')

.PHONY: fmt
fmt: ## Go fmt your code
	hack/go-fmt.sh .

.PHONY: vet
vet: ## Apply go vet to all go files
	hack/go-vet.sh ./...

.PHONY: help
help:
	@grep -E '^[a-zA-Z/0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
