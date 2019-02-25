#!/usr/bin/env bash

set -eu


function print_info {
	echo "INFO: $1" >&2
}

REPO=${REPO:-"openshift"}

if [ -z ${VERSION_OVERRIDE+a} ]; then
        print_info "Using version from git..."
        VERSION_OVERRIDE=$(git describe --abbrev=8 --dirty --always)
fi

set -x
podman push "cluster-version-operator:${VERSION_OVERRIDE}" "${REPO}/origin-cluster-version-operator:${VERSION_OVERRIDE}"
podman push "cluster-version-operator:${VERSION_OVERRIDE}" "${REPO}/origin-cluster-version-operator:latest"
