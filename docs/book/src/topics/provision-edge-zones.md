# Manage Local Zone subnets

## Overview

[AWS Local Zones](https://aws.amazon.com/about-aws/global-infrastructure/localzones/)
extends the cloud infrastructure to metropolitan region,
allowing to deliver applications closer to the end-users, decreasing the
network latency.

CAPA provides the option to manage network resources required to provision compute nodes
to Local Zone locations.

## Requirements and defaults

- Subnets in AWS Local Zones is _not_ created by default.
- When you choose to CAPA manage AWS Local Zones, you also must specify the
  regular zones (Availability Zones) you will create the cluster.
- The subnets in AWS Local Zones will never be used by CAPA to create NAT Gateways,
  Network Load Balancers, and provision Control Plane or Compute nodes by default.
- Nat Gateways is not globaly available to AWS Local Zones locations, the CAPA will use
  the Parent Zone for the Local Zone to create the NAT Gateway to allow the instances on
  private subnets to egress traffic to the internet.
- The subnets holds the attribute of ZoneType and the ParentZoneName where the subnet is created,
  those fields are used to ensure subnets, for example: only subnets with `ZoneType` with
  value `availability-zone` will be used to create Network Load Balancers for control planes.
- It is required to opt-in every zone group of each AWS Local Zones you are planning to create subnets.
    - To check the zone group name for a Local Zone, you can use the [EC2 API `DescribeAvailabilityZones`](describe-availability-zones):
```sh
aws --region "<value_of_AWS_Region>" ec2 describe-availability-zones \
    --query 'AvailabilityZones[].[{ZoneName: ZoneName, GroupName: GroupName, Status: OptInStatus}]' \
    --filters Name=zone-type,Values=local-zone \
    --all-availability-zones
```

    - To opt-int the zone group, you can use the [EC2 API `ModifyZoneAttributes`](https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_ModifyAvailabilityZoneGroup.html):
```sh
aws ec2 modify-availability-zone-group \
    --group-name "<value_of_GroupName>" \
    --opt-in-status opted-in
```

## Creating the VPC subnets with Local Zones

To create a cluster with AWS Local Zones, add the `Subnets` stanza to your `AWSCluster.NetworkSpec`.

```yaml
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
kind: AWSCluster
metadata:
  name: test-aws-cluster
spec:
  region: us-east-1
  sshKeyName: local-zone
  networkSpec:
    vpc:
      cidrBlock: "10.0.0.0/20"
    subnets:
    # regular zones (availability zones)
    - availabilityZone: us-east-1a
      cidrBlock: "10.0.0.0/24"
      id: "cluster-subnet-private-us-east-1a"
      isPublic: false
    - availabilityZone: us-east-1a
      cidrBlock: "10.0.1.0/24"
      id: "cluster-subnet-public-us-east-1a"
      isPublic: true
    - availabilityZone: us-east-1b
      cidrBlock: "10.0.3.0/24"
      id: "cluster-subnet-private-us-east-1b"
      isPublic: false
    - availabilityZone: us-east-1b
      cidrBlock: "10.0.4.0/24"
      id: "cluster-subnet-public-us-east-1b"
      isPublic: true
    - availabilityZone: us-east-1c
      cidrBlock: "10.0.5.0/24"
      id: "cluster-subnet-private-us-east-1c"
      isPublic: false
    - availabilityZone: us-east-1c
      cidrBlock: "10.0.6.0/24"
      id: "cluster-subnet-public-us-east-1c"
      isPublic: true
    # Subnets in Local Zones of New York location (public and private)
    - availabilityZone: us-east-1-nyc-1a
      cidrBlock: "10.0.128.0/25"
      id: "cluster-subnet-private-us-east-1-nyc-1a"
      isPublic: false
    - availabilityZone: us-east-1-nyc-1a
      cidrBlock: "10.0.128.128/25"
      id: "cluster-subnet-public-us-east-1-nyc-1a"
      isPublic: true
```
