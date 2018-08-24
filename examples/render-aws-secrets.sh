#!/bin/bash

set -e

if [ $# -lt 1 ]; then
    echo "usage: $0 <filename>"
    exit 1
fi

if [ -z "$AWS_ACCESS_KEY_ID" ]; then
    echo "error: AWS_ACCESS_KEY_ID is not set in the environment" 2>&1
    exit 1
fi

if [ -z "$AWS_SECRET_ACCESS_KEY" ]; then
    echo "error: AWS_SECRET_ACCESS_KEY is not set in the environment" 2>&1
    exit 1
fi

x=$(echo -n "$AWS_ACCESS_KEY_ID" | base64)
y=$(echo -n "$AWS_SECRET_ACCESS_KEY" | base64)

sed -e "s/awsAccessKeyId:.*/awsAccessKeyId: $x/" \
    -e "s/awsSecretAccessKey:.*/awsSecretAccessKey: $y/" \
    "$1"
