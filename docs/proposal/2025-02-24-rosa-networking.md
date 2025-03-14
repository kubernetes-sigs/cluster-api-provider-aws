---
title: Networking for ROSA HCP
authors:
  - "@mzazrivec"
reviewers:
  -
creation-date: 2025-02-24
last-updated: 2025-03-10
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
- [Implementation History](#implementation-history)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Glossary

Refer to the [Cluster API Book Glossary](https://cluster-api.sigs.k8s.io/reference/glossary.html).

## Summary

This proposal defines implementation of networking infrastructure in CAPA for ROSA Hosted Control Plane.

## Motivation

To be able to provision a new ROSA HCP kubernetes cluster using CAPA, one has to create and setup the underlying network infrastructure first: VPC, public and private subnets, internet gateway, routing tables for both subnets, elastic IP address.

All of the above can be currently provisioned and configured via AWS CLI, AWS Management Console or Terraform. Motivation for this work is to be able to provision and configure the network infrastructure for ROSA HCP using CAPI.

### Goals

1. Implement a new custom resource `RosaNetwork` representing the networking stack for ROSA HCP.
2. It will be possible to reference the new `RosaNetwork` resource from ROSA control plane resource
3. Implement standard CRUD for the new `RosaNetwork` resource.
4. Support the same networking scenarios as [ROSA CLI](https://github.com/openshift/rosa) using the same embeded AWS CloudFormation template that ROSA CLI uses.

### Non-Goals/Future Work

- Modify current networking code in AWS / EKS clusters.
- Support custom CloudFormation template.

## Proposal

The goal of this proposal is to be able to provision the networking infrastructure required for a ROSA HCP cluster.

[ROSA CLI](https://github.com/openshift/rosa) supports creation of the networking infrastructure for ROSA HCP and uses [AWS CloudFormation](https://aws.amazon.com/cloudformation/) template under the hood. The [CloudFormation template used by ROSA CLI](https://github.com/openshift/rosa/blob/master/cmd/create/network/templates/rosa-quickstart-default-vpc/cloudformation.yaml) allows to specify four parameters: CloudFormation stack name, AZ count (between 1 and 3), region and CIDR block for the VPC. The created CloudFormation stack then contains a VPC, public and private subnets (each pair created in separate AZ), internet gateway attached to VPC, elastic IPs, NAT gateways, public and private routing tables and a security group. 

Adopting the CloudFormation template used by rosa-cli would mean that CAPA and the `RosaNetwork` custom resource would be relying on a mechanism that is know to work well and any changes or fixes implemented in ROSA CLI would be picked up automatically in CAPA.

In practical terms, implementation of the proposal would mean:
1. A new custom resource definition `RosaNetwork` in CAPA
2. A new reconciler for the new custom resource, implementing CRUD.
3. Modifications in the ROSA control plane CRD & reconciler so that it would be possible to reference the `RosaNetwork` resource from control plane.
4. New tests.

### User Stories

1. As a CAPA user, I want to be able to provision the network infrastructure for ROSA HCP.
2. As a CAPA user, I want to be able to use the provisioned network infrastructure in ROSA HCP control plane.
3. As a CAPA user, I want to be able to modify the network inrastructure previously provisioned by CAPA.
4. As a CAPA user, I want to be able to delete the network infrastcture previously provisioned by CAPA.

#### Functional Requirements

1. Ability to create a new custom resource `RosaNetwork` with four attributes: name, AZ count, region and CIDR block for VPC.
2. Reconciler implementing CRUD for the `RosaNetwork` resource.
3. Ability to reference the new custom resource from ROSA HCP control plane.

## Alternatives

1. Implement CRDs and reconcilers for each of the atoms of network infrastructure (VPCs, subnets, etc.).
2. Implement the network infrasructure similar to EKS, the network parameters being attributes of the EKS control plane.
3. Not implement anything and rely purely on AWS CLI or Terraform.

## Upgrade Strategy

The implementation will not affect CAPA upgrades.

## Implementation History

- [ ] 2025-03-17: Proposed idea in an issue or [community meeting]

<!-- Links -->
[community meeting]: https://docs.google.com/document/d/1ushaVqAKYnZ2VN_aa3GyKlS4kEd6bSug13xaXOakAQI/edit#heading=h.pxsq37pzkbdq
