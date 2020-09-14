# EKS Support in the AWS Provider

## Overview

Experimental support for EKS has been introduced in 0.6.0 of the provider. Initially the following features are supported:

* Provisioning/managing an AWS EKS Cluster
* Upgrading the Kubernetes version of the EKS Cluster
* Creating a self-managed node group

The implementation introduces 3 new CRD kinds:

* AWSManagedControlPlane - specifies the EKS Cluster in AWS and used by the Cluster API AWS Managed Control plane (MACP)
* AWSManagedCluster - holds details of the EKS cluster for use by CAPI
* EKSConfig - used by Cluster API bootstrap provider EKS (CABPE)

And a new template is available in the templates folder for creating a managed workload cluster.

## Prerequisites

To use EKS you must give the controller the required permissions. The easiest way to do this is by using `clusterawasadm`. For instructions on how to do this see the [AWS prerequisites](prerequisites.md).

When using `clusterawsadm` and enabling the managed control plane a new IAM role will be created for you called **eks-controlplane.cluster-api-provider-aws.sigs.k8s.io**. This role is the IAM role that will be used for EKS if you don't specify your own role and if **EKSEnableIAM** isn't enabled.

Additionally using `clusterawsadm` will add permissions to the **controllers.cluster-api-provider-aws.sigs.k8s.io** policy for EKS to function properly.  

## Enabling EKS Support

You must explicitly enable the EKS support in the provider by doing the following:

* Enabling support in the infrastructure manager (capa-controller-manager) by enabling the **EKS** feature flag
* Add the AWS Managed Control Plane Provider
* Add the Cluster API bootstrap provider EKS

### Enabling the **EKS** feature 

Enabling the **EKS** feature on the core infrastructure managercan be done using `clusterctl` by setting the following environment variables to **true** (they all default to **false**):

* **EXP_EKS** - this is used to set the value of the **EKS** feature flag

As an example:

```bash
export EXP_EKS=true
clusterctl --infrastructure=aws
```

### Adding the AWS Managed Control Plane and Bootstrap Providers

Create the` ~/.cluster-api/clusterctl.yaml` file with the following contents:

```yaml
providers:
  - name: "eks"
    url: "https://github.com/kubernetes-sigs/cluster-api-provider-aws/releases/latest/eks-bootstrap-components.yaml"
    type: "BootstrapProvider"
  - name: "eks"
    url: "https://github.com/kubernetes-sigs/cluster-api-provider-aws/releases/latest/eks-controlplane-components.yaml"
    type: "ControlPlaneProvider"
```

The EKS functionality in the control plane has 2 feature flags you can enable:

* **EKSEnableIAM** - by enabling this the controller will create the IAM role required by the EKS control plane. If this isn't enabled then you will need to manually create a role and specify the role name in the AWSManagedControlPlane.
* **EKSAllowAddRoles** - by enabling this you can add additional roles to the control plane role that is created. This has no affect unless used wtih __EKSEnableIAM__

The feature flags can be enabled when using `clusterctl` by setting the following environment variables to **true** (they all default to **false**):

* **EXP_EKS_IAM** - this is used to set the value of the **EKSEnableIAM** feature flag
* **EXP_EKS_ADD_ROLES** - this is used to set the value of the **EKSAllowAddRoles** feature flag

As an example, to enable the control plane and bootstrap providers with IAM role creation:

```bash
export EXP_EKS_IAM=true
clusterctl --infrastructure=aws --boostrap=eks --controlplane=eks
```

## Creating a EKS cluster

A new "managed" cluster template has been created that you can use with `clusterctl` to create a EKS cluster. To use the template:

```bash
clusterctl config cluster capi-eks-quickstart --flavor managed --kubernetes-version v1.17.3 --worker-machine-count=3 > capi-eks-quickstart.yaml
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

### CAPI kubeconfig

This kubeconfig is used internally by CAPI and shouldn't be used outside of the management server. It is used by CAPI to perform operations, such as draining a node. The name of the secret that contains the kubeconfig will be `[cluster-name]-kubeconfig` where you need to replace **[cluster-name]** with the name of your cluster. Note that there is NO `-user` in the name.

The kubeconfig is regenerated every `sync-period` as the token that is embedded in the kubeconfig is only valid for a short period of time. When EKS support is enabled the maximum sync period is 10 minutes. If you try to set `--sync-period` to greater than 10 minutes then an error will be raised.

## Control Plane Upgrade

Upgrading the kubernetes version of the control plane is supported by the provider. To perform an upgrade you need to update the `version` in the spec of the `AWSManagedControlPlane`. Once the version has changed the provider will handle the upgrade for you.

You can only upgrade a EKS cluster by 1 minor version at a time. If you attemp to upgrade the version by more then 1 minor version the provider will ensure the upgrade is done in multiple steps of 1 minor version. For example upgrading from v1.15 to v1.17 would result in your cluster being upgraded v1.15 -> v1.16 first and then v1.16 to v1.17.
