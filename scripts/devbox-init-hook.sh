#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
SRC_DIR="$(dirname $DIR)"

# Check if envsubst.old exists
if [[ ! -f "/usr/bin/envsubst.old" ]]; then
    echo "envsubst.old does not exist. Proceeding with installation..."

    # Set install path
    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p "$INSTALL_DIR"

    # Download and install envsubst
    URL="https://github.com/a8m/envsubst/releases/download/v1.2.0/envsubst-`uname -s`-`uname -m`"
    echo "Downloading: $URL"

    curl -sL $URL -o $HOME/.local/bin/envsubst

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
    if ! command -v envsubst &> /dev/null; then
        echo "Installation failed."
        exit 1
    fi
fi

if [[ ! -e "$HOME/.local/bin/clusterawsadm" ]]; then
    ln -s "$SRC_DIR/bin/clusterawsadm" $HOME/.local/bin/clusterawsadm
fi
