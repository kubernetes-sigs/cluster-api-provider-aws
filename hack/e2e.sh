#!/bin/bash

set -euo pipefail

unset GOFLAGS

git clone "https://github.com/openshift/cluster-api-actuator-pkg.git" cluster-api-actuator-pkg

exec make -C cluster-api-actuator-pkg test-e2e
