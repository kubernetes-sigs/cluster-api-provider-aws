# Developing Cluster API Provider AWS <!-- omit in toc -->

## Contents <!-- omit in toc -->

<!-- Below is generated using VSCode yzhang.markdown-all-in-one >

<!-- TOC depthFrom:2 -->

- [Setting up](#setting-up)
  - [Base requirements](#base-requirements)
  - [Get the source](#get-the-source)
  - [Get familiar with basic concepts](#get-familiar-with-basic-concepts)
  - [Dev manifest files](#dev-manifest-files)
  - [Dev images](#dev-images)
    - [Container registry](#container-registry)
      - [Setting up a Google Container Registry](#setting-up-a-google-container-registry)
- [Developing](#developing)
  - [Manual Testing](#manual-testing)
    - [Setting up the environment](#setting-up-the-environment)
    - [Building and pushing dev images to GCR](#building-and-pushing-dev-images-to-gcr)
    - [Building and pushing dev images to custom (non GCR) container registry](#building-and-pushing-dev-images-to-custom-non-gcr-container-registry)
    - [Build manifests](#build-manifests)
    - [Running clusterctl](#running-clusterctl)
    - [Executing unit tests](#executing-unit-tests)
    - [Executing integration tests](#executing-integration-tests)
    - [Executing e2e tests](#executing-e2e-tests)
    - [Executing e2e tests with Boskos](#executing-e2e-tests-with-boskos)
  - [Automated Testing](#automated-testing)
    - [Mocks](#mocks)
- [Troubleshooting](#troubleshooting)
  - [`make docker-dev-build` fails](#make-docker-dev-build-fails)

<!-- /TOC -->

## Setting up

### Base requirements

1. Install [go][go]
   - Get the latest patch version for go v1.11.
2. Install [jq][jq]
   - `brew install jq` on MacOS.
3. Install [gettext][gettext] package
   - `brew install gettext && brew link --force gettext` on MacOS.
4. Install [KIND][kind]
   - `go get sigs.k8s.io/kind`.
5. Install [bazel][bazel]
6. Configure Python 2.7+ with [pyenv][pyenv] if your default is Python 3.x.

[go]: https://golang.org/doc/install

### Get the source

`go get -d sigs.k8s.io/cluster-api-provider-aws`

Ensure you have updated the vendor directory with:

``` shell
cd "$(go env GOPATH)/src/sigs.k8s.io/cluster-api-provider-aws"
make dep-ensure
```

### Get familiar with basic concepts

This provider is modeled after the upstream cluster-api project. To get familiar
with resources, concepts and conventions refer to the [upstream gitbook](https://kubernetes-sigs.github.io/cluster-api/).

### Dev manifest files

Part of running cluster-api-provider-aws is generating manifests to run.
Generating dev manifests allows you to test dev images instead of the default
releases.

### Dev images

#### Container registry

Instructions for setting up a Google Container Registry (GCR) are below, but
GCR is not required. The requirement is any public container registry.

AWS Elastic Container Registry (ECR) isn't recommended as public images are
required for the provider manifests.

##### Setting up a Google Container Registry

This project provides a [Terraform module](../hack/terraform-gcr-init/README.md)
that creates a new Google Cloud project suitable for use by Cluster API.

## Developing

Change some code!

### Manual Testing

#### Setting up the environment

Your environment must have your the AWS credentials as outlined in the [getting
started prerequisites section](./getting-started.md#Prerequisites)


#### Building and pushing dev images to GCR

1. [Install the gcloud cli][gcloud_sdk].
2. Set project: `gcloud config set project YOUR_PROJECT_NAME`.
3. To build images with custom tags, run the `make docker-push` as follows:

```(bash)
make docker-push MANAGER_IMAGE_TAG=<YOUR_TAG_HERE>
```

#### Building and pushing dev images to custom (non GCR) container registry

1. Login to your container registry using `docker login`.

   E.g. `docker login quay.io`
2. To build images with custom tags and push to your custom image registry,
   run the `make docker-build` as follows::

```(bash)
make docker-push REGISTRY="your repo"
```

3. Push your docker images as `docker push <ContainerImage>:<YourTag>`

#### Build manifests

Whenever you are working on a branch, you will need to generate manifests
using:

```(bash)
make manifests
```

It's expected that some set of AWS credentials are available at the time, either
as environment variable, a CLI profile or other SDK-supported method.

You will then have a sample cluster and machine manifest in:
`/cmd/clusterctl/examples/aws/out` and a provider components file to use with clusterctl in
`cmd/clusterctl/examples/aws/out/provider-components.yaml`

#### Running clusterctl

`make create-cluster` will launch a bootstrap cluster and then run the generated
manifests creating a target cluster in AWS. After this is finished you will have
a kubeconfig copied locally. You can debug most issues by SSHing into the
instances that have been created and reading `/var/log/cloud-init-output.log`.

#### Executing unit tests

`make test` executes the project's unit tests. These tests do not stand up a
Kubernetes cluster, nor do they have external dependencies.

#### Executing integration tests
`make integration` executes the project's integration tests. These tests stand
up a local Kubernetes cluster using Kind in order to deploy the project's CRDs.
The tested controller is **not** used to deploy Kubernetes to AWS.

These tests depend on the following binaries in the system path:
* `kind`

#### Executing e2e tests
`make e2e` executes the project's end-to-end tests with AWS account
information parsed from the environment.

These tests stand up a local Kubernetes cluster using Kind. The project's CRDs
and controllers are deployed to the Kind cluster and are used to deploy
Kubernetes to AWS.

The AWS janitor is disabled by default. `JANITOR_ENABLED=1 make e2e` executes
janitor immediately after running the e2e tests.

Please keep in mind that the janitor is highly destructive and should not
be executed against shared AWS accounts or preferrably AWS accounts not
dedicated to testing this project.

These tests depend on the following binaries in the system path:
* `kind`
* `aws-janitor`

#### Executing e2e tests with Boskos
`BOSKOS_HOST=http://boskos make e2e` executes the project's end-to-end tests
with AWS account information acquired from a Boskos host.

These tests stand up a local Kubernetes cluster using Kind. The project's CRDs
and controllers are deployed to the Kind cluster and are used to deploy
Kubernetes to AWS.

The AWS janitor is disabled by default.
`BOSKOS_HOST=http://boskos JANITOR_ENABLED=1 make e2e` executes
the janitor immediately after running the e2e tests.

Please keep in mind that the janitor is highly destructive and should not
be executed against shared AWS accounts or preferrably AWS accounts not
dedicated to testing this project.

`BOSKOS_HOST=http://boskos make e2e` executes the project's
end-to-end tests with the janitor disabled.

These tests depend on the following binaries in the system path:
* `kind`
* `aws-janitor`

### Automated Testing

#### Mocks

Mocks are set up using Bazel, see [build](../../build)

If you then want to use these mocks with `go test ./...`, run

`make copy-genmocks`

## Troubleshooting

Troubleshooting steps and workarounds for commonly encountered errors
encountered in doing development work.

### `make docker-dev-build` fails

```(bash)
ERROR: Analysis of target '//cmd/manager:manager-image-dev' failed; build aborted: no such package '@golang-image//image': Pull command failed
```

This is caused by Python2 not being the version of active python in the shell. Bazel
requires [Python 2](https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/master/docs/development.md#base-requirements)
for its docker rules to succeed.

See Issue: [624](https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/624)

<!-- References -->

[jq]: https://stedolan.github.io/jq/download/
[image_pull_secrets]: https://kubernetes.io/docs/concepts/containers/images/#specifying-imagepullsecrets-on-a-pod
[ecr_credential_helper]: https://github.com/awslabs/amazon-ecr-credential-helper
[aws_vault]: https://github.com/99designs/aws-vault
[gcloud_sdk]: https://cloud.google.com/sdk/install
[gettext]: https://www.gnu.org/software/gettext/
[kind]: https://sigs.k8s.io/kind
[aws_cli]: https://docs.aws.amazon.com/cli/latest/userguide/installing.html
[bazel]: https://docs.bazel.build/versions/master/install.html
[pyenv]: https://github.com/pyenv/pyenv
