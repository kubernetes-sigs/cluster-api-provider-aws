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

API_DIRS := cmd/clusterawsadm/api/bootstrap/v1beta1 api/v1beta1 exp/api/v1beta1 controlplane/eks/api/v1beta1 bootstrap/eks/api/v1beta1 iam/api/v1beta1
API_SRCS := $(foreach dir, $(API_DIRS), $(call rwildcard,../../$(dir),*.go))
GENERATED_SRCS := $(foreach dir, $(API_DIRS),../../$(dir)/zz_generated.deepcopy.go ../../$(dir)/zz_generated.conversion.go ../../$(dir)/zz_generated.defaults.go)
API_FILTERED_SRCS := $(filter-out $(GENERATED_SRCS), $(API_SRCS))

OS := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)
PATH := $(abspath $(TOOLS_BIN_DIR)):$(PATH)
export PATH

src/crd:
	mkdir -p src/crd

src/crd/index.md: $(API_FILTERED_SRCS) src/crd
	$(MAKE) gen_crd_docs API_DIR="../.." OUT_FILE=$@

.PHONY: gen_crd_docs
gen_crd_docs: $(GENCRDAPIREFERENCEDOCS) src/crd $(GENCRDAPIREFERENCEDOCS_SRCS)
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

BOOK_DEPS := $(BOOK_SRCS) $(MDBOOK) $(MDBOOK_EMBED) $(MDBOOK_RELEASELINK) $(MDBOOK_TABULATE) release-policies src/SUMMARY.md src/crd/index.md

.PHONY: build
build: $(BOOK_DEPS) ## Build the book
	$(MDBOOK) build

.PHONY: serve
serve: $(BOOK_DEPS) ## Run a local webserver with the compiled book
	$(MDBOOK) serve

.PHONY: clean
clean:
	rm -rf book
