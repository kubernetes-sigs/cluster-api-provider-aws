# Release process

## Manual

1. Make sure your repo is clean by git's standards
2. If this is a new minor release, create a new release branch and push to github, otherwise switch to it, for example `release-0.2`
3. Run `make release-notes` to gather changes since the last revision. If you need to specify a specific tag to look for changes
   since, use `make release-notes ARGS="--from <tag>"` Pay close attention to the `## :question: Sort these by hand` section, as it contains items that need to be manually sorted.
4. Tag the repository and push the tag `git tag -s -m $VERSION $VERSION`
5. Create a draft release in github and associate it with the tag that was just created, copying the generated release notes into
   the draft.
6. Checkout the tag you've just created and make sure git is in a clean state
7. Run `make release`
8. Attach the files to the drafted release:
    1. `./out/clusterawsadm-darwin-amd64`
    2. `./out/clusterawsadm-linux-amd64`
    3. `./out/infrastructure-components.yaml`
    4. `./templates/cluster-template.yaml`
9.  Perform the [image promotion process](https://github.com/kubernetes/k8s.io/tree/master/k8s.gcr.io#image-promoter).
    The staging repository is at https://console.cloud.google.com/gcr/images/k8s-staging-cluster-api-aws/GLOBAL. Be
    sure to choose the top level `cluster-api-aws-controller`, which will provide the multi-arch manifest, rather than one for a specific architecture.
10.  Finalise the release notes
11.  Publish release. Use the pre-release option for release
    candidate versions of Cluster API Provider AWS.
12.  Email `kubernetes-sig-cluster-lifecycle@googlegroups.com` to announce the release
