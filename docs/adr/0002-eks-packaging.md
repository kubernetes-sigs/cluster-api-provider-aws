# 2. EKS Controller Packaging

* Status: accepted
* Date: 2020-11-02 [YYY-MM-DD - date of the decision] <!-- mandatory -->
* Authors:@richardcase
* Deciders: @richardcase @randomvariable

## Context

The EKS controllers are implemented in a number of separate managers:

* **Infrastructure manager** - includes the controllers for AWSManagedCluster, AWSManagedMachinePools, AWSMachinePools (used for both managed/unmanaged clusters), AWSMachine (used for both managed/unmanaged clusters)
* **Control Plane manager** - handles the the AWSManagedControlPlane kind that creates the EKS control plane
* **Bootstrap manager** - handles EKSConfig which generates the bootstrap config the EC2 instances when creating an EKS cluster

To create EKS clusters using CAPA you need to do **ALL** of the following:

1. Enable the EKS functionality in the **Infrastructure manager** by using the **EKS** feature flag (which can be controlled via env var `EXP_EKS`)
2. Install the `aws-eks` control plane provider (i.e. the **Control Plane manager**)
3. Install the `aws-eks` bootstrap provider (i.e. **Bootstrap manager** )

An error occurs if you enable EKS functionality in the **infrastructure manager** but don't install the **control plane manager** (for example if you have the EKS environment variables set):

```bash
E1028 20:55:12.614531       9 source.go:116] controller-runtime/source "msg"="if kind is a CRD, it should be installed before calling Start" "error"="no matches for kind \"AWSManagedControlPlane\" in version \"controlplane.cluster.x-k8s.io/v1alpha3\""  "kind"={"Group":"controlplane.cluster.x-k8s.io","Kind":"AWSManagedControlPlane"}
```

This was reported in issue [#2078](https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/2078) and on investigation the cause was that both the AWSManagedClusterReconciler and AWSManagedMachinePoolReconciler in the **infrastructure manager** watch AWSManagedControlPlane which is only installed when you enable the aws-eks control plane (which applies eks-controlplane-components.yaml).

A number of potential solutions where suggested on the issue.

## Decision

During the CAPA office hours call on 2nd November 2020 it was decided that the **infrastructure manager** would be updated to test for the existance of the **AWSManagedControlPlane** CRD if the EKS feature flag is specified. If the **AWSManagedControlPlane** CRD isn't available an error will be reported stating that the `aws-eks` control plane needs to be installed. This change will be carried under issue [#2078](https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/2078).

## Consequences

The **infrastructure manager** will need to be updated to include this test.
