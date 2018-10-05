# Getting started with cluster-api-provider-aws <!-- omit in toc -->

## Contents <!-- omit in toc -->

<!-- Below is generated using VSCode yzhang.markdown-all-in-one >

<!-- TOC depthFrom:2 -->

- [Requirements](#requirements)
  - [Optional](#optional)
- [Prerequisites](#prerequisites)
  - [Bootstrapping AWS Identity and Access Management with CloudFormation](#bootstrapping-aws-identity-and-access-management-with-cloudformation)
  - [Installing clusterawsadm](#installing-clusterawsadm)
  - [SSH key pair](#ssh-key-pair)
    - [Create a fresh key pair](#create-a-fresh-key-pair)
    - [Using an existing key](#using-an-existing-key)
  - [Setting up Minikube](#setting-up-minikube)
    - [Customising for Cluster API](#customising-for-cluster-api)
  - [Setting up the environment](#setting-up-the-environment)
- [Deploying a cluster](#deploying-a-cluster)
  - [Generating cluster manifests](#generating-cluster-manifests)
  - [Starting Cluster API](#starting-cluster-api)
- [Troubleshooting](#troubleshooting)

<!-- /TOC -->

## Requirements

- A set of AWS credentials sufficient to bootstrap the cluster (see [#aws-cloudformation-setup](#aws-cloudformation-setup)).
- A AWS IAM role to give to the Cluster API control plane.
- [Minikube][minikube]
- [kubectl][kubectl]
- Linux or OS X (WSL isn't supported at present)

### Optional

- The [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/installing.html) [[Homebrew][aws_brew]] + [jq][jq] or [AWS Tools for PowerShell][aws_powershell]

## Prerequisites

Get the [latest version](../releases) of the Cluster API Provider AWS release of `clusterctl` and place it in your path.

### Bootstrapping AWS Identity and Access Management with CloudFormation

**NOTE**:
> Your credentials must let you make changes in AWS Identity and Access Management,
> and use CloudFormation.

An [AWS CloudFormation template is provided](../clusterctl/examples/aws/bootstrap-cloudformation.yaml)
to be able to create an AWS IAM User with sufficient authorisations to use
set up a cluster, as well as AWS IAM roles and instance profiles for use within
the cluster.

### Installing clusterawsadm

We provide an additional utility, `clusterawsadm` that provides AWS specific
tooling that is not directly part of Cluster API, e.g. applying IAM policies
via CloudFormation.

Using a set of AWS credentials already present, via either environment variables,
a default AWS CLI profile or via another method, apply the default template using
the following command:

``` bash
clusterawsadm alpha bootstrap create-stack
```

### SSH key pair

You will need to specify the name of an SSH key pair that is stored on EC2
within the region you want to deploy into. If you don't have one set up,
you can do the following.

#### Create a fresh key pair

*Bourne shell:*

``` bash
# Save the output to a secure location
aws ec2 create-key-pair --key-name cluster-api-provider-aws.sigs.k8s.io | jq .KeyMaterial -r
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEApyAOuq47vLNRph47VeL5IUrla6rL2V3mATitbk8rq9MmGVKo5iilJOzJMJ/R
7c+CYhPDT3GGdWYsJnxhSxSLi5rFXc4UyinRfRVcfKlqdzriUQ1k6qgYMdRm3hS9crCKLIYGzzri
...
XJi9MQ2vgTuUKRRczWbsOqFqd5Sl6ZPrrf/SWp65sP35h+BM2oSQmDHJcAT8Bqy0NMDRI2VZHLZF
1S+Ll6Dv539G8s88vlE35SU+DJOioeJEE979b92MYNjownDAQUZENJKddHepYtNvwPAK
-----END RSA PRIVATE KEY-----
```

*PowerShell:*

``` powershell
(New-EC2KeyPair -KeyName cluster-api-provider-aws.sigs.k8s.io).KeyMaterial
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEApyAOuq47vLNRph47VeL5IUrla6rL2V3mATitbk8rq9MmGVKo5iilJOzJMJ/R
7c+CYhPDT3GGdWYsJnxhSxSLi5rFXc4UyinRfRVcfKlqdzriUQ1k6qgYMdRm3hS9crCKLIYGzzri
...
XJi9MQ2vgTuUKRRczWbsOqFqd5Sl6ZPrrf/SWp65sP35h+BM2oSQmDHJcAT8Bqy0NMDRI2VZHLZF
1S+Ll6Dv539G8s88vlE35SU+DJOioeJEE979b92MYNjownDAQUZENJKddHepYtNvwPAK
-----END RSA PRIVATE KEY-----
```

If you want to save the private key directly into AWS Systems Manager Parameter
Store with KMS encryption for security, you can use the following command:

*Bourne shell:*

``` bash
aws ssm put-parameter --name "/sigs.k8s.io/cluster-api-provider-aws/ssh-key \
  --type SecureString \
  --value "$(aws ec2 create-key-pair --key-name cluster-api-provider-aws.sigs.k8s.io | jq .KeyMaterial -r)"
{
"Version": 1
}
```

*PowerShell:*

``` powershell
Write-SSMParameter -Name "/sigs.k8s.io/cluster-api-provider-aws/ssh-key" `
  -Type SecureString `
  -Value (New-EC2KeyPair -KeyName cluster-api-provider-aws.sigs.k8s.io).KeyMaterial
1
```

#### Using an existing key

*Bourne shell:*

``` bash
# Replace with your own public key
aws ec2 import-key-pair --key-name cluster-api-provider-aws.sigs.k8s.io \
  --public-key-material $(cat ~/.ssh/id_rsa.pub)
```

*PowerShell:*

``` powershell
$publicKey = [System.Convert]::ToBase64String( `
  [System.Text.Encoding]::UTF8.GetBytes(((get-content ~/.ssh/id_rsa.pub))))
Import-EC2KeyPair -KeyName cluster-api-provider-aws.sigs.k8s.io -PublicKeyMaterial $publicKey
```

**NOTE**:
> Only RSA keys are supported by AWS.

### Setting up Minikube

Minikube needs to be installed on your local machine, as this is what will
be used by the Cluster API to bootstrap your cluster in AWS.

[Instructions for setting up minikube][minikube] are on the Kubernetes website.

#### Customising for Cluster API

At present, the Cluster API provider runs minikube to create a new instance,
but requires Kubernetes 1.9 and the kubeadm bootstrap method to work properly,
so we configure Minikube as follows:

``` bash
minikube config set kubernetes-version v1.9.4
minikube config set bootstrapper kubeadm
```

You should also select the hypervisor that works best for your environment.

People have had the success with the following:

- OS X: `minikube config set vm-driver hyperkit`
- Linux: `minikube config set vm-driver kvm2`

### Setting up the environment

The current iteration of the Cluster API Provider AWS relies on credentials
being present in your environment. These then get written into the cluster
manifests for use by the controllers.

You therefore need to set the following environment variables:

*Bourne shell:*

``` bash
# The region you will be running the cluster components in
export AWS_REGION=us-east-1
# The AWS credentials for bootstrap
export AWS_ACCESS_KEY_ID="AKIAIOSFODNN7EXAMPLE"
export AWS_SECRET_ACCESS_KEY="wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
# The name of the SSH key stored in the EC2 API
export SSH_KEY_NAME="cluster-api-provider-aws.sigs.k8s.io"
```

*PowerShell:*

``` powershell
$ENV:AWS_REGION = "us-east-1"
$ENV:AWS_ACCESS_KEY_ID="AKIAIOSFODNN7EXAMPLE"
$ENV:AWS_SECRET_ACCESS_KEY="wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
$ENV:SSH_KEY_NAME="cluster-api-provider-aws.sigs.k8s.io"
```

If you applied the bootstrap AWS CloudFormation, then you can use the IAM user
that it created to scope down the permissions:

*Bourne shell:*

``` shell
export AWS_CREDENTIALS=$(aws iam create-access-key \
  --user-name bootstrapper.cluster-api-provider-aws.sigs.k8s.io)
export AWS_ACCESS_KEY_ID=$(echo $AWS_CREDENTIALS | jq .AccessKey.AccessKeyId -r)
export AWS_SECRET_ACCESS_KEY=$(echo $AWS_CREDENTIALS | jq .AccessKey.SecretAccessKey -r)
```

*PowerShell:*

``` powershell
$awsCredentials = New-IAMAccessKey -UserName bootstrapper.cluster-api-provider-aws.sigs.k8s.io
$ENV:AWS_ACCESS_KEY_ID=$awsCredentials.AccessKeyId
$ENV:AWS_SECRET_ACCESS_KEY=$awsCredentials.SecretAccessKey
```

**NOTE**:
> If you want to save these credentials securely, use a tool like [aws-vault][aws-vault]
> to store them in your operating system key store, which has commands to set
> environment variables for you.

## Deploying a cluster

### Generating cluster manifests

Now that the environment is set up, you can generate the cluster manifests.

Use the `generate-yaml.sh` shell script in [clusterctl/examples](clusterctl/examples)
to use this:

``` bash
# This works in PowerShell too
sh -c "cd ./clusterctl/examples/aws && ./generate-yaml.sh"
```

**NOTE**:
> The manifests will contain a copy of the AWS credentials. Secure storage of
> these is slated for a future release, but please ensure you do not check in
> these files into public repositories.

### Starting Cluster API

You can now start the Cluster API controllers and deploy a new cluster in AWS:

*Bourne shell:*

``` bash
clusterctl create cluster -v2 --provider aws \
  -m ./clusterctl/examples/aws/out/machines.yaml \
  -c ./clusterctl/examples/aws/out/cluster.yaml \
  -p ./clusterctl/examples/aws/out/provider-components.yaml
I1005 16:18:54.768403   22094 clusterdeployer.go:95] Creating bootstrap cluster
I1005 16:20:23.611501   22094 minikube.go:50] Ran: minikube [start --bootstrapper=kubeadm] Output: Starting local Kubernetes v1.9.4 cluster...
Starting VM...
Getting VM IP address...
Moving files into cluster...
Setting up certs...
Connecting to cluster...
Setting up kubeconfig...
Starting cluster components...
Kubectl is now configured to use the cluster.
Loading cached images from config file.
I1005 16:20:23.636837   22094 clusterdeployer.go:112] Applying Cluster API stack to bootstrap cluster
I1005 16:20:23.636871   22094 clusterdeployer.go:301] Applying Cluster API APIServer
I1005 16:20:24.048512   22094 clusterclient.go:511] Waiting for kubectl apply...
I1005 16:20:24.776374   22094 clusterclient.go:539] Waiting for Cluster v1alpha resources to become available...
I1005 16:20:34.777530   22094 clusterclient.go:539] Waiting for Cluster v1alpha resources to become available...
I1005 16:20:44.777572   22094 clusterclient.go:539] Waiting for Cluster v1alpha resources to become available...
I1005 16:20:54.777565   22094 clusterclient.go:539] Waiting for Cluster v1alpha resources to become available...
I1005 16:21:04.777556   22094 clusterclient.go:539] Waiting for Cluster v1alpha resources to become available...
I1005 16:21:04.786146   22094 clusterclient.go:552] Waiting for Cluster v1alpha resources to be listable...
I1005 16:21:15.668990   22094 clusterdeployer.go:307] Applying Cluster API Provider Components
I1005 16:21:15.669022   22094 clusterclient.go:511] Waiting for kubectl apply...
I1005 16:21:15.883375   22094 clusterdeployer.go:117] Provisioning target cluster via bootstrap cluster
I1005 16:21:15.883402   22094 clusterdeployer.go:119] Creating cluster object test1 on bootstrap cluster in namespace "default"
I1005 16:21:15.902219   22094 clusterdeployer.go:124] Creating master  in namespace "default"
I1005 16:21:15.926980   22094 clusterclient.go:563] Waiting for Machine aws-controlplane-x7gxx to become ready...
```

*PowerShell:*

``` powershell
clusterctl create cluster -v2 --provider aws `
  -m ./clusterctl/examples/aws/out/machines.yaml `
  -c ./clusterctl/examples/aws/out/cluster.yaml `
  -p ./clusterctl/examples/aws/out/provider-components.yaml
I1005 16:18:54.768403   22094 clusterdeployer.go:95] Creating bootstrap cluster
I1005 16:20:23.611501   22094 minikube.go:50] Ran: minikube [start --bootstrapper=kubeadm] Output: Starting local Kubernetes v1.9.4 cluster...
...
```

## Troubleshooting

You can get more detailed logs of what is happening in the cluster using [`kubectl`][kubectl]:

*Bourne shell:*

``` bash
kubectl get po -o name | grep clusterapi-controllers \
  | xargs kubectl logs -c aws-cluster-controller -f
```

*PowerShell:*

``` powershell
kubectl logs -c aws-cluster-controller -f `
  $(kubectl get po -o name | Select-String -Pattern "clusterapi-controllers")
```

<!-- References -->

[aws_brew]: https://stedolan.github.io/jq/download/
[jq]: https://stedolan.github.io/jq/download/
[minikube]: https://kubernetes.io/docs/tasks/tools/install-minikube/
[kubectl]: https://kubernetes.io/docs/tasks/tools/install-kubectl/
[aws_powershell]: (https://docs.aws.amazon.com/powershell/index.html#lang/en_us)
[aws-vault]: https://github.com/99designs/aws-vault
