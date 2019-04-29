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

GIT_COMMIT="$(git describe --always --dirty --abbrev=14)"

if git_status=$(git status --porcelain 2>/dev/null) && [[ -z ${git_status} ]]; then
    GIT_TREE_STATE="clean"
else
    GIT_TREE_STATE="dirty"
fi

# stolen from k8s.io/hack/lib/version.sh
# Use git describe to find the version based on tags.
if GIT_VERSION=$(git describe --tags --abbrev=14 2>/dev/null); then
    # This translates the "git describe" to an actual semver.org
    # compatible semantic version that looks something like this:
    #   v1.1.0-alpha.0.6+84c76d1142ea4d
    #
    # TODO: We continue calling this "git version" because so many
    # downstream consumers are expecting it there.
    DASHES_IN_VERSION=$(echo "${GIT_VERSION}" | sed "s/[^-]//g")
    if [[ "${DASHES_IN_VERSION}" == "---" ]] ; then
        # We have distance to subversion (v1.1.0-subversion-1-gCommitHash)
        GIT_VERSION=$(echo "${GIT_VERSION}" | sed "s/-\([0-9]\{1,\}\)-g\([0-9a-f]\{14\}\)$/.\1\-\2/")
    elif [[ "${DASHES_IN_VERSION}" == "--" ]] ; then
        # We have distance to base tag (v1.1.0-1-gCommitHash)
        GIT_VERSION=$(echo "${GIT_VERSION}" | sed "s/-g\([0-9a-f]\{14\}\)$/-\1/")
    fi
    if [[ "${GIT_TREE_STATE}" == "dirty" ]]; then
        # git describe --dirty only considers changes to existing files, but
        # that is problematic since new untracked .go files affect the build,
        # so use our idea of "dirty" from git status instead.
        GIT_VERSION+="-dirty"
    fi


    # Try to match the "git describe" output to a regex to try to extract
    # the "major" and "minor" versions and whether this is the exact tagged
    # version or whether the tree is between two tagged versions.
    if [[ "${GIT_VERSION}" =~ ^v([0-9]+)\.([0-9]+)(\.[0-9]+)?([-].*)?([+].*)?$ ]]; then
        GIT_MAJOR=${BASH_REMATCH[1]}
        GIT_MINOR=${BASH_REMATCH[2]}
    fi

    # If GIT_VERSION is not a valid Semantic Version, then refuse to build.
    if ! [[ "${GIT_VERSION}" =~ ^v([0-9]+)\.([0-9]+)(\.[0-9]+)?(-[0-9A-Za-z.-]+)?(\+[0-9A-Za-z.-]+)?$ ]]; then
        echo "GIT_VERSION should be a valid Semantic Version. Current value: ${GIT_VERSION}"
        echo "Please see more details here: https://semver.org"
        exit 1
    fi
fi

GIT_BRANCH=$(git branch | grep \* | cut -d ' ' -f2)
GIT_RELEASE_TAG=$(git describe --abbrev=0 --tags)
GIT_RELEASE_COMMIT=$(git rev-list -n 1  ${GIT_RELEASE_TAG} | head -c 14)

cat <<EOF
GIT_COMMIT ${GIT_COMMIT-}
GIT_TREE_STATE ${GIT_TREE_STATE-}
GIT_MAJOR ${GIT_MAJOR-}
GIT_MINOR ${GIT_MINOR-}
GIT_VERSION ${GIT_VERSION-}
GIT_BRANCH ${GIT_BRANCH-}
GIT_RELEASE_TAG ${GIT_RELEASE_TAG-}
GIT_RELEASE_COMMIT ${GIT_RELEASE_COMMIT-}
STABLE_AWS_ACCESS_KEY_ID ${AWS_ACCESS_KEY_ID-}
STABLE_AWS_SECRET_ACCESS_KEY ${AWS_SECRET_ACCESS_KEY-}
STABLE_AWS_DEFAULT_REGION ${AWS_DEFAULT_REGION:-${AWS_REGION-}}
STABLE_AWS_REGION ${AWS_DEFAULT_REGION:-${AWS_REGION-}}
STABLE_AWS_SESSION_TOKEN ${AWS_SESSION_TOKEN-}
EOF
