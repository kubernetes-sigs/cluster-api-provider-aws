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

install bazel-bin/cmd/clusterctl/darwin_amd64_pure_stripped/clusterctl /usr/local/bin
if [ $? -eq 0 ]; then
    echo -e "\033[0;32m[$(date)]\033[0;0m Installed clusterctl"
else
    echo -e "\033[0;31m[$(date)]\033[0;0m failed to install clusterctl"
fi
install bazel-bin/cmd/clusterawsadm/darwin_amd64_pure_stripped/clusterawsadm /usr/local/bin
if [ $? -eq 0 ]; then
    echo -e "\033[0;32m[$(date)]\033[0;0m Installed clusterawsadm"
else
    echo -e "\033[0;31m[$(date)]\033[0;0m failed to install clusterawsadm"
fi


