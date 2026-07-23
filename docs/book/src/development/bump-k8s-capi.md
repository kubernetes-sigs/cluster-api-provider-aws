# Bumping Kubernetes and Cluster API

This document describes how to bump Kubernetes and Cluster API (CAPI) versions
across the project. It is primarily intended to be consumed by an AI coding
agent (e.g. via `/bump-k8s-capi 1.35 1.13`), but the steps can also be followed
manually.

The first argument is the target **Kubernetes minor** version (e.g. `1.35`) and
the second is the target **CAPI minor** version (e.g. `1.13`).

## Step 1: Research

Perform these lookups before making any changes. Call the first argument
`K8S_MINOR` and the second `CAPI_MINOR`.

### 1a. Latest CAPI patch

Find the latest patch release for the target CAPI minor:

```bash
gh api repos/kubernetes-sigs/cluster-api/releases \
  --paginate --jq '.[].tag_name' \
  | grep "^v${CAPI_MINOR}\." | head -1
```

Call this `CAPI_VERSION` (e.g. `v1.13.2`).

### 1b. Dependency versions from CAPI

From CAPI's `go.mod` at that tag, extract the **minor** version of
`controller-runtime` and `cert-manager`. Then find the **latest patch** for each
minor — we do not need to use the exact same patch CAPI uses, only the same
minor.

```bash
# controller-runtime minor used by CAPI
gh api repos/kubernetes-sigs/cluster-api/contents/go.mod?ref=${CAPI_VERSION} \
  --jq '.content' | base64 -d | grep 'sigs.k8s.io/controller-runtime'

# latest patch for that minor
gh api repos/kubernetes-sigs/controller-runtime/releases \
  --paginate --jq '.[].tag_name' \
  | grep "^v0.XX\." | head -1
```

Call this `CONTROLLER_RUNTIME_VERSION`.

For cert-manager, check `versions.mk` or the CAPI release notes for the minor
they target, then find the latest patch:

```bash
gh api repos/cert-manager/cert-manager/releases \
  --paginate --jq '.[].tag_name' \
  | grep "^v1.XX\." | head -1
```

Call this `CERT_MANAGER_VERSION`.

Also extract from CAPI's `go.mod` the `setup-envtest` release branch (typically
`release-0.XX` matching controller-runtime minor). Call this
`SETUP_ENVTEST_VER`.

### 1c. Kubernetes patch version

Find the latest stable patch for the target k8s minor:

```bash
gh api repos/kubernetes/kubernetes/releases \
  --paginate --jq '.[].tag_name' \
  | grep "^v${K8S_MINOR}\." | grep -v -E '(alpha|beta|rc)' | head -1
```

Call this `K8S_VERSION` (e.g. `v1.35.3`).

### 1d. KIND version and management cluster k8s version

KIND does not publish `kindest/node` images for every k8s patch — only for
specific patches chosen at each KIND release. Find the latest KIND release that
has a `kindest/node` image for the target k8s minor:

```bash
gh api repos/kubernetes-sigs/kind/releases/latest --jq '.body'
```

Look for `kindest/node:v${K8S_MINOR}.X` in the release notes. The patch version
may differ from `K8S_VERSION`. Call this `KIND_K8S_VERSION` (e.g. `v1.35.5`) and
the KIND release itself `KIND_VERSION` (e.g. `v0.32.0`).

If the latest KIND release does not have an image for the target k8s minor,
check earlier releases or consider whether the bump should wait for a new KIND
release.

### 1e. AMI availability

Verify that AMIs are published for the target workload k8s version(s). The
project has an automated AMI pipeline (`auto-publish-ami` workflow) that builds
AMIs for new k8s releases, but there may be a lag.

Check the AMI build config:

```bash
cat hack/tools/release-tools/internal/ami/k8srelease/data/AMIK8sBuildConfig.json
```

Or, if you have AWS credentials configured:

```bash
clusterawsadm ami list
```

If AMIs for the target version are missing, warn the user. The
`detect-k8s-releases` workflow should pick them up automatically, but the user
may need to merge the resulting automated PR first to trigger builds.

### 1f. Ancillary versions

Derive the etcd and coredns versions that ship with the target k8s version.
These are needed for e2e upgrade test configuration.

Check the Kubernetes changelog or component versions:

```bash
gh api repos/kubernetes/kubernetes/contents/CHANGELOG/CHANGELOG-${K8S_MINOR}.md?ref=${K8S_VERSION} \
  --jq '.content' | base64 -d | grep -i -E 'etcd|coredns' | head -10
```

Call these `ETCD_VERSION` and `COREDNS_VERSION`.

Also determine suitable values for the "upgrade from" versions:
`KUBERNETES_VERSION_UPGRADE_FROM` (typically the latest patch of the previous
k8s minor that has a published AMI).

### 1g. CAPI migration guide

Read the CAPI migration guide for breaking changes:

```bash
gh api repos/kubernetes-sigs/cluster-api/contents/docs/release/releases \
  --jq '.[].name' | sort
```

Look for a file matching the target version transition (e.g.
`release-1.12-to-1.13.md`). Fetch and read it:

