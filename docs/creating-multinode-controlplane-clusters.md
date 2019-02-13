# Creating Clusters With Multi-Node Control Plane  <!-- omit in toc -->
Currently, clusterctl only supports creating clusters with a single instance control plane, but cluster-api-provider-aws has the capability to create clusters with multiple control plane nodes. This document outlines the steps for standing up a cluster with a multi-node control plane.

## Contents <!-- omit in toc -->
- [Prerequisites](#prerequisites)
  - [Creating cluster with single instance controlplane.](#creating-cluster-with-single-instance-controlplane)
- [Growing controlplane into multi-node controlplane.](#growing-controlplane-into-multi-node-controlplane)
- [Troubleshooting](#troubleshooting)
  - [Instances for the new machines not getting created](#instances-for-the-new-machines-not-getting-created)
  - [Instances were created and initialized, but nodes don't show up on running 'kubectl get no'](#instances-were-created-and-initialized-but-nodes-dont-show-up-on-running-kubectl-get-no)

## Prerequisites
### Creating cluster with single instance controlplane.
Following the [getting started guide](docs/getting-started.md), please set up a cluster with a single control plane instance.

## Growing controlplane into multi-node controlplane.
At this point, a cluster with a single instance control plane has been created and is responding to `kubectl`. Please refer to the 'Using the Cluster' section of the [getting started guide](docs/getting-started.md) for instructions. If you are unable to run `kubectl` against the cluster, please follow the 'Troubleshooting' section of the [getting started guide](docs/getting-started.md) to resolve the problem. 

Now, that the cluster with the single instance control plane has been created, more controlplane instances can be added by creating a machine object for each desired controlplane machine, one machine at a time. 
1. Copy the following yaml into a file, say `controlplane-machine.yaml`, and for each desired controlplane machine, update `controlplane-machine.yaml` with the machine and cluster name.
   
  ```bash
  apiVersion: "cluster.k8s.io/v1alpha1"
  kind: Machine
  metadata:
    name: <CONTROLPLANE_MACHINE_NAME>
    labels:
      cluster.k8s.io/cluster-name: <CLUSTER_NAME>
      set: controlplane
  spec:
    versions:
      kubelet: v1.13.2
      controlPlane: v1.13.2
    providerSpec:
      value:
        apiVersion: awsprovider/v1alpha1
        kind: AWSMachineProviderSpec
        instanceType: "t2.medium"
        iamInstanceProfile: "control-plane.cluster-api-provider-aws.sigs.k8s.io"
        keyName: "cluster-api-provider-aws.sigs.k8s.io"
  ```
  *Pro tip* 💡: You may refer to the machine's yaml in the release to create the above yaml for yourself.

2. For each desired controlplane machine, create the machine object, alongside other machines of the cluster, by applying the `controlplane-machine.yaml`, from above.
```bash
kubectl apply -f controlplane-machines.yaml
``` 
Note: You may have to use the `-n <NAMESPACE>` option to kubectl incase your cluster and machines were created in the non-default namespace.

3. This step is optional. Watch these newly created machine objects being reconciled by following the logs from the `aws-provider-controller-manager-0` pod using the below command
```bash
kubectl -n aws-provider-system logs -f aws-provider-controller-manager-0
```
Additionally, view the new instances initializing in the AWS console.

4. Once the machine objects have been reconciled and the instances in AWS have been created and initialized, `kubectl get nodes` should show the two new machines with master roles.

Just like that, you have now transformed your cluster from a single instance controlplane to one with a multi-node controlplane. 🎉

## Troubleshooting
### Instances for the new machines not getting created
Inspect the logs from the `aws-provider-controller-manager-0` pod to determine why the controlplane is unable to create the necessary resources. Follow the logs from the `aws-provider-controller-manager-0` pod using the command below:
```bash
kubectl -n aws-provider-system logs -f aws-provider-controller-manager-0
```
### Instances were created and initialized, but nodes don't show up on running 'kubectl get no'
Inspect the system log on the newly created instances for information about why this instance failed to join the existing controlplane. Instructions on how to view the system log are available [here](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-console.html)
