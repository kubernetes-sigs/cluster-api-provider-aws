#!/bin/bash
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SRC_DIR="$(dirname $DIR)"

# Set install path to devbox packages directory
INSTALL_DIR="$DEVBOX_PACKAGES_DIR/bin"

if [[ ! -f "$INSTALL_DIR/envsubst" ]]; then
    ENVSUBST_VERSION="v1.4.3"
    URL="https://github.com/a8m/envsubst/releases/download/$ENVSUBST_VERSION/envsubst-$(uname -s)-$(uname -m)"
    echo "Downloading: $URL"
    echo "Installing github.com/a8m/envsubst into Devbox shell"
    sudo curl -sL $URL -o $INSTALL_DIR/envsubst
    sudo chmod +x "$INSTALL_DIR/envsubst"

    # Verify installation
    if ! command -v envsubst &>/dev/null; then
        echo "Installation failed."
        exit 1
    fi
fi

# Use build location by default
if [[ ! -L "$INSTALL_DIR/clusterawsadm" ]]; then
    echo "Linking [$SRC_DIR/bin/clusterawsadm] [$INSTALL_DIR/clusterawsadm]"
    sudo ln -s "$SRC_DIR/bin/clusterawsadm" $INSTALL_DIR/clusterawsadm
fi
