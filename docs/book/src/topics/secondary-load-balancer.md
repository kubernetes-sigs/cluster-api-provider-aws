# Enabling a Secondary Control Plane Load Balancer

## Overview

It is possible to use a second control plane load balancer within a CAPA cluster.
This secondary control plane load balancer is primarily meant to be used for internal cluster traffic, for use cases where traffic between nodes and pods should be kept internal to the VPC network.
This adds a layer of privacy to traffic, as well as potentially saving on egress costs for traffic to the Kubernetes API server.

A dual load balancer topology is not used as a default in order to maintain backward compatibility with existing CAPA clusters.

## Requirements and defaults

- A secondary control plane load balancer is _not_ created by default.
- The secondary control plane load balancer _must_ be a [Network Load Balancer](https://docs.aws.amazon.com/elasticloadbalancing/latest/network/introduction.html), and will default to this type.
- The secondary control plane load balancer must also be provided a name.
- The secondary control plane's `Scheme` defaults to `internal`, and _must_ be different from the `spec.controlPlaneLoadBalancer`'s `Scheme`.

The secondary load balancer will use the same Security Group information as the primary control plane load balancer.

## Creating a secondary load balancer

To create a secondary load balancer, add the `secondaryControlPlaneLoadBalancer` stanza to your `AWSCluster`.

```yaml
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
kind: AWSCluster
metadata:
  name: test-aws-cluster
spec:
  region: us-east-2
  sshKeyName: nrb-default
  secondaryControlPlaneLoadBalancer:
    name: internal-apiserver
    scheme: internal     # optional
```

## Target Group IP Address Type

**Note:** The `targetGroupIPType` field is only available when using Network Load Balancers (NLB), Application Load Balancers (ALB), or Gateway Load Balancers (ELB). It **cannot** be configured when using Classic Load Balancers.

The secondary load balancer supports the same `targetGroupIPType` configuration as the primary load balancer. By default, the target group IP address type is set based on the VPC configuration:
- If the VPC has IPv6 enabled, the target group uses `ipv6`
- Otherwise, it defaults to `ipv4`

You can explicitly configure the IP address type for the secondary load balancer's target group:

```yaml
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
kind: AWSCluster
metadata:
  name: test-aws-cluster
spec:
  region: us-east-2
  sshKeyName: nrb-default
  secondaryControlPlaneLoadBalancer:
    name: internal-apiserver
    scheme: internal
    targetGroupIPType: ipv6
```

Valid values are:
- `ipv4`: Routes traffic to targets using IPv4 addresses
- `ipv6`: Routes traffic to targets using IPv6 addresses
