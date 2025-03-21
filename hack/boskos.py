#!/usr/bin/env python3

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

import argparse
import json
import os

import requests
import time

BOSKOS_HOST = os.environ.get("BOSKOS_HOST", "boskos")
BOSKOS_RESOURCE_NAME = os.environ.get('BOSKOS_RESOURCE_NAME')
# Retry up to 3 times, with 10 seconds between tries. This is the same as defaults
# on https://github.com/kubernetes-sigs/boskos/
MAX_RETRIES = 3
RETRY_WAIT = 10


def checkout_account(resource_type, user, input_state="free", tries=1):
    url = f'http://{BOSKOS_HOST}/acquire?type={resource_type}&state={input_state}&dest=busy&owner={user}'

    r = requests.post(url)

    if r.status_code == 200:
        content = r.content.decode()
        result = json.loads(content)

        print(f"export BOSKOS_RESOURCE_NAME={result['name']}")
        print(f"export AWS_ACCESS_KEY_ID={result['userdata']['access-key-id']}")
        print(f"export AWS_SECRET_ACCESS_KEY={result['userdata']['secret-access-key']}")
    # The http API has two possible meanings of 404s - the named resource type cannot be found or there no available resources of the type.
    # For our purposes, we don't need to differentiate.
    elif r.status_code == 404:
        print(f"could not find available host, retrying in {RETRY_WAIT}s")
        if tries > MAX_RETRIES:
            raise Exception(f"could not allocate host after {MAX_RETRIES} tries")
        tries = tries + 1
        time.sleep(RETRY_WAIT)
        return checkout_account(resource_type, user, input_state, tries)
    else:
        raise Exception(f"Got invalid response {r.status_code}: {r.reason}")


def release_account(user):
    url = f'http://{BOSKOS_HOST}/release?name={BOSKOS_RESOURCE_NAME}&dest=dirty&owner={user}'

    r = requests.post(url)
    
    if r.status_code != 200:
        raise Exception(f"Got invalid response {r.status_code}: {r.reason}")


def send_heartbeat(user):
    url = f'http://{BOSKOS_HOST}/update?name={BOSKOS_RESOURCE_NAME}&state=busy&owner={user}'

    while True:
        print(f"POST-ing heartbeat for resource {BOSKOS_RESOURCE_NAME} to {BOSKOS_HOST}")
        r = requests.post(url)

        if r.status_code == 200:
            print(f"response status: {r.status_code}")
        else:
            print(f"Got invalid response {r.status_code}: {r.reason}")

        time.sleep(60)


def main():
    parser = argparse.ArgumentParser(description='Boskos AWS Account Management')

    parser.add_argument(
        '--get', dest='checkout_account', action="store_true",
        help='Checkout a Boskos AWS Account'
    )

    parser.add_argument(
        '--release', dest='release_account', action="store_true",
        help='Release a Boskos AWS Account'
    )

    parser.add_argument(
        '--heartbeat', dest='send_heartbeat', action="store_true",
        help='Send heartbeat for the checked out a Boskos AWS Account'
    )

    parser.add_argument(
        '--resource-type', dest="resource_type", type=str,
        default="aws-account",
        help="Type of Boskos resource to manage"
    )

    parser.add_argument(
        '--user', dest="user", type=str,
        default="cluster-api-provider-aws",
        help="username"
    )

    args = parser.parse_args()

    if args.checkout_account:
        checkout_account(args.resource_type, args.user)

    elif args.release_account:
        release_account(args.user)

    elif args.send_heartbeat:
        send_heartbeat(args.user)


if __name__ == "__main__":
    main()
