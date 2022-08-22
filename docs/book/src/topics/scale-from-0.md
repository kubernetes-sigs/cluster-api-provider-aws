# Scaling from 0

With the changes introduce into `cluster-api` described in [this](https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20210310-opt-in-autoscaling-from-zero.md#upgrade-strategy) proposal, a user can now opt in to scaling nodes from 0.

To do that, simply define some values to the new field called `capacity` in the `AWSMachineTemplate` like this:

```yaml
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: AWSMachineTemplate
metadata:
  name: "${CLUSTER_NAME}-md-0"
spec:
  template:
    spec:
      instanceType: "${AWS_NODE_MACHINE_TYPE}"
      iamInstanceProfile: "nodes.cluster-api-provider-aws.sigs.k8s.io"
      sshKeyName: "${AWS_SSH_KEY_NAME}"
status:
  capacity:
    memory: "500m"
    cpu: "1"
    nvidia.com/gpu: "1"
```

To read more about what values are available, consult the proposal. These values can be overridden by selected annotations
on the MachineTemplate.
