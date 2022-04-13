# Cluster API Provider AWS Roadmap

This roadmap is a constant work in progress, subject to frequent revision. Dates are approximations.

## v1.5.x (v1beta1) - April/May 2022
- [Network load balancer support](https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/3088)
- [Graduating EventBridge experimental feature](https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/3414)
- [EFS CSI driver support](https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/3384)
- [AWSManagedMachinePool - Launch Template support](https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/2055)

## v1.6.x (v1beta1) - June/July 2022

- [Spot instance support for AWSMachinePools](https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/2523)
- [Node draining support for AWSMachinePools](https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/2574)
- [IPv6 Support](https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/2420)
- [Security group customization support](https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/392)

## v2.0.x (v1beta2) - End of 2022

- [Support for multiple topologies](https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/1484)

## TBD
- [AWS Fault injector integration to improve resiliency](https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/2173)
- AWSMachinePool implementation backed by Spot Fleet and EC2 Fleet
- [Dual stack IPv4/IPv6 support](https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/3381)
- Windows nodes
- FIPS/NIST/STIG compliance
- Workload identity support to CAPA-managed clusters
- [Use ACK/CrossPlane as backend for AWS SDK calls](https://github.com/kubernetes-sigs/cluster-api-provider-aws/discussions/3306)
- Karpenter support
- [Draining resources created by CCM/CSI like LBs, SGs](https://github.com/kubernetes-sigs/cluster-api/issues/3075)
- [OpenTelemetry integration](https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/2178)