```bash
gh api repos/kubernetes-sigs/cluster-api/contents/docs/release/releases/release-1.XX-to-1.YY.md \
  --jq '.content' | base64 -d
```

Note any required code changes: import path updates, API type renames,
controller-runtime interface changes, predicate renames, webhook changes,
utility migrations, etc.

### 1h. vSphere provider reference (advisory)

Check if `kubernetes-sigs/cluster-api-provider-vsphere` has already done a bump
to the same CAPI/k8s versions. If so, skim their diff for non-obvious gotchas:

```bash
gh pr list --repo kubernetes-sigs/cluster-api-provider-vsphere \
  --search "bump CAPI ${CAPI_MINOR} OR bump kubernetes ${K8S_MINOR}" \
  --state all --limit 5
```

If a matching PR exists, review it:

```bash
gh pr diff <PR_NUMBER> --repo kubernetes-sigs/cluster-api-provider-vsphere
```

Treat this as **advisory only** — the providers differ structurally, so not all
changes will apply. Look for patterns like: unexpected new dependencies, new
required RBAC, generator flag changes, or test framework migrations.

## Step 2: Update dependency files

### go.mod (root)

Bump the following direct dependencies:

```bash
go get sigs.k8s.io/cluster-api@${CAPI_VERSION}
go get sigs.k8s.io/cluster-api/test@${CAPI_VERSION}
go get sigs.k8s.io/controller-runtime@${CONTROLLER_RUNTIME_VERSION}
```

