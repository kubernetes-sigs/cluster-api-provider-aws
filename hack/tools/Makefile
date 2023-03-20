# Copyright 2020 The Kubernetes Authors.
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

ROOT_DIR_RELATIVE := ../..
include $(ROOT_DIR_RELATIVE)/common.mk

UNAME := $(shell uname -s)

# Directories.
BIN_DIR := bin
BIN_DIR_ABS := $(abspath $(BIN_DIR))
SHARE_DIR := share

OS := $(shell go env GOOS)
RUST_TARGET := unknown-$(OS)-gnu

ifeq ($(OS), darwin)
	RUST_TARGET := apple-darwin
	GH_ARCH_SUFFIX := macOS_amd64
	GTAR := gtar
endif

ifeq ($(OS), linux)
	GH_ARCH_SUFFIX := linux_amd64
	GTAR := tar
endif

MDBOOK_EXTRACT_COMMAND := tar xfvz $(SHARE_DIR)/mdbook.tar.gz -C bin
MDBOOK_ARCHIVE_EXT := .tar.gz

ifeq ($(OS), windows)
	RUST_TARGET := pc-windows-msvc
	MDBOOK_ARCHIVE_EXT := .zip
	MDBOOK_EXTRACT_COMMAND := unzip -d /tmp
endif

GOLANGCI_LINT_VERSION := v1.52.2
## --------------------------------------
## Tooling Binaries
## --------------------------------------

$(BIN_DIR):
	mkdir -p $@

$(SHARE_DIR):
	mkdir -p $@

.PHONY: $(GTAR)
$(GTAR):
	@$(GTAR) --version > /dev/null || (echo Install GNU Tar with brew install gnu-tar && exit -1)

$(BIN_DIR)/start.sh:
	curl https://raw.githubusercontent.com/windmilleng/rerun-process-wrapper/master/start.sh -o $@

$(BIN_DIR)/restart.sh:
	curl https://raw.githubusercontent.com/windmilleng/rerun-process-wrapper/master/restart.sh -o $@

GO_APIDIFF := $(BIN_DIR)/go-apidiff
$(GO_APIDIFF): $(BIN_DIR) go.mod go.sum
	go build -tags=tools -o $@ github.com/joelanford/go-apidiff

CONTROLLER_GEN := $(BIN_DIR)/controller-gen
$(CONTROLLER_GEN): $(BIN_DIR) go.mod go.sum # Build controller-gen from tools folder.
	go build -tags=tools -o $@ sigs.k8s.io/controller-tools/cmd/controller-gen

