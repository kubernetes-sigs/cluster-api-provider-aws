# Creating a EKS cluster

New "eks" cluster templates have been created that you can use with `clusterctl` to create a EKS cluster. To create a EKS cluster with self-managed nodes (a.k.a machines):

```bash
clusterctl generate cluster capi-eks-quickstart --flavor eks --kubernetes-version v1.22.9 --worker-machine-count=3 > capi-eks-quickstart.yaml
```

To create a EKS cluster with a managed node group (a.k.a managed machine pool):

```bash
clusterctl generate cluster capi-eks-quickstart --flavor eks-managedmachinepool --kubernetes-version v1.22.9 --worker-machine-count=3 > capi-eks-quickstart.yaml
```

NOTE: When creating an EKS cluster only the **MAJOR.MINOR** of the `-kubernetes-version` is taken into consideration.

By default CAPA relies on the default EKS cluster upgrade policy, which at the moment of writing is EXTENDED support.
See more info about [cluster upgrade policy](https://docs.aws.amazon.com/eks/latest/userguide/view-upgrade-policy.html)

## Choosing a Bootstrap Provider: EKSConfig vs. NodeadmConfig

With the introduction of Amazon Linux 2023 (AL2023), the bootstrapping method for EKS nodes has changed. Cluster API Provider AWS (CAPA) supports two bootstrap providers for EKS:

1.  **`EKSConfig`**: The original bootstrap provider. It uses the legacy `bootstrap.sh` script and is intended for use with **Amazon Linux 2 (AL2)** AMIs.
2.  **`NodeadmConfig`**: The new bootstrap provider. It uses the modern `nodeadm` tool and is **required** for **Amazon Linux 2023 (AL2023)** AMIs.

### When to use which provider

The provider you must use depends on the Amazon Machine Image (AMI) and Kubernetes version you are targeting. Amazon Linux 2 AMIs are only supported for Kubernetes v1.32 and older.

| Bootstrap Provider | AMI Type | Kubernetes Version |
| --- | --- | --- |
| `EKSConfig` | Amazon Linux 2 (AL2) | $\le$ v1.32 |
| `NodeadmConfig` | Amazon Linux 2023 (AL2023) | $\ge$ v1.33 |

When you generate a cluster, you will need to ensure your `MachineDeployment` or `MachinePool` references the correct bootstrap template `kind`.

NOTE:

- [The EKS team stopped publishing Al2 AMIs for Kubernetes versions 1.33 and higher.](https://awslabs.github.io/amazon-eks-ami/usage/al2/)
- [Amazon Linux 2 end of support date (End of Life, or EOL) will be on 2026-06-30.](https://aws.amazon.com/amazon-linux-2/faqs/)

**For AL2 / K8s $\le$ v1.32, use `EKSConfigTemplate`:**
```yaml
apiVersion: cluster.x-k8s.io/v1beta1
kind: MachineDeployment
metadata:
  name: default
spec:
  template:
    spec:
      bootstrap:
        configRef:
          apiVersion: bootstrap.cluster.x-k8s.io/v1beta2
          kind: EKSConfigTemplate  # <-- Uses bootstrap.sh
          name: default-132
      version: v1.32.0
```

### Secrets Manager

Amazon Linux 2023 does not have the proper tooling to use the secrets manager flow for bootstrapping. CAPA uses a [custom cloud-init datasource](https://github.com/kubernetes-sigs/image-builder/pull/1583) to fetch the secure contents like the `kubeadm` tokens from secrets manager. Crucially, there is no current support for publishing CAPA-compatible AL2023 AMIs that include this necessary custom cloud-init datasource.

Therefore, whenever creating `AWSMachineTemplate` objects  `insecureSkipSecretsManager` must be set to true.

```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
kind: AWSMachineTemplate
metadata:
  name: default
spec:
  template:
    spec:
      cloudInit:
        insecureSkipSecretsManager: true
      ami:
        eksLookupType: AmazonLinux2023
      iamInstanceProfile: nodes.cluster-api-provider-aws.sigs.k8s.io
      instanceType: m5a.16xlarge
```

## Kubeconfig

When creating an EKS cluster 2 kubeconfigs are generated and stored as secrets in the management cluster. This is different to when you create a non-managed cluster using the AWS provider.

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

There are three keys in the CAPI kubeconfig for eks clusters:

| keys        | purpose                                                                                                                                                                            |
|-------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| value       | contains a complete kubeconfig with the cluster admin user and token embedded                                                                                                      |
| relative    | contains a kubeconfig with the cluster admin user, referencing the token file in a relative path - assumes you are mounting all the secret keys in the same dir                    |
| single-file | contains the same token embedded in the complete kubeconfig, it is separated into a single file so that existing APIMachinery can reload the token file when the secret is updated |

The secret contents are regenerated every `sync-period` as the token that is embedded in the kubeconfig and token file is only valid for a short period of time. When EKS support is enabled the maximum sync period is 10 minutes. If you try to set `--sync-period` to greater than 10 minutes then an error will be raised.
