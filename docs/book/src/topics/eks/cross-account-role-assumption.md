# Cross-account role assumption with IRSA

When the management cluster runs on EKS, CAPA can use IAM Roles for Service
Accounts (IRSA) for its credentials and assume a role in another AWS
account. This avoids storing long-lived AWS credentials in the management
cluster.

The credential chain in this example is:

```text
CAPA service account
  -> sts:AssumeRoleWithWebIdentity
management account CAPA controller role
  -> sts:AssumeRole
workload account CAPA controller role
```

The management account hosts the EKS management cluster. The workload account
is where CAPA creates the EKS workload cluster.

## Prerequisites

- An EKS management cluster with an [IAM OIDC
  provider](https://docs.aws.amazon.com/eks/latest/userguide/enable-iam-roles-for-service-accounts.html)
- `clusterawsadm` and `clusterctl`
- Permission to create IAM roles and policies in both AWS accounts

Set the values used in the examples:

```bash
export AWS_REGION=us-east-1
export MANAGEMENT_ACCOUNT_ID=111122223333
export WORKLOAD_ACCOUNT_ID=444455556666
export OIDC_PROVIDER=oidc.eks.us-east-1.amazonaws.com/id/EXAMPLE
export MANAGEMENT_ROLE=controllers.cluster-api-provider-aws.sigs.k8s.io
export WORKLOAD_ROLE=controllers.cluster-api-provider-aws.sigs.k8s.io
```

`OIDC_PROVIDER` is the management cluster's issuer URL without the `https://`
prefix.

## Configure the management account role

Create the CAPA controller role in the management account with:

- a trust policy that permits the CAPA service accounts to call
  `sts:AssumeRoleWithWebIdentity`; and
- permission to call `sts:AssumeRole` on the CAPA controller role in the
  workload account.

The relevant trust policy statement is:

```json
{
  "Effect": "Allow",
  "Principal": {
    "Federated": "arn:aws:iam::111122223333:oidc-provider/oidc.eks.us-east-1.amazonaws.com/id/EXAMPLE"
  },
  "Action": "sts:AssumeRoleWithWebIdentity",
  "Condition": {
    "StringEquals": {
      "oidc.eks.us-east-1.amazonaws.com/id/EXAMPLE:aud": "sts.amazonaws.com",
      "oidc.eks.us-east-1.amazonaws.com/id/EXAMPLE:sub": [
        "system:serviceaccount:capa-system:capa-controller-manager",
        "system:serviceaccount:capa-eks-control-plane-system:capa-eks-control-plane-controller-manager"
      ]
    }
  }
}
```

Restricting both `aud` and `sub` ensures that only the expected Kubernetes
service accounts can use the role. If CAPA is installed in different
namespaces, update the subjects accordingly.

Attach the normal CAPA controller permissions to this role. Also attach the
following policy for the second step in the credential chain:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "sts:AssumeRole",
      "Resource": "arn:aws:iam::444455556666:role/controllers.cluster-api-provider-aws.sigs.k8s.io"
    }
  ]
}
```

The controller permissions and the additional `sts:AssumeRole` permission can
also be generated with `clusterawsadm`. Set `allowAssumeRole: true` in the
management account's `AWSIAMConfiguration`; see [Cross Account Role
Assumption](../using-clusterawsadm-to-fulfill-prerequisites.md#cross-account-role-assumption).

## Configure the workload account role

In the workload account, create a role with the CAPA controller permissions.
Its trust policy must allow the management account role to assume it:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::111122223333:role/controllers.cluster-api-provider-aws.sigs.k8s.io"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
```

The workload account role can also be created with `clusterawsadm` by adding
this statement to `clusterAPIControllers.trustStatements` in its
`AWSIAMConfiguration`.

## Install CAPA with IRSA

Set the management account role before installing CAPA. `clusterctl` uses this
value to annotate the CAPA controller service accounts with
`eks.amazonaws.com/role-arn`:

```bash
export AWS_CONTROLLER_IAM_ROLE="arn:aws:iam::${MANAGEMENT_ACCOUNT_ID}:role/${MANAGEMENT_ROLE}"
clusterctl init --infrastructure aws
```

Do not provide static bootstrap credentials to the CAPA controllers. If the
provider was previously installed with static credentials, remove them by
following [Using IAM roles in the management cluster](../using-iam-roles-in-mgmt-cluster.md).

Verify the annotation after installation:

```bash
kubectl -n capa-system get serviceaccount capa-controller-manager \
  -o jsonpath='{.metadata.annotations.eks\.amazonaws\.com/role-arn}{"\n"}'
```

## Create the cross-account identity

Create an `AWSClusterRoleIdentity` that chains from the controller's IRSA
credentials to the role in the workload account:

```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
kind: AWSClusterControllerIdentity
metadata:
  name: default
spec:
  allowedNamespaces:
    list:
      - default
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
kind: AWSClusterRoleIdentity
metadata:
  name: workload-account
spec:
  allowedNamespaces:
    list:
      - default
  roleARN: arn:aws:iam::444455556666:role/controllers.cluster-api-provider-aws.sigs.k8s.io
  sourceIdentityRef:
    kind: AWSClusterControllerIdentity
    name: default
```

If the `default` `AWSClusterControllerIdentity` already exists, do not include
it again. Confirm that its `allowedNamespaces` permits the namespaces that will
use it.

Set the identity on the `AWSManagedControlPlane` in the generated EKS cluster
manifest. `AWSManagedCluster` does not have an `identityRef`; it obtains the
resulting infrastructure information from the managed control plane:

```yaml
apiVersion: controlplane.cluster.x-k8s.io/v1beta2
kind: AWSManagedControlPlane
metadata:
  name: workload-cluster-control-plane
spec:
  identityRef:
    kind: AWSClusterRoleIdentity
    name: workload-account
  # Other control plane settings are omitted.
```

Apply the identity resources before applying the workload cluster manifest.
CAPA will use IRSA for its initial credentials and `AWSClusterRoleIdentity` for
the cross-account role assumption.

For identity scoping, external IDs, and nested role chains, see
[Multi-tenancy](../multitenancy.md).
