# Cluster API Provider AWS

This is Giant Swarm's fork. See the upstream [cluster-api-provider-aws README](https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/main/README.md) for official documentation.

## How to work with this repo

Currently, we try to follow the upstream `release-X.Y` branch to always get the latest stable release and fixes, but not untested commits from `main`. Our only differences against upstream should be in this `README.md`, `.circleci/` and `.github/workflows/release.yml`. Other changes should be opened as PR for the upstream project first.

We release cluster-api-provider-aws versions with [cluster-api-provider-aws-app](https://github.com/giantswarm/cluster-api-provider-aws-app/). To provide the YAML manifests, we use GitHub releases as the upstream project. The scripts in `cluster-api-provider-aws-app` convert them into the final manifests.

### Repo setup

Since we follow upstream, add their Git repo as remote from which we merge commits:

```sh
git clone git@github.com:giantswarm/cluster-api-provider-aws.git
cd cluster-api-provider-aws
git remote add upstream https://github.com/kubernetes-sigs/cluster-api-provider-aws.git
```

### Test and release

If you have a non-urgent fix, create an upstream PR and wait until it gets released. We call this release `vX.Y.Z` in the below instructions, so please fill in the desired tag.

Please follow the development workflow:

- Ensure a stable release branch exists in our fork repo. For example with a desired upstream release v2.2.1, the branch is `release-2.2`. If it does not exist on our side yet, copy the branch from upstream and add our changes such as `README.md` and `.circleci/` on top.
- Create a working branch for your changes
- We want to use stable upstream release tags unless a hotfix is absolutely required ([decision](https://intranet.giantswarm.io/docs/product/pdr/010_fork_management/)). Please decide what type of change you're making:

  - Either: you want to merge and test the latest upstream tag

    ```sh
    # Get latest changes on our release branch
    git checkout release-X.Y
    git pull

    git fetch upstream

    git checkout -b my-working-branch release-X.Y

    # Create a merge commit using upstream's desired release tag (the one we want
    # to upgrade to).
    # This creates a commit message such as "Merge tag 'v2.2.1' into release-2.2".
    git merge --no-ff vX.Y.Z

    # Since we want the combined content of our repo and the upstream Git tag,
    # we need to create our own tag on the merge commit
    git tag "vX.Y.Z-gs-$(git rev-parse --short HEAD)"

    # Push your working branch. This triggers image build in CircleCI.
    git push

    # Push your Giant Swarm tag (assuming `origin` is the Giant Swarm fork).
    # This triggers the GitHub release action - please continue reading below!
    git push origin "vX.Y.Z-gs-$(git rev-parse --short HEAD)"
    ```

  - Or: you want to implement something else, such as working on some issue that we have which is not fixed in upstream yet. Note that for testing changes to upstream, you probably better base your work on the `upstream/main` branch and try your change together with the latest commits from upstream. This also avoids merge conflicts. Maintainers can then help you cherry-pick into their release branches. The latest release branch is usually a bit behind `main`.

    ```sh
    # Get latest changes on our release branch
    git checkout release-X.Y
    git pull

    git checkout -b my-working-branch release-X.Y # or based on `main` instead of `release-X.Y`, see hint above

    # Make some changes and commit as usual
    git commit

    git tag "vX.Y.Z-gs-$(git rev-parse --short HEAD)"

    # Push your working branch. This triggers image build in CircleCI
    git push

    # Push your Giant Swarm tag (assuming `origin` is the Giant Swarm fork).
    # This triggers the GitHub release action - please continue reading below!
    git push origin "vX.Y.Z-gs-$(git rev-parse --short HEAD)"
    ```

- Check that the [CircleCI pipeline](https://app.circleci.com/pipelines/github/giantswarm/cluster-api-provider-aws) succeeds for the desired Git tag in order to produce images. If the tag build fails, fix it.
- Check that the [GitHub release action](https://github.com/giantswarm/cluster-api-provider-aws/actions) for the `vX.Y.Z-gs-...` tag succeeds
- Edit [that draft GitHub release](https://github.com/giantswarm/cluster-api-provider-aws/releases) and turn it from draft to released. This makes the release's manifest files available on the internet, as used in [cluster-api-provider-aws-app](https://github.com/giantswarm/cluster-api-provider-aws-app).
- Test the changes in the app

  - Replace `.tag` in [cluster-api-provider-aws-app's `values.yaml`](https://github.com/giantswarm/cluster-api-provider-aws-app/blob/master/helm/cluster-api-provider-aws/values.yaml) with the new tag `vX.Y.Z-gs-...`.
  - Run `cd cluster-api-provider-aws-app && make generate` to update manifests
  - Commit and push your working branch for `cluster-api-provider-aws-app` to trigger CircleCI pipeline
  - Install and test the app thoroughly on a management cluster. Continue with the next step only once you're confident.
- Open PR for `cluster-api-provider-aws` fork (your working branch)

  - If you merged an upstream release tag, we should target our `release-X.Y` branch with the PR.
  - On the other hand, if you implemented something else which is not in upstream yet, we should target `upstream/main` so that it first lands in the upstream project, officially approved, tested and released. Afterwards, you would repeat this whole procedure and merge the release that includes your fix. For a quick in-house hotfix, you can alternatively do a quicker PR targeted against our `release-X.Y` branch.
- Also open PR for `cluster-api-provider-aws-app` change
- Once merged, manually bump the version in the respective collection to deploy it for one provider (e.g. [capa-app-collection](https://github.com/giantswarm/capa-app-collection/))

### Keep fork customizations up to date

Only `README.md`, `.circleci/` and `.github/workflows/release.yml` should differ between upstream and our fork, so the diff of everything else should be empty, or at worst, contain hotfixes that are not in upstream yet:

```sh
git fetch upstream
git diff `# the upstream tag we merged recently` vX.Y.Z..origin/release-X.Y `# our release branch` -- ':!.circleci/' "!.github/workflows/release.yml' ':!README.md'
```

And we should also keep our `main` and `release-X.Y` branches in sync, so this diff should be empty:

```sh
git diff main..release-X.Y -- .circleci/ .github/workflows/release.yml README.md
```

If this shows any output, please align the `main` branch with the release branches.
