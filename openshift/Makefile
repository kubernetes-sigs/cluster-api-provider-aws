BIN_DIR := bin
TOOLS_DIR := tools

$(RELEASE_DIR):
	mkdir -p $(RELEASE_DIR)/

KUSTOMIZE := tools/bin/kustomize
MANIFESTS_GEN := go run ./vendor/github.com/openshift/cluster-capi-operator/manifests-gen/

$(KUSTOMIZE):
	./tools/ensure-kustomize.sh

.PHONY: check-env
check-env:
ifndef PROVIDER_VERSION
	$(error PROVIDER_VERSION is undefined)
endif

.PHONY: ocp-manifests
ocp-manifests: $(RELEASE_DIR) $(KUSTOMIZE) check-env ## Builds openshift specific manifests
	# Build core-components.
	$(KUSTOMIZE) build ../config/default > infrastructure-components.yaml
	# Generate provider manifests.
	# TODO: load the provider-version dynamically at rebase time when this is invoked by the Rebase Bot during one of its lifecycle hooks.
	cd tools && $(MANIFESTS_GEN) --provider-name "aws" --provider-type "InfrastructureProvider" --provider-version "${PROVIDER_VERSION}" --base-path "../../" --manifests-path "../manifests"
