# IAM Permissions

## Required to use clusterawasadm to provision IAM roles via CloudFormation

If using `clusterawsadm` to automate deployment of IAM roles via CloudFormation,
you must have IAM administrative access as `clusterawsadm` will provision IAM
roles and policies.

## Required by Cluster API Provider AWS controllers

The Cluster API Provider AWS controller requires permissions to use EC2, ELB
Autoscaling and optionally EKS. If provisioning IAM roles using `clusterawsadm`,
these will be set up as the `controllers.cluster-api-provider-aws.sigs.k8s.io`
IAM Policy, and attached to the `controllers.cluster-api-provider-aws.sigs.k8s.io`
and `control-plane.cluster-api-provider-aws.sigs.k8s.io` IAM roles.

### EC2 Provisioned Kubernetes Clusters

``` json
{{#include ../../../../out/AWSIAMManagedPolicyControllers.json}}
```

### With EKS Support

``` json
{{#include ../../../../out/AWSIAMManagedPolicyControllersWithEKS.json}}
```

## Required by the Kubernetes AWS Cloud Provider

These permissions are used by the Kubernetes AWS Cloud Provider. If you are
running with the in-tree cloud provider, this will typically be used by the
`controller-manager` pod in the `kube-system` namespace.

If provisioning IAM roles using `clusterawsadm`,
these will be set up as the `control-plane.cluster-api-provider-aws.sigs.k8s.io`
IAM Policy, and attached to the `control-plane.cluster-api-provider-aws.sigs.k8s.io`
IAM role.

``` json
{{#include ../../../../out/AWSIAMManagedPolicyCloudProviderControlPlane.json}}
```
## Required by all nodes

All nodes require these permissions in order to run, and are used by the AWS
cloud provider run by kubelet.

If provisioning IAM roles using `clusterawsadm`,
these will be set up as the `nodes.cluster-api-provider-aws.sigs.k8s.io`
IAM Policy, and attached to the `nodes.cluster-api-provider-aws.sigs.k8s.io`
IAM role.


``` json
{{#include ../../../../out/AWSIAMManagedPolicyCloudProviderNodes.json}}
```

When using EKS, the `AmazonEKSWorkerNodePolicy` and `AmazonEKS_CNI_Policy`
AWS managed policies will also be attached to
`nodes.cluster-api-provider-aws.sigs.k8s.io` IAM role.
