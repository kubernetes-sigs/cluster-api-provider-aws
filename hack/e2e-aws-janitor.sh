#!/bin/bash

# Copyright 2019 The Kubernetes Authors.
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

################################################################################
# usage: e2e-aws-janitor.sh [FLAGS]
#  This program is a wrapper for running the aws-janitor command with a check
#  that prevents disallowed AWS keys from being used.
#
# FLAGS
#  To see a full list of flags supported by this program, run "aws-janitor -h"
################################################################################

set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
cd "${REPO_ROOT}" || exit 1

# Require the aws-janitor command.
command -v aws-janitor >/dev/null 2>&1 || \
  { echo "aws-janitor not found" 1>&2; exit 1; }

# Prevent a disallowed AWS key from being used.
if grep -iqF "$(echo "${AWS_ACCESS_KEY_ID-}" | \
  { md5sum 2>/dev/null || md5; } | \
  awk '{print $1}')" hack/e2e-aws-disallowed.txt; then
  echo "The provided AWS key is not allowed" 1>&2
  exit 1
fi

exec aws-janitor -all "${@-}"
