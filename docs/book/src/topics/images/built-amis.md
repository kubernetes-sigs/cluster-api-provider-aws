# Pre-built Kubernetes AMIs

New AMIs are built on a best effort basis when a new Kubernetes version is released.  for each supported OS distribution and then published to supported regions.

`clusterawsadm ami list` command lists pre-built reference AMIs by Kubernetes version, OS, or AWS region. See [clusterawsadm ami list](https://cluster-api-aws.sigs.k8s.io/clusterawsadm/clusterawsadm_ami_list.html) for details.

If you are using a version of clusterawsadm prior to v2.6.2 then you will need to explicitly specify the owner-id for the community account: `clusterawsadm ami list --owner-id 819546954734`.

> **Note:**  These images are not updated for security fixes and it is recommended to always use the latest patch version for the Kubernetes version you want to run. For production environments, it is highly recommended to build and use your own custom images.

## Supported OS Distributions
- Amazon Linux 2 (amazon-2)
- Ubuntu (ubuntu-20.04, ubuntu-22.04, ubuntu-24.04)
- Centos (centos-7)
- Flatcar (flatcar-stable)

## Supported AWS Regions
- ap-northeast-1
- ap-northeast-2
- ap-south-1
- ap-southeast-1
- ap-southeast-2
- ca-central-1
- eu-central-1
- eu-west-1
- eu-west-2
- eu-west-3
- sa-east-1
- us-east-1
- us-east-2
- us-west-1
- us-west-2
