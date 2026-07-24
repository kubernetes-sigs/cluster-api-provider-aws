# Bumping Go

This document describes how to bump the Go version across the project. It is
primarily intended to be consumed by an AI coding agent (e.g. via
`/bump-go 1.26`), but the steps can also be followed manually.

## Convention

- `go` directive in `go.mod`: use `X.Y.0` (the minor version with patch `0`)
- `GO_VERSION` in Makefile: use `X.Y.Z` where Z is the **latest available
  patch** (e.g. `1.26.4`)

## Step 1: Research

Perform these lookups before making any changes.

### 1a. Latest patch version

Find the latest Go `X.Y.x` patch release at <https://go.dev/doc/devel/release>.
Determine the full version string (e.g. `1.26.4`). Call this `FULL_VERSION`.

### 1b. Upstream cluster-api references

Look up what `kubernetes-sigs/cluster-api` uses at the latest release tag on the
main CAPI minor version this project depends on (check `go.mod` for
`sigs.k8s.io/cluster-api` to find the CAPI minor, then find the latest tag for
that minor):

```bash
# GCB image digest + tag comment
curl -sL https://raw.githubusercontent.com/kubernetes-sigs/cluster-api/<TAG>/cloudbuild.yaml

# golangci-lint version
curl -sL https://raw.githubusercontent.com/kubernetes-sigs/cluster-api/<TAG>/.github/workflows/pr-golangci-lint.yaml \
  | grep 'version:'
```

Call these `GCB_DIGEST`, `GCB_TAG_COMMENT`, and `CAPI_GOLANGCI_VER`.

### 1c. golangci-lint compatibility

The golangci-lint version must be **built with the target Go version** or newer.
Versions built with an older Go will refuse to lint. Test candidate versions:

```bash
go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@<VERSION> version
```

Look for `built with goX.Y` in the output. Pick the latest stable version that
reports the target Go minor or newer. Call this `GOLANGCI_VERSION` (e.g.
`v2.12.2`).

## Step 2: Update files

### go.mod (root), hack/tools/go.mod, and hack/tools/release-tools/go.mod

Update the `go` directive:

```
go X.Y.0
```

### Makefile

```makefile
GO_VERSION ?= FULL_VERSION
GO_DIRECTIVE_VERSION ?= X.Y.0
```

### .golangci-kal.yml

Update the `go:` field under `run:`:

```yaml
run:
  go: "X.Y"
```

### hack/tools/Makefile

Update `GOLANGCI_LINT_VERSION` and the comment:

```makefile
# Use a golangci-lint built with Go >= X.Y to match repo go.mod
GOLANGCI_LINT_VERSION := GOLANGCI_VERSION
```

### hack/tools/.custom-gcl.yaml

Update the `version:` field used to build the custom KAL linter:

```yaml
version: GOLANGCI_VERSION
```

### .github/workflows/pr-golangci-lint.yaml

Update the golangci-lint `version:`:

```yaml
version: GOLANGCI_VERSION
```

(The `go-version` in this workflow is already derived dynamically via
`make go-version`, so it does not need a manual update.)

### .github/workflows/dependabot.yml

Update the `go-version:` field:

```yaml
go-version: 'X.Y'
```

### .github/workflows/release.yaml

Update the `go-version:` field:

```yaml
go-version: 'X.Y'
```

### netlify.toml

```toml
GO_VERSION = "FULL_VERSION"
```

### cloudbuild.yaml and cloudbuild-nightly.yaml

Only update if the upstream CAPI release tag uses a **different** GCB image
digest than what this repo already has. If they match, skip this file.

```yaml
- name: 'gcr.io/k8s-staging-test-infra/gcb-docker-gcloud@sha256:GCB_DIGEST' # GCB_TAG_COMMENT
```

## Step 3: go mod tidy

Run in all module directories:

```bash
go mod tidy
cd hack/tools && go mod tidy
cd release-tools && go mod tidy
```

## Step 4: Lint and test

```bash
make lint
make lint-api
make test
```

If `make lint` or `make lint-api` fails:
- If the failure is `the Go language version used to build golangci-lint is
  lower than the targeted Go version`, the chosen `GOLANGCI_VERSION` is wrong.
  Go back to Step 1c. For `make lint-api`, also check
  `hack/tools/.custom-gcl.yaml`.
- If the failure shows new lint findings, fix them. These are pre-existing
  issues surfaced by the newer linter version, not caused by the Go bump itself.
- If a new linter version produces noisy false positives (e.g. `prealloc` on
  idiomatic `var allErrs field.ErrorList` patterns), add targeted exclusion
  rules in `.golangci.yml` rather than suppressing the linter entirely.

## Step 5: Verify build

```bash
go build ./...
```

## Step 6: Commit

Create **three separate commits** in this order:

1. **Lint fixes** (if any): `fix(lint): resolve <linter-names> lint issues`
   - Only the source files with lint fixes
2. **Linter bump** (if version changed):
   `chore(bump): bump golangci-lint from <old> to <new>`
   - `hack/tools/Makefile`, `hack/tools/.custom-gcl.yaml`,
     `.github/workflows/pr-golangci-lint.yaml`, and `.golangci.yml` if
     exclusion rules were added
3. **Go version bump**: `chore(bump): bump Go to FULL_VERSION`
   - All remaining files: `go.mod`, `hack/tools/go.mod`,
     `hack/tools/release-tools/go.mod`, `Makefile`,
     `.golangci-kal.yml`, `.github/workflows/dependabot.yml`,
     `.github/workflows/release.yaml`, `netlify.toml`,
     `cloudbuild.yaml`, `cloudbuild-nightly.yaml`

If there are no lint fixes or no linter version change, skip that commit.

All commits must be signed off (`git commit -s`).

Do NOT push or create a PR unless the user asks. If a PR is created, use the
pull request template defined in `.github/PULL_REQUEST_TEMPLATE.md`.
