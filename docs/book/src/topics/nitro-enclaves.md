# AWS Nitro Enclaves

[AWS Nitro Enclaves](https://aws.amazon.com/ec2/nitro/nitro-enclaves/) are isolated compute environments built on the Nitro hypervisor. They allow you to process highly sensitive data by providing a strongly isolated, highly constrained environment with no persistent storage, no interactive access, and no external networking.

Nitro Enclaves are enabled per-instance via the `enclaveOptions` field on the launch template.

## Prerequisites

- The instance type must support Nitro Enclaves. See the [AWS documentation](https://docs.aws.amazon.com/enclaves/latest/user/nitro-enclave.html#nitro-enclave-reqs) for supported instance types.
- The instance must use a Nitro-based AMI.

## Configuration

Nitro Enclave support is configured via the `enclaveOptions` field in `AWSManagedMachinePool.spec.awsLaunchTemplate`:

```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
kind: AWSManagedMachinePool
metadata:
  name: my-node-pool
spec:
  awsLaunchTemplate:
    name: my-launch-template
    instanceType: m5.xlarge
    enclaveOptions:
      enabled: true
```

## Notes

- Enabling Nitro Enclaves on an existing launch template will trigger a new launch template version and a rolling update of the node group.
- Nitro Enclaves cannot be used with hibernation or with instance types that do not support the Nitro hypervisor.

## Limitations

- **Local Zones and Wavelength Zones are not supported.** AWS does not support Nitro Enclaves in Local Zones or Wavelength Zones. If `enclaveOptions.enabled: true` is set and the `AWSMachinePool` specifies a subnet ID or availability zone that resolves to an edge zone, the controller will block reconciliation and set the `LaunchTemplateReady` condition to `False` with reason `NitroEnclaveEdgeZoneUnsupported`.

  To resolve this, either remove the edge-zone subnet or availability zone from the pool spec, or disable `enclaveOptions`.
