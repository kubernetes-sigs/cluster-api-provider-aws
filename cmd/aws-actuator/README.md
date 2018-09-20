# AWS actuator

The command allows to directly interact with the aws actuator.

## To build the `aws-actuator` binary:

```sh
go build -o bin/aws-actuator -a sigs.k8s.io/cluster-api-provider-aws/cmd/aws-actuator
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

### All in one

Alternatively, you can run the `aws-actuator bootstrap` that does all the above (up to step 2.):

```sh
./bin/aws-actuator bootstrap --manifests examples --environment-id UNIQUE_ID
INFO[0000] Reading cluster manifest from examples/cluster.yaml
INFO[0000] Reading master machine manifest from examples/master-machine.yaml
INFO[0000] Reading master user data manifest from examples/master-userdata.yaml
INFO[0000] Reading worker machine manifest from examples/worker-machine.yaml
INFO[0000] Reading worker user data manifest from examples/worker-userdata.yaml
INFO[0000] Creating master machine                      
DEBU[0000] Describing AMI based on filters               bootstrap=create-master-machine machine=test/UNIQUE_ID-aws-actuator-testing-machine-master
DEBU[0007] Describing security groups based on filters   bootstrap=create-master-machine machine=test/UNIQUE_ID-aws-actuator-testing-machine-master
DEBU[0008] Describing subnets based on filters           bootstrap=create-master-machine machine=test/UNIQUE_ID-aws-actuator-testing-machine-master
WARN[0008] More than one subnet id returned, only first one will be used  bootstrap=create-master-machine machine=test/UNIQUE_ID-aws-actuator-testing-machine-master
INFO[0009] Master machine created with ipv4: 10.0.102.149, InstanceId: i-0cd65d6ce5640d343
INFO[0009] Generating worker user data for master listening at 10.0.102.149
INFO[0009] Creating worker machine                      
INFO[0009] no stopped instances found for machine UNIQUE_ID-aws-actuator-testing-machine-worker  bootstrap=create-worker-machine machine=test/UNIQUE_ID-aws-actuator-testing-machine-worker
DEBU[0009] Describing AMI based on filters               bootstrap=create-worker-machine machine=test/UNIQUE_ID-aws-actuator-testing-machine-worker
DEBU[0014] Describing security groups based on filters   bootstrap=create-worker-machine machine=test/UNIQUE_ID-aws-actuator-testing-machine-worker
DEBU[0014] Describing subnets based on filters           bootstrap=create-worker-machine machine=test/UNIQUE_ID-aws-actuator-testing-machine-worker
WARN[0015] More than one subnet id returned, only first one will be used  bootstrap=create-worker-machine machine=test/UNIQUE_ID-aws-actuator-testing-machine-worker
INFO[0016] Worker machine created with InstanceId: i-0763fb7fafc607ecf
```

## Bootstrapping cluster API stack

Running the `aws-actuator bootstrap` with `--cluster-api-stack` will deploy the cluster API stack as well.
It's assumed both `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` are set.

```sh
$ ./bin/aws-actuator bootstrap --manifests examples --environment-id UNIQUE_ID --cluster-api-stack
INFO[0000] Reading cluster manifest from examples/cluster.yaml
INFO[0000] Reading master machine manifest from examples/master-machine.yaml
INFO[0000] Reading master user data manifest from examples/master-userdata.yaml
INFO[0000] Reading worker machine manifest from examples/worker-machine.yaml
INFO[0000] Reading worker user data manifest from examples/worker-userdata.yaml
INFO[0000] Creating master machine                      
DEBU[0000] Describing AMI based on filters               bootstrap=create-master-machine machine=test/UNIQUE_ID-aws-actuator-testing-machine-master
DEBU[0007] Describing security groups based on filters   bootstrap=create-master-machine machine=test/UNIQUE_ID-aws-actuator-testing-machine-master
DEBU[0007] Describing subnets based on filters           bootstrap=create-master-machine machine=test/UNIQUE_ID-aws-actuator-testing-machine-master
WARN[0007] More than one subnet id returned, only first one will be used  bootstrap=create-master-machine machine=test/UNIQUE_ID-aws-actuator-testing-machine-master
INFO[0008] Master machine created with ipv4: 10.0.101.159, InstanceId: i-04c41ad24e885a8c6
INFO[0008] Generating worker user data for master listening at 10.0.101.159
INFO[0008] Creating worker machine                      
INFO[0009] no stopped instances found for machine UNIQUE_ID-aws-actuator-testing-machine-worker  bootstrap=create-worker-machine machine=test/UNIQUE_ID-aws-actuator-testing-machine-worker
DEBU[0009] Describing AMI based on filters               bootstrap=create-worker-machine machine=test/UNIQUE_ID-aws-actuator-testing-machine-worker
DEBU[0012] Describing security groups based on filters   bootstrap=create-worker-machine machine=test/UNIQUE_ID-aws-actuator-testing-machine-worker
DEBU[0013] Describing subnets based on filters           bootstrap=create-worker-machine machine=test/UNIQUE_ID-aws-actuator-testing-machine-worker
WARN[0013] More than one subnet id returned, only first one will be used  bootstrap=create-worker-machine machine=test/UNIQUE_ID-aws-actuator-testing-machine-worker
INFO[0014] Worker machine created with InstanceId: i-0d548c5592e4e78a7
INFO[0019] Waiting for master machine PublicDNS         
DEBU[0019] checking if machine exists                    bootstrap=create-worker-machine machine=test/UNIQUE_ID-aws-actuator-testing-machine-master
INFO[0019] PublicDnsName: ec2-34-239-226-191.compute-1.amazonaws.com

