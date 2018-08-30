# Building Images

This directory contains tooling for building base images for use as nodes in Kubernetes Clusters. [Packer](https://www.packer.io) is used for building these images.

## Prerequisites

### Prerequisites for all images

- [Packer](https://www.packer.io/docs/installation.html)
- [Ansible](http://docs.ansible.com/ansible/latest/intro_installation.html) version >= 2.4.0.0

### Prerequisites for Amazon Web Services

- An AWS account
- The AWS CLI installed and configured

## Building Images

### Build Variables

The following variables can be overriden when building images using the `-var` option when calling `packer build`:

| Variable | Default | Description |
|----------|---------|-------------|
| kubernetes_version | 1.11.2-00 | Kubernetes Version to install |
| kubernetes_cni_version | 0.6.0-00 | CNI Version to install |

For example, to build all images for use with Kubernetes 1.11.2 for build version 1:

```sh
packer build -var kubernetes_version=1.11.2-00
```

There are additional variables that may be set that affect the behavior of specific builds or packer post-processors. `packer inspect packer.json` will list all available variables and their default values.

### Limiting Images to Build

If packer build is run without specifying which images to build, then it will attempt to build all configured images. `packer inspect packer.json` will list the configured builders. The `--only` option can be specified when running `packer build` to limit the images built.

For example, to build only the Ubuntu image:

```sh
packer build --only=ami-ubuntu packer.json
```

### Required Permissions to Build the AWS AMIs

The [Packer documentation for the Amazon AMI builder](https://www.packer.io/docs/builders/amazon.html) supplies a suggested set of minimum permissions.

```json
{
  "Version": "2012-10-17",
  "Statement": [{
      "Effect": "Allow",
      "Action" : [
        "ec2:AttachVolume",
        "ec2:AuthorizeSecurityGroupIngress",
        "ec2:CopyImage",
        "ec2:CreateImage",
        "ec2:CreateKeypair",
        "ec2:CreateSecurityGroup",
        "ec2:CreateSnapshot",
        "ec2:CreateTags",
        "ec2:CreateVolume",
        "ec2:DeleteKeyPair",
        "ec2:DeleteSecurityGroup",
        "ec2:DeleteSnapshot",
        "ec2:DeleteVolume",
        "ec2:DeregisterImage",
        "ec2:DescribeImageAttribute",
        "ec2:DescribeImages",
        "ec2:DescribeInstances",
        "ec2:DescribeRegions",
        "ec2:DescribeSecurityGroups",
        "ec2:DescribeSnapshots",
        "ec2:DescribeSubnets",
        "ec2:DescribeTags",
        "ec2:DescribeVolumes",
        "ec2:DetachVolume",
        "ec2:GetPasswordData",
        "ec2:ModifyImageAttribute",
        "ec2:ModifyInstanceAttribute",
        "ec2:ModifySnapshotAttribute",
        "ec2:RegisterImage",
        "ec2:RunInstances",
        "ec2:StopInstances",
        "ec2:TerminateInstances"
      ],
      "Resource" : "*"
  }]
}
```

### Building the AMIs

Building images requires setting additional variables not set by default. The `base-images-us-east-1.json` file is provided as an example.

To build both the Ubuntu and CentOS AMIs:

```sh
packer build -var-file base-images-us-east-1.json packer.json
```

By default images are copied to all available AWS regions. The list can be obtained running:
```sh
aws ec2 describe-regions --query "Regions[].{Name:RegionName}" --output text | paste -sd "," -
```

To limit the regions, provide the `ami_regions` variable as a comma-delimited list of AWS regions. 

For example, to build all images in us-east-1 and copy only to us-west-2:
```sh
packer build -var-file base-images-us-east-1.json -var ami_regions='us-west-2'
```

## Testing Images

Connect remotely to an instance created from the image and run the Node Conformance tests using the following commands:

```sh
wget https://dl.k8s.io/$(< /etc/kubernetes_community_ami_version)/kubernetes-test.tar.gz
tar -zxvf kubernetes-test.tar.gz kubernetes/platforms/linux/amd64
cd kubernetes/platforms/linux/amd64
sudo ./ginkgo --nodes=8 --flakeAttempts=2 --focus="\[Conformance\]" --skip="\[Flaky\]|\[Serial\]|\[sig-network\]|Container Lifecycle Hook" ./e2e_node.test -- --k8s-bin-dir=/usr/bin
```