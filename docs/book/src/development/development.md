# Developer Guide

## Initial setup for development environment

### Install prerequisites

1. Install [go][go]
    - Get the latest patch version for go v1.22.
2. Install [jq][jq]
    - `brew install jq` on macOS.
    - `chocolatey install jq` on Windows.
    - `sudo apt install jq` on Ubuntu Linux.
3. Install [KIND][kind]
    - `GO111MODULE="on" go get sigs.k8s.io/kind@v0.12.0`.
4. Install [Kustomize][kustomize]
    - [install instructions](https://kubectl.docs.kubernetes.io/installation/kustomize/)
5. Install [envsubst][envsubst]
6. Install make.
7. Install direnv
    - `brew install direnv` on macOS.
8. Set AWS Environment variable for an IAM Admin user
    - ```bash
      export AWS_ACCESS_KEY_ID=ADMID
      export AWS_SECRET_ACCESS_KEY=ADMKEY
      export AWS_REGION=eu-west-1
      ```

### Get the source

Fork the [cluster-api-provider-aws repo](https://github.com/kubernetes-sigs/cluster-api-provider-aws):

```bash
cd "$(go env GOPATH)"/src
mkdir sigs.k8s.io
cd sigs.k8s.io/
git clone git@github.com:<GITHUB USERNAME>/cluster-api-provider-aws.git
cd cluster-api-provider-aws
git remote add upstream git@github.com:kubernetes-sigs/cluster-api-provider-aws.git
git fetch upstream
```

### Build clusterawsadm

Build `clusterawsadm` in `cluster-api-provider-aws`:

```bash
cd "$(go env GOPATH)"/src/sigs.k8s.io/cluster-api-provider-aws/
make clusterawsadm
sudo mv ./bin/clusterawsadm /usr/local/bin/clusterawsadm
```

### Setup AWS Environment

Create bootstrap file and bootstrap IAM roles and policies using `clusterawsadm`

```bash
$ cat config-bootstrap.yaml

apiVersion: bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1
kind: AWSIAMConfiguration
spec:
  bootstrapUser:
    enable: true

$ clusterawsadm bootstrap iam create-cloudformation-stack
Attempting to create AWS CloudFormation stack cluster-api-provider-aws-sigs-k8s-io
```

#### Customizing the bootstrap permission

The IAM permissions can be customized by using a configuration file with **clusterawsadm**. For example, to create the default IAM role for use with managed machine pools:

```bash
$ cat config-bootstrap.yaml
apiVersion: bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1
kind: AWSIAMConfiguration
spec:
  bootstrapUser:
    enable: true
  eks:
    iamRoleCreation: false # Set to true if you plan to use the EKSEnableIAM feature flag to enable automatic creation of IAM roles
    managedMachinePool:
      disable: false # Set to false to enable creation of the default node role for managed machine pools
```

Use the configuration file to create the additional IAM role:

```bash
$ clusterawsadm bootstrap iam create-cloudformation-stack --config=config-bootstrap.yaml
Attempting to create AWS CloudFormation stack cluster-api-provider-aws-sigs-k8s-io
```

> If you don't plan on using EKS then see the [documentation on disabling EKS support](../topics/eks/disabling.md).

#### Sample Output

When creating the CloudFormation stack using **clusterawsadm** you will see output similar to this:

```bash
Following resources are in the stack:

Resource                  |Type                                                                                |Status
AWS::IAM::Group           |cluster-api-provider-aws-s-AWSIAMGroupBootstrapper-ME9XZVCO2491                     |CREATE_COMPLETE
AWS::IAM::InstanceProfile |control-plane.cluster-api-provider-aws.sigs.k8s.io                                  |CREATE_COMPLETE
AWS::IAM::InstanceProfile |controllers.cluster-api-provider-aws.sigs.k8s.io                                    |CREATE_COMPLETE
AWS::IAM::InstanceProfile |nodes.cluster-api-provider-aws.sigs.k8s.io                                          |CREATE_COMPLETE
AWS::IAM::ManagedPolicy   |arn:aws:iam::xxx:policy/control-plane.cluster-api-provider-aws.sigs.k8s.io |CREATE_COMPLETE
AWS::IAM::ManagedPolicy   |arn:aws:iam::xxx:policy/nodes.cluster-api-provider-aws.sigs.k8s.io         |CREATE_COMPLETE
AWS::IAM::ManagedPolicy   |arn:aws:iam::xxx:policy/controllers.cluster-api-provider-aws.sigs.k8s.io   |CREATE_COMPLETE
AWS::IAM::Role            |control-plane.cluster-api-provider-aws.sigs.k8s.io                                  |CREATE_COMPLETE
AWS::IAM::Role            |controllers.cluster-api-provider-aws.sigs.k8s.io                                    |CREATE_COMPLETE
AWS::IAM::Role            |eks-controlplane.cluster-api-provider-aws.sigs.k8s.io                               |CREATE_COMPLETE
AWS::IAM::Role            |eks-nodegroup.cluster-api-provider-aws.sigs.k8s.io                                  |CREATE_COMPLETE
AWS::IAM::Role            |nodes.cluster-api-provider-aws.sigs.k8s.io                                          |CREATE_COMPLETE
AWS::IAM::User            |bootstrapper.cluster-api-provider-aws.sigs.k8s.io                                   |CREATE_COMPLETE
```

### Set Environment Variables

- Create a security credentials in the `bootstrapper.cluster-api-provider-aws.sigs.k8s.io` IAM user that is created by cloud-formation stack and copy the `AWS_ACCESS_KEY_ID` and `AWS_SECRETS_ACCESS_KEY`.
  (Or use admin user credentials instead)

- Set AWS_B64ENCODED_CREDENTIALS environment variable

   ```bash
   export AWS_ACCESS_KEY_ID=AKIATEST
   export AWS_SECRET_ACCESS_KEY=TESTTEST
   export AWS_REGION=eu-west-1
   export AWS_B64ENCODED_CREDENTIALS=$(clusterawsadm bootstrap credentials encode-as-profile)
   ```

## Running local management cluster for development

Before the next steps, make sure [initial setup for development environment][Initial-setup-for-development-environment] steps are complete.

[Initial-setup-for-development-environment]: development.md#initial-setup-for-development-environment

There are two ways to build aws manager from local cluster-api-provider-aws source and run it in local kind cluster:

### Option 1: Setting up Development Environment with Tilt

[Tilt][tilt] is a tool for quickly building, pushing, and reloading Docker containers as part of a Kubernetes deployment.
Many of the Cluster API engineers use it for quick iteration. Please see our [Tilt instructions][Tilt instructions] to get started.

[tilt]: https://tilt.dev
[Tilt instructions]: ../development/tilt-setup.md

### Option 2: The Old-fashioned way

Running cluster-api and cluster-api-provider-aws controllers in a kind cluster:

1. Create a local kind cluster
   - `kind create cluster`
2. Install core cluster-api controllers (the version must match the cluster-api version in [go.mod][go.mod])
   - `clusterctl init --core cluster-api:v0.3.16 --bootstrap kubeadm:v0.3.16 --control-plane kubeadm:v0.3.16`
3. Build cluster-api-provider-aws docker images
   - `make e2e-image`
4. Release manifests under `./out` directory
   - `RELEASE_TAG="e2e" make release-manifests`
5. Apply the manifests
   - `kubectl apply -f ./out/infrastructure.yaml`

[go]: https://golang.org/doc/install
[jq]: https://stedolan.github.io/jq/download/
[go.mod]: https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/master/go.mod
[kind]: https://sigs.k8s.io/kind
[kustomize]: https://github.com/kubernetes-sigs/kustomize
[kustomizelinux]: https://github.com/kubernetes-sigs/kustomize/blob/master/docs/INSTALL.md
[envsubst]: https://github.com/a8m/envsubst
