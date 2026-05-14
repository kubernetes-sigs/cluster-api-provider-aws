## Kubernetes Release

This section lists the Kubernetes versions tracked by this repository following the [CAPA AMI publication policy](https://cluster-api-aws.sigs.k8s.io/topics/images/built-amis#ami-publication-policy): the latest three supported minor releases and all their stable patch versions as detected by `release-tool ami detect-k8s-release`.

| Minor Version | Patch Versions |
| --- | --- |
| `1.36` | `1.36.0` |
| `1.35` | `1.35.4`, `1.35.3`, `1.35.2`, `1.35.1`, `1.35.0` |
| `1.34` | `1.34.7`, `1.34.6`, `1.34.5`, `1.34.4`, `1.34.3`, `1.34.2`, `1.34.1`, `1.34.0` |

## CAPA AMI

The table below lists all AMIs currently published in AWS account `027487054958`, as returned by `clusterawsadm ami list --owner-id 027487054958`.

| AMI Name | Kubernetes Version | OS | Region | AMI ID | Created |
| --- | --- | --- | --- | --- | --- |
| `capa-ami-ubuntu-24.04-v1.35.3-1778395590` | `v1.35.3` | ubuntu-24.04 | ap-southeast-2 | `ami-05c4b2006be0a8557` | 2026-05-10T06:58:16Z |
| `capa-ami-ubuntu-22.04-v1.35.2-1773582653` | `v1.35.2` | ubuntu-22.04 | ap-southeast-2 | `ami-04136e76627f61a0b` | 2026-03-15T14:04:21Z |
| `capa-ami-ubuntu-22.04-v1.35.1-1773640135` | `v1.35.1` | ubuntu-22.04 | ap-southeast-2 | `ami-01da8306336074553` | 2026-03-16T06:01:02Z |

## Missing AMI

### Default OS

List of OS for which AMIs should be published (default):

- ubuntu-24.04
- ubuntu-22.04

### Default Region

List of regions for which AMIs should be published (default):

- ap-southeast-2

### List of Missing AMIs

| Kubernetes Version | OS | Region |
| --- | --- | --- |
| `v1.36.0` | ubuntu-24.04 | ap-southeast-2 |
| `v1.36.0` | ubuntu-22.04 | ap-southeast-2 |
| `v1.35.4` | ubuntu-24.04 | ap-southeast-2 |
| `v1.35.4` | ubuntu-22.04 | ap-southeast-2 |
| `v1.35.3` | ubuntu-22.04 | ap-southeast-2 |
| `v1.35.2` | ubuntu-24.04 | ap-southeast-2 |
| `v1.35.1` | ubuntu-24.04 | ap-southeast-2 |
| `v1.35.0` | ubuntu-24.04 | ap-southeast-2 |
| `v1.35.0` | ubuntu-22.04 | ap-southeast-2 |
| `v1.34.7` | ubuntu-24.04 | ap-southeast-2 |
| `v1.34.7` | ubuntu-22.04 | ap-southeast-2 |
| `v1.34.6` | ubuntu-24.04 | ap-southeast-2 |
| `v1.34.6` | ubuntu-22.04 | ap-southeast-2 |
| `v1.34.5` | ubuntu-24.04 | ap-southeast-2 |
| `v1.34.5` | ubuntu-22.04 | ap-southeast-2 |
| `v1.34.4` | ubuntu-24.04 | ap-southeast-2 |
| `v1.34.4` | ubuntu-22.04 | ap-southeast-2 |
| `v1.34.3` | ubuntu-24.04 | ap-southeast-2 |
| `v1.34.3` | ubuntu-22.04 | ap-southeast-2 |
| `v1.34.2` | ubuntu-24.04 | ap-southeast-2 |
| `v1.34.2` | ubuntu-22.04 | ap-southeast-2 |
| `v1.34.1` | ubuntu-24.04 | ap-southeast-2 |
| `v1.34.1` | ubuntu-22.04 | ap-southeast-2 |
| `v1.34.0` | ubuntu-24.04 | ap-southeast-2 |
| `v1.34.0` | ubuntu-22.04 | ap-southeast-2 |
