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
        ipv6Pool: pool-id
        ipv6CidrBlock: "2009:1234:ff00::/56"
```

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

Now comes the tricky part. If you wish, it's possible to set up IPv6 with unmanaged clusters, however, that requires
a lot of extra manual steps and some extra configuration settings.

_Note_: I DO NOT recommend doing this in production. AWS provides a more robust and easier approach on doing IPv6.
This approach is brittle and has many manual steps that need to be performed in order to get things working.

### Extra Config

The extra configs are on the kubeadm side. These are `node-ip` and `bind-address`. These need to be set as follows:

```yaml
  kubeadmConfigSpec:
    initConfiguration:
      nodeRegistration:
        name: '{{ ds.meta_data.local_hostname }}'
        kubeletExtraArgs:
          cloud-provider: aws
          node-ip: '::'
    clusterConfiguration:
      apiServer:
        extraArgs:
          cloud-provider: aws
          bind-address: '::'
      controllerManager:
        extraArgs:
          cloud-provider: aws
          bind-address: '::'
      scheduler:
        extraArgs:
          bind-address: '::'
    joinConfiguration:
      nodeRegistration:
        name: '{{ ds.meta_data.local_hostname }}'
        kubeletExtraArgs:
          node-ip: '::'
          cloud-provider: aws

```

This will tell kubeadm to bind to a specific address type which should be IPv6.

Next, it's pod CIDRs and service CIDRs. This is a bit more tricky. You need to know your IPv6 CIDR beforehand.
Having your own IPv6 pool is most of the times, impractical. But there is a way to get you started up quickly
and with low effort. You can ask CAPA to create the network topology for you with a simple cluster config such
as this one:

```yaml
---
apiVersion: cluster.x-k8s.io/v1beta1
kind: Cluster
metadata:
  name: "ipv6-network-topology"
spec:
  infrastructureRef:
    apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
    kind: AWSCluster
    name: "ipv6-network-topology"
  controlPlaneRef:
    kind: KubeadmControlPlane
    apiVersion: controlplane.cluster.x-k8s.io/v1beta1
    name: "ipv6-network-topology-control-plane"
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: AWSCluster
metadata:
  name: "ipv6-network-topology"
spec:
  network:
    vpc:
      ipv6: {}
  region: "eu-central-1"
```

This will create a VPC with proper load-balancing and elastic ips and everything. This can be fine-tuned as
desired. This is the bare minimum where CAPA creates the whole topology.

Once this is done, we can acquire the IPv6 CIDR and we can continue by using the vpc id and subnets in the
unmanaged setting like this:

```yaml
---
apiVersion: cluster.x-k8s.io/v1beta1
kind: Cluster
metadata:
  name: "unmanaged-ipv6"
spec:
  clusterNetwork:
    services:
      cidrBlocks: ["192.168.0.0/16", "2a05:d014:852:f::/56"]
    pods:
      cidrBlocks: ["192.168.0.0/16", "2a05:d014:852:f::/56"]
  infrastructureRef:
    apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
    kind: AWSCluster
    name: "unmanaged-ipv6"
  controlPlaneRef:
    kind: KubeadmControlPlane
    apiVersion: controlplane.cluster.x-k8s.io/v1beta1
    name: "test-ipv6-unmanaged-2-control-plane"
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: AWSCluster
metadata:
  name: "unmanaged-ipv6"
spec:
  network:
    subnets:
    - id: "subnet-0812d970f75867a72"
    - id: "subnet-0aa534118f62cab86"
    - id: "subnet-0380be0501f2c8bb0"
    - id: "subnet-09d083492229f6281"
    - id: "subnet-0b12a78a3b2cebdec"
    - id: "subnet-0a0b33746595ecc89"
    vpc:
      id: "vpc-024ae81c0ca3b7209"
      ipv6: {}
  region: "eu-central-1"
---
```

### Cilium

Since we are on an unmanaged cluster, we need a cni installed. We'll use Cilium as an example. There are two settings
that are needed for cilium to work. (3 really, you also have to set up ipv4 cidr to the proper cidr the VPC provides).

These are:

```
--cluster-pool-ipv6-cidr="2a05:d014:852:f02::/112"
--cluster-pool-ipv6-cidr-size="128"
```

Note that Cilium doesn't allow any CIDR above `112`. So the CIDR you've got from AWS with size `64` needs to be cut down
to `112`.

Once this is done, we can install Cilium into the workload cluster and restart all Pods so they can acquire IPv6
addresses.

### Putting it all together

All new pods will acquire an IPv6 address, however, existing pods, such as kube-proxy, will remain with IPv4 addresses.
