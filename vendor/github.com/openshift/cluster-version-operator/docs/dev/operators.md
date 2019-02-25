# Operator integration with CVO

The CVO installs other operators onto a cluster. It is responsible for applying the manifests each
operator uses (without any parameterization) and for ensuring an order that installation and
updates follow.

## What is the order that resources get created/updated in?

The CVO will load a release image and look at the contents of two directories: `/manifests` and
`/release-manifests` within that image. The `/manifests` directory is part of the CVO image and
contains the basic deployment and other roles for the CVO. The `/release-manifests` directory is
created by the `oc adm release new` command from the `/manifests` directories of all other
candidate operators.

The contents of `/release-manifests` are applied in lexigraphic order, exactly as `ls` would
return on a standard Linux or Unix derivative. The CVO supports the idea of "run levels" by
defining a convention for how operators that wish to run before other operators should name
their manifests. A run level is of the form `0000_\d\d_[a-z0-9\-]+_<filename>` where the first
digits are the level (see [below for a list of assigned levels](#how-do-i-get-added-as-a-special-run-level)), the second chunk is the component name that
usually matches your operator name (e.g. `kube-apiserver` or `cluster-monitoring-operator`)
and the filename is a local name.

A few special optimizations are applied above linear ordering - if the CVO sees two different
components that have the same run level - for instance, `0000_70_cluster-monitoring-operator_*` and
`0000_70_cluster-samples-operator_*` - each component will execute in parallel to the others,
preserving the order of tasks within the component.

Ordering is most important during upgrades, where some components rely on another component
being updated first. As a convenience, the CVO guarantees that components at an earlier
run level will be created or updated before your component is invoked. Note however that
components without `ClusterOperator` objects defined may not be fully deployed when your
component is executed, so always ensure your prerequisites know that they must correctly
obey the `ClusterOperator` protocol to be available. More sophisticated components should
observe the prerequisite `ClusterOperator`s directly and use the `versions` field to
enforce safety.

## How do I get added to the release image?

Add the following to your Dockerfile

```Dockerfile
FROM …

ADD manifests-for-operator/ /manifests
LABEL io.openshift.release.operator=true
```

Ensure your image is published into the cluster release tag by ci-operator
Wait for a new release image to be created (usually once you push to master in your operator).

## What do I put in /manifests?

You need the following:

1..N manifest yaml or JSON files (preferably YAML for readability) that deploy your operator, including:

- Namespace for your operator
- Roles your operator needs
- A service account and a service account role binding
- Deployment for your operator
- A ClusterOperator CR [more info here](clusteroperator.md)
- Any other config objects your operator might need
- An image-references file (See below)

In your deployment you can reference the latest development version of your operator image (quay.io/openshift/origin-machine-api-operator:latest).  If you have other hard-coded image strings, try to put them as environment variables on your deployment or as a config map.

### Names of manifest files

Your manifests will be applied in alphabetical order by the CVO, so name your files in the order you want them run.
If you are a normal operator (don’t need to run before the kube apiserver), you should name your manifest files in a way that feels easy:

```
/manifests/
  deployment.yaml
  roles.yaml
```

If you’d like to ensure your manifests are applied in order to the cluster add a numeric prefix to sort in the directory:

```
/manifests/
  01_roles.yaml
  02_deployment.yaml
```

When your manifests are added to the release image, they’ll be given a prefix that corresponds to the name of your repo/image:

```
/release-manifests/
  99_ingress-operator_01_roles.yaml
  99_ingress-operator_02_deployment.yaml
```

Only manifests with the extensions `.yaml`, `.yml`, or `.json` will be applied, like `kubectl create -f DIR`.

### How do I get added as a special run level?

Some operators need to run at a specific time in the release process (OLM, kube, openshift core operators, network, service CA).  These components can ensure they run in a specific order across operators by prefixing their manifests with:

    0000_<runlevel>_<dash-separated_component>-<manifest_filename>

For example, the Kube core operators run in runlevel 10-19 and have filenames like

    0000_13_cluster-kube-scheduler-operator_03_crd.yaml

Assigned runlevels

- 00-04 - CVO
- 05 - cluster-config-operator
- 07 - Network operator
- 08 - DNS operator
- 09 - Service certificate authority and machine approver
- 10-29 - Kubernetes operators (master team)
- 30-39 - Machine API
- 50-59 - Operator-lifecycle manager
- 60-69 - OpenShift core operators (master team)

## How do I ensure the right images get used by my manifests?

Your manifests can contain a tag to the latest development image published by Origin.  You’ll annotate your manifests by creating a file that identifies those images.

Assume you have two images in your manifests - `quay.io/openshift/origin-ingress-operator:latest` and `quay.io/openshift/origin-haproxy-router:latest`.  Those correspond to the following tags `ingress-operator` and `haproxy-router` when the CI runs.

Create a file `image-references` in the /manifests dir with the following contents:

```yaml
kind: ImageStream
apiVersion: image.openshift.io/v1
spec:
  tags:
  - name: ingress-operator
    from:
      kind: DockerImage
      Name: quay.io/openshift/origin-ingress-operator
  - name: haproxy-router
    from:
      kind: DockerImage
      Name: quay.io/openshift/origin-haproxy-router
```

The release tooling will read image-references and do the following operations:

Verify that the tags `ingress-operator` and `haproxy-router` exist from the release / CI tooling (in the image stream `openshift/origin-v4.0` on api.ci).  If they don’t exist, you’ll get a build error.
Do a find and replace in your manifests (effectively a sed)  that replaces `quay.io/openshift/origin-haproxy-router(:.*|@:.*)` with `registry.svc.ci.openshift.org/openshift/origin-v4.0@sha256:<latest SHA for :haproxy-router>`
Store the fact that operator ingress-operator uses both of those images in a metadata file alongside the manifests
Bundle up your manifests and the metadata file as a docker image and push them to a registry

Later on, when someone wants to mirror a particular release, there will be tooling that can take the list of all images used by operators and mirror them to a new repo.

This pattern tries to balance between having the manifests in your source repo be able to deploy your latest upstream code *and* allowing us to get a full listing of all images used by various operators.
