<!-- NB: This page is meant to be embedded in Cluster API book -->

## Requirements

- Linux or MacOS (Windows isn't supported at the moment).
- AWS credentials.
- [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/installing.html)
- [jq](https://stedolan.github.io/jq/download/)

### With `clusterawsadm`

Get the latest [clusterawsadm](https://github.com/kubernetes-sigs/cluster-api-provider-aws/releases)
and place it in your path.

Cluster API Provider AWS ships with clusterawsadm, a utility to help you manage
IAM objects for this project.

In order to use clusterawsadm you must have an administrative user in an AWS account.
Once you have that administrator user you need to set your environment variables:

* `AWS_REGION`
* `AWS_ACCESS_KEY_ID`
* `AWS_SECRET_ACCESS_KEY`
* `AWS_SESSION_TOKEN` (if you are using Multi-factor authentication)

After these are set run this command to get you up and running:

```bash
clusterawsadm bootstrap iam create-cloudformation-stack
```

You can optionally specify additional AWS policies (they can be user or AWS managed, but must already exists) in the call to `create-formation-stack` (and `generate-cloudformation`) if required, e.g.

```
clusterawsadm alpha bootstrap create-stack \
  --extra-controlplane-policies arn:aws:iam::<AWS_ACCOUNT>:policy/my-policy,arn:aws:iam::aws:policy/AmazonEC2FullAccess \
  --extra-node-policies arn:aws:iam::<AWS_ACCOUNT>:policy/my-other-policy
```

These will be added to the control plane and node roles respectively when they are created.

### Without `clusterawsadm`

This is not a recommended route as the policies are very specific and will
change with new features.

If you do not wish to use the `clusteradwsadm` tool then you will need to
understand exactly which IAM policies and groups we are expecting. There are
several policies, roles and users that need to be created. Please see our
[controller policy][controllerpolicy] file to understand the permissions that are necessary.

[controllerpolicy]: https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/0e543e0eb30a7065c967f5df8d6abd872aa4ff0c/pkg/cloud/aws/services/cloudformation/bootstrap.go#L149-L188

## SSH Key pair

You will need to specify the name of an existing SSH key pair within the region
you plan on using. If you don't have one yet, a new one needs to be created.

### Create a new key pair

```bash
# Save the output to a secure location
aws ec2 create-key-pair --key-name default | jq .KeyMaterial -r
-----BEGIN RSA PRIVATE KEY-----
[... contents omitted ...]
-----END RSA PRIVATE KEY-----
```

If you want to save the private key directly into AWS Systems Manager Parameter
Store with KMS encryption for security, you can use the following command:

```bash
aws ssm put-parameter --name "/sigs.k8s.io/cluster-api-provider-aws/ssh-key" \
  --type SecureString \
  --value "$(aws ec2 create-key-pair --key-name default | jq .KeyMaterial -r)"
```

### Adding an existing public key to AWS

```bash
# Replace with your own public key
aws ec2 import-key-pair \
  --key-name default \
  --public-key-material "$(base64 -i ~/.ssh/id_rsa.pub)"
```

> **NB**: Only RSA keys are supported by AWS.

## Setting up the environment

The current iteration of the Cluster API Provider AWS relies on credentials
being present in your environment. These then get written into the cluster
manifests for use by the controllers.

If you used `clusterawsadm` to set up IAM resources for you then you can run
these commands to prepare your environment.

Your `AWS_REGION` must already be set.

```bash
export AWS_CREDENTIALS=$(aws iam create-access-key \
  --user-name bootstrapper.cluster-api-provider-aws.sigs.k8s.io)
export AWS_ACCESS_KEY_ID=$(echo $AWS_CREDENTIALS | jq .AccessKey.AccessKeyId -r)
export AWS_SECRET_ACCESS_KEY=$(echo $AWS_CREDENTIALS | jq .AccessKey.SecretAccessKey -r)
```

If you did not use `clusterawsadm` to provision your user, you will need to set
these environment variables in your own way.

> To save credentials securely in your environment, [aws-vault](https://github.com/99designs/aws-vault) uses
> the OS keystore as permanent storage, and offers shell features to securely
> expose and setup local AWS environments.
