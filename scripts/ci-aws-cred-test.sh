#!/bin/bash

# Copyright 2018 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail


GROUP_NAME=capa-cred-test-group-$(openssl rand -hex 8)
USERNAME=capa-cred-test-user-$(openssl rand -hex 8)

echo "Group: ${GROUP_NAME}\nUser: ${USERNAME}"

echo "Creating IAM group"
aws iam create-group --group-name "${GROUP_NAME}"

echo "Attaching Policies"
aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/AmazonEC2FullAccess --group-name "${GROUP_NAME}"
aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/IAMFullAccess --group-name "${GROUP_NAME}"
aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/AmazonVPCFullAccess --group-name "${GROUP_NAME}"

echo "Creating User"
aws iam create-user --user-name "${USERNAME}"

echo "Adding User to Group"
aws iam add-user-to-group --user-name "${USERNAME}" --group-name "${GROUP_NAME}"

echo '### Clean up ###\n'
echo '\n'

echo 'Remove user from group'
aws iam remove-user-from-group --user-name "${USERNAME}" --group-name "${GROUP_NAME}"

echo 'Delete user'
aws iam delete-user --user-name "${USERNAME}"

echo 'Detach policies'
aws iam detach-group-policy --policy-arn arn:aws:iam::aws:policy/AmazonEC2FullAccess --group-name "${GROUP_NAME}"
aws iam detach-group-policy --policy-arn arn:aws:iam::aws:policy/IAMFullAccess --group-name "${GROUP_NAME}"
aws iam detach-group-policy --policy-arn arn:aws:iam::aws:policy/AmazonVPCFullAccess --group-name "${GROUP_NAME}"

echo 'Delete group'
aws iam delete-group --group-name "${GROUP_NAME}"

echo 'All tasks done'