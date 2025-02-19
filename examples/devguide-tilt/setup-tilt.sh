#!/bin/bash
# https://cluster-api-aws.sigs.k8s.io/development/development

# Check AWS auth
aws sts get-caller-identity

# Build clusterawsadm
make generate
make clusterawsadm

# ? change to $HOME/.local/bin ?
sudo ln -s $PWD/bin/clusterawsadm /usr/local/bin/clusterawsadm

clusterawsadm bootstrap iam create-cloudformation-stack

export AWS_B64ENCODED_CREDENTIALS=$(clusterawsadm bootstrap credentials encode-as-profile)
# echo  $AWS_B64ENCODED_CREDENTIALS | base64 -d

kind create cluster --name=capi-test


cd "$(go env GOPATH)"/src
mkdir -p sigs.k8s.io
cd sigs.k8s.io/
git clone git@github.com:kubernetes-sigs/cluster-api.git
cd cluster-api
# git fetch upstream

# TODO: create tilt-settings.json from template

# setup tilt
go install github.com/tilt-dev/ctlptl/cmd/ctlptl@latest
curl -fsSL https://raw.githubusercontent.com/tilt-dev/tilt/master/scripts/install.sh | bash

tilt up