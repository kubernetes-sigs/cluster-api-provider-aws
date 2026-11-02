# Cluster Autoscaler

Every ROSA HCP cluster runs a cluster autoscaler in its hosted control plane. It adds nodes when pods can't be
scheduled and removes nodes that are underused, within the limits of each machine pool.

Two settings work together:

- **Machine pool autoscaling** (`ROSAMachinePool.spec.autoscaling`, or `ROSAControlPlane.spec.defaultMachinePoolSpec.autoscaling`
  for the default machine pool) sets the minimum and maximum number of replicas of each machine pool.
- **Cluster autoscaler configuration** (`ROSAControlPlane.spec.clusterAutoscaler`) sets cluster-wide behaviour, such as
  the maximum number of nodes across all machine pools or how long to wait for a node to be provisioned.

## Configuration

```yaml
---
apiVersion: controlplane.cluster.x-k8s.io/v1beta2
kind: ROSAControlPlane
metadata:
  name: "capi-rosa-quickstart-control-plane"
spec:
  clusterAutoscaler:
    resourceLimits:
      maxNodesTotal: 50
    maxPodGracePeriod: 600
    maxNodeProvisionTime: 15m
    podPriorityThreshold: -10
  ....
```

| Field | Description | Allowed values | Default |
|---|---|---|---|
| `resourceLimits.maxNodesTotal` | Maximum number of worker nodes across all machine pools. | `0` (no limit) or at least `2` | `0` |
| `maxPodGracePeriod` | Maximum time, in seconds, pods get to terminate when a node is removed. | `0` or at least `600` | `600` |
| `maxNodeProvisionTime` | Maximum time to wait for a new node before the scale up is considered failed. | `15m` to `60m` | `15m` |
| `podPriorityThreshold` | Pods with a lower priority don't trigger a scale up. | at most `1000000000` | `-10` |

OCM also checks that the cluster can't grow beyond the size limit of your organization. If it can, the configuration is
rejected and the error is reported in the `ClusterAutoscalerConfigured` condition.

## How the configuration is managed

- **`clusterAutoscaler` not set**: CAPA doesn't change the cluster autoscaler. Its current configuration in OCM is kept,
  including values set with `rosa edit autoscaler` or the OpenShift Cluster Manager console.
- **`clusterAutoscaler` set**: CAPA manages the whole cluster autoscaler configuration. Fields you omit are set to the
  defaults in the table above, and changes made outside of CAPA are reverted on the next reconcile.
- **`clusterAutoscaler: {}`**: resets every field to its default.
- **Removing `clusterAutoscaler`**: CAPA stops managing it and the current values in OCM are kept.

To keep the values a cluster already has before letting CAPA manage it, copy them from `status.clusterAutoscaler`
into `spec.clusterAutoscaler`.

The configuration is sent to OCM once the cluster is ready. While the cluster is being installed, the cluster autoscaler
runs with the default values. OCM applies accepted changes to the hosted cluster asynchronously, so they take effect
shortly after OCM accepts them.

## Status

`status.clusterAutoscaler` always shows the configuration currently stored in OCM once the cluster is ready, whether
or not `spec.clusterAutoscaler` is set:

```yaml
status:
  clusterAutoscaler:
    resourceLimits:
      maxNodesTotal: 50
    maxPodGracePeriod: 600
    maxNodeProvisionTime: 15m0s
    podPriorityThreshold: -10
```

When `spec.clusterAutoscaler` is set, the `ClusterAutoscalerConfigured` condition reports whether OCM accepted it:

| Status | Reason | Meaning |
|---|---|---|
| `True` | | OCM accepted and stores the requested configuration. |
| `False` | `WaitingForCluster` | The cluster isn't ready yet; the configuration will be sent to OCM once it is. |
| `False` | `ValidationFailed` | OCM rejected the configuration. The message contains the OCM error. |
| `False` | `NotAvailable` | The cluster autoscaler API isn't available for this cluster. |
| `False` | `ReconciliationFailed` | A temporary error occurred; the controller retries. |

The condition doesn't affect `ROSAControlPlaneReady`. A rejected cluster autoscaler configuration doesn't make the
control plane unavailable.
