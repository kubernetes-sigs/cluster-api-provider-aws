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

## Target Group IP Address Type

**Note:** The `targetGroupIPType` field is only available when using Network Load Balancers (NLB), Application Load Balancers (ALB), or Gateway Load Balancers (ELB). It **cannot** be configured when using Classic Load Balancers.

By default, the target group IP address type is automatically determined from the control plane subnet configuration:
- If control plane subnets are IPv4-only, the target group uses `ipv4`
- If control plane subnets are IPv6-only or dualstack, the target group uses `ipv6`

You can explicitly override the IP address type for the target group using the `targetGroupIPType` field:

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
    targetGroupIPType: ipv6
```

Valid values are:
- `ipv4`: Routes traffic to targets using IPv4 addresses
- `ipv6`: Routes traffic to targets using IPv6 addresses

### Additional Listeners

The `targetGroupIPType` can also be configured independently for each additional listener:

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
    targetGroupIPType: ipv4
    additionalListeners:
      - port: 8443
        protocol: TCP
        targetGroupIPType: ipv6
```

This allows you to have different IP address types for different target groups within the same load balancer.

**Note:** The `targetGroupIPType` field is only applicable when using Network Load Balancers (NLB), Application Load Balancers (ALB), or Gateway Load Balancers (ELB). It **cannot** be set when using Classic Load Balancers.

## Extension of the code

Right now, only NLBs and a Classic Load Balancer is supported. However, the code has been written in a way that it
should be easy to extend with an ALB or a GLB.
