#!/usr/bin/env bash
# Copyright 2022 The Kubernetes Authors.
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

# This has been copied from https://github.com/kubernetes-sigs/cluster-api/blob/release-1.1/hack/utils.sh

# get_root_path returns the root path of the project source tree
get_root_path() {
    git rev-parse --show-toplevel
}

# cd_root_path cds to the root path of the project source tree
cd_root_path() {
    cd "$(get_root_path)" || exit
}

