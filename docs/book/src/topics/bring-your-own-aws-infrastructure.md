# Bring Your Own AWS Infrastructure

Normally, Cluster API will create infrastructure on AWS when standing up a new workload cluster. However, it is possible to have Cluster API re-use external AWS infrastructure instead of creating its own infrastructure. 

There are two possible ways to do this:
* By consuming existing AWS infrastructure
* By using externally managed AWS infrastructure
> **IMPORTANT NOTE**: This externally managed AWS infrastructure should not be confused with EKS-managed clusters.

Follow the instructions below to configure Cluster API to consume existing AWS infrastructure.

## Consuming Existing AWS Infrastructure

### Overview

CAPA supports using existing AWS resources while creating AWS Clusters which gives flexibility to the users to bring their own existing resources into the cluster instead of creating new resources again.

Follow the instructions below to configure Cluster API to consume existing AWS infrastructure.

### Prerequisites

In order to have Cluster API consume existing AWS infrastructure, you will need to have already created the following resources:

* A VPC
* One or more private subnets (subnets that do not have a route to an Internet gateway)
* A NAT gateway for each private subnet, along with associated Elastic IP addresses (only needed if the nodes require access to the Internet, i.e. pulling public images)
  * A public subnet in the same Availability Zone (AZ) for each private subnet (this is required for NAT gateways to function properly)
* An Internet gateway for all public subnets (only required if the workload cluster is set to use an Internet facing load balancer or one or more NAT gateways exist in the VPC)
* Route table associations that provide connectivity to the Internet through a NAT gateway (for private subnets) or the Internet gateway (for public subnets)
* VPC endpoints for `ec2`, `elasticloadbalancing`, `secretsmanager` an `autoscaling` (if using MachinePools) when the private Subnets do not have a NAT gateway

You will need the ID of the VPC and subnet IDs that Cluster API should use. This information is available via the AWS Management Console or the AWS CLI.

Note that there is no need to create an Elastic Load Balancer (ELB), security groups, or EC2 instances; Cluster API will take care of these items.

If you want to use existing security groups, these can be specified and new ones will not be created.

If you want to use an existing control load load balancer, specify its name.

### Tagging AWS Resources

Cluster API itself does tag AWS resources it creates. The `sigs.k8s.io/cluster-api-provider-aws/cluster/<cluster-name>` (where `<cluster-name>` matches the `metadata.name` field of the Cluster object) tag, with a value of `owned`, tells Cluster API that it has ownership of the resource. In this case, Cluster API will modify and manage the lifecycle of the resource.

When consuming existing AWS infrastructure, the Cluster API AWS provider does not require any tags to be present. The absence of the tags on an AWS resource indicates to Cluster API that it should not modify the resource or attempt to manage the lifecycle of the resource.

However, the built-in Kubernetes AWS cloud provider _does_ require certain tags in order to function properly. Specifically, all subnets where Kubernetes nodes reside should have the `kubernetes.io/cluster/<cluster-name>` tag present. Private subnets should also have the `kubernetes.io/role/internal-elb` tag with a value of 1, and public subnets should have the `kubernetes.io/role/elb` tag with a value of 1. These latter two tags help the cloud provider understand which subnets to use when creating load balancers.

Finally, if the controller manager isn't started with the `--configure-cloud-routes: "false"` parameter, the route table(s) will also need the `kubernetes.io/cluster/<cluster-name>` tag. (This parameter can be added by customizing the `KubeadmConfigSpec` object of the `KubeadmControlPlane` object.)

> **Note**: All the tagging of resources should be the responsibility of the users and are not managed by CAPA controllers.

### Configuring the AWSCluster Specification

Specifying existing infrastructure for Cluster API to use takes place in the specification for the AWSCluster object. Specifically, you will need to add an entry with the VPC ID and the IDs of all applicable subnets into the `network` field. Here is an example:

