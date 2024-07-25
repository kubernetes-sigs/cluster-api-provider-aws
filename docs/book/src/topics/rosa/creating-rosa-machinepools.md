# Creating MachinePools

Cluster API Provider AWS (CAPA) has experimental support for managed ROSA MachinePools through the infrastructure type `ROSAMachinePool`. A `ROSAMachinePool` is responsible for orchestrating and bootstraping a group of EC2 machines into kubernetes nodes.

### Using `clusterctl` to deploy

To deploy a MachinePool / ROSAMachinePool via `clusterctl generate` use the template located [here](https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/main/templates/cluster-template-rosa-machinepool.yaml).

Make sure to set up your environment as described [here](./creating-a-cluster.md#creating-the-cluster).

```shell
clusterctl generate cluster my-cluster --from templates/cluster-template-rosa-machinepool > my-cluster.yaml
```

## Example

Below is an example of the resources needed to create a ROSA MachinePool.

```yaml
---
apiVersion: cluster.x-k8s.io/v1beta1
kind: MachinePool
metadata:
  name: "${CLUSTER_NAME}-pool-0"
spec:
  clusterName: "${CLUSTER_NAME}"
  replicas: 1
  template:
    spec:
      clusterName: "${CLUSTER_NAME}"
      bootstrap:
        dataSecretName: ""
      infrastructureRef:
        name: "${CLUSTER_NAME}-pool-0"
        apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
        kind: ROSAMachinePool
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
kind: ROSAMachinePool
metadata:
  name: "${CLUSTER_NAME}-pool-0"
spec:
  nodePoolName: "nodepool-0"
  instanceType: "m5.xlarge"
  subnet: "${PRIVATE_SUBNET_ID}"
  version: "${OPENSHIFT_VERSION}"
```

see [ROSAMachinePool CRD Reference](https://cluster-api-aws.sigs.k8s.io/crd/#infrastructure.cluster.x-k8s.io/v1beta2.ROSAMachinePool) for all possible configurations.
