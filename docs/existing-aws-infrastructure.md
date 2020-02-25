# Consuming Existing AWS Infrastructure

Normally, Cluster API will create infrastructure on AWS when standing up a new workload cluster. However, it is possible to have Cluster API re-use existing AWS infrastructure instead of creating its own infrastructure. Follow the instructions below to configure Cluster API to consume existing AWS infrastructure.

## Prerequisites

In order to have Cluster API consume existing AWS infrastructure, you will need to have already created the following resources:

* A VPC
* One or more private subnets (subnets that do not automatically assign a public IP address to EC2 instances)
* A public subnet in the same Availability Zone (AZ) for each private subnet (this is required for NAT gateways to function properly)
* A NAT gateway for each private subnet, along with associated Elastic IP addresses
* An Internet gateway for all public subnets
* Route table associations that provide connectivity to the Internet through a NAT gateway (for private subnets) or the Internet gateway (for public subnets)

You will need the following information, which can be obtained either via the AWS Management Console or using the AWS CLI:

1. The ID of the VPC
2. The IDs of the subnets in the VPC that you want Cluster API to use

Note that there is no need to create an Elastic Load Balancer (ELB), security groups, or EC2 instances; Cluster API will take care of these items.

## Tagging AWS Resources

In order for Cluster API to properly recognize and consume existing AWS resources, these resources must be tagged with a specific set of tags. The list below provides the tags that must be present; in this list, `<cluster-name>` refers to the name of cluster as specified in the `metadata.name` field of the Cluster object (or manifest). Here are the tags that are required:

    sigs.k8s.io/cluster-api-provider-aws/cluster/<cluster-name>
    sigs.k8s.io/cluster-api-provider-aws/role

These two tags should be present on the following AWS resources:

* VPC
* All subnets that Cluster API will utilize (private and public)
* Route tables
* Internet gateway
* NAT gateways
* Elastic IPs assigned to NAT gateways

For the `sigs.k8s.io/cluster-api-provider-aws/cluster/<cluster-name>` tag, a value of "owned" is used by Cluster API itself to indicate resources are owned by a specific cluster, and that the lifecycle of the resource is tied to the lifecycle of the cluster. Cluster API recognizes a value of "shared" to indicate a resource may be shared between multiple clusters, and the lifecycle of said resource is not tied to any particular cluster. In this situation---where AWS infrastructure is _not_ being managed by Cluster API---the value should be something **other** than "owned". Using "shared" is acceptable.

For the `sigs.k8s.io/cluster-api-provider-aws/role` tag, Cluster API uses a variety of values to indicate resources are associated with a particular role. Here are the tag values that should be used:

* VPC, route tables, Internet gateway, and NAT gateways should use the value "common".
* Public subnets should use the value "public".
* Private subnets should use the value "private".
* Elastic IPs (for use by the NAT gateways) should use the value "apiserver".

Additionally, for proper operation of the AWS cloud provider, all subnets where Kubernetes nodes reside should also have the `kubernetes.io/cluster/<cluster-name>` tag present.

## Configuring the AWSCluster Specification

Specifying existing infrastructure for Cluster API to use takes place in the specification for the AWSCluster object. Specifically, you will need to add an entry with the VPC ID and the IDs of all applicable subnets into the `networkSpec` field. Here is an example:

```yaml
spec:
  networkSpec:
    vpc:
      id: vpc-0425c335226437144
    subnets:
      - id: subnet-07758a9bc904d06af
      - id: subnet-0a3507a5ad2c5c8c3
      - id: subnet-02ad6429dd0532452
      - id: subnet-02b300779e9d895cf
      - id: subnet-03d8f353b289b025f
      - id: subnet-0a2fe03b0d88fa078
```

When you use `kubectl apply` to apply the Cluster and AWSCluster specifications to the management cluster, Cluster API will use the specified VPC ID and specified subnet IDs, and will not create a new VPC, new subnets, or other associated resources. It _will_, however, create a new ELB and new security groups.

## Configuring the AWSMachine Specification (Optional)

To distribute EC2 instances across different subnets or different AZs, you can add information to the AWSMachine specification. This is optional; without it, Cluster API will deploy to the first private subnet it discovers.

To specify that an EC2 instance should be placed in a specific subnet, add this to the AWSMachine specification:

```yaml
spec:
  subnet:
    id: subnet-0a3507a5ad2c5c8c3
```

To tell Cluster API that an EC2 instance should be placed in a particular AZ but allow Cluster API to select which subnet in that AZ can be used, add this to the AWSMachine specification:

```yaml
spec:
  availabilityZone: "us-west-2a"
```

These two settings are mutually exclusive. You may use only one or the other, not both.

## Caveats/Notes

* When both public and private subnets are available in an AZ, CAPI will choose the private subnet in the AZ over the public subnet for placing EC2 instances.
* If you configure CAPI to use existing infrastructure as outlined above, CAPI will _not_ create an SSH bastion host. Combined with the previous bullet, this means you must make sure you have established some form of connectivity to the instances that CAPI will create.
