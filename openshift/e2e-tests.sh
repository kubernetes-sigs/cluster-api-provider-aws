#!/bin/bash

set -euo pipefail

echo "Running e2e-tests.sh"

unset GOFLAGS
tmp="$(mktemp -d)"

if [ "${PULL_BASE_REF}" == "master" ]; then
  # the default branch for cluster-capi-operator is main.
  CCAPIO_BASE_REF="main"
else
  CCAPIO_BASE_REF=$PULL_BASE_REF
fi

echo "cloning github.com/openshift/cluster-capi-operator at branch '$CCAPIO_BASE_REF'"
git clone --single-branch --branch="$CCAPIO_BASE_REF" --depth=1 "https://github.com/openshift/cluster-capi-operator.git" "$tmp"

echo "running cluster-capi-operator's: make e2e"
exec make -C "$tmp" e2e
