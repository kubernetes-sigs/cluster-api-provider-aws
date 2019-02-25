#!/usr/bin/env bash

set -eu

# Print errors to stderr
function print_error {
	echo "ERROR: $1" >&2
}

function print_info {
	echo "INFO: $1" >&2
}

# Warn when unprivileged
if [ `id --user` -ne 0 ]; then
	print_error "Note: Building unprivileged may fail due to permissions"
fi

if [ -z ${VERSION_OVERRIDE+a} ]; then
        print_info "Using version from git..."
        VERSION_OVERRIDE=$(git describe --abbrev=8 --dirty --always)
fi

set -x
podman build -t "cluster-version-operator:${VERSION_OVERRIDE}" -f Dockerfile --no-cache