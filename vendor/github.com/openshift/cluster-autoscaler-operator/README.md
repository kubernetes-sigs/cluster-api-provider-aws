# Cluster Autoscaler Operator

The cluster-autoscaler-operator manages deployments of the OpenShift
[Cluster Autoscaler][1] using the [cluster-api][2] provider.

[1]: https://github.com/openshift/kubernetes-autoscaler/tree/master/cluster-autoscaler
[2]: https://github.com/kubernetes-sigs/cluster-api


## Custom Resource Definitions

The operator manages the following custom resources:

- __ClusterAutoscaler__: This is a singleton resource which controls the
  configuration of the cluster's autoscaler instance.  The operator will
  only respond to the ClusterAutoscaler resource named "default" in the
  managed namespace, i.e. the value of the `WATCH_NAMESPACE` environment
  variable.  ([Example][ClusterAutoscaler])

  The fields in the spec for ClusterAutoscaler resources correspond to
  command-line arguments to the cluster-autoscaler.  The example
  linked above results in the following invocation:

  ```
    Command:
      cluster-autoscaler
    Args:
      --logtostderr
      --cloud-provider=cluster-api
      --namespace=openshift-machine-api
      --expendable-pods-priority-cutoff=-10
      --max-nodes-total=24
      --cores-total=8:128
      --memory-total=4:256
      --gpu-total=nvidia.com/gpu:0:16
      --gpu-total=amd.com/gpu:0:4
      --scale-down-enabled=true
      --scale-down-delay-after-add=10s
      --scale-down-delay-after-delete=10s
      --scale-down-delay-after-failure=10s
  ```

- __MachineAutoscaler__: This resource targets a node group and manages
  the annotations to enable and configure autoscaling for that group,
  e.g. the min and max size.  Currently only `MachineSet` objects can be
  targeted.  ([Example][MachineAutoscaler])

[ClusterAutoscaler]: https://github.com/openshift/cluster-autoscaler-operator/blob/master/examples/clusterautoscaler.yaml
[MachineAutoscaler]: https://github.com/openshift/cluster-autoscaler-operator/blob/master/examples/machineautoscaler.yaml


## Development

```sh-session
## Build, Test, & Run
$ make build
$ make test

$ export WATCH_NAMESPACE=openshift-machine-api
$ ./bin/cluster-autoscaler-operator -alsologtostderr
```

The Cluster Autoscaler Operator is designed to be deployed on
OpenShift by the [Cluster Version Operator][CVO], but it's possible to
run it directly on any vanilla Kubernetes cluster that has the
[machine-api][machine-api] components available.  To do so, apply the
manifests in the install directory: `kubectl apply -f ./install`

This will create the `openshift-machine-api` namespace, register the
custom resource definitions, configure RBAC policies, and create a
deployment for the operator.

[CVO]: https://github.com/openshift/cluster-version-operator
[machine-api]: https://github.com/openshift/cluster-api
[cluster-api]: https://github.com/kubernetes-sigs/cluster-api


### End-to-End Tests

You can run the e2e test suite with `make test-e2e`.  These tests
assume the presence of a cluster not already running the operator, and
that the `KUBECONFIG` environment variable points to a configuration
granting admin rights on said cluster.
