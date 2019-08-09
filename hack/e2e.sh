#!/bin/bash

set -euo pipefail

GO111MODULE = off
export GO111MODULE
GOFLAGS =
export GOFLAGS

GOPATH="$(mktemp -d)"
export GOPATH

ACTUATOR_PKG="github.com/openshift/cluster-api-actuator-pkg"

go get -u -d "${ACTUATOR_PKG}/..."

exec make -C "${GOPATH}/src/${ACTUATOR_PKG}" test-e2e
