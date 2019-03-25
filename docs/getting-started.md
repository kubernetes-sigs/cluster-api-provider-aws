# Getting started with cluster-api-provider-aws <!-- omit in toc -->

## Contents <!-- omit in toc -->

<!-- Below is generated using VSCode yzhang.markdown-all-in-one >

<!-- TOC depthFrom:2 -->

- [Requirements](#requirements)
  - [Optional](#optional)
- [Prerequisites](#prerequisites)
  - [Install release binaries](#install-release-binaries)
  - [Setting up AWS](#setting-up-aws)
    - [`clusterawsadm`](#clusterawsadm)
    - [non-`clusterawsadm`](#non-clusterawsadm)
  - [SSH Key pair](#ssh-key-pair)
    - [Create a new key pair](#create-a-new-key-pair)
    - [Using an existing key](#using-an-existing-key)
- [Deploying a cluster](#deploying-a-cluster)
  - [Setting up the environment](#setting-up-the-environment)
  - [Generating cluster manifests and example cluster](#generating-cluster-manifests-and-example-cluster)
    - [Cluster name](#cluster-name)
    - [Using an existing VPC](#using-an-existing-vpc)
  - [Creating a cluster](#creating-a-cluster)
- [Using the cluster](#using-the-cluster)
- [Troubleshooting](#troubleshooting)
- [Bootstrap running, but resources aren't being created](#bootstrap-running-but-resources-arent-being-created)

<!-- /TOC -->

## Requirements

- Linux or MacOS (Windows isn't supported at the moment)
- A set of AWS credentials sufficient to bootstrap the cluster (see [bootstrapping-aws-identity-and-access-management-with-cloudformation](#bootstrapping-aws-identity-and-access-management-with-cloudformation)).
- An AWS IAM role to give to the Cluster API control plane.
- [KIND >= v0.1](https://sigs.k8s.io/kind)
- gettext (with `envsubst` in your PATH)

### Optional

- [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/installing.html)
- [Homebrew][brew] (MacOS)
- [jq][jq]
- [Go >= v1.11](https://golang.org/dl/)

## Prerequisites

### Install release binaries

Get the latest [release of `clusterctl` and
`clusterawsadm`](https://github.com/kubernetes-sigs/cluster-api-provider-aws/releases)
and place them in your path.

If you prefer to build the latest version from master you can use `go get
sigs.k8s.io/cluster-api-provider-aws/...` – the trailing `...` will ask for both
`clusterctl` and `clusterawsadm` to be built.

### Setting up AWS

#### `clusterawsadm`

Cluster-API-Provider-AWS provides a tool, `clusterawsadm` to help you manage
your AWS IAM objects for this project. In order to use `clusterawsadm`
you must have an administrative user in an AWS account. Once you have that
administrator user you need to set your environment variables:

* `AWS_REGION`
* `AWS_ACCESS_KEY_ID`
* `AWS_SECRET_ACCESS_KEY`

After these are set run this command to get you up and running:

`clusterawsadm alpha bootstrap create-stack`

#### non-`clusterawsadm`

This is not a recommended route as the policies are very specific and will
change with new features.

If you do not wish to use the `clusteradwsadm` tool then you will need to
understand exactly which IAM policies and groups we are expecting. There are
several policies, roles and users that need to be created. Please see our
[controller policy][controllerpolicy] file to understand the permissions that are necessary.

[controllerpolicy]: https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/0e543e0eb30a7065c967f5df8d6abd872aa4ff0c/pkg/cloud/aws/services/cloudformation/bootstrap.go#L149-L188

### SSH Key pair

You will need to specify the name of an existing SSH key pair within the region
you plan on using. If you don't have one yet, a new one needs to be created.

#### Create a new key pair

```bash
# Save the output to a secure location
aws ec2 create-key-pair --key-name default | jq .KeyMaterial -r
-----BEGIN RSA PRIVATE KEY-----
[... contents omitted ...]
-----END RSA PRIVATE KEY-----
```

If you want to save the private key directly into AWS Systems Manager Parameter
Store with KMS encryption for security, you can use the following command:

```bash
aws ssm put-parameter --name "/sigs.k8s.io/cluster-api-provider-aws/ssh-key" \
  --type SecureString \
  --value "$(aws ec2 create-key-pair --key-name default | jq .KeyMaterial -r)"
{
"Version": 1
}
```

#### Using an existing key

```bash
# Replace with your own public key
aws ec2 import-key-pair \
  --key-name default \
  --public-key-material "$(cat ~/.ssh/id_rsa.pub)"
```

> Only RSA keys are supported by AWS.

## Deploying a cluster

### Setting up the environment

The current iteration of the Cluster API Provider AWS relies on credentials
being present in your environment. These then get written into the cluster
manifests for use by the controllers.

If you used `clusterawsadm` to set up IAM resources for you then you can run
these commands to prepare your environment.

Your `AWS_REGION` must already be set.

```bash
export AWS_CREDENTIALS=$(aws iam create-access-key \
  --user-name bootstrapper.cluster-api-provider-aws.sigs.k8s.io)
export AWS_ACCESS_KEY_ID=$(echo $AWS_CREDENTIALS | jq .AccessKey.AccessKeyId -r)
export AWS_SECRET_ACCESS_KEY=$(echo $AWS_CREDENTIALS | jq .AccessKey.SecretAccessKey -r)
```

If you did not use `clusterawsadm` to provision your user you will need to set
these environment variables in your own way.

> To save credentials securely in your environment, [aws-vault][aws-vault] uses
> the OS keystore as permanent storage, and offers shell features to securely
> expose and setup local AWS environments.

### Generating cluster manifests and example cluster

Download the cluster-api-provider-aws-examples.tar file and unpack it.

```bash
tar xfv cluster-api-provider-aws-examples.tar
```

Then run `./generate_yaml.sh` to generate manifests:

```bash
./aws/generate_yaml.sh
```

You should not need to edit the generated manifests, but if you want to do any
customization now is the time to do it. Take a look at
`./aws/out/cluster.yaml` and
`./aws/out/machine.yaml`.

Ensure the `region` and `keyName` are set to what you expect.

#### Cluster name

By default the cluster name is set `test1` when generating manifests, `CLUSTER_NAME` environment variable is used in `make manifests` to customize the cluster name.

```bash
export CLUSTER_NAME="<pony-unicorns>"
```

#### Using an existing VPC

By default the provider creates network resources needed to spin up a cluster on AWS. Users can bring their own network infrastructure and disable the default behavior declaring the `VPC_ID` environment variable before the `make manifests` step.

Features and limitations:
- VPCs are required to have at least 1 private and 1 public subnet available.
- Cluster name must be unique within a region: the limitation comes from the ELB naming scheme which is required to be unique within a region.
- Security groups cannot be customized at the time of writing.


### Creating a cluster

You can now start the Cluster API controllers and deploy a new cluster in AWS:

```bash
clusterctl create cluster -v 3 \
  --bootstrap-type kind \
  --provider aws \
  -m ./aws/out/machines.yaml \
  -c ./aws/out/cluster.yaml \
  -p ./aws/out/provider-components.yaml \
  -a ./aws/out/addons.yaml

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

## Using the cluster
Kubeconfig for the new cluster is created in the directory from where the above `clusterctl create` was run.
Run the following command to point `kubectl` to the kubeconfig of the new cluster
`export KUBECONFIG=$(PWD)/kubeconfig`

Alternatively, move the kubeconfig file to a desired location and set the `KUBECONFIG` environment variable accordingly.

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
