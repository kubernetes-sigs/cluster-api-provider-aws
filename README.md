# OpenShift cluster-api-provider-aws

This repository hosts an implementation of a provider for AWS for the
OpenShift [machine-api](https://github.com/openshift/cluster-api).

This provider runs as a machine-controller deployed by the
[machine-api-operator](https://github.com/openshift/machine-api-operator)

### How to build the images in the RH infrastructure
The Dockerfiles use `as builder` in the `FROM` instruction which is not currently supported
by the RH's docker fork (see [https://github.com/kubernetes-sigs/kubebuilder/issues/268](https://github.com/kubernetes-sigs/kubebuilder/issues/268)).
One needs to run the `imagebuilder` command instead of the `docker build`.

Note: this info is RH only, it needs to be backported every time the `README.md` is synced with the upstream one.

## How to deploy and test the machine controller with minikube

1. **Install kvm**

    Depending on your virtualization manager you can choose a different [driver](https://github.com/kubernetes/minikube/blob/master/docs/drivers.md).
    In order to install kvm, you can run (as described in the [drivers](https://github.com/kubernetes/minikube/blob/master/docs/drivers.md#kvm2-driver) documentation):

    ```sh
    $ sudo yum install libvirt-daemon-kvm qemu-kvm libvirt-daemon-config-network
    $ systemctl start libvirtd
    $ sudo usermod -a -G libvirt $(whoami)
    $ newgrp libvirt
    ```

    To install to kvm2 driver:

    ```sh
    curl -Lo docker-machine-driver-kvm2 https://storage.googleapis.com/minikube/releases/latest/docker-machine-driver-kvm2 \
    && chmod +x docker-machine-driver-kvm2 \
    && sudo cp docker-machine-driver-kvm2 /usr/local/bin/ \
    && rm docker-machine-driver-kvm2
    ```

2. **Deploying the cluster**

    Because of [cluster-api#475](https://github.com/kubernetes-sigs/cluster-api/issues/475) the minikube version can't be higher than `0.28.0`.
    To install minikube `v0.28.0`, you can run:

    ```sg
    $ curl -Lo minikube https://storage.googleapis.com/minikube/releases/v0.28.0/minikube-linux-amd64 && chmod +x minikube && sudo mv minikube /usr/local/bin/
    ```

    To deploy the cluster:

    ```
    minikube start --vm-driver kvm2
    eval $(minikube docker-env)
    ```

3. **Building the machine controller**

    ```
    $ make -C cmd/machine-controller
    ```

4. **Deploying the cluster-api stack manifests**

    Add your AWS credentials to the `addons.yaml` file (in base64
    format). You can either do this manually or use the
    `examples/render-aws-secrets.sh`.

    The easy deployment is:

    ```sh
    ./examples/render-aws-secrets.sh examples/addons.yaml | kubectl apply -f -
    ```

    The manual deployment is:

    ``` sh
    $ echo -n 'your_id' | base64
    $ echo -n 'your_key' | base64
    $ kubectl apply -f examples/addons.yaml
    ```

    Deploy CRDs:

    ```sh
    $ kubectl apply -f config/crd/machine.crd.yaml
    $ kubectl apply -f config/crd/machineset.crd.yaml
    $ kubectl apply -f config/crd/machinedeployment.crd.yaml
    $ kubectl apply -f config/crd/cluster.crd.yaml
    ```

    Deploy machine API controllers:

    ```sh
    $ kubectl apply -f config/rbac/rbac_role.yaml
    $ kubectl apply -f config/rbac/rbac_role_binding.yaml
    $ kubectl apply -f config/controllers/deployment.yaml
    ```

    Deploy the cluster manifest:
    ```sh
    $ kubectl apply -f examples/cluster.yaml
    ```

    Deploy the machines:

    ```sh
    $ kubectl apply -f examples/machine.yaml --validate=false
    ```

    or alternatively:

    ```sh
    $ kubectl apply -f examples/machine-set.yaml --validate=false
    ```

# Upstream Implementation
Other branches of this repository may choose to track the upstream
Kubernetes [Cluster-API AWS provider](https://github.com/kubernetes-sigs/cluster-api-provider-aws/)

In the future, we may align the master branch with the upstream project as it
stabilizes within the community.
