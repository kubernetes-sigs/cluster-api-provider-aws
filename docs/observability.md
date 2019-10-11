# Observability

## Events

The following events are published by this provider:

### AWSMachines

* `FailedTerminate`: The provider failed to terminate an instance during machine
  deletion.
* `SuccessfulTerminate`: The provider successfully terminated the EC2 instance
* `InvalidUpdate`: An attempt to mutate the machine object was made. This includes
  changing the EC2 instance type, SSH or IAM profile, root device size and
  changing from public IP address to not. This will eventually be enforced
  by validating webhooks, so will not remain an event in the long term.
* `NoInstanceFound`: No instance was found matching the machine.
* `FailedAttachControlPlaneELB`: Couldn't attach the EC2 instance to the Elastic
  Load Balancer.
