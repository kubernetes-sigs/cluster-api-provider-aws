# Getting started with cluster-api-provider-aws <!-- omit in toc -->

## Contents <!-- omit in toc -->

<!-- Below is generated using VSCode yzhang.markdown-all-in-one >

<!-- TOC depthFrom:2 -->

- [Requirements](#requirements)
  - [Optional](#optional)
- [Prerequisites](#prerequisites)
  - [Installing clusterawsadm](#installing-clusterawsadm)
  - [Bootstrapping AWS Identity and Access Management with CloudFormation](#bootstrapping-aws-identity-and-access-management-with-cloudformation)
  - [SSH Key pair](#ssh-key-pair)
    - [Create a new key pair](#create-a-new-key-pair)
    - [Using an existing key](#using-an-existing-key)
  - [Setting up KIND](#setting-up-kind)
  - [Setting up the environment](#setting-up-the-environment)
- [Deploying a cluster](#deploying-a-cluster)
  - [Generating cluster manifests](#generating-cluster-manifests)
  - [Starting Cluster API](#starting-cluster-api)
- [Troubleshooting](#troubleshooting)
  - [Bootstrap running, but resources aren't being created](#bootstrap-running-but-resources-aren't-being-created)

<!-- /TOC -->

## Requirements

- Linux or MacOS (Windows isn't supported at the moment)
- A set of AWS credentials sufficient to bootstrap the cluster (see [bootstrapping-aws-identity-and-access-management-with-cloudformation](#bootstrapping-aws-identity-and-access-management-with-cloudformation)).
- An AWS IAM role to give to the Cluster API control plane.
- [KIND](https://sigs.k8s.io/kind)
- [kubectl][kubectl]
- [kustomize][kustomize]
- make
- gettext (with `envsubst` in your PATH)
- bazel

### Optional

- [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/installing.html)
- [Homebrew][brew] (MacOS)
- [jq][jq]
- [Go](https://golang.org/dl/)

## Prerequisites

Get the latest [release of `clusterctl` and `clusterawsadm`](https://github.com/kubernetes-sigs/cluster-api-provider-aws/releases) and place it in your path. If a release isn't available, or you might prefer to build the latest version from master you can use `go get sigs.k8s.io/cluster-api-provider-aws/...` – the trailing `...` will ask for both `clusterctl` and `clusterawsadm` to be built.

Before launching clusterctl, you need to define a few environment variables (`AWS_REGION`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`). You thus need an AWS user with sufficient permissions:

1. You can create that user and assign the permissions manually.
2. Or you can use the `clusterawsadm` tool.

### Installing clusterawsadm

`clusterawsadm`, is a helper utlity that users might choose to use to quickly setup prerequisites. It can be installed as per the previous section, by either downloading a release or using `go get` to build it.

> NOTE: The `clusterawsadm` command requires to have a working AWS environment.


### Bootstrapping AWS IAM with CloudFormation

> Your credentials must let you make changes in AWS Identity and Access Management (IAM),
> and use CloudFormation.

```bash
export AWS_REGION=us-east-1
clusterawsadm alpha bootstrap create-stack
```

### SSH Key pair

You will need to specify the name of an existing SSH key pair within the region you plan on using. If you don't have one yet, a new one needs to be created.

#### Create a new key pair

```bash
# Save the output to a secure location
aws ec2 create-key-pair --key-name cluster-api-provider-aws.sigs.k8s.io | jq .KeyMaterial -r
-----BEGIN RSA PRIVATE KEY-----
[... contents omitted ...]
-----END RSA PRIVATE KEY-----
```

If you want to save the private key directly into AWS Systems Manager Parameter
Store with KMS encryption for security, you can use the following command:

```bash
aws ssm put-parameter --name "/sigs.k8s.io/cluster-api-provider-aws/ssh-key" \
  --type SecureString \
  --value "$(aws ec2 create-key-pair --key-name cluster-api-provider-aws.sigs.k8s.io | jq .KeyMaterial -r)"
{
"Version": 1
}
```

#### Using an existing key

```bash
# Replace with your own public key
aws ec2 import-key-pair \
  --key-name cluster-api-provider-aws.sigs.k8s.io \
  --public-key-material $(cat ~/.ssh/id_rsa.pub)
```

> Only RSA keys are supported by AWS.

### Setting up KIND

[KIND](https://sigs.k8s.io/kind) is used by clusterctl as an ephemeral bootstrap cluster. Follow these [instructions][https://github.com/kubernetes-sigs/kind#installation-and-usage] to install it on your local machine.

### Setting up the environment

The current iteration of the Cluster API Provider AWS relies on credentials being present in your environment.
These then get written into the cluster manifests for use by the controllers.

```bash
# Region used to deploy the cluster in.
export AWS_REGION=us-east-1

# User access credentials.
export AWS_ACCESS_KEY_ID="AKIAIOSFODNN7EXAMPLE"
export AWS_SECRET_ACCESS_KEY="wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"

# SSH Key to be used to run instances.
export SSH_KEY_NAME="cluster-api-provider-aws.sigs.k8s.io"
```

If you applied the CloudFormation template above, an IAM user was created for you:

```bash
export AWS_REGION=us-east-1
export AWS_CREDENTIALS=$(aws iam create-access-key \
  --user-name bootstrapper.cluster-api-provider-aws.sigs.k8s.io)
export AWS_ACCESS_KEY_ID=$(echo $AWS_CREDENTIALS | jq .AccessKey.AccessKeyId -r)
export AWS_SECRET_ACCESS_KEY=$(echo $AWS_CREDENTIALS | jq .AccessKey.SecretAccessKey -r)
```

> To save credentials securely in your environment, [aws-vault][aws-vault] uses the OS keystore as permanent storage,
> and offers shell features to securely expose and setup local AWS environments.

## Deploying a cluster

### Generating cluster manifests

There is a make target `manifests` that can be used to generate the
cluster manifests.

```bash
make manifests
```

Then edit `cmd/clusterctl/examples/aws/out/cluster.yaml` and `cmd/clusterctl/examples/aws/out/machine.yaml`. Ensure that the `keyName` is set to the `cluster-api-provider-aws.sigs.k8s.io` we set up above. This is also an opportunity to edit the AWS
region and any apply other customisations you want to make.

> Note: The generated manifests may refer to a keypair named `default`, which differs from the keypair created in this guide. That can be overridden by setting the `SSH_KEY_NAME` env var before running `make manifests`.

### Creating a cluster

You can now start the Cluster API controllers and deploy a new cluster in AWS:

```bash
clusterctl create cluster -v 3 \
  --bootstrap-type kind \
  --provider aws \
  -m ./cmd/clusterctl/examples/aws/out/machines.yaml \
  -c ./cmd/clusterctl/examples/aws/out/cluster.yaml \
  -p ./cmd/clusterctl/examples/aws/out/provider-components.yaml \
  -a ./cmd/clusterctl/examples/aws/out/addons.yaml

I0119 12:16:07.521123   38557 plugins.go:39] Registered cluster provisioner "aws"
I0119 12:16:07.522563   38557 createbootstrapcluster.go:27] Creating bootstrap cluster
I0119 12:16:07.522573   38557 kind.go:53] Running: kind [create cluster --name=clusterapi]
I0119 12:16:40.661674   38557 kind.go:56] Ran: kind [create cluster --name=clusterapi] Output: Creating cluster 'kind-clusterapi' ...
 • Ensuring node image (kindest/node:v1.12.2) 🖼  ...
 ✓ Ensuring node image (kindest/node:v1.12.2) 🖼
 • [kind-clusterapi-control-plane] Creating node container 📦  ...
 ✓ [kind-clusterapi-control-plane] Creating node container 📦
 • [kind-clusterapi-control-plane] Fixing mounts 🗻  ...
 ✓ [kind-clusterapi-control-plane] Fixing mounts 🗻
 • [kind-clusterapi-control-plane] Starting systemd 🖥  ...
 ✓ [kind-clusterapi-control-plane] Starting systemd 🖥
 • [kind-clusterapi-control-plane] Waiting for docker to be ready 🐋  ...
 ✓ [kind-clusterapi-control-plane] Waiting for docker to be ready 🐋
 • [kind-clusterapi-control-plane] Starting Kubernetes (this may take a minute) ☸  ...
 ✓ [kind-clusterapi-control-plane] Starting Kubernetes (this may take a minute) ☸
Cluster creation complete. You can now use the cluster with:

export KUBECONFIG="$(kind get kubeconfig-path --name="clusterapi")"
kubectl cluster-info
I0119 12:16:40.661740   38557 kind.go:53] Running: kind [get kubeconfig-path --name=clusterapi]
I0119 12:16:40.686496   38557 kind.go:56] Ran: kind [get kubeconfig-path --name=clusterapi] Output: /path/to/.kube/kind-config-clusterapi
I0119 12:16:40.688189   38557 clusterdeployer.go:95] Applying Cluster API stack to bootstrap cluster
I0119 12:16:40.688199   38557 applyclusterapicomponents.go:26] Applying Cluster API Provider Components
I0119 12:16:40.688207   38557 clusterclient.go:520] Waiting for kubectl apply...
I0119 12:16:40.981186   38557 clusterclient.go:549] Waiting for Cluster v1alpha resources to become available...
I0119 12:16:40.989350   38557 clusterclient.go:562] Waiting for Cluster v1alpha resources to be listable...
I0119 12:16:40.997829   38557 clusterdeployer.go:100] Provisioning target cluster via bootstrap cluster
I0119 12:16:41.002232   38557 applycluster.go:36] Creating cluster object test1 in namespace "default"
I0119 12:16:41.007516   38557 clusterdeployer.go:109] Creating control plane controlplane-0 in namespace "default"
I0119 12:16:41.011616   38557 applymachines.go:36] Creating machines in namespace "default"
I0119 12:16:41.021539   38557 clusterclient.go:573] Waiting for Machine controlplane-0 to become ready...
```

The created KIND cluster is ephemeral and is cleaned up automatically when done. During the cluster creation, the KIND configuration is written to a local directory and can be retrieved using `kind get kubeconfig-path --name="clusterapi"`.

For a more in-depth look into what `clusterctl` is doing during this create step, please see the [clusterctl document](/docs/clusterctl.md).

## Troubleshooting

## Bootstrap running, but resources aren't being created

Logs can be tailed using [`kubectl`][kubectl]:

```bash
export KUBECONFIG=$(kind get kubeconfig-path --name="clusterapi")
kubectl logs -f -n aws-provider-system aws-provider-controller-manager-0
```

<!-- References -->

[brew]: https://brew.sh/
[jq]: https://stedolan.github.io/jq/download/
[kubectl]: https://kubernetes.io/docs/tasks/tools/install-kubectl/
[aws-vault]: https://github.com/99designs/aws-vault
[kustomize]: https://github.com/kubernetes-sigs/kustomize
