# AWS actuator

The command allows to directly interact with the aws actuator.

## To build the `aws-actuator` binary:

```sh
$ make aws-actuator
```

## Prerequisities

All the machine manifests assume existence for various AWS resources such as vpc, subnet,
security groups, etc. In case you are starting from scratch, you can run `hack/aws-provision.sh` to
create all the required resources:

```sh
$ cd hack
$ ENVIRONMENT_ID=UNIQUE_ID ./aws-provision.sh install
```

where `UNIQUE_ID` is unique identification of your environment.

Don't forget to run `./aws-provision.sh destroy` once you are done.

## Create aws instance based on machine manifest

The `examples/userdata.yml` secret encodes the following user data:
```sh
#!/bin/bash
echo "Ahoj" > /tmp/test
```

```sh
$ ./bin/aws-actuator create --logtostderr -m examples/machine-with-user-data.yaml -c examples/cluster.yaml -u examples/userdata.yml --environment-id UNIQUE_ID
DEBU[0000] Describing AMI ami-a9acbbd6                   example=create-machine machine=test/aws-actuator-testing-machine
Machine creation was successful! InstanceID: i-027681ebf9a842183
```

Once the aws instance is created you can run `$ cat /tmp/test` to verify it contains the `Ahoj` string.

## Test if aws instance exists based on machine manifest

```sh
$ ./bin/aws-actuator exists --logtostderr -m examples/machine-with-user-data.yaml -c examples/cluster.yaml --environment-id UNIQUE_ID
DEBU[0000] checking if machine exists                    example=create-machine machine=test/aws-actuator-testing-machine
DEBU[0000] instance exists as "i-027681ebf9a842183"      example=create-machine machine=test/aws-actuator-testing-machine
Underlying machine's instance exists.
```

## Delete aws instance based on machine manifest

```sh
$ ./bin/aws-actuator delete --logtostderr -m examples/machine-with-user-data.yaml -c examples/cluster.yaml --environment-id UNIQUE_ID
WARN[0000] cleaning up extraneous instance for machine   example=create-machine instanceID=i-027681ebf9a842183 launchTime="2018-08-18 15:50:54 +0000 UTC" machine=test/aws-actuator-testing-machine state=running
INFO[0000] terminating instance                          example=create-machine instanceID=i-027681ebf9a842183 machine=test/aws-actuator-testing-machine
Machine delete operation was successful.
```

## Bootstrapping kubernetes cluster with kubeadm via user data

1. Bootstrap the control plane

```sh
./bin/aws-actuator create --logtostderr -m examples/master-machine.yaml -c examples/cluster.yaml -u examples/master-userdata.yaml --environment-id UNIQUE_ID
```

By default networking is enabled on the master machine. You can
interact with the master by copying the kubeconfig and using
`kubectl`. You'll need to find the external IP address of your master
and then:

```sh
ssh ec2-user@ec2-184-73-119-192.compute-1.amazonaws.com "sudo cat /etc/kubernetes/admin.conf" > kubeconfig
export KUBECONFIG=$PWD/kubeconfig
kubectl config set-cluster kubernetes --server=https://ec2-184-73-119-192.compute-1.amazonaws.com:8443
kubectl get nodes
NAME                           STATUS    ROLES     AGE       VERSION
ip-172-31-34-42.ec2.internal   Ready     master    6m        v1.11.2
```

2. Once the ip address of the master node is known (e.g. 172.31.34.2), update the `kubeadm join` line in `examples/worker-user-data.sh` to:
```sh
kubeadm join 172.31.34.2:8443 --token 2iqzqm.85bs0x6miyx1nm7l --discovery-token-unsafe-skip-ca-verification
```

You can get the internal IP address of the cluster dynamically by
running:

```sh
echo $(ssh ec2-user@ec2-184-73-119-192.compute-1.amazonaws.com wget -qO - http://169.254.169.254/latest/meta-data/local-ipv4)
172.31.34.42
```

3. Encode the `examples/worker-user-data.sh` by running:

```sh
$ cat examples/worker-user-data.sh | base64
```

4. Update the `data.userData` of `examples/worker-user-data.yaml` with the generated hash

5. Create the worker node by running:

```sh
$ ./bin/aws-actuator create --logtostderr -m examples/worker-machine.yaml -c examples/cluster.yaml -u examples/worker-userdata.yaml --environment-id UNIQUE_ID
```