INFO[0024] Pulling kubeconfig from ec2-34-239-226-191.compute-1.amazonaws.com:8443
INFO[0093] Unable to pull kubeconfig: exit status 1, Warning: Permanently added 'ec2-34-239-226-191.compute-1.amazonaws.com,34.239.226.191' (ECDSA) to the list of known hosts.
cat: /etc/kubernetes/admin.conf: No such file or directory

INFO[0094] Pulling kubeconfig from ec2-34-239-226-191.compute-1.amazonaws.com:8443
INFO[0096] Unable to pull kubeconfig: exit status 1, cat: /etc/kubernetes/admin.conf: No such file or directory

INFO[0099] Pulling kubeconfig from ec2-34-239-226-191.compute-1.amazonaws.com:8443
INFO[0101] Unable to pull kubeconfig: exit status 1, cat: /etc/kubernetes/admin.conf: No such file or directory

INFO[0104] Pulling kubeconfig from ec2-34-239-226-191.compute-1.amazonaws.com:8443
INFO[0106] Unable to pull kubeconfig: exit status 1, cat: /etc/kubernetes/admin.conf: No such file or directory

INFO[0109] Pulling kubeconfig from ec2-34-239-226-191.compute-1.amazonaws.com:8443
INFO[0111] Unable to pull kubeconfig: exit status 1, cat: /etc/kubernetes/admin.conf: No such file or directory

INFO[0114] Pulling kubeconfig from ec2-34-239-226-191.compute-1.amazonaws.com:8443
INFO[0116] Unable to pull kubeconfig: exit status 1, cat: /etc/kubernetes/admin.conf: No such file or directory

INFO[0119] Pulling kubeconfig from ec2-34-239-226-191.compute-1.amazonaws.com:8443
INFO[0121] Unable to pull kubeconfig: exit status 1, cat: /etc/kubernetes/admin.conf: No such file or directory

INFO[0124] Pulling kubeconfig from ec2-34-239-226-191.compute-1.amazonaws.com:8443
INFO[0126] Unable to pull kubeconfig: exit status 1, cat: /etc/kubernetes/admin.conf: No such file or directory

INFO[0129] Pulling kubeconfig from ec2-34-239-226-191.compute-1.amazonaws.com:8443
INFO[0131] Unable to pull kubeconfig: exit status 1, cat: /etc/kubernetes/admin.conf: No such file or directory

INFO[0134] Pulling kubeconfig from ec2-34-239-226-191.compute-1.amazonaws.com:8443
INFO[0136] Unable to pull kubeconfig: exit status 1, cat: /etc/kubernetes/admin.conf: No such file or directory

INFO[0139] Pulling kubeconfig from ec2-34-239-226-191.compute-1.amazonaws.com:8443
INFO[0146] Running kubectl config set-cluster kubernetes --server=https://ec2-34-239-226-191.compute-1.amazonaws.com:8443
INFO[0151] Waiting for all nodes to come up             
INFO[0156] Waiting for all nodes to come up             
INFO[0161] Waiting for all nodes to come up             
INFO[0166] Waiting for all nodes to come up             
INFO[0171] Waiting for all nodes to come up             
INFO[0179] Is node "ip-10-0-101-159.ec2.internal" ready?: true

INFO[0179] Deploying cluster-api stack                  
INFO[0179] Deploying aws credentials                    
INFO[0179] Creating "test" namespace...                 
INFO[0179] Creating "test/aws-credentials-secret" secret...
INFO[0185] Deploying cluster-api server                 
INFO[0197] Deploying cluster-api controllers
```
