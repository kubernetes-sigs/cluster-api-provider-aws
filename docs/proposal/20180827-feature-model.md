# Cluster API Provider AWS Feature Set

## Introduction

We wish to make a [feature model](https://en.wikipedia.org/wiki/Feature_model) which allows us to see the common functionality that will need to be shared between AWS implementations.

Give each feature a unique number. Reference a feature by its number. If possible, provide a justification or requirement for the feature, which will help with prioritisation for a minimum viable product.

You can also write constraints:
Feature B is optional sub-feature of A
Feature C is a mandatory sub-feature of B
Feature E is an alternative feature to C, D
Feature F is mutually exclusive of Feature E
Feature G requires Feature B

A minimum viable product will be a configuration of features, which when solved for all constraints provides the minimum list of features that will need to be developed.

Different MVPs may be possible - e.g. EKS vs. not EKS, but they may rely on shared components which will become the critical path.

## Feature Set

### 0: AWS Cluster Provider

### 1: VPC Selection

#### 2: The provider will need the ability to create a new VPC

* Constraint: Mandatory sub-feature of [1](#1-vpc-selection)

#### 3: The provider provisions in an existing VPC, selecting the default VPC if none is specified

* Constraint: Optional alternative to [2](#2-the-provider-will-need-the-ability-to-create-a-new-vpc)
* Requirement: Some customers may wish to reuse VPCs in which they have existing infrastructure.

### 45: Etcd location

#### 46: The provider deploys etcd as part of the control plane

* Constraint: Mandatory sub-feature of [45](#45-etcd-location)
* Requirement: For simple clusters a colocated etcd is the easiest way to operate a cluster

#### 47: The provider deploys etcd externally to the control plane

* Constraint: Alternative to [46](#46-the-provider-deploys-etcd-as-part-of-the-control-plane)
* Requirement: For larger clusters etcd placed external to the control plane allows for independent control plane and datastore scaling
  
#### 48: The provider can connect to a pre-existing etcd cluster

* Constraint: Optional subfeature of [45](#45-etcd-location)
* Requirement: An existing etcd store could be used to replace a cluster during upgrade or a complete cluster restore.

### 5: Control plane placement

#### 6: The provider deploys the control plane in public subnets

* Constraint: Mandatory sub-feature of [5](#5-control-plane-placement)
* Requirement: For simple clusters without bastion hosts, allowing users to break-glass SSH to control plane nodes
  
#### 7: The provider deploys the control plane in private subnets

* Constraint: Alternative to [6](#6-the-provider-deploys-the-control-plane-in-public-subnets)
* Requirement: [AWS Well-Architected SEC 5](https://d1.awsstatic.com/whitepapers/architecture/AWS_Well-Architected_Framework.pdf), security requirements may require access via bastion hosts

#### 8: The provider deploys control plane components to a single AZ

* Constraint: Mandatory sub-feature of [5](#5-control-plane-placement)
* Requirement: Architectural requirement for a particular customer workload

#### 9: The provider deploys control plane components across multiple AZs

* Constraint: Alternative to [8](#8-the-provider-deploys-control-plane-components-to-a-single-az)
* Requirement: Robustness of control plane components

### 10: Worker node placement

#### 11: Provider deploys worker nodes deployed to public subnets

* Constraint: Mandatory sub-feature of [10](#10-worker-node-placement)
* Requirement: For simple clusters without bastion hosts, allowing users to break-glass SSH to worker nodes
  
#### 12: Provider deploys worker nodes to private subnets

* Constraint: Alternative to [11](#11-provider-deploys-worker-nodes-deployed-to-public-subnets)
* Requirement: AWS Well-Architected SEC 5, security requirements may require access via bastion hosts / VPN / Direct Connect
  
#### 13: Provider deploys worker nodes to single AZ

* Constraint: Mandatory sub-feature of [10](#10-worker-node-placement)
* Requirement: Architectural requirement for a particular customer workload

#### 14: Provider deploys worker nodes across multiple AZs

* Constraint: Alternative to [13](#13-provider-deploys-worker-nodes-to-single-az)
* Requirement: Robustness of cluster

#### 15: Deploy worker nodes to a placement group

* Constraint: Optional sub-feature of [10](#10-worker-node-placement)
* Requirement: HPC type workload that requires fast interconnect between nodes

#### 16: The provider deploys worker nodes to shared instances

* Constraint: Mandatory sub-feature of [10](#10-worker-node-placement)
* Requirement: Default behaviour / cost / availability of instances

#### 17: The provider deploys worker nodes to dedicated EC2 instances

* Constraint: Optional alternative to [16](#16-the-provider-deploys-worker-nodes-to-shared-instances)
* Requirement: License requirements for a particular workload (e.g. Oracle) may require a dedicated instance

### 18: Worker node scaling methodology

#### 19: Worker nodes are deployed individually or in batches not using auto-scaling groups

* Constraint: Mandatory sub-feature of [18](#18-worker-node-scaling-methodology)

#### 20: Worker nodes are deployed via Auto-Scaling Groups using MachineSets

* Constraint: Alternative to [19](#19-worker-nodes-are-deployed-individually-or-in-batches-not-using-auto-scaling-groups)
* Note: The implementation here would be significantly different to [19](#19-worker-nodes-are-deployed-individually-or-in-batches-not-using-auto-scaling-groups).

### 21: API Server Access

#### 22: The API server is publicly accessible

* Constraint: Mandatory sub-feature of [21](#21-api-server-access)
* Requirement: Standard way of accessing k8s

#### 23: The API server is not publicly accessible

* Constraint: Alternative to [22](#22-the-api-server-is-publicly-accessible)
* Requirement: Security requirement (e.g. someone’s interpretation of UK OFFICIAL) prohibits making API server endpoint publicly accessible

#### 31: The API server is connected to a VPC via PrivateLink

* Constraint: Sub-feature of [23](#23-the-api-server-is-not-publicly-accessible) & [25](#25-the-control-plane-is-eks)
* Requirement: Compliance requirements for API traffic to not transit public internet, e.g. UK OFFICIAL-SENSITIVE workloads. AWS recommend for FedRamp(?) and UK-OFFICIAL to use VPC or PrivateLink endpoints to connect publicly accessible regional services to VPCs to prevent traffic exiting the internal AWS network. The actual EKS endpoint, for example may still present itself with a public load balancer endpoint even if it’s connected by PrivateLink to the VPC.

#### 43: The API server is accessible via a load balancer

* Constraint: Sub-feature of [22](#22-the-api-server-is-publicly-accessible) & [23](#23-the-api-server-is-not-publicly-accessible)
* Requirement: For potential HA access OR public/private subnet distinctions, the API server is accessed via an AWS load balancer.

#### 44: The API server is accessed directly via the IP address of cluster nodes hosting the API server

* Constraint: Alternative to [43](#43-the-api-server-is-accessible-via-a-load-balancer)
* Requirement: The IP addresses of each node hosting an API server is registered in DNS

### 24: Type of control plane

#### 25: The control plane is EKS

* Constraint: Mandatory sub-feature of [24](#25-type-of-control-plane)
* Requirement: Leverage AWS to provide heavy lifting of control plane deployment & operations. Also, meets compliance requirements: [UK OFFICIAL IL2 & OFFICIAL-SENSITIVE/HIGH IL3](https://www.digitalmarketplace.service.gov.uk/g-cloud/services/760447139328659)

#### 26: The control plane is managed within the provider

* Constraint: Alternative to [25](#25-the-control-plane-is-eks)
* Requirement: Customer requires functionality not provided by EKS (e.g. admission controller, non-public API endpoint)

### 33: CRI

#### 34: The provider deploys a [credential helper](https://github.com/awslabs/amazon-ecr-credential-helper) for ECR

* Constraint: Mandatory sub-feature of [33](#33-cri)
* Requirement: [AWS Well-Architected SEC 3](https://d1.awsstatic.com/whitepapers/architecture/AWS_Well-Architected_Framework.pdf)

### 35: Container Hosts

#### 36: The provider deploys to Amazon Linux 2

* Constraint: Mandatory sub-feature of [35](#35-container-hosts)
* Requirement: Parity with AWS recommendations

#### 37: The provider deploys to CentOS / Ubuntu

* Constraint: Alternative to [36](#36-the-provider-deploys-to-amazon-linux-2)
* Requirement: Greater familiarity in the community (particularly Ubuntu), organisational requirements?

#### 38: The provider deploys from arbitrary AMIs

* Constraint: Alternative to [36](#36-the-provider-deploys-to-amazon-linux-2). Sub-feature of [37](#37-the-provider-deploys-to-centos--ubuntu)
* Requirement: Compliance requirement may require AMI that passes [CIS Distribution Independent Linux Benchmark](https://www.cisecurity.org/benchmark/distribution_independent_linux/), or EC2 instances should have encrypted EBS root volumes for data loss prevention, requiring AMIs in the customer’s account.

#### 39: The provider allows kubelet configuration to be customised, e.g. “--allow-privileged”

* Constraint: Sub-feature of [35](#35-container-hosts)
* Requirement: [NIST 800-190, control 4.4.3](https://nvlpubs.nist.gov/nistpubs/specialpublications/nist.sp.800-190.pdf) recommends disabling “allow-privileged”

#### 40: Arbitrary customisation of bootstrap script

* Constraint: Sub-feature of [35](#35-container-hosts)
* Requirement: Organisational or security requirements, e.g. NIST 800-190 & AWS Well-Architected controls recommend the installation of additional file integrity tools such as OSSEC, Tripwire etc…, some organisations may even mandate antivirus, etc… Cannot encode all of this as additional options, so some mix of 38 plus some ability to customise bootstrap would satisfy this without bringing too much variability into scope.

### 41: API Server configuration

#### 42: The provider allows customisation of API Server

* Constraint: Sub-feature of [41](#41-api-server-configuration)
* Requirement: Example - Would need to enable AdmissionController for Istio automatic sidecar injection (note, EKS doesn’t allow customisation of webhooks at present, but may in the future).

## TODO

* HA / non-HA installs?

## Out of Scope

Anything that can be applied after the cluster has come up with kubectl is not a cluster-api responsibility, including:

* Monitoring / Logging
* Many of the CNI options (at least Calico & AWS VPC CNI)
* IAM identity for pods (e.g. kube2iam, kiam etc…)
* ALB ingress

These should be addressed with documentation - we do not want the cluster-api provider to be a package manager for Kubernetes manifests. In addition, @roberthbailey has stated on 2018/08/14 that Google is working on a declarative cluster add on manager to be presented to sig-cluster-lifecycle for discussion later.