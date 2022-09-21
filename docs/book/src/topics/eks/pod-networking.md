# Pod Networking

When creating a EKS cluster the Amazon VPC CNI will be used by default for Pod Networking.

> When using the AWS Console to create an EKS cluster with a Kubernetes version of v1.18 or greater you are required to select a specific version of the VPC CNI to use.

## Using the VPC CNI Addon
You can use an explicit version of the Amazon VPC CNI by using the **vpc-cni** EKS addon. See the [addons](./addons.md) documentation for further details of how to use addons.

## Using an alternative CNI

There may be scenarios where you do not want to use the Amazon VPC CNI. EKS supports a number of alternative CNIs such as Calico, Cilium, and Weave Net (see [docs](https://docs.aws.amazon.com/eks/latest/userguide/alternate-cni-plugins.html) for full list).

There are a number of ways to install an alternative CNI into the cluster. One option is to use a [ClusterResourceSet](https://cluster-api.sigs.k8s.io/tasks/experimental-features/cluster-resource-set.html) to apply the required artifacts to a newly provisioned cluster.

When using an alternative CNI you will want to delete the Amazon VPC CNI, especially for a cluster using v1.17 or less. This can be done via the **disableVPCCNI** property of the **AWSManagedControlPlane**:

```yaml
kind: AWSManagedControlPlane
apiVersion: controlplane.cluster.x-k8s.io/v1beta1
metadata:
  name: "capi-managed-test-control-plane"
spec:
  region: "eu-west-2"
  sshKeyName: "capi-management"
  version: "v1.18.0"
  disableVPCCNI: true
```

> You cannot set **disableVPCCNI** to true if you are using the VPC CNI addon.

Some alternative CNIs provide for the replacement of kube-proxy, such as in [Calico](https://projectcalico.docs.tigera.io/maintenance/ebpf/enabling-ebpf#configure-kube-proxy) and [Cilium](https://docs.cilium.io/en/stable/gettingstarted/kubeproxy-free/). When enabling the kube-proxy alternative, the kube-proxy installed by EKS must be deleted. This can be done via the **disable** property of **kubeProxy** in **AWSManagedControlPlane**:

```yaml
kind: AWSManagedControlPlane
apiVersion: controlplane.cluster.x-k8s.io/v1beta1
metadata:
  name: "capi-managed-test-control-plane"
spec:
  region: "eu-west-2"
  sshKeyName: "capi-management"
  version: "v1.18.0"
  disableVPCCNI: true
  kubeProxy:
    disable: true
```

> You cannot set **disable** to true in **kubeProxy** if you are using the kube-proxy addon.

## Additional Information

See the [AWS documentation](https://docs.aws.amazon.com/eks/latest/userguide/pod-networking.html) for further details of EKS pod networking.
