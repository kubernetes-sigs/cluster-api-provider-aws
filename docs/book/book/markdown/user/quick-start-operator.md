# Cluster API Operator Quickstart

This section provides a quickstart guide for using the Cluster API Operator to create a Kubernetes cluster.
To use the `clusterctl` quickstart path, visit [this quickstart guide](./quick-start.md).

# Quickstart

This is a quickstart guide for getting Cluster API Operator up and running on your Kubernetes cluster.

For more detailed information, please refer to the full documentation.

## Prerequisites

- [Running Kubernetes cluster](https://cluster-api.sigs.k8s.io/user/quick-start#install-andor-configure-a-kubernetes-cluster).
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) for interacting with the management cluster.
- [Helm](https://helm.sh/docs/intro/install/) for installing operator on the cluster (optional).

## Install and configure Cluster API Operator

### Configuring credential for cloud providers

Instead of using environment variables as clusterctl does, Cluster API Operator uses Kubernetes secrets to store credentials for cloud providers. Refer to [provider documentation](https://cluster-api.sigs.k8s.io/user/quick-start#initialization-for-common-providers) on which credentials are required.

This example uses AWS provider, but the same approach can be used for other providers.

```bash
export CREDENTIALS_SECRET_NAME="credentials-secret"
export CREDENTIALS_SECRET_NAMESPACE="default"

kubectl create secret generic "${CREDENTIALS_SECRET_NAME}" --from-literal=AWS_B64ENCODED_CREDENTIALS="${AWS_B64ENCODED_CREDENTIALS}" --namespace "${CREDENTIALS_SECRET_NAMESPACE}"
```

### Installing Cluster API Operator

Add CAPI Operator & cert manager helm repository:

```bash
helm repo add capi-operator https://kubernetes-sigs.github.io/cluster-api-operator
helm repo add jetstack https://charts.jetstack.io --force-update
helm repo update
```

Install cert manager:

```bash
helm install cert-manager jetstack/cert-manager --namespace cert-manager --create-namespace --set installCRDs=true
```

Deploy Cluster API components with docker provider using a single command during operator installation

<aside class="note warning">

<h1> Warning </h1>

The `--wait` flag is REQUIRED for the helm install command to work. If the --wait flag is not used, the helm install command will not wait for the resources to be created and will return immediately. This will cause the helm install command to fail because the webhooks will not be ready in time. The --timeout flag is optional and can be used to specify the amount of time to wait for the resources to be created.

</aside>

```bash
helm install capi-operator capi-operator/cluster-api-operator --create-namespace -n capi-operator-system --set infrastructure.docker.enabled=true --set configSecret.name=${CREDENTIALS_SECRET_NAME} --set configSecret.namespace=${CREDENTIALS_SECRET_NAMESPACE}  --wait --timeout 90s
```

Docker provider can be replaced by any provider supported by [clusterctl](https://cluster-api.sigs.k8s.io/reference/providers.html#infrastructure).

Other options for installing Cluster API Operator are described in [full documentation](https://github.com/kubernetes-sigs/cluster-api-operator/blob/main/docs/README.md#installation).

# Example API Usage

Deploy latest version of core Cluster API components:

```yaml
apiVersion: operator.cluster.x-k8s.io/v1alpha2
kind: CoreProvider
metadata:
  name: cluster-api
  namespace: capi-system
```

Deploy Cluster API AWS provider with specific version, custom manager options and flags:

```yaml
---
apiVersion: operator.cluster.x-k8s.io/v1alpha2
kind: InfrastructureProvider
metadata:
 name: aws
 namespace: capa-system
spec:
 version: v2.1.4
 configSecret:
   name: credentials-secret
```
