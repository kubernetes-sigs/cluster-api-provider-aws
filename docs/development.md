# Developing Cluster API Provider AWS <!-- omit in toc -->

## Contents <!-- omit in toc -->

<!-- Below is generated using VSCode yzhang.markdown-all-in-one >

<!-- TOC depthFrom:2 -->
- [Committing code](#committing-code)
- [Mocks](#mocks)
- [Set up](#set-up)
  - [Base requirements](#base-requirements)
  - [Using Google Cloud](#using-google-cloud)
    - [Using images on Google Cloud](#using-images-on-google-cloud)
- [cluster-api-dev-helper](#cluster-api-dev-helper)

<!-- /TOC -->

## Committing code

Ensure you have updated the vendor directory with:

``` shell
make dep-ensure
```

## Mocks

Mocks are set up using Bazel, see [build](../../build)

If you then want to use these mocks with `go test ./...`, run
`make copy-genmocks`

## Set up

You can use Google Cloud or just use local minikube storage.

AWS Elastic Container Registry isn't recommended as public images are
required for the provider manifests.

### Base requirements

1. Install [jq][jq]
   - `brew install jq` on MacOS.
2. Install [gettext][gettext] package
   - `brew install gettext && brew link --force gettext` on MacOS.
3. Install [minikube][minikube]
   - `brew install minikube` on MacOS.
4. Install [bazel][bazel]
5. Configure Python 2.7+ with [pyenv][pyenv] if your default is Python 3.x.
6. Configure `minikube`:
    1. Use Kubernetes v1.12.1 `minikube config set kubernetes-version v1.12.1`.
    2. Use kubeadm as bootstrapper `minikube config set bootstrapper kubeadm`.

### Using Google Cloud

To create a new project suitable for use by Cluster API, there is a
[Terraform module](../hack/terraform-gcr-init/README.md) you can use to bootstrap
Google Container Registry.

1. [Install the gcloud cli][gcloud_sdk].
1. Set project: `gcloud config set project YOUR_PROJECT_NAME`.
1. Pushing dev images: `make docker-push-dev`.

#### Using images on Google Cloud

``` bash
export DEV_DOCKER_REPO=gcr.io/$(gcloud config get-value project)
```

Then generate the [example configuration](../README.md#running-clusterctl) as normal.

## cluster-api-dev-helper

Some command development tasks have been put into a cluster-api-dev-helper
utility in the /hack directory.

To build it, run:

``` bash
make cluster-api-dev-helper
```


<!-- References -->

[jq]: https://stedolan.github.io/jq/download/
[image_pull_secrets]: https://kubernetes.io/docs/concepts/containers/images/#specifying-imagepullsecrets-on-a-pod
[ecr_credential_helper]: https://github.com/awslabs/amazon-ecr-credential-helper
[aws_vault]: https://github.com/99designs/aws-vault
[gcloud_sdk]: https://cloud.google.com/sdk/install
[gettext]: https://www.gnu.org/software/gettext/
[minikube]: https://kubernetes.io/docs/setup/minikube/
[aws_cli]: https://docs.aws.amazon.com/cli/latest/userguide/installing.html
[bazel]: https://docs.bazel.build/versions/master/install.html
[pyenv]: https://github.com/pyenv/pyenv
