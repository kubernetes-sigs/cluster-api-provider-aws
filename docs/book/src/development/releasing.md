# Release process

## Manual

1. Make sure your repo is clean by git's standards
2. Set environment variable `GITHUB_TOKEN` to a GitHub personal access token
3. If this is a new minor release, create a new release branch and push to GitHub, otherwise switch to it, for example `release-0.6`
4. Tag the repository and push the tag `git tag -s -m $VERSION $VERSION`. `-s` flag is for GNU Privacy Guard (GPG) signing
5. Push the commit of the tag to the release branch: `git push origin HEAD:release-0.6`
6. Set environment variables `PREVIOUS_VERSION` which is the last release tag and `VERSION` which is the current release version.
7. Checkout the tag you've just created and make sure git is in a clean state
8. Export the current branch `BRANCH=release-0.6` (`BRANCH=main`)and run `make release`
9. A prow job will start running to push images to the staging repo, can be seen [here](https://testgrid.k8s.io/sig-cluster-lifecycle-image-pushes#post-cluster-api-provider-aws-push-images).
10. Run `make create-gh-release` to create a draft release on Github, copying the generated release notes from `out/CHANGELOG.md` into the draft.
11. Run `make upload-gh-artifacts` to upload artifacts from .out/ directory, however you may run into API limit errors, so verify artifacts at next step
12. Verify that all the files below are attached to the drafted release:
    1. `clusterawsadm-darwin-amd64`
    2. `clusterawsadm-linux-amd64`
    3. `infrastructure-components.yaml`
    4. `cluster-template.yaml`
    5. `cluster-template-machinepool.yaml`
    6. `cluster-template-eks.yaml`
    7. `cluster-template-eks-managedmachinepool.yaml`
    8. `cluster-template-eks-managedmachinepool-vpccni.yaml`
    9. `cluster-template-eks-managedmachinepool-gpu.yaml`
    10. `eks-controlplane-components.yaml`
    11. `eks-bootstrap-components.yaml`
    12. `metadata.yaml`
13. Perform the image promotion process to promote the images from the staging to production registry (`registry.k8s.io/cluster-api-provider-aws`):
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
14. Finalise the release notes. Add image locations `<ADD_IMAGE_HERE>` (e.g., registry.k8s.io/cluster-api-aws/cluster-api-aws-controller:v0.6.4) and replace `<RELEASE_VERSION>` and `<PREVIOUS_VERSION>`.
15. Make sure image promotion is complete before publishing the release draft. The promotion job logs can be found [here](https://testgrid.k8s.io/sig-k8s-infra-k8sio#post-k8sio-image-promo) and you can also try and pull the images (i.e. ``docker pull registry.k8s.io/cluster-api-aws/cluster-api-aws-controller:v0.6.4`) 
16. Publish release. Use the pre-release option for release
     candidate versions of Cluster API Provider AWS.
17. Email `kubernetes-sig-cluster-lifecycle@googlegroups.com` to announce the release
