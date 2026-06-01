#!/usr/bin/env bash
# Copyright 2026 The Kubernetes Authors.
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
#
# Prerequisites:
#   GITHUB_TOKEN       — GitHub personal access token
#   AWS credentials    — for fetching published AMIs
#   bin/release-tool   — build: go build -o bin/release-tool ./hack/tools/release-tools
#   bin/clusterawsadm  — build: go build -o bin/clusterawsadm ./cmd/clusterawsadm
#   jq                 — brew install jq

set -euo pipefail

# --- Constants --- 

AWS_ACCOUNT_ID="819546954734"
DEFAULT_REGIONS="ap-south-1,eu-west-3,eu-west-2,eu-west-1,ap-northeast-2,ap-northeast-1,sa-east-1,ca-central-1,ap-southeast-1,ap-southeast-2,eu-central-1,us-east-1,us-east-2,us-west-1,us-west-2"
DEFAULT_OS="ubuntu-24.04,ubuntu-22.04,flatcar-stable"

TARGET_FILE="docs/book/src/development/ami-inventory.md"

# --- Setup --- 

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${SCRIPT_DIR}/.."

die()  { echo "ERROR: $*" >&2; exit 1; }
info() { echo "==> $*" >&2; }

command -v jq &>/dev/null   || die "'jq' is required. Install: brew install jq"
[ -n "${GITHUB_TOKEN:-}" ]  || die "GITHUB_TOKEN is not set"
[ -f "${TARGET_FILE}" ]     || die "Target file not found: ${TARGET_FILE}"

[ -x "bin/release-tool" ]  || { info "Building bin/release-tool...";  go build -o bin/release-tool  ./hack/tools/release-tools; }
[ -x "bin/clusterawsadm" ] || { info "Building bin/clusterawsadm..."; go build -o bin/clusterawsadm ./cmd/clusterawsadm; }

# --- replace_section <marker-name> <content> --- 
# Replaces the lines after <!-- $<name> --> up to the next ## or ### heading.
# Idempotent — reruns overwrite previous output.

replace_section() {
  local marker="<!-- \$$1 -->"
  local content="$2"
  local tmp
  tmp="$(mktemp)"
  printf '%s\n' "${content}" > "${tmp}"
  awk -v m="${marker}" -v f="${tmp}" '
    skip && /^#{2,3} / { skip=0 }
    skip               { next }
    index($0, m)       { print; while ((getline line < f) > 0) print line; skip=1; next }
    { print }
  ' "${TARGET_FILE}" > "${TARGET_FILE}.tmp"
  mv "${TARGET_FILE}.tmp" "${TARGET_FILE}"
  rm -f "${tmp}"
}

# --- 1. Report date --- 

REPORT_DATE="$(date -u +%Y-%m-%d)"
info "Setting report date to ${REPORT_DATE}..."
sed -i.bak "s/^# CAPA AMIs Inventory Report .*/# CAPA AMIs Inventory Report ${REPORT_DATE}/" "${TARGET_FILE}" \
  && rm -f "${TARGET_FILE}.bak"

# --- 2. Defaults sections --- 

info "Populating defaults..."

os_content=""
for os in $(echo "${DEFAULT_OS}" | tr ',' ' '); do
  os_content="${os_content}- \`${os}\`
"
done
replace_section "default_os" "${os_content%?}"

replace_section "default_aws_account_id" "- \`${AWS_ACCOUNT_ID}\`"

region_content=""
for region in $(echo "${DEFAULT_REGIONS}" | tr ',' ' '); do
  region_content="${region_content}- \`${region}\`
"
done
replace_section "default_regions" "${region_content%?}"

# --- 3. Supported Kubernetes versions --- 

info "Fetching supported Kubernetes versions..."
k8s_table="$(./bin/release-tool ami detect-k8s-release \
  --token "${GITHUB_TOKEN}" \
  -o json | jq -r '
    "| Minor Version | Patch Versions |",
    "| --- | --- |",
    (.versions[] |
      "| `" + .minor + "` | " + (.patches | map("`" + . + "`") | join(", ")) + " |"
    )
  ')"
replace_section "k8s_release_table" "${k8s_table}"

# --- 4. Published AMIs --- 

info "Fetching published AMIs from account ${AWS_ACCOUNT_ID}..."
AMI_JSON="$(./bin/clusterawsadm ami list --owner-id "${AWS_ACCOUNT_ID}" -o json)"
AMI_COUNT="$(echo "${AMI_JSON}" | jq '.items | length')"
info "  found ${AMI_COUNT} AMI(s)."

if [ "${AMI_COUNT}" -eq 0 ]; then
  ami_content="_No AMIs found._"
else
  ami_content="$(echo "${AMI_JSON}" | jq -r '
    "| AMI Name | Kubernetes Version | OS | Region | AMI ID |",
    "| --- | --- | --- | --- | --- |",
    (.items[]? |
      "| `" + (.metadata.name           // "-") + "`" +
      " | `" + (.spec.kubernetesVersion // "-") + "`" +
      " | "  + (.spec.os                // "-") +
      " | "  + (.spec.region            // "-") +
      " | `" + (.spec.imageID           // "-") + "` |"
    )
  ')"
fi
replace_section "ami-table" "${ami_content}"

# --- 5. Missing AMIs --- 

info "Computing missing AMIs..."
missing_content="$(echo "${AMI_JSON}" | ./bin/release-tool ami find-missing-ami \
  --token "${GITHUB_TOKEN}" \
  --os "${DEFAULT_OS}" \
  --region "${DEFAULT_REGIONS}" \
  -o json | jq -r '
    if (.items | length) == 0 then "_No missing AMIs._"
    else
      "| Kubernetes Version | OS | Region |",
      "| --- | --- | --- |",
      (.items[] | "| `" + .kubernetesVersion + "` | " + .os + " | " + .region + " |")
    end
  ')"
replace_section "missing_ami_table" "${missing_content}"

# --- 6. EOL AMIs --- 

info "Computing EOL AMIs..."
supported_minors="$(echo "${AMI_JSON}" | jq -r '.items[]?.spec.kubernetesVersion' \
  | sed 's/^v//' | cut -d'.' -f1,2 | sort -u)"

eol_rows=()
if [ "${AMI_COUNT}" -gt 0 ]; then
  while IFS='|' read -r name version os region image_id; do
    minor="$(echo "${version}" | sed 's/^v//' | cut -d'.' -f1,2)"
    if ! grep -qF "${minor}" <<< "${supported_minors}"; then
      eol_rows+=("| \`${name}\` | \`${version}\` | ${os} | ${region} | \`${image_id}\` |")
    fi
  done < <(echo "${AMI_JSON}" | jq -r '
    .items[]? |
    (.metadata.name           // "-") + "|" +
    (.spec.kubernetesVersion  // "-") + "|" +
    (.spec.os                 // "-") + "|" +
    (.spec.region             // "-") + "|" +
    (.spec.imageID            // "-")
  ')
fi

if [ "${#eol_rows[@]}" -eq 0 ]; then
  eol_content="_No EOL AMIs._"
else
  eol_content="| AMI Name | Kubernetes Version | OS | Region | AMI ID |
| --- | --- | --- | --- | --- |
$(printf '%s\n' "${eol_rows[@]}")"
fi
replace_section "eol_ami_table" "${eol_content}"

info "Done — updated ${TARGET_FILE}"
