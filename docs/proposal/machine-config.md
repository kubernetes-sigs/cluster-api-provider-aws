# AWS Provider Machine Configuration and Status

## Overview
This proposal aims to provide an initial API for AWS machine configuration and status for the AWS machine actuator.
Cluster configuration for the AWS provider is out of scope for this proposal.

## AWSMachineProviderConfig
Defines the spec portion of a Machine specific to AWS. It is meant to be embedded in the
Machine.Spec.ProviderConfig field of a cluster API Machine.

Following are design decisions around this type:

- No AWS SDK types are included. For the following reasons:
  - Clients would also need to have the AWS SDK vendored.
  - The type would be tied to a specific AWS client version.

- The machine ProviderConfig includes all the configuration needed to create a Machine.
  Until we resolve issues related to linking a cluster and machine, it is not clear
  whether the machine controller should be required to pass a Cluster resource to the actuator:
  - https://github.com/kubernetes-sigs/cluster-api/issues/41
  - https://github.com/kubernetes-sigs/cluster-api/issues/145
  - https://github.com/kubernetes-sigs/cluster-api/issues/177

- Secrets are referenced for UserData and Account credentials, but are optional.
  If the secrets are not present, then no user data is set on the instance, and if no
  account credentials are present, it is assumed the IAMRole of the machine where the
  actuator is running has enough permission to create a new instance.

- The list of ELBs that should be associated with a machine are set explicitly in the config
  and are attached when the instance is created.

## AWSMachineProviderStatus
Defines the status portion of a Machine specific to AWS. It is meant to be embedded in the
Machine.Status.ProviderStatus field of a cluster API Machine.

Design decisions for this type:

- Conditions are used to report on problems creating an instance, or attaching ELBs.
  Reference: https://github.com/kubernetes/kubernetes/issues/7856
  Although the issue seems to imply that we want to get rid of conditions, Brian Grant's
  [conclusion](https://github.com/kubernetes/kubernetes/issues/7856#issuecomment-335687733)
  states that conditions should be kept to report extended status. Knowing why we failed
  to create an instance or whether attaching an ELB did not succeed is still very valuable
  information.

## Open Issues

### Load Balancer Management
- Currently we include a set of IDs for ELBs to attach to the machine on creation. However, what happens if the ELB
  doesn't exist, or if after the machine is created the ELB fails to attach? Should this be handled by a separate
  controller?

### BlockDevice Mapping
- The RunInstances EC2 call allows specifying a BlockDeviceMapping to allow changing volume size or type from what
  is specified in the AMI. Is this something that we need? or is it enough to say that if you want something different
  to simply create a different AMI?

### Extensibility
- If the current types don't meet all the needs for different providers, do we allow some kind of extension mechanism
  to the types? How do we extend the actuator itself? Or should we simply allow the base implementation to be vendored
  into an implementation that can augment it. What hooks would we provide to allow this?
