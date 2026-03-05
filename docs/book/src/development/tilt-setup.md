# Developing Cluster API Provider AWS with Tilt

This document describes how to use kind and [Tilt][tilt] for a simplified workflow that offers easy deployments and rapid iterative builds.
Before the next steps, make sure [initial setup for development environment][Initial-setup-for-development-environment] steps are complete.

Also, visit the [Cluster API documentation on Tilt][cluster_api_tilt] for more information on how to set up your development environment.

[tilt]: https://tilt.dev
[cluster_api_tilt]: https://cluster-api.sigs.k8s.io/developer/tilt.html
[Initial-setup-for-development-environment]: ./development.md/#initial-setup-for-development-environment

## Create a kind cluster

First, make sure you have a kind cluster and that your `KUBECONFIG` is set up correctly:

```bash
kind create cluster --name=capi-test
```

This local cluster will be running all the cluster api controllers and become the management cluster which then can be used to spin up workload clusters on AWS.

## Get the source

Get the source for core cluster-api for development with Tilt along with cluster-api-provider-aws.

```bash
cd "$(go env GOPATH)"/src
mkdir sigs.k8s.io
cd sigs.k8s.io/
git clone git@github.com:kubernetes-sigs/cluster-api.git
cd cluster-api
```

## Create a tilt-settings.json file

Next, create a `tilt-settings.json` file and place it in your local copy of `cluster-api`. Here is an example:

**Example `tilt-settings.json` for CAPA clusters:**

```json
{
  "enable_providers": [
    "kubeadm-bootstrap", 
    "kubeadm-control-plane", 
    "aws"
  ],
  "default_registry": "gcr.io/your-project-name-here",
  "provider_repos": [
    "/Users/username/go/src/sigs.k8s.io/cluster-api-provider-aws"
  ],
  "kustomize_substitutions": {
    "EXP_CLUSTER_RESOURCE_SET": "true",
    "EXP_MACHINE_POOL": "true",
    "EXP_MACHINE_POOL_MACHINES": "true",
    "EVENT_BRIDGE_INSTANCE_STATE": "true",
    "AWS_B64ENCODED_CREDENTIALS": "YOUR_B64_CREDENTIALS", // paste output from "clusterawsadm bootstrap credentials encode-as-profile".
    "EXP_EKS_FARGATE": "false",
    "CAPA_EKS_IAM": "false",
    "CAPA_EKS_ADD_ROLES": "false",
    "EXP_BOOTSTRAP_FORMAT_IGNITION": "true"
  },
  "extra_args": {
    "aws": ["--v=2"]
  }
}
```

**Example `tilt-settings.json` for EKS managed clusters prior to CAPA v0.7.0:**

```json
{
  "default_registry": "gcr.io/your-project-name-here",
  "provider_repos": ["../cluster-api-provider-aws"],
  "enable_providers": [
    "eks-bootstrap",
    "eks-controlplane",
    "kubeadm-bootstrap",
    "kubeadm-control-plane",
    "aws"
  ],
  "kustomize_substitutions": {
    "AWS_B64ENCODED_CREDENTIALS": "YOUR_B64_CREDENTIALS", // paste output from "clusterawsadm bootstrap credentials encode-as-profile".
    "EXP_EKS": "true",
    "EXP_EKS_IAM": "true",
    "EXP_MACHINE_POOL": "true"
  },
  "extra_args": {
    "aws": ["--v=2"],
    "eks-bootstrap": ["--v=2"],
    "eks-controlplane": ["--v=2"]
  }
}
```

