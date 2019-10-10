#!/usr/bin/env bash

# Copyright 2019 The Kubernetes Authors.
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

type aws || exit

echo "fetching region list and AZ list"
for REGION in $(aws ec2 describe-regions --query "Regions[].{Name:RegionName}" --output text | paste -sd " " -);
do
  AZLIST=$(aws ec2 describe-availability-zones --region $REGION --query "AvailabilityZones[].{Name:ZoneName}" --output text | paste -sd " " -)
  echo "$REGION:$AZLIST,"
done