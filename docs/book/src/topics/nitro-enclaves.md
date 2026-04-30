# AWS Nitro Enclaves

[AWS Nitro Enclaves](https://aws.amazon.com/ec2/nitro/nitro-enclaves/) are isolated compute environments built on the Nitro hypervisor. They allow you to process highly sensitive data by providing a strongly isolated, highly constrained environment with no persistent storage, no interactive access, and no external networking.

Nitro Enclaves are enabled per-instance via the `enclaveOptions` field on the launch template.

## Prerequisites

- The instance type must support Nitro Enclaves. See the [AWS documentation](https://docs.aws.amazon.com/enclaves/latest/user/nitro-enclave.html#nitro-enclave-reqs) for supported instance types.
- The instance must use a Nitro-based AMI.

## Configuration

Three resources work together: a `MachinePool` (CAPI core) that owns the bootstrap and infrastructure references, an `AWSManagedMachinePool` with `enclaveOptions` enabled on the launch template, and a `NodeadmConfig` that installs the Nitro Enclaves CLI and configures the allocator as pre-nodeadm commands.

```yaml
apiVersion: cluster.x-k8s.io/v1beta1
kind: MachinePool
metadata:
  name: my-node-pool
spec:
  clusterName: my-cluster
  replicas: 1
  template:
    spec:
      bootstrap:
        configRef:
          apiVersion: bootstrap.cluster.x-k8s.io/v1beta2
          kind: NodeadmConfig
          name: my-node-pool-bootstrap
      clusterName: my-cluster
      infrastructureRef:
        apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
        kind: AWSManagedMachinePool
        name: my-node-pool
      version: "${KUBERNETES_VERSION}"
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
kind: AWSManagedMachinePool
metadata:
  name: my-node-pool
spec:
  amiType: CUSTOM
  labels:
    aws-nitro-enclaves-k8s-dp: enabled
  awsLaunchTemplate:
    instanceType: m5.xlarge
    ami:
      eksLookupType: AmazonLinux2023
    enclaveOptions:
      enabled: true
---
apiVersion: bootstrap.cluster.x-k8s.io/v1beta2
kind: NodeadmConfig
metadata:
  name: my-node-pool-bootstrap
spec:
  preNodeadmCommands:
    - dnf install -y aws-nitro-enclaves-cli
    - sed -i "s/cpu_count:.*/cpu_count: 2/" /etc/nitro_enclaves/allocator.yaml
    - sed -i "s/memory_mib:.*/memory_mib: 768/" /etc/nitro_enclaves/allocator.yaml
    - systemctl enable --now nitro-enclaves-allocator.service
```

Adjust `cpu_count` and `memory_mib` to match the resources you intend to allocate to enclaves on each node. These values must be consistent with the hugepage size requested in pod specs (see Step 5 below). The example above uses Amazon Linux 2023 with nodeadm; for Amazon Linux 2 replace the `dnf` command with `amazon-linux-extras install -y aws-nitro-enclaves-cli`.

## Using with EKS

CAPA handles the launch template and node provisioning. The remaining steps to run enclave workloads must be completed after the cluster is ready.

The full workflow is documented in [Using Nitro Enclaves with Amazon EKS](https://docs.aws.amazon.com/enclaves/latest/user/kubernetes.html). The high-level steps are:

**Step 1 — Launch template (CAPA):** The `AWSManagedMachinePool` and `NodeadmConfig` resources in the Configuration section above cover this step. CAPA creates the launch template with enclaves enabled and injects the allocator bootstrap commands via `preNodeadmCommands`.

**Step 2 — Cluster and node provisioning (CAPA):** CAPA creates the EKS cluster and managed node group using the launch template. No additional action required here beyond normal cluster provisioning.

**Step 3 — Install the device plugin:** Deploy the [AWS Nitro Enclaves Kubernetes Device Plugin](https://github.com/aws/aws-nitro-enclaves-k8s-device-plugin) as a DaemonSet. This exposes the `aws.ec2.nitro/nitro_enclaves` resource to the Kubernetes scheduler. The device plugin targets nodes with the `aws-nitro-enclaves-k8s-dp=enabled` label, which is set via `spec.labels` in the configuration above.

**Step 4 — Prepare the application image:** Enclave applications must be packaged as a container image that embeds an Enclave Image File (EIF). The [aws-nitro-enclaves-with-k8s](https://github.com/aws/aws-nitro-enclaves-with-k8s) repository provides the `enclavectl` tool to build these images.

**Step 5 — Deploy the application:** Pod specs must request the enclave device and hugepages. The hugepage size must match the memory allocated in the node user data. See the [AWS documentation](https://docs.aws.amazon.com/enclaves/latest/user/kubernetes.html#deploy-app) for a full deployment example.

```yaml
resources:
  limits:
    aws.ec2.nitro/nitro_enclaves: "1"
    hugepages-2Mi: 768Mi
    cpu: 250m
  requests:
    aws.ec2.nitro/nitro_enclaves: "1"
    hugepages-2Mi: 768Mi
```

Pods also require a `HugePages-2Mi` volume mounted at `/dev/hugepages`.

## Limitations

- **Existing launch templates:** Enabling Nitro Enclaves on an existing launch template triggers a new launch template version and a rolling update of the node group.
- **Hibernation:** Nitro Enclaves cannot be used with hibernation.
- **Local Zones and Wavelength Zones are not supported.** AWS does not support Nitro Enclaves in Local Zones or Wavelength Zones. If `enclaveOptions.enabled: true` is set and the `AWSManagedMachinePool` specifies a subnet ID or availability zone that resolves to an edge zone, the controller will block reconciliation and set the `LaunchTemplateReady` condition to `False` with reason `NitroEnclaveEdgeZoneUnsupported`.

To resolve this, either remove the edge-zone subnet or availability zone from the pool spec, or disable `enclaveOptions`.
