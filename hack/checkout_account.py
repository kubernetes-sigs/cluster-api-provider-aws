#!/usr/bin/env python

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

"""Checks out a AWS account from E2E"""

import urllib
import httplib
import os
import sys
import json

BOSKOS_HOST=os.environ.get("BOSKOS_HOST", "boskos")

RESOURCE_TYPE = "aws-account"
USER = "cluster-api-provider-aws"

if __name__ == "__main__":
    conn = httplib.HTTPConnection(BOSKOS_HOST)

    conn.request("POST", "/acquire?%s" % urllib.urlencode({
        'type': RESOURCE_TYPE,
        'owner': USER,
        'state': 'free',
        'dest': 'busy',
    })

    )
    resp = conn.getresponse()
    if resp.status != 200:
        sys.exit("Got invalid response %d: %s" % (resp.status, resp.reason))

    body = json.load(resp)
    conn.close()
    print 'export BOSKOS_RESOURCE_NAME="%s";' % body['name']
    print 'export AWS_ACCESS_KEY_ID="%s";' % body['userdata']['access-key-id']
    print 'export AWS_SECRET_ACCESS_KEY="%s";' % body['userdata']['secret-access-key']