Bump all `k8s.io/*` modules to the version matching `K8S_MINOR` (the exact
patch is determined by what CAPI depends on — check CAPI's `go.mod`):

```bash
go get k8s.io/api@v0.XX.Y
go get k8s.io/apimachinery@v0.XX.Y
go get k8s.io/client-go@v0.XX.Y
go get k8s.io/apiextensions-apiserver@v0.XX.Y
go get k8s.io/apiserver@v0.XX.Y
go get k8s.io/cli-runtime@v0.XX.Y
go get k8s.io/component-base@v0.XX.Y
go get k8s.io/kubectl@v0.XX.Y
```

Also bump `sigs.k8s.io/kind` if needed:

```bash
go get sigs.k8s.io/kind@${KIND_VERSION}
```

Review and remove any `replace` directives that are no longer needed (e.g.
workarounds for bugs fixed in the new versions).

### hack/tools/go.mod

Bump `k8s.io/apimachinery` to match the root module:

```bash
cd hack/tools && go get k8s.io/apimachinery@v0.XX.Y
```

### go mod tidy

Run in both module directories:

```bash
go mod tidy
cd hack/tools && go mod tidy
```

## Step 3: Update version references

### versions.mk

```makefile
CAPI_VERSION := CAPI_VERSION
CERT_MANAGER_VERSION := CERT_MANAGER_VERSION
```

### Makefile

```makefile
KUBEBUILDER_ENVTEST_KUBERNETES_VERSION := K8S_MINOR.0
```

Update `SETUP_ENVTEST_VER` to the controller-runtime release branch:

```makefile
SETUP_ENVTEST_VER := release-0.XX
```

### metadata.yaml

Add a new CAPA release series entry. Determine the next CAPA minor by
incrementing the last entry in the file. Use the appropriate CAPI contract
version (check the CAPI release notes — e.g. `v1beta1` or `v1beta2`):

```yaml
  - major: 2
    minor: XX
    contract: v1beta2
```

## Step 4: Update e2e configuration

### test/e2e/data/e2e_conf.yaml

Update:
- All CAPI component image URLs (search for
  `cluster-api/releases/download/v1.XX` and update the version)
- `KUBERNETES_VERSION_MANAGEMENT`: set to `KIND_K8S_VERSION` (must match a
  version with a published `kindest/node` image)
- `KUBERNETES_VERSION`: set to a recent k8s patch with published AMIs
- `KUBERNETES_VERSION_UPGRADE_TO`: typically same as `KUBERNETES_VERSION`
- `KUBERNETES_VERSION_UPGRADE_FROM`: latest patch of the previous k8s minor
- `INIT_WITH_KUBERNETES_VERSION`: set to `KIND_K8S_VERSION`
- `ETCD_VERSION_UPGRADE_TO`: set to `ETCD_VERSION`
- `COREDNS_VERSION_UPGRADE_TO`: set to `COREDNS_VERSION`
- `CONFORMANCE_CI_ARTIFACTS_KUBERNETES_VERSION`: set to match `KUBERNETES_VERSION`
- cert-manager image tags: update to `CERT_MANAGER_VERSION`

### test/e2e/data/e2e_eks_conf.yaml

Apply the same version updates as above where the variables appear.

### test/e2e/data/shared/capi/ metadata files

Each CAPI minor version gets its own metadata directory named after that version
(e.g. `capi/v1.13/` for CAPI v1.13). This follows the same convention used by
cluster-api-provider-vsphere.

The metadata file lists all known CAPI release series up to and including that
minor, mapping each to its contract version (`v1beta1` or `v1beta2`).

When bumping to a new CAPI minor N:
1. Create `test/e2e/data/shared/capi/v{N}/metadata.yaml` with all release
   series up to minor N
2. Point the new CAPI version entries in both `e2e_conf.yaml` and
   `e2e_eks_conf.yaml` at `./shared/capi/v{N}/metadata.yaml`
3. Delete the oldest `capi/v{old}/` directory if it is no longer referenced
4. Update `capi/v1.2/metadata.yaml` to include the new release series (this
   catch-all file is used for v1.2.0 upgrade test entries and provider-level
   defaults)

## Step 5: Apply migration-guide-driven changes

Using the migration guide from Step 1g, apply any required code changes. This
step varies per bump and cannot be fully prescribed. Common categories from
historical bumps include:

- **Import path updates**: CAPI may reorganize module paths (e.g.
  `api/v1beta1` -> `api/core/v1beta2`, conditions/patch to
  `util/deprecated/v1beta1/`)
- **API type changes**: fields may change types (e.g. `*errors.MachineStatusError`
  -> `*string`, `FailureDomainSpec` -> `FailureDomain`, pointer vs value changes)
- **controller-runtime interface changes**: webhook interfaces, predicate
  signatures, and builder patterns may change between minor versions
- **Utility migrations**: CAPI may make public packages internal or rename
  utilities (e.g. `cmd.LongDesc` -> `templates.LongDesc`, `gofuzz` ->
  `randfill`)
- **`.golangci.yml` alias updates**: import alias rules may need updating when
  CAPI reorganizes module paths
- **Makefile generator flags**: `--extra-peer-dirs` flags for `DEFAULTER_GEN` and
  `CONVERSION_GEN` may need updating
- **Predicate renames**: e.g.
  `ClusterUnpausedAndInfrastructureReady` ->
  `ClusterPausedTransitionsOrInfrastructureReady` ->
  `ClusterPausedTransitionsOrInfrastructureProvisioned`

## Step 6: Regenerate and validate

```bash
make generate
make lint
make test
```

Fix any failures before proceeding. See the Troubleshooting section below for
common issues.

## Step 7: Commit

Create separate commits for distinct concerns:

1. **Migration/adaptation fixes** (if any):
   `fix: apply CAPI vCAPI_MINOR migration changes`
   - Only the source files with migration-driven code changes
2. **The version bump**: `chore(bump): bump to Kubernetes K8S_MINOR and Cluster
   API vCAPI_MINOR`
   - All remaining files: `go.mod`, `go.sum`, `hack/tools/go.mod`,
     `hack/tools/go.sum`, `versions.mk`, `Makefile`, `metadata.yaml`, e2e
     configs, generated files

If the migration changes are trivial (e.g. only import path renames), they can
be folded into the bump commit.

All commits must be signed off (`git commit -s`).

Do NOT push or create a PR unless the user asks. If a PR is created, use the
pull request template defined in `.github/PULL_REQUEST_TEMPLATE.md`.

## Step 8: CI validation

After the PR is pushed, suggest kicking off the **optional** full e2e and full
EKS e2e CI presubmit jobs to verify nothing is broken:

- `pull-cluster-api-provider-aws-e2e` (full e2e suite)
- `pull-cluster-api-provider-aws-e2e-eks` (full EKS e2e suite)

These are optional presubmits that must be triggered explicitly (e.g. via a
`/test` comment on the PR).

If they fail, compare presubmit results against periodic job history on main at:

<https://testgrid.k8s.io/sig-cluster-lifecycle-cluster-api-provider-aws>

This helps distinguish genuine regressions introduced by the bump from
pre-existing flakes on main.

## Troubleshooting

### `go mod tidy` failures

If `go mod tidy` fails due to incompatible transitive dependencies, check if
CAPI has `replace` directives in its `go.mod` at the target tag that need to be
mirrored temporarily:

```bash
gh api repos/kubernetes-sigs/cluster-api/contents/go.mod?ref=${CAPI_VERSION} \
  --jq '.content' | base64 -d | grep '^replace'
```

### `make generate` failures

If CRD or conversion generation fails, check if `--extra-peer-dirs` flags in
the Makefile need updating. This happens when CAPI reorganizes module paths
(e.g. PR #5720 changed paths from `sigs.k8s.io/cluster-api/api/v1beta1` to
`sigs.k8s.io/cluster-api/api/core/v1beta2`).

### `make test` failures

Compare against CAPI's own test changes at the target tag. Test helpers,
builders, and framework utilities may have changed signatures:

```bash
gh pr diff --repo kubernetes-sigs/cluster-api <CAPI_PR_NUMBER>
```

Check the CAPI release notes for any test migration guidance.

### Compilation errors after the bump

If there are widespread compile errors from CAPI API changes that the migration
guide did not cover, inspect the changed types directly:

```bash
go doc sigs.k8s.io/cluster-api/api/core/v1beta2.SomeType
```

And grep the codebase for usages to find all call sites that need updating.

### E2e test failures

Use TestGrid to compare presubmit results against periodic results on main:

<https://testgrid.k8s.io/sig-cluster-lifecycle-cluster-api-provider-aws>

If the same tests fail on main periodics, the failures are pre-existing flakes,
not caused by the bump.