CONVERSION_GEN := $(BIN_DIR)/conversion-gen
$(CONVERSION_GEN): $(BIN_DIR) go.mod go.sum $(call rwildcard,.,third_party/conversion-gen/*.*)
	go build -tags=tools -o $@ ./third_party/conversion-gen

CONVERSION_VERIFIER := $(BIN_DIR)/conversion-verifier
$(CONVERSION_VERIFIER): $(BIN_DIR) go.mod go.sum
	go build -tags=tools -o $@ sigs.k8s.io/cluster-api/hack/tools/conversion-verifier

DEFAULTER_GEN := $(BIN_DIR)/defaulter-gen
$(DEFAULTER_GEN): $(BIN_DIR) go.mod go.sum
	go build -tags=tools -o $@ k8s.io/code-generator/cmd/defaulter-gen

ENVSUBST := $(BIN_DIR)/envsubst
$(ENVSUBST): $(BIN_DIR) go.mod go.sum # Build envsubst from tools folder.
	go build -tags=tools -o $@ github.com/a8m/envsubst/cmd/envsubst

GENCRDAPIREFERENCEDOCS := $(BIN_DIR)/gen-crd-api-reference-docs
$(GENCRDAPIREFERENCEDOCS): $(BIN_DIR) go.mod go.sum # Build gen-crd-api-reference-docs from tools folder.
	go build -tags=tools -o $@ github.com/ahmetb/gen-crd-api-reference-docs

GH_SHARE := $(SHARE_DIR)/gh

$(GH_SHARE): $(SHARE_DIR)
	mkdir -p $@

$(GH_SHARE)/gh.tar.gz: $(GH_SHARE)
	curl -L "https://github.com/cli/cli/releases/download/v$(GH_VERSION)/gh_$(GH_VERSION)_$(GH_ARCH_SUFFIX).tar.gz" -o $@

GH := $(BIN_DIR)/gh
$(GH): $(GTAR) $(GH_SHARE)/gh.tar.gz
	$(GTAR) -xvf share/gh/gh.tar.gz gh_$(GH_VERSION)_$(GH_ARCH_SUFFIX)/bin/gh --strip-components 1 --directory $(BIN_DIR)
	chmod +x $@
	touch -m $@

GOJQ := $(BIN_DIR)/gojq
$(GOJQ): $(BIN_DIR) go.mod go.sum
	go build -tags=tools -o $@ github.com/itchyny/gojq/cmd/gojq

GOLANGCI_LINT := $(BIN_DIR)/golangci-lint
$(GOLANGCI_LINT): $(BIN_DIR) # Install golangci-lint in tools folder.
	${REPO_ROOT}/hack/ensure-golangci-lint.sh -b $(BIN_DIR) ${GOLANGCI_LINT_VERSION}

KIND := $(BIN_DIR)/kind
$(KIND): $(BIN_DIR) go.mod go.sum
	go build -tags tools -o $@ sigs.k8s.io/kind

KUSTOMIZE := $(BIN_DIR)/kustomize
$(KUSTOMIZE): $(BIN_DIR) go.mod go.sum # Build kustomize from tools folder.
	CGO_ENABLED=0 go build -tags=tools -o $@ sigs.k8s.io/kustomize/kustomize/v4

MDBOOK_SHARE := $(SHARE_DIR)/mdbook$(MDBOOK_ARCHIVE_EXT)
$(MDBOOK_SHARE): ../../versions.mk $(SHARE_DIR)
	curl -sL -o $(MDBOOK_SHARE) "https://github.com/rust-lang/mdBook/releases/download/$(MDBOOK_VERSION)/mdBook-$(MDBOOK_VERSION)-x86_64-$(RUST_TARGET)$(MDBOOK_ARCHIVE_EXT)"

MDBOOK := $(BIN_DIR)/mdbook
$(MDBOOK): $(BIN_DIR) $(MDBOOK_SHARE)
	$(MDBOOK_EXTRACT_COMMAND)
	chmod +x $@
	touch -m $@

MDBOOK_EMBED := $(BIN_DIR)/mdbook-embed
$(MDBOOK_EMBED): $(BIN_DIR) go.mod go.sum
	go build -tags=tools -o $(BIN_DIR)/mdbook-embed sigs.k8s.io/cluster-api/hack/tools/mdbook/embed

MDBOOK_RELEASELINK := $(BIN_DIR)/mdbook-releaselink
$(MDBOOK_RELEASELINK): $(BIN_DIR) go.mod go.sum
	go build -tags=tools -o $(BIN_DIR)/mdbook-releaselink sigs.k8s.io/cluster-api/hack/tools/mdbook/releaselink

MDBOOK_TABULATE := $(BIN_DIR)/mdbook-tabulate
$(MDBOOK_TABULATE): $(BIN_DIR) go.mod go.sum
	go build -tags=tools -o $(BIN_DIR)/mdbook-tabulate sigs.k8s.io/cluster-api/hack/tools/mdbook/tabulate

MOCKGEN := $(BIN_DIR)/mockgen
$(MOCKGEN): $(BIN_DIR) go.mod go.sum # Build mockgen from tools folder.
	go build -tags=tools -o $@ github.com/golang/mock/mockgen

PLANTUML := $(BIN_DIR)/plantuml-sentinal
$(PLANTUML): plantuml.Dockerfile ../../versions.mk
	docker build --build-arg PLANTUML_VERSION=$(PLANTUML_VERSION)  . -f plantuml.Dockerfile -t "plantuml-builder"
	touch $@

KPROMO_BIN := kpromo
KPROMO :=  $(BIN_DIR)/$(KPROMO_BIN)
$(KPROMO):
	CGO_ENABLED=0 go build -o $@ sigs.k8s.io/promo-tools/v3/cmd/kpromo

YQ_BIN := yq
YQ :=  $(BIN_DIR)/$(YQ_BIN)
$(YQ):
	CGO_ENABLED=0 go build -tags=tools -o $@ github.com/mikefarah/yq/v4

## ------------------------------------------------------------------------------------------------
## AWS Session Manager Plugin Installation. Currently support Linux and MacOS AMD64 architectures.
## ------------------------------------------------------------------------------------------------

SSM_SHARE := $(SHARE_DIR)/ssm

$(SSM_SHARE): $(SHARE_DIR)
	mkdir -p $@

$(SSM_SHARE)/session-manager-plugin.deb: $(SSM_SHARE)
	curl "https://s3.amazonaws.com/session-manager-downloads/plugin/latest/ubuntu_64bit/session-manager-plugin.deb" -o $@

$(SSM_SHARE)/sessionmanager-bundle.zip: $(SSM_SHARE)
	curl "https://s3.amazonaws.com/session-manager-downloads/plugin/latest/mac/sessionmanager-bundle.zip" -o $@

$(SSM_SHARE)/data.tar.gz: $(SSM_SHARE)/session-manager-plugin.deb
	cd $(SSM_SHARE) && ar x session-manager-plugin.deb data.tar.gz

SSM_PLUGIN := $(BIN_DIR)/session-manager-plugin
$(SSM_PLUGIN): $(BIN_DIR)
ifeq ($(UNAME), Linux)
	$(MAKE) $(SSM_SHARE)/data.tar.gz
	cd $(BIN_DIR) && $(GTAR) -xvf ../share/ssm/data.tar.gz usr/local/sessionmanagerplugin/bin/session-manager-plugin --strip-components 4 --directory $(BIN_DIR)
endif
ifeq ($(UNAME), Darwin)
	$(MAKE) $(SSM_SHARE)/sessionmanager-bundle.zip
	cd $(BIN_DIR) && unzip -o -j ../share/ssm/sessionmanager-bundle.zip sessionmanager-bundle/bin/session-manager-plugin
endif

.PHONY: clean
clean: ## Remove all tools
	rm -rf bin
	rm -rf share
