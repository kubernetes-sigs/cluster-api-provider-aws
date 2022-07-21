# External Resource Garbage Collection

- **Feature status:** Experimental
- **Feature gate (required):** ExternalResourceGC=true

## Overview

Workload clusters that have been created by CAPA may have additional resources in AWS that need to be deleted when the cluster is deleted.

For example, if the workload cluster has `Services` of type `LoadBalancer` then there will be AWS ELB/NLB provisioned. If you try to delete the workload cluster in this example it will fail as the VPC is still being used by these load balancers.

This feature enables deleting these external resources as part of cluster deletion. It works by annotating the AWS infra cluster / control plane on creation. When a CAPI `Cluster` is requested to be deleted the deletion of CAPA resources is blocked depending on the status of this annotation. When the resources have been garbage collected (i.e. deleted) then the annotation is updated and normal CAPA deletion starts. This is not related to [externally managed infrastructure](https://cluster-api-aws.sigs.k8s.io/topics/bring-your-own-aws-infrastructure.html).

Currently we support cleaning up the following:

- AWS ELB/NLB - by deleting `Services` of type `LoadBalancer` from the workload cluster

We will look to potentially supporting deleting EBS volumes in the future.

> Note: this feature will likely be superseded by an upstream CAPI feature in the future when [this issue](https://github.com/kubernetes-sigs/cluster-api/issues/3075) is resolved.

## Enabling

To enable this you must set the `ExternalResourceGC` to `true` on the controller manager. The easiest way to do this is via an environment variable:

```bash
export EXP_EXTERNAL_RESOURCE_GC=true
clusterctl init --infrastructure aws
```

> Note: if you enable this feature **ALL** clusters will be marked as requiring garbage collection.

## Operations

### Manually Disabling Garbage Collection for a Cluster

There are 2 ways to manually disable garbage collection for an individual cluster:

#### Using `clusterawsadm`

You can disable garbage collection for a cluster by running the following command:

```bash
clusterawsadm gc disable --cluster-name mycluster
```

See the command help for more examples.

#### Editing `AWSCluster\AWSManagedControlPlane`

You can disable garbage collection for a cluster by editing your `AWSCluster` or `AWSManagedControlPlane`.

Either remove the `aws.cluster.x-k8s.io/external-resource-gc` or set its value to **true**.

```yaml
apiVersion: controlplane.cluster.x-k8s.io/v1beta1
kind: AWSManagedControlPlane
metadata:
  annotations:
    aws.cluster.x-k8s.io/external-resource-gc: "true"
```

### Manually Enablind Garbage Collection for a Cluster

There are 2 ways to manually enable garbage collection for an individual cluster:

#### Using `clusterawsadm`

You can enable garbage collection for a cluster by running the following command:

```bash
clusterawsadm gc enable --cluster-name mycluster
```

See the command help for more examples.

#### Editing `AWSCluster\AWSManagedControlPlane`

You can enable garbage collection for a cluster by editing your `AWSCluster` or `AWSManagedControlPlane`.

Add the `aws.cluster.x-k8s.io/external-resource-gc` annotation if it doesn't exist and set its value to **false**.

```yaml
apiVersion: controlplane.cluster.x-k8s.io/v1beta1
kind: AWSManagedControlPlane
metadata:
  annotations:
    aws.cluster.x-k8s.io/external-resource-gc: "false"
```
