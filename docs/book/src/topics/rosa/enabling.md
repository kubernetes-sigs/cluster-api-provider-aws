# Enabling ROSA Support

To enable support for ROSA clusters, the ROSA feature flag must be set to true. This can be done using the **EXP_ROSA** environment variable:

```shell
export EXP_ROSA="true"
export EXP_MACHINE_POOL="true"
clusterctl init --infrastructure aws
```

## Troubleshooting
To check the feature-gates for the Cluster API controller run the following command:

```shell
$ kubectl get deploy capi-controller-manager -n capi-system -o yaml
```
the feature gate container arg should have `MachinePool=true` as shown below.

```yaml
spec:
  containers:
  - args:
    - --feature-gates=MachinePool=true,ClusterTopology=true,...
```

To check the feature-gates for the Cluster API AWS controller run the following command:
```shell
$ kubectl get deploy capa-controller-manager -n capa-system -o yaml
```
the feature gate arg should have `ROSA=true` as shown below.

```yaml
spec:
  containers:
  - args:
    - --feature-gates=ROSA=true,...
```