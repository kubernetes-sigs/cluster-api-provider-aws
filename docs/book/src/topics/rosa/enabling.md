# Enabling ROSA Support

To enable support for ROSA clusters, the ROSA feature flag must be set to true. This can be done using the **EXP_ROSA** environment variable:

```shell
export EXP_ROSA="true"
clusterctl init --infrastructure aws
```