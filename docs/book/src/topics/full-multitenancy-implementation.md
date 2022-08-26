# Multitenancy setup with EKS and Service Account

See [multitenancy](./multitenancy.md) for more
details on enabling the functionality and the various options you can use.

In this example, we are going to see how to create the following architecture with cluster API:

```
                                  AWS Account 1
                                 +--------------------+
                                 |                    |
                 +---------------+->EKS - (Managed)   |
                 |               |                    |
                 |               +--------------------+
 AWS Account 0   |                AWS Account 2
+----------------+---+           +--------------------+
|                |   |           |                    |
|  EKS - (Manager)---+-----------+->EKS - (Managed)   |
|                |   |           |                    |
+----------------+---+           +--------------------+
                 |                AWS Account 3
                 |               +--------------------+
                 |               |                    |
                 +---------------+->EKS - (Managed)   |
                                 |                    |
                                 +--------------------+
```

And specifically, we will only include:

- AWS Account 0 (aka Manager account used by management cluster where cluster API controllers reside)
- AWS Account 1 (aka Managed account used for EKS-managed workload clusters)

## Prerequisites

- A bootstrap cluster (kind)
- AWS CLI installed
- 2 (or more) AWS accounts
- [clusterawsadm](https://github.com/kubernetes-sigs/cluster-api-provider-aws/releases)
- [clusterctl](https://github.com/kubernetes-sigs/cluster-api/releases)

## Set variables

**Note:** the credentials below are the ones of the manager account

Export the following environment variables:

- AWS_REGION
- AWS_ACCESS_KEY_ID
- AWS_SECRET_ACCESS_KEY
- AWS_SESSION_TOKEN (if you are using Multi-factor authentication)
- AWS_MANAGER_ACCOUNT_ID
- AWS_MANAGED_ACCOUNT_ID
- OIDC_PROVIDER_ID="WeWillReplaceThisLater"

## Prepare the manager account

As explained in the [EKS prerequisites page](./eks/prerequisites.md), we need a couple of roles in the account to build the cluster, `clusterawsadm` CLI can take care of it.

We know that the CAPA provider in the Manager account should be able to assume roles in the Managed account (AWS Account 1).

We can create a clusterawsadm configuration that adds an inline policy to the `controllers.cluster-api-provider-aws.sigs.k8s.io` role.

```bash
envsubst > bootstrap-manager-account.yaml << EOL
apiVersion: bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1
kind: AWSIAMConfiguration
spec:
  eks: # This section should be changed accordingly to your requirements
    iamRoleCreation: false
    managedMachinePool:
      disable: true
    fargate:
      disable: false
  clusterAPIControllers: # This is the section that really matter
    disabled: false
    extraStatements:
    - Action:
      - "sts:AssumeRole"
      Effect: "Allow"
      Resource: ["arn:aws:iam::${AWS_MANAGED_ACCOUNT_ID}:role/controllers.cluster-api-provider-aws.sigs.k8s.io"]
    trustStatements:
    - Action:
      - "sts:AssumeRoleWithWebIdentity"
      Effect: "Allow"
      Principal:
        Federated:
        - "arn:aws:iam::${AWS_MANAGER_ACCOUNT_ID}:oidc-provider/oidc.eks.${AWS_REGION}.amazonaws.com/id/${OIDC_PROVIDER_ID}"
      Condition:
        "ForAnyValue:StringEquals":
          "oidc.eks.${AWS_REGION}.amazonaws.com/id/${OIDC_PROVIDER_ID}:sub":
            - system:serviceaccount:capa-system:capa-controller-manager
            - system:serviceaccount:capa-eks-control-plane-system:capa-eks-control-plane-controller-manager # Include if also using EKS
EOL
```

Let's provision the Manager role with:

```
clusterawsadm bootstrap iam create-cloudformation-stack --config bootstrap-manager-account.yaml
```

## Manager cluster

The following commands assume you have the AWS credentials for the Manager account exposed, and your kube context is pointing to the bootstrap cluster.

### Install cluster API provider in the bootstrap cluster

```bash
export AWS_B64ENCODED_CREDENTIALS=$(clusterawsadm bootstrap credentials encode-as-profile)
export EKS=true
export EXP_MACHINE_POOL=true
clusterctl init --infrastructure aws --target-namespace capi-providers
```

### Generate the cluster configuration

**NOTE:** You might want to update the Kubernetes and VPC addon versions to one of the available versions when running this command.

- [Kubernetes versions](https://docs.aws.amazon.com/eks/latest/userguide/kubernetes-versions.html)
- [VPC CNI add-on versions](https://docs.aws.amazon.com/eks/latest/userguide/managing-vpc-cni.html) don't forget to add the `v` prefix

```bash
export AWS_SSH_KEY_NAME=default
export VPC_ADDON_VERSION="v1.10.2-eksbuild.1"
clusterctl generate cluster manager --flavor eks-managedmachinepool-vpccni --kubernetes-version v1.20.2 --worker-machine-count=3 > manager-cluster.yaml
```

### Apply the cluster configuration

```bash
kubectl apply -f manager-cluster.yaml
```

**WAIT**: time to have a drink, the cluster is creating and we will have to wait for it to be there before continuing.

### IAM OIDC Identity provider

Follow AWS documentation to create an OIDC provider https://docs.aws.amazon.com/eks/latest/userguide/enable-iam-roles-for-service-accounts.html

### Update the TrustStatement above

```bash
export OIDC_PROVIDER_ID=<OIDC_ID_OF_THE_CLUSTER>
```

run the [Prepare the manager account](./full-multitenancy-implementation.md#prepare-the-manager-aws-account-0-account) step again

### Get manager cluster credentials

```bash
kubectl --namespace=default get secret manager-user-kubeconfig \
   -o jsonpath={.data.value} | base64 --decode \
   > manager.kubeconfig
```

### Install the CAPA provider in the manager cluster

Here we install the Cluster API providers into the manager cluster and create a service account to use the `controllers.cluster-api-provider-aws.sigs.k8s.io` role for the Management Components.

```bash
export AWS_B64ENCODED_CREDENTIALS=$(clusterawsadm bootstrap credentials encode-as-profile)
export EKS=true
export EXP_MACHINE_POOL=true
export AWS_CONTROLLER_IAM_ROLE=arn:aws:iam::${AWS_MANAGER_ACCOUNT_ID}:role/controllers.cluster-api-provider-aws.sigs.k8s.io
clusterctl init --kubeconfig manager.kubeconfig --infrastructure aws --target-namespace capi-providers
```

## Managed cluster

Time to build the managed cluster for pivoting the bootstrap cluster.

### Generate the cluster configuration

**NOTE:** As for the manager cluster you might want to update the Kubernetes and VPC addon versions.

```bash
export AWS_SSH_KEY_NAME=default
export VPC_ADDON_VERSION="v1.10.2-eksbuild.1"
clusterctl generate cluster manager --flavor eks-managedmachinepool-vpccni --kubernetes-version v1.20.2 --worker-machine-count=3 > managed-cluster.yaml
```

Edit the file and add the following to the `AWSManagedControlPlane` resource spec to point the controller to the manager account when creating the cluster.

```yaml
identityRef:
  kind: AWSClusterRoleIdentity
  name: managed-account
```

### Create the identities

```bash
envsubst > cluster-role-identity.yaml << EOL
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: AWSClusterRoleIdentity
metadata:
  name: managed-account
spec:
  allowedNamespaces: {} # This is unsafe since every namespace is allowed to use the role identity
  roleARN: arn:aws:iam::${AWS_MANAGED_ACCOUNT_ID}:role/controllers.cluster-api-provider-aws.sigs.k8s.io
  sourceidentityRef:
    kind: AWSClusterControllerIdentity
    name: default
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: AWSClusterControllerIdentity
metadata:
  name: default
spec:
  allowedNamespaces:{}
EOL
```

### Prepare the managed account

**NOTE:** Expose the **managed** account credentials before running the following commands.

This configuration is adding the trustStatement in the cluster api controller role to allow the `controllers.cluster-api-provider-aws.sigs.k8s.io` in the manager account to assume it.

```bash
envsubst > bootstrap-managed-account.yaml << EOL
apiVersion: bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1
kind: AWSIAMConfiguration
spec:
  eks:
    iamRoleCreation: false # Set to true if you plan to use the EKSEnableIAM feature flag to enable automatic creation of IAM roles
    managedMachinePool:
      disable: true # Set to false to enable creation of the default node role for managed machine pools
    fargate:
      disable: false # Set to false to enable creation of the default role for the fargate profiles
  clusterAPIControllers:
    disabled: false
    trustStatements:
    - Action:
      - "sts:AssumeRole"
      Effect: "Allow"
      Principal:
        AWS:
        - "arn:aws:iam::${AWS_MANAGER_ACCOUNT_ID}:role/controllers.cluster-api-provider-aws.sigs.k8s.io"
EOL
```

Let's provision the Managed account with:

```bash
clusterawsadm bootstrap iam create-cloudformation-stack --config bootstrap-managed-account.yaml
```

### Apply the cluster configuration

**Note:** Back to the **manager** account credentials

```
kubectl --kubeconfig manager.kubeconfig apply -f cluster-role-identity.yaml
kubectl --kubeconfig manager.kubeconfig apply -f managed-cluster.yaml
```

Time for another drink, enjoy your multi-tenancy setup.
