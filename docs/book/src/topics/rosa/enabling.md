# Enabling ROSA Support

To enable support for ROSA clusters, the ROSA feature flag must be set to true. This can be done using the **EXP_ROSA** environment variable:

```shell
export EXP_ROSA="true"
export EXP_MachinePool="true"
clusterctl init --infrastructure aws
```

## TroubleShoot
In case the cluster is initalized and the MachinePool Feature did not become enable. Run the following command.
```shell
$ kubectl edit deploy capi-controller-manager -n capi-system
```
edit the featue gate container args to enable the MachinePool=true as below.

```yaml
spec:
      containers:
      - args:
        - --leader-elect
        - --diagnostics-address=:8443
        - --insecure-diagnostics=false
        - --feature-gates=MachinePool=true,ClusterResourceSet=true,ClusterTopology=true,RuntimeSDK=false,MachineSetPreflightChecks=false
```
