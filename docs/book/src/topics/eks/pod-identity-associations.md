# EKS Pod Identity Associations

[EKS Pod Identity Associations](https://aws.amazon.com/blogs/containers/introducing-amazon-eks-add-ons/) can be used with EKS clusters created using Cluster API Provider AWS.

## Prerequisites

### Setting up the IAM role in AWS

Outside of CAPI/CAPA, you must first create an IAM Role which allows the `pods.eks.amazonaws.com` service principal in the trust policy. EKS Identities trust relationships must also include the `sts:TagSession` permission (on top of the `sts:AssumeRole` permission).

This is a sample trust policy which allows a kubernetes service account to assume this role. We'll call the role `capi-test-role` in the next steps.

```yaml
{
  "Version": "2012-10-17",
  "Statement":
    [
      {
        "Effect": "Allow",
        "Principal": { "Service": "pods.eks.amazonaws.com" },
        "Action": ["sts:AssumeRole", "sts:TagSession"],
      },
    ],
}
```

### Installing the EKS Pod Identity Agent

The EKS Pod Identity Agent can be installed as a Managed Add-on through the AWS Console, or through CAPA.
To install the addon through CAPA, add it to `AWSManagedControlPlane`. Please ensure that the version is up to date, according to the [addons section](addons.md).

```yaml
# [...]
kind: AWSManagedControlPlane
spec:
  # [...]
  addons:
    # [...]
    - conflictResolution: overwrite
      name: eks-pod-identity-agent
      version: v1.1.0-eksbuild.1
```

You can verify that this is running on your kubernetes cluster with `kubectl get deploy -A | grep eks`

## Mapping a service account to an IAM role

Now that you have created a role `capi-test-role` in AWS, and have added the EKS agent to your cluster, we must add the following to our `AWSManagedControlPlane` under `.spec.podIdentityAssociations`

```yaml
# [...]
kind: AWSManagedControlPlane
spec:
  # [...]
  podIdentityAssociations:
    - serviceAccount:
        namespace: default
        name: myserviceaccount
        roleARN: arn:aws:iam::012345678901:role/capi-test-role
```

- `serviceAccount.namespace` and `serviceAccount.name` refer to the [`ServiceAccount`](https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/) object in the Kubernetes cluster
- `serviceAccount.roleARN` is the AWS ARN for the IAM role you created in step 1 (named `capi-test-role` in this tutorial). Make sure to copy this exactly from your AWS console (`IAM > Roles`).

To use the same IAM role across multiple service accounts/namespaces, you must create multiple associations.

A full CAPA example of everything mentioned above, including 2 role mappings, is shown below:

```yaml
kind: AWSManagedControlPlane
apiVersion: controlplane.cluster.x-k8s.io/v1beta1
metadata:
  name: "capi-managed-test"
spec:
  region: "eu-west-2"
  sshKeyName: "capi-management"
  version: "v1.27.0"
  addons:
    - conflictResolution: overwrite
      name: eks-pod-identity-agent
      version: v1.1.0-eksbuild.1
  podIdentityAssociations:
    - serviceAccountNamespace: default
      serviceAccountName: myserviceaccount
      serviceAccountRoleARN: arn:aws:iam::012345678901:role/capi-test-role
    - serviceAccountNamespace: another-namespace
      serviceAccountName: another-service-account
      serviceAccountRoleARN: arn:aws:iam::012345678901:role/capi-test-role
```
