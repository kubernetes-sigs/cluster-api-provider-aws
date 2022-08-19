# Pod Networking

When creating a EKS cluster the Amazon VPC CNI will be used by default for Pod Networking.

> When using the AWS Console to create an EKS cluster with a Kubernetes version of v1.18 or greater you are required to select a specific version of the VPC CNI to use.

## Using the VPC CNI Addon
You can use an explicit version of the Amazon VPC CNI by using the **vpc-cni** EKS addon. See the [addons](./addons.md) documentation for further details of how to use addons.

## Using Custom VPC CNI Configuration
If your use case demands [custom networking](https://docs.aws.amazon.com/eks/latest/userguide/cni-custom-network.html) VPC CNI configuration you might already be familiar with the [helm chart](https://github.com/aws/amazon-vpc-cni-k8s) which helps with the process. This gives you access to ENI Configs and you can set Environment Variables on the `aws-node` DaemonSet where the VPC CNI runs. CAPA is able to tune the same DaemonSet through Kubernetes.

The following example shows how to turn on custom network config and set a [label definition](https://github.com/aws/amazon-vpc-cni-k8s#eni_config_label_def).

```yaml
kind: AWSManagedControlPlane
apiVersion: controlplane.cluster.x-k8s.io/v1beta1
metadata:
  name: "capi-managed-test-control-plane"
spec:
  vpcCni:
    env:
    - name: AWS_VPC_K8S_CNI_CUSTOM_NETWORK_CFG
      value: "true" 
    - name: ENABLE_PREFIX_DELEGATION
      value: "true"
```

### Increase node pod limit
You can increase the pod limit per-node as [per the upstream AWS documentation](https://aws.amazon.com/blogs/containers/amazon-vpc-cni-increases-pods-per-node-limits/). You'll need to enable the `vpc-cni` plugin addon on your EKS cluster as well as enable prefix assignment mode through the `ENABLE_PREFIX_DELEGATION` environment variable.

```yaml
kind: AWSManagedControlPlane
apiVersion: controlplane.cluster.x-k8s.io/v1beta1
metadata:
  name: "capi-managed-test-control-plane"
spec:
  vpcCni:
    env:
    - name: AWS_VPC_K8S_CNI_CUSTOM_NETWORK_CFG
      value: "true" 
    - name: ENABLE_PREFIX_DELEGATION
      value: "true"
  addons:
  - name: vpc-cni
    version: <replace_with_version>
    conflictResolution: overwrite
  associateOIDCProvider: true
  disableVPCCNI: false
```

## Using an alternative CNI

There may be scenarios where you do not want to use the Amazon VPC CNI. EKS supports a number of alternative CNIs such as Calico and Weave Net (see [docs](https://docs.aws.amazon.com/eks/latest/userguide/alternate-cni-plugins.html) for full list).

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

> You cannot set **disableVPCCNI** to true if you are using the VPC CNI addon or if you have specified a secondary CIDR block.

## Additional Information

See the [AWS documentation](https://docs.aws.amazon.com/eks/latest/userguide/pod-networking.html) for further details of EKS pod networking.
