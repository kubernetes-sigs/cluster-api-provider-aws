#!/usr/bin/env bash

# Copyright 2021 The Kubernetes Authors.
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
#set -o verbose

root=$(dirname "${BASH_SOURCE[0]}")/..
kustomize="${root}/hack/tools/bin/kustomize"
test_dir="${root}/test/e2e/data/infrastructure-aws/"

find "${test_dir}"* -maxdepth 1 -type d -print0 | xargs -0 -I {} basename {} | grep -v patches |  grep -v addons | grep -v cni |  grep -v base | xargs -I {} sh -c "${kustomize} build --load-restrictor LoadRestrictionsNone --reorder none ${test_dir}{} > ${test_dir}cluster-template-{}.yaml"
## move the default template to the default file expected by clusterctl
mv "${test_dir}/cluster-template-default.yaml" "${test_dir}/cluster-template.yaml"
