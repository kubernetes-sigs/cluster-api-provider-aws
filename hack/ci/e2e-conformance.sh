#!/usr/bin/env bash

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

# hack script for running a cluster-api-provider-gcp e2e

set -o errexit -o nounset -o pipefail

REGISTRY=${REGISTRY:-"gcr.io/"$(gcloud config get-value project)}
AWS_REGION=${AWS_REGION:-"us-east-1"}
KUBERNETES_VERSION=${KUBERNETES_VERSION:-"v1.16.1"}
TIMESTAMP=$(date +"%Y-%m-%dT%H:%M:%SZ")

ARTIFACTS="${ARTIFACTS:-${PWD}/_artifacts}"
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"

# dump logs from kind and all the nodes
dump-logs() {

  # log version information
  echo "=== versions ==="
  echo "kind : $(kind version)" || true
  echo "bootstrap cluster:"
  kubectl --kubeconfig=$(kind get kubeconfig-path --name="clusterapi") version || true
  echo "deployed cluster:"
  kubectl --kubeconfig=${PWD}/kubeconfig version || true
  echo ""

  # dump all the info from the CAPI related CRDs
  kubectl --kubeconfig=$(kind get kubeconfig-path --name="clusterapi") get \
  clusters,gcpclusters,machines,gcpmachines,kubeadmconfigs,machinedeployments,gcpmachinetemplates,kubeadmconfigtemplates,machinesets \
  --all-namespaces -o yaml >> "${ARTIFACTS}/logs/capg.info" || true

  # dump images info
  echo "images in docker" >> "${ARTIFACTS}/logs/images.info"
  docker images >> "${ARTIFACTS}/logs/images.info"
  echo "images from bootstrap using containerd CLI" >> "${ARTIFACTS}/logs/images.info"
  docker exec clusterapi-control-plane ctr -n k8s.io images list >> "${ARTIFACTS}/logs/images.info" || true
  echo "images in bootstrap cluster using kubectl CLI" >> "${ARTIFACTS}/logs/images.info"
  (kubectl --kubeconfig=$(kind get kubeconfig-path --name="clusterapi") get pods --all-namespaces -o json \
   | jq --raw-output '.items[].spec.containers[].image' | sort)  >> "${ARTIFACTS}/logs/images.info" || true
  echo "images in deployed cluster using kubectl CLI" >> "${ARTIFACTS}/logs/images.info"
  (kubectl --kubeconfig=${PWD}/kubeconfig get pods --all-namespaces -o json \
   | jq --raw-output '.items[].spec.containers[].image' | sort)  >> "${ARTIFACTS}/logs/images.info" || true

  # dump cluster info for kind
  kubectl --kubeconfig=$(kind get kubeconfig-path --name="clusterapi") cluster-info dump > "${ARTIFACTS}/logs/kind-cluster.info" || true

  # export all logs from kind
  kind "export" logs --name="clusterapi" "${ARTIFACTS}/logs" || true
}

# cleanup all resources we use
cleanup() {
  # KIND_IS_UP is true once we: kind create
  if [[ "${KIND_IS_UP:-}" = true ]]; then
    timeout 600 kubectl \
      --kubeconfig=$(kind get kubeconfig-path --name="clusterapi") \
      delete cluster test1 || true
     timeout 600 kubectl \
      --kubeconfig=$(kind get kubeconfig-path --name="clusterapi") \
      wait --for=delete cluster/test1 || true
    make kind-reset || true
  fi
  # clean up e2e.test symlink
  (cd "$(go env GOPATH)/src/k8s.io/kubernetes" && rm -f _output/bin/e2e.test) || true
}

# our exit handler (trap)
exit-handler() {
  dump-logs
  cleanup
}

