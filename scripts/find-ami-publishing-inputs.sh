#!/bin/bash
#
# Script to find various Kubernetes-related package versions for CAPA AMI publishing
# based on a target Kubernetes release series (e.g., v1.31 or 1.31) or a
# specific patch version (e.g., v1.31.5 or 1.31.5).
#
# Usage: ./script.sh [vX.Y | X.Y | vX.Y.Z | X.Y.Z]
# Example: ./script.sh v1.31
#          ./script.sh 1.31
#          ./script.sh v1.31.5
#          ./script.sh 1.31.5
#          ./script.sh
#

# Exit immediately if a command exits with a non-zero status.
# Treat unset variables as an error.
# The return value of a pipeline is the status of the last command to exit with a non-zero status,
# or zero if no command exited with a non-zero status.
set -euo pipefail

# --- Dependency Checks ---
# Check for container runtime
if command -v podman &> /dev/null; then
    CONTAINER_CMD="podman"
elif command -v docker &> /dev/null; then
    CONTAINER_CMD="docker"
else
    echo "ERROR: Neither 'podman' nor 'docker' command found. Please install one." >&2
    exit 1
fi
echo "INFO: Using container runtime: $CONTAINER_CMD"

# Check for other required host tools
REQUIRED_TOOLS=( "curl" "grep" "sed" "awk" "cut" "sort" "tail" )
MISSING_TOOLS=()
for tool in "${REQUIRED_TOOLS[@]}"; do
    if ! command -v "$tool" &> /dev/null; then
        MISSING_TOOLS+=("$tool")
    fi
done

