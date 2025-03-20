# Problem

EKS and ROSA are implemented in (relatively) standalone sets of controllers. These are currently grouped with the associated feature gates.

However, the feature gate mechanism has no specifier for GA, meaning that both implementations have stayed in beta.

This proposal supercedes `docs/proposal/20250314-optional-controllers.md`.

# Background

Cluster API Provider for AWS offers a `--feature-flags` argument to manage introduction of new features.

This allows features to move through `alpha` and `beta` phases, but there is not currently a defined path to general availability.

Specifically of interest is controllers for the `EKS` and ROSA`, which offer optional features to the program.

This proposal focuses on how these controllers could be grouped in order to remain optional, while also giving a path to GA.

## Goals

 - Allow features spanning groups of controllers to graduate, while also remaining optional

## Non-goals

 - Provide a generalized path to general availbility


## Proposed solution

Introduce a new argument, `--disable-controllers`, which allows controllers to be logically grouped as feature sets and turned on and off independently.

The groups and their disabled status will be tracked in a map within a private module, and can be queried via exposed functions.

The map's structure and functions would be as follows.

```go
var disabledControllers = map[string]bool{
    ControllerGroupName: false,
}

// IsDisabled checks if a controller is disabled.
// If the name provided is not in the map, this will return 'false'.
func IsDisabled(name string) bool

// GetValidNames returns a list of controller names that are valid to disable.
// Note: these are the entries in the `disabledControllers` variable.
// Used for error and help messages 
func GetValidNames() []string

// ValidateNamesAndDisable validates a list of controller names against the known set, and disables valid names.
func ValidateNamesAndDisable(names []string) error
```

Within `main.go`, `ValidateNamesAndDisable` will check against the contents of the `--disable-controllers` slice.
Valid entires will then be marked as `true`, indicating they are disabled.

Before initializing a controller or group of controllers, `IsDisabled` can be checked to determine whether or not they should register with the manager and start.

If a controller is disabled, a log message indicating it is disabled should be emitted.
This will aid users in troubleshooting, should the deployment behave unexpectedly.

## Logical groups

EKS and ROSA are currently behind feature gate checks.
These checks can be updated to instead use `IsDisabled` and entries within the `disabledControllers` map can be made.

## Core controllers and alternatives

The proposal `2025-01-07-aws-self-managed-feature-gates.md` was merged and planned to add `AWSMachine` and `AWSCluster` feature gates.
This merged with lazy concensus and the implementation was proposed in [PR #5284](https://github.com/kubernetes-sigs/cluster-api-provider-aws/pull/5284).

However, maintainers were opposed to moving these core controllers into feature gates, which would have put them into a permanent pre-GA phase.
This leads to a conflict between the implementation and the proposal, which was indeed accepted.

This current proposal includes the ability to disable `AWSMachine` and `AWSCluster` as a compromise, using the `unmanaged` set; it achieves the goals of the original proposal, while addressing the objections to the proposed implementation.
