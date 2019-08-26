#!/bin/bash

set -euo pipefail

unset GOFLAGS
tmp="$(mktemp -d)"

git clone "https://github.com/openshift/cluster-api-actuator-pkg.git" "$tmp"

exec make -C "$tmp" test-e2e
