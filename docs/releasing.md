# Release process

## Manual

1. Make sure your repo is clean by git's standards
2. If this is a new minor release, create a new release branch and push to github, for example `release-0.2`
3. If using the release notes tool (documented below), run it now prior to creating the tag
4. Tag the repository and push the tag `git tag -s -m $VERSION $VERSION`
5. Create a draft release in github and associate it with the tag that was just created
6. Checkout the tag you've just created and make sure git is in a clean state
7. Run `make release`
8. Attach the files to the drafted release:
    1. `./out/clusterawsadm-darwin-amd64`
    2. `./out/clusterawsadm-linux-amd64`
    3. `./out/infrastructure-components.yaml`
    4. `./templates/cluster-template.yaml`
9. Perform the [image promotion process](https://github.com/kubernetes/k8s.io/tree/master/k8s.gcr.io#image-promoter)
10. Write the release notes (see note below on release notes)
11. Publish release
12. Email `kubernetes-sig-cluster-lifecycle@googlegroups.com` to announce the release

### Release Notes

Running `make release-notes` will generate an output that can be copied to the drafted release.
Pay close attention to the `## :question: Sort these by hand` section, as it contains items that need to be manually sorted.


