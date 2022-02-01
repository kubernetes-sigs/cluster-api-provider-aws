# 1. Use ADRs to record decisions

* Status: proposed
* Date: 2020-10-29
* Authors: @richardcase
* Deciders: [list of GitHub handles for those that made the decision]  <!-- mandatory -->

## Context

Currently decisions that affect ongoing CAPA development are being made via:

1. CAEP for larger or more impactful pieces of work
2. Discussions on issue, PRs, in Slack or in the bi-weekly community call.

The decisions made via 2) are not easily discoverable and may not be documented. The CAPA project needs a way to record these decisions so they are easily discoverable by the contributors when they are looking to understand why something is done a particular way.

## Decision

The project will use [Architectural Decision Records (ADR)](https://adr.github.io/) to record decisions that are applicable across the project and that are not explicitly covered by a CAEP. Additionally, in the implementation of a CAEP decisions may still be made via discussions and so ADRs should be created in this instance as well.

A [template](./0000-template.md) has been created based on prior work:

* https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions
* https://adr.github.io/madr/

## Consequences

When decisions are made that affect the entire project then these need to be recorded as an ADR. For decisions that are made as part of discussion on issues or PRs a new label could be added `needs-adr` (or something similar) so that its explicit. For decisions made on slack or via Zoom call someone will need to take an action to create the PR with the ADR. Maintainers and contributors will need to decide when a "decision" has been made and ensure an ADR is created.
