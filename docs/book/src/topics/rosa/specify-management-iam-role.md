# Specifying the IAM Role for ROSA HCP Management Components

When using a management cluster (OCP or ROSA-HCP) created using AWS credentials with CAPI and CAPA installed, you can configure the CAPA controller to use IAM roles instead of storing AWS credentials. This uses OIDC federation to allow the CAPA controller service account to assume an IAM role.

## Prerequisites

- A management cluster (OCP or ROSA-HCP) created using AWS credentials with CAPI and CAPA installed
- The management cluster must have an OIDC provider configured

## Retrieve the OIDC Provider

1. Extract the OIDC provider from the management cluster and set your AWS account ID:

    ```shell
    export OIDC_PROVIDER=$(kubectl get authentication.config.openshift.io cluster -ojson | jq -r .spec.serviceAccountIssuer | sed 's/https:\/\///')
    export AWS_ACCOUNT_ID=<your-aws-account-id>
    ```

## Create the Trust Policy

1. Save the following trust policy to a file `trust.json`. This allows the `capa-controller-manager` service account to assume the IAM role:

    ```json
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": {
            "Federated": "arn:aws:iam::${AWS_ACCOUNT_ID}:oidc-provider/${OIDC_PROVIDER}"
          },
          "Action": "sts:AssumeRoleWithWebIdentity",
          "Condition": {
            "StringEquals": {
              "${OIDC_PROVIDER}:sub": "system:serviceaccount:capa-system:capa-controller-manager"
            }
          }
        }
      ]
    }
    ```

    **Note:** Replace `${AWS_ACCOUNT_ID}` and `${OIDC_PROVIDER}` with the values exported in the previous step.

## Create the IAM Role

1. Create the IAM role using the trust policy:

    ```shell
    aws iam create-role --role-name "capa-manager-role" \
      --assume-role-policy-document file://trust.json \
      --description "IAM role for CAPA to assume"
    ```

2. Attach the required AWS policies to the role:

    ```shell
    aws iam attach-role-policy --role-name capa-manager-role \
      --policy-arn arn:aws:iam::aws:policy/AWSCloudFormationFullAccess

    aws iam attach-role-policy --role-name capa-manager-role \
      --policy-arn arn:aws:iam::aws:policy/AmazonVPCFullAccess
    ```

## Annotate the Service Account

1. Retrieve the IAM role ARN:

    ```shell
    export APP_IAM_ROLE_ARN=$(aws iam get-role --role-name=capa-manager-role --query Role.Arn --output text)
    ```

2. Annotate the CAPA controller service account with the role ARN:

    ```shell
    kubectl annotate serviceaccount -n capa-system capa-controller-manager \
      eks.amazonaws.com/role-arn=$APP_IAM_ROLE_ARN
    ```

3. Restart the CAPA controller to pick up the new role:

    ```shell
    kubectl rollout restart deployment capa-controller-manager -n capa-system
    ```

After this configuration, the CAPA controller will use the IAM role to manage AWS resources, and you can provision ROSA HCP clusters without storing AWS credentials in the management cluster.
