# Release Process Guide

**Important:** Before you start, make sure all [periodic tests](https://testgrid.k8s.io/sig-cluster-lifecycle-cluster-api-provider-aws) are passing on the most recent commit that will be included in the release. Check for consistency by scrolling to the right to view older test runs.
    Examples:
    - <https://testgrid.k8s.io/sig-cluster-lifecycle-cluster-api-provider-aws-1.5#periodic-e2e-release-1-5>
    - <https://testgrid.k8s.io/sig-cluster-lifecycle-cluster-api-provider-aws-1.5#periodic-eks-e2e-release-1-5>

## Create tag, and build staging container images

1. Please fork <https://github.com/kubernetes-sigs/cluster-api-provider-aws> and clone your own repository with e.g. `git clone git@github.com:YourGitHubUsername/cluster-api-provider-aws.git`. `kpromo` uses the fork to build images from.
1. Add a git remote to the upstream project. `git remote add upstream git@github.com:kubernetes-sigs/cluster-api-provider-aws.git`
1. If this is a major or minor release, create a new release branch and push to GitHub, otherwise switch to it, e.g. `git checkout release-1.5`.
1. If this is a major or minor release, update `metadata.yaml` by adding a new section with the version, and make a commit.
1. Update the release branch on the repository, e.g. `git push origin HEAD:release-1.5`. `origin` refers to the remote git reference to your fork.
1. Update the release branch on the repository, e.g. `git push upstream HEAD:release-1.5`. `upstream` refers to the upstream git reference.
1. Make sure your repo is clean by git standards.
1. Set environment variables which is the last release tag and `VERSION` which is the current release version, e.g. `export VERSION=v1.5.0`, or `export VERSION=v1.5.1`).
    _**Note**_: the version MUST contain a `v` in front.
    _**Note**_: you must have a gpg signing configured with git and registered with GitHub.

1. Create a tag `git tag -s -m $VERSION $VERSION`. `-s` flag is for GNU Privacy Guard (GPG) signing.
1. Make sure you have push permissions to the upstream CAPA repo. Push tag you've just created (`git push <upstream-repo-remote> $VERSION`). Pushing this tag will kick off a GitHub Action that will create the release and attach the binaries and YAML templates to it.
1. A prow job will start running to push images to the staging repo, can be seen [here](https://testgrid.k8s.io/sig-cluster-lifecycle-image-pushes#post-cluster-api-provider-aws-push-images). The job is called "post-cluster-api-provider-aws-push-images," and is defined in <https://github.com/kubernetes/test-infra/blob/master/config/jobs/image-pushing/k8s-staging-cluster-api.yaml>.
1. When the job is finished, wait for the images to be created: `docker pull gcr.io/k8s-staging-cluster-api-aws/cluster-api-aws-controller:$VERSION`. You can also wrap this with a command to retry periodically, until the job is complete, e.g. `watch --interval 30 --chgexit docker pull <...>`.

## Promote container images from staging to production

Promote the container images from the staging registry to the production registry (`registry.k8s.io/cluster-api-provider-aws`) by following the steps below.

1. Navigate to the staging repository [dashboard](https://console.cloud.google.com/gcr/images/k8s-staging-cluster-api-aws/GLOBAL).
2. Choose the _top level_ [cluster-api-aws-controller](https://console.cloud.google.com/gcr/images/k8s-staging-cluster-api-aws/GLOBAL/cluster-api-aws-controller?gcrImageListsize=30) image. Only the top level image provides the multi-arch manifest, rather than one for a specific architecture.
3. Wait for an image to appear with the tagged release version.
4. If you don't have a GitHub token, create one by going to your GitHub settings, in [Personal access tokens](https://github.com/settings/tokens). Make sure you give the token the `repo` scope.
5. Create a PR to promote the images:

    ```bash
    export GITHUB_TOKEN=<your GH token>
    make promote-images
    ```

    **Notes**:
     *`make promote-images` target tries to figure out your Github user handle in order to find the forked [k8s.io](https://github.com/kubernetes/k8s.io) repository.
          If you have not forked the repo, please do it before running the Makefile target.
     * if `make promote-images` fails with an error like `FATAL while checking fork of kubernetes/k8s.io` you may be able to solve it by manually setting the USER_FORK variable i.e.  `export USER_FORK=<personal GitHub handle>`
     * `kpromo` uses `git@github.com:...` as remote to push the branch for the PR. If you don't have `ssh` set up you can configure
       git to use `https` instead via `git config --global url."https://github.com/".insteadOf git@github.com:`.
     * This will automatically create a PR in [k8s.io](https://github.com/kubernetes/k8s.io) and assign the CAPA maintainers.
6. Wait for the PR to be approved (typically by CAPA maintainers authorized to merge PRs into the k8s.io repository) and merged.
7. Verify the images are available in the production registry:

    ```bash
    docker pull registry.k8s.io/cluster-api-aws/cluster-api-aws-controller:${VERSION}
    ```


## Verify and Publish the draft release

1. Verify that all the files below are attached to the drafted release:
    1. `clusterawsadm-darwin-amd64`
    1. `clusterawsadm-darwin-arm64`
    1. `clusterawsadm-linux-amd64`
    1. `clusterawsadm-linux-arm64`
    1. `clusterawsadm-windows-amd64.exe`
    1. `clusterawsadm-windows-arm64.exe`
    1. `infrastructure-components.yaml`
    1. `cluster-template.yaml`
    1. `cluster-template-machinepool.yaml`
    1. `cluster-template-eks.yaml`
    1. `cluster-template-eks-ipv6.yaml`
    1. `cluster-template-eks-fargate.yaml`
    1. `cluster-template-eks-managedmachinepool.yaml`
    1. `cluster-template-eks-managedmachinepool-vpccni.yaml`
    1. `cluster-template-eks-managedmachinepool-gpu.yaml`
    1. `cluster-template-external-cloud-provider.yaml`
    1. `cluster-template-flatcar.yaml`
    1. `cluster-template-machinepool.yaml`
    1. `cluster-template-multitenancy-clusterclass.yaml`
    1. `cluster-template-rosa-machinepool.yaml`
    1. `cluster-template-rosa.yaml`
    1. `cluster-template-simple-clusterclass.yaml`
    1. `metadata.yaml`
1. Update the release description to link to the promotion image.
1. Publish release. Use the pre-release option for release candidate versions of Cluster API Provider AWS.
1. Email `kubernetes-sig-cluster-lifecycle@googlegroups.com` to announce the release. You can use this template for the email:

    ```
    Subject: [ANNOUNCE] cluster-api-provider-aws v2.1.0 is released
    Body:
    The cluster-api-provider-aws (CAPA) project has published a new release. Please see here for more details:
    Release v2.1.0 · kubernetes-sigs/cluster-api-provider-aws (github.com)

    If you have any questions about this release or CAPA, please join us on our Slack channel:
    https://kubernetes.slack.com/archives/CD6U2V71N
    ```

1. Update the Title and Description of the Slack channel to point to the new version.

## Post Release Steps

### Create new Prow jobs

If the release is for a new MAJOR.MINOR version (i.e. not a patch release) then you will need to create a new set of prow jobs for the release branch.

This is done by updating the [test-infra](https://github.com/kubernetes/test-infra) repo. For an example of PR see [this](https://github.com/kubernetes/test-infra/pull/33751) for the v2.7 release series.

Consider removing jobs from an old release as well. We should only keep jobs for 4 release branches.