For EC2
```yaml
apiVersion: controlplane.cluster.x-k8s.io/v1beta1
kind: AWSCluster
```
For EKS
```yaml
apiVersion: controlplane.cluster.x-k8s.io/v1beta1
kind: AWSManagedControlPlane
```

```yaml
spec:
  network:
    vpc:
      id: vpc-0425c335226437144
    subnets:
    - id: subnet-0261219d564bb0dc5
    - id: subnet-0fdcccba78668e013
```

When you use `kubectl apply` to apply the Cluster and AWSCluster specifications to the management cluster, Cluster API will use the specified VPC ID and subnet IDs, and will not create a new VPC, new subnets, or other associated resources. It _will_, however, create a new ELB and new security groups.

### Placing EC2 Instances in Specific AZs

To distribute EC2 instances across multiple AZs, you can add information to the Machine specification. This is optional and only necessary if control over AZ placement is desired.

To tell Cluster API that an EC2 instance should be placed in a particular AZ but allow Cluster API to select which subnet in that AZ can be used, add this to the Machine specification:

```yaml
spec:
  failureDomain: "us-west-2a"
```

If using a MachineDeployment, specify AZ placement like so:

```yaml
spec:
  template:
    spec:
      failureDomain: "us-west-2b"
```

Note that all replicas within a MachineDeployment will reside in the same AZ.

### Placing EC2 Instances in Specific Subnets

To specify that an EC2 instance should be placed in a specific subnet, add this to the AWSMachine specification:

```yaml
spec:
  subnet:
    id: subnet-0a3507a5ad2c5c8c3
```

When using MachineDeployments, users can control subnet selection by adding information to the AWSMachineTemplate associated with that MachineDeployment, like this:

```yaml
spec:
  template:
    spec:
      subnet:
        id: subnet-0a3507a5ad2c5c8c3
```

Users may either specify `failureDomain` on the Machine or MachineDeployment objects, _or_ users may explicitly specify subnet IDs on the AWSMachine or AWSMachineTemplate objects. If both are specified, the subnet ID is used and the `failureDomain` is ignored.

### Placing EC2 Instances in Specific External VPCs

CAPA clusters are deployed within a single VPC, but it's possible to place machines that live in external VPCs. For this kind of configuration, we assume that all the VPCs have the ability to communicate, either through external peering, a transit gateway, or some other mechanism already established outside of CAPA. CAPA will not create a tunnel or manage the network configuration for any secondary VPCs.

The AWSMachineTemplate `subnet` field allows specifying filters or specific subnet ids for worker machine placement. If the filters or subnet id is specified in a secondary VPC, CAPA will place the machine in that VPC and subnet.

```yaml
spec:
  template:
    spec:
      subnet:
        filters:
          name: "vpc-id"
          values:
            - "secondary-vpc-id"
      securityGroupOverrides:
        node: sg-04e870a3507a5ad2c5c8c2
        node-eks-additional: sg-04e870a3507a5ad2c5c8c1
```

#### Caveats/Notes

CAPA helpfully creates security groups for various roles in the cluster and automatically attaches them to workers. However, security groups are tied to a specific VPC, so workers placed in a VPC outside of the cluster will need to have these security groups created by some external process first and set in the `securityGroupOverrides` field, otherwise the ec2 creation will fail.

### Security Groups

To use existing security groups for instances for a cluster, add this to the AWSCluster specification:

```yaml
spec:
  network:
    securityGroupOverrides:
      bastion: sg-0350a3507a5ad2c5c8c3
      controlplane: sg-0350a3507a5ad2c5c8c3
      apiserver-lb: sg-0200a3507a5ad2c5c8c3
      node: sg-04e870a3507a5ad2c5c8c3
      lb: sg-00a3507a5ad2c5c8c3
```

Any additional security groups specified in an AWSMachineTemplate will be applied in addition to these overriden security groups.

To specify additional security groups for the control plane load balancer for a cluster, add this to the AWSCluster specification:

