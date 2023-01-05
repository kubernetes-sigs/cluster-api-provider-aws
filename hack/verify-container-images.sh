#!/usr/bin/env bash
# Copyright 2022 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This refers https://github.com/kubernetes-sigs/cluster-api/blob/main/hack/verify-container-images.sh

set -o errexit
set -o nounset
set -o pipefail

if [[ "${TRACE-0}" == "1" ]]; then
    set -o xtrace
fi

TRIVY_VERSION=0.35.0

GO_OS="$(go env GOOS)"
if [[ "${GO_OS}" == "linux" ]]; then
  TRIVY_OS="Linux"
elif [[ "${GO_OS}" == "darwin"* ]]; then
  TRIVY_OS="macOS"
fi

GO_ARCH="$(go env GOARCH)"
if [[ "${GO_ARCH}" == "amd" ]]; then
  TRIVY_ARCH="32bit"
elif [[ "${GO_ARCH}" == "amd64"* ]]; then
  TRIVY_ARCH="64bit"
elif [[ "${GO_ARCH}" == "arm" ]]; then
  TRIVY_ARCH="ARM"
elif [[ "${GO_ARCH}" == "arm64" ]]; then
  TRIVY_ARCH="ARM64"
fi

TOOL_BIN=hack/tools/bin
mkdir -p ${TOOL_BIN}

# Downloads trivy scanner
curl -L -o ${TOOL_BIN}/trivy.tar.gz \
    https://github.com/aquasecurity/trivy/releases/download/v${TRIVY_VERSION}/trivy_${TRIVY_VERSION}_${TRIVY_OS}-${TRIVY_ARCH}.tar.gz \

tar xfO ${TOOL_BIN}/trivy.tar.gz trivy > ${TOOL_BIN}/trivy
chmod +x ${TOOL_BIN}/trivy
rm ${TOOL_BIN}/trivy.tar.gz

## Builds the container images to be scanned
make REGISTRY=gcr.io/k8s-staging-cluster-api-aws PULL_POLICY=IfNotPresent TAG=dev docker-build

BRed='\033[1;31m'
BGreen='\033[1;32m'
NC='\033[0m' # No

# Scan the images
echo -e "\n${BGreen}List of dependencies that can bumped to fix the vulnerabilities:${NC}"
${TOOL_BIN}/trivy image -q --exit-code 1 --ignore-unfixed --severity MEDIUM,HIGH,CRITICAL gcr.io/k8s-staging-cluster-api-aws/cluster-api-aws-controller-${GO_ARCH}:dev && R1=$? || R1=$?
echo -e "\n${BGreen}List of dependencies having fixes/no fixes for review only:${NC}"
${TOOL_BIN}/trivy image -q  --severity MEDIUM,HIGH,CRITICAL gcr.io/k8s-staging-cluster-api-aws/cluster-api-aws-controller-${GO_ARCH}:dev

if [ "$R1" -ne "0" ]
then
  echo -e "\n${BRed}Container images check failed! There are vulnerability to be fixed${NC}"
  exit 1
fi

echo -e "\n${BGreen}Container images check passed! No unfixed vulnerability found${NC}"

