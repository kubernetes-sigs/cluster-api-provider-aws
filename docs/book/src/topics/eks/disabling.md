# Disabling EKS Support

Support for EKS is enabled by default when you use the AWS infrastructure provider. But if you never plan to use EKS then you can disable the EKS support. The following sections describe the process.

## Disabling IAM objects for EKS

To ensure that there are no IAM objects created for EKS you will need to use a configuration file with `clusterawsadm` and specify that EKS is disabled:

```yaml
apiVersion: bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1
kind: AWSIAMConfiguration
spec:
  eks:
    disable: true
```

and then use that configuration file:

```bash
clusterawsadm bootstrap iam create-cloudformation-stack --config bootstrap-config.yaml
```

## Disable EKS in the provider

Disbling EKS support is done via the **EKS** feature flag by setting it to false. This can be done before running `clusterctl init` by using the **CAPA_EKS** environment variable:

```shell
export CAPA_EKS=false
clusterctl init --infrastructure aws
```
