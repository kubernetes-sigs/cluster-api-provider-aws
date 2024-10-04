# AWS Machine Images for CAPA Clusters

CAPA requires a “machine image” containing pre-installed, matching versions of kubeadm and kubelet.

## EKS Clusters

For an EKS cluster the default behaviour is to retieve the AMI to use from SSM. This is so the recommended Amazon Linux AMI is used (see [here](https://docs.aws.amazon.com/eks/latest/userguide/retrieve-ami-id.html)).

Instead of using the auto resolved AMIs an appropriate custom image ID for the Kubernetes version can be set in `AWSMachineTemplate` spec.

## Non-EKS Clusters

By default the machine image is auto-resolved by CAPA to a public AMI that matches the Kubernetes version in `KubeadmControlPlane` or `MachineDeployment` spec. These AMIs are published in a community owned AWS account. See [pre-built public AMIs](built-amis.md) for details of the CAPA project published images.

> IMPORTANT:
> The project doesn't recommend using the public AMIs for production use. Instead its recommended that you build your own AMIs for the Kubernetes versions you want to use. The AMI can then be specified in the `AWSMachineTemplate` spec. [Custom images](custom-amis.md) can be created using [image-builder][image-builder] project.

[image-builder]: https://github.com/kubernetes-sigs/image-builder
