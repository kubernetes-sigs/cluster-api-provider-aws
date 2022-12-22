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
