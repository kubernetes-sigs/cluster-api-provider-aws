#!/usr/bin/env bash
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

echo RELEASE_TAG $(git describe --tags)
echo BUILD_TIME $(date -u '+%Y-%m-%d_%I:%M:%S%p')
GIT_DIRTY=$(git status --porcelain 2> /dev/null)
if [[ -z "${GIT_DIRTY}" ]]; then
  GIT_TREE_STATE=clean
else
  GIT_TREE_STATE=dirty
fi
echo GIT_TREE_STATE $GIT_TREE_STATE
echo GIT_BRANCH $(git branch | grep \* | cut -d ' ' -f2)
echo GIT_COMMIT $(git rev-parse HEAD)
