# AMI Housekeeping

Tools for managing AWS AMIs tagged with `kubernetes_version` across multiple regions.

- **`list_amis.py`** — list matching AMIs with per-region and total counts (read-only)
- **`delete_amis.py`** — deregister matching AMIs and their backing EBS snapshots

> NOTE: this is a temporary measure until the AMI maintenance policy has been fully automated.

## Prerequisites

- [uv](https://docs.astral.sh/uv/) installed
- AWS credentials configured (env vars, `~/.aws/credentials`, or instance profile). When running via GitHub Actions the credentials are automatically set.

## Local development

```sh
uv sync
source .venv/bin/activate   # Windows: .venv\Scripts\activate
```

---

## list-amis

Lists AMIs whose `kubernetes_version` tag value contains a given substring, across multiple regions. Read-only — no changes made.

### Usage

```sh
python3 list_amis.py \
  --regions us-east-1,us-west-2,eu-west-1 \
  --version-contains 1.28
```

### Arguments

| Argument | Required | Default | Description |
|---|---|---|---|
| `--regions` | no | all major regions | Comma-separated AWS region names |
| `--version-contains` | yes | — | Substring matched against `kubernetes_version` tag value |
| `--owner` | no | `self` | AMI owner filter |

### Example output

```text
Filter: kubernetes_version contains '1.28'
Regions: ['us-east-1', 'us-west-2']

--- us-east-1 ---
  ami-0abc123  kubernetes_version=1.28.5  name=my-eks-node-1.28.5
  ami-0def456  kubernetes_version=1.28.7  name=my-eks-node-1.28.7
  Count: 2

--- us-west-2 ---
  No matching AMIs found.
  Count: 0

--- Summary ---
  us-east-1: 2
  us-west-2: 0
  Total: 2
```

### GitHub Actions

See `.github/workflows/list-amis.yml` for a manual workflow that accepts `kubernetes_version` and an optional `regions` override.

Required repository secret: `AWS_ROLE_ARN` — IAM role ARN to assume via OIDC.

The IAM role needs only: `ec2:DescribeImages`

---

## delete-amis

Deregisters AMIs (and their backing EBS snapshots) across multiple regions where the `kubernetes_version` tag value contains a given substring.

### Usage

```sh
python3 delete_amis.py \
  --regions us-east-1,us-west-2,eu-west-1 \
  --version-contains 1.28
```

**Default is dry-run** — no AMIs are deleted. You will see a preview of what would be removed.

To apply changes, add `--execute`:

```sh
python3 delete_amis.py \
  --regions us-east-1,us-west-2,eu-west-1 \
  --version-contains 1.28 \
  --execute
```

### Arguments

| Argument | Required | Default | Description |
|---|---|---|---|
| `--regions` | no | all major regions | Comma-separated AWS region names |
| `--version-contains` | yes | — | Substring matched against `kubernetes_version` tag value |
| `--execute` | no | false | Apply deletions (omit for dry-run preview) |
| `--owner` | no | `self` | AMI owner filter |

### Example output (dry run)

```text
Filter: kubernetes_version contains '1.28'
Regions: ['us-east-1', 'us-west-2']
Mode: DRY RUN (pass --execute to apply changes)

--- us-east-1 ---
  [DRY RUN] Found: ami-0abc123  kubernetes_version=1.28.5  name=my-eks-node-1.28.5
    [DRY RUN] Would deregister ami-0abc123
    [DRY RUN] Would delete snapshot snap-0def456

--- us-west-2 ---
  No matching AMIs found.

--- Summary ---
  us-east-1: 1
  us-west-2: 0
  Total: 1

Dry run complete — nothing deleted. Pass --execute to apply changes.
```

### GitHub Actions

See `.github/workflows/delete-amis.yml` for a manual workflow that accepts `kubernetes_version`, `execute`, and an optional `regions` override.

Required repository secret: `AWS_ROLE_ARN` — IAM role ARN to assume via OIDC.

The IAM role needs these permissions:

- `ec2:DescribeImages`
- `ec2:DeregisterImage`
- `ec2:DeleteSnapshot`