If you use Kind with Docker Desktop and see an error such as `ERROR: failed to load image` or `command "docker exec"` failures, image preloading may be failing for multi-arch images (View the [known kind issue here](https://github.com/kubernetes-sigs/cluster-api/blob/main/hack/tools/internal/tilt-prepare/main.go#L545-L549)). Add `"preload_images": false` at the top level of your `tilt-settings.json`. Tilt will then pull images at runtime instead of preloading them into the Kind cluster.

### Debugging

If you would like to debug CAPA (or core CAPI / another provider) you can run the provider with delve. This will then allow you to attach to delve and debug.

To do this you need to use the **debug** configuration in **tilt-settings.json**. Full details of the options can be seen [here](https://cluster-api.sigs.k8s.io/developer/tilt.html).

An example **tilt-settings.json**:

```json
{
  "enable_providers": [
    "kubeadm-bootstrap",
    "kubeadm-control-plane",
    "aws"
  ],
  "default_registry": "gcr.io/your-project-name-here",
  "provider_repos": [
    "/Users/username/go/src/sigs.k8s.io/cluster-api-provider-aws/v2"
  ],
  "kustomize_substitutions": {
    "EXP_CLUSTER_RESOURCE_SET": "true",
    "EXP_MACHINE_POOL": "true",
    "EVENT_BRIDGE_INSTANCE_STATE": "true",
    "AWS_B64ENCODED_CREDENTIALS": "YOUR_B64_CREDENTIALS", // paste output from "clusterawsadm bootstrap credentials encode-as-profile".
    "EXP_EKS_FARGATE": "false",
    "CAPA_EKS_IAM": "false",
    "CAPA_EKS_ADD_ROLES": "false"
  },
  "extra_args": {
    "aws": ["--v=2"]
  }
  "debug": {
    "aws": {
      "continue": true,
      "port": 30000
    }
  }
}
```

Once you have run tilt (see section below) you will be able to connect to the running instance of delve.

For vscode, you can use the a launch configuration like this:

```json
{
  "name": "Connect to CAPA",
  "type": "go",
  "request": "attach",
  "mode": "remote",
  "remotePath": "",
  "port": 30000,
  "host": "127.0.0.1",
  "showLog": true,
  "trace": "log",
  "logOutput": "rpc"
}
```

For GoLand/IntelliJ add a new run configuration following [these instructions](https://www.jetbrains.com/help/go/attach-to-running-go-processes-with-debugger.html#step-3-create-the-remote-run-debug-configuration-on-the-client-computer).

Or you could use delve directly from the CLI using a command similar to this:

```bash
dlv-dap connect 127.0.0.1:3000
```

## Run Tilt!

To launch your development environment, run:

```bash
tilt up
```

kind cluster becomes a management cluster after this point, check the pods running on the kind cluster `kubectl get pods -A`.

## Create workload clusters

Set the following variables for both CAPA and EKS managed clusters:

```bash
export AWS_SSH_KEY_NAME=<sshkeypair>
export KUBERNETES_VERSION=v1.X.Y  # Use a patch version with published CAPA AMIs (e.g., v1.30.1). Avoid v1.x.0; run `clusterawsadm ami list` to see available versions.
export CLUSTER_NAME=capi-<test-clustename>
export CONTROL_PLANE_MACHINE_COUNT=1
export AWS_CONTROL_PLANE_MACHINE_TYPE=t3.large
export WORKER_MACHINE_COUNT=1
export AWS_NODE_MACHINE_TYPE=t3.large
```

Set the following variables for only EKS managed clusters:

```bash
export AWS_EKS_ROLE_ARN=arn:aws:iam::<accountid>:role/aws-service-role/eks.amazonaws.com/AWSServiceRoleForAmazonEKS
export EKS_KUBERNETES_VERSION=v1.X  # Use any EKS-supported version, including first releases (x.x.0) (e.g., v1.31.0).
```

**Create CAPA managed workload cluster:**

```bash
cat templates/cluster-template.yaml
cat templates/cluster-template.yaml | $HOME/go/bin/envsubst > test-cluster.yaml
kubectl apply -f test-cluster.yaml
```

**Create EKS workload cluster:**

```bash
cat templates/cluster-template-eks.yaml
cat templates/cluster-template-eks.yaml | $HOME/go/bin/envsubst > test-cluster.yaml
kubectl apply -f test-cluster.yaml
```

Check the tilt logs and wait for the clusters to be created.

## Clean up

Before deleting the kind cluster, make sure you delete all the workload clusters.

```bash
kubectl delete cluster <clustername>
tilt up (ctrl-c)
kind delete cluster
```

## Troubleshooting

- Ensure you have at least three Elastic IPs (EIPs) available in your AWS region. CAPA creates one NAT Gateway per availability zone (typically 3), and each NAT Gateway requires an EIP. 
- If your git starts throwing this error

```bash
flag provided but not defined: -variables
Usage: envsubst [options...] <input>
```

you might need to reinstall the system `envsubst`

```bash
brew install gettetxt
# or
brew reinstall gettext
```

Make sure you specify which `envsubst` you are using
