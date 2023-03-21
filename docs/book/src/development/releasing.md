# Release Process Guide

**Important:** Before you start, make sure all [periodic tests](https://testgrid.k8s.io/sig-cluster-lifecycle-cluster-api-provider-aws) are passing on the most recent commit that will be included in the release. Check for consistency by scrolling to the right to view older test runs.
    Examples:
    - <https://testgrid.k8s.io/sig-cluster-lifecycle-cluster-api-provider-aws-1.5#periodic-e2e-release-1-5>
    - <https://testgrid.k8s.io/sig-cluster-lifecycle-cluster-api-provider-aws-1.5#periodic-eks-e2e-release-1-5>

## Create tag, and build staging container images

1. Create a new local repository of <https://github.com/kubernetes-sigs/cluster-api-provider-aws> (e.g. using `git clone`).
1. If this is a major or minor release, create a new release branch and push to GitHub, otherwise switch to it, e.g. `git checkout release-1.5`.
1. If this is a major or minor release, update `metadata.yaml` and make a commit.
1. Update the release branch on the repository, e.g. `git push origin HEAD:release-1.5`.
1. Make sure your repo is clean by git standards.
1. Set environment variable `GITHUB_TOKEN` to a GitHub personal access token. The token must have write access to the `kubernetes-sigs/cluster-api-provider-aws` repository.
1. Set environment variables `PREVIOUS_VERSION` which is the last release tag and `VERSION` which is the current release version, e.g. `export PREVIOUS_VERSION=v1.4.0 VERSION=v1.5.0`, or `export PREVIOUS_VERSION=v1.5.0 VERSION=v1.5.1`).
1. Create a tag `git tag -s -m $VERSION $VERSION`. `-s` flag is for GNU Privacy Guard (GPG) signing.
1. Make sure you have push permissions to the upstream CAPA repo. Push tag you've just created (`git push <upstream-repo-remote> $VERSION`).
1. A prow job will start running to push images to the staging repo, can be seen [here](https://testgrid.k8s.io/sig-cluster-lifecycle-image-pushes#post-cluster-api-provider-aws-push-images). The job is called "post-cluster-api-provider-aws-push-images," and is defined in <https://github.com/kubernetes/test-infra/blob/master/config/jobs/image-pushing/k8s-staging-cluster-api.yaml>.
1. When the job is finished, wait for the images to be created: `docker pull gcr.io/k8s-staging-cluster-api-aws/cluster-api-aws-controller:$VERSION`. You can also wrap this with a command to retry periodically, until the job is complete, e.g. `watch --interval 30 --chgexit docker pull <...>`.

## Promote container images from staging to production

Promote the container images from the staging registry to the production registry (`registry.k8s.io/cluster-api-provider-aws`) by following the steps below.

1. Navigate to the the staging repository [dashboard](https://console.cloud.google.com/gcr/images/k8s-staging-cluster-api-aws/GLOBAL).
2. Choose the _top level_ [cluster-api-aws-controller](https://console.cloud.google.com/gcr/images/k8s-staging-cluster-api-aws/GLOBAL/cluster-api-aws-controller?gcrImageListsize=30) image. Only the top level image provides the multi-arch manifest, rather than one for a specific architecture.
3. Wait for an image to appear with the tagged release version.
4. If you don't have a GitHub token, create one by going to your GitHub settings, in [Personal access tokens](https://github.com/settings/tokens). Make sure you give the token the `repo` scope.
5. Create a PR to promote the images:
    ```bash
    export GITHUB_TOKEN=<your GH token>
    make promote-images
    ```
    **Notes**:
     * `kpromo` uses `git@github.com:...` as remote to push the branch for the PR. If you don't have `ssh` set up you can configure
       git to use `https` instead via `git config --global url."https://github.com/".insteadOf git@github.com:`.
     * This will automatically create a PR in [k8s.io](https://github.com/kubernetes/k8s.io) and assign the CAPA maintainers.
6. Wait for the PR to be approved (typically by CAPA maintainers authorized to merge PRs into the k8s.io repository) and merged.
7. Verify the images are available in the production registry:
    ```bash
    docker pull registry.k8s.io/cluster-api-provider-aws/cluster-api-aws-controller:${RELEASE_TAG}
    ```

## Create release artifacts, and a GitHub draft release

1. Again, make sure your repo is clean by git standards.
1. Export the current branch `export BRANCH=release-1.5` (`export BRANCH=main`)and run `make release`.
1. Run `make create-gh-release` to create a draft release on Github, copying the generated release notes from `out/CHANGELOG.md` into the draft.
1. Run `make upload-gh-artifacts` to upload artifacts from .out/ directory. You may run into API limit errors, so verify artifacts at next step.
1. Verify that all the files below are attached to the drafted release:
    1. `clusterawsadm-darwin-amd64`
    1. `clusterawsadm-linux-amd64`
    1. `infrastructure-components.yaml`
    1. `cluster-template.yaml`
    1. `cluster-template-machinepool.yaml`
    1. `cluster-template-eks.yaml`
    1. `cluster-template-eks-managedmachinepool.yaml`
    1. `cluster-template-eks-managedmachinepool-vpccni.yaml`
    1. `cluster-template-eks-managedmachinepool-gpu.yaml`
    1. `eks-controlplane-components.yaml`
    1. `eks-bootstrap-components.yaml`
    1. `metadata.yaml`
1. Finalise the release notes by editing the draft release.

## Publish the draft release

1. Make sure image promotion is complete before publishing the release draft. The promotion job logs can be found [here](https://testgrid.k8s.io/sig-k8s-infra-k8sio#post-k8sio-image-promo) and you can also try and pull the images (i.e. ``docker pull registry.k8s.io/cluster-api-aws/cluster-api-aws-controller:v0.6.4`).
1. Publish release. Use the pre-release option for release
     candidate versions of Cluster API Provider AWS.
1. Email `kubernetes-sig-cluster-lifecycle@googlegroups.com` to announce the release.
