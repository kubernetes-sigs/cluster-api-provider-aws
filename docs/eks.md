# EKS Support in the AWS Provider

## Overview

Experimental support for EKS has been introduced in 0.6.0 of the provider. Initially the following features are supported:

* Provisioning/managing an AWS EKS Cluster
* Upgrading the Kubernetes version of the EKS Cluster
* Creating a self-managed node group

The implementation introduces 3 new CRD kinds:

* AWSManagedControlPlane - specifies the EKS Cluster in AWS
* AWSManagedCluster - holds details of the EKS cluster for use by CAPI
* EKSConfig - used by Cluster API bootstrap provider EKS (CABPE)

And a new template is available in the templates folder for creating a managed workload cluster.

## Enabling EKS Support

You must explicitly enable the EKS support in the provider using feature flags. The following feature flags are supported:

* **EKS** - this will enable the EKS support in the manager (capa-controller-manager)
* **EKSEnableIAM** - by enabling this the controller will create the IAM role required by the EKS control plane. If this isn't enabled then you will need to manually create a role and specify the role name in the AWSManagedControlPlane.
* **EKSAllowAddRoles** - by enabling this you can add additional roles to the control plane role that is created. This has no affect unless used wtih __EKSEnableIAM__

The feature flags can be enabled when using `clusterctl` by setting the following environment variables to **true** (they all default to **false**):

* **EXP_EKS** - this is used to set the value of the **EKS** feature flag
* **EXP_EKS_IAM** - this is used to set the value of the **EKSEnableIAM** feature flag
* **EXP_EKS_ADD_ROLES** - this is used to set the value of the **EKSAllowAddRoles** feature flag

As an example, to enable EKS with IAM role creation:

```bash
export EXP_EKS=true
export EXP_EKS_IAM=true
clusterctl --infrastructure=aws
```

## Creating a EKS cluster

A new "managed" cluster template has been created that you can use with `clusterctl` to create a EKS cluster. To use the template:

```bash
clusterctl config cluster capi-eks-quickstart --flavour managed --kubernetes-version v1.17.3 --control-plane-machine-count=3 --worker-machine-count=3 > capi-eks-quickstart.yaml
```

NOTE: When creating an EKS cluster only the **MAJOR.MINOR** of the `-kubernetes-version` are taken into consideration. 

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

