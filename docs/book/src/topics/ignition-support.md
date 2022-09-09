# Ignition support

- **Feature status:** Experimental
- **Feature gate:** BootstrapFormatIgnition=true

The default configuration engine for bootstrapping workload cluster machines is [cloud-init][cloud-init].
**Ignition** is an alternative engine used by Linux distributions such as [Flatcar Container Linux][flatcar]
and [Fedora CoreOS][fedora-coreos] and therefore should be used when choosing an Ignition-based distribution as
the underlying OS for workload clusters.

<aside class="note warning">

<h1>Note</h1>

This initial implementation uses Ignition **v2** and was tested with **Flatcar Container Linux** only.
Future releases are expected to add Ignition **v3** support and cover more Linux distributions.

</aside>

This document explains how Ignition support works.

For more generic information, see [Cluster API documentation on Ignition Bootstrap configuration][cabpk].

## Overview

By default machine controller stores EC2 instance user data using SSM to store it encrypted, which underneath
use multi part mime types, which are [unlikely to be supported](https://github.com/coreos/ignition/issues/1072)
by Ignition.

EC2 user data is also limited to 64 KB, which is often not enough to provision Kubernetes controlplane because
of the size of required certificates and configuration files.

To address those limitations CAPA can create and use S3 Bucket to store encrypted user data, which will be then
pulled by the instances during provisioning.

## IAM Permissions

To manage S3 Buckets and objects inside them, CAPA controllers require additional IAM permissions.

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

## Enabling EXP_BOOTSTRAP_FORMAT_IGNITION feature gate

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

## Bucket and object management

When you want to use Ignition user data format for you machines, you need to configure your cluster to
specify which S3 bucket to use. Controller will then make sure that the bucket exists and that required policies
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

After successful machine provisioning, bootstrap data is removed from the bucket.

During cluster removal, if S3 bucket is empty, it will be removed as well.

## Bucket naming

Bucket naming must follow [S3 Bucket naming rules][bucket-naming-rules].

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
