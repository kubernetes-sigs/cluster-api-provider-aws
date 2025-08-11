# Enabling EKS Support

Support for EKS is enabled by default when you use the AWS infrastructure provider. For example:

```shell
clusterctl init --infrastructure aws
```

## Enabling optional **EKS** features

There are additional EKS experimental features that are disabled by default. The sections below cover how to enable these features.

### Machine Pools

To enable support for machine pools the **MachinePool** feature flag must be set to to **true**. This can be done using the **EXP_MACHINE_POOL** environment variable:

```shell
export EXP_MACHINE_POOL=true
clusterctl init --infrastructure aws
```

See the [machine pool documentation](../machinepools.md) for further information.

NOTE: you will need to enable the creation of the default IAM role. The easiest way is using `clusterawsadm`, for instructions see the [prerequisites](../using-clusterawsadm-to-fulfill-prerequisites.md).

### IAM Roles Per Cluster

By default EKS clusters will use the same IAM roles (i.e. control plane, node group roles). There is a feature that allows each cluster to have its own IAM roles. This is done by enabling the **EKSEnableIAM** feature flag. This can be done before running `clusterctl init` by using the the **CAPA_EKS_IAM** environment variable:

```shell
export CAPA_EKS_IAM=true
clusterctl init --infrastructure aws
```

NOTE: you will need the correct prerequisities for this. The easiest way is using `clusterawsadm` and setting `iamRoleCreation` to true, for instructions see the [prerequisites](../using-clusterawsadm-to-fulfill-prerequisites.md).

### Additional Control Plane Roles

You can add additional roles to the control plane role that is created for an EKS cluster. To use this you must enable the **EKSAllowAddRoles** feature flag. This can be done before running `clusterctl init` by using the **CAPA_EKS_ADD_ROLES** environment variable:

```shell
export CAPA_EKS_IAM=true
export CAPA_EKS_ADD_ROLES=true
clusterctl init --infrastructure aws
```

NOTE: to use this feature you must also enable the **CAPA_EKS_IAM** feature.

### EKS Fargate Profiles

You can use Fargate Profiles with EKS. To use this you must enable the **EKSFargate** feature flag. This can be done before running `clusterctl init` by using the **EXP_EKS_FARGATE** environmnet variable:

```shell
export EXP_EKS_FARGATE=true
clusterctl init --infrastructure aws
```

NOTE: you will need to enable the creation of the default Fargate IAM role. The easiest way is using `clusterawsadm` and using the `fargate` configuration option, for instructions see the [prerequisites](../using-clusterawsadm-to-fulfill-prerequisites.md).

### Amazon Linux 2023

Amazon EKS will end support for EKS optimized AL2 AMIs on November 26, 2025.

With AL2023, [nodeadm](https://github.com/awslabs/amazon-eks-ami/tree/main/nodeadm) is used to join EKS cluster.
Starting with v2.9.0, it's possible to set the node type in `EKSConfig` and `EKSConfigTemplate` like this:

```yaml
apiVersion: bootstrap.cluster.x-k8s.io/v1beta2
kind: EKSConfigTemplate
metadata:
  name: al2023
spec:
  template:
    spec:
      nodeType: al2023
```

AL2023 AMI can also be set in `AWSMAchineTemplate`.
The use of Secrets Manager trick should be disabled because
nodeadm expect the `NodeConfig` in plain text in EC2 instance's userdata.


```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
kind: AWSMachineTemplate
metadata:
  name: al2023
spec:
  template:
    spec:
      ami:
        eksLookupType: AmazonLinux2023
      cloudInit:
        insecureSkipSecretsManager: true
```
