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
# usage: e2e.sh
#  This program runs the e2e tests.
#
# ENVIRONMENT VARIABLES
#  JANITOR_ENABLED
#    Set to 1 to run the aws-janitor command after running the e2e tests.
################################################################################

set -o nounset
set -o pipefail

REPO_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
cd "${REPO_ROOT}" || exit 1

# shellcheck source=../hack/ensure-go.sh
source "${REPO_ROOT}/hack/ensure-go.sh"
# shellcheck source=../hack/ensure-kind.sh
source "${REPO_ROOT}/hack/ensure-kind.sh"
# shellcheck source=../hack/ensure-kubectl.sh
source "${REPO_ROOT}/hack/ensure-kubectl.sh"

ARTIFACTS="${ARTIFACTS:-${PWD}/_artifacts}"
mkdir -p "$ARTIFACTS/logs/"

# our exit handler (trap)
cleanup() {
  # stop boskos heartbeat
  [[ -z ${HEART_BEAT_PID:-} ]] || kill -9 "${HEART_BEAT_PID}"
}
trap cleanup EXIT

#Install requests module explicitly for HTTP calls
python3 -m pip install requests

# If BOSKOS_HOST is set then acquire an AWS account from Boskos.
if [ -n "${BOSKOS_HOST:-}" ]; then
  # Check out the account from Boskos and store the produced environment
  # variables in a temporary file.
  account_env_var_file="$(mktemp)"
  python3 hack/boskos.py --get 1>"${account_env_var_file}"
  checkout_account_status="${?}"

  # If the checkout process was a success then load the account's
  # environment variables into this process.
  # shellcheck disable=SC1090
  [ "${checkout_account_status}" = "0" ] && . "${account_env_var_file}"

  # Always remove the account environment variable file. It contains
  # sensitive information.
  rm -f "${account_env_var_file}"

  if [ ! "${checkout_account_status}" = "0" ]; then
    echo "error getting account from boskos" 1>&2
    exit "${checkout_account_status}"
  fi

  # run the heart beat process to tell boskos that we are still
  # using the checked out account periodically
  python3 -u hack/boskos.py --heartbeat >>$ARTIFACTS/logs/boskos.log 2>&1 &
  HEART_BEAT_PID=$(echo $!)
fi

# Prevent a disallowed AWS key from being used.
if grep -iqF "$(echo "${AWS_ACCESS_KEY_ID-}" |
  { md5sum 2>/dev/null || md5; } |
  awk '{print $1}')" hack/e2e-aws-disallowed.txt; then
  echo "The provided AWS key is not allowed" 1>&2
  exit 1
fi

EXP_EXTERNAL_RESOURCE_GC="true" GC_WORKLOAD="../../data/gcworkload.yaml" make test-e2e-gc ARTIFACTS=$ARTIFACTS

test_status="${?}"

# If Boskos is being used then release the AWS account back to Boskos.
[ -z "${BOSKOS_HOST:-}" ] || python3 -u hack/boskos.py --release

# The janitor is typically not run as part of the e2e process, but rather
# in a parallel process via a service on the same cluster that runs Prow and
# Boskos.
#
# However, setting JANITOR_ENABLED=1 tells this program to run the janitor
# after the e2e test is executed.
if [ "${JANITOR_ENABLED:-0}" = "1" ]; then
  if ! command -v aws-janitor >/dev/null 2>&1; then
    echo "skipping janitor; aws-janitor not found" 1>&2
  else
    aws-janitor -all -v 2
  fi
else
  echo "skipping janitor; JANITOR_ENABLED=${JANITOR_ENABLED:-0}" 1>&2
fi

exit "${test_status}"
