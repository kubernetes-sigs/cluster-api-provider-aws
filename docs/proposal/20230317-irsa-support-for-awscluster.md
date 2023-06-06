---
title: IRSA Support for Self-Managed Clusters
authors:
  - "@luthermonson"
reviewers:
  - "@richardcase"
  - "@Skarlso"
creation-date: 2023-03-17
last-updated: 2023-03-17
status: provisional
see-also: []
replaces: []
superseded-by: []
---

# Add Support for IRSA to Non-Managed Clusters

## Table of Contents

- [Add Support for IRSA to Non-Managed Clusters](#launch-templates-for-managed-machine-pools)
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
- [CAPI](https://github.com/kubernetes-sigs/cluster-api) - Cluster API.
- [IRSA](https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts.html) - IAM Roles for Service Accounts
- [pod-identity-webhook](https://github.com/aws/amazon-eks-pod-identity-webhook) - Pod Identity Webhook Repo

## Summary
The IAM Roles for Service Accounts take the access control enabled by IAM and bridge the gap to Kubernetes by adding role-based access to service accounts. CAPA users of self-managed clusters can now give granular role-based access to the AWS API at a pod level.      

## Motivation
This functionality is currently built into EKS, with a simple boolean in the AWSManagedCluster API called `AssociateOIDCProvider` CAPA will build an IAM OIDC provider for the cluster and create a trust policy template in a config map to be used for created IAM Roles. Self-managed clusters can use IRSA but require additional manual steps already done in Managed Clusters, including patching kube-api-server, creating an OIDC provider and deploying the `pod-identity-webhook`, which is documented in their [self-hosted setup](https://github.com/aws/amazon-eks-pod-identity-webhook/blob/master/SELF_HOSTED_SETUP.md) walkthrough but with CAPA style ingredients like using the management cluster, kubeadm config modification and the built-in serving certs' OpenID Configuration API endpoints.

The pieces to IRSA are easily created with the existing access for CAPA. By adding `AssociateOIDCProvider` to `AWSCluster` we can kick off a reconciliation process to generate all pieces necessary to utilize IRSA in your self-managed cluster. 

### Goals

1. On cluster creation, add all components to self-managed clusters to use IAM Roles for Service Accounts.
2. On cluster deletion, remove all external dependencies from the AWS account.

### Non-Goals/Future Work
- Migrate all IAM work for Managed cluster to the IAM service. 
- S3 bucket code currently dies when the bucket exists, needs to see if the bucket exists, we can write to it to reuse one bucket for multiple clusters.
- S3 bucket code creates a client that is locked to the region chosen for the cluster, not all regions support S3 and the code should be smarter and here are some options.
  - Add a region to the s3 bucket configs and reconfigure the client is set, default to the AWS default of us-east-1 if empty string
  - S3 enabled regions is a finite list, we could take the cluster region and see if s3 enabled and default to us-east-1 if no match
  - Force all buckets to S3 default region us-east-1

## Proposal
- Create a boolean on `AWSCluster` called `AssociateOIDCProvider` to match the `AWSManagedCluster` API and have a default value of `false`.
- Migrate the status types for `OIDCProvider` out of the experimental EKS APIs and into the v1beta2 APIs. 
- Build an IAM cloud service and add a reconciler to work to persist all components required for IRSA; the logic is as follows.
  1. Create a self-signed issuer for the workload cluster namespace to be used to make the pod identity webhook serving cert.
  2. Generate the patch file and update kubeadm configs to write the patch to disk for the control plane nodes.
  3. Create the Identity Provider in IAM pointed to the S3 bucket.
  4. Pause the reconciler until the workload cluster is online, as we have created all the pieces we can without a working kube api, the `AWSMachine` controller has additional code to annotate the `AWSCluster` if a control plane node is up and if the management cluster has a kubeconfig which will unpause our reconciler.
  5. Copy the [JWKS](https://auth0.com/docs/secure/tokens/json-web-tokens/json-web-key-sets) and OpenID Configuration from the kubeapi to the S3 bucket.
  6. Create all kube components in the workload cluster to run the pod-identity-webhook
  7. Create the trust policy boilerplate configmap in the workload cluster 

Identical to the EKS implementation, a trust policy document boilerplate will reference the ARN for the Identity Provider created in step 3. This can be used to generate IAM roles, and the ARNs for those roles can be annotated on a service account. The pod-identity-webhook works by watching all service accounts and pods. When it finds a pod using a service account with the annotation, it will inject AWS STS Tokens via environment variables generated from the role ARN. 

### S3 Bucket
A previous implementation for ignition support added an S3 bucket to support the configuration needed for ignition boots. The original functionality used two sub-folders, `control-plane` and `node`. These remain the same in this proposal with an addition of a new folder which matches the CAPA cluster name and makes a directory structure like the following. 

```
unique-s3-bucket-name/
|-- cluster1
|   |-- .well-known
|   `-- openid
|       `-- v1
|-- cluster2
|   |-- .well-known
|   `-- openid
|       `-- v1
|-- control-plane
`-- node
```

**Note**: today the code does not support reusing an S3 bucket as it errors if the bucket exists but support can be added to catch the exist error and attempt to write to the bucket to confirm access and reuse it for another cluster.

### Sample YAML
To add IRSA Support to an self-managed cluster your AWSCluster YAML will look something like the following.

```
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
kind: AWSCluster
metadata:
  name: capi-quickstart
  namespace: default
spec:
  region: us-west-2
  sshKeyName: luther
  associateOIDCProvider: true
  s3Bucket:
    name: capi-quickstart-1234 # regionally unique, be careful of name clashes with other AWS users
    nodesIAMInstanceProfiles:
      - nodes.cluster-api-provider-aws.sigs.k8s.io
    controlPlaneIAMInstanceProfile: control-plane.cluster-api-provider-aws.sigs.k8s.io
```

### User Stories

Story 1:
As an EKS cluster user who uses IRSA I want to...
- Migrate to self-managed clusters and maintain the same AWS API access

Story 2:
As a self-managed cluster user I want to...
- Give pods granular access to the AWS API based on IAM Roles

### Security Model

Access to the necessary CRDs is already declared for the controllers, and we are not adding any new kinds, so there is no change.

Since the jwks and openid config need public access the S3 Bucket config will need to be modified to allow both private and public access to objects. This is done by setting `PublicAccessBlockConfiguration` to false setting bucket ownership to `BucketOwnerPreferred`

Additional Permissions granted to the IAM Policies as follows

**Controllers Policy**
- iam:CreateOpenIDConnectProvider
- iam:DeleteOpenIDConnectProvider
- iam:ListOpenIDConnectProviders
- iam:GetOpenIDConnectProvider
- iam:TagOpenIDConnectProvider
- s3:PutBucketOwnershipControls
- s3:PutObjectAcl
- s3:PutBucketPublicAccessBlock

### Risks and Mitigations


## Alternatives

The process to install everything to use IRSA is documented and could be done by hand if necessary, but CAPA has complete control over the pieces needed and auto-mating this through a reconciler would make the feature on par with the existing functionality for Managed Clusters.

#### Benefits

This approach makes IRSA in self-managed clusters relatively trivial. The kube-api-server patch is tricky to manage by hand, and CAPA already has access to all the AWS Infrastructure it needs to auto-manage this problem.

#### Downsides

- Might be too much for CAPA to manage and not worth the complexity.

#### Decision

## Upgrade Strategy
Moving the OIDCProvider type from the experimental EKS API to the v1beta2 API for both cluster types will have converters for upgrading and downgrading. Through testing we can confirm but IRSA should be able to be added to a cluster after the fact, CAPA will need to patch kube-apiserver and create new control planes and the upgrade process should make this process seamless.

## Additional Details

### Test Plan
* Test creating a cluster, confirm all pieces work and have a simple AWS CLI example with a service account attached to a pod and exec commands successfully gaining auth through STS tokens attached via environment variables.
* Test deleting a cluster and confirm all AWS components are removed (s3 bucket contents, management cluster configmaps, etc.)
* Test upgrading a cluster with no IRSA to add the feature and confirm all components deployed successfully and test the AWS CLI example.
 
### Graduation Criteria

## Implementation History

- [x] 2023-03-22: Open proposal (PR)
- [x] 2023-02-22: WIP Implementation (PR)[https://github.com/kubernetes-sigs/cluster-api-provider-aws/pull/4094]

<!-- Links -->
[community meeting]: https://docs.google.com/document/d/1iW-kqcX-IhzVGFrRKTSPGBPOc-0aUvygOVoJ5ETfEZU/edit#
[discussion]: https://github.com/kubernetes-sigs/cluster-api-provider-aws/discussions/4153
