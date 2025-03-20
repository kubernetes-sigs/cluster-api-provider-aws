#!/bin/bash
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SRC_DIR="$(dirname $DIR)"

# Set install path
INSTALL_DIR="$HOME/.local/bin"
mkdir -p "$INSTALL_DIR"

# Check if envsubst.old exists
if [[ ! -f "/usr/bin/envsubst.old" ]]; then
    echo "envsubst.old does not exist. Proceeding with installation..."

    ENVSUBST_VERSION="v1.4.3"
    URL="https://github.com/a8m/envsubst/releases/download/$ENVSUBST_VERSION/envsubst-$(uname -s)-$(uname -m)"
    echo "Downloading: $URL"

    curl -sL $URL -o $INSTALL_DIR/envsubst
    chmod +x "$INSTALL_DIR/envsubst"

    # Replace existing envsubst if present
    if [[ -f "/usr/bin/envsubst" ]]; then
        sudo mv /usr/bin/envsubst /usr/bin/envsubst.old
        sudo ln -s "$INSTALL_DIR/envsubst" /usr/bin/envsubst
        echo "Replaced existing envsubst with the new one at /usr/bin/envsubst"
    else
        touch /usr/bin/envsubst.old
    fi

    # Verify installation
    if ! command -v envsubst &>/dev/null; then
        echo "Installation failed."
        exit 1
    fi
fi

# Use build location by default
if [[ ! -e "$INSTALL_DIR/clusterawsadm" ]]; then
    echo "Linking [$SRC_DIR/bin/clusterawsadm] [$INSTALL_DIR/clusterawsadm]"
    mkdir -p "$INSTALL_DIR"
    ln -s "$SRC_DIR/bin/clusterawsadm" $INSTALL_DIR/clusterawsadm
fi
