#!/bin/bash
set -o errexit
set -o nounset

# Directories.
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
OUTPUT_DIR=${OUTPUT_DIR:-${DIR}/out}

# Manager image.
export MANAGER_IMAGE="${MANAGER_IMAGE:-gcr.io/cluster-api-provider-aws/cluster-api-aws-controller:latest}"

# Machine settings.
export CONTROL_PLANE_MACHINE_TYPE="${CONTROL_PLANE_MACHINE_TYPE:-t2.medium}"
export NODE_MACHINE_TYPE="${CONTROL_PLANE_MACHINE_TYPE:-t2.medium}"
export SSH_KEY_NAME="${SSH_KEY_NAME:-default}"

# Templates.
CLUSTER_TEMPLATE_FILE=${DIR}/cluster.yaml.template
CLUSTER_GENERATED_FILE=${OUTPUT_DIR}/cluster.yaml
MACHINES_TEMPLATE_FILE=${DIR}/machines.yaml.template
MACHINES_GENERATED_FILE=${OUTPUT_DIR}/machines.yaml
MANAGER_PATCH_TEMPLATE_FILE=${DIR}/aws_manager_image_patch.yaml.template
MANAGER_PATCH_GENERATED_FILE=${OUTPUT_DIR}/aws_manager_image_patch.yaml
ADDONS_FILE=${OUTPUT_DIR}/addons.yaml

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

if [ $OVERWRITE -ne 1 ] && [ -f $MACHINES_GENERATED_FILE ]; then
  echo File $MACHINES_GENERATED_FILE already exists. Delete it manually before running this script.
  exit 1
fi

if [ $OVERWRITE -ne 1 ] && [ -f $CLUSTER_GENERATED_FILE ]; then
  echo File $CLUSTER_GENERATED_FILE already exists. Delete it manually before running this script.
  exit 1
fi

mkdir -p ${OUTPUT_DIR}

envsubst < $CLUSTER_TEMPLATE_FILE > "${CLUSTER_GENERATED_FILE}"
echo "Done generating ${CLUSTER_GENERATED_FILE}"

envsubst < $MACHINES_TEMPLATE_FILE > "${MACHINES_GENERATED_FILE}"
echo "Done generating ${MACHINES_GENERATED_FILE}"

envsubst < $MANAGER_PATCH_TEMPLATE_FILE > "${MANAGER_PATCH_GENERATED_FILE}"
echo "Done generating ${MANAGER_PATCH_GENERATED_FILE}"

cp  ${DIR}/addons.yaml ${ADDONS_FILE}
echo "Done copying ${ADDONS_FILE}"
