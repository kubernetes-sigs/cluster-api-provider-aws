# Multi-AZ Control Plane

## Overview

By default, the control plane of a workload cluster created by CAPA will not span multiple availability zones (AZs) (also referred to as "failure domains"), even when using multiple control plane nodes. This is because CAPA will, by default, only create resources in the first discovered AZ of a region. However, if CAPA is specifically configured to create resources in multiple AZs, then control plane nodes will be automatically distributed across multiple AZs. This document explains how to instruct CAPA to create resources in multiple AZs.

## Configuring CAPA to Use Multiple AZs

To explicitly instruct CAPA to create resources in multiple AZs, users must add a `networkSpec` object to the AWSCluster specification. Here is an example `networkSpec` that creates resources across three AZs in the "us-west-2" region:

```yaml
spec:
  networkSpec:
    vpc:
      cidrBlock: 10.50.0.0/16
    subnets:
    - availabilityZone: us-west-2a
      cidrBlock: 10.50.0.0/20
      isPublic: true
    - availabilityZone: us-west-2a
      cidrBlock: 10.50.16.0/20
    - availabilityZone: us-west-2b
      cidrBlock: 10.50.32.0/20
      isPublic: true
    - availabilityZone: us-west-2b
      cidrBlock: 10.50.48.0/20
    - availabilityZone: us-west-2c
      cidrBlock: 10.50.64.0/20
      isPublic: true
    - availabilityZone: us-west-2c
      cidrBlock: 10.50.80.0/20
```

Specifying the CIDR block alone for the VPC is not enough; users must also supply a list of subnets that provides the desired AZ, the CIDR for the subnet, and whether the subnet is public (has a route to an Internet gateway) or is private (does not have a route to an Internet gateway).

Note that CAPA insists that there must be a public subnet (and associated Internet gateway), even if no public load balancer is requested for the control plane. Therefore, for every AZ where a control plane node should be placed, the `networkSpec` object must define both a public and private subnet.

Once CAPA is provided with a `networkSpec` that spans multiple AZs, the KubeadmControlPlane controller will automatically distribute control plane nodes across multiple AZs. No further configuration from the user is required.

## Caveats

Deploying control plane nodes across multiple AZs is not a panacea to cure all availability concerns. The sizing and overall utilization of the cluster will greatly affect the behavior of the cluster and the workloads hosted there in the event of an AZ failure. Careful planning is needed to maximize the availability of the cluster even in the face of an AZ failure. There are also other considerations, like cross-AZ traffic charges, that should be taken into account.
