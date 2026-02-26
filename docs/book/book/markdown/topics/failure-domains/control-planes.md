# Failure domains in control-plane nodes

By default, the control plane of a workload cluster created by CAPA will span multiple availability zones (AZs) (also referred to as "failure domains") when using multiple control plane nodes. This is because CAPA will, by default, create public and private subnets in all the AZs of a region (up to a maximum of 3 AZs by default). If a region has more than 3 AZs then CAPA will pick 3 AZs to use.

## Configuring CAPA to Use Specific AZs

The Cluster API controller will look for the **FailureDomain** status field and will set the **FailureDomain** field in a `Machine` if a value hasn't already been explicitly set. It will try to ensure that the machines are spread across all the failure domains.

The `AWSMachine` controller looks for a failure domain (i.e. Availability Zone) first in the `Machine` before checking in the `network` specification of `AWSMachine`. This failure domain is then used when provisioning the `AWSMachine`.

### Using FailureDomain in Machine/MachineDeployment spec

To control the placement of `AWSMachine` into a failure domain (i.e. Availability Zones), we can explicitly state the failure domain in `Machine`. The best way is to specify this using the **FailureDomain** field within the `Machine` (or `MachineDeployment`) spec.

For example:

```yaml
apiVersion: cluster.x-k8s.io/v1beta1
kind: Machine
metadata:
  labels:
    cluster.x-k8s.io/cluster-name: my-cluster
    cluster.x-k8s.io/control-plane: "true"
  name: controlplane-0
  namespace: default
spec:
  version: "v1.22.1"
  clusterName: my-cluster
  failureDomain: "1"
  bootstrap:
    configRef:
        apiVersion: bootstrap.cluster.x-k8s.io/v1beta1
        kind: KubeadmConfigTemplate
        name: my-cluster-md-0
  infrastructureRef:
    apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
    kind: AWSMachineTemplate
    name: my-cluster-md-0
```
>**IMPORTANT WARNING:** All the replicas within a `MachineDeployment` will reside in the same Availability Zone.

### Using FailureDomain in network object of AWSMachine

Another way to explicitly instruct CAPA to create resources in specific AZs (and not by random), users can add a `network` object to the AWSCluster specification. Here is an example `network` that creates resources across three AZs in the "us-west-2" region:

```yaml
spec:
  network:
    vpc:
      cidrBlock: 10.50.0.0/16
    subnets:
    - availabilityZone: us-west-2a
      cidrBlock: 10.50.0.0/20
      isPublic: true
    - availabilityZone: us-west-2a
      cidrBlock: 10.50.16.0/20
    - availabilityZone: us-west-2b
      cidrBlock: 10.50.32.0/20
      isPublic: true
    - availabilityZone: us-west-2b
      cidrBlock: 10.50.48.0/20
    - availabilityZone: us-west-2c
      cidrBlock: 10.50.64.0/20
      isPublic: true
    - availabilityZone: us-west-2c
      cidrBlock: 10.50.80.0/20
```

> Note: This method can also be used with worker nodes as well.

Specifying the CIDR block alone for the VPC is not enough; users must also supply a list of subnets that provides the desired AZ, the CIDR for the subnet, and whether the subnet is public (has a route to an Internet gateway) or is private (does not have a route to an Internet gateway).

Note that CAPA insists that there must be a public subnet (and associated Internet gateway), even if no public load balancer is requested for the control plane. Therefore, for every AZ where a control plane node should be placed, the `network` object must define both a public and private subnet.

Once CAPA is provided with a `network` that spans multiple AZs, the KubeadmControlPlane controller will automatically distribute control plane nodes across multiple AZs. No further configuration from the user is required.

> Note: This method can also be used if you do not want to split your EC2 instances across multiple AZs.

## Changing AZ defaults

When creating default subnets by default a maximum of 3 AZs will be used. If you are creating a cluster in a region that has more than 3 AZs then 3 AZs will be picked based on alphabetical from that region.

If this default behavior for maximum number of AZs and ordered selection method doesn't suit your requirements you can use the following to change the behaviour:

* `availabilityZoneUsageLimit` - specifies the maximum number of availability zones (AZ) that should be used in a region when automatically creating subnets.
* `availabilityZoneSelection` - specifies how AZs should be selected if there are more AZs in a region than specified by availabilityZoneUsageLimit. There are 2 selection schemes:
  * `Ordered` - selects based on alphabetical order
  * `Random` - selects AZs randomly in a region

For example if you wanted have a maximum of 2 AZs using a random selection scheme:

```yaml
spec:
  network:
    vpc:
      availabilityZoneUsageLimit: 2
      availabilityZoneSelection: Random
```

## Caveats

Deploying control plane nodes across multiple AZs is not a panacea to cure all availability concerns. The sizing and overall utilization of the cluster will greatly affect the behavior of the cluster and the workloads hosted there in the event of an AZ failure. Careful planning is needed to maximize the availability of the cluster even in the face of an AZ failure. There are also other considerations, like cross-AZ traffic charges, that should be taken into account.
