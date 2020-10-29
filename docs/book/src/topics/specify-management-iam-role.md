# Specifying the IAM Role to use for Management Components

## Prerequisites

To be able to specify the IAM role that the management components should run as your cluster must be set up with [IAM Roles for Service Accounts](https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts.html).

## Setting IAM Role

Set the `AWS_CONTROLLER_IAM_ROLE` environment variable to the ARN of the IAM role to use when performing the `clustercrl init` command.

For example:

```bash
export AWS_CONTROLLER_IAM_ROLE=arn:aws:iam::1234567890:role/capa-management-components

clusterctl init --infrastructure=aws
```

## IAM Role Trust Policy

When creating the IAM role the following trust policy will need to be used with the `AWS_ACCOUNT_ID`, `AWS_REGION` and `OIDC_PROVIDER_ID` environment variables replaced.

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Federated": "arn:aws:iam::${AWS_ACCOUNT_ID}:oidc-provider/oidc.eks.${AWS_REGION}.amazonaws.com/id/${OIDC_PROVIDER_ID}"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "ForAnyValue:StringEquals": {
          "oidc.eks.${AWS_REGION}.amazonaws.com/id/${OIDC_PROVIDER_ID}:sub": [
            "system:serviceaccount:capa-system:capa-controller-manager",
            "system:serviceaccount:capi-system:capi-controller-manager",
            "system:serviceaccount:capa-eks-control-plane-system:capa-eks-control-plane-controller-manager",
            "system:serviceaccount:capa-eks-bootstrap-system:capa-eks-bootstrap-controller-manager",
          ]
        }
      }
    }
  ]
}
```
