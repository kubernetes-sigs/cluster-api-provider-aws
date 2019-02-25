#!/usr/bin/env bash

PROJECT_ROOT="$(dirname "${BASH_SOURCE[0]}")/.."

# This project performs no code generation
#${PROJECT_ROOT}/vendor/k8s.io/code-generator/generate-groups.sh \
#	all \
#	github.com/openshift/cluster-version-operator/pkg/generated \
#	github.com/openshift/cluster-version-operator/pkg/apis \
#	"config.openshift.io:v1 operatorstatus.openshift.io:v1" \
#	$@ \