if [ ${#MISSING_TOOLS[@]} -ne 0 ]; then
    echo "ERROR: Missing required command-line tools: ${MISSING_TOOLS[*]}" >&2
    echo "Please install them and try again." >&2
    exit 1
fi
echo "INFO: All required host tools found."


# --- Configuration & Input Handling ---
DEFAULT_K8S_RELEASE_SERIES_NO_V="1.31"
K8S_RELEASE_ARG="${1:-}" # Get potential argument
K8S_TARGET_VERSION_NO_V="" # Specific X.Y.Z requested, if any
K8S_TARGET_VERSION_V=""    # Specific vX.Y.Z requested, if any

if [[ -n "$K8S_RELEASE_ARG" ]]; then
    # Argument provided, validate and normalize
    if [[ "$K8S_RELEASE_ARG" =~ ^v([0-9]+\.[0-9]+)\.([0-9]+)$ ]]; then # Matches vX.Y.Z
        K8S_RELEASE_SERIES="v${BASH_REMATCH[1]}" # vX.Y
        K8S_RELEASE_SERIES_NO_V="${BASH_REMATCH[1]}" # X.Y
        K8S_TARGET_VERSION_NO_V="${BASH_REMATCH[1]}.${BASH_REMATCH[2]}" # X.Y.Z
        K8S_TARGET_VERSION_V="$K8S_RELEASE_ARG" # vX.Y.Z
        echo "INFO: Specific version requested: ${K8S_TARGET_VERSION_V}"
    elif [[ "$K8S_RELEASE_ARG" =~ ^([0-9]+\.[0-9]+)\.([0-9]+)$ ]]; then # Matches X.Y.Z
        K8S_RELEASE_SERIES="v${BASH_REMATCH[1]}" # vX.Y
        K8S_RELEASE_SERIES_NO_V="${BASH_REMATCH[1]}" # X.Y
        K8S_TARGET_VERSION_NO_V="$K8S_RELEASE_ARG" # X.Y.Z
        K8S_TARGET_VERSION_V="v$K8S_RELEASE_ARG" # vX.Y.Z
        echo "INFO: Specific version requested: ${K8S_TARGET_VERSION_V}"
    elif [[ "$K8S_RELEASE_ARG" =~ ^v([0-9]+\.[0-9]+)$ ]]; then # Matches vX.Y
        K8S_RELEASE_SERIES="$K8S_RELEASE_ARG"
        K8S_RELEASE_SERIES_NO_V="${BASH_REMATCH[1]}"
        echo "INFO: Release series requested: ${K8S_RELEASE_SERIES}"
    elif [[ "$K8S_RELEASE_ARG" =~ ^([0-9]+\.[0-9]+)$ ]]; then # Matches X.Y
        K8S_RELEASE_SERIES="v$K8S_RELEASE_ARG"
        K8S_RELEASE_SERIES_NO_V="$K8S_RELEASE_ARG"
        echo "INFO: Release series requested: ${K8S_RELEASE_SERIES}"
    else
        echo "ERROR: Invalid format for Kubernetes version/series argument: '$K8S_RELEASE_ARG'." >&2
        echo "       Expected format 'vX.Y', 'X.Y', 'vX.Y.Z', or 'X.Y.Z'." >&2
        exit 1
    fi
else
    # No argument, use default series
    K8S_RELEASE_SERIES_NO_V="$DEFAULT_K8S_RELEASE_SERIES_NO_V"
    K8S_RELEASE_SERIES="v$K8S_RELEASE_SERIES_NO_V"
    echo "INFO: No version specified, using default series: ${K8S_RELEASE_SERIES}"
fi

echo "--- Finding versions for Kubernetes release series: ${K8S_RELEASE_SERIES} ---"
echo "INFO: Using Kubernetes major.minor: ${K8S_RELEASE_SERIES_NO_V}"

# --- Get Deb Package Versions ---
echo "INFO: Running container (${CONTAINER_CMD}) to fetch deb package versions from pkgs.k8s.io..."

# Construct the command for kubectl lookup based on whether a specific version was requested
if [[ -n "$K8S_TARGET_VERSION_NO_V" ]]; then
    # Look for the exact X.Y.Z version (allow any deb revision like -1.1)
    # Need to escape special characters for grep and ensure we match start of version
    KUBECTL_LOOKUP_PATTERN="${K8S_TARGET_VERSION_NO_V}-[0-9]"
    KUBECTL_NOT_FOUND_MSG="kubectl version ${K8S_TARGET_VERSION_NO_V} not found"
    # We still sort and take tail -n 1 in case multiple deb revisions exist for the same k8s version
    KUBECTL_MADISON_CMD="apt-cache madison kubectl | grep \"${KUBECTL_LOOKUP_PATTERN}\" | sort -V | tail -n 1 || echo \"${KUBECTL_NOT_FOUND_MSG}\""
else
    # Look for the latest within the X.Y series
    KUBECTL_LOOKUP_PATTERN="${K8S_RELEASE_SERIES_NO_V}"
    KUBECTL_NOT_FOUND_MSG="kubectl version not found for series ${K8S_RELEASE_SERIES_NO_V}"
    KUBECTL_MADISON_CMD="apt-cache madison kubectl | grep \"${KUBECTL_LOOKUP_PATTERN}\" | sort -V | tail -n 1 || echo \"${KUBECTL_NOT_FOUND_MSG}\""
fi

PACKAGE_INFO=$(${CONTAINER_CMD} run --platform=linux/amd64 --rm -i ubuntu:22.04 /bin/bash <<EOF
set -e # Exit on error inside the container script
export DEBIAN_FRONTEND=noninteractive
apt-get update > /dev/null
apt-get install -y gpg curl sudo > /dev/null

# Add Kubernetes APT repository using the vX.Y format
echo "deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/${K8S_RELEASE_SERIES}/deb/ /" | tee /etc/apt/sources.list.d/kubernetes.list > /dev/null
mkdir -p -m 755 /etc/apt/keyrings
curl -fsSL "https://pkgs.k8s.io/core:/stable:/${K8S_RELEASE_SERIES}/deb/Release.key" | gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg > /dev/null
apt-get update > /dev/null

# Output marker and package info for kubectl using the constructed command
echo "---KUBECTL_START---"
eval "$KUBECTL_MADISON_CMD" # Use eval to run the command string built on the host
echo "---KUBECTL_END---"

# Output marker and package info for CNI (always latest for the series)
echo "---CNI_START---"
apt-cache madison kubernetes-cni | sort -V | tail -n 1 || echo "kubernetes-cni version not found"
echo "---CNI_END---"
EOF
)

# Check if container command was successful
if [ $? -ne 0 ]; then
    echo "ERROR: Failed to run ${CONTAINER_CMD} container." >&2
    exit 1
fi

echo "INFO: Processing package information..."

# Extract kubectl version found by the container
K8S_DEB_VERSION_LINE=$(echo "${PACKAGE_INFO}" | awk '/---KUBECTL_START---/{flag=1; next} /---KUBECTL_END---/{flag=0} flag')
if [[ -z "$K8S_DEB_VERSION_LINE" || "$K8S_DEB_VERSION_LINE" == *"not found"* ]]; then
    # Use the specific "not found" message from the container if available
    NOT_FOUND_MSG=$(echo "$K8S_DEB_VERSION_LINE" | grep "not found" || echo "Could not find requested kubectl deb package version.")
    echo "ERROR: ${NOT_FOUND_MSG}" >&2
    exit 1
fi
# Extract the version part (e.g., 1.31.5-1.1) which is the 3rd field
K8S_DEB_VERSION=$(echo "$K8S_DEB_VERSION_LINE" | awk '{print $3}')

# Extract latest CNI version (always latest for the series)
CNI_DEB_VERSION_LINE=$(echo "${PACKAGE_INFO}" | awk '/---CNI_START---/{flag=1; next} /---CNI_END---/{flag=0} flag')
if [[ -z "$CNI_DEB_VERSION_LINE" || "$CNI_DEB_VERSION_LINE" == *"not found"* ]]; then
    echo "ERROR: Could not find kubernetes-cni deb package version for series ${K8S_RELEASE_SERIES}." >&2
    exit 1
fi
# Extract the version part (e.g., 1.5.1-1.1) which is the 3rd field
CNI_DEB_VERSION=$(echo "$CNI_DEB_VERSION_LINE" | awk '{print $3}')

# --- Derive K8s Versions ---
# K8s Semver (vX.Y.Z)
if [[ -n "$K8S_TARGET_VERSION_V" ]]; then
    # If specific version was requested, use that as the Semver base
    K8S_SEMVER="$K8S_TARGET_VERSION_V"
    # Sanity check: ensure the found deb package matches the requested semver
    FOUND_SEMVER_BASE="v$(echo "$K8S_DEB_VERSION" | cut -d'-' -f1)"
    if [[ "$FOUND_SEMVER_BASE" != "$K8S_SEMVER" ]]; then
         echo "WARN: Found kubectl deb version (${K8S_DEB_VERSION}) does not match requested K8s Semver (${K8S_SEMVER}). Using found version's base." >&2
         K8S_SEMVER="$FOUND_SEMVER_BASE" # Trust the found version more
    fi
else
    # If only series was requested, derive Semver from the found deb package
    K8S_SEMVER="v$(echo "$K8S_DEB_VERSION" | cut -d'-' -f1)"
fi

# K8s Release Series (vX.Y) - Always derived from the input/default series
K8S_RELEASE_SERIES_OUTPUT="${K8S_RELEASE_SERIES}"

# K8s RPM Package Version (X.Y.Z) - Derived from the final K8S_SEMVER
K8S_RPM_VERSION="${K8S_SEMVER#v}"

# --- Derive CNI Version ---
# CNI Semver (vX.Y.Z) - Derived from the latest CNI deb package found for the series
CNI_SEMVER="v$(echo "$CNI_DEB_VERSION" | cut -d'-' -f1)"

# --- Determine CRICTL Version ---
# Always based on the K8S_RELEASE_SERIES (vX.Y)
echo "INFO: Fetching crictl release tags from GitHub API for series ${K8S_RELEASE_SERIES}..."
CRICTL_RAW_TAG=$(curl -s https://api.github.com/repos/kubernetes-sigs/cri-tools/releases \
  | grep '"tag_name":' \
  | grep "\"${K8S_RELEASE_SERIES}." \
  | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/' \
  | sort -V \
  | tail -n 1 \
  || echo "" ) # Avoid script exit if curl or grep fails in the pipe

if [ -z "$CRICTL_RAW_TAG" ]; then
    echo "WARN: Could not automatically determine latest crictl tag for ${K8S_RELEASE_SERIES} from GitHub API." >&2
    echo "WARN: Please verify the crictl version manually at https://github.com/kubernetes-sigs/cri-tools/releases" >&2
    CRICTL_VERSION="<Check Manually for ${K8S_RELEASE_SERIES}>"
else
    CRICTL_VERSION="${CRICTL_RAW_TAG#v}"
fi


# --- Output Results ---
echo
echo "--- Derived Versions for K8s ${K8S_RELEASE_SERIES} (Target: ${K8S_TARGET_VERSION_V:-Latest}) ---"
# Output order from previous request is maintained
printf "%-25s : %s\n" "K8s Semver" "${K8S_SEMVER}"
printf "%-25s : %s\n" "K8s Release Series" "${K8S_RELEASE_SERIES_OUTPUT}"
printf "%-25s : %s\n" "K8s RPM Package Version" "${K8S_RPM_VERSION}"
printf "%-25s : %s\n" "K8s Deb Package Version" "${K8S_DEB_VERSION}"
printf "%-25s : %s\n" "CNI Semver" "${CNI_SEMVER}"
printf "%-25s : %s\n" "CNI Deb Package Version" "${CNI_DEB_VERSION}"
printf "%-25s : %s\n" "CRICTL Version" "${CRICTL_VERSION}"
echo "---------------------------------------"

exit 0
