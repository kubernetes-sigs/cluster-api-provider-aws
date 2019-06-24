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
$ ./bin/aws-actuator create --logtostderr -m examples/machine-with-user-data.yaml -u examples/userdata.yml --environment-id UNIQUE_ID
DEBU[0000] Describing AMI ami-a9acbbd6                   example=create-machine machine=test/aws-actuator-testing-machine
Machine creation was successful! InstanceID: i-027681ebf9a842183
```

Once the aws instance is created you can run `$ cat /tmp/test` to verify it contains the `Ahoj` string.

## Test if aws instance exists based on machine manifest

```sh
$ ./bin/aws-actuator exists --logtostderr -m examples/machine-with-user-data.yaml --environment-id UNIQUE_ID
DEBU[0000] checking if machine exists                    example=create-machine machine=test/aws-actuator-testing-machine
DEBU[0000] instance exists as "i-027681ebf9a842183"      example=create-machine machine=test/aws-actuator-testing-machine
Underlying machine's instance exists.
```

## Delete aws instance based on machine manifest

```sh
$ ./bin/aws-actuator delete --logtostderr -m examples/machine-with-user-data.yaml --environment-id UNIQUE_ID
WARN[0000] cleaning up extraneous instance for machine   example=create-machine instanceID=i-027681ebf9a842183 launchTime="2018-08-18 15:50:54 +0000 UTC" machine=test/aws-actuator-testing-machine state=running
INFO[0000] terminating instance                          example=create-machine instanceID=i-027681ebf9a842183 machine=test/aws-actuator-testing-machine
Machine delete operation was successful.
```

## Bootstrapping kubernetes cluster with kubeadm via user data

1. **Generate secret**

   AWS actuator assumes existence of a secret file (references in machine object) with base64 encoded credentials:

   ```yaml
   apiVersion: v1
   kind: Secret
   metadata:
     name: aws-credentials-secret
     namespace: default
   type: Opaque
   data:
     aws_access_key_id: FILLIN
     aws_secret_access_key: FILLIN
   ```

   You can use `examples/render-aws-secrets.sh` script to generate the secret:

   ```sh
   $ ./examples/render-aws-secrets.sh examples/addons.yaml > secret.yaml
   ```

1. **Generate bootstrap user data**

   To generate bootstrap script for machine api plane, simply run:

   ```sh
   $ ./config/generate-bootstrap.sh
   ```

   The script requires `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environment variables to be set.
   It generates `config/bootstrap.yaml` secret for master machine
   under `config/master-machine.yaml`.

   The generated bootstrap secret contains user data responsible for:
   - deployment of kube-apiserver
   - deployment of machine API plane with aws machine controllers
   - generating worker machine user data script secret deploying a node
   - deployment of worker machineset

1. **Create master machine with bootstrapping user data**

   ```sh
   $ ./bin/aws-actuator create -m config/master-machine.yaml -u config/bootstrap.yaml -a secret.yaml
   E0624 15:08:07.983868   30446 utils.go:186] NodeRef not found in machine master-machine
   Machine creation was successful! InstanceID: i-02e6c5f9d1ba3c743
   ```

1. **Pull kubeconfig from created master machine**

   The master public IP can be accessed from AWS Portal. Once done, you
   can collect the kube config by running:

   ```sh
   $ ssh -i SSHPMKEY ec2-user@PUBLICIP 'sudo cat /root/.kube/config' > kubeconfig
   $ kubectl --kubeconfig=kubeconfig config set-cluster kubernetes --server=https://PUBLICIP:8443
   ```

   Once done, you can access the cluster via `kubectl`. E.g.

   ```sh
   $ kubectl --kubeconfig=kubeconfig get nodes
   ```
