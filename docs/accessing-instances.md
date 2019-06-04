# Accessing cluster instances

## Overview

After running `clusterctl`, the new cluster will be deployed. This document
explains how to access it's nodes through SSH.

## Prerequisites

* `clusterctl` was successfully executed and the cluster is up and running in
an AWS environment.
* The `cluster-api-provider-aws.sigs.k8s.io` SSH key has been created in AWS and saved
as `$HOME/.ssh/cluster-api-provider-aws`

## AWS architecture

After cluster creation none of the cluster instances are exposed to the
internet, i.e. cannot be accessed. To make it accessible, `clusterctl` also
creates a bastion node.

### Bastion node

The Bastion node is created in a public subnet and provides SSH access from the
world. It runs the official Ubuntu 18.04 Linux image.

### Cluster nodes

Cluster nodes are either control plane or worker nodes. They all run the
official Ubuntu 18.04 Linux image and are deployed in a private subnet.

## Accessing cluster nodes

Cluster nodes should be accessed through the bastion node that is created
along with the cluster. The `cluster-api-provider-aws.sigs.k8s.io` SSH key
should be used for authentication.

### Setting up the SSH key path

> Assumming that the `cluster-api-provider-aws.sigs.k8s.io` SSH key is stored in
`$HOME/.ssh/cluster-api-provider-aws`

```bash
export CLUSTER_SSH_KEY=$HOME/.ssh/cluster-api-provider-aws
```

### Obtain public IP address of the bastion node

> Your credentials must let you query the EC2 API.

```bash
export BASTION_HOST=$(aws ec2 describe-instances --filter='Name=tag:Name,Values=<CLUSTER_NAME>-bastion' \
	| jq '.Reservations[].Instances[].PublicIpAddress' -r)
```

**NOTE**: If `make manifests` was used to generate manifests, by default the
**CLUSTER_NAME** is set to `test1`

### Get private IP addresses of nodes in the cluster

> Your credentials must let you query the EC2 API.

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
