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

echo "# Release notes for Cluster API Provider AWS (CAPA) $VERSION"
echo "[Documentation](https://cluster-api-aws.sigs.k8s.io/)"
echo "# Changelog since $PREVIOUS_VERSION"
$GH api repos/$GH_ORG_NAME/$GH_REPO_NAME/releases/generate-notes -F tag_name=$VERSION -F previous_tag_name=$PREVIOUS_VERSION --jq '.body'
echo "**The image for this release is**: $CORE_CONTROLLER_PROMOTED_IMG:$VERSION"
echo "Thanks to all our contributors!"

