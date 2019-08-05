#!/bin/bash

set -euo pipefail

exec make -C $(go mod download -json github.com/openshift/cluster-api-actuator-pkg | grep '"Dir"' | cut -d '"' -f 4) test-e2e
