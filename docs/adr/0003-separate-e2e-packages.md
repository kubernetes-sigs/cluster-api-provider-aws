# 3. E2E tests organised as packages

* Status: accepted
* Date: 2020-11-09
* Authors: @richardcase
* Deciders: @randomvariable

## Context

The e2e tests for CAPA where initially in a single package and this package contained functional e2e and conformance tests for the unmanaged (i.e. purely EC2 with CAPBK) side of the provider.

With the addition of the EKS functionality to CAPA new e2e tests ([#1907](https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/1907)) are required and we need the ability to maintain and run e2e tests for the different parts of the provider separately.

## Decision

The e2e tests will be split into separate packages that will represent a distinct suite of e2e tests with accompying targets to enable running each suite. To reduce code duplication there will also be a shared e2e package that will contain functionality that is common to all the test package suites.

## Consequences

The existing e2e tests will be refactored into separate test suites. As this ADR is written retrospectively there is no issue but [#2102](https://github.com/kubernetes-sigs/cluster-api-provider-aws/pull/2102) is the initial consequnce of this and then [#1907](https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/1907) will be implemented as a new package.
