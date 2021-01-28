#  Setting up Development Environment for Cluster API Provider AWS 

In this post we will be deep diving into how to setup and contribute into the ClusterAPI provider AWS.

## Setup Development Environment for EKS Control Plane

login to github and fork both

- https://github.com/kubernetes-sigs/cluster-api
- https://github.com/kubernetes-sigs/cluster-api-provider-aws

### Install Prerequisites

- Golang Version 1.13 or higher
- direnv
- envsubst
- kubectl
- Working Development Environment

install `direnv`

```bash
brew install direnv
```

install `envsubst`

```bash
curl -L https://github.com/a8m/envsubst/releases/download/v1.2.0/envsubst-`uname -s`-`uname -m` -o envsubst
chmod +x envsubst
sudo mv envsubst $HOME/go/bin/

# or use go get
export GOPATH=~go
go get -v github.com/a8m/envsubst/cmd/envsubst
```

### Setup Repos and GOPATH

setup this in `~/go/src/sigs.k8s.io` and `export GOPATH=~go`

```bash
$ export GOPATH=<HOME ABS PATH>/go

$ mkdir ~/go/src/sigs.k8s.io/

$ cd ~/go/src/sigs.k8s.io/

$ git clone git@github.com:<GITHUB USERNAME>/cluster-api.git

$ git clone git@github.com:<GITHUB USERNAME>/cluster-api-provider-aws.git
```

and add upstream for both repos

```bash
$ cd cluster-api 

$ git remote add upstream git@github.com:kubernetes-sigs/cluster-api.git

$ git fetch upstream

$ cd ..

$ cd cluster-api-provider-aws

$ git remote add upstream git@github.com:kubernetes-sigs/cluster-api.git

$ git fetch upstream
```

### Build clusterawsadm and setup your AWS Environment

build the `clusterawsadm` in `cluster-api-provider-aws`

```bash
$ make clusterawsadm
```

create bootstrap file and bootstrap IAM roles and policies using `clusterawsadm`

```bash
$ cat config-bootstrap.yaml
apiVersion: bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1
kind: AWSIAMConfiguration
spec:
  bootstrapUser:
    enable: true
  eks:
    enable: true
    iamRoleCreation: false # Set to true if you plan to use the EKSEnableIAM feature flag to enable automatic creation of IAM roles
    defaultControlPlaneRole:
      disable: false # Set to false to enable creation of the default control plane role
    managedMachinePool:
      disable: false # Set to false to enable creation of the default node role for managed machine pools
```

create IAM Resources that will be needed for bootstrapping EKS 

```bash
$ ./bin/clusterawsadm bootstrap iam create-cloudformation-stack --config=config-bootstrap.yaml
Attempting to create AWS CloudFormation stack cluster-api-provider-aws-sigs-k8s-io
```

this will create cloudformation stack for those IAM resources

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

create a security credentials in the `bootstrapper.cluster-api-provider-aws.sigs.k8s.io` IAM user and copy the `AWS_ACCESS_KEY_ID` and `AWS_SECRETS_ACCESS_KEY`

```bash
$ brew install direnv
$ touch .envrc #add the aws key

unset AWS_SESSION_TOKEN
unset AWS_SECURITY_TOKEN
export AWS_ACCESS_KEY_ID=AKIATEST
export AWS_SECRET_ACCESS_KEY=TESTTEST
export AWS_REGION=eu-west-1
```

then run `direnv allow` for each change done in `.envrc`

create a kind cluster where we can deploy the kubernetes manifests. This kind cluster will be a temp cluster which we will use to create the EKS cluster. 

```bash
$ kind create cluster
```

setup cluster-api with the same `.envrc` file and then allow direnv 

```bash
cp .envrc ../cluster-api
cd ../cluster-api
direnv allow

```

create `tilt-settings.json` change the value of `AWS_B64ENCODED_CREDENTIALS`

```bash
{
"default_registry": "gcr.io/<GITHUB USERNAME>",
    "provider_repos": ["../cluster-api-provider-aws"],
    "enable_providers": ["eks-bootstrap", "eks-controlplane", "kubeadm-bootstrap", "kubeadm-control-plane", "aws"],
    "kustomize_substitutions": {
        "AWS_B64ENCODED_CREDENTIALS": "W2RlZmFZSZnRg==",
        "EXP_EKS": "true",
        "EXP_EKS_IAM": "true",
        "EXP_MACHINE_POOL": "true"
    },
    "extra_args": {
        "aws": ["--v=2"],
        "eks-bootstrap": ["--v=2"],
        "eks-controlplane": ["--v=2"]
    }
  }
```

run dev env using `tilt` let it run and press space once its running to open the web browser

```bash
tilt up
```

create`cd cluster-api-provider-aws`  and edit `.envrc`

```bash
export AWS_EKS_ROLE_ARN=arn:aws:iam::<accountid>:role/aws-service-role/eks.amazonaws.com/AWSServiceRoleForAmazonEKS
export AWS_SSH_KEY_NAME=<sshkeypair>
export KUBERNETES_VERSION=v1.15.2
export EKS_KUBERNETES_VERSION=v1.15
export CLUSTER_NAME=capi-<test-clustename>
export CONTROL_PLANE_MACHINE_COUNT=1
export AWS_CONTROL_PLANE_MACHINE_TYPE=t3.large
export WORKER_MACHINE_COUNT=1
export AWS_NODE_MACHINE_TYPE=t3.large
```

check and pipe output of template into a file

```bash
cat templates/cluster-template-eks.yaml 
cat templates/cluster-template-eks.yaml | $HOME/go/bin/envsubst > test-cluster.yaml
```

apply generate `test-cluster.yaml` file in the `kind cluster`

```bash
kubectx
kubectx kind-kind
kubectl apply -f test-cluster.yaml
```

Check the tilt logs and wait for the EKS Cluster to be created 

## Retry if theres an error when creating the test-cluster

To retry apis and services again delete the cluster and recreate it again 

```bash
tilt up (ctrl-c)
press space to see the logs

kubectl delete -f test-cluster.yaml

kind delete cluster

kind create cluster
```

try again 

```bash
tilt up
```

### Clean up

To clean make sure you delete the Kubernetes Resources first before deleting the Kind Cluster

### Troubleshooting

- Make sure you have at least three available spaces EIP and NAT Gateways to be created
- If your git starts throwing this error

```bash
flag provided but not defined: -variables
Usage: envsubst [options...] <input>
```

you might need to reinstall the system `envsubst`

```bash
brew install gettetxt
# or
brew reinstall gettext
```

Make sure you specify which `envsubst` you are using
