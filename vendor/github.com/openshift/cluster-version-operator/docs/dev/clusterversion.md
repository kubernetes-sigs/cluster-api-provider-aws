# ClusterVersion Custom Resource

The `ClusterVersion` is a custom resource object which holds the current version of the cluster.
This object is used by the administrator to declare their target cluster state, which the cluster-version operator (CVO) then works to transition the cluster to that target state.

## Finding your current update image

You can extract the current update image from the `ClusterVersion` object:

```console
$ oc get clusterversion -o jsonpath='{.status.current.image}{"\n"}' version
registry.svc.ci.openshift.org/openshift/origin-release@sha256:c1f11884c72458ffe91708a4f85283d591b42483c2325c3d379c3d32c6ac6833
```

## Setting objects unmanaged

For testing operators, it is sometimes helpful to disable CVO management so you can alter objects without the CVO stomping on your changes.
To get a list of objects managed by the CVO, run:

```console
$ oc adm release extract --from=registry.svc.ci.openshift.org/openshift/origin-release@sha256:c1f11884c72458ffe91708a4f85283d591b42483c2325c3d379c3d32c6ac6833 --to=release-image
$ ls release-image | head -n5
0000_07_cluster-network-operator_00_namespace.yaml
0000_07_cluster-network-operator_01_crd.yaml
0000_07_cluster-network-operator_02_rbac.yaml
0000_07_cluster-network-operator_03_daemonset.yaml
0000_08_cluster-dns-operator_00-cluster-role.yaml
```

To get a list of current overrides, run:

```console
$ oc get -o json clusterversion version | jq .spec.overrides
[
  {
    "kind": "APIService",
    "name": "v1alpha1.packages.apps.redhat.com",
    "unmanaged": true
  }
]
```

To add an entry to that list, you can use a [JSON Patch][json-patch] to add a [`ComponentOverride`][ComponentOverride].
For example, to set the network operator's daemonset unmanaged:

```console
$ head -n5 release-image/0000_07_cluster-network-operator_03_daemonset.yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: cluster-network-operator
  namespace: openshift-cluster-network-operator
$ cat <<EOF >version-patch.yaml
> - op: add
>   path: /spec/overrides/-
>   value:
>     kind: DaemonSet
>     name: cluster-network-operator
>     namespace: openshift-cluster-network-operator
>     unmanaged: true
> EOF
$ oc patch clusterversion version --type json -p "$(cat version-patch.yaml)"
```

You can verify the update with:

```console
$ oc get -o json clusterversion version | jq .spec.overrides
[
  {
    "kind": "APIService",
    "name": "v1alpha1.packages.apps.redhat.com",
    "unmanaged": true
  },
  {
    "kind": "DaemonSet",
    "name": "cluster-network-operator",
    "namespace": "openshift-cluster-network-operator",
    "unmanaged": true
  }
]
```

After updating the `ClusterVersion`, you can make your desired edits to the unmanaged object.

## Disabling the cluster-version operator

When you just want to turn off the cluster-version operator instead of fiddling with per-object overrides, you can:

```console
$ oc scale --replicas 0 -n openshift-cluster-version deployments/cluster-version-operator
```

[ComponentOverride]: https://godoc.org/github.com/openshift/api/config/v1#ComponentOverride
[json-patch]: https://tools.ietf.org/html/rfc6902
