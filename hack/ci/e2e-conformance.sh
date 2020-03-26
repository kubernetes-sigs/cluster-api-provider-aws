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

# hack script for running a cluster-api-provider-aws e2e

set -o errexit -o nounset -o pipefail

REGISTRY=${REGISTRY:-"gcr.io/"$(gcloud config get-value project)}
AWS_REGION=${AWS_REGION:-"us-east-1"}
CLUSTER_NAME=${CLUSTER_NAME:-"test-$(date +%s)"}
AWS_SSH_KEY_NAME=${AWS_SSH_KEY_NAME:-"${CLUSTER_NAME}-key"}
KUBERNETES_VERSION=${KUBERNETES_VERSION:-"v1.17.3"}

ARTIFACTS="${ARTIFACTS:-${PWD}/_artifacts}"
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"

# dump logs from kind and all the nodes
dump-logs() {
  # log version information
  echo "=== versions ==="
  echo "kind : $(kind version)" || true
  echo "bootstrap cluster:"
  kubectl version || true
  echo "deployed cluster:"
  kubectl --kubeconfig="${PWD}"/kubeconfig version || true
  echo ""

  # dump all the info from the CAPI related CRDs
  kubectl get \
  clusters,awsclusters,machines,awsmachines,kubeadmconfigs,machinedeployments,awsmachinetemplates,kubeadmconfigtemplates,machinesets,kubeadmcontrolplanes \
  --all-namespaces -o yaml >> "${ARTIFACTS}/logs/capa.info" || true

  # dump images info
  echo "images in docker" >> "${ARTIFACTS}/logs/images.info"
  docker images >> "${ARTIFACTS}/logs/images.info"
  echo "images from bootstrap using containerd CLI" >> "${ARTIFACTS}/logs/images.info"
  docker exec clusterapi-control-plane ctr -n k8s.io images list >> "${ARTIFACTS}/logs/images.info" || true
  echo "images in bootstrap cluster using kubectl CLI" >> "${ARTIFACTS}/logs/images.info"
  (kubectl get pods --all-namespaces -o json \
   | jq --raw-output '.items[].spec.containers[].image' | sort)  >> "${ARTIFACTS}/logs/images.info" || true
  echo "images in deployed cluster using kubectl CLI" >> "${ARTIFACTS}/logs/images.info"
  (kubectl --kubeconfig="${PWD}"/kubeconfig get pods --all-namespaces -o json \
   | jq --raw-output '.items[].spec.containers[].image' | sort)  >> "${ARTIFACTS}/logs/images.info" || true

  # dump cluster info for kind
  kubectl cluster-info dump > "${ARTIFACTS}/logs/kind-cluster.info" || true

  # dump cluster info for kind
  echo "=== aws ec2 describe-instances ===" >> "${ARTIFACTS}/logs/capa-cluster.info" || true
  aws ec2 describe-instances --region "${AWS_REGION}" >> "${ARTIFACTS}/logs/capa-cluster.info" || true
  echo "=== cluster-info dump ===" >> "${ARTIFACTS}/logs/capa-cluster.info" || true
  kubectl --kubeconfig="${PWD}"/kubeconfig cluster-info dump >> "${ARTIFACTS}/logs/capa-cluster.info" || true

  # export all logs from kind
  kind "export" logs --name="clusterapi" "${ARTIFACTS}/logs" || true

  node_filters="Name=tag:sigs.k8s.io/cluster-api-provider-aws/cluster/${CLUSTER_NAME},Values=owned"
  bastion_filters="${node_filters} Name=tag:sigs.k8s.io/cluster-api-provider-aws/role,Values=bastion"
  jump_node=$(aws ec2 describe-instances --region "$AWS_REGION" --filters "${bastion_filters}" --query "Reservations[*].Instances[*].PublicIpAddress" --output text | head -1)

  # We used to pipe this output to 'tail -n +2' but for some reason this was sometimes (all the time?) only finding the
  # bastion host. For now, omit the tail and gather logs for all VMs that have a private IP address. This will include
  # the bastion, but that's better than not getting logs from all the VMs.
  for node in $(aws ec2 describe-instances --region "$AWS_REGION" --filters "${node_filters}" --query "Reservations[*].Instances[*].PrivateIpAddress" --output text)
  do
    echo "collecting logs from ${node} using jump host ${jump_node}"
    dir="${ARTIFACTS}/logs/${node}"
    mkdir -p "${dir}"
    ssh-to-node "${node}" "${jump_node}" "sudo journalctl --output=short-precise -k" > "${dir}/kern.log" || true
    ssh-to-node "${node}" "${jump_node}" "sudo journalctl --output=short-precise" > "${dir}/systemd.log" || true
    ssh-to-node "${node}" "${jump_node}" "sudo crictl version && sudo crictl info" > "${dir}/containerd.info" || true
    ssh-to-node "${node}" "${jump_node}" "sudo journalctl --no-pager -u cloud-final" > "${dir}/cloud-final.log" || true
    ssh-to-node "${node}" "${jump_node}" "sudo journalctl --no-pager -u kubelet.service" > "${dir}/kubelet.log" || true
    ssh-to-node "${node}" "${jump_node}" "sudo journalctl --no-pager -u containerd.service" > "${dir}/containerd.log" || true
  done
}

