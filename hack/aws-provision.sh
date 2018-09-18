#!/usr/bin/env bash

# Your AWS user account
export TF_VAR_aws_user=$(aws iam get-user | jq --raw-output '.User.UserName')

export TF_VAR_cluster_domain="aos-cloud.eu"
export TF_VAR_cluster_name=$(whoami)
export TF_VAR_cluster_namespace="dev-${TF_VAR_cluster_name}"

export TF_IN_AUTOMATION="true"

cd ./prebuild
echo "***  starting terraform"

terraform init -input=false
if [ $? == "0" ]; then
  terraform plan -input=false -out=tfplan.out && terraform apply -input=false -auto-approve tfplan.out
fi
