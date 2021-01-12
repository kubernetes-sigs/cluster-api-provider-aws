# EKS Cluster Upgrades

## Control Plane Upgrade

Upgrading the Kubernetes version of the control plane is supported by the provider. To perform an upgrade you need to update the `version` in the spec of the `AWSManagedControlPlane`. Once the version has changed the provider will handle the upgrade for you.

You can only upgrade a EKS cluster by 1 minor version at a time. If you attempt to upgrade the version by more then 1 minor version the provider will ensure the upgrade is done in multiple steps of 1 minor version. For example upgrading from v1.15 to v1.17 would result in your cluster being upgraded v1.15 -> v1.16 first and then v1.16 to v1.17.