# SSH to a node by name ($1) via jump server ($2) and run a command ($3).
function ssh-to-node() {
  local node="$1"
  local jump="$2"
  local cmd="$3"

  ssh_key_pem="/tmp/${AWS_SSH_KEY_NAME}.pem"
  ssh_params="-o LogLevel=quiet -o ConnectTimeout=30 -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no"
  scp "$ssh_params" -i "$ssh_key_pem" "$ssh_key_pem" "ubuntu@${jump}:$ssh_key_pem"
  ssh "$ssh_params" -i "$ssh_key_pem" \
    -o "ProxyCommand ssh $ssh_params -W %h:%p -i $ssh_key_pem ubuntu@${jump}" \
    ubuntu@"${node}" "${cmd}"
}

# cleanup all resources we use
cleanup() {
  # KIND_IS_UP is true once we: kind create
  if [[ "${KIND_IS_UP:-}" = true ]]; then
    timeout 600 kubectl \
      delete cluster "${CLUSTER_NAME}" || true
     timeout 600 kubectl \
      wait --for=delete cluster/"${CLUSTER_NAME}" || true
    make kind-reset || true
  fi
  # clean up e2e.test symlink
  (cd "$(go env GOPATH)/src/k8s.io/kubernetes" && rm -f _output/bin/e2e.test) || true

  delete_key_pair
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
    if [[ -n "$image" ]]; then
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
    wget --quiet -O packer.zip "$url"  && \
      unzip packer.zip && \
      rm packer.zip && \
      ln -s "$PWD"/packer /usr/local/bin/packer
  fi

  tracestate="$(shopt -po xtrace)"
  set +o xtrace

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
    # Ensure go is available in PATH
    ln -s /usr/local/go/bin/go /usr/bin/go
    # install goss plugin
    su - packer -c "bash -c 'cd /home/prow/go/src/sigs.k8s.io/image-builder/images/capi/packer/ami && make plugins'"
    # use the packer user to run the build
    su - packer -c "bash -c 'cd /home/prow/go/src/sigs.k8s.io/image-builder/images/capi && AWS_REGION=$AWS_REGION AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID:-""} AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY:-""} make build-ami-default'"
  fi

  eval "$tracestate"
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

  if [[ -z ${IMAGE_ID:-} ]]; then
    # default lookup org hould be the same as defaultMachineAMIOwnerID
    IMAGE_LOOKUP_ORG=${IMAGE_LOOKUP_ORG:-"258751437250"}
    filter="capa-ami-ubuntu-18.04-1.17.*"
    image_id=$(aws ec2 describe-images --query 'Images[*].[ImageId,Name]' \
      --filters "Name=name,Values=$filter" "Name=owner-id,Values=$IMAGE_LOOKUP_ORG" \
      --region "${AWS_REGION}" --output json | jq -r '.[0][0] | select (.!=null)')
    if [[ -z "$image_id" ]]; then
      echo "unable to find image using : $filter $IMAGE_LOOKUP_ORG ... bailing out!"
      exit 1
    fi
  else
    image_id=$(aws ec2 describe-images --image-ids "$IMAGE_ID" \
      --query 'Images[*].[ImageId,Name]' --output json | jq -r '.[0][0] | select (.!=null)')
    echo "using specified image id : ${IMAGE_ID}"
    if [[ -z "$image_id" ]]; then
      echo "unable to find image using id : $IMAGE_ID ... bailing out!"
      exit 1
    fi
  fi


  # Enable the bits to inject a script that can pull newer versions of kubernetes
  if [[ -n ${CI_VERSION:-} || -n ${USE_CI_ARTIFACTS:-} ]]; then
    if ! grep -i -wq "patchesStrategicMerge" "templates/kustomization.yaml"; then
      echo "patchesStrategicMerge:" >> "templates/kustomization.yaml"
      echo "- kustomizeversions.yaml" >> "templates/kustomization.yaml"
    fi
  fi

  PULL_POLICY=IfNotPresent \
    make modules docker-build clusterawsadm
}

# install cloud formation templates, iam objects etc
create_stack() {
  "${REPO_ROOT}/bin/clusterawsadm" alpha bootstrap create-stack
}

create_key_pair() {
  (aws ec2 create-key-pair --key-name "${AWS_SSH_KEY_NAME}" --region "${AWS_REGION}" > /tmp/keypair-"${AWS_SSH_KEY_NAME}".json \
   && KEY_PAIR_CREATED="true" \
   && jq -r '.KeyMaterial' /tmp/keypair-"${AWS_SSH_KEY_NAME}".json > /tmp/"${AWS_SSH_KEY_NAME}".pem \
   && chmod 600 /tmp/"${AWS_SSH_KEY_NAME}".pem) || true
}

