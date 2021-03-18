# Consuming Existing AWS Infrastructure

Normally, Cluster API will create infrastructure on AWS when standing up a new workload cluster. However, it is possible to have Cluster API re-use existing AWS infrastructure instead of creating its own infrastructure. Follow the instructions below to configure Cluster API to consume existing AWS infrastructure.

## Prerequisites

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

## Tagging AWS Resources

Cluster API itself does tag AWS resources it creates. The `sigs.k8s.io/cluster-api-provider-aws/cluster/<cluster-name>` (where `<cluster-name>` matches the `metadata.name` field of the Cluster object) tag, with a value of `owned`, tells Cluster API that it has ownership of the resource. In this case, Cluster API will modify and manage the lifecycle of the resource.

When consuming existing AWS infrastructure, the Cluster API AWS provider does not require any tags to be present. The absence of the tags on an AWS resource indicates to Cluster API that it should not modify the resource or attempt to manage the lifecycle of the resource.

However, the built-in Kubernetes AWS cloud provider _does_ require certain tags in order to function properly. Specifically, all subnets where Kubernetes nodes reside should have the `kubernetes.io/cluster/<cluster-name>` tag present. Private subnets should also have the `kubernetes.io/role/internal-elb` tag with a value of 1, and public subnets should have the `kubernetes.io/role/elb` tag with a value of 1. These latter two tags help the cloud provider understand which subnets to use when creating load balancers.

Finally, if the controller manager isn't started with the `--configure-cloud-routes: "false"` parameter, the route table(s) will also need the `kubernetes.io/cluster/<cluster-name>` tag. (This parameter can be added by customizing the `KubeadmConfigSpec` object of the `KubeadmControlPlane` object.)

## Configuring the AWSCluster Specification

Specifying existing infrastructure for Cluster API to use takes place in the specification for the AWSCluster object. Specifically, you will need to add an entry with the VPC ID and the IDs of all applicable subnets into the `networkSpec` field. Here is an example:

For EC2
```yaml
apiVersion: controlplane.cluster.x-k8s.io/v1alpha3
kind: AWSCluster
```
For EKS
```yaml
apiVersion: controlplane.cluster.x-k8s.io/v1alpha3
kind: AWSManagedControlPlane
```

```yaml
spec:
  networkSpec:
    vpc:
      id: vpc-0425c335226437144
    subnets:
    - id: subnet-0261219d564bb0dc5
    - id: subnet-0fdcccba78668e013
```

When you use `kubectl apply` to apply the Cluster and AWSCluster specifications to the management cluster, Cluster API will use the specified VPC ID and subnet IDs, and will not create a new VPC, new subnets, or other associated resources. It _will_, however, create a new ELB and new security groups.

## Placing EC2 Instances in Specific AZs

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

## Placing EC2 Instances in Specific Subnets

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

## Security Groups

To use existing security groups for instances for a cluster, add this to the AWSCluster specification:

```yaml
spec:
  networkSpec:
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
    AdditionalsecurityGroups:
    - sg-0200a3507a5ad2c5c8c3
    - ...   
```

## Caveats/Notes

* When both public and private subnets are available in an AZ, CAPI will choose the private subnet in the AZ over the public subnet for placing EC2 instances.
* If you configure CAPI to use existing infrastructure as outlined above, CAPI will _not_ create an SSH bastion host. Combined with the previous bullet, this means you must make sure you have established some form of connectivity to the instances that CAPI will create.