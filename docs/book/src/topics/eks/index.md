# EKS Support in the AWS Provider

- **Feature status:** Experimental
- **Feature gate (required):** EKS=true
- **Feature gate (optional):** EKSEnableIAM=true,EKSAllowAddRoles=true

## Overview

Experimental support for EKS has been introduced in the AWS provider. Currently the following features are supported:

- Provisioning/managing an Amazon EKS Cluster
- Upgrading the Kubernetes version of the EKS Cluster
- Attaching a self-managed machines as nodes to the EKS cluster
- Creating a machine pool and attaching it to the EKS cluster. See [machine pool docs for details](../machinepools.md)
- Creating a managed machine pool and attaching it to the EKS cluster. See [machine pool docs for details](../machinepools.md)

The implementation introduces new CRD kinds:

- AWSManagedControlPlane - specifies the EKS Cluster in AWS and used by the Cluster API AWS Managed Control plane (MACP)
- AWSManagedCluster - holds details of the EKS cluster for use by CAPI
- AWSManagedMachinePool - defines the managed node pool for the cluster
- EKSConfig - used by Cluster API bootstrap provider EKS (CABPE)

And a number of new templates are available in the templates folder for creating a managed workload cluster.


## SEE ALSO

* [Prerequisites](prerequisites.md)
* [Enabling EKS Support](enabling.md)
* [Creating a cluster](creating-a-cluster.md)
* [Using EKS Console](eks-console.md)
* [Cluster Upgrades](cluster-upgrades.md)