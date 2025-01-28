# Reference

## Table of feature gates and their corresponding environment variables

| Feature Gate | Environment Variable | Default |
| ------------ | -------------------- |---------|
| EKS                           | CAPA_EKS                          | true    |
| EKSEnableIAM                  | CAPA_EKS_IAM	                    | false   |
| EKSAllowAddRoles              | CAPA_EKS_ADD_ROLES                | false   |
| EKSFargate                    | EXP_EKS_FARGATE                   | false   |
| MachinePool                   | EXP_MACHINE_POOL                  | false   |
| EventBridgeInstanceState      | EVENT_BRIDGE_INSTANCE_STATE       | false   |
| AutoControllerIdentityCreator | AUTO_CONTROLLER_IDENTITY_CREATOR  | true    |
| BootstrapFormatIgnition       | EXP_BOOTSTRAP_FORMAT_IGNITION     | false   |
| ExternalResourceGC            | EXP_EXTERNAL_RESOURCE_GC          | false   |
| AlternativeGCStrategy         | EXP_ALTERNATIVE_GC_STRATEGY       | false   |
| TagUnmanagedNetworkResources  | TAG_UNMANAGED_NETWORK_RESOURCES   | true    |
| ROSA                          | EXP_ROSA                          | false   |