---
title: Networking for ROSA HCP
authors:
  - "@mzazrivec"
reviewers:
  -
creation-date: 2025-02-24
last-updated: 2025-03-21
status: provisional
---

# Networking for ROSA HCP

## Table of Contents

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Glossary](#glossary)
- [Summary](#summary)
- [Motivation](#motivation)
  - [Goals](#goals)
  - [Non-Goals/Future Work](#non-goalsfuture-work)
- [Proposal](#proposal)
  - [User Stories](#user-stories)
    - [Functional Requirements](#functional-requirements)
- [Alternatives](#alternatives)
- [Upgrade Strategy](#upgrade-strategy)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Glossary

Refer to the [Cluster API Book Glossary](https://cluster-api.sigs.k8s.io/reference/glossary.html).

## Summary

This proposal defines implementation of networking infrastructure in CAPA for ROSA Hosted Control Plane.

## Motivation

To be able to provision a new ROSA HCP kubernetes cluster using CAPA, one has to create and setup the underlying network infrastructure first: VPC, public and private subnets, internet gateway, routing tables for both subnets, elastic IP address.

All of the above can be currently provisioned and configured via AWS CLI, AWS Management Console or Terraform. Motivation for this work is to be able to provision and configure the network infrastructure for ROSA HCP using CAPI.

### Goals

1. Implement a namespaced new custom resource `RosaNetwork` representing the networking stack for ROSA HCP.
2. It will be possible to reference the new `RosaNetwork` resource from ROSA control plane resource
3. Implement creation and deletion for the new `RosaNetwork` resource.
4. Support the same networking scenarios as [ROSA CLI](https://github.com/openshift/rosa) using the same embeded AWS CloudFormation template that ROSA CLI uses.

### Non-Goals/Future Work

- Modify current networking code in AWS / EKS clusters.
- Support custom CloudFormation template.

## Proposal

The goal of this proposal is to be able to provision the networking infrastructure required for a ROSA HCP cluster.

[ROSA CLI](https://github.com/openshift/rosa) supports creation of the networking infrastructure for ROSA HCP and uses [AWS CloudFormation](https://aws.amazon.com/cloudformation/) template under the hood. The [CloudFormation template used by ROSA CLI](https://github.com/openshift/rosa/blob/master/cmd/create/network/templates/rosa-quickstart-default-vpc/cloudformation.yaml) allows to specify five parameters: CloudFormation stack name, AZ count or list of availability zones, region and CIDR block for the VPC. The created CloudFormation stack then contains a VPC, public and private subnets (each pair created in separate AZ), internet gateway attached to VPC, elastic IPs, NAT gateways, public and private routing tables and a security group.

Adopting the CloudFormation template used by rosa-cli would mean that CAPA and the `RosaNetwork` custom resource would be relying on a mechanism that is know to work well and any changes or fixes implemented in ROSA CLI would be picked up automatically in CAPA.

In practical terms, implementation of the proposal would mean:
1. A new namespaced custom resource definition `RosaNetwork` in CAPA with five attributes: name, AZ count, list of availability zones, region and CIDR block for VPC. `availabilityZoneCount`, `availabilityZones`, `region` and `cidrBlock` will become the `spec` part of the new `RosaNetwork` type, name of the cloudformation stack will be the same as `metadata.name`. The `availabilityZoneCount` and `availabilityZones` parameters will be mutually exclusive.

   `RosaNetwork` spec example with AZ count specified:
    ```
    kind: RosaNetwork
    metadata:
      name: rosa-network-01
      namespace: default
    spec:
      availabilityZoneCount: 3
      region: us-west-2
      cidrBlock: 10.0.0.0/16
    ```

   `RosaNetwork` spec example with specified availability zones:
    ```
    kind: RosaNetwork
    metadata:
      name: rosa-network-01
      namespace: default
    spec:
      availabilityZones:
      - us-west-2a
      - us-west-2d
      region: us-west-2
      cidrBlock: 10.0.0.0/16
    ```

1. A new reconciler for the new custom resource, implementing creation and deletion. The reconciler will be using an existing [CloudFormation template from ROSA CLI](https://github.com/openshift/rosa/blob/master/cmd/create/network/templates/rosa-quickstart-default-vpc/cloudformation.yaml) and will use [AWS CloudFormation API](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/cloudformation) to [create](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/cloudformation#Client.CreateStack) and [delete](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/cloudformation#Client.DeleteStack) the AWS CloudFormation stack.

   Outputs and resources created in the cloudformation stack will be tracked under `status` of the `RosaNetwork` type. In particular, the `status` will contain the list of public and private subnets and availability zones, grouped together by the availability zones.

   Example:
   ```
   kind: RosaNetwork
   metadata:
     name: rosa-network-01
     namespace: default
   status:
     subnets:
     - availabilityZone: us-west-2a
       publicSubnet: subnet-0d9f28ba991b93514
       privateSubnet: subnet-1d9f28ba992a83514
     - availabilityZone: us-west-2b
       publicSubnet: subnet-2d7f18c09f1b43512
       privateSubnet: subnet-2d7f58c09f1b43512
     - availabilityZone: us-west-2c
       publicSubnet: subnet-1d7e19c0af1c4c57f
       privateSubnet: subnet-7d7e19c0af1f4d57f
   ```

   All resources created in the cloudformation stack will be tracked under `status.resources` array:
   ```
   kind: RosaNetwork
   metadata:
     name: rosa-network-01
     namespace: default
   status:
     resources:
       - resource: NATGateway1
         id:
         status: CREATE_IN_PROGRESS
         reason: Eventual consistency check initiated
       - resource: VPC
         id: vpc-055edf3ebf27f18d6 
         status: CREATE_COMPLETE
         reason:
       - resource: SecurityGroup
         id:
         status: CREATE_IN_PROGRESS
         reason: Resource creation Initiated
   ```
   and will be reflecting the the values coming from [AWS CloudFormation API](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/cloudformation#Client.DescribeStackEvents) (`resource`, `id`, `status` and `reason`).

   `status.conditions` of the `RosaNetwork` resource will be consistent with the CAPA conventions, example of a successful network stack creation:
   ```
   kind: RosaNetwork
   metadata:
     name: rosa-network-01
     namespace: default
   status:
     conditions:
     - lastTransitionTime: "2025-03-20T14:45:26Z"
       status: "True"
       type: RosaNetworkReady
   ```
   Example of failed network stack creation:
   ```
   kind: RosaNetwork
   metadata:
     name: rosa-network-01
     namespace: default
   status:
     conditions:
     - lastTransitionTime: "2025-03-18T13:25:16Z"
       status: "False"
       type: RosaNetworkReady
       severity: Error
       reason: ReconciliationFailed
       message: Insufficient privileges for ...
   ```
   Failed deletion:
   ```
   kind: RosaNetwork
   metadata:
     name: rosa-network-01
     namespace: default
   status:
     conditions:
     - lastTransitionTime: "2025-03-18T13:25:16Z"
       status: "False"
       type: RosaNetworkReady
       severity: Error
       reason: DeletionFailed
       message: ...
   ```
   
1. Modifications in the ROSA control plane CRD & reconciler so that it would be possible to reference the `RosaNetwork` resource from control plane:
   ```
   kind: ROSAControlPlane
     metadata:
       name: hcp01-control-plane
       namespace: default
   spec:
     rosaNetworkRef:
       name: hcp01-rosa-network
   ```
   Should the ROSA control plane CR contain reference to ROSA network, the reconciler will read the AWS region, AZ and subnet ids parameters from the ROSA network CR. The ROSA control plane should also be validated through a webhook so that it does not contain both the reference to `RosaNetwork` and the subnet ids and / or availability zones.

1. New tests.

### User Stories

1. As a CAPA user, I want to be able to provision the network infrastructure for ROSA HCP.
2. As a CAPA user, I want to be able to use the provisioned network infrastructure in ROSA HCP control plane.
3. As a CAPA user, I want to be able to delete the network infrastructure previously provisioned by CAPA.

#### Functional Requirements

1. Ability to create a new namespaced custom resource `RosaNetwork` with four attributes: name, AZ count or AZ list, region and CIDR block for VPC.
2. Reconciler implementing creation and deletion of the `RosaNetwork` resource.
3. Ability to reference the new custom resource from ROSA HCP control plane.

## Alternatives

1. Implement CRDs and reconcilers for each of the atoms of network infrastructure (VPCs, subnets, etc.).
2. Implement the network infrasructure similar to EKS, the network parameters being attributes of the EKS control plane.
3. Not implement anything and rely purely on AWS CLI or Terraform.

## Upgrade Strategy

The implementation will not affect CAPA upgrades.

<!-- Links -->
[community meeting]: https://docs.google.com/document/d/1ushaVqAKYnZ2VN_aa3GyKlS4kEd6bSug13xaXOakAQI/edit#heading=h.pxsq37pzkbdq
