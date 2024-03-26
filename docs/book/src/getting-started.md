# Getting Started

In this tutorial we'll cover the basics of how to use Cluster API Provider AWS to create one or more Kubernetes clusters.

## Installation

### Common Prerequisites

- Install and setup [kubectl] in your local environment
- Install [kind] and [Docker]

### Install and/or configure a Kubernetes cluster

Cluster API requires an existing Kubernetes cluster accessible via kubectl. During the installation process the
Kubernetes cluster will be transformed into a [management cluster] by installing the Cluster API [provider components], so it
is recommended to keep it separated from any application workload.

It is a common practice to create a temporary, local bootstrap cluster which is then used to provision
a target [management cluster] on the selected [infrastructure provider].

**Choose one of the options below:**

1. **Existing Management Cluster**

   For production use-cases a "real" Kubernetes cluster should be used with appropriate backup and disaster recovery policies and procedures in place. The Kubernetes cluster must be at least v1.20.0.

   ```bash
   export KUBECONFIG=<...>
   ```
**OR**

2. **Kind**

   <aside class="note warning">

   <h1>Warning</h1>

   [kind] is not designed for production use.

   **Minimum [kind] supported version**: v0.17.0

   **Help with common issues can be found in the [Troubleshooting Guide](./topics/troubleshooting.md).**

   </aside>

   [kind] can be used for creating a local Kubernetes cluster for development environments or for
   the creation of a temporary [bootstrap cluster] used to provision a target [management cluster] on the selected infrastructure provider.

   The installation procedure depends on the version of kind; if you are planning to use the Docker infrastructure provider,
   please follow the additional instructions in the dedicated tab:

   {{#tabs name:"install-kind" tabs:"Default,Docker,KubeVirt"}}
   {{#tab Default}}

   Create the kind cluster:
   ```bash
   kind create cluster
   ```
   Test to ensure the local kind cluster is ready:
   ```bash
   kubectl cluster-info
   ```

   Then follow the instruction for your kind version using  `kind create cluster --config kind-cluster-with-extramounts.yaml`
   to create the management cluster using the above file.

   {{#/tab }}
   {{#tab Docker}}

   Run the following command to create a kind config file for allowing the Docker provider to access Docker on the host:

   ```bash
   cat > kind-cluster-with-extramounts.yaml <<EOF
   kind: Cluster
   apiVersion: kind.x-k8s.io/v1alpha4
   nodes:
   - role: control-plane
     extraMounts:
       - hostPath: /var/run/docker.sock
         containerPath: /var/run/docker.sock
   EOF
   ```

   Then follow the instruction for your kind version using  `kind create cluster --config kind-cluster-with-extramounts.yaml`
   to create the management cluster using the above file.

   {{#/tab }}
   {{#tab KubeVirt}}

   #### Create the Kind Cluster
   [KubeVirt][KubeVirt] is a cloud native virtualization solution. The virtual machines we're going to create and use for
   the workload cluster's nodes, are actually running within pods in the management cluster. In order to communicate with
   the workload cluster's API server, we'll need to expose it. We are using Kind which is a limited environment. The
   easiest way to expose the workload cluster's API server (a pod within a node running in a VM that is itself running
   within a pod in the management cluster, that is running inside a docker container), is to use a LoadBalancer service.

   To allow using a LoadBalancer service, we can't use the kind's default CNI (kindnet), but we'll need to install
   another CNI, like Calico. In order to do that, we'll need first to initiate the kind cluster with two modifications:
   1. Disable the default CNI
   2. Add the docker credentials to the cluster, to avoid the docker hub pull rate limit of the calico images; read more
      about it in the [docker documentation](https://docs.docker.com/docker-hub/download-rate-limit/), and in the
      [kind documentation](https://kind.sigs.k8s.io/docs/user/private-registries/#mount-a-config-file-to-each-node).

   Create a configuration file for kind. Please notice the docker config file path, and adjust it to your local setting:
   ```bash
   cat <<EOF > kind-config.yaml
   kind: Cluster
   apiVersion: kind.x-k8s.io/v1alpha4
   networking:
   # the default CNI will not be installed
     disableDefaultCNI: true
   nodes:
   - role: control-plane
     extraMounts:
      - containerPath: /var/lib/kubelet/config.json
        hostPath: <YOUR DOCKER CONFIG FILE PATH>
   EOF
   ```
   Now, create the kind cluster with the configuration file:
   ```bash
   kind create cluster --config=kind-config.yaml
   ```
   Test to ensure the local kind cluster is ready:
   ```bash
   kubectl cluster-info
   ```

   #### Install the Calico CNI
   Now we'll need to install a CNI. In this example, we're using calico, but other CNIs should work as well. Please see
   [calico installation guide](https://projectcalico.docs.tigera.io/getting-started/kubernetes/self-managed-onprem/onpremises#install-calico)
   for more details (use the "Manifest" tab). Below is an example of how to install calico version v3.24.4.

   Use the Calico manifest to create the required resources; e.g.:
   ```bash
   kubectl create -f  https://raw.githubusercontent.com/projectcalico/calico/v3.24.4/manifests/calico.yaml
   ```

   {{#/tab }}
   {{#/tabs }}

### Install clusterctl
The clusterctl CLI tool handles the lifecycle of a Cluster API management cluster.

{{#tabs name:"install-clusterctl" tabs:"Linux,macOS,homebrew,Windows"}}
{{#tab Linux}}

#### Install clusterctl binary with curl on Linux
If you are unsure you can determine your computers architecture by running `uname -a`

Download for AMD64:
```bash
curl -L {{#releaselink gomodule:"sigs.k8s.io/cluster-api" asset:"clusterctl-linux-amd64" version:"1.3.x"}} -o clusterctl
```

Download for ARM64:
```bash
curl -L {{#releaselink gomodule:"sigs.k8s.io/cluster-api" asset:"clusterctl-linux-arm64" version:"1.3.x"}} -o clusterctl
```

Download for PPC64LE:
```bash
curl -L {{#releaselink gomodule:"sigs.k8s.io/cluster-api" asset:"clusterctl-linux-ppc64le" version:"1.3.x"}} -o clusterctl
```

Install clusterctl:
```bash
sudo install -o root -g root -m 0755 clusterctl /usr/local/bin/clusterctl
```
Test to ensure the version you installed is up-to-date:
```bash
clusterctl version
```

{{#/tab }}
{{#tab macOS}}

#### Install clusterctl binary with curl on macOS
If you are unsure you can determine your computers architecture by running `uname -a`

Download for AMD64:
```bash
curl -L {{#releaselink gomodule:"sigs.k8s.io/cluster-api" asset:"clusterctl-darwin-amd64" version:"1.3.x"}} -o clusterctl
```

Download for M1 CPU ("Apple Silicon") / ARM64:
```bash
curl -L {{#releaselink gomodule:"sigs.k8s.io/cluster-api" asset:"clusterctl-darwin-arm64" version:"1.3.x"}} -o clusterctl
```

Make the clusterctl binary executable.
```bash
chmod +x ./clusterctl
```
Move the binary in to your PATH.
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
```powershell
curl.exe -L {{#releaselink gomodule:"sigs.k8s.io/cluster-api" asset:"clusterctl-windows-amd64.exe" version:"1.3.x"}} -o clusterctl.exe
```
Append or prepend the path of that directory to the `PATH` environment variable.

Test to ensure the version you installed is up-to-date:
```powershell
clusterctl.exe version
```

{{#/tab }}
{{#/tabs }}

### Initialize the management cluster

Now that we've got clusterctl installed and all the prerequisites in place, let's transform the Kubernetes cluster
into a management cluster by using `clusterctl init`.

The command accepts as input a list of providers to install; when executed for the first time, `clusterctl init`
automatically adds to the list the `cluster-api` core provider, and if unspecified, it also adds the `kubeadm` bootstrap
and `kubeadm` control-plane providers.

#### Enabling Feature Gates

Feature gates can be enabled by exporting environment variables before executing `clusterctl init`.
For example, the `ClusterTopology` feature, which is required to enable support for managed topologies and ClusterClass,
can be enabled via:
```bash
export CLUSTER_TOPOLOGY=true
```
Additional documentation about experimental features can be found in [Experimental Features].

#### Initialization for AWS provider

Depending on the infrastructure provider you are planning to use, some additional prerequisites should be satisfied
before getting started with Cluster API. See below for the expected settings for AWS provider.

Download the latest binary of `clusterawsadm` from the [AWS provider releases]. The [clusterawsadm] command line utility assists with identity and access management (IAM) for [Cluster API Provider AWS][capa].

{{#tabs name:"install-clusterawsadm" tabs:"Linux,macOS,homebrew,Windows"}}
{{#tab Linux}}

Download the latest release; on Linux, type:
```
curl -L {{#releaselink gomodule:"sigs.k8s.io/cluster-api-provider-aws" asset:"clusterawsadm-linux-amd64" version:">=2.0.0"}} -o clusterawsadm
```

Make it executable
```
chmod +x clusterawsadm
```

Move the binary to a directory present in your PATH
```
sudo mv clusterawsadm /usr/local/bin
```

Check version to confirm installation
```
clusterawsadm version
```

**Example Usage**
```bash
export AWS_REGION=us-east-1 # This is used to help encode your environment variables
export AWS_ACCESS_KEY_ID=<your-access-key>
export AWS_SECRET_ACCESS_KEY=<your-secret-access-key>
export AWS_SESSION_TOKEN=<session-token> # If you are using Multi-Factor Auth.

# The clusterawsadm utility takes the credentials that you set as environment
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

Download the latest release; on macOs, type:
```
curl -L {{#releaselink gomodule:"sigs.k8s.io/cluster-api-provider-aws" asset:"clusterawsadm-darwin-amd64" version:">=2.0.0"}} -o clusterawsadm
```

Or if your Mac has an M1 CPU (”Apple Silicon”):
```
curl -L {{#releaselink gomodule:"sigs.k8s.io/cluster-api-provider-aws" asset:"clusterawsadm-darwin-arm64" version:">=2.0.0"}} -o clusterawsadm
```

Make it executable
```
chmod +x clusterawsadm
```

Move the binary to a directory present in your PATH
```
sudo mv clusterawsadm /usr/local/bin
```

Check version to confirm installation
```
clusterawsadm version
```

**Example Usage**
```bash
export AWS_REGION=us-east-1 # This is used to help encode your environment variables
export AWS_ACCESS_KEY_ID=<your-access-key>
export AWS_SECRET_ACCESS_KEY=<your-secret-access-key>
export AWS_SESSION_TOKEN=<session-token> # If you are using Multi-Factor Auth.

# The clusterawsadm utility takes the credentials that you set as environment
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

Check version to confirm installation
```
clusterawsadm version
```

**Example Usage**
```bash
export AWS_REGION=us-east-1 # This is used to help encode your environment variables
export AWS_ACCESS_KEY_ID=<your-access-key>
export AWS_SECRET_ACCESS_KEY=<your-secret-access-key>
export AWS_SESSION_TOKEN=<session-token> # If you are using Multi-Factor Auth.

# The clusterawsadm utility takes the credentials that you set as environment
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
curl.exe -L {{#releaselink gomodule:"sigs.k8s.io/cluster-api-provider-aws" asset:"clusterawsadm-windows-amd64" version:">=2.0.0"}} -o clusterawsadm.exe
```

Append or prepend the path of that directory to the `PATH` environment variable.
Check version to confirm installation
```
clusterawsadm.exe version
```

**Example Usage in Powershell**
```bash
$Env:AWS_REGION="us-east-1" # This is used to help encode your environment variables
$Env:AWS_ACCESS_KEY_ID="<your-access-key>"
$Env:AWS_SECRET_ACCESS_KEY="<your-secret-access-key>"
$Env:AWS_SESSION_TOKEN="<session-token>" # If you are using Multi-Factor Auth.

# The clusterawsadm utility takes the credentials that you set as environment
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

The output of `clusterctl init --infrastructure aws` is similar to this:

```bash
Fetching providers
Installing cert-manager Version="v1.11.0"
Waiting for cert-manager to be available...
Installing Provider="cluster-api" Version="v1.3.3" TargetNamespace="capi-system"
Installing Provider="bootstrap-kubeadm" Version="v1.3.3" TargetNamespace="capi-kubeadm-bootstrap-system"
Installing Provider="control-plane-kubeadm" Version="v1.3.3" TargetNamespace="capi-kubeadm-control-plane-system"
Installing Provider="infrastructure-aws" Version="v2.0.2" TargetNamespace="capa-system"

Your management cluster has been initialized successfully!

You can now create your first workload cluster by running the following:

  clusterctl generate cluster [name] --kubernetes-version [version] | kubectl apply -f -

```

<aside class="note">

<h1>Alternatives to environment variables</h1>

Throughout this quickstart guide we've given instructions on setting parameters using environment variables. For most
environment variables in the rest of the guide, you can also set them in ~/.cluster-api/clusterctl.yaml

See [`clusterctl init`](https://cluster-api.sigs.k8s.io/clusterctl/commands/init.html) for more details.

</aside>

### Create your first workload cluster

Once the management cluster is ready, you can create your first workload cluster.

#### Preparing the workload cluster configuration

The `clusterctl generate cluster` command returns a YAML template for creating a [workload cluster].

<aside class="note">

<h1> Which provider will be used for my cluster? </h1>

The `clusterctl generate cluster` command uses smart defaults in order to simplify the user experience; for example,
if only the `aws` infrastructure provider is deployed, it detects and uses that when creating the cluster.

</aside>

<aside class="note">

<h1> What topology will be used for my cluster? </h1>

The `clusterctl generate cluster` command by default uses cluster templates which are provided by the infrastructure
providers. See the provider's documentation for more information.

See the `clusterctl generate cluster` [command][clusterctl generate cluster] documentation for
details about how to use alternative sources. for cluster templates.

</aside>

#### Required configuration for AWS provider

Depending on the infrastructure provider you are planning to use, some additional prerequisites should be satisfied
before configuring a cluster with Cluster API. Instructions are provided for common providers below.

Otherwise, you can look at the `clusterctl generate cluster` [command][clusterctl generate cluster] documentation for details about how to
discover the list of variables required by a cluster templates.

```bash
export AWS_REGION=us-east-1
export AWS_SSH_KEY_NAME=default
# Select instance types
export AWS_CONTROL_PLANE_MACHINE_TYPE=t3.large
export AWS_NODE_MACHINE_TYPE=t3.large
```

See the [AWS provider prerequisites] document for more details.


#### Generating the cluster configuration

For the purpose of this tutorial, we'll name our cluster capi-quickstart.

```bash
clusterctl generate cluster capi-quickstart \
  --kubernetes-version v1.26.0 \
  --control-plane-machine-count=3 \
  --worker-machine-count=3 \
  > capi-quickstart.yaml
```

This creates a YAML file named `capi-quickstart.yaml` with a predefined list of Cluster API objects; Cluster, Machines,
Machine Deployments, etc.

The file can be eventually modified using your editor of choice.

See [clusterctl generate cluster] for more details.

#### Apply the workload cluster

When ready, run the following command to apply the cluster manifest.

```bash
kubectl apply -f capi-quickstart.yaml
```

The output is similar to this:

```bash
cluster.cluster.x-k8s.io/capi-quickstart created
dockercluster.infrastructure.cluster.x-k8s.io/capi-quickstart created
kubeadmcontrolplane.controlplane.cluster.x-k8s.io/capi-quickstart-control-plane created
dockermachinetemplate.infrastructure.cluster.x-k8s.io/capi-quickstart-control-plane created
machinedeployment.cluster.x-k8s.io/capi-quickstart-md-0 created
dockermachinetemplate.infrastructure.cluster.x-k8s.io/capi-quickstart-md-0 created
kubeadmconfigtemplate.bootstrap.cluster.x-k8s.io/capi-quickstart-md-0 created
```

#### Accessing the workload cluster

The cluster will now start provisioning. You can check status with:

```bash
kubectl get cluster
```

You can also get an "at glance" view of the cluster and its resources by running:

```bash
clusterctl describe cluster capi-quickstart
```

To verify the first control plane is up:

```bash
kubectl get kubeadmcontrolplane
```

You should see an output is similar to this:

```bash
NAME                    CLUSTER           INITIALIZED   API SERVER AVAILABLE   REPLICAS   READY   UPDATED   UNAVAILABLE   AGE    VERSION
capi-quickstart-g2trk   capi-quickstart   true                                 3                  3         3             4m7s   v1.26.0
```

<aside class="note warning">

<h1> Warning </h1>

The control plane won't be `Ready` until we install a CNI in the next step.

</aside>

After the first control plane node is up and running, we can retrieve the [workload cluster] Kubeconfig.

```bash
clusterctl get kubeconfig capi-quickstart > capi-quickstart.kubeconfig
```


### Deploy a CNI solution

Calico is used here as an example.

```bash
kubectl --kubeconfig=./capi-quickstart.kubeconfig \
  apply -f https://raw.githubusercontent.com/projectcalico/calico/v3.24.1/manifests/calico.yaml
```

After a short while, our nodes should be running and in `Ready` state,
let's check the status using `kubectl get nodes`:

```bash
kubectl --kubeconfig=./capi-quickstart.kubeconfig get nodes
```

```bash
NAME                                          STATUS   ROLES           AGE   VERSION
capi-quickstart-g2trk-9xrjv                   Ready    control-plane   12m   v1.26.0
capi-quickstart-g2trk-bmm9v                   Ready    control-plane   11m   v1.26.0
capi-quickstart-g2trk-hvs9q                   Ready    control-plane   13m   v1.26.0
capi-quickstart-md-0-55x6t-5649968bd7-8tq9v   Ready    <none>          12m   v1.26.0
capi-quickstart-md-0-55x6t-5649968bd7-glnjd   Ready    <none>          12m   v1.26.0
capi-quickstart-md-0-55x6t-5649968bd7-sfzp6   Ready    <none>          12m   v1.26.0
```

### Clean Up

Delete workload cluster.
```bash
kubectl delete cluster capi-quickstart
```
<aside class="note warning">

IMPORTANT: In order to ensure a proper cleanup of your infrastructure you must always delete the cluster object. Deleting the entire cluster template with `kubectl delete -f capi-quickstart.yaml` might lead to pending resources to be cleaned up manually.
</aside>

Delete management cluster
```bash
kind delete cluster
```

## Next steps

See the [clusterctl] documentation for more detail about clusterctl supported actions.

<!-- links -->
[Experimental Features]: https://cluster-api.sigs.k8s.io/tasks/experimental-features/experimental-features.html
[AWS provider prerequisites]: https://cluster-api-aws.sigs.k8s.io/topics/using-clusterawsadm-to-fulfill-prerequisites.html
[AWS provider releases]: https://github.com/kubernetes-sigs/cluster-api-provider-aws/releases
[bootstrap cluster]: ./topics/reference/glossary.md
[capa]: https://cluster-api-aws.sigs.k8s.io
[clusterawsadm]: https://cluster-api-aws.sigs.k8s.io/clusterawsadm/clusterawsadm.html
[clusterctl generate cluster]: https://cluster-api.sigs.k8s.io/clusterctl/commands/generate-cluster.html
[clusterctl]: https://cluster-api.sigs.k8s.io/clusterctl/overview.html
[Docker]: https://www.docker.com/
[infrastructure provider]: ./topics/reference/glossary.md
[kind]: https://kind.sigs.k8s.io/
[kubectl]: https://kubernetes.io/docs/tasks/tools/install-kubectl/
[management cluster]: ./topics/reference/glossary.md
[KubeVirt]: https://kubevirt.io/
[provider components]: ./topics/reference/glossary.md
[workload cluster]: ./topics/reference/glossary.md
