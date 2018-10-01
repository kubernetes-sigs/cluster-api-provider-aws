# AWS actuator

The command allows to directly interact with the aws actuator.

## To build the `aws-actuator` binary:

```sh
go build -o bin/aws-actuator sigs.k8s.io/cluster-api-provider-aws/cmd/aws-actuator
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
$ ./bin/aws-actuator create -m examples/machine-with-user-data.yaml -c examples/cluster.yaml -u examples/userdata.yml --environment-id UNIQUE_ID
DEBU[0000] Describing AMI ami-a9acbbd6                   example=create-machine machine=test/aws-actuator-testing-machine
Machine creation was successful! InstanceID: i-027681ebf9a842183
```

Once the aws instance is created you can run `$ cat /tmp/test` to verify it contains the `Ahoj` string.

## Test if aws instance exists based on machine manifest

```sh
$ ./bin/aws-actuator exists -m examples/machine-with-user-data.yaml -c examples/cluster.yaml --environment-id UNIQUE_ID
DEBU[0000] checking if machine exists                    example=create-machine machine=test/aws-actuator-testing-machine
DEBU[0000] instance exists as "i-027681ebf9a842183"      example=create-machine machine=test/aws-actuator-testing-machine
Underlying machine's instance exists.
```

## Delete aws instance based on machine manifest

```sh
$ ./bin/aws-actuator delete -m examples/machine-with-user-data.yaml -c examples/cluster.yaml --environment-id UNIQUE_ID
WARN[0000] cleaning up extraneous instance for machine   example=create-machine instanceID=i-027681ebf9a842183 launchTime="2018-08-18 15:50:54 +0000 UTC" machine=test/aws-actuator-testing-machine state=running
INFO[0000] terminating instance                          example=create-machine instanceID=i-027681ebf9a842183 machine=test/aws-actuator-testing-machine
Machine delete operation was successful.
```

## Bootstrapping kubernetes cluster with kubeadm via user data

1. Bootstrap the control plane

```sh
./bin/aws-actuator create -m examples/master-machine.yaml -c examples/cluster.yaml -u examples/master-userdata.yaml --environment-id UNIQUE_ID
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
./bin/aws-actuator create -m examples/worker-machine.yaml -c examples/cluster.yaml -u examples/worker-userdata.yaml --environment-id UNIQUE_ID
```

After some time the kubernetes cluster with the control plane (master node) and the worker node gets provisioned
and the worker joins the cluster.

## Bootstrapping cluster API stack

The following command will deploy kubernetes cluster with cluster API stack
deployed inside. Worker nodes are deployed with a machineset.
It's assumed both `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` are set.
It takes some time before the worker node joins the cluster (~3 minutes).

