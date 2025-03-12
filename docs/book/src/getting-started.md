# Getting Started

This guide covers the basics of using Cluster API (CAPI) to create one or more Kubernetes clusters on Amazon Web Services (AWS).
This guide is **specific for the Cluster API AWS provider** (CAPA) for clarity and reliability. For a more comprehensive guide covering other providers, you can also refer to the generic [Cluster API book](https://cluster-api.sigs.k8s.io/user/quick-start.html). 


## Installation

This guide describes a path that uses the `clusterctl` CLI tool to handle the lifecycle of a Cluster API [management cluster](https://cluster-api.sigs.k8s.io/reference/glossary#management-cluster). The Cluster API Operator can also be used for that purpose, but we'll focus on the CLI tools in this guide.

The clusterctl command-line interface is designed to provide a simple experience and a quick start with Cluster API. It automates fetching the YAML files defining [provider components](https://cluster-api.sigs.k8s.io/reference/glossary#provider-components) and installing them.
Additionally, it encodes a set of best practices for managing providers that help the user avoid misconfigurations or manage subsequent operations, such as upgrades.

The Cluster API Operator is a Kubernetes Operator built on top of clusterctl and designed to empower cluster administrators to handle the lifecycle of Cluster API providers within a management cluster using a declarative approach. It aims to improve user experience in deploying and managing Cluster API, making it easier to handle day-to-day tasks and automate workflows with GitOps. Visit the [CAPI Operator quickstart] for more information about this tool.

### Common Prerequisites

Install and setup:
- [kubectl] in your local environment or cloud devcontainer
- [kind] and [Docker]
- [Helm]

The [CAPA repository](https://github.com/kubernetes-sigs/cluster-api-provider-aws) includes configuration for [devbox](https://www.jetify.com/devbox/) and [devcontainers](https://containers.dev) so that it can be built and executed in any development environment, local or cloud, such as [GitHub Codespaces](https://github.com/codespaces/new?hide_repo_select=true&ref=main&repo=141609301&skip_quickstart=true&machine=standardLinux32gb&devcontainer_path=.devcontainer%2Fcapa-devbox%2Fdevcontainer.json&geo=UsEast). 

### Setup a Kubernetes Cluster for Management

Cluster API requires a Kubernetes cluster ready and accessible via `kubectl`. During the installation process, this
Installing the Cluster API [provider components] will configure the Kubernetes cluster as a [management cluster]. It is recommended that this cluster be used exclusively to manage other clusters, isolating it from any application workload.

It is a common practice to create a temporary, local bootstrap cluster, which is then used to provision
a target [management cluster] for the selected [infrastructure provider], especially for development. 
A reliable and secure Kubernetes cluster should be used for production use cases, with appropriate backup and disaster recovery policies and procedures in place. The management Kubernetes cluster must be at least version v1.20.0.

In this guide, we´ll use [kind] (v0.25.0 or more recent) to provision a management cluster.

**Help with common issues can be found in the [Troubleshooting Guide](./troubleshooting.md).**

<aside class="note warning">

<h1>Warning</h1>

[kind] is not recommended for production use. It is targeted at local development, testing, and CI/CD and has limitations in reliability, security, and performance.

</aside>

[kind] can be used to create a local Kubernetes cluster for development environments or a temporary [bootstrap cluster] to provision a target [management cluster] on the selected infrastructure provider. 
*The versions indicated in this guide are tested and known to work as documented, even though newer versions are available.*

Create the kind cluster at the desired Kubernetes version:
```bash
export KUBERNETES_VERSION="v1.30.8"
kind create cluster --image kindest/node:$KUBERNETES_VERSION
```

Test to ensure the local kind cluster is ready:
```bash
kubectl cluster-info
```

### Install clusterctl
The clusterctl CLI tool handles the lifecycle of a Cluster API management cluster.

{{#tabs name:"install-clusterctl" tabs:"Linux,macOS,homebrew,Windows"}}
{{#tab Linux}}

#### Install clusterctl binary with curl on Linux
If you are unsure, you can determine your computer's architecture by running `uname -a`

Download for AMD64:
```bash
curl -L {{#releaselink repo:"https://github.com/kubernetes-sigs/cluster-api" gomodule:"sigs.k8s.io/cluster-API" asset:"clusterctl-linux-amd64" version:"1.9.x"}} -o clusterctl
```

Download for ARM64:
```bash
curl -L {{#releaselink repo:"https://github.com/kubernetes-sigs/cluster-api" gomodule:"sigs.k8s.io/cluster-API" asset:"clusterctl-linux-arm64" version:"1.9.x"}} -o clusterctl
```

Download for PPC64LE:
```bash
curl -L {{#releaselink repo:"https://github.com/kubernetes-sigs/cluster-api" gomodule:"sigs.k8s.io/cluster-API" asset:"clusterctl-linux-ppc64le" version:"1.9.x"}} -o clusterctl
```

Alternativaly, using uname to detect architecture and os:
```bash
LOCAL_OS=$(uname -s | tr '[:upper:]' '[:lower:]')
LOCAL_ARCH=$(uname -m)

if [ "$LOCAL_ARCH" == "x86_64" ]; then
  LOCAL_ARCH="amd64"
elif [ "$LOCAL_ARCH" == "aarch64" ]; then
  LOCAL_ARCH="arm64"
fi
echo $LOCAL_OS $LOCAL_ARCH

CLUSTERCTL_VERSION=$(grep 'sigs.k8s.io/cluster-api v' go.mod | awk '{print $2}')
echo "CLUSTERCTL_VERSION=$CLUSTERCTL_VERSION"
curl -L "https://github.com/kubernetes-sigs/cluster-api/releases/download/$CLUSTERCTL_VERSION/clusterctl-$LOCAL_OS-$LOCAL_ARCH" -o clusterctl

```

Install clusterctl for the current user:
```bash
install -D -m 0755 clusterctl "$HOME/.local/bin/clusterctl"
```

Alternatively, for all users:
```bash
sudo install -o root -g root -m 0755 clusterctl /usr/local/bin/clusterctl
```

Verify that the version you installed is up-to-date:
```bash
clusterctl version
```

{{#/tab }}
{{#tab macOS}}

#### Install clusterctl binary with curl on macOS
If you are unsure, you can determine your computer's architecture by running `uname -a`

Download for AMD64:
```bash
curl -L {{#releaselink repo:"https://github.com/kubernetes-sigs/cluster-api" gomodule:"sigs.k8s.io/cluster-API" asset:"clusterctl-darwin-amd64" version:"1.9.x"}} -o clusterctl
```

Download for M1 CPU ("Apple Silicon") / ARM64:
```bash
curl -L {{#releaselink repo:"https://github.com/kubernetes-sigs/cluster-api" gomodule:"sigs.k8s.io/cluster-API" asset:"clusterctl-darwin-arm64" version:"1.9.x"}} -o clusterctl
```

Make the clusterctl binary executable.
```bash
chmod +x ./clusterctl
```
Move the binary into your PATH.
```bash
sudo mv ./clusterctl /usr/local/bin/clusterctl
```
Test to ensure the version you installed is up-to-date:
```bash
clusterctl version
```
{{#/tab }}
{{#tab homebrew}}

#### Install clusterctl with homebrew on macOS and Linux

Install the latest release using homebrew:

```bash
brew install clusterctl
```

Test to ensure the version you installed is up-to-date:
```bash
clusterctl version
```

{{#/tab }}
{{#tab windows}}

#### Install clusterctl binary with curl on Windows using PowerShell
Go to the working directory where you want clusterctl downloaded.

Download the latest release; on Windows, type:
```PowerShell
curl.exe -L {{#releaselink repo:"https://github.com/kubernetes-sigs/cluster-api" gomodule:"sigs.k8s.io/cluster-API" asset:"clusterctl-windows-amd64.exe" version:"1.9.x"}} -o clusterctl.exe
```
Append or prepend the path of that directory to the `PATH` environment variable.

Test to ensure the version you installed is up-to-date:
```powershell
clusterctl.exe version
```

{{#/tab }}
{{#/tabs }}

### Initialize the management cluster

Now that we've got clusterctl installed and all the prerequisites in place let's initialize the Kubernetes cluster
as a management cluster by using the `clusterctl init` command.

The command accepts as input a list of providers to install; when executed for the first time, `clusterctl init`
automatically adds to the list the `cluster-api` core provider, and if unspecified, it also adds the `kubeadm` bootstrap
and `kubeadm` control-plane providers.

Download the latest binary of `clusterawsadm` from the [AWS provider releases]. The [clusterawsadm] command line utility assists with identity and access management (IAM) for [Cluster API Provider AWS][capa].

{{#tabs name:"install-clusterawsadm" tabs:"Linux,macOS,homebrew,Windows"}}
{{#tab Linux}}

Download the latest release; on Linux, type:
```
curl -L {{#releaselink repo:"https://github.com/kubernetes-sigs/cluster-api-provider-aws" gomodule:"sigs.k8s.io/cluster-api-provider-aws" asset:"clusterawsadm-linux-amd64" version:">=2.0.0"}} -o clusterawsadm
```

Alternatively:
```bash
CLUSTERAWSADM_VERSION=v2.7.1
curl -L "https://github.com/kubernetes-sigs/cluster-api-provider-aws/releases/download/$CLUSTERAWSADM_VERSION/clusterawsadm-$LOCAL_OS-$LOCAL_ARCH" -o clusterawsadm
```

Install clusterawsadm for the current user.
```bash
install -D -m 0755 clusterawsadm "$HOME/.local/bin/clusterawsadm"
```

Alternatively, system-wide:
```bash
sudo install -D -o root -g root -m 0755 clusterawsadm "$HOME/.local/bin/clusterawsadm"
```

Check the version to confirm the installation.
```bash
clusterawsadm version
```

**Example Usage**
```bash
# AWS region is used to help encode your environment variables
export AWS_REGION=us-east-1 

export AWS_ACCESS_KEY_ID=<your-access-key>
export AWS_SECRET_ACCESS_KEY=<your-secret-access-key>
export AWS_SESSION_TOKEN=<session-token> # If you are using Multi-Factor Auth.

# Verify that you are authenticated with the expected AWS account
aws sts get-caller-identity

# The clusterawsadm utility takes the credentials you set as environment
# variables and uses them to create a CloudFormation stack in your AWS account
# with the correct IAM resources.
clusterawsadm bootstrap iam create-cloudformation-stack

# Create the base64 encoded credentials using clusterawsadm.
# This command uses your environment variables and encodes
# them in a value to be stored in a Kubernetes Secret.
export AWS_B64ENCODED_CREDENTIALS=$(clusterawsadm bootstrap credentials encode-as-profile)

# Finally, initialize the management cluster
clusterctl init --infrastructure aws
```

{{#/tab }}
{{#tab macOS}}

Download the latest release; on macOS, type:
```
curl -L {{#releaselink repo:"https://github.com/kubernetes-sigs/cluster-api-provider-aws" gomodule:"sigs.k8s.io/cluster-api-provider-aws" asset:"clusterawsadm-darwin-amd64" version:">=2.0.0"}} -o clusterawsadm
```

Or if your Mac has an M1 CPU (”Apple Silicon”):
```
curl -L {{#releaselink repo:"https://github.com/kubernetes-sigs/cluster-api-provider-aws" gomodule:"sigs.k8s.io/cluster-api-provider-aws" asset:"clusterawsadm-darwin-arm64" version:">=2.0.0"}} -o clusterawsadm
```

Make it executable
```
chmod +x clusterawsadm
```

Move the binary to a directory present in your PATH.
```
sudo mv clusterawsadm /usr/local/bin
```

Check the version to confirm the installation.
```
clusterawsadm version
```

**Example Usage**
```bash
export AWS_REGION=us-east-1 # This is used to help encode your environment variables
export AWS_ACCESS_KEY_ID=<your-access-key>
export AWS_SECRET_ACCESS_KEY=<your-secret-access-key>
export AWS_SESSION_TOKEN=<session-token> # If you are using Multi-Factor Auth.

# The clusterawsadm utility takes the credentials you set as environment
# variables and uses them to create a CloudFormation stack in your AWS account
# with the correct IAM resources.
clusterawsadm bootstrap iam create-cloudformation-stack

# Create the base64 encoded credentials using clusterawsadm.
# This command uses your environment variables and encodes
# them in a value to be stored in a Kubernetes Secret.
export AWS_B64ENCODED_CREDENTIALS=$(clusterawsadm bootstrap credentials encode-as-profile)

# Finally, initialize the management cluster
clusterctl init --infrastructure aws
```
{{#/tab }}
{{#tab homebrew}}

Install the latest release using homebrew:
```
brew install clusterawsadm
```

Check the version to confirm the installation.
```
clusterawsadm version
```

**Example Usage**
```bash
export AWS_REGION=us-east-1 # This is used to help encode your environment variables
export AWS_ACCESS_KEY_ID=<your-access-key>
export AWS_SECRET_ACCESS_KEY=<your-secret-access-key>
export AWS_SESSION_TOKEN=<session-token> # If you are using Multi-Factor Auth.

# The clusterawsadm utility takes the credentials you set as environment
# variables and uses them to create a CloudFormation stack in your AWS account
# with the correct IAM resources.
clusterawsadm bootstrap iam create-cloudformation-stack

# Create the base64 encoded credentials using clusterawsadm.
# This command uses your environment variables and encodes
# them in a value to be stored in a Kubernetes Secret.
export AWS_B64ENCODED_CREDENTIALS=$(clusterawsadm bootstrap credentials encode-as-profile)

# Finally, initialize the management cluster
clusterctl init --infrastructure aws
```

{{#/tab }}
{{#tab Windows}}

Download the latest release; on Windows, type:
```
curl.exe -L {{#releaselink repo:"https://github.com/kubernetes-sigs/cluster-api-provider-aws" gomodule:"sigs.k8s.io/cluster-api-provider-aws" asset:"clusterawsadm-windows-amd64.exe" version:">=2.0.0"}} -o clusterawsadm.exe
```

Append or prepend the path of that directory to the `PATH` environment variable.
Check the version to confirm the installation.
```
clusterawsadm.exe version
```

**Example Usage in Powershell**
```bash
$Env:AWS_REGION="us-east-1" # This is used to help encode your environment variables
$Env:AWS_ACCESS_KEY_ID="<your-access-key>"
$Env:AWS_SECRET_ACCESS_KEY="<your-secret-access-key>"
$Env:AWS_SESSION_TOKEN="<session-token>" # If you are using Multi-Factor Auth.

# The clusterawsadm utility takes the credentials you set as environment
# variables and uses them to create a CloudFormation stack in your AWS account
# with the correct IAM resources.
clusterawsadm bootstrap iam create-cloudformation-stack

# Create the base64 encoded credentials using clusterawsadm.
# This command uses your environment variables and encodes
# them in a value to be stored in a Kubernetes Secret.
$Env:AWS_B64ENCODED_CREDENTIALS=$(clusterawsadm bootstrap credentials encode-as-profile)

# Finally, initialize the management cluster
clusterctl init --infrastructure aws
```
{{#/tab }}
{{#/tabs }}

See the [AWS provider prerequisites] document for more details.


The output of `clusterctl init` is similar to this:

```bash
Fetching providers
Installing cert-manager Version="v1.11.0"
Waiting for the manager to be available...
Installing Provider="cluster-api" Version="v1.0.0" TargetNamespace="capi-system"
Installing Provider="bootstrap-kubeadm" Version="v1.0.0" TargetNamespace="capi-kubeadm-bootstrap-system"
Installing Provider="control-plane-kubeadm" Version="v1.0.0" TargetNamespace="capi-kubeadm-control-plane-system"
Installing Provider="infrastructure-docker" Version="v1.0.0" TargetNamespace="capd-system"

Your management cluster has been initialized successfully!

You can now create your first workload cluster by running the following:

  clusterctl generate cluster [name] --Kubernetes-version [version] | kubectl apply -f -
```

<aside class="note">

<h1>Alternatives to environment variables</h1>

Throughout this quickstart guide, we've given instructions on setting parameters using environment variables. For most
environment variables in the rest of the guide, you can also set them in `$XDG_CONFIG_HOME/cluster-api/clusterctl.yaml`

See [`clusterctl init`](https://cluster-api.sigs.k8s.io/clusterctl/commands/init) for more details.

</aside>

### Create your first workload cluster

Once the management cluster is ready, you can create your first workload cluster.

#### Preparing the workload cluster configuration

The `clusterctl generate cluster` command returns a YAML template for creating a [workload cluster].

<aside class="note">

<h1> Which provider will be used for my cluster? </h1>

The `clusterctl generate cluster` command uses smart defaults to simplify the user experience; for example,
If only the `aws` infrastructure provider is deployed, it will detect and use that when creating the cluster.

</aside>

<aside class="note">

<h1> What topology will be used for my cluster? </h1>

The `clusterctl generate cluster` command defaults to using cluster templates provided by the infrastructure providers. For more information, see the provider's documentation.

See the `clusterctl generate cluster` [command][clusterctl generate cluster] documentation for
details about using alternative sources for cluster templates.

</aside>

#### Required configuration

Depending on the infrastructure provider you plan to use, some additional prerequisites should be satisfied
before configuring a cluster with Cluster API. Instructions are provided for the AWS provider.

Otherwise, you can look at the `clusterctl generate cluster` [command][clusterctl generate cluster] documentation for details about how to
discover the list of variables required by a cluster template.

```bash
# Set AWS SSH keypair name
export AWS_SSH_KEY_NAME=default

# Verify that the keypair exists
aws ec2 describe-key-pairs --query "KeyPairs[?KeyName=='$AWS_SSH_KEY_NAME']" 


# Select instance types
export AWS_CONTROL_PLANE_MACHINE_TYPE=t3.large
export AWS_NODE_MACHINE_TYPE=t3.large
```

See the [AWS provider prerequisites] document for more details.

#### Generating the cluster configuration
The following command generates a file with all yaml manifests for the workload cluster creation:

```bash
clusterctl generate cluster capa-quickstart \
  --kubernetes-version "$KUBERNETES_VERSION" \
  --control-plane-machine-count=3 \
  --worker-machine-count=3 \
 > capa-quickstart.yaml
```

This command creates a YAML file named `capa-quickstart.yaml` with a predefined list of Cluster API objects: Cluster, Machines,
Machine Deployments, etc.

The file can be eventually modified using your editor of choice.

See [clusterctl generate cluster] for more details.

#### Apply the workload cluster

When ready, run the following command to apply the cluster manifest.

```bash
kubectl apply -f capa-quickstart.yaml
```

The output is similar to this:

```bash
cluster.cluster.x-k8s.io/capa-quickstart created
dockercluster.infrastructure.cluster.x-k8s.io/capa-quickstart created
kubeadmcontrolplane.controlplane.cluster.x-k8s.io/capa-quickstart-control-plane created
dockermachinetemplate.infrastructure.cluster.x-k8s.io/capa-quickstart-control-plane created
machinedeployment.cluster.x-k8s.io/capa-quickstart-md-0 created
dockermachinetemplate.infrastructure.cluster.x-k8s.io/capa-quickstart-md-0 created
kubeadmconfigtemplate.bootstrap.cluster.x-k8s.io/capa-quickstart-md-0 created
```

#### Accessing the workload cluster

The cluster will now start provisioning. You can check status with:

```bash
watch -n 15 kubectl get cluster
```

And see an output similar to this:

```bash
NAME              PHASE         AGE   VERSION
capa-quickstart   Provisioned   8s    v1.30.8
```


You can also get an "at a glance" view of the cluster and its resources by running:

```bash
clusterctl describe cluster capa-quickstart
```

To verify the first control plane is up:

```bash
kubectl get kubeadmcontrolplane
```

You should see an output similar to this, with `INITIALIZED = true`:

```bash
NAME                            CLUSTER           INITIALIZED   API SERVER AVAILABLE   REPLICAS   READY   UPDATED   UNAVAILABLE   AGE   VERSION
capa-quickstart-control-plane   capa-quickstart   true                                 1                  1         1             14m   v1.30.8
```

<aside class="note warning">

<h1> Warning </h1>

The control plane won't be `Ready` until we install a CNI in the next step.

</aside>

After the first control plane node is running successfully, we can retrieve the [workload cluster] kubeconfig.

```bash
clusterctl get kubeconfig capa-quickstart > capa-quickstart.kubeconfig
```

### Deploy a CNI solution

Calico is used here as an example.


```bash
kubectl --kubeconfig=./capa-quickstart.kubeconfig \
 apply -f https://raw.githubusercontent.com/projectcalico/calico/v3.26.1/manifests/calico.yaml
```

After a short while, our nodes should be running and in the `Ready` state,
let's check the status using `kubectl get nodes`:

```bash
kubectl --kubeconfig=./capa-quickstart.kubeconfig get nodes
```
```bash
NAME                                          STATUS   ROLES           AGE    VERSION
capa-quickstart-vs89t-gmbld                   Ready    control-plane   5m33s  v1.32.0
capa-quickstart-vs89t-kf9l5                   Ready    control-plane   6m20s  v1.32.0
capa-quickstart-vs89t-t8cfn                   Ready    control-plane   7m10s  v1.32.0
capa-quickstart-md-0-55x6t-5649968bd7-8tq9v   Ready <none>          6m5s   v1.32.0
capa-quickstart-md-0-55x6t-5649968bd7-glnjd   Ready <none>          6m9s   v1.32.0
capa-quickstart-md-0-55x6t-5649968bd7-sfzp6   Ready <none>          6m9s   v1.32.0
```

Your management and workload clusters are now ready to use.

### Troubleshooting

#### Troubleshooting node availability

If the nodes don't become ready after a long period, read the pods in the `kube-system` namespace
```bash
kubectl --kubeconfig=./capa-quickstart.kubeconfig get pod -n kube-system
```

If the Calico pods are in an image pull error state (`ErrImagePull`), it's probably because of the Docker Hub pull rate limit.
We can try to fix that by adding a secret with our Docker Hub credentials and using it;
see [here](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/#registry-secret-existing-credentials)
for details.

First, create the secret. Please note the Docker config file path and adjust it to your local setting.
```bash
kubectl --kubeconfig=./capa-quickstart.kubeconfig create secret generic docker-creds \
 --from-file=.dockerconfigjson=<YOUR DOCKER CONFIG FILE PATH> \
    --type=kubernetes.io/dockerconfigjson \
 -n kube-system
```

Now, if the `calico-node` pods have a status of `ErrImagePull`, patch their DaemonSet to make them use the new secret to pull images:
```bash
kubectl --kubeconfig=./capa-quickstart.kubeconfig patch daemonset \
 -n kube-system calico-node \
    -p '{"spec":{"template":{"spec":{"imagePullSecrets":[{"name":"docker-creds"}]}}}}'
```

After a short while, the calico-node pods will have `Running` status. Now, if the calico-kube-controllers pod is also
in `ErrImagePull` status, patch its deployment to fix the problem:
```bash
kubectl --kubeconfig=./capa-quickstart.kubeconfig patch deployment \
 -n kube-system calico-kube-controllers \
    -p '{"spec":{"template":{"spec":{"imagePullSecrets":[{"name":"docker-creds"}]}}}}'
```

Read the pods again
```bash
kubectl --kubeconfig=./capa-quickstart.kubeconfig get pod -n kube-system
```

Eventually, all the pods in the kube-system namespace will run, and the result should be similar to this:
```text
NAME                                                          READY   STATUS    RESTARTS   AGE
calico-kube-controllers-c969cf844-dgld6                       1/1     Running   0          50s
calico-node-7zz7c                                             1/1     Running   0          54s
calico-node-jmjd6                                             1/1     Running   0          54s
coredns-64897985d-dspjm                                       1/1     Running   0          3m49s
coredns-64897985d-pgtgz                                       1/1     Running   0          3m49s
etcd-capa-quickstart-control-plane-kjjbb                      1/1     Running   0          3m57s
kube-apiserver-capa-quickstart-control-plane-kjjbb            1/1     Running   0          3m57s
kube-controller-manager-capa-quickstart-control-plane-kjjbb   1/1     Running   0          3m57s
kube-proxy-b9g5m                                              1/1     Running   0          3m12s
kube-proxy-p6xx8                                              1/1     Running   0          3m49s
kube-scheduler-capa-quickstart-control-plane-kjjbb            1/1     Running   0          3m57s
```

#### Troubleshooting load balancers
If the default load balancer fails, using a network load balancer instead may solve the issue:
In the `AWSCLuster` section of the manifests, add the `loadBalancerType` attribute like so:
```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
kind: AWSCluster
metadata:
  name: capa-quickstart
  namespace: default
spec:
  region: us-east-1
  sshKeyName: capa-quickstart
  controlPlaneLoadBalancer:
    loadBalancerType: nlb
```
<!-- Until https://github.com/kubernetes-sigs/cluster-api-provider-aws/pull/5345 is published -->

### Clean Up

Delete workload cluster.
```bash
kubectl delete cluster capa-quickstart
```
<aside class="note warning">

IMPORTANT: To ensure a proper cleanup of your infrastructure, you must always delete the cluster object. Deleting the entire cluster template with `kubectl delete—f capa-quickstart.yaml` might result in pending resources that need to be cleaned up manually.
</aside>

Delete management cluster
```bash
kind delete cluster
```

## Next steps

Create a second workload cluster. Follow the steps outlined above, but remember to provide a different name for your second workload cluster.
Deploy applications to your workload cluster. For pointers, use the [CNI deployment steps](#deploy-a-cni-solution).
- See the [clusterctl] documentation for more detail about clusterctl supported actions.

<!-- links -->
[Experimental Features]: https://cluster-api.sigs.k8s.io/tasks/experimental-features/experimental-features
[Akamai (Linode) provider]: https://linode.github.io/cluster-api-provider-linode/introduction.html
[AWS provider prerequisites]: https://cluster-api-aws.sigs.k8s.io/topics/using-clusterawsadm-to-fulfill-prerequisites.html
[AWS provider releases]: https://github.com/kubernetes-sigs/cluster-api-provider-aws/releases
[Azure Provider Prerequisites]: https://capz.sigs.k8s.io/getting-started.html#prerequisites
[bootstrap cluster]: https://cluster-api.sigs.k8s.io/reference/glossary#bootstrap-cluster
[capa]: https://cluster-api-aws.sigs.k8s.io
[capv-upload-images]: https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/blob/master/docs/getting_started.md#uploading-the-machine-images
[clusterawsadm]: https://cluster-api-aws.sigs.k8s.io/clusterawsadm/clusterawsadm.html
[clusterctl generate cluster]: https://cluster-api.sigs.k8s.io/clusterctl/commands/generate-cluster
[clusterctl get kubeconfig]: https://cluster-api.sigs.k8s.io/clusterctl/commands/get-kubeconfig
[clusterctl]: https://cluster-api.sigs.k8s.io/clusterctl/overview
[Docker]: https://www.docker.com/
[GCP provider]: https://cluster-api-gcp.sigs.k8s.io/
[Helm]: https://helm.sh/docs/intro/install/
[Harvester provider]: https://github.com/rancher-sandbox/cluster-api-provider-harvester
[Hetzner provider]: https://github.com/syself/cluster-api-provider-hetzner
[Hivelocity provider]: https://github.com/hivelocity/cluster-api-provider-hivelocity
[Huawei Cloud provider]: https://github.com/HuaweiCloudDeveloper/cluster-api-provider-huawei
[IBM Cloud provider]: https://github.com/kubernetes-sigs/cluster-api-provider-ibmcloud
[infrastructure provider]: https://cluster-api.sigs.k8s.io/reference/glossary#infrastructure-provider
[ionoscloud provider]: https://github.com/ionos-cloud/cluster-api-provider-ionoscloud
[kind]: https://kind.sigs.k8s.io/
[KubeadmControlPlane]: https://cluster-api.sigs.k8s.io/tasks/control-plane/kubeadm-control-plane
[kubectl]: https://kubernetes.io/docs/tasks/tools/install-kubectl/
[management cluster]: https://cluster-api.sigs.k8s.io/reference/glossary#management-cluster
[Metal3 getting started guide]: https://github.com/metal3-io/cluster-api-provider-metal3/blob/master/docs/getting-started.md
[Metal3 provider]: https://github.com/metal3-io/cluster-api-provider-metal3/
[K0smotron provider]: https://github.com/k0smotron/k0smotron
[KubeKey provider]: https://github.com/kubesphere/kubekey
[KubeVirt provider]: https://github.com/kubernetes-sigs/cluster-api-provider-kubevirt/
[KubeVirt]: https://kubevirt.io/
[oci-provider]: https://oracle.github.io/cluster-api-provider-oci/#getting-started
[openstack-resource-controller]: https://k-orc.cloud/
[Equinix Metal getting started guide]: https://github.com/kubernetes-sigs/cluster-api-provider-packet#using
[provider]: https://cluster-api.sigs.k8s.io/reference/glossary#provider
[provider components]: https://cluster-api.sigs.k8s.io/reference/glossary#provider-components
[vSphere getting started guide]: https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/blob/master/docs/getting_started.md
[workload cluster]: https://cluster-api.sigs.k8s.io/reference/glossary#workload-cluster
[CAPI Operator quickstart]: https://cluster-api.sigs.k8s.io/user/quick-start-operator.md
[Proxmox getting started guide]: https://github.com/ionos-cloud/cluster-api-provider-proxmox/blob/main/docs/Usage.md
[Tinkerbell getting started guide]: https://github.com/tinkerbell/cluster-api-provider-tinkerbell/blob/main/docs/QUICK-START.md
[CAPONE Wiki]: https://github.com/OpenNebula/cluster-api-provider-opennebula/wiki
