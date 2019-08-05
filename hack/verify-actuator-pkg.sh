#!/bin/bash

set -e

go get -mod=readonly -u github.com/openshift/cluster-api-actuator-pkg
exec git diff --exit-code
