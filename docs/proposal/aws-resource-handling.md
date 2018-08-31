# AWS Resource Handling

## Problem

Since the AWS APIs do not provide consistent support for tagging resources on creation or idempotent operations, we need to ensure that we can provide as much of an external guarantee to users as we can. Otherwise we risk creating and orphaning resources on users leading to extraneous costs and the potential to exhaust resource quotas prematurely.

## Existing solutions and drawbacks for our use

### kops

kops relies heavily on tagging of resources. This provides the benefit of recreating the state of the cluster resources by querying with filters. However, in many cases tagging requires a second API call and there are some resources that are not able to be tagged. If we succeed in creating a resource but fail to tag the resource, then we risk the chance of orphaning that resource for the user.

### Kubicorn

In contrast to the kops approach, Kubicorn mainly relies on recording the resource IDs as part of state. However, since we rely on an external API server for recording the state there is still a possibility of creating the resource and failing to record the resource ID, which still exposes the possible risk of orphaning that resource for the user.

## Summary of edge cases for creating an individual resource

- resource create succeeds, but subsequent tagging fails
- resource creates succeeds, but update of cluster/machine object fails
- attempting to delete resource fails after an attempt to rollback due to a failure to record the ID of the created resource to the cluster/machine object for resources that do not support tagging on create.
- the controller/actuator dies after creating a resource but before tagging and or recording the resource

## Proposed workflow

### Create resource with tag on create support

- Query resource by tags to determine if resource already exists
- Create the resource if it doesn't already exist
- Attempt to record the ID to the cluster/machine object
- Enque update for available/ready state if not already available/ready

![Create Resource](create-resource-with-tags.png)

### Create resource with separate tagging required

- Create resource
- Attempt to record ID to cluster/machine object
  - Attempt rollback of created resource, logging if unsuccessful
- Attempt to tag resource
- Enque update for available/ready state if not already available/ready

![Create Resource Separate Tagging](create-resource-separate-tags.png)

## Using client tokens

Where possible use [client tokens](https://docs.aws.amazon.com/AWSEC2/latest/APIReference/Run_Instance_Idempotency.html) in the create request so that subsequent requests will return the same response.

## Tagging of resources

Resources that are managed by the controllers/actuators should be tagged with: `kubernetes.io/cluster/<name or id>=owned`

TODO: Define additional tags that can be used to provide additional metadata about the resource configuration/usage by the actuator.

## Handling of errors

Each resource has specific error codes that it will return and these can be used to differentiate fatal errors from retryable errors. These errors are well documented in some cases (elbv2 api), and poorly in others (ec2 api). We should provide a best effort to [properly handle these errors](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/handling-errors.html) in the correct manner.