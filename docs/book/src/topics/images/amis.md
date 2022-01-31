# AWS Machine Images for CAPA Clusters

CAPA requires a “machine image” containing pre-installed, matching versions of kubeadm and kubelet.
Machine image is either auto-resolved by CAPA to a public AMI that matches the Kubernetes version in `KubeadmControlPlane` or `MachineDeployment` spec,
or an appropriate custom image ID for the Kubernetes version can be set in `AWSMachineTemplate` spec.

[Pre-built public AMIs](built-amis.md) are published by the maintainers regularly for each new Kubernetes version.

[Custom images](custom-amis.md) can be created using [image-builder][image-builder] project.

[image-builder]: https://github.com/kubernetes-sigs/image-builder
