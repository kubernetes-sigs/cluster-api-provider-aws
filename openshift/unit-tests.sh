
#!/bin/bash

set -o errexit
set -o pipefail

echo "Running unit-tests.sh"

# Ensure that some home var is set and that it's not the root
export HOME=${HOME:=/tmp/kubebuilder/testing}
if [ $HOME == "/" ]; then
  export HOME=/tmp/kubebuilder/testing
fi

GOFLAGS='-mod=readonly' make test 
