# Release process

## Manual
1. Make sure your repo is clean by git's standards
2. If this is a new minor release, create a new release branch and push to github, otherwise switch to it, for example `release-0.2`
3. Tag the repository and push the tag `git tag -s -m $VERSION $VERSION`. `-s` flag is for GNU Privacy Guard (GPG) signing
4. Push the commit of the tag to the release branch: `git push origin HEAD:release-0.2`
5. Set environment variables `PREVIOUS_VERSION` which is the last release tag and `VERSION` which is the current release version.
6. Checkout the tag you've just created and make sure git is in a clean state
7. Run `make release`
8.  A prow job will start running to push images to the staging repo, can be seen here: https://testgrid.k8s.io/sig-cluster-lifecycle-image-pushes#post-cluster-api-provider-aws-push-images
9. Run `make create-gh-release` to create a draft release on Github, copying the generated release notes from out/CHANGELOG.md into the draft.
10. Run `make upload-gh-artifacts` to upload artifacts from .out/ directory, however you may run into API limit errors, so verify artifacts at next step
11. Verify that all the files below are attached to the drafted release:
    1. `clusterawsadm-darwin-amd64`
    2. `clusterawsadm-linux-amd64`
    3. `infrastructure-components.yaml`
    4. `cluster-template.yaml`
    5. `cluster-template-machinepool.yaml`
    6. `cluster-template-eks.yaml`
    7. `cluster-template-eks-managedmachinepool.yaml`
    8. `cluster-template-eks-managedmachinepool-vpccni.yaml`
    9. `eks-controlplane-components.yaml`
    10.`eks-bootstrap-components.yaml`
    11. `metadata.yaml`
12. Perform the [image promotion process](https://github.com/kubernetes/k8s.io/tree/master/k8s.gcr.io#image-promoter).
    The staging repository is at https://console.cloud.google.com/gcr/images/k8s-staging-cluster-api-aws/GLOBAL. Be
    sure to choose the top level `cluster-api-aws-controller`, which will provide the multi-arch manifest, rather than one for a specific architecture.
13.  Finalise the release notes. Add image locations `<ADD_IMAGE_HERE>` (e.g., us.gcr.io/k8s-artifacts-prod/cluster-api-aws/cluster-api-aws-controller:v0.6.4)
14.  Make sure image promotion is complete before publishing the release draft. The promotion job logs can be found here: https://testgrid.k8s.io/wg-k8s-infra-k8sio#post-k8sio-image-promo
15.  Publish release. Use the pre-release option for release
     candidate versions of Cluster API Provider AWS.
16.  Email `kubernetes-sig-cluster-lifecycle@googlegroups.com` to announce the release
