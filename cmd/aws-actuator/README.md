# AWS actuator

The command allows to directly interact with the aws actuator.

## To build the `aws-actuator` binary:

```sh
go build -o bin/aws-actuator -a sigs.k8s.io/cluster-api-provider-aws/cmd/aws-actuator
```

## Create aws instance based on machine manifest

The `examples/userdata.yml` secret encodes the following user data:
```sh
#!/bin/bash
echo "Ahoj" > /tmp/test
```

```sh
$ ./bin/aws-actuator create -m examples/machine-with-user-data.yaml -c examples/cluster.yaml -u examples/userdata.yml
DEBU[0000] Describing AMI ami-a9acbbd6                   example=create-machine machine=test/aws-actuator-testing-machine
Machine creation was successful! InstanceID: i-027681ebf9a842183
```

Once the aws instance is created you can run `$ cat /tmp/test` to verify it contains the `Ahoj` string.

## Test if aws instance exists based on machine manifest

```sh
$ ./bin/aws-actuator exists -m examples/machine-with-user-data.yaml -c examples/cluster.yaml
DEBU[0000] checking if machine exists                    example=create-machine machine=test/aws-actuator-testing-machine
DEBU[0000] instance exists as "i-027681ebf9a842183"      example=create-machine machine=test/aws-actuator-testing-machine
Underlying machine's instance exists.
```

## Delete aws instance based on machine manifest

```sh
$ ./bin/aws-actuator delete -m examples/machine-with-user-data.yaml -c examples/cluster.yaml
WARN[0000] cleaning up extraneous instance for machine   example=create-machine instanceID=i-027681ebf9a842183 launchTime="2018-08-18 15:50:54 +0000 UTC" machine=test/aws-actuator-testing-machine state=running
INFO[0000] terminating instance                          example=create-machine instanceID=i-027681ebf9a842183 machine=test/aws-actuator-testing-machine
Machine delete operation was successful.
```

## Bootstrapping kubernetes cluster with kubeadm via user data

1. Bootstrap the control plane

```sh
./bin/aws-actuator create -m examples/master-machine.yaml -c examples/cluster.yaml -u examples/master-userdata.yaml
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
./bin/aws-actuator create -m examples/worker-machine.yaml -c examples/cluster.yaml -u examples/worker-userdata.yaml
```

After some time the kubernetes cluster with the control plane (master node) and the worker node gets provisioned
and the worker joins the cluster.
