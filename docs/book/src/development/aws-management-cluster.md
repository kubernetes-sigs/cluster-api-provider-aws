# AWS-hosted management cluster

For tests that require the management cluster itself to run in AWS, the
repository provides a CAPI-native bootstrap flow. It does not require `eksctl`.

The flow creates a temporary kind cluster, initializes CAPI and CAPA, applies
an AWS cluster definition, and waits for the new cluster. It then initializes
the providers in the new cluster, pivots the Cluster API resources with
`clusterctl move`, and removes kind. The AWS cluster is then self-managing.

## Prerequisites

- AWS credentials and an existing EC2 SSH key pair
- Docker, `kind`, `kubectl`, `clusterctl`, and `clusterawsadm`
- Enough AWS quota for the control-plane, worker, and networking resources

## Create the cluster

```bash
make bootstrap
```

The management-cluster kubeconfig defaults to `$(ARTIFACTS)/self-hosted-management-cluster.kubeconfig`
(where `ARTIFACTS` defaults to `_artifacts/`).

## Delete the cluster

```bash
make teardown
```

Teardown pivots the objects back to a temporary kind cluster before deleting
the AWS Cluster object. This keeps the controllers available until AWS
infrastructure cleanup finishes. The `make teardown` command must use the same
`ARTIFACTS` folder as `make bootstrap`.

Configuration is controlled by environment variables used by `make test-e2e`:
- `AWS_REGION` — AWS region (defaults to `us-east-1`)
- `AWS_SSH_KEY_NAME` — existing EC2 SSH key pair name
- `ARTIFACTS` — output directory (defaults to `_artifacts/`)
- Other e2e configuration variables defined in the e2e test framework

The underlying e2e flags (`-provision-self-hosted-management-cluster` and
`-teardown-self-hosted-management-cluster`) are defined in
`test/e2e/shared/defaults.go`.
