# External Resource Garbage Collection

- **Feature status:** Experimental
- **Feature gate (required):** ExternalResourceGC=true

## Overview

Workload clusters that have been created by CAPA may have additional resources in AWS that need to be deleted when the cluster is deleted. 

For example, if the workload cluster has `Services` of type `LoadBalancer` then there will be AWS ELB/NLB provisioned. If you try to delete the workload cluster in this example it will fail as the VPC is still being used by these load balancers.

This feature enables deleting these external resources as part of cluster deletion. It works by annotating the AWS infra cluster / control plane on creation. When a CAPI `Cluster` is requested to be deleted the deletion of CAPA resources is blocked depending on the status of this annotation. When the resources have been garbage collected (i.e. deleted) then the annotation is updated and normal CAPA deletion starts.

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
