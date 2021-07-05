# Creating a EKS cluster

New "eks" cluster templates have been created that you can use with `clusterctl` to create a EKS cluster. To create a EKS cluster with self-managed nodes (a.k.a machines):

```bash
clusterctl generate cluster capi-eks-quickstart --flavor eks --kubernetes-version v1.17.3 --worker-machine-count=3 > capi-eks-quickstart.yaml
```

To create a EKS cluster with a managed node group (a.k.a managed machine pool):

```bash
clusterctl generate cluster capi-eks-quickstart --flavor eks-managedmachinepool --kubernetes-version v1.17.3 --worker-machine-count=3 > capi-eks-quickstart.yaml
```

NOTE: When creating an EKS cluster only the **MAJOR.MINOR** of the `-kubernetes-version` is taken into consideration.

## Kubeconfig

When creating an EKS cluster 2 kubeconfigs are generated and stored as secrets in the managmenet cluster. This is different to when you create a non-managed cluster using the AWS provider.

### User kubeconfig

This should be used by users that want to connect to the newly created EKS cluster. The name of the secret that contains the kubeconfig will be `[cluster-name]-user-kubeconfig` where you need to replace **[cluster-name]** with the name of your cluster. The **-user-kubeconfig** in the name indicates that the kubeconfig is for the user use.

To get the user kubeconfig for a cluster named `managed-test` you can run a command similar to:

```bash
kubectl --namespace=default get secret managed-test-user-kubeconfig \
   -o jsonpath={.data.value} | base64 --decode \
   > managed-test.kubeconfig
```

### Cluster API (CAPI) kubeconfig

This kubeconfig is used internally by CAPI and shouldn't be used outside of the management server. It is used by CAPI to perform operations, such as draining a node. The name of the secret that contains the kubeconfig will be `[cluster-name]-kubeconfig` where you need to replace **[cluster-name]** with the name of your cluster. Note that there is NO `-user` in the name.

The kubeconfig is regenerated every `sync-period` as the token that is embedded in the kubeconfig is only valid for a short period of time. When EKS support is enabled the maximum sync period is 10 minutes. If you try to set `--sync-period` to greater than 10 minutes then an error will be raised.
