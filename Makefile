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

.PHONY: all
all: generate build images

.PHONY: depend
depend:
	dep version || go get -u github.com/golang/dep/cmd/dep
	dep ensure

.PHONY: vendor
vendor:
	dep version || go get -u github.com/golang/dep/cmd/dep
	dep ensure -v -update
	patch -p1 < 0001-Delete-annotated-machines-first-when-scaling-down.patch
	patch -p1 < 0002-Sort-machines-before-syncing.patch

.PHONY: generate
generate: gendeepcopy

.PHONY: gendeepcopy
gendeepcopy:
	go build -o $$GOPATH/bin/deepcopy-gen sigs.k8s.io/cluster-api-provider-aws/vendor/k8s.io/code-generator/cmd/deepcopy-gen
	deepcopy-gen \
	  -i ./cloud/aws/providerconfig,./cloud/aws/providerconfig/v1alpha1 \
	  -O zz_generated.deepcopy \
	  -h boilerplate.go.txt

build:
	CGO_ENABLED=0 go install -a -ldflags '-extldflags "-static"' sigs.k8s.io/cluster-api-provider-aws/cmd/cluster-controller
	CGO_ENABLED=0 go install -a -ldflags '-extldflags "-static"' sigs.k8s.io/cluster-api-provider-aws/cmd/machine-controller

aws-actuator:
	go build -o bin/aws-actuator sigs.k8s.io/cluster-api-provider-aws/cmd/aws-actuator

.PHONY: images
images: ## Create images
	$(MAKE) -C cmd/cluster-controller image
	$(MAKE) -C cmd/machine-controller image

.PHONY: push
push:
	$(MAKE) -C cmd/cluster-controller push
	$(MAKE) -C cmd/machine-controller push

.PHONY: check
check: fmt vet lint test ## Check your code

.PHONY: test
test: # Run unit test
	go test -race -cover ./cmd/... ./cloud/...

.PHONY: integration
integration: ## Run integration test
	go test -v sigs.k8s.io/cluster-api-provider-aws/test/integration

.PHONY: lint
lint: ## Go lint your code
	hack/go-lint.sh -min_confidence 0.3 $(go list -f '{{ .ImportPath }}' ./...)

.PHONY: fmt
fmt: ## Go fmt your code
	hack/verify-gofmt.sh

.PHONY: vet
vet: ## Apply go vet to all go files
	hack/go-vet.sh ./...

.PHONY: help
help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
