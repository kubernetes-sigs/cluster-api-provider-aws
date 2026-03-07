# Agent Development Guide

## Project Context

This project implements a [Cluster API](https://github.com/kubernetes-sigs/cluster-api) (CAPI) provider for provisioning Kubernetes clusters in AWS. It supports pure EC2 based clusters (non-managed) and EKS clusters (managed). It implements the following provider contracts:

- infrastructure
- bootstrap
- control plane

## Rules & Constraints

- Ask clarifying questions and propose a plan
- Always use iterative development and keep changes as small as possible
- Add unit tests for all changes
- Major new features should have new e2e tests
- All unit tests should pass
- There should be no linting issues
- Project uses "scopes" and "services" pattern
- Do not update the Go version, Kubernetes or controller-runtime dependencies in go.mod
- Do not manually edit generated files, run `make generate` instead
- Do not automatically submit PRs

## Architecture

### Core Components

1. **Infrastructure Provider API Definitions (`/api` and `/exp/api`)**
   - `v1beta2`: Stable API version for core resources
   - `v1beta1`: Deprecated, use v1beta2
   - `/exp`: Experimental features (MachinePools, Fargate, ROSA)
   - Key resources: AWSCluster, AWSMachine, AWSManagedCluster, AWSMachinePool, AWSManagedMachinePool
   - Implements the CAPI infrastructure provider contract

2. **Infrastructure Provider Controllers (`/controllers` and `/exp/controllers`)**
   - Each controller reconciles a specific custom resource type
   - Controllers use scopes and services
   - Reconciliation pattern: observe state → determine actions → apply changes → update status
   - Key controllers: AWSClusterReconciler, AWSMachineReconciler, AWSMachinePoolReconciler

3. **EKS Bootstrap Provider API Definitions (`/bootstrap/eks/api`)**
   - `v1beta2`: Stable API version for core resources
   - `v1beta1`: Deprecated, use v1beta2
   - Key resources: EKSConfig, EKSConfigTemplate
   - Implements the CAPI bootstrap provider contract

4. **EKS Bootstrap Provider Controllers (`/bootstrap/eks/controllers`)**
   - Each controller reconciles a specific custom resource type
   - Controllers use scopes and services
   - Reconciliation pattern: observe state → determine actions → apply changes → update status
   - Key controllers: EKSConfigReconciler

5. **EKS Control Plane Provider API Definitions (`/controlplane/eks/api`)**
   - `v1beta2`: Stable API version for core resources
   - `v1beta1`: Deprecated, use v1beta2
   - Key resources: AWSManagedControlPlane
   - Implements the CAPI control plane provider contract

6. **EKS Control Plane Provider Controllers (`/controlplane/eks/controllers`)**
   - Each controller reconciles a specific custom resource type
   - Controllers use scopes and services
   - Reconciliation pattern: observe state → determine actions → apply changes → update status
   - Key controllers: AWSManagedControlPlaneReconciler

7. **Services Layer (`/pkg/cloud/services`)**
   - AWS Service-specific API clients organized by functional area
   - Examples: `ec2`, `s3`, `eks`
   - Each service has a `Reconcile*` and `Delete*` function called from the controllers

8. **Scope Package (`/pkg/cloud/scope`)**
   - Provides context and configuration created in the controllers and used in reconcilers
   - Scopes encapsulate cluster/machine specs, credentials, and AWS clients
   - Key scopes: ClusterScope, MachineScope, ManagedControlPlaneScode

5. **Feature Gates (`/feature`)**
   - Controls experimental/optional functionality
   - Important gates: `MachinePool`, `Fargate`

### Data Flow

```
User creates K8s resource → Controller watches → Reconciler triggered →
Scope created → Service methods called → AWS API interactions →
Status updated → Requeue if needed
```

## Commands

### Building

To build all binaries:

```bash
make binaries
```

To build a docker container for current architecture:

```bash
make doker-build
```

To build the docker containers for all supported architectures:

```bash
make docker build-all
```

### Testing & linting

To run the unit tests:

```bash
make test
```

To run the linters:

```bash
make lint
```

### Run code generators

To run all the code generators:

```bash
make generate
```