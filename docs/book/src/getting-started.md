# Getting Started Guide
In this guide we'll cover the basics of how to use [Cluster API](https://github.com/kubernetes-sigs/cluster-api) to create one or more Kubernetes clusters on [Amazon Web Services](https://aws.amazon.com).


## Installation
This guide describes how to use the `clusterctl` CLI tool to handle the lifecycle of a Cluster API [management cluster](https://cluster-api.sigs.k8s.io/reference/glossary#management-cluster). It is also possible to use the [Cluster API Operator](https://cluster-api.sigs.k8s.io/user/quick-start-operator) as an alternative, but we're focused on the `clusterctl`on this guide.

The `clusterctl` command-line interface is designed to simplify cluster management tasks using Cluster API. It automates fetching the YAML files defining [provider components](https://cluster-api.sigs.k8s.io/reference/glossary#provider-components) and installing them. Additionally, it encodes a set of best practices for managing providers, helping the user avoid misconfigurations or issues when managing subsequent operations such as upgrades.  

The Cluster API Operator is a Kubernetes Operator built on top of `clusterctl` and designed to empower cluster administrators to handle the lifecycle of Cluster API providers within a management cluster using a declarative approach. It aims to improve user experience in deploying and managing Cluster API, making it easier to handle day-to-day tasks and automate workflows with GitOps. Visit the [CAPI Operator quickstart](https://cluster-api.sigs.k8s.io/user/quick-start-operator) if you want to experiment with this tool.


### Prerequisites
Install ans setup the following tools:
- [kubectl]
- [Docker]
- [Kind] (v0.25.0 or more recent)
- [Helm]

For your convenience, you can also use the provided devcontainer with all dependencies installed and ready to run on [GitHub Codespaces](https://github.com/codespaces/new?hide_repo_select=true&ref=main&repo=141609301&skip_quickstart=true&machine=standardLinux32gb&geo=UsEast&devcontainer_path=.devcontainer%2Fcapa-devbox%2Fdevcontainer.json) or other cloud development environments.

### Install a management cluster

Cluster API requires an Kubernetes cluster accessible via kubectl. During the installation process the
Kubernetes cluster will be transformed into a [management cluster] by installing the Cluster API [provider components]. This management cluster will be used to create and manage workload clusters, it is recommended to keep it dedicated to this purpose, separated from application workloads.

It is a common practice to create a temporary, local bootstrap cluster which is then used to provision
a target [management cluster] on the selected [infrastructure provider]. In this guide, we'll create a new cluster, using [kind] (Kubernetes in Docker), but you can also use an existing cluster case you have one. Keep in mind that kind is not recommended for production deployments, with strict reliabity and security requirements.

[kind] can be used for creating a local Kubernetes cluster for development environments or for the creation of a temporary [bootstrap cluster] used to provision a target [management cluster] on the selected infrastructure provider.

Create the kind cluster:
```bash
kind create cluster
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
If you are unsure you can determine your computers architecture by running `uname -a`

Download for AMD64:
```bash
curl -L {{#releaselink repo:"https://github.com/kubernetes-sigs/cluster-api" gomodule:"sigs.k8s.io/cluster-api" asset:"clusterctl-linux-amd64" version:"1.9.x"}} -o clusterctl
```

Download for ARM64:
```bash
curl -L {{#releaselink repo:"https://github.com/kubernetes-sigs/cluster-api" gomodule:"sigs.k8s.io/cluster-api" asset:"clusterctl-linux-arm64" version:"1.9.x"}} -o clusterctl
```

Download for PPC64LE:
```bash
curl -L {{#releaselink repo:"https://github.com/kubernetes-sigs/cluster-api" gomodule:"sigs.k8s.io/cluster-api" asset:"clusterctl-linux-ppc64le" version:"1.9.x"}} -o clusterctl
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
curl -L {{#releaselink repo:"https://github.com/kubernetes-sigs/cluster-api" gomodule:"sigs.k8s.io/cluster-api" asset:"clusterctl-darwin-amd64" version:"1.9.x"}} -o clusterctl
```

Download for M1 CPU ("Apple Silicon") / ARM64:
```bash
curl -L {{#releaselink repo:"https://github.com/kubernetes-sigs/cluster-api" gomodule:"sigs.k8s.io/cluster-api" asset:"clusterctl-darwin-arm64" version:"1.9.x"}} -o clusterctl
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

#### Install clusterctl with homebrew on macOS or Linux

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
curl.exe -L {{#releaselink repo:"https://github.com/kubernetes-sigs/cluster-api" gomodule:"sigs.k8s.io/cluster-api" asset:"clusterctl-windows-amd64.exe" version:"1.9.x"}} -o clusterctl.exe
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

The command accepts as input a list of providers to install. When executed for the first time, `clusterctl init`
automatically adds to the list the `cluster-api` core provider, and if unspecified, it also adds the `kubeadm` bootstrap
and `kubeadm` control-plane providers. This guide focuses on the AWS provider.

<!-- 
#### Enabling Feature Gates

Feature gates can be enabled by exporting environment variables before executing `clusterctl init`.
For example, the `ClusterTopology` feature, which is required to enable support for managed topologies and ClusterClass,
can be enabled via:
```bash
export CLUSTER_TOPOLOGY=true
```
Additional documentation about experimental features can be found in [Experimental Features].
-->

Download the latest binary of `clusterawsadm` from the [AWS provider releases]. The [clusterawsadm] command line utility assists with identity and access management (IAM) for [Cluster API Provider AWS][capa].

{{#tabs name:"install-clusterawsadm" tabs:"Linux,macOS,homebrew,Windows"}}
{{#tab Linux}}

Download the latest release; on Linux, type:
<!--
```
curl -L {{#releaselink repo:"https://github.com/kubernetes-sigs/cluster-api-provider-aws" gomodule:"sigs.k8s.io/cluster-api-provider-aws" asset:"clusterawsadm-linux-amd64" version:">=2.0.0"}} -o clusterawsadm
```
-->
```
ARCH=$(uname -m)
if [[ "$ARCH" == "x86_64" ]]; then ARCH="amd64"; fi
if [[ "$ARCH" == "aarch64" ]]; then ARCH="arm64"; fi

VERSION=$(curl -s https://api.github.com/repos/kubernetes-sigs/cluster-api-provider-aws/releases/latest | grep '"tag_name":' | cut -d '"' -f 4)

curl -L "https://github.com/kubernetes-sigs/cluster-api-provider-aws/releases/download/$VERSION/clusterawsadm-linux-$ARCH" -o clusterawsadm
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
export AWS_ACCESS_KEY_ID=<your-access-key>
export AWS_SECRET_ACCESS_KEY=<your-secret-access-key>
export AWS_SESSION_TOKEN=<session-token> # If you are using Multi-Factor Auth.

export AWS_REGION=us-east-1 # This is used to help encode your environment variables
```

Before creating cloud resources, verify that authentication is set to the desired account:
```bash
aws sts get-caller-identity
```

Bootstrap the account with the required IAM resources using cloudformation
```bash
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
curl.exe -L {{#releaselink repo:"https://github.com/kubernetes-sigs/cluster-api-provider-aws" gomodule:"sigs.k8s.io/cluster-api-provider-aws" asset:"clusterawsadm-windows-amd64.exe" version:">=2.0.0"}} -o clusterawsadm.exe
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


The output of `clusterctl init` is similar to this:

```bash
Fetching providers
Installing cert-manager Version="v1.11.0"
Waiting for cert-manager to be available...
Installing Provider="cluster-api" Version="v1.0.0" TargetNamespace="capi-system"
Installing Provider="bootstrap-kubeadm" Version="v1.0.0" TargetNamespace="capi-kubeadm-bootstrap-system"
Installing Provider="control-plane-kubeadm" Version="v1.0.0" TargetNamespace="capi-kubeadm-control-plane-system"
Installing Provider="infrastructure-docker" Version="v1.0.0" TargetNamespace="capd-system"

Your management cluster has been initialized successfully!

You can now create your first workload cluster by running the following:

  clusterctl generate cluster [name] --kubernetes-version [version] | kubectl apply -f -
```

<aside class="note">

<h1>Alternatives to environment variables</h1>

Throughout this quickstart guide we've given instructions on setting parameters using environment variables. For most
environment variables in the rest of the guide, you can also set them in `$XDG_CONFIG_HOME/cluster-api/clusterctl.yaml`

See [`clusterctl init`](../clusterctl/commands/init.md) for more details.

</aside>

### Create the first workload cluster

Once the management cluster is ready, you can create the first workload cluster.

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

#### Required configuration for common providers

Depending on the infrastructure provider you are planning to use, some additional prerequisites should be satisfied
before configuring a cluster with Cluster API. Instructions are provided for common providers below.

Otherwise, you can look at the `clusterctl generate cluster` [command][clusterctl generate cluster] documentation for details about how to
discover the list of variables required by a cluster templates.

{{#tab AWS}}

```bash
export AWS_REGION=us-east-1

export AWS_SSH_KEY_NAME=default

export AWS_CONTROL_PLANE_MACHINE_TYPE=t3.large
export AWS_NODE_MACHINE_TYPE=t3.large

export KUBERNETES_VERSION=$(clusterawsadm ami list --os="ubuntu-24.04" -o json | jq -r '.items[0].spec.kubernetesVersion' | cut -c2-)
echo $KUBERNETES_VERSION
```

See the [AWS provider prerequisites] document for more details.


#### Generating the cluster configuration

For the purpose of this tutorial, we'll name our cluster capa-quickstart.


```bash
clusterctl generate cluster capa-quickstart \
    --kubernetes-version $KUBERNETES_VERSION \
    --control-plane-machine-count=3 \
    --worker-machine-count=3 \
    --from - \
    < templates/cluster-template-quickstart.yaml \
    > "capa-quickstart.yaml"
```

This creates a YAML file named `capa-quickstart.yaml` with a predefined list of Cluster API objects; Cluster, Machines,
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
kubectl get cluster
```

You can also get an "at glance" view of the cluster and its resources by running:

```bash
clusterctl describe cluster capa-quickstart
```

and see an output similar to this:

```bash
NAME              PHASE         AGE   VERSION
capa-quickstart   Provisioned   8s    v1.32.0
```

To verify the first control plane is up:

```bash
kubectl get kubeadmcontrolplane
```

You should see an output is similar to this:

```bash
NAME                    CLUSTER           INITIALIZED   API SERVER AVAILABLE   REPLICAS   READY   UPDATED   UNAVAILABLE   AGE    VERSION
capa-quickstart-g2trk   capa-quickstart   true                                 3                  3         3             4m7s   v1.32.0
```

<aside class="note warning">

<h1> Warning </h1>

The control plane won't be `Ready` until we install a CNI in the next step.

</aside>

After the first control plane node is up and running, we can retrieve the [workload cluster] Kubeconfig.


```bash
clusterctl get kubeconfig capa-quickstart > capa-quickstart.kubeconfig
```

<!--
### Install a Cloud Provider

The Kubernetes in-tree cloud provider implementations are being [removed](https://github.com/kubernetes/enhancements/tree/master/keps/sig-cloud-provider/2395-removing-in-tree-cloud-providers) in favor of external cloud providers (also referred to as "out-of-tree"). This requires deploying a new component called the cloud-controller-manager which is responsible for running all the cloud specific controllers that were previously run in the kube-controller-manager. To learn more, see [this blog post](https://kubernetes.io/blog/2019/04/17/the-future-of-cloud-providers-in-kubernetes/).

{{#tabs name:"tab-install-cloud-provider" tabs:"Azure,OpenStack"}}
{{#tab Azure}}

Install the official cloud-provider-azure Helm chart on the workload cluster:

```bash
helm install --kubeconfig=./capa-quickstart.kubeconfig --repo https://raw.githubusercontent.com/kubernetes-sigs/cloud-provider-azure/master/helm/repo cloud-provider-azure --generate-name --set infra.clusterName=capa-quickstart --set cloudControllerManager.clusterCIDR="192.168.0.0/16"
```

For more information, see the [CAPZ book](https://capz.sigs.k8s.io/self-managed/addons.html).

{{#/tab }}
{{#tab OpenStack}}

Before deploying the OpenStack external cloud provider, configure the `cloud.conf` file for integration with your OpenStack environment:

```bash
cat > cloud.conf <<EOF
[Global]
auth-url=<your_auth_url>
application-credential-id=<your_credential_id>
application-credential-secret=<your_credential_secret>
region=<your_region>
domain-name=<your_domain_name>
EOF
```

For more detailed information on configuring the `cloud.conf` file, see the [OpenStack Cloud Controller Manager documentation](https://github.com/kubernetes/cloud-provider-openstack/blob/master/docs/openstack-cloud-controller-manager/using-openstack-cloud-controller-manager.md#config-openstack-cloud-controller-manager).

Next, create a Kubernetes secret using this configuration to securely store your cloud environment details.
You can create this secret for example with:

```bash
kubectl --kubeconfig=./capa-quickstart.kubeconfig -n kube-system create secret generic cloud-config --from-file=cloud.conf
```

Now, you are ready to deploy the external cloud provider!

```bash
kubectl apply --kubeconfig=./capa-quickstart.kubeconfig -f https://raw.githubusercontent.com/kubernetes/cloud-provider-openstack/master/manifests/controller-manager/cloud-controller-manager-roles.yaml
kubectl apply --kubeconfig=./capa-quickstart.kubeconfig -f https://raw.githubusercontent.com/kubernetes/cloud-provider-openstack/master/manifests/controller-manager/cloud-controller-manager-role-bindings.yaml
kubectl apply --kubeconfig=./capa-quickstart.kubeconfig -f https://raw.githubusercontent.com/kubernetes/cloud-provider-openstack/master/manifests/controller-manager/openstack-cloud-controller-manager-ds.yaml
```

Alternatively, refer to the [helm chart](https://github.com/kubernetes/cloud-provider-openstack/tree/master/charts/openstack-cloud-controller-manager).

{{#/tab }}
{{#/tabs }}
-->

### Deploy a CNI solution

Calico is used here as an example.


```bash
kubectl --kubeconfig=./capa-quickstart.kubeconfig \
  apply -f https://raw.githubusercontent.com/projectcalico/calico/v3.26.1/manifests/calico.yaml
```

After a short while, our nodes should be running and in `Ready` state,
let's check the status using `kubectl get nodes`:

```bash
kubectl --kubeconfig=./capa-quickstart.kubeconfig get nodes
```
```bash
NAME                                          STATUS   ROLES           AGE    VERSION
capa-quickstart-vs89t-gmbld                   Ready    control-plane   5m33s  v1.32.0
capa-quickstart-vs89t-kf9l5                   Ready    control-plane   6m20s  v1.32.0
capa-quickstart-vs89t-t8cfn                   Ready    control-plane   7m10s  v1.32.0
capa-quickstart-md-0-55x6t-5649968bd7-8tq9v   Ready    <none>          6m5s   v1.32.0
capa-quickstart-md-0-55x6t-5649968bd7-glnjd   Ready    <none>          6m9s   v1.32.0
capa-quickstart-md-0-55x6t-5649968bd7-sfzp6   Ready    <none>          6m9s   v1.32.0
```

<aside class="note">

<h1>Troubleshooting</h1>

If the nodes don't become ready after a long period, read the pods in the `kube-system` namespace
```bash
kubectl --kubeconfig=./capa-quickstart.kubeconfig get pod -n kube-system
```

If the Calico pods are in image pull error state (`ErrImagePull`), it's probably because of the Docker Hub pull rate limit.
We can try to fix that by adding a secret with our Docker Hub credentials, and use it;
see [here](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/#registry-secret-existing-credentials)
for details.

First, create the secret. Please notice the Docker config file path, and adjust it to your local setting.
```bash
kubectl --kubeconfig=./capa-quickstart.kubeconfig create secret generic docker-creds \
    --from-file=.dockerconfigjson=<YOUR DOCKER CONFIG FILE PATH> \
    --type=kubernetes.io/dockerconfigjson \
    -n kube-system
```

Now, if the `calico-node` pods are with status of `ErrImagePull`, patch their DaemonSet to make them use the new secret to pull images:
```bash
kubectl --kubeconfig=./capa-quickstart.kubeconfig patch daemonset \
    -n kube-system calico-node \
    -p '{"spec":{"template":{"spec":{"imagePullSecrets":[{"name":"docker-creds"}]}}}}'
```

After a short while, the calico-node pods will be with `Running` status. Now, if the calico-kube-controllers pod is also
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
</aside>

### Clean Up

Delete workload cluster.
```bash
kubectl delete cluster capa-quickstart
```
<aside class="note warning">

IMPORTANT: In order to ensure a proper cleanup of your infrastructure you must always delete the cluster object. Deleting the entire cluster template with `kubectl delete -f capa-quickstart.yaml` might lead to pending resources to be cleaned up manually.
</aside>

Delete management cluster
```bash
kind delete cluster
```

## Troubleshooting
Help with common issues can be found in the [Troubleshooting Guide](./troubleshooting.md).

## Next steps

- Create a second workload cluster. Simply follow the steps outlined above, but remember to provide a different name for your second workload cluster.
- Deploy applications to your workload cluster. Use the [CNI deployment steps](#deploy-a-cni-solution) for pointers.
- See the [clusterctl] documentation for more detail about clusterctl supported actions.

<!-- links -->
[Experimental Features]: ../tasks/experimental-features/experimental-features.md
[Akamai (Linode) provider]: https://linode.github.io/cluster-api-provider-linode/introduction.html
[AWS provider prerequisites]: https://cluster-api-aws.sigs.k8s.io/topics/using-clusterawsadm-to-fulfill-prerequisites.html
[AWS provider releases]: https://github.com/kubernetes-sigs/cluster-api-provider-aws/releases
[Azure Provider Prerequisites]: https://capz.sigs.k8s.io/getting-started.html#prerequisites
[bootstrap cluster]: ../reference/glossary.md#bootstrap-cluster
[capa]: https://cluster-api-aws.sigs.k8s.io
[capv-upload-images]: https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/blob/master/docs/getting_started.md#uploading-the-machine-images
[clusterawsadm]: https://cluster-api-aws.sigs.k8s.io/clusterawsadm/clusterawsadm.html
[clusterctl generate cluster]: ../clusterctl/commands/generate-cluster.md
[clusterctl get kubeconfig]: ../clusterctl/commands/get-kubeconfig.md
[clusterctl]: ../clusterctl/overview.md
[Docker]: https://www.docker.com/
[GCP provider]: https://cluster-api-gcp.sigs.k8s.io/
[Helm]: https://helm.sh/docs/intro/install/
[Harvester provider]: https://github.com/rancher-sandbox/cluster-api-provider-harvester
[Hetzner provider]: https://github.com/syself/cluster-api-provider-hetzner
[Hivelocity provider]: https://github.com/hivelocity/cluster-api-provider-hivelocity
[Huawei Cloud provider]: https://github.com/HuaweiCloudDeveloper/cluster-api-provider-huawei
[IBM Cloud provider]: https://github.com/kubernetes-sigs/cluster-api-provider-ibmcloud
[infrastructure provider]: ../reference/glossary.md#infrastructure-provider
[ionoscloud provider]: https://github.com/ionos-cloud/cluster-api-provider-ionoscloud
[kind]: https://kind.sigs.k8s.io/
[KubeadmControlPlane]: ../developer/core/controllers/control-plane.md
[kubectl]: https://kubernetes.io/docs/tasks/tools/install-kubectl/
[management cluster]: ../reference/glossary.md#management-cluster
[Metal3 getting started guide]: https://github.com/metal3-io/cluster-api-provider-metal3/blob/master/docs/getting-started.md
[Metal3 provider]: https://github.com/metal3-io/cluster-api-provider-metal3/
[K0smotron provider]: https://github.com/k0smotron/k0smotron
[KubeKey provider]: https://github.com/kubesphere/kubekey
[KubeVirt provider]: https://github.com/kubernetes-sigs/cluster-api-provider-kubevirt/
[KubeVirt]: https://kubevirt.io/
[oci-provider]: https://oracle.github.io/cluster-api-provider-oci/#getting-started
[openstack-resource-controller]: https://k-orc.cloud/
[Equinix Metal getting started guide]: https://github.com/kubernetes-sigs/cluster-api-provider-packet#using
[provider]:../reference/providers.md
[provider components]: ../reference/glossary.md#provider-components
[vSphere getting started guide]: https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/blob/master/docs/getting_started.md
[workload cluster]: ../reference/glossary.md#workload-cluster
[CAPI Operator quickstart]: ./quick-start-operator.md
[Proxmox getting started guide]: https://github.com/ionos-cloud/cluster-api-provider-proxmox/blob/main/docs/Usage.md
[Tinkerbell getting started guide]: https://github.com/tinkerbell/cluster-api-provider-tinkerbell/blob/main/docs/QUICK-START.md
[CAPONE Wiki]: https://github.com/OpenNebula/cluster-api-provider-opennebula/wiki