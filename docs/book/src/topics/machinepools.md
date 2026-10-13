# MachinePools

- **Feature status:** Experimental
- **Feature gate:** MachinePool=true

MachinePool allows users to manage many machines as a single entity. Infrastructure providers implement a separate CRD that handles infrastructure side of the feature.

## AWSMachinePool

Cluster API Provider AWS (CAPA) has experimental support for `MachinePool` though the infrastructure type `AWSMachinePool`. An `AWSMachinePool` corresponds to an [AWS AutoScaling Groups](https://docs.aws.amazon.com/autoscaling/ec2/userguide/AutoScalingGroup.html), which provides the cloud provider specific resource for orchestrating a group of EC2 machines.

The AWSMachinePool controller creates and manages an AWS AutoScaling Group using launch templates so users don't have to manage individual machines. You can use Autoscaling health checks for replacing instances and it will maintain the number of instances specified.

### Using `clusterctl` to deploy

To deploy a MachinePool / AWSMachinePool via `clusterctl generate` there's a [flavor](https://cluster-api.sigs.k8s.io/clusterctl/commands/generate-cluster.html#flavors) for that.

Make sure to set up your AWS environment as described [here](https://cluster-api.sigs.k8s.io/user/quick-start.html).

```shell
export EXP_MACHINE_POOL=true
clusterctl init --infrastructure aws
clusterctl generate cluster my-cluster --kubernetes-version v1.25.0 --flavor machinepool > my-cluster.yaml
```

The template used for this [flavor](https://cluster-api.sigs.k8s.io/clusterctl/commands/generate-cluster.html#flavors) is located [here](https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/main/templates/cluster-template-machinepool.yaml).

## AWSManagedMachinePool

Cluster API Provider AWS (CAPA) has experimental support for [EKS Managed Node Groups](https://docs.aws.amazon.com/eks/latest/userguide/managed-node-groups.html) using `MachinePool` through the infrastructure type `AWSManagedMachinePool`. An `AWSManagedMachinePool` corresponds to an [AWS AutoScaling Groups](https://docs.aws.amazon.com/autoscaling/ec2/userguide/AutoScalingGroup.html) that is used for an EKS managed node group. .

The AWSManagedMachinePool controller creates and manages an EKS managed node group which in turn manages an AWS AutoScaling Group of managed EC2 instance types.

To use the managed machine pools certain IAM permissions are needed. The easiest way to ensure the required IAM permissions are in place is to use `clusterawsadm` to create them. To do this, follow the EKS instructions in [using clusterawsadm to fulfill prerequisites](using-clusterawsadm-to-fulfill-prerequisites.md).

### Using `clusterctl` to deploy

To deploy an EKS managed node group using AWSManagedMachinePool via `clusterctl generate` you can use a [flavor](https://cluster-api.sigs.k8s.io/clusterctl/commands/generate-cluster.html#flavors).

Make sure to set up your AWS environment as described [here](https://cluster-api.sigs.k8s.io/user/quick-start.html).

```shell
export EXP_MACHINE_POOL=true
clusterctl init --infrastructure aws
clusterctl generate cluster my-cluster --kubernetes-version v1.22.0 --flavor eks-managedmachinepool > my-cluster.yaml
```

The template used for this [flavor](https://cluster-api.sigs.k8s.io/clusterctl/commands/generate-cluster.html#flavors) is located [here](https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/main/templates/cluster-template-eks-managedmachinepool.yaml).

## Machine pools in a ClusterClass

A [ClusterClass](https://cluster-api.sigs.k8s.io/tasks/experimental-features/cluster-class/) can define machine pool classes under `spec.workers.machinePools`. Each class references a bootstrap template and an infrastructure template. CAPA provides two infrastructure templates:

| Template kind | Creates | Use it for |
|---|---|---|
| `AWSMachinePoolTemplate` | `AWSMachinePool` (AWS Auto Scaling Group) | Self-managed machine pools |
| `AWSManagedMachinePoolTemplate` | `AWSManagedMachinePool` (EKS managed node group) | EKS managed node groups |

The `spec.template.spec` of each template takes the same fields as the `spec` of the object it creates.

Requirements:

- The `ClusterTopology` feature gate must be enabled in Cluster API. It is disabled by default. With `clusterctl`, export `CLUSTER_TOPOLOGY=true` before running `clusterctl init`.
- The `MachinePool` feature gate must be enabled in Cluster API and CAPA. The CAPA manifest disables it unless you export `EXP_MACHINE_POOL=true` before running `clusterctl init`.

### Example: EKS managed node groups

The ClusterClass below defines one machine pool class, `default-pool`, backed by an EKS managed node group:

```yaml
apiVersion: cluster.x-k8s.io/v1beta2
kind: ClusterClass
metadata:
  name: eks-mmp
spec:
  controlPlane:
    templateRef:
      apiVersion: controlplane.cluster.x-k8s.io/v1beta2
      kind: AWSManagedControlPlaneTemplate
      name: eks-mmp-control-plane
  infrastructure:
    templateRef:
      apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
      kind: AWSManagedClusterTemplate
      name: eks-mmp
  workers:
    machinePools:
      - class: default-pool
        bootstrap:
          templateRef:
            apiVersion: bootstrap.cluster.x-k8s.io/v1beta2
            kind: NodeadmConfigTemplate
            name: eks-mmp-pool-0
        infrastructure:
          templateRef:
            apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
            kind: AWSManagedMachinePoolTemplate
            name: eks-mmp-pool-0
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
kind: AWSManagedMachinePoolTemplate
metadata:
  name: eks-mmp-pool-0
spec:
  template:
    spec:
      amiType: AL2023_x86_64_STANDARD
```

The `AWSManagedClusterTemplate`, `AWSManagedControlPlaneTemplate` and `NodeadmConfigTemplate` objects are left out here. A complete ClusterClass is in the flavors listed below.

A `Cluster` uses the class through `spec.topology`:

```yaml
apiVersion: cluster.x-k8s.io/v1beta2
kind: Cluster
metadata:
  name: my-cluster
spec:
  topology:
    classRef:
      name: eks-mmp
    version: v1.35.0
    workers:
      machinePools:
        - class: default-pool
          name: mp-0
          replicas: 3
```

Cluster API creates the `MachinePool`, the `AWSManagedMachinePool` and the `NodeadmConfig` for `mp-0` from the templates of `default-pool`.

### Flavors

| Flavor | Infrastructure template |
|---|---|
| `eks-clusterclass-machinepool` | `AWSMachinePoolTemplate` |
| `eks-clusterclass-managedmachinepool` | `AWSManagedMachinePoolTemplate` |

These flavors create the ClusterClass and its templates. They do not create a `Cluster`. Set `CLUSTER_CLASS_NAME`, `AWS_REGION`, `AWS_SSH_KEY_NAME` and `KUBERNETES_VERSION` before you generate them. The `eks-clusterclass-machinepool` flavor also needs `AWS_NODE_MACHINE_TYPE`.

### Template rules

- `spec.template.spec.providerIDList` cannot be set in either template. `spec.template.spec.providerID` cannot be set in an `AWSMachinePoolTemplate`.
- `AWSManagedMachinePoolTemplate` does not set a default `eksNodegroupName`. Each `AWSManagedMachinePool` created from the template gets its own generated node group name.
- An `AWSManagedMachinePoolTemplate` has the same immutable fields as an `AWSManagedMachinePool`: `eksNodegroupName` (once set), `subnetIDs`, `roleName` (once set), `diskSize`, `amiType`, `remoteAccess`, `capacityType` (once set), `availabilityZones`, `availabilityZoneSubnetType`, and whether `awsLaunchTemplate` is set together with its `name`.
- All fields of an `AWSMachinePoolTemplate` can be changed, as with an `AWSMachinePool`.

### Changing a machine pool class

When the ClusterClass changes, Cluster API updates the existing `AWSMachinePool` or `AWSManagedMachinePool` in place. Unlike machine deployments, it does not create a new infrastructure object. The `AWSManagedMachinePool` webhook therefore rejects a change to one of the immutable fields above. To change an immutable field, follow these steps:

1. Create a new `AWSManagedMachinePoolTemplate` with the new value.
2. Add a machine pool class that uses the new template.
3. Add a machine pool to `spec.topology.workers.machinePools` that uses the new class.
4. Remove the old machine pool once the new one is ready.

## Examples

### Example: MachinePool, AWSMachinePool and KubeadmConfig Resources

Below is an example of the resources needed to create a pool of EC2 machines orchestrated with
an AWS Auto Scaling Group.

```yaml
---
apiVersion: cluster.x-k8s.io/v1beta1
kind: MachinePool
metadata:
  name: capa-mp-0
spec:
  clusterName: capa
  replicas: 2
  template:
    spec:
      bootstrap:
        configRef:
          apiVersion: bootstrap.cluster.x-k8s.io/v1beta1
          kind: KubeadmConfig
          name: capa-mp-0
      clusterName: capa
      infrastructureRef:
        apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
        kind: AWSMachinePool
        name: capa-mp-0
      version: v1.25.0
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: AWSMachinePool
metadata:
  name: capa-mp-0
spec:
  minSize: 1
  maxSize: 10
  availabilityZones:
    - "${AWS_AVAILABILITY_ZONE}"
  awsLaunchTemplate:
    instanceType: "${AWS_CONTROL_PLANE_MACHINE_TYPE}"
    sshKeyName: "${AWS_SSH_KEY_NAME}"
  subnets:
    - id : "${AWS_SUBNET_ID}"
---
apiVersion: bootstrap.cluster.x-k8s.io/v1beta1
kind: KubeadmConfig
metadata:
  name: capa-mp-0
  namespace: default
spec:
  joinConfiguration:
    nodeRegistration:
      name: '{{ ds.meta_data.local_hostname }}'
      kubeletExtraArgs:
        cloud-provider: aws
```

## Autoscaling

[`cluster-autoscaler`](https://github.com/kubernetes/autoscaler/tree/master/cluster-autoscaler) can be used to scale MachinePools up and down.
Two providers are possible to use with CAPA MachinePools: `clusterapi`, or `aws`.

If the `AWS` autoscaler provider is used, each MachinePool needs to have an annotation set to prevent scale up/down races between
cluster-autoscaler and cluster-api. Example:

```yaml
apiVersion: cluster.x-k8s.io/v1beta1
kind: MachinePool
metadata:
  name: capa-mp-0
  annotations:
    cluster.x-k8s.io/replicas-managed-by: "external-autoscaler"
spec:
  clusterName: capa
  replicas: 2
  template:
    spec:
      bootstrap:
        configRef:
          apiVersion: bootstrap.cluster.x-k8s.io/v1beta1
          kind: KubeadmConfig
          name: capa-mp-0
      clusterName: capa
      infrastructureRef:
        apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
        kind: AWSMachinePool
        name: capa-mp-0
      version: v1.25.0
```

When using GitOps, make sure to ignore differences in `spec.replicas` on MachinePools. Example when using ArgoCD:

```yaml
  ignoreDifferences:
    - group: cluster.x-k8s.io
      kind: MachinePool
      jsonPointers:
        - /spec/replicas
```

## Machine pool machines

With the feature gate `MachinePoolMachines=true`, you can enable creation of `Machine`/`AWSMachine` objects for nodes created by a `AWSMachinePool`. This is experimental and will be used to introduce features such as per-node health checks.
