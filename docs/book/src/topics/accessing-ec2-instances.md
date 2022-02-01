# Accessing cluster instances

## Overview

After running `clusterctl generate cluster` to generate the configuration for a new workload cluster (and then redirecting that output to a file for use with `kubectl apply`, or piping it directly to `kubectl apply`), the new workload cluster will be deployed. This document explains how to access the new workload cluster's nodes.

## Prerequisites

1. `clusterctl generate cluster` was successfully executed to generate the configuration for a new workload cluster
2. The configuration for the new workload cluster was applied to the management cluster using `kubectl apply` and the cluster is up and running in an AWS environment.
3. The SSH key referenced by `clusterctl` in step 1 exists in AWS and is stored in the correct location locally for use by SSH (on macOS/Linux systems, this is typically `$HOME/.ssh`). This document will refer to this key as `cluster-api-provider-aws.sigs.k8s.io`.
4. _(If using AWS Session Manager)_ The AWS CLI and the Session Manager plugin have been installed and configured.

## Methods for accessing nodes

There are two ways to access cluster nodes once the workload cluster is up and running:

* via SSH
* via AWS Session Manager

### Accessing nodes via SSH

By default, workload clusters created in AWS will _not_ support access via SSH apart from AWS Session Manager (see the section titled "Accessing nodes via AWS Session Manager"). However, the manifest for a workload cluster can be modified to include an SSH bastion host, created and managed by the management cluster, to enable SSH access to cluster nodes. The bastion node is created in a public subnet and provides SSH access from the world. It runs the official Ubuntu Linux image.

#### Enabling the bastion host

To configure the Cluster API Provider for AWS to create an SSH bastion host, add this line to the AWSCluster spec:

```yaml
spec:
  bastion:
    enabled: true
```

#### Obtain public IP address of the bastion node

Once the workload cluster is up and running after being configured for an SSH bastion host, you can use the `kubectl get awscluster` command to look up the public IP address of the bastion host (make sure the `kubectl` context is set to the management cluster). The output will look something like this:

```bash
NAME   CLUSTER   READY   VPC                     BASTION IP
test   test      true    vpc-1739285ed052be7ad   1.2.3.4
```

#### Setting up the SSH key path

Assumming that the `cluster-api-provider-aws.sigs.k8s.io` SSH key is stored in
`$HOME/.ssh/cluster-api-provider-aws`, use this command to set up an environment variable for use in a later command:

```bash
export CLUSTER_SSH_KEY=$HOME/.ssh/cluster-api-provider-aws
```

#### Get private IP addresses of nodes in the cluster

To get the private IP addresses of nodes in the cluster (nodes may be control plane nodes or worker nodes), use this `kubectl` command with the context set to the management cluster:

```bash
kubectl get nodes -o custom-columns=NAME:.metadata.name,\
IP:"{.status.addresses[?(@.type=='InternalIP')].address}"
```

This will produce output that looks like this:

```bash
NAME                                         IP
ip-10-0-0-16.us-west-2.compute.internal   10.0.0.16
ip-10-0-0-68.us-west-2.compute.internal   10.0.0.68
```

The above command returns IP addresses of the nodes in the cluster. In this
case, the values returned are `10.0.0.16` and `10.0.0.68`.

### Connecting to the nodes via SSH

To access one of the nodes (either a control plane node or a worker node) via the SSH bastion host, use this command if you are using a non-EKS cluster:

```bash
ssh -i ${CLUSTER_SSH_KEY} ubuntu@<NODE_IP> \
	-o "ProxyCommand ssh -W %h:%p -i ${CLUSTER_SSH_KEY} ubuntu@${BASTION_HOST}"
```

And use this command if you are using a EKS based cluster:

```bash
ssh -i ${CLUSTER_SSH_KEY} ec2-user@<NODE_IP> \
	-o "ProxyCommand ssh -W %h:%p -i ${CLUSTER_SSH_KEY} ubuntu@${BASTION_HOST}"
```


If the whole document is followed, the value of `<NODE_IP>` will be either
10.0.0.16 or 10.0.0.68.

Alternately, users can add a configuration stanza to their SSH configuration file (typically found on macOS/Linux systems as `$HOME/.ssh/config`):

```text
Host 10.0.*
  User ubuntu
  IdentityFile <CLUSTER_SSH_KEY>
  ProxyCommand ssh -W %h:%p ubuntu@<BASTION_HOST>
```

### Accessing nodes via AWS Session Manager

All CAPA-published AMIs based on Ubuntu have the AWS SSM Agent pre-installed (as a Snap package; this was added in June 2018 to the base Ubuntu Server image for all 16.04 and later AMIs). This allows users to access cluster nodes directly, without the need for an SSH bastion host, using the AWS CLI and the Session Manager plugin.

To access a cluster node (control plane node or worker node), you'll need the instance ID. You can retrieve the instance ID using this `kubectl` command with the context set to the management cluster:

```bash
kubectl get awsmachines -o custom-columns=NAME:.metadata.name,INSTANCEID:.spec.providerID
```

This will produce output similar to this:

```bash
NAME                      INSTANCEID
test-controlplane-52fhh   aws:////i-112bac41a19da1819
test-controlplane-lc5xz   aws:////i-99aaef2381ada9228
```

Users can then use the instance ID (everything after the `aws:////` prefix) to connect to the cluster node with this command:

```bash
aws ssm start-session --target <INSTANCE_ID>
```

This will log you into the cluster node as the `ssm-user` user ID.

## Additional Notes

### Using the AWS CLI instead of `kubectl`

It is also possible to use AWS CLI commands instead of `kubectl` to gather information about the cluster nodes.

For example, to use the AWS CLI to get the public IP address of the SSH bastion host, use this AWS CLI command:

```bash
export BASTION_HOST=$(aws ec2 describe-instances --filter='Name=tag:Name,Values=<CLUSTER_NAME>-bastion' \
	| jq '.Reservations[].Instances[].PublicIpAddress' -r)
```

You should substitute the correct cluster name for `<CLUSTER_NAME>` in the above command. (**NOTE**: If `make manifests` was used to generate manifests, by default the `<CLUSTER_NAME>` is set to `test1`.)

Similarly, to obtain the list of private IP addresses of the cluster nodes, use this AWS CLI command:

```bash
for type in control-plane node
do
	aws ec2 describe-instances \
    --filter="Name=tag:sigs.k8s.io/cluster-api-provider-aws/role,\
    Values=${type}" \
		| jq '.Reservations[].Instances[].PrivateIpAddress' -r
done
10.0.0.16
10.0.0.68
```

Finally, to obtain AWS instance IDs for cluster nodes, you can use this AWS CLI command:

```bash
for type in control-plane node
do
	aws ec2 describe-instances \
    --filter="Name=tag:sigs.k8s.io/cluster-api-provider-aws/role,\
    Values=${type}" \
		| jq '.Reservations[].Instances[].InstanceId' -r
done
i-112bac41a19da1819
i-99aaef2381ada9228
```

Note that your AWS CLI must be configured with credentials that enable you to query the AWS EC2 API.
