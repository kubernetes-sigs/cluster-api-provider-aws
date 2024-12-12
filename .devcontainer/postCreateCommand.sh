#!/bin/bash
set -x

echo "Verifying versions"
go version
docker --version
kind --version
kubectl version --client
clusterctl version
tilt version
aws --version


echo "Verifying system"
whoami
pwd
docker ps


echo "Run a build"
make clusterawsadm
sudo ln -sf ${PWD}/bin/clusterawsadm /usr/local/bin/clusterawsadm
clusterawsadm version

echo "Dev container is ready at $(date)"
