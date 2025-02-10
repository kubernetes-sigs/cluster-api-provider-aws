#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

# Load cluster resource configuration 
source $DIR/cluster.envrc

# Create management cluster
kind create cluster
kubectl cluster-info

# Check AWS Authentication and Account Preparation
aws sts get-caller-identity

export AWS_SSH_KEY_NAME=${AWS_SSH_KEY_NAME:-"capa-dedicated-hosts"}
if aws ec2 describe-key-pairs --key-names $AWS_SSH_KEY_NAME 2>/dev/null; then
    echo "Key pair [$AWS_SSH_KEY_NAME] already exists."
else
    KEYPAIR_FILE="$CLUSTER_DIR/${AWS_SSH_KEY_NAME}.pem"
    aws ec2 create-key-pair --key-name $AWS_SSH_KEY_NAME --query 'KeyMaterial' --output text > $KEYPAIR_FILE
    chmod 400 $KEYPAIR_FILE
    ls -l $KEYPAIR_FILE
    echo "Key pair [$AWS_SSH_KEY_NAME] created [$KEYPAIR_FILE]."

fi

clusterawsadm bootstrap iam create-cloudformation-stack

export AWS_B64ENCODED_CREDENTIALS=$(clusterawsadm bootstrap credentials encode-as-profile)
echo $AWS_B64ENCODED_CREDENTIALS 

# Initialize the management cluster with AWS provider
clusterctl init --infrastructure aws

# Create a workload cluster
CLUSTER_DIR=${CLUSTER_DIR:-"tmp"}
mkdir -p $CLUSTER_DIR

export AWS_CONTROL_PLANE_MACHINE_TYPE=t3.large
export AWS_NODE_MACHINE_TYPE=t3.large

export AWS_HOST_AZ="us-east-1a"
export AWS_HOST_FAMILY="t3"

# Allocate dedicated host
aws ec2 allocate-hosts \
    --availability-zone "$AWS_HOST_AZ" \
    --auto-placement "off" \
    --host-recovery "off" \
    --host-maintenance "on" \
    --quantity 1 \
    --instance-family "$AWS_HOST_FAMILY" | tee "$CLUSTER_DIR/host.json"

export AWS_HOST_ID=$(jq -r '.HostIds[0]' "$CLUSTER_DIR/host.json")
export AWS_HOST_AFFINITY="Default"
echo $AWS_HOST_ID

export KUBERNETES_VERSION_DEFAULT=$(clusterawsadm ami list -o json | jq -r '.items[0].spec.kubernetesVersion')
export KUBERNETES_VERSION=${KUBERNETES_VERSION:-$KUBERNETES_VERSION_DEFAULT}
echo $KUBERNETES_VERSION

export CLUSTER_NAME=${CLUSTER_NAME:-"capa-dedicated-hosts"}

export AWS_HOST_AZ="us-east-1a"
export AWS_HOST_FAMILY="t3"

# Allocate dedicated host
aws ec2 allocate-hosts \
    --availability-zone "$AWS_HOST_AZ" \
    --auto-placement "off" \
    --host-recovery "off" \
    --host-maintenance "on" \
    --quantity 1 \
    --instance-family "$AWS_HOST_FAMILY" | tee "$CLUSTER_DIR/host.json"

clusterctl generate cluster $CLUSTER_NAME \
    --from - \
    --kubernetes-version $KUBERNETES_VERSION \
    --control-plane-machine-count=3 \
    --worker-machine-count=3 \
    < templates/cluster-template-dedicated-hosts.yaml \
    > "$CLUSTER_DIR/capa-dedicated-hosts.yaml"

kubectl apply -f "$CLUSTER_DIR/capa-dedicated-hosts.yaml"

kubectl get cluster

watch -n 15 clusterctl describe cluster capa-dedicated-hosts

kubectl get kubeadmcontrolplane

# Function to check if kubeadmcontrolplane is initialized
check_initialized() {
    kubectl get kubeadmcontrolplane -o json | jq -e '.items[] | select(.status.initialized == true)' > /dev/null 2>&1
}

# Loop until the kubeadmcontrolplane is initialized
while true; do
    if check_initialized; then
        echo "kubeadmcontrolplane is initialized."
        break
    else
        echo "Waiting for kubeadmcontrolplane to be initialized..."
        sleep 30
    fi
done

echo "Fetching workload cluster kubeconfig"
WORKLOAD_KUBECONFIG="$CLUSTER_DIR/capa-dedicated-hosts.kubeconfig"
clusterctl get kubeconfig capa-dedicated-hosts > "$WORKLOAD_KUBECONFIG"

# Authenticate on docker hub
kubectl create secret docker-registry docker-creds \
    --docker-server='https://index.docker.io/v1/' \
    --docker-username=$DOCKER_USERNAME \
    --docker-password=$DOCKER_PASSWORD \
    --docker-email=$DOCKER_EMAIL 


echo "Installing Calico CNI"
helm repo add projectcalico https://docs.tigera.io/calico/charts \
    --kubeconfig=$WORKLOAD_KUBECONFIG

helm install calico projectcalico/tigera-operator \
    --kubeconfig=$WORKLOAD_KUBECONFIG \
    -f https://raw.githubusercontent.com/kubernetes-sigs/cluster-api-provider-azure/main/templates/addons/calico/values.yaml \
    --namespace tigera-operator \
    --create-namespace


# Patch Calico to use Docker Hub credentials
kubectl --kubeconfig=$WORKLOAD_KUBECONFIG patch daemonset \
    -n kube-system calico-node \
    -p '{"spec":{"template":{"spec":{"imagePullSecrets":[{"name":"docker-creds"}]}}}}'


# Verify that the workload cluster is up and running

kubectl --kubeconfig=$WORKLOAD_KUBECONFIG cluster-info
kubectl --kubeconfig=$WORKLOAD_KUBECONFIG get nodes
kubectl --kubeconfig=$WORKLOAD_KUBECONFIG get pods -n kube-system

