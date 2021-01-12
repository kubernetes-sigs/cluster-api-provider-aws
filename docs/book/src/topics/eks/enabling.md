# Enabling EKS Support

You must explicitly enable the EKS support in the provider by doing the following:

- Enabling support in the infrastructure manager (capa-controller-manager) by enabling the **EKS** feature flags (see below)
- Add the EKS Control Plane Provider (aws-eks)
- Add the EKS Bootstrap Provider (aws-eks)

## Enabling the **EKS** features

Enabling the **EKS** functionality is done using the following feature flags:

- **EKS** - this enables the core EKS functionality and is required for the other EKS feature flags
- **EKSEnableIAM** - by enabling this the controllers will create any IAM roles required by EKS and the roles will be cluster specific. If this isn't enabled then you can manually create a role and specify the role name in the AWSManagedControlPlane spec otherwise the default rolename will be used.
- **EKSAllowAddRoles** - by enabling this you can add additional roles to the control plane role that is created. This has no affect unless used wtih __EKSEnableIAM__

Enabling the feature flags can be done using `clusterctl` by setting the following environment variables to **true** (they all default to **false**):

- **EXP_EKS** - this is used to set the value of the **EKS** feature flag
- **EXP_EKS_IAM** - this is used to set the value of the **EKSEnableIAM** feature flag
- **EXP_EKS_ADD_ROLES** - this is used to set the value of the **EKSAllowAddRoles** feature flag

As an example:

```bash
export EXP_EKS=true
export EXP_EKS_IAM=true
export EXP_EKS_ADD_ROLES=true

clusterctl init --infrastructure=aws --control-plane aws-eks --bootstrap aws-eks
```

