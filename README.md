Kubernetes cluster-api-provider-aws Project
===========================================

This repository hosts an implementation of a provider for AWS for the [cluster-api project](https://sigs.k8s.io/cluster-api).

Note: This repository is currently a skeleton implementation of a cluster-api provider, implementation will begin once there is agreement on the [Design Spec](https://docs.google.com/document/d/1G7DRQccoTY5YBrinQb6sz_fRLB9zFbCnI1O984XFk7Q).

## Development

### Updating mocks

When you update the mocks, please make sure the imports are not from the dependencies vendor directory. Instead, make them an explicit dependency of this project.

For example, if you see `types "sigs.k8s.io/cluster-api/vendor/k8s.io/apimachinery/pkg/types"` in the import path, replace it with `types "k8s.io/apimachinery/pkg/types"`.

### Requirements

* A Google Cloud project or 2 AWS Elastic Container Registry repositories

### Set up

You can use Google Cloud or AWS ECR to host repositories.

#### Pushing development images to Google Cloud

1. [Install the gcloud cli](https://cloud.google.com/sdk/install)
2. Set project: `gcloud config set project YOUR_PROJECT_NAME`
3. Pushing dev images: `make dev_push`

#### Pushing development images to AWS Elastic Container Registry
1. Install the [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/installing.html)
2. Install [jq](https://stedolan.github.io/jq/download/)
3. Set `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` and `AWS_DEFAULT_REGION` variables
   * you may want to use [`aws-vault exec <profile> bash`](https://github.com/99designs/aws-vault)
     to store credentials or assume role in MFA scenarios.
4. Create two ECR repositories

  ``` shell
  aws ecr create-repository --repository-name aws-machine-controller
  aws ecr create-repository --repository-name aws-cluster-controller
  ```

5. Run `eval $(aws ecr get-login --no-include-email)` or use the [ECR Credential Helper](https://github.com/awslabs/amazon-ecr-credential-helper) to set up Docker
6. Push dev images with `DEV_REPO_TYPE=ECR make dev_push`

### clusterctl

#### Using images on Google Cloud
``` shell
export CLUSTER_CONTROLLER_IMAGE=gcr.io/$(gcloud config get-value project)/aws-cluster-controller:0.0.1-dev
export MACHINE_CONTROLLER_IMAGE=gcr.io/$(gcloud config get-value project)/aws-machine-controller:0.0.1-dev
```

#### Using images on Elastic Container Registry
``` shell
export AWS_ACCOUNT_ID=$(aws sts get-caller-identity | jq -r .Account)
export CLUSTER_CONTROLLER_IMAGE=${AWS_ACCOUNT_ID}.dkr.ecr.${DEFAULT_AWS_REGION}.amazonaws.com/aws-cluster-controller:0.0.1-dev
export MACHINE_CONTROLLER_IMAGE=${AWS_ACCOUNT_ID}.dkr.ecr.${DEFAULT_AWS_REGION}.amazonaws.com/aws-machine-controller:0.0.1-dev
```

You will also need to configure minikube or your bootstrap cluster with credentials to access ECR, using [image pull secrets](https://kubernetes.io/docs/concepts/containers/images/#specifying-imagepullsecrets-on-a-pod)
or another mechanism.

#### Running clusterctl
3. Generate environment file with `make envfile` and fill out the desired variables.
3. Generate the input files with `make example`.
4. Instantiate the cluster with:
``` shell
clusterctl create cluster -v 2 -c clusterctl/examples/aws/out/cluster.yaml -m clusterctl/examples/aws/out/machines.yaml -p clusterctl/examples/aws/out/provider-components.yaml --provider aws
```

At this point clusterctl will spin up a local minikube and create a new cluster, the output should look like:
```bash
Starting VM...
Getting VM IP address...
Moving files into cluster...
Setting up certs...
Connecting to cluster...
Setting up kubeconfig...
Starting cluster components...
Kubectl is now configured to use the cluster.
Loading cached images from config file.
I0926 14:05:25.555809   15312 clusterdeployer.go:112] Applying Cluster API stack to bootstrap cluster
I0926 14:05:25.555827   15312 clusterdeployer.go:301] Applying Cluster API APIServer
I0926 14:05:25.720092   15312 clusterclient.go:511] Waiting for kubectl apply...
I0926 14:05:25.971457   15312 clusterclient.go:539] Waiting for Cluster v1alpha resources to become available...
I0926 14:05:55.978996   15312 clusterclient.go:552] Waiting for Cluster v1alpha resources to be listable...
I0926 14:05:56.024979   15312 clusterdeployer.go:307] Applying Cluster API Provider Components
I0926 14:05:56.025001   15312 clusterclient.go:511] Waiting for kubectl apply...
I0926 14:05:56.168724   15312 clusterdeployer.go:117] Provisioning target cluster via bootstrap cluster
I0926 14:05:56.168752   15312 clusterdeployer.go:119] Creating cluster object test1 on bootstrap cluster in namespace "default"
I0926 14:05:56.173557   15312 clusterdeployer.go:124] Creating master  in namespace "default"
I0926 14:05:56.183134   15312 clusterclient.go:563] Waiting for Machine aws-controlplane-65vtk to become ready...
```

#### Pre-built Kubernetes AMIs for testing

| Kubernetes Version | OS | Region | AMI ID |
|--------------------|----|--------|--------|
| v1.12.0 | CentOS 7 | ap-northeast-1: ami-0d99fdd8448ccc341 |
| | | ap-northeast-2: ami-0e54afe205a335cc1 |
| | | ap-south-1: ami-05f18fed80df3ffea |
| | | ap-southeast-1: ami-040ab1574387cf664 |
| | | ap-southeast-2: ami-04b4aa021f4125c7b |
| | | ca-central-1: ami-06f5aede29daa4ef1 |
| | | eu-central-1: ami-03e8c6bff63cf0c56 |
| | | eu-west-1: ami-0a3c2624667e4ca93 |
| | | eu-west-2: ami-028ea3af39a4a7661 |
| | | eu-west-3: ami-02ee45f96ce686f98 |
| | | sa-east-1: ami-0b4b29cd52371486d |
| | | us-east-1: ami-04ab2c6566a35dfe1 |
| | | us-east-2: ami-02c853738d42d0cc9 |
| | | us-west-1: ami-0b1994ea6e6f57f4b |
| | | us-west-2: ami-018e347dfe0c44534 |
| | Amazon Linux 2 | ap-northeast-1: ami-03e497af31dfd38c9 |
| | | ap-northeast-2: ami-01eb44929bedf4708 |
| | | ap-south-1: ami-06aea13984317ba0b |
| | | ap-southeast-1: ami-09fd6b717ce611d23 |
| | | ap-southeast-2: ami-0ad1751f1d49320e2 |
| | | ca-central-1: ami-09a82ebea7b4c5194 |
| | | eu-central-1: ami-0bbc9065e95f06ec0 |
| | | eu-west-1: ami-0d929d41d1f580dd5 |
| | | eu-west-2: ami-0cbb336c0ee02142b |
| | | eu-west-3: ami-02daf41cab4e66639 |
| | | sa-east-1: ami-08434c76fb33b7dea |
| | | us-east-1: ami-06d15e7a5aa71bd81 |
| | | us-east-2: ami-002c5b8ed125b204d |
| | | us-west-1: ami-0f14c305334254b09 |
| | | us-west-2: ami-0cbf198692b35cc70 |
| | Ubuntu Bionic | ap-northeast-1: ami-01b9b95604fddade1 |
| | | ap-northeast-2: ami-0429d18e76a0b3705 |
| | | ap-south-1: ami-09930fca077b07e18 |
| | | ap-southeast-1: ami-0b2e5665eda719758 |
| | | ap-southeast-2: ami-01c2df5ce9bc2f573 |
| | | ca-central-1: ami-0f5fbb71a0c65e000 |
| | | eu-central-1: ami-0fad7824ed21125b1 |
| | | eu-west-1: ami-0da760e590e7de0e8 |
| | | eu-west-2: ami-04137690b33cb5a8e |
| | | eu-west-3: ami-028272d8c8e9ff369 |
| | | sa-east-1: ami-041a08e5511d25535 |
| | | us-east-1: ami-0de61b6929e8f091c |
| | | us-east-2: ami-0a2463ac1e1f46b95 |
| | | us-west-1: ami-05dc1567db5bd869a |
| | | us-west-2: ami-0f33a1d90f189e0a1 |

## Community, discussion, contribution, and support

Learn how to engage with the Kubernetes community on the [community page](http://kubernetes.io/community/).

You can reach the maintainers of this project at:

- [#cluster-api on Kubernetes Slack](http://slack.k8s.io/messages/cluster-api)
- [SIG-Cluster-Lifecycle Mailing List](https://groups.google.com/forum/#!forum/kubernetes-sig-cluster-lifecycle)

### Code of conduct

Participation in the Kubernetes community is governed by the [Kubernetes Code of Conduct](code-of-conduct.md).

