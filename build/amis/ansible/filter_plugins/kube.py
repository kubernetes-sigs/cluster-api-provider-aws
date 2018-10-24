# Copyright 2018 The Kubernetes Authors.

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

# http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import re

class FilterModule(object):

    def filters(self):
        return {
            'kube_platform_version': self.kube_platform_version,
        }

    def kube_platform_version(self, version, platform):
        if version == "latest":
            return version

        match = re.match('(\d+\.\d+.\d+)\-(\d+)', version)
        if not match:
            raise Exception("Version '%s' does not appear to be a "
                            "kubernetes version." % version)
        sub = match.groups(1)[1]
        if len(sub) == 1:
            if platform.lower() == "debian":
                return "%s-%s" % (match.groups(1)[0], '{:02d}'.format(sub))
            else:
                return version
        if len(sub) == 2:
            if platform.lower() == "redhat":
                return "%s-%s" % (match.groups(1)[0], int(sub))
            else:
                return version

        raise Exception("Could not parse kubernetes version")
