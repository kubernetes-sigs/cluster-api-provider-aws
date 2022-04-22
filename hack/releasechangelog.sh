#!/bin/bash
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

set -o errexit
set -o nounset
set -o pipefail

function show_help()
{
  cat << EOF
Usage: ${0##*/} -v VERSION -p PREVIOUS_VERSION -o GH_ORG_NAME -r GH_REPO_NAME -c CORE_CONTROLLER_PROMOTED_IMG

This generates the release notes for the new CAPA version being released.

Required Arguments:
    -v VERSION                      Version of the Cluster API Provider AWS (CAPA) being released
    -p PREVIOUS_VERSION             Current CAPA version previously released
    -o GH_ORG_NAME                  GitHub organization name
    -r GH_REPO_NAME                 GitHub repository name
    -c CORE_CONTROLLER_PROMOTED_IMG Image used for this release
EOF
}

while getopts "v:p:o:r:c:h" opt; do
  case $opt in
    v)
      VERSION=${OPTARG}
      ;;
    p)
      PREVIOUS_VERSION=${OPTARG}
      ;;
    o)
      GH_ORG_NAME=${OPTARG}
      ;;
    r)
      GH_REPO_NAME=${OPTARG}
      ;;
    c)
      CORE_CONTROLLER_PROMOTED_IMG=${OPTARG}
      ;;
    h)
      show_help
      exit 0
      ;;
    *)
      show_help >&2
      exit 1
      ;;
  esac
done

echo "# Release notes for Cluster API Provider AWS (CAPA) $VERSION"
echo "[Documentation](https://cluster-api-aws.sigs.k8s.io/)"
echo "# Changelog since $PREVIOUS_VERSION"
$GH api repos/$GH_ORG_NAME/$GH_REPO_NAME/releases/generate-notes -F tag_name=$VERSION -F previous_tag_name=$PREVIOUS_VERSION --jq '.body'
echo "**The image for this release is**: $CORE_CONTROLLER_PROMOTED_IMG:$VERSION"
echo "Thanks to all our contributors!"

