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

## Community, discussion, contribution, and support

Learn how to engage with the Kubernetes community on the [community page](http://kubernetes.io/community/).

You can reach the maintainers of this project at:

- [#cluster-api on Kubernetes Slack](http://slack.k8s.io/messages/cluster-api)
- [SIG-Cluster-Lifecycle Mailing List](https://groups.google.com/forum/#!forum/kubernetes-sig-cluster-lifecycle)

### Code of conduct

Participation in the Kubernetes community is governed by the [Kubernetes Code of Conduct](code-of-conduct.md).

