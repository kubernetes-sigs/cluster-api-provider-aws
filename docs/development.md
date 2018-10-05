Developing Cluster API Provider AWS
===================================

## Updating mocks

When you update the mocks, please make sure the imports are not from the dependencies vendor directory. Instead, make them an explicit dependency of this project.

For example, if you see `types "sigs.k8s.io/cluster-api/vendor/k8s.io/apimachinery/pkg/types"` in the import path, replace it with `types "k8s.io/apimachinery/pkg/types"`.

## Set up

You can use Google Cloud or AWS ECR to host repositories, or just use local minikube
storage.

### Base requirements

1. Install [jq](https://stedolan.github.io/jq/download/) (`brew install jq` on Mac).
1. Install [gettext](https://www.gnu.org/software/gettext/) package (`brew install gettext && brew link --force gettext` on Mac).
1. Install [minikube](https://kubernetes.io/docs/setup/minikube/) (`brew install minikube` on Mac).
1. Configure `minikube`:
	- Use Kubernetes v1.9.4 `minikube config set kubernetes-version v1.9.4`.
	- Use kubeadm as bootstrapper `minikube config set bootstrapper kubeadm`.

### Using Google Cloud

To create a new project suitable for use by Cluster API, there is a
[Terraform module](../hack/terraform-gcr-init/README.md) you can use to bootstrap ECR.

1. [Install the gcloud cli](https://cloud.google.com/sdk/install).
1. Set project: `gcloud config set project YOUR_PROJECT_NAME`.
. Pushing dev images: `make dev_push`.

#### Using images on Google Cloud
``` shell
export CLUSTER_CONTROLLER_IMAGE=gcr.io/$(gcloud config get-value project)/aws-cluster-controller:0.0.1-dev
export MACHINE_CONTROLLER_IMAGE=gcr.io/$(gcloud config get-value project)/aws-machine-controller:0.0.1-dev
```

Then generate the [example configuration](../README.md#running-clusterctl) as normal.

### Using AWS Elastic Container Registry
1. Install the [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/installing.html).
2. Set `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` and `AWS_DEFAULT_REGION` variables
   * you may want to use [`aws-vault exec <profile> bash`](https://github.com/99designs/aws-vault)
     to store credentials or assume role in MFA scenarios.
3. Create two ECR repositories.

  ``` shell
  aws ecr create-repository --repository-name aws-machine-controller
  aws ecr create-repository --repository-name aws-cluster-controller
  ```

4. Run `eval $(aws ecr get-login --no-include-email)` or use the [ECR Credential Helper](https://github.com/awslabs/amazon-ecr-credential-helper) to set up Docker
5. Push dev images with `DEV_REPO_TYPE=ECR make dev_push`

#### Using images on Elastic Container Registry
``` shell
export AWS_ACCOUNT_ID=$(aws sts get-caller-identity | jq -r .Account)
export CLUSTER_CONTROLLER_IMAGE=${AWS_ACCOUNT_ID}.dkr.ecr.${DEFAULT_AWS_REGION}.amazonaws.com/aws-cluster-controller:0.0.1-dev
export MACHINE_CONTROLLER_IMAGE=${AWS_ACCOUNT_ID}.dkr.ecr.${DEFAULT_AWS_REGION}.amazonaws.com/aws-machine-controller:0.0.1-dev
```

Then generate the [example configuration](../README.md#running-clusterctl) as normal.

You will also need to configure minikube or your bootstrap cluster with credentials to access ECR, using [image pull secrets](https://kubernetes.io/docs/concepts/containers/images/#specifying-imagepullsecrets-on-a-pod)
or another mechanism.

### Using minikube's local Docker storage

If you are using the default minikube as the bootstrap cluster, you can also use Minikube's local storage and not upload
Docker images to an external repository

``` shell
make minikube_build
```

## Convenience Makefile targets for development

### `make test_cluster`

Runs `clusterctl`, applying the generator example configuration.


### `make tail_cluster_actuator_logs`

Tails cluster actuator logs

### `make tail_machine_actuator_logs`

Tails machine actuator logs

### `make test_cluster_deletion`

Deletes the cluster object, testing the cluster actuator's deletion pathway.

### `make destroy_test_cluster`

Forcibly tears down the objects from the example configuration so you can reuse
the same bootstrap cluster. Will not exercise deletion on the actuators,
so you may need to manually clean up AWS objects if needed.

### `make minikube_set`

Sets necessary parameters if running your own minikube as an external
bootstrap cluster.

### `make minikube_unset`

Resets the Kubernetes version for minikube
