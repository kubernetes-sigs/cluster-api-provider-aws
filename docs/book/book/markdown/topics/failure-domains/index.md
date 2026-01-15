# Failure Domains

A failure domain in the AWS provider corresponds to an availability zone within an AWS region. 

In AWS, Availability Zones are distinct locations within an AWS Region that are engineered to be isolated from failures in other Availability Zones. They provide inexpensive, low-latency network connectivity to other Availability Zones in the same AWS Region, to ensure a cluster (or any application) is resilient to failure. 

If a zone goes down, your cluster will continue to run as the other 2 zones are physically separated and can continue to run.

More details of availability zones and regions can be found in the [AWS docs](https://aws.amazon.com/about-aws/global-infrastructure/regions_az/).

The usage of failure domains for control-plane and worker nodes can be found below in detail:

- [Control Plane](control-planes.md)
- [Worker nodes](worker-nodes.md)