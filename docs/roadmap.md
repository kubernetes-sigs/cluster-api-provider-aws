# Cluster API Provider AWS Roadmap

This roadmap is a constant work in progress, subject to frequent revision. Dates are approximations.

## Ongoing

- Documentation improvements

## v0.5.x (v1alpha3+) - June/July 2020
- [Consume Cluster API e2e testing framework](https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/1435)
- [Supporting multiple AWS accounts as Kubernetes objects](https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/1552)
- [Bootstrap failure detection](https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/972)
- [Improved status conditions](https://github.com/kubernetes-sigs/cluster-api/issues/1658)
- [Spot instances support](https://github.com/kubernetes-sigs/cluster-api/issues/1876)

## v0.6 (v1alpha4) ~ Q4 2020

- [Multiple topologies](https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/1484)
- [Machine load balancer implementation](https://github.com/kubernetes-sigs/cluster-api/issues/1250)
- SNS/SQS-based updates from CloudWatch Events so we don’t have to poll AWS for updates

## TBD

- Implement MachinePools - Autoscaling groups and instances
- Dual stack IPv4/IPv6 support
- Windows nodes
- Support for GPU instances and Elastic GPU
- FIPS/NIST/STIG compliance
