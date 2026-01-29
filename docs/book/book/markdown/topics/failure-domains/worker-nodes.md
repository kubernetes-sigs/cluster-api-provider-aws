# Failure domains in worker nodes

To ensure that the worker machines are spread across failure domains, we need to create N `MachineDeployment` for your N failure domains, scaling them independently. Resiliency to failures comes from having multiple `MachineDeployment`.
For example:

```yaml
apiVersion: cluster.x-k8s.io/v1beta1
kind: MachineDeployment
metadata:
  name: ${CLUSTER_NAME}-md-0
  namespace: default
spec:
  clusterName: ${CLUSTER_NAME}
  replicas: ${WORKER_MACHINE_COUNT}
  selector:
    matchLabels: null
  template:
    spec:
      bootstrap:
        configRef:
          apiVersion: bootstrap.cluster.x-k8s.io/v1beta1
          kind: KubeadmConfigTemplate
          name: ${CLUSTER_NAME}-md-0
      clusterName: ${CLUSTER_NAME}
      infrastructureRef:
        apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
        kind: AWSMachineTemplate
        name: ${CLUSTER_NAME}-md-0
      version: ${KUBERNETES_VERSION}
      failureDomain: "1"
---
apiVersion: cluster.x-k8s.io/v1beta1
kind: MachineDeployment
metadata:
  name: ${CLUSTER_NAME}-md-1
  namespace: default
spec:
  clusterName: ${CLUSTER_NAME}
  replicas: ${WORKER_MACHINE_COUNT}
  selector:
    matchLabels: null
  template:
    spec:
      bootstrap:
        configRef:
          apiVersion: bootstrap.cluster.x-k8s.io/v1beta1
          kind: KubeadmConfigTemplate
          name: ${CLUSTER_NAME}-md-1
      clusterName: ${CLUSTER_NAME}
      infrastructureRef:
        apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
        kind: AWSMachineTemplate
        name: ${CLUSTER_NAME}-md-1
      version: ${KUBERNETES_VERSION}
      failureDomain: "2"
---
apiVersion: cluster.x-k8s.io/v1beta1
kind: MachineDeployment
metadata:
  name: ${CLUSTER_NAME}-md-2
  namespace: default
spec:
  clusterName: ${CLUSTER_NAME}
  replicas: ${WORKER_MACHINE_COUNT}
  selector:
    matchLabels: null
  template:
    spec:
      bootstrap:
        configRef:
          apiVersion: bootstrap.cluster.x-k8s.io/v1beta1
          kind: KubeadmConfigTemplate
          name: ${CLUSTER_NAME}-md-2
      clusterName: ${CLUSTER_NAME}
      infrastructureRef:
        apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
        kind: AWSMachineTemplate
        name: ${CLUSTER_NAME}-md-2
      version: ${KUBERNETES_VERSION}
      failureDomain: "3"
```

>**IMPORTANT WARNING:** All the replicas within a `MachineDeployment` will reside in the same Availability Zone.

### Using AWSMachinePool

You can use an `AWSMachinePool` object which automatically distributes worker machines across the configured availability zones.
Set the **FailureDomains** field to the list of availability zones that you want to use. Be aware that not all regions have the same availability zones.

```yaml
apiVersion: cluster.x-k8s.io/v1beta1
kind: MachinePool
metadata:
  labels:
    cluster.x-k8s.io/cluster-name: my-cluster
  name: ${CLUSTER_NAME}-mp-0
  namespace: default
spec:
  clusterName: my-cluster
  failureDomains:
    - "1"
    - "3"
  replicas: 3
  template:
    spec:
      clusterName: my-cluster
      bootstrap:
        configRef:
          apiVersion: bootstrap.cluster.x-k8s.io/v1beta1
          kind: KubeadmConfigTemplate
          name: ${CLUSTER_NAME}-mp-0
      infrastructureRef:
        apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
        kind: AWSMachinePool
        name: ${CLUSTER_NAME}-mp-0
      version: ${KUBERNETES_VERSION}
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: AWSMachinePool
metadata:
  labels:
    cluster.x-k8s.io/cluster-name: my-cluster
  name: ${CLUSTER_NAME}-mp-0
  namespace: default
spec:
  minSize: 1
  maxSize: 4
  awsLaunchTemplate:
    instanceType: ${AWS_NODE_MACHINE_TYPE}
    iamInstanceProfile: "nodes.cluster-api-provider-aws.sigs.k8s.io"
    sshKeyName: ${AWS_SSH_KEY_NAME}
```