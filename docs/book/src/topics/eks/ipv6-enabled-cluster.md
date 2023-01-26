# IPv6 Enabled Cluster

CAPA supports IPv6 enabled clusters. Dual stack clusters are not yet supported, but
dual VPC, meaning both ipv6 and ipv4 are defined, is supported and in fact, it's the
only mode of operation at the writing of this doc.

Upcoming feature will be IPv6 _only_.

## Managed Clusters

### How to set up

Two modes of operations are supported. Request AWS to generate and assign an address
or BYOIP which is Bring Your Own IP. There must already be a provisioned pool and a
set of IPv6 CIDRs for that.

#### Automatically Generated IP

To request AWS to assign a set of IPv6 addresses from an AWS defined address pool,
use the following setting:

```yaml
kind: AWSManagedControlPlane
apiVersion: controlplane.cluster.x-k8s.io/v1beta1
metadata:
  name: "${CLUSTER_NAME}-control-plane"
spec:
  network:
    vpc:
      ipv6: {}
```

#### BYOIP ( Bring Your Own IP )

To define your own IPv6 address pool and CIDR set the following values:

```yaml
spec:
  network:
    vpc:
      ipv6:
        poolId: pool-id
        cidrBlock: "2009:1234:ff00::/56"
```

If you have a VPC that is IPv6 enabled and you would like to use it, please define it in the config:

```yaml
spec:
  network:
    vpc:
      ipv6: {}
```

This has to be done explicitly because otherwise, it would break in the following two scenarios:
- During an upgrade from 1.5 to >=2.0 where the VPC is ipv6 enabled, but CAPA was only recently made aware
- During a migration on the VPC, switching it from only IPv4 to Dual Stack ( it would see that ipv6 is enabled and
  enforce it while doing that would not have been the intention of the user )


### Requirements

The use of a Nitro enabled instance is required. To see a list of nitro instances in your region
run the following command:

```bash
aws ec2 describe-instance-types --filters Name=hypervisor,Values=nitro --region us-west-2  | grep "InstanceType"
```

This will list all available Nitro hypervisor based instances in your region.

All addons **must** be enabled. A working cluster configuration looks like this:

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
      version: "v1.11.0-eksbuild.1"
      conflictResolution: "overwrite" # this is important, otherwise environment property update will not work
    - name: "coredns"
      version: "v1.8.7-eksbuild.1"
    - name: "kube-proxy"
      version: "v1.22.6-eksbuild.1"
```

You can't define custom POD CIDRs on EKS with IPv6. EKS automatically assigns an address range from a unique local
address range of `fc00::/7`.

## Unmanaged Clusters

Unmanaged clusters are not supported at this time.
