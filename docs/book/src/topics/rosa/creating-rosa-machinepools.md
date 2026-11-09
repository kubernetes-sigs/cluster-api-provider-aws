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
  instanceType: "m7i.xlarge"
  subnet: "${PRIVATE_SUBNET_ID}"
  version: "${OPENSHIFT_VERSION}"
```

**Note:** `m7i.xlarge` is not available in all AWS regions. Choose an instance type available in your target region.

## Spot Instances

Setting `spotMarketOptions` runs the MachinePool's nodes on AWS Spot instances, which cost less than on-demand instances but can be reclaimed by AWS at any time.

The simplest form bids the on-demand price, so you pay no more than you would for on-demand capacity while still benefiting from lower Spot pricing:

```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
kind: ROSAMachinePool
metadata:
  name: "${CLUSTER_NAME}-pool-0"
spec:
  nodePoolName: "nodepool-0"
  instanceType: "m7i.xlarge"
  subnet: "${PRIVATE_SUBNET_ID}"
  version: "${OPENSHIFT_VERSION}"
  spotMarketOptions: {}
```

To cap the hourly bid instead, set `maxPrice`:

```yaml
  spotMarketOptions:
    maxPrice: "0.5"
```

`maxPrice` accepts either a price in USD per hour or the string `"on-demand"`, and defaults to `"on-demand"` when omitted. A numeric price must be greater than zero, so `"0"` and `"0.0"` are rejected. Setting `maxPrice` below the current Spot price means the node pool will not get capacity.

### Constraints

`spotMarketOptions` cannot be added to or removed from an existing MachinePool. Whether a node pool uses Spot instances is fixed when it is created; switching requires replacing the MachinePool.

It also cannot be combined with `capacityReservationID`, since a capacity reservation applies to on-demand instances.

### Changing the price later

`maxPrice` is the one part of `spotMarketOptions` that can be changed after creation, and only in one direction:

| Change | Supported |
| --- | --- |
| `"on-demand"` → a price (e.g. `"0.5"`) | Yes |
| One price → another (e.g. `"0.5"` → `"0.4"`) | Yes |
| A price → `"on-demand"` | **No** |

Once a numeric `maxPrice` is set it cannot be returned to `"on-demand"` — the ROSA API rejects clearing the bid, so the change is refused by the `ROSAMachinePool` webhook rather than failing later during reconciliation. Removing the `maxPrice` field is treated the same way as setting it back to `"on-demand"`. To go back to on-demand pricing, replace the MachinePool.

see [ROSAMachinePool CRD Reference](https://cluster-api-aws.sigs.k8s.io/crd/#infrastructure.cluster.x-k8s.io/v1beta2.ROSAMachinePool) for all possible configurations.
