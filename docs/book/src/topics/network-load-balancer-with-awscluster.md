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

NLBs cannot use Security Groups. Therefore, the following steps have been taken to increase security for nodes
communication. NLBs need access to the node in order to send traffic its way. A port has to be opened using an ip
address range instead of a security group as a _source_. There are two scenarios and CIDRs that can be enabled.

First, if client IP preservation is _disabled_ we only add the VPC's private CIDR range as allowed source for the API
server's port (usually 6443). This will work because then the NLB will use its dynamically allocated internal IP
address as source.

Second, if client IP preservation is _enabled_ we MUST set `0.0.0.0/0` as communication source because then the
incoming IP address will be that of the client's that might not be in the current VPC. This shouldn't be too much of a
problem, but user's need to be aware of this restriction.

## Extension of the code

Right now, only NLBs and a Classic Load Balancer is supported. However, the code has been written in a way that it
should be easy to extend with an ALB or a GLB.
