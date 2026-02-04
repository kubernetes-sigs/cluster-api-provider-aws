# Enabling IPv6

## Glossary

- **CAPA** - Cluster API Provider AWS
- **NLB** - Network Load Balancer
- **CCM** - Cloud Controller Manager
- **Dualstack** - In this document, dualstack (IPv4 + IPv6) refers to network infrastructure (VPC, subnets, and load balancers) that supports both IPv4 and IPv6 address families. This is distinct from Kubernetes dualstack, where pods and services are assigned addresses from both IP families.

## Overview

CAPA enables you to create IPv6 only and dualstack (IPv4 + IPv6) Kubernetes clusters on Amazon Web Services (AWS) on a dualstack network infrastructure.

**Important**: CAPA does not support in-place migration on self-managed and managed clusters from IPv4 to dualstack or IPv6. You must create a new cluster.

## Prerequisites

The instance types for control plane and worker machines must support IPv6. To see a list of instance types that support IPv6 in your region, run the following command:

```bash
aws ec2 describe-instance-types \
  --region <region> \
  --filters "Name=network-info.ipv6-supported,Values=true" \
  --query 'InstanceTypes[].InstanceType'
```

If you want to check whether a specific instance type supports IPv6, run the following command (output `true`):

```bash
aws ec2 describe-instance-types \
  --region <region> \
  --instance-types <instance-type> \
  --query 'InstanceTypes[0].NetworkInfo.Ipv6Supported'
```

If you are launching machines in IPv6-only subnets, the instance type must be Nitro-based. To check whether the instance type is Nitro-based, run the following command (output `"nitro"`):

```bash
aws ec2 describe-instance-types \
  --region <region> \
  --instance-types <instance-type> \
  --query 'InstanceTypes[0].Hypervisor'
```

## Enabling IPv6 capabilities

To instruct CAPA to configure IPv6 capabilities for the cluster's network infrastructure, you must explicitly define `spec.network.vpc.ipv6` in either `AWSCluster` or `AWSManagedControlPlane`. See [IPv6 CIDR Allocations](#ipv6-cidr-allocations) for different IPv6 CIDR configuration options.

```yaml
spec:
  network:
    vpc:
      ipv6: {}
```

**Note:** CAPA, by default, will provision a dualstack network infrastructure (i.e. dualstack VPC and subnets). However, your Kubernetes cluster can be configured as either Kubernetes IPv6-only or dualstack depending on your pod/service CIDR configuration.

## Supported Network Configurations

CAPA supports various network configuration combinations for creating clusters with different IPv4/IPv6 requirements. The following table shows all possible combinations of subnet and load balancer configurations (assuming internet-facing load balancer):

| Public Subnet | Private Subnet | Load Balancer IP Type | Target Group IP Type | Status | Notes |
|---------------|----------------|----------------------|---------------------|--------|-------|
| IPv4 | IPv4 | ipv4 | ipv4 | ✅ Supported | Traditional IPv4 cluster |
| dualstack | dualstack | dualstack | ipv6 | ✅ Supported | Dualstack LB with IPv6 as primary |
| dualstack | dualstack | dualstack | ipv4 | ✅ Supported | Dualstack LB with IPv4 as primary |
| dualstack | dualstack | ipv4 | ipv4 | ✅ Supported | IPv4-only LB on dualstack infrastructure |
| dualstack | IPv4 | dualstack | ipv4 | ✅ Supported | Dualstack LB with IPv4-only control plane |
| dualstack | IPv6-only | dualstack | ipv6 | ✅ Supported | Dualstack LB with IPv6-only control plane |
| dualstack | IPv4 | dualstack | ipv6 | ❌ Invalid | Cannot use IPv6 targets when control plane has no IPv6 addresses |
| IPv6-only | * | * | * | ❌ Invalid | Load Balancers cannot use IPv6-only subnets |
| * | * | ipv6 | * | ❌ Invalid | `ipv6` is not a valid load balancer IP type |

