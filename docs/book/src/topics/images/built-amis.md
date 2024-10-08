# Pre-built Kubernetes AMIs

New AMIs are built on a best effort basis when a new Kubernetes version is released for each supported OS distribution and then published to supported regions.

## AMI Publication Policy

- AMIs should only be used for non-production usage. For production environments we recommend that you build and maintain your own AMIs using the image-builder project.
- AMIs will only be published for the latest release series and 2 previous release series. For example, if the current release series is v1.31 then AMIs will only be published for v1.31.x, v1.30.x, v1.29.x.
- When there is a new k8s release series then any AMIs no longer covered by the previous point will be deleted. For example, when v1.32.0 is pubslished then any AMIs for the v1.29 release series will be deleted.
- Existing AMIs are not updated for security fixes and it is recommended to always use the latest patch version for the Kubernetes version you want to run.

## Finding AMIs

`clusterawsadm ami list` command lists pre-built reference AMIs by Kubernetes version, OS, or AWS region. See [clusterawsadm ami list](https://cluster-api-aws.sigs.k8s.io/clusterawsadm/clusterawsadm_ami_list.html) for details.

If you are using a version of clusterawsadm prior to v2.6.2 then you will need to explicitly specify the owner-id for the community account: `clusterawsadm ami list --owner-id 819546954734`.

## Supported OS Distributions
- Amazon Linux 2 (amazon-2)
- Ubuntu (ubuntu-22.04, ubuntu-24.04)
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
