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
- [Minikube][minikube] version v0.30.0 or later
- [kubectl][kubectl]
- [kustomize][kustomize]

### Optional

- [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/installing.html)
- [Homebrew][brew] (MacOS)
- [jq][jq]
- [PowerShell AWS Tools][aws_powershell]
- [Go](https://golang.org/dl/)

## Prerequisites

Get the latest [clusterctl release](https://github.com/kubernetes-sigs/cluster-api-provider-aws/releases) and place it in your path. If a release isn't available, or you might prefer to build the latest version from master you can use `go get sigs.k8s.io/cluster-api/provider-aws/...`.

### Bootstrapping AWS Identity and Access Management with CloudFormation

**NOTE**:
> Your credentials must let you make changes in AWS Identity and Access Management (IAM),
> and use CloudFormation.

The [example](../clusterctl/examples/aws/bootstrap-cloudformation.yaml) CloudFormation template can be used to create IAM roles, users and instance profiles required to bootstrap the cluster.

### Installing clusterawsadm

`clusterawsadm`, is a helper utlity that users might choose to use to quickly setup prerequisites.

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
aws ssm put-parameter --name "/sigs.k8s.io/cluster-api-provider-aws/ssh-key" \
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
but requires Kubernetes 1.12 and the kubeadm bootstrap method to work properly,
so we configure Minikube as follows:

```bash
minikube config set kubernetes-version v1.12.1
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

If you applied the CloudFormation template above, an IAM user was created for you:

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

### Starting Minikube
For security reasons, it's best to start a Minikube instance before using
clusterctl to store your AWS credentials.

``` shell
minikube start
```

Then write the AWS credentials into the Minikube instance for use by Cluster API:

``` shell
minikube ssh 'mkdir -p .aws'
clusterawsadm alpha bootstrap generate-aws-default-profile | minikube ssh 'cat > .aws/credentials'
[default]
aws_access_key_id = AKIAIOSFODNN7EXAMPLE
aws_secret_access_key = wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
region = eu-west-1


^C
```

**NOTE**:
> This will print the AWS credentials to the screen and then hang
> (due to a restriction in Minikube).
> This is OK, and you can cancel with Ctrl+C, but make sure to do this privately.

### Generating cluster manifests

The shell script `config/generate-yaml.sh` can be used to generate the
cluster manifests.

```bash
./config/generate-yaml.sh
```

Then edit `out/cluster.yaml` and `out/machine.yaml` for your SSH key and AWS
region and any other customisations you want to make.

### Starting Cluster API

If you haven't already, set up your [environment](#setting-up-the-environment)
in the [terminal session] you're working in.

You can now start the Cluster API controllers and deploy a new cluster in AWS:

*Bash:*

```bash
clusterctl create cluster -v2 --provider aws \
  -m ./out/machines.yaml \
  -c ./out/cluster.yaml \
  -p ./out/provider-components.yaml \
  --existing-bootstrap-cluster-kubeconfig ~/.kube/config

I1018 01:21:12.079384   16367 clusterdeployer.go:94] Creating bootstrap cluster
I1018 01:21:12.106882   16367 clusterdeployer.go:111] Applying Cluster API stack to bootstrap cluster
I1018 01:21:12.106901   16367 clusterdeployer.go:300] Applying Cluster API Provider Components
I1018 01:21:12.106909   16367 clusterclient.go:505] Waiting for kubectl apply...
I1018 01:21:12.460755   16367 clusterclient.go:533] Waiting for Cluster v1alpha resources to become available...
I1018 01:21:12.464840   16367 clusterclient.go:546] Waiting for Cluster v1alpha resources to be listable...
I1018 01:21:12.517706   16367 clusterdeployer.go:116] Provisioning target cluster via bootstrap cluster
I1018 01:21:12.517722   16367 clusterdeployer.go:118] Creating cluster object aws-provider-test1 on bootstrap cluster in namespace "aws-provider-system"
I1018 01:21:12.524912   16367 clusterdeployer.go:123] Creating master  in namespace "aws-provider-system"
```

*PowerShell:*

```powershell
clusterctl create cluster -v2 --provider aws `
  -m ./config/samples/out/machines.yaml `
  -c ./config/samples/out/cluster.yaml `
  -p ./config/samples/out/provider-components.yaml `
  --existing-bootstrap-cluster-kubeconfig ~/.kube/config

I1018 01:21:12.079384   16367 clusterdeployer.go:94] Creating bootstrap cluster
I1018 01:21:12.106882   16367 clusterdeployer.go:111] Applying Cluster API stack to bootstrap cluster
I1018 01:21:12.106901   16367 clusterdeployer.go:300] Applying Cluster API Provider Components
...
```

## Troubleshooting

Controller logs can be tailed using [`kubectl`][kubectl]:

*Bash:*

```bash
kubectl get po -o name -n aws-provider-system | grep aws-provider-controller-manager | xargs kubectl logs -n aws-provider-system -c manager -f
```

*PowerShell:*

```powershell
kubectl logs -n aws-provider-system -c manager -f `
  $(kubectl get po -o name | Select-String -Pattern "aws-provider-controller-manager")
```

<!-- References -->

[brew]: https://brew.sh/
[jq]: https://stedolan.github.io/jq/download/
[minikube]: https://kubernetes.io/docs/tasks/tools/install-minikube/
[kubectl]: https://kubernetes.io/docs/tasks/tools/install-kubectl/
[aws_powershell]: (https://docs.aws.amazon.com/powershell/index.html#lang/en_us)
[aws-vault]: https://github.com/99designs/aws-vault
[kustomize]: https://github.com/kubernetes-sigs/kustomize
