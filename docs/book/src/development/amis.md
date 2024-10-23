# Publish AMIs

Publishing new AMIs is done via manually invoking a GitHub Actions workflow. 

> NOTE: the plan is to ultimately fully automate the process in the future (see [this issue](https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/1982) for progress).

> NOTE: there are some issues with the RHEL based images at present.

## Get build inputs

For a new Kubernetes version that you want to build an AMI for you will need to determine the following values:

| Input             | Description |
| ----------------- | ----------- |
| kubernetes_semver      | The semver version of k8s you want to build an AMI for. In format vMAJOR.MINOR.PATCH. |
| kubernetes_series      | The release series for the Kubernetes version. In format vMAJOR.MINOR.                |
| kubernetes_deb_version | The version of the debian package for the release.                                    |
| kubernetes_rpm_version | The version of the rpm package for the release                                        |
| kubernetes_cni_semver  | The version of CNI to include. It needs to match the k8s release.                     |
| kubernetes_cni_deb_version | The version of the debian package for the CNI release to use                      |
| crictl_version         | The vesion of the cri-tools package to install into the AMI                           |

You can determine these values directly or by looking at the publish debian apt repositories for the k8s release.

## Build

### Using GitHub Actions Workflow

To build the AMI using GitHub actions you must have write access to the CAPA repository (i.e. be a maintainer or part of release team).

To build the new version:

1. Got to the GitHub Action
2. Click the **Start Workflow** button
3. Fill in the details of the build
4. Click **Run**

### Manually

> **WARNING: the manual process should only be followed in exceptional circumstances.

To build manually you must have admin access to the CNCF AWS account used for the AMIs.

The steps to build manually are:

1. Clone [image-builder](https://github.com/kubernetes-sigs/image-builder)
2. Open a terminal
3. Set the AWS environment variables for the CAPA AMI account
4. Change directory into `images/capi`
5. Create a new file called `vars.json` with the following content (substituing the values with the build inputs):

```json
{
    "kubernetes_rpm_version": "<INSERT_INPUT_VALUE>",
    "kubernetes_semver": "<INSERT_INPUT_VALUE>",
    "kubernetes_series": "<INSERT_INPUT_VALUE>",
    "kubernetes_deb_version": "<INSERT_INPUT_VALUE>",
    "kubernetes_cni_semver": "<INSERT_INPUT_VALUE>",
    "kubernetes_cni_deb_version": "<INSERT_INPUT_VALUE>",
    "crictl_version": "<INSERT_INPUT_VALUE>"
}
```
6. Install dependencies by running:

```shell
make deps-ami
```

7. Build the AMIs using:

```shell
PACKER_VAR_FILES=vars.json make build-ami-ubuntu-2204
PACKER_VAR_FILES=vars.json make build-ami-ubuntu-2404
PACKER_VAR_FILES=vars.json make build-ami-flatcar
PACKER_VAR_FILES=vars.json make build-ami-rhel-8
```
## Additional Information

- The AMIs are hosted in a CNCF owned AWS account (819546954734).
- The AWS resources that are needed to support the GitHub Actions workflow are created via terraform. Source is [here](https://github.com/kubernetes/k8s.io/tree/main/infra/aws/terraform/cncf-k8s-infra-aws-capa-ami).
- OIDC and IAM Roles are used to grant access via short lived credentials to the GitHub Action workflow instance when it runs.

