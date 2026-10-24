# release-tool

A Go CLI collecting automation helpers used by the CAPA release process, exposed as `release-tool <group> <command>`.

## Local development

```sh
go build -o release-tool .
./release-tool ami --help
```

or without building a binary:

```sh
go run . ami --help
```

## `ami` commands

### detect-k8s-release

Queries the `kubernetes/kubernetes` GitHub repository for the latest supported stable Kubernetes minor versions (the CAPA AMI build policy). Used by `.github/workflows/detect-k8s-releases.yml`. See `release-tool ami detect-k8s-release --help` for flags.

### find-missing-ami

Computes the AMIs that should be published (k8s versions × OS × regions) but aren't yet, based on a published-AMI inventory piped in via stdin (`clusterawsadm ami list -o json`). Used by `.github/workflows/auto-publish-ami.yml`. See `release-tool ami find-missing-ami --help` for flags.

### remove-duplicates

Finds — and optionally removes — duplicate AMIs owned by the CAPA AMI build pipeline.

Discovery is scoped to CAPA-built AMIs only: the command queries EC2 for images owned by `--owner-id` that carry a `kubernetes_version` tag (every CAPA AMI build sets one) and are `x86_64`, `available`, and `hvm` (also true of every CAPA AMI build). An image without that tag, such as an unrelated AMI that happens to live in the same account, is never fetched, so it can never be classified or removed.

Of the images returned, AMIs are grouped per-region by their `distribution`, `distribution_version` and `kubernetes_version` tags. Within each group, the AMI with the highest `build_timestamp` tag is **kept**; the rest are **duplicates**. An AMI missing one of the grouping tags, or with a missing/invalid `build_timestamp` tag, is **ungroupable** rather than silently skipped.

This command is **dry-run by default**: it always prints the full report (kept, duplicate, and ungroupable AMIs), then previews what removal would do. Pass `--delete` to actually deregister the targeted AMIs.

Only `DUPLICATE` AMIs are removal targets by default; pass `--include-ungroupable` to also target `UNGROUPABLE` AMIs. `KEEP` AMIs are never removed.

#### Usage

```sh
# Dry run across all CAPA regions
release-tool ami remove-duplicates

# Dry run in a single region, JSON report
release-tool ami remove-duplicates --region us-west-2 -o json

# Actually remove duplicate AMIs
release-tool ami remove-duplicates --delete

# Also remove AMIs that couldn't be evaluated for duplication
release-tool ami remove-duplicates --delete --include-ungroupable
```

#### Flags

| Flag | Required | Default | Description |
|---|---|---|---|
| `--region` | no | all CAPA regions | Comma-separated list of AWS regions to search |
| `--owner-id` | no | `819546954734` | AWS account ID that owns the AMIs to search |
| `--output`, `-o` | no | `table` | Output format: `table`, `json`, or `yaml` |
| `--delete` | no | `false` | Actually remove the AMIs (omit for dry-run preview) |
| `--include-ungroupable` | no | `false` | Also remove `UNGROUPABLE` AMIs (default: only `DUPLICATE`) |

#### Example output (dry run)

```text
REGION     GROUP                      STATUS       AMI ID       BUILD TIMESTAMP / REASON
us-east-1  ubuntu|22.04|v1.36.3       KEEP         ami-0abc123  1785150077
us-east-1  ubuntu|22.04|v1.36.3       DUPLICATE    ami-0def456  1785063677
us-east-1  -                          UNGROUPABLE  ami-0ghi789  missing "distribution" tag

Would remove 1 AMI(s):
  [DRY RUN] us-east-1 ami-0def456 (DUPLICATE)

Dry run — nothing removed. Pass --delete to remove them.
```

`--delete` deregisters each targeted AMI with `DeleteAssociatedSnapshots: true` — one API call removes both the AMI and its backing EBS snapshots.

#### GitHub Actions

See `.github/workflows/ami-housekeeping-remove-duplicates.yml` for a manual (`workflow_dispatch`) workflow that accepts `region`, `include_ungroupable`, and `delete` inputs, all with defaults.

Authentication is via OIDC (no stored credentials/secrets) — the workflow assumes `arn:aws:iam::819546954734:role/gh-image-builder`. That role needs:

- `ec2:DescribeImages`
- `ec2:DeregisterImage`
- `ec2:DeleteSnapshot`
