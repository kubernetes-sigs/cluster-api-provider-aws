#!/bin/bash
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
