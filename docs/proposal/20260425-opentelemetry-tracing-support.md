# Propose OpenTelemetry Support in CAPA

## Table of Contents
- Summary  
- Motivation  
- Problem Statement  
- Goals  
- Non Goals
- Proposal  
- User Stories  
- Implementation Details  
- Risks and Mitigations  
- Alternatives  
- Implementation Plan  
- Resources  

## Summary
With the increasing adoption of CAPA we want to improve the supportability of CAPA especially in production. CAPA currently lacks OpenTelemetry support for tracing in controllers. This proposal suggests adding OpenTelemetry support to its controllers enabling tracing for better observability.

## Motivation

### Problem Statement
With the increasing adoption of CAPA we want to improve the supportability of CAPA especially in production. CAPA currently lacks tracing support for its controllers making it difficult to track failures in controllers during reconciliation. Adding support for tracing would also benefit maintainers by enabling them to use tracing while investigating issues, flows etc.

## Goals
- OpenTelemetry SDK initialization and configuration in CAPA  
- Root spans at reconciliation boundaries for core CAPA controllers  
- Enable tracing for debugging and performance analysis  
- Add configurable sampling via flags  
- Enable trace collection locally when using tilt

## Non-Goals
- Enable/disable tracing at runtime without restarting  
- Tail-based sampling  

## Proposal

## User Stories

### 1. Controller Observability
As a CAPA operator, I want to trace controller reconciliation workflows so that I can diagnose cluster provisioning issues and understand execution behavior.

### 2. Debugging Infrastructure Failures
As a developer, I want visibility into AWS API calls and controller operations to quickly identify failures during cluster lifecycle events.

### 3. Platform / Vendor Observability
As a vendor or platform team using CAPA to manage Kubernetes clusters at scale, I want consistent and configurable tracing across controllers and AWS interactions so that I can debug issues across tenants, ensure reliability, and integrate with my observability stack.

## Implementation Details

### Telemetry Initialization

OpenTelemetry initialization logic will be encapsulated within a dedicated `otel` package (e.g., `pkg/otel`).

The CAPA controller entry point (`main.go`) will expose configuration flags to control tracing behavior (e.g., enable/disable, endpoint, sampling rate).

During startup:
- If tracing is enabled via flags, the `otel` package will be initialized.
- The initialization will set up the TracerProvider, resources, and exporters.
- The configured tracer provider will then be registered globally.

### Tracing Controller Reconciliation

Each controller’s reconciliation loop will be instrumented with spans.

**Example controllers:**
- awscluster_controller  
- awsmachine_controller 

**Conceptual flow:**

    Reconcile()
       ↓
    start tracing span
       ↓
    perform reconciliation logic
       ↓
    end span

**Spans will capture:**
- Request metadata  
- Execution duration  
- Result status  

### Shared Telemetry Utilities
To avoid duplication across multiple controllers, a shared telemetry helper package may be introduced.

**Proposed location:**

    pkg/otel

**Responsibilities:**
- Initialize tracing  
- Provide reusable helpers  
- Configure exporters/resources  

This approach ensures maintainability and avoids repetitive instrumentation code in each controller.

### Trace Export

Tracing data will be exported using the OpenTelemetry Protocol (OTLP) exporter.

CAPA will configure a single OTLP exporter over grpc (`otlptracegrpc`) within the controller process to send trace data to a user-defined endpoint. This endpoint will be configurable via controller flags.

The configured endpoint may be:
- An OpenTelemetry Collector (recommended for production), which can process, batch, and route traces to one or more backends.
- A compatible backend (e.g., Grafana Tempo or Jaeger) that supports OTLP ingestion.

### High-Level Architecture

    CAPA Controllers
          ↓
    OpenTelemetry Go SDK
          ↓
    OTLP Exporter
          ↓
    OpenTelemetry Collector
          ↓
    Jaeger / Grafana Tempo / Other Backends

### Sampling Strategy

To control the volume of traces and reduce performance overhead, sampling will be applied.
Sampling will be implemented using a ParentBased TraceIDRatioBased sampler, with the sampling ratio configurable via flags. It balances risk of OOM (which is a very possible risk in case of multiple large traces), operational visibility and system performance.

Sampling will be:
- Configurable via controller flags
- Applied at the trace level

