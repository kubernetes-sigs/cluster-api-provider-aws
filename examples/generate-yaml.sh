#!/bin/bash
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

# Directories.
SOURCE_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
DIR=${DIR:=${SOURCE_DIR}}
OUTPUT_DIR=${OUTPUT_DIR:-${DIR}/out}
ENVSUBST=${ENVSUBST:-envsubst}
CLUSTERAWSADM=${CLUSTERAWSADM:-clusterawsadm}

# Cluster name.
export CLUSTER_NAME="${CLUSTER_NAME:-test1}"

# Machine settings.
export CONTROL_PLANE_MACHINE_TYPE="${CONTROL_PLANE_MACHINE_TYPE:-t2.medium}"
export NODE_MACHINE_TYPE="${CONTROL_PLANE_MACHINE_TYPE:-t2.medium}"
export SSH_KEY_NAME="${SSH_KEY_NAME:-default}"
export VPC_ID="${VPC_ID:-}"

# Templates.
CLUSTER_TEMPLATE_FILE=${DIR}/cluster.yaml.template
CLUSTER_NETWORKSPEC_TEMPLATE_FILE=${DIR}/cluster-network-spec.yaml.template
CLUSTER_GENERATED_FILE=${OUTPUT_DIR}/cluster.yaml
MACHINES_TEMPLATE_FILE=${DIR}/machines.yaml.template
MACHINES_GENERATED_FILE=${OUTPUT_DIR}/machines.yaml
CONTROLPLANE_MACHINES_HA_TEMPLATE_FILE=${DIR}/controlplane-machines-ha.yaml.template
CONTROLPLANE_MACHINES_HA_GENERATED_FILE=${OUTPUT_DIR}/controlplane-machines-ha.yaml
CONTROLPLANE_MACHINE_TEMPLATE_FILE=${DIR}/controlplane-machine.yaml.template
CONTROLPLANE_MACHINE_GENERATED_FILE=${OUTPUT_DIR}/controlplane-machine.yaml
DEPLOYMENT_MACHINES_TEMPLATE_FILE=${DIR}/machine-deployment.yaml.template
DEPLOYMENT_MACHINES_GENERATED_FILE=${OUTPUT_DIR}/machine-deployment.yaml
ADDONS_FILE=${OUTPUT_DIR}/addons.yaml
PROVIDER_COMPONENTS_SRC=${DIR}/provider-components-base.yaml
PROVIDER_COMPONENTS_FILE=${OUTPUT_DIR}/provider-components.yaml
CREDENTIALS_FILE=${OUTPUT_DIR}/aws-credentials.yaml

# Overwrite flag.
OVERWRITE=0

SCRIPT=$(basename $0)
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

if [ $OVERWRITE -ne 1 ] && [ -d $OUTPUT_DIR ]; then
  echo "ERR: Folder ${OUTPUT_DIR} already exists. Delete it manually before running this script."
  exit 1
fi

mkdir -p ${OUTPUT_DIR}

if [ -z "$VPC_ID" ]; then
  $ENVSUBST < $CLUSTER_TEMPLATE_FILE > "${CLUSTER_GENERATED_FILE}"
  echo "Done generating ${CLUSTER_GENERATED_FILE}"
else
  $ENVSUBST < $CLUSTER_NETWORKSPEC_TEMPLATE_FILE > "${CLUSTER_GENERATED_FILE}"
  echo "Done generating ${CLUSTER_GENERATED_FILE}"
fi

$ENVSUBST < $MACHINES_TEMPLATE_FILE > "${MACHINES_GENERATED_FILE}"
echo "Done generating ${MACHINES_GENERATED_FILE}"

$ENVSUBST < $CONTROLPLANE_MACHINES_HA_TEMPLATE_FILE > "${CONTROLPLANE_MACHINES_HA_GENERATED_FILE}"
echo "Done generating ${CONTROLPLANE_MACHINES_HA_GENERATED_FILE}"

$ENVSUBST < $CONTROLPLANE_MACHINE_TEMPLATE_FILE > "${CONTROLPLANE_MACHINE_GENERATED_FILE}"
echo "Done generating ${CONTROLPLANE_MACHINE_GENERATED_FILE}"

$ENVSUBST < $DEPLOYMENT_MACHINES_TEMPLATE_FILE > "${DEPLOYMENT_MACHINES_GENERATED_FILE}"
echo "Done generating ${DEPLOYMENT_MACHINES_GENERATED_FILE}"

cp  ${DIR}/addons.yaml ${ADDONS_FILE}
echo "Done copying ${ADDONS_FILE}"

CREDENTIALS="$(${CLUSTERAWSADM} alpha bootstrap encode-aws-credentials)"
echo "Generated credentials"

PROVIDER_COMPONENTS="$(cat ${PROVIDER_COMPONENTS_SRC})"

echo -e "${PROVIDER_COMPONENTS}\n${CREDENTIALS}" > "${PROVIDER_COMPONENTS_FILE}"
echo "Done writing ${PROVIDER_COMPONENTS_FILE}"
echo "WARNING: ${PROVIDER_COMPONENTS_FILE} includes credentials"
