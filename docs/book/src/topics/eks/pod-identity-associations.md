# EKS Pod Identity Associations

[EKS Pod Identities](https://docs.aws.amazon.com/eks/latest/eksctl/pod-identity-associations.html) allow pods in an EKS cluster to assume an IAM role by associating the IAM role to a Kubernetes Service Account.

## Pre-requisites

The **Amazon EKS Pod Identity Agent** addon must be installed in the cluster. See [EKS Addons](addons.md).

## Declaring associations

Add entries to `spec.podIdentityAssociations`:

```yaml
kind: AWSManagedControlPlane
apiVersion: controlplane.cluster.x-k8s.io/v1beta2
metadata:
  name: "capi-managed-test-control-plane"
spec:
  region: "eu-west-2"
  version: "v1.30.0"
  podIdentityAssociations:
    - serviceAccountNamespace: "kube-system"
      serviceAccountName: "aws-load-balancer-controller"
      roleARN: "arn:aws:iam::123456789012:role/AmazonEKSLoadBalancerControllerRole"
    - serviceAccountNamespace: "default"
      serviceAccountName: "my-app"
      roleARN: "arn:aws:iam::123456789012:role/AppAccessRole"
      targetRoleARN: "arn:aws:iam::987654321098:role/CrossAccountRole"
```

Fields:

| Field | Required | Description |
| --- | --- | --- |
| `serviceAccountNamespace` | yes | Namespace of the Service Account. Defaults to `default`. |
| `serviceAccountName` | yes | Name of the Service Account. |
| `roleARN` | yes | IAM role the Service Account assumes. The role's trust policy must allow the `pods.eks.amazonaws.com` service principal. |
| `targetRoleARN` | no | Role to chain to via `sts:AssumeRole` after assuming `roleARN`, e.g. for cross-account access. |
