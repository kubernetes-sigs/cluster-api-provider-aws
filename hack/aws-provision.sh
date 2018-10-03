#!/usr/bin/env bash

script_dir="$(cd $(dirname "${BASH_SOURCE[0]}") && pwd -P)"

# Your AWS user account
export TF_VAR_aws_user=$(aws iam get-user | jq --raw-output '.User.UserName')

export TF_VAR_cluster_domain="aos-cloud.eu"
export TF_VAR_environment_id=${ENVIRONMENT_ID:-$(uuidgen)}
export TF_VAR_cluster_namespace="dev-${TF_VAR_cluster_name}"

export TF_IN_AUTOMATION="true"

pushd $script_dir/prebuild

case ${1} in
 "install")
   echo "***  starting terraform with TF_VAR_environment_id=${TF_VAR_environment_id}"

   terraform init -input=false
   if [ $? == "0" ]; then
     terraform plan -input=false -out=tfplan.out && terraform apply -input=false -auto-approve tfplan.out
   fi
   ;;
 "destroy")
   # Terminate all running instances in the vpc first
   VPC_ID=$(terraform output vpc_id || echo "")
   if [[ "${VPC_ID}" != "" ]]; then
     instances=$(aws ec2 describe-instances --filters Name=vpc-id,Values=${VPC_ID} | jq '.Reservations[].Instances[].InstanceId' --raw-output || echo "")
     if [[ "${instances}" != "" ]]; then
       aws ec2 terminate-instances --instance-ids ${instances} || true
     fi
   fi
   terraform destroy -input=false -auto-approve
   ;;
 *)
   echo "Use $0 install or $0 destroy. Thanks!"
   ;;
esac
