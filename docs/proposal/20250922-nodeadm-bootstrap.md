---
title: Proposal EKS Support in CAPA for nodeadm
authors:
  - "@faiq"
reviewers:
creation-date: 2025-09-22
last-updated: 2025-09-22
status: proposed
see-also:
- https://github.com/kubernetes-sigs/cluster-api-provider-aws/pull/3518
replaces: []
superseded-by: []
---


## Table of Contents
- [Summary](#summary)
- [Motivation](#motivation)
  - [Goals](#goals)
  - [Non-Goals](#non-goals)
- [Proposal](#proposal)
  - [User Stories](#user-stories)
- [Alternatives](#alternatives)


 ## Summary

 Currently, EKS support in the Cluster API Provider for AWS (CAPA) is broken for Amazon Linux 2023 (AL2023) because the `bootstrap.sh` script is no longer supported. This proposal introduces a new Custom Resource Definition (CRD), `NodeadmConfig`, to handle the new `nodeadm` bootstrapping method required by AL2023. This approach is favored over modifying the existing `EKSConfig` type to maintain a cleaner API, avoid fields that are not cross-compatible between bootstrapping methods, and simplify the future deprecation of the `bootstrap.sh` implementation.

 -----

 ## Motivation

 Currently EKS support in CAPA is broken for AL2023 (Amazon Linux 2023) because the bootstrapping method that was previously being used to provision EKS nodes is no longer supported `bootstrap.sh`. Users who are using AL2023 see errors like this on the worker nodes:

 ```bash
 [root@localhost bin]# /etc/eks/bootstrap.sh default_dk-eks-133-control-plane


 \!\!\!\!\!\!\!\!\!\!
 \!\!\!\!\!\!\!\!\!\! ERROR: bootstrap.sh has been removed from AL2023-based EKS AMIs.
 \!\!\!\!\!\!\!\!\!\!
 \!\!\!\!\!\!\!\!\!\! EKS nodes are now initialized by nodeadm.
 \!\!\!\!\!\!\!\!\!\!
 \!\!\!\!\!\!\!\!\!\! To migrate your user data, see:
 \!\!\!\!\!\!\!\!\!\!
 \!\!\!\!\!\!\!\!\!\!     https://awslabs.github.io/amazon-eks-ami/nodeadm/
 \!\!\!\!\!\!\!\!\!\!

 ````

 In CAPA our implementation of the EKS bootstrapping method is currently tied to the `bootstrap.sh` script and is implemented by the `EKSConfig` type.

 Additionally, the EKS team is not publishing any more AmazonLinux (AL2) AMIs after November 26th, 2025, and Kubernetes version 1.32 is the last version for which AL2 AMIs will be released. This makes the transition to a new bootstrapping method for AL2023 urgent.

 ### Goals

   * Restore the ability to provision EKS nodes using CAPA with AL2023 AMIs.
   * Introduce a new, clean API (`NodeadmConfig`) specifically for the `nodeadm` bootstrap method.
   * Provide a clear upgrade path for users moving from `EKSConfig` (`bootstrap.sh`) to `NodeadmConfig` (`nodeadm`).
   * Make future deprecation of the `bootstrap.sh` implementation in `EKSConfig` easier.

 ### Non-Goals
   * Create a metatype that can handle both bootstrap.sh and nodeadm.
   * Handle Operating Systems with different bootstrapping mechanisms like bottlerocket.
 -----

 ## Proposal

This KEP proposes a new type that handles bootstrapping with `nodeadm` alone. This new type, `NodeadmConfig`, will wrap the API implementation for the Nodeadm option as a bootstrap provider.

 This approach is proposed due to drawbacks with the alternative of modifying the existing `EKSConfig` type, which would involve the introduction of new fields to distinguish between bootstrap methods and lead to a confusing API where some fields are only valid for one method.

 Examples of fields in the existing API that are no longer valid with `nodeadm`:

   * `ContainerRuntime`
   * `DNSClusterIP`
   * `DockerConfigJSON`
   * `APIRetryAttempts`
   * `PostBootstrapCommands`
   * `BootstrapCommandOverride`

 The **pros** of this approach are:

   * A cleaner API that’s more descriptive for each bootstrap method.
   * A new implementation will make deprecating EKSConfig’s `bootstrap.sh` implementation easier.

 The **cons** are:

   * The scope of work to support EKS nodes grows significantly and is pushed out.

 ### User Stories

   * As a cluster admin I need to provision nodes to my EKS cluster using Kubernetes 1.33 or higher
   * As a cluster admin I need to provision EKS worker nodes using the latest AL2023 AMIs
   * As a cluster admin I need to upgrade my existing EKS cluster nodes from an AL2-based version (e.g., 1.32) to an AL2023-based version (e.g., 1.33) with minimal disruption.

 ### API Design

 On a high level this new type `NodeadmConfig` wraps the API implementation for the Nodeadm option as a bootstrap provider.

 ```go
 // NodeadmConfigSpec defines the desired state of NodeadmConfig.
 type NodeadmConfigSpec struct {
 // Kubelet contains options for kubelet.
 // +optional
 Kubelet *KubeletOptions `json:"kubelet,omitempty"`

 // Containerd contains options for containerd.
 // +optional
 Containerd *ContainerdOptions `json:"containerd,omitempty"`

 // FeatureGates holds key-value pairs to enable or disable application features.
 // +optional
 FeatureGates map[Feature]bool `json:"featureGates,omitempty"`

 // PreBootstrapCommands specifies extra commands to run before bootstrapping nodes.
 // +optional
 PreBootstrapCommands []string `json:"preBootstrapCommands,omitempty"`

 // Files specifies extra files to be passed to user_data upon creation.
 // +optional
 Files []File `json:"files,omitempty"`

 // Users specifies extra users to add.
 // +optional
 Users []User `json:"users,omitempty"`

 // NTP specifies NTP configuration.
 // +optional
 NTP *NTP `json:"ntp,omitempty"`

 // DiskSetup specifies options for the creation of partition tables and file systems on devices.
 // +optional
 DiskSetup *DiskSetup `json:"diskSetup,omitempty"`

 // Mounts specifies a list of mount points to be setup.
 // +optional
 Mounts []MountPoints `json:"mounts,omitempty"`
 }
 ```

 -----

 ## Design Details

 ### Upgrade Strategy

A valid concern that CAPA users will have is upgrading existing clusters to machines that use the new bootstrap `Nodeadm` CRD. This KEP does not change the process. As before, the user will reference a new BootstrapConfigTemplate. However, the kind will change from EKSConfigTemplate to NodeadmConfigTemplate. 

 #### MachineDeployment Upgrade Example

 A user with a `MachineDeployment` using `EKSConfig` for Kubernetes v1.32 would upgrade to v1.33 by creating a new `NodeadmConfigTemplate` and updating the `MachineDeployment` to reference it and the new Kubernetes version. New machines are rolled out according to the `MachineDeployment` update strategy.

 **Before (v1.32 with `EKSConfigTemplate`):**

 ```yaml
 apiVersion: bootstrap.cluster.x-k8s.io/v1beta2
 kind: EKSConfigTemplate
 metadata:
   name: default132
 spec:
   template:
     spec:
       postBootstrapCommands:
         - "echo \"bye world\""
 ---
 apiVersion: cluster.x-k8s.io/v1beta1
 kind: MachineDeployment
 metadata:
   name: default
 spec:
   clusterName: default
   template:
     spec:
       bootstrap:
         configRef:
           apiVersion: bootstrap.cluster.x-k8s.io/v1beta2
           kind: EKSConfigTemplate
           name: default132
       infrastructureRef:
         kind: AWSMachineTemplate
         name: default132
       version: v1.32.0
 ````

 **After (v1.33 with `NodeadmConfigTemplate`):**

 ```yaml
 apiVersion: bootstrap.cluster.x-k8s.io/v1beta2
 kind: NodeadmConfigTemplate
 metadata:
   name: default
 spec:
   template:
     spec:
       preBootstrapCommands:
         - "echo \"hello world\""
 ---
 apiVersion: cluster.x-k8s.io/v1beta1
 kind: MachineDeployment
 metadata:
   name: default
 spec:
   clusterName: default
   template:
     spec:
       bootstrap:
         configRef:
           apiVersion: bootstrap.cluster.x-k8s.io/v1beta2
           kind: NodeadmConfigTemplate
           name: default
       infrastructureRef:
         kind: AWSMachineTemplate
         name: default
       version: v1.33.0
 ```

 #### MachinePool Upgrade Example

 The flow would be very similar for `MachinePools`. A user would update the `MachinePool` resource to reference a new `NodeadmConfigTemplate` and the target Kubernetes version.

 **Before (v1.32 with `EKSConfigTemplate`):**

 ```yaml
 apiVersion: bootstrap.cluster.x-k8s.io/v1beta2
 kind: EKSConfigTemplate
 metadata:
   name: default-132
 spec:
   template:
     spec: {}
 ---
 apiVersion: cluster.x-k8s.io/v1beta1
 kind: MachinePool
 metadata:
   name: default
 spec:
   clusterName: default
   template:
     spec:
       bootstrap:
         configRef:
           apiVersion: bootstrap.cluster.x-k8s.io/v1beta2
           kind: EKSConfigTemplate
           name: default-132
       infrastructureRef:
         kind: AWSMachinePool
         name: default
       version: v1.32.0
 ```

 **After (v1.33 with `NodeadmConfigTemplate`):**

 ```yaml
 apiVersion: bootstrap.cluster.x-k8s.io/v1beta2
 kind: NodeadmConfigTemplate
 metadata:
   name: default-133
 spec:
   template:
     spec:
       preBootstrapCommands:
         - "echo \"hello from v1.33.0\""
 ---
 apiVersion: cluster.x-k8s.io/v1beta1
 kind: MachinePool
 metadata:
   name: default
 spec:
   clusterName: default
   template:
     spec:
       bootstrap:
         configRef:
           apiVersion: bootstrap.cluster.x-k8s.io/v1beta2
           kind: NodeadmConfigTemplate
           name: default-133
       infrastructureRef:
         apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
         kind: AWSMachinePool
         name: default
       version: v1.33.0
 ```

 ### Test Plan

    * Unit tests for the new code.
    * Integration tests for new Nodeadm Controller.
    * E2e tests exercising the migration from EKSConfig to NodeadmConfig,


 -----

 ## Alternatives

 The primary alternative considered was to modify the existing `EKSConfig` type to support `nodeadm`. Currently, there’s work being done upstream to address this gap. On a high level, [this PR](https://github.com/kubernetes-sigs/cluster-api-provider-aws/pull/5553) is adding a new bootstrapping implementation to the existing `EKSConfig` type with some additional API fields to distinguish between bootstrap methods.

 However, there are some drawbacks with this implementation regarding the API design:

   * **Introduction of new fields to distinguish between bootstrap methods**: This complicates the API.
   * **Fields that are valid for `bootstrap.sh` are not valid for `nodeadm` and vice versa**: This would lead to a confusing user experience where users could set fields that have no effect for their chosen bootstrap method.
