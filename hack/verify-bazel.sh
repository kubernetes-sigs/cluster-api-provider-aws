#!/usr/bin/env bash
# Copyright 2018 The Kubernetes Authors.
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
set -o verbose

export KUBE_ROOT=$(dirname "${BASH_SOURCE}")/..
source "${KUBE_ROOT}/hack/lib/init.sh"

_tmpdir="$(kube::realpath $(mktemp -d -t verify-bazel.XXXXXX))"
kube::util::trap_add "rm -rf ${_tmpdir}" EXIT

_tmp_gopath="${_tmpdir}/go"
_tmp_kuberoot="${_tmp_gopath}/src/sigs.k8s.io/cluster-api-provider-aws"
mkdir -p "${_tmp_kuberoot}/.."
cp -a "${KUBE_ROOT}" "${_tmp_kuberoot}/.."

cd "${_tmp_kuberoot}"
GOPATH="${_tmp_gopath}" PATH="${_tmp_gopath}/bin:${PATH}" "${KUBE_ROOT}/hack/update-bazel.sh"

diff=$(diff -Naupr --exclude=".git" --exclude="bazel-*" -x '_output' "${KUBE_ROOT}" "${_tmp_kuberoot}" || true)

if [[ -n "${diff}" ]]; then
  echo "${diff}" >&2
  echo >&2
  echo "Run ./hack/update-bazel.sh" >&2
  exit 1
fi