delete_key_pair() {
  # Delete only if we created it
  if [[ "${KEY_PAIR_CREATED:-}" = true ]]; then
    aws ec2 delete-key-pair --key-name "${AWS_SSH_KEY_NAME}" --region "${AWS_REGION}" || true
    rm /tmp/keypair-"${AWS_SSH_KEY_NAME}".json || true
    rm /tmp/"${AWS_SSH_KEY_NAME}".pem || true
  fi
}

# up a cluster with kind
create_cluster() {
  # actually create the cluster
  KIND_IS_UP=true

  tracestate="$(shopt -po xtrace)"
  set +o xtrace

  if [[ -n ${USE_CI_ARTIFACTS:-} ]]; then
    # TODO: revert to https://dl.k8s.io/ci/latest-green.txt once https://github.com/kubernetes/release/issues/897 is fixed.
    CI_VERSION=${CI_VERSION:-$(curl -sSL https://dl.k8s.io/ci/k8s-master.txt)}
  fi

  # Load the newly built image into kind and start the cluster
  (AWS_CREDENTIALS=$(aws iam create-access-key --user-name bootstrapper.cluster-api-provider-aws.sigs.k8s.io) \
  AWS_ACCESS_KEY_ID=$(echo "$AWS_CREDENTIALS" | jq .AccessKey.AccessKeyId -r) \
  AWS_SECRET_ACCESS_KEY=$(echo "$AWS_CREDENTIALS" | jq .AccessKey.SecretAccessKey -r) \
  AWS_REGION=${AWS_REGION} \
  CONTROL_PLANE_MACHINE_COUNT=1 \
  WORKER_MACHINE_COUNT=2 \
  KUBERNETES_VERSION=${KUBERNETES_VERSION} \
  CI_VERSION=${CI_VERSION:-} \
  IMAGE_ID=${image_id} \
  AWS_SSH_KEY_NAME=$AWS_SSH_KEY_NAME \
  AWS_CONTROL_PLANE_MACHINE_TYPE=m5.large \
  AWS_NODE_MACHINE_TYPE=m5.large \
  AWS_B64ENCODED_CREDENTIALS=$("${REPO_ROOT}"/bin/clusterawsadm alpha bootstrap encode-aws-credentials) \
  LOAD_IMAGE="${REGISTRY}/cluster-api-aws-controller-amd64:dev" CLUSTER_NAME="${CLUSTER_NAME}" \
    make create-cluster)

  eval "$tracestate"

  # Wait till all machines are running (bail out at 30 mins)
  attempt=0
  while true; do
    kubectl get machines
    read running total <<< $(kubectl get machines \
      -o json | jq -r '.items[].status.phase' | awk 'BEGIN{count=0} /(r|R)unning/{count++} END{print count " " NR}') ;
    if [[ $total == "3" && $running == "3" ]]; then
      return 0
    fi
    read failed total <<< $(kubectl get machines \
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
  NUM_NODES="$(kubectl get nodes --kubeconfig="$KUBECONFIG" \
    -o=jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.spec.taints}{"\n"}{end}' \
    | grep -cv "node-role.kubernetes.io/master" )"

  # wait for all the nodes to be ready
  kubectl wait --for=condition=Ready node --kubeconfig="$KUBECONFIG" --all || true

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
  for arg in "$@"
  do
    if [[ "$arg" == "--verbose" ]]; then
      set -o xtrace
    fi
    if [[ "$arg" == "--clean" ]]; then
      cleanup
      return 0
    fi
    if [[ "$arg" == "--use-ci-artifacts" ]]; then
      USE_CI_ARTIFACTS="1"
    fi
    if [[ "$arg" == "--skip-init-image" ]]; then
      SKIP_INIT_IMAGE="1"
    fi
  done

  if [[ -z "$AWS_ACCESS_KEY_ID" ]]; then
    cat <<EOF
AWS_ACCESS_KEY_ID is not set.
EOF
    return 2
  fi
  if [[ -z "$AWS_SECRET_ACCESS_KEY" ]]; then
    cat <<EOF
AWS_SECRET_ACCESS_KEY is not set.
EOF
    return 2
  fi
  if [[ -z "$AWS_REGION" ]]; then
    cat <<EOF
AWS_REGION is not set.
Please specify which the AWS region to use.
EOF
    return 2
  fi

  # create temp dir and setup cleanup
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
  generate_manifests
  if [[ -n "${SKIP_INIT_IMAGE:-}" ]]; then
    echo "Skipping image initialization..."
  else
    init_image
  fi

  create_stack
  create_key_pair
  create_cluster

  if [[ -z "${SKIP_RUN_TESTS:-}" ]]; then
    run_tests
  fi
}

main "$@"