```yaml
spec:
  controlPlaneLoadBalancer:
    additionalSecurityGroups:
    - sg-0200a3507a5ad2c5c8c3
    - ...
```

It's also possible to override the cluster security groups for an individual AWSMachine or AWSMachineTemplate:

```yaml
spec:
  SecurityGroupOverrides:
    node: sg-04e870a3507a5ad2c5c8c2
    node-eks-additional: sg-04e870a3507a5ad2c5c8c1
```

### Control Plane Load Balancer

The cluster control plane is accessed through a Classic ELB. By default, Cluster API creates the Classic ELB. To use an existing Classic ELB, add its name to the AWSCluster specification:

```yaml
spec:
  controlPlaneLoadBalancer:
    name: my-classic-elb-name
```

As control plane instances are added or removed, Cluster API will register and deregister them, respectively, with the Classic ELB.

It's also possible to specify custom ingress rules for the control plane load balancer. To do so, add this to the AWSCluster specification:

```yaml
spec:
  controlPlaneLoadBalancer:
    ingressRules:
      - description: "example ingress rule"
        protocol: "-1" # all
        fromPort: 7777
        toPort: 7777
```

> **WARNING:** Using an existing Classic ELB is an advanced feature. **If you use an existing Classic ELB, you must correctly configure it, and attach subnets to it.**
> 
>An incorrectly configured Classic ELB can easily lead to a non-functional cluster. We strongly recommend you let Cluster API create the Classic ELB.

### Control Plane ingress rules

It's possible to specify custom ingress rules for the control plane itself. To do so, add this to the AWSCluster specification:

```yaml
spec:
  network:
    additionalControlPlaneIngressRules:
    - description: "example ingress rule"
      protocol: "-1" # all
      fromPort: 7777
      toPort: 7777
```
### Caveats/Notes

* When both public and private subnets are available in an AZ, CAPI will choose the private subnet in the AZ over the public subnet for placing EC2 instances.
* If you configure CAPI to use existing infrastructure as outlined above, CAPI will _not_ create an SSH bastion host. Combined with the previous bullet, this means you must make sure you have established some form of connectivity to the instances that CAPI will create.

## Using Externally managed AWS Clusters

### Overview

Alternatively, CAPA supports externally managed cluster infrastructure which is useful for scenarios where a different persona is managing the cluster infrastructure out-of-band(external system) while still wanting to use CAPI for automated machine management.
Users can make use of existing AWSCluster CRDs in their externally managed clusters.

### How to use externally managed clusters?

