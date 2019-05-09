# AWS Resource Handling

## Problem

Since the AWS APIs do not provide consistent support for tagging resources on creation or idempotent operations, we need to ensure that we can provide as much of an external guarantee to users as we can. Otherwise we risk creating and orphaning resources on users leading to extraneous costs and the potential to exhaust resource quotas prematurely.

## Existing solutions and drawbacks for our use

### kops

kops relies heavily on tagging of resources. This provides the benefit of recreating the state of the cluster resources by querying with filters. However, in many cases tagging requires a second API call and there are some resources that are not able to be tagged. If we succeed in creating a resource but fail to tag the resource, then we risk the chance of orphaning that resource for the user.

### Kubicorn

In contrast to the kops approach, Kubicorn mainly relies on recording the resource IDs as part of state. However, since we rely on an external API server for recording the state there is still a possibility of creating the resource and failing to record the resource ID, which still exposes the possible risk of orphaning that resource for the user.

## Summary of edge cases for creating an individual resource

1. resource create succeeds, but subsequent tagging fails
2. resource creates succeeds, but update of cluster/machine object fails
3. attempting to delete resource fails after an attempt to rollback due to a failure to record the ID of the created resource to the cluster/machine object for resources that do not support tagging on create.
4. the controller/actuator dies after creating a resource but before tagging and or recording the resource

## Misc TODOs

- Solicit feedback on whether aws-sdk-go and client-go retry defaults are sufficient:
  - aws-sdk-go
    - https://docs.aws.amazon.com/general/latest/gr/api-retries.html
    - https://docs.aws.amazon.com/sdk-for-go/api/aws/client/#DefaultRetryer
  - client-go
    - https://github.com/kubernetes/client-go/blob/master/rest/request.go#L659
- Identify which resources fall into which class of workflow based on tagging support, client token support, etc.
- Better define mutatable/non-mutatable attributes for objects during update

## Using client tokens

