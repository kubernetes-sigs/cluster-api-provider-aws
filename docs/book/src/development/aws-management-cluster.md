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
cp .env.template .env
# Set AWS_REGION, AWS_SSH_KEY_NAME, and a unique MANAGEMENT_CLUSTER_NAME.
make bootstrap
```

The management-cluster kubeconfig defaults to
`_artifacts/${MANAGEMENT_CLUSTER_NAME}.kubeconfig`.

To preview the pivot and enable detailed `clusterctl move` logging, use:

```bash
make bootstrap MANAGEMENT_CLUSTER_ARGS="--dry-run --verbose 5"
```

A dry-run retains the bootstrap cluster because the resources remain there.
The same options are supported by `make teardown`; teardown dry-run does not
delete the management cluster.

## Delete the cluster

```bash
make teardown
```

Teardown pivots the objects back to a temporary kind cluster before deleting
the AWS Cluster object. This keeps the controllers available until AWS
infrastructure cleanup finishes. The lifecycle-script README documents all
configuration options.