Users have to use `cluster.x-k8s.io/managed-by: "<name-of-system>"` annotation to depict that AWS resources are managed externally. If CAPA controllers come across this annotation in any of the AWS resources while reconciliation, then it will ignore the resource and not perform any reconciliation(including creating/modifying any of the AWS resources, or it's status).

A predicate `ResourceIsNotExternallyManaged` is exposed by Cluster API which allows CAPA controllers to differentiate between externally managed vs CAPA managed resources. For example:
```go
c, err := ctrl.NewControllerManagedBy(mgr).
        For(&providerv1.InfraCluster{}).
        Watches(...).
        WithOptions(options).
        WithEventFilter(predicates.ResourceIsNotExternallyManaged(mgr.GetScheme(),logger.FromContext(ctx))).
        Build(r)
if err != nil {
	return errors.Wrap(err, "failed setting up with a controller manager")
}
```
The external system must provide all required fields within the spec of the AWSCluster and must adhere to the CAPI provider contract and set the AWSCluster status to be ready when it is appropriate to do so.

> **IMPORTANT NOTE**: Users should take care of skipping reconciliation in external controllers within mapping function while enqueuing requests. For example:
> ```go
> err := c.Watch(
>		&source.Kind{Type: &infrav1.AWSCluster{}},
>		handler.EnqueueRequestsFromMapFunc(func(a client.Object) []reconcile.Request {
>		   if annotations.IsExternallyManaged(awsCluster) {
>			    log.Info("AWSCluster is externally managed, skipping mapping.")
>			    return nil
>		   }
>            return []reconcile.Request{
>	           {
>		         NamespacedName: client.ObjectKey{Namespace: c.Namespace, Name: c.Spec.InfrastructureRef.Name},
>	           },
>	        }}))
> if err != nil {
>    // handle it
> }
> ```

### Caveats
Once the user has created externally managed AWSCluster, it is not allowed to convert it to CAPA managed cluster. However, converting from managed to externally managed is allowed.

User should only use this feature if their cluster infrastructure lifecycle management has constraints that the reference implementation does not support. See [user stories](https://github.com/kubernetes-sigs/cluster-api/blob/10d89ceca938e4d3d94a1d1c2b60515bcdf39829/docs/proposals/20210203-externally-managed-cluster-infrastructure.md#user-stories) for more details.


## Bring your own (BYO) Public IPv4 addresses

Cluster API also provides a mechanism to allocate Elastic IP from the existing Public IPv4 Pool that you brought to AWS[1].

Bringing your own Public IPv4 Pool (BYOIPv4) can be used as an alternative to buying Public IPs from AWS, also considering the changes in charging for this since February 2024[2].

Supported resources to BYO Public IPv4 Pool (`BYO Public IPv4`):
- NAT Gateways
- Network Load Balancer for API server
- Machines

Use `BYO Public IPv4` when you have brought to AWS custom IPv4 CIDR blocks and want the cluster to automatically use IPs from the custom pool instead of Amazon-provided pools.

### Prerequisites and limitations for BYO Public IPv4 Pool

- BYOIPv4 is limited to AWS to selected regions. See more in [AWS Documentation for Regional availability](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-byoip.html#byoip-reg-avail)
- The IPv4 address must be provisioned and advertised to the AWS account before the cluster is installed
- The public IPv4 addresses is limited to the network border group that the CIDR block have been advertised[3][4], and the `NetworkSpec.ElasticIpPool.PublicIpv4Pool` must be the same of the cluster will be installed.
- Only NAT Gateways and the Network Load Balancer for API server will consume from the IPv4 pool defined in the network scope.
- The public IPv4 pool must be assigned to each machine to consume public IPv4 from a custom IPv4 pool.

### Steps to set BYO Public IPv4 Pool to core infrastructure

Currently, CAPA supports BYO Public IPv4 to core components NAT Gateways and Network Load Balancer for the internet-facing API server.

To specify a Public IPv4 Pool for core components you must set the `spec.elasticIpPool` as follows:

```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
kind: AWSCluster
metadata:
  name: aws-cluster-localzone
spec:
  region: us-east-1
  networkSpec:
    vpc:
      elasticIpPool:
        publicIpv4Pool: ipv4pool-ec2-0123456789abcdef0
        publicIpv4PoolFallbackOrder: amazon-pool
```

Then all the Elastic IPs will be created by consuming from the pool `ipv4pool-ec2-0123456789abcdef0`.

### Steps to BYO Public IPv4 Pool to machines

To create a machine consuming from a custom Public IPv4 Pool you must set the pool ID to the AWSMachine spec, then set the `PublicIP` to `true`:

```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
kind: AWSMachine
metadata:
  name: byoip-s55p4-bootstrap
spec:
  # placeholder for AWSMachine spec
  elasticIpPool:
    publicIpv4Pool: ipv4pool-ec2-0123456789abcdef0
    publicIpv4PoolFallbackOrder: amazon-pool
  publicIP: true
```

[1] https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-byoip.html
[2] https://aws.amazon.com/blogs/aws/new-aws-public-ipv4-address-charge-public-ip-insights/
[3] https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-byoip.html#byoip-onboard
[4] https://docs.aws.amazon.com/cli/latest/reference/ec2/advertise-byoip-cidr.html
