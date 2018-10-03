#!/bin/bash

image=${1}
sshhost=${2}
sshkey=${3}

#docker save docker.io/centos:centos7 | bzip2 | ssh -i /home/jchaloup/.ssh/libra.pem ec2-user@ec2-34-200-213-218.compute-1.amazonaws.com 'bunzip2 > /tmp/sideimage && sudo docker image import /tmp/sideimage docker.io/centos:centos7'
docker save ${image} | bzip2 | ssh -i ${sshkey} ec2-user@${sshhost} "bunzip2 > /tmp/tempimage.bz2 && sudo docker load -i /tmp/tempimage.bz2"
echo "image: $image"
echo "sshhost: $sshhost"
