# Specifying the IAM Role for ROSA HCP Management Components

When using an OpenShift or ROSA-HCP cluster as the management cluster, you can configure the CAPA controller to use IAM roles instead of storing AWS credentials. This uses OIDC federation to allow the CAPA controller service account to assume an IAM role.

## Prerequisites

- A management cluster (OpenShift or ROSA-HCP) with CAPI and CAPA installed.
  Follow the [Quick Start Guide](https://cluster-api-aws.sigs.k8s.io/quick-start) to install CAPI and CAPA using `clusterctl init --infrastructure aws`. For the initial installation, you can use temporary AWS credentials (e.g. via `aws sts get-session-token` or environment variables). Once the IAM role is configured below, the CAPA controller will use the role instead of stored credentials.

    **Note:** The ROSA and MachinePool feature gates must be enabled before running `clusterctl init`:
    ```shell
    export EXP_ROSA=true
    export EXP_MACHINE_POOL=true
    ```
- The management cluster must have an OIDC provider configured

## Retrieve the OIDC Provider

Extract the OIDC provider from the management cluster and set your AWS account ID:

```shell
export OIDC_PROVIDER=$(kubectl get authentication.config.openshift.io cluster -ojson | jq -r .spec.serviceAccountIssuer | sed 's/https:\/\///')
export AWS_ACCOUNT_ID=<your-aws-account-id>
```

## Create the Trust Policy

Create a trust policy that allows the `capa-controller-manager` service account to assume the IAM role:

```shell
cat <<EOF > trust.json
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
EOF
```


## Create the IAM Role

Create the IAM role and attach the required AWS policies:

```shell
aws iam create-role --role-name "capa-manager-role" \
  --assume-role-policy-document file://trust.json \
  --description "IAM role for CAPA to assume"

aws iam attach-role-policy --role-name capa-manager-role \
  --policy-arn arn:aws:iam::aws:policy/AWSCloudFormationFullAccess

aws iam attach-role-policy --role-name capa-manager-role \
  --policy-arn arn:aws:iam::aws:policy/AmazonVPCFullAccess
```

## Annotate the Service Account

Retrieve the IAM role ARN and annotate the CAPA controller service account:

```shell
export APP_IAM_ROLE_ARN=$(aws iam get-role --role-name=capa-manager-role --query Role.Arn --output text)

kubectl annotate serviceaccount -n capa-system capa-controller-manager \
  eks.amazonaws.com/role-arn=$APP_IAM_ROLE_ARN
```

Restart the CAPA controller to pick up the new role:

```shell
kubectl rollout restart deployment capa-controller-manager -n capa-system
```

After this configuration, the CAPA controller will use the IAM role to manage AWS resources, and you can provision ROSA HCP clusters without storing AWS credentials in the management cluster.
