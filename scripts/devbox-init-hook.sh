#!/bin/bash
# Copyright 2025 The Kubernetes Authors.
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

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SRC_DIR="$(dirname $DIR)"

# Set install path
INSTALL_DIR="$SRC_DIR/bin"
mkdir -p "$INSTALL_DIR"

if [[ ! -f "$INSTALL_DIR/envsubst" ]]; then
    echo "Installing github.com/a8m/envsubst into bin"
    make -C "$SRC_DIR/hack/tools" bin/envsubst
    ln -s "$SRC_DIR/hack/tools/bin/envsubst" "$INSTALL_DIR/envsubst"

    # Verify installation
    if ! command -v envsubst &>/dev/null; then
        echo "Installation failed."
        exit 1
    fi
fi

# Use build location by default
if [[ ! -f "$INSTALL_DIR/clusterawsadm" ]]; then
    echo "Installing clusterawsadm into bin"
	make -C "$SRC_DIR" clusterawsadm
fi
