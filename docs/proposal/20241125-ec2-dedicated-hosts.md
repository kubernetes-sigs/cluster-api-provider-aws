---
title: Add Support for EC2 Dedicated Hosts
authors:
  - "@faermanj"
reviewers:
creation-date: 2024-11-25
last-updated: 2024-11-25
status: draft
see-also: []
replaces: []
superseded-by: []
---

# Add Support for EC2 Dedicated Hosts

## Table of Contents

- [Add Support for EC2 Dedicated Hosts](#start)
  - [Table of Contents](#table-of-contents)
  - [Glossary](#glossary)
  - [Summary](#summary)
  - [Motivation](#motivation)
    - [Goals](#goals)
    - [Non-Goals/Future Work](#non-goalsfuture-work)
  - [Proposal](#proposal)
    - [User Stories](#user-stories)
      - [Story 1](#story-1)
    - [Requirements](#requirements)
      - [Functional Requirements](#functional-requirements)
      - [Non-Functional Requirements](#non-functional-requirements)
    - [Implementation Details/Notes/Constraints](#implementation-detailsnotesconstraints)
    - [Security Model](#security-model)
    - [Risks and Mitigations](#risks-and-mitigations)
  - [Alternatives](#alternatives)
  - [Upgrade Strategy](#upgrade-strategy)
  - [Additional Details](#additional-details)
    - [Test Plan](#test-plan)
    - [Graduation Criteria](#graduation-criteria)
  - [Implementation History](#implementation-history)

## Glossary

- [CAPA](https://cluster-api.sigs.k8s.io/reference/glossary.html#capa) - Cluster API Provider AWS. 
- [CAPI](https://github.com/kubernetes-sigs/cluster-api) - Cluster API

## Summary
The "Dedicated Hosts" feature of Amazon EC2 lets customers allocate physical hosts, with explicit hardware capacity, and allocate instances on those hosts. Also, instances on dedicated hosts have an "affinity" setting (Default or Host affinity), specifying the instance behavior on stopping and restarting the instance.
More information can be found on the [dedicated hosts feature documentation](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/dedicated-hosts-overview.html).
This proposal is for adding support of that feature on CAPA, allowing users to leverage the provider feature.

## Motivation
Dedicated hosts are primarily used for licensing compliance, when the required software (i.e. enterprise databases) needs to account for the physical host capacity for licensing. Also, it is a mechanism that can control explicitly which instances are executed on the same hosts, that might be relevant in scenarios related to security and performance predictability. 

### Goals
1. On cluster creation, let users specify the host id and affinity for instances.

### Non-Goals/Future Work
- Dedicated Host allocation. For now, dedicated hosts must be pre-allocated and explicitly configured. In future work, we may want to auto-allocate / release dedicated hosts.
- Custom replica allocation policy. For now, replicas are allocated in the first dedicated host that accepts the instance. In future work, we may let users configure other policies (round-robin, least-utilized, ...)

## Proposal
- Add list host id and affinity to the proper kind spec (`AWSMachineTemplate`?).
- Define and document the initial policy for instance allocation
- Add an E2E test case


### Risks and Mitigations


## Alternatives

Manually configuring the cluster, setting the appropriate host affinities to match the desired dedicated hosts.

#### Benefits

- Let enterprise customers stay compliant with licensing policies and other dedicated hosts applications.

#### Downsides

- "Instance to host" mapping might be too much for CAPA to manage.

#### Decision

## Upgrade Strategy
No impact on upgrades.

## Additional Details

### Test Plan
* Test creating a cluster, confirm all instances are executed on their respective dedicated hosts.
 
### Graduation Criteria

## Implementation History

- [x] 2024-11-25: Open proposal (PR)

<!-- Links -->
[ec2 dedicated hosts]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/dedicated-hosts-overview.html
[discussion]: https://github.com/kubernetes-sigs/cluster-api-provider-aws/discussions/5213
