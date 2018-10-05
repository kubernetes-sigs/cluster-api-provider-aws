# Developing Cluster API Provider AWS <!-- omit in toc -->

## Contents <!-- omit in toc -->

<!-- Below is generated using VSCode yzhang.markdown-all-in-one >

<!-- TOC depthFrom:2 -->
- [Updating mocks](#updating-mocks)
- [Set up](#set-up)
  - [Base requirements](#base-requirements)
  - [Using Google Cloud](#using-google-cloud)
    - [Using images on Google Cloud](#using-images-on-google-cloud)
  - [Using AWS Elastic Container Registry](#using-aws-elastic-container-registry)
    - [Using images on Elastic Container Registry](#using-images-on-elastic-container-registry)

<!-- /TOC -->

## Updating mocks

When you update the mocks, please make sure the imports are not from the dependencies vendor directory. Instead, make them an explicit dependency of this project.

For example, if you see `types "sigs.k8s.io/cluster-api/vendor/k8s.io/apimachinery/pkg/types"` in the import path, replace it with `types "k8s.io/apimachinery/pkg/types"`.

## Set up

You can use Google Cloud or AWS ECR to host repositories, or just use local minikube
storage.

### Base requirements

1. Install [jq][jq]
   - `brew install jq` on MacOS.
1. Install [gettext][gettext] package
   - `brew install gettext && brew link --force gettext` on MacOS.
1. Install [minikube][minikube]
   - `brew install minikube` on MacOS.
1. Configure `minikube`:
    1. Use Kubernetes v1.9.4 `minikube config set kubernetes-version v1.9.4`.
    1. Use kubeadm as bootstrapper `minikube config set bootstrapper kubeadm`.

### Using Google Cloud

To create a new project suitable for use by Cluster API, there is a
[Terraform module](../hack/terraform-gcr-init/README.md) you can use to bootstrap
Google Container Registry.

1. [Install the gcloud cli][gcloud_sdk].
1. Set project: `gcloud config set project YOUR_PROJECT_NAME`.
1. Pushing dev images: `make dev_push`.

#### Using images on Google Cloud

``` bash
export CLUSTER_CONTROLLER_IMAGE=gcr.io/$(gcloud config get-value project)/aws-cluster-controller:0.0.1-dev
export MACHINE_CONTROLLER_IMAGE=gcr.io/$(gcloud config get-value project)/aws-machine-controller:0.0.1-dev
```

Then generate the [example configuration](../README.md#running-clusterctl) as normal.

### Using AWS Elastic Container Registry

1. Install the [AWS CLI][aws_cli].
1. Set `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` and `AWS_DEFAULT_REGION` variables
   - you may want to use [`aws-vault exec <profile> bash`][aws_vault]
     to store credentials or assume role in MFA scenarios.
1. Create two ECR repositories.

    ``` bash
    aws ecr create-repository --repository-name aws-machine-controller
    aws ecr create-repository --repository-name aws-cluster-controller
    ```
1. Run `eval $(aws ecr get-login --no-include-email)` or use the [ECR Credential Helper[ecr_credential_helper] to set up Docker
1. Push dev images with `DEV_REPO_TYPE=ECR make dev_push`

#### Using images on Elastic Container Registry

``` bash
export AWS_ACCOUNT_ID=$(aws sts get-caller-identity | jq -r .Account)
export CLUSTER_CONTROLLER_IMAGE=${AWS_ACCOUNT_ID}.dkr.ecr.${DEFAULT_AWS_REGION}.amazonaws.com/aws-cluster-controller:0.0.1-dev
export MACHINE_CONTROLLER_IMAGE=${AWS_ACCOUNT_ID}.dkr.ecr.${DEFAULT_AWS_REGION}.amazonaws.com/aws-machine-controller:0.0.1-dev
```

Then generate the [example configuration](getting-started.md#generating-cluster-manifests) as normal.

You will also need to configure minikube or your bootstrap cluster with credentials to access ECR, using [image pull secrets][image_pull_secrets]
or another mechanism.

<!-- References -->

[jq]: https://stedolan.github.io/jq/download/
[image_pull_secrets]: https://kubernetes.io/docs/concepts/containers/images/#specifying-imagepullsecrets-on-a-pod
[ecr_credential_helper]: https://github.com/awslabs/amazon-ecr-credential-helper
[aws_vault]: https://github.com/99designs/aws-vault
[gcloud_sdk]: https://cloud.google.com/sdk/install
[gettext]: https://www.gnu.org/software/gettext/
[minikube]: https://kubernetes.io/docs/setup/minikube/
[aws_cli]: https://docs.aws.amazon.com/cli/latest/userguide/installing.html
