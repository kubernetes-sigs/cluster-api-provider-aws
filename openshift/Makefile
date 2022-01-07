TOOLS_DIR := ../hack/tools
TOOLS_DIR_DEPS := $(TOOLS_DIR)/go.sum $(TOOLS_DIR)/go.mod $(TOOLS_DIR)/Makefile
TOOLS_BIN_DIR := $(TOOLS_DIR)/bin

$(TOOLS_BIN_DIR)/%: $(TOOLS_DIR_DEPS)
	make -C $(TOOLS_DIR) $(subst $(TOOLS_DIR)/,,$@)
	
KUSTOMIZE := $(TOOLS_BIN_DIR)/kustomize

.PHONY: ocp-manifests
ocp-manifests:  $(KUSTOMIZE)  ## Builds openshift specific manifests
	$(KUSTOMIZE) build ../config/default > infrastructure-components.yaml
