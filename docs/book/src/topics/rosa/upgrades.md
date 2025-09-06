# Upgrades

## Control Plane Upgrade

Upgrading the OpenShift version of the control plane is supported by the provider. To perform an upgrade you need to update the `version` in the spec of the `ROSAControlPlane`. Once the version has changed the provider will handle the upgrade for you.

Upgrading y-stream version ex; v4.16.x to v4.17.x required the version gate acknowledgement. By default the versionGate is set to WaitForAcknowledge in the `ROSAControlPlane` CR. When upgrading to y-stream version the versionGate should be set to Acknowledge or AlwaysAcknowledge.

##### Note:
When the versionGate is set to 'Acknowledge', it will revert to 'WaitForAcknowledge' once the upgrade is successfully completed. However, if the versionGate is set to 'AlwaysAcknowledge', it will remain set to 'AlwaysAcknowledge' after the upgrade is successfully completed.

The available upgrades versions for the `ROSAControlPlane` will be listed under `ROSAControlPlane.status.availableUpgrades`

The version channel group `ROSAControlPlane.spec.channelGroup` defaults to stable. However, it can be set to eus, fast, candidate, or nightly. Changing the version channel group will change the `ROSAControlPlane.status.availableUpgrades` accordingly. Note that the use of channel groups other than stable may require additional permissions.

The Upgrade state can be checked in the conditions under `ROSAControlPlane.status`.

## MachinePool Upgrade

Upgrading the OpenShift version of the MachinePools is supported by the provider and can be performed independently from the Control Plane upgrades. To perform an upgrade you need to update the `version` in the spec of the `ROSAMachinePool`. Once the version has changed the provider will handle the upgrade for you.

The available upgrades versions for the `ROSAMachinePool` will be listed under `ROSAMachinePool.status.availableUpgrades`

The Upgrade state can be checked in the conditions under `ROSAMachinePool.status`.

The version of the ROSAMachinePool can't be greater than its ROSAControlPlane version.