init_image() {
  if [[ "${REUSE_OLD_IMAGES:-false}" == "true" ]]; then
    filter="capa-ami-ubuntu-18.04-1.16.*"
    image=$(aws ec2 describe-images --query 'Images[*].[ImageId,Name]' --filters "Name=name,Values=${filter}" \
      --region "${AWS_REGION}" --output json | jq -r '.[0][0]')
    if [[ ! -z "$image" ]]; then
      return
    fi
  fi

  if ! command -v ansible &> /dev/null; then
    if [[ $EUID -ne 0 ]]; then
      echo "Please install ansible and try again."
      exit 1
    else
      # we need pip to install ansible
      curl -L https://bootstrap.pypa.io/get-pip.py -o get-pip.py
      python get-pip.py --user
      rm -f get-pip.py

      # install ansible needed by packer
      version="2.8.5"
      python -m pip install "ansible==${version}"
    fi
  fi
  if ! command -v packer &> /dev/null; then
    hostos=$(go env GOHOSTOS)
    hostarch=$(go env GOHOSTARCH)
    version="1.4.3"
    url="https://releases.hashicorp.com/packer/${version}/packer_${version}_${hostos}_${hostarch}.zip"
    echo "Downloading packer from $url"
    wget --quiet -O packer.zip $url  && \
      unzip packer.zip && \
      rm packer.zip && \
      ln -s $PWD/packer /usr/local/bin/packer
  fi
  (cd "$(go env GOPATH)/src/sigs.k8s.io/image-builder/images/capi" && \
    sed -i 's/1\.15\.4/1.16.1/' packer/config/kubernetes.json && \
    sed -i 's/1\.15/1.16/' packer/config/kubernetes.json)
  if [[ $EUID -ne 0 ]]; then
    # install goss plugin
    (cd "$(go env GOPATH)/src/sigs.k8s.io/image-builder/images/capi/packer/ami" && \
      make plugins)
    (cd "$(go env GOPATH)/src/sigs.k8s.io/image-builder/images/capi" && \
      AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID:-""}
      AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY:-""} \
      AWS_REGION=$AWS_REGION \
      make build-ami-default)
  else
    # assume we are running in the CI environment as root
    # Add a user for ansible to work properly
    groupadd -r packer && useradd -m -s /bin/bash -r -g packer packer
    # install goss plugin
    su - packer -c "bash -c 'cd /go/src/sigs.k8s.io/image-builder/images/capi/packer/ami && make plugins'"
    # use the packer user to run the build
    su - packer -c "bash -c 'cd /go/src/sigs.k8s.io/image-builder/images/capi && AWS_REGION=$AWS_REGION AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID:-""} AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY:-""} make build-ami-default'"
  fi
}

# build kubernetes / node image, e2e binaries
build() {
  # possibly enable bazel build caching before building kubernetes
  if [[ "${BAZEL_REMOTE_CACHE_ENABLED:-false}" == "true" ]]; then
    create_bazel_cache_rcs.sh || true
  fi

  pushd "$(go env GOPATH)/src/k8s.io/kubernetes"

  # make sure we have e2e requirements
  bazel build //cmd/kubectl //test/e2e:e2e.test //vendor/github.com/onsi/ginkgo/ginkgo

  # ensure the e2e script will find our binaries ...
  mkdir -p "${PWD}/_output/bin/"
  cp "${PWD}/bazel-bin/test/e2e/e2e.test" "${PWD}/_output/bin/e2e.test"
  PATH="$(dirname "$(find "${PWD}/bazel-bin/" -name kubectl -type f)"):${PATH}"
  export PATH

  # attempt to release some memory after building
  sync || true
  echo 1 > /proc/sys/vm/drop_caches || true

  popd
}

# generate manifests needed for creating the GCP cluster to run the tests
generate_manifests() {
  if ! command -v kustomize >/dev/null 2>&1; then
    (cd ./hack/tools/ && GO111MODULE=on go install sigs.k8s.io/kustomize/kustomize/v3)
  fi

  filter="capa-ami-ubuntu-18.04-1.16.*"
  image_id=$(aws ec2 describe-images --query 'Images[*].[ImageId,Name]' --filters "Name=name,Values=$filter" --region ${AWS_REGION} --output json | jq -r '.[0][0] | select (.!=null)')

  PULL_POLICY=IfNotPresent \
  AWS_REGION=${AWS_REGION} \
  KUBERNETES_VERSION=$KUBERNETES_VERSION \
  IMAGE_ID=$image_id \
    make modules docker-build generate-examples
}

# install cloud formation templates, iam objects etc
create_stack() {
  "${REPO_ROOT}/bin/clusterawsadm" alpha bootstrap create-stack
}

