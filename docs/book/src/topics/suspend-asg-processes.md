# Suspend ASG Processes

- **Feature status:** Experimental
- **Feature gate:** MachinePool=true

MachinePool allows users to manage many machines as a single entity. Infrastructure providers implement a separate CRD that handles infrastructure side of the feature.

## Suspend Processes

It's possible to suspend certain processes for ASG. The list of processes can be found [here](https://docs.aws.amazon.com/autoscaling/ec2/APIReference/API_SuspendProcesses.html).

To utilize this feature, simply denote the list of processes that are desired to be suspended.

```yaml
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
kind: AWSMachinePool
metadata:
  name: capa-mp-0
spec:
  minSize: 1
  maxSize: 10
  availabilityZones:
    - "${AWS_AVAILABILITY_ZONE}"
  awsLaunchTemplate:
    instanceType: "${AWS_CONTROL_PLANE_MACHINE_TYPE}"
    sshKeyName: "${AWS_SSH_KEY_NAME}"  
  suspendProcesses:
    processes:
        launch: true
        alarmNotification: true
        azRebalance: true
---
```

## Resume Processes

If a service is desired to be resumed, simply remove it from the list of suspended processes. The reconciler will then
resume any process that is not part of the desired suspended processes list.

```yaml
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
kind: AWSMachinePool
metadata:
  name: capa-mp-0
spec:
  minSize: 1
  maxSize: 10
  availabilityZones:
    - "${AWS_AVAILABILITY_ZONE}"
  awsLaunchTemplate:
    instanceType: "${AWS_CONTROL_PLANE_MACHINE_TYPE}"
    sshKeyName: "${AWS_SSH_KEY_NAME}"  
  suspendProcesses:
    processes:
      launch: true
---
```

_Note_ that now `AlarmNotification` and `AZRebalance` will be resumed, but the reconciler will not try to suspend
`Launch` again. So it doesn't incur additional expensive, redundant API calls.

## Optional `All`

An option is also provided to suspend all processes without having to set each of them to `true`. Simply use `all` like
this:

```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
kind: AWSMachinePool
metadata:
  name: capa-mp-0
spec:
  minSize: 1
  maxSize: 10
  availabilityZones:
    - "${AWS_AVAILABILITY_ZONE}"
  awsLaunchTemplate:
    instanceType: "${AWS_CONTROL_PLANE_MACHINE_TYPE}"
    sshKeyName: "${AWS_SSH_KEY_NAME}"  
  suspendProcesses:
    all: true
```

To exclude individual processes from `all` simply add them with value `false`:

```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
kind: AWSMachinePool
metadata:
  name: capa-mp-0
spec:
  minSize: 1
  maxSize: 10
  availabilityZones:
    - "${AWS_AVAILABILITY_ZONE}"
  awsLaunchTemplate:
    instanceType: "${AWS_CONTROL_PLANE_MACHINE_TYPE}"
    sshKeyName: "${AWS_SSH_KEY_NAME}"  
  suspendProcesses:
    all: true
    processes:
      launch: false
```
