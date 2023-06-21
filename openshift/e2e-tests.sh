#!/bin/bash

set -euo pipefail

echo "Running e2e-tests.sh"

unset GOFLAGS
tmp="$(mktemp -d)"

git clone --depth=1 "https://github.com/openshift/cluster-capi-operator.git" "$tmp"

exec make -C "$tmp" e2e
