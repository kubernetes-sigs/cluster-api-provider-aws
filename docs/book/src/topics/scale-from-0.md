# Scaling from 0

With the changes introduce into `cluster-api` described in [this](https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20210310-opt-in-autoscaling-from-zero.md#upgrade-strategy) proposal, a user can now opt in to scaling nodes from 0.

This entails a number of things which I will describe in detail.

The following actions need to be taken to enabled cluster autoscaling:

## Set Capacity field

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

## Add two necessary annotations to MachineDeployment

There are two annotations which need to be applied to the MachineDeployment like this:

```yaml
apiVersion: cluster.x-k8s.io/v1beta1
kind: MachineDeployment
metadata:
  name: "managed-cluster-md-0"
  annotations:
    cluster.x-k8s.io/cluster-api-autoscaler-node-group-max-size: "5"
    cluster.x-k8s.io/cluster-api-autoscaler-node-group-min-size: "0"
```

These are necessary for the autoscaler to be able to pick up the deployment and scale it. Read more about these [here](https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/cloudprovider/clusterapi/README.md#enabling-autoscaling).

## Install and start cluster-autoscaler

Now comes the tricky part. In order for this to work, you need the cluster-autoscaler binary located [here](https://github.com/kubernetes/autoscaler).
You have to options. Use Helm to install autoscaler, or use the command line ( which is faster in if you are testing ).

In either cases, you need the following options:
- namespace
- cloud-provider
- scale-down-delay-after-add
- scale-down-delay-after-delete
- scale-down-delay-after-failure
- scale-down-unneeded-time
- expander
- kubeconfig
- cloud-config

These last two values are crucial for the autoscaler to work. `cloud-config` is the kubeconfig of the management cluster.
If you are using a service account to access it, you also have an option to define that. Read more about it on the
autoscaler's repository. The second one is the workload cluster. It needs both because the MachineDeployment is in the
control-plane while the actual node and pods are in the workload cluster.

Therefor, you have to install cluster-autoscaler into the _control-plane_ cluster.

I have a handy script to launch autoscaler which looks like this:

```bash
#!/bin/sh
# usage: start-autoscaler management.kubeconfig workload.kubeconfig
cluster-autoscaler \
	--cloud-provider=clusterapi \
	--v=4 \
	--namespace=default \
	--max-nodes-total=30 \
	--scale-down-delay-after-add=10s \
	--scale-down-delay-after-delete=10s \
	--scale-down-delay-after-failure=10s \
	--scale-down-unneeded-time=23s \
	--max-node-provision-time=2m \
	--balance-similar-node-groups \
	--expander=random \
	--kubeconfig=$2 \
	--cloud-config=$1
```

Courtesy of [@elmiko](https://github.com/elmiko).

The Helm equivalent is a bit more complex and either needs to mount in the kubeconfig from somewhere or be pointed to it.

## Permissions

This depends on your scenario. Read about it more [here](https://github.com/kubernetes/autoscaler/tree/master/cluster-autoscaler).
Since this is Cluster API Provider AWS, you would need to look for the AWS provider settings [here](https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/cloudprovider/aws/README.md).

Further, the service account associated with cluster-autoscaler requires permissions to access `get` and `list` the
Cluster API machine template infrastructure objects.

## Putting it together

The whole yaml looks like this:

```yaml
---
apiVersion: cluster.x-k8s.io/v1beta1
kind: Cluster
metadata:
  name: "managed-cluster"
spec:
  infrastructureRef:
    kind: AWSManagedControlPlane
    apiVersion: controlplane.cluster.x-k8s.io/v1beta1
    name: "managed-cluster-control-plane"
  controlPlaneRef:
    kind: AWSManagedControlPlane
    apiVersion: controlplane.cluster.x-k8s.io/v1beta1
    name: "managed-cluster-control-plane"
---
kind: AWSManagedControlPlane
apiVersion: controlplane.cluster.x-k8s.io/v1beta1
metadata:
  name: "managed-cluster-control-plane"
spec:
  region: "eu-central-1"
  version: "v1.22.0"
---
apiVersion: cluster.x-k8s.io/v1beta1
kind: MachineDeployment
metadata:
  name: "managed-cluster-md-0"
  annotations:
    cluster.x-k8s.io/cluster-api-autoscaler-node-group-max-size: "5"
    cluster.x-k8s.io/cluster-api-autoscaler-node-group-min-size: "0"
spec:
  clusterName: "managed-cluster"
  replicas: 0 # _NOTE_ that we set the initial replicas size to *ZERO*.
  selector:
    matchLabels:
  template:
    spec:
      clusterName: "managed-cluster"
      version: "v1.22.0"
      bootstrap:
        configRef:
          name: "managed-cluster-md-0"
          apiVersion: bootstrap.cluster.x-k8s.io/v1beta1
          kind: EKSConfigTemplate
      infrastructureRef:
        name: "managed-cluster-md-0"
        apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
        kind: AWSMachineTemplate
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: AWSMachineTemplate
metadata:
  name: "managed-cluster-md-0"
spec:
  template:
    spec:
      instanceType: "t3.small"
      iamInstanceProfile: "nodes.cluster-api-provider-aws.sigs.k8s.io"
status:
  capacity:
    memory: "500m"
    cpu: "1"
---
apiVersion: bootstrap.cluster.x-k8s.io/v1beta1
kind: EKSConfigTemplate
metadata:
  name: "managed-cluster-md-0"
spec:
  template: {}
```

## When will it not scale?

There is a document describing under what circumstances it won't be able to scale located [here](https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/FAQ.md#what-types-of-pods-can-prevent-ca-from-removing-a-node). Read this carefully.

It has some ramifications when scaling back down to 0. Which will only work if all pods are removed from the node and
the node cannot schedule even the aws-node and kube-proxy pods. There is this tiny manual step of cordoning off the last
node in order to scale back down to 0.

## Conclusion

Once the cluster-autoscaler is running, you will start seeing nodes pop-in as soon as there is some load on the cluster.
To test it, simply create and inflate a deployment like this:

```bash
kubectl create deployment inflate --image=public.ecr.aws/eks-distro/kubernetes/pause:3.2 --kubeconfig workload.kubeconfig
kubectl scale deployment inflate --replicas=50 --kubeconfig workload.kubeconfig
```
