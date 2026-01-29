# Jobs

This document intends to provide an overview over our jobs running via Prow, GitHub actions and Google Cloud Build.

## Builds and Tests running on the main branch

> NOTE: To see which test jobs execute which tests or e2e tests, you can click on the links which lead to the respective test overviews in [test-grid].

### Presubmits

Prow Presubmits:

* [pull-cluster-api-provider-aws-test] `./scripts/ci-test.sh`
* [pull-cluster-api-provider-aws-build] `./scripts/ci-build.sh`
* [pull-cluster-api-provider-aws-verify] `make verify`
* [pull-cluster-api-provider-aws-e2e-conformance] `./scripts/ci-conformance.sh`
* [pull-cluster-api-provider-aws-e2e-conformance-with-ci-artifacts] `./scripts/ci-conformance.sh`
  * E2E_ARGS: `-kubetest.use-ci-artifacts`
* [pull-cluster-api-provider-aws-e2e-blocking] `./scripts/ci-e2e.sh`
  * GINKGO_FOCUS: `[PR-Blocking]`
* [pull-cluster-api-provider-aws-e2e] `./scripts/ci-e2e.sh`
* [pull-cluster-api-provider-aws-e2e-eks] `./scripts/ci-e2e-eks.sh`

[pull-cluster-api-provider-aws-e2e-eks]: https://testgrid.k8s.io/sig-cluster-lifecycle-cluster-api-provider-aws#pr-e2e-eks-main&show-stale-tests=
[pull-cluster-api-provider-aws-e2e]: https://testgrid.k8s.io/sig-cluster-lifecycle-cluster-api-provider-aws#pr-e2e-main&show-stale-tests=
[pull-cluster-api-provider-aws-e2e-blocking]: https://testgrid.k8s.io/sig-cluster-lifecycle-cluster-api-provider-aws#pr-quick-e2e-main&show-stale-tests=
[pull-cluster-api-provider-aws-e2e-conformance-with-ci-artifacts]: https://testgrid.k8s.io/sig-cluster-lifecycle-cluster-api-provider-aws#pr-conformance-main-k8s-main&show-stale-tests=
[pull-cluster-api-provider-aws-e2e-conformance]: https://testgrid.k8s.io/sig-cluster-lifecycle-cluster-api-provider-aws#pr-conformance&show-stale-tests=
[pull-cluster-api-provider-aws-verify]: https://testgrid.k8s.io/sig-cluster-lifecycle-cluster-api-provider-aws#pr-verify&show-stale-tests=
[pull-cluster-api-provider-aws-test]: https://testgrid.k8s.io/sig-cluster-lifecycle-cluster-api-provider-aws#pr-test&show-stale-tests=
[pull-cluster-api-provider-aws-build]: https://testgrid.k8s.io/sig-cluster-lifecycle-cluster-api-provider-aws#pr-build&show-stale-tests=

### Postsubmits

Prow Postsubmits:

* [ci-cluster-api-provider-aws-e2e] `./scripts/ci-e2e.sh`
* [ci-cluster-api-provider-aws-eks-e2e] `./scripts/ci-e2e-eks.sh`
* [ci-cluster-api-provider-aws-e2e-conformance] `./scripts/ci-conformance.sh`
  
[ci-cluster-api-provider-aws-e2e-conformance]: https://testgrid.k8s.io/sig-cluster-lifecycle-cluster-api-provider-aws#postsubmit-conformance-main&show-stale-tests=
[ci-cluster-api-provider-aws-eks-e2e]: https://testgrid.k8s.io/sig-cluster-lifecycle-cluster-api-provider-aws#postsubmit-eks-e2e-main&show-stale-tests=
[ci-cluster-api-provider-aws-e2e]: https://testgrid.k8s.io/sig-cluster-lifecycle-cluster-api-provider-aws#postsubmit-e2e-main&show-stale-tests=

* [post-cluster-api-provider-aws-push-images] Google Cloud Build: `make release-staging`

[post-cluster-api-provider-aws-push-images]: https://testgrid.k8s.io/sig-cluster-lifecycle-image-pushes#post-cluster-api-provider-aws-push-images

### Periodics

Prow Periodics:
* [periodic-cluster-api-provider-aws-e2e] `./scripts/ci-e2e.sh`
* [periodic-cluster-api-provider-aws-eks-e2e] `/scripts/ci-e2e-eks.sh`
* [periodic-cluster-api-provider-aws-e2e-conformance] `./scripts/ci-conformance.sh`
* [periodic-cluster-api-provider-aws-e2e-conformance-with-k8s-ci-artifacts] `./scripts/ci-conformance.sh`
  * E2E_ARGS: `-kubetest.use-ci-artifacts`
* [periodic-cluster-api-provider-aws-coverage] `./scripts/ci-test-coverage.sh`

[periodic-cluster-api-provider-aws-e2e-conformance-with-k8s-ci-artifacts]: https://testgrid.k8s.io/sig-cluster-lifecycle-cluster-api-provider-aws#periodic-conformance-main-k8s-main
[periodic-cluster-api-provider-aws-coverage]: https://testgrid.k8s.io/sig-cluster-lifecycle-cluster-api-provider-aws#periodic-test-coverage
[periodic-cluster-api-provider-aws-e2e-conformance]: https://testgrid.k8s.io/sig-cluster-lifecycle-cluster-api-provider-aws#periodic-conformance-main
[periodic-cluster-api-provider-aws-eks-e2e]: https://testgrid.k8s.io/sig-cluster-lifecycle-cluster-api-provider-aws#periodic-eks-e2e-main
[periodic-cluster-api-provider-aws-e2e]: https://testgrid.k8s.io/sig-cluster-lifecycle-cluster-api-provider-aws#periodic-e2e-main

* [cluster-api-provider-aws-push-images-nightly] Google Cloud Build: `make release-staging-nightly`

[cluster-api-provider-aws-push-images-nightly]: https://testgrid.k8s.io/sig-cluster-lifecycle-image-pushes#cluster-api-provider-aws-push-images-nightly
