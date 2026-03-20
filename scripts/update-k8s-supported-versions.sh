#!/bin/bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
OUTPUT_PATH="${1:-${REPO_ROOT}/hack/k8s-supported-versions.json}"
GITHUB_API="https://api.github.com/repos/kubernetes/kubernetes"
DL_K8S="https://dl.k8s.io/release"

# --- Dependency checks ---
for tool in curl jq; do
    if ! command -v "$tool" &>/dev/null; then
        echo "ERROR: Required tool '$tool' not found. Please install it." >&2
        exit 1
    fi
done

# --- Step 1: Fetch supported series from release-1.* branches ---
BRANCHES_JSON=$(curl -sSL "${GITHUB_API}/branches?per_page=100")

SUPPORTED_SERIES=()
while IFS= read -r minor; do
    [[ -z "$minor" ]] && continue
    SUPPORTED_SERIES+=("v1.${minor}")
done < <(echo "$BRANCHES_JSON" | jq -r '.[].name | select(test("^release-1\\.[0-9]+$")) | sub("^release-1\\."; "")' | sort -n -r | head -3)

# Fallback: if no branches found, use stable.txt
if [[ ${#SUPPORTED_SERIES[@]} -eq 0 ]]; then
    echo "WARN: No release-1.* branches found, falling back to dl.k8s.io stable.txt" >&2
    LATEST=$(curl -sSL "${DL_K8S}/stable.txt")
    LATEST="${LATEST//[$'\r']}"
    if [[ "$LATEST" =~ ^v([0-9]+)\.([0-9]+)\.([0-9]+)$ ]]; then
        MAJOR="${BASH_REMATCH[1]}"
        MINOR="${BASH_REMATCH[2]}"
        [[ "${BASH_REMATCH[3]}" == "0" ]] && MINOR=$((MINOR - 1))
        for i in 0 1 2; do
            M=$((MINOR - i))
            [[ $M -lt 0 ]] && break
            SUPPORTED_SERIES+=("v${MAJOR}.${M}")
        done
    fi
fi

echo "Supported series: ${SUPPORTED_SERIES[*]:-none}"

# Build regex for supported series: ^v1\.(35|34|33)\.[0-9]+$
SERIES_REGEX="^v1\\.($(IFS='|'; echo "${SUPPORTED_SERIES[*]#v1.}"))\\.[0-9]+\$"

# --- Step 2: Fetch all patch releases from GitHub (paginated) ---
echo "Fetching patch releases from GitHub..."
ALL_RELEASES_JSON="[]"
PAGE=1
MAX_PAGES=15

while [[ $PAGE -le $MAX_PAGES ]]; do
    PAGE_JSON=$(curl -sSL "${GITHUB_API}/releases?per_page=100&page=${PAGE}")
    COUNT=$(echo "$PAGE_JSON" | jq 'length')
    [[ "$COUNT" -eq 0 ]] && break

    FILTERED=$(echo "$PAGE_JSON" | jq --arg re "$SERIES_REGEX" '
        [.[] | select(.prerelease == false) | select(.tag_name | test($re)) |
         {version: .tag_name, release_date: .published_at, series: (.tag_name | capture("^v(?<major>[0-9]+)\\.(?<minor>[0-9]+)\\.") | "v\(.major).\(.minor)")}]
    ')
    ALL_RELEASES_JSON=$(echo "$ALL_RELEASES_JSON" "$FILTERED" | jq -s 'add')

    [[ "$COUNT" -lt 100 ]] && break
    PAGE=$((PAGE + 1))
done

# Sort by version descending
ALL_RELEASES_JSON=$(echo "$ALL_RELEASES_JSON" | jq 'sort_by(.version) | reverse')

# --- Step 3: Write output ---
SUPPORTED_SERIES_JSON=$(printf '%s\n' "${SUPPORTED_SERIES[@]}" | jq -R . | jq -s .)
LAST_UPDATED=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

OUTPUT=$(jq -n \
    --arg lastUpdated "$LAST_UPDATED" \
    --argjson supported_kubernetes_series "$SUPPORTED_SERIES_JSON" \
    --argjson latest_releases "$ALL_RELEASES_JSON" \
    '{lastUpdated: $lastUpdated, supported_kubernetes_series: $supported_kubernetes_series, latest_releases: $latest_releases}')

mkdir -p "$(dirname "$OUTPUT_PATH")"
echo "$OUTPUT" | jq '.' > "$OUTPUT_PATH"
echo "Wrote $(basename "$OUTPUT_PATH") to $OUTPUT_PATH ($(echo "$ALL_RELEASES_JSON" | jq 'length') releases)"
