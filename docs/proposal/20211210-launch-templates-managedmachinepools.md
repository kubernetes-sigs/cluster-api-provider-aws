---
title: Launch Templates for Managed Machine Pools
authors:
  - "@richardcase"
reviewers:
  - "@sedefsavas"
  - "@richardchen331"
creation-date: 2021-12-10
last-updated: 2022-03-29
status: provisional
see-also: []
replaces: []
superseded-by: []
---

# Launch Templates for Managed Machine Pools

## Table of Contents

- [Launch Templates for Managed Machine Pools](#launch-templates-for-managed-machine-pools)
  - [Table of Contents](#table-of-contents)
  - [Glossary](#glossary)
  - [Summary](#summary)
  - [Motivation](#motivation)
    - [Goals](#goals)
    - [Non-Goals/Future Work](#non-goalsfuture-work)
  - [Proposal](#proposal)
    - [User Stories](#user-stories)
      - [Story 1](#story-1)
      - [Story 2](#story-2)
      - [Story 3](#story-3)
      - [Story 4](#story-4)
      - [Story 5](#story-5)
    - [Requirements](#requirements)
      - [Functional Requirements](#functional-requirements)
      - [Non-Functional Requirements](#non-functional-requirements)
    - [Implementation Details/Notes/Constraints](#implementation-detailsnotesconstraints)
    - [Security Model](#security-model)
    - [Risks and Mitigations](#risks-and-mitigations)
  - [Alternatives](#alternatives)
    - [New `AWSLaunchTemplate` CRD & Controller](#new-awslaunchtemplate-crd--controller)
      - [Benefits](#benefits)
      - [Downsides](#downsides)
      - [Decision](#decision)
  - [Upgrade Strategy](#upgrade-strategy)
  - [Additional Details](#additional-details)
    - [Test Plan](#test-plan)
    - [Graduation Criteria](#graduation-criteria)
  - [Implementation History](#implementation-history)

## Glossary

- [CAPA](https://cluster-api.sigs.k8s.io/reference/glossary.html#capa) - Cluster API Provider AWS. 
- [CAPI](https://github.com/kubernetes-sigs/cluster-api) - Cluster API.
- [Launch Template](https://docs.aws.amazon.com/autoscaling/ec2/userguide/LaunchTemplates.html) - a configuration template that is used to configure an AWS EC2 instance when its created.
- [ASG](https://docs.aws.amazon.com/autoscaling/ec2/userguide/AutoScalingGroup.html) - an auto-scale group that represents a pool of EC2 instances that can scale up & down automatically.

## Summary

Currently, with CAPA we have 2 varieties of  **Machine Pools** implemented called `AWSMachinePool` and `AWSManagedMachinePool`. Each variety has a differing level of support for [launch templates](https://docs.aws.amazon.com/autoscaling/ec2/userguide/LaunchTemplates.html).

The `AWSMachinePool` is used to create an **ASG** who's EC2 instances are used as worker nodes for the Kubernetes cluster. The specification for `AWSMachinePool` exposes settings that are ultimately used to create a EC2 launch template (and version of it thereafter) via the `AWSLaunchTemplate` field and struct:

```go
// AWSLaunchTemplate specifies the launch template and version to use when an instance is launched.
// +kubebuilder:validation:Required
AWSLaunchTemplate AWSLaunchTemplate `json:"awsLaunchTemplate"`
```

([source](https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/main/exp/api/v1beta1/awsmachinepool_types.go#L67))

The `AWSManagedMachinePool` is used to create a [EKS managed node group](https://docs.aws.amazon.com/eks/latest/userguide/managed-node-groups.html) which results in an AWS managed **ASG** being created that utilises AWS managed EC2 instances. In the spec for `AWSManagedMachinePool` we expose details of the pool to create but we don't support using a launch template, and we don't automatically create launch templates (like we do for `AWSMachinePool`). There have been a number of requests from users of CAPA that have wanted to use `AWSManagedMachinePool` but we don't expose required functionality that only comes with using launch templates.

This proposal outlines changes to CAPA that will introduce new capabilities to utilise launch templates for `AWSManagedMachinePool` and brings its functionality in line with `AWSMachinePool`.

## Motivation

We are increasingly hearing requests from users of CAPA that a particular feature / configuration option isn't exposed by CAPAs implementation of managed machine pools (i.e. `AWSManagedMachinePool`) and on investigation the feature is available via a launch template (nitro enclaves or placement as an example). In some instances, users of CAPA have had to use unmanaged machine pools (i.e. `AWSMachinePool`) instead.

The motivation is to improve consistency between the 2 varieties of machine pools and expose to the user features of launch templates.

> Note: it may not be completely consistent in the initial implementation as we may need to deprecate some API definitions over time but the plan will be to be eventually consistent ;)

### Goals

- Consistent API to use launch templates for `AWSMachinePool` and `AWSManagedMachinePool`
- Single point of reconciliation of launch templates
- Guide to the deprecation of certain API elements in `AWSManagedMachinePool`

### Non-Goals/Future Work

- Add non-existent controller unit tests for `AWSMachinePool` and `AWSManagedMachinePool`

## Proposal

At a high level, the plan is to:

1. Add a new `AWSLaunchTemplate` field to [AWSManagedMachinePoolSpec](https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/main/exp/api/v1beta1/awsmanagedmachinepool_types.go#L65) that uses the existing [AWSLaunchTemplate](https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/ec057ad6e613a6578f67bf68a6c77fbe772af933/exp/api/v1beta1/types.go#L58) struct. For example:

```go
// AWSLaunchTemplate specifies the launch template and version to use when an instance is launched. This field 
// will become mandatory in the future and its recommended you use this over fields AMIType,AMIVersion,InstanceType,DiskSize,InstanceProfile.
// +optional
AWSLaunchTemplate AWSLaunchTemplate `json:"awsLaunchTemplate"`
```

2. Update the comments on the below fields of [AWSManagedMachinePoolSpec](https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/9bc29570614aa7123d79f042b6e6efc2aaf3e490/exp/api/v1beta1/awsmanagedmachinepool_types.go#L65) to indicate that the fields is deprecated and that `AWSlaunchTemplate` should be used.
    - AMIVersion
    - AMIType
    - DiskSize
    - InstanceType
3. Add new `LaunchTemplateID` and `LaunchTemplateVersion` fields to [AWSManagedMachinePoolStatus](https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/9bc29570614aa7123d79f042b6e6efc2aaf3e490/exp/api/v1beta1/awsmanagedmachinepool_types.go#L171) to store details of the launch template and version used.
4. Add a new `LaunchTemplateVersion` field to [AWSMachinePoolStatus](https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/main/exp/api/v1beta1/awsmachinepool_types.go#L112) to store the version of the launch template used.
5. [Refactor the code](https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/ec057ad6e613a6578f67bf68a6c77fbe772af933/exp/controllers/awsmachinepool_controller.go#L383) from the `AWSMachinePool` controller that reconciles `AWSLaunchTemplate` into a common location so that it can be shared.
6. Update the controller for `AWSManagedMachinePool` to use the `AWSLaunchTemplate` reconciliation logic.
7. Add checks in the `AWSManagedMachinePool` create/update validation webhooks that stops users specifying `AWSLaunchTemplate` if fields `AMIType,AMIVersion,InstanceType,DiskSize,InstanceProfile` are set
8. Add warning logs to the `AWSManagedMachinePool` create/update validation webhooks if fields `AMIType,AMIVersion,InstanceType,DiskSize,InstanceProfile` stating that these fields will be deprecated in the future and that `AWSLaunchTemplate` should be used instead
> An area that is undecided upon is should we auto convert the `AMIType,AMIVersion,InstanceType,DiskSize,InstanceProfile` fields if specified into a `AWSLaunchTemplate`. We should investigate this as part of implementation. 
10. Update the cluster templates that use `AWSManagedMachinePool` so that they use `AWSLaunchTemplate`
11. Update the API version roundtrip tests for v1alpha4<->v1beta1 conversions of `AWSManagedMachinePool`
12. Update the EKS e2e tests to add an additional test step where we create an additional managed machine pool using  `AWSLaunchTemplate`.
13. Update any relevant documentation
14. Release note must mention that "action is required" in the future, as fields are being deprecated.
15. Ensure that we capture the field deprecations for future removal in an API version bump.

### User Stories

#### Story 1

AS a CAPA user
I want to create a managed machine pool using a launch template
So that I can use functionality from the AWS launch template

#### Story 2

As a CAPA user
I want to have consistency between managed and unmanaged machine pools
So that I can choose which to use based on whether I want managed and not based on missing functionality

#### Story 3

As a CAPA user
I want to ensure that changes to the pool result in a new version of the launch templates
So that I can see a history of the changes in the console

#### Story 4

As a CAPA user
I want the controller to clean up old launch templates / launch template versions
So that I don't have to worry about cleaning up old versions and so i don't exceed the AWS limits
(see [AWS docs](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-launch-templates.html) for limits)

#### Story 5

As a CAPA user
I want to be able to use the output of a bootstrap provider in my launch template
So that i can bootstrap Kubernetes on the nodes

### Requirements

#### Functional Requirements

**FR1:** CAPA MUST continue to support using launch templates with non-managed ASG based machine pools (i.e. `AWSMachinePool`).

**FR2:** CAPA MUST support using launch templates with EKS managed nodegroup based machine pools (i.e. `AWSManagedMachinePool`).

**FR3:** CAPA MUST provide a consistent declarative API to expose Launch Template configuration to the machine pool implementations.

**FR4:** CAPA MUST manage the lifecycle of a launch template in AWS based on its declaration.

**FR5:** CAPA MUST version launch templates in AWS.

**FR6:** CAPA MUST allow keeping a configurable number of previous versions of launch templates.

**FR7:** CAPA MUST validate the declarations for `AWSLaunchTemplate`

#### Non-Functional Requirements

**NFR1:** CAPA MUST provide logging and tracing to expose the progress of reconciliation of `AWSLaunhcTemplate`.

**NFR2:** CAPA MUST raise events at important milestones during reconciliation.

**NFR3:** CAPA MUST requeue where possible and not wait during reconciliation so as to free up the reconciliation loop

**NFR4:** CAPA must have e2e tests that cover usage of launch templates with BOTH variants of machine pools.

### Implementation Details/Notes/Constraints

The code in [reconcileLaunchTemplate](https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/ec057ad6e613a6578f67bf68a6c77fbe772af933/exp/controllers/awsmachinepool_controller.go#L383) must be refactored into a package that can be use by the `AWSManagedMachinePool` controller as well. We could think about shifting more of this functionality into the "ec2" service.

Cleaning up old versions of launch templates is currently handled by [PruneLaunchTemplateVersions](https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/ec057ad6e613a6578f67bf68a6c77fbe772af933/pkg/cloud/services/ec2/launchtemplate.go#L265) which is sufficient for this change. We may want to make the minimum number of versions to keep configurable in the future but this can be covered by a different change.

### Security Model

There are no changes required to the security model. Access to the required CRDs is already declared for the controllers and as we are not adding any new kinds this doesn't need to change.

No change is required to the AWS permissions the controller requires for reconciliation.

### Risks and Mitigations

The risk is that we are being constrained by the existing API definition used in unmanaged machine pools. This may raise unforeseen issues.

## Alternatives

### New `AWSLaunchTemplate` CRD & Controller

The idea is that a `AWSLaunchTemplate` CRD would be created with an associated controller. The controller would then be responsible for reconciling the definition and managing the lifecycle of launch templates on AWS.

#### Benefits

- Single point of reconciliation and lifecycle management of launch templates in AWS. 
- Separate lifecycle per launch template. So, we can change the number of previous instances to keep etc.

#### Downsides

- Additional complexity of orchestrating the creation of the launch template with the bootstrap data. The machine pool reconcilers would need to wait for the bootstrap data and the launch template.
- Would require deprecation of fields in 2 CRDs (i.e both machine pool varieties). 

#### Decision

As `AWSMachinePool` already managed launch templates, it was felt that we should follow the same approach for consistency and it would be a smaller change.

We can revisit the idea of a separate launch template kind in the future. The proposed change in this proposal will not preclude implementing this alternative in the future.

## Upgrade Strategy

The changes we are making to `AWSManagedMachinePool` are optional. Therefore, current users do not have to use the new `AWSLaunchTemplate` field. On upgrading there will be a new log entry written that informs the user that certain fields will be deprecated in the future.

## Additional Details

### Test Plan

- There are currently no controller unit tests for the machine pools in CAPA. We do need to add tests, but this can be done as part of a separate change.
- The EKS e2e tests will need to be updated so that a managed machine pool is created with a launch template specified.

### Graduation Criteria

With this proposal, we are planning to deprecate a number of fields on `AWSManagedMachinePool`

The current API version is **beta level** and this normally means:

- We must support the beta API for 9 months or 3 releases (whichever is longer). See [rule 4a](https://kubernetes.io/docs/reference/using-api/deprecation-policy/)

However, the machine pools feature is marked as experimental in CAPI/CAPA and as such it has to be explicitly enabled via a feature flag. Therefore its proposed that we remove the deprecated fields when we bump the api version from v1beta. As part of the field removal we will update the API conversion functions to automatically populate `AWSLaunchTemplate` on create.

## Implementation History

- [x] 2021-12-10: Initial WIP proposal created
- [x] 2021-12-13: Discussed in [community meeting]
- [x] 2022-01-14: Discussions between richardcase and richardchen331 on slack
- [x] 2022-02-04: Updated proposal based on discussions
- [x] 2022-02-05: Created proposal [discussion]
- [x] 2022-02-07: Present proposal at a [community meeting]
- [x] 2022-02-05: Open proposal PR
- [x] 2022-03-29: Updated based on review feedback

<!-- Links -->
[community meeting]: https://docs.google.com/document/d/1iW-kqcX-IhzVGFrRKTSPGBPOc-0aUvygOVoJ5ETfEZU/edit#
[discussion]: https://github.com/kubernetes-sigs/cluster-api-provider-aws/discussions/3154
