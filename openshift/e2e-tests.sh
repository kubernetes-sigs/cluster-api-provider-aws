#!/bin/bash

set -euo pipefail

echo "Running e2e-tests.sh"

unset GOFLAGS
tmp="$(mktemp -d)"

# Default branch for both CCAPIO and this repo should be `main`.
CCAPIO_BASE_REF=$PULL_BASE_REF

echo "cloning github.com/openshift/cluster-capi-operator at branch '$CCAPIO_BASE_REF'"
git clone --single-branch --branch="$CCAPIO_BASE_REF" --depth=1 "https://github.com/openshift/cluster-capi-operator.git" "$tmp"

echo "running cluster-capi-operator's: make e2e"
exec make -C "$tmp" e2e