After some time the kubernetes cluster with the control plane (master node) and the worker node gets provisioned
and the worker joins the cluster.

## Bootstrapping cluster API stack

The following command will deploy kubernetes cluster with cluster API stack
deployed inside. Worker nodes are deployed with a machineset.
It's assumed both `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` are set.
It takes some time before the worker node joins the cluster (~3 minutes).

```sh
$ ./bin/aws-actuator bootstrap --logtostderr --environment-id UNIQUE_ID --master-machine-private-key AWS_MACHINE_PRIVATE_KEY -v 3
I1126 00:27:50.704676   13319 main.go:280] Creating master machine
I1126 00:27:52.525208   13319 actuator.go:188] no stopped instances found for machine UNIQUE_ID-master-machine-b9aee4
I1126 00:27:52.525248   13319 actuator.go:306] Describing AMI based on filters
I1126 00:27:54.523474   13319 actuator.go:217] Describing security groups based on filters
I1126 00:27:54.978100   13319 actuator.go:252] Describing subnets based on filters
E1126 00:27:56.524317   13319 main.go:284] <nil>
I1126 00:27:56.524363   13319 main.go:289] Master machine created with ipv4: 10.0.101.124, InstanceId: i-0a2f757591c6c7fab
I1126 00:28:01.524596   13319 main.go:294] Waiting for master machine PublicDNS
I1126 00:28:01.524638   13319 actuator.go:581] checking if machine exists
I1126 00:28:01.947078   13319 main.go:301] PublicDnsName: ec2-34-200-235-67.compute-1.amazonaws.com
I1126 00:28:01.948117   13319 main.go:331] Collecting master kubeconfig
I1126 00:28:12.268962   13319 machines.go:293] Pulling kubeconfig from ec2-34-200-235-67.compute-1.amazonaws.com:8443
I1126 00:30:22.074665   13319 machines.go:301] Unable to pull kubeconfig: failed to dial: dial tcp 34.200.235.67:22: connect: connection timed out
I1126 00:30:22.268985   13319 machines.go:293] Pulling kubeconfig from ec2-34-200-235-67.compute-1.amazonaws.com:8443
I1126 00:30:25.343047   13319 machines.go:301] Unable to pull kubeconfig: failed to collect kubeconfig: Process exited with status 1, cat: /root/.kube/config: No such file or directory
I1126 00:30:27.269053   13319 machines.go:293] Pulling kubeconfig from ec2-34-200-235-67.compute-1.amazonaws.com:8443
I1126 00:30:28.890164   13319 machines.go:301] Unable to pull kubeconfig: failed to collect kubeconfig: Process exited with status 1, cat: /root/.kube/config: No such file or directory
I1126 00:30:32.269160   13319 machines.go:293] Pulling kubeconfig from ec2-34-200-235-67.compute-1.amazonaws.com:8443
I1126 00:30:34.130056   13319 machines.go:301] Unable to pull kubeconfig: failed to collect kubeconfig: Process exited with status 1, cat: /root/.kube/config: No such file or directory
I1126 00:30:37.269024   13319 machines.go:293] Pulling kubeconfig from ec2-34-200-235-67.compute-1.amazonaws.com:8443
I1126 00:30:38.923598   13319 machines.go:301] Unable to pull kubeconfig: failed to collect kubeconfig: Process exited with status 1, cat: /root/.kube/config: No such file or directory
I1126 00:30:42.269020   13319 machines.go:293] Pulling kubeconfig from ec2-34-200-235-67.compute-1.amazonaws.com:8443
I1126 00:30:43.873202   13319 machines.go:301] Unable to pull kubeconfig: failed to collect kubeconfig: Process exited with status 1, cat: /root/.kube/config: No such file or directory
I1126 00:30:47.269133   13319 machines.go:293] Pulling kubeconfig from ec2-34-200-235-67.compute-1.amazonaws.com:8443
I1126 00:30:49.110626   13319 machines.go:301] Unable to pull kubeconfig: failed to collect kubeconfig: Process exited with status 1, cat: /root/.kube/config: No such file or directory
I1126 00:30:52.269046   13319 machines.go:293] Pulling kubeconfig from ec2-34-200-235-67.compute-1.amazonaws.com:8443
I1126 00:30:54.229656   13319 machines.go:301] Unable to pull kubeconfig: failed to collect kubeconfig: Process exited with status 1, cat: /root/.kube/config: No such file or directory
I1126 00:30:57.269052   13319 machines.go:293] Pulling kubeconfig from ec2-34-200-235-67.compute-1.amazonaws.com:8443
I1126 00:30:59.099948   13319 machines.go:301] Unable to pull kubeconfig: failed to collect kubeconfig: Process exited with status 1, cat: /root/.kube/config: No such file or directory
I1126 00:31:02.269091   13319 machines.go:293] Pulling kubeconfig from ec2-34-200-235-67.compute-1.amazonaws.com:8443
I1126 00:31:04.182657   13319 machines.go:301] Unable to pull kubeconfig: failed to collect kubeconfig: Process exited with status 1, cat: /root/.kube/config: No such file or directory
I1126 00:31:07.269189   13319 machines.go:293] Pulling kubeconfig from ec2-34-200-235-67.compute-1.amazonaws.com:8443
I1126 00:31:09.290235   13319 machines.go:301] Unable to pull kubeconfig: failed to collect kubeconfig: Process exited with status 1, cat: /root/.kube/config: No such file or directory
I1126 00:31:12.269123   13319 machines.go:293] Pulling kubeconfig from ec2-34-200-235-67.compute-1.amazonaws.com:8443
I1126 00:31:14.151180   13319 machines.go:301] Unable to pull kubeconfig: failed to collect kubeconfig: Process exited with status 1, cat: /root/.kube/config: No such file or directory
I1126 00:31:17.269023   13319 machines.go:293] Pulling kubeconfig from ec2-34-200-235-67.compute-1.amazonaws.com:8443
I1126 00:31:19.373448   13319 machines.go:301] Unable to pull kubeconfig: failed to collect kubeconfig: Process exited with status 1, cat: /root/.kube/config: No such file or directory
I1126 00:31:22.269634   13319 machines.go:293] Pulling kubeconfig from ec2-34-200-235-67.compute-1.amazonaws.com:8443
I1126 00:31:24.905696   13319 machines.go:301] Unable to pull kubeconfig: failed to collect kubeconfig: Process exited with status 1, cat: /root/.kube/config: No such file or directory
I1126 00:31:27.269061   13319 machines.go:293] Pulling kubeconfig from ec2-34-200-235-67.compute-1.amazonaws.com:8443
I1126 00:31:28.999498   13319 machines.go:301] Unable to pull kubeconfig: failed to collect kubeconfig: Process exited with status 1, cat: /root/.kube/config: No such file or directory
I1126 00:31:32.269036   13319 machines.go:293] Pulling kubeconfig from ec2-34-200-235-67.compute-1.amazonaws.com:8443
I1126 00:31:34.325456   13319 machines.go:301] Unable to pull kubeconfig: failed to collect kubeconfig: Process exited with status 1, cat: /root/.kube/config: No such file or directory
I1126 00:31:37.268993   13319 machines.go:293] Pulling kubeconfig from ec2-34-200-235-67.compute-1.amazonaws.com:8443
I1126 00:31:39.242166   13319 machines.go:301] Unable to pull kubeconfig: failed to collect kubeconfig: Process exited with status 1, cat: /root/.kube/config: No such file or directory
I1126 00:31:42.269179   13319 machines.go:293] Pulling kubeconfig from ec2-34-200-235-67.compute-1.amazonaws.com:8443
I1126 00:31:44.257477   13319 machines.go:301] Unable to pull kubeconfig: failed to collect kubeconfig: Process exited with status 1, cat: /root/.kube/config: No such file or directory
I1126 00:31:47.269127   13319 machines.go:293] Pulling kubeconfig from ec2-34-200-235-67.compute-1.amazonaws.com:8443
I1126 00:31:49.276483   13319 machines.go:301] Unable to pull kubeconfig: failed to collect kubeconfig: Process exited with status 1, cat: /root/.kube/config: No such file or directory
I1126 00:31:52.269210   13319 machines.go:293] Pulling kubeconfig from ec2-34-200-235-67.compute-1.amazonaws.com:8443
I1126 00:31:54.394890   13319 machines.go:301] Unable to pull kubeconfig: failed to collect kubeconfig: Process exited with status 1, cat: /root/.kube/config: No such file or directory
I1126 00:31:57.269117   13319 machines.go:293] Pulling kubeconfig from ec2-34-200-235-67.compute-1.amazonaws.com:8443
I1126 00:31:59.209442   13319 machines.go:301] Unable to pull kubeconfig: failed to collect kubeconfig: Process exited with status 1, cat: /root/.kube/config: No such file or directory
I1126 00:32:02.269029   13319 machines.go:293] Pulling kubeconfig from ec2-34-200-235-67.compute-1.amazonaws.com:8443
I1126 00:32:03.969797   13319 machines.go:311] Master running on https://ec2-34-200-235-67.compute-1.amazonaws.com:8443
I1126 00:32:03.982206   13319 main.go:361] Waiting for all nodes to come up
I1126 00:32:09.451771   13319 machines.go:222] unable to list nodes: Get https://ec2-34-200-235-67.compute-1.amazonaws.com:8443/api/v1/nodes: x509: certificate has expired or is not yet valid
I1126 00:32:14.365702   13319 machines.go:222] unable to list nodes: Get https://ec2-34-200-235-67.compute-1.amazonaws.com:8443/api/v1/nodes: x509: certificate has expired or is not yet valid
I1126 00:32:19.383442   13319 machines.go:222] unable to list nodes: Get https://ec2-34-200-235-67.compute-1.amazonaws.com:8443/api/v1/nodes: x509: certificate has expired or is not yet valid
I1126 00:32:24.400808   13319 machines.go:222] unable to list nodes: Get https://ec2-34-200-235-67.compute-1.amazonaws.com:8443/api/v1/nodes: x509: certificate has expired or is not yet valid
I1126 00:32:29.274112   13319 machines.go:222] unable to list nodes: Get https://ec2-34-200-235-67.compute-1.amazonaws.com:8443/api/v1/nodes: x509: certificate has expired or is not yet valid
I1126 00:32:34.536479   13319 machines.go:222] unable to list nodes: Get https://ec2-34-200-235-67.compute-1.amazonaws.com:8443/api/v1/nodes: x509: certificate has expired or is not yet valid
I1126 00:32:39.349722   13319 machines.go:222] unable to list nodes: Get https://ec2-34-200-235-67.compute-1.amazonaws.com:8443/api/v1/nodes: x509: certificate has expired or is not yet valid
I1126 00:32:44.786880   13319 machines.go:239] Node "ip-10-0-101-124.ec2.internal" is ready
I1126 00:32:44.786916   13319 main.go:367] Creating "test" namespace
I1126 00:32:44.981447   13319 main.go:358] Deploying cluster API stack components
I1126 00:32:44.981463   13319 main.go:358] Deploying cluster CRD manifest
I1126 00:32:50.206519   13319 framework.go:302] create.err: <nil>
I1126 00:32:50.408804   13319 framework.go:311] get.err: <nil>
I1126 00:32:50.408827   13319 main.go:358] Deploying machine CRD manifest
I1126 00:32:55.632340   13319 framework.go:302] create.err: <nil>
I1126 00:32:55.751620   13319 framework.go:311] get.err: <nil>
I1126 00:32:55.752952   13319 main.go:358] Deploying machineset CRD manifest
I1126 00:33:00.957146   13319 framework.go:302] create.err: <nil>
I1126 00:33:01.162107   13319 framework.go:311] get.err: <nil>
I1126 00:33:01.162204   13319 main.go:358] Deploying machinedeployment CRD manifest
I1126 00:33:06.383504   13319 framework.go:302] create.err: <nil>
I1126 00:33:06.589079   13319 framework.go:311] get.err: <nil>
I1126 00:33:06.589126   13319 main.go:358] Deploying cluster role
I1126 00:33:06.947393   13319 main.go:358] Deploying controller manager
I1126 00:33:07.239718   13319 main.go:358] Deploying machine controller
I1126 00:33:07.412079   13319 main.go:358] Waiting until cluster objects can be listed
I1126 00:33:12.636542   13319 main.go:358] Cluster API stack deployed
I1126 00:33:12.636656   13319 main.go:358] Creating "UNIQUE_ID" cluster
gI1126 00:33:33.622648   13319 main.go:358] Creating "UNIQUE_ID-worker-machineset-861404" machineset
I1126 00:33:43.973580   13319 main.go:358] Verify machineset's underlying instances is running
I1126 00:33:44.099488   13319 main.go:358] Waiting for "UNIQUE_ID-worker-machineset-861404-qh6sx" machine
I1126 00:33:50.459899   13319 main.go:358] Verify machine's underlying instance is running
I1126 00:33:55.460140   13319 machines.go:80] Waiting for instance to come up
I1126 00:33:57.176449   13319 machines.go:88] Machine is running
```
