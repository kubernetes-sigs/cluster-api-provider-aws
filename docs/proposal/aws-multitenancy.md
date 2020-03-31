
---
title: AWS Account Multi-Tenancy
authors:
  - "@andrewmyhre"
reviewers:
  - ""
creation-date: 2020-03-30
last-updated: 2020-03-30
status: provisional
see-also:
  - "[https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/1552](https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/1552)"
replaces:
  - ""
superseded-by:
  - ""
---

# Federated AWS Account Support

<!-- BEGIN Remove before PR -->
To get started with this template:
1. **Make a copy of this template.**
  Copy this template into `docs/enhacements` and name it `YYYYMMDD-my-title.md`, where `YYYYMMDD` is the date the proposal was first drafted.
1. **Fill out the required sections.**
1. **Create a PR.**
  Aim for single topic PRs to keep discussions focused.
  If you disagree with what is already in a document, open a new PR with suggested changes.

The canonical place for the latest set of instructions (and the likely source of this file) is [here](/proposals/YYYYMMDD-proposal-template.md).

The `Metadata` section above is intended to support the creation of tooling around the proposal process.
This will be a YAML section that is fenced as a code block.
See the proposal process for details on each of these items.

<!-- END Remove before PR -->

## Table of Contents

A table of contents is helpful for quickly jumping to sections of a proposal and for highlighting
any additional information provided beyond the standard proposal template.
[Tools for generating](https://github.com/ekalinin/github-markdown-toc) a table of contents from markdown are available.

- [Title](#title)
  - [Table of Contents](#table-of-contents)
  - [Summary](#summary)
  - [Motivation](#motivation)
    - [Goals](#goals)
    - [Non-Goals](#non-goals)
  - [Proposal](#proposal)
    - [User Stories [optional]](#user-stories-optional)
      - [Story 1](#story-1)
      - [Story 2](#story-2)
    - [Implementation Details/Notes/Constraints [optional]](#implementation-detailsnotesconstraints-optional)
    - [Risks and Mitigations](#risks-and-mitigations)
  - [Design Details](#design-details)
    - [Test Plan](#test-plan)
    - [Graduation Criteria](#graduation-criteria)
      - [Examples](#examples)
        - [Alpha -> Beta Graduation](#alpha---beta-graduation)
        - [Beta -> GA Graduation](#beta---ga-graduation)
        - [Removing a deprecated flag](#removing-a-deprecated-flag)
    - [Upgrade / Downgrade Strategy](#upgrade--downgrade-strategy)
    - [Version Skew Strategy](#version-skew-strategy)
  - [Implementation History](#implementation-history)
  - [Drawbacks [optional]](#drawbacks-optional)
  - [Alternatives [optional]](#alternatives-optional)
  - [Infrastructure Needed [optional]](#infrastructure-needed-optional)

## Glossary

Refer to the [Cluster API Book Glossary](https://cluster-api.sigs.k8s.io/reference/glossary.html).

If this proposal adds new terms, or defines some, make the changes to the book's glossary when in PR stage.

## Summary
The CAPA operator is able to manage cloud infrastructure within the permission scope of the AWS principle it is initialized with. It is expected that the CAPA operator will be provided credentials via the deployment, either explicitly or implicitly, which will be static for the lifetime of the pod. The CAPA operator is unable to manage cloud infrastructure in a different AWS account than the one it is initially configured to use.

This proposal outlines new capabilities for CAPA to use IAM role assumption to assume a permission set in a different AWS account, at runtime, on a per-cluster basis. The proposed changes would be fully backwards compatible and maintain the existing behavior with no changes to user configuration required.

## Motivation
For large organizations, especially highly-regulated organizations, there is a need to be able separate duties at various levels of infrastructure - permissions, networks and accounts. [VPC sharing]([https://aws.amazon.com/blogs/networking-and-content-delivery/vpc-sharing-a-new-approach-to-multiple-accounts-and-vpc-management/](https://aws.amazon.com/blogs/networking-and-content-delivery/vpc-sharing-a-new-approach-to-multiple-accounts-and-vpc-management/)) is a model which provides separation at the AWS account level. Within this model it is appropriate for tooling running within the 'Owner' account to manage infrastructure within the 'Participant' accounts, which requires a principle in the Owner account which can assume a principle within a Participant account. For CAPA to be most useful within these organizations it will need to support multi-account models.

Because a single deployment of the CAPA operator may reconcile many different clusters in its lifetime, it is necessary to modify the CAPA operator to scope its AWS client instances to within the reconciliation process.

It follows that an organization may wish to provision control planes and worker groups in separate accounts (including each worker group in a separate account). There is a desire to support this configuration also.

### Goals

1. To enable cluster reconciliation across AWS account boundaries
2. To maintain backwards compatibility and cause no impact for users who don't intend to make use of this capability

### Non-Goals/Future Work

- Enabling Machines to be provisioned in AWS accounts different than their control planes. This would require adding various AWS infrastructure specification to the AWSMachineTemplate type which currently does not exist.

## Proposal

![Diagram](https://user-images.githubusercontent.com/208667/78036771-b299ad00-7338-11ea-90b2-1a88702fb1dd.png)

The current implementation of CAPA requests a new AWS EC2 and ELB service per cluster and per machine and stores these in fields on the ClusterScope struct. ClusterScopes are reference values which are created per-reconciliation:

```golang
type ClusterScope struct {  
   logr.Logger  
  client      client.Client  
  patchHelper *patch.Helper  
  
  AWSClients  
  Cluster    *clusterv1.Cluster  
  AWSCluster *infrav1.AWSCluster  
}
```

The field `AWSClients` holds the actual AWS service clients, and is defined like so:
```golang
type AWSClients struct {  
  EC2             ec2iface.EC2API  
  ELB             elbiface.ELBAPI  
  SecretsManager  secretsmanageriface.SecretsManagerAPI  
  ResourceTagging resourcegroupstaggingapiiface.ResourceGroupsTaggingAPIAPI  
}
```

The signatures for the functions which create these instances are as follows:

```golang
func NewClusterScope(params ClusterScopeParams) (*ClusterScope, error) {
  ...
  return &ClusterScope{
    ...
  }, nil
}
```

### Proposed changes:
#### Data Model
api/v1alpha3/awscluster_types.go:
```golang
type  AWSClusterSpec  struct {
  + IAMRoleARN *string `json:"iamRoleArn,omitempty"`
```

#### Code
- Add a field to the AWSClients struct: `IAMRoleARN *string`
- Add an optional field to the AWSClusterSpec type named `IAMRoleARN`
- When creating a ClusterScope:
	- if the IAMRoleARN has been specified, create EC2 and ELB services with an [STS credentials object configured to use the provided IAM role ARN]([https://docs.aws.amazon.com/sdk-for-go/api/aws/credentials/stscreds/](https://docs.aws.amazon.com/sdk-for-go/api/aws/credentials/stscreds/))
	- if the IAMRoleARN has not been specified then create EC2 and ELB services with no credentials specified (preserve existing behavior)

#### Examples
cluster.yaml
```yaml
apiVersion: cluster.x-k8s.io/v1alpha3
kind: Cluster
metadata:
  name: my-cluster
spec:
  infrastructureRef:
    apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
    kind: AWSCluster
    name: my-cluster
```
awscluster.yaml
```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: AWSCluster
metadata:
  name: my-cluster
spec:
  iamRoleArn: arn:aws:iam::123456789:role/CAPARole
  networkSpec:
    ...
```

### User Stories

- Detail the things that people will be able to do if this proposal is implemented.
- Include as much detail as possible so that people can understand the "how" of the system.
- The goal here is to make this feel real for users without getting bogged down.

#### Story 1
Alex is an engineer in a large organization which has a strict AWS account architecture. This architecture dictates that Kubernetes clusters must be hosted in dedicated AWS accounts. The organization has adopted Cluster API in order to manage Kubernetes infrastructure, and expects 'Owner' clusters running the Cluster API operators to manage 'Participant' clusters in dedicated AWS accounts. 

The current configuration exists:
AWS Account 'Owner':
- Vpc, subnets shared with 'Participant' AWS Account
- IAM role 'ClusterAPI-Owner'
- Kubernetes Cluster 'Owner' running Cluster API operators
	- CAPA operator is using the 'ClusterAPI-Owner' IAM principle

AWS Account 'Participant':
- Vpc and subnets provided by 'Owner' AWS Account
- IAM Role 'ClusterAPI-Participant' which trusts IAM role 'ClusterAPI-Owner'

Alex can provision a new cluster in the 'Participant' AWS Account by creating new Cluster API resources in Kubernetes Cluster 'Owner. Alex specifies the IAM role 'ClusterAPI-Participant' in the AWSCluster spec. The CAPA operator in Kubernetes Cluster 'Owner' assumes the role 'ClusterAPI-Participant' when reconciling the AWSCluster so that it can create/use/destroy resources in the 'Participant' AWS Account.

#### Story 2
Dascha is an engineer in a smaller, less strict organization with a few AWS accounts intended to host all infrastructure. There is a single AWS account named 'dev', and Dascha wants to provision a new cluster in this account. An existing Kubernetes cluster is already running the Cluster API operators and managing resources in the dev account. Dascha can provision a new cluster by creating Cluster API resources in the existing cluster, omitting the IAMRoleARN field in the AWSCluster spec. The CAPA operator will not attempt to assume an IAM role and instead will use the AWS credentials provided in its deployment template (using Kiam, environment variables or some other method of obtaining credentials).

### Implementation Details/Notes/Constraints
- Pre-flight checks for subnets, security groups and IAM roles would be performed as normal, though possibly in another AWS account than the operator resides in. 

### Risks and Mitigations
- This change adds a new expectation that there is network access between the account that the CAPI operator resides in to the account where a reconciled cluster resides. Existing pre-flight checks would not confirm this. The existing pattern for the CAPA operator to create security groups if they do not exist may need to account for this network access requirement. Currently, when CAPA creates a security group for cluster control plane load balancer it allows ingress from any CIDR block. However the security groups constraining the CAPI operator would require appropriate egress rules to be able to access load balancers in other AWS accounts. The extent to which CAPA can solve for this needs to be determined.


## Alternatives
None.

## Upgrade Strategy
The data changes are additive and optional, so existing AWSCluster specifications will continue to reconcile as before. These changes will only come into play when an IAM role is provided in the new field in `AWSClusterSpec`. Upgrades to versions with this new field will be break.

It may be worthwhile to create a new structure for the IAMRoleARN field, rather than adding it directly to the `AWSClusterSpec` type. This would allow extending the spec in the future for other optional configuration to be provided to the `aws.Config{}` type and reduce the risk of breakages caused by future changes to this new field.

## Additional Details

### Test Plan [optional]

- Unit tests to validate that the cluster controller can reconcile an AWSCluster when IAMRoleARN field is nil, or provided
- Unit tests to ensure pre-flight checks are performed relating to IAM role assumption when IAMRoleARN is provided
	- Propose performing an initial sts:AssumeRole call and fail pre-flight if this fails

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



## Implementation History

- [ ] 03/30/2020: Open proposal PR

<!-- Links -->
[community meeting]: https://docs.google.com/document/d/1Ys-DOR5UsgbMEeciuG0HOgDQc8kZsaWIWJeKJ1-UfbY
