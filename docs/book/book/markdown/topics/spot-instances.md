# Spot Instances

[AWS Spot Instances](https://aws.amazon.com/ec2/spot/?cards.sort-by=item.additionalFields.startDateTime&cards.sort-order=asc&trk=a9b30b20-d23f-4d61-9452-c51a7e407fcd&sc_channel=ps&sc_campaign=acquisition&sc_medium=ACQ-P|PS-GO|Brand|Desktop|SU|Compute|EC2%20Spot|IN|EN|Text&s_kwcid=AL!4422!3!517651795636!e!!g!!amazon%20ec2%20spot&ef_id=Cj0KCQiA95aRBhCsARIsAC2xvfxB17BKyQFcn9UUKZ1GT2sfvxKyhboEKa87gl8wBO37fSrNXmx52cIaAtqwEALw_wcB:G:s&s_kwcid=AL!4422!3!517651795636!e!!g!!amazon%20ec2%20spot) allows user to reduce the costs of their compute resources by utilising AWS spare capacity for a lower price.

Because Spot Instances are tightly integrated with AWS services such as Auto Scaling, ECS and CloudFormation, users can choose how to launch and maintain their applications running on Spot Instances.

Although, with this lower cost, comes the risk of preemption. When capacity within a particular Availability Zone is increased, AWS may need to reclaim Spot instances to satisfy the demand on their data centres.

## When to use spot instances? 

Spot instances are ideal for workloads that can be interrupted. For example, short jobs or stateless services that can be rescheduled quickly, without data loss, and resume operation with limited degradation to a service.

## Using Spot Instances with AWSMachine

To enable AWS Machine to be backed by a Spot Instance, users need to add `spotMarketOptions` to AWSMachineTemplate:
```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: AWSMachineTemplate
metadata:
  name: ${CLUSTER_NAME}-md-0
spec:
  template:
    spec:
      iamInstanceProfile: nodes.cluster-api-provider-aws.sigs.k8s.io
      instanceType: ${AWS_NODE_MACHINE_TYPE}
      spotMarketOptions:
        maxPrice: ""
      sshKeyName: ${AWS_SSH_KEY_NAME}
```
Users may also add a `maxPrice` to the options to limit the maximum spend for the instance. It is however, recommended not to set a maxPrice as AWS will cap your spending at the on-demand price if this field is left empty, and you will experience fewer interruptions.
```yaml
spec:
  template:
    spotMarketOptions:
      maxPrice: 0.02 # Price in USD per hour (up to 5 decimal places)
```

## Using Spot Instances with AWSManagedMachinePool
To use spot instance in EKS managed node groups for a EKS cluster, set `capacityType` to `spot` in `AWSManagedMachinePool`.
```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: AWSManagedMachinePool
metadata:
  name: ${CLUSTER_NAME}-pool-0
spec:
  capacityType: spot
  ...
```

See [AWS doc](https://docs.aws.amazon.com/eks/latest/userguide/managed-node-groups.html) for more details.

## Using Spot Instances with AWSMachinePool
To enable AWSMachinePool to be backed by a Spot Instance, users need to add `spotMarketOptions` to AWSLaunchTemplate:
```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: AWSMachinePool
metadata:
  name: ${CLUSTER_NAME}-mp-0
spec:
  minSize: 1
  maxSize: 4
  awsLaunchTemplate:
    instanceType: "${AWS_CONTROL_PLANE_MACHINE_TYPE}"
    iamInstanceProfile: "nodes.cluster-api-provider-aws.sigs.k8s.io"
    sshKeyName: "${AWS_SSH_KEY_NAME}"
    spotMarketOptions:
       maxPrice: ""
```

> **IMPORTANT WARNING**: The experimental feature `AWSMachinePool` supports using spot instances, but the graceful shutdown of machines in `AWSMachinePool` is not supported and has to be handled externally by users.
