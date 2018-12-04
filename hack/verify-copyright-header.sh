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

year=$(date +'%Y')

GIT_BRANCH=$(git branch | grep \* | cut -d ' ' -f2)
newly_added_files=($(git diff origin/master ${GIT_BRANCH} --name-only --diff-filter=A))

if [ -n "$newly_added_files" ]
then
    for naf in $newly_added_files; do
        files_without_header+=$(grep -L "Copyright $year" $naf)
    done

    if [ -n "$files_without_header" ]
    then
        echo "Copyright $year license header not found in the following newly added files:"
        for f in "${files_without_header[@]}"; do
            echo "    - $f"
        done
        exit 1;
    fi
fi
exit 0;
