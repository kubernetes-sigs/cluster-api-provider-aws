# Ignition support

- **Feature status:** Experimental
- **Feature gate:** BootstrapFormatIgnition=true

The default configuration engine for bootstrapping workload cluster machines is [cloud-init][cloud-init].
**Ignition** is an alternative engine used by Linux distributions such as [Flatcar Container Linux][flatcar]
and [Fedora CoreOS][fedora-coreos] and therefore should be used when choosing an Ignition-based distribution as
the underlying OS for workload clusters.

<aside class="note warning">

<h1>Note</h1>

This initial implementation used Ignition **v2** and was tested with **Flatcar Container Linux** only.
Further releases added Ignition **v3** support.

</aside>

This document explains how Ignition support works.

For more generic information, see [Cluster API documentation on Ignition Bootstrap configuration][cabpk].

## Overview

When using CloudInit for bootstrapping, by default the awsmachine controller stores EC2 instance user data using SSM to store it encrypted, which underneath uses multi part mime types.
Unfortunately multi part mime types are [not supported](https://github.com/coreos/ignition/issues/1072) by Ignition. Moreover EC2 instance user data storage is also limited to 64 KB, which might not always be enough to provision Kubernetes controlplane because of the size of required certificates and configuration files.

To address these limitations, when using Ignition for bootstrapping, by default the awsmachine controller uses a Cluster Object Store (e.g. S3 Bucket), configured in the AWSCluster, to store user data,
which will be then pulled by the instances during provisioning.

Optionally, when using Ignition for bootstrapping, users can optionally choose an alternative storageType for user data.
For now the single available alternative is to store user data unencrypted directly in the EC2 instance user data.
This storageType option is although discouraged unless strictly necessary, as it is not considered as safe as storing it in the S3 Bucket.

## Prerequirements for enabling Ignition bootstrapping

### Enabling EXP_BOOTSTRAP_FORMAT_IGNITION feature gate

In order to activate Ignition bootstrap you first need to enable its feature gate.

When deploying CAPA using `clusterctl`, make sure you set `BOOTSTRAP_FORMAT_IGNITION=true` and
`EXP_KUBEADM_BOOTSTRAP_FORMAT_IGNITION=true `environment variables to enable experimental Ignition bootstrap
support.

``` sh
# Enable the feature gates controlling Ignition bootstrap.
export EXP_KUBEADM_BOOTSTRAP_FORMAT_IGNITION=true # Used by the kubeadm bootstrap provider.
export EXP_BOOTSTRAP_FORMAT_IGNITION=true # Used by the AWS provider.

# Initialize the management cluster.
clusterctl init --infrastructure aws
```

## Choosing a storage type for Ignition user data

S3 is the default storage type when Ignition is enabled for managing machine's bootstrapping.
But other methods can be choosen for storing Ignition user data.

### Store Ignition config in a Cluster Object Store (e.g. S3 bucket)

To explicitly set ClusterObjectStore as the storage type, provide the following config in the `AWSMachineTemplate`:
```yaml
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: AWSMachineTemplate
metadata:
  name: "test"
spec:
  template:
    spec:
      ignition:
        storageType: ClusterObjectStore
```

#### Cluster Object Store and object management

When you want to use Ignition user data format for you machines, you need to configure your cluster to
specify which Cluster Object Store to use. Controller will then check that the bucket already exists and that required policies
are in place.

See the configuration snippet below to learn how to configure `AWSCluster` to manage S3 bucket.

``` yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: AWSCluster
spec:
  s3Bucket:
    controlPlaneIAMInstanceProfile: control-plane.cluster-api-provider-aws.sigs.k8s.io
    name: cluster-api-provider-aws-unique-suffix
    nodesIAMInstanceProfiles:
    - nodes.cluster-api-provider-aws.sigs.k8s.io
```

Buckets are safe to be reused between clusters.

After successful machine provisioning, the bootstrap data is removed from the object store.

During cluster removal, if the Cluster Object Store is empty, it will be deleted as well.

#### S3 IAM Permissions

If you choose to use an S3 bucket as the Cluster Object Store, CAPA controllers require additional IAM permissions.

If you use `clusterawsadm` for managing the IAM roles, you can use the configuration below to create S3-related
IAM permissions.

``` yaml
apiVersion: bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1
kind: AWSIAMConfiguration
spec:
  s3Buckets:
    enable: true
```

See [Using clusterawsadm to fulfill prerequisites](./using-clusterawsadm-to-fulfill-prerequisites.md) for more
details.

#### Cluster Object Store naming

Cluster Object Store and bucket naming must follow [S3 Bucket naming rules][bucket-naming-rules].

In addition, by default `clusterawsadm` creates IAM roles to only allow interacting with buckets with
`cluster-api-provider-aws-` prefix to reduce the permissions of CAPA controller, so all bucket names should
use this prefix.

To change it, use `spec.s3Buckets.namePrefix` field in `AWSIAMConfiguration`.

```yaml
apiVersion: bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1
kind: AWSIAMConfiguration
spec:
  s3Buckets:
    namePrefix: my-custom-secure-bucket-prefix-
```

### Store Ignition config as UnencryptedUserData

<aside class="note warning">
<h1>WARNING</h1>
**This is discouraged as is not considered as secure as other storage types.**
</aside>

To instruct the controllers to store the user data directly in the EC2 instance user data unencrypted,
 provide the following config in the `AWSMachineTemplate`:
```yaml
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: AWSMachineTemplate
metadata:
  name: "test"
spec:
  template:
    spec:
      ignition:
        storageType: UnencryptedUserData
```

No further requirements are necessary.

## Supported bootstrap providers

At the moment only [CABPK][cabpk] is known to support producing bootstrap data in Ignition format.

## Trying it out

If you want to test Ignition support, use `flatcar` cluster flavor.

## Other bootstrap providers

If you want to use Ignition support with custom bootstrap provider which supports producing Ignition bootstrap
data, ensure that bootstrap provider sets the `format` field in machine bootstrap secret to `ignition`. This
information is used by the machine controller to determine which user data format to use for the instances.

[bucket-naming-rules]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/bucketnamingrules.html
[cloud-init]: https://cloudinit.readthedocs.io/
[flatcar]: https://www.flatcar.org/docs/latest/provisioning/ignition/
[fedora-coreos]: https://docs.fedoraproject.org/en-US/fedora-coreos/producing-ign/
[cabpk]: https://cluster-api.sigs.k8s.io/tasks/experimental-features/ignition.html
