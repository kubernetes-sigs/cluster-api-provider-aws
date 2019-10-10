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

export AWS_REGION_AZA=${AWS_REGION}a
export AWS_REGION_AZB=${AWS_REGION}a
export AWS_REGION_AZC=${AWS_REGION}a

REGIONMAPAZLIST="
eu-north-1:eu-north-1a eu-north-1b eu-north-1c,
ap-south-1:ap-south-1a ap-south-1b ap-south-1c,
eu-west-3:eu-west-3a eu-west-3b eu-west-3c,
eu-west-2:eu-west-2a eu-west-2b eu-west-2c,
eu-west-1:eu-west-1a eu-west-1b eu-west-1c,
ap-northeast-2:ap-northeast-2a ap-northeast-2b ap-northeast-2c,
ap-northeast-1:ap-northeast-1a ap-northeast-1c ap-northeast-1d,
sa-east-1:sa-east-1a sa-east-1b sa-east-1c,
ca-central-1:ca-central-1a ca-central-1b,
ap-east-1:ap-east-1a ap-east-1b ap-east-1c,
ap-southeast-1:ap-southeast-1a ap-southeast-1b ap-southeast-1c,
ap-southeast-2:ap-southeast-2a ap-southeast-2b ap-southeast-2c,
eu-central-1:eu-central-1a eu-central-1b eu-central-1c,
us-east-1:us-east-1a us-east-1b us-east-1c us-east-1d us-east-1e us-east-1f,
us-east-2:us-east-2a us-east-2b us-east-2c,
us-west-1:us-west-1b us-west-1c,
us-west-2:us-west-2a us-west-2b us-west-2c us-west-2d,
"

function set_region_az_info() {
  declare -a AZLIST='('$(echo $REGIONMAPAZLIST | tr ',' '\n' | grep $1 | cut -d":" -f2)')'
  if [ ${#AZLIST[@]} -lt 3 ];
  then
    echo "this region not three AZ"
    AWS_REGION_AZA=${AZLIST[0]}
    AWS_REGION_AZB=${AZLIST[0]}
    AWS_REGION_AZC=${AZLIST[0]}
  else
    AWS_REGION_AZA=${AZLIST[0]}
    AWS_REGION_AZB=${AZLIST[1]}
    AWS_REGION_AZC=${AZLIST[2]}
  fi
}