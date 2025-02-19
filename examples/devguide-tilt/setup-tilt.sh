#!/bin/bash
# https://cluster-api-aws.sigs.k8s.io/development/development

# Check AWS auth
aws sts get-caller-identity

# Build clusterawsadm
make clusterawsadm

clusterawsadm bootstrap iam create-cloudformation-stack

export AWS_B64ENCODED_CREDENTIALS=$(clusterawsadm bootstrap credentials encode-as-profile)

kind create cluster --name=capi-test


cd "$(go env GOPATH)"/src
mkdir sigs.k8s.io
cd sigs.k8s.io/
git clone git@github.com:kubernetes-sigs/cluster-api.git
cd cluster-api
git fetch upstream