Where possible use [client tokens](https://docs.aws.amazon.com/AWSEC2/latest/APIReference/Run_Instance_Idempotency.html) in the create request so that subsequent requests will return the same response.

## Tagging of resources

Resources handled by these components fall into one of three categories:

1. Fully-managed resources whose lifecycle is tied to the cluster. These resources should be tagged with `sigs.k8s.io/cluster-api-provider-aws/cluster/<name or id>=owned`, and the actuator is expected to keep these resources as closely in sync with the spec as possible.
2. Resources whose management is shared with the in-cluster aws cloud provider, such as a security group for load balancer ingress rules. These resources should be tagged with `sigs.k8s.io/cluster-api-provider-aws/cluster/<name or id>=owned` and `kubernetes.io/cluster/<name or id>=owned`, with the latter being the tag defined by the cloud provider. These resources are create/delete only: that is to say their ongoing management is "handed off" to the cloud provider.
3. Unmanaged resources that are provided by config (such as a common VPC). The provider will avoid changing these resources as much as is possible.

TODO: Define additional tags that can be used to provide additional metadata about the resource configuration/usage by the actuator. This is would allow us to rebuild status without relying on polluting the object config.

## Handling of AWS api errors

Each resource has specific error codes that it will return and these can be used to differentiate fatal errors from retryable errors. These errors are well documented in some cases (elbv2 api), and poorly in others (ec2 api). We should provide a best effort to [properly handle these errors](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/handling-errors.html) in the correct manner.

## Proposed workflows

### Resources that support tag on create support (with or without client tokens)

#### Create

- Query resource by tags to determine if resource already exists
- Create the resource if it doesn't already exist
- Update the cluster/machine object config and status
  - If update fails return a retryable error to requeue the create
- Enqueue cluster/machine update if not already available/ready

TODO: flowchart

#### Update

- Query resource by ID
- Update object status
- Enqueue cluster/machine update if not already available/ready

#### Edge case coverage

1. Yes - tagging is handled on creation
2. Yes - since resources are tagged on creation, returning an error and requeueing the create will find the tagged resource and attempt to retry the object update.
3. Yes - there is no delete attempt since we can re-query the resource by tags.
4. Yes - the next attempt to create the resource will find the already created resource by tags.

### Resources that support client tokens but require separate tagging

#### Create

- Create the resource using object uid as the client token
- Update the cluster/machine object config and status
  - If update fails return a retryable error to requeue the create
- Tag AWS resource
- Enqueue cluster/machine update if not already available/ready

TODO: flowchart

#### Update

- Query resource by ID
- tag resource if missing tags
- Update object status
- Enqueue cluster/machine update if not already available/ready

#### Edge case coverage

1. Yes - If the update was successful, tagging will be reconciled the next time the object is updated on reconiliation. If the update was not successful, the next call using the same client token will return the same object as previously created.
2. Yes - since we are using a client token, subsequent requests will return the same result.
3. Yes - there is no delete attempt since we can repeat the request to create the resource safely.
4. Yes - the next attempt to create the resource will return the already created resource.

### Resources that require separate tagging without client token support

#### Create - option 1

- Create resource
- Update cluster/machine object config and status
  - If update fails attempt delete of created resource
    - If delete fails log delete failure and return non-retryable error
- Tag AWS resource
- Enqueue cluster/machine update if not already available/ready or tagging fails

TODO: flowchart

#### Update - option 1

- Query resource by ID
- tag resource if missing tags
- Update object status
- Enqueue cluster/machine update if not already available/ready

#### Edge case coverage - option 1

1. Yes - Since the resource ID is already recorded, the update process will reconcile missing tags
2. Yes, with caveat - If the object update fails, we attempt to rollback the creation but edge case 3 comes into play
3. Minor mitigation - If we fail to delete the resource, we will still orphan the resource, but output a log message for querying/followup and return a non-retryable error
4. Minor mitigation - If the process dies before recording the ID the resource is orphaned. If the process dies after recording the ID, but before tagging it is reconciled through update.

#### Create - option 2

- Query resource by tags to determine if resource already exists
- Create the resource if it doesn't already exist
- Update cluster/machine object config and status
  - Note failure but do not return error yet
- Tag AWS resource if needed
  - If both update and tagging fails, delete resource
    - If delete fails log failure prominently, return non-retryable error
  - If only tagging fails, return retryable error
- If update failed, return retryable error
- Enqueue cluster/machine update if not already available/ready or tagging fails

TODO: flowchart

#### Update - option 2

- Query resource by ID
- Update object status
- Enqueue cluster/machine update if not already available/ready

#### Edge case coverage - option 2

1. Yes - If only the tagging fails, then we will reconcile tags on update. If the update also fails, then we attempt to delete the resource.
2. Yes, with caveat - If only the object update fails, then we throw a retryable error that will requeue the create operation and attempt to update the object after discovering the existing resource. If the tagging fails as well, then we attempt to delete the resource and edge case 3 will still apply.
3. Minor mitigation - If we fail to delete the resource, we will still orphan the resource, but output a log message for querying/followup and return a non-retryable error
4. Minor mitigation - If the process dies before recording the ID the resource is orphaned. If the process dies after recording the ID, but before tagging it is reconciled through update.

### Resources without tag support with client token support

#### Create

- Create the resource using the object uid as the client token
- Update cluster/machine object config and status
- Enqueue cluster/machine update if not already available/ready or tagging fails

TODO: flowchart

#### Update

- Query resource by ID
- Update object status
- Enqueue cluster/machine update if not already available/ready

#### Edge case coverage

1. Yes - There is no tagging
2. Yes - If the update fails, subsequent calls will return the same resource.
3. Yes - No delete is used
4. Yes - If the process dies before recording the ID subsequent calls to create the resource return the same resource.

### Resources without tag support without client token support

#### Create

- Create the resource
- Update cluster/machine object config and status
  - If update fails, delete resource
    - If delete fails log failure prominently, return non-retryable error
- Enqueue cluster/machine update if not already available/ready or tagging fails

TODO: flowchart

#### Update

- Query resource by ID
- Update object status
- Enqueue cluster/machine update if not already available/ready

#### Edge case coverage

1. Yes - There is no tagging
2. Yes, with caveat - If the update fails, then we attempt to delete the resource and edge case 3 will still apply.
3. Minor mitigation - If we fail to delete the resource, we will still orphan the resource, but output a log message for querying/followup and return a non-retryable error
4. No - If the process dies before recording the ID the resource is orphaned.
