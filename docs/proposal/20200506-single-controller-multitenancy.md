---
title: Single Controller Multitenancy
authors:
  - "@randomvariable"
  - "@andrewmyhre"
reviewers:
  - "@ashoksekar07"
  - "@bagnaram"
  - "@CecileRobertMichon"
  - "@detiber"
  - "@devigned"
  - "@nader-ziada"
  - "@ncdc"
  - "@richardcase"
  - "@rudoi"
  - "@sethp-nr"
creation-date: 2020-05-05
last-updated: 2020-05-05
status: implementable
see-also:
- https://github.com/kubernetes-sigs/cluster-api-provider-azure/issues/586
replaces: []
superseded-by: []
---

# Single Controller Multitenancy

## Table of Contents

- [Single Controller Multitenancy](#single-controller-multitenancy)
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
  - [Requirements](#requirements)
    - [Functional](#functional)
    - [Non-Functional](#non-functional)
    - [Implementation Details/Notes/Constraints](#implementation-detailsnotesconstraints)
      - [Proposed Changes](#proposed-changes)
        - [Cluster API Provider AWS v1alpha3 types](#cluster-api-provider-aws-v1alpha3-types)
    - [Controller Changes](#controller-changes)
    - [Clusterctl changes](#clusterctl-changes)
      - [Validating webhook changes](#validating-webhook-changes)
      - [Identity Type Credential Provider Behaviour](#identity-type-credential-provider-behaviour)
    - [Security Model](#security-model)
      - [Roles](#roles)
    - [RBAC](#rbac)
        - [Write Permissions](#write-permissions)
    - [Namespace Restrictions](#namespace-restrictions)
    - [CAPA Controller Requirements](#capa-controller-requirements)
    - [Alternative Approaches Considered](#alternative-approaches-considered)
      - [Using only secrets or RoleARN field on AWSCluster](#using-only-secrets-or-rolearn-field-on-awscluster)
      - [1:1 mapping one namespace to one AWSIdentity](#11-mapping-one-namespace-to-one-awsidentity)
    - [Risks and Mitigations](#risks-and-mitigations)
      - [Network assumptions are made explicit](#network-assumptions-are-made-explicit)
      - [Caching and handling refresh of credentials](#caching-and-handling-refresh-of-credentials)
      - [AWS Cloud Provider behaviour with regards to cluster names](#aws-cloud-provider-behaviour-with-regards-to-cluster-names)
  - [Upgrade Strategy](#upgrade-strategy)
  - [Additional Details](#additional-details)
    - [Test Plan](#test-plan)
    - [Graduation Criteria](#graduation-criteria)
      - [Alpha](#alpha)
      - [Beta](#beta)
      - [Stable](#stable)
    - [Version Skew Strategy](#version-skew-strategy)
  - [Implementation History](#implementation-history)

## Glossary

* Identity Type - One of several ways to provide a form of identity that is ultimately resolved to an AWS access key ID,
  secret access key and optional session token tuple.
* Credential Provider - An implementation of the interface specified in the [AWS SDK for
  Go][aws-sdk-go-credential-provider].
* CAPA - An abbreviation of Cluster API Provider AWS.

## Summary

The CAPA operator is able to manage cloud infrastructure within the permission scope of the AWS principle it is
initialized with. It is expected that the CAPA operator will be provided credentials via the deployment, either
explicitly via environment variables or implicitly via the default SDK credential provider chain, including EC2 instance
profiles or metadata proxies such as [kiam][kiam].

Currently, without a custom deployment, role assumption in CAPA can only take place using a metadata proxy such as
[kiam][kiam], or using OIDC based service account token volume projection, requiring Kubernetes 1.16, and special setup
of certificates with an appropriate OIDC trust.

It is technically possible to provide an IAM role via a shared configuration file, but the UX has been found to be
complicated, and additional breaking code changes have been suggested to resolve these issues.

In addition, where role assumption is configured, that role is used for the entire lifetime of the CAPA deployment. This
also means that an AWSCluster could be broken if the instance of CAPA that created it is misconfigured for another set
of credentials.

This proposal outlines new capabilities for CAPA to use IAM role assumption to assume a permission set in a different
AWS account, at runtime, on a per-cluster basis. The proposed changes would be fully backwards compatible and maintain
the existing behavior with no changes to user configuration required.


## Motivation

For large organizations, especially highly-regulated organizations, there is a need to be able to perform separate
duties at various levels of infrastructure - permissions, networks and accounts. VPC sharing is a model which provides
separation at the AWS account level. Within this model it is appropriate for tooling running within the 'management' account
to manage infrastructure within the 'workload' accounts, which requires a identity in the management account which can
assume a identity within a workload account. For CAPA to be most useful within these organizations it will need to
support multi-account models.

Some organizations may also delegate the management of clusters to another third-party. In that case, the boundary
between organizations needs to be secured, and issues such as confused deputy scenarios need to be mitigated.

Because a single deployment of the CAPA operator may reconcile many different clusters in its lifetime, it is necessary
to modify the CAPA operator to scope its AWS client instances to within the reconciliation process.

AWS provides a number of mechanisms for assume roles across account boundaries:

* Role assumption using credentials provided in environment variables
* Role assumption using credentials provided by a metadata service, e.g. IMDS or ECS
* Role assumption using an OIDC JWT token
* Role assumption using explicitly provided credentials

The first three methods are supported per instance of the CAPA controller today. Of special note is OIDC, which was
introduced in 2019 as a way to leverage projected service account tokens from Kubernetes to assume role. This allows for
fine-grained specification of IAM roles using Kubernetes primitives, but is still limited to a single role assumption
per instance of CAPA.

### Goals

1. To enable AWSCluster resources reconciliation across AWS account boundaries
2. To maintain backwards compatibility and cause no impact for users who don't intend to make use of this capability

### Non-Goals/Future Work

- To enable Machines to be provisioned in AWS accounts different than their control planes. This would require adding
  various AWS infrastructure specification to the AWSMachineTemplate type which currently does not exist.
- To enable control plane and worker machines to be provisioned in separate accounts.

## Proposal

### User Stories

#### Story 1

Alex is an engineer in a large organization which has a strict AWS account architecture. This architecture dictates that
Kubernetes clusters must be hosted in dedicated AWS accounts. The organization has adopted Cluster API in order to
manage Kubernetes infrastructure, and expects 'management' clusters running the Cluster API controllers to manage 'workload'
clusters in dedicated AWS accounts.

The current configuration exists: AWS Account 'management':
* Vpc, subnets shared with 'workload' AWS Account
* EC2 instance profile linked to the IAM role 'ClusterAPI-Mgmt'
* A management Kubernetes cluster running Cluster API Provider AWS controllers using the 'ClusterAPI-Owner' IAM identity.

AWS Account 'workload':
* Vpc and subnets provided by 'Owner' AWS Account
* IAM Role 'ClusterAPI-Participant' which trusts IAM role 'ClusterAPI-Owner'

Alex can provision a new cluster in the 'workload' AWS Account by creating new Cluster API resources in the management cluster.
Cluster 'Owner’ Alex specifies the IAM role 'ClusterAPI-Participant' in the AWSCluster spec. The CAPA controller in the management cluster assumes the role 'ClusterAPI-Participant' when reconciling the AWSCluster so that it can
create/use/destroy resources in the 'workload' AWS Account.


#### Story 2

Dascha is an engineer in a smaller, less strict organization with a few AWS accounts intended to host all
infrastructure. There is a single AWS account named 'dev', and Dascha wants to provision a new cluster in this account.
An existing Kubernetes cluster is already running the Cluster API operators and managing resources in the dev account.
Dascha can provision a new cluster by creating Cluster API resources in the existing cluster, omitting the IAMRoleARN
field in the AWSCluster spec. The CAPA operator will not attempt to assume an IAM role and instead will use the AWS
credentials provided in its deployment template (using Kiam, environment variables or some other method of obtaining
credentials).

#### Story 3

ACME Industries is offering Kubernetes as a service to other organizations, and follows AWS guidelines for
Software-as-a-Service. This means they want to use cross-account role assumptions for access to customer systems. ACME
Industries also wants to prevent cross-organisation attacks, so they use “external IDs” to prevent confused deputy
scenarios. ACME Industries wants to minimise the memory footprint of managing many clusters, and wants to move to having
a single instance of CAPA to cover multiple organisations.

## Requirements

### Functional

<a name="FR1">FR1.</a> CAPA MUST support IAM role assumption using the STS::AssumeRole API.

<a name="FR2">FR2.</a> CAPA MUST support session names and external ID to prevent confused deputy attacks.

<a name="FR3">FR3.</a> CAPA MUST support static credentials.

<a name="FR4">FR4.</a> CAPA MUST prevent privilege escalation allowing users to create clusters in AWS accounts they should
  not be able to.

<a name="FR5">FR5.</a> CAPA SHOULD support credential refreshing when identity data is modified.

<a name="FR6">FR6.</a> CAPA SHOULD provide validation for identity data submitted by users.

<a name="FR7">FR7.</a> CAPA COULD support role assumption using OIDC projected volume service account tokens.

<a name="FR9">FR9.</a> CAPA MUST support clusterctl move scenarios.

### Non-Functional

<a name="NFR8">NFR8.</a> Each instance of CAPA SHOULD be able to support 200 clusters using role assumption.

<a name="NFR8">NFR9.</a> CAPA MUST call STS APIs only when necessary to prevent rate limiting.

<a name="NFR8">NFR10.</a> Unit tests MUST exist for all credential provider code.

<a name="NFR8">NFR11.</a> e2e tests SHOULD exist for all credential provider code.

<a name="NFR8">NFR12.</a> Credential provider code COULD be audited by security engineers.

### Implementation Details/Notes/Constraints

The current implementation of CAPA requests a new AWS EC2 and ELB service per cluster and per machine and stores these
in fields on the ClusterScope struct. ClusterScopes are reference values which are created per-reconciliation:

```go
type ClusterScope struct {
  logr.Logger
  client      client.Client
  patchHelper *patch.Helper

  AWSClients
  Cluster    *clusterv1.Cluster
  AWSCluster *infrav1.AWSCluster
}
```

The field AWSClients holds the actual AWS service clients, and is defined like so:

```go
type AWSClients struct {
  EC2             ec2iface.EC2API
  ELB             elbiface.ELBAPI
  SecretsManager  secretsmanageriface.SecretsManagerAPI
  ResourceTagging resourcegroupstaggingapiiface.ResourceGroupsTaggingAPIAPI
}
```

The signatures for the functions which create these instances are as follows:

```go
func NewClusterScope(params ClusterScopeParams) (*ClusterScope, error) {
  ...
  return &ClusterScope{
    ...
  }, nil
}
```

#### Proposed Changes

This proposal borrows heavily from the Service APIs specification.

In the initial implementation, all new resources will be scoped at the cluster level, this is to enable delegation of
AWS accounts whilst preventing privilege escalation as per [FR4](#FR4). Reasons for this are documented in the
[alternatives](#alternatives) section.

##### Cluster API Provider AWS v1alpha3 types

<strong><em>Changed Resources</strong></em>
* `AWSCluster`

<strong><em>New Resources</strong></em>

<em>Cluster scoped resources</em>

* `AWSClusterControllerIdentity` represents an intent to use Cluster API Provider AWS Controller credentials for management cluster.
* `AWSClusterStaticIdentity` represents a static AWS tuple of credentials.
* `AWSClusterRoleIdentity` represents an intent to assume an AWS role for cluster management.

<em>Namespace scoped resources</em>

* `AWSServiceAccountIdentity` represents the use of a Kubernetes service account for access
  to AWS using federated identity.

<strong><em>Changes to AWSCluster</em></strong>

A new field is added to the `AWSClusterSpec` to reference a identity. We intend to use `corev1.LocalObjectReference` in
order to ensure that the only objects that can be references are either in the same namespace or are scoped to the
entire cluster.

```go
// AWSIdentityKind defines allowed AWS identity types
type AWSIdentityKind string

type AWSIdentityRef struct {
  Kind AWSIdentityKind `json:"kind"`
  Name string `json:"name"`
}

type  AWSClusterSpec  struct {
  ...
  // +optional
  IdentityRef *AWSIdentityRef `json:"identityRef,omitempty"`
  // AccountID is the AWS Account ID for this cluster
  // +optional
  AccountID *string `json:"identityRef,omitempty"`
```

An example usage would be:

```yaml
---
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: AWSCluster
metadata:
  name: "test"
  namespace: "test"
spec:
  region: "eu-west-1"
  identityRef:
    kind: AWSClusterRoleIdentity
    name: test-account-role
---
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: AWSClusterRoleIdentity
metadata:
  name: "test-account-role"
spec:
...
```

The `IdentityRef` field will be mutable in order to support `clusterctl move`
scenarios where a user instantiates a cluster on their laptop and then makes
the cluster self-managed.

<strong><em>Identity CRDs</em></strong>

<em>Common elements</em>

All AWSCluster*Identity types will have Spec structs with an `AllowedNamespaces`
field as follows:

```go

type AWSClusterIdentitySpec struct {
// AllowedNamespaces is used to identify which namespaces are allowed to use the identity from.
// Namespaces can be selected either using an array of namespaces or with label selector.
// An empty allowedNamespaces object indicates that AWSClusters can use this identity from any namespace.
// If this object is nil, no namespaces will be allowed (default behaviour, if this field is not provided)
// A namespace should be either in the NamespaceList or match with Selector to use the identity.
//
// +optional
AllowedNamespaces *AllowedNamespaces `json:"allowedNamespaces"`
}

type AllowedNamespaces struct {
// An nil or empty list indicates that AWSClusters cannot use the identity from any namespace.
//
// +optional
// +nullable
NamespaceList []string `json:"list"`

// AllowedNamespaces is a selector of namespaces that AWSClusters can
// use this ClusterIdentity from. This is a standard Kubernetes LabelSelector,
// a label query over a set of resources. The result of matchLabels and
// matchExpressions are ANDed.
//
// An empty selector indicates that AWSClusters cannot use this
// AWSClusterIdentity from any namespace.
// +optional
Selector metav1.LabelSelector `json:"selector"`
}

```

All identities based around AWS roles will have the following fields in their
spec, as per FR1, :

```go
type AWSRoleSpec struct {
  // The Amazon Resource Name (ARN) of the role to assume.
  // +kubebuilder:validation:Pattern:=[\u0009\u000A\u000D\u0020-\u007E\u0085\u00A0-\uD7FF\uE000-\uFFFD\u10000-\u10FFFF]+
  RoleArn string `json:"roleARN"`
  // An identifier for the assumed role session
  // +kubebuilder:validation:Pattern:=[\w+=,.@-]*
  SessionName string `json:"sessionName,omitempty"`
  // The duration, in seconds, of the role session before it is renewed.
  // +kubebuilder:validation:Minimum:=900
  // +kubebuilder:validation:Maximum:=43200
  DurationSeconds uint `json:"durationSeconds,omitempty"`
  // An IAM policy in JSON format that you want to use as an inline session policy.
  // https://docs.aws.amazon.com/STS/latest/APIReference/API_AssumeRole.html
  // +kubebuilder:validation:Pattern:=[\u0009\u000A\u000D\u0020-\u00FF]+
  InlinePolicy string `json:"inlinePolicy,omitempty"`

  // The Amazon Resource Names (ARNs) of the IAM managed policies that you want
  // to use as managed session policies.
  // The policies must exist in the same account as the role.
  PolicyARNs []string `json:"policyARNs,omitempty"`
}
```

<em>AWSClusterControllerIdentity</em>

Supporting [FR4](#FR4) by restricting the usage of controller credentials only from `allowedNamespaces`.
`AWSClusterControllerIdentity` resource will be a singleton and this will be enforced by OpenAPI checks.
This instance's creation is automated by a controller for not affecting existing AWSClusters. For details, see [Upgrade Strategy](#upgrade-strategy)

```go

// AWSClusterControllerIdentity represents an intent to use Cluster API Provider AWS Controller credentials for management cluster
// and restricts the usage of it by namespaces.

type AWSClusterControllerIdentity struct {
  metav1.TypeMeta   `json:",inline"`
  metav1.ObjectMeta `json:"metadata,omitempty"`

  // Spec for this AWSClusterControllerIdentity.
  Spec AWSClusterControllerIdentitySpec `json:"spec,omitempty""`
}

type AWSClusterControllerIdentitySpec struct {
  AWSClusterIdentitySpec
}
```

Example:

```yaml
---
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: AWSCluster
metadata:
  name: "test"
  namespace: "test"
spec:
  region: "eu-west-1"
  identityRef:
    kind: AWSClusterControllerIdentity
    name: default
---
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: AWSClusterControllerIdentity
metadata:
  name: "default"
spec:
  allowedNamespaces:{}  # matches all namespaces
```


<em>AWSClusterStaticIdentity</em>

Supporting [FR3](#FR3).

```go

// AWSClusterStaticIdentity represents a reference to an AWS access key ID and
// secret access key, stored in a secret.
type AWSClusterStaticIdentity struct {
  metav1.TypeMeta   `json:",inline"`
  metav1.ObjectMeta `json:"metadata,omitempty"`

  // Spec for this AWSClusterStaticIdentity.
  Spec AWSClusterStaticIdentitySpec `json:"spec,omitempty""`
}

type AWSClusterSecretReference struct {
  // Namspace is where the Secret is located
  Namespace string `json:"string"`
  // Name is the resource name of the secret
  Name string `json:"name"`
}

type AWSClusterStaticIdentitySpec struct {
  AWSClusterIdentitySpec
  // Reference to a secret containing the credentials. The secret should
  // contain the following data keys:
  //  AccessKeyID: AKIAIOSFODNN7EXAMPLE
  //  SecretAccessKey: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
  //  SessionToken: Optional
  SecretRef AWSClusterSecretReference `json:"secretRef"`
}
```

Example:

```yaml
---
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: AWSCluster
metadata:
  name: "test"
  namespace: "test"
spec:
  region: "eu-west-1"
  identityRef:
    kind: AWSClusterStaticIdentity
    name: test-account
---
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: AWSClusterStaticIdentity
metadata:
  name: "test-account"
spec:
  secretRef:
    name: test-account-creds
    namespace: capa-system
      clusterSelector:
  allowedNamespaces:
    selector:
      matchLabels:
        ns: "testlabel"
---
apiVersion: v1
kind: Namespace
metadata:
  labels:
    cluster.x-k8s.io/ns: "testlabel"
  name: "test"
---
apiVersion: v1
kind: Secret
metadata:
  name: "test-account-creds"
  namespace: capa-system
stringData:
 accessKeyID: AKIAIOSFODNN7EXAMPLE
 secretAccessKey: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
```

<em>AWSClusterRoleIdentity</em>

`AWSClusterRoleIdentity` allows CAPA to assume a role either in the same
or another AWS account, using the STS::AssumeRole API, meeting [FR 1](#FR1).

```go
type AWSClusterRoleIdentity struct {
  metav1.TypeMeta   `json:",inline"`
  metav1.ObjectMeta `json:"metadata,omitempty"`

  // Spec for this AWSClusterRoleIdentity.
  Spec AWSClusterRoleIdentitySpec `json:"spec,omitempty""`
}

// AWSClusterIdentityKind defines allowed cluster-scoped AWS identity types
type AWSClusterIdentityKind AWSIdentityKind

type AWSClusterIdentityReference struct {
  Kind AWSClusterIdentityKind `json:"kind"`
  Name string `json:"name"`
}

type AWSClusterRoleIdentitySpec struct {
  AWSClusterIdentitySpec
  AWSRoleSpec
  // A unique identifier that might be required when you assume a role in another account.
  // If the administrator of the account to which the role belongs provided you with an
  // external ID, then provide that value in the ExternalId parameter. This value can be
  // any string, such as a passphrase or account number. A cross-account role is usually
  // set up to trust everyone in an account. Therefore, the administrator of the trusting
  // account might send an external ID to the administrator of the trusted account. That
  // way, only someone with the ID can assume the role, rather than everyone in the
  // account. For more information about the external ID, see How to Use an External ID
  // When Granting Access to Your AWS Resources to a Third Party in the IAM User Guide.
  // +optional
  // +kubebuilder:validation:Pattern:=[\w+=,.@:\/-]*
  ExternalID *string `json:"externalID,omitempty"`

  // SourceIdentityRef is a reference to another identity which will be chained to do
  // role assumption.
  SourceIdentityRef AWSClusterIdentityReference `json:"sourceIdentityRef,omitempty"`
}

```

Implementation of `ExternalID` supports [FR2](#FR2).

Example:

```yaml
---
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: AWSCluster
metadata:
  name: "test"
  namespace: "test"
spec:
  region: "eu-west-1"
  identityRef:
    kind: AWSClusterRoleIdentity
    name: test-account-role
---
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: AWSClusterRoleIdentity
metadata:
  name: "test-account-role"
spec:
  allowedNamespaces:
    list: # allows only "test" namespace to use this identity
      "test"
  roleARN: "arn:aws:iam::123456789:role/CAPARole"
  # Optional settings
  sessionName: "cluster-spinner.acme.com"
  externalID:  "7a5b816a-7743-4377-a382-2d695bf1f172"
  policyARNs: '["RestrictedCAPAPolicy"]'
  inlinePolicy: '{"Version": "2012-10-17","Statement":[...]}'
  tags: '["company": "acme industries", "product": "cluster-spinner"]'
  transitiveTags: '["company", "product"]'
  sourceIdentityRef:
    kind: AWSClusterStaticIdentity
    name: test-account-creds
```

<em>Future implementation: AWSServiceAccountIdentity</em>

This would not be implemented in the first instance, but opens the possibility to use Kubernetes service accounts
together with `STS::AssumeRoleWithWebIdentity`, supporting [FR7](#FR7).

Definition:

```go
type AWSServiceAccountIdentity struct {
  metav1.TypeMeta   `json:",inline"`
  metav1.ObjectMeta `json:"metadata,omitempty"`

  // Spec for this AWSServiceAccountIdentity.
  Spec AWSServiceAccountIdentitySpec `json:"spec,omitempty""`
}

type AWSServiceAccountIdentitySpec struct {
  AWSRoleSpec

  // Audience is the intended audience of the token. A recipient of a token
  // must identify itself with an identifier specified in the audience of the
  // token, and otherwise should reject the token. The audience defaults to
  // sts.amazonaws.com
  // +default="[sts.amazonaws.com]"
  Audiences []string

  // ExpirationSeconds is the requested duration of validity of the request. The
  // token issuer may return a token with a different validity duration so a
  // client needs to check the 'expiration' field in a response.
  // +optional
  // +default=86400
  ExpirationSeconds int `json:ExpirationSeconds`
}
```

Because service account subjects are necessarily scoped to a namespace within the Kubernetes RBAC model, it's therefore
ideal to leverage a namespace-scoped CRD for permissions.

Example:

The account owner would be expected to set up an appropriate IAM role with the following trust policy:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Identity": {
        "Federated": "<Provider ARN for the management cluster OIDC configuration>"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "sub": "system:serviceaccount:test:test-service-account",
          "audience": "sts.amazonaws.com",
        }
      }
    }
  ]
}
```

and then apply the following:

```yaml
---
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: AWSCluster
metadata:
  name: "test"
  namespace: "test"
spec:
  region: "eu-west-1"
  identityRef:
    kind: AWSServiceAccountIdentity
    name: test-service-account
---
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: AWSServiceAccountIdentity
metadata:
  name: "test-service-account"
  namespace: "test"
spec:
  audiences:
  - "sts.amazonaws.com"
  roleARN: "arn:aws:iam::123456789:role/CAPARole"
```

With the CAPA controller running with the appropriate permissions to create service accounts in the `"test"` namespace,
the controller would then request a POST request to `/api/v1/namespaces/test/serviceaccounts/test-service-account/token`
with the requisite parameters and receive a new JWT token to use with `STS::AssumeRoleWithWebIdentity`.

### Controller Changes

* If identityRef is specified, the CRD is fetched and unmarshalled into a  AWS SDK credential.Provider for the identity type.
* The controller will compare the hash of the credential provider against the same secret’s provider in a cache ([NFR 8](#NFR8)).
* The controller will take the newer of the two and instantiate AWSClients with the selected credential provider.
* The controller will set an the identity resource as one of the OwnerReferences of the AWSCluster.
* The controller and defaulting webhook will default `nil` `identityRef` field in AWSClusters to `AWSClusterControllerIdentity`.

This flow is shown below:

![](20200506-single-controller-multitenancy-flow.svg)

### Clusterctl changes

Today, `clusterctl move` operates by tracking objectreferences within the same namespace, since we are now proposing to
use cluster-scoped resources, we will need to add requisite support to clusterctl's object graph to track ownerReferences
pointing at cluster-scoped resources, and ensure they are moved. We will naively not delete cluster-scoped resources
during a move, as they maybe referenced across namespaces. Because we will be tracking the AWS Account ID for a given identity, it is expected, that for this proposal, this provides sufficient protection against the possibility of a cluster-scoped
identity after one move, and being copied again.


#### Validating webhook changes

A validating webhook could potentially handle some of the cross-resource validation necessary for the [security
model](#security-model) and provide more immediate feedback to end users. However, it would be imperfect. For example, a
change to a `AWSCluster*Identity` could affect the validity of corresponding AWSCluster.
The singleton `AWSClusterControllerIdentity` resource will be immutable to avoid any unwanted overrides to the allowed namespaces, especially during upgrading clusters.

#### Identity Type Credential Provider Behaviour

Implementations for all identity types will implement the `credentials.Provider` interface in the AWS SDK to support [FR5](#FR5) as well as an
additional function signature to support caching:

```go
type Provider interface {
 // Retrieve returns nil if it successfully retrieved the value.
 // Error is returned if the value were not obtainable, or empty.
 Retrieve() (Value, error)

 // IsExpired returns if the credentials are no longer valid, and need
 // to be retrieved.
 IsExpired() bool
}

type AWSIdentityTypeProvider interface {
 credentials.Provider
 // Hash returns a unique hash of the data forming the credentials
 // for this identity
 Hash() (string, error)
}
```

AWS client sessions are structs implementing the provider interface. Every SDK call will call `IsExpired()` on the
credential provider. If `IsExpired()` returns true, `Retrieve()` is called to refresh the credential. Where the controller
is using a temporary credential, the “real” credential provider will be nested, such that `IsExpired()` calls `IsExpired()`
on the nested provider. This also allows the end user to chain role identities together if needed.

The controller will maintain a cache of all the custom credential providers that are referenced by an AWSCluster,
similar as is done to sessions at present.

In the first instance, the controller will refresh secrets on each reconciliation run, and store a hash of the secret
contents in the cache. If the hash changes, then the credential will be invalidated, forcing `IsExpired()` to return
true and a consuming EC2, ELB or SecretsManager client will therefore issue a `Retrieve()` on the credential provider
forcing a refresh.

The authors have implemented a similar mechanism [based on the Ruby AWS SDK][aws-assume-role].

This could be further optimised by having the controller maintain a watch on all Secrets matching the identity types.
Upon receiving an update event, the controller will update lookup the key in the cache and update the relevant provider.
This may be implemented as its own interface. Mutexes will ensure in-flight updates are completed prior to SDK calls are
made. This would require changes to RBAC, and maintaining a watch on secrets of a specific type will require further
investigation as to feasibility.

### Security Model

The intended RBAC model mirrors that for Service APIs:

#### Roles

For the purposes of this security model, 3 common roles have been identified:

* **Infrastructure provider**: The infrastructure provider (infra) is responsible for the overall environment that
  the cluster(s) are operating in or the PaaS provider in a company.

* **Management cluster operator**: The cluster operator (ops) is responsible for
  administration of the Cluster API management cluster. They manage policies, network access,
  application permissions.

* **Workload cluster operator**: The workload cluster operator (dev) is responsible for
  management of the cluster relevant to their particular applications .

There are two primary components to the Service APIs security model: RBAC and namespace restrictions.

### RBAC
RBAC (role-based access control) is the standard used for Kubernetes
authorization. This allows users to configure who can perform actions on
resources in specific scopes. RBAC can be used to enable each of the roles
defined above. In most cases, it will be desirable to have all resources be
readable by most roles, so instead we'll focus on write access for this model.

##### Write Permissions
|                              | AWSCluster*Identity | AWSServiceAccountIdentity | AWS IAM API | Cluster |
| ---------------------------- | ------------------- | ------------------------- | ----------- | ------- |
| Infrastructure Provider      | Yes                 | Yes                       | Yes         | Yes     |
| Management Cluster Operators | Yes                 | Yes                       | Yes         | Yes     |
| Workload Cluster Operator    | No                  | Yes                       | No          | Yes     |

Because Kubernetes service accounts necessarily encode the namespace into the JWT subject, we can allow workload cluster
operators to create their own `AWSServiceAccountIdentities`. Whether they have actual permissions on AWS IAM to set up
the trust policy and the management cluster as the federated identity provider then becomes a problem external to
Kubernetes.

### Namespace Restrictions
The extra configuration options are not possible to control with RBAC. Instead,
they will be controlled with configuration fields on GatewayClasses:

* **allowedNamespaces**: This field is a selector of namespaces that
  Gateways can use this `AWSCluster*Identity` from. This is a standard Kubernetes
  LabelSelector, a label query over a set of resources. The result of
  matchLabels and matchExpressions are ANDed. CAPA will not support
  AWSClusters in namespaces outside this selector. An empty selector (default)
  indicates that AWSCluster can use this `AWSCluster*Identity` from any namespace. This
  field is intentionally not a pointer because the nil behavior (no namespaces)
  is undesirable here.


### CAPA Controller Requirements
The CAPA controller will need to:

* Populate condition fields on AWSClusters and indicate if it is
  compatible with `AWS*Identity`.
* Not implement invalid configuration. Fore example, if a `AWSCluster*Identity` is referenced in
  an invalid namespace for it, it should be ignored.
* Respond to changes in `AWS*Identity` configuration that may change.

### Alternative Approaches Considered

#### Using only secrets or RoleARN field on AWSCluster

In an earlier iteration of this proposal, it was proposed to only use secretReferences in the `AWSCluster` object, and
then using the `Type` field on the Secret to disambiguate between the different orincipal types.

**Benefits**

* Re-using secrets ensures encryption by default where it is KMS encryption is used, and additionally provides a clear UX
signal to end users that the data is meant to be kept secure.

<em>Mitigations for current proposal</em>

This is traded off in this proposal with ensuring static credentials are stored in secrets, whilst
allowing role ARNs and other parameters to be stored in CRDs.

Where cluster operators are using encryption providers for the Kubernetes API Server, they can
optionally specify that these new resources are encrypted at rest through configuration of the
API server KMS.

**Downsides**

By allowing workload cluster operators to create the various identitys, or referencing them directly in the
AWSCluster as a field they could potentially escalate privilege when assuming role across AWS accounts as the CAPA
controller itself may be running with privileged trust.

#### 1:1 mapping one namespace to one AWSIdentity

The mapping of a singular AWSIdentity to a single namespace such that there is a 1:1 mapping
was considered, via either some implicitly named secret or other metadata on the namespace.

**Benefits**

* Can potentially resolve an issue where multiple clusters with the same name in an AWS account cause problems for the [AWS Cloud Provider](#aws-cloud-provider-behaviour-with-regards-to-cluster-names).

**Downsides**

* Implicit link of credentials to clusters is harder for introspection.
* Doesn't necessarily guarantee uniqueness of cluster names in an account

### Risks and Mitigations

#### Network assumptions are made explicit

This change maintains but makes explicit the expectation that there is network access between the account that the CAPI
operator resides in to an API server endpoint in an account where a reconciled cluster resides. Existing pre-flight
checks would not confirm this. The existing pattern for the CAPA operator to create security groups if they do not exist
may need to account for this network access requirement. Currently, when CAPA creates a security group for cluster
control plane load balancer it allows ingress from any CIDR block. However the security groups constraining the CAPI
operator would require appropriate egress rules to be able to access load balancers in other AWS accounts. The extent to
which CAPA can solve for this needs to be determined. However, this risk has already been present and exposed doing role
assumption using another method such as Kiam.

#### Caching and handling refresh of credentials

For handling many accounts, the number of calls to the STS service must be minimised. This is currently implemented as a
cache on cluster key in Cluster API Provider AWS.

#### AWS Cloud Provider behaviour with regards to cluster names

At present, multiple clusters deployed into the same AWS account with the same
cluster name cause the AWS Cloud Provider to malfunction due to the way
tags are looked up. This has driven some concern to allowing multiple
namespaces access to the same AWS account.

However, since it is possible to configure multiple credentials against the same
account, there's no easy way to guarantee uniqueness.

This could be mitigated by having CAPA call `STS::GetCallerIdentity` on entry
to the cluster, and then comparing against all current clusters declared across
the management cluster. We would probably need to store the AccountID into
the AWSCluster object and index the field.

A better approach maybe to modify the cloud provider to support customisation
of tag lookups such that multiple clusters with the same name, but in
different namespaces do not break the AWS Cloud Provider.

## Upgrade Strategy

The data changes are additive and optional, except `AWSClusterControllerIdentity`.
`AWSClusterControllerIdentity` singleton instance restricts the usage of controller credentials only from `allowedNamespaces`.
AWSClusters that do not have an assigned `IdentityRef` is defaulted to use `AWSClusterControllerIdentity`, hence existing clusters needs to have
`AWSClusterControllerIdentity` instance. In order to make existing AWSClusters to continue to reconcile as before, a new controller is added as experimental feature
and gated with **Feature gate:** AutoControllerIdentityCreator=true. By default, this feature is enabled. This controller creates `AWSClusterControllerIdentity` singleton instance (if missing) that allows all namespaces to use the identity.
During v1alpha4 releases, since breaking changes will be allowed, this feature will become obsolete.

## Additional Details

### Test Plan

* Unit tests to validate that the cluster controller can reconcile an AWSCluster when IAMRoleARN field is nil, or provided.
* Unit tests to ensure pre-flight checks are performed relating to IAM role assumption when IAMRoleARN is provided.
* Propose performing an initial sts:AssumeRole call and fail pre-flight if this fails.
* e2e test for role assumption ([NFR12](#NFR12)).
* If it can be supported in Prow environment, additional e2e test for OIDC-based role assumption.
* clusterctl e2e test with a move of a self-hosted cluster using a identityRef.

### Graduation Criteria

#### Alpha

* Support role assumption with external ID as a aws-sdk-go `credentials.Provider`.
* Ensure `clusterctl move` works with the mechanism.

#### Beta

* Support OIDC and chained role assumption.
* Admission controller validation for secrets of type.
* Full e2e coverage.

#### Stable

* Two releases since beta.

These may be defined in terms of API maturity, or as something else. Initial proposal should keep this high-level with a
focus on what signals will be looked at to determine graduation.

### Version Skew Strategy

Most of the version skew is contained within the secret type. We explicitly add a `v1alpha3` to the secret type allowing
the format to change over time. However, we do not have the benefit of hub and spoke storage conversion that would exist
with a CRD. We propose, that in the event of a secret type signature change:

* The old version is accepted for 2 releases.
* For the n+3 and n+4 releases, the Cluster API Provider AWS controller will convert the storage version.
* On release n+5, remaining code for the old type signatures are deleted.

## Implementation History

- [ ] 2020/03/30: Initial proposal
- [ ] 2020/04/02: Revised proposal as [Google doc][google-doc]
- [ ] 2020/05/01: Presented proposal at meeting with stakeholders
- [ ] 2020/05/06: Open proposal PR

<!-- Links -->
[aws-sdk-go-credential-provider]: https://github.com/aws/aws-sdk-go/blob/master/aws/credentials/credentials.go#L35
[kiam]: https://github.com/uswitch/kiam
[aws-assume-role]: https://github.com/scalefactory/aws-assume-role
[google-doc]: https://docs.google.com/document/d/1vwjvWc-RIZfwDXb-3tW4gKBOaPCKhripfhXshOmB5z4/edit#heading=h.13bi7vgop9tn
