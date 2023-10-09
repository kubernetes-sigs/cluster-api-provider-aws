# Setting up a Network Load Balancer

## Overview

It's possible to set up and use a [Network Load Balancer](https://docs.aws.amazon.com/elasticloadbalancing/latest/network/introduction.html) with `AWSCluster` instead of the
Classic Load Balancer that is created by default.

## `AWSCluster` setting

To make CAPA create a network load balancer simply set the load balancer type to `network` like this:

```yaml
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
kind: AWSCluster
metadata:
  name: "test-aws-cluster"
spec:
  region: "eu-central-1"
  controlPlaneLoadBalancer:
    loadBalancerType: nlb
```

This will create the following objects:

- A network load balancer
- Listeners
- A target group

It will also take into consideration IPv6 enabled clusters and create an IPv6 aware load balancer.

## Preserve Client IPs

By default, client ip preservation is disabled. This is to avoid [hairpinning](https://docs.aws.amazon.com/elasticloadbalancing/latest/network/load-balancer-troubleshooting.html#loopback-timeout) issues between kubelet and the node
registration process. To enable client IP preservation, you can set it to enable with the following flag:

```yaml
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
kind: AWSCluster
metadata:
  name: "test-aws-cluster"
spec:
  region: "eu-central-1"
  sshKeyName: "capa-key"
  controlPlaneLoadBalancer:
    loadBalancerType: nlb
    preserveClientIP: true
```

## Security

NLBs can use security groups, but only if one is associated at the time of creation.
CAPA will associate the default control plane security groups with a new NLB by default.

For more information, see AWS's [Network Load Balancer and Security Groups](https://docs.aws.amazon.com/elasticloadbalancing/latest/network/load-balancer-security-groups.html) documentation.

## Extension of the code

Right now, only NLBs and a Classic Load Balancer is supported. However, the code has been written in a way that it
should be easy to extend with an ALB or a GLB.
