#!/usr/bin/env bash

set -o errexit
set -o nounset

DIR="$(cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd)"
OUTDIR="${DIR}/../out"

mkdir -p "${OUTDIR}"
{ echo "---"; cat "${DIR}/manager/namespace.yaml"; echo "---"; } > "${OUTDIR}/provider-components.yaml"
kustomize build "${DIR}/default/provider-components/" >> "${OUTDIR}/provider-components.yaml"
echo "Finished creating provider components manifest"
cat "${DIR}/samples/machines.yaml" > "${OUTDIR}/machines.yaml"
echo "Copied clusterctl machines manifest"
kustomize build "${DIR}/default/cluster/" > "${OUTDIR}/cluster.yaml"
echo "Finished creating cluster manifest"
