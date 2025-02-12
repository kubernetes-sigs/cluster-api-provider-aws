# External Resource Garbage Collection

- **Feature status:** Stable
- **Feature gate (required):** ExternalResourceGC=true

## Overview

Workload clusters that CAPA has created may have additional resources in AWS that need to be deleted when the cluster is deleted.

For example, if the workload cluster has `Services` of type `LoadBalancer` then AWS ELB/NLB are provisioned. If you try to delete the workload cluster in this example, it will fail as these load balancers are still using the VPC.

This feature enables deletion of these external resources as part of cluster deletion. During the deletion of a workload cluster the external AWS resources that where created by the Cloud Controller Manager (CCM) in the workload cluster will be identified and deleted.

> NOTE: This is not related to [externally managed infrastructure](https://cluster-api-aws.sigs.k8s.io/topics/bring-your-own-aws-infrastructure.html).

Currently, we support cleaning up the following:

- AWS ELB/NLB - by deleting `Services` of type `LoadBalancer` from the workload cluster

We will look to support deleting EBS volumes in the future potentially.

> Note: this feature will likely be superseded by an upstream CAPI feature in the future when [this issue](https://github.com/kubernetes-sigs/cluster-api/issues/3075) is resolved.

## Disabling

The garbage collection feature is enabled by default. If you want to disable the feature then you must set the `ExternalResourceGC` feature gate to `false` on the controller manager. The easiest way to do this is via an environment variable:

```bash
export EXTERNAL_RESOURCE_GC=false
clusterctl init --infrastructure aws
```

> Note: if you disablw this feature **ALL** clusters will be marked as not requiring garbage collection.

## Operations

### Manually Disabling Garbage Collection for a Cluster

There are 2 ways to manually disable garbage collection for an individual cluster:

#### Using `clusterawsadm`

By running the following command:

```bash
clusterawsadm gc disable --cluster-name mycluster
```

See the command help for more examples.

#### Editing `AWSCluster\AWSManagedControlPlane`

Or, by editing your `AWSCluster` or `AWSManagedControlPlane` so that the annotation `aws.cluster.x-k8s.io/external-resource-gc` is set to **false**.

```yaml
apiVersion: controlplane.cluster.x-k8s.io/v1beta1
kind: AWSManagedControlPlane
metadata:
  annotations:
    aws.cluster.x-k8s.io/external-resource-gc: "false"
```

### Manually Enabling Garbage Collection for a Cluster

There are 2 ways to manually enable garbage collection for an individual cluster:

#### Using `clusterawsadm`

By running the following command:

```bash
clusterawsadm gc enable --cluster-name mycluster
```

See the command help for more examples.

#### Editing `AWSCluster\AWSManagedControlPlane`

Or, by editing your `AWSCluster` or `AWSManagedControlPlane` o that the annotation `aws.cluster.x-k8s.io/external-resource-gc` is either removed or set to **true**.

```yaml
apiVersion: controlplane.cluster.x-k8s.io/v1beta1
kind: AWSManagedControlPlane
metadata:
  annotations:
    aws.cluster.x-k8s.io/external-resource-gc: "true"
```
