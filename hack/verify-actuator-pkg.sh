#!/bin/bash

set -e

if ! command -v dep 1>/dev/null 2>&1; then
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
fi

if dep check | grep -q cluster-api-actuator-pkg; then
	exit 1
fi
dep ensure -update github.com/openshift/cluster-api-actuator-pkg
git diff --exit-code
