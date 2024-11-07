# Manage Local Zone subnets

## Overview

CAPA provides the option to manage network resources required to provision compute nodes
to Local Zone and Wavelength Zone locations.

[AWS Local Zones](https://aws.amazon.com/about-aws/global-infrastructure/localzones/)
extends the cloud infrastructure to metropolitan regions,
allowing to deliver applications closer to the end-users, decreasing the
network latency.

[AWS Wavelength Zones](https://aws.amazon.com/wavelength/)
extends the AWS infrastructure deployments infrastructure to carrier infrastructure,
allowing to deploy within communications service providers’ (CSP) 5G networks.

When "edge zones" is mentioned in this document, it is referencing to AWS Local Zones and AWS Wavelength Zones.

## Requirements and defaults

For both Local Zones and Wavelength Zones ('edge zones'):

- Subnets in edge zones are _not_ created by default.
- When you choose to CAPA manage edge zone's subnets, you also must specify the
  regular zones (Availability Zones) you will create the cluster.
- IPv6 is not globally supported by AWS across Local Zones,
  and is not supported in Wavelength zones, CAPA support is limited to IPv4
  subnets in edge zones.
- The subnets in edge zones will not be used by CAPA to create NAT Gateways,
  Network Load Balancers, or provision Control Plane or Compute nodes by default.
- NAT Gateways are not globally available to edge zone's locations, the CAPA uses
  the Parent Zone for the edge zone to create the NAT Gateway to allow the instances on
  private subnets to egress traffic to the internet.
- The CAPA subnet controllers discovers the zone attributes `ZoneType` and
  `ParentZoneName` for each subnet on creation, those fields are used to ensure subnets for
  it's role. For example: only subnets with `ZoneType` with value `availability-zone`
  can be used to create a load balancer for API.
- It is required to manually opt-in to each zone group for edge zones you are planning to create subnets.

The following steps are example to describe the zones and opt-into an zone group for an Local Zone:

    - To check the zone group name for a Local Zone, you can use the [EC2 API `DescribeAvailabilityZones`][describe-availability-zones]. For example:
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

## Installing managed clusters extending subnets to Local Zones

To create a cluster with support of subnets on AWS Local Zones, add the `Subnets` stanza to your `AWSCluster.NetworkSpec`. Example:

```yaml
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
kind: AWSCluster
metadata:
  name: aws-cluster-localzone
spec:
  region: us-east-1
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

## Installing managed clusters extending subnets to Wavelength Zones

To create a cluster with support of subnets on AWS Wavelength Zones, add the `Subnets` stanza to your `AWSCluster.NetworkSpec`. Example:

```yaml
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
kind: AWSCluster
metadata:
  name: aws-cluster-wavelengthzone
spec:
  region: us-east-1
  networkSpec:
    vpc:
      cidrBlock: "10.0.0.0/20"
    subnets:
    # <placeholder for regular zones (availability zones)>
    - availabilityZone: us-east-1-wl1-was-wlz-1
      cidrBlock: "10.0.128.0/25"
      id: "cluster-subnet-private-us-east-1-wl1-was-wlz-1"
      isPublic: false
    - availabilityZone: us-east-1-wl1-was-wlz-1
      cidrBlock: "10.0.128.128/25"
      id: "cluster-subnet-public-us-east-1-wl1-was-wlz-1"
      isPublic: true
```

## Installing managed clusters extending subnets to Local and Wavelength Zones

It is also possible to mix the creation across both Local and Wavelength zones.

To create a cluster with support of edge zones, add the `Subnets` stanza to your `AWSCluster.NetworkSpec`. Example:

```yaml
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
kind: AWSCluster
metadata:
  name: aws-cluster-edge
spec:
  region: us-east-1
  networkSpec:
    vpc:
      cidrBlock: "10.0.0.0/20"
    subnets:
    # <placeholder for regular zones (availability zones)>
    - availabilityZone: us-east-1-nyc-1a
      cidrBlock: "10.0.128.0/25"
      id: "cluster-subnet-private-us-east-1-nyc-1a"
      isPublic: false
    - availabilityZone: us-east-1-nyc-1a
      cidrBlock: "10.0.128.128/25"
      id: "cluster-subnet-public-us-east-1-nyc-1a"
      isPublic: true
    - availabilityZone: us-east-1-wl1-was-wlz-1
      cidrBlock: "10.0.129.0/25"
      id: "cluster-subnet-private-us-east-1-wl1-was-wlz-1"
      isPublic: false
    - availabilityZone: us-east-1-wl1-was-wlz-1
      cidrBlock: "10.0.129.128/25"
      id: "cluster-subnet-public-us-east-1-wl1-was-wlz-1"
      isPublic: true
```


[describe-availability-zones]: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeAvailabilityZones.html
