---
title: Feature gates for AWSCluster and AWSMachine
authors:
  - "@mzazrivec"
reviewers:
  - "@serngawy"
  - "@nrb"
creation-date: 2025-01-07
last-updated: 2025-01-13
status: provisional
---

# Feature gates for AWS self managed clusters

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

Feature gates are a mechanism that CAPA uses to enable or disable particular features components of CAPA. This document proposes to implement two new feature gates, one for AWSMachine controller, another one for AWSCluster controller.

## Motivation

Motivation for the two new feature gates is the option to turn on and off CAPA's ability to reconcile AWSMachine and AWSCluster resources. Currently, controllers for these two resource types are always on.

The possibility to turn the respective controllers on and off becomes important in multi-tenant CAPA installations (more than one CAPA installed in a single kubernetes cluster). The two new feature gates will help avoid interference and conflicts between the controllers during resource reconciliations.

### Goals

1. Implement a feature gate for AWSMachine controller
2. Implement a feature gate for AWSCluster controller

### Non-Goals/Future Work

- change any of the existing feature gates or their semantics

## Proposal

- Introduce AWSMachine feature gate which would turn on and off the AWSMachineReconciler. The feature gate would be enabled by default.
- Introduce AWSCluster feature gate which would turn on and off the AWSClusterReconciler. The feature gate would be enabled by default.

Both feature gates would be defined in [feature.go](https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/main/feature/feature.go#L25), the default values would be set in the same file [feature.go](https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/main/feature/feature.go#L97). The logic turning the particular reconcilers on and off would happen in [main.go](https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/main/main.go).

Having two separate feature gates for AWSMachine and AWSCluster would mean that user could enable AWSMachine without enabling AWSCluster. This setup may lead to a situation where the reconciliation may not work correctly and the manager would warn the user of this during the start. Alternatively, the manager could just print an error and exit during start.

On the other hand, user should be able to enable the AWSCluster feature gate along with the MachinePool feature gate and have the AWSMachine feature gate disabled. This setup should be working correctly.

### User Stories

1.As a CAPA user and a cluster admin, I want to be able to install two (or more) CAPA instances on my cluster. The first instance would have the feature gates for AWS self managed cluster enabled, the other instances would have those feature gates disabled.
2. As a CAPA user, I want to be able to provision only EKS clusters.
3. As a CAPA user, I want to be able to provision only ROSA clusters.
4. As a CAPA user, I want to be able to provision only EKS and ROSA clusters.

#### Functional Requirements

1. Ability to enable / disable feature gates for AWSCluster and AWSMachine controllers.
2. Both feature gates would be enabled by default.
3. In case the AWSMachine feature gate is enabled and the AWSCluster feature gates is disabled, the manager would print a warning informing the user to enable both feature gates, as things may not work correctly.

## Alternatives

1. The alternative here would be simply not to implement the two new feature gates and rely just on namespace separation.
2. Another alternative would be to use the combination of `--watch-filter` on the manager side and object labeling on the user side.

## Upgrade Strategy

Users upgrading should expect two new feature gates in CAPA, both of them being enabled by default.

## Implementation History

- [ ] 2025-01-20: Proposed idea in an issue or [community meeting]

<!-- Links -->
[community meeting]: https://docs.google.com/document/d/1ushaVqAKYnZ2VN_aa3GyKlS4kEd6bSug13xaXOakAQI/edit#heading=h.pxsq37pzkbdq
