# Building a dev environment using terraform
This directory holds the recipes required in order to deploy a working AWS dev environment using [terraform](https://www.terraform.io/downloads.html).

## Deployment Instructions
1. Download terraform from the link above and place it into your `$PATH`.
2. The following environment variables need to be set:
  1. `TF_VAR_aws_user` your AWS username (`aws iam get-user | jq --raw-output '.User.UserName'`)
  2. `TF_VAR_cluster_domain` The Route53 domain name to be used for the cluster.
  3. `TF_VAR_cluster_name` (self explanatory...)
  4. `TF_VAR_cluster_namespace`: the namespace to deploy the cluster components to (suggestion: `dev-${TF_VAR_cluster_name}`)

4. Run terraform:
```
> terraform init
> terraform plan
> terraform apply
```

To destroy the environment, simply run `terraform destroy`.
