# Copyright 2019 The Kubernetes Authors.
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

HOSTOS := $(shell go env GOHOSTOS)
HOSTARCH := $(shell go env GOARCH)
BINARYPATHPATTERN :=${HOSTOS}-${HOSTARCH}
VERSION := v0.3.0

PLUGINS_GOSS_URL := https://github.com/YaleUniversity/packer-provisioner-goss/releases/download/$(VERSION)/packer-provisioner-goss-$(VERSION)-$(BINARYPATHPATTERN)
PLUGINS_GOSS_SHA256-linux-amd64 := 28be39d0ddf9ad9c14e432818261abed2f2bd83257cfba213e19d5c59b710d03
PLUGINS_GOSS_SHA256-darwin-amd64 := 7ae43b5dbd26a166c8673fc7299e91d1c2244c7d2b3b558ce04e2e53acfa6f88


${HOME}/.packer.d/plugins/packer-provisioner-goss:
	mkdir -p ${HOME}/.packer.d
	mkdir -p ${HOME}/.packer.d/plugins
	mkdir -p /tmp/packer-plugins
	curl -L -o /tmp/packer-plugins/packer-provisioner-goss ${PLUGINS_GOSS_URL}
	echo ${PLUGINS_GOSS_SHA256-$(BINARYPATHPATTERN)}  */tmp/packer-plugins/packer-provisioner-goss > /tmp/packer-plugins/SHA256SUMS
	shasum -c /tmp/packer-plugins/SHA256SUMS
	mv /tmp/packer-plugins/packer-provisioner-goss ${HOME}/.packer.d/plugins
	chmod +x ${HOME}/.packer.d/plugins/packer-provisioner-goss

.PHONY: plugins
plugins: ${HOME}/.packer.d/plugins/packer-provisioner-goss
