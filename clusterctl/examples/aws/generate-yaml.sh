#!/bin/bash
set -o errexit
set -o nounset

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

OUTPUT_DIR=${DIR}/out

# provider-components
export CLUSTER_CONTROLLER_IMAGE="${CLUSTER_CONTROLLER_IMAGE:-gcr.io/k8s-cluster-api/aws-cluster-controller:0.0.1}"
export MACHINE_CONTROLLER_IMAGE="${MACHINE_CONTROLLER_IMAGE:-gcr.io/k8s-cluster-api/aws-machine-controller:0.0.1}"

# aws credentials
# TODO generate a secret instead of embedding directly in the YAML
export AWS_REGION="${AWS_REGION:-us-west-2}"
export AWS_ACCESS_KEY_ID="${AWS_ACCESS_KEY_ID}"
export AWS_SECRET_ACCESS_KEY="${AWS_SECRET_ACCESS_KEY}"

PROVIDERCOMPONENT_TEMPLATE_FILE=${DIR}/provider-components.yaml.template
PROVIDERCOMPONENT_GENERATED_FILE=${OUTPUT_DIR}/provider-components.yaml

CLUSTER_TEMPLATE_FILE=${DIR}/cluster.yaml.template
CLUSTER_GENERATED_FILE=${OUTPUT_DIR}/cluster.yaml

MACHINES_TEMPLATE_FILE=${DIR}/machines.yaml.template
MACHINES_GENERATED_FILE=${OUTPUT_DIR}/machines.yaml

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

if [ $OVERWRITE -ne 1 ] && [ -f $PROVIDERCOMPONENT_GENERATED_FILE ]; then
  echo File $PROVIDERCOMPONENT_GENERATED_FILE already exists. Delete it manually before running this script.
  exit 1
fi

mkdir -p ${OUTPUT_DIR}

envsubst < $PROVIDERCOMPONENT_TEMPLATE_FILE \
  > "${PROVIDERCOMPONENT_GENERATED_FILE}"
echo "Done generating ${PROVIDERCOMPONENT_GENERATED_FILE}"

envsubst < $CLUSTER_TEMPLATE_FILE \
  > "${CLUSTER_GENERATED_FILE}"
echo "Done generating ${CLUSTER_GENERATED_FILE}"

envsubst < $MACHINES_TEMPLATE_FILE \
  > "${MACHINES_GENERATED_FILE}"
echo "Done generating ${MACHINES_GENERATED_FILE}"
