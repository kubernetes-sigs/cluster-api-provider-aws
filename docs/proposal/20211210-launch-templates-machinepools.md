---
title: Launch Templates for Machine Pools
authors:
  - "@richardcase"
reviewers:
  - "@randomvariable"
  - "@sedefsavas"
creation-date: 2021-12-10
last-updated: 2021-12-10
status: provisional
see-also: []
replaces: []
superseded-by: []
---

# Launch Templates for Machine Pools

## Table of Contents

[Tools for generating](https://github.com/ekalinin/github-markdown-toc) a table of contents from markdown are available.

## Glossary

* [CAPA](https://cluster-api.sigs.k8s.io/reference/glossary.html#capa) - Cluster API Provider AWS.
* [Launch Template](https://docs.aws.amazon.com/autoscaling/ec2/userguide/LaunchTemplates.html) - a configuration template that is used to configure an AWS EC2 instance when its created.
* [ASG](https://docs.aws.amazon.com/autoscaling/ec2/userguide/AutoScalingGroup.html) - an auto scale group that represents a pool of EC2 instances that can scale up & down automatically.


## Summary

Currently with CAPA we have 2 varieties of  **Machine Pools** implemented. Each variety has a differing level of support for [launch templates](https://docs.aws.amazon.com/autoscaling/ec2/userguide/LaunchTemplates.html).

The `AWSMachinePool` is used to create an **ASG** who's EC2 instances are used as worker nodes fro the Kubernetes cluster. The specification for `AWSMachinePool` exposes settings that are ultimately used to create a EC2 launch template (and version of it thereafter). 

The `AWSManagedMachine` is used to create a [EKS managed node group](https://docs.aws.amazon.com/eks/latest/userguide/managed-node-groups.html) which results in an AWS managed **ASG** being created that utilizses AWS managed EC2 instances. In the spec for `AWSManagedMachinePool` we expose details of the pool to create but we don't support using a launch template and we don't automatically create laucnh templates (like we do for `AWSMachinePool`).

This proposal outlines changes to CAPA that will introduce new capabilities to define and manage the lifecycle of **Launch Templates** in AWS and changes to both varieties of Machine Pools to utilise laucnh templates in a consistent way.

## Motivation

We are increasingly hearing requests from users of CAPA that a particular feature / configuration option isn't exposed by CAPAs implementation of machine pools and on investigation the feature is available via a launch template (nitro enclaves or placement as an example). In this scenario we can expose options on the machine pool spec and map to launch templates in the controller.  As we have 2 varieties of machine pools this means adding to both API definitions and we end up with a certain level of duplication.

With `AWSManagedMachinePools` currently we don't support using a launch template at all and this creates inconsistency with our `AWSMachinePool` implementation and limits the functionalitry we have exposed to the users.

The motivation is to improve consistency and provider a sinmgle way to manage the lifecycle of launch templates across both varieties of machine pools.

> Note: it may not be completely consistent in the initial implementation as we may need to deprecate some API definitions over time but the plan will be to be eventually consistent ;)


### Goals

- Consistent API to declare an use launch templates for `AWSMachinePool` and `AWSManagedMachinePool`
- Single point of reconciliation of launch templates
- Guide to the deprecation of certain API elements in `AWSMachinePool`

### Non-Goals/Future Work

- TBD

## Proposal

At a high-level the plan is to introduce 2 new API kinds called `AWSLaunchConfig` and `AWSLaunchConfigTemplate`. The usage of the 2 different kinds will follow the pattern adopted by Machine/MachineTemplate and EKSConfig/EKSConfigTemplate.

When creating a machine pool the user will supply a reference to a `AWSLaunchConfigTemplate` which specifies the configuration options that will ultimately be used to create a launch template. During reconciliation of the machine pools (both varieties) the controller will wait for the `AWSLaunchConfig` to be **Ready**. At which point it will create the ASG or EKS node group with a reference to the launch template.

The `AWSLaunchConfigTemplate` will be reconciled and a new instance of a `AWSLaunchConfig` that is specific to a `Cluster`. This means that the same `AWSLaunchConfigTemplate` can be used as the source across multiple clusters. 

### User Stories

#### Story 1

As a CAPA user
I want to create a machine pool using a launch template
So that i can use functionality from AWS launch templates

#### Story 2

AS a CAPA user
I want to create a managed machine pool using a launch template
So that i can use functionality from the AWS launch template

#### Story 3

As a CAPA user
I want to ensure launch templates are versioned
So that i can see a histiory of the changes in the console

#### Story 5

As a CAPA user
I want to be able to set the maximum number of versions to keep for launch templates
So that i don't have to worry about cleaning up old versions.

#### Story 6

As a CAPA user
I want to be able to use the output of a bootstrap provider in my launch template
So that i can bootstrap Kubernetes on the nodes

### Requirements

#### Functional Requirements

**FR1:** CAPA MUST support using launch templates with non-managed ASG based machine pools (i.e. `AWSMachinePool`).

**FR2:** CAPA MUST support using launch templates with EKS managed nodegroup based machine pools (i.e. `AWSManagedMachinePool`).

**FR3:** CAPA MUST provide a consistent declaritive API to expose Launch Template configuration to the machine pool implementations.

**FR4:** CAPA MUST manage the lifecycle of a launch template in AWS based on its declaration.

**FR5:** CAPA SHOULD allow using the same template for a launch template across different clusters.

**FR6:** CAPA MUST version launch templates in AWS.

**FR7:** CAPA MUST allow keeping a configurable number of previous versions of launch templates.

**FR8:** CAPA MUST support `clusterctl move` with machine pool launch templates.

**FR9:** CAPA MUST validate the declarations for `AWSLaunchConfig` & `AWSLaunchConfigTemplate`

#### Non-Functional Requirements

**NFR1:** CAPA MUST provide logging and tracing to expose the progress of reconciliation of `AWSLaunchConfig[Tenmplate`.

**NFR2:** CAPA MUST raise events at important milestine during reconciliation.

**NFR3:** CAPA MUST requeue where possible and not wait during reconciliation so as to free up the reconciliation loop

**NFR4:** CAPA must have e2e tests that cover usage of launch templates with BOTH variants of machine pools.


### Implementation Details/Notes/Constraints

- What are some important details that didn't come across above.
- What are the caveats to the implementation?
- Go in to as much detail as necessary here.
- Talk about core concepts and how they releate.

### Security Model

Document the intended security model for the proposal, including implications
on the Kubernetes RBAC model. Questions you may want to answer include:

* Does this proposal implement security controls or require the need to do so?
  * If so, consider describing the different roles and permissions with tables.
* Are their adequate security warnings where appropriate (see https://adam.shostack.org/ReederEtAl_NEATatMicrosoft.pdf for guidance).
* Are regex expressions going to be used, and are their appropriate defenses against DOS.
* Is any sensitive data being stored in a secret, and only exists for as long as necessary?

### Risks and Mitigations

- What are the risks of this proposal and how do we mitigate? Think broadly.
- How will UX be reviewed and by whom?
- How will security be reviewed and by whom?
- Consider including folks that also work outside the SIG or subproject.

## Alternatives

The `Alternatives` section is used to highlight and record other possible approaches to delivering the value proposed by a proposal.

## Upgrade Strategy

If applicable, how will the component be upgraded? Make sure this is in the test plan.

Consider the following in developing an upgrade strategy for this enhancement:
- What changes (in invocations, configurations, API use, etc.) is an existing cluster required to make on upgrade in order to keep previous behavior?
- What changes (in invocations, configurations, API use, etc.) is an existing cluster required to make on upgrade in order to make use of the enhancement?

## Additional Details

### Test Plan [optional]

**Note:** *Section not required until targeted at a release.*

Consider the following in developing a test plan for this enhancement:
- Will there be e2e and integration tests, in addition to unit tests?
- How will it be tested in isolation vs with other components?

No need to outline all of the test cases, just the general strategy.
Anything that would count as tricky in the implementation and anything particularly challenging to test should be called out.

All code is expected to have adequate tests (eventually with coverage expectations).
Please adhere to the [Kubernetes testing guidelines][testing-guidelines] when drafting this test plan.

[testing-guidelines]: https://git.k8s.io/community/contributors/devel/sig-testing/testing.md

### Graduation Criteria [optional]

**Note:** *Section not required until targeted at a release.*

Define graduation milestones.

These may be defined in terms of API maturity, or as something else. Initial proposal should keep
this high-level with a focus on what signals will be looked at to determine graduation.

Consider the following in developing the graduation criteria for this enhancement:
- [Maturity levels (`alpha`, `beta`, `stable`)][maturity-levels]
- [Deprecation policy][deprecation-policy]

Clearly define what graduation means by either linking to the [API doc definition](https://kubernetes.io/docs/concepts/overview/kubernetes-api/#api-versioning),
or by redefining what graduation means.

In general, we try to use the same stages (alpha, beta, GA), regardless how the functionality is accessed.

[maturity-levels]: https://git.k8s.io/community/contributors/devel/sig-architecture/api_changes.md#alpha-beta-and-stable-versions
[deprecation-policy]: https://kubernetes.io/docs/reference/using-api/deprecation-policy/

### Version Skew Strategy [optional]

If applicable, how will the component handle version skew with other components? What are the guarantees? Make sure
this is in the test plan.

Consider the following in developing a version skew strategy for this enhancement:
- Does this enhancement involve coordinating behavior in the control plane and in the kubelet? How does an n-2 kubelet without this feature available behave when this feature is used?
- Will any other components on the node change? For example, changes to CSI, CRI or CNI may require updating that component before the kubelet.

## Implementation History

- [ ] MM/DD/YYYY: Proposed idea in an issue or [community meeting]
- [ ] MM/DD/YYYY: Compile a Google Doc following the CAEP template (link here)
- [ ] MM/DD/YYYY: First round of feedback from community
- [ ] MM/DD/YYYY: Present proposal at a [community meeting]
- [ ] MM/DD/YYYY: Open proposal PR

<!-- Links -->
[community meeting]: https://docs.google.com/document/d/1Ys-DOR5UsgbMEeciuG0HOgDQc8kZsaWIWJeKJ1-UfbY

