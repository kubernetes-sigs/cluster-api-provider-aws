# Cluster API Provider AWS Roadmap

This roadmap is a constant work in progress, subject to frequent revision. Dates are approximations.

## Ongoing

- Documentation improvements

## v0.5 (v1alpha3) ~ March 2020

- Updates for Cluster API v1alpha3
- Add support for `clusterctl` v2
- [Disable creation of bastion host by default](https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/947)
- [Ability to configure node disks](https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/1403)
- [SSM integration (alternative to bastion host)](https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/1553)
- [Encrypted user-data](https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/1387)
- [Support for cross-zone load balancing](https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/1416)
- [Private API-server load balancer](https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/873)
- Consume Cluster API e2e testing framework

## v0.6 (v1alpha4) ~ July 2020

- [Supporting multiple AWS accounts as Kubernetes objects](https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/1552)
- [Multiple topologies](https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/1484)
- [Bootstrap failure detection](https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/972)
- [Improved status conditions](https://github.com/kubernetes-sigs/cluster-api/issues/1658)
- [Machine load balancer implementation](https://github.com/kubernetes-sigs/cluster-api/issues/1250)
- SNS/SQS-based updates from CloudWatch Events so we don’t have to poll AWS for updates

## v0.6 (v1alpha5? v1beta1?) ~ November 2020

- ?

## TBD

- Implement MachinePools - Autoscaling groups and spot instances
- Dual stack IPv4/IPv6 support
- Windows nodes
- Support for GPU instances and Elastic GPU
- FIPS/NIST/STIG compliance
