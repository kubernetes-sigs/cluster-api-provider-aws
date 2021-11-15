# Prerequisites

To use EKS you must give the controller the required permissions. The easiest way to do this is by using `clusterawasadm`. For instructions on how to do this see the [prerequisites](../using-clusterawsadm-to-fulfill-prerequisites.md).

When using `clusterawsadm` and enabling EKS support a new IAM role will be created for you called **eks-controlplane.cluster-api-provider-aws.sigs.k8s.io**. This role is the IAM role that will be used for the EKS control plane if you don't specify your own role and if **EKSEnableIAM** isn't enabled (see the [enabling docs](enabling.md) for further information).

Additionally using `clusterawsadm` will add permissions to the **controllers.cluster-api-provider-aws.sigs.k8s.io** policy for EKS to function properly.
