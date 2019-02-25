#!/usr/bin/env bash
#
# This script updates Go vendoring using dep.

set -euo pipefail

# Go to the root of the repo
cd "$(git rev-parse --show-cdup)"

# Run dep.
dep ensure

(cd hack && ./update-codegen.sh)
