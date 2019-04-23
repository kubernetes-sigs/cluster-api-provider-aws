# Creating clusters with a multi-node control plane  <!-- omit in toc -->

`clusterctl` only supports creating clusters with a single instance control
plane. However, cluster-api-provider-aws has the capability to create clusters
with multiple control plane nodes. This document outlines the steps for standing
up a cluster with a multi-node control plane without using the standard
`clusterctl create cluster` command.

## High availability vs multi-node control plane <!-- omit in toc -->

We don't call a multi-node control plane highly available in this document
because all the control plane nodes live in the same availability zone. This
means that an availability zone is a single point of failure for this cluster.

cluster-api-provider-aws does actually account for this and allows you to create
a highly available control plane with a little bit of extra work. Please see
[setting up a highly available cluster](#setting-up-a-highly-available-cluster)
for more details.

## Contents <!-- omit in toc -->

- [Prerequisites](#prerequisites)
  - [Creating a cluster with a single instance control plane](#creating-a-cluster-with-a-single-instance-control-plane)
- [Growing the control plane into a multi-node control plane](#growing-the-control-plane-into-a-multi-node-control-plane)
- [Troubleshooting](#troubleshooting)
  - [Instances for the new machines not getting created](#instances-for-the-new-machines-not-getting-created)
  - [Instances were created and initialized, but nodes don't show up on running 'kubectl get nodes'](#instances-were-created-and-initialized-but-nodes-dont-show-up-on-running-kubectl-get-nodes)
- [Setting up a highly available cluster](#setting-up-a-highly-available-cluster)

## Prerequisites

### Creating a cluster with a single instance control plane

Set up a single instance control plane cluster by following the
[getting started guide](docs/getting-started.md).

## Growing the control plane into a multi-node control plane

At this point you should already have a cluster with a single instance control
plane and is responding to `kubectl`.

The next step is to create YAML representing a new Machine object for each
control plane you want to add to your existing cluster.

1. Create the Machine YAML

  Copy the following YAML into a file named `control-plane-machine.yaml`.

```yaml
apiVersion: "cluster.k8s.io/v1alpha1"
kind: Machine
metadata:
  name: <CONTROLPLANE_MACHINE_NAME>
  namespace: default # Edit this if necessary
  labels:
    cluster.k8s.io/cluster-name: <CLUSTER_NAME>
    set: controlplane
spec:
  versions:
    kubelet: v1.13.3
    controlPlane: v1.13.3
  providerSpec:
    value:
      apiVersion: awsprovider/v1alpha1
      kind: AWSMachineProviderSpec
      instanceType: "t2.medium"
      iamInstanceProfile: "control-plane.cluster-api-provider-aws.sigs.k8s.io"
      keyName: "cluster-api-provider-aws.sigs.k8s.io"
```

  *Pro tip* 💡: You may refer to the machine's YAML in the release to create the
  above YAML for yourself.

  Fill in the cluster name and update the namespace if necessary. Then copy
  `control-plane-machine.yaml` as many times as the number of *new* control
  planes you want to add to your cluster.

  For each copy, edit the file with a unique machine name and create the machine object:

```bash
kubectl apply -f control-plane-machine.yaml
```

  Note: If you did not specify namespace in the `control-plane-machine.yaml`
  file(s) and your cluster is not in the `default` namespace, you will have to
  use the `-n <NAMESPACE>` option to kubectl.

2. Optionally watch machines get created

  Watch the newly created machine objects being reconciled by following the logs
  from the `aws-provider-controller-manager-0` pod using this command:

```bash
kubectl -n aws-provider-system logs -f aws-provider-controller-manager-0
```

   Additionally you can view the new instances initializing in the AWS console.

3. Verify nodes are ready

  Once the machine objects have been reconciled and the instances in AWS have
  been created and initialized, `kubectl get nodes` will show the new
  machines with master roles.

Your cluster now has a multi-node control plane! 🎉

## Troubleshooting

### Instances for the new machines not getting created

Inspect the logs from the `aws-provider-controller-manager-0` pod to determine
why the control plane is unable to create the necessary resources. Follow the
logs from the `aws-provider-controller-manager-0` pod using the command below:

```bash
kubectl -n aws-provider-system logs -f aws-provider-controller-manager-0
```

### Instances were created and initialized, but nodes don't show up on running 'kubectl get nodes'

Inspect `/var/log/cloud-init-output.log` on the newly created instances for information about why
this instance failed to join the existing control plane. Instructions on how to
view the system log are available
[here](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-console.html).

You can ssh to any node in the private subnet that was created through the
bastion host using a command like:

`ssh -J ubuntu@<PUBLIC_BASTION_IP> ubuntu@<PRIVATE_INSTANCE_IP>`

## Setting up a highly available cluster

If you pre-provision subnets across multiple availability zones in the same region and define
the different subnet ids in the control plane YAML copies, then each node will
be created in the availability zone that the subnet is in.

For example:

```yaml
kind: Machine
metadata:
    name: aws-controlplane-0
    labels:
        cluster.k8s.io/cluster-name: <CLUSTER_NAME>
        set: controlplane
spec:
    versions:
        kubelet: v1.13.3
        controlPlane: v1.13.3
    providerSpec:
        value:
            apiVersion: awsprovider/v1alpha1
            kind: AWSMachineProviderSpec
            instanceType: "t2.medium"
            iamInstanceProfile: "control-plane.cluster-api-provider-aws.sigs.k8s.io"
            keyName: "laptop"
            subnet:
                id: <SUBNET_IN_AZ1_ID>
---
kind: Machine
metadata:
    name: aws-controlplane-1
    labels:
        cluster.k8s.io/cluster-name: <CLUSTER_NAME>
        set: controlplane
spec:
    versions:
        kubelet: v1.13.3
        controlPlane: v1.13.3
    providerSpec:
        value:
            apiVersion: awsprovider/v1alpha1
            kind: AWSMachineProviderSpec
            instanceType: "t2.medium"
            iamInstanceProfile: "control-plane.cluster-api-provider-aws.sigs.k8s.io"
            keyName: "laptop"
            subnet:
                id: <SUBNET_IN_AZ2_ID>
```