Example configuration:
- Low sampling rate in production to minimize overhead
- Higher sampling rate during debugging or incident analysis

This ensures a balance between observability and system performance.
 
### Tracing AWS API Calls

AWS API calls can be traced using OpenTelemetry instrumentation for AWS SDK (`otelaws` middleware). This enables automatic creation of spans for all outgoing AWS API requests. Each AWS API span will be a child of the reconciliation span. it will be optional to trace AWS API call and can be enabled by flags.

However, this approach may generate excessive spans and introduce noise. Instead we can also create spans for important AWS API Calls For Ex- `CreateInstance`

To address this we can manually trace instrument important AWS operations such as CreateInstance or DeleteLoadBalancer etc.

### Functions to Trace
We start tracing our controllers from Reconcile function and then trace function which are most relevant for us

    Reconcile (root)
     ├── ReconcileDelete
     │    ├── DeleteLoadBalancer
     │    ├── DeleteInstances
     │    ├── CleanupNetwork
     │    └── CleanupSecurityGroups
     │
     └── ReconcileNormal
          ├── ReconcileLoadBalancer
          ├── ReconcileNetwork
          ├── ReconcileSecurityGroups
          └── ReconcileBastion

**We will primilarly trace:**
- Resource creation  
- Resource deletion  
- Reconciliation workflows  

### Context Propagation
The tracing context will be propagated via `context.Context` across all controller and AWS service calls. This ensures that all spans (controller + AWS calls) are linked under the same trace.

For Example:

    ctx, span := tracer.Start(ctx, "Reconcile")
    defer span.End()

    awsService.CreateInstance(ctx, ...)// passing context to the functions under the trace

### Additional Trace Data
We can attach some useful data to traces which will give us meaningful review of the traces

**Span Attributes:**

- **Kubernetes metadata**
  - `cluster.name`
  - `namespace`
  - `resource.name`
  - `resource.kind`

- **Reconciliation metadata**
  - `reconcile.id`
  - `reconcile.result`
  - `reconcile.duration`

- **AWS metadata**
  - `aws.resource.id`
  - `aws.region`

### Combined Event Recorder

We can explore a unified approach where Kubernetes events are recorded using the standard event recorder and also attached to OpenTelemetry spans.

This would allow better correlation between Kubernetes events and traces. Further investigation is required to determine the best implementation approach.

### Handling logs in relation to spans

We will explore correlating logs with traces by including trace and span IDs in structured logs. This will allow easier navigation between logs and traces.

Further investigation is required to finalize this approach.

### Exposed Configurations

| Flag                    | Description               |
|-------------------------|---------------------------|
| `--otel-enabled`        | Enable/disable tracing    |
| `--otel-endpoint`       | OTLP endpoint (collector or backend)   |
| `--otel-sampling-rate`  | Sampling ratio            |
| `--otelaws-enabled`     | Enable/disable aws api calls |

## Documentation Updates
To ensure usability and adoption, documentation will be updated to cover tracing setup and usage.

### Enabling Tracing
Documentation will describe:
- How to enable via flags  
- Configure endpoints and sampling  

### Collecting Traces via AWS ADOT
Guidance will be provided for integrating CAPA with AWS Distro for OpenTelemetry (ADOT), including:
- Use ADOT Collector  
- Export to AWS X-Ray  

### Collecting Traces via Self-Hosted Solutions
Documentation will include steps for using:
- OpenTelemetry Collector  
- Jaeger / Tempo  

### Local Development
Instructions will be provided for developers to:

- Run a local OpenTelemetry Collector
- Export traces to a local backend (e.g., Jaeger)
- Enable tracing during local development and testing

This ensures that tracing can be easily adopted in production, staging, and development environments.

## Risks and Mitigations

### Performance Overhead
Tracing introduces additional runtime overhead.
 
**Mitigation:**
- use configurable sampling
- allow tracing to be disabled through configuration
- ensure spans remain lightweight

### Complexity of Integration
Adding telemetry may increase code complexity.

**Mitigation:**
- isolate tracing logic into a dedicated observability package
- minimize modifications to existing controller logic 

## Resources
- [OpenTelemetry Documentation ](https://opentelemetry.io/docs/) 
- [Cluster API Provider AWS](https://cluster-api-aws.sigs.k8s.io/)