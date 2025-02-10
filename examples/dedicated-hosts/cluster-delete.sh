#!/bin/bash

kubectl delete cluster capa-quickstart
kind delete cluster
aws ec2 release-hosts --host-ids $AWS_HOST_ID