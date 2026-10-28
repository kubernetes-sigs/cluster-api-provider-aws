#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
REPO_ROOT=${REPO_ROOT_OVERRIDE:-$REPO_ROOT}
ROOT_GOMOD=${ROOT_GOMOD:-${REPO_ROOT}/go.mod}
TOOLS_GOMOD=${TOOLS_GOMOD:-${REPO_ROOT}/hack/tools/go.mod}

mod_version() {
  awk -v m="$2" '
    $1 == "require" && $2 == m { print $3; exit }
    $1 == m && $2 ~ /^v[0-9]/ { print $2; exit }
  ' "$1"
}

semver_ge() {
  local -a a b
  IFS=. read -r -a a <<<"${1#v}"
  IFS=. read -r -a b <<<"${2#v}"
  local i x y
  for i in 0 1 2; do
    x=${a[i]%%[-+]*}; y=${b[i]%%[-+]*}
    if ((10#${x} > 10#${y})); then return 0; fi
    if ((10#${x} < 10#${y})); then return 1; fi
  done
  return 0
}

if grep -Eq '^[[:space:]]*(replace[[:space:]]+)?k8s\.io/[^[:space:]]+([[:space:]]+v[^[:space:]]+)?[[:space:]]+=>' "${ROOT_GOMOD}" "${TOOLS_GOMOD}"; then
  echo "k8s.io/* replace directives are not supported by $(basename "$0"); update the script." >&2
  exit 1
fi

rc=0
check() { # tools_module root_module
  local tools root
  tools=$(mod_version "${TOOLS_GOMOD}" "$1")
  root=$(mod_version "${ROOT_GOMOD}" "$2")
  if [[ -z "${tools}" || -z "${root}" ]]; then
    echo "ERROR: cannot compare ${1} (hack/tools: '${tools}') with ${2} (root: '${root}')" >&2; rc=1; return
  fi
  if semver_ge "${tools}" "${root}"; then
    echo "OK: hack/tools ${1} ${tools} >= root ${2} ${root}"
  else
    echo "ERROR: hack/tools ${1} ${tools} < root ${2} ${root}. Code generators must not be older than the k8s.io libraries the root module builds against; bump hack/tools/go.mod." >&2
    rc=1
  fi
}
check k8s.io/code-generator k8s.io/apimachinery
for m in k8s.io/apimachinery k8s.io/api k8s.io/apiextensions-apiserver; do check "$m" "$m"; done
exit "${rc}"
