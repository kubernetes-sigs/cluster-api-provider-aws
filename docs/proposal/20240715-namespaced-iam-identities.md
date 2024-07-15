---
title: Namespaced Multitenancy
authors:
  - "@tjamet"
reviewers:
creation-date: 2024-07-15
last-updated: 2024-07-15
status: draft
see-also:
- https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/main/docs/proposal/20200506-single-controller-multitenancy.md
replaces: []
superseded-by: []
---

# Namespaced Multitenancy

- [Glossary](#glossary)
- [Summary](#summary)
- [Motivation](#motivation)
- [Goals](#goals)
- [Proposal](#proposal)
   * [User Story](#user-story)
- [Functional Requirements](#functional-requirements)
- [Implementation Details/Notes/Constraints](#implementation-detailsnotesconstraints)
   * [New namespaced resources](#new-namespaced-resources)
- [Security considerations](#security-considerations)
   * [Privilege escalation prevention deep dive](#privilege-escalation-prevention-deep-dive)
      + [AWSRoleIdentity case](#awsroleidentity-case)
      + [AWSStaticIdentity case](#awsstaticidentity-case)

## Glossary

* Identity Type - One of several ways to provide a form of identity that is ultimately resolved to an AWS access key ID,
  secret access key and optional session token tuple.
* Credential Provider - An implementation of the interface specified in the [AWS SDK for
  Go][aws-sdk-go-credential-provider].
* CAPA - An abbreviation of Cluster API Provider AWS.
* CAPA owners - The team responsible to operate the CAPA provider
* Cluster administrators - The team or individuals creating cluster objects to run clusters in their own accounts

## Summary

The CAPA operator is currently capable of offering multi-tenancy at the cluster level.
With the latest changes of the AWS STS AssumeRole API, it is now possible to provide a unique and dynamic identifier to refer to
the external unique identity of the requester ( external username, external resource ID, ... ) as documented in the AWS [SourceIdentity documentation],
and later grant accesses based on it.

This proposal shapes a new capability for CAPA to use namespaced identities while preventing privilege escalation.


## Motivation

In the [single cluster multitenancy] proposal, the functional requirement [FR4] introduced the use of cluster-wide resources, managed by the CAPI maintainers and
hence preventing privilege escalation, through administrator review.

In large organisations favouring autonomy, this brings high responsibility on the team operating CAPA. They need to judge which roles can be used in which namespaces.
This breaks the autonomy principle those organisations have.

In this situation, the current model introduces two sources to trust (the CAPA operator and the team operating it) and reduces the cluster operator autonomy to create
clusters in new accounts.

## Goals

1. To enable AWSIdentity resources granting autonomy to cluster administrators to deploy clusters in their own accounts
2. To enable cluster administrators to allow of forbid AWSIdentities in their accounts

## Proposal

### User Story

Manuela is an infrastructure engineer in a large corporation. The corporation Manuela works in values autonomy and prefers
that the different areas are autonomous to deploy new clusters in their accounts.

Manuela was provided with a cluster where a CAPI installation is maintained for her. She has access to a single namespace of
this cluster. Yet, Manuela needs to isolate the production and non-production workload they run into separate accounts they own.

To respect the autonomy the company management is asking for, Manuela needs to be able to create on her own all the CAPI objects
so she can deploy clusters end-to-end.

## Functional Requirements

<a name="FR1">FR1.</a> CAPA MUST support cluster administrators to autonomously deploy clusters in their own accounts without the need of CAPA owners.

<a name="FR2">FR2.</a> CAPA MUST use the SourceIdentity field to uniquely identify the AWSIdentity objects.

<a name="FR3">FR3.</a> CAPA MUST support static credentials.

<a name="FR4">FR4.</a> CAPA MUST prevent privilege escalation allowing users to create clusters in AWS accounts they should
  not be able to.

<a name="FR5">FR5.</a> CAPA MUST guarantee namespace isolation of namespaced identities. Cross-namespace reference of those objects
  should be denied.

<a name="FR6">FR6.</a> CAPA MUST be backward compatible with cluster wide identities introduced in [single cluster multitenancy].
  Namespaced and Cluster-wide identies must work together.

## Implementation Details/Notes/Constraints


### New namespaced resources

In this proposal we introduce 2 new namespaced resources

* `AWSStaticIdentity` represents a static AWS tuple of credentials.
* `AWSRoleIdentity` represents an intent to assume an AWS role for cluster management.

Those resources **must** only be used in the current namespace and **must not** be usable by `AWSCluster*Identity` to chain AssumeRoles.

They would follow the folowing schemas.

```golang

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// AWSStaticIdentity is the Schema for the awsstaticidentities API
// It represents a reference to an AWS access key ID and secret access key, stored in a secret.
type AWSStaticIdentity struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec for this AWSStaticIdentity
	Spec AWSStaticIdentitySpec `json:"spec,omitempty"`
}

// AWSStaticIdentitySpec defines the specifications for AWSStaticIdentity.
type AWSStaticIdentitySpec struct {
	// Selector allows to restrict the usage of this identity to certain objects based
	// on their labels.
	// This applies to all possible object kinds (AWSRoleIdentity, AWSCluster, AWSManagedControlPlane, ...)
	Selector metav1.LabelSelector `json:"selector"`
	// Reference to a secret containing the credentials. The secret should
	// contain the following data keys:
	//  AccessKeyID: AKIAIOSFODNN7EXAMPLE
	//  SecretAccessKey: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
	//  SessionToken: Optional
	SecretRef string `json:"secretRef"`
}

// AWSRoleIdentity is the Schema for the awsroleidentities API
// It is used to assume a role using the provided sourceRef.
type AWSRoleIdentity struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec for this AWSRoleIdentity.
	Spec AWSRoleIdentitySpec `json:"spec,omitempty"`
}

// AWSRoleIdentitySpec defines the specifications for AWSRoleIdentity.
type AWSRoleIdentitySpec struct {
	// Selector allows to restrict the usage of this identity to certain objects based
	// on their labels.
	// This applies to all possible object kinds (AWSRoleIdentity, AWSCluster, AWSManagedControlPlane, ...)
	Selector metav1.LabelSelector `json:"selector"`
	AWSRoleSpec            `json:",inline"`
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
	ExternalID string `json:"externalID,omitempty"`

	// SourceIdentityRef is a reference to another identity which will be chained to do
	// role assumption. All identity types are accepted.
	SourceIdentityRef *AWSIdentityReference `json:"sourceIdentityRef,omitempty"`
}
```

## Security considerations

This proposal relies on AWS SourceIdentity field which goal is to identify the principal on behalf of which the AssumeRole action is called,
as defined in the [aws-sdk-go-credential-provider].

```golang
// AssumeRoleOptions is the configurable options for AssumeRoleProvider
type AssumeRoleOptions struct {
	// [...]

	// The source identity specified by the principal that is calling the AssumeRole
	// operation. You can require users to specify a source identity when they assume a
	// role. You do this by using the sts:SourceIdentity condition key in a role trust
	// policy. You can use source identity information in CloudTrail logs to determine
	// who took actions with a role. You can use the aws:SourceIdentity condition key
	// to further control access to Amazon Web Services resources based on the value of
	// source identity. For more information about using source identity, see Monitor
	// and control actions taken with assumed roles
	// (https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_temp_control-access_monitor.html)
	// in the IAM User Guide.
	SourceIdentity *string
}
```

By the definition of this field, the `SourceIdentity` field should be managed by CAPA and not exposed to the cluster administrators in any mean.

The `SourceIdentity` may be customisable by the CAPA owners to customise a certain prefix and hence increase the unicity of the requests.
The default `SourceIdentity` field may look like `CAPA:provider:aws:AWSRoleIdentity:identity-namespace:identity-name`. The values `AWSRoleIdentity`, `identity-namespace`
and `identity-name` refer to kubernetes resources and must be injected by the CAPA controller without any posibility to be changed by neither the CAPA owners or the cluster administrators.

By default, when allowing an [AWS principal] to assume a role, it is not allowed to `SetIdentitySource`, and this must be explicitely allowed with the following statement:

```json
{
    "Sid": "AllowRoleToAssumeIdentity",
    "Effect": "Allow",
    "Principal": {
        "AWS": "arn:aws:iam::0123456789:role/my-role"
    },
    "Action": [
        "sts:SetSourceIdentity"
    ]
},
```

Without this statement, any AssumeRole action will be denied with a message similar to `AccessDenied: User: arn:aws:sts::123456789:role/cluster-provider-aws is not authorized to perform: sts:SetSourceIdentity on resource: arn:aws:iam::987654321:role/cluster-provider`.

After setting the field, the role owner will be able to allow a speficic `AWSRoleIdentity` to assume a role with the following [trust relationship policy] statement, as mentioned in the [SourceIdentity documentation] and makes this proposal compliant with [FR4](#FR4).

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "Allow users to set the source identity when using ",
            "Effect": "Allow",
            "Principal": {
                "AWS": "arn:aws:iam::0123456789:role/my-role"
            },
            "Action": [
                "sts:SetSourceIdentity"
            ]
        },
        {
            "Sid": "Only for my namespace",
            "Effect": "Allow",
            "Principal": {
                "AWS": "arn:aws:iam::0123456789:role/my-role"
            },
            "Condition": {
                "StringEquals": {
                    "sts:ExternalId": "the-external-id-set-by-the-identity-object",
                    "sts:SourceIdentity": "CAPA:provider:aws:AWSRoleIdentity:identity-namespace:identity-name"
                }
            },
            "Action": [
                "sts:AssumeRole"
            ]
        }
    ]
}
```

### Privilege escalation prevention deep dive

In both deep dives, we will consider two namespaces `legit-team` and `hacker-team` owned respectively by a team legitimate to manage clusters in account `987654321` and a team trying
to elevate their privileges and break into the `987654321` account where they are not legitimate to manage clusters.

#### AWSRoleIdentity case

The `legit-team` has used the standard `clusterawsadm bootstrap iam create-cloudformation-stack` and hence uses the standard `controllers.cluster-api-provider-aws.sigs.k8s.io` role
to deploy their clusters.

Hence, they have configured an `AWSRoleIdentity` object with the following content

```yaml
kind: AWSRoleIdentity
metadata:
  namespace: legit-team
  name: legit-team-account
spec:
  roleARN: arn:aws:iam::987654321:role/controllers.cluster-api-provider-aws.sigs.k8s.io
  externalID: legit-team-in-kubernetes
```

In their managed controlplane definition, they have referenced the role Identity

```yaml
kind: AWSManagedControlPlane
metadata:
  namespace: legit-team
  name: legit-cluster
spec:
  identityRef:
    name: legit-team-account
    kind: AWSRoleIdentity
```

Because CAPA uses `SourceIdentity` and they are concerned about security concerns, they have set-up their role trust relationships to only allow this `AWSRoleIdentity`
to assume the `controllers.cluster-api-provider-aws.sigs.k8s.io` role.

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "Allow users to set the source identity when using ",
            "Effect": "Allow",
            "Principal": {
                "AWS": "arn:aws:iam::0123456789:role/my-role"
            },
            "Action": [
                "sts:SetSourceIdentity"
            ]
        },
        {
            "Sid": "Only allow the legit-team identity",
            "Effect": "Allow",
            "Principal": {
                "AWS": "arn:aws:iam::0123456789:role/my-role"
            },
            "Condition": {
                "StringEquals": {
                    "sts:ExternalId": "the-external-id-set-by-the-identity-object",
                    "sts:SourceIdentity": "CAPA:provider:aws:AWSRoleIdentity:legit-team:legit-team-account"
                }
            },
            "Action": [
                "sts:AssumeRole"
            ]
        }
    ]
}
```

Meanwhile, the `hacker-team` has discovered the `legit-team` account ID, and is trying to break in.
They considered the `legit-team` would eventually use the default settings and would create both objects in their namespaces:

```yaml
kind: AWSRoleIdentity
metadata:
  namespace: hacker-team
  name: legit-team-account
spec:
  roleARN: arn:aws:iam::987654321:role/controllers.cluster-api-provider-aws.sigs.k8s.io
  externalID: legit-team-in-kubernetes
---
kind: AWSManagedControlPlane
metadata:
  namespace: hacker-team
  name: hacker-cluster
spec:
  identityRef:
    name: legit-team-account
    kind: AWSRoleIdentity
```

Because CAPA has set the `SourceIdentity` field, and the `legit-team` has set the `sts:SourceIdentity` condition, the CAPA operator will not be able to assume the `arn:aws:iam::987654321:role/controllers.cluster-api-provider-aws.sigs.k8s.io` role to deploy the `hacker-cluster` of the `hacker-team` in the `legit-team` account, fulfiling [FR4](#FR4).
The assume role will error with a message like `AccessDenied: User: arn:aws:iam::0123456789:role/capa-role is not authorized to perform: sts:AssumeRole on resource: arn:aws:iam::987654321:role/controllers.cluster-api-provider-aws.sigs.k8s.io`.

Knowing the CAPA implementation, the next thing the `hacker-team` tries is to use directly the `legit-team` `AWSRoleIdentity` in their cluster.

```yaml
kind: AWSManagedControlPlane
metadata:
  namespace: hacker-team
  name: hacker-cluster
spec:
  identityRef:
    name: legit-team-account
    namespace: legit-team
    kind: AWSRoleIdentity
```

Because the `AWSIdentityReference` object does not accept any `namespace` field the `AWSManagedControlPlane` will be denied by the Kubernetes API with the error `strict decoding error: unknown field "spec.identityRef.namespace"` error, complying with [FR4](#FR4).

```golang
type AWSIdentityReference struct {
	// Name of the identity.
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`

	// Kind of the identity.
	// +kubebuilder:validation:Enum=AWSClusterControllerIdentity;AWSClusterRoleIdentity;AWSClusterStaticIdentity
	Kind AWSIdentityKind `json:"kind"`
}
```

#### AWSStaticIdentity case

The `legit-team` uses static credentials to provision cluster and hence creates a `Secret` and `AWSStaticIdentity` in their namespaces to create an `AWSManagedControlPlane`.

```yaml
type: Secret
metadata:
  namespace: legit-team
  name: legit-account-access-keys
dataString:
  AccessKeyID: AKIAIOSFODNN7EXAMPLE
  SecretAccessKey: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
---
kind: AWSStaticIdentity
metadata:
  namespace: legit-team
  name: legit-team-account
spec:
  SecretRef: legit-account-access-keys
---
kind: AWSManagedControlPlane
metadata:
  namespace: legit-team
  name: legit-cluster
spec:
  identityRef:
    name: legit-team-account
    kind: AWSStaticIdentity
```

As they are using plain credentials, CAPA will be authenticated with those credentials and be allowed to create clusters in the aws account.

Meanwhile, the `hacker-team` is trying to break into the `legit-team` account.
Its first attempt is to re-use the `AWSStaticIdentity` from the `legit-team` namespace.

```yaml
kind: AWSManagedControlPlane
metadata:
  namespace: hacker-team
  name: hacker-cluster
spec:
  identityRef:
    name: legit-team-account
    namespace: legit-team
    kind: AWSStaticIdentity
```

Because the `AWSIdentityReference` object does not accept any `namespace` field the `AWSManagedControlPlane` will be denied by the Kubernetes API with the error `strict decoding error: unknown field "spec.identityRef.namespace"` error, complying with [FR4](#FR4).

```golang
type AWSIdentityReference struct {
	// Name of the identity.
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`

	// Kind of the identity.
	// +kubebuilder:validation:Enum=AWSClusterControllerIdentity;AWSClusterRoleIdentity;AWSClusterStaticIdentity
	Kind AWSIdentityKind `json:"kind"`
}
```

The next attent of the `hacker-team` is to use the `legit-team` secret in a `AWSStaticIdentity` in their own namespace.

```yaml
kind: AWSStaticIdentity
metadata:
  namespace: legit-team
  name: legit-team-account
spec:
  SecretRef: legit-account-access-keys
  SecretNamespace: legit-team
---
kind: AWSManagedControlPlane
metadata:
  namespace: hacker-team
  name: hacker-cluster
spec:
  identityRef:
    name: legit-team-account
    kind: AWSStaticIdentity
```

Similarly, the hacker team can't use the legit team secret as it the `AWSStaticIdentity` does not have any Namespace field for the secret.

<!-- Links -->
[aws-sdk-go-credential-provider]: https://github.com/aws/aws-sdk-go-v2/blob/03768e0d0276b360a6abaa4d30318d4aedc44995/credentials/stscreds/assume_role_provider.go#L163
[SourceIdentity documentation]: https://aws.amazon.com/blogs/security/how-to-integrate-aws-sts-sourceidentity-with-your-identity-provider/
[trust relationship policy]: https://aws.amazon.com/blogs/security/how-to-use-trust-policies-with-iam-roles/
[single cluster multitenancy]: https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/main/docs/proposal/20200506-single-controller-multitenancy.md
[FR4]: https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/main/docs/proposal/20200506-single-controller-multitenancy.md#FR4
[AWS principal]: https://medium.com/@reach2shristi.81/aws-principal-vs-identity-3d8eacc5377f
