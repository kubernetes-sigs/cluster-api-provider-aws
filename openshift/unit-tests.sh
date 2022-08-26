#!/bin/bash

set -o errexit
set -o pipefail

echo "Running unit-tests.sh"

# Ensure that some home var is set and that it's not the root
export HOME=${HOME:=/tmp/kubebuilder/testing}
if [ $HOME == "/" ]; then
  export HOME=/tmp/kubebuilder/testing
fi

export GOFLAGS='-mod=readonly'

source ./openshift/fetch_ext_bins.sh
fetch_tools
setup_envs

go test ./api/...
go test ./bootstrap/...
go test ./cmd/...
go test ./controllers/...
go test ./controlplane/...
go test ./exp/...
go test ./pkg/...
