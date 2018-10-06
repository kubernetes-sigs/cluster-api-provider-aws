# Getting started with cluster-api-provider-aws <!-- omit in toc -->

## Contents <!-- omit in toc -->

<!-- Below is generated using VSCode yzhang.markdown-all-in-one >

<!-- TOC depthFrom:2 -->

- [Requirements](#requirements)
  - [Optional](#optional)
- [Prerequisites](#prerequisites)
  - [Bootstrapping AWS Identity and Access Management with CloudFormation](#bootstrapping-aws-identity-and-access-management-with-cloudformation)
  - [Installing clusterawsadm](#installing-clusterawsadm)
  - [SSH Key pair](#ssh-key-pair)
    - [Create a new key pair](#create-a-new-key-pair)
    - [Using an existing key](#using-an-existing-key)
  - [Setting up Minikube](#setting-up-minikube)
    - [Customizing for Cluster API](#customizing-for-cluster-api)
  - [Setting up the environment](#setting-up-the-environment)
- [Deploying a cluster](#deploying-a-cluster)
  - [Generating cluster manifests](#generating-cluster-manifests)
  - [Starting Cluster API](#starting-cluster-api)
- [Troubleshooting](#troubleshooting)

<!-- /TOC -->

## Requirements

- Linux or MacOS (Windows isn't supported at the moment)
- A set of AWS credentials sufficient to bootstrap the cluster (see [bootstrapping-aws-identity-and-access-management-with-cloudformation](#bootstrapping-aws-identity-and-access-management-with-cloudformation)).
- An AWS IAM role to give to the Cluster API control plane.
- [Minikube][minikube]
- [kubectl][kubectl]

### Optional

- [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/installing.html)
- [Homebrew][brew] (MacOS)
- [jq][jq]
- [PowerShell AWS Tools][aws_powershell]
- [Go](https://golang.org/dl/)

## Prerequisites

Get the latest [clusterctl release](https://github.com/kubernetes-sigs/cluster-api-provider-aws/releases) and place it in your path. If a release isn't available, or you might prefer to build the latest version from master you can build it with `go get sigs.k8s.io/cluster-api/provider-aws/...`.

### Bootstrapping AWS Identity and Access Management with CloudFormation

**NOTE**:
> Your credentials must let you make changes in AWS Identity and Access Management (IAM),
> and use CloudFormation.

The [example](../clusterctl/examples/aws/bootstrap-cloudformation.yaml)) CloudFormation can be used to create IAM roles, users and instance profiles required to bootstrap the cluster.

### Installing clusterawsadm

`clusterawsadm`, can be used to create IAM

> NOTE: This command requires to have a working AWS environment.

```bash
clusterawsadm alpha bootstrap create-stack
```

### SSH Key pair

You will need to specify the name of an existing SSH key pair within the region you plan on using. If you don't have one yet, a new one needs to be created.

#### Create a new key pair

*Bash:*

```bash
# Save the output to a secure location
aws ec2 create-key-pair --key-name cluster-api-provider-aws.sigs.k8s.io | jq .KeyMaterial -r
-----BEGIN RSA PRIVATE KEY-----
[... contents omitted ...]
-----END RSA PRIVATE KEY-----
```

*PowerShell:*

```powershell
(New-EC2KeyPair -KeyName cluster-api-provider-aws.sigs.k8s.io).KeyMaterial
-----BEGIN RSA PRIVATE KEY-----
[... contents omitted ...]
-----END RSA PRIVATE KEY-----
```

If you want to save the private key directly into AWS Systems Manager Parameter
Store with KMS encryption for security, you can use the following command:

*Bash:*

```bash
aws ssm put-parameter --name "/sigs.k8s.io/cluster-api-provider-aws/ssh-key \
  --type SecureString \
  --value "$(aws ec2 create-key-pair --key-name cluster-api-provider-aws.sigs.k8s.io | jq .KeyMaterial -r)"
{
"Version": 1
}
```

*PowerShell:*

```powershell
Write-SSMParameter -Name "/sigs.k8s.io/cluster-api-provider-aws/ssh-key" `
  -Type SecureString `
  -Value (New-EC2KeyPair -KeyName cluster-api-provider-aws.sigs.k8s.io).KeyMaterial
1
```

#### Using an existing key

*Bash:*

```bash
# Replace with your own public key
aws ec2 import-key-pair --key-name cluster-api-provider-aws.sigs.k8s.io \
  --public-key-material $(cat ~/.ssh/id_rsa.pub)
```

*PowerShell:*

```powershell
$publicKey = [System.Convert]::ToBase64String( `
  [System.Text.Encoding]::UTF8.GetBytes(((get-content ~/.ssh/id_rsa.pub))))
Import-EC2KeyPair -KeyName cluster-api-provider-aws.sigs.k8s.io -PublicKeyMaterial $publicKey
```

**NOTE**:
> Only RSA keys are supported by AWS.

### Setting up Minikube

Minikube needs to be installed on your local machine, as this is what will be used by the Cluster API to bootstrap your cluster in AWS.

[Instructions for setting up minikube][minikube] are available on the Kubernetes website.

#### Customizing for Cluster API

At present, the Cluster API provider runs minikube to create a new instance,
but requires Kubernetes 1.9 and the kubeadm bootstrap method to work properly,
so we configure Minikube as follows:

```bash
minikube config set kubernetes-version v1.9.4
minikube config set bootstrapper kubeadm
```

### Setting up the environment

The current iteration of the Cluster API Provider AWS relies on credentials being present in your environment.
These then get written into the cluster manifests for use by the controllers.

*Bash:*

```bash
# Region used to deploy the cluster in.
export AWS_REGION=us-east-1

# User access credentials.
export AWS_ACCESS_KEY_ID="AKIAIOSFODNN7EXAMPLE"
export AWS_SECRET_ACCESS_KEY="wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"

# SSH Key to be used to run instances.
export SSH_KEY_NAME="cluster-api-provider-aws.sigs.k8s.io"
```

*PowerShell:*

```powershell
$ENV:AWS_REGION = "us-east-1"
$ENV:AWS_ACCESS_KEY_ID="AKIAIOSFODNN7EXAMPLE"
$ENV:AWS_SECRET_ACCESS_KEY="wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
$ENV:SSH_KEY_NAME="cluster-api-provider-aws.sigs.k8s.io"
```

If you applied the CloudFormation template above, use the IAM user that it created:

*Bash:*

```bash
export AWS_CREDENTIALS=$(aws iam create-access-key \
  --user-name bootstrapper.cluster-api-provider-aws.sigs.k8s.io)
export AWS_ACCESS_KEY_ID=$(echo $AWS_CREDENTIALS | jq .AccessKey.AccessKeyId -r)
export AWS_SECRET_ACCESS_KEY=$(echo $AWS_CREDENTIALS | jq .AccessKey.SecretAccessKey -r)
```

*PowerShell:*

```powershell
$awsCredentials = New-IAMAccessKey -UserName bootstrapper.cluster-api-provider-aws.sigs.k8s.io
$ENV:AWS_ACCESS_KEY_ID=$awsCredentials.AccessKeyId
$ENV:AWS_SECRET_ACCESS_KEY=$awsCredentials.SecretAccessKey
```

**NOTE**:
> To save credentials securely in your environment, [aws-vault][aws-vault] uses the OS keystore as permanent storage,
> and offers shell features to securely expose and setup local AWS environments.

## Deploying a cluster

### Generating cluster manifests

Use the shell script `generate-yaml.sh` in [clusterctl/examples](clusterctl/examples) to generate the manifests.

> The following command is valid in both Bash and PowerShell.

```bash
sh -c "cd ./clusterctl/examples/aws && ./generate-yaml.sh"
```

**NOTE**:
> The generated manifests contain a copy of the AWS credentials.
> Secure credentials storage is slated for a future release.

### Starting Cluster API

You can now start the Cluster API controllers and deploy a new cluster in AWS:

*Bash:*

```bash
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

```powershell
clusterctl create cluster -v2 --provider aws `
  -m ./clusterctl/examples/aws/out/machines.yaml `
  -c ./clusterctl/examples/aws/out/cluster.yaml `
  -p ./clusterctl/examples/aws/out/provider-components.yaml

I1005 16:18:54.768403   22094 clusterdeployer.go:95] Creating bootstrap cluster
I1005 16:20:23.611501   22094 minikube.go:50] Ran: minikube [start --bootstrapper=kubeadm] Output: Starting local Kubernetes v1.9.4 cluster...
...
```

## Troubleshooting

Controller logs can be tailed using [`kubectl`][kubectl]:

*Bash:*

```bash
kubectl get po -o name | grep clusterapi-controllers | xargs kubectl logs -c aws-cluster-controller -f
```

*PowerShell:*

```powershell
kubectl logs -c aws-cluster-controller -f `
  $(kubectl get po -o name | Select-String -Pattern "clusterapi-controllers")
```

<!-- References -->

[brew]: https://brew.sh/
[jq]: https://stedolan.github.io/jq/download/
[minikube]: https://kubernetes.io/docs/tasks/tools/install-minikube/
[kubectl]: https://kubernetes.io/docs/tasks/tools/install-kubectl/
[aws_powershell]: (https://docs.aws.amazon.com/powershell/index.html#lang/en_us)
[aws-vault]: https://github.com/99designs/aws-vault
