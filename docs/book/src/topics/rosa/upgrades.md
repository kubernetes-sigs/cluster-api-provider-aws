# Upgrades

## Control Plane Upgrade

Upgrading the OpenShift version of the control plane is supported by the provider. To perform an upgrade you need to update the `version` in the spec of the `ROSAControlPlane`. Once the version has changed the provider will handle the upgrade for you.

The Upgrade state can be checked in the conditions under `ROSAControlPlane.status`.

## MachinePool Upgrade

Upgrading the OpenShift version of the MachinePools is supported by the provider and can be performed independently from the Control Plane upgrades. To perform an upgrade you need to update the `version` in the spec of the `ROSAMachinePool`. Once the version has changed the provider will handle the upgrade for you.

The Upgrade state can be checked in the conditions under `ROSAMachinePool.status`.

The version of the MachinePool can't be greater than Control Plane version.
