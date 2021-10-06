# Cluster API bootstrap provider EKS

Cluster API bootstrap provider EKS (CABPE) is a component of [Cluster API](https://github.com/kubernetes-sigs/cluster-api/blob/master/README.md) that is responsible for generating a cloud-init script to turn a Machine into a Kubernetes Node; this implementation uses the [AWS-provided EKS bootstrap script](https://github.com/awslabs/amazon-eks-ami/blob/master/files/bootstrap.sh) for joining Kubernetes Nodes to EKS clusters.

CABPE is the bootstrap component of Cluster API Provider AWS' (CAPA) EKS ecosystem. This ecosystem is comprised of:
- EKS controlplane provider (AWSManagedControlPlane)
- AWS infrastructure provider (AWSMachine, AWSMachineTemplate, AWSMachinePool)
- EKS bootstrap provider (EKSConfig, EKSConfigTemplate)

## How does CABPE work?

CABPE generates cloud-init configuration that CAPA will use as the [user-data](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/user-data.html) for new EC2 instances. The generated configuration relies on the `bootstrap.sh` script that is present on all official AWS EKS AMIs.

The output script looks something like this (assuming an EKS cluster with the name `my-cluster`):

```bash
#!/bin/bash
/etc/eks/bootstrap.sh my-cluster
```

CAPA will look up the official EKS AMI for a particular Kubernetes version if the AMI is not overriden in the AWSMachine object. In order for this bootstrap provider to function, any user-specified AMIs will need to have `/etc/eks/bootstrap.sh` present (script source linked above).

Because the bootstrap script has no required fields other than the EKS cluster's name, which can be derived, an `EKSConfig` object with an empty `spec` is acceptable:

```yaml
kind: EKSConfig
apiVersion: bootstrap.cluster.x-k8s.io/v1beta1
metadata:
  name: my-config
spec: {}
```

The only configuration option available is `kubeletExtraArgs`, which is a `map[string]string`:

```yaml
kind: EKSConfig
apiVersion: bootstrap.cluster.x-k8s.io/v1beta1
metadata:
  name: my-config
spec:
  kubeletExtraArgs:
    image-gc-high-threshold: "50"
    image-gc-low-threshold: "40"
    minimum-container-ttl-duration: "1m"
```

This will produce the following:

```bash
#!/bin/bash
/etc/eks/bootstrap.sh my-cluster --kubelet-extra-args '--image-gc-high-threshold=50 --image-gc-low-threshold=45 --minimum-container-ttl-duration=1m'
```

This script will join the Node using the specified `kubelet` arguments.
