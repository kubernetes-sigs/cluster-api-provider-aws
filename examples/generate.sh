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

set -o errexit
set -o nounset

# Directories.
SOURCE_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
OUTPUT_DIR=${OUTPUT_DIR:-${SOURCE_DIR}/_out}

# Binaries
envsubst() {
  python -c 'import os,sys;[sys.stdout.write(os.path.expandvars(l)) for l in sys.stdin]'
}

CLUSTERAWSADM=${CLUSTERAWSADM:-${SOURCE_DIR}/../bin/clusterawsadm}
command -v "${CLUSTERAWSADM}" >/dev/null 2>&1 || echo -v "Cannot find ${CLUSTERAWSADM} in path, build it using 'make binaries' in this repository."

# Cluster.
export CLUSTER_NAME="${CLUSTER_NAME:-test1}"
export KUBERNETES_VERSION="${KUBERNETES_VERSION:-v1.16.1}"

# Machine settings.
export CONTROL_PLANE_MACHINE_TYPE="${CONTROL_PLANE_MACHINE_TYPE:-t2.medium}"
export NODE_MACHINE_TYPE="${CONTROL_PLANE_MACHINE_TYPE:-t2.medium}"
export SSH_KEY_NAME="${SSH_KEY_NAME:-default}"

# Outputs.
COMPONENTS_CLUSTER_API_GENERATED_FILE=${SOURCE_DIR}/provider-components/provider-components-cluster-api.yaml
COMPONENTS_AWS_GENERATED_FILE=${SOURCE_DIR}/provider-components/provider-components-aws.yaml

PROVIDER_COMPONENTS_GENERATED_FILE=${OUTPUT_DIR}/provider-components.yaml
CERTMANAGER_COMPONENTS_GENERATED_FILE=${OUTPUT_DIR}/cert-manager.yaml
CLUSTER_GENERATED_FILE=${OUTPUT_DIR}/cluster.yaml
CONTROLPLANE_GENERATED_FILE=${OUTPUT_DIR}/controlplane.yaml
MACHINEDEPLOYMENT_GENERATED_FILE=${OUTPUT_DIR}/machinedeployment.yaml

# Overwrite flag.
OVERWRITE=0

SCRIPT=$(basename "$0")
while test $# -gt 0; do
        case "$1" in
          -h|--help)
            echo "$SCRIPT - generates input yaml files for Cluster API on aws"
            echo " "
            echo "$SCRIPT [options]"
            echo " "
            echo "options:"
            echo "-h, --help                show brief help"
            echo "-f, --force-overwrite     if file to be generated already exists, force script to overwrite it"
            exit 0
            ;;
          -f)
            OVERWRITE=1
            shift
            ;;
          --force-overwrite)
            OVERWRITE=1
            shift
            ;;
          *)
            break
            ;;
        esac
done

if [ $OVERWRITE -ne 1 ] && [ -d "$OUTPUT_DIR" ]; then
  echo "ERR: Folder ${OUTPUT_DIR} already exists. Delete it manually before running this script."
  exit 1
fi

mkdir -p "${OUTPUT_DIR}"

# Download cert-manager component
curl -sL https://github.com/jetstack/cert-manager/releases/download/v0.11.0/cert-manager.yaml > "${CERTMANAGER_COMPONENTS_GENERATED_FILE}"
echo "Generated ${CERTMANAGER_COMPONENTS_GENERATED_FILE}"

# Generate AWS Credentials.
AWS_B64ENCODED_CREDENTIALS="$(${CLUSTERAWSADM} alpha bootstrap encode-aws-credentials)"
export AWS_B64ENCODED_CREDENTIALS

# Generate cluster resources.
kustomize build "${SOURCE_DIR}/cluster" | envsubst > "${CLUSTER_GENERATED_FILE}"
echo "Generated ${CLUSTER_GENERATED_FILE}"

# Generate controlplane resources.
kustomize build "${SOURCE_DIR}/controlplane" | envsubst > "${CONTROLPLANE_GENERATED_FILE}"
echo "Generated ${CONTROLPLANE_GENERATED_FILE}"

# Generate machinedeployment resources.
kustomize build "${SOURCE_DIR}/machinedeployment" | envsubst >> "${MACHINEDEPLOYMENT_GENERATED_FILE}"
echo "Generated ${MACHINEDEPLOYMENT_GENERATED_FILE}"

# Generate Cluster API provider components file.
CAPI_BRANCH=${CAPI_BRANCH:-"master"}
if [[ ${CAPI_BRANCH} == "stable" ]]; then
  # TODO(vincepri): Fix the version once the first v0.3.x is released.
  curl -sL https://github.com/kubernetes-sigs/cluster-api/releases/download/v0.2.5/cluster-api-components.yaml > "${COMPONENTS_CLUSTER_API_GENERATED_FILE}"
  echo "Downloaded ${COMPONENTS_CLUSTER_API_GENERATED_FILE} from cluster-api stable branch - v0.2.5"
elif [[ ${CAPI_BRANCH} == "dev" ]]; then
  kustomize build ../cluster-api/config > "${COMPONENTS_CLUSTER_API_GENERATED_FILE}"
  kustomize build ../cluster-api/bootstrap/kubeadm/config > "${COMPONENTS_CLUSTER_API_GENERATED_FILE}"
  kustomize build ../cluster-api/controlplane/kubeadm/config > "${COMPONENTS_CLUSTER_API_GENERATED_FILE}"
  echo "Generated ${COMPONENTS_CLUSTER_API_GENERATED_FILE} from cluster-api - ${CAPI_BRANCH}"
else
  kustomize build "github.com/kubernetes-sigs/cluster-api/config/?ref=${CAPI_BRANCH}" > "${COMPONENTS_CLUSTER_API_GENERATED_FILE}"
  kustomize build "github.com/kubernetes-sigs/cluster-api/bootstrap/kubeadm/config/?ref=${CAPI_BRANCH}" > "${COMPONENTS_CLUSTER_API_GENERATED_FILE}"
  kustomize build "github.com/kubernetes-sigs/cluster-api/controlplane/kubeadm/config/?ref=${CAPI_BRANCH}" > "${COMPONENTS_CLUSTER_API_GENERATED_FILE}"
  echo "Generated ${COMPONENTS_CLUSTER_API_GENERATED_FILE} from cluster-api - ${CAPI_BRANCH}"
fi

# Generate AWS Infrastructure Provider components file.
kustomize build "${SOURCE_DIR}/../config" | envsubst > "${COMPONENTS_AWS_GENERATED_FILE}"
echo "Generated ${COMPONENTS_AWS_GENERATED_FILE}"

# Generate a single provider components file.
kustomize build "${SOURCE_DIR}/provider-components" | envsubst > "${PROVIDER_COMPONENTS_GENERATED_FILE}"
echo "Generated ${PROVIDER_COMPONENTS_GENERATED_FILE}"
echo "WARNING: ${PROVIDER_COMPONENTS_GENERATED_FILE} includes AWS credentials"

# Patch kubernetes version
sed -i'' -e 's|kubernetesVersion: .*|kubernetesVersion: '$KUBERNETES_VERSION'|' ${OUTPUT_DIR}/controlplane.yaml

echo "NOTE: Ensure that the cert-manager components are running before creating the provider-components, cluster and control-plane."
