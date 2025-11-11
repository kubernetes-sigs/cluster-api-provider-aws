# EKS Cluster Upgrades

## Control Plane Upgrade

Upgrading the Kubernetes version of the control plane is supported by the provider. To perform an upgrade you need to update the `version` in the spec of the `AWSManagedControlPlane`. Once the version has changed the provider will handle the upgrade for you.

You can only upgrade a EKS cluster by 1 minor version at a time. If you attempt to upgrade the version by more then 1 minor version the provider will ensure the upgrade is done in multiple steps of 1 minor version. For example upgrading from v1.15 to v1.17 would result in your cluster being upgraded v1.15 -> v1.16 first and then v1.16 to v1.17.

## Upgrading Nodes from AL2 (EKSConfig) to AL2023 (NodeadmConfig)

Amazon Linux 2 (AL2) AMIs are only supported up to Kubernetes v1.32. To upgrade cluster nodes to v1.33 or newer, you **must** migrate them to Amazon Linux 2023 (AL2023) AMIs. This migration also requires changing the bootstrap provider from `EKSConfig` to the new `NodeadmConfig`.

The upgrade process follows the standard Cluster API rolling update strategy. You will create a new bootstrap template (using `NodeadmConfigTemplate`) and update your `MachineDeployment` or `MachinePool` to reference it, along with the new Kubernetes version and an AL2023-based AMI.

### MachineDeployment Upgrade Example

Here is an example of upgrading a `MachineDeployment` from Kubernetes v1.32 (using `EKSConfig`) to v1.33 (using `NodeadmConfig`).

**Before (v1.32 with `EKSConfigTemplate`):**

```yaml
apiVersion: bootstrap.cluster.x-k8s.io/v1beta2
kind: EKSConfigTemplate
metadata:
  name: default132
spec:
  template:
    spec:
      postBootstrapCommands:
        - "echo \"bye world\""
---
apiVersion: cluster.x-k8s.io/v1beta1
kind: MachineDeployment
metadata:
  name: default
spec:
  clusterName: default
  template:
    spec:
      bootstrap:
        configRef:
          apiVersion: bootstrap.cluster.x-k8s.io/v1beta2
          kind: EKSConfigTemplate
          name: default132
      infrastructureRef:
        kind: AWSMachineTemplate
        name: default132
      version: v1.32.0
```

After (v1.33 with NodeadmConfigTemplate):

A new NodeadmConfigTemplate is created, and the MachineDeployment is updated to reference it and the new version.
YAML

```yaml
apiVersion: bootstrap.cluster.x-k8s.io/v1beta2
kind: NodeadmConfigTemplate
metadata:
  name: default
spec:
  template:
    spec:
      preNodeadmCommands:
        - "echo \"hello world\""
---
apiVersion: cluster.x-k8s.io/v1beta1
kind: MachineDeployment
metadata:
  name: default
spec:
  clusterName: default
  template:
    spec:
      bootstrap:
        configRef:
          apiVersion: bootstrap.cluster.x-k8s.io/v1beta2
          kind: NodeadmConfigTemplate
          name: default
      infrastructureRef:
        kind: AWSMachineTemplate
        name: default
      version: v1.33.0
```
