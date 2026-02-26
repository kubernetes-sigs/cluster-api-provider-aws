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

[capi-images]: https://image-builder.sigs.k8s.io/capi/capi.html
[image-builder]: https://github.com/kubernetes-sigs/image-builder
[image-builder-aws]: https://github.com/kubernetes-sigs/image-builder/tree/master/images/capi/packer/ami
[aws-capi-images]: https://image-builder.sigs.k8s.io/capi/providers/aws.html
[upgrading-workload-clusters]: https://cluster-api.sigs.k8s.io/tasks/kubeadm-control-plane.html#upgrading-workload-clusters

