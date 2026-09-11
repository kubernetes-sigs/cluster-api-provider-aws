# Test a self-hosted management cluster

The unmanaged end-to-end suite can run with a self-managed CAPA cluster in AWS
instead of using kind as its management cluster for the entire run. This uses
the same CAPA templates and provisioning path as the AWS workload clusters
created by the suite.

The e2e framework:

1. creates a CAPA workload cluster in AWS using the
   `remote-management-cluster` flavor;
2. installs the providers under test in the workload cluster;
3. moves the Cluster API objects from kind to the workload cluster with
   `clusterctl move`;
4. runs the selected e2e specs with the AWS cluster as their management cluster;
5. moves the management-cluster objects back to kind and deletes both clusters.

The mode is opt-in. The normal e2e workflow continues to use kind as its
management cluster.

## Prerequisites

- AWS credentials
- The standard dependencies required by `make test-e2e`
- Enough AWS quota for the control-plane, worker, and networking resources

## Provision a self-hosted management cluster

```bash
make provision-self-hosted-management-cluster
```

The e2e framework initially creates its normal kind management cluster. It then
uses the existing `remote-management-cluster` flavor to provision an AWS
workload cluster, installs the providers under test in it, and moves the Cluster
API objects into the AWS cluster. It then skips all e2e specs and exits without
running suite cleanup.

Both the AWS management cluster and the temporary kind bootstrap cluster remain
running. They create chargeable resources and must be cleaned up explicitly
after manual testing.

## Run e2e tests

```bash
make test-self-hosted-management-cluster
```

This uses `_artifacts/self-hosted-management-cluster.kubeconfig` in
existing-cluster mode. The e2e specs run independently of provisioning, and
suite cleanup leaves both management clusters available for the teardown
stage.

## Tear down the environment

```bash
make teardown-self-hosted-management-cluster
```

Teardown reads the lifecycle state from the artifact directory, moves the CAPA
objects back to kind, deletes the AWS cluster, and then deletes kind.