When Kubernetes pods have only IPv6 addresses and need to communicate with IPv4-only internet services, [NAT64/DNS64](https://docs.aws.amazon.com/vpc/latest/userguide/nat-gateway-nat64-dns64.html) must be enabled for the subnets where nodes reside. CAPA automatically enables NAT64/DNS64 for IPv6-only subnets that it creates (see [Mixing subnets of different IP families](#mixing-subnets-of-different-ip-families) on how to tell CAPA to create IPv6-only subnets). For pre-existing subnets, you must enable NAT64/DNS64 manually.

## IPv6 CIDR Allocations

CAPA supports various methods to allocate an IPv6 CIDR to the cluster VPC.

### AWS-assigned IPv6 VPC CIDR

To request AWS to automatically assign an IPv6 CIDR from an AWS defined address pool, use the following setting:

```yaml
spec:
  network:
    vpc:
      ipv6: {}
```

By default, Amazon provides one fixed size (`/56`) IPv6 CIDR block to a VPC.

### Bring-your-own IPv6 VPC CIDR (EC2)

If you own an IPv6 address space, you can import it into AWS EC2 IPv6 address pool (See [guide](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-byoip.html#byoip-requirements)). After importing it, you can assign /56 ranges from the space to individual VPCs in the same account.

To define your own IPv6 address pool and CIDR set the following values:

```yaml
spec:
  network:
    vpc:
      ipv6:
        poolId: pool-id
        cidrBlock: "2009:1234:ff00::/56"
```

### Bring-your-own IPv6 VPC CIDR via VPC Address Manager (VPC IPAM)

If you want to allocate an IPv6 CIDR to the VPC from an existing VPC IPAM pool, define the pool ID and a prefix length as follows:

```yaml
spec:
  network:
    vpc:
      ipv6:
        ipamPool:
          id: ipam-pool-id
          netmaskLength: 56
```

By default, if you omit `netmaskLength`, CAPA will set it to the default `56`.

### Bring-your-own IPv6 VPC

If you have an existing dualstack VPC that you would like to use, you must explicitly provide the IPv6 CIDR block and egress-only internet gateway ID specs:

```yaml
spec:
  network:
    vpc:
      id: vpc-1234567890abcdefg
      cidrBlock: 10.0.0.0/16
      ipv6:
        cidrBlock: "2001:1234:ff00::/56"
        egressOnlyInternetGatewayId: eigw-1234567890abcdefg
```

This has to be done to explicitly express the user intention to use the IPv6 capabilities of the VPC.

## Creating IPv6 EKS-managed Clusters

An example configuration for deploying an IPv6 EKS cluster can be found here: [IPv6 EKS cluster template](https://raw.githubusercontent.com/kubernetes-sigs/cluster-api-provider-aws/refs/heads/main/templates/cluster-template-eks-ipv6.yaml).

<aside class="note warning">

<h1>Warning</h1>

EKS currently only supports IPv6-only clusters (not dualstack). You can't define custom Pod CIDRs on EKS with IPv6. EKS automatically assigns an address range from a unique local address range of `fc00::/7`.

</aside>

**Notes**: All addons **must** be enabled. A working IPv6 cluster configuration defines `spec.network.vpc.ipv6` and all addons as follows:

```yaml
kind: AWSManagedControlPlane
apiVersion: controlplane.cluster.x-k8s.io/v1beta1
metadata:
  name: "${CLUSTER_NAME}-control-plane"
spec:
  network:
    vpc:
      ipv6: {}
  region: "${AWS_REGION}"
  sshKeyName: "${AWS_SSH_KEY_NAME}"
  version: "${KUBERNETES_VERSION}"
  addons:
    - name: "vpc-cni"
      version: "v1.11.0-eksbuild.1"  # Note: Check for latest compatible version
      # this is important, otherwise environment property update will not work
      conflictResolution: "overwrite"
    - name: "coredns"
      version: "v1.8.7-eksbuild.1"  # Note: Check for latest compatible version
    - name: "kube-proxy"
      version: "v1.22.6-eksbuild.1"  # Note: Check for latest compatible version
```

## Creating IPv6 Self-managed Clusters

An example configuration for deploying an IPv6 self-managed cluster can be found here: [IPv6 cluster template](https://raw.githubusercontent.com/kubernetes-sigs/cluster-api-provider-aws/refs/heads/main/templates/cluster-template-ipv6.yaml).

When creating a self-managed cluster, you can define the IPv6 Pod and Service CIDR. For example, you can define ULA IPv6 range `fd01::/48` for pod networking and `fd02::/112` for service networking.

```yaml
apiVersion: cluster.x-k8s.io/v1beta1
kind: Cluster
metadata:
  name: "${CLUSTER_NAME}"
spec:
  clusterNetwork:
    pods:
      cidrBlocks:
      - fd01::/48
    services:
      cidrBlocks:
      - fd02::/112
```

<aside class="note warning">

<h1>Warning</h1>

**CoreDNS Resolution Issue**: In Kubernetes IPv6-only clusters, CoreDNS pods are assigned only IPv6 addresses. When resolving external domain names (e.g. `quay.io`), they forward queries to upstream nameservers listed in `/etc/resolv.conf`. However, these nameservers might be defined with only IPv4 addresses, which CoreDNS pods cannot reach since they only have IPv6 connectivity. This causes DNS resolution failures for non-cluster domains.

**Solution**: Configure CoreDNS to use AWS Route53's IPv6 nameserver at `fd00:ec2::253`. Edit the `coredns` ConfigMap in namespace `kube-system` and change the forwarding configuration from `forward . /etc/resolv.conf` to `forward . fd00:ec2::253 /etc/resolv.conf`:

```bash
kubectl -n kube-system edit cm/coredns
```

This tells CoreDNS to try the IPv6 nameserver first, then fall back to addresses in `/etc/resolv.conf`.

**Note**: This workaround is only required for Kubernetes IPv6-only clusters. In Kubernetes dualstack clusters where pods have both IPv4 and IPv6 addresses, CoreDNS can use its IPv4 address to reach the IPv4 nameservers, so no modification is needed.

</aside>

## Creating Dualstack Self-managed Clusters

An example configuration for deploying a dualstack self-managed cluster can be found here:
- [Dualstack IPv6 primary cluster template](https://raw.githubusercontent.com/kubernetes-sigs/cluster-api-provider-aws/refs/heads/main/templates/cluster-template-dualstack-ipv6-primary.yaml)
- [Dualstack IPv4 primary cluster template](https://raw.githubusercontent.com/kubernetes-sigs/cluster-api-provider-aws/refs/heads/main/templates/cluster-template-dualstack-ipv4-primary.yaml)

When creating a self-managed cluster, you can define both IPv4 and IPv6 Pod and Service CIDRs. The first entry defines the primary IP family. For example, the below configuration sets up a dualstack IPv6 primary cluster:

```yaml
apiVersion: cluster.x-k8s.io/v1beta1
kind: Cluster
metadata:
  name: "${CLUSTER_NAME}"
spec:
  clusterNetwork:
    pods:
      cidrBlocks:
      - fd01::/48
      - 192.168.0.0/16
    services:
      cidrBlocks:
      - fd02::/112
      - 172.30.0.0/16
```

## Cloud Controller Manager (CCM) Node IP Configurations

You need to provide cloud-config to the CCM via a ConfigMap to set the `NodeIPFamilies` to include IPv6. This instructs the CCM to consider IPv6 in the machine network interface, if any. If not configured, the CCM will only consider the node's IPv4 address. This causes nodes to have only IPv4 and new pods with `hostNetwork: true` will only pick up the node's IPv4 address.

For example, provide the following ConfigMap to `cloud-controller-manager-addon`:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: cloud-config
  namespace: kube-system
data:
  cloud-config.conf: |
    [Global]
    NodeIPFamilies=ipv6
    NodeIPFamilies=ipv4
```

And then provide the `cloud-config.conf` to the CCM DaemonSet as follows:

```yaml
spec:
  containers:
  - name: aws-cloud-controller-manager
    image: registry.k8s.io/provider-aws/cloud-controller-manager:v1.28.3
    args:
      - --v=2
      - --cloud-provider=aws
      - --use-service-account-credentials=true
      - --configure-cloud-routes=false
      - --cloud-config=/etc/kubernetes/cloud-config.conf # Define cloud-config file path
    volumeMounts:
    - name: cloud-config
      mountPath: /etc/kubernetes/cloud-config.conf
      subPath: cloud-config.conf
  hostNetwork: true
  volumes:
  - name: cloud-config
    configMap:
      name: cloud-config
```

## Cloud Controller Manager Load Balancer Limitations

<aside class="note warning">
<h1>Warning</h1>

The AWS Cloud Controller Manager (CCM) does **not** currently support dualstack Load Balancers. When creating Services of type LoadBalancer in a dualstack cluster, the Load Balancers will be created with **only** IPv4.

Please see [kubernetes/cloud-provider-aws/issues/1219](https://github.com/kubernetes/cloud-provider-aws/issues/1219) to track CCM dualstack support progress.
</aside>

## CNI IPv6 Support for Self-managed Clusters

By default, no CNI plugin is installed when provisioning a self-managed cluster. You need to install your own CNI solution that supports IPv6, for example, Calico with VXLAN. You can find the guides to enable [IPv6](https://docs.tigera.io/calico/latest/networking/ipam/ipv6) and [VXLAN](https://docs.tigera.io/calico/latest/networking/configuring/vxlan-ipip) support for Calico on their official documentation.

**Important notes for Calico with IPv6**:
- Calico supports IPv6 with VXLAN encapsulation only (IP-in-IP is not supported for IPv6)
- VXLAN for IPv6 requires kernel version ≥ 4.19.1 (or Red Hat kernel ≥ 4.18.0)
- If you are using Calico as the CNI provider, ensure the CNI ingress rule allows VXLAN for cross-subnet communications. You can set the rule in the `AWSCluster` resource, for example:
```yaml
spec:
  network:
    cni:
      cniIngressRules:
      # If using Calico as CNI provider, this rule is required.
      # Note: Calico currently supports IPv6 with VXLAN.
      - description: "VXLAN (calico)"
        protocol: udp
        fromPort: 4789
        toPort: 4789
```

## Mixing subnets of different IP families

CAPA allows you to define the AZs the subnets should be created in, the number of subnets per AZ and whether a subnet is IPv4, dualstack, or IPv6-only. For example:

```yaml
spec:
  network:
    subnets:
      # This creates a dualstack public subnet in us-east-1a
      # Both cidrBlock + isIpv6==true
    - cidrBlock: 10.0.0.0/20
      isIpv6: true
      isPublic: true
      availabilityZone: us-east-1a
      id: ${CLUSTER_NAME}-subnet-public-us-east-1a
    # This creates a dualstack public subnet in us-east-1b
    # Both cidrBlock + isIpv6==true
    - cidrBlock: 10.0.16.0/20
      isIpv6: true
      isPublic: true
      availabilityZone: us-east-1b
      id: ${CLUSTER_NAME}-subnet-public-us-east-1b
    # This creates an IPv4 private subnet in us-east-1a
    # Only cidrBlock defined + isIpv6==false (default)
    - cidrBlock: 10.0.128.0/20
      isPublic: false
      availabilityZone: us-east-1a
      id: ${CLUSTER_NAME}-subnet-private-us-east-1a
    # This creates an IPv6-only private subnet in us-east-1a
    # cidrBlock is undefined + isIpv6==true
    - isPublic: false
      isIpv6: true
      availabilityZone: us-east-1a
      id: ${CLUSTER_NAME}-subnet-private-1-us-east-1a
    # This creates an IPv4 private subnet in us-east-1b
    # Only cidrBlock defined + isIpv6==false (default)
    - cidrBlock: 10.0.144.0/20
      isPublic: false
      availabilityZone: us-east-1b
      id: ${CLUSTER_NAME}-subnet-private-us-east-1b
    # This creates an IPv6-only private subnet in us-east-1b
    # cidrBlock is undefined + isIpv6==true
    - isPublic: false
      isIpv6: true
      availabilityZone: us-east-1b
      id: ${CLUSTER_NAME}-subnet-private-1-us-east-1b
    vpc:
      cidrBlock: 10.0.0.0/16
      # The VPC IPv6 CIDR will be allocated by AWS.
      ipv6: {}
  region: us-east-1
```

A subnet IP specification is defined as follows (applied to CAPA-managed VPC only):

| Subnet Type | `isIpv6` | `cidrBlock` | `ipv6CidrBlock` | Notes |
|-------------|----------|-------------|-----------------|-------|
| **IPv4-only** | `false` or omitted | Required | N/A | Traditional IPv4 subnet |
| **Dualstack** | `true` | Required | Optional | Auto-assigned from VPC CIDR if omitted |
| **IPv6-only** | `true` | Omitted/empty | Optional | Auto-assigned from VPC CIDR if omitted |

## IPv6 support for Local and Wavelength zones

According to the AWS docs, the state of IPv6 support is as follows:

- ❌ No IPv6 support for Wavelength zones. See [reference](https://docs.aws.amazon.com/wavelength/latest/developerguide/wavelength-quotas.html#vpc-considerations).
- ⚠️ Limited support for Local zones, which requires a dedicated IPv6 CIDR for local zone network border group. See [reference](https://docs.aws.amazon.com/local-zones/latest/ug/how-local-zones-work.html#considerations).

Thus, CAPA currently does not support creating IPv6-enabled subnets in Local and Wavelength zones.

However, if you have an existing VPC with IPv6-only or dualstack subnets in Local zones, you can define them in the cluster spec.


```yaml
spec:
  network:
    subnets:
    - id: "cluster-subnet-private-us-east-1a"
    - id: "cluster-subnet-public-us-east-1a"
    - id: "cluster-subnet-private-us-east-1b"
    - id: "cluster-subnet-public-us-east-1b"
    - id: "cluster-subnet-private-us-east-1-nyc-1a"
    - id: "cluster-subnet-public-us-east-1-nyc-1a"
    - id: "cluster-subnet-private-us-east-1-wl1-was-wlz-1"
    - id: "cluster-subnet-public-us-east-1-wl1-was-wlz-1"
    vpc:
      id: vpc-1234567890abcdefg
      cidrBlock: 10.0.0.0/16
      ipv6:
        cidrBlock: "2001:1234:ff00::/56"
        egressOnlyInternetGatewayId: eigw-1234567890abcdefg
```
