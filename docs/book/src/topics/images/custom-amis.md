# Custom Kubernetes AMIs

Cluster API uses the Kubernetes [Image Builder][image-builder] tools. You should use the [AWS images][image-builder-aws] from that project as a starting point for your custom image.

[The Image Builder Book][capi-images] explains how to build the images defined in that repository, with instructions for [AWS CAPI Images][aws-capi-images] in particular.

## Operating system requirements

For custom images to work with Cluster API, it must meet the operating system requirements of the bootstrap provider. For example, the default `kubeadm` bootstrap provider has a set of [`preflight checks`][kubeadm-preflight-checks] that a VM is expected to pass before it can join the cluster.

## Kubernetes version requirements

The pre-built public images are each built to support a specific version of Kubernetes. When using custom images, make sure to match the image to the `version:` field of the `KubeadmControlPlane` and `MachineDeployment` in the YAML template for your workload cluster.

To upgrade to a new Kubernetes release with custom images requires this preparation:

- create a new custom image which supports the Kubernetes release version
- copy the existing `AWSMachineTemplate` and change its `ami:` section to reference the new custom image
- create the new `AWSMachineTemplate` on the management cluster
- modify the existing `KubeadmControlPlane` and `MachineDeployment` to reference the new `AWSMachineTemplate` and update the `version:` field to match

See [Upgrading workload clusters][upgrading-workload-clusters] for more details.

## Creating a cluster from a custom image

To use a custom image, it needs to be referenced in an `ami:` section of your `AWSMachineTemplate`.

```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: AWSMachineTemplate
metadata:
  name: capa-image-id-example
  namespace: default
spec:
  template:
    spec:
      ami:
        id: ami-09709369c53539c11
      iamInstanceProfile: control-plane.cluster-api-provider-aws.sigs.k8s.io
      instanceType: m5.xlarge
      sshKeyName: default
```

## Selecting a custom image by filters

Instead of pinning an explicit AMI `id`, you can let the provider resolve the AMI at provisioning time from a set of [EC2 image filters][ec2-describe-images]. When more than one image matches, the most recently created one (by `CreationDate`) is selected.

```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
kind: AWSMachineTemplate
metadata:
  name: capa-image-filters-example
  namespace: default
spec:
  template:
    spec:
      ami:
        filters:
          - name: owner-id
            values:
              - "123456789012"
          - name: name
            values:
              - "my-custom-ami-*"
      iamInstanceProfile: control-plane.cluster-api-provider-aws.sigs.k8s.io
      instanceType: m5.xlarge
      sshKeyName: default
```

<aside class="note warning">

<h1>Warning</h1>

**Restrict ownership.** If the filters do not constrain ownership, the result set includes **all public AMIs** in the region; thus, a third party could publish an image matching your filters and have it selected instead of yours because the newest match wins. Always include an `owner-id` (or `owner-alias`) filter that scopes the lookup to accounts you trust.

</aside>

Additional caveats:

- **The AMI is not pinned.** Because the newest matching image is selected, the resolved AMI can change over time as new images are published. Use `ami.id` when you need a stable, reproducible image.
- **Mutually exclusive.** `ami.filters` cannot be combined with `ami.id`, `ami.eksLookupType`, or the `imageLookupFormat`, `imageLookupOrg`, and `imageLookupBaseOS` fields.

[capi-images]: https://image-builder.sigs.k8s.io/capi/capi.html
[image-builder]: https://github.com/kubernetes-sigs/image-builder
[image-builder-aws]: https://github.com/kubernetes-sigs/image-builder/tree/master/images/capi/packer/ami
[aws-capi-images]: https://image-builder.sigs.k8s.io/capi/providers/aws.html
[upgrading-workload-clusters]: https://cluster-api.sigs.k8s.io/tasks/kubeadm-control-plane.html#upgrading-workload-clusters
[ec2-describe-images]: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeImages.html
