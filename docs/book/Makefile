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

# Directories.
MDBOOK := $(TOOLS_BIN_DIR)/mdbook
MDBOOK_EMBED := $(TOOLS_BIN_DIR)/mdbook-embed
MDBOOK_RELEASELINK := $(TOOLS_BIN_DIR)/mdbook-releaselink
MDBOOK_TABULATE := $(TOOLS_BIN_DIR)/mdbook-tabulate
GENCRDAPIREFERENCEDOCS := $(TOOLS_BIN_DIR)/gen-crd-api-reference-docs
GENCRDAPIREFERENCEDOCS_SRCS := $(call rwildcard,.,gen-crd-api-reference-docs/*.*)
CRD_DOCS := src/crd/cluster-api-aws.md src/crd/experimental.md src/crd/eks-control-plane.md
BOOK_SRCS :=  $(filter-out $(CRD_DOCS), $(call rwildcard,.,*.*))
API_SRCS := $(call rwildcard,.,../../cmd/clusterawsadm/api) $(call rwildcard,.,../../api) $(call rwildcard,.,../../exp/api) $(call rwildcard,.,../../controlplane/eks/api)

OS := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)
PATH := $(abspath $(TOOLS_BIN_DIR)):$(PATH)
export PATH

src/crd:
	mkdir -p src/crd

src/crd/cluster-api-aws.md:
	$(MAKE) gen_crd_docs API_DIR="../../api" OUT_FILE=$@

src/crd/experimental.md:
	$(MAKE) gen_crd_docs API_DIR="../../exp/api" OUT_FILE=$@

src/crd/eks-control-plane.md:
	$(MAKE) gen_crd_docs API_DIR="../../controlplane/eks/api/" OUT_FILE=$@

src/crd/clusterawsadm.md:
	$(MAKE) gen_crd_docs API_DIR="../../cmd/clusterawsadm/api" OUT_FILE=$@

.PHONY: gen_crd_docs
gen_crd_docs: $(GENCRDAPIREFERENCEDOCS) src/crd $(API_SRCS) $(GENCRDAPIREFERENCEDOCS_SRCS)
	$(GENCRDAPIREFERENCEDOCS) -template-dir gen-crd-api-reference-docs/template -config "gen-crd-api-reference-docs/config.json" \
		-api-dir $(API_DIR) \
		-out-file $(OUT_FILE)

src/clusterawsadm:
	mkdir -p src/clusterawsadm

terraform/out:
	mkdir -p terraform/out


terraform/out/amilist: terraform/out
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $@ ./cmd/amilist && chmod 755 $@

.PHONY: generate
generate: ## Generate any dependent files outside of a build
	$(MAKE) -C ../ diagrams

.PHONY: release-policies
release-policies:
	$(MAKE) -C ../.. release-policies

src/SUMMARY_CLUSTERAWSADM.md: src/clusterawsadm $(CLUSTERAWSADMCMDREF_SRCS)
	go run ./cmd/clusterawsadmdocs > src/SUMMARY_CLUSTERAWSADM.md

# SUMMARY.md doesn't support processor statements, so construct it using cat
src/SUMMARY.md: src/SUMMARY_PREFIX.md src/SUMMARY_SUFFIX.md src/SUMMARY_CLUSTERAWSADM.md
	cat src/SUMMARY_PREFIX.md > src/SUMMARY.md
	cat src/SUMMARY_CLUSTERAWSADM.md >> src/SUMMARY.md
	cat src/SUMMARY_SUFFIX.md >> src/SUMMARY.md

verify: generate ## Verify that dependent artifacts are up to date
	@if !(git diff --quiet HEAD); then \
		git diff; \
		echo "generated files are out of date, run make generate"; exit 1; \
	fi

BOOK_DEPS := $(BOOK_SRCS) $(MDBOOK) $(MDBOOK_EMBED) $(MDBOOK_RELEASELINK) $(MDBOOK_TABULATE) release-policies src/SUMMARY.md $(CRD_DOCS)

.PHONY: build
build: $(BOOK_DEPS) ## Build the book
	$(MDBOOK) build

.PHONY: serve
serve: $(BOOK_DEPS) ## Run a local webserver with the compiled book
	$(MDBOOK) serve

.PHONY: clean
clean:
	rm -rf book