```sh
$ ./bin/aws-actuator bootstrap --manifests examples --environment-id UNIQUE_ID
INFO[0000] Reading cluster manifest from examples/cluster.yaml
INFO[0000] Reading master machine manifest from examples/master-machine.yaml
INFO[0000] Reading master user data manifest from examples/master-userdata.yaml
INFO[0000] Creating master machine                      
DEBU[0000] Describing AMI based on filters               bootstrap=create-master-machine machine=test/jchaloup-cama-aws-actuator-testing-machine-master
DEBU[0007] Describing security groups based on filters   bootstrap=create-master-machine machine=test/jchaloup-cama-aws-actuator-testing-machine-master
DEBU[0007] Describing subnets based on filters           bootstrap=create-master-machine machine=test/jchaloup-cama-aws-actuator-testing-machine-master
WARN[0008] More than one subnet id returned, only first one will be used  bootstrap=create-master-machine machine=test/jchaloup-cama-aws-actuator-testing-machine-master
INFO[0009] Master machine created with ipv4: 10.0.102.217, InstanceId: i-0eea29823ae5d50e8
INFO[0014] Waiting for master machine PublicDNS         
DEBU[0014] checking if machine exists                    bootstrap=create-master-machine machine=test/jchaloup-cama-aws-actuator-testing-machine-master
INFO[0014] PublicDnsName: ec2-34-207-227-3.compute-1.amazonaws.com

INFO[0019] Pulling kubeconfig from ec2-34-207-227-3.compute-1.amazonaws.com:8443
INFO[0150] Unable to pull kubeconfig: exit status 255, ssh: connect to host ec2-34-207-227-3.compute-1.amazonaws.com port 22: Connection timed out

INFO[0154] Pulling kubeconfig from ec2-34-207-227-3.compute-1.amazonaws.com:8443
INFO[0158] Unable to pull kubeconfig: exit status 1, Warning: Permanently added 'ec2-34-207-227-3.compute-1.amazonaws.com,34.207.227.3' (ECDSA) to the list of known hosts.
cat: /etc/kubernetes/admin.conf: No such file or directory

INFO[0159] Pulling kubeconfig from ec2-34-207-227-3.compute-1.amazonaws.com:8443
INFO[0162] Unable to pull kubeconfig: exit status 1, cat: /etc/kubernetes/admin.conf: No such file or directory

INFO[0164] Pulling kubeconfig from ec2-34-207-227-3.compute-1.amazonaws.com:8443
INFO[0167] Unable to pull kubeconfig: exit status 1, cat: /etc/kubernetes/admin.conf: No such file or directory

INFO[0169] Pulling kubeconfig from ec2-34-207-227-3.compute-1.amazonaws.com:8443
INFO[0172] Unable to pull kubeconfig: exit status 1, cat: /etc/kubernetes/admin.conf: No such file or directory

INFO[0174] Pulling kubeconfig from ec2-34-207-227-3.compute-1.amazonaws.com:8443
INFO[0177] Unable to pull kubeconfig: exit status 1, cat: /etc/kubernetes/admin.conf: No such file or directory

INFO[0179] Pulling kubeconfig from ec2-34-207-227-3.compute-1.amazonaws.com:8443
INFO[0182] Unable to pull kubeconfig: exit status 1, cat: /etc/kubernetes/admin.conf: No such file or directory

INFO[0184] Pulling kubeconfig from ec2-34-207-227-3.compute-1.amazonaws.com:8443
INFO[0187] Unable to pull kubeconfig: exit status 1, cat: /etc/kubernetes/admin.conf: No such file or directory

INFO[0189] Pulling kubeconfig from ec2-34-207-227-3.compute-1.amazonaws.com:8443
INFO[0191] Unable to pull kubeconfig: exit status 1, cat: /etc/kubernetes/admin.conf: No such file or directory

INFO[0194] Pulling kubeconfig from ec2-34-207-227-3.compute-1.amazonaws.com:8443
INFO[0197] Unable to pull kubeconfig: exit status 1, cat: /etc/kubernetes/admin.conf: No such file or directory

INFO[0199] Pulling kubeconfig from ec2-34-207-227-3.compute-1.amazonaws.com:8443
INFO[0202] Unable to pull kubeconfig: exit status 1, cat: /etc/kubernetes/admin.conf: No such file or directory

INFO[0204] Pulling kubeconfig from ec2-34-207-227-3.compute-1.amazonaws.com:8443
INFO[0207] Unable to pull kubeconfig: exit status 1, cat: /etc/kubernetes/admin.conf: No such file or directory

INFO[0209] Pulling kubeconfig from ec2-34-207-227-3.compute-1.amazonaws.com:8443
INFO[0212] Running kubectl --kubeconfig=kubeconfig config set-cluster kubernetes --server=https://ec2-34-207-227-3.compute-1.amazonaws.com:8443
INFO[0217] Waiting for all nodes to come up             
INFO[0222] Waiting for all nodes to come up             
INFO[0227] Waiting for all nodes to come up             
INFO[0232] Waiting for all nodes to come up             
INFO[0237] Waiting for all nodes to come up             
INFO[0242] Waiting for all nodes to come up             
INFO[0248] Is node "ip-10-0-102-217.ec2.internal" ready?: true

INFO[0248] Deploying cluster-api stack                  
INFO[0248] Deploying aws credentials                    
INFO[0248] Creating "test" namespace...                 
INFO[0248] Creating "test/aws-credentials-secret" secret...
INFO[0254] Deploying cluster-api server                 
INFO[0271] Deploying cluster-api controllers            
INFO[0277] Deploying cluster resource                   
INFO[0277] Creating "test/tb-asg-35" cluster...         
INFO[0277] Unable to deploy cluster manifest: unable to create cluster: an error on the server ("service unavailable") has prevented the request from succeeding (post clusters.cluster.k8s.io)
INFO[0282] Deploying cluster resource                   
INFO[0282] Creating "test/tb-asg-35" cluster...         
INFO[0282] Unable to deploy cluster manifest: unable to create cluster: an error on the server ("service unavailable") has prevented the request from succeeding (post clusters.cluster.k8s.io)
INFO[0287] Deploying cluster resource                   
INFO[0287] Creating "test/tb-asg-35" cluster...         
INFO[0287] Unable to deploy cluster manifest: unable to create cluster: an error on the server ("service unavailable") has prevented the request from succeeding (post clusters.cluster.k8s.io)
INFO[0292] Deploying cluster resource                   
INFO[0292] Creating "test/tb-asg-35" cluster...         
INFO[0292] Unable to deploy cluster manifest: unable to create cluster: an error on the server ("service unavailable") has prevented the request from succeeding (post clusters.cluster.k8s.io)
INFO[0297] Deploying cluster resource                   
INFO[0297] Creating "test/tb-asg-35" cluster...         
INFO[0297] Unable to deploy cluster manifest: unable to create cluster: an error on the server ("service unavailable") has prevented the request from succeeding (post clusters.cluster.k8s.io)
INFO[0302] Deploying cluster resource                   
INFO[0302] Creating "test/tb-asg-35" cluster...         
INFO[0302] Reading worker user data manifest from examples/worker-userdata.yaml
INFO[0302] Generating worker machine set user data for master listening at 10.0.102.217
INFO[0302] Creating "test/aws-actuator-node-user-data-secret" secret...
INFO[0302] Reading worker machine manifest from examples/worker-machineset.yaml
INFO[0307] Deploying worker machineset                  
INFO[0307] Creating "test/jchaloup-cama-default-worker-machineset" machineset...
```
