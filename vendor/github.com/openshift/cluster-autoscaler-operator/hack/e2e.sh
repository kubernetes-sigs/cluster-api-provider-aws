#!/bin/bash

set -euo pipefail

GOPATH="$(mktemp -d)"
export GOPATH

ACTUATOR_PKG="github.com/openshift/cluster-api-actuator-pkg"

go get -u -d "${ACTUATOR_PKG}/..."

exec make --directory="${GOPATH}/src/${ACTUATOR_PKG}" test-e2e
