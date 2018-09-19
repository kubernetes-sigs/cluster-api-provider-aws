#!/usr/bin/env bash

# Your AWS user account
export TF_VAR_aws_user=$(aws iam get-user | jq --raw-output '.User.UserName')

export TF_VAR_cluster_domain="aos-cloud.eu"
export TF_VAR_environment_id=${ENVIRONMENT_ID:-$(uuidgen)}
export TF_VAR_cluster_namespace="dev-${TF_VAR_cluster_name}"

export TF_IN_AUTOMATION="true"

cd ./prebuild
case ${1} in
 "install")
   echo "***  starting terraform with TF_VAR_environment_id=${TF_VAR_environment_id}"

   terraform init -input=false
   if [ $? == "0" ]; then
     terraform plan -input=false -out=tfplan.out && terraform apply -input=false -auto-approve tfplan.out
   fi
   ;;
 "destroy")
   terraform destroy -input=false -auto-approve
   ;;
 *)
   echo "Use $0 install or $0 destroy. Thanks!"
   ;;
esac
