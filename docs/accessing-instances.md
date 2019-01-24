# Accessing cluster instances

## Overview

After running `clusterctl` the new cluster will be deployed. This document
explains how to access its nodes throught SSH.

## Prerequisites

* `clusterctl` was successfully executed and cluster is up and running in AWS
environment
* the `cluster-api-provider-aws.sigs.k8s.io` has been created and saved as
`$HOME/.ssh/cluster-api-provider-aws`


## AWS architecture

After cluster creation none of the cluster instances is exposed to the internet,
hence cannot be accessed. To make it accessible, bastion node is created.

### bastion node

bastion node is created in public subnet and is available from the world and
runs official Ubuntu 18.04 Linux image.

### cluster nodes

cluster nodes are whether control plane or worker nodes. They all run official
Ubuntu 18.04 Linux image and are deployed in private subnet.

## Accessing cluster nodes

Cluster nodes should be accessed through the bastion node that is created
along with the cluster. `cluster-api-provider-aws.sigs.k8s.io` SSH key
should be used for authentication.

### Setting up the SSH key path

> Assummig the `cluster-api-provider-aws.sigs.k8s.io` is stored in
`$HOME/.ssh/cluster-api-provider-aws`

```bash
export CLUSTER_SSH_KEY=$HOME/.ssh/cluster-api-provider-aws
```

### Obtain public IP address of the bastion node

> Your credentials must let you query EC2 API

```bash
export BASTION_HOST=$(aws ec2 describe-instances --filter='Name=tag:Name,Values=<CLUSTER_NAME>-bastion' \
	| jq '.Reservations[].Instances[].PublicIpAddress' -r)
```

**NOTE**: If `make manifests` was used to generate manifests, by default the
**CLUSTER_NAME** is set to `test1`

### Get private IP addresses of nodes in the cluster

> Your credentials must let you query EC2 API

```bash
for type in controlplane node
do
	aws ec2 describe-instances --filter="Name=tag:Name,Values=${type}-*" \
		| jq '.Reservations[].Instances[].PrivateIpAddress' -r
done
10.0.0.16
10.0.0.68
```

The above command returns IP addressess of the nodes in the cluster. In this
case **10.0.0.16** and **10.0.0.68**.

### SSHing to the nodes

```bash
ssh -i ${CLUSTER_SSH_KEY} ubuntu@<NODE_IP> \
	-o "ProxyCommand ssh -W %h:%p -i ${CLUSTER_SSH_KEY} ubuntu@${BASTION_HOST}"
```

If the whole document is followed, the value of **NODE_IP** will be either
10.0.0.16 or 10.0.0.16.
