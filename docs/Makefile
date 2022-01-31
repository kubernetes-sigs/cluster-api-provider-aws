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

# If you update this file, please follow
# https://suva.sh/posts/well-documented-makefiles

ROOT_DIR_RELATIVE := ..
include $(ROOT_DIR_RELATIVE)/common.mk

DIAGRAM_SRCS := $(call rwildcard,.,*.md,*.plantuml)

diagrams: $(TOOLS_BIN_DIR)/plantuml-sentinal ## Generate plantuml diagrams
	-docker run -u $(UID):$(GID) -v $(abspath .):/docs$(DOCKER_VOL_OPTS) plantuml-builder -tsvg /docs/**/*.md
	-docker run -u $(UID):$(GID) -v $(abspath .):/docs$(DOCKER_VOL_OPTS) plantuml-builder -tsvg /docs/**/*.plantuml