# fix manifests to use k/k from CI
fix_manifests() {
  CI_VERSION=${CI_VERSION:-$(curl -sSL https://dl.k8s.io/ci/latest-green.txt)}
  echo "Overriding Kubernetes version to : ${CI_VERSION}"
  sed -i 's|kubernetesVersion: .*|kubernetesVersion: "ci/'${CI_VERSION}'"|' examples/_out/controlplane.yaml
  sed -i 's|CI_VERSION=.*|CI_VERSION='$CI_VERSION'|' examples/_out/controlplane.yaml
  sed -i 's|CI_VERSION=.*|CI_VERSION='$CI_VERSION'|' examples/_out/machinedeployment.yaml
}

# up a cluster with kind
create_cluster() {
  # actually create the cluster
  KIND_IS_UP=true

  # Load the newly built image into kind and start the cluster
  LOAD_IMAGE="${REGISTRY}/cluster-api-aws-controller-amd64:dev" \
    make create-cluster-management

  # Wait till all machines are running (bail out at 30 mins)
  attempt=0
  while true; do
    kubectl get machines --kubeconfig=$(kind get kubeconfig-path --name="clusterapi")
    read running total <<< $(kubectl get machines --kubeconfig=$(kind get kubeconfig-path --name="clusterapi") \
      -o json | jq -r '.items[].status.phase' | awk 'BEGIN{count=0} /(r|R)unning/{count++} END{print count " " NR}') ;
    if [[ $total == "4" && $running == "4" ]]; then
      return 0
    fi
    read failed total <<< $(kubectl get machines --kubeconfig=$(kind get kubeconfig-path --name="clusterapi") \
      -o json | jq -r '.items[].status.phase' | awk 'BEGIN{count=0} /(f|F)ailed/{count++} END{print count " " NR}') ;
    if [[ ! $failed -eq 0 ]]; then
      echo "$failed machines (out of $total) in cluster failed ... bailing out"
      exit 1
    fi
    timestamp=$(date +"[%H:%M:%S]")
    if [ $attempt -gt 180 ]; then
      echo "cluster did not start in 30 mins ... bailing out!"
      exit 1
    fi
    echo "$timestamp Total machines : $total / Running : $running .. waiting for 10 seconds"
    sleep 10
    attempt=$((attempt+1))
  done
}

# run e2es with kubetest
run_tests() {
  # export the KUBECONFIG
  KUBECONFIG="${PWD}/kubeconfig"
  export KUBECONFIG

  # ginkgo regexes
  SKIP="${SKIP:-}"
  FOCUS="${FOCUS:-"\\[Conformance\\]"}"
  # if we set PARALLEL=true, skip serial tests set --ginkgo-parallel
  if [[ "${PARALLEL:-false}" == "true" ]]; then
    export GINKGO_PARALLEL=y
    if [[ -z "${SKIP}" ]]; then
      SKIP="\\[Serial\\]"
    else
      SKIP="\\[Serial\\]|${SKIP}"
    fi
  fi

  # get the number of worker nodes
  # TODO(bentheelder): this is kinda gross
  NUM_NODES="$(kubectl get nodes --kubeconfig=$KUBECONFIG \
    -o=jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.spec.taints}{"\n"}{end}' \
    | grep -cv "node-role.kubernetes.io/master" )"

  # wait for all the nodes to be ready
  kubectl wait --for=condition=Ready node --kubeconfig=$KUBECONFIG --all || true

  # setting this env prevents ginkg e2e from trying to run provider setup
  export KUBERNETES_CONFORMANCE_TEST="y"
  # run the tests
  (cd "$(go env GOPATH)/src/k8s.io/kubernetes" && ./hack/ginkgo-e2e.sh \
    '--provider=skeleton' "--num-nodes=${NUM_NODES}" \
    "--ginkgo.focus=${FOCUS}" "--ginkgo.skip=${SKIP}" \
    "--report-dir=${ARTIFACTS}" '--disable-log-dump=true')

  unset KUBECONFIG
  unset KUBERNETES_CONFORMANCE_TEST
}

# setup kind, build kubernetes, create a cluster, run the e2es
main() {
  if [[ ${1:-} == "--verbose" ]]; then
     set -o xtrace
  fi

  if [[ ${1:-} == "--clean" ]]; then
    cleanup
    return 0
  fi

  # create temp dir and setup cleanup
  TMP_DIR=$(mktemp -d)
  SKIP_CLEANUP=${SKIP_CLEANUP:-""}
  if [[ -z "${SKIP_CLEANUP}" ]]; then
    trap exit-handler EXIT
  fi
  # ensure artifacts exists when not in CI
  export ARTIFACTS
  mkdir -p "${ARTIFACTS}/logs"

  source "${REPO_ROOT}/hack/ensure-go.sh"
  source "${REPO_ROOT}/hack/ensure-kind.sh"

  build
  init_image
  generate_manifests
  if [[ ${1:-} == "--use-ci-artifacts" ]]; then
    fix_manifests
  fi
  create_stack
  create_cluster

  SKIP_RUN_TESTS=${SKIP_RUN_TESTS:-""}
  if [[ -z "${SKIP_RUN_TESTS}" ]]; then
    run_tests
  fi
}

main "$@